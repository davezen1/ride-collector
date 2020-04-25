![Go](https://github.com/davezen1/ride-collector/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/davezen1/ride-collector)](https://goreportcard.com/report/davezen1/ride-collector)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/davezen1/ride-collector)
![GitHub](https://img.shields.io/github/license/davezen1/ride-collector)

# Ride Collector

Initial project to find roller coaster ride information especially height from websites. Uses[Go Colly](http://go-colly.org/) for screen scraping. Eventually, will include writers such as CSV, Firebase.

## Development 


```
go run main.go
```

## RoadMap

- Proof of concept 
- introduce interface to collect ride info and write
- separate out interface by ride for contributions
- tests are a good thing
- add a writer to csv
- add a writer to firebase
- add a writer to json

Sample JSON output (snipper)

```
{
        "park": "Busch Gardens Williamsburg",
        "name": "Apollo's Chariot",
        "description": "A guest favorite. Hold on as this coaster drops from 210 feet!",
        "image": "https://buschgardens.com/williamsburg/rides//williamsburg/-/media/busch-gardens-williamsburg/images/rides/apollos-chariot/listing-images/apollos-chariot.ashx?version=1_201905015448",
        "height": "52",
        "url": "https://buschgardens.com/williamsburg/roller-coasters/apollos-chariot/",
        "tags": [
            "High Thrill",
            "Roller Coaster"
        ]
    },
```

