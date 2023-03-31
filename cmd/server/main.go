package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	GaugeMap   map[string]float64
	CounterMap map[string]int64
}

func (ms *MemStorage) Init() {
	(*ms).GaugeMap = make(map[string]float64)
	(*ms).CounterMap = make(map[string]int64)
}

func (ms *MemStorage) AddGauge(key string, value float64) {
	ms.GaugeMap[key] = value
}

func (ms *MemStorage) AddCounter(key string, value int64) {
	ms.CounterMap[key] += value
}

func (ms *MemStorage) GetGauge(key string) (float64, error) {
	value, ok := (*ms).GaugeMap[key]
	if ok {
		return value, nil
	}
	return 0, errors.New("Missing key " + key)
}

func (ms *MemStorage) GetCounter(key string) (int64, error) {
	value, ok := (*ms).CounterMap[key]
	if ok {
		return value, nil
	}
	return 0, errors.New("Missing key " + key)
}

var MS MemStorage

func PostGaugeHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not a valid HTTP method.", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Not a valid content-type.", http.StatusUnsupportedMediaType)
		return
	}

	urlSplit := strings.Split(r.URL.String(), "/")
	if len(urlSplit) != 5 {
		http.Error(w, "Not a valid URL.", http.StatusBadRequest)
		return
	}

	key := urlSplit[len(urlSplit)-2]
	value, err := strconv.ParseFloat(urlSplit[len(urlSplit)-1], 64)
	if err != nil {
		http.Error(w, "Not a valid metric value.", http.StatusBadRequest)
		return
	}
	MS.AddGauge(key, value)
	w.Header().Set("content-type", "text/plain; charset=utf-8")
	gaugeValue, _ := MS.GetGauge(key)
	response := fmt.Sprintf("Gauge value updated. Metric named %s is now %5f.", key, gaugeValue)
	w.Write([]byte(response))
}

func PostCounterHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not a valid HTTP method.", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Not a valid content-type.", http.StatusUnsupportedMediaType)
		return
	}

	urlSplit := strings.Split(r.URL.String(), "/")
	if len(urlSplit) != 5 {
		http.Error(w, "Not a valid URL.", http.StatusBadRequest)
		return
	}

	key := urlSplit[len(urlSplit)-2]
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
}

func MainHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<!DOCTYPE HTML>\n<html>\n <head>\n <title>Empty Page</title>\n<p> This page is empty for now</p> <meta charset=\"utf-8\">\n </head>\n <body>\n\n</body>\n</html>"))
	return
}

func main() {
	MS.Init()
	mux := http.NewServeMux()

	mux.HandleFunc(`/`, MainHandle)
	mux.HandleFunc(`/update/counter/`, PostCounterHandle)
	mux.HandleFunc(`/update/gauge/`, PostGaugeHandle)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
