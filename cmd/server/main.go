package main

import (
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/compress"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/config"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/handlers"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/logger"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatalf(`cannot create config: %s`, err)
	}
	sug := logger.Initialize()

	rout := chi.NewRouter()

	rout.Get(`/`, logger.HandlerLogger(handlers.MainHandler, sug))

	rout.Post(`/update/`, logger.HandlerLogger(handlers.UpdateMetricsHandler, sug))
	rout.Post(`/value/`, logger.HandlerLogger(handlers.GetMetricHandler, sug))
	rout.Post(`/update/{metric_type}/{metric_name}/{metric_value}`, logger.HandlerLogger(handlers.UpdateMetricsPlainHandler, sug))
	rout.Get(`/value/{metric_type}/{metric_name}`, logger.HandlerLogger(handlers.GetMetricPlainHandler, sug))

	if err = http.ListenAndServe(cfg.Addr, compress.GzipMiddleware(rout)); err != nil {
		log.Fatalf(`cannot start server: %s`, err)
	}
}
