package handlers

import (
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/storage"
	"html/template"
	"net/http"
)

func MainHandle(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("internal/static/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, struct {
		Gauges   []storage.Gauge
		Counters []storage.Counter
	}{storage.GetGauges(), storage.GetCounters()})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
}
