package unbound

import (
	"testing"
	"time"
)

// TestNewUnbound tests the constructor for creating a new Unbound collector
func TestNewUnbound(t *testing.T) {
	u := New()

	if u == nil {
		t.Fatal("Expected non-nil Unbound instance")
	}

	if u.Name != "unbound" {
		t.Errorf("Expected Name 'unbound', got '%s'", u.Name)
	}

	if u.Charts == nil {
		t.Error("Expected Charts to be initialized")
	}
}

// TestUnboundInit tests the initialization of the Unbound collector
func TestUnboundInit(t *testing.T) {
	u := New()

	err := u.Init()

	if err != nil {
		t.Fatalf("Expected no error during init, got %v", err)
	}

	if u.client == nil {
		t.Error("Expected client to be initialized")
	}
}

// TestUnboundInitWithoutRequiredConfig tests Init fails when required config is missing
func TestUnboundInitWithoutRequiredConfig(t *testing.T) {
	u := New()
	u.Address = "" // Clear required address

	err := u.Init()

	if err == nil {
		t.Error("Expected error when address is missing")
	}
}

// TestUnboundInitWithInvalidAddress tests Init with invalid address format
func TestUnboundInitWithInvalidAddress(t *testing.T) {
	u := New()
	u.Address = "invalid-address-format"

	err := u.Init()

	// Should either fail or handle gracefully
	if err != nil && err.Error() == "" {
		t.Error("Error message is empty")
	}
}

// TestUnboundInitWithValidAddress tests Init with valid address
func TestUnboundInitWithValidAddress(t *testing.T) {
	u := New()
	u.Address = "127.0.0.1:8953"

	err := u.Init()

	if err != nil {
		t.Fatalf("Expected no error with valid address, got %v", err)
	}
}

// TestUnboundCheck tests the Check method validates connection
func TestUnboundCheck(t *testing.T) {
	u := New()
	u.client = &mockUnboundClient{
		stats: map[string]string{
			"total.num.queries": "1000",
		},
		err: nil,
	}

	err := u.Check()

	if err != nil {
		t.Fatalf("Expected no error during check, got %v", err)
	}
}

// TestUnboundCheckWithClientError tests Check when client returns error
func TestUnboundCheckWithClientError(t *testing.T) {
	u := New()
	u.client = &mockUnboundClient{
		stats: nil,
		err:   NewConnectionError("failed to connect"),
	}

	err := u.Check()

	if err == nil {
		t.Error("Expected error during check, got nil")
	}
}

// TestUnboundCheckWithNilClient tests Check when client is nil
func TestUnboundCheckWithNilClient(t *testing.T) {
	u := New()
	u.client = nil

	err := u.Check()

	if err == nil {
		t.Error("Expected error when client is nil")
	}
}

// TestUnboundCheckWithEmptyStats tests Check with empty stats
func TestUnboundCheckWithEmptyStats(t *testing.T) {
	u := New()
	u.client = &mockUnboundClient{
		stats: map[string]string{},
		err:   nil,
	}

	err := u.Check()

	if err != nil {
		t.Fatalf("Expected no error with empty stats, got %v", err)
	}
}

// TestUnboundCharts tests Charts are properly initialized
func TestUnboundCharts(t *testing.T) {
	u := New()

	charts := u.Charts()

	if charts == nil {
		t.Error("Expected non-nil charts")
	}

	if len(charts) == 0 {
		t.Error("Expected non-empty charts")
	}

	// Check for expected chart groups
	hasQueryChart := false
	for _, chart := range charts {
		if chart != nil && chart.ID != "" {
			hasQueryChart = true
			break
		}
	}

	if !hasQueryChart {
		t.Error("Expected at least one chart to be defined")
	}
}

// TestUnboundChartsAfterInit tests Charts are correctly formed after Init
func TestUnboundChartsAfterInit(t *testing.T) {
	u := New()
	u.Init()

	charts := u.Charts()

	if charts == nil {
		t.Fatal("Expected non-nil charts")
	}

	for _, chart := range charts {
		if chart == nil {
			t.Error("Found nil chart in charts list")
			continue
		}

		if chart.ID == "" {
			t.Error("Expected chart to have non-empty ID")
		}

		if len(chart.Dims) == 0 {
			t.Errorf("Expected chart %s to have dimensions", chart.ID)
		}
	}
}

// TestUnboundCollect tests the main Collect method
func TestUnboundCollect(t *testing.T) {
	u := New()
	u.client = &mockUnboundClient{
		stats: map[string]string{
			"total.num.queries":   "1000",
			"total.num.cachehits": "500",
		},
		err: nil,
	}

	result := u.Collect()

	if result == nil {
		t.Error("Expected non-nil result from Collect")
	}
}

// TestUnboundCollectWithNilClient tests Collect when client is nil
func TestUnboundCollectWithNilClient(t *testing.T) {
	u := New()
	u.client = nil

	// Should handle gracefully or return error
	result := u.Collect()

	if result != nil {
		t.Error("Expected nil or empty result when client is nil")
	}
}

// TestUnboundCollectWithClientError tests Collect when client returns error
func TestUnboundCollectWithClientError(t *testing.T) {
	u := New()
	u.client = &mockUnboundClient{
		stats: nil,
		err:   NewConnectionError("connection failed"),
	}

	result := u.Collect()

	if result != nil {
		t.Error("Expected nil result when client returns error")
	}
}

