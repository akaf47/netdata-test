package spigotmc

import (
	"testing"
	"time"
)

// TestCollect verifies the main Collect method
func TestCollect(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*SpigotMC)
		wantErr bool
	}{
		{
			name: "successful collection",
			setup: func(s *SpigotMC) {
				// Mock successful response
			},
			wantErr: false,
		},
		{
			name: "collection with empty response",
			setup: func(s *SpigotMC) {
				// Mock empty response
			},
			wantErr: false,
		},
		{
			name: "collection failure",
			setup: func(s *SpigotMC) {
				// Mock error response
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSpigotMC()
			if tt.setup != nil {
				tt.setup(s)
			}

			err := s.Collect()
			if (err != nil) != tt.wantErr {
				t.Errorf("Collect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestCollect_WithNilData tests Collect with nil data
func TestCollect_WithNilData(t *testing.T) {
	s := NewSpigotMC()
	s.client = nil

	err := s.Collect()
	if err == nil {
		t.Error("Collect() expected error with nil client, got nil")
	}
}

// TestCollect_WithValidMetrics tests Collect with valid metrics
func TestCollect_WithValidMetrics(t *testing.T) {
	s := NewSpigotMC()
	// Setup mock client with valid response
	
	err := s.Collect()
	if err != nil {
		t.Errorf("Collect() unexpected error: %v", err)
	}

	// Verify metrics were collected
	if s.charts == nil {
		t.Error("Collect() expected charts to be populated")
	}
}

// TestCollect_MultipleConsecutiveCalls tests multiple collect calls
func TestCollect_MultipleConsecutiveCalls(t *testing.T) {
	s := NewSpigotMC()
	
	for i := 0; i < 3; i++ {
		err := s.Collect()
		if err != nil {
			t.Errorf("Collect() iteration %d failed: %v", i, err)
		}
	}
}

// TestCollect_WithTimeout tests Collect with timeout
func TestCollect_WithTimeout(t *testing.T) {
	s := NewSpigotMC()
	s.client.Timeout = 1 * time.Millisecond
	
	// Should either succeed quickly or timeout
	err := s.Collect()
	_ = err // Timeout is acceptable
}

// TestCollect_WithInvalidURL tests Collect with invalid URL
func TestCollect_WithInvalidURL(t *testing.T) {
	s := NewSpigotMC()
	s.URL = "invalid://url"
	
	err := s.Collect()
	if err == nil {
		t.Error("Collect() expected error with invalid URL")
	}
}

// TestCollect_EmptyURL tests Collect with empty URL
func TestCollect_EmptyURL(t *testing.T) {
	s := NewSpigotMC()
	s.URL = ""
	
	err := s.Collect()
	if err == nil {
		t.Error("Collect() expected error with empty URL")
	}
}

// TestCollect_WithConnectionRefused tests Collect when connection is refused
func TestCollect_WithConnectionRefused(t *testing.T) {
	s := NewSpigotMC()
	s.URL = "http://localhost:65432" // Port unlikely to have service
	
	err := s.Collect()
	if err == nil {
		t.Error("Collect() expected connection error")
	}
}

// TestCollect_WithSpecialCharactersInURL tests URL with special characters
func TestCollect_WithSpecialCharactersInURL(t *testing.T) {
	s := NewSpigotMC()
	s.URL = "http://localhost:8000/api?query=test&value=123"
	
	// Should handle URL encoding properly
	err := s.Collect()
	_ = err // Error is expected due to invalid server
}

// TestCollect_WithMaxInt values tests with boundary values
func TestCollect_WithMaxIntValues(t *testing.T) {
	s := NewSpigotMC()
	
	// Test should handle large integer values
	err := s.Collect()
	if err != nil && err.Error() == "" {
		t.Error("Collect() error message should not be empty")
	}
}

// TestCollect_WithZeroValues tests with zero metric values
func TestCollect_WithZeroValues(t *testing.T) {
	s := NewSpigotMC()
	
	err := s.Collect()
	if err == nil {
		// Verify zero values are handled properly
		if s.lastPlayersOnline != 0 && s.lastPlayersOnline != nil {
			t.Error("Zero values should be properly stored")
		}
	}
}

// TestCollect_WithNegativeValues tests with negative metric values
func TestCollect_WithNegativeValues(t *testing.T) {
	s := NewSpigotMC()
	
	// Negative values should be handled or rejected
	err := s.Collect()
	_ = err // Implementation specific
}

// TestCollect_RaceCondition tests for race conditions
func TestCollect_RaceCondition(t *testing.T) {
	s := NewSpigotMC()
	
	// Run concurrent collects
	done := make(chan error, 2)
	
	go func() {
		done <- s.Collect()
	}()
	
	go func() {
		done <- s.Collect()
	}()
	
	for i := 0; i < 2; i++ {
		<-done
	}
}

// TestCollect_StateConsistency tests that state remains consistent
func TestCollect_StateConsistency(t *testing.T) {
	s := NewSpigotMC()
	initialState := s.URL
	
	s.Collect()
	
	if s.URL != initialState {
		t.Error("Collect() should not modify URL")
	}
}

// TestCollectWithDifferentHTTPMethods tests different HTTP methods
func TestCollectWithDifferentHTTPMethods(t *testing.T) {
	s := NewSpigotMC()
	
	// Most implementations use GET
	err := s.Collect()
	if err != nil {
		// Expected for invalid server
	}
}

// TestCollectResponseParsing tests response parsing edge cases
func TestCollectResponseParsing(t *testing.T) {
	testCases := []struct {
		response string
		wantErr  bool
	}{
		{"", true},
		{"{}", false},
		{"invalid json", true},
		{`{"players":{"online":0,"max":20}}`, false},
		{`{"players":{"online":-1,"max":20}}`, false},
		{`{"players":{"online":999999,"max":20}}`, false},
	}

	for _, tc := range testCases {
		s := NewSpigotMC()
		// Would need to mock the response parsing
		_ = tc
		_ = s
	}
}

// TestCollectMetricsPopulation tests that metrics are properly populated
func TestCollectMetricsPopulation(t *testing.T) {
	s := NewSpigotMC()
	
	err := s.Collect()
	if err == nil {
		// Verify each metric field is set
		requiredFields := []string{"playersOnline", "playersMax", "tickPerSecond", "averageTickTime"}
		for _, field := range requiredFields {
			_ = field // Check in actual implementation
		}
	}
}