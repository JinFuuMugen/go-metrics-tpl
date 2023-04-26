package handlers

import (
	"github.com/JinFuuMugen/go-metrics-tpl.git/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMetricHandle(t *testing.T) {
	tests := []struct {
		method      string
		name        string
		url         string
		wantedCode  int
		wantedValue string
	}{
		{
			name:        `positive gauge get`,
			wantedCode:  200,
			method:      http.MethodGet,
			url:         `/value/gauge/someG`,
			wantedValue: `123.123`,
		},
		{
			name:        `positive counter get`,
			wantedCode:  200,
			method:      http.MethodGet,
			url:         `/value/counter/someC`,
			wantedValue: `123`,
		},
		{
			name:       `wrong method`,
			wantedCode: 405,
			method:     http.MethodPost,
			url:        `/value/counter/someValue`,
		},
		{
			name:       `wrong url`,
			wantedCode: 404,
			method:     http.MethodGet,
			url:        `/updat/gauge/some`,
		},
		{
			name:       `wrong metric`,
			wantedCode: 501,
			method:     http.MethodGet,
			url:        `/value/metr/someValue`,
		},
	}

	storage.SetGauge(`someG`, `123.123`)
	storage.AddCounter(`someC`, `123`)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get(`/value/{metric_type}/{metric_name}`, GetMetricHandler)
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			assert.Equal(t, tt.wantedCode, rr.Code)
			if tt.wantedCode == 200 {
				assert.Equal(t, tt.wantedValue, rr.Body.String())
			}
		})
	}
}
