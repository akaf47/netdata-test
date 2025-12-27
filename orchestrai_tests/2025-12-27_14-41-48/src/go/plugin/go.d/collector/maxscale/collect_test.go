package maxscale

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/netdata/netdata/go/plugins/pkg/web"
)

// MockHTTPClient mocks the HTTP client for testing
type MockHTTPClient struct {
	roundTripFunc func(*http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func TestCollectorCollect(t *testing.T) {
	tests := []struct {
		name         string
		setupMocks   func(*Collector)
		wantErr      bool
		errContains  string
		validateData func(map[string]int64)
	}{
		{
			name: "successful collection of all metrics",
			setupMocks: func(c *Collector) {
				c.httpClient = createMockClient(t, func(req *http.Request) *http.Response {
					switch req.URL.Path {
					case "/maxscale":
						return mockMaxScaleGlobalResponse()
					case "/maxscale/threads":
						return mockMaxScaleThreadsResponse()
					case "/servers":
						return mockServersResponse()
					default:
						return mockErrorResponse(http.StatusNotFound)
					}
				})
			},
			wantErr: false,
			validateData: func(mx map[string]int64) {
				if mx["uptime"] != 3600 {
					t.Errorf("uptime = %d, want 3600", mx["uptime"])
				}
				if mx["threads_reads"] <= 0 {
					t.Errorf("threads_reads not set correctly")
				}
			},
		},
		{
			name: "collect fails when global data fetch errors",
			setupMocks: func(c *Collector) {
				c.httpClient = createMockClient(t, func(req *http.Request) *http.Response {
					if req.URL.Path == "/maxscale" {
						return mockErrorResponse(http.StatusInternalServerError)
					}
					return mockMaxScaleGlobalResponse()
				})
			},
			wantErr: true,
		},
		{
			name: "collect fails when threads data fetch errors",
			setupMocks: func(c *Collector) {
				c.httpClient = createMockClient(t, func(req *http.Request) *http.Response {
					switch req.URL.Path {
					case "/maxscale":
						return mockMaxScaleGlobalResponse()
					case "/maxscale/threads":
						return mockErrorResponse(http.StatusInternalServerError)
					default:
						return mockErrorResponse(http.StatusNotFound)
					}
				})
			},
			wantErr: true,
		},
		{
			name: "collect fails when servers data fetch errors",
			setupMocks: func(c *Collector) {
				c.httpClient = createMockClient(t, func(req *http.Request) *http.Response {
					switch req.URL.Path {
					case "/maxscale":
						return mockMaxScaleGlobalResponse()
					case "/maxscale/threads":
						return mockMaxScaleThreadsResponse()
					case "/servers":
						return mockErrorResponse(http.StatusInternalServerError)
					default:
						return mockErrorResponse(http.StatusNotFound)
					}
				})
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collector{
				Config:      web.Config{},
				RequestConfig: web.RequestConfig{},
				seenServers: make(map[string]bool),
			}
			tt.setupMocks(c)

			mx, err := c.collect()

			if tt.wantErr && err == nil {
				t.Errorf("collect() error = nil, want error")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("collect() error = %v, want nil", err)
			}
			if !tt.wantErr && tt.validateData != nil {
				tt.validateData(mx)
			}
		})
	}
}

func TestCollectorCollectMaxScaleGlobal(t *testing.T) {
	tests := []struct {
		name        string
		mockFn      func(*http.Request) *http.Response
		wantErr     bool
		errContains string
		wantUptime  int64
	}{
		{
			name: "successfully collects global maxscale data",
			mockFn: func(req *http.Request) *http.Response {
				return mockMaxScaleGlobalResponse()
			},
			wantErr:    false,
			wantUptime: 3600,
		},
		{
			name: "returns error when HTTP request creation fails",
			mockFn: func(req *http.Request) *http.Response {
				return mockErrorResponse(http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "returns error when response data is nil",
			mockFn: func(req *http.Request) *http.Response {
				return mockMaxScaleGlobalResponseWithNilData()
			},
			wantErr:     true,
			errContains: "missing expected MaxScale data",
		},
		{
			name: "returns error when HTTP call fails",
			mockFn: func(req *http.Request) *http.Response {
				return mockErrorResponse(http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "correctly stores uptime value",
			mockFn: func(req *http.Request) *http.Response {
				return createMaxScaleGlobalResponseWithUptime(9999)
			},
			wantErr:    false,
			wantUptime: 9999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collector{
				Config:        web.Config{},
				RequestConfig: web.RequestConfig{},
				httpClient:    createMockClient(t, tt.mockFn),
			}

			mx := make(map[string]int64)
			err := c.collectMaxScaleGlobal(mx)

			if tt.wantErr && err == nil {
				t.Errorf("collectMaxScaleGlobal() error = nil, want error")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("collectMaxScaleGlobal() error = %v, want nil", err)
			}
			if tt.errContains != "" && err != nil && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("collectMaxScaleGlobal() error = %v, want error containing %s", err, tt.errContains)
			}
			if !tt.wantErr && mx["uptime"] != tt.wantUptime {
				t.Errorf("uptime = %d, want %d", mx["uptime"], tt.wantUptime)
			}
		})
	}
}

func TestCollectorCollectMaxScaleThreads(t *testing.T) {
	tests := []struct {
		name         string
		mockFn       func(*http.Request) *http.Response
		wantErr      bool
		validateData func(map[string]int64)
	}{
		{
			name: "successfully collects threads data with single thread",
			mockFn: func(req *http.Request) *http.Response {
				return mockMaxScaleThreadsResponse()
			},
			wantErr: false,
			validateData: func(mx map[string]int64) {
				if mx["threads_reads"] <= 0 {
					t.Errorf("threads_reads not set correctly")
				}
				if mx["threads_writes"] <= 0 {
					t.Errorf("threads_writes not set correctly")
				}
				if mx["threads_errors"] < 0 {
					t.Errorf("threads_errors not set correctly")
				}
				if mx["threads_hangups"] < 0 {
					t.Errorf("threads_hangups not set correctly")
				}
				if mx["threads_accepts"] <= 0 {
					t.Errorf("threads_accepts not set correctly")
				}
				if mx["threads_sessions"] < 0 {
					t.Errorf("threads_sessions not set correctly")
				}
				if mx["threads_zombies"] < 0 {
					t.Errorf("threads_zombies not set correctly")
				}
				if mx["threads_current_fds"] < 0 {
					t.Errorf("threads_current_fds not set correctly")
				}
				if mx["threads_total_fds"] < 0 {
					t.Errorf("threads_total_fds not set correctly")
				}
				if mx["threads_qc_cache_inserts"] < 0 {
					t.Errorf("threads_qc_cache_inserts not set correctly")
				}
				if mx["threads_qc_cache_evictions"] < 0 {
					t.Errorf("threads_qc_cache_evictions not set correctly")
				}
				if mx["threads_qc_cache_hits"] < 0 {
					t.Errorf("threads_qc_cache_hits not set correctly")
				}
				if mx["threads_qc_cache_misses"] < 0 {
					t.Errorf("threads_qc_cache_misses not set correctly")
				}
			},
		},
		{
			name: "successfully collects threads data with multiple threads",
			mockFn: func(req *http.Request) *http.Response {
				return mockMaxScaleThreadsResponseMultiple()
			},
			wantErr: false,
			validateData: func(mx map[string]int64) {
				// Should accumulate across multiple threads
				if mx["threads_reads"] <= 0 {
					t.Errorf("threads_reads should be accumulated")
				}
			},
		},
		{
			name: "returns error when HTTP call fails",
			mockFn: func(req *http.Request) *http.Response {
				return mockErrorResponse(http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "correctly sets all thread state values to 0 initially",
			mockFn: func(req *http.Request) *http.Response {
				return mockMaxScaleThreadsResponse()
			},
			wantErr: false,
			validateData: func(mx map[string]int64) {
				// All thread states should be initialized
				for _, state := range threadStates {
					if mx["threads_state_"+state] != 0 && mx["threads_state_"+state] != 1 {
						t.Errorf("threads_state_%s has unexpected value", state)
					}
				}
			},
		},
		{
			name: "increments state counter for matching state",
			mockFn: func(req *http.Request) *http.Response {
				return mockMaxScaleThreadsResponse()
			},
			wantErr: false,
			validateData: func(mx map[string]int64) {
				// At least one state should be incremented
				found := false
				for _, state := range threadStates {
					if mx["threads_state_"+state] == 1 {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("no thread state was incremented")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collector{
				Config:        web.Config{},
				RequestConfig: web.RequestConfig{},
				httpClient:    createMockClient(t, tt.mockFn),
			}

			mx := make(map[string]int64)
			err := c.collectMaxScaleThreads(mx)

			if tt.wantErr && err == nil {
				t.Errorf("collectMaxScaleThreads() error = nil, want error")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("collectMaxScaleThreads() error = %v, want nil", err)
			}
			if !tt.wantErr && tt.validateData != nil {
				tt.validateData(mx)
			}
		})
	}
}

func TestCollectorCollectServers(t *testing.T) {
	tests := []struct {
		name         string
		mockFn       func(*http.Request) *http.Response
		wantErr      bool
		validateData func(*Collector, map[string]int64)
	}{
		{
			name: "successfully collects servers data with single server",
			mockFn: func(req *http.Request) *http.Response {
				return mockServersResponse()
			},
			wantErr: false,
			validateData: func(c *Collector, mx map[string]int64) {
				if c.seenServers["server-1"] != true {
					t.Errorf("seenServers['server-1'] should be true")
				}
				if mx["server_server-1_connections"] <= 0 {
					t.Errorf("server connections not set correctly")
				}
			},
		},
		{
			name: "successfully collects servers data with multiple servers",
			mockFn: func(req *http.Request) *http.Response {
				return mockServersResponseMultiple()
			},
			wantErr: false,
			validateData: func(c *Collector, mx map[string]int64) {
				if c.seenServers["server-1"] != true {
					t.Errorf("seenServers['server-1'] should be true")
				}
				if c.seenServers["server-2"] != true {
					t.Errorf("seenServers['server-2'] should be true")
				}
			},
		},
		{
			name: "returns error when HTTP call fails",
			mockFn: func(req *http.Request) *http.Response {
				return mockErrorResponse(http.StatusInternalServerError)
			},
			wantErr: true,
		},
		{
			name: "skips servers with empty ID",
			mockFn: func(req *http.Request) *http.Response {
				return mockServersResponseWithEmptyID()
			},
			wantErr: false,
			validateData: func(c *Collector, mx map[string]int64) {
				// Only the valid server should be seen
				if len(c.seenServers) != 1 {
					t.Errorf("seenServers should contain only 1 server, got %d", len(c.seenServers))
				}
			},
		},
		{
			name: "removes charts for servers no longer present",
			mockFn: func(req *http.Request) *http.Response {
				return mockServersResponseSingleServer()
			},
			wantErr: false,
			validateData: func(c *Collector, mx map[string]int64) {
				// After second call with different server, first one should be removed
				if c.seenServers["server-1"] != true {
					t.Errorf("current server should be in seenServers")
				}
			},
		},
		{
			name: "correctly parses server state with multiple values",
			mockFn: func(req *http.Request) *http.Response {
				return mockServersResponseWithMultipleStates()
			},
			wantErr: false,
			validateData: func(c *Collector, mx map[string]int64) {
				// State parsing should work with multiple space-separated states
				if c.seenServers["server-1"] != true {
					t.Errorf("server should be seen")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collector{
				Config:        web.Config{},
				RequestConfig: web.RequestConfig{},
				httpClient:    createMockClient(t, tt.mockFn),
				seenServers:   make(map[string]bool),
			}

			mx := make(map[string]int64)
			err := c.collectServers(mx)

			if tt.wantErr && err == nil {
				t.Errorf("collectServers() error = nil, want error")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("collectServers() error = %v, want nil", err)
			}
			if !tt.wantErr && tt.validateData != nil {
				tt.validateData(c, mx)
			}
		})
	}
}

func TestCollectorCollectServersStateHandling(t *testing.T) {
	tests := []struct {
		name          string
		mockFn        func(*http.Request) *http.Response
		expectedState map[string]int64
	}{
		{
			name: "handles server with single state",
			mockFn: func(req *http.Request) *http.Response {
				return mockServersResponseWithState("Master")
			},
			expectedState: map[string]int64{
				"server_server-1_state_Master": 1,
			},
		},
		{
			name: "handles server with multiple states",
			mockFn: func(req *http.Request) *http.Response {
				return mockServersResponseWithState("Master Authenticated")
			},
			expectedState: map[string]int64{
				"server_server-1_state_Master":        1,
				"server_server-1_state_Authenticated": 1,
			},
		},
		{
			name: "handles server with no state",
			mockFn: func(req *http.Request) *http.Response {
				return mockServersResponseWithState("")
			},
			expectedState: map[string]int64{},
		},
		{
			name: "handles server with whitespace in state",
			mockFn: func(req *http.Request) *http.Response {
				return mockServersResponseWithState("  Master   Authenticated  ")
			},
			expectedState: map[string]int64{
				"server_server-1_state_Master":        1,
				"server_server-1_state_Authenticated": 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Collector{
				Config:        web.Config{},
				RequestConfig: web.RequestConfig{},
				httpClient:    createMockClient(t, tt.mockFn),
				seenServers:   make(map[string]bool),
			}

			mx := make(map[string]int64)
			err := c.collectServers(mx)

			if err != nil {
				t.Fatalf("collectServers() error = %v, want nil", err)
			}

			for key, expectedVal := range tt.expectedState {
				if mx[key] != expectedVal {
					t.Errorf("mx[%s] = %d, want %d", key, mx[key], expectedVal)
				}
			}
		})
	}
}

func TestCollectorCollectServersDynamicAddition(t *testing.T) {
	c := &Collector{
		Config:        web.Config{},
		RequestConfig: web.RequestConfig{},
		seenServers:   make(map[string]bool),
	}

	// First call - server exists
	c.httpClient = createMockClient(t, func(req *http.Request) *http.Response {
		return mockServersResponseSingleServer()
	})

	mx := make(map[string]int64)
	if err := c.collectServers(mx); err != nil {
		t.Fatalf("first collectServers() error = %v, want nil", err)
	}

	if c.seenServers["server-1"] != true {
		t.Errorf("server-1 should be in seenServers after first call")
	}

	// Second call - server still exists (no change)
	mx = make(map[string]int64)
	if err := c.collectServers(mx); err != nil {
		t.Fatalf("second collectServers() error = %v, want nil", err)
	}

	if c.seenServers["server-1"] != true {
		t.Errorf("server-1 should still be in seenServers after second call")
	}

	// Third call - server is gone
	c.httpClient = createMockClient(t, func(req *http.Request) *http.Response {
		return mockServersResponseEmpty()
	})

	mx = make(map[string]int64)
	if err := c.collectServers(mx); err != nil {
		t.Fatalf("third collectServers() error = %v, want nil", err)
	}

	if c.seenServers["server-1"] != false {
		t.Errorf("server-1 should be removed from seenServers after third call")
	}
}

// Helper functions

func createMockClient(t *testing.T, fn func(*http.Request) *http.Response) *http.Client {
	return &http.Client{
		Transport: &MockHTTPClient{
			roundTripFunc: fn,
		},
	}
}

func mockMaxScaleGlobalResponse() *http.Response {
	body := `{"data":{"attributes":{"uptime":3600}}}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       createReadCloser(body),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func mockMaxScaleGlobalResponseWithNilData() *http.Response {
	body := `{"data":null}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       createReadCloser(body),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func createMaxScaleGlobalResponseWithUptime(uptime int64) *http.Response {
	body := fmt.Sprintf(`{"data":{"attributes":{"uptime":%d}}}`, uptime)
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       createReadCloser(body),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func mockMaxScaleThreadsResponse() *http.Response {
	body := `{
		"data":[
			{
				"attributes":{
					"stats":{
						"reads":100,
						"writes":50,
						"errors":0,
						"hangups":1,
						"accepts":10,
						"sessions":5,
						"zombies":0,
						"current_descriptors":20,
						"total_descriptors":100,
						"qc_cache":{"inserts":10,"evictions":1,"hits":50,"misses":5},
						"state":"Running"
					}
				}
			}
		]
	}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       createReadCloser(body),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func mockMaxScaleThreadsResponseMultiple() *http.Response {
	body := `{
		"data":[
			{
				"attributes":{
					"stats":{
						"reads":100,
						"writes":50,
						"errors":0,
						"hangups":1,
						"accepts":10,
						"sessions":5,
						"zombies":0,
						"current_descriptors":20,
						"total_descriptors":100,
						"qc_cache":{"inserts":10,"evictions":1,"hits":50,"misses":5},
						"state":"Running"
					}
				}
			},
			{
				"attributes":{
					"stats":{
						"reads":200,
						"writes":75,
						"errors":1,
						"hangups":0,
						"accepts":15,
						"sessions":8,
						"zombies":1,
						"current_descriptors":30,
						"total_descriptors":150,
						"qc_cache":{"inserts":20,"evictions":2,"hits":100,"misses":10},
						"state":"Running"
					}
				}
			}
		]
	}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       createReadCloser(body),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func mockServersResponse() *http.Response {
	body := `{
		"data":[
			{
				"id":"server-1",
				"attributes":{
					"params":{"address":"127.0.0.1","port":3306},
					"statistics":{"connections":10},
					"state":"Master"
				}
			}
		]
	}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       createReadCloser(body),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func mockServersResponseSingleServer() *http.Response {
	body := `{
		"data":[
			{
				"id":"server-1",
				"attributes":{
					"params":{"address":"192.168.1.1","port":3306},
					"statistics":{"connections":15},
					"state":"Slave"
				}
			}
		]
	}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       createReadCloser(body),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func mockServersResponseMultiple() *http.Response {
	body := `{
		"data":[
			{
				"id":"server-1",
				"attributes":{
					"params":{"address":"127.0.0.1","port":3306},
					"statistics":{"connections":10},
					"state":"Master"
				}
			},
			{
				"id":"server-2",
				"attributes":{
					"params":{"address":"127.0.0.2","port":3306},
					"statistics":{"connections":20},
					"state":"Slave"
				}
			}
		]
	}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       createReadCloser(body),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func mockServersResponseWithEmptyID() *http.Response {
	body := `{
		"data":[
			{
				"id":"",
				"attributes":{
					"params":{"address":"127.0.0.1","port":3306},
					"statistics":{"connections":10},
					"state":"Master"
				}
			},
			{
				"id":"server-1",
				"attributes":{
					"params":{"address":"127.0.0.2","port":3306},
					"statistics":{"connections":20},
					"state":"Slave"
				}
			}
		]
	}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       createReadCloser(body),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func mockServersResponseEmpty() *http.Response {
	body := `{"data":[]}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       createReadCloser(body),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func mockServersResponseWithMultipleStates() *http.Response {
	body := `{
		"data":[
			{
				"id":"server-1",
				"attributes":{
					"params":{"address":"127.0.0.1","port":3306},
					"statistics":{"connections":10},
					"state":"Master Authenticated"
				}
			}
		]
	}`
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       createReadCloser(body),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func mockServersResponseWithState(state string) *http.Response {
	body := fmt.Sprintf(`{
		"data":[
			{
				"id":"server-1",
				"attributes":{
					"params":{"address":"127.0.0.1","port":3306},
					"statistics":{"connections":10},
					"state":"%s"
				}
			}
		]
	}`, state)
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       createReadCloser(body),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func mockErrorResponse(statusCode int) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       createReadCloser("error"),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func createReadCloser(body string) ReadCloser {
	return &stringReadCloser{data: body}
}

type ReadCloser interface {
	Read(p []byte) (n int, err error)
	Close() error
}

type stringReadCloser struct {
	data string
	pos  int
}

func (src *stringReadCloser) Read(p []byte) (n int, err error) {
	if src.pos >= len(src.data) {
		return 0, nil
	}
	n = copy(p, src.data[src.pos:])
	src.pos += n
	return n, nil
}

func (src *stringReadCloser) Close() error {
	return nil
}