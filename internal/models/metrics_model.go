package models

import "strconv"

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (m *Metrics) SetDelta(deltaString string) error {
	delta, err := strconv.ParseInt(deltaString, 10, 64)
	if err != nil {
		m.Delta = &delta
	}
	return err
}
func (m *Metrics) SetValue(valueString string) error {
	value, err := strconv.ParseFloat(valueString, 64)
	if err != nil {
		m.Value = &value
	}
	return err
}
