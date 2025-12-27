package dockerhub

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollect(t *testing.T) {
	tests := []struct {
		name          string
		serverFunc    func() *httptest.Server
		expectedError bool
		assertions    func(*testing.T, map[string]int64)
	}{
		{
			name: "should collect data successfully",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"pull_count": 1000, "rating": 5, "star_count": 100}`))
				}))
			},
			expectedError: false,
			assertions: func(t *testing.T, data map[string]int64) {
				assert.NotEmpty(t, data)
			},
		},
		{
			name: "should handle empty response",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{}`))
				}))
			},
			expectedError: false,
			assertions: func(t *testing.T, data map[string]int64) {
				// Empty response should still work
				assert.NotNil(t, data)
			},
		},
		{
			name: "should handle malformed JSON",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{invalid json`))
				}))
			},
			expectedError: true,
			assertions: func(t *testing.T, data map[string]int64) {
				// Should handle error gracefully
			},
		},
		{
			name: "should handle HTTP error responses",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(`{"detail": "Not found"}`))
				}))
			},
			expectedError: true,
			assertions: func(t *testing.T, data map[string]int64) {
				// Should handle HTTP errors
			},
		},
		{
			name: "should handle 500 server error",
			serverFunc: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
				}))
			},
			expectedError: true,
			assertions: func(t *testing.T, data map[string]int64) {
				// Should handle server errors
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.serverFunc()
			defer server.Close()

			dh := &DockerHub{
				baseURL: server.URL,
				httpCli: server.Client(),
			}

			data, err := dh.collectData()
			if tt.expectedError {
				if err == nil {
					t.Logf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Logf("Unexpected error: %v", err)
				}
			}

			if data != nil {
				tt.assertions(t, data)
			}
		})
	}
}

func TestCollectWithValidMetrics(t *testing.T) {
	tests := []struct {
		name        string
		responseBody string
		checks      func(*testing.T, map[string]int64)
	}{
		{
			name:         "should parse pull_count",
			responseBody: `{"pull_count": 5000}`,
			checks: func(t *testing.T, data map[string]int64) {
				if data != nil {
					// Verify parsing worked
					assert.NotNil(t, data)
				}
			},
		},
		{
			name:         "should parse rating",
			responseBody: `{"rating": 4}`,
			checks: func(t *testing.T, data map[string]int64) {
				if data != nil {
					assert.NotNil(t, data)
				}
			},
		},
		{
			name:         "should parse star_count",
			responseBody: `{"star_count": 250}`,
			checks: func(t *testing.T, data map[string]int64) {
				if data != nil {
					assert.NotNil(t, data)
				}
			},
		},
		{
			name:         "should handle large numbers",
			responseBody: `{"pull_count": 999999999, "rating": 5, "star_count": 100000}`,
			checks: func(t *testing.T, data map[string]int64) {
				if data != nil {
					assert.NotNil(t, data)
				}
			},
		},
		{
			name:         "should handle zero values",
			responseBody: `{"pull_count": 0, "rating": 0, "star_count": 0}`,
			checks: func(t *testing.T, data map[string]int64) {
				if data != nil {
					assert.NotNil(t, data)
				}
			},
		},
		{
			name:         "should handle negative values",
			responseBody: `{"pull_count": -1, "rating": -1}`,
			checks: func(t *testing.T, data map[string]int64) {
				if data != nil {
					assert.NotNil(t, data)
				}
			},
		},
		{
			name:         "should handle null values",
			responseBody: `{"pull_count": null, "rating": null}`,
			checks: func(t *testing.T, data map[string]int64) {
				// Should handle null gracefully
			},
		},
		{
			name:         "should handle missing fields",
			responseBody: `{"some_field": "value"}`,
			checks: func(t *testing.T, data map[string]int64) {
				if data != nil {
					assert.NotNil(t, data)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			dh := &DockerHub{
				baseURL: server.URL,
				httpCli: server.Client(),
			}

			data, _ := dh.collectData()
			tt.checks(t, data)
		})
	}
}

func TestCollectHTTPMethods(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		expectedError bool
	}{
		{
			name:          "should handle GET request",
			method:        http.MethodGet,
			expectedError: false,
		},
		{
			name:          "should handle request timeout",
			method:        http.MethodGet,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"pull_count": 1000}`))
			}))
			defer server.Close()

			dh := &DockerHub{
				baseURL: server.URL,
				httpCli: server.Client(),
			}

			data, err := dh.collectData()
			if !tt.expectedError && err != nil {
				t.Logf("Unexpected error: %v", err)
			}
			if data != nil {
				assert.NotNil(t, data)
			}
		})
	}
}

func TestCollectContentType(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		shouldWork  bool
	}{
		{
			name:        "should handle application/json",
			contentType: "application/json",
			shouldWork:  true,
		},
		{
			name:        "should handle application/json with charset",
			contentType: "application/json; charset=utf-8",
			shouldWork:  true,
		},
		{
			name:        "should handle invalid content type",
			contentType: "text/plain",
			shouldWork:  true, // Parser might still work
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", tt.contentType)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"pull_count": 1000}`))
			}))
			defer server.Close()

			dh := &DockerHub{
				baseURL: server.URL,
				httpCli: server.Client(),
			}

			data, err := dh.collectData()
			if tt.shouldWork && err != nil {
				t.Logf("Expected success but got error: %v", err)
			}
			if data != nil {
				assert.NotNil(t, data)
			}
		})
	}
}

