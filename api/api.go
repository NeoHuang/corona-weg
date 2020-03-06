package api

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/NeoHuang/bit-hedge/core"
	"github.com/PuerkitoBio/goquery"
)

type Api interface {
	GetCurrent() (core.EpidemicMap, error)
}

type ExtractFunc func(*goquery.Document) core.EpidemicMap

type GeneralApi struct {
	url       string
	extractFn ExtractFunc

	cachePeriod      time.Duration
	lastEpidemicMap  core.EpidemicMap
	lastDetectedTime time.Time
	mutex            sync.Mutex
}

func NewGeneralApi(endpointUrl string, extractFn ExtractFunc, cachePeriod time.Duration) *GeneralApi {
	return &GeneralApi{
		url:         endpointUrl,
		extractFn:   extractFn,
		cachePeriod: cachePeriod,
	}
}

func (api *GeneralApi) GetCurrent() (core.EpidemicMap, error) {
	api.mutex.Lock()
	defer api.mutex.Unlock()

	now := time.Now()
	if api.lastEpidemicMap != nil &&
		api.lastDetectedTime.Add(api.cachePeriod).After(now) {
		return api.lastEpidemicMap, nil
	}

	res, err := http.Get(api.url)
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

	epidemicMap := api.extractFn(doc)
	api.lastEpidemicMap = epidemicMap
	api.lastDetectedTime = now

	return epidemicMap, nil
}

func NewRkiApi(cachePeriod time.Duration) Api {
	return NewGeneralApi("https://www.rki.de/DE/Content/InfAZ/N/Neuartiges_Coronavirus/Fallzahlen.html", RkiExtractFunc, cachePeriod)
}

func NewJetztApi(cachePeriod time.Duration) Api {
	return NewGeneralApi("https://www.coronavirus.jetzt/karten/deutschland/", JetztExtractFunc, cachePeriod)
}
