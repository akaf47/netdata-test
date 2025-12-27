package metrix

import (
	"testing"
	"math"
)

// Test histogram creation
func TestNewHistogram(t *testing.T) {
	tests := []struct {
		name string
		id   string
	}{
		{
			name: "should create histogram with valid id",
			id:   "test_histogram",
		},
		{
			name: "should create histogram with empty id",
			id:   "",
		},
		{
			name: "should create histogram with numeric id",
			id:   "histogram_123",
		},
		{
			name: "should create histogram with special characters",
			id:   "hist-_.-",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHistogram(tt.id)
			if h == nil {
				t.Errorf("NewHistogram() returned nil")
			}
			if h.id != tt.id {
				t.Errorf("NewHistogram() id = %v, want %v", h.id, tt.id)
			}
		})
	}
}

// Test histogram initialization
func TestHistogramInit(t *testing.T) {
	h := NewHistogram("test")
	if h.count != 0 {
		t.Errorf("histogram initial count = %v, want 0", h.count)
	}
	if h.sum != 0 {
		t.Errorf("histogram initial sum = %v, want 0", h.sum)
	}
	if h.min != 0 {
		t.Errorf("histogram initial min = %v, want 0", h.min)
	}
	if h.max != 0 {
		t.Errorf("histogram initial max = %v, want 0", h.max)
	}
}

// Test histogram Observe method - single observation
func TestHistogramObserveSingle(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		expected float64
	}{
		{
			name:     "should observe positive value",
			value:    42.5,
			expected: 42.5,
		},
		{
			name:     "should observe zero value",
			value:    0,
			expected: 0,
		},
		{
			name:     "should observe negative value",
			value:    -10.5,
			expected: -10.5,
		},
		{
			name:     "should observe large value",
			value:    1e10,
			expected: 1e10,
		},
		{
			name:     "should observe small value",
			value:    1e-10,
			expected: 1e-10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHistogram("test")
			h.Observe(tt.value)
			
			if h.sum != tt.expected {
				t.Errorf("Observe(%v) resulted in sum %v, want %v", tt.value, h.sum, tt.expected)
			}
			if h.count != 1 {
				t.Errorf("Observe(%v) resulted in count %v, want 1", tt.value, h.count)
			}
			if h.min != tt.expected {
				t.Errorf("Observe(%v) resulted in min %v, want %v", tt.value, h.min, tt.expected)
			}
			if h.max != tt.expected {
				t.Errorf("Observe(%v) resulted in max %v, want %v", tt.value, h.max, tt.expected)
			}
		})
	}
}

// Test histogram Observe method - multiple observations
func TestHistogramObserveMultiple(t *testing.T) {
	h := NewHistogram("test")
	values := []float64{10, 20, 15, 5, 25}
	
	for _, v := range values {
		h.Observe(v)
	}
	
	expectedCount := int64(5)
	expectedSum := 75.0
	expectedMin := 5.0
	expectedMax := 25.0
	
	if h.count != expectedCount {
		t.Errorf("Multiple Observe() resulted in count %v, want %v", h.count, expectedCount)
	}
	if h.sum != expectedSum {
		t.Errorf("Multiple Observe() resulted in sum %v, want %v", h.sum, expectedSum)
	}
	if h.min != expectedMin {
		t.Errorf("Multiple Observe() resulted in min %v, want %v", h.min, expectedMin)
	}
	if h.max != expectedMax {
		t.Errorf("Multiple Observe() resulted in max %v, want %v", h.max, expectedMax)
	}
}

// Test histogram with all equal values
func TestHistogramObserveEqualValues(t *testing.T) {
	h := NewHistogram("test")
	value := 42.0
	
	for i := 0; i < 10; i++ {
		h.Observe(value)
	}
	
	if h.count != 10 {
		t.Errorf("Observe equal values resulted in count %v, want 10", h.count)
	}
	if h.sum != 420.0 {
		t.Errorf("Observe equal values resulted in sum %v, want 420", h.sum)
	}
	if h.min != value {
		t.Errorf("Observe equal values resulted in min %v, want %v", h.min, value)
	}
	if h.max != value {
		t.Errorf("Observe equal values resulted in max %v, want %v", h.max, value)
	}
}

// Test histogram average calculation
func TestHistogramAverage(t *testing.T) {
	h := NewHistogram("test")
	h.Observe(10)
	h.Observe(20)
	h.Observe(30)
	
	expectedAverage := 20.0
	actualAverage := h.Average()
	
	if actualAverage != expectedAverage {
		t.Errorf("Average() = %v, want %v", actualAverage, expectedAverage)
	}
}

// Test histogram average with zero observations
func TestHistogramAverageZeroObservations(t *testing.T) {
	h := NewHistogram("test")
	average := h.Average()
	
	if average != 0 {
		t.Errorf("Average() with no observations = %v, want 0", average)
	}
}

// Test histogram average with single observation
func TestHistogramAverageSingleObservation(t *testing.T) {
	h := NewHistogram("test")
	h.Observe(42)
	
	if h.Average() != 42 {
		t.Errorf("Average() with single observation = %v, want 42", h.Average())
	}
}

