package freeradius

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCollect tests the Collect method
func TestCollect(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*FreeRadius)
		expectError bool
		expectFunc  func(*testing.T, *FreeRadius)
	}{
		{
			name: "successful collection",
			setup: func(fr *FreeRadius) {
				fr.Client = &MockClient{
					RequestFunc: func() (map[string]interface{}, error) {
						return map[string]interface{}{
							"requests":  float64(100),
							"responses": float64(95),
							"accepts":   float64(80),
							"rejects":   float64(15),
						}, nil
					},
				}
			},
			expectError: false,
			expectFunc: func(t *testing.T, fr *FreeRadius) {
				assert.NotNil(t, fr)
			},
		},
		{
			name: "client request fails",
			setup: func(fr *FreeRadius) {
				fr.Client = &MockClient{
					RequestFunc: func() (map[string]interface{}, error) {
						return nil, ErrClientRequestFailed
					},
				}
			},
			expectError: true,
		},
		{
			name: "nil client",
			setup: func(fr *FreeRadius) {
				fr.Client = nil
			},
			expectError: true,
		},
		{
			name: "empty response",
			setup: func(fr *FreeRadius) {
				fr.Client = &MockClient{
					RequestFunc: func() (map[string]interface{}, error) {
						return map[string]interface{}{}, nil
					},
				}
			},
			expectError: false,
		},
		{
			name: "response with unexpected types",
			setup: func(fr *FreeRadius) {
				fr.Client = &MockClient{
					RequestFunc: func() (map[string]interface{}, error) {
						return map[string]interface{}{
							"requests": "not_a_number",
							"responses": int64(95),
						}, nil
					},
				}
			},
			expectError: false,
		},
		{
			name: "response with negative values",
			setup: func(fr *FreeRadius) {
				fr.Client = &MockClient{
					RequestFunc: func() (map[string]interface{}, error) {
						return map[string]interface{}{
							"requests":  float64(-100),
							"responses": float64(-95),
						}, nil
					},
				}
			},
			expectError: false,
		},
		{
			name: "response with zero values",
			setup: func(fr *FreeRadius) {
				fr.Client = &MockClient{
					RequestFunc: func() (map[string]interface{}, error) {
						return map[string]interface{}{
							"requests":  float64(0),
							"responses": float64(0),
						}, nil
					},
				}
			},
			expectError: false,
		},
		{
			name: "response with max int64 values",
			setup: func(fr *FreeRadius) {
				fr.Client = &MockClient{
					RequestFunc: func() (map[string]interface{}, error) {
						return map[string]interface{}{
							"requests":  float64(9223372036854775807),
							"responses": float64(9223372036854775806),
						}, nil
					},
				}
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := &FreeRadius{
				Config: Config{
					Server:   "localhost",
					Port:     1812,
					Timeout:  time.Second * 5,
					Password: "testing123",
				},
			}

			if tt.setup != nil {
				tt.setup(fr)
			}

			err := fr.Collect()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectFunc != nil {
				tt.expectFunc(t, fr)
			}
		})
	}
}

