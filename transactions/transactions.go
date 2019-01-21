package main

import (
	"fmt"
	"github.com/danielrahman/ambassadorsscraper/ambassadors"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

func main() {
	var db ambassadors.DbAmbassadors
	_, err := db.ConnectDatabase()
	if err != nil {
		log.Error(err.Error())
	}

	products := db.GetData("product_id, url", "products")
	defer products.Close()
	for products.Next() {
		c := colly.NewCollector()
		var productId string
		var url string
		if err := products.Scan(&productId, &url); err != nil {
			log.Fatal(err)
		}
		c.OnHTML(".content-box", func(e *colly.HTMLElement) {
			quantityDirty := e.ChildText(`#content > div > div:nth-child(1) > div.col-sm-4 > ul:nth-child(2) > li:nth-child(4) > span`)
			quantityClean := strings.Replace(quantityDirty, " ks", "", -1)
			if quantityClean == "Vyprodáno" {
				quantityClean = "0"
			}
			quantityNew, err := strconv.Atoi(quantityClean)
			if err != nil {
				log.Error(err.Error())
			}

			transactions := db.GetData("adjustment", "transactions WHERE product_id = "+productId)
			sum := 0
			for transactions.Next() {
				var adjustment int
				if err := transactions.Scan(&adjustment); err != nil {
					log.Fatal(err)
				}
				sum += adjustment
			}

			println("Product url: ", url)
			println("Product id: ", productId)
			println("Product new quantity: ", quantityNew)
			println("Product sum transactions: ", sum)

			if sum == quantityNew {
				return
			} else {
				adjustment := quantityNew - sum
				quantityActual := sum + adjustment
				println("Adjustment (quantityNew - sum): ", adjustment)
				println("quantityActual (sum + adjustment): ", quantityActual)
				db.UpdateDatabase(fmt.Sprintf(`INSERT INTO transactions (product_id, date, Adjustment, Quantity)
			VALUES ("%s", "%s","%d", "%d" )`,
					productId, time.Now().Format("20060102150405"), adjustment, quantityActual))
			}

		})

		c.Visit(url)
	}

}
