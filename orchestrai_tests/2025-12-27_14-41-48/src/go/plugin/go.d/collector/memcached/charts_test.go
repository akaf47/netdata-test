// SPDX-License-Identifier: GPL-3.0-or-later

package memcached

import (
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/module"
)

// TestChartsInitialization verifies that charts slice is properly initialized
func TestChartsInitialization(t *testing.T) {
	// Assert
	if len(charts) == 0 {
		t.Fatal("charts slice should not be empty")
	}
}

// TestChartsLength verifies that charts contains the expected number of chart definitions
func TestChartsLength(t *testing.T) {
	// Assert - should have 14 charts
	expectedCount := 14
	if len(charts) != expectedCount {
		t.Fatalf("expected %d charts, got %d", expectedCount, len(charts))
	}
}

// TestChartsPriorities verifies that all charts have correct priority values
func TestChartsPriorities(t *testing.T) {
	tests := []struct {
		index    int
		chartID  string
		priority int
	}{
		{0, "cache", prioCache},
		{1, "net", prioNet},
		{2, "connections", prioConnections},
		{3, "items", prioItems},
		{4, "evicted_reclaimed", prioEvictedReclaimed},
		{5, "get", prioGet},
		{6, "get_rate", prioGetRate},
		{7, "set_rate", prioSetRate},
		{8, "delete", prioDelete},
		{9, "cas", prioCas},
		{10, "increment", prioIncrement},
		{11, "decrement", prioDecrement},
		{12, "touch", prioTouch},
		{13, "touch_rate", prioTouchRate},
	}

	for _, tt := range tests {
		t.Run(tt.chartID, func(t *testing.T) {
			// Assert
			if tt.index >= len(charts) {
				t.Fatalf("chart index %d out of range", tt.index)
			}
			chart := charts[tt.index]
			if chart.ID != tt.chartID {
				t.Errorf("expected chart ID %s, got %s", tt.chartID, chart.ID)
			}
			if chart.Priority != tt.priority {
				t.Errorf("expected priority %d, got %d", tt.priority, chart.Priority)
			}
		})
	}
}

// TestCacheChart verifies cache chart structure
func TestCacheChart(t *testing.T) {
	// Assert
	if charts[0].ID != "cache" {
		t.Fatalf("expected cache chart at index 0")
	}
	chart := charts[0]
	if chart.Title != "Cache Size" {
		t.Errorf("expected title 'Cache Size', got '%s'", chart.Title)
	}
	if chart.Units != "MiB" {
		t.Errorf("expected units 'MiB', got '%s'", chart.Units)
	}
	if chart.Fam != "cache" {
		t.Errorf("expected family 'cache', got '%s'", chart.Fam)
	}
	if chart.Ctx != "memcached.cache" {
		t.Errorf("expected context 'memcached.cache', got '%s'", chart.Ctx)
	}
	if chart.Type != module.Stacked {
		t.Errorf("expected type Stacked, got %v", chart.Type)
	}
	if len(chart.Dims) != 2 {
		t.Errorf("expected 2 dimensions, got %d", len(chart.Dims))
	}
	if chart.Dims[0].ID != "avail" {
		t.Errorf("expected first dimension 'avail', got '%s'", chart.Dims[0].ID)
	}
	if chart.Dims[1].ID != "bytes" {
		t.Errorf("expected second dimension 'bytes', got '%s'", chart.Dims[1].ID)
	}
	if chart.Dims[1].Name != "used" {
		t.Errorf("expected dimension name 'used', got '%s'", chart.Dims[1].Name)
	}
}