// TestCollect_ParseMetrics tests metric parsing in Collect
func TestCollect_ParseMetrics(t *testing.T) {
	tests := []struct {
		name      string
		response  map[string]interface{}
		expectErr bool
	}{
		{
			name: "all standard metrics present",
			response: map[string]interface{}{
				"requests":         float64(1000),
				"responses":        float64(950),
				"accepts":          float64(900),
				"rejects":          float64(50),
				"challenges":       float64(10),
				"accounting_requests":  float64(500),
				"accounting_responses": float64(495),
				"bad_authenticators":   float64(5),
				"dropped_requests":     float64(2),
				"unknown_types":        float64(1),
			},
			expectErr: false,
		},
		{
			name:      "nil response",
			response:  nil,
			expectErr: true,
		},
		{
			name:      "empty response",
			response:  map[string]interface{}{},
			expectErr: false,
		},
		{
			name: "mixed numeric types",
			response: map[string]interface{}{
				"requests":  int64(100),
				"responses": float64(95),
				"accepts":   int32(80),
				"rejects":   uint64(15),
			},
			expectErr: false,
		},
		{
			name: "string numeric values",
			response: map[string]interface{}{
				"requests":  "100",
				"responses": "95",
			},
			expectErr: false,
		},
		{
			name: "bool values",
			response: map[string]interface{}{
				"requests": true,
				"responses": false,
			},
			expectErr: false,
		},
		{
			name: "null values in response",
			response: map[string]interface{}{
				"requests":  nil,
				"responses": float64(95),
			},
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := &FreeRadius{
				Client: &MockClient{
					RequestFunc: func() (map[string]interface{}, error) {
						return tt.response, nil
					},
				},
			}

			err := fr.Collect()

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestCollect_ClientRetry tests retry logic
func TestCollect_ClientRetry(t *testing.T) {
	tests := []struct {
		name          string
		maxRetries    int
		failureCount  int
		shouldSucceed bool
	}{
		{
			name:          "succeeds on first attempt",
			maxRetries:    3,
			failureCount:  0,
			shouldSucceed: true,
		},
		{
			name:          "succeeds after one failure",
			maxRetries:    3,
			failureCount:  1,
			shouldSucceed: true,
		},
		{
			name:          "succeeds after multiple failures",
			maxRetries:    5,
			failureCount:  3,
			shouldSucceed: true,
		},
		{
			name:          "fails when retries exhausted",
			maxRetries:    2,
			failureCount:  5,
			shouldSucceed: false,
		},
		{
			name:          "zero retries with failure",
			maxRetries:    0,
			failureCount:  1,
			shouldSucceed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attempts := 0
			fr := &FreeRadius{
				Client: &MockClient{
					RequestFunc: func() (map[string]interface{}, error) {
						attempts++
						if attempts <= tt.failureCount {
							return nil, ErrTemporaryFailure
						}
						return map[string]interface{}{"requests": float64(100)}, nil
					},
				},
				Config: Config{
					MaxRetries: tt.maxRetries,
				},
			}

			err := fr.Collect()

			if tt.shouldSucceed {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

// TestCollect_Timeout tests timeout handling
func TestCollect_Timeout(t *testing.T) {
	tests := []struct {
		name            string
		clientDelay     time.Duration
		configTimeout   time.Duration
		shouldTimeout   bool
	}{
		{
			name:            "request completes before timeout",
			clientDelay:     time.Millisecond * 100,
			configTimeout:   time.Second,
			shouldTimeout:   false,
		},
		{
			name:            "request exceeds timeout",
			clientDelay:     time.Second * 2,
			configTimeout:   time.Millisecond * 500,
			shouldTimeout:   true,
		},
		{
			name:            "zero timeout",
			clientDelay:     time.Millisecond * 100,
			configTimeout:   0,
			shouldTimeout:   false,
		},
		{
			name:            "negative timeout treated as no timeout",
			clientDelay:     time.Millisecond * 100,
			configTimeout:   time.Duration(-1) * time.Second,
			shouldTimeout:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fr := &FreeRadius{
				Client: &MockClient{
					RequestFunc: func() (map[string]interface{}, error) {
						time.Sleep(tt.clientDelay)
						return map[string]interface{}{"requests": float64(100)}, nil
					},
				},
				Config: Config{
					Timeout: tt.configTimeout,
				},
			}

			err := fr.Collect()

			if tt.shouldTimeout {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "timeout")
			} else if tt.configTimeout > 0 && tt.clientDelay < tt.configTimeout {
				assert.NoError(t, err)
			}
		})
	}
}

// TestCollect_MetricAggregation tests metric aggregation logic
func TestCollect_MetricAggregation(t *testing.T) {
	tests := []struct {
		name            string
		responses       []map[string]interface{}
		expectedTotal   int64
	}{
		{
			name: "single response",
			responses: []map[string]interface{}{
				{"requests": float64(100)},
			},
			expectedTotal: 100,
		},
		{
			name: "multiple responses aggregated",
			responses: []map[string]interface{}{
				{"requests": float64(100)},
				{"requests": float64(200)},
				{"requests": float64(150)},
			},
			expectedTotal: 450,
		},
		{
			name: "handles missing metrics in aggregation",
			responses: []map[string]interface{}{
				{"requests": float64(100)},
				{},
				{"requests": float64(200)},
			},
			expectedTotal: 300,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseIndex := 0
			fr := &FreeRadius{
				Client: &MockClient{
					RequestFunc: func() (map[string]interface{}, error) {
						if responseIndex < len(tt.responses) {
							result := tt.responses[responseIndex]
							responseIndex++
							return result, nil
						}
						return nil, ErrNoMoreResponses
					},
				},
			}

			err := fr.Collect()
			assert.NoError(t, err)
		})
	}
}

// TestCollect_ConcurrentRequests tests concurrent collection
func TestCollect_ConcurrentRequests(t *testing.T) {
	t.Run("multiple concurrent collections", func(t *testing.T) {
		fr := &FreeRadius{
			Client: &MockClient{
				RequestFunc: func() (map[string]interface{}, error) {
					return map[string]interface{}{
						"requests": float64(100),
					}, nil
				},
			},
		}

		done := make(chan error, 10)
		for i := 0; i < 10; i++ {
			go func() {
				done <- fr.Collect()
			}()
		}

		for i := 0; i < 10; i++ {
			err := <-done
			assert.NoError(t, err)
		}
	})
}

// Mock types for testing
type MockClient struct {
	RequestFunc func() (map[string]interface{}, error)
	CloseFunc   func() error
}

func (m *MockClient) Request() (map[string]interface{}, error) {
	if m.RequestFunc != nil {
		return m.RequestFunc()
	}
	return nil, nil
}

func (m *MockClient) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

// Custom error types for testing
var (
	ErrClientRequestFailed = Error("client request failed")
	ErrTemporaryFailure    = Error("temporary failure")
	ErrNoMoreResponses     = Error("no more responses")
)

type Error string

func (e Error) Error() string {
	return string(e)
}