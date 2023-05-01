package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/models"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/storage"
	"net/http"
	"strconv"
)

func UpdateMetricsHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, `Can't read request body.`, http.StatusBadRequest)
		return
	}

	var metric models.Metrics

	err = json.Unmarshal(buf.Bytes(), &metric)
	if err != nil {
		http.Error(w, `Can't process body.`, http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case storage.MetricTypeCounter:
		storage.AddCounter(metric.ID, *metric.Delta)
		tmpCounter, _ := storage.GetCounter(metric.ID)
		delta, _ := strconv.ParseInt(tmpCounter.GetValueString(), 10, 64)
		metric.Delta = &delta
	case storage.MetricTypeGauge:
		storage.SetGauge(metric.ID, *metric.Value)
	default:
		http.Error(w, `Unsupported metric type.`, http.StatusNotImplemented)
		return
	}
	jsonBytes, err := json.Marshal(metric)
	if err != nil {
		http.Error(w, `Internal server error`, http.StatusInternalServerError)
		return
	}
	_, err = w.Write(jsonBytes)
	if err != nil {
		http.Error(w, `Can't write persponse`, http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
