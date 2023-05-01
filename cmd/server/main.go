package main

import (
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
	rout.Post(`/update/`, logger.HandlerLogger(handlers.UpdateMetricsHandler, sug))
	rout.Get(`/`, logger.HandlerLogger(handlers.MainHandler, sug))
	rout.Post(`/value/`, logger.HandlerLogger(handlers.GetMetricHandler, sug))

	if err = http.ListenAndServe(cfg.Addr, rout); err != nil {
		log.Fatalf(`cannot start server: %s`, err)
	}
}
