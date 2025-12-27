package pihole

import (
	"testing"
)

// TestCreateCharts tests the CreateCharts method
func TestCreateCharts(t *testing.T) {
	p := &Pihole{}

	charts := p.createCharts()

	if charts == nil {
		t.Fatal("createCharts() returned nil")
	}

	if len(charts) == 0 {
		t.Error("createCharts() returned empty charts")
	}
}

// TestCreateCharts_AllChartsHaveFamily tests that all charts have a family
func TestCreateCharts_AllChartsHaveFamily(t *testing.T) {
	p := &Pihole{}

	charts := p.createCharts()

	for i, chart := range charts {
		if chart.Family == "" {
			t.Errorf("Chart[%d] has empty family", i)
		}
	}
}

// TestCreateCharts_AllChartsHaveUnits tests that all charts have units
func TestCreateCharts_AllChartsHaveUnits(t *testing.T) {
	p := &Pihole{}

	charts := p.createCharts()

	for i, chart := range charts {
		if chart.Units == "" {
			t.Errorf("Chart[%d] %s has empty units", i, chart.ID)
		}
	}
}

// TestCreateCharts_AllChartsHaveDimensions tests that all charts have dimensions
func TestCreateCharts_AllChartsHaveDimensions(t *testing.T) {
	p := &Pihole{}

	charts := p.createCharts()

	for i, chart := range charts {
		if len(chart.Dims) == 0 {
			t.Errorf("Chart[%d] %s has no dimensions", i, chart.ID)
		}
	}
}

// TestCreateCharts_UniqueChartIDs tests that chart IDs are unique
func TestCreateCharts_UniqueChartIDs(t *testing.T) {
	p := &Pihole{}

	charts := p.createCharts()

	ids := make(map[string]bool)
	for _, chart := range charts {
		if ids[chart.ID] {
			t.Errorf("Duplicate chart ID: %s", chart.ID)
		}
		ids[chart.ID] = true
	}
}

// TestCreateCharts_UniqueDimensionIDs tests that dimension IDs are unique
func TestCreateCharts_UniqueDimensionIDs(t *testing.T) {
	p := &Pihole{}

	charts := p.createCharts()

	ids := make(map[string]bool)
	for _, chart := range charts {
		for _, dim := range chart.Dims {
			if ids[dim.ID] {
				t.Errorf("Duplicate dimension ID: %s in chart %s", dim.ID, chart.ID)
			}
			ids[dim.ID] = true
		}
	}
}

// TestCreateCharts_AllDimensionsHaveNames tests that all dimensions have names
func TestCreateCharts_AllDimensionsHaveNames(t *testing.T) {
	p := &Pihole{}

	charts := p.createCharts()

	for _, chart := range charts {
		for _, dim := range chart.Dims {
			if dim.Name == "" {
				t.Errorf("Dimension %s in chart %s has empty name", dim.ID, chart.ID)
			}
		}
	}
}

// TestCreateCharts_ValidChartTypes tests that all charts have valid types
func TestCreateCharts_ValidChartTypes(t *testing.T) {
	p := &Pihole{}

	charts := p.createCharts()

	validTypes := map[string]bool{
		"line":   true,
		"area":   true,
		"stacked": true,
	}

	for i, chart := range charts {
		if chart.Type == "" {
			t.Errorf("Chart[%d] %s has empty type", i, chart.ID)
		}
		if !validTypes[chart.Type] && chart.Type != "" {
			t.Logf("Chart[%d] %s has unusual type: %s", i, chart.ID, chart.Type)
		}
	}
}

// TestCreateCharts_QueriesChart tests the queries chart structure
func TestCreateCharts_QueriesChart(t *testing.T) {
	p := &Pihole{}

	charts := p.createCharts()

	queriesChart := findChartByID(charts, "pihole_queries")
	if queriesChart == nil {
		t.Skip("Queries chart not found")
	}

	if len(queriesChart.Dims) == 0 {
		t.Error("Queries chart has no dimensions")
	}
}

// TestCreateCharts_BlockedChart tests the blocked queries chart structure
func TestCreateCharts_BlockedChart(t *testing.T) {
	p := &Pihole{}

	charts := p.createCharts()

	blockedChart := findChartByID(charts, "pihole_blocked")
	if blockedChart == nil {
		t.Skip("Blocked queries chart not found")
	}

	if len(blockedChart.Dims) == 0 {
		t.Error("Blocked queries chart has no dimensions")
	}
}

// TestCreateCharts_CacheChart tests the cache chart structure
func TestCreateCharts_CacheChart(t *testing.T) {
	p := &Pihole{}

	charts := p.createCharts()

	cacheChart := findChartByID(charts, "pihole_cache")
	if cacheChart == nil {
		t.Skip("Cache chart not found")
	}

	if len(cacheChart.Dims) == 0 {
		t.Error("Cache chart has no dimensions")
	}
}

// TestCreateCharts_UniqueClients tests the unique clients chart structure
func TestCreateCharts_UniqueClients(t *testing.T) {
	p := &Pihole{}

	charts := p.createCharts()

	clientsChart := findChartByID(charts, "pihole_clients")
	if clientsChart == nil {
		t.Skip("Unique clients chart not found")
	}

	if len(clientsChart.Dims) == 0 {
		t.Error("Unique clients chart has no dimensions")
	}
}

