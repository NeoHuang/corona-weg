package api

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/NeoHuang/bit-hedge/core"
	"github.com/PuerkitoBio/goquery"
)

// crawl https://www.rki.de/DE/Content/InfAZ/N/Neuartiges_Coronavirus/Fallzahlen.html
type RkiApi struct {
	cachePeriod      time.Duration
	lastEpidemicMap  core.EpidemicMap
	lastDetectedTime time.Time
	mutex            sync.Mutex
}

func NewRkiApi(cachePeriod time.Duration) *RkiApi {
	return &RkiApi{
		cachePeriod: cachePeriod,
	}
}

func (api *RkiApi) GetCurrent() (core.EpidemicMap, error) {
	api.mutex.Lock()
	defer api.mutex.Unlock()

	now := time.Now()
	if api.lastEpidemicMap != nil &&
		api.lastDetectedTime.Add(api.cachePeriod).After(now) {
		return api.lastEpidemicMap, nil
	}

	res, err := http.Get("https://www.rki.de/DE/Content/InfAZ/N/Neuartiges_Coronavirus/Fallzahlen.html")
	if err != nil {
		return nil, fmt.Errorf("failed to get http request:%s", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse body:%s", err)
	}

	epidemicMap := core.EpidemicMap{}
	doc.Find("tbody").EachWithBreak(func(i int, s *goquery.Selection) bool {
		// For each item found, get the band and title
		var bundesland string
		s.Find("td").EachWithBreak(func(i int, s *goquery.Selection) bool {
			if i%2 == 0 {
				bundesland = s.Text()
			} else {
				infections, _ := strconv.Atoi(s.Text())
				epidemicMap[bundesland] = core.Epidemic{
					Infections: infections,
					Deaths:     0,
					Timestamp:  now,
				}
			}

			if i == 1 {
				return false
			}
			return true

		})
		return false
	})

	api.lastEpidemicMap = epidemicMap
	api.lastDetectedTime = now

	return epidemicMap, nil
}