// Test histogram Count method
func TestHistogramCount(t *testing.T) {
	h := NewHistogram("test")
	
	if h.Count() != 0 {
		t.Errorf("Count() on new histogram = %v, want 0", h.Count())
	}
	
	h.Observe(10)
	if h.Count() != 1 {
		t.Errorf("Count() after 1 observation = %v, want 1", h.Count())
	}
	
	h.Observe(20)
	if h.Count() != 2 {
		t.Errorf("Count() after 2 observations = %v, want 2", h.Count())
	}
}

// Test histogram Sum method
func TestHistogramSum(t *testing.T) {
	h := NewHistogram("test")
	
	if h.Sum() != 0 {
		t.Errorf("Sum() on new histogram = %v, want 0", h.Sum())
	}
	
	h.Observe(10)
	if h.Sum() != 10 {
		t.Errorf("Sum() after observing 10 = %v, want 10", h.Sum())
	}
	
	h.Observe(15)
	if h.Sum() != 25 {
		t.Errorf("Sum() after observing 15 = %v, want 25", h.Sum())
	}
}

// Test histogram Min method
func TestHistogramMin(t *testing.T) {
	h := NewHistogram("test")
	
	if h.Min() != 0 {
		t.Errorf("Min() on new histogram = %v, want 0", h.Min())
	}
	
	h.Observe(50)
	if h.Min() != 50 {
		t.Errorf("Min() after observing 50 = %v, want 50", h.Min())
	}
	
	h.Observe(30)
	if h.Min() != 30 {
		t.Errorf("Min() after observing 30 = %v, want 30", h.Min())
	}
	
	h.Observe(40)
	if h.Min() != 30 {
		t.Errorf("Min() after observing 40 = %v, want 30", h.Min())
	}
}

// Test histogram Max method
func TestHistogramMax(t *testing.T) {
	h := NewHistogram("test")
	
	if h.Max() != 0 {
		t.Errorf("Max() on new histogram = %v, want 0", h.Max())
	}
	
	h.Observe(10)
	if h.Max() != 10 {
		t.Errorf("Max() after observing 10 = %v, want 10", h.Max())
	}
	
	h.Observe(50)
	if h.Max() != 50 {
		t.Errorf("Max() after observing 50 = %v, want 50", h.Max())
	}
	
	h.Observe(30)
	if h.Max() != 50 {
		t.Errorf("Max() after observing 30 = %v, want 50", h.Max())
	}
}

// Test histogram Reset method
func TestHistogramReset(t *testing.T) {
	h := NewHistogram("test")
	h.Observe(10)
	h.Observe(20)
	h.Observe(30)
	
	h.Reset()
	
	if h.Count() != 0 {
		t.Errorf("Count() after Reset() = %v, want 0", h.Count())
	}
	if h.Sum() != 0 {
		t.Errorf("Sum() after Reset() = %v, want 0", h.Sum())
	}
	if h.Min() != 0 {
		t.Errorf("Min() after Reset() = %v, want 0", h.Min())
	}
	if h.Max() != 0 {
		t.Errorf("Max() after Reset() = %v, want 0", h.Max())
	}
}

// Test histogram Type method
func TestHistogramType(t *testing.T) {
	h := NewHistogram("test")
	histogramType := h.Type()
	if histogramType != "histogram" {
		t.Errorf("Type() = %v, want 'histogram'", histogramType)
	}
}

// Test histogram ID method
func TestHistogramID(t *testing.T) {
	id := "my_histogram_id"
	h := NewHistogram(id)
	if h.ID() != id {
		t.Errorf("ID() = %v, want %v", h.ID(), id)
	}
}

// Test histogram with negative values tracking
func TestHistogramNegativeValues(t *testing.T) {
	h := NewHistogram("test")
	h.Observe(-10)
	h.Observe(-5)
	h.Observe(5)
	h.Observe(10)
	
	expectedSum := 0.0
	if h.Sum() != expectedSum {
		t.Errorf("Sum with negative values = %v, want %v", h.Sum(), expectedSum)
	}
	
	expectedMin := -10.0
	if h.Min() != expectedMin {
		t.Errorf("Min with negative values = %v, want %v", h.Min(), expectedMin)
	}
	
	expectedMax := 10.0
	if h.Max() != expectedMax {
		t.Errorf("Max with negative values = %v, want %v", h.Max(), expectedMax)
	}
}

// Test histogram with zero values
func TestHistogramZeroValues(t *testing.T) {
	h := NewHistogram("test")
	for i := 0; i < 5; i++ {
		h.Observe(0)
	}
	
	if h.Count() != 5 {
		t.Errorf("Count with zero values = %v, want 5", h.Count())
	}
	if h.Sum() != 0 {
		t.Errorf("Sum with zero values = %v, want 0", h.Sum())
	}
	if h.Min() != 0 {
		t.Errorf("Min with zero values = %v, want 0", h.Min())
	}
	if h.Max() != 0 {
		t.Errorf("Max with zero values = %v, want 0", h.Max())
	}
}

