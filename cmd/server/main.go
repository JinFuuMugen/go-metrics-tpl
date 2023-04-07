package main

import (
	"github.com/JinFuuMugen/go-metrics-tpl.git/cmd/server/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	rout := chi.NewRouter()

	rout.HandleFunc(`/update/{metric_type}/{metric_name}/{metric_value}`, handlers.UpdateMetricsHandle)
	rout.HandleFunc(`/`, handlers.MainHandle)
	rout.HandleFunc(`/value/{metric_type}/{metric_name}`, handlers.GetMetricHandle)

	err := http.ListenAndServe(`:8080`, rout)
	if err != nil {
		panic(err)
	}
}
