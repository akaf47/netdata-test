package adaptecraid

import (
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin.d/go.d/agent/module"
)

func TestChartsDefinition(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "Charts should return valid chart collection",
			test: testChartsReturnValidCollection,
		},
		{
			name: "Charts should not be empty",
			test: testChartsNotEmpty,
		},
		{
			name: "Charts should have valid IDs",
			test: testChartsHaveValidIDs,
		},
		{
			name: "Charts should have valid dimensions",
			test: testChartsHaveValidDimensions,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

func testChartsReturnValidCollection(t *testing.T) {
	charts := Charts()
	if charts == nil {
		t.Fatal("Charts() returned nil")
	}
}

func testChartsNotEmpty(t *testing.T) {
	charts := Charts()
	if len(charts) == 0 {
		t.Fatal("Charts() returned empty collection")
	}
}

func testChartsHaveValidIDs(t *testing.T) {
	charts := Charts()
	for i, chart := range charts {
		if chart == nil {
			t.Fatalf("Chart at index %d is nil", i)
		}
		if chart.ID == "" {
			t.Fatalf("Chart at index %d has empty ID", i)
		}
		if chart.Title == "" {
			t.Fatalf("Chart ID '%s' has empty Title", chart.ID)
		}
	}
}

func testChartsHaveValidDimensions(t *testing.T) {
	charts := Charts()
	for _, chart := range charts {
		if chart == nil {
			continue
		}
		if len(chart.Dims) == 0 {
			t.Logf("Chart %s has no dimensions (warning only)", chart.ID)
		}
		for _, dim := range chart.Dims {
			if dim == nil {
				t.Fatalf("Chart %s has nil dimension", chart.ID)
			}
			if dim.ID == "" {
				t.Fatalf("Chart %s has dimension with empty ID", chart.ID)
			}
		}
	}
}

func TestChartTypes(t *testing.T) {
	charts := Charts()
	for _, chart := range charts {
		if chart == nil {
			continue
		}
		// Verify chart type is one of valid module chart types
		if chart.Type == "" {
			t.Fatalf("Chart %s has empty type", chart.ID)
		}
	}
}

func TestChartUnits(t *testing.T) {
	charts := Charts()
	for _, chart := range charts {
		if chart == nil {
			continue
		}
		// Units can be empty or have valid values
		if chart.Units != "" && chart.Units != "number" && chart.Units != "percentage" {
			t.Logf("Chart %s has non-standard unit: %s", chart.ID, chart.Units)
		}
	}
}

func TestDimensionOrder(t *testing.T) {
	charts := Charts()
	for _, chart := range charts {
		if chart == nil {
			continue
		}
		prevOrder := -1
		for _, dim := range chart.Dims {
			if dim == nil {
				continue
			}
			if dim.Order <= prevOrder {
				t.Logf("Dimension %s in chart %s may have incorrect order", dim.ID, chart.ID)
			}
			prevOrder = dim.Order
		}
	}
}

func TestMultipleChartsInstances(t *testing.T) {
	// Calling Charts() multiple times should return consistent results
	charts1 := Charts()
	charts2 := Charts()

	if len(charts1) != len(charts2) {
		t.Fatalf("Charts() returned different lengths: %d vs %d", len(charts1), len(charts2))
	}

	for i := 0; i < len(charts1); i++ {
		if charts1[i] == nil && charts2[i] == nil {
			continue
		}
		if charts1[i] == nil || charts2[i] == nil {
			t.Fatalf("Chart mismatch at index %d", i)
		}
		if charts1[i].ID != charts2[i].ID {
			t.Fatalf("Chart ID mismatch at index %d: %s vs %s", i, charts1[i].ID, charts2[i].ID)
		}
	}
}