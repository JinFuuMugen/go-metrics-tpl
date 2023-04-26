package sender

import (
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/config"
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
	url := `http://` + s.Addr + `/update/` + m.GetType() + `/` + m.GetName() + `/` + m.GetValueString()
	_, err := s.client.R().Post(url)
	return err
}
