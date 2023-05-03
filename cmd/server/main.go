package main

import (
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/compress"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/config"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/fileio"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/handlers"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.LoadServerConfig()
	if err != nil {
		log.Fatalf(`cannot create config: %s`, err)
	}

	storeTicker := cfg.StoreTicker()

	if cfg.FileStoragePath != `` {

		if cfg.Restore {
			err := fileio.LoadMetrics(cfg.FileStoragePath)
			if err != nil {
				log.Fatalf(`cannot read metrics: %s`, err)
			}
		}

		go func() {
			for {
				select {
				case <-storeTicker.C:
					err := fileio.SaveMetrics(cfg.FileStoragePath, storage.GetCounters(), storage.GetGauges())
					if err != nil {
						log.Fatalf(`cannot write metrics: %s`, err)
					}
				}
			}
		}()
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
