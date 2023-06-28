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
	Process(storage.Metric) error
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

func (v *values) Process(m storage.Metric) error {
	var err error
	name := m.GetName()
	mType := m.GetType()
	var value float64
	var delta int64

	switch mType {
	case storage.MetricTypeGauge:
		value = m.GetValue().(float64)
	case storage.MetricTypeCounter:
		delta = m.GetValue().(int64)

	}

	data, err := json.Marshal(models.Metrics{
		ID:    name,
		MType: mType,
		Delta: &delta,
		Value: &value,
	})
	if err != nil {
		return fmt.Errorf("cannot serialize metric: %w", err)
	}
	compressedData, err := v.Compress(data)
	if err != nil {
		return fmt.Errorf("error while compressing data: %w", err)
	}

	url := "http://" + v.addr + "/update/"

	if v.key != "" {
		hash := cryptography.GetHMACSHA256(compressedData, v.key)
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
