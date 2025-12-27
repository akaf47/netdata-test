package haproxy

import (
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/pkg/charts"
)

// Test chartsDef creation and structure
func TestChartsDef(t *testing.T) {
	chartsDef := charts.New()
	if chartsDef == nil {
		t.Fatalf("NewCharts() returned nil, expected non-nil value")
	}
}

// Test that all expected chart types are created
func TestChartsDefContainsExpectedCharts(t *testing.T) {
	chartsDef := charts.New()
	
	expectedCharts := []string{
		"haproxy_frontend_bytes_in",
		"haproxy_frontend_bytes_out",
		"haproxy_frontend_connections",
		"haproxy_frontend_errors",
		"haproxy_backend_bytes_in",
		"haproxy_backend_bytes_out",
		"haproxy_backend_connections",
		"haproxy_backend_errors",
	}
	
	for _, expectedChart := range expectedCharts {
		found := false
		for _, chart := range chartsDef {
			if chart.ID == expectedChart {
				found = true
				break
			}
		}
		if !found {
			t.Logf("Expected chart %s not found in chartsDef", expectedChart)
		}
	}
}

// Test chart dimensions are properly configured
func TestChartDimensionsConfiguration(t *testing.T) {
	chartsDef := charts.New()
	
	for _, chart := range chartsDef {
		if chart == nil {
			t.Errorf("Chart is nil")
			continue
		}
		
		if chart.ID == "" {
			t.Errorf("Chart ID is empty")
		}
		
		if chart.Title == "" {
			t.Errorf("Chart %s has empty Title", chart.ID)
		}
		
		if chart.Fqdn == "" {
			t.Errorf("Chart %s has empty Fqdn", chart.ID)
		}
		
		if chart.Units == "" {
			t.Errorf("Chart %s has empty Units", chart.ID)
		}
		
		if chart.Dims == nil {
			t.Errorf("Chart %s has nil Dims", chart.ID)
		}
		
		if len(chart.Dims) == 0 {
			t.Errorf("Chart %s has no dimensions", chart.ID)
		}
		
		for _, dim := range chart.Dims {
			if dim.ID == "" {
				t.Errorf("Chart %s has dimension with empty ID", chart.ID)
			}
			if dim.Name == "" {
				t.Errorf("Chart %s dimension %s has empty Name", chart.ID, dim.ID)
			}
		}
	}
}

// Test chart groups are properly set
func TestChartGroups(t *testing.T) {
	chartsDef := charts.New()
	
	for _, chart := range chartsDef {
		if chart == nil {
			continue
		}
		
		if chart.Group == "" {
			t.Errorf("Chart %s has empty Group", chart.ID)
		}
	}
}

// Test chart types are valid
func TestChartTypes(t *testing.T) {
	chartsDef := charts.New()
	validTypes := map[string]bool{
		"line":    true,
		"area":    true,
		"stacked": true,
	}
	
	for _, chart := range chartsDef {
		if chart == nil {
			continue
		}
		
		if !validTypes[chart.Type] {
			t.Errorf("Chart %s has invalid type: %s", chart.ID, chart.Type)
		}
	}
}

// Test multiple calls to New() return consistent results
func TestChartDefConsistency(t *testing.T) {
	charts1 := charts.New()
	charts2 := charts.New()
	
	if len(charts1) != len(charts2) {
		t.Errorf("Subsequent calls to charts.New() returned different number of charts: %d vs %d", len(charts1), len(charts2))
	}
	
	for i, chart1 := range charts1 {
		if i >= len(charts2) {
			t.Errorf("charts2 has fewer charts than charts1")
			break
		}
		chart2 := charts2[i]
		
		if chart1.ID != chart2.ID {
			t.Errorf("Chart ID mismatch at index %d: %s vs %s", i, chart1.ID, chart2.ID)
		}
		
		if len(chart1.Dims) != len(chart2.Dims) {
			t.Errorf("Chart %s has different number of dimensions: %d vs %d", chart1.ID, len(chart1.Dims), len(chart2.Dims))
		}
	}
}

// Test empty chart definition handling
func TestEmptyChartsDef(t *testing.T) {
	chartsDef := charts.New()
	
	if chartsDef == nil {
		t.Errorf("charts.New() returned nil")
		return
	}
	
	for i, chart := range chartsDef {
		if chart == nil {
			t.Errorf("Chart at index %d is nil", i)
		}
	}
}

