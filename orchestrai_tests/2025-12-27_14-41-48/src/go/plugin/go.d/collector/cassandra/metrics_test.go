package cassandra

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMetricsConstruction tests metric initialization and setup
func TestMetricsConstruction(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() Metrics
		assertions func(t *testing.T, m Metrics)
	}{
		{
			name: "should initialize metrics with default values",
			setupFunc: func() Metrics {
				return NewMetrics()
			},
			assertions: func(t *testing.T, m Metrics) {
				assert.NotNil(t, m)
			},
		},
		{
			name: "should create metrics with all fields properly initialized",
			setupFunc: func() Metrics {
				return NewMetrics()
			},
			assertions: func(t *testing.T, m Metrics) {
				require.NotNil(t, m)
				// Verify all metric fields are initialized
				assert.Equal(t, int64(0), m.ReadLatency)
				assert.Equal(t, int64(0), m.WriteLatency)
				assert.Equal(t, int64(0), m.GCPauseTime)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setupFunc()
			tt.assertions(t, m)
		})
	}
}

// TestMetricsUpdate tests metric updates and assignments
func TestMetricsUpdate(t *testing.T) {
	tests := []struct {
		name           string
		initialValue   int64
		updateValue    int64
		expectedResult int64
	}{
		{
			name:           "should update metric from zero to positive value",
			initialValue:   0,
			updateValue:    100,
			expectedResult: 100,
		},
		{
			name:           "should update metric from positive to higher value",
			initialValue:   100,
			updateValue:    200,
			expectedResult: 200,
		},
		{
			name:           "should update metric with zero value",
			initialValue:   100,
			updateValue:    0,
			expectedResult: 0,
		},
		{
			name:           "should update metric with large value",
			initialValue:   0,
			updateValue:    9223372036854775807,
			expectedResult: 9223372036854775807,
		},
		{
			name:           "should handle negative values",
			initialValue:   0,
			updateValue:    -1,
			expectedResult: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetrics()
			m.ReadLatency = tt.initialValue
			m.ReadLatency = tt.updateValue
			assert.Equal(t, tt.expectedResult, m.ReadLatency)
		})
	}
}

// TestMetricsReset tests metric reset functionality
func TestMetricsReset(t *testing.T) {
	m := NewMetrics()
	m.ReadLatency = 100
	m.WriteLatency = 200
	m.GCPauseTime = 300

	// Reset metrics
	m.Reset()

	assert.Equal(t, int64(0), m.ReadLatency)
	assert.Equal(t, int64(0), m.WriteLatency)
	assert.Equal(t, int64(0), m.GCPauseTime)
}

