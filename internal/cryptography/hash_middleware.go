package cryptography

import (
	"encoding/hex"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/config"
	"net/http"
)

func ValidateHashMiddleware(cfg *config.ServerConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			hashString := r.Header.Get("HashSHA256")
			if hashString != "" {
				body := r.Method + r.URL.Path + r.Host

				hash := GetHMACSHA256([]byte(body), cfg.Key)
				calculatedHashString := hex.EncodeToString(hash)

				if hashString != calculatedHashString {
					http.Error(w, "Несоответствие хэша", http.StatusBadRequest)
					return
				}
			}

			next.ServeHTTP(w, r)

			if cfg.Key != "" {
				responseHash := GetHMACSHA256([]byte(""), cfg.Key)
				responseHashString := hex.EncodeToString(responseHash)
				w.Header().Set("HashSHA256", responseHashString)
			}
		})
	}
}
