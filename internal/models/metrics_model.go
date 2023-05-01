package models

type Metrics struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (m *Metrics) SetDelta(delta int64) {
	m.Delta = &delta
}
func (m *Metrics) SetValue(value float64) {
	m.Value = &value
}