// Test histogram with extreme values
func TestHistogramExtremeValues(t *testing.T) {
	h := NewHistogram("test")
	
	maxFloat := 1.7976931348623157e+308
	minFloat := -1.7976931348623157e+308
	
	h.Observe(maxFloat)
	h.Observe(minFloat)
	h.Observe(0)
	
	if h.Max() != maxFloat {
		t.Errorf("Max with extreme values = %v, want %v", h.Max(), maxFloat)
	}
	if h.Min() != minFloat {
		t.Errorf("Min with extreme values = %v, want %v", h.Min(), minFloat)
	}
}

// Test histogram concurrent observations
func TestHistogramConcurrentObserve(t *testing.T) {
	h := NewHistogram("test")
	done := make(chan bool, 100)
	
	for i := 0; i < 100; i++ {
		go func(val float64) {
			h.Observe(val)
			done <- true
		}(float64(i))
	}
	
	for i := 0; i < 100; i++ {
		<-done
	}
	
	expectedCount := int64(100)
	if h.Count() != expectedCount {
		t.Errorf("Concurrent Observe() count = %v, want %v", h.Count(), expectedCount)
	}
}

// Test histogram min/max tracking with sequence
func TestHistogramSequentialMinMax(t *testing.T) {
	h := NewHistogram("test")
	
	// Ascending sequence
	for i := 1; i <= 10; i++ {
		h.Observe(float64(i))
	}
	
	if h.Min() != 1 {
		t.Errorf("Min() with ascending sequence = %v, want 1", h.Min())
	}
	if h.Max() != 10 {
		t.Errorf("Max() with ascending sequence = %v, want 10", h.Max())
	}
}

// Test histogram min/max tracking descending sequence
func TestHistogramDescendingSequence(t *testing.T) {
	h := NewHistogram("test")
	
	// Descending sequence
	for i := 10; i >= 1; i-- {
		h.Observe(float64(i))
	}
	
	if h.Min() != 1 {
		t.Errorf("Min() with descending sequence = %v, want 1", h.Min())
	}
	if h.Max() != 10 {
		t.Errorf("Max() with descending sequence = %v, want 10", h.Max())
	}
}

// Test histogram precision with very small differences
func TestHistogramPrecision(t *testing.T) {
	h := NewHistogram("test")
	h.Observe(0.1)
	h.Observe(0.2)
	h.Observe(0.3)
	
	expectedSum := 0.6
	if math.Abs(h.Sum()-expectedSum) > 1e-10 {
		t.Errorf("Sum() precision = %v, want %v", h.Sum(), expectedSum)
	}
}

// Test histogram large dataset
func TestHistogramLargeDataset(t *testing.T) {
	h := NewHistogram("test")
	
	for i := 1; i <= 10000; i++ {
		h.Observe(float64(i))
	}
	
	expectedCount := int64(10000)
	expectedSum := 50005000.0 // sum of 1 to 10000
	expectedMin := 1.0
	expectedMax := 10000.0
	expectedAverage := 5000.5
	
	if h.Count() != expectedCount {
		t.Errorf("Large dataset Count() = %v, want %v", h.Count(), expectedCount)
	}
	if h.Sum() != expectedSum {
		t.Errorf("Large dataset Sum() = %v, want %v", h.Sum(), expectedSum)
	}
	if h.Min() != expectedMin {
		t.Errorf("Large dataset Min() = %v, want %v", h.Min(), expectedMin)
	}
	if h.Max() != expectedMax {
		t.Errorf("Large dataset Max() = %v, want %v", h.Max(), expectedMax)
	}
	if h.Average() != expectedAverage {
		t.Errorf("Large dataset Average() = %v, want %v", h.Average(), expectedAverage)
	}
}

// Test histogram reset after observations
func TestHistogramResetAfterObservations(t *testing.T) {
	h := NewHistogram("test")
	
	h.Observe(100)
	h.Reset()
	
	h.Observe(5)
	
	if h.Count() != 1 {
		t.Errorf("Count() after reset and new observation = %v, want 1", h.Count())
	}
	if h.Sum() != 5 {
		t.Errorf("Sum() after reset and new observation = %v, want 5", h.Sum())
	}
	if h.Min() != 5 {
		t.Errorf("Min() after reset and new observation = %v, want 5", h.Min())
	}
	if h.Max() != 5 {
		t.Errorf("Max() after reset and new observation = %v, want 5", h.Max())
	}
}

// Test histogram very large count
func TestHistogramVeryLargeCount(t *testing.T) {
	h := NewHistogram("test")
	
	for i := 0; i < 1000000; i++ {
		h.Observe(1)
	}
	
	if h.Count() != 1000000 {
		t.Errorf("Very large count = %v, want 1000000", h.Count())
	}
	if h.Sum() != 1000000 {
		t.Errorf("Sum with very large count = %v, want 1000000", h.Sum())
	}
	if h.Average() != 1 {
		t.Errorf("Average with very large count = %v, want 1", h.Average())
	}
}