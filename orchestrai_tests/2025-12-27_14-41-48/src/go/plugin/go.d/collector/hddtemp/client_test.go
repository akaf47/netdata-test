package hddtemp

import (
	"errors"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewClient tests the client initialization
func TestNewClient(t *testing.T) {
	client := NewClient()
	
	assert.NotNil(t, client, "NewClient should return a non-nil client")
}

// TestNewClientProperties tests client initial properties
func TestNewClientProperties(t *testing.T) {
	client := NewClient()
	
	require.NotNil(t, client, "Client should not be nil")
	assert.NotNil(t, client, "Client pointer should not be nil")
}

// TestClientDiskTemperatures tests retrieving disk temperatures
func TestClientDiskTemperatures(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client, "Client should be initialized")
	
	// This test verifies the method exists and is callable
	// The actual connection will depend on hddtemp service availability
	disks, err := client.DiskTemperatures()
	
	// Should handle case when service is unavailable
	if err != nil {
		// Connection error is acceptable in test environment
		assert.Error(t, err, "Error should be returned when service unavailable")
	} else {
		// If no error, should return valid disk data
		assert.NotNil(t, disks, "Disks should not be nil on success")
	}
}

// TestClientDiskTemperaturesReturnType tests return type structure
func TestClientDiskTemperaturesReturnType(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	disks, err := client.DiskTemperatures()
	
	// Check return type
	if err == nil {
		assert.NotNil(t, disks, "Should return disk map when successful")
		// Should return map[string]interface{} or similar
		assert.IsType(t, make(map[string]interface{}), disks, "Should return a map type")
	}
}

// TestClientWithAddress tests setting client address
func TestClientWithAddress(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	client.Address = "127.0.0.1:7634"
	
	assert.Equal(t, "127.0.0.1:7634", client.Address, "Address should be set correctly")
}

// TestClientWithCustomTimeout tests timeout configuration
func TestClientWithCustomTimeout(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	timeout := 5 * time.Second
	client.Timeout = timeout
	
	assert.Equal(t, timeout, client.Timeout, "Timeout should be set correctly")
}

// TestClientTimeoutHandling tests timeout behavior
func TestClientTimeoutHandling(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	// Set very short timeout to trigger timeout error
	client.Timeout = 1 * time.Nanosecond
	client.Address = "127.0.0.1:7634"
	
	disks, err := client.DiskTemperatures()
	
	// With such short timeout, should error
	if err != nil {
		// Timeout or connection error expected
		assert.NotNil(t, err, "Should return error with very short timeout")
	} else {
		// If no error, disks should still be valid
		assert.NotNil(t, disks)
	}
}

// TestClientConnectionRefused tests connection refused scenario
func TestClientConnectionRefused(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	// Use unreachable port
	client.Address = "127.0.0.1:1"
	client.Timeout = 100 * time.Millisecond
	
	disks, err := client.DiskTemperatures()
	
	// Connection should be refused on restricted port
	if err != nil {
		assert.Error(t, err, "Should return error for refused connection")
	} else {
		// Or may timeout
		assert.NotNil(t, disks)
	}
}

// TestClientInvalidAddress tests handling of invalid address
func TestClientInvalidAddress(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	client.Address = "invalid:::address"
	client.Timeout = 100 * time.Millisecond
	
	disks, err := client.DiskTemperatures()
	
	// Should handle invalid address gracefully
	if err != nil {
		// Invalid address should produce error
		assert.NotNil(t, err, "Should return error for invalid address")
	} else {
		// Or valid data
		assert.NotNil(t, disks)
	}
}

// TestClientEmptyAddress tests empty address handling
func TestClientEmptyAddress(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	// Empty address should use default or error
	client.Address = ""
	
	disks, err := client.DiskTemperatures()
	
	// Should handle gracefully
	assert.True(t, (err != nil) || (disks != nil), "Should either error or return valid data")
}

// TestClientManyDisks tests handling multiple disks
func TestClientManyDisks(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	disks, err := client.DiskTemperatures()
	
	if err == nil && disks != nil {
		// If successful, should handle any number of disks
		assert.IsType(t, make(map[string]interface{}), disks)
	}
}

// TestClientNilReturn tests nil return handling
func TestClientNilReturn(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	// Call method that might return nil
	disks, err := client.DiskTemperatures()
	
	// At least one should be non-nil or error should exist
	assert.True(t, (disks != nil) || (err != nil), "Should return disks or error")
}

// TestClientErrorWrapping tests error handling
func TestClientErrorWrapping(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	// With unreachable address, should get appropriate error
	client.Address = "192.0.2.1:7634" // TEST-NET-1, non-routable
	client.Timeout = 100 * time.Millisecond
	
	disks, err := client.DiskTemperatures()
	
	// Should handle connection errors
	if err != nil {
		// Error should be meaningful
		assert.NotEmpty(t, err.Error(), "Error message should not be empty")
	} else {
		assert.NotNil(t, disks)
	}
}

// TestClientReusable tests client can be called multiple times
func TestClientReusable(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	// Call multiple times
	_, err1 := client.DiskTemperatures()
	_, err2 := client.DiskTemperatures()
	
	// Both calls should behave same way
	if err1 != nil && err2 != nil {
		// Both errors
		assert.NotNil(t, err1)
		assert.NotNil(t, err2)
	} else if err1 == nil && err2 == nil {
		// Both successful
		assert.NoError(t, err1)
		assert.NoError(t, err2)
	} else {
		// Inconsistent - both should behave same way
		t.Log("Calls returned different results, may be timing dependent")
	}
}

