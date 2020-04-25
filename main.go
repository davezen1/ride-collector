package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

//Ride details
type Ride struct {
	Park                      string   `json:"park"`
	Name                      string   `json:"name"`
	ShortDescription          string   `json:"description"`
	ImageURL                  string   `json:"image"`
	HeightRequirementInInches string   `json:"height"`
	RideURL                   string   `json:"url"`
	Tags                      []string `json:"tags"`
}

//Rides represents slice of Ride structs
type Rides []Ride

func buschgardens(c colly.Collector) ([]Ride, [][]string) {
	var ridesList Rides
	var ridesSlices = [][]string{{"name", "image", "description", "tags", "height", "url"}}
	parkName := "Busch Gardens Williamsburg"

	c.OnHTML("#page-content > div > div > ul", func(e *colly.HTMLElement) {

		e.ForEach("li", func(_ int, el *colly.HTMLElement) {

			image := e.Request.URL.String() + el.ChildAttr("span > a > img", "src")
			name := el.ChildText("div > h2 > a")
			shortDescription := el.ChildText("div > p")
			tags := el.ChildTexts("div > ul > li")

			if name != "" {
				//follow link
				rideURL := e.Request.URL.Scheme + "://" + e.Request.URL.Host + e.Request.URL.Opaque + el.ChildAttr("span > a", "href")
				c2 := colly.NewCollector()
				c2.OnHTML("#page-content", func(e2 *colly.HTMLElement) {
					height := e2.ChildText("#page-content > div.container > div > div.col-sm-8 > div > ul > li:nth-child(1) > dl > dd")
					height = strings.TrimSuffix(height, "\"")

					ridesList = append(ridesList, Ride{Park: parkName, Name: name, ImageURL: image, ShortDescription: shortDescription, Tags: tags, HeightRequirementInInches: height, RideURL: rideURL})

					ridesSlices = append(ridesSlices, []string{name, image, shortDescription, strings.Join(tags, "|"), height, rideURL})
				})

				c2.Visit(rideURL)

			}

		})

	})
	c.Visit("https://buschgardens.com/williamsburg/rides")

	return ridesList, ridesSlices
}

func main() {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	//make generic with interface to collect ride info (extract) and write
	bgRidesList, _ := buschgardens(*c)

	write(bgRidesList)
}

func write(r Rides) {

	var fileName strings.Builder
	ridesJSON, err := json.MarshalIndent(r, "", "    ")

	if err != nil {
		log.Fatal("Cannot encode to JSON ", err)
	}
	fmt.Println(string(ridesJSON))

	fileName.WriteString("rides-")
	fileName.WriteString(time.Now().Format("20060102150405"))
	fileName.WriteString(".json")

	_ = ioutil.WriteFile(fileName.String(), ridesJSON, 0644)
}
