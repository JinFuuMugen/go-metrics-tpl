package main

import (
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/config"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/monitors"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/sender"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/storage"
	"log"
	"time"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("cannot create config: %s", err)
	}
	pollInterval := time.Duration(cfg.PollInterval) * time.Second
	reportInterval := time.Duration(cfg.ReportInterval) * time.Second
	pollTimer := cfg.PollTimer()
	reportTimer := cfg.ReportTimer()

	m := monitors.NewMonitor(storage.NewStorage(), sender.NewSender(cfg.Addr))
	for {
		select {
		case <-pollTimer.C:
			m.CollectMetrics()
			pollTimer.Reset(pollInterval)
		case <-reportTimer.C:
			m.Dump()
			reportTimer.Reset(reportInterval)
		}
	}
}
