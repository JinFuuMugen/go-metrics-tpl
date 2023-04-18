package config

import (
	"flag"
	"github.com/caarlos0/env"
	"log"
)

type ServerConfig struct {
	Addr string `env:"ADDRESS"`
}

func LoadServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{}
	flag.StringVar(&cfg.Addr, `a`, cfg.Addr, `server address`)
	flag.Parse()

	if envErr := env.Parse(cfg); envErr != nil {
		log.Printf("cannot read env config: %s", envErr)
	}

	if cfg.Addr == "" {
		cfg.Addr = "localhost:8080"
	}

	return cfg, nil
}