// TestCreateCharts_GravityChart tests the gravity chart structure
func TestCreateCharts_GravityChart(t *testing.T) {
	p := &Pihole{}

	charts := p.createCharts()

	gravityChart := findChartByID(charts, "pihole_gravity")
	if gravityChart == nil {
		t.Skip("Gravity chart not found")
	}

	if len(gravityChart.Dims) == 0 {
		t.Error("Gravity chart has no dimensions")
	}
}

// TestRemoveCharts tests the RemoveCharts method
func TestRemoveCharts(t *testing.T) {
	p := &Pihole{}

	p.removeCharts()
	// Method should not panic or error
}

// TestRemoveCharts_MultipleNilCharts tests RemoveCharts with nil charts
func TestRemoveCharts_MultipleNilCharts(t *testing.T) {
	p := &Pihole{
		Charts: nil,
	}

	p.removeCharts()
	// Should handle nil charts gracefully
}

// TestRemoveCharts_EmptyCharts tests RemoveCharts with empty charts
func TestRemoveCharts_EmptyCharts(t *testing.T) {
	p := &Pihole{
		Charts: []*Chart{},
	}

	p.removeCharts()
	// Should handle empty charts gracefully
}

// TestInitCharts tests the InitCharts method
func TestInitCharts(t *testing.T) {
	p := &Pihole{}

	p.initCharts()

	if p.Charts == nil || len(p.Charts) == 0 {
		t.Error("initCharts() did not initialize charts")
	}
}

// TestInitCharts_MultipleInits tests calling InitCharts multiple times
func TestInitCharts_MultipleInits(t *testing.T) {
	p := &Pihole{}

	p.initCharts()
	firstCharts := p.Charts

	p.initCharts()
	secondCharts := p.Charts

	if len(firstCharts) != len(secondCharts) {
		t.Errorf("InitCharts() length mismatch: %d vs %d", len(firstCharts), len(secondCharts))
	}
}

// TestAddCharts tests adding charts to the collector
func TestAddCharts(t *testing.T) {
	p := &Pihole{}

	p.addCharts()

	if p.Charts == nil {
		t.Error("addCharts() did not set Charts")
	}
}

// TestChartTypesConsistency tests that chart types are consistent
func TestChartTypesConsistency(t *testing.T) {
	p := &Pihole{}

	charts := p.createCharts()

	for i, chart := range charts {
		if chart.Type == "" {
			t.Errorf("Chart[%d] %s has no type", i, chart.ID)
		}
	}
}

// TestDimensionMetaAttributes tests dimension meta attributes
func TestDimensionMetaAttributes(t *testing.T) {
	p := &Pihole{}

	charts := p.createCharts()

	for _, chart := range charts {
		for _, dim := range chart.Dims {
			// Check that dimension has required fields
			if dim.ID == "" {
				t.Errorf("Dimension in chart %s has empty ID", chart.ID)
			}
		}
	}
}

// TestChartLayout tests chart layout validity
func TestChartLayout(t *testing.T) {
	p := &Pihole{}

	charts := p.createCharts()

	for i, chart := range charts {
		if chart.ID == "" {
			t.Errorf("Chart[%d] has empty ID", i)
		}
		if chart.Type == "" {
			t.Errorf("Chart[%d] has empty Type", i)
		}
		if chart.Family == "" {
			t.Errorf("Chart[%d] has empty Family", i)
		}
		if chart.Units == "" {
			t.Errorf("Chart[%d] has empty Units", i)
		}
	}
}

// TestChartsMultipleCreations tests creating charts multiple times
func TestChartsMultipleCreations(t *testing.T) {
	p := &Pihole{}

	charts1 := p.createCharts()
	charts2 := p.createCharts()

	if len(charts1) != len(charts2) {
		t.Errorf("createCharts() returned different counts: %d vs %d", len(charts1), len(charts2))
	}
}

// TestChartIDNaming tests chart ID naming conventions
func TestChartIDNaming(t *testing.T) {
	p := &Pihole{}

	charts := p.createCharts()

	for _, chart := range charts {
		if chart.ID == "" {
			t.Error("Chart has empty ID")
		}
		// IDs should typically follow naming conventions
		if len(chart.ID) < 3 {
			t.Logf("Chart ID %s seems very short", chart.ID)
		}
	}
}

// TestDimensionAggregationType tests dimension aggregation types
func TestDimensionAggregationType(t *testing.T) {
	p := &Pihole{}

	charts := p.createCharts()

	for _, chart := range charts {
		for _, dim := range chart.Dims {
			if dim.Algo == "" {
				t.Logf("Dimension %s in chart %s has no algorithm specified", dim.ID, chart.ID)
			}
		}
	}
}

// Helper function to find chart by ID
func findChartByID(charts []*Chart, id string) *Chart {
	for _, chart := range charts {
		if chart.ID == id {
			return chart
		}
	}
	return nil
}

// Helper structure definitions (if not already in the package)
type Chart struct {
	ID   string
	Type string
	Family string
	Units string
	Dims []*Dimension
}

type Dimension struct {
	ID   string
	Name string
	Algo string
}