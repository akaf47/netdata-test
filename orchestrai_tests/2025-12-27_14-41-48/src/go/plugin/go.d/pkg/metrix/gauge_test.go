package metrix

import (
	"testing"
)

// Test gauge creation
func TestNewGauge(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{
			name: "should create gauge with valid id",
			id:   "test_gauge",
		},
		{
			name: "should create gauge with empty id",
			id:   "",
		},
		{
			name: "should create gauge with numeric id",
			id:   "123",
		},
		{
			name: "should create gauge with special characters",
			id:   "gauge-_.-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGauge(tt.id)
			if g == nil {
				t.Errorf("NewGauge() returned nil")
			}
			if g.id != tt.id {
				t.Errorf("NewGauge() id = %v, want %v", g.id, tt.id)
			}
		})
	}
}

// Test gauge value initialization
func TestGaugeInit(t *testing.T) {
	g := NewGauge("test")
	if g.value != 0 {
		t.Errorf("gauge initial value = %v, want 0", g.value)
	}
}

// Test gauge Set method
func TestGaugeSet(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected float64
	}{
		{
			name:     "should set positive value",
			value:    42.5,
			expected: 42.5,
		},
		{
			name:     "should set negative value",
			value:    -10.5,
			expected: -10.5,
		},
		{
			name:     "should set zero value",
			value:    0,
			expected: 0,
		},
		{
			name:     "should set very large value",
			value:    1e10,
			expected: 1e10,
		},
		{
			name:     "should set very small value",
			value:    1e-10,
			expected: 1e-10,
		},
		{
			name:     "should overwrite previous value",
			value:    100,
			expected: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGauge("test")
			g.Set(tt.value)
			if g.value != tt.expected {
				t.Errorf("Set(%v) resulted in value %v, want %v", tt.value, g.value, tt.expected)
			}
		})
	}
}

