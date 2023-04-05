package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateMetricsHandle(t *testing.T) {

	tests := []struct {
		method     string
		name       string
		url        string
		wantedCode int
	}{
		{
			name:       "positive gauge post",
			wantedCode: 200,
			method:     http.MethodPost,
			url:        "localhost:8080/update/gauge/someValue/120.414",
		},
		{
			name:       "positive counter post",
			wantedCode: 200,
			method:     http.MethodPost,
			url:        "localhost:8080/update/counter/someValue/120",
		},
		{
			name:       "wrong method",
			wantedCode: 405,
			method:     http.MethodGet,
			url:        "localhost:8080/update/counter/someValue/120",
		},
		{
			name:       "wrong url",
			wantedCode: 404,
			method:     http.MethodPost,
			url:        "localhost:8080/update",
		},
		{
			name:       "wrong metric",
			wantedCode: 501,
			method:     http.MethodPost,
			url:        "localhost:8080/update/metr/someValue/900.009",
		},
		{
			name:       "bad metric value",
			wantedCode: 400,
			method:     http.MethodPost,
			url:        "localhost:8080/update/counter/someValue/120.321",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()
			UpdateMetricsHandle(w, request)
			res := w.Result()
			assert.Equal(t, tt.wantedCode, res.StatusCode)
		})
	}
}
