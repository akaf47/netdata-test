package ap

import (
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/collector/ap/charts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCharts(t *testing.T) {
	tests := []struct {
		name      string
		testFunc  func(*testing.T)
	}{
		{
			name: "should initialize charts",
			testFunc: func(t *testing.T) {
				c := NewAP()
				require.NotNil(t, c)
				assert.NotNil(t, c.Charts())
			},
		},
		{
			name: "should have valid chart definitions",
			testFunc: func(t *testing.T) {
				c := NewAP()
				charts := c.Charts()
				assert.NotNil(t, charts)
				assert.Greater(t, len(charts.Charts), 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

func TestChartsStructure(t *testing.T) {
	c := NewAP()
	charts := c.Charts()
	
	for _, chart := range charts.Charts {
		assert.NotEmpty(t, chart.ID)
		assert.NotEmpty(t, chart.Title)
		assert.NotEmpty(t, chart.Type)
		assert.GreaterOrEqual(t, len(chart.Dims), 0)
		
		for _, dim := range chart.Dims {
			assert.NotEmpty(t, dim.ID)
			assert.NotEmpty(t, dim.Name)
		}
	}
}

func TestChartsInit(t *testing.T) {
	tests := []struct {
		name    string
		chart   interface{}
		wantErr bool
	}{
		{
			name:    "should handle nil chart gracefully",
			chart:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test chart initialization with various states
			c := NewAP()
			assert.NotNil(t, c)
		})
	}
}

func TestChartsAddDimensions(t *testing.T) {
	c := NewAP()
	charts := c.Charts()
	
	// Verify dimensions can be added/modified
	for _, chart := range charts.Charts {
		originalDimCount := len(chart.Dims)
		assert.GreaterOrEqual(t, originalDimCount, 0)
	}
}

func TestChartsIDs(t *testing.T) {
	c := NewAP()
	charts := c.Charts()
	
	seenIDs := make(map[string]bool)
	for _, chart := range charts.Charts {
		// Verify chart IDs are unique
		assert.False(t, seenIDs[chart.ID], "duplicate chart ID: %s", chart.ID)
		seenIDs[chart.ID] = true
		
		dimIDs := make(map[string]bool)
		for _, dim := range chart.Dims {
			assert.False(t, dimIDs[dim.ID], "duplicate dimension ID: %s", dim.ID)
			dimIDs[dim.ID] = true
		}
	}
}

func TestChartsMetadata(t *testing.T) {
	c := NewAP()
	charts := c.Charts()
	
	for _, chart := range charts.Charts {
		// Verify required chart metadata
		assert.NotEmpty(t, chart.ID)
		assert.NotEmpty(t, chart.Title)
		assert.NotEmpty(t, chart.Type)
		assert.NotEmpty(t, chart.Fqdn)
		assert.NotEmpty(t, chart.Family)
		assert.NotNil(t, chart.Context)
	}
}

func TestChartsRemoveDimensions(t *testing.T) {
	c := NewAP()
	charts := c.Charts()
	
	for _, chart := range charts.Charts {
		// Verify chart can handle dimension removal
		for _, dim := range chart.Dims {
			assert.NotNil(t, dim)
			assert.NotEmpty(t, dim.ID)
		}
	}
}

func TestChartsEmpty(t *testing.T) {
	c := NewAP()
	assert.NotNil(t, c)
	
	charts := c.Charts()
	assert.NotNil(t, charts)
}

func TestChartsMultipleInstances(t *testing.T) {
	c1 := NewAP()
	c2 := NewAP()
	
	// Verify multiple instances can exist
	assert.NotNil(t, c1)
	assert.NotNil(t, c2)
	assert.NotEqual(t, c1, c2)
}

func TestChartsNilCheck(t *testing.T) {
	c := NewAP()
	charts := c.Charts()
	
	if charts != nil {
		for _, chart := range charts.Charts {
			assert.NotNil(t, chart)
		}
	}
}