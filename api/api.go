package api

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/NeoHuang/corona-weg/core"
	"github.com/PuerkitoBio/goquery"
)

type Api interface {
	Name() string
	GetCurrent() (core.EpidemicMap, error)
}

type ExtractFunc func(*goquery.Document, string) core.EpidemicMap

type GeneralApi struct {
	url       string
	name      string
	extractFn ExtractFunc

	cachePeriod      time.Duration
	lastEpidemicMap  core.EpidemicMap
	lastDetectedTime time.Time
	mutex            sync.Mutex
}

func NewGeneralApi(endpointUrl string, name string, extractFn ExtractFunc, cachePeriod time.Duration) *GeneralApi {
	return &GeneralApi{
		url:         endpointUrl,
		name:        name,
		extractFn:   extractFn,
		cachePeriod: cachePeriod,
	}
}

func (api *GeneralApi) Name() string {
	return api.name
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
		return nil, fmt.Errorf("Api %q failed to get http request:%s", api.name, err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Api %q status code error: %d %s", api.name, res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Api %q failed to parse body:%s", api.name, err)
	}

	epidemicMap := api.extractFn(doc, api.name)
	api.lastEpidemicMap = epidemicMap
	api.lastDetectedTime = now

	return epidemicMap, nil
}

func NewRkiApi(cachePeriod time.Duration) Api {
	return NewGeneralApi("https://www.rki.de/DE/Content/InfAZ/N/Neuartiges_Coronavirus/Fallzahlen.html", "RKI", RkiExtractFunc, cachePeriod)
}

func NewJetztApi(cachePeriod time.Duration) Api {
	return NewGeneralApi("https://www.coronavirus.jetzt/karten/deutschland/", "Jetzt", JetztExtractFunc, cachePeriod)
}