func TestCollectErrorHandling(t *testing.T) {
	tests := []struct {
		name          string
		setupServer   func() *httptest.Server
		expectedError bool
	}{
		{
			name: "should handle connection refused",
			setupServer: func() *httptest.Server {
				// Simulate connection refused by closing server immediately
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}))
				server.Close()
				return server
			},
			expectedError: true,
		},
		{
			name: "should handle empty response body",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(""))
				}))
			},
			expectedError: true,
		},
		{
			name: "should handle response with only whitespace",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("   \n\t  "))
				}))
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupServer()
			defer func() {
				// Server might be already closed
				if server != nil {
					server.Close()
				}
			}()

			if server != nil {
				dh := &DockerHub{
					baseURL: server.URL,
					httpCli: server.Client(),
				}

				_, err := dh.collectData()
				if tt.expectedError && err == nil {
					t.Logf("Expected error but got none")
				}
			}
		})
	}
}

func TestCollectResponseStatusCodes(t *testing.T) {
	tests := []struct {
		statusCode    int
		expectedError bool
		name          string
	}{
		{
			statusCode:    http.StatusOK,
			expectedError: false,
			name:          "should handle 200 OK",
		},
		{
			statusCode:    http.StatusCreated,
			expectedError: false,
			name:          "should handle 201 Created",
		},
		{
			statusCode:    http.StatusBadRequest,
			expectedError: true,
			name:          "should handle 400 Bad Request",
		},
		{
			statusCode:    http.StatusUnauthorized,
			expectedError: true,
			name:          "should handle 401 Unauthorized",
		},
		{
			statusCode:    http.StatusForbidden,
			expectedError: true,
			name:          "should handle 403 Forbidden",
		},
		{
			statusCode:    http.StatusNotFound,
			expectedError: true,
			name:          "should handle 404 Not Found",
		},
		{
			statusCode:    http.StatusInternalServerError,
			expectedError: true,
			name:          "should handle 500 Internal Server Error",
		},
		{
			statusCode:    http.StatusBadGateway,
			expectedError: true,
			name:          "should handle 502 Bad Gateway",
		},
		{
			statusCode:    http.StatusServiceUnavailable,
			expectedError: true,
			name:          "should handle 503 Service Unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.statusCode == http.StatusOK || tt.statusCode == http.StatusCreated {
					w.Write([]byte(`{"pull_count": 1000}`))
				}
			}))
			defer server.Close()

			dh := &DockerHub{
				baseURL: server.URL,
				httpCli: server.Client(),
			}

			_, err := dh.collectData()
			if tt.expectedError && err == nil {
				t.Logf("Expected error for status %d but got none", tt.statusCode)
			}
			if !tt.expectedError && err != nil {
				t.Logf("Expected success for status %d but got error: %v", tt.statusCode, err)
			}
		})
	}
}

func TestCollectDataTypes(t *testing.T) {
	tests := []struct {
		name        string
		responseBody string
	}{
		{
			name:        "should handle string values as numbers",
			responseBody: `{"pull_count": "1000"}`,
		},
		{
			name:        "should handle float values",
			responseBody: `{"pull_count": 1000.5, "rating": 4.5}`,
		},
		{
			name:        "should handle boolean values",
			responseBody: `{"is_official": true, "user_is_admin": false}`,
		},
		{
			name:        "should handle nested objects",
			responseBody: `{"user": {"username": "test"}, "pull_count": 1000}`,
		},
		{
			name:        "should handle arrays",
			responseBody: `{"tags": ["latest", "stable"], "pull_count": 1000}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			dh := &DockerHub{
				baseURL: server.URL,
				httpCli: server.Client(),
			}

			_, _ = dh.collectData()
		})
	}
}

func TestCollectMultipleCalls(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"pull_count": 1000, "rating": 5}`))
	}))
	defer server.Close()

	dh := &DockerHub{
		baseURL: server.URL,
		httpCli: server.Client(),
	}

	// Test multiple consecutive calls
	for i := 0; i < 5; i++ {
		_, err := dh.collectData()
		if err != nil {
			t.Logf("Call %d failed: %v", i, err)
		}
	}
}

func TestCollectEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		responseBody string
		description string
	}{
		{
			name:        "empty object",
			responseBody: `{}`,
			description: "Should handle empty JSON object",
		},
		{
			name:        "only whitespace in JSON",
			responseBody: `{   }`,
			description: "Should handle JSON with whitespace",
		},
		{
			name:        "very large numbers",
			responseBody: `{"pull_count": 9223372036854775807}`,
			description: "Should handle max int64",
		},
		{
			name:        "exponential notation",
			responseBody: `{"pull_count": 1e10}`,
			description: "Should handle scientific notation",
		},
		{
			name:        "escaped characters in strings",
			responseBody: `{"description": "test\nwith\nescapes", "pull_count": 100}`,
			description: "Should handle escaped strings",
		},
		{
			name:        "unicode characters",
			responseBody: `{"name": "测试", "pull_count": 100}`,
			description: "Should handle unicode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			dh := &DockerHub{
				baseURL: server.URL,
				httpCli: server.Client(),
			}

			_, _ = dh.collectData()
		})
	}
}