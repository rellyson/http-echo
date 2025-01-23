package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHealthCheck(t *testing.T) {
	// Create a new request
	req, err := http.NewRequest("GET", "/health", nil)

	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler := http.HandlerFunc(HealthCheck)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the content type
	expectedContentType := "application/json"

	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v", contentType, expectedContentType)
	}

	// Parse the response body
	var response HealthResponse

	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	// Check the response fields
	if response.Message != "Ok" {
		t.Errorf("handler returned unexpected message: got %v want %v", response.Message, "Ok")
	}

	if response.Status != http.StatusOK {
		t.Errorf("handler returned unexpected status: got %v want %v", response.Status, http.StatusOK)
	}

	// Verify timestamp is in correct format and recent
	parsedTime, err := time.Parse(time.RFC3339, response.Timestamp)

	if err != nil {
		t.Errorf("handler returned invalid timestamp format: %v", err)
	}

	// Check if timestamp is within the last second
	if time.Since(parsedTime) > time.Second {
		t.Errorf("handler returned timestamp too old: %v", parsedTime)
	}
}
