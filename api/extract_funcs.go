package api

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/NeoHuang/corona-weg/core"
	"github.com/PuerkitoBio/goquery"
)

func JetztExtractFunc(doc *goquery.Document, apiName string) core.EpidemicMap {
	now := time.Now()
	epidemicMap := core.EpidemicMap{}
	doc.Find("tr[class^='row']").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		bundesland := s.Find(".column-1").Text()
		if bundesland != "" && bundesland != "Bundesland" {
			infections, _ := strconv.Atoi(s.Find(".column-2").Text())
			deaths, _ := strconv.Atoi(s.Find(".column-3").Text())
			epidemicMap[bundesland] = core.Epidemic{
				Infections: infections,
				Deaths:     deaths,
				Timestamp:  now,
				SourceApi:  apiName,
			}
		}
	})

	return epidemicMap
}

func RkiExtractFunc(doc *goquery.Document, apiName string) core.EpidemicMap {
	now := time.Now()
	epidemicMap := core.EpidemicMap{}
	doc.Find("tbody").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		s.Find("tr").Each(func(row int, s *goquery.Selection) {
			var bundesland string
			epidemic := core.Epidemic{
				Infections: 0,
				Deaths:     0,
				Timestamp:  now,
				SourceApi:  apiName,
			}

			s.Find("td").Each(func(i int, s *goquery.Selection) {
				switch i {
				case 0:
					bundesland = strings.Replace(s.Text(), " ", "-", -1)
					// remove soft hypen
					bundesland = strings.Replace(bundesland, "\u00AD", "", -1)
					// remove line break
					bundesland = strings.Replace(bundesland, "\n", "", -1)
				case 1:
					infectionsString := strings.Replace(s.Text(), ".", "", -1)
					infections, err := strconv.Atoi(infectionsString)
					if err != nil {
						log.Printf("failed to parse infections number for %s:%s", bundesland, err)
					}
					epidemic.Infections = infections
				case 5:
					deathString := strings.Replace(s.Text(), ".", "", -1)
					deaths, err := strconv.Atoi(deathString)
					if err != nil {
						log.Printf("failed to parse deaths number for %s:%s", bundesland, err)
					}
					epidemic.Deaths = deaths
				}

			})
			if bundesland != "" {
				epidemicMap[bundesland] = epidemic
			} else {
				log.Printf("can't parse bundesland in row %d", row)
			}
		})
		return false
	})

	return epidemicMap
}