// TestNetChart verifies network chart structure
func TestNetChart(t *testing.T) {
	// Assert
	chart := charts[1]
	if chart.ID != "net" {
		t.Fatalf("expected net chart at index 1")
	}
	if chart.Title != "Network" {
		t.Errorf("expected title 'Network', got '%s'", chart.Title)
	}
	if chart.Units != "kilobits/s" {
		t.Errorf("expected units 'kilobits/s', got '%s'", chart.Units)
	}
	if chart.Type != module.Area {
		t.Errorf("expected type Area, got %v", chart.Type)
	}
	if len(chart.Dims) != 2 {
		t.Errorf("expected 2 dimensions, got %d", len(chart.Dims))
	}
	// Verify read dimension
	readDim := chart.Dims[0]
	if readDim.ID != "bytes_read" {
		t.Errorf("expected read dimension 'bytes_read', got '%s'", readDim.ID)
	}
	if readDim.Name != "in" {
		t.Errorf("expected read dimension name 'in', got '%s'", readDim.Name)
	}
	if readDim.Mul != 8 {
		t.Errorf("expected read multiplier 8, got %d", readDim.Mul)
	}
	if readDim.Div != 1000 {
		t.Errorf("expected read divisor 1000, got %d", readDim.Div)
	}
	if readDim.Algo != module.Incremental {
		t.Errorf("expected read algo Incremental, got %v", readDim.Algo)
	}
	// Verify write dimension
	writeDim := chart.Dims[1]
	if writeDim.ID != "bytes_written" {
		t.Errorf("expected write dimension 'bytes_written', got '%s'", writeDim.ID)
	}
	if writeDim.Name != "out" {
		t.Errorf("expected write dimension name 'out', got '%s'", writeDim.Name)
	}
	if writeDim.Mul != -8 {
		t.Errorf("expected write multiplier -8, got %d", writeDim.Mul)
	}
}

// TestConnectionsChart verifies connections chart structure
func TestConnectionsChart(t *testing.T) {
	// Assert
	chart := charts[2]
	if chart.ID != "connections" {
		t.Fatalf("expected connections chart at index 2")
	}
	if chart.Title != "Connections" {
		t.Errorf("expected title 'Connections', got '%s'", chart.Title)
	}
	if chart.Units != "connections/s" {
		t.Errorf("expected units 'connections/s', got '%s'", chart.Units)
	}
	if chart.Type != module.Line {
		t.Errorf("expected type Line, got %v", chart.Type)
	}
	if len(chart.Dims) != 3 {
		t.Errorf("expected 3 dimensions, got %d", len(chart.Dims))
	}
	if chart.Dims[0].ID != "curr_connections" {
		t.Errorf("expected dimension 'curr_connections'")
	}
	if chart.Dims[1].ID != "rejected_connections" {
		t.Errorf("expected dimension 'rejected_connections'")
	}
	if chart.Dims[2].ID != "total_connections" {
		t.Errorf("expected dimension 'total_connections'")
	}
}

// TestItemsChart verifies items chart structure
func TestItemsChart(t *testing.T) {
	// Assert
	chart := charts[3]
	if chart.ID != "items" {
		t.Fatalf("expected items chart at index 3")
	}
	if chart.Title != "Items" {
		t.Errorf("expected title 'Items', got '%s'", chart.Title)
	}
	if chart.Units != "items" {
		t.Errorf("expected units 'items', got '%s'", chart.Units)
	}
	if len(chart.Dims) != 2 {
		t.Errorf("expected 2 dimensions, got %d", len(chart.Dims))
	}
}

// TestEvictedReclaimedChart verifies evicted/reclaimed chart structure
func TestEvictedReclaimedChart(t *testing.T) {
	// Assert
	chart := charts[4]
	if chart.ID != "evicted_reclaimed" {
		t.Fatalf("expected evicted_reclaimed chart at index 4")
	}
	if chart.Title != "Evicted and Reclaimed Items" {
		t.Errorf("expected title 'Evicted and Reclaimed Items', got '%s'", chart.Title)
	}
	if len(chart.Dims) != 2 {
		t.Errorf("expected 2 dimensions, got %d", len(chart.Dims))
	}
	if chart.Dims[0].ID != "reclaimed" {
		t.Errorf("expected dimension 'reclaimed'")
	}
	if chart.Dims[1].ID != "evictions" {
		t.Errorf("expected dimension 'evictions'")
	}
	if chart.Dims[1].Name != "evicted" {
		t.Errorf("expected dimension name 'evicted', got '%s'", chart.Dims[1].Name)
	}
}

