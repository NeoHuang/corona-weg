package main

import (
	"log"
	"net/http"
	"time"

	"github.com/NeoHuang/corona-weg/api"
	"github.com/NeoHuang/corona-weg/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	metricsSuffix = "/metrics"
	listenAddress = ":8404"
)

func main() {
	infectionsGauge := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "conrona_weg",
		Subsystem: "deutschland",
		Name:      "status_total",
		Help:      "total number of infections/death",
	}, []string{"bundesland", "type"})
	prometheus.MustRegister(infectionsGauge)
	api := api.NewMultiApi(5 * time.Minute)
	go func() {
		for {
			epidemicMap, err := api.GetCurrent()
			if err != nil {
				log.Printf("Failed to get ticker from huobi %s", err)
			} else {
				for bundesland, epidemic := range epidemicMap {
					infectionsGauge.With(prometheus.Labels{"bundesland": bundesland, "type": "infections"}).Set(float64(epidemic.Infections))
					infectionsGauge.With(prometheus.Labels{"bundesland": bundesland, "type": "deaths"}).Set(float64(epidemic.Deaths))
				}
			}

			time.Sleep(9 * time.Minute)
		}
	}()

	// Handle Metrics endpoint
	http.Handle(metricsSuffix, promhttp.Handler())

	http.Handle("/api/epidemic", server.NewEpidemicHandler(api))

	log.Printf("Metrics exported at http://%s%s", listenAddress, metricsSuffix)
	log.Fatal(http.ListenAndServe(listenAddress, nil))
	log.Printf("Shutting down......")
}
