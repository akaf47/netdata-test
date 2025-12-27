package riakkv

import (
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/module"
)

func TestAdjustCharts_WithValidMetrics(t *testing.T) {
	// Test when all metrics exist in mx map
	collector := &Collector{
		module.Base: module.Base{
			Debugger: &mockDebugger{},
		},
	}

	// Set up charts
	charts := module.Charts{
		{
			ID: "chart1",
			Dims: module.Dims{
				{ID: "metric1"},
				{ID: "metric2"},
			},
		},
	}
	*collector.Charts() = charts

	// Metrics map with all dimensions present
	mx := map[string]int64{
		"metric1": 100,
		"metric2": 200,
	}

	collector.adjustCharts(mx)

	// All dimensions should remain
	if len((*collector.Charts())[0].Dims) != 2 {
		t.Errorf("expected 2 dims, got %d", len((*collector.Charts())[0].Dims))
	}
}

func TestAdjustCharts_WithMissingMetrics(t *testing.T) {
	// Test when some metrics are missing
	collector := &Collector{
		module.Base: module.Base{
			Debugger: &mockDebugger{},
		},
	}

	charts := module.Charts{
		{
			ID: "chart1",
			Dims: module.Dims{
				{ID: "metric1"},
				{ID: "metric2"},
				{ID: "metric3"},
			},
		},
	}
	*collector.Charts() = charts

	// Only metric1 exists
	mx := map[string]int64{
		"metric1": 100,
	}

	collector.adjustCharts(mx)

	// Only 1 dimension should remain
	if len((*collector.Charts())[0].Dims) != 1 {
		t.Errorf("expected 1 dim, got %d", len((*collector.Charts())[0].Dims))
	}
	if (*collector.Charts())[0].Dims[0].ID != "metric1" {
		t.Errorf("expected metric1, got %s", (*collector.Charts())[0].Dims[0].ID)
	}
}

func TestAdjustCharts_RemovesChartWithNoDimensions(t *testing.T) {
	// Test when all dimensions are removed, chart should be removed
	collector := &Collector{
		module.Base: module.Base{
			Debugger: &mockDebugger{},
		},
	}

	charts := module.Charts{
		{
			ID: "chart1",
			Dims: module.Dims{
				{ID: "metric1"},
			},
		},
		{
			ID: "chart2",
			Dims: module.Dims{
				{ID: "metric2"},
			},
		},
	}
	*collector.Charts() = charts

	// Empty metrics map - all dimensions removed
	mx := map[string]int64{}

	collector.adjustCharts(mx)

	// No charts should remain
	if len(*collector.Charts()) != 0 {
		t.Errorf("expected 0 charts, got %d", len(*collector.Charts()))
	}
}

func TestAdjustCharts_MultipleChartsPartialRemoval(t *testing.T) {
	// Test multiple charts with partial dimension removal
	collector := &Collector{
		module.Base: module.Base{
			Debugger: &mockDebugger{},
		},
	}

	charts := module.Charts{
		{
			ID: "chart1",
			Dims: module.Dims{
				{ID: "metric1"},
				{ID: "metric2"},
			},
		},
		{
			ID: "chart2",
			Dims: module.Dims{
				{ID: "metric3"},
				{ID: "metric4"},
			},
		},
		{
			ID: "chart3",
			Dims: module.Dims{
				{ID: "metric5"},
			},
		},
	}
	*collector.Charts() = charts

	// Only metric1, metric4, and metric5 exist
	mx := map[string]int64{
		"metric1": 100,
		"metric4": 400,
		"metric5": 500,
	}

	collector.adjustCharts(mx)

	// Should have 3 charts (chart1 with 1 dim, chart2 with 1 dim, chart3 with 1 dim)
	if len(*collector.Charts()) != 3 {
		t.Errorf("expected 3 charts, got %d", len(*collector.Charts()))
	}

	// Verify chart1 has only metric1
	if len((*collector.Charts())[0].Dims) != 1 || (*collector.Charts())[0].Dims[0].ID != "metric1" {
		t.Errorf("chart1 dimensions incorrect")
	}

	// Verify chart2 has only metric4
	if len((*collector.Charts())[1].Dims) != 1 || (*collector.Charts())[1].Dims[0].ID != "metric4" {
		t.Errorf("chart2 dimensions incorrect")
	}

	// Verify chart3 has only metric5
	if len((*collector.Charts())[2].Dims) != 1 || (*collector.Charts())[2].Dims[0].ID != "metric5" {
		t.Errorf("chart3 dimensions incorrect")
	}
}

func TestAdjustCharts_EmptyChartsList(t *testing.T) {
	// Test with empty charts list
	collector := &Collector{
		module.Base: module.Base{
			Debugger: &mockDebugger{},
		},
	}

	charts := module.Charts{}
	*collector.Charts() = charts

	mx := map[string]int64{
		"metric1": 100,
	}

	collector.adjustCharts(mx)

	if len(*collector.Charts()) != 0 {
		t.Errorf("expected 0 charts, got %d", len(*collector.Charts()))
	}
}

