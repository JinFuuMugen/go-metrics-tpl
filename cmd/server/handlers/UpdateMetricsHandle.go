package handlers

import (
	"github.com/JinFuuMugen/go-metrics-tpl.git/cmd/server/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func UpdateMetricsHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not a valid HTTP method.", http.StatusMethodNotAllowed)
		return
	}
	metricType := chi.URLParam(r, "metric_type")
	metricName := chi.URLParam(r, "metric_name")
	metricValue := chi.URLParam(r, "metric_value")

	switch metricType {
	case "counter":
		value, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, "Not a valid metric value.", http.StatusBadRequest)
			return
		}
		storage.MS.AddCounter(metricName, value)
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		return
	case "gauge":
		value, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, "Not a valid metric value.", http.StatusBadRequest)
			return
		}
		storage.MS.AddGauge(metricName, value)
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		return
	default:
		http.Error(w, "Not a valid metric.", http.StatusNotImplemented)
		return
	}
}
