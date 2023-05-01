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
	if err != nil {
		err := fmt.Errorf(`metric is not found: %s`, err)
		response := models.ErrorResponse{
			Message: err.Error(),
			Status:  http.StatusNotFound,
		}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.Status)
		w.Write(jsonResponse)
		return
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
	w.Header().Set("Content-Type", "application/json")
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
	w.WriteHeader(http.StatusOK)
}
