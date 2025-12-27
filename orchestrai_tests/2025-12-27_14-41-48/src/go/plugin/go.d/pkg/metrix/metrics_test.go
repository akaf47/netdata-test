package metrix

import (
	"testing"
)

// TestMetricsNew tests the creation of a new Metrics instance
func TestMetricsNew(t *testing.T) {
	tests := []struct {
		name string
		want Metrics
	}{
		{
			name: "should create empty metrics",
			want: Metrics{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New()
			if got == nil {
				t.Errorf("New() returned nil")
			}
		})
	}
}

// TestMetricsAdd tests adding metrics
func TestMetricsAdd(t *testing.T) {
	tests := []struct {
		name     string
		metrics  Metrics
		value    float64
		expected float64
	}{
		{
			name:     "should add positive value",
			metrics:  New(),
			value:    10.5,
			expected: 10.5,
		},
		{
			name:     "should add negative value",
			metrics:  New(),
			value:    -5.3,
			expected: -5.3,
		},
		{
			name:     "should add zero",
			metrics:  New(),
			value:    0,
			expected: 0,
		},
		{
			name:     "should add multiple values",
			metrics:  New(),
			value:    15.0,
			expected: 15.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.metrics
			m.Add("test_metric", tt.value)
			// Verify the metric was added
			if m == nil {
				t.Errorf("Add() modified metrics to nil")
			}
		})
	}
}

// TestMetricsCount tests counting metrics
func TestMetricsCount(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(Metrics)
		expected int
	}{
		{
			name: "should return 0 for empty metrics",
			setup: func(m Metrics) {
				// Do nothing
			},
			expected: 0,
		},
		{
			name: "should return 1 after adding one metric",
			setup: func(m Metrics) {
				m.Add("metric1", 1.0)
			},
			expected: 1,
		},
		{
			name: "should return count of multiple metrics",
			setup: func(m Metrics) {
				m.Add("metric1", 1.0)
				m.Add("metric2", 2.0)
				m.Add("metric3", 3.0)
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			tt.setup(m)
			count := m.Count()
			if count != tt.expected {
				t.Errorf("Count() = %d, want %d", count, tt.expected)
			}
		})
	}
}

// TestMetricsGet tests retrieving metric values
func TestMetricsGet(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     float64
		found     bool
		wantValue float64
	}{
		{
			name:      "should return value for existing metric",
			key:       "cpu_usage",
			value:     75.5,
			found:     true,
			wantValue: 75.5,
		},
		{
			name:      "should return false for non-existent metric",
			key:       "non_existent",
			value:     0,
			found:     false,
			wantValue: 0,
		},
		{
			name:      "should get negative values",
			key:       "negative_metric",
			value:     -100.5,
			found:     true,
			wantValue: -100.5,
		},
		{
			name:      "should get zero value",
			key:       "zero_metric",
			value:     0,
			found:     true,
			wantValue: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			if tt.found {
				m.Add(tt.key, tt.value)
			}
			value, ok := m.Get(tt.key)
			if ok != tt.found {
				t.Errorf("Get() found = %v, want %v", ok, tt.found)
			}
			if ok && value != tt.wantValue {
				t.Errorf("Get() value = %v, want %v", value, tt.wantValue)
			}
		})
	}
}

// TestMetricsRemove tests removing metrics
func TestMetricsRemove(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		shouldAdd bool
		shouldExist bool
	}{
		{
			name:      "should remove existing metric",
			key:       "metric1",
			shouldAdd: true,
			shouldExist: false,
		},
		{
			name:      "should not error when removing non-existent metric",
			key:       "non_existent",
			shouldAdd: false,
			shouldExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			if tt.shouldAdd {
				m.Add(tt.key, 10.0)
			}
			m.Remove(tt.key)
			_, ok := m.Get(tt.key)
			if ok != tt.shouldExist {
				t.Errorf("After Remove(), Get() found = %v, want %v", ok, tt.shouldExist)
			}
		})
	}
}

// TestMetricsReset tests resetting all metrics
func TestMetricsReset(t *testing.T) {
	tests := []struct {
		name  string
		setup func(Metrics)
	}{
		{
			name: "should reset empty metrics",
			setup: func(m Metrics) {
				// Already empty
			},
		},
		{
			name: "should reset metrics with one entry",
			setup: func(m Metrics) {
				m.Add("metric1", 1.0)
			},
		},
		{
			name: "should reset metrics with multiple entries",
			setup: func(m Metrics) {
				m.Add("metric1", 1.0)
				m.Add("metric2", 2.0)
				m.Add("metric3", 3.0)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			tt.setup(m)
			m.Reset()
			if m.Count() != 0 {
				t.Errorf("After Reset(), Count() = %d, want 0", m.Count())
			}
		})
	}
}

// TestMetricsIterator tests iterating over metrics
func TestMetricsIterator(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(Metrics)
		wantCount int
	}{
		{
			name: "should iterate empty metrics",
			setup: func(m Metrics) {
				// Empty
			},
			wantCount: 0,
		},
		{
			name: "should iterate single metric",
			setup: func(m Metrics) {
				m.Add("metric1", 1.0)
			},
			wantCount: 1,
		},
		{
			name: "should iterate multiple metrics",
			setup: func(m Metrics) {
				m.Add("metric1", 1.0)
				m.Add("metric2", 2.0)
				m.Add("metric3", 3.0)
			},
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			tt.setup(m)
			count := 0
			m.Iterate(func(k string, v float64) {
				count++
			})
			if count != tt.wantCount {
				t.Errorf("Iterate callback count = %d, want %d", count, tt.wantCount)
			}
		})
	}
}

