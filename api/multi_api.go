package api

import (
	"log"
	"time"

	"github.com/NeoHuang/corona-weg/core"
)

// multi api merge all apis to find the max of all bundesland and calculate the
// total number
type MultiApi struct {
	apis []Api
}

var Bundeslaender = map[string]struct{}{
	"Baden-Württemberg":      struct{}{},
	"Bayern":                 struct{}{},
	"Berlin":                 struct{}{},
	"Brandenburg":            struct{}{},
	"Bremen":                 struct{}{},
	"Hamburg":                struct{}{},
	"Hessen":                 struct{}{},
	"Mecklenburg-Vorpommern": struct{}{},
	"Niedersachsen":          struct{}{},
	"Nordrhein-Westfalen":    struct{}{},
	"Rheinland-Pfalz":        struct{}{},
	"Saarland":               struct{}{},
	"Sachsen":                struct{}{},
	"Schleswig-Holstein":     struct{}{},
	"Thüringen":              struct{}{},
	"Sachsen-Anhalt":         struct{}{},
}

func NewMultiApi(cachePeriod time.Duration) *MultiApi {
	apis := []Api{NewRkiApi(cachePeriod), NewJetztApi(cachePeriod)}
	return &MultiApi{
		apis: apis,
	}
}

func (multiApi *MultiApi) Name() string {
	return "Multi"
}

func (multiApi *MultiApi) GetCurrent() (core.EpidemicMap, error) {
	mergedEpidemicMap := core.EpidemicMap{}
	for _, api := range multiApi.apis {
		epidemicMap, err := api.GetCurrent()
		if err != nil {
			continue
		}
		for bundesland, _ := range Bundeslaender {
			epidemic, ok := epidemicMap[bundesland]
			if !ok {
				log.Printf("Api %s miss bundesland %s", api.Name(), bundesland)
				continue
			}

			currentEpidemic := mergedEpidemicMap[bundesland]
			if currentEpidemic.Infections < epidemic.Infections {
				mergedEpidemicMap[bundesland] = epidemic
			}
		}
		for bundesland, _ := range epidemicMap {
			if _, ok := Bundeslaender[bundesland]; !ok {
				log.Printf("Api %s found unknown bundesland %s", api.Name(), bundesland)
			}
		}
	}

	var totalInfections, totalDeaths int
	for _, epidemic := range mergedEpidemicMap {
		totalInfections += epidemic.Infections
		totalDeaths += epidemic.Deaths
	}

	mergedEpidemicMap["Gesamt"] = core.Epidemic{
		Infections: totalInfections,
		Deaths:     totalDeaths,
		Timestamp:  time.Now(),
		SourceApi:  "Merged",
	}

	return mergedEpidemicMap, nil
}