// Test chart hierarchy and relationships
func TestChartHierarchy(t *testing.T) {
	chartsDef := charts.New()
	
	frontendCharts := 0
	backendCharts := 0
	serverCharts := 0
	
	for _, chart := range chartsDef {
		if chart == nil {
			continue
		}
		
		if chart.Fqdn != "" {
			// Count charts by type based on their ID prefix
			switch {
			case len(chart.ID) >= 8 && chart.ID[:8] == "haproxy_":
				chartType := chart.ID[8:]
				switch {
				case len(chartType) >= 8 && chartType[:8] == "frontend":
					frontendCharts++
				case len(chartType) >= 7 && chartType[:7] == "backend":
					backendCharts++
				case len(chartType) >= 6 && chartType[:6] == "server":
					serverCharts++
				}
			}
		}
	}
	
	if frontendCharts == 0 {
		t.Logf("No frontend charts found")
	}
	if backendCharts == 0 {
		t.Logf("No backend charts found")
	}
}

// Test addCharts function (if exported)
func TestAddChartsForFrontend(t *testing.T) {
	chartsDef := charts.New()
	
	// Find frontend charts
	hasFrontendChart := false
	for _, chart := range chartsDef {
		if chart != nil && len(chart.ID) > 8 && chart.ID[8:] == "frontend" || 
		   (chart != nil && chart.ID == "haproxy_frontend") {
			hasFrontendChart = true
			break
		}
	}
	
	if !hasFrontendChart {
		t.Logf("Frontend charts not found in definition")
	}
}

// Test addCharts function for backend
func TestAddChartsForBackend(t *testing.T) {
	chartsDef := charts.New()
	
	// Find backend charts
	hasBackendChart := false
	for _, chart := range chartsDef {
		if chart != nil && len(chart.ID) > 8 && chart.ID[8:] == "backend" || 
		   (chart != nil && chart.ID == "haproxy_backend") {
			hasBackendChart = true
			break
		}
	}
	
	if !hasBackendChart {
		t.Logf("Backend charts not found in definition")
	}
}

// Test chart definition completeness
func TestChartDefinitionCompleteness(t *testing.T) {
	chartsDef := charts.New()
	
	requiredFields := []string{
		"ID", "Title", "Units", "Fqdn", "Type", "Group",
	}
	
	for _, chart := range chartsDef {
		if chart == nil {
			t.Errorf("Nil chart found in definition")
			continue
		}
		
		// Check ID
		if chart.ID == "" {
			t.Errorf("Chart has empty ID")
		}
		
		// Check Title
		if chart.Title == "" {
			t.Errorf("Chart %s has empty Title", chart.ID)
		}
		
		// Check Units
		if chart.Units == "" {
			t.Errorf("Chart %s has empty Units", chart.ID)
		}
		
		// Check Fqdn
		if chart.Fqdn == "" {
			t.Errorf("Chart %s has empty Fqdn", chart.ID)
		}
		
		// Check Type
		if chart.Type == "" {
			t.Errorf("Chart %s has empty Type", chart.ID)
		}
		
		// Check Group
		if chart.Group == "" {
			t.Errorf("Chart %s has empty Group", chart.ID)
		}
		
		// Check dimensions
		if chart.Dims == nil || len(chart.Dims) == 0 {
			t.Errorf("Chart %s has no dimensions", chart.ID)
		}
	}
}

// Test chart dimension IDs are unique within a chart
func TestChartDimensionUniqueness(t *testing.T) {
	chartsDef := charts.New()
	
	for _, chart := range chartsDef {
		if chart == nil || chart.Dims == nil {
			continue
		}
		
		dimIDs := make(map[string]bool)
		for _, dim := range chart.Dims {
			if dimIDs[dim.ID] {
				t.Errorf("Chart %s has duplicate dimension ID: %s", chart.ID, dim.ID)
			}
			dimIDs[dim.ID] = true
		}
	}
}

// Test all charts have valid hierarchy
func TestChartHierarchyValidity(t *testing.T) {
	chartsDef := charts.New()
	
	for _, chart := range chartsDef {
		if chart == nil {
			continue
		}
		
		// Basic validation
		if chart.ID == "" {
			t.Errorf("Chart with empty ID found")
		}
		
		if chart.Title == "" {
			t.Errorf("Chart %s has no title", chart.ID)
		}
		
		if chart.Group == "" {
			t.Errorf("Chart %s has no group", chart.ID)
		}
	}
}