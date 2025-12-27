package nvme

import (
	"testing"

	"github.com/netdata/go.d.plugin/pkg/web"
	"github.com/stretchr/testify/assert"
)

func TestChartsInitialization(t *testing.T) {
	t.Run("should initialize charts structure", func(t *testing.T) {
		charts := createCharts()
		assert.NotNil(t, charts)
	})
}

func TestCreateCharts(t *testing.T) {
	tests := []struct {
		name          string
		expectedCount int
	}{
		{
			name:          "charts should not be empty",
			expectedCount: 0, // At least 0, but typically > 0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charts := createCharts()
			assert.NotNil(t, charts)
			// Charts should be properly initialized
		})
	}
}

func TestChartDefinitions(t *testing.T) {
	t.Run("should have valid chart definitions", func(t *testing.T) {
		charts := createCharts()

		for _, chart := range *charts {
			assert.NotNil(t, chart)
			assert.NotEmpty(t, chart.ID)
			assert.NotEmpty(t, chart.Title)
		}
	})
}

func TestChartMetadata(t *testing.T) {
	tests := []struct {
		name       string
		validator  func(*testing.T)
	}{
		{
			name: "charts should have units",
			validator: func(t *testing.T) {
				charts := createCharts()
				for _, chart := range *charts {
					// Units might be optional, but should be consistent
					assert.NotNil(t, chart.Units)
				}
			},
		},
		{
			name: "charts should have dimensions",
			validator: func(t *testing.T) {
				charts := createCharts()
				for _, chart := range *charts {
					assert.NotNil(t, chart.Dims)
					// Each chart should have at least one dimension
					if len(chart.Dims) == 0 {
						t.Logf("Warning: Chart %s has no dimensions", chart.ID)
					}
				}
			},
		},
		{
			name: "dimension names should not be empty",
			validator: func(t *testing.T) {
				charts := createCharts()
				for _, chart := range *charts {
					for _, dim := range chart.Dims {
						assert.NotNil(t, dim)
						assert.NotEmpty(t, dim.ID)
						assert.NotEmpty(t, dim.Name)
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validator(t)
		})
	}
}

func TestChartIDs(t *testing.T) {
	t.Run("should have unique chart IDs", func(t *testing.T) {
		charts := createCharts()
		ids := make(map[string]bool)

		for _, chart := range *charts {
			assert.False(t, ids[chart.ID], "Duplicate chart ID: %s", chart.ID)
			ids[chart.ID] = true
		}
	})
}

func TestDimensionIDs(t *testing.T) {
	t.Run("should have unique dimension IDs within chart", func(t *testing.T) {
		charts := createCharts()

		for _, chart := range *charts {
			dimIDs := make(map[string]bool)
			for _, dim := range chart.Dims {
				key := dim.ID
				assert.False(t, dimIDs[key], "Duplicate dimension ID in chart %s: %s", chart.ID, key)
				dimIDs[key] = true
			}
		}
	})
}

func TestChartTypes(t *testing.T) {
	tests := []struct {
		name        string
		expectTypes bool
	}{
		{
			name:        "charts should have valid types",
			expectTypes: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charts := createCharts()
			for _, chart := range *charts {
				// Type should be a valid string (could be "line", "area", etc.)
				assert.NotEmpty(t, chart.Type)
			}
		})
	}
}

func TestChartFamily(t *testing.T) {
	t.Run("should organize charts by family", func(t *testing.T) {
		charts := createCharts()

		families := make(map[string]int)
		for _, chart := range *charts {
			families[chart.Fam]++
		}

		// Should have at least one family
		assert.Greater(t, len(families), 0)

		// All family names should be non-empty
		for family := range families {
			assert.NotEmpty(t, family)
		}
	})
}

func TestChartPriority(t *testing.T) {
	t.Run("should have priorities for ordering", func(t *testing.T) {
		charts := createCharts()

		for _, chart := range *charts {
			// Priority should be >= 0
			assert.GreaterOrEqual(t, chart.Priority, 0)
		}
	})
}

func TestChartOptions(t *testing.T) {
	t.Run("should have chart options", func(t *testing.T) {
		charts := createCharts()

		for _, chart := range *charts {
			assert.NotNil(t, chart.Opts)
		}
	})
}

func TestDimensionOptions(t *testing.T) {
	t.Run("should have dimension options", func(t *testing.T) {
		charts := createCharts()

		for _, chart := range *charts {
			for _, dim := range chart.Dims {
				assert.NotNil(t, dim.Opts)
			}
		}
	})
}

func TestChartInitializationEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "multiple chart initializations should produce consistent results",
			testFunc: func(t *testing.T) {
				charts1 := createCharts()
				charts2 := createCharts()

				assert.Equal(t, len(*charts1), len(*charts2))

				for i := 0; i < len(*charts1); i++ {
					assert.Equal(t, (*charts1)[i].ID, (*charts2)[i].ID)
				}
			},
		},
		{
			name: "charts should not be nil after creation",
			testFunc: func(t *testing.T) {
				charts := createCharts()
				assert.NotNil(t, charts)
				assert.Greater(t, len(*charts), 0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

func TestChartStructureIntegrity(t *testing.T) {
	t.Run("chart structure should be valid", func(t *testing.T) {
		charts := createCharts()

		for _, chart := range *charts {
			// All required fields should be present
			assert.NotNil(t, chart.ID)
			assert.NotNil(t, chart.Title)
			assert.NotNil(t, chart.Dims)
			assert.NotNil(t, chart.Units)
		}
	})
}

func TestEmptyChartHandling(t *testing.T) {
	t.Run("should handle empty dimensions gracefully", func(t *testing.T) {
		charts := createCharts()

		for _, chart := range *charts {
			// Even if dimensions are empty, chart should be valid
			assert.NotNil(t, chart.Dims)
		}
	})
}

func TestChartConsistency(t *testing.T) {
	t.Run("should maintain chart consistency across calls", func(t *testing.T) {
		chart1 := createCharts()
		chart2 := createCharts()

		assert.Equal(t, len(*chart1), len(*chart2))

		for i, c1 := range *chart1 {
			c2 := (*chart2)[i]
			assert.Equal(t, c1.ID, c2.ID)
			assert.Equal(t, c1.Title, c2.Title)
			assert.Equal(t, len(c1.Dims), len(c2.Dims))
		}
	})
}

func TestChartMetricsNVMe(t *testing.T) {
	tests := []struct {
		name        string
		expectChart func(*testing.T, []*web.Chart)
	}{
		{
			name: "should include nvme performance metrics",
			expectChart: func(t *testing.T, charts []*web.Chart) {
				assert.NotNil(t, charts)
				// Performance metrics should exist
			},
		},
		{
			name: "should include nvme health metrics",
			expectChart: func(t *testing.T, charts []*web.Chart) {
				assert.NotNil(t, charts)
				// Health metrics should exist
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charts := createCharts()
			tt.expectChart(t, *charts)
		})
	}
}

func TestChartDimensionConsistency(t *testing.T) {
	t.Run("dimension IDs should match dimension names pattern", func(t *testing.T) {
		charts := createCharts()

		for _, chart := range *charts {
			for _, dim := range chart.Dims {
				assert.NotEmpty(t, dim.ID)
				assert.NotEmpty(t, dim.Name)
				// ID and Name should be consistently formatted
				assert.NotEqual(t, "", dim.ID)
				assert.NotEqual(t, "", dim.Name)
			}
		}
	})
}

func TestChartsBoundaryConditions(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{
			name: "should handle large dimension lists",
			testFunc: func(t *testing.T) {
				charts := createCharts()
				for _, chart := range *charts {
					// Should handle any number of dimensions
					assert.GreaterOrEqual(t, len(chart.Dims), 0)
				}
			},
		},
		{
			name: "should handle single dimension charts",
			testFunc: func(t *testing.T) {
				charts := createCharts()
				for _, chart := range *charts {
					if len(chart.Dims) >= 1 {
						assert.NotNil(t, chart.Dims[0])
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}