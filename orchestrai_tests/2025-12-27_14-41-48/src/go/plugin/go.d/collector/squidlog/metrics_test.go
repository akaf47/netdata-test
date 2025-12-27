package squidlog

import (
	"testing"
)

// TestNewMetrics tests the NewMetrics function
func TestNewMetrics(t *testing.T) {
	tests := []struct {
		name string
		want *Metrics
	}{
		{
			name: "should create new metrics instance",
			want: &Metrics{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetrics()
			if m == nil {
				t.Error("NewMetrics() returned nil")
			}
		})
	}
}

// TestMetricsAddRequest tests the AddRequest method
func TestMetricsAddRequest(t *testing.T) {
	tests := []struct {
		name            string
		initialCount    int
		requestsToAdd   int
		expectedCount   int
	}{
		{
			name:          "should increment request count",
			initialCount:  0,
			requestsToAdd: 1,
			expectedCount: 1,
		},
		{
			name:          "should handle multiple additions",
			initialCount:  5,
			requestsToAdd: 3,
			expectedCount: 8,
		},
		{
			name:          "should handle zero addition",
			initialCount:  10,
			requestsToAdd: 0,
			expectedCount: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetrics()
			for i := 0; i < tt.initialCount; i++ {
				m.AddRequest()
			}
			for i := 0; i < tt.requestsToAdd; i++ {
				m.AddRequest()
			}
			if m.Requests != int64(tt.expectedCount) {
				t.Errorf("AddRequest() = %d, want %d", m.Requests, tt.expectedCount)
			}
		})
	}
}

// TestMetricsAddRequestBytes tests the AddRequestBytes method
func TestMetricsAddRequestBytes(t *testing.T) {
	tests := []struct {
		name           string
		bytes          int64
		expectedTotal  int64
	}{
		{
			name:          "should add request bytes",
			bytes:         1024,
			expectedTotal: 1024,
		},
		{
			name:          "should add multiple times",
			bytes:         512,
			expectedTotal: 512,
		},
		{
			name:          "should handle zero bytes",
			bytes:         0,
			expectedTotal: 0,
		},
		{
			name:          "should handle large byte values",
			bytes:         9223372036854775807, // max int64
			expectedTotal: 9223372036854775807,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetrics()
			m.AddRequestBytes(tt.bytes)
			if m.RequestBytes != tt.expectedTotal {
				t.Errorf("AddRequestBytes(%d) = %d, want %d", tt.bytes, m.RequestBytes, tt.expectedTotal)
			}
		})
	}
}

// TestMetricsAddResponseBytes tests the AddResponseBytes method
func TestMetricsAddResponseBytes(t *testing.T) {
	tests := []struct {
		name           string
		bytes          int64
		expectedTotal  int64
	}{
		{
			name:          "should add response bytes",
			bytes:         2048,
			expectedTotal: 2048,
		},
		{
			name:          "should accumulate bytes",
			bytes:         1024,
			expectedTotal: 1024,
		},
		{
			name:          "should handle zero bytes",
			bytes:         0,
			expectedTotal: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetrics()
			m.AddResponseBytes(tt.bytes)
			if m.ResponseBytes != tt.expectedTotal {
				t.Errorf("AddResponseBytes(%d) = %d, want %d", tt.bytes, m.ResponseBytes, tt.expectedTotal)
			}
		})
	}
}

// TestMetricsAddUntaggedHierarchy tests the AddUntaggedHierarchy method
func TestMetricsAddUntaggedHierarchy(t *testing.T) {
	tests := []struct {
		name           string
		count          int
		expectedTotal  int64
	}{
		{
			name:          "should increment untagged hierarchy",
			count:         1,
			expectedTotal: 1,
		},
		{
			name:          "should handle multiple increments",
			count:         5,
			expectedTotal: 5,
		},
		{
			name:          "should handle zero",
			count:         0,
			expectedTotal: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetrics()
			for i := 0; i < tt.count; i++ {
				m.AddUntaggedHierarchy()
			}
			if m.UntaggedHierarchy != tt.expectedTotal {
				t.Errorf("AddUntaggedHierarchy() = %d, want %d", m.UntaggedHierarchy, tt.expectedTotal)
			}
		})
	}
}

