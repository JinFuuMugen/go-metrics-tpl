package handlers

import (
	"html/template"
	"net/http"
)

func MainHandle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Not a valid method", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("cmd/server/static/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, struct {
		GaugeMap   map[string]float64
		CounterMap map[string]int64
	}{MS.GaugeMap, MS.CounterMap})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
}