package api

import (
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
		var bundesland string
		s.Find("tr").Each(func(_ int, s *goquery.Selection) {
			s.Find("td").EachWithBreak(func(i int, s *goquery.Selection) bool {
				if i%2 == 0 {
					bundesland = strings.Replace(s.Text(), " ", "-", -1)
				} else {
					infections, _ := strconv.Atoi(s.Text())
					epidemicMap[bundesland] = core.Epidemic{
						Infections: infections,
						Deaths:     0,
						Timestamp:  now,
						SourceApi:  apiName,
					}
				}

				if i == 1 {
					return false
				}
				return true

			})
		})
		return false
	})

	return epidemicMap
}
