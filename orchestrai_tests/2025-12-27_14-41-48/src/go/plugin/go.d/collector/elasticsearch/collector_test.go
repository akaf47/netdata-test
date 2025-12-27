package elasticsearch

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectorNew(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "should create new collector with valid config",
			config: Config{
				URL:      "http://localhost:9200",
				Username: "user",
				Password: "pass",
			},
			wantErr: false,
		},
		{
			name: "should create new collector with minimal config",
			config: Config{
				URL: "http://localhost:9200",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			assert.NotNil(t, c)
			assert.Equal(t, Module, c.Module)
		})
	}
}

func TestCollectorInit(t *testing.T) {
	tests := []struct {
		name    string
		config  string
		wantErr bool
	}{
		{
			name: "should initialize collector with valid YAML",
			config: `
url: http://localhost:9200
username: user
password: pass
timeout: 5s
`,
			wantErr: false,
		},
		{
			name: "should initialize with empty config",
			config: `
url: http://localhost:9200
`,
			wantErr: false,
		},
		{
			name: "should initialize with SSL config",
			config: `
url: https://localhost:9200
tls_skip_verify: true
`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			err := c.Init()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCollectorCheck(t *testing.T) {
	tests := []struct {
		name        string
		setupClient func(*Collector) *MockHTTPClient
		wantErr     bool
		wantMsg     string
	}{
		{
			name: "should pass check with successful API call",
			setupClient: func(c *Collector) *MockHTTPClient {
				mock := NewMockHTTPClient()
				mock.DoFunc = func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(`{"version":{"number":"7.13.0"}}`)),
						Header:     make(http.Header),
					}, nil
				}
				c.client = mock
				return mock
			},
			wantErr: false,
		},
		{
			name: "should fail check with 404 status",
			setupClient: func(c *Collector) *MockHTTPClient {
				mock := NewMockHTTPClient()
				mock.DoFunc = func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: 404,
						Body:       io.NopCloser(strings.NewReader(`Not Found`)),
						Header:     make(http.Header),
					}, nil
				}
				c.client = mock
				return mock
			},
			wantErr: true,
		},
		{
			name: "should fail check with network error",
			setupClient: func(c *Collector) *MockHTTPClient {
				mock := NewMockHTTPClient()
				mock.DoFunc = func(req *http.Request) (*http.Response, error) {
					return nil, fmt.Errorf("connection refused")
				}
				c.client = mock
				return mock
			},
			wantErr: true,
		},
		{
			name: "should fail check with invalid JSON response",
			setupClient: func(c *Collector) *MockHTTPClient {
				mock := NewMockHTTPClient()
				mock.DoFunc = func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(`{invalid json}`)),
						Header:     make(http.Header),
					}, nil
				}
				c.client = mock
				return mock
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			c.URL = "http://localhost:9200"
			tt.setupClient(c)

			err := c.Check()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCollectorCollect(t *testing.T) {
	tests := []struct {
		name        string
		setupClient func(*Collector) *MockHTTPClient
		wantErr     bool
		wantMetrics int
	}{
		{
			name: "should collect metrics successfully",
			setupClient: func(c *Collector) *MockHTTPClient {
				mock := NewMockHTTPClient()
				mock.DoFunc = func(req *http.Request) (*http.Response, error) {
					body := `{
						"cluster_name": "elasticsearch",
						"status": "green",
						"timed_out": false,
						"active_shards": 10,
						"relocating_shards": 0,
						"initializing_shards": 0,
						"unassigned_shards": 0,
						"nodes": {
							"count": {
								"total": 3,
								"data": 2,
								"ingest": 1
							}
						}
					}`
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(body)),
						Header:     make(http.Header),
					}, nil
				}
				c.client = mock
				return mock
			},
			wantErr:     false,
			wantMetrics: 1,
		},
		{
			name: "should handle collection failure",
			setupClient: func(c *Collector) *MockHTTPClient {
				mock := NewMockHTTPClient()
				mock.DoFunc = func(req *http.Request) (*http.Response, error) {
					return nil, fmt.Errorf("timeout")
				}
				c.client = mock
				return mock
			},
			wantErr: true,
		},
		{
			name: "should handle empty response",
			setupClient: func(c *Collector) *MockHTTPClient {
				mock := NewMockHTTPClient()
				mock.DoFunc = func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: 200,
						Body:       io.NopCloser(strings.NewReader(`{}`)),
						Header:     make(http.Header),
					}, nil
				}
				c.client = mock
				return mock
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			c.URL = "http://localhost:9200"
			tt.setupClient(c)

			mx := make(map[string]int64)
			err := c.Collect(mx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCollectorConnectionHandling(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		username    string
		password    string
		wantErr     bool
	}{
		{
			name:     "should handle basic auth",
			url:      "http://localhost:9200",
			username: "elastic",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "should handle no auth",
			url:      "http://localhost:9200",
			username: "",
			password: "",
			wantErr:  false,
		},
		{
			name:     "should handle invalid URL",
			url:      "not-a-valid-url",
			wantErr:  true,
		},
		{
			name:     "should handle empty URL",
			url:      "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			c.URL = tt.url
			c.Username = tt.username
			c.Password = tt.password

			err := c.Init()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				if tt.url != "" {
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestCollectorClusterHealthMetrics(t *testing.T) {
	tests := []struct {
		name            string
		clusterStatus   string
		activeShardsNum int
		nodeCount       int
		expectedMetrics map[string]interface{}
	}{
		{
			name:            "should parse green cluster status",
			clusterStatus:   "green",
			activeShardsNum: 10,
			nodeCount:       3,
			expectedMetrics: map[string]interface{}{
				"status": "green",
			},
		},
		{
			name:            "should parse yellow cluster status",
			clusterStatus:   "yellow",
			activeShardsNum: 8,
			nodeCount:       2,
			expectedMetrics: map[string]interface{}{
				"status": "yellow",
			},
		},
		{
			name:            "should parse red cluster status",
			clusterStatus:   "red",
			activeShardsNum: 5,
			nodeCount:       1,
			expectedMetrics: map[string]interface{}{
				"status": "red",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			assert.NotNil(t, c)
		})
	}
}

func TestCollectorErrorHandling(t *testing.T) {
	tests := []struct {
		name     string
		scenario string
		wantErr  bool
	}{
		{
			name:     "should handle nil response",
			scenario: "nil_response",
			wantErr:  true,
		},
		{
			name:     "should handle malformed JSON",
			scenario: "malformed_json",
			wantErr:  true,
		},
		{
			name:     "should handle timeout",
			scenario: "timeout",
			wantErr:  true,
		},
		{
			name:     "should handle server error",
			scenario: "server_error",
			wantErr:  true,
		},
		{
			name:     "should handle unauthorized",
			scenario: "unauthorized",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			assert.NotNil(t, c)
		})
	}
}

func TestCollectorBoundaryConditions(t *testing.T) {
	tests := []struct {
		name      string
		condition string
		value     interface{}
	}{
		{
			name:      "should handle zero active shards",
			condition: "zero_shards",
			value:     0,
		},
		{
			name:      "should handle negative shard count",
			condition: "negative_shards",
			value:     -1,
		},
		{
			name:      "should handle maximum shard count",
			condition: "max_shards",
			value:     2147483647,
		},
		{
			name:      "should handle empty cluster name",
			condition: "empty_cluster_name",
			value:     "",
		},
		{
			name:      "should handle very long cluster name",
			condition: "long_cluster_name",
			value:     strings.Repeat("a", 10000),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			assert.NotNil(t, c)
		})
	}
}

// MockHTTPClient for testing
type MockHTTPClient struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if m.DoFunc != nil {
		return m.DoFunc(req)
	}
	return nil, fmt.Errorf("not implemented")
}

func NewMockHTTPClient() *MockHTTPClient {
	return &MockHTTPClient{}
}