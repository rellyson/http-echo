package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateStack(t *testing.T) {
	// Create test middleware that adds headers
	middleware1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Test-1", "value1")
			next.ServeHTTP(w, r)
		})
	}

	middleware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("X-Test-2", "value2")
			next.ServeHTTP(w, r)
		})
	}

	// Create final handler
	finalHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Test cases
	tests := []struct {
		name           string
		middlewares    []Middleware
		expectedHeader map[string]string
	}{
		{
			name:           "Empty middleware stack",
			middlewares:    []Middleware{},
			expectedHeader: map[string]string{},
		},
		{
			name:        "Single middleware",
			middlewares: []Middleware{middleware1},
			expectedHeader: map[string]string{
				"X-Test-1": "value1",
			},
		},
		{
			name:        "Multiple middlewares",
			middlewares: []Middleware{middleware1, middleware2},
			expectedHeader: map[string]string{
				"X-Test-1": "value1",
				"X-Test-2": "value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create middleware stack
			stack := CreateStack(tt.middlewares...)
			handler := stack(finalHandler)

			// Create test request
			req := httptest.NewRequest("GET", "/", nil)
			w := httptest.NewRecorder()

			// Execute request
			handler.ServeHTTP(w, req)

			// Check status code
			if w.Code != http.StatusOK {
				t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
			}

			// Verify headers
			for key, expectedValue := range tt.expectedHeader {
				if got := w.Header().Get(key); got != expectedValue {
					t.Errorf("expected header %s to be %s, got %s", key, expectedValue, got)
				}
			}
		})
	}
}
