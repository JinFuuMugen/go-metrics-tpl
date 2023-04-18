package handlers

import (
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func UpdateMetricsHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, `metric_type`)
	metricName := chi.URLParam(r, `metric_name`)
	metricValue := chi.URLParam(r, `metric_value`)

	switch metricType {
	case storage.MetricTypeCounter:
		if err := storage.AddCounter(metricName, metricValue); err != nil {
			http.Error(w, `Not a valid metric value.`, http.StatusBadRequest)
			return
		}
	case storage.MetricTypeGauge:
		if err := storage.SetGauge(metricName, metricValue); err != nil {
			http.Error(w, `Not a valid metric value.`, http.StatusBadRequest)
			return
		}
	default:
		http.Error(w, `Unsupported metric type`, http.StatusNotImplemented)
		return
	}
	w.Header().Set(`content-type`, `text/plain; charset=utf-8`)
	return
}
