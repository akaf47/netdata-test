package pcf

import (
	"testing"
)

// TestResetQueueStatsNilInput tests handling of nil input
func TestResetQueueStatsNilInput(t *testing.T) {
	result := ResetQueueStats(nil)
	if result != nil {
		t.Errorf("expected nil result for nil input, got %v", result)
	}
}

// TestResetQueueStatsEmptySlice tests handling of empty slice
func TestResetQueueStatsEmptySlice(t *testing.T) {
	emptyStats := []QueueStats{}
	result := ResetQueueStats(emptyStats)
	if result != nil && len(result) != 0 {
		t.Errorf("expected empty or nil result for empty input, got %v", result)
	}
}

// TestResetQueueStatsSingleElement tests with single element
func TestResetQueueStatsSingleElement(t *testing.T) {
	stats := []QueueStats{
		{
			QueueName: "TestQueue",
			GetCount:  100,
			PutCount:  50,
		},
	}
	result := ResetQueueStats(stats)
	if result == nil {
		t.Fatalf("expected non-nil result, got nil")
	}
	if len(result) != 1 {
		t.Errorf("expected 1 element, got %d", len(result))
	}
	if result[0].GetCount != 0 {
		t.Errorf("expected GetCount to be reset to 0, got %d", result[0].GetCount)
	}
	if result[0].PutCount != 0 {
		t.Errorf("expected PutCount to be reset to 0, got %d", result[0].PutCount)
	}
	if result[0].QueueName != "TestQueue" {
		t.Errorf("expected QueueName to remain unchanged, got %s", result[0].QueueName)
	}
}

// TestResetQueueStatsMultipleElements tests with multiple elements
func TestResetQueueStatsMultipleElements(t *testing.T) {
	stats := []QueueStats{
		{QueueName: "Queue1", GetCount: 100, PutCount: 50},
		{QueueName: "Queue2", GetCount: 200, PutCount: 100},
		{QueueName: "Queue3", GetCount: 300, PutCount: 150},
	}
	result := ResetQueueStats(stats)
	if result == nil {
		t.Fatalf("expected non-nil result, got nil")
	}
	if len(result) != 3 {
		t.Errorf("expected 3 elements, got %d", len(result))
	}
	for i, stat := range result {
		if stat.GetCount != 0 {
			t.Errorf("element %d: expected GetCount to be 0, got %d", i, stat.GetCount)
		}
		if stat.PutCount != 0 {
			t.Errorf("element %d: expected PutCount to be 0, got %d", i, stat.PutCount)
		}
	}
}

// TestResetQueueStatsZeroValues tests with already zero values
func TestResetQueueStatsZeroValues(t *testing.T) {
	stats := []QueueStats{
		{QueueName: "Queue1", GetCount: 0, PutCount: 0},
	}
	result := ResetQueueStats(stats)
	if result == nil {
		t.Fatalf("expected non-nil result, got nil")
	}
	if result[0].GetCount != 0 || result[0].PutCount != 0 {
		t.Errorf("expected counts to remain zero, got GetCount=%d, PutCount=%d", result[0].GetCount, result[0].PutCount)
	}
}

// TestResetQueueStatsLargeValues tests with large values
func TestResetQueueStatsLargeValues(t *testing.T) {
	stats := []QueueStats{
		{QueueName: "Queue1", GetCount: 999999999, PutCount: 888888888},
	}
	result := ResetQueueStats(stats)
	if result == nil {
		t.Fatalf("expected non-nil result, got nil")
	}
	if result[0].GetCount != 0 || result[0].PutCount != 0 {
		t.Errorf("expected large values to be reset to 0, got GetCount=%d, PutCount=%d", result[0].GetCount, result[0].PutCount)
	}
}

// TestResetQueueStatsPreservesQueueName tests that queue names are preserved
func TestResetQueueStatsPreservesQueueName(t *testing.T) {
	queueNames := []string{"SYSTEM.QUEUE", "MY.APP.QUEUE", "DLQ", ""}
	for _, name := range queueNames {
		stats := []QueueStats{{QueueName: name, GetCount: 100, PutCount: 50}}
		result := ResetQueueStats(stats)
		if result[0].QueueName != name {
			t.Errorf("queue name not preserved: expected %q, got %q", name, result[0].QueueName)
		}
	}
}

// TestResetQueueStatsPreservesOtherFields tests that other fields are preserved
func TestResetQueueStatsPreservesOtherFields(t *testing.T) {
	stats := []QueueStats{
		{
			QueueName:         "TestQueue",
			GetCount:          100,
			PutCount:          50,
			BrowseCount:       25,
			DeleteCount:       10,
			OpenCount:         5,
			CloseCount:        3,
		},
	}
	result := ResetQueueStats(stats)
	// Only Get and Put counts should be reset; verify other fields remain unchanged
	if result[0].BrowseCount != 25 {
		t.Errorf("expected BrowseCount to be unchanged, got %d", result[0].BrowseCount)
	}
	if result[0].DeleteCount != 10 {
		t.Errorf("expected DeleteCount to be unchanged, got %d", result[0].DeleteCount)
	}
}

// TestResetQueueStatsModifiesOnlyRelevantCounters tests precise counter reset
func TestResetQueueStatsModifiesOnlyRelevantCounters(t *testing.T) {
	originalGetCount := int64(100)
	originalPutCount := int64(50)
	originalBrowseCount := int64(25)

	stats := []QueueStats{
		{
			QueueName:   "Queue1",
			GetCount:    originalGetCount,
			PutCount:    originalPutCount,
			BrowseCount: originalBrowseCount,
		},
	}

	result := ResetQueueStats(stats)

	if result[0].GetCount != 0 {
		t.Errorf("GetCount should be reset to 0, got %d", result[0].GetCount)
	}
	if result[0].PutCount != 0 {
		t.Errorf("PutCount should be reset to 0, got %d", result[0].PutCount)
	}
	if result[0].BrowseCount != originalBrowseCount {
		t.Errorf("BrowseCount should remain %d, got %d", originalBrowseCount, result[0].BrowseCount)
	}
}

// TestResetQueueStatsReturnValue tests that a new slice is returned (or properly modified)
func TestResetQueueStatsReturnValue(t *testing.T) {
	originalStats := []QueueStats{
		{QueueName: "Queue1", GetCount: 100, PutCount: 50},
	}
	resultStats := ResetQueueStats(originalStats)

	if resultStats == nil {
		t.Fatalf("expected non-nil result")
	}

	// Verify the return value has reset counts
	if resultStats[0].GetCount != 0 || resultStats[0].PutCount != 0 {
		t.Errorf("returned stats should have reset counts")
	}
}

// TestResetQueueStatsNegativeValues tests with negative values (edge case)
func TestResetQueueStatsNegativeValues(t *testing.T) {
	stats := []QueueStats{
		{QueueName: "Queue1", GetCount: -100, PutCount: -50},
	}
	result := ResetQueueStats(stats)
	if result == nil {
		t.Fatalf("expected non-nil result, got nil")
	}
	// Even negative values should be reset to 0
	if result[0].GetCount != 0 || result[0].PutCount != 0 {
		t.Errorf("expected negative values to be reset to 0, got GetCount=%d, PutCount=%d", result[0].GetCount, result[0].PutCount)
	}
}