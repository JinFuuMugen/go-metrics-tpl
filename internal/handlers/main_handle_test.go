package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMainHandle(t *testing.T) {
	tests := []struct {
		method     string
		url        string
		name       string
		wantedCode int
	}{
		{
			name:       "positive get",
			wantedCode: 200,
			method:     http.MethodGet,
			url:        "/",
		},
		{
			name:       "wrong method",
			wantedCode: 405,
			method:     http.MethodPost,
			url:        "/",
		},
		{
			name:       "wrong url",
			wantedCode: 404,
			method:     http.MethodGet,
			url:        "/123/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			r.Get(`/`, MainHandle)
			req, err := http.NewRequest(tt.method, tt.url, nil)
			if err != nil {
				t.Fatal(err)
			}
			rr := httptest.NewRecorder()
			r.ServeHTTP(rr, req)
			assert.Equal(t, tt.wantedCode, rr.Code)
		})
	}
}
