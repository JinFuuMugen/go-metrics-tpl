package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env"
)

type ServerConfig struct {
	Addr string `env:"ADDRESS"`
}

func LoadServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{}
	flag.StringVar(&cfg.Addr, `a`, cfg.Addr, `server address`)
	flag.Parse()

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf(`cannot read env config: %w`, err)
	}

	if cfg.Addr == `` {
		cfg.Addr = `localhost:8080`
	}

	return cfg, nil
}
