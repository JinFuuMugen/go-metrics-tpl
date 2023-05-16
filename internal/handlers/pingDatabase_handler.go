package handlers

import (
	"context"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/database"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/logger"
	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"net/http"
	"time"
)

func PingDBHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := database.DB.PingContext(ctx); err != nil {
			logger.Errorf("error pinging database: %s", err)
			http.Error(w, "error pinging database", http.StatusInternalServerError)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			return
		}
	}
}
