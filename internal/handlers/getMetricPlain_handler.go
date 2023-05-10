package handlers

import (
	"fmt"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/logger"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func GetMetricPlainHandler(w http.ResponseWriter, r *http.Request) {

	zapLogger := logger.GetLogger()

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
		zapLogger.Errorf("unsupported metric type")
		http.Error(w, "unsupported metric type", http.StatusNotImplemented)
		return
	}

	if err != nil {
		zapLogger.Errorf("metric is not found: %s", err)
		http.Error(w, fmt.Sprintf("metric is not found: %s", err), http.StatusNotFound)
		return
	}
	_, err = w.Write([]byte(m.GetValueString()))
	if err != nil {
		zapLogger.Errorf("cannot write response: %s", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	w.Header().Add("Content-Type", "text/plain")
}
