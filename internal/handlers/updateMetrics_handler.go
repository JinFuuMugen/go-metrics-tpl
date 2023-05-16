package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/logger"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/models"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/storage"
	"net/http"
)

func UpdateMetricsHandler(w http.ResponseWriter, r *http.Request) {

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		logger.Errorf("cannot read request body: %s", err)
		http.Error(w, fmt.Sprintf("cannot read request body: %s", err), http.StatusBadRequest)
		return
	}

	var metric models.Metrics
	err = json.Unmarshal(buf.Bytes(), &metric)
	if err != nil {
		logger.Errorf("cannot process body: %s", err)
		http.Error(w, fmt.Sprintf("cannot process body: %s", err), http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case storage.MetricTypeCounter:
		delta, err := metric.GetDelta()
		if err != nil {
			logger.Errorf("cannot get delta: %s", err)
			http.Error(w, fmt.Sprintf("bad request: %s", err), http.StatusBadRequest)
			return
		}
		storage.AddCounter(metric.ID, delta)
		tmpCounter, _ := storage.GetCounter(metric.ID)
		deltaNew := tmpCounter.GetValue().(int64)
		metric.SetDelta(deltaNew)
	case storage.MetricTypeGauge:
		value, err := metric.GetValue()
		if err != nil {
			logger.Errorf("cannot get value: %s", err)
			http.Error(w, fmt.Sprintf("bad request: %s", err), http.StatusBadRequest)
			return
		}
		storage.SetGauge(metric.ID, value)
	default:
		logger.Errorf("unsupported metric type")
		http.Error(w, "unsupported metric type", http.StatusNotImplemented)
		return
	}

	jsonBytes, err := json.Marshal(metric)
	if err != nil {
		logger.Errorf("cannot serialize metric to json: %s", err)
		http.Error(w, fmt.Sprintf("internal server error: %s", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonBytes)
	if err != nil {
		logger.Fatalf("cannot write response: %s", err)
	}
}
