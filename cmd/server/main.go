package main

import (
	"flag"
	"github.com/JinFuuMugen/go-metrics-tpl.git/cmd/server/handlers"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	serverAddr := flag.String("a", "localhost:8080", "server address")
	flag.Parse()

	rout := chi.NewRouter()
	rout.HandleFunc(`/update/{metric_type}/{metric_name}/{metric_value}`, handlers.UpdateMetricsHandle)
	rout.HandleFunc(`/`, handlers.MainHandle)
	rout.HandleFunc(`/value/{metric_type}/{metric_name}`, handlers.GetMetricHandle)

	err := http.ListenAndServe(*serverAddr, rout)
	if err != nil {
		panic(err)
	}
}
