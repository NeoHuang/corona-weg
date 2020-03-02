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

// crawl https://www.coronavirus.jetzt/karten/deutschland/
type JetztApi struct {
	cachePeriod      time.Duration
	lastEpidemicMap  core.EpidemicMap
	lastDetectedTime time.Time
	mutex            sync.Mutex
}

func NewJetztApi(cachePeriod time.Duration) *JetztApi {
	return &JetztApi{
		cachePeriod: cachePeriod,
	}
}

func (api *JetztApi) GetCurrent() (core.EpidemicMap, error) {
	api.mutex.Lock()
	defer api.mutex.Unlock()

	now := time.Now()
	if api.lastEpidemicMap != nil &&
		api.lastDetectedTime.Add(api.cachePeriod).After(now) {
		return api.lastEpidemicMap, nil
	}

	res, err := http.Get("https://www.coronavirus.jetzt/karten/deutschland/")
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
			}
		}
	})

	api.lastEpidemicMap = epidemicMap
	api.lastDetectedTime = now

	return epidemicMap, nil
}
