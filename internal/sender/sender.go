package sender

import (
	"encoding/json"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/config"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/models"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/storage"
	"github.com/go-resty/resty/v2"
)

type Sender interface {
	Process(storage.Metric) error
}

type sender struct {
	Addr   string
	client *resty.Client
}

func NewSender(cfg config.Config) *sender {
	return &sender{cfg.Addr, resty.New()}
}

func (s *sender) Process(m storage.Metric) error {
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

	url := `http://` + s.Addr + `/update/`
	data, err := json.Marshal(models.Metrics{
		ID:    name,
		MType: mType,
		Delta: &delta,
		Value: &value,
	})
	if err != nil {
		return err
	}
	_, err = s.client.R().SetHeader("Content-Type", "application/json").SetBody(data).Post(url)
	return err
}
