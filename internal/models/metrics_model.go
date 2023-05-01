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

//func (m *Metrics) UnmarshalJSON(data []byte) (err error) {
//	type MetricsAlias Metrics
//
//	aliasValue := &struct {
//		*MetricsAlias
//		Delta string `json:"delta,omitempty"`
//		Value string `json:"value,omitempty"`
//	}{
//		MetricsAlias: (*MetricsAlias)(m),
//	}
//	if err = json.Unmarshal(data, aliasValue); err != nil {
//		return
//	}
//	switch aliasValue.MType {
//	case storage.MetricTypeCounter:
//		var tmp int64
//		tmp, err = strconv.ParseInt(aliasValue.Delta, 10, 64)
//		m.Delta = &tmp
//	case storage.MetricTypeGauge:
//		var tmp float64
//		tmp, err = strconv.ParseFloat(aliasValue.Value, 64)
//		m.Value = &tmp
//	default:
//		err = errors.New(`error processing json: Can't find metric type`)
//	}
//	return
//}