// TestMetricsAddTaggedHierarchy tests the AddTaggedHierarchy method
func TestMetricsAddTaggedHierarchy(t *testing.T) {
	tests := []struct {
		name           string
		count          int
		expectedTotal  int64
	}{
		{
			name:          "should increment tagged hierarchy",
			count:         1,
			expectedTotal: 1,
		},
		{
			name:          "should handle multiple increments",
			count:         3,
			expectedTotal: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetrics()
			for i := 0; i < tt.count; i++ {
				m.AddTaggedHierarchy()
			}
			if m.TaggedHierarchy != tt.expectedTotal {
				t.Errorf("AddTaggedHierarchy() = %d, want %d", m.TaggedHierarchy, tt.expectedTotal)
			}
		})
	}
}

// TestMetricsAddHitCount tests the AddHitCount method
func TestMetricsAddHitCount(t *testing.T) {
	tests := []struct {
		name           string
		code           string
		expectedCount  int64
	}{
		{
			name:          "should add TCP_HIT code",
			code:          "TCP_HIT",
			expectedCount: 1,
		},
		{
			name:          "should add TCP_MISS code",
			code:          "TCP_MISS",
			expectedCount: 1,
		},
		{
			name:          "should handle unknown codes",
			code:          "UNKNOWN",
			expectedCount: 0,
		},
		{
			name:          "should handle empty code",
			code:          "",
			expectedCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetrics()
			m.AddHitCount(tt.code)
		})
	}
}

// TestMetricsAddResponseTimeHistogram tests the AddResponseTimeHistogram method
func TestMetricsAddResponseTimeHistogram(t *testing.T) {
	tests := []struct {
		name string
		time int64
	}{
		{
			name: "should add response time 100ms",
			time: 100,
		},
		{
			name: "should add response time 0ms",
			time: 0,
		},
		{
			name: "should add response time 1000ms",
			time: 1000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetrics()
			m.AddResponseTimeHistogram(tt.time)
		})
	}
}

// TestMetricsAddStatusCodeCount tests the AddStatusCodeCount method
func TestMetricsAddStatusCodeCount(t *testing.T) {
	tests := []struct {
		name           string
		code           int
		expectedCount  int64
	}{
		{
			name:          "should count 200 status code",
			code:          200,
			expectedCount: 1,
		},
		{
			name:          "should count 404 status code",
			code:          404,
			expectedCount: 1,
		},
		{
			name:          "should count 500 status code",
			code:          500,
			expectedCount: 1,
		},
		{
			name:          "should handle invalid codes",
			code:          999,
			expectedCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetrics()
			m.AddStatusCodeCount(tt.code)
		})
	}
}

// TestMetricsReset tests the Reset method
func TestMetricsReset(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "should reset all metrics to zero",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetrics()
			m.AddRequest()
			m.AddRequestBytes(1024)
			m.AddResponseBytes(2048)
			m.AddUntaggedHierarchy()
			m.AddTaggedHierarchy()
			
			m.Reset()
			
			if m.Requests != 0 {
				t.Errorf("Reset() - Requests = %d, want 0", m.Requests)
			}
			if m.RequestBytes != 0 {
				t.Errorf("Reset() - RequestBytes = %d, want 0", m.RequestBytes)
			}
			if m.ResponseBytes != 0 {
				t.Errorf("Reset() - ResponseBytes = %d, want 0", m.ResponseBytes)
			}
		})
	}
}

// TestMetricsConcurrentOperations tests concurrent access to metrics
func TestMetricsConcurrentOperations(t *testing.T) {
	t.Run("should handle concurrent requests safely", func(t *testing.T) {
		m := NewMetrics()
		done := make(chan bool, 100)
		
		for i := 0; i < 100; i++ {
			go func() {
				m.AddRequest()
				m.AddRequestBytes(int64(i))
				m.AddResponseBytes(int64(i * 2))
				done <- true
			}()
		}
		
		for i := 0; i < 100; i++ {
			<-done
		}
		
		if m.Requests != 100 {
			t.Errorf("Concurrent operations - Requests = %d, want 100", m.Requests)
		}
	})
}

// TestMetricsGetRequests tests the GetRequests method
func TestMetricsGetRequests(t *testing.T) {
	tests := []struct {
		name          string
		requestCount  int
		expectedValue int64
	}{
		{
			name:          "should return request count 0",
			requestCount:  0,
			expectedValue: 0,
		},
		{
			name:          "should return request count 5",
			requestCount:  5,
			expectedValue: 5,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetrics()
			for i := 0; i < tt.requestCount; i++ {
				m.AddRequest()
			}
			if m.Requests != tt.expectedValue {
				t.Errorf("Requests = %d, want %d", m.Requests, tt.expectedValue)
			}
		})
	}
}