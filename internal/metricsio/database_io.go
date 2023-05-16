package metricsio

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/database"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/models"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/storage"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func saveMetricsDB(counters []storage.Counter, gauges []storage.Gauge) error {
	for _, c := range counters {
		_, err := database.DB.Exec("INSERT INTO metrics (type, value, delta, name) VALUES ($1, $2, $3, $4) ON CONFLICT (type, name) DO UPDATE SET value = $2, delta = $3;", c.Type, sql.NullFloat64{Valid: false, Float64: 0}, c.Value, c.Name)
		if err != nil {
			return fmt.Errorf("cannot execute query to save counters: %w", err)
		}
	}

	for _, g := range gauges {
		_, err := database.DB.Exec("INSERT INTO metrics (type, value, delta, name) VALUES ($1, $2, $3, $4) ON CONFLICT (type, name) DO UPDATE SET value = $2, delta = $3;", g.Type, g.Value, sql.NullInt64{Valid: false, Int64: 0}, g.Name)
		if err != nil {
			return fmt.Errorf("cannot execute query to save gauges: %w", err)
		}
	}
	return nil
}

func loadMetricsDB() error {
	var metrics []models.Metrics

	rows, err := database.DB.Query("SELECT name, type, value, delta from metrics")
	if err != nil {
		return fmt.Errorf("cannot read metrics from db: %w", err)
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
