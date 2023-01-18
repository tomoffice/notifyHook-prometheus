package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
	go func() {
		for {
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}
func switchFlip() {
	go func() {
		sec := time.Duration(10)
		i := 0
		for {
			gf.Set(float64(i))
			i++
			time.Sleep(sec * time.Second)
			if i == 5 {
				i = 0
				gf.Set(float64(99))
				time.Sleep(sec * time.Second)
			}
		}

	}()
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "goApp_processed_ops_total",
		Help: "The total number of processed events",
	})
	gf = promauto.NewGauge(prometheus.GaugeOpts{
		Name:        "goApp_flip",
		Help:        "開開關關翻翻樂",
		ConstLabels: map[string]string{},
	})
)

func main() {
	recordMetrics()
	switchFlip()
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9001", nil))
}