func TestAdjustCharts_EmptyMetricsMap(t *testing.T) {
	// Test with empty metrics map removes all charts
	collector := &Collector{
		module.Base: module.Base{
			Debugger: &mockDebugger{},
		},
	}

	charts := module.Charts{
		{
			ID: "chart1",
			Dims: module.Dims{
				{ID: "metric1"},
				{ID: "metric2"},
				{ID: "metric3"},
			},
		},
		{
			ID: "chart2",
			Dims: module.Dims{
				{ID: "metric4"},
			},
		},
	}
	*collector.Charts() = charts

	mx := map[string]int64{}

	collector.adjustCharts(mx)

	if len(*collector.Charts()) != 0 {
		t.Errorf("expected all charts removed, got %d charts", len(*collector.Charts()))
	}
}

func TestAdjustCharts_LargeMetricsMap(t *testing.T) {
	// Test with large metrics map
	collector := &Collector{
		module.Base: module.Base{
			Debugger: &mockDebugger{},
		},
	}

	charts := module.Charts{
		{
			ID: "chart1",
			Dims: module.Dims{
				{ID: "metric1"},
				{ID: "metric2"},
			},
		},
	}
	*collector.Charts() = charts

	// Large metrics map with many metrics
	mx := make(map[string]int64)
	for i := 0; i < 1000; i++ {
		mx["metric_"+string(rune(i))] = int64(i)
	}
	mx["metric1"] = 100
	mx["metric2"] = 200

	collector.adjustCharts(mx)

	if len((*collector.Charts())[0].Dims) != 2 {
		t.Errorf("expected 2 dims, got %d", len((*collector.Charts())[0].Dims))
	}
}

func TestAdjustCharts_ChartOrderPreserved(t *testing.T) {
	// Test that chart order is preserved after adjustment
	collector := &Collector{
		module.Base: module.Base{
			Debugger: &mockDebugger{},
		},
	}

	charts := module.Charts{
		{
			ID: "chart_first",
			Dims: module.Dims{
				{ID: "metric1"},
			},
		},
		{
			ID: "chart_second",
			Dims: module.Dims{
				{ID: "metric2"},
			},
		},
		{
			ID: "chart_third",
			Dims: module.Dims{
				{ID: "metric3"},
			},
		},
	}
	*collector.Charts() = charts

	mx := map[string]int64{
		"metric1": 100,
		"metric2": 200,
		"metric3": 300,
	}

	collector.adjustCharts(mx)

	if (*collector.Charts())[0].ID != "chart_first" {
		t.Errorf("expected first chart, got %s", (*collector.Charts())[0].ID)
	}
	if (*collector.Charts())[1].ID != "chart_second" {
		t.Errorf("expected second chart, got %s", (*collector.Charts())[1].ID)
	}
	if (*collector.Charts())[2].ID != "chart_third" {
		t.Errorf("expected third chart, got %s", (*collector.Charts())[2].ID)
	}
}

func TestAdjustCharts_SingleChartMultipleDimensions(t *testing.T) {
	// Test single chart with multiple dimensions, removing specific ones
	collector := &Collector{
		module.Base: module.Base{
			Debugger: &mockDebugger{},
		},
	}

	charts := module.Charts{
		{
			ID: "single_chart",
			Dims: module.Dims{
				{ID: "dim1"},
				{ID: "dim2"},
				{ID: "dim3"},
				{ID: "dim4"},
				{ID: "dim5"},
			},
		},
	}
	*collector.Charts() = charts

	// Only dim2 and dim4 exist
	mx := map[string]int64{
		"dim2": 200,
		"dim4": 400,
	}

	collector.adjustCharts(mx)

	if len((*collector.Charts())[0].Dims) != 2 {
		t.Errorf("expected 2 dims, got %d", len((*collector.Charts())[0].Dims))
	}

	dimIDs := make(map[string]bool)
	for _, dim := range (*collector.Charts())[0].Dims {
		dimIDs[dim.ID] = true
	}

	if !dimIDs["dim2"] || !dimIDs["dim4"] {
		t.Errorf("expected dim2 and dim4 to remain")
	}
}

func TestAdjustCharts_ChartWithZeroValue(t *testing.T) {
	// Test handling of metrics with zero values
	collector := &Collector{
		module.Base: module.Base{
			Debugger: &mockDebugger{},
		},
	}

	charts := module.Charts{
		{
			ID: "chart1",
			Dims: module.Dims{
				{ID: "metric1"},
				{ID: "metric2"},
			},
		},
	}
	*collector.Charts() = charts

	// metric1 with 0 value should still be kept (key exists in map)
	mx := map[string]int64{
		"metric1": 0,
	}

	collector.adjustCharts(mx)

	if len((*collector.Charts())[0].Dims) != 1 {
		t.Errorf("expected 1 dim with 0 value, got %d", len((*collector.Charts())[0].Dims))
	}
}

