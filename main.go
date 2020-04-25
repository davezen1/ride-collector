package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
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

// POC - most likely websites will contain a single page with all info
// OR we will have to go to each ride page to get info
// this could be made generic with Config args for park, url, query selectors for each attribute
func buschgardens(c colly.Collector) []Ride {
	var ridesList Rides
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
				})

				c2.Visit(rideURL)
			}

		})

	})
	c.Visit("https://buschgardens.com/williamsburg/rides")

	return ridesList
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
	bgRidesList := buschgardens(*c)

	filenameJSON := write(bgRidesList)
	writeCSV(filenameJSON)
}

func write(r Rides) string {

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

	return fileName.String()
}

func writeCSV(f string) {

	jsonDataFromFile, err := ioutil.ReadFile(f)

	if err != nil {
		fmt.Println(err)
	}

	// Unmarshal JSON data
	var jsonData []Ride
	err = json.Unmarshal([]byte(jsonDataFromFile), &jsonData)

	if err != nil {
		fmt.Println(err)
	}

	csvFile, err := os.Create("./rides.csv")

	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)

	headerRow := []string{"park", "name", "description", "image", "url", "height", "tags"}

	writer.Write(headerRow)

	for _, ride := range jsonData {
		var row []string
		row = append(row, ride.Park)
		row = append(row, ride.Name)
		row = append(row, ride.ShortDescription)
		row = append(row, ride.ImageURL)
		row = append(row, ride.RideURL)
		row = append(row, ride.HeightRequirementInInches)
		row = append(row, strings.Join(ride.Tags, " "))
		writer.Write(row)
	}

	// remember to flush!
	writer.Flush()
}
