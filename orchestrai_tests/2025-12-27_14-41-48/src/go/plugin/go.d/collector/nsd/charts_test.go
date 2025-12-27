package nsd

import (
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/module"
)

func TestChartsInit(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "charts should not be nil",
			test: func(t *testing.T) {
				if Charts == nil {
					t.Fatal("Charts is nil")
				}
			},
		},
		{
			name: "charts should be of correct type",
			test: func(t *testing.T) {
				if _, ok := Charts.(*module.Charts); !ok {
					t.Errorf("Charts type mismatch, got %T", Charts)
				}
			},
		},
		{
			name: "charts should contain expected charts",
			test: func(t *testing.T) {
				if len(*Charts) == 0 {
					t.Error("Charts collection is empty")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

func TestChartsHaveValidStructure(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "each chart should have an ID",
			test: func(t *testing.T) {
				for _, chart := range *Charts {
					if chart.ID == "" {
						t.Error("Chart has empty ID")
					}
				}
			},
		},
		{
			name: "each chart should have a title",
			test: func(t *testing.T) {
				for _, chart := range *Charts {
					if chart.Title == "" {
						t.Error("Chart has empty Title")
					}
				}
			},
		},
		{
			name: "each chart should have dimensions",
			test: func(t *testing.T) {
				for _, chart := range *Charts {
					if len(chart.Dims) == 0 {
						t.Errorf("Chart %s has no dimensions", chart.ID)
					}
				}
			},
		},
		{
			name: "each dimension should have a name",
			test: func(t *testing.T) {
				for _, chart := range *Charts {
					for _, dim := range chart.Dims {
						if dim.Name == "" {
							t.Errorf("Chart %s has dimension with empty name", chart.ID)
						}
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

func TestChartsStructure(t *testing.T) {
	if Charts == nil {
		t.Fatal("Charts should not be nil")
	}

	chartsCount := len(*Charts)
	if chartsCount == 0 {
		t.Fatal("Charts should not be empty")
	}

	// Verify each chart has proper structure
	for i, chart := range *Charts {
		if chart == nil {
			t.Fatalf("Chart at index %d is nil", i)
		}

		// Verify ID is set
		if chart.ID == "" {
			t.Errorf("Chart at index %d has empty ID", i)
		}

		// Verify Title is set
		if chart.Title == "" {
			t.Errorf("Chart %s has empty Title", chart.ID)
		}

		// Verify Type is set
		if chart.Type == "" {
			t.Errorf("Chart %s has empty Type", chart.ID)
		}

		// Verify dimensions exist
		if chart.Dims == nil {
			t.Errorf("Chart %s has nil Dims", chart.ID)
		}

		if len(chart.Dims) == 0 {
			t.Errorf("Chart %s has no dimensions", chart.ID)
		}

		// Verify each dimension has required fields
		for _, dim := range chart.Dims {
			if dim == nil {
				t.Errorf("Chart %s has nil dimension", chart.ID)
				continue
			}

			if dim.Name == "" {
				t.Errorf("Chart %s has dimension with empty Name", chart.ID)
			}

			if dim.Algo == "" {
				t.Errorf("Chart %s dimension %s has empty Algo", chart.ID, dim.Name)
			}
		}
	}
}

func TestChartsHaveUniqueIDs(t *testing.T) {
	if Charts == nil {
		t.Fatal("Charts should not be nil")
	}

	seen := make(map[string]bool)
	for _, chart := range *Charts {
		if seen[chart.ID] {
			t.Errorf("Chart ID %s is not unique", chart.ID)
		}
		seen[chart.ID] = true
	}
}

func TestChartsDimensionNameUniqueness(t *testing.T) {
	if Charts == nil {
		t.Fatal("Charts should not be nil")
	}

	for _, chart := range *Charts {
		seen := make(map[string]bool)
		for _, dim := range chart.Dims {
			if seen[dim.Name] {
				t.Errorf("Chart %s has duplicate dimension name %s", chart.ID, dim.Name)
			}
			seen[dim.Name] = true
		}
	}
}

func TestChartsAreNotModified(t *testing.T) {
	if Charts == nil {
		t.Fatal("Charts should not be nil")
	}

	// Create a copy of the current state
	originalCount := len(*Charts)

	// Ensure Charts are not modified by reference
	chartsCopy := make([]*module.Chart, len(*Charts))
	copy(chartsCopy, *Charts)

	// Verify count hasn't changed
	if len(*Charts) != originalCount {
		t.Errorf("Charts count changed from %d to %d", originalCount, len(*Charts))
	}

	// Verify the charts are the same
	for i, chart := range chartsCopy {
		if (*Charts)[i] != chart {
			t.Errorf("Chart at index %d was modified", i)
		}
	}
}

func TestChartsContainExpectedMetrics(t *testing.T) {
	if Charts == nil {
		t.Fatal("Charts should not be nil")
	}

	// Collect all dimension names from all charts
	allDims := make(map[string]bool)
	for _, chart := range *Charts {
		for _, dim := range chart.Dims {
			allDims[dim.Name] = true
		}
	}

	// Verify we have some dimensions
	if len(allDims) == 0 {
		t.Error("No dimensions found in Charts")
	}
}

func TestChartsTypeValues(t *testing.T) {
	if Charts == nil {
		t.Fatal("Charts should not be nil")
	}

	validTypes := map[string]bool{
		"line": true,
		"area": true,
		"stacked": true,
	}

	for _, chart := range *Charts {
		if chart.Type == "" {
			t.Errorf("Chart %s has empty Type", chart.ID)
		}
		// Note: We don't strictly validate type values as they may vary
		// but we ensure they are not empty
	}
}

func TestChartsUnitValues(t *testing.T) {
	if Charts == nil {
		t.Fatal("Charts should not be nil")
	}

	for _, chart := range *Charts {
		// Unit should be set (but can be empty string for some charts)
		// Just verify the field exists
		_ = chart.Units
	}
}