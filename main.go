package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

type ride struct {
	name                           string
	shortDescription               string
	imageURL                       string
	longDescription                string
	tags                           []string
	heightRequirementInInches      int
	heightRequirementInCentimeters int
	rideURL                        string
}

func main() {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	//make generic with interface to collect ride info (scrape) and write
	buschgardens(*c)

}

func buschgardens(c colly.Collector) {
	c.OnHTML("#page-content > div > div > ul", func(e *colly.HTMLElement) {

		e.ForEach("li", func(_ int, el *colly.HTMLElement) {

			image := e.Request.URL.String() + el.ChildAttr("span > a > img", "src")
			name := el.ChildText("div > h2 > a")
			shortDescription := el.ChildText("div > p")
			tags := el.ChildTexts("div > ul > li")

			if name != "" {
				log.Printf("name %v sdesc %v \n", name, shortDescription)
				log.Printf("image %v \n", image)
				log.Printf("tags %v \n", tags)

				//follow link
				rideURL := e.Request.URL.Scheme + "://" + e.Request.URL.Host + e.Request.URL.Opaque + el.ChildAttr("span > a", "href")
				log.Printf("rideURL %v \n", rideURL)
				c2 := colly.NewCollector()
				c2.OnHTML("#page-content", func(e2 *colly.HTMLElement) {
					height := e2.ChildText("#page-content > div.container > div > div.col-sm-8 > div > ul > li:nth-child(1) > dl > dd")
					height = strings.TrimSuffix(height, "\"")
					log.Printf("height %v \n", height)
				})

				c2.Visit(rideURL)

			}

		})

	})
	c.Visit("https://buschgardens.com/williamsburg/rides")
}