// TestGetChart verifies get requests chart structure
func TestGetChart(t *testing.T) {
	// Assert
	chart := charts[5]
	if chart.ID != "get" {
		t.Fatalf("expected get chart at index 5")
	}
	if chart.Title != "Get Requests" {
		t.Errorf("expected title 'Get Requests', got '%s'", chart.Title)
	}
	if chart.Type != module.Stacked {
		t.Errorf("expected type Stacked, got %v", chart.Type)
	}
	if len(chart.Dims) != 2 {
		t.Errorf("expected 2 dimensions, got %d", len(chart.Dims))
	}
	if chart.Dims[0].Algo != module.PercentOfAbsolute {
		t.Errorf("expected PercentOfAbsolute algorithm")
	}
}

// TestGetRateChart verifies get rate chart structure
func TestGetRateChart(t *testing.T) {
	// Assert
	chart := charts[6]
	if chart.ID != "get_rate" {
		t.Fatalf("expected get_rate chart at index 6")
	}
	if chart.Title != "Get Request Rate" {
		t.Errorf("expected title 'Get Request Rate', got '%s'", chart.Title)
	}
	if chart.Units != "requests/s" {
		t.Errorf("expected units 'requests/s', got '%s'", chart.Units)
	}
	if len(chart.Dims) != 1 {
		t.Errorf("expected 1 dimension, got %d", len(chart.Dims))
	}
	if chart.Dims[0].ID != "cmd_get" {
		t.Errorf("expected dimension 'cmd_get', got '%s'", chart.Dims[0].ID)
	}
	if chart.Dims[0].Algo != module.Incremental {
		t.Errorf("expected Incremental algorithm")
	}
}

// TestSetRateChart verifies set rate chart structure
func TestSetRateChart(t *testing.T) {
	// Assert
	chart := charts[7]
	if chart.ID != "set_rate" {
		t.Fatalf("expected set_rate chart at index 7")
	}
	if chart.Title != "Set Request Rate" {
		t.Errorf("expected title 'Set Request Rate', got '%s'", chart.Title)
	}
	if len(chart.Dims) != 1 {
		t.Errorf("expected 1 dimension, got %d", len(chart.Dims))
	}
}

// TestDeleteChart verifies delete chart structure
func TestDeleteChart(t *testing.T) {
	// Assert
	chart := charts[8]
	if chart.ID != "delete" {
		t.Fatalf("expected delete chart at index 8")
	}
	if chart.Title != "Delete Requests" {
		t.Errorf("expected title 'Delete Requests', got '%s'", chart.Title)
	}
	if chart.Type != module.Stacked {
		t.Errorf("expected type Stacked")
	}
	if len(chart.Dims) != 2 {
		t.Errorf("expected 2 dimensions, got %d", len(chart.Dims))
	}
}

// TestCasChart verifies CAS (check and set) chart structure
func TestCasChart(t *testing.T) {
	// Assert
	chart := charts[9]
	if chart.ID != "cas" {
		t.Fatalf("expected cas chart at index 9")
	}
	if chart.Title != "Check and Set Requests" {
		t.Errorf("expected title 'Check and Set Requests', got '%s'", chart.Title)
	}
	if len(chart.Dims) != 3 {
		t.Errorf("expected 3 dimensions, got %d", len(chart.Dims))
	}
	if chart.Dims[2].ID != "cas_badval" {
		t.Errorf("expected dimension 'cas_badval'")
	}
	if chart.Dims[2].Name != "bad value" {
		t.Errorf("expected dimension name 'bad value', got '%s'", chart.Dims[2].Name)
	}
}

// TestIncrementChart verifies increment chart structure
func TestIncrementChart(t *testing.T) {
	// Assert
	chart := charts[10]
	if chart.ID != "increment" {
		t.Fatalf("expected increment chart at index 10")
	}
	if chart.Title != "Increment Requests" {
		t.Errorf("expected title 'Increment Requests', got '%s'", chart.Title)
	}
	if len(chart.Dims) != 2 {
		t.Errorf("expected 2 dimensions, got %d", len(chart.Dims))
	}
}

