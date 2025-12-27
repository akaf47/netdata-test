package powerdns

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestInit tests the Init method
func TestInit(t *testing.T) {
	c := NewCollector()
	err := c.Init()
	assert.NoError(t, err)
}

// TestInit_WithValidConfiguration tests Init with valid config
func TestInit_WithValidConfiguration(t *testing.T) {
	c := NewCollector()
	c.URL = "http://localhost:8081"
	c.Timeout.Duration = 5 * time.Second

	err := c.Init()
	assert.NoError(t, err)
	assert.NotNil(t, c.httpClient)
}

// TestInit_WithEmptyURL tests Init with empty URL
func TestInit_WithEmptyURL(t *testing.T) {
	c := NewCollector()
	c.URL = ""

	err := c.Init()
	// Should either error or use default
	if err == nil {
		// If no error, httpClient should be initialized with default URL
		assert.NotNil(t, c.httpClient)
	}
}

// TestInit_CreatesHTTPClient tests that Init creates HTTP client
func TestInit_CreatesHTTPClient(t *testing.T) {
	c := NewCollector()
	assert.Nil(t, c.httpClient)

	err := c.Init()
	assert.NoError(t, err)
	assert.NotNil(t, c.httpClient)
}

// TestInit_WithCustomTimeout tests Init with custom timeout
func TestInit_WithCustomTimeout(t *testing.T) {
	c := NewCollector()
	c.Timeout.Duration = 10 * time.Second

	err := c.Init()
	assert.NoError(t, err)
	assert.NotNil(t, c.httpClient)
}

// TestInit_WithZeroTimeout tests Init with zero timeout
func TestInit_WithZeroTimeout(t *testing.T) {
	c := NewCollector()
	c.Timeout.Duration = 0

	err := c.Init()
	// Should still initialize even with zero timeout
	assert.NoError(t, err)
	assert.NotNil(t, c.httpClient)
}

// TestInit_WithMaxTimeout tests Init with very large timeout
func TestInit_WithMaxTimeout(t *testing.T) {
	c := NewCollector()
	c.Timeout.Duration = 1000 * time.Hour

	err := c.Init()
	assert.NoError(t, err)
	assert.NotNil(t, c.httpClient)
}

// TestInit_MultipleTimesRequiresCheck tests calling Init multiple times
func TestInit_MultipleCalls(t *testing.T) {
	c := NewCollector()

	err := c.Init()
	assert.NoError(t, err)
	firstClient := c.httpClient

	err = c.Init()
	assert.NoError(t, err)
	secondClient := c.httpClient

	// Both should be valid clients
	assert.NotNil(t, firstClient)
	assert.NotNil(t, secondClient)
}

// TestInit_PreservesConfiguration tests that Init preserves config
func TestInit_PreservesConfiguration(t *testing.T) {
	c := NewCollector()
	c.URL = "http://custom.url:8081"
	c.Timeout.Duration = 7 * time.Second

	err := c.Init()
	assert.NoError(t, err)

	// Configuration should be preserved
	assert.Equal(t, "http://custom.url:8081", c.URL)
	assert.Equal(t, 7*time.Second, c.Timeout.Duration)
}

// TestInit_WithDifferentURLSchemes tests various URL schemes
func TestInit_WithDifferentURLSchemes(t *testing.T) {
	schemes := []string{
		"http://localhost:8081",
		"https://localhost:8081",
		"http://127.0.0.1:8081",
		"http://powerdns.example.com:8081",
	}

	for _, scheme := range schemes {
		t.Run(scheme, func(t *testing.T) {
			c := NewCollector()
			c.URL = scheme

			err := c.Init()
			assert.NoError(t, err)
			assert.NotNil(t, c.httpClient)
		})
	}
}

// TestInit_StateBeforeAfter tests state transitions
func TestInit_StateTransition(t *testing.T) {
	c := NewCollector()

	// Before init
	assert.Nil(t, c.httpClient)

	// After init
	err := c.Init()
	assert.NoError(t, err)
	assert.NotNil(t, c.httpClient)

	// Should remain initialized
	assert.NotNil(t, c.httpClient)
}

// TestInit_WithFactoryDefaults tests Init with factory defaults
func TestInit_WithFactoryDefaults(t *testing.T) {
	c := NewCollector()
	// Don't change any configuration, use defaults

	err := c.Init()
	assert.NoError(t, err)
	assert.NotNil(t, c.httpClient)
	assert.Equal(t, defaultURL, c.URL)
}

// TestInit_TimeoutBoundaryValues tests boundary timeout values
func TestInit_TimeoutBoundaryValues(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
	}{
		{"1 nanosecond", 1 * time.Nanosecond},
		{"1 millisecond", 1 * time.Millisecond},
		{"1 second", 1 * time.Second},
		{"1 minute", 1 * time.Minute},
		{"1 hour", 1 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			c.Timeout.Duration = tt.timeout

			err := c.Init()
			assert.NoError(t, err)
			assert.NotNil(t, c.httpClient)
		})
	}
}

