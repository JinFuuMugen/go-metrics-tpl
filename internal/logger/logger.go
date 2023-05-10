package logger

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var log zap.SugaredLogger

func Init() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return fmt.Errorf("can't initialize zap logger: %w", err)
	}
	defer logger.Sync()
	sug := *logger.Sugar()
	log = sug
	return nil
}

func GetLogger() zap.SugaredLogger {
	return log
}

func HandlerLogger(h http.HandlerFunc, log zap.SugaredLogger) http.HandlerFunc {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		uri := r.RequestURI
		method := r.Method
		duration := time.Since(start)

		responseData := &responseData{
			status: http.StatusOK,
			size:   0,
		}
		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		log.Infoln(
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
	if err != nil {
		return 0, fmt.Errorf("cannot implement ResponseWriter: %w", err)
	}
	r.responseData.size += size
	return size, nil
}
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
