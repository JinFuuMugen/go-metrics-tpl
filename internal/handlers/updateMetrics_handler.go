package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/models"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/storage"
	"net/http"
	"strconv"
)

func UpdateMetricsHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)

	if err != nil {
		err := fmt.Errorf(`bad request`)
		response := models.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.Status)
		w.Write(jsonResponse)
		return
	}

	var metric models.Metrics

	err = json.Unmarshal(buf.Bytes(), &metric)
	if err != nil {
		err := fmt.Errorf(`bad request`)
		response := models.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.Status)
		w.Write(jsonResponse)
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
		if err != nil {
			err := fmt.Errorf(`unsupported metric type`)
			response := models.ErrorResponse{
				Message: err.Error(),
				Status:  http.StatusNotImplemented,
			}
			jsonResponse, _ := json.Marshal(response)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(response.Status)
			w.Write(jsonResponse)
			return
		}
	}
	jsonBytes, err := json.Marshal(metric)
	if err != nil {
		err := fmt.Errorf(`internal server error`)
		response := models.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.Status)
		w.Write(jsonResponse)
		return
	}
	_, err = w.Write(jsonBytes)
	if err != nil {
		err := fmt.Errorf(`internal server error`)
		response := models.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.Status)
		w.Write(jsonResponse)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
