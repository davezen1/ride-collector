package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)

//Ride details
type Ride struct {
	name, shortDescription, imageURL, heightRequirementInInches, rideURLstring, tags string
}

//Rides represents slice of Ride structs
type Rides []Ride

func buschgardens(c colly.Collector) ([]Ride, [][]string) {
	var ridesList Rides
	var ridesSlices = [][]string{{"name", "image", "description", "tags", "height"}}

	c.OnHTML("#page-content > div > div > ul", func(e *colly.HTMLElement) {

		e.ForEach("li", func(_ int, el *colly.HTMLElement) {

			image := e.Request.URL.String() + el.ChildAttr("span > a > img", "src")
			name := el.ChildText("div > h2 > a")
			shortDescription := el.ChildText("div > p")
			tags := strings.Join(el.ChildTexts("div > ul > li"), ",")

			if name != "" {
				//follow link
				rideURL := e.Request.URL.Scheme + "://" + e.Request.URL.Host + e.Request.URL.Opaque + el.ChildAttr("span > a", "href")
				c2 := colly.NewCollector()
				c2.OnHTML("#page-content", func(e2 *colly.HTMLElement) {
					height := e2.ChildText("#page-content > div.container > div > div.col-sm-8 > div > ul > li:nth-child(1) > dl > dd")
					height = strings.TrimSuffix(height, "\"")

					ridesList = append(ridesList, Ride{name: name, imageURL: image, shortDescription: shortDescription, tags: tags, heightRequirementInInches: height})

					ridesSlices = append(ridesSlices, []string{name, image, shortDescription, tags, height})
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
	bgRidesList, bgRidesSlices := buschgardens(*c)

	write(bgRidesList, bgRidesSlices)
}

//experimenting with multidimensional slices and structs writing to csv
func write(r Rides, rs [][]string) {
	log.Println(r)
	ridesJSON, err := json.Marshal(r)
	if err != nil {
		log.Fatal("Cannot encode to JSON ", err)
	}
	// file, _ := json.MarshalIndent(r, "", " ")

	_ = ioutil.WriteFile("bg.json", ridesJSON, 0644)

	//write using [][]string
	file, err := os.Create("bgslice.csv")
	checkError("Cannot create file", err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, value := range rs {
		err := writer.Write(value)
		checkError("Cannot write to file", err)
	}

}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
