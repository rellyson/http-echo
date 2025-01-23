package middlewares

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogging(t *testing.T) {
	var buf bytes.Buffer
	log.SetOutput(&buf)

	// Create a test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test response"))
	})

	// Create a test server with the logging middleware
	handler := Logging(testHandler)
	server := httptest.NewServer(handler)
	defer server.Close()

	// Make a test request
	resp, err := http.Get(server.URL)

	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Check log output
	logOutput := buf.String()
	expectedParts := []string{
		"[INFO]",
		"method=GET",
		"status=200",
	}

	for _, part := range expectedParts {
		if !strings.Contains(logOutput, part) {
			t.Errorf("Expected log to contain '%s', got: %s", part, logOutput)
		}
	}
}

func TestResponseWriter(t *testing.T) {
	// Create a test response writer
	rec := httptest.NewRecorder()
	rw := &responseWriter{
		ResponseWriter: rec,
		statusCode:     http.StatusOK,
	}

	// Test WriteHeader
	testCode := http.StatusNotFound
	rw.WriteHeader(testCode)

	if rw.statusCode != testCode {
		t.Errorf("Expected status code %d, got %d", testCode, rw.statusCode)
	}

	if rec.Code != testCode {
		t.Errorf("Expected recorder status code %d, got %d", testCode, rec.Code)
	}
}
