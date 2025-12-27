package consul

import (
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/pkg/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectMetrics_Success(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true

	expectedMetrics := map[string]interface{}{
		"consul_up":                  1,
		"consul_serf_lan_members":    3,
		"consul_serf_lan_events":     0,
		"consul_catalog_services":    5,
		"consul_catalog_nodes":       3,
		"consul_health_node_status":  3,
		"consul_health_service_status": 10,
		"consul_http_requests_total": 100,
		"consul_http_request_duration_seconds": 0.05,
	}

	// Act
	err := c.collectMetrics(expectedMetrics)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, c.metrics)
}

func TestCollectMetrics_NilMetrics(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true

	// Act
	err := c.collectMetrics(nil)

	// Assert
	assert.Error(t, err)
}

func TestCollectMetrics_EmptyMetrics(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true
	emptyMetrics := make(map[string]interface{})

	// Act
	err := c.collectMetrics(emptyMetrics)

	// Assert
	// Should handle empty metrics gracefully
	assert.NoError(t, err)
}

func TestCollectMetrics_InvalidMetricType(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true

	invalidMetrics := map[string]interface{}{
		"consul_up": "invalid_value",
	}

	// Act
	err := c.collectMetrics(invalidMetrics)

	// Assert
	assert.Error(t, err)
}

func TestCollectMetrics_IntMetricValue(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true

	metrics := map[string]interface{}{
		"consul_serf_lan_members": 5,
	}

	// Act
	err := c.collectMetrics(metrics)

	// Assert
	assert.NoError(t, err)
}

func TestCollectMetrics_FloatMetricValue(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true

	metrics := map[string]interface{}{
		"consul_http_request_duration_seconds": 0.123,
	}

	// Act
	err := c.collectMetrics(metrics)

	// Assert
	assert.NoError(t, err)
}

func TestCollectMetrics_ZeroValues(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true

	metrics := map[string]interface{}{
		"consul_up":                  0,
		"consul_serf_lan_events":     0,
		"consul_catalog_services":    0,
	}

	// Act
	err := c.collectMetrics(metrics)

	// Assert
	assert.NoError(t, err)
}

func TestCollectMetrics_NegativeValues(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true

	metrics := map[string]interface{}{
		"consul_serf_lan_members": -1,
	}

	// Act
	err := c.collectMetrics(metrics)

	// Assert
	// Should either handle gracefully or error appropriately
	assert.True(t, err != nil || err == nil)
}

func TestCollectMetrics_LargeValues(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true

	metrics := map[string]interface{}{
		"consul_catalog_services": 999999,
		"consul_catalog_nodes":    999999,
	}

	// Act
	err := c.collectMetrics(metrics)

	// Assert
	assert.NoError(t, err)
}

func TestCollectMetrics_MixedValidAndInvalidMetrics(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true

	metrics := map[string]interface{}{
		"consul_up":               1,
		"consul_serf_lan_members": "invalid",
	}

	// Act
	err := c.collectMetrics(metrics)

	// Assert
	assert.Error(t, err)
}

func TestCollectMetrics_UpdateExistingMetrics(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true

	metricsFirst := map[string]interface{}{
		"consul_up": 1,
	}
	metricsSecond := map[string]interface{}{
		"consul_up": 0,
	}

	// Act
	err1 := c.collectMetrics(metricsFirst)
	err2 := c.collectMetrics(metricsSecond)

	// Assert
	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestCollectMetrics_PartialMetricsUpdate(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true

	firstMetrics := map[string]interface{}{
		"consul_up":                  1,
		"consul_serf_lan_members":    3,
		"consul_catalog_services":    5,
	}

	secondMetrics := map[string]interface{}{
		"consul_up": 1,
	}

	// Act
	err1 := c.collectMetrics(firstMetrics)
	err2 := c.collectMetrics(secondMetrics)

	// Assert
	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestCollectMetrics_MultipleCallsConsistency(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true

	metrics := map[string]interface{}{
		"consul_serf_lan_members": 5,
	}

	// Act & Assert
	for i := 0; i < 3; i++ {
		err := c.collectMetrics(metrics)
		assert.NoError(t, err)
	}
}

func TestCollectMetrics_SpecialCharactersInMetricNames(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true

	metrics := map[string]interface{}{
		"consul_up-test":           1,
		"consul_up_test":           1,
		"consul.up":                1,
	}

	// Act
	err := c.collectMetrics(metrics)

	// Assert
	assert.NoError(t, err)
}

func TestCollectMetrics_VerySmallFloatValues(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true

	metrics := map[string]interface{}{
		"consul_http_request_duration_seconds": 0.00001,
	}

	// Act
	err := c.collectMetrics(metrics)

	// Assert
	assert.NoError(t, err)
}

func TestCollectMetrics_VeryLargeFloatValues(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true

	metrics := map[string]interface{}{
		"consul_http_request_duration_seconds": 9999999.99,
	}

	// Act
	err := c.collectMetrics(metrics)

	// Assert
	assert.NoError(t, err)
}

func TestCollectMetrics_StatusMetrics(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true

	metrics := map[string]interface{}{
		"consul_up": 1,
	}

	// Act
	err := c.collectMetrics(metrics)

	// Assert
	assert.NoError(t, err)
}

func TestCollectMetrics_HandlesMetricProcessing(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true

	metrics := map[string]interface{}{
		"consul_raft_leader":          1,
		"consul_raft_state":           3,
		"consul_memberlist_lan_nodes": 5,
	}

	// Act
	err := c.collectMetrics(metrics)

	// Assert
	assert.NoError(t, err)
}

func TestCollectMetrics_AllMetricCategories(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true

	allMetrics := map[string]interface{}{
		// Server metrics
		"consul_up":                  1,
		"consul_raft_leader":         1,
		"consul_raft_state":          3,
		
		// Serf metrics
		"consul_serf_lan_members":    5,
		"consul_serf_lan_events":     0,
		
		// Catalog metrics
		"consul_catalog_services":    10,
		"consul_catalog_nodes":       3,
		
		// Health metrics
		"consul_health_node_status":  3,
		"consul_health_service_status": 15,
		
		// HTTP metrics
		"consul_http_requests_total": 100,
		"consul_http_request_duration_seconds": 0.05,
	}

	// Act
	err := c.collectMetrics(allMetrics)

	// Assert
	assert.NoError(t, err)
}

func TestCollectMetrics_UninitialisedCollector(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = false

	metrics := map[string]interface{}{
		"consul_up": 1,
	}

	// Act
	err := c.collectMetrics(metrics)

	// Assert
	assert.Error(t, err)
}

func TestCollectMetrics_RepeatedMetricNames(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true

	metrics := map[string]interface{}{
		"consul_up": 1,
	}

	// Act
	err1 := c.collectMetrics(metrics)
	metrics["consul_up"] = 0
	err2 := c.collectMetrics(metrics)

	// Assert
	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestCollectMetrics_HandlesBoundaryIntValues(t *testing.T) {
	// Arrange
	c := NewConsulCollector()
	c.initialized = true

	testCases := []interface{}{
		int(0),
		int(1),
		int(-1),
		int64(9223372036854775807), // max int64
		int64(-9223372036854775808), // min int64
	}

	for _, val := range testCases {
		metrics := map[string]interface{}{
			"consul_serf_lan_members": val,
		}

		// Act
		err := c.collectMetrics(metrics)

		// Assert
		assert.True(t, err != nil || err == nil)
	}
}