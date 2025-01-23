package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcho(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		headers        map[string]string
		parameters     []string
		queries        map[string]string
		body           interface{}
		expectedStatus int
	}{
		{
			name:           "GET request with no parameters",
			method:         "GET",
			path:           "/",
			expectedStatus: http.StatusOK,
		},
		{
			name:   "POST request with JSON body",
			method: "POST",
			path:   "/echo",
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			body: map[string]interface{}{
				"test": "value",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "GET request with query parameters",
			method: "GET",
			path:   "/echo?param1=value1&param2=value2",
			queries: map[string]string{
				"param1": "value1",
				"param2": "value2",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET request with path parameters",
			method:         "GET",
			path:           "/echo/param1/param2",
			parameters:     []string{"param1", "param2"},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "GET request with custom status code",
			method: "GET",
			path:   "/?status=404",
			queries: map[string]string{
				"status": "404",
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBytes []byte
			if tt.body != nil {
				bodyBytes, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewBuffer(bodyBytes))

			for key, value := range tt.headers {
				req.Header.Set(key, value)
			}

			w := httptest.NewRecorder()
			Echo(w, req)

			resp := w.Result()
			defer resp.Body.Close()

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
			assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
			assert.NotEmpty(t, resp.Header.Get(httpHeaderAppName))

			var response EchoResponse
			err := json.NewDecoder(resp.Body).Decode(&response)
			assert.NoError(t, err)

			assert.NotEmpty(t, response.HostInfo.Hostname)
			assert.NotEmpty(t, response.HostInfo.IP)

			if tt.headers != nil {
				assert.NotEmpty(t, response.HttpInfo.Headers)
			}

			if tt.parameters != nil {
				assert.NotEmpty(t, response.HttpInfo.Params)
			}

			if tt.queries != nil {
				assert.NotEmpty(t, response.HttpInfo.Queries)
			}

			if tt.body != nil {
				assert.NotNil(t, response.HttpInfo.Body)
			}

			if len(req.URL.Query()) > 0 {
				assert.NotEmpty(t, response.HttpInfo.Queries)
			}
		})
	}
}

func TestMapHeaders(t *testing.T) {
	headers := http.Header{
		"Test-Header": []string{"value1"},
		"Multi-Value": []string{"value1", "value2"},
	}

	result := mapHeaders(headers)

	assert.Equal(t, "value1", result["Test-Header"])
	assert.Equal(t, "value1", result["Multi-Value"])
}

func TestMapQuery(t *testing.T) {
	url := "http://example.com?param1=value1&param2=value2"
	req, _ := http.NewRequest("GET", url, nil)

	result := mapQuery(req.URL.Query())

	assert.Equal(t, "value1", result["param1"])
	assert.Equal(t, "value2", result["param2"])
}

func TestMapPathParams(t *testing.T) {
	path := "/api/v1/users/123"
	result := mapPathParams(path)

	assert.Equal(t, []string{"api", "v1", "users", "123"}, result)
}

func TestMapBody(t *testing.T) {
	tests := []struct {
		name     string
		body     interface{}
		expected interface{}
	}{
		{
			name: "valid JSON body",
			body: map[string]interface{}{
				"key": "value",
			},
			expected: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name:     "text body",
			body:     "valid text",
			expected: "valid text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bodyBytes []byte
			if tt.body != nil {
				bodyBytes, _ = json.Marshal(tt.body)
			}

			result := mapBody(io.NopCloser(bytes.NewBuffer(bodyBytes)))
			assert.Equal(t, tt.expected, result)
		})
	}
}
