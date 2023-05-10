package compress

import (
	"bytes"
	"compress/gzip"
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/handlers"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGzipMiddleware(t *testing.T) {

	// Gzip encode the string
	var gzippedBytes bytes.Buffer
	gzipWriter := gzip.NewWriter(&gzippedBytes)
	_, err := gzipWriter.Write([]byte(`{"id":"BuckHashSys","type":"gauge","delta":0,"value":6347}`))
	if err != nil {
		log.Fatal("error gzipping data: ", err)
	}
	err = gzipWriter.Close()
	if err != nil {
		log.Fatal("error closing gzip writer: ", err)
	}

	testCases := []struct {
		name                  string
		method                string
		contentTypeHeader     string
		contentEncodingHeader string
		acceptEncodingHeader  string
		body                  string
		expectedCode          int
		expectedBody          string
		url                   string
	}{
		{
			name:                  "Valid JSON no encoding",
			url:                   "http://localhost:8080/update/",
			method:                http.MethodPost,
			contentTypeHeader:     "application/json",
			contentEncodingHeader: "",
			acceptEncodingHeader:  "",
			body:                  `{"id":"testValue","type":"gauge","value":123.123}`,
			expectedCode:          200,
			expectedBody:          `{"id":"testValue","type":"gauge","value":123.123}`,
		},
		{
			name:                  "Valid JSON gzip encoding",
			url:                   "http://localhost:8080/update/",
			method:                http.MethodPost,
			contentTypeHeader:     "application/json",
			contentEncodingHeader: "",
			acceptEncodingHeader:  "gzip",
			body:                  `{"id":"testValue","type":"gauge","value":123.123}`,
			expectedCode:          200,
			expectedBody:          `{"id":"testValue","type":"gauge","value":123.123}`,
		},

		{
			name:                  "Invalid content type for gzip encoding",
			url:                   "http://localhost:8080/update/",
			method:                http.MethodPost,
			contentTypeHeader:     "image/png",
			contentEncodingHeader: "gzip",
			acceptEncodingHeader:  "gzip",
			body:                  `{"id":"testValue","type":"gauge","value":123.123}`,
			expectedCode:          400,
			expectedBody:          "invalid content type for gzip encoding\n",
		},
		{
			name:                  "Unsupported encoding",
			url:                   "http://localhost:8080/update/",
			method:                http.MethodPost,
			contentTypeHeader:     "application/json",
			contentEncodingHeader: "deflate",
			acceptEncodingHeader:  "deflate",
			body:                  `{"id":"testCnt","type":"counter","delta":123}`,
			expectedCode:          200,
			expectedBody:          `{"id":"testCnt","type":"counter","delta":123}`,
		},
		{
			name:                  "GZIP decoding",
			url:                   "http://localhost:8080/update/",
			method:                http.MethodPost,
			contentTypeHeader:     "application/json",
			contentEncodingHeader: "gzip",
			acceptEncodingHeader:  "gzip",
			body:                  gzippedBytes.String(),
			expectedCode:          200,
			expectedBody:          `{"id":"BuckHashSys","type":"gauge","delta":0,"value":6347}`,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			rout := chi.NewRouter()

			rout.Get(`/`, handlers.MainHandler)
			rout.Post(`/update/`, handlers.UpdateMetricsHandler)
			rout.Post(`/value/`, handlers.GetMetricHandler)
			rout.Post(`/update/{metric_type}/{metric_name}/{metric_value}`, handlers.UpdateMetricsPlainHandler)
			rout.Get(`/value/{metric_type}/{metric_name}`, handlers.GetMetricPlainHandler)

			req, err := http.NewRequest(tt.method, tt.url, bytes.NewBufferString(tt.body))
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			req.Header.Set("Content-Type", tt.contentTypeHeader)
			req.Header.Set("Content-Encoding", tt.contentEncodingHeader)
			req.Header.Set("Accept-Encoding", tt.acceptEncodingHeader)

			rr := httptest.NewRecorder()
			handler := GzipMiddleware(rout)

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.expectedCode {
				t.Errorf("expected code %d, but got %d", tt.expectedCode, rr.Code)
			}
			if strings.Contains(rr.Header().Get("Content-Encoding"), "gzip") {
				reader, _ := gzip.NewReader(rr.Body)
				data, _ := io.ReadAll(reader)
				if string(data) != tt.expectedBody {
					t.Errorf("expected body to be %q, but got %q", tt.expectedBody, string(data))
				}
			} else if rr.Body.String() != tt.expectedBody {
				t.Errorf("expected body to be %q, but got %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}
