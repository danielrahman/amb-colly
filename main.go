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
	Category string
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
	c.OnHTML(`#content > div > div .caption a[href]`, func(e *colly.HTMLElement) {
		courseURL := e.Request.AbsoluteURL(e.Attr(`href`))
		detailCollector.Visit(courseURL)
	})

	// Extract details from products
	detailCollector.OnHTML(`.content-box`, func(e *colly.HTMLElement) {

		title := html.EscapeString(e.ChildText(`h1`))
		log.Println(`Product:`, title)
		vendor := html.EscapeString(e.ChildText(`#content > div > div:nth-child(1) > div.col-sm-4 > ul:nth-child(2) > li:nth-child(1) > a`))
		quantityDirty := e.ChildText(`#content > div > div:nth-child(1) > div.col-sm-4 > ul:nth-child(2) > li:nth-child(4) > span`)
		quantityClean := strings.Replace(quantityDirty, " ks", "", -1)
		priceDirty := html.EscapeString(e.ChildText(`li.big-price`))
		priceClean := strings.Replace(priceDirty, " Kč", "", -1)
		priceClean = strings.Replace(priceClean, " ", "", -1)
		code := html.EscapeString(e.ChildText(`#content > div > div:nth-child(1) > div.col-sm-4 > ul:nth-child(2) > li:nth-child(2) > span`))

		productUrl := e.Request.URL.String()
		category := getCategory(productUrl)

		if quantityClean == "Vyprodáno" || quantityClean == "" {
			quantityClean = "0"
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
			Category: category,
			Vendor:   vendor,
			Quantity: quantity,
			Price:    price,
			Code:     code,
			URL:      e.Request.URL.String(),
		}
		// Iterate over rows of the table which contains different information
		// about the product
		db.UpdateDatabase(fmt.Sprintf(`INSERT INTO products (product_id, Title, Category, Vendor, Quantity, Price, Code, Url)
			VALUES ("%d", "%s", "%s", "%s","%d", "%d","%s", "%s" )
			ON DUPLICATE KEY UPDATE product_id=VALUES(product_id), Title=VALUES(Title), Category=VALUES(Category), Vendor=VALUES(Vendor), Quantity=VALUES(Quantity), Price=VALUES(Price), Code=VALUES(Code), Url=VALUES(Url)`,
			product.Id, product.Title, product.Category, product.Vendor, product.Quantity, product.Price, product.Code, product.URL))
	})

	c.Visit(`https://www.ambassadors.eu/skate/skateboard-desky?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/skate/komplety-a-jine-desky?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/skate/gripy-na-skateboard?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/skate/trucky-pro-skateboard?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/skate/loziska-pro-skateboard?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/skate/vosky?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/skate/skateboard-bushings?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/skate/ostatni-hardware?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/skate/chranice?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/skate/kolecka-pro-skateboard?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/obleceni/jacket?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/obleceni/trika?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/obleceni/kalhoty-teplaky-kratasy?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/obleceni/mikiny?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/obleceni/ksiltovky?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/obleceni/kulichy?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/obleceni/pasky?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/obleceni/ponozky?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/obleceni/tilka?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/obleceni/ostatni?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/doplnky/penezenka?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/doplnky/slunecni-bryle?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/doplnky/samolepky-a-nasivky?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/doplnky/ostatnidoplnky?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/doplnky/plakaty?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/doplnky/DVD-CD?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/doplnky/thrasher-magazine?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/boty/obuv?limit=1000`)
	c.Visit(`https://www.ambassadors.eu/boty/vlozky-do-bot?limit=1000`)
}

func getCategory(productUrl string) string {
	if strings.Contains(productUrl, "/skate/") {
		category := strings.Split(productUrl, "/skate/")
		category = strings.Split(category[1], "?")
		category = strings.Split(category[0], "/")
		return category[0]
	} else if strings.Contains(productUrl, "/obleceni/") {
		category := strings.Split(productUrl, "/obleceni/")
		category = strings.Split(category[1], "?")
		category = strings.Split(category[0], "/")
		return category[0]
	} else if strings.Contains(productUrl, "/doplnky/") {
		category := strings.Split(productUrl, "/doplnky/")
		category = strings.Split(category[1], "?")
		category = strings.Split(category[0], "/")
		return category[0]
	} else if strings.Contains(productUrl, "/boty/") {
		category := strings.Split(productUrl, "/boty/")
		category = strings.Split(category[1], "?")
		category = strings.Split(category[0], "/")
		return category[0]
	}
	return "Nezařazené"
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
