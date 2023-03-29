package main

import (
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	gaugeMap   map[string]float64
	counterMap map[string]int64
}

func (ms *MemStorage) Init() {
	(*ms).gaugeMap = make(map[string]float64)
	(*ms).counterMap = make(map[string]int64)
}

func (ms *MemStorage) addGauge(key string, value float64) {
	ms.gaugeMap[key] = value
}

func (ms *MemStorage) addCounter(key string, value int64) {
	ms.counterMap[key] += value
}

var MS MemStorage

func updateGauge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only post method required", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Content-Type must be text/plain", http.StatusPartialContent)
		return
	}

	urlSplit := strings.Split(r.URL.String(), "/")
	if len(urlSplit) != 5 {
		http.Error(w, "Wrong URL", http.StatusBadRequest)
		return
	}

	key := urlSplit[len(urlSplit)-2]
	value, err := strconv.ParseFloat(urlSplit[len(urlSplit)-1], 64)
	if err != nil {
		panic(err)
	}
	MS.addGauge(key, value)
	w.Header().Set("content-type", "text/plain; charset=utf-8")
}

func updateCounter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only post method required", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Content-Type must be text/plain", http.StatusPartialContent)
		return
	}

	urlSplit := strings.Split(r.URL.String(), "/")
	if len(urlSplit) != 5 {
		http.Error(w, "Wrong URL", http.StatusBadRequest)
		return
	}

	key := urlSplit[len(urlSplit)-2]
	value, err := strconv.ParseInt(urlSplit[len(urlSplit)-1], 10, 64)
	if err != nil {
		panic(err)
	}
	MS.addCounter(key, value)
	w.Header().Set("content-type", "text/plain; charset=utf-8")
}

func main() {
	MS.Init()
	mux := http.NewServeMux()

	mux.HandleFunc(`/update/counter/`, updateCounter)
	mux.HandleFunc(`/update/gauge/`, updateGauge)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
