package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecover(t *testing.T) {
	tests := []struct {
		name         string
		handler      http.Handler
		expectedCode int
		shouldPanic  bool
		expectedBody string
	}{
		{
			name: "normal request - no panic",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("OK"))
			}),
			expectedCode: http.StatusOK,
			shouldPanic:  false,
			expectedBody: "OK",
		},
		{
			name: "panic request",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("panic hour!!!")
			}),
			expectedCode: http.StatusInternalServerError,
			shouldPanic:  true,
			expectedBody: "Internal Server Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new request
			req := httptest.NewRequest("GET", "/test", nil)
			rec := httptest.NewRecorder()

			// Create the recover middleware
			handler := Recover(tt.handler)

			// Serve the request
			handler.ServeHTTP(rec, req)

			// Check status code
			if rec.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, rec.Code)
			}

			// Check response body
			if rec.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, rec.Body.String())
			}
		})
	}
}