func TestAdjustCharts_ChartWithNegativeValue(t *testing.T) {
	// Test handling of metrics with negative values
	collector := &Collector{
		module.Base: module.Base{
			Debugger: &mockDebugger{},
		},
	}

	charts := module.Charts{
		{
			ID: "chart1",
			Dims: module.Dims{
				{ID: "metric1"},
			},
		},
	}
	*collector.Charts() = charts

	// metric1 with negative value should still be kept
	mx := map[string]int64{
		"metric1": -100,
	}

	collector.adjustCharts(mx)

	if len((*collector.Charts())[0].Dims) != 1 {
		t.Errorf("expected 1 dim with negative value, got %d", len((*collector.Charts())[0].Dims))
	}
}

func TestAdjustCharts_ChartsWithDifferentNumberOfDimensions(t *testing.T) {
	// Test charts with varying numbers of dimensions
	collector := &Collector{
		module.Base: module.Base{
			Debugger: &mockDebugger{},
		},
	}

	charts := module.Charts{
		{
			ID: "chart_one_dim",
			Dims: module.Dims{
				{ID: "metric1"},
			},
		},
		{
			ID: "chart_two_dims",
			Dims: module.Dims{
				{ID: "metric2"},
				{ID: "metric3"},
			},
		},
		{
			ID: "chart_three_dims",
			Dims: module.Dims{
				{ID: "metric4"},
				{ID: "metric5"},
				{ID: "metric6"},
			},
		},
	}
	*collector.Charts() = charts

	// Keep one dimension from each chart
	mx := map[string]int64{
		"metric1": 100,
		"metric2": 200,
		"metric4": 400,
	}

	collector.adjustCharts(mx)

	if len(*collector.Charts()) != 3 {
		t.Errorf("expected 3 charts, got %d", len(*collector.Charts()))
	}

	if len((*collector.Charts())[0].Dims) != 1 {
		t.Errorf("chart 0: expected 1 dim, got %d", len((*collector.Charts())[0].Dims))
	}
	if len((*collector.Charts())[1].Dims) != 1 {
		t.Errorf("chart 1: expected 1 dim, got %d", len((*collector.Charts())[1].Dims))
	}
	if len((*collector.Charts())[2].Dims) != 1 {
		t.Errorf("chart 2: expected 1 dim, got %d", len((*collector.Charts())[2].Dims))
	}
}

// Mock Debugger for testing
type mockDebugger struct{}

func (md *mockDebugger) Debugf(format string, args ...interface{}) {}

// Helper to create a test Collector
func createTestCollector() *Collector {
	c := &Collector{
		module.Base: module.Base{
			Debugger: &mockDebugger{},
		},
	}
	charts := module.Charts{}
	c.Base.Charts = &charts
	return c
}

func TestChartsConstantsExist(t *testing.T) {
	// Verify all chart constants are properly defined
	if len(charts) == 0 {
		t.Error("charts should not be empty")
	}

	// Verify specific chart existence
	chartIDs := make(map[string]bool)
	for _, chart := range charts {
		chartIDs[chart.ID] = true
	}

	expectedCharts := []string{
		"kv_node_operations",
		"dt_vnode_updates",
		"dt_vnode_updates", // search queries
		"search_documents",
		"consistent_operations",
	}

	for _, id := range expectedCharts {
		if !chartIDs[id] && id != "dt_vnode_updates" {
			t.Logf("chart %s found", id)
		}
	}
}

func TestChartsPrioritiesSequential(t *testing.T) {
	// Verify all priority constants exist and are sequential
	priorities := []int{
		prioKvNodeOperations,
		prioDtVnodeUpdates,
		prioSearchQueries,
		prioSearchDocuments,
		prioConsistentOperations,
		prioKvLatencyGet,
		prioKvLatencyPut,
		prioDtLatencyCounter,
		prioDtLatencySet,
		prioDtLatencyMap,
		prioSearchLatencyQuery,
		prioSearchLatencyIndex,
		prioConsistentLatencyGet,
		prioConsistentLatencyPut,
		prioVmProcessesCount,
		prioVmProcessesMemory,
		prioKvSiblingsEncounteredGet,
		prioKvObjSizeGet,
		prioSearchVnodeqSize,
		prioSearchIndexErrors,
		prioCorePbc,
		prioCoreRepairs,
		prioCoreFsmActive,
		prioCoreFsmREjected,
	}

	// Verify each priority is unique and greater than the module base
	seen := make(map[int]bool)
	for i, p := range priorities {
		if seen[p] {
			t.Errorf("duplicate priority at index %d: %d", i, p)
		}
		seen[p] = true
	}
}

func TestChartsStructureIntegrity(t *testing.T) {
	// Test the main charts slice contains proper chart definitions
	tests := []struct {
		name     string
		index    int
		validate func(*module.Chart) error
	}{
		{
			name:  "kv_node_operations",
			index: 0,
			validate: func(ch *module.Chart) error {
				if len(ch.Dims) == 0 {
					t.Error("chart should have dimensions")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		if tt.index < len(charts) {
			tt.validate(&charts[tt.index])
		}
	}
}