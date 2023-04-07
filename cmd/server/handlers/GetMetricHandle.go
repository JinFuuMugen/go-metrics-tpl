package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetMetricHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Not a valid method", http.StatusMethodNotAllowed)
		return
	}

	metricType := chi.URLParam(r, "metric_type")
	metricName := chi.URLParam(r, "metric_name")

	var answer string
	switch metricType {
	case "gauge":
		v, err := MS.GetGauge(metricName)
		if err != nil {
			http.Error(w, "Metric is not found.", http.StatusNotFound)
			return
		}
		answer = fmt.Sprintf("Metric of type gauge named %s is %f", metricName, v)
		w.Write([]byte(answer))
		w.Header().Add("Content-Type", "text/plain")
		return
	case "counter":
		v, err := MS.GetCounter(metricName)
		if err != nil {
			http.Error(w, "Metric is not found.", http.StatusNotFound)
			return
		}
		answer = fmt.Sprintf("Metric of type counter named %s is %d", metricName, v)
		w.Write([]byte(answer))
		w.Header().Add("Content-Type", "text/plain")
		return
	default:
		http.Error(w, "Metric is not found.", http.StatusNotFound)
		return
	}
}