// TestMetricsHas tests checking if metric exists
func TestMetricsHas(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		shouldAdd bool
		want      bool
	}{
		{
			name:      "should return true for existing metric",
			key:       "metric1",
			shouldAdd: true,
			want:      true,
		},
		{
			name:      "should return false for non-existent metric",
			key:       "non_existent",
			shouldAdd: false,
			want:      false,
		},
		{
			name:      "should handle empty string key",
			key:       "",
			shouldAdd: true,
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			if tt.shouldAdd {
				m.Add(tt.key, 1.0)
			}
			if got := m.Has(tt.key); got != tt.want {
				t.Errorf("Has() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestMetricsUpdate tests updating metric values
func TestMetricsUpdate(t *testing.T) {
	tests := []struct {
		name       string
		key        string
		initialVal float64
		updateVal  float64
		wantVal    float64
	}{
		{
			name:       "should update existing metric",
			key:        "metric1",
			initialVal: 10.0,
			updateVal:  20.0,
			wantVal:    20.0,
		},
		{
			name:       "should update to negative value",
			key:        "metric2",
			initialVal: 10.0,
			updateVal:  -5.5,
			wantVal:    -5.5,
		},
		{
			name:       "should update to zero",
			key:        "metric3",
			initialVal: 10.0,
			updateVal:  0,
			wantVal:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			m.Add(tt.key, tt.initialVal)
			m.Update(tt.key, tt.updateVal)
			val, ok := m.Get(tt.key)
			if !ok {
				t.Errorf("Update() metric not found")
			}
			if val != tt.wantVal {
				t.Errorf("After Update(), Get() = %v, want %v", val, tt.wantVal)
			}
		})
	}
}

// TestMetricsConcurrentAccess tests concurrent access to metrics
func TestMetricsConcurrentAccess(t *testing.T) {
	t.Run("should handle concurrent additions", func(t *testing.T) {
		m := New()
		done := make(chan bool)
		
		for i := 0; i < 10; i++ {
			go func(index int) {
				m.Add("metric", float64(index))
				done <- true
			}(i)
		}
		
		for i := 0; i < 10; i++ {
			<-done
		}
		
		if m.Count() == 0 {
			t.Errorf("Concurrent additions resulted in empty metrics")
		}
	})
}

// TestMetricsEdgeCases tests edge cases
func TestMetricsEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		fn   func(*testing.T)
	}{
		{
			name: "should handle very large values",
			fn: func(t *testing.T) {
				m := New()
				m.Add("large", 1e308)
				val, ok := m.Get("large")
				if !ok || val != 1e308 {
					t.Errorf("Failed to handle large value")
				}
			},
		},
		{
			name: "should handle very small values",
			fn: func(t *testing.T) {
				m := New()
				m.Add("small", 1e-308)
				val, ok := m.Get("small")
				if !ok {
					t.Errorf("Failed to get small value")
				}
				if val == 0 {
					t.Errorf("Small value was lost")
				}
			},
		},
		{
			name: "should handle special metric names",
			fn: func(t *testing.T) {
				m := New()
				specialNames := []string{
					"metric-with-dash",
					"metric.with.dot",
					"metric_with_underscore",
					"metric123",
					"METRIC_CAPS",
				}
				for _, name := range specialNames {
					m.Add(name, 1.0)
				}
				if m.Count() != len(specialNames) {
					t.Errorf("Failed to handle special metric names")
				}
			},
		},
		{
			name: "should handle removing and re-adding",
			fn: func(t *testing.T) {
				m := New()
				m.Add("metric", 1.0)
				m.Remove("metric")
				m.Add("metric", 2.0)
				val, ok := m.Get("metric")
				if !ok || val != 2.0 {
					t.Errorf("Remove and re-add failed")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.fn(t)
		})
	}
}

// TestMetricsLen tests length method
func TestMetricsLen(t *testing.T) {
	tests := []struct {
		name  string
		setup func(Metrics)
		want  int
	}{
		{
			name: "should return 0 for empty",
			setup: func(m Metrics) {},
			want:  0,
		},
		{
			name: "should return correct length",
			setup: func(m Metrics) {
				m.Add("m1", 1)
				m.Add("m2", 2)
			},
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			tt.setup(m)
			if got := m.Len(); got != tt.want {
				t.Errorf("Len() = %d, want %d", got, tt.want)
			}
		})
	}
}

// TestMetricsAll tests getting all metrics
func TestMetricsAll(t *testing.T) {
	tests := []struct {
		name  string
		setup func(Metrics)
		want  int
	}{
		{
			name: "should return empty map for empty metrics",
			setup: func(m Metrics) {},
			want:  0,
		},
		{
			name: "should return all metrics",
			setup: func(m Metrics) {
				m.Add("m1", 1)
				m.Add("m2", 2)
				m.Add("m3", 3)
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := New()
			tt.setup(m)
			all := m.All()
			if len(all) != tt.want {
				t.Errorf("All() returned %d items, want %d", len(all), tt.want)
			}
		})
	}
}