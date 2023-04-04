package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func UpdateMetricsHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not a valid HTTP method.", http.StatusMethodNotAllowed)
		return
	}
	urlSplit := strings.Split(r.URL.String(), "/")
	if len(urlSplit) != 5 {
		http.Error(w, "Not a valid URL.", http.StatusNotFound)
		return
	}
	key := urlSplit[len(urlSplit)-2]
	switch urlSplit[2] {
	case "counter":
		value, err := strconv.ParseInt(urlSplit[len(urlSplit)-1], 10, 64)
		if err != nil {
			http.Error(w, "Not a valid metric value.", http.StatusBadRequest)
			return
		}
		MS.AddCounter(key, value)
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		counterValue, _ := MS.GetCounter(key)
		response := fmt.Sprintf("Counter value updated. Metric named %s is now %d.", key, counterValue)
		w.Write([]byte(response))
		return
	case "gauge":
		value, err := strconv.ParseFloat(urlSplit[len(urlSplit)-1], 64)
		if err != nil {
			http.Error(w, "Not a valid metric value.", http.StatusBadRequest)
			return
		}
		MS.AddGauge(key, value)
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		gaugeValue, _ := MS.GetGauge(key)
		response := fmt.Sprintf("Gauge value updated. Metric named %s is now %f.", key, gaugeValue)
		w.Write([]byte(response))
		return
	default:
		http.Error(w, "Not a valid metric.", http.StatusNotImplemented)
		return
	}
}
