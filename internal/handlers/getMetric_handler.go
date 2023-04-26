package handlers

import (
	"fmt"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetMetricHandler(w http.ResponseWriter, r *http.Request) {

	metricType := chi.URLParam(r, `metric_type`)
	metricName := chi.URLParam(r, `metric_name`)

	var m storage.Metric
	var err error
	switch metricType {
	case storage.MetricTypeGauge:
		m, err = storage.GetGauge(metricName)
	case storage.MetricTypeCounter:
		m, err = storage.GetCounter(metricName)
	default:
		http.Error(w, `Unsupported metric type`, http.StatusNotImplemented)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf(`Metric is not found: %s`, err), http.StatusNotFound)
		return
	}
	w.Write([]byte(m.GetValueString()))
	w.Header().Add(`Content-Type`, `text/plain`)
}
