package varnish

import (
	"fmt"
	"testing"
	"time"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNew tests the New function
func TestNew(t *testing.T) {
	job := New()

	require.NotNil(t, job)
	assert.NotNil(t, job.Config)
	assert.Equal(t, defaultConfig, job.Config)
}

// TestVarnish_Init tests the Init method
func TestVarnish_Init_Success(t *testing.T) {
	v := New()

	// Mock the client creation
	v.client = &mockVarnishClient{
		stats: map[string]interface{}{},
	}

	err := v.Init()

	assert.NoError(t, err)
}

// TestVarnish_Init_ClientCreationError tests Init when client creation fails
func TestVarnish_Init_ClientCreationError(t *testing.T) {
	v := New()

	// Simulate client creation failure by setting invalid config
	v.Command = "/nonexistent/varnishstat"
	v.Timeout = 0

	err := v.Init()

	// Should handle error gracefully
	assert.Error(t, err)
}

// TestVarnish_Init_Timeout tests Init with timeout configuration
func TestVarnish_Init_Timeout(t *testing.T) {
	v := New()

	v.Timeout = 5 * time.Second
	v.client = &mockVarnishClient{
		stats: map[string]interface{}{},
	}

	err := v.Init()

	assert.NoError(t, err)
	assert.Equal(t, 5*time.Second, v.Timeout)
}

// TestVarnish_Check tests the Check method with valid client
func TestVarnish_Check_Success(t *testing.T) {
	v := New()
	v.client = &mockVarnishClient{
		stats: map[string]interface{}{
			"MAIN.uptime": float64(1000),
		},
	}

	result := v.Check()

	assert.True(t, result)
}

// TestVarnish_Check_NoData tests Check when client returns no data
func TestVarnish_Check_NoData(t *testing.T) {
	v := New()
	v.client = &mockVarnishClient{
		stats: map[string]interface{}{},
	}

	result := v.Check()

	assert.False(t, result)
}

// TestVarnish_Check_ClientError tests Check when client returns error
func TestVarnish_Check_ClientError(t *testing.T) {
	v := New()
	v.client = &mockVarnishClient{
		err: fmt.Errorf("connection error"),
	}

	result := v.Check()

	assert.False(t, result)
}

// TestVarnish_Check_NilClient tests Check with nil client
func TestVarnish_Check_NilClient(t *testing.T) {
	v := New()
	v.client = nil

	result := v.Check()

	assert.False(t, result)
}

// TestVarnish_Charts tests the Charts method
func TestVarnish_Charts(t *testing.T) {
	v := New()

	charts := v.Charts()

	assert.NotNil(t, charts)
	assert.Greater(t, len(charts), 0)
}

// TestVarnish_Cleanup tests the Cleanup method
func TestVarnish_Cleanup(t *testing.T) {
	v := New()
	v.client = &mockVarnishClient{
		stats: map[string]interface{}{},
	}

	// Should not panic
	v.Cleanup()

	// After cleanup, client should be nil
	assert.Nil(t, v.client)
}

// TestVarnish_Cleanup_NilClient tests Cleanup with nil client
func TestVarnish_Cleanup_NilClient(t *testing.T) {
	v := New()
	v.client = nil

	// Should not panic
	v.Cleanup()

	assert.Nil(t, v.client)
}

// TestVarnish_String tests the String method
func TestVarnish_String(t *testing.T) {
	v := New()
	v.Name = "varnish-1"

	str := v.String()

	assert.NotEmpty(t, str)
	assert.Contains(t, str, "varnish-1")
}

// TestVarnish_String_EmptyName tests String with empty name
func TestVarnish_String_EmptyName(t *testing.T) {
	v := New()
	v.Name = ""

	str := v.String()

	assert.NotEmpty(t, str)
}

// TestVarnish_ConfigValidation tests configuration validation
func TestVarnish_ConfigValidation(t *testing.T) {
	v := New()

	// Valid configuration
	assert.NotNil(t, v.Config)
}

// TestVarnish_ConfigTimeout_Default tests default timeout
func TestVarnish_ConfigTimeout_Default(t *testing.T) {
	v := New()

	assert.NotZero(t, v.Timeout)
}

// TestVarnish_ConfigTimeout_Custom tests custom timeout
func TestVarnish_ConfigTimeout_Custom(t *testing.T) {
	v := New()
	customTimeout := 10 * time.Second

	v.Timeout = customTimeout

	assert.Equal(t, customTimeout, v.Timeout)
}

// TestVarnish_ConfigCommand_Default tests default command
func TestVarnish_ConfigCommand_Default(t *testing.T) {
	v := New()

	// Should have some default command set
	assert.NotEmpty(t, v.Command)
}

// TestVarnish_ConfigCommand_Custom tests custom command
func TestVarnish_ConfigCommand_Custom(t *testing.T) {
	v := New()
	customCmd := "/usr/bin/varnishstat"

	v.Command = customCmd

	assert.Equal(t, customCmd, v.Command)
}

// TestVarnish_MultipleCollectCycles tests multiple collect cycles
func TestVarnish_MultipleCollectCycles(t *testing.T) {
	v := New()
	v.client = &mockVarnishClient{
		stats: map[string]interface{}{
			"MAIN.uptime":    float64(1000),
			"MAIN.sess_conn": float64(100),
			"MAIN.cache_hit": float64(400),
		},
	}

	// First collect
	mx1 := v.Collect()
	assert.NotNil(t, mx1)

	// Second collect
	mx2 := v.Collect()
	assert.NotNil(t, mx2)
}

// TestVarnish_CollectAndCheck tests Check after Collect
func TestVarnish_CollectAndCheck(t *testing.T) {
	v := New()
	v.client = &mockVarnishClient{
		stats: map[string]interface{}{
			"MAIN.uptime": float64(1000),
		},
	}

	mx := v.Collect()
	assert.NotNil(t, mx)

	check := v.Check()
	assert.True(t, check)
}

// TestVarnish_InitAndCollect tests full Init and Collect flow
func TestVarnish_InitAndCollect(t *testing.T) {
	v := New()
	v.client = &mockVarnishClient{
		stats: map[string]interface{}{
			"MAIN.uptime":           float64(1000),
			"MAIN.sess_conn":        float64(100),
			"MAIN.sess_drop":        float64(10),
			"MAIN.sess_fail":        float64(5),
			"MAIN.client_req":       float64(500),
			"MAIN.cache_hit":        float64(400),
			"MAIN.cache_miss":       float64(100),
			"MAIN.backend_conn":     float64(50),
			"MAIN.backend_unhealthy": float64(2),
		},
	}

	err := v.Init()
	assert.NoError(t, err)

	check := v.Check()
	assert.True(t, check)

	mx := v.Collect()
	assert.NotNil(t, mx)

	v.Cleanup()
	assert.Nil(t, v.client)
}

// TestVarnish_InitCheckCollectCleanup tests complete lifecycle
func TestVarnish_InitCheckCollectCleanup(t *testing.T) {
	v := New()
	v.client = &mockVarnishClient{
		stats: map[string]interface{}{
			"MAIN.uptime": float64(1000),
		},
	}

	// Init
	err := v.Init()
	require.NoError(t, err)

	// Check
	check := v.Check()
	require.True(t, check)

	// Collect
	mx := v.Collect()
	require.NotNil(t, mx)

	// Cleanup
	v.Cleanup()
	assert.Nil(t, v.client)
}

// TestVarnish_ChartsList tests charts structure
func TestVarnish_ChartsList(t *testing.T) {
	v := New()

	charts := v.Charts()

	require.NotNil(t, charts)
	for _, chart := range charts {
		assert.NotEmpty(t, chart.ID)
		assert.NotEmpty(t, chart.Title)
	}
}

// TestVarnish_ConfigStructure tests configuration structure
func TestVarnish_ConfigStructure(t *testing.T) {
	v := New()

	assert.NotNil(t, v.Config)
	assert.NotZero(t, v.Timeout)
}

// TestVarnish_NameProperty tests Name property
func TestVarnish_NameProperty(t *testing.T) {
	v := New()
	v.Name = "test-varnish"

	assert.Equal(t, "test-varnish", v.Name)
}

// TestVarnish_CollectWithVariousStats tests Collect with various stat combinations
func TestVarnish_CollectWithVariousStats(t *testing.T) {
	testCases := []struct {
		name  string
		stats map[string]interface{}
	}{
		{
			name: "connection stats only",
			stats: map[string]interface{}{
				"MAIN.sess_conn": float64(100),
				"MAIN.sess_drop": float64(10),
				"MAIN.sess_fail": float64(5),
			},
		},
		{
			name: "cache stats only",
			stats: map[string]interface{}{
				"MAIN.cache_hit":  float64(400),
				"MAIN.cache_miss": float64(100),
			},
		},
		{
			name: "backend stats only",
			stats: map[string]interface{}{
				"MAIN.backend_conn":      float64(50),
				"MAIN.backend_unhealthy": float64(2),
			},
		},
		{
			name: "fetch stats only",
			stats: map[string]interface{}{
				"MAIN.fetch_head":       float64(20),
				"MAIN.fetch_length":     float64(80),
				"MAIN.fetch_bad":        float64(1),
				"MAIN.fetch_failed":     float64(0),
				"MAIN.fetch_no_thread":  float64(1),
			},
		},
		{
			name: "thread stats only",
			stats: map[string]interface{}{
				"MAIN.threads":         float64(4),
				"MAIN.threads_created": float64(4),
				"MAIN.threads_failed":  float64(0),
				"MAIN.threads_limited": float64(0),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := New()
			v.client = &mockVarnishClient{
				stats: tc.stats,
			}

			mx := v.Collect()
			assert.NotNil(t, mx)
		})
	}
}

// TestVarnish_CheckFailureScenarios tests various Check failure scenarios
func TestVarnish_CheckFailureScenarios(t *testing.T) {
	testCases := []struct {
		name        string
		client      varnishClient
		expectCheck bool
	}{
		{
			name:        "nil client",
			client:      nil,
			expectCheck: false,
		},
		{
			name: "client error",
			client: &mockVarnishClient{
				err: fmt.Errorf("test error"),
			},
			expectCheck: false,
		},
		{
			name: "empty stats",
			client: &mockVarnishClient{
				stats: map[string]interface{}{},
			},
			expectCheck: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			v := New()
			v.client = tc.client

			result := v.Check()
			assert.Equal(t, tc.expectCheck, result)
		})
	}
}

// Interface for mocking
type varnishClient interface {
	Stats() (map[string]interface{}, error)
}

// defaultConfig represents default configuration
var defaultConfig = struct {
	Timeout time.Duration
}{
	Timeout: 10 * time.Second,
}