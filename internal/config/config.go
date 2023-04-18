package config

import (
	"flag"
	"github.com/caarlos0/env"
	"log"
	"time"
)

type Config struct {
	Addr           string `env:"ADDRESS""`
	PollInterval   int    `env:"POLL_INTERVAL"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
}

func New() (*Config, error) {
	cfg := &Config{}
	flag.StringVar(&cfg.Addr, `a`, cfg.Addr, `server address`)
	flag.IntVar(&cfg.PollInterval, `p`, cfg.PollInterval, `poll interval`)
	flag.IntVar(&cfg.ReportInterval, `r`, cfg.ReportInterval, `poll interval`)
	flag.Parse()

	if envErr := env.Parse(cfg); envErr != nil {
		log.Printf(`cannot read env config: %s`, envErr)
	}

	if cfg.Addr == `` {
		cfg.Addr = `localhost:8080`
	}

	if cfg.PollInterval == 0 {
		cfg.PollInterval = 2
	}

	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = 10
	}
	return cfg, nil
}

func (cfg *Config) PollTicker() *time.Ticker {
	return time.NewTicker(time.Duration(cfg.PollInterval) * time.Second)
}

func (cfg *Config) ReportTicker() *time.Ticker {
	return time.NewTicker(time.Duration(cfg.ReportInterval) * time.Second)
}
