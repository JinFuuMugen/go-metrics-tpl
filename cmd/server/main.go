package main

import (
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/compress"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/config"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/cryptography"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/database"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/handlers"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/io"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/logger"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatalf("cannot create config: %s", err)
	}

	if err := logger.Init(); err != nil {
		log.Fatalf("cannot create logger: %s", err)
	}

	if cfg.DatabaseDSN != "" {
		err := database.InitDatabase(cfg.DatabaseDSN)
		if err != nil {
			log.Fatalf("cannot create database connection: %s", err)
		}
	}

	if err := io.Run(cfg); err != nil {
		logger.Fatalf("cannot load preload metrics: %s", err)
	}

	rout := chi.NewRouter()

	rout.Get("/", handlers.MainHandler)

	rout.Get("/ping", handlers.PingDBHandler())

	rout.Route("/updates", func(r chi.Router) {
		r.Use(io.GetDumperMiddleware(cfg))
		r.Use(cryptography.ValidateHashMiddleware(cfg))
		r.Post("/", handlers.UpdateBatchMetricsHandler)
	})

	rout.Route("/update", func(r chi.Router) {
		r.Use(io.GetDumperMiddleware(cfg))
		r.Use(cryptography.ValidateHashMiddleware(cfg))
		r.Post("/", handlers.UpdateMetricsHandler)
		r.Post("/{metric_type}/{metric_name}/{metric_value}", handlers.UpdateMetricsPlainHandler)
	})

	rout.Post("/value/", handlers.GetMetricHandler)
	rout.Get("/value/{metric_type}/{metric_name}", handlers.GetMetricPlainHandler)

	if err = http.ListenAndServe(cfg.Addr, compress.GzipMiddleware(rout)); err != nil {
		logger.Fatalf("cannot start server: %s", err)
	}
}
