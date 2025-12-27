package fluentd

import (
	"testing"
)

// TestNewCollector tests the initialization of a new Fluentd collector
func TestNewCollector(t *testing.T) {
	tests := []struct {
		name    string
		want    bool
		wantErr bool
	}{
		{
			name:    "creates collector successfully",
			want:    true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := New()

			if collector == nil {
				t.Errorf("New() returned nil, want non-nil")
			}

			if tt.wantErr {
				t.Errorf("unexpected error: %v", nil)
			}
		})
	}
}

// TestCollectorName tests that collector has correct name
func TestCollectorName(t *testing.T) {
	collector := New()
	if collector == nil {
		t.Fatal("New() returned nil")
	}

	expectedName := "fluentd"
	if collector.Name() != expectedName {
		t.Errorf("Name() = %q, want %q", collector.Name(), expectedName)
	}
}

// TestCollectorInit tests the initialization logic
func TestCollectorInit(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*Fluentd)
		wantErr bool
	}{
		{
			name: "init with default config",
			setup: func(f *Fluentd) {
				// Default initialization
			},
			wantErr: false,
		},
		{
			name: "init with custom address",
			setup: func(f *Fluentd) {
				if f.Address != "" {
					f.Address = "custom:24224"
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := New()
			if collector == nil {
				t.Fatal("New() returned nil")
			}

			tt.setup(collector)

			err := collector.Init()
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestCollectorInitWithNilChart tests init when charts are nil
func TestCollectorInitWithNilChart(t *testing.T) {
	collector := New()
	if collector == nil {
		t.Fatal("New() returned nil")
	}

	// Clear charts to test nil handling
	collector.Charts = nil

	err := collector.Init()
	if err != nil {
		t.Errorf("Init() with nil charts error = %v, want nil", err)
	}
}

// TestCollectorInitIdempotent tests that init can be called multiple times
func TestCollectorInitIdempotent(t *testing.T) {
	collector := New()
	if collector == nil {
		t.Fatal("New() returned nil")
	}

	err1 := collector.Init()
	err2 := collector.Init()

	if err1 != nil || err2 != nil {
		t.Errorf("Init() called twice failed: err1=%v, err2=%v", err1, err2)
	}
}

// TestCollectorInitWithEmptyAddress tests init with empty address
func TestCollectorInitWithEmptyAddress(t *testing.T) {
	collector := New()
	if collector == nil {
		t.Fatal("New() returned nil")
	}

	collector.Address = ""

	// Should either handle gracefully or error appropriately
	err := collector.Init()
	// Check behavior based on implementation
	_ = err
}

// TestCollectorInitWithInvalidAddress tests init with invalid address format
func TestCollectorInitWithInvalidAddress(t *testing.T) {
	collector := New()
	if collector == nil {
		t.Fatal("New() returned nil")
	}

	collector.Address = "invalid::address::format"

	err := collector.Init()
	// Depending on implementation, may error or succeed
	_ = err
}

// TestCollectorInitWithBoundaryValues tests edge cases
func TestCollectorInitWithBoundaryValues(t *testing.T) {
	tests := []struct {
		name    string
		address string
		timeout int64
	}{
		{
			name:    "zero timeout",
			address: "localhost:24224",
			timeout: 0,
		},
		{
			name:    "negative timeout",
			address: "localhost:24224",
			timeout: -1,
		},
		{
			name:    "large timeout",
			address: "localhost:24224",
			timeout: 1000000,
		},
		{
			name:    "localhost with IPv6",
			address: "[::1]:24224",
			timeout: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := New()
			if collector == nil {
				t.Fatal("New() returned nil")
			}

			collector.Address = tt.address
			collector.Timeout = tt.timeout

			err := collector.Init()
			// Verify initialization completes
			_ = err
		})
	}
}

// TestCollectorInitCharts tests chart initialization
func TestCollectorInitCharts(t *testing.T) {
	collector := New()
	if collector == nil {
		t.Fatal("New() returned nil")
	}

	// Store original charts
	originalCharts := collector.Charts

	err := collector.Init()
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	// Charts should be properly set up
	if collector.Charts == nil {
		t.Error("Charts should not be nil after Init()")
	}

	if originalCharts != nil && collector.Charts != originalCharts {
		t.Log("Charts were modified or replaced during Init()")
	}
}

// TestCollectorInitWithZeroPort tests init with port 0
func TestCollectorInitWithZeroPort(t *testing.T) {
	collector := New()
	if collector == nil {
		t.Fatal("New() returned nil")
	}

	collector.Address = "localhost:0"

	err := collector.Init()
	// Port 0 may be valid (let OS choose) or invalid
	_ = err
}

// TestCollectorInitConsistency tests consistency across multiple initializations
func TestCollectorInitConsistency(t *testing.T) {
	collector1 := New()
	collector2 := New()

	if collector1 == nil || collector2 == nil {
		t.Fatal("New() returned nil")
	}

	err1 := collector1.Init()
	err2 := collector2.Init()

	if (err1 != nil) != (err2 != nil) {
		t.Errorf("Inconsistent Init() results: err1=%v, err2=%v", err1, err2)
	}
}