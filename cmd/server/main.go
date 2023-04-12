package main

import (
	"flag"
	"github.com/JinFuuMugen/go-metrics-tpl.git/cmd/server/handlers"
	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Config struct {
	Addr string `env:"ADDRESS"`
}

func main() {
	var cfg Config
	envParseError := env.Parse(&cfg)
	if envParseError != nil {
		panic(envParseError)
	}

	var serverAddr string
	if cfg.Addr != "" {
		serverAddr = cfg.Addr
	} else {
		serverAddr = *flag.String("a", "localhost:8080", "server address")
	}

	rout := chi.NewRouter()
	rout.HandleFunc(`/update/{metric_type}/{metric_name}/{metric_value}`, handlers.UpdateMetricsHandle)
	rout.HandleFunc(`/`, handlers.MainHandle)
	rout.HandleFunc(`/value/{metric_type}/{metric_name}`, handlers.GetMetricHandle)

	err := http.ListenAndServe(serverAddr, rout)
	if err != nil {
		panic(err)
	}
}
