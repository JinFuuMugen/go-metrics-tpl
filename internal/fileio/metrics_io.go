package fileio

import (
	"bufio"
	"encoding/json"
	"errors"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/models"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/storage"
	"os"
)

func SaveMetrics(filepath string, counters []storage.Counter, gauges []storage.Gauge) error {
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
		return err
	}

	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	if _, err = file.Write(jsonData); err != nil {
		return err
	}
	return nil
}

func LoadMetrics(filepath string) error {
	var metrics []models.Metrics

	file, err := os.OpenFile(filepath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return scanner.Err()
	}

	fileData := scanner.Bytes()
	err = json.Unmarshal(fileData, &metrics)
	if err != nil {
		return err
	}

	return func(metrics []models.Metrics) error {
		for _, m := range metrics {
			switch m.MType {
			case storage.MetricTypeCounter:
				storage.AddCounter(m.ID, *m.Delta)
			case storage.MetricTypeGauge:
				storage.SetGauge(m.ID, *m.Value)
			default:
				err := errors.New(`wrong metric type`)
				return err
			}
		}
		return nil
	}(metrics)
}
