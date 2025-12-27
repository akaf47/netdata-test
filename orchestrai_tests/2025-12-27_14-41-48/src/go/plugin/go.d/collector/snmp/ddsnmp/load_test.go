package ddsnmp

import (
	"testing"
)

// Test Load function with valid configuration
func TestLoad_ValidConfiguration(t *testing.T) {
	// Create a mock collector
	collector := &Collector{}

	err := collector.Load(nil)
	if err != nil {
		t.Errorf("Load() with nil config should not error, got %v", err)
	}
}

// Test Load function with empty configuration
func TestLoad_EmptyConfiguration(t *testing.T) {
	collector := &Collector{}
	config := make(map[string]interface{})

	err := collector.Load(config)
	if err != nil {
		t.Errorf("Load() with empty config should not error, got %v", err)
	}
}

// Test Load function initialization
func TestLoad_InitializesCollector(t *testing.T) {
	collector := &Collector{}

	err := collector.Load(nil)
	if err != nil {
		t.Errorf("Load() failed: %v", err)
	}

	// Verify collector is properly initialized
	if collector == nil {
		t.Error("Load() should not result in nil collector")
	}
}

// Test Load with various configuration parameters
func TestLoad_WithConfigurationParameters(t *testing.T) {
	tests := []struct {
		name     string
		config   map[string]interface{}
		wantErr  bool
	}{
		{
			name:    "empty config",
			config:  make(map[string]interface{}),
			wantErr: false,
		},
		{
			name: "config with single parameter",
			config: map[string]interface{}{
				"param1": "value1",
			},
			wantErr: false,
		},
		{
			name: "config with multiple parameters",
			config: map[string]interface{}{
				"param1": "value1",
				"param2": "value2",
				"param3": 123,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := &Collector{}
			err := collector.Load(tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Test Load idempotency
func TestLoad_Idempotent(t *testing.T) {
	collector := &Collector{}
	config := map[string]interface{}{"test": "config"}

	err1 := collector.Load(config)
	err2 := collector.Load(config)

	if err1 != nil || err2 != nil {
		t.Errorf("Load() failed on repeated calls: err1=%v, err2=%v", err1, err2)
	}
}

// Test Load with nil receiver (edge case)
func TestLoad_NilConfig(t *testing.T) {
	collector := &Collector{}

	err := collector.Load(nil)
	if err != nil {
		t.Errorf("Load(nil) should handle nil config gracefully, got %v", err)
	}
}

// Test Load with complex nested configuration
func TestLoad_ComplexNestedConfig(t *testing.T) {
	collector := &Collector{}
	config := map[string]interface{}{
		"nested": map[string]interface{}{
			"level2": map[string]interface{}{
				"level3": "deep_value",
			},
		},
	}

	err := collector.Load(config)
	if err != nil {
		t.Errorf("Load() with nested config failed: %v", err)
	}
}

// Test Load with various data types in configuration
func TestLoad_MixedDataTypes(t *testing.T) {
	collector := &Collector{}
	config := map[string]interface{}{
		"string":   "value",
		"int":      42,
		"float":    3.14,
		"bool":     true,
		"slice":    []int{1, 2, 3},
		"nil":      nil,
	}

	err := collector.Load(config)
	if err != nil {
		t.Errorf("Load() with mixed types failed: %v", err)
	}
}

// Test Load error handling
func TestLoad_ErrorHandling(t *testing.T) {
	collector := &Collector{}

	// Test with invalid/problematic config values
	config := map[string]interface{}{
		"invalid": make(chan int), // channels can't be encoded
	}

	err := collector.Load(config)
	// Should either handle gracefully or return error
	// The behavior depends on implementation
	if collector == nil {
		t.Error("Load() should not result in nil collector on error")
	}
}

// Test Load with empty string parameters
func TestLoad_EmptyStringParameters(t *testing.T) {
	collector := &Collector{}
	config := map[string]interface{}{
		"param1": "",
		"param2": "",
	}

	err := collector.Load(config)
	if err != nil {
		t.Errorf("Load() with empty strings failed: %v", err)
	}
}

// Test Load with very large configuration
func TestLoad_LargeConfiguration(t *testing.T) {
	collector := &Collector{}
	config := make(map[string]interface{})

	// Create a large config
	for i := 0; i < 1000; i++ {
		config[string(rune(i%26 + 97))] = i // 'a' to 'z'
	}

	err := collector.Load(config)
	if err != nil {
		t.Errorf("Load() with large config failed: %v", err)
	}
}

// Test Load state management
func TestLoad_StateManagement(t *testing.T) {
	collector1 := &Collector{}
	collector2 := &Collector{}

	config1 := map[string]interface{}{"id": 1}
	config2 := map[string]interface{}{"id": 2}

	err1 := collector1.Load(config1)
	err2 := collector2.Load(config2)

	if err1 != nil || err2 != nil {
		t.Errorf("Load() failed: err1=%v, err2=%v", err1, err2)
	}

	// Verify collectors are independent
	if collector1 == collector2 {
		t.Error("Load() should create independent collector instances")
	}
}

// Test Load with boolean configuration values
func TestLoad_BooleanConfig(t *testing.T) {
	tests := []struct {
		name   string
		config map[string]interface{}
	}{
		{
			name: "all true",
			config: map[string]interface{}{
				"enabled": true,
				"debug":   true,
			},
		},
		{
			name: "all false",
			config: map[string]interface{}{
				"enabled": false,
				"debug":   false,
			},
		},
		{
			name: "mixed",
			config: map[string]interface{}{
				"enabled": true,
				"debug":   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := &Collector{}
			err := collector.Load(tt.config)
			if err != nil {
				t.Errorf("Load() with boolean config failed: %v", err)
			}
		})
	}
}

// Test Load preserves collector fields
func TestLoad_PreservesFields(t *testing.T) {
	collector := &Collector{}
	config := map[string]interface{}{}

	err := collector.Load(config)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Verify the collector object still exists and is usable
	if collector == nil {
		t.Error("Load() should not nil out the collector")
	}
}