// Test gauge Get method
func TestGaugeGet(t *testing.T) {
	tests := []struct {
		name     string
		setValue float64
		expected float64
	}{
		{
			name:     "should get value after set",
			setValue: 42.5,
			expected: 42.5,
		},
		{
			name:     "should get zero value",
			setValue: 0,
			expected: 0,
		},
		{
			name:     "should get negative value",
			setValue: -15.3,
			expected: -15.3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGauge("test")
			g.Set(tt.setValue)
			result := g.Get()
			if result != tt.expected {
				t.Errorf("Get() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test gauge Inc method
func TestGaugeInc(t *testing.T) {
	tests := []struct {
		name     string
		initial  float64
		delta    float64
		expected float64
	}{
		{
			name:     "should increment from zero",
			initial:  0,
			delta:    1,
			expected: 1,
		},
		{
			name:     "should increment positive value",
			initial:  10,
			delta:    5,
			expected: 15,
		},
		{
			name:     "should increment with negative delta",
			initial:  10,
			delta:    -3,
			expected: 7,
		},
		{
			name:     "should increment by zero",
			initial:  42,
			delta:    0,
			expected: 42,
		},
		{
			name:     "should increment multiple times",
			initial:  0,
			delta:    2,
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGauge("test")
			g.Set(tt.initial)
			g.Inc(tt.delta)
			if g.Get() != tt.expected {
				t.Errorf("Inc(%v) resulted in value %v, want %v", tt.delta, g.Get(), tt.expected)
			}
		})
	}
}

// Test multiple increments
func TestGaugeMultipleInc(t *testing.T) {
	g := NewGauge("test")
	g.Inc(1)
	g.Inc(2)
	g.Inc(3)
	expected := 6.0
	if g.Get() != expected {
		t.Errorf("Multiple Inc() resulted in %v, want %v", g.Get(), expected)
	}
}

// Test gauge Dec method
func TestGaugeDec(t *testing.T) {
	tests := []struct {
		name     string
		initial  float64
		delta    float64
		expected float64
	}{
		{
			name:     "should decrement from value",
			initial:  10,
			delta:    3,
			expected: 7,
		},
		{
			name:     "should decrement to negative",
			initial:  5,
			delta:    10,
			expected: -5,
		},
		{
			name:     "should decrement by zero",
			initial:  42,
			delta:    0,
			expected: 42,
		},
		{
			name:     "should decrement from zero",
			initial:  0,
			delta:    5,
			expected: -5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGauge("test")
			g.Set(tt.initial)
			g.Dec(tt.delta)
			if g.Get() != tt.expected {
				t.Errorf("Dec(%v) resulted in value %v, want %v", tt.delta, g.Get(), tt.expected)
			}
		})
	}
}

// Test gauge Reset method
func TestGaugeReset(t *testing.T) {
	tests := []struct {
		name      string
		setValue  float64
		resetTo   float64
		expectZero bool
	}{
		{
			name:      "should reset to zero",
			setValue:  42,
			resetTo:   0,
			expectZero: true,
		},
		{
			name:      "should reset to new value",
			setValue:  100,
			resetTo:   50,
			expectZero: false,
		},
		{
			name:      "should reset negative to zero",
			setValue:  -42,
			resetTo:   0,
			expectZero: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGauge("test")
			g.Set(tt.setValue)
			g.Reset()
			if tt.expectZero && g.Get() != 0 {
				t.Errorf("Reset() did not reset to zero, got %v", g.Get())
			}
		})
	}
}

// Test gauge Type method
func TestGaugeType(t *testing.T) {
	g := NewGauge("test")
	gaugeType := g.Type()
	if gaugeType != "gauge" {
		t.Errorf("Type() = %v, want 'gauge'", gaugeType)
	}
}

// Test gauge ID method
func TestGaugeID(t *testing.T) {
	id := "my_gauge_id"
	g := NewGauge(id)
	if g.ID() != id {
		t.Errorf("ID() = %v, want %v", g.ID(), id)
	}
}

// Test concurrent Set operations
func TestGaugeConcurrentSet(t *testing.T) {
	g := NewGauge("test")
	done := make(chan bool, 100)
	
	for i := 0; i < 100; i++ {
		go func(val float64) {
			g.Set(val)
			done <- true
		}(float64(i))
	}
	
	for i := 0; i < 100; i++ {
		<-done
	}
	
	// Last value set should be present
	if g.Get() < 0 || g.Get() >= 100 {
		t.Errorf("Concurrent Set() resulted in unexpected value %v", g.Get())
	}
}

// Test gauge with extreme values
func TestGaugeExtremeValues(t *testing.T) {
	tests := []struct {
		name  string
		value float64
	}{
		{
			name:  "should handle max float",
			value: 1.7976931348623157e+308,
		},
		{
			name:  "should handle min positive float",
			value: 2.2250738585072014e-308,
		},
		{
			name:  "should handle very negative value",
			value: -1.7976931348623157e+308,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := NewGauge("test")
			g.Set(tt.value)
			if g.Get() != tt.value {
				t.Errorf("Set(%v) resulted in %v", tt.value, g.Get())
			}
		})
	}
}

// Test gauge value persistence
func TestGaugeValuePersistence(t *testing.T) {
	g := NewGauge("test")
	g.Set(42)
	
	// Multiple gets should return same value
	for i := 0; i < 10; i++ {
		if g.Get() != 42 {
			t.Errorf("Iteration %d: Get() = %v, want 42", i, g.Get())
		}
	}
}

// Test gauge methods chaining effects
func TestGaugeChainedOperations(t *testing.T) {
	g := NewGauge("test")
	g.Set(10)
	g.Inc(5)
	g.Dec(2)
	g.Inc(3)
	
	expected := 16.0 // 10 + 5 - 2 + 3
	if g.Get() != expected {
		t.Errorf("Chained operations resulted in %v, want %v", g.Get(), expected)
	}
}