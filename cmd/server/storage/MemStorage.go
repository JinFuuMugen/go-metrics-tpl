package storage

import "errors"

type MemStorage struct {
	GaugeMap   map[string]float64
	CounterMap map[string]int64
}

func (ms *MemStorage) AddGauge(key string, value float64) {
	if ms.GaugeMap == nil {
		ms.GaugeMap = make(map[string]float64)
	}
	ms.GaugeMap[key] = value
}

func (ms *MemStorage) AddCounter(key string, value int64) {
	if ms.CounterMap == nil {
		ms.CounterMap = make(map[string]int64)
	}
	ms.CounterMap[key] += value
}

func (ms *MemStorage) GetGauge(key string) (float64, error) {
	value, ok := (*ms).GaugeMap[key]
	if ok {
		return value, nil
	}
	return 0, errors.New("Missing key " + key)
}

func (ms *MemStorage) GetCounter(key string) (int64, error) {
	value, ok := (*ms).CounterMap[key]
	if ok {
		return value, nil
	}
	return 0, errors.New("Missing key " + key)
}
