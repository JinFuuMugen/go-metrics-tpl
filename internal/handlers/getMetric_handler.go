package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/models"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/storage"
	"net/http"
)

func GetMetricHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var metric models.Metrics
	err := decoder.Decode(&metric)
	if err != nil {
		http.Error(w, `Bad request`, http.StatusBadRequest)
		return
	}
	var m storage.Metric
	switch metric.MType {
	case storage.MetricTypeGauge:
		m, err = storage.GetGauge(metric.ID)
		value := m.GetValue().(float64)
		metric.Value = &value
	case storage.MetricTypeCounter:
		m, err = storage.GetCounter(metric.ID)
		delta := m.GetValue().(int64)
		metric.Delta = &delta
	default:
		http.Error(w, `Unsupported metric type`, http.StatusNotImplemented)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf(`Metric is not found: %s`, err), http.StatusNotFound)
		return
	}
	jsonBytes, err := json.Marshal(metric)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonBytes)
	if err != nil {
		http.Error(w, `Can't write response`, http.StatusInternalServerError)
	}
}