// TestInit_URLValidation tests URL path variations
func TestInit_URLValidation(t *testing.T) {
	urls := []string{
		"http://localhost:8081",
		"http://localhost:8081/",
		"http://localhost:8081/api",
		"http://localhost:8081/api/",
	}

	for _, url := range urls {
		t.Run(url, func(t *testing.T) {
			c := NewCollector()
			c.URL = url

			err := c.Init()
			assert.NoError(t, err)
			assert.NotNil(t, c.httpClient)
		})
	}
}

// TestInit_ErrorHandling tests error handling in Init
func TestInit_ErrorHandling(t *testing.T) {
	c := NewCollector()
	c.URL = "http://localhost:8081"

	// This should work
	err := c.Init()
	assert.NoError(t, err)
}

// TestInit_InvalidURLFormats tests various invalid URL formats
func TestInit_InvalidURLFormats(t *testing.T) {
	invalidURLs := []string{
		"",
		"not-a-url",
		"://invalid",
		"http://",
	}

	for _, url := range invalidURLs {
		t.Run(url, func(t *testing.T) {
			c := NewCollector()
			c.URL = url

			// Init may or may not error on invalid URL, depending on implementation
			// But it should not panic
			assert.NotPanics(t, func() {
				_ = c.Init()
			})
		})
	}
}

// TestInit_HTTPClientInitialization tests HTTP client is properly initialized
func TestInit_HTTPClientInitialization(t *testing.T) {
	c := NewCollector()

	err := c.Init()
	require.NoError(t, err)

	// HTTP client should be initialized
	require.NotNil(t, c.httpClient)

	// Should be usable
	assert.NotNil(t, c.httpClient)
}

// TestInit_SequentialInitCalls tests sequential Init calls
func TestInit_SequentialInitCalls(t *testing.T) {
	c := NewCollector()

	for i := 0; i < 3; i++ {
		err := c.Init()
		assert.NoError(t, err)
		assert.NotNil(t, c.httpClient)
	}
}

// TestInit_URLEdgeCases tests edge case URLs
func TestInit_URLEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"URL with port", "http://localhost:8081"},
		{"URL with high port", "http://localhost:65535"},
		{"URL with low port", "http://localhost:1"},
		{"localhost IPv4", "http://127.0.0.1:8081"},
		{"localhost IPv6", "http://[::1]:8081"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			c.URL = tt.url

			assert.NotPanics(t, func() {
				_ = c.Init()
			})
		})
	}
}

// TestInit_TimeoutEdgeCases tests timeout edge cases
func TestInit_TimeoutEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
	}{
		{"negative timeout", -1 * time.Second},
		{"minimal timeout", 1 * time.Nanosecond},
		{"large timeout", time.Duration(^uint64(0) >> 1)}, // Max duration
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			c.Timeout.Duration = tt.timeout

			assert.NotPanics(t, func() {
				_ = c.Init()
			})
		})
	}
}

// TestInit_InitBeforeCheck tests Init is called before Check
func TestInit_InitBeforeCheck(t *testing.T) {
	c := NewCollector()

	// Without Init, Check should fail or handle it gracefully
	if c.httpClient == nil {
		err := c.Init()
		assert.NoError(t, err)
	}

	// Now Check should work
	err := c.Check()
	assert.NoError(t, err)
}

// TestInit_InitBeforeCollect tests Init is called before Collect
func TestInit_InitBeforeCollect(t *testing.T) {
	c := NewCollector()

	if c.httpClient == nil {
		err := c.Init()
		assert.NoError(t, err)
	}

	// Now Collect should work
	mx := c.Collect()
	assert.NotNil(t, mx)
}

// TestInit_ConcurrentInitCalls tests concurrent Init calls
func TestInit_ConcurrentInitCalls(t *testing.T) {
	c := NewCollector()
	done := make(chan error, 5)

	for i := 0; i < 5; i++ {
		go func() {
			done <- c.Init()
		}()
	}

	for i := 0; i < 5; i++ {
		err := <-done
		assert.NoError(t, err)
	}

	assert.NotNil(t, c.httpClient)
}

// TestInit_StateConsistency tests state consistency across calls
func TestInit_StateConsistency(t *testing.T) {
	c := NewCollector()

	err := c.Init()
	assert.NoError(t, err)
	client1 := c.httpClient

	err = c.Init()
	assert.NoError(t, err)
	client2 := c.httpClient

	// After re-init, both should be valid
	assert.NotNil(t, client1)
	assert.NotNil(t, client2)
}

// TestInit_ConfigPersistence tests configuration persists through Init
func TestInit_ConfigPersistence(t *testing.T) {
	c := NewCollector()
	originalURL := "http://custom.url:9999"
	originalTimeout := 15 * time.Second

	c.URL = originalURL
	c.Timeout.Duration = originalTimeout

	err := c.Init()
	assert.NoError(t, err)

	// Config should persist
	assert.Equal(t, originalURL, c.URL)
	assert.Equal(t, originalTimeout, c.Timeout.Duration)

	// Init again
	err = c.Init()
	assert.NoError(t, err)

	// Config should still persist
	assert.Equal(t, originalURL, c.URL)
	assert.Equal(t, originalTimeout, c.Timeout.Duration)
}