package logger

import (
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

func Initialize() zap.SugaredLogger {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	defer logger.Sync()
	sug := *logger.Sugar()
	return sug
}

func HandlerLogger(h http.HandlerFunc, sug zap.SugaredLogger) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method
		duration := time.Since(start)

		responseData := &responseData{
			status: 0,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		sug.Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
			"status", responseData.status,
			"size", responseData.size,
		)

	}
	return logFn
}

type (
	responseData struct {
		status int
		size   int
	}
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.responseData.status = statusCode
}