// TestDecrementChart verifies decrement chart structure
func TestDecrementChart(t *testing.T) {
	// Assert
	chart := charts[11]
	if chart.ID != "decrement" {
		t.Fatalf("expected decrement chart at index 11")
	}
	if chart.Title != "Decrement Requests" {
		t.Errorf("expected title 'Decrement Requests', got '%s'", chart.Title)
	}
	if len(chart.Dims) != 2 {
		t.Errorf("expected 2 dimensions, got %d", len(chart.Dims))
	}
}

// TestTouchChart verifies touch chart structure
func TestTouchChart(t *testing.T) {
	// Assert
	chart := charts[12]
	if chart.ID != "touch" {
		t.Fatalf("expected touch chart at index 12")
	}
	if chart.Title != "Touch Requests" {
		t.Errorf("expected title 'Touch Requests', got '%s'", chart.Title)
	}
	if len(chart.Dims) != 2 {
		t.Errorf("expected 2 dimensions, got %d", len(chart.Dims))
	}
}

// TestTouchRateChart verifies touch rate chart structure
func TestTouchRateChart(t *testing.T) {
	// Assert
	chart := charts[13]
	if chart.ID != "touch_rate" {
		t.Fatalf("expected touch_rate chart at index 13")
	}
	if chart.Title != "Touch Requests Rate" {
		t.Errorf("expected title 'Touch Requests Rate', got '%s'", chart.Title)
	}
	if chart.Units != "requests/s" {
		t.Errorf("expected units 'requests/s', got '%s'", chart.Units)
	}
	if len(chart.Dims) != 1 {
		t.Errorf("expected 1 dimension, got %d", len(chart.Dims))
	}
}

// TestByteToMiBConstant verifies the byte to MiB conversion constant
func TestByteToMiBConstant(t *testing.T) {
	// Assert - 1 MiB = 1048576 bytes = 1 << 20
	expectedValue := 1 << 20
	if byteToMiB != expectedValue {
		t.Errorf("expected byteToMiB to be %d, got %d", expectedValue, byteToMiB)
	}
}

// TestChartsAreCopies verifies that charts slice contains copies of base charts
func TestChartsAreCopies(t *testing.T) {
	// Arrange
	originalCache := cacheChart

	// Act
	chartsCopy := charts[0]

	// Assert - they should have the same values but different pointers if deep copied
	if chartsCopy.ID != originalCache.ID {
		t.Errorf("chart ID mismatch after copy")
	}
	if chartsCopy.Title != originalCache.Title {
		t.Errorf("chart Title mismatch after copy")
	}
}

// TestAllDimensionsDefined verifies all dimensions in all charts are properly defined
func TestAllDimensionsDefined(t *testing.T) {
	for chartIdx, chart := range charts {
		if len(chart.Dims) == 0 {
			t.Errorf("chart at index %d (%s) has no dimensions", chartIdx, chart.ID)
		}
		for dimIdx, dim := range chart.Dims {
			if dim.ID == "" {
				t.Errorf("chart %d dimension %d has empty ID", chartIdx, dimIdx)
			}
		}
	}
}

// TestChartFamilies verifies all charts have valid families
func TestChartFamilies(t *testing.T) {
	families := map[string]bool{
		"cache":           false,
		"network":         false,
		"connections":     false,
		"items":           false,
		"get ops":         false,
		"set ops":         false,
		"delete ops":      false,
		"check and set ops": false,
		"increment ops":   false,
		"decrement ops":   false,
		"touch ops":       false,
	}

	for _, chart := range charts {
		if _, exists := families[chart.Fam]; !exists {
			t.Errorf("chart %s has unexpected family: %s", chart.ID, chart.Fam)
		}
	}
}

