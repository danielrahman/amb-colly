package main

import (
	"fmt"
	"github.com/danielrahman/ambassadorsscraper/ambassadors"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"hash/fnv"
	"html"
	"strconv"
	"strings"
)

// Product stores information about a coursera course
type Product struct {
	Id       uint32
	Title    string
	Vendor   string
	Quantity int64
	Price    int64
	Code     string
	URL      string
}

func main() {

	var db ambassadors.DbAmbassadors
	_, err := db.ConnectDatabase()
	if err != nil {
		log.Error(err.Error())
	}
	// Instantiate default collector
	c := colly.NewCollector()

	// Create another collector to scrape course details
	detailCollector := c.Clone()

	// Before making a request print `Visiting ...`
	c.OnRequest(func(r *colly.Request) {
		log.Println(`visiting`, r.URL.String())
	})

	// On every a HTML element which has name attribute call callback
	c.OnHTML(`#content > div > div:nth-child(10) a[href]`, func(e *colly.HTMLElement) {
		courseURL := e.Request.AbsoluteURL(e.Attr(`href`))
		detailCollector.Visit(courseURL)
	})

	// Extract details from products
	detailCollector.OnHTML(`.content-box`, func(e *colly.HTMLElement) {
		log.Println(`Product found`, e.Request.URL)

		title := html.EscapeString(e.ChildText(`h1`))
		vendor := html.EscapeString(e.ChildText(`#content > div > div:nth-child(1) > div.col-sm-4 > ul:nth-child(2) > li:nth-child(1) > a`))
		quantityDirty := e.ChildText(`#content > div > div:nth-child(1) > div.col-sm-4 > ul:nth-child(2) > li:nth-child(4) > span`)
		quantityClean := strings.Replace(quantityDirty, " ks", "", -1)
		priceDirty := html.EscapeString(e.ChildText(`li.big-price`))
		priceClean := strings.Replace(priceDirty, " Kč", "", -1)
		priceClean = strings.Replace(priceClean, " ", "", -1)
		code := html.EscapeString(e.ChildText(`#content > div > div:nth-child(1) > div.col-sm-4 > ul:nth-child(2) > li:nth-child(2) > span`))

		if quantityClean == "Vyprodáno" {
			priceClean = "0"
		}
		quantity, err := strconv.ParseInt(quantityClean, 10, 64)
		if err != nil {
			log.Error(err.Error())
		}

		price, err := strconv.ParseInt(priceClean, 10, 64)
		if err != nil {
			log.Error(err.Error())
		}
		product := Product{
			Id:       hash(code),
			Title:    title,
			Vendor:   vendor,
			Quantity: quantity,
			Price:    price,
			Code:     code,
			URL:      e.Request.URL.String(),
		}
		// Iterate over rows of the table which contains different information
		// about the product
		db.UpdateDatabase(fmt.Sprintf(`INSERT INTO products (product_id, Title, Vendor, Quantity, Price, Code, Url)
			VALUES ("%d", "%s", "%s","%d", "%d","%s", "%s" )
			ON DUPLICATE KEY UPDATE product_id=VALUES(product_id), Title=VALUES(Title), Vendor=VALUES(Vendor), Quantity=VALUES(Quantity), Price=VALUES(Price), Code=VALUES(Code), Url=VALUES(Url)`,
			product.Id, product.Title, product.Vendor, product.Quantity, product.Price, product.Code, product.URL))
	})

	c.Visit(`https://www.ambassadors.eu/skate?limit=1000`)

}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
