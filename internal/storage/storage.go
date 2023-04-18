package storage

import (
	"fmt"
	"strconv"
	"strings"
)

const MetricTypeGauge = `gauge`
const MetricTypeCounter = `counter`

type (
	Metric interface {
		GetType() string
		GetName() string
		GetValueString() string
	}

	Storage interface {
		SetGauge(string, float64)
		AddCounter(string, int64)
		GetCounters() []Counter
		GetGauges() []Gauge
		GetCounter(string) (Counter, error)
		GetGauge(string) (Gauge, error)
	}

	Counter struct {
		Name  string
		Type  string
		Value int64
	}
	Gauge struct {
		Name  string
		Type  string
		Value float64
	}
)

func (c Counter) GetType() string {
	return c.Type
}

func (c Counter) GetName() string {
	return c.Name
}

func (c Counter) GetValueString() string {
	return strconv.FormatInt(c.Value, 10)
}

func (g Gauge) GetType() string {
	return g.Type
}

func (g Gauge) GetName() string {
	return g.Name
}

func (g Gauge) GetValueString() string {
	f := func(num float64) string {
		s := fmt.Sprintf("%.4f", num)
		return strings.TrimRight(strings.TrimRight(s, "0"), ".")
	}
	return f(g.Value)
}

func NewStorage() Storage {
	return &MemStorage{
		GaugeMap:   make(map[string]float64),
		CounterMap: make(map[string]int64),
	}
}

var defaultStorage = NewStorage()

func GetCounter(k string) (Counter, error) {
	return defaultStorage.GetCounter(k)
}

func GetGauge(k string) (Gauge, error) {
	return defaultStorage.GetGauge(k)
}

func AddCounter(k, v string) error {
	value, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return fmt.Errorf("cannot parse counter value: %w", err)
	}
	defaultStorage.AddCounter(k, value)
	return nil
}

func SetGauge(k, v string) error {
	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return fmt.Errorf("cannot parse gauge value: %w", err)
	}
	defaultStorage.SetGauge(k, value)
	return nil
}
func GetCounters() []Counter {
	return defaultStorage.GetCounters()
}

func GetGauges() []Gauge {
	return defaultStorage.GetGauges()
}
