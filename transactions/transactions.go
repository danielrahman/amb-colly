package main

import (
	"fmt"
	"github.com/danielrahman/amb-colly/ambassadors"
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
	db.UpdateDatabase(fmt.Sprintf(`INSERT INTO log (type, status, date) VALUES ("transactions", "start", "%s")`, time.Now().Format("20060102150405")))

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
			if quantityClean == "Vyprodáno" || quantityClean == "Skladem" || quantityClean == "Není skladem" || quantityClean == "" {
				quantityClean = "0"
			}
			quantityNew, err := strconv.Atoi(quantityClean)
			if err != nil {
				log.Error(err.Error())
				log.Info(productId)
				log.Info(url)
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

			adjustment := quantityNew - sum
			quantityActual := sum + adjustment

			if sum == quantityNew {
				return
			} else {
				log.Println(productId, "/", adjustment)

				db.UpdateDatabase(fmt.Sprintf(`INSERT INTO products (product_id, quantity) VALUES ("%s", "%d") ON DUPLICATE KEY UPDATE product_id=VALUES(product_id), Quantity=VALUES(Quantity) `, productId, quantityActual))

				db.UpdateDatabase(fmt.Sprintf(`INSERT INTO transactions (product_id, date, Adjustment, Quantity)
			VALUES ("%s", "%s","%d", "%d" )`,
					productId, time.Now().Format("20060102150405"), adjustment, quantityActual))
			}

		})

		c.Visit(url)
	}
	db.UpdateDatabase(fmt.Sprintf(`INSERT INTO log (type, status, date) VALUES ("transactions", "end", "%s")`, time.Now().Format("20060102150405")))
}