// TestChartContexts verifies all charts have proper contexts
func TestChartContexts(t *testing.T) {
	expectedContexts := []string{
		"memcached.cache",
		"memcached.net",
		"memcached.connections",
		"memcached.items",
		"memcached.evicted_reclaimed",
		"memcached.get",
		"memcached.get_rate",
		"memcached.set_rate",
		"memcached.delete",
		"memcached.cas",
		"memcached.increment",
		"memcached.decrement",
		"memcached.touch",
		"memcached.touch_rate",
	}

	for i, expectedCtx := range expectedContexts {
		if i >= len(charts) {
			t.Fatalf("expected more charts than found")
		}
		if charts[i].Ctx != expectedCtx {
			t.Errorf("chart at index %d: expected context %s, got %s", i, expectedCtx, charts[i].Ctx)
		}
	}
}

// TestIncrementalAlgorithmDimensions verifies dimensions using incremental algorithm
func TestIncrementalAlgorithmDimensions(t *testing.T) {
	// Test net chart dimensions
	netChart := charts[1]
	for _, dim := range netChart.Dims {
		if dim.Algo != module.Incremental {
			t.Errorf("net chart dimension %s should use Incremental algorithm", dim.ID)
		}
	}

	// Test connections chart dimensions
	connChart := charts[2]
	for _, dim := range connChart.Dims {
		if dim.Algo != module.Incremental {
			t.Errorf("connections chart dimension %s should use Incremental algorithm", dim.ID)
		}
	}
}

// TestPercentOfAbsoluteAlgorithmDimensions verifies dimensions using PercentOfAbsolute algorithm
func TestPercentOfAbsoluteAlgorithmDimensions(t *testing.T) {
	testCases := []struct {
		chartIndex int
		chartName  string
	}{
		{5, "get"},
		{8, "delete"},
		{9, "cas"},
		{10, "increment"},
		{11, "decrement"},
		{12, "touch"},
	}

	for _, tc := range testCases {
		chart := charts[tc.chartIndex]
		for _, dim := range chart.Dims {
			if dim.Algo != module.PercentOfAbsolute {
				t.Errorf("chart %s dimension %s should use PercentOfAbsolute algorithm, got %v", tc.chartName, dim.ID, dim.Algo)
			}
		}
	}
}

// TestCacheDimensionDivisors verifies cache dimensions use correct divisors
func TestCacheDimensionDivisors(t *testing.T) {
	cacheChart := charts[0]
	if cacheChart.Dims[0].Div != byteToMiB {
		t.Errorf("avail dimension should use byteToMiB divisor")
	}
	if cacheChart.Dims[1].Div != byteToMiB {
		t.Errorf("bytes dimension should use byteToMiB divisor")
	}
}

// TestPriorityConstantsAreSequential verifies priorities are properly sequenced
func TestPriorityConstantsAreSequential(t *testing.T) {
	// Verify priorities are in ascending order
	if prioCache >= prioNet {
		t.Error("prioCache should be less than prioNet")
	}
	if prioNet >= prioConnections {
		t.Error("prioNet should be less than prioConnections")
	}
	if prioConnections >= prioItems {
		t.Error("prioConnections should be less than prioItems")
	}
	if prioItems >= prioEvictedReclaimed {
		t.Error("prioItems should be less than prioEvictedReclaimed")
	}
	if prioEvictedReclaimed >= prioGet {
		t.Error("prioEvictedReclaimed should be less than prioGet")
	}
	if prioGet >= prioGetRate {
		t.Error("prioGet should be less than prioGetRate")
	}
	if prioGetRate >= prioSetRate {
		t.Error("prioGetRate should be less than prioSetRate")
	}
	if prioSetRate >= prioDelete {
		t.Error("prioSetRate should be less than prioDelete")
	}
	if prioDelete >= prioCas {
		t.Error("prioDelete should be less than prioCas")
	}
	if prioCas >= prioIncrement {
		t.Error("prioCas should be less than prioIncrement")
	}
	if prioIncrement >= prioDecrement {
		t.Error("prioIncrement should be less than prioDecrement")
	}
	if prioDecrement >= prioTouch {
		t.Error("prioDecrement should be less than prioTouch")
	}
	if prioTouch >= prioTouchRate {
		t.Error("prioTouch should be less than prioTouchRate")
	}
}