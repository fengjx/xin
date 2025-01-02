package xin_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fengjx/xin"
)

func TestMuxBasicRouting(t *testing.T) {
	mux := xin.NewMux()

	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
	}{
		{
			name:           "basic GET route",
			method:         "GET",
			path:           "/hello",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "route not found",
			method:         "GET",
			path:           "/notfound",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "method not allowed",
			method:         "POST",
			path:           "/hello",
			expectedStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d; got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
