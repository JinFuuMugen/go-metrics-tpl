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
	zapLogger := logger.GetLogger()
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		zapLogger.Errorf("can't read request body: %s", err)
		http.Error(w, fmt.Sprintf("can't read request body: %s", err), http.StatusBadRequest)
		return
	}

	var metric models.Metrics

	err = json.Unmarshal(buf.Bytes(), &metric)
	if err != nil {
		zapLogger.Errorf("can't process body: %s", err)
		http.Error(w, fmt.Sprintf("can't process body: %s", err), http.StatusBadRequest)
		return
	}

	switch metric.MType {
	case storage.MetricTypeCounter:
		delta, err := metric.GetDelta()
		if err != nil {
			zapLogger.Errorf("cannot get delta: %s", err)
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
			zapLogger.Errorf("can't get value: %s", err)
			http.Error(w, fmt.Sprintf("bad request: %s", err), http.StatusBadRequest)
			return
		}
		storage.SetGauge(metric.ID, value)
	default:
		zapLogger.Errorf("unsupported metric type")
		http.Error(w, "unsupported metric type", http.StatusNotImplemented)
		return
	}
	jsonBytes, err := json.Marshal(metric)
	if err != nil {
		zapLogger.Errorf("can't serialize metric to json: %s", err)
		http.Error(w, fmt.Sprintf("internal server error: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(jsonBytes)
	if err != nil {
		zapLogger.Errorf("can't write response: %s", err)
		http.Error(w, fmt.Sprintf("can't write response: %s", err), http.StatusInternalServerError)
	}
}
