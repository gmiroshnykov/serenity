package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	fooGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "foo",
	}, []string{"region"})

	barGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "bar",
	}, []string{"cluster"})
)

func index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "404 Not Found")
}

func toggler(gauge *prometheus.GaugeVec, val float64, labelValue string) (func(http.ResponseWriter, *http.Request)) {
	// a quick hack to initialize some metrics on startup
	gauge.WithLabelValues(labelValue).Set(val)

	return func (w http.ResponseWriter, r *http.Request) {
		gauge.WithLabelValues(labelValue).Set(val)
		w.WriteHeader(200)
	}
}

func main() {
	r := prometheus.NewRegistry()
	r.MustRegister(fooGauge, barGauge)
	http.Handle("/metrics", promhttp.HandlerFor(r, promhttp.HandlerOpts{}))

	http.HandleFunc("/foo/on", toggler(fooGauge, 1, "eu-west-1"))
	http.HandleFunc("/foo/off", toggler(fooGauge, 0, "eu-west-1"))
	http.HandleFunc("/bar/on", toggler(barGauge, 1, "europe-01"))
	http.HandleFunc("/bar/off", toggler(barGauge, 0, "europe-01"))
	http.HandleFunc("/", index)

	fmt.Println("Listening on http://0.0.0.0:8080")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
