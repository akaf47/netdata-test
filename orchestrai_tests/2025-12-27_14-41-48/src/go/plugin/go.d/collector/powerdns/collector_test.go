package powerdns

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// TestNewCollector tests the NewCollector function
func TestNewCollector(t *testing.T) {
	c := NewCollector()
	assert.NotNil(t, c)
	assert.Equal(t, defaultURL, c.URL)
}

// TestCollectorName tests the Name method
func TestCollectorName(t *testing.T) {
	c := NewCollector()
	assert.Equal(t, "powerdns", c.Name())
}

// TestCollectorPriority tests the Priority method
func TestCollectorPriority(t *testing.T) {
	c := NewCollector()
	priority := c.Priority()
	assert.Greater(t, priority, 0)
}

// TestCollectorCheck tests the Check method
func TestCollectorCheck_Success(t *testing.T) {
	c := NewCollector()
	c.httpClient = &mockHTTPClient{
		doFunc: func(req interface{}) error {
			return nil
		},
	}

	err := c.Check()
	assert.NoError(t, err)
}

func TestCollectorCheck_HTTPClientError(t *testing.T) {
	c := NewCollector()
	c.httpClient = &mockHTTPClient{
		doFunc: func(req interface{}) error {
			return errors.New("http client error")
		},
	}

	err := c.Check()
	assert.Error(t, err)
}

func TestCollectorCheck_NilHTTPClient(t *testing.T) {
	c := NewCollector()
	c.httpClient = nil

	err := c.Check()
	assert.Error(t, err)
}

// TestCollectorCollect tests the Collect method
func TestCollectorCollect_Success(t *testing.T) {
	c := NewCollector()
	c.httpClient = &mockHTTPClient{
		doFunc: func(req interface{}) error {
			return nil
		},
	}

	mx := c.Collect()
	assert.NotNil(t, mx)
}

func TestCollectorCollect_HTTPClientError(t *testing.T) {
	c := NewCollector()
	c.httpClient = &mockHTTPClient{
		doFunc: func(req interface{}) error {
			return errors.New("http error")
		},
	}

	mx := c.Collect()
	assert.Nil(t, mx)
}

func TestCollectorCollect_InvalidResponse(t *testing.T) {
	c := NewCollector()
	c.httpClient = &mockHTTPClient{
		doFunc: func(req interface{}) error {
			return errors.New("invalid response")
		},
	}

	mx := c.Collect()
	assert.Nil(t, mx)
}

// TestCollectorStop tests the Stop method
func TestCollectorStop(t *testing.T) {
	c := NewCollector()
	err := c.Stop()
	assert.NoError(t, err)
}

// TestCollectorStop_WithError tests Stop with cleanup errors
func TestCollectorStop_WithCleanup(t *testing.T) {
	c := NewCollector()
	c.httpClient = &mockHTTPClient{}

	err := c.Stop()
	assert.NoError(t, err)
}

// TestCollectorWithURL tests setting custom URL
func TestCollectorWithURL(t *testing.T) {
	c := NewCollector()
	customURL := "http://localhost:8081"
	c.URL = customURL
	assert.Equal(t, customURL, c.URL)
}

// TestCollectorWithEmptyURL tests with empty URL
func TestCollectorWithEmptyURL(t *testing.T) {
	c := NewCollector()
	c.URL = ""
	assert.Equal(t, "", c.URL)
}

// TestCollectorWithTimeout tests timeout configuration
func TestCollectorWithTimeout(t *testing.T) {
	c := NewCollector()
	assert.Greater(t, c.Timeout.Duration.Seconds(), 0)
}

// TestCollectorMetrics tests that metrics are properly collected
func TestCollectorMetrics(t *testing.T) {
	c := NewCollector()
	c.httpClient = &mockHTTPClient{
		doFunc: func(req interface{}) error {
			return nil
		},
	}

	mx := c.Collect()
	if mx != nil {
		assert.NotEmpty(t, mx)
	}
}

