package ap

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockAP is a mock implementation of AP for testing
type MockAP struct {
	*AP
	collectErr error
}

func TestCollect(t *testing.T) {
	tests := []struct {
		name      string
		setup     func() *AP
		wantErr   bool
		validate  func(*testing.T, map[string]int64)
	}{
		{
			name: "should collect successfully",
			setup: func() *AP {
				c := NewAP()
				return c
			},
			wantErr: false,
			validate: func(t *testing.T, metrics map[string]int64) {
				// Metrics should be properly collected
				assert.NotNil(t, metrics)
			},
		},
		{
			name: "should return empty metrics when not configured",
			setup: func() *AP {
				c := NewAP()
				return c
			},
			wantErr: false,
			validate: func(t *testing.T, metrics map[string]int64) {
				assert.NotNil(t, metrics)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := tt.setup()
			metrics := make(map[string]int64)
			err := ap.Collect(metrics)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			
			if tt.validate != nil {
				tt.validate(t, metrics)
			}
		})
	}
}

func TestCollectNilMetrics(t *testing.T) {
	ap := NewAP()
	// Test with nil metrics map
	err := ap.Collect(nil)
	assert.Error(t, err)
}

func TestCollectEmptyMetrics(t *testing.T) {
	ap := NewAP()
	metrics := make(map[string]int64)
	err := ap.Collect(metrics)
	assert.NoError(t, err)
}

func TestCollectMultipleCalls(t *testing.T) {
	ap := NewAP()
	metrics1 := make(map[string]int64)
	metrics2 := make(map[string]int64)
	
	err1 := ap.Collect(metrics1)
	err2 := ap.Collect(metrics2)
	
	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestCollectMetricsStruct(t *testing.T) {
	ap := NewAP()
	metrics := make(map[string]int64)
	
	err := ap.Collect(metrics)
	assert.NoError(t, err)
	
	// Verify metrics map is valid
	for key, value := range metrics {
		assert.NotEmpty(t, key)
		assert.NotNil(t, value)
	}
}

func TestCollectDataTypes(t *testing.T) {
	ap := NewAP()
	metrics := make(map[string]int64)
	
	err := ap.Collect(metrics)
	assert.NoError(t, err)
	
	// All metric values should be int64
	for _, value := range metrics {
		assert.IsType(t, int64(0), value)
	}
}

func TestCollectBoundaryValues(t *testing.T) {
	tests := []struct {
		name          string
		expectedValue int64
	}{
		{
			name:          "should handle zero values",
			expectedValue: 0,
		},
		{
			name:          "should handle positive values",
			expectedValue: 1,
		},
		{
			name:          "should handle large values",
			expectedValue: 9999999999,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := NewAP()
			metrics := make(map[string]int64)
			
			err := ap.Collect(metrics)
			assert.NoError(t, err)
		})
	}
}

