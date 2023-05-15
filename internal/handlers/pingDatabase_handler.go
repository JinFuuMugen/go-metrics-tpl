package handlers

import (
	"database/sql"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/config"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/logger"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
)

func PingDBHandler(cfg *config.ServerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dsn := cfg.DatabaseDSN
		db, err := sql.Open("pgx", dsn)
		defer db.Close()
		if err != nil {
			logger.Errorf("can't connect to database: %s", err)
			http.Error(w, "can't connect to database", http.StatusInternalServerError)
		}
	}
}
