package compress

import (
	"compress/gzip"
	"net/http"
	"strings"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") && strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			if strings.Contains(r.Header.Get("Content-Type"), "application/json") || strings.Contains(r.Header.Get("Content-Type"), "text/html") {
				gz, err := gzip.NewReader(r.Body)
				if err != nil {
					http.Error(w, err.Error(), http.StatusBadRequest)
					next.ServeHTTP(w, r)
					return
				}
				defer gz.Close()
				r.Body = gz
				r.Header.Del("Content-Length")
				//r.Header.Set("Content-Encoding", "gzip")
			}
		}
		next.ServeHTTP(w, r)
	})
}