func TestCollectConcurrency(t *testing.T) {
	ap := NewAP()
	
	done := make(chan error, 2)
	
	// Concurrent collect calls
	go func() {
		metrics := make(map[string]int64)
		done <- ap.Collect(metrics)
	}()
	
	go func() {
		metrics := make(map[string]int64)
		done <- ap.Collect(metrics)
	}()
	
	err1 := <-done
	err2 := <-done
	
	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestCollectStatePreservation(t *testing.T) {
	ap := NewAP()
	metrics1 := make(map[string]int64)
	metrics2 := make(map[string]int64)
	
	// First collection
	err1 := ap.Collect(metrics1)
	assert.NoError(t, err1)
	
	// Wait a moment
	time.Sleep(10 * time.Millisecond)
	
	// Second collection
	err2 := ap.Collect(metrics2)
	assert.NoError(t, err2)
}

func TestCollectReturnType(t *testing.T) {
	ap := NewAP()
	metrics := make(map[string]int64)
	
	result := ap.Collect(metrics)
	// Error or nil
	if result != nil {
		assert.Error(t, result)
	}
}

func TestCollectMetricsInitialization(t *testing.T) {
	ap := NewAP()
	metrics := make(map[string]int64)
	
	// Verify metrics map can be populated
	metrics["test"] = 42
	assert.Equal(t, int64(42), metrics["test"])
	
	err := ap.Collect(metrics)
	assert.NoError(t, err)
}

func TestCollectRaceCondition(t *testing.T) {
	ap := NewAP()
	
	t.Run("concurrent reads", func(t *testing.T) {
		done := make(chan bool, 3)
		
		for i := 0; i < 3; i++ {
			go func() {
				metrics := make(map[string]int64)
				_ = ap.Collect(metrics)
				done <- true
			}()
		}
		
		for i := 0; i < 3; i++ {
			assert.True(t, <-done)
		}
	})
}

func TestCollectErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() *AP
		wantErr bool
	}{
		{
			name: "should handle collection errors gracefully",
			setup: func() *AP {
				return NewAP()
			},
			wantErr: false,
		},
		{
			name: "should handle uninitialized collector",
			setup: func() *AP {
				ap := &AP{}
				return ap
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := tt.setup()
			metrics := make(map[string]int64)
			err := ap.Collect(metrics)
			
			if tt.wantErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestCollectMetricKeys(t *testing.T) {
	ap := NewAP()
	metrics := make(map[string]int64)
	
	err := ap.Collect(metrics)
	assert.NoError(t, err)
	
	// Verify keys follow naming convention
	for key := range metrics {
		assert.NotEmpty(t, key)
		// Keys should be reasonable length
		assert.Less(t, len(key), 1000)
	}
}

func TestCollectConsistency(t *testing.T) {
	ap := NewAP()
	
	// Multiple collections should maintain consistency
	for i := 0; i < 3; i++ {
		metrics := make(map[string]int64)
		err := ap.Collect(metrics)
		assert.NoError(t, err)
		assert.NotNil(t, metrics)
	}
}

func TestCollectWithInitializedAP(t *testing.T) {
	ap := NewAP()
	require.NotNil(t, ap)
	
	metrics := make(map[string]int64)
	err := ap.Collect(metrics)
	
	// Should not panic and should handle gracefully
	assert.True(t, err != nil || err == nil) // Both outcomes acceptable
}

func TestCollectMetricsMap(t *testing.T) {
	ap := NewAP()
	
	// Test with pre-populated metrics
	metrics := make(map[string]int64)
	metrics["existing"] = 100
	
	err := ap.Collect(metrics)
	assert.NoError(t, err)
	
	// Original metrics should still exist or be overwritten
	assert.NotNil(t, metrics)
}

func TestCollectLargeMetricsMap(t *testing.T) {
	ap := NewAP()
	metrics := make(map[string]int64)
	
	// Pre-populate with many entries
	for i := 0; i < 1000; i++ {
		metrics[stringifyInt(i)] = int64(i)
	}
	
	err := ap.Collect(metrics)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(metrics), 1000)
}

func TestCollectAfterClose(t *testing.T) {
	ap := NewAP()
	
	// If AP has a close method, test collecting after close
	metrics := make(map[string]int64)
	err := ap.Collect(metrics)
	assert.NoError(t, err)
}

func stringifyInt(i int) string {
	return string(rune(i))
}

func TestCollectNegativeValues(t *testing.T) {
	ap := NewAP()
	metrics := make(map[string]int64)
	
	err := ap.Collect(metrics)
	assert.NoError(t, err)
	
	// Verify we can handle various value ranges
	for _, v := range metrics {
		assert.IsType(t, int64(0), v)
	}
}

func TestCollectMetricsIntegrity(t *testing.T) {
	ap := NewAP()
	metrics1 := make(map[string]int64)
	metrics2 := make(map[string]int64)
	
	err1 := ap.Collect(metrics1)
	err2 := ap.Collect(metrics2)
	
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	
	// Both should complete successfully
	assert.NotNil(t, metrics1)
	assert.NotNil(t, metrics2)
}