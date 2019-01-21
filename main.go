package main

import (
	"github.com/danielrahman/ambassadorsscraper/ambassadors"
	"github.com/gocolly/colly"
	log "github.com/sirupsen/logrus"
	"strings"
)

// Product stores information about a coursera course
type Product struct {
	Title        string
	Category     string
	Vendor       string
	Availability string
	Price        string
	Code         string
	URL          string
}

func main() {
	// Instantiate default collector
	c := colly.NewCollector()

	// Create another collector to scrape course details
	detailCollector := c.Clone()

	products := make([]Product, 0, 200)

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		log.Println("visiting", r.URL.String())
	})

	// On every a HTML element which has name attribute call callback
	c.OnHTML(`#content > div > div:nth-child(10) a[href]`, func(e *colly.HTMLElement) {
		// Activate detailCollector if the link contains "ambassadors.eu/skate/skateboard-desky"
		courseURL := e.Request.AbsoluteURL(e.Attr("href"))
		if strings.Index(courseURL, "ambassadors.eu/skate/skateboard-desky") != -1 {
			detailCollector.Visit(courseURL)
		}
	})

	// Extract details of the course
	detailCollector.OnHTML(`.content-box`, func(e *colly.HTMLElement) {
		log.Println("Product found", e.Request.URL)

		title := e.ChildText("h1")
		category := e.ChildText("body > div:nth-child(9) > ul > li:nth-child(3) > a")
		vendor := e.ChildText("#content > div > div:nth-child(1) > div.col-sm-4 > ul:nth-child(2) > li:nth-child(1) > a")
		availability := e.ChildText(`#content > div > div:nth-child(1) > div.col-sm-4 > ul:nth-child(2) > li:nth-child(4) > span`)
		price := e.ChildText(".big-price")
		code := e.ChildText("#content > div > div:nth-child(1) > div.col-sm-4 > ul:nth-child(2) > li:nth-child(2) > span")

		product := Product{
			Title:        title,
			Category:     category,
			Vendor:       vendor,
			Availability: availability,
			Price:        price,
			Code:         code,
			URL:          e.Request.URL.String(),
		}
		// Iterate over rows of the table which contains different information
		// about the product
		products = append(products, product)
	})

	c.OnScraped(func(r *colly.Response) {
		var db ambassadors.DbAmbassadors
		_, err := db.ConnectDatabase()
		if err != nil {
			log.Error(err.Error())
		}
		log.Println(products)
	})

	c.Visit("https://www.ambassadors.eu/skate/skateboard-desky")

}