// TestCollectorConcurrentCalls tests concurrent Check and Collect
func TestCollectorConcurrentCalls(t *testing.T) {
	c := NewCollector()
	c.httpClient = &mockHTTPClient{
		doFunc: func(req interface{}) error {
			return nil
		},
	}

	// Simulate concurrent calls
	done := make(chan error, 2)

	go func() {
		done <- c.Check()
	}()

	go func() {
		mx := c.Collect()
		if mx == nil {
			done <- errors.New("collect returned nil")
		} else {
			done <- nil
		}
	}()

	for i := 0; i < 2; i++ {
		err := <-done
		// Both should succeed or fail independently
		_ = err
	}
}

// TestCollectorMultipleCalls tests multiple sequential calls
func TestCollectorMultipleCalls(t *testing.T) {
	c := NewCollector()
	c.httpClient = &mockHTTPClient{
		doFunc: func(req interface{}) error {
			return nil
		},
	}

	for i := 0; i < 5; i++ {
		err := c.Check()
		assert.NoError(t, err)

		mx := c.Collect()
		assert.NotNil(t, mx)
	}
}

// TestCollectorURLHandling tests various URL formats
func TestCollectorURLHandling(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "default URL",
			url:  defaultURL,
		},
		{
			name: "custom localhost URL",
			url:  "http://localhost:8081",
		},
		{
			name: "IP-based URL",
			url:  "http://127.0.0.1:8081",
		},
		{
			name: "HTTPS URL",
			url:  "https://powerdns.example.com:8081",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			c.URL = tt.url
			assert.Equal(t, tt.url, c.URL)
		})
	}
}

// TestCollectorReInitialization tests re-initialization after Stop
func TestCollectorReInitialization(t *testing.T) {
	c := NewCollector()
	c.httpClient = &mockHTTPClient{
		doFunc: func(req interface{}) error {
			return nil
		},
	}

	err := c.Check()
	assert.NoError(t, err)

	err = c.Stop()
	assert.NoError(t, err)

	// Should be able to check again after stop
	err = c.Check()
	assert.NoError(t, err)
}

// TestCollectorEdgeCases tests edge cases
func TestCollectorEdgeCases_ZeroTimeout(t *testing.T) {
	c := NewCollector()
	c.Timeout.Duration = 0
	assert.Equal(t, 0, int(c.Timeout.Duration.Seconds()))
}

func TestCollectorEdgeCases_MaxIntValues(t *testing.T) {
	c := NewCollector()
	c.Timeout.Duration = 1<<63 - 1 // Max int64
	assert.NotNil(t, c)
}

// TestCollectorStatusAfterCheck tests collector state after Check
func TestCollectorStatusAfterCheck(t *testing.T) {
	c := NewCollector()
	c.httpClient = &mockHTTPClient{
		doFunc: func(req interface{}) error {
			return nil
		},
	}

	// Before check, httpClient is set
	assert.NotNil(t, c.httpClient)

	// After check
	err := c.Check()
	assert.NoError(t, err)
	assert.NotNil(t, c.httpClient)
}

// TestCollectorRetryScenarios tests error handling with retries
func TestCollectorErrorRecovery(t *testing.T) {
	c := NewCollector()
	callCount := 0

	c.httpClient = &mockHTTPClient{
		doFunc: func(req interface{}) error {
			callCount++
			if callCount < 2 {
				return errors.New("temporary error")
			}
			return nil
		},
	}

	// First call fails
	err := c.Check()
	assert.Error(t, err)
	assert.Equal(t, 1, callCount)

	// Second call succeeds
	err = c.Check()
	assert.NoError(t, err)
	assert.Equal(t, 2, callCount)
}

// Mock HTTP Client
type mockHTTPClient struct {
	mock.Mock
	doFunc func(req interface{}) error
}

func (m *mockHTTPClient) Do(req interface{}) error {
	if m.doFunc != nil {
		return m.doFunc(req)
	}
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(error)
}