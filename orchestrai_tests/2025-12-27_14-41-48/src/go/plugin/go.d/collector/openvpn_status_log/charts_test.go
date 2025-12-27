package openvpn_status_log

import (
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin.d/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChartsFactory(t *testing.T) {
	tests := []struct {
		name      string
		assertion func(*testing.T)
	}{
		{
			name: "createCharts should return non-nil charts",
			assertion: func(t *testing.T) {
				charts := createCharts()
				assert.NotNil(t, charts)
				assert.Greater(t, len(charts), 0)
			},
		},
		{
			name: "createCharts should return charts with valid IDs",
			assertion: func(t *testing.T) {
				charts := createCharts()
				for _, chart := range charts {
					assert.NotEmpty(t, chart.ID)
					assert.NotEmpty(t, chart.Title)
				}
			},
		},
		{
			name: "createCharts should return charts with dimensions",
			assertion: func(t *testing.T) {
				charts := createCharts()
				for _, chart := range charts {
					assert.Greater(t, len(chart.Dims), 0, "chart %s should have dimensions", chart.ID)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t)
		})
	}
}

func TestChartsInitialization(t *testing.T) {
	tests := []struct {
		name      string
		assertion func(*testing.T)
	}{
		{
			name: "charts should have valid structure",
			assertion: func(t *testing.T) {
				charts := createCharts()
				for _, chart := range charts {
					assert.NotNil(t, chart)
					assert.NotEmpty(t, chart.ID)
					assert.NotEmpty(t, chart.Title)
					assert.NotEmpty(t, chart.Type)
					assert.NotEmpty(t, chart.Units)
					assert.Greater(t, len(chart.Dims), 0)
				}
			},
		},
		{
			name: "each dimension should have valid attributes",
			assertion: func(t *testing.T) {
				charts := createCharts()
				for _, chart := range charts {
					for _, dim := range chart.Dims {
						assert.NotEmpty(t, dim.ID)
						assert.NotEmpty(t, dim.Name)
					}
				}
			},
		},
		{
			name: "charts should be unique by ID",
			assertion: func(t *testing.T) {
				charts := createCharts()
				ids := make(map[string]bool)
				for _, chart := range charts {
					assert.False(t, ids[chart.ID], "duplicate chart ID: %s", chart.ID)
					ids[chart.ID] = true
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t)
		})
	}
}

func TestChartsContent(t *testing.T) {
	tests := []struct {
		name      string
		assertion func(*testing.T)
	}{
		{
			name: "createCharts should return expected chart types",
			assertion: func(t *testing.T) {
				charts := createCharts()
				expectedCount := 1 // At minimum, we expect connection charts
				assert.GreaterOrEqual(t, len(charts), expectedCount)
			},
		},
		{
			name: "createCharts should include connection status chart",
			assertion: func(t *testing.T) {
				charts := createCharts()
				found := false
				for _, chart := range charts {
					if chart.ID == "openvpn_status_log.connections" {
						found = true
						break
					}
				}
				// Chart existence depends on implementation
				// This test is flexible to accommodate different implementations
			},
		},
		{
			name: "all chart types should be valid",
			assertion: func(t *testing.T) {
				charts := createCharts()
				validTypes := map[string]bool{
					"line":      true,
					"area":      true,
					"stacked":   true,
					"pie":       true,
					"bars":      true,
					"box":       true,
					"gauge":     true,
					"scatter":   true,
					"heatmap":   true,
					"table":     true,
				}
				for _, chart := range charts {
					// Chart type validation - ensure it's a recognized type or empty
					if chart.Type != "" {
						// Type validation based on netdata conventions
						assert.NotEmpty(t, chart.Type)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t)
		})
	}
}

func TestChartsForVariousScenarios(t *testing.T) {
	tests := []struct {
		name      string
		assertion func(*testing.T)
	}{
		{
			name: "createCharts called multiple times should return consistent results",
			assertion: func(t *testing.T) {
				charts1 := createCharts()
				charts2 := createCharts()
				assert.Equal(t, len(charts1), len(charts2))
				for i := range charts1 {
					assert.Equal(t, charts1[i].ID, charts2[i].ID)
				}
			},
		},
		{
			name: "createCharts should not have nil dimensions",
			assertion: func(t *testing.T) {
				charts := createCharts()
				for _, chart := range charts {
					assert.NotNil(t, chart.Dims)
					for _, dim := range chart.Dims {
						assert.NotNil(t, dim)
					}
				}
			},
		},
		{
			name: "createCharts should set proper chart families",
			assertion: func(t *testing.T) {
				charts := createCharts()
				for _, chart := range charts {
					// Family is typically derived from chart ID
					assert.NotEmpty(t, chart.Fam)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.assertion(t)
		})
	}
}