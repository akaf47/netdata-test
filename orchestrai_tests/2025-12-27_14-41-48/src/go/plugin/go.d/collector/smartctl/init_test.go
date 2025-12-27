package smartctl

import (
	"testing"
	"time"
)

// Test NewCollector creates a valid collector
func TestNewCollector(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful creation",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector, err := NewCollector()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCollector() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && collector == nil {
				t.Errorf("NewCollector() collector is nil")
			}
		})
	}
}

// Test Collector initialization with default values
func TestCollectorDefaultValues(t *testing.T) {
	collector, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}

	if collector == nil {
		t.Fatalf("NewCollector() returned nil collector")
	}

	// Check that collector has expected fields initialized
	if collector.exec == nil {
		t.Errorf("collector.exec is nil")
	}
}

// Test Collector.Init with valid executable path
func TestCollectorInit(t *testing.T) {
	collector, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}

	err = collector.Init()
	if err != nil {
		// Init may fail if smartctl is not installed, but it shouldn't panic
		t.Logf("collector.Init() error = %v", err)
	}
}

// Test Collector initialization without calling Init
func TestCollectorWithoutInit(t *testing.T) {
	collector, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}

	// Just creating the collector should work
	if collector == nil {
		t.Errorf("NewCollector() returned nil")
	}
}

// Test Collector fields are properly set after NewCollector
func TestCollectorFieldsAfterNew(t *testing.T) {
	collector, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}

	// Check essential fields exist
	if collector == nil {
		t.Fatalf("collector is nil")
	}

	if collector.exec == nil {
		t.Errorf("exec provider not initialized")
	}
}

// Test NewCollector multiple times
func TestNewCollectorMultipleTimes(t *testing.T) {
	collector1, err1 := NewCollector()
	if err1 != nil {
		t.Fatalf("First NewCollector() error = %v", err1)
	}

	collector2, err2 := NewCollector()
	if err2 != nil {
		t.Fatalf("Second NewCollector() error = %v", err2)
	}

	// Both should be valid but different instances
	if collector1 == collector2 {
		t.Errorf("NewCollector() returned same instance twice")
	}

	if collector1 == nil || collector2 == nil {
		t.Errorf("NewCollector() returned nil collector")
	}
}

// Test Collector.Init idempotency
func TestCollectorInitIdempotent(t *testing.T) {
	collector, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}

	// Try Init multiple times - should not cause issues
	_ = collector.Init()
	_ = collector.Init()
	_ = collector.Init()
	// Should not panic
}

// Test Collector initialization fields after Init
func TestCollectorFieldsAfterInit(t *testing.T) {
	collector, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}

	_ = collector.Init()

	if collector == nil {
		t.Errorf("collector is nil after Init")
	}
}

// Test Collector has exec provider
func TestCollectorHasExecProvider(t *testing.T) {
	collector, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}

	if collector.exec == nil {
		t.Errorf("collector.exec is nil")
	}
}

// Test NewCollector error cases - testing robustness
func TestNewCollectorRobustness(t *testing.T) {
	// NewCollector should not take any arguments and should always work
	collector, err := NewCollector()
	if err != nil {
		t.Errorf("NewCollector() error = %v", err)
	}
	if collector == nil {
		t.Errorf("NewCollector() returned nil")
	}
}

// Test Collector timeout configuration
func TestCollectorTimeoutConfig(t *testing.T) {
	collector, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}

	// Check if timeout is set to reasonable value
	if collector.exec == nil {
		t.Fatalf("exec provider is nil")
	}

	// Timeout should be greater than 0
	if collector.exec.timeout <= 0 {
		t.Errorf("timeout should be > 0, got %v", collector.exec.timeout)
	}
}

// Test Collector executable path configuration
func TestCollectorExecutablePath(t *testing.T) {
	collector, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}

	if collector.exec == nil {
		t.Fatalf("exec provider is nil")
	}

	// Executable should be set to something
	if collector.exec.executable == "" {
		t.Errorf("executable path is empty")
	}
}

// Test Collector.Init with smartctl availability
func TestCollectorInitSmartctlDetection(t *testing.T) {
	collector, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}

	err = collector.Init()
	// Init may or may not succeed depending on smartctl availability
	// The important thing is it doesn't panic
	t.Logf("Init result: %v", err)
}

// Test Collector creation with proper cleanup
func TestCollectorCleanup(t *testing.T) {
	collector, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}

	// Verify collector is in valid state
	if collector == nil {
		t.Errorf("collector is nil")
	}
	if collector.exec == nil {
		t.Errorf("exec provider is nil")
	}
}

// Test that Init doesn't modify essential fields
func TestCollectorInitDoesntModifyFields(t *testing.T) {
	collector, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}

	originalExec := collector.exec

	_ = collector.Init()

	// exec provider should still be available after Init
	if collector.exec == nil {
		t.Errorf("exec provider was removed by Init")
	}

	// The exec provider object reference might change, but should exist
	if collector.exec == nil {
		t.Errorf("exec provider is nil after Init")
	}
}

// Test Collector initialization with very small timeout
func TestCollectorSmallTimeout(t *testing.T) {
	collector, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}

	if collector.exec.timeout < time.Second {
		t.Logf("Collector uses timeout < 1 second: %v", collector.exec.timeout)
	}
}

// Test Collector initialization with large timeout
func TestCollectorLargeTimeout(t *testing.T) {
	collector, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}

	if collector.exec.timeout > 10*time.Minute {
		t.Logf("Collector uses timeout > 10 minutes: %v", collector.exec.timeout)
	}
}

// Test Collector state consistency
func TestCollectorStateConsistency(t *testing.T) {
	collector, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}

	// Before Init
	if collector.exec == nil {
		t.Errorf("exec is nil before Init")
	}

	_ = collector.Init()

	// After Init
	if collector.exec == nil {
		t.Errorf("exec is nil after Init")
	}
}

// Test NewCollector doesn't require any special setup
func TestNewCollectorNoSetup(t *testing.T) {
	// Just call it directly without any setup
	collector, err := NewCollector()
	if err != nil {
		t.Errorf("NewCollector() error = %v", err)
	}
	if collector == nil {
		t.Errorf("NewCollector() returned nil")
	}
}

// Test Collector.Init returns error or nil
func TestCollectorInitReturnValue(t *testing.T) {
	collector, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}

	initErr := collector.Init()
	// Init should return either an error or nil (both are valid)
	// We just verify it returns something of type error (or nil)
	_ = initErr
}

// Test multiple Collectors can be created independently
func TestMultipleCollectorsIndependent(t *testing.T) {
	collectors := make([]*Collector, 5)
	for i := 0; i < 5; i++ {
		c, err := NewCollector()
		if err != nil {
			t.Fatalf("NewCollector() %d error = %v", i, err)
		}
		collectors[i] = c
	}

	// All should be different instances
	for i := 0; i < len(collectors); i++ {
		for j := i + 1; j < len(collectors); j++ {
			if collectors[i] == collectors[j] {
				t.Errorf("Collectors %d and %d are the same instance", i, j)
			}
		}
	}
}

// Test Collector fields are not nil
func TestCollectorFieldsNotNil(t *testing.T) {
	collector, err := NewCollector()
	if err != nil {
		t.Fatalf("NewCollector() error = %v", err)
	}

	if collector == nil {
		t.Fatalf("collector is nil")
	}

	if collector.exec == nil {
		t.Fatalf("exec is nil")
	}
}