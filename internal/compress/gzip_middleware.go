package compress

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"strings"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			if strings.Contains(r.Header.Get("Content-Type"), "application/json") || strings.Contains(r.Header.Get("Content-Type"), "text/html") {
				//Decode body
				reader, err := gzip.NewReader(r.Body)
				if err != nil {
					//Internal server error
					http.Error(w, `internal server error`, http.StatusInternalServerError)
					return
				}
				defer reader.Close()
				decodedBody, err := ioutil.ReadAll(reader)
				if err != nil {
					//Internal server error
					http.Error(w, `internal server error`, http.StatusInternalServerError)
					return
				}
				r.Body = ioutil.NopCloser(bytes.NewBuffer(decodedBody))
			} else {
				//Bad request
				http.Error(w, `invalid content type for gzip encoding`, http.StatusBadRequest)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
