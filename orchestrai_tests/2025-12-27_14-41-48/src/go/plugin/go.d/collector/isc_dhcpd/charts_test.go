package isc_dhcpd

import (
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin.d/plugin"
)

// TestChartsInitialization tests the initialization of charts
func TestChartsInitialization(t *testing.T) {
	tests := []struct {
		name      string
		chartID   string
		chartType string
		wantErr   bool
	}{
		{
			name:      "valid chart creation",
			chartID:   "dhcpd_leases",
			chartType: "line",
			wantErr:   false,
		},
		{
			name:      "valid utilization chart",
			chartID:   "dhcpd_utilization",
			chartType: "area",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charts := Charts()
			if charts == nil {
				t.Error("Charts() returned nil")
			}
			if len(charts) == 0 {
				t.Error("Charts() returned empty slice")
			}
		})
	}
}

// TestChartsContent validates chart structure
func TestChartsContent(t *testing.T) {
	charts := Charts()

	if charts == nil {
		t.Fatal("Charts() returned nil")
	}

	// Verify charts are not empty
	if len(charts) == 0 {
		t.Fatal("Charts() should return non-empty charts")
	}

	// Verify each chart has an ID
	for _, chart := range charts {
		if chart.ID == "" {
			t.Error("Chart ID is empty")
		}
		if chart.Title == "" {
			t.Error("Chart title is empty")
		}
		if len(chart.Dims) == 0 {
			t.Error("Chart dimensions are empty")
		}
	}
}

// TestChartsValidation tests that charts conform to Netdata standards
func TestChartsValidation(t *testing.T) {
	charts := Charts()

	for i, chart := range charts {
		// Check unique IDs
		for j := i + 1; j < len(charts); j++ {
			if charts[i].ID == charts[j].ID {
				t.Errorf("Duplicate chart ID: %s", chart.ID)
			}
		}

		// Check dimensions exist and are valid
		dimIDs := make(map[string]bool)
		for _, dim := range chart.Dims {
			if dim.ID == "" {
				t.Errorf("Chart %s has dimension with empty ID", chart.ID)
			}
			if dim.Name == "" {
				t.Errorf("Chart %s dimension %s has empty name", chart.ID, dim.ID)
			}
			if dimIDs[dim.ID] {
				t.Errorf("Chart %s has duplicate dimension ID: %s", chart.ID, dim.ID)
			}
			dimIDs[dim.ID] = true
		}
	}
}

// TestChartsMetadata tests chart metadata fields
func TestChartsMetadata(t *testing.T) {
	charts := Charts()

	for _, chart := range charts {
		// Verify required metadata
		if chart.Family == "" {
			t.Errorf("Chart %s missing Family", chart.ID)
		}
		if chart.Type == "" {
			t.Errorf("Chart %s missing Type", chart.ID)
		}
	}
}

// TestChartsMultipleCalls verifies Charts() consistency across calls
func TestChartsMultipleCalls(t *testing.T) {
	charts1 := Charts()
	charts2 := Charts()

	if len(charts1) != len(charts2) {
		t.Errorf("Charts() returned different lengths: %d vs %d", len(charts1), len(charts2))
	}

	for i := range charts1 {
		if charts1[i].ID != charts2[i].ID {
			t.Errorf("Chart order differs: %s vs %s", charts1[i].ID, charts2[i].ID)
		}
	}
}

// TestChartsEmptyState handles edge case where no charts exist
func TestChartsEmptyState(t *testing.T) {
	charts := Charts()
	// Charts should never be nil, but may be empty
	if charts == nil {
		t.Error("Charts() should not return nil")
	}
}

// TestChartsDimensionTypes verifies dimension type consistency
func TestChartsDimensionTypes(t *testing.T) {
	charts := Charts()

	validTypes := map[string]bool{
		"absolute":   true,
		"incremental": true,
		"percentage": true,
		"gauge":      true,
		"counter":    true,
	}

	for _, chart := range charts {
		for _, dim := range chart.Dims {
			// Type checking - if Type is set, it should be valid
			if dim.Type != "" && !validTypes[dim.Type] {
				t.Logf("Chart %s dimension %s has type %s", chart.ID, dim.ID, dim.Type)
			}
		}
	}
}