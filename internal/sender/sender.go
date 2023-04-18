package sender

import (
	"fmt"
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

func NewSender(addr string) *sender {
	return &sender{addr, resty.New()}
}

func (s *sender) Process(m storage.Metric) error {
	r, err := s.client.R().Post(`http://` + s.Addr + `/update/` + m.GetType() + `/` + m.GetName() + `/` + m.GetValueString())
	fmt.Println(r.StatusCode())
	return err
}
