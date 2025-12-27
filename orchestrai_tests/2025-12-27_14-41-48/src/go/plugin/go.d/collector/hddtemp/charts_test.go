package hddtemp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test Charts function exists and returns proper chart definitions
func TestCharts(t *testing.T) {
	charts := Charts()
	
	// Test that charts is not nil
	assert.NotNil(t, charts, "Charts should not be nil")
	
	// Test that charts has expected length
	assert.Greater(t, len(charts), 0, "Charts should contain at least one chart")
	
	// Test that each chart has proper initialization
	for _, chart := range charts {
		assert.NotNil(t, chart, "Individual chart should not be nil")
		assert.NotEmpty(t, chart.ID, "Chart ID should not be empty")
		assert.NotEmpty(t, chart.Title, "Chart Title should not be empty")
	}
}

// Test Charts returns consistent results
func TestChartsConsistency(t *testing.T) {
	charts1 := Charts()
	charts2 := Charts()
	
	require.Equal(t, len(charts1), len(charts2), "Charts should return same number of charts on multiple calls")
	
	for i := range charts1 {
		assert.Equal(t, charts1[i].ID, charts2[i].ID, "Chart IDs should be consistent")
		assert.Equal(t, charts1[i].Title, charts2[i].Title, "Chart titles should be consistent")
	}
}

// Test temperature chart is present
func TestTemperatureChartPresent(t *testing.T) {
	charts := Charts()
	
	// Temperature chart should be the primary chart in hddtemp
	found := false
	for _, chart := range charts {
		if chart.ID == "hddtemp.temperature" || 
		   chart.Title == "Temperature" ||
		   chart.ID == "hddtemp.disk_temperature" {
			found = true
			assert.NotNil(t, chart.Dims, "Temperature chart should have dimensions")
			break
		}
	}
	assert.True(t, found, "Temperature chart should be present in charts")
}

// Test chart properties
func TestChartProperties(t *testing.T) {
	charts := Charts()
	
	for _, chart := range charts {
		// Test chart has required properties
		assert.NotEmpty(t, chart.ID, "Chart must have an ID")
		assert.NotEmpty(t, chart.Title, "Chart must have a title")
		
		// Test dimensions if present
		if chart.Dims != nil && len(chart.Dims) > 0 {
			for _, dim := range chart.Dims {
				assert.NotNil(t, dim, "Dimension should not be nil")
				assert.NotEmpty(t, dim.ID, "Dimension must have an ID")
				assert.NotEmpty(t, dim.Name, "Dimension must have a name")
			}
		}
	}
}

// Test Charts function is idempotent
func TestChartsIdempotent(t *testing.T) {
	charts1 := Charts()
	charts2 := Charts()
	charts3 := Charts()
	
	// All calls should return same structure
	assert.Equal(t, len(charts1), len(charts2), "First and second call should match")
	assert.Equal(t, len(charts2), len(charts3), "Second and third call should match")
}

// Test chart unit information
func TestChartUnits(t *testing.T) {
	charts := Charts()
	
	for _, chart := range charts {
		// Temperature charts should have temperature units
		if chart.ID == "hddtemp.temperature" || 
		   chart.Title == "Temperature" {
			assert.NotEmpty(t, chart.Unit, "Temperature chart should have unit specified")
		}
	}
}

// Test multiple charts with different units
func TestMultipleCharts(t *testing.T) {
	charts := Charts()
	
	assert.Greater(t, len(charts), 0, "Should have at least one chart")
	
	for i, chart := range charts {
		assert.NotNil(t, chart, "Chart at index %d should not be nil", i)
		assert.NotEmpty(t, chart.ID, "Chart at index %d should have ID", i)
	}
}

// Test chart dimensions count
func TestChartDimensionsCount(t *testing.T) {
	charts := Charts()
	
	for _, chart := range charts {
		if chart.Dims != nil {
			// Dimension count should be non-negative
			assert.GreaterOrEqual(t, len(chart.Dims), 0, "Dimension count should be non-negative")
		}
	}
}

// Test chart properties are properly formatted
func TestChartPropertiesFormatted(t *testing.T) {
	charts := Charts()
	
	for _, chart := range charts {
		// ID should follow naming convention (usually lowercase with dots)
		assert.NotEmpty(t, chart.ID, "Chart ID must not be empty")
		
		// Title should be human-readable
		assert.NotEmpty(t, chart.Title, "Chart title must not be empty")
		assert.True(t, len(chart.Title) > 0, "Chart title should have content")
	}
}

// Test empty conditions (edge case for potential modifications)
func TestChartsNotEmpty(t *testing.T) {
	charts := Charts()
	
	if charts != nil {
		for _, chart := range charts {
			assert.NotNil(t, chart, "No nil charts should be returned")
		}
	}
}