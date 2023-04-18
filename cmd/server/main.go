package main

import (
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/config"
	handlers2 "github.com/JinFuuMugen/go-metrics-tpl.git/internal/handlers"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func main() {
	cfg, cfgErr := config.LoadServerConfig()
	if cfgErr != nil {
		log.Fatalf("cannot create config: %s", cfgErr)
	}

	rout := chi.NewRouter()
	rout.Post(`/update/{metric_type}/{metric_name}/{metric_value}`, handlers2.UpdateMetricsHandle)
	rout.Get(`/`, handlers2.MainHandle)
	rout.Get(`/value/{metric_type}/{metric_name}`, handlers2.GetMetricHandle)

	servErr := http.ListenAndServe(cfg.Addr, rout)
	if servErr != nil {
		panic(servErr)
	}
}