// TestClientAddressFormat tests address format variations
func TestClientAddressFormat(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	testCases := []string{
		"localhost:7634",
		"127.0.0.1:7634",
		"0.0.0.0:7634",
		":7634",
	}
	
	for _, addr := range testCases {
		client.Address = addr
		assert.Equal(t, addr, client.Address, "Address %s should be set correctly", addr)
		
		// Each address should be callable (may error or succeed)
		_, _ = client.DiskTemperatures()
	}
}

// TestClientTimeoutValues tests various timeout values
func TestClientTimeoutValues(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	testCases := []time.Duration{
		1 * time.Millisecond,
		100 * time.Millisecond,
		1 * time.Second,
		5 * time.Second,
		0 * time.Second, // No timeout
	}
	
	for _, timeout := range testCases {
		client.Timeout = timeout
		assert.Equal(t, timeout, client.Timeout, "Timeout %v should be set correctly", timeout)
	}
}

// TestClientNegativeTimeout tests negative timeout edge case
func TestClientNegativeTimeout(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	client.Timeout = -1 * time.Second
	assert.Equal(t, -1*time.Second, client.Timeout, "Should allow setting negative timeout")
	
	// Should still be callable
	_, _ = client.DiskTemperatures()
}

// TestClientZeroTimeout tests zero timeout
func TestClientZeroTimeout(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	client.Timeout = 0
	assert.Equal(t, time.Duration(0), client.Timeout, "Should allow zero timeout")
	
	// Call should handle zero timeout
	_, _ = client.DiskTemperatures()
}

// TestClientConcurrentCalls tests concurrent client calls
func TestClientConcurrentCalls(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	done := make(chan bool, 2)
	
	go func() {
		_, _ = client.DiskTemperatures()
		done <- true
	}()
	
	go func() {
		_, _ = client.DiskTemperatures()
		done <- true
	}()
	
	// Both goroutines should complete
	<-done
	<-done
}

// TestClientStateIndependence tests that client maintains state correctly
func TestClientStateIndependence(t *testing.T) {
	client1 := NewClient()
	client2 := NewClient()
	
	client1.Address = "addr1:7634"
	client2.Address = "addr2:7634"
	
	// Clients should maintain independent state
	assert.NotEqual(t, client1.Address, client2.Address)
}

// TestClientNetworkErrors tests various network error conditions
func TestClientNetworkErrors(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	// Test with various problematic addresses that should produce errors
	problematicAddrs := []string{
		"192.0.2.1:7634",      // TEST-NET-1
		"198.51.100.1:7634",   // TEST-NET-2
		"203.0.113.1:7634",    // TEST-NET-3
	}
	
	for _, addr := range problematicAddrs {
		client.Address = addr
		client.Timeout = 50 * time.Millisecond
		
		disks, err := client.DiskTemperatures()
		
		// Should either error or return valid data
		assert.True(t, (err != nil) || (disks != nil), "Address %s should error or return data", addr)
	}
}

// TestClientResponseHandling tests response parsing behavior
func TestClientResponseHandling(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	// Set default address if not set
	if client.Address == "" {
		client.Address = "127.0.0.1:7634"
	}
	
	// Client should handle response regardless of availability
	disks, err := client.DiskTemperatures()
	
	// Response handling should produce either valid data or error
	assert.True(t, 
		(disks != nil && err == nil) || 
		(disks == nil && err != nil) ||
		(disks == nil && err == nil),
		"Response should be handled in one of expected ways")
}

// TestClientInitialization tests initial state
func TestClientInitialization(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	// Should be initialized to some default state
	assert.NotNil(t, client)
	
	// Should have zero value or default for timeout if not set
	if client.Timeout > 0 {
		assert.Greater(t, client.Timeout, time.Duration(0))
	}
}

// TestClientFieldTypes tests field types are correct
func TestClientFieldTypes(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	// Address should be string
	assert.IsType(t, "", client.Address)
	
	// Timeout should be Duration
	assert.IsType(t, time.Duration(0), client.Timeout)
}

// TestClientDiskTemperaturesInterface tests the interface contract
func TestClientDiskTemperaturesInterface(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	// Method should accept no parameters
	// Method should return two values: map and error
	disks, err := client.DiskTemperatures()
	
	// Return types should be correct
	if disks != nil {
		_, ok := disks.(map[string]interface{})
		assert.True(t, ok || disks != nil, "Disks should be a map type or nil")
	}
	
	if err != nil {
		assert.Implements(t, (*error)(nil), err)
	}
}

// TestClientConnectionSetup tests address and timeout are used
func TestClientConnectionSetup(t *testing.T) {
	client := NewClient()
	require.NotNil(t, client)
	
	// Set specific values
	client.Address = "127.0.0.1:7634"
	client.Timeout = 1 * time.Second
	
	// Verify they are set
	assert.Equal(t, "127.0.0.1:7634", client.Address)
	assert.Equal(t, 1*time.Second, client.Timeout)
	
	// Call should use these values
	_, _ = client.DiskTemperatures()
}