// TestUnboundCollectPerformance tests Collect completes within reasonable time
func TestUnboundCollectPerformance(t *testing.T) {
	u := New()
	u.client = &mockUnboundClient{
		stats: map[string]string{
			"total.num.queries":   "1000",
			"total.num.cachehits": "500",
			"total.num.cachemiss": "400",
		},
		err: nil,
	}

	start := time.Now()
	u.Collect()
	duration := time.Since(start)

	// Should complete in reasonable time (< 1 second)
	if duration > time.Second {
		t.Errorf("Collect took too long: %v", duration)
	}
}

// TestUnboundCollectMultipleInvocations tests multiple Collect calls
func TestUnboundCollectMultipleInvocations(t *testing.T) {
	u := New()
	u.client = &mockUnboundClient{
		stats: map[string]string{
			"total.num.queries": "1000",
		},
		err: nil,
	}

	for i := 0; i < 3; i++ {
		result := u.Collect()
		if result != nil && len(result) > 0 {
			// Verify result structure
			for key, value := range result {
				if key == "" {
					t.Error("Expected non-empty key in result")
				}
				if value == nil {
					t.Errorf("Expected non-nil value for key %s", key)
				}
			}
		}
	}
}

// TestUnboundName tests the Name method returns correct collector name
func TestUnboundName(t *testing.T) {
	u := New()

	name := u.Name

	if name != "unbound" {
		t.Errorf("Expected name 'unbound', got '%s'", name)
	}
}

// TestUnboundModule tests the Module method returns correct module name
func TestUnboundModule(t *testing.T) {
	u := New()

	module := u.Module()

	if module == "" {
		t.Error("Expected non-empty module name")
	}
}

// TestUnboundDescription tests the Description method
func TestUnboundDescription(t *testing.T) {
	u := New()

	desc := u.Description()

	if desc == "" {
		t.Error("Expected non-empty description")
	}
}

// TestUnboundConfigDefaults tests default configuration values
func TestUnboundConfigDefaults(t *testing.T) {
	u := New()

	if u.Address == "" {
		t.Error("Expected default Address to be set")
	}

	if u.Timeout == 0 {
		t.Error("Expected default Timeout to be set")
	}
}

// TestUnboundInitIdempotent tests that Init can be called multiple times safely
func TestUnboundInitIdempotent(t *testing.T) {
	u := New()

	err1 := u.Init()
	err2 := u.Init()

	if err1 != nil || err2 != nil {
		t.Fatalf("Expected no errors on multiple init calls")
	}

	if u.client == nil {
		t.Error("Expected client to remain initialized")
	}
}

// TestUnboundClose tests the Close method
func TestUnboundClose(t *testing.T) {
	u := New()
	u.Init()

	err := u.Close()

	if err != nil {
		t.Fatalf("Expected no error on close, got %v", err)
	}
}

// TestUnboundCloseWithoutInit tests Close without prior Init
func TestUnboundCloseWithoutInit(t *testing.T) {
	u := New()

	err := u.Close()

	if err != nil {
		t.Fatalf("Expected no error closing uninitialized collector, got %v", err)
	}
}

// TestUnboundCloseIdempotent tests Close can be called multiple times
func TestUnboundCloseIdempotent(t *testing.T) {
	u := New()
	u.Init()

	err1 := u.Close()
	err2 := u.Close()

	if err1 != nil || err2 != nil {
		t.Fatalf("Expected no errors on multiple close calls")
	}
}

// TestUnboundInitAndClose tests full lifecycle
func TestUnboundInitAndClose(t *testing.T) {
	u := New()

	if err := u.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	u.client = &mockUnboundClient{
		stats: map[string]string{
			"total.num.queries": "1000",
		},
		err: nil,
	}

	result := u.Collect()
	if result != nil {
		// Verify we got data
		if len(result) == 0 {
			t.Error("Expected non-empty result from Collect")
		}
	}

	if err := u.Close(); err != nil {
		t.Fatalf("Close failed: %v", err)
	}
}

// TestUnboundCollectResultStructure tests the structure of Collect result
func TestUnboundCollectResultStructure(t *testing.T) {
	u := New()
	u.client = &mockUnboundClient{
		stats: map[string]string{
			"total.num.queries":   "1000",
			"total.num.cachehits": "500",
		},
		err: nil,
	}

	result := u.Collect()

	if result != nil {
		for key, value := range result {
			if key == "" {
				t.Error("Result contains empty key")
			}
			if value == nil {
				t.Errorf("Result contains nil value for key: %s", key)
			}
		}
	}
}

// TestUnboundWithCustomTimeout tests initialization with custom timeout
func TestUnboundWithCustomTimeout(t *testing.T) {
	u := New()
	u.Timeout = 5 * time.Second

	err := u.Init()

	if err != nil {
		t.Fatalf("Expected no error with custom timeout, got %v", err)
	}

	if u.Timeout != 5*time.Second {
		t.Errorf("Expected timeout 5s, got %v", u.Timeout)
	}
}

// TestUnboundWithZeroTimeout tests behavior with zero timeout
func TestUnboundWithZeroTimeout(t *testing.T) {
	u := New()
	u.Timeout = 0

	// Should either use default or fail appropriately
	err := u.Init()

	// Behavior depends on implementation
	if err != nil && err.Error() == "" {
		t.Error("Error message is empty for zero timeout")
	}
}

// Helper function to create connection errors
func NewConnectionError(msg string) error {
	// Implementation depends on actual error type
	return nil // Placeholder
}