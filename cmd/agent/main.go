package main

import (
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/config"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/logger"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/monitors"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/sender"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/storage"
	"log"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("cannot create config: %s", err)
	}

	err = logger.Init()
	if err != nil {
		log.Fatalf("cannot initialize logger: %s", err)
	}

	pollTicker := cfg.PollTicker()
	reportTicker := cfg.ReportTicker()

	m := monitors.NewMonitor(storage.NewStorage(), sender.NewSender(*cfg))
	for {
		select {
		case <-pollTicker.C:
			m.CollectMetrics()
		case <-reportTicker.C:
			err := m.Dump()
			if err != nil {
				logger.Warnf("error dumping metrics: %w", err)
			}
		}
	}
}
