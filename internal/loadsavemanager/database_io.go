package loadsavemanager

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/database"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/models"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/storage"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"strings"
)

func saveMetricsDB(counters []storage.Counter, gauges []storage.Gauge) error {
	counterValues := make([]string, len(counters))
	counterParams := make([]interface{}, 0, len(counters)*4)

	for i, c := range counters {
		counterValues[i] = fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4)
		counterParams = append(counterParams, c.Type, sql.NullFloat64{Valid: false, Float64: 0}, c.Value, c.Name)
	}

	gaugeValues := make([]string, len(gauges))
	gaugeParams := make([]interface{}, 0, len(gauges)*4)

	for i, g := range gauges {
		gaugeValues[i] = fmt.Sprintf("($%d, $%d, $%d, $%d)", (i+len(counters))*4+1, (i+len(counters))*4+2, (i+len(counters))*4+3, (i+len(counters))*4+4)
		gaugeParams = append(gaugeParams, g.Type, g.Value, sql.NullInt64{Valid: false, Int64: 0}, g.Name)
	}

	query := fmt.Sprintf("INSERT INTO metrics (type, value, delta, name) VALUES %s ON CONFLICT (type, name) DO UPDATE SET value = excluded.value, delta = excluded.delta;", strings.Join(append(counterValues, gaugeValues...), ", "))
	params := append(counterParams, gaugeParams...)

	_, err := database.DB.Exec(query, params...)
	if err != nil {
		return fmt.Errorf("cannot execute query to save metrics: %w", err)
	}

	return nil
}

func loadMetricsDB() error {
	var metrics []models.Metrics

	rows, err := database.DB.Query("SELECT id, type, value, delta from metrics")
	if err != nil {
		return fmt.Errorf("cannot read metrics from db: %w", err)
	}
	if rows.Err() != nil {
		return fmt.Errorf("cannot read metrics from db: %w", rows.Err())
	}
	defer rows.Close()

	for rows.Next() {
		var m models.Metrics
		err = rows.Scan(&m.ID, &m.MType, &m.Value, &m.Delta)
		if err != nil {
			return fmt.Errorf("cannot scan values from db: %w", err)
		}

		metrics = append(metrics, m)
	}

	for _, m := range metrics {
		switch m.MType {
		case storage.MetricTypeCounter:
			storage.AddCounter(m.ID, *m.Delta)
		case storage.MetricTypeGauge:
			storage.SetGauge(m.ID, *m.Value)
		default:
			return fmt.Errorf("cannot opperate metric: %w", errors.New("unsupported metric type"))
		}
	}
	return nil
}
