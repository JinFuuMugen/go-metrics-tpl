package main

import (
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/compress"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/config"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/handlers"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatalf(`cannot create config: %s`, err)
	}

	rout := chi.NewRouter()

	rout.Get(`/`, handlers.MainHandler)

	rout.Post(`/update/`, handlers.UpdateMetricsHandler)
	rout.Post(`/value/`, handlers.GetMetricHandler)
	rout.Post(`/update/{metric_type}/{metric_name}/{metric_value}`, handlers.UpdateMetricsPlainHandler)
	rout.Get(`/value/{metric_type}/{metric_name}`, handlers.GetMetricPlainHandler)

	if err = http.ListenAndServe(cfg.Addr, compress.GzipMiddleware(rout)); err != nil {
		log.Fatalf(`cannot start server: %s`, err)
	}
}
