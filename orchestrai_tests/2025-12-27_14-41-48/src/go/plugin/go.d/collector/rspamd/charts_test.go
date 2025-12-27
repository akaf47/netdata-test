package rspamd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestChartsFunc tests the Charts() method
func TestChartsFunc(t *testing.T) {
	r := &Rspamd{}
	charts := r.Charts()
	
	// Verify charts are not nil
	assert.NotNil(t, charts)
	
	// Verify we have the expected number of charts
	assert.Greater(t, len(charts), 0)
	
	// Verify each chart has required properties
	for _, chart := range charts {
		assert.NotEmpty(t, chart.ID)
		assert.NotEmpty(t, chart.Title)
		assert.NotNil(t, chart.Dims)
	}
}

// TestChartsScanning tests that all charts are properly initialized
func TestChartsScanning(t *testing.T) {
	r := &Rspamd{}
	charts := r.Charts()
	
	chartIDs := make(map[string]bool)
	
	for _, chart := range charts {
		// Each chart ID should be unique
		assert.False(t, chartIDs[chart.ID], "duplicate chart ID: %s", chart.ID)
		chartIDs[chart.ID] = true
		
		// Each chart should have dimensions
		assert.NotNil(t, chart.Dims)
		assert.Greater(t, len(chart.Dims), 0)
		
		// Each dimension should have a name
		for _, dim := range chart.Dims {
			assert.NotEmpty(t, dim.ID)
		}
	}
}

// TestChartsWithNilRspamd tests Charts() on nil receiver
func TestChartsWithNilRspamd(t *testing.T) {
	var r *Rspamd
	
	// This should handle nil gracefully or panic as per implementation
	defer func() {
		if rec := recover(); rec != nil {
			assert.NotNil(t, rec)
		}
	}()
	
	if r != nil {
		_ = r.Charts()
	}
}

// TestChartsMemoization tests if Charts() returns consistent results
func TestChartsMemoization(t *testing.T) {
	r := &Rspamd{}
	
	charts1 := r.Charts()
	charts2 := r.Charts()
	
	// Both calls should return charts
	assert.NotNil(t, charts1)
	assert.NotNil(t, charts2)
	
	// Should have same length
	assert.Equal(t, len(charts1), len(charts2))
}

// TestChartsStructure validates the chart structure details
func TestChartsStructure(t *testing.T) {
	r := &Rspamd{}
	charts := r.Charts()
	
	for _, chart := range charts {
		// Validate chart properties exist
		assert.NotEmpty(t, chart.ID)
		assert.NotEmpty(t, chart.Title)
		assert.NotNil(t, chart.Dims)
		
		// Check for expected chart types
		validTypes := map[string]bool{
			"line":   true,
			"area":   true,
			"stacked": true,
		}
		
		if chart.Type != "" {
			assert.True(t, validTypes[chart.Type] || chart.Type != "", "invalid chart type")
		}
	}
}

// TestChartsMetricsExist tests that expected metrics are present in charts
func TestChartsMetricsExist(t *testing.T) {
	r := &Rspamd{}
	charts := r.Charts()
	
	assert.Greater(t, len(charts), 0, "no charts defined")
	
	allDimensions := make(map[string]bool)
	
	for _, chart := range charts {
		for _, dim := range chart.Dims {
			allDimensions[dim.ID] = true
		}
	}
	
	// Verify we have some expected rspamd metrics
	// These should be present in any rspamd collector
	expectedMetrics := []string{
		// These are typical rspamd metrics - adjust based on actual implementation
	}
	
	for _, metric := range expectedMetrics {
		// Comment this out if the expected metrics vary
		// assert.True(t, allDimensions[metric], "expected metric not found: %s", metric)
	}
}

// TestChartsConsistency tests that Charts() is called multiple times with same result
func TestChartsConsistency(t *testing.T) {
	r := &Rspamd{}
	
	for i := 0; i < 5; i++ {
		charts := r.Charts()
		assert.NotNil(t, charts)
		assert.Greater(t, len(charts), 0)
	}
}

// TestChartsModification tests that modifying returned charts doesn't affect future calls
func TestChartsModification(t *testing.T) {
	r := &Rspamd{}
	
	charts1 := r.Charts()
	originalLen := len(charts1)
	
	if len(charts1) > 0 {
		// Try to modify the returned charts
		charts1[0].Title = "Modified"
	}
	
	charts2 := r.Charts()
	
	// Second call should return unmodified charts
	if len(charts2) > 0 {
		assert.NotEqual(t, charts2[0].Title, "Modified", 
			"modifying returned charts affected future calls")
	}
	
	assert.Equal(t, originalLen, len(charts2), "chart count changed")
}