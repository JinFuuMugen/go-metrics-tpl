package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/models"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/storage"
	"net/http"
	"strconv"
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
		value, _ := strconv.ParseFloat(m.GetValueString(), 64)
		metric.Value = &value
	case storage.MetricTypeCounter:
		m, err = storage.GetCounter(metric.ID)
		delta, _ := strconv.ParseInt(m.GetValueString(), 10, 64)
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
	_, err = w.Write(jsonBytes)
	if err != nil {
		http.Error(w, `Can't write persponse`, http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

}