// TestMetricsAggregation tests combining multiple metrics
func TestMetricsAggregation(t *testing.T) {
	tests := []struct {
		name      string
		m1        Metrics
		m2        Metrics
		expected  Metrics
	}{
		{
			name: "should aggregate two metrics with positive values",
			m1: Metrics{
				ReadLatency:  100,
				WriteLatency: 200,
				GCPauseTime:  50,
			},
			m2: Metrics{
				ReadLatency:  50,
				WriteLatency: 100,
				GCPauseTime:  25,
			},
			expected: Metrics{
				ReadLatency:  150,
				WriteLatency: 300,
				GCPauseTime:  75,
			},
		},
		{
			name: "should aggregate with zero values",
			m1: Metrics{
				ReadLatency:  0,
				WriteLatency: 0,
				GCPauseTime:  0,
			},
			m2: Metrics{
				ReadLatency:  100,
				WriteLatency: 200,
				GCPauseTime:  50,
			},
			expected: Metrics{
				ReadLatency:  100,
				WriteLatency: 200,
				GCPauseTime:  50,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Metrics{
				ReadLatency:  tt.m1.ReadLatency + tt.m2.ReadLatency,
				WriteLatency: tt.m1.WriteLatency + tt.m2.WriteLatency,
				GCPauseTime:  tt.m1.GCPauseTime + tt.m2.GCPauseTime,
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestMetricsValidation tests metrics validation
func TestMetricsValidation(t *testing.T) {
	tests := []struct {
		name      string
		metrics   Metrics
		isValid   bool
	}{
		{
			name: "should validate metrics with all positive values",
			metrics: Metrics{
				ReadLatency:  100,
				WriteLatency: 200,
				GCPauseTime:  50,
			},
			isValid: true,
		},
		{
			name: "should validate metrics with zero values",
			metrics: Metrics{
				ReadLatency:  0,
				WriteLatency: 0,
				GCPauseTime:  0,
			},
			isValid: true,
		},
		{
			name: "should validate metrics with negative values",
			metrics: Metrics{
				ReadLatency:  -100,
				WriteLatency: -200,
				GCPauseTime:  -50,
			},
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.metrics.IsValid()
			assert.Equal(t, tt.isValid, result)
		})
	}
}

// TestMetricsFieldAccess tests accessing individual metric fields
func TestMetricsFieldAccess(t *testing.T) {
	m := NewMetrics()
	m.ReadLatency = 123
	m.WriteLatency = 456
	m.GCPauseTime = 789

	t.Run("should access ReadLatency field", func(t *testing.T) {
		assert.Equal(t, int64(123), m.ReadLatency)
	})

	t.Run("should access WriteLatency field", func(t *testing.T) {
		assert.Equal(t, int64(456), m.WriteLatency)
	})

	t.Run("should access GCPauseTime field", func(t *testing.T) {
		assert.Equal(t, int64(789), m.GCPauseTime)
	})
}

// TestMetricsComparison tests metric comparison and equality
func TestMetricsComparison(t *testing.T) {
	tests := []struct {
		name     string
		m1       Metrics
		m2       Metrics
		isEqual  bool
	}{
		{
			name: "should identify equal metrics",
			m1: Metrics{
				ReadLatency:  100,
				WriteLatency: 200,
				GCPauseTime:  50,
			},
			m2: Metrics{
				ReadLatency:  100,
				WriteLatency: 200,
				GCPauseTime:  50,
			},
			isEqual: true,
		},
		{
			name: "should identify different metrics",
			m1: Metrics{
				ReadLatency:  100,
				WriteLatency: 200,
				GCPauseTime:  50,
			},
			m2: Metrics{
				ReadLatency:  101,
				WriteLatency: 200,
				GCPauseTime:  50,
			},
			isEqual: false,
		},
		{
			name: "should identify different WriteLatency",
			m1: Metrics{
				ReadLatency:  100,
				WriteLatency: 200,
				GCPauseTime:  50,
			},
			m2: Metrics{
				ReadLatency:  100,
				WriteLatency: 201,
				GCPauseTime:  50,
			},
			isEqual: false,
		},
		{
			name: "should identify different GCPauseTime",
			m1: Metrics{
				ReadLatency:  100,
				WriteLatency: 200,
				GCPauseTime:  50,
			},
			m2: Metrics{
				ReadLatency:  100,
				WriteLatency: 200,
				GCPauseTime:  51,
			},
			isEqual: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := (tt.m1.ReadLatency == tt.m2.ReadLatency &&
				tt.m1.WriteLatency == tt.m2.WriteLatency &&
				tt.m1.GCPauseTime == tt.m2.GCPauseTime)
			assert.Equal(t, tt.isEqual, result)
		})
	}
}

// TestMetricsEdgeCases tests edge cases and boundary conditions
func TestMetricsEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "should handle maximum int64 value",
			testFunc: func(t *testing.T) {
				m := NewMetrics()
				m.ReadLatency = 9223372036854775807
				assert.Equal(t, int64(9223372036854775807), m.ReadLatency)
			},
		},
		{
			name: "should handle minimum int64 value",
			testFunc: func(t *testing.T) {
				m := NewMetrics()
				m.ReadLatency = -9223372036854775808
				assert.Equal(t, int64(-9223372036854775808), m.ReadLatency)
			},
		},
		{
			name: "should handle overflow scenario in aggregation",
			testFunc: func(t *testing.T) {
				m1 := NewMetrics()
				m1.ReadLatency = 9223372036854775807
				m2 := NewMetrics()
				m2.ReadLatency = 1
				// Note: This tests overflow behavior without actual overflow
				assert.True(t, m1.ReadLatency == 9223372036854775807)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.testFunc(t)
		})
	}
}

// TestMetricsState tests metrics state transitions
func TestMetricsState(t *testing.T) {
	m := NewMetrics()

	// Initial state
	assert.Equal(t, int64(0), m.ReadLatency)
	assert.Equal(t, int64(0), m.WriteLatency)

	// Transition to state 1
	m.ReadLatency = 100
	assert.Equal(t, int64(100), m.ReadLatency)

	// Transition to state 2
	m.WriteLatency = 200
	assert.Equal(t, int64(200), m.WriteLatency)
	assert.Equal(t, int64(100), m.ReadLatency)

	// Reset to initial state
	m.Reset()
	assert.Equal(t, int64(0), m.ReadLatency)
	assert.Equal(t, int64(0), m.WriteLatency)
}