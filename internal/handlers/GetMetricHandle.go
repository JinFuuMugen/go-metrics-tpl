package handlers

import (
	"fmt"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strings"
)

func cutZeroes(num float64) string {
	s := fmt.Sprintf("%.3f", num)
	return strings.TrimRight(strings.TrimRight(s, "0"), ".")
}

func GetMetricHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Not a valid method", http.StatusMethodNotAllowed)
		return
	}
	metricType := chi.URLParam(r, "metric_type")
	metricName := chi.URLParam(r, "metric_name")

	switch metricType {
	case "gauge":
		v, err := storage.MS.GetGauge(metricName)
		if err != nil {
			http.Error(w, "Metric is not found.", http.StatusNotFound)
			return
		}
		w.Write([]byte(cutZeroes(v)))
		w.Header().Add("Content-Type", "text/plain")
		return
	case "counter":
		v, err := storage.MS.GetCounter(metricName)
		if err != nil {
			http.Error(w, "Metric is not found.", http.StatusNotFound)
			return
		}
		w.Write([]byte(fmt.Sprintf("%d", v)))
		w.Header().Add("Content-Type", "text/plain")
		return
	default:
		http.Error(w, "Metric is not found.", http.StatusNotFound)
		return
	}
}
