package main

import (
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/config"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/logger"
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

	err = logger.Init()
	if err != nil {
		log.Fatalf("cannot initialize logger: %s", err)
	}

	pollTicker := cfg.PollTicker()
	reportTicker := cfg.ReportTicker()

	str := storage.NewStorage()
	snd := sender.NewSender(*cfg)

	m := monitors.NewRuntimeMonitor(str, snd)
	g := monitors.NewGopsutilMonitor(str, snd)

	rateLimit := cfg.RateLimit
	semaphore := make(chan bool, rateLimit)

	rateLimitTicker := time.NewTicker(time.Second / time.Duration(rateLimit))
	defer rateLimitTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			m.CollectRuntimeMetrics()
			g.CollectGopsutil()

		case <-reportTicker.C:
			select {
			case semaphore <- true:
				go func() {
					err := m.Dump()
					if err != nil {
						logger.Warnf("error dumping metrics: %w", err)
					}
					err = g.Dump()
					if err != nil {
						logger.Warnf("error dumping metrics: %w", err)
					}
					<-semaphore
				}()
			default:
				logger.Warnf("maximum concurrent Dump executions reached, skipping current dump")
			}
		}

		<-rateLimitTicker.C
	}
}
