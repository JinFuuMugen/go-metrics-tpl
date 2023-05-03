package config

import (
	"flag"
	"os"
	"strconv"
	"time"
)

type ServerConfig struct {
	Addr            string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
}

func LoadServerConfig() (*ServerConfig, error) {
	cfg := &ServerConfig{
		Addr:            "localhost:8080",
		StoreInterval:   300,
		FileStoragePath: "tmp/metrics-db.json",
		Restore:         true,
	}
	flag.StringVar(&cfg.Addr, `a`, cfg.Addr, `server address`)
	flag.IntVar(&cfg.StoreInterval, `i`, cfg.StoreInterval, `metrics store interval(0 to sync)`)
	flag.StringVar(&cfg.FileStoragePath, `f`, cfg.FileStoragePath, `path of storage file`)
	flag.BoolVar(&cfg.Restore, `r`, cfg.Restore, `boolean to load/not saved values(also requires -f or FILE_STORAGE_PATH)`)
	flag.Parse()
	if envAddr := os.Getenv("ADDRESS"); envAddr != "" {
		cfg.Addr = envAddr
	}
	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		storeInterval, err := strconv.Atoi(envStoreInterval)
		if err == nil {
			cfg.StoreInterval = storeInterval
		}
	}
	if envFileStoragePath := os.Getenv("FILE_STORAGE_PATH"); envFileStoragePath != "" {
		cfg.FileStoragePath = envFileStoragePath
	}
	if envRestore := os.Getenv("RESTORE"); envRestore != "" {
		restore, err := strconv.ParseBool(envRestore)
		if err == nil {
			cfg.Restore = restore
		}
	}
	return cfg, nil
}

func (cfg *ServerConfig) StoreTicker() *time.Ticker {
	if cfg.StoreInterval != 0 {
		return time.NewTicker(time.Duration(cfg.StoreInterval) * time.Second)
	} else {
		return time.NewTicker(time.Duration(1) * time.Nanosecond)
	}
}
