package sender

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/config"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/cryptography"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/models"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/storage"
	"github.com/go-resty/resty/v2"
)

type Sender interface {
	Process([]storage.Counter, []storage.Gauge) error
	Compress(data []byte) ([]byte, error)
}

type values struct {
	addr   string
	client *resty.Client
	key    string
}

func NewSender(cfg config.Config) *values {
	return &values{cfg.Addr, resty.New(), cfg.Key}
}

func (v *values) Compress(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)

	_, err := w.Write(data)
	if err != nil {
		return nil, fmt.Errorf("failed write data to compress temporary buffer: %v", err)
	}

	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed compress data: %v", err)
	}
	return b.Bytes(), nil
}

func (v *values) Process(counters []storage.Counter, gauges []storage.Gauge) error {
	var err error

	var metrics []models.Metrics

	for _, c := range counters {
		cDelta := c.GetValue().(int64)
		metrics = append(metrics, models.Metrics{
			ID:    c.GetName(),
			MType: c.GetType(),
			Delta: &cDelta,
			Value: nil,
		})
	}
	for _, g := range gauges {
		gValue := g.GetValue().(float64)
		metrics = append(metrics, models.Metrics{
			ID:    g.GetName(),
			MType: g.GetType(),
			Delta: nil,
			Value: &gValue,
		})
	}
	jsonData, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("cannot serialize metric to json: %w", err)
	}
	compressedData, err := v.Compress(jsonData)
	if err != nil {
		return fmt.Errorf("error while compressing data: %w", err)
	}

	url := "http://" + v.addr + "/updates/"

	if v.key != "" {
		hash := cryptography.GetHMACSHA256(jsonData, v.key)
		hashString := hex.EncodeToString(hash)
		_, err = v.client.R().SetHeader("Content-Type", "application/json").SetHeader("Content-Encoding", "gzip").SetHeader("HashSHA256", hashString).SetBody(compressedData).Post(url)
	} else {
		_, err = v.client.R().SetHeader("Content-Type", "application/json").SetHeader("Content-Encoding", "gzip").SetBody(compressedData).Post(url)
	}

	if err != nil {
		return fmt.Errorf("cannot send HTTP-Request: %w", err)
	}
	return nil
}
