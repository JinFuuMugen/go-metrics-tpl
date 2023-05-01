package handlers

import (
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func UpdateMetricsPlainHandler(w http.ResponseWriter, r *http.Request) {
	metricType := chi.URLParam(r, `metric_type`)
	metricName := chi.URLParam(r, `metric_name`)
	metricValue := chi.URLParam(r, `metric_value`)

	switch metricType {
	case storage.MetricTypeCounter:
		v, err := strconv.ParseInt(metricValue, 10, 64)
		if err != nil {
			http.Error(w, `Not a valid metric value.`, http.StatusBadRequest)
			return
		}
		storage.AddCounter(metricName, v)
	case storage.MetricTypeGauge:
		v, err := strconv.ParseFloat(metricValue, 64)
		if err != nil {
			http.Error(w, `Not a valid metric value.`, http.StatusBadRequest)
			return
		}
		storage.SetGauge(metricName, v)
	default:
		http.Error(w, `Unsupported metric type`, http.StatusNotImplemented)
		return
	}
	w.Header().Set(`content-type`, `text/plain; charset=utf-8`)
	w.WriteHeader(http.StatusOK)
}
