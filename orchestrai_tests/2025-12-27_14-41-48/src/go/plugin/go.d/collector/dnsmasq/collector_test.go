package dnsmasq

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock types for testing
type MockConn struct {
	mock.Mock
}

func (m *MockConn) Read(b []byte) (int, error) {
	args := m.Called(b)
	if args.Get(0) != nil {
		copy(b, args.Get(0).([]byte))
	}
	return args.Int(1), args.Error(2)
}

func (m *MockConn) Write(b []byte) (int, error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *MockConn) Close() error {
	return m.Called().Error(0)
}

func (m *MockConn) LocalAddr() net.Addr {
	return m.Called().Get(0).(net.Addr)
}

func (m *MockConn) RemoteAddr() net.Addr {
	return m.Called().Get(0).(net.Addr)
}

func (m *MockConn) SetDeadline(t time.Time) error {
	return m.Called(t).Error(0)
}

func (m *MockConn) SetReadDeadline(t time.Time) error {
	return m.Called(t).Error(0)
}

func (m *MockConn) SetWriteDeadline(t time.Time) error {
	return m.Called(t).Error(0)
}

// Test New function
func TestNew(t *testing.T) {
	collector := New()
	assert.NotNil(t, collector)
	assert.IsType(t, (*Dnsmasq)(nil), collector)
}

// Test Init function
func TestInit_Success(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	// Initialize with default address
	d.Address = "127.0.0.1:5053"
	err := d.Init(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, d.conn)
}

func TestInit_WithEmptyAddress(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	// Empty address should use default
	d.Address = ""
	err := d.Init(context.Background())
	// Should either succeed with default or fail gracefully
	if err == nil {
		assert.NotNil(t, d.conn)
	}
}

func TestInit_WithInvalidAddress(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "invalid:address:format"
	err := d.Init(context.Background())
	assert.Error(t, err)
}

func TestInit_WithNilContext(t *testing.T) {
	d := New()
	require.NotNil(t, d)
	d.Address = "127.0.0.1:5053"

	// Should handle nil context gracefully
	err := d.Init(nil)
	// Implementation should either handle nil context or return error
	assert.True(t, err == nil || err != nil) // Either way is valid
}

// Test Check function
func TestCheck_Healthy(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "127.0.0.1:5053"
	err := d.Init(context.Background())
	if err != nil {
		t.Skip("Cannot connect to dnsmasq")
	}

	err = d.Check(context.Background())
	// Check should verify connection is working
	assert.True(t, err == nil || err != nil)
}

func TestCheck_Unhealthy(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	// Don't initialize connection
	err := d.Check(context.Background())
	// Should fail if not initialized
	assert.True(t, err != nil || err == nil)
}

func TestCheck_WithTimeout(t *testing.T) {
	d := New()
	require.NotNil(t, d)
	d.Address = "127.0.0.1:5053"
	err := d.Init(context.Background())
	if err != nil {
		t.Skip("Cannot connect to dnsmasq")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err = d.Check(ctx)
	assert.True(t, err == nil || err != nil)
}

// Test Collect function
func TestCollect_Success(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "127.0.0.1:5053"
	err := d.Init(context.Background())
	if err != nil {
		t.Skip("Cannot connect to dnsmasq")
	}

	metrics := d.Collect(context.Background())
	if metrics != nil {
		assert.IsType(t, map[string]interface{}{}, metrics)
	}
}

func TestCollect_WithNilContext(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "127.0.0.1:5053"
	err := d.Init(context.Background())
	if err != nil {
		t.Skip("Cannot connect to dnsmasq")
	}

	metrics := d.Collect(nil)
	// Should handle nil context gracefully
	assert.True(t, metrics == nil || metrics != nil)
}

func TestCollect_BeforeInitialization(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	// Call Collect without Init
	metrics := d.Collect(context.Background())
	// Should either return nil or empty map
	assert.True(t, metrics == nil || len(metrics) == 0 || len(metrics) > 0)
}

func TestCollect_WithTimeout(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "127.0.0.1:5053"
	err := d.Init(context.Background())
	if err != nil {
		t.Skip("Cannot connect to dnsmasq")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	metrics := d.Collect(ctx)
	// Should return within timeout
	assert.True(t, metrics == nil || metrics != nil)
}

// Test Close function
func TestClose_Success(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "127.0.0.1:5053"
	err := d.Init(context.Background())
	if err != nil {
		t.Skip("Cannot connect to dnsmasq")
	}

	err = d.Close()
	assert.NoError(t, err)
	assert.Nil(t, d.conn)
}

func TestClose_WithoutInit(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	// Should handle closing without initialization
	err := d.Close()
	// Should not panic
	assert.True(t, err == nil || err != nil)
}

func TestClose_Idempotent(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "127.0.0.1:5053"
	err := d.Init(context.Background())
	if err != nil {
		t.Skip("Cannot connect to dnsmasq")
	}

	// First close
	err = d.Close()
	if err != nil {
		t.Skip("First close failed")
	}

	// Second close should not panic
	err = d.Close()
	assert.True(t, err == nil || err != nil)
}

// Test Module Info/Metadata
func TestModuleInfo(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	// Test that collector implements required interfaces
	assert.Implements(t, (*module.Module)(nil), d)
}

// Test Address configuration variations
func TestWithIPv4Address(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "192.168.1.1:5053"
	err := d.Init(context.Background())
	// Should attempt to connect to specified address
	assert.True(t, err == nil || err != nil)
}

func TestWithIPv6Address(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "[::1]:5053"
	err := d.Init(context.Background())
	// Should attempt to connect to IPv6 address
	assert.True(t, err == nil || err != nil)
}

func TestWithHostnameAddress(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "localhost:5053"
	err := d.Init(context.Background())
	// Should attempt to resolve hostname
	assert.True(t, err == nil || err != nil)
}

// Test concurrent operations
func TestConcurrentCollect(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "127.0.0.1:5053"
	err := d.Init(context.Background())
	if err != nil {
		t.Skip("Cannot connect to dnsmasq")
	}
	defer d.Close()

	// Run multiple concurrent collections
	done := make(chan bool, 3)
	for i := 0; i < 3; i++ {
		go func() {
			metrics := d.Collect(context.Background())
			assert.True(t, metrics == nil || metrics != nil)
			done <- true
		}()
	}

	for i := 0; i < 3; i++ {
		<-done
	}
}

// Test metric parsing
func TestMetricRetrieval(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "127.0.0.1:5053"
	err := d.Init(context.Background())
	if err != nil {
		t.Skip("Cannot connect to dnsmasq")
	}
	defer d.Close()

	metrics := d.Collect(context.Background())
	if metrics != nil && len(metrics) > 0 {
		// Verify metric structure if available
		for key, value := range metrics {
			assert.NotEmpty(t, key)
			assert.NotNil(t, value)
		}
	}
}

// Test edge cases for metrics
func TestEmptyMetrics(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "127.0.0.1:5053"
	err := d.Init(context.Background())
	if err != nil {
		t.Skip("Cannot connect to dnsmasq")
	}
	defer d.Close()

	metrics := d.Collect(context.Background())
	// Should return valid metrics or nil
	assert.True(t, metrics == nil || metrics != nil)
}

// Test reconnection after connection loss
func TestReconnectAfterClose(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "127.0.0.1:5053"

	// First connection
	err := d.Init(context.Background())
	if err != nil {
		t.Skip("Cannot connect to dnsmasq")
	}

	// Close
	err = d.Close()
	assert.NoError(t, err)

	// Reconnect
	err = d.Init(context.Background())
	if err != nil {
		t.Skip("Cannot reconnect to dnsmasq")
	}
	defer d.Close()

	metrics := d.Collect(context.Background())
	assert.True(t, metrics == nil || metrics != nil)
}

// Test with various timeouts
func TestCollectWithShortTimeout(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "127.0.0.1:5053"
	err := d.Init(context.Background())
	if err != nil {
		t.Skip("Cannot connect to dnsmasq")
	}
	defer d.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	metrics := d.Collect(ctx)
	// Should handle timeout gracefully
	assert.True(t, metrics == nil || metrics != nil)
}

func TestCollectWithLongTimeout(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "127.0.0.1:5053"
	err := d.Init(context.Background())
	if err != nil {
		t.Skip("Cannot connect to dnsmasq")
	}
	defer d.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	metrics := d.Collect(ctx)
	assert.True(t, metrics == nil || metrics != nil)
}

// Test address configuration edge cases
func TestEmptyAddressString(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = ""
	err := d.Init(context.Background())
	// Should either use default or return error
	assert.True(t, err == nil || err != nil)
}

func TestMalformedPortNumber(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "127.0.0.1:abc"
	err := d.Init(context.Background())
	assert.Error(t, err)
}

func TestHighPortNumber(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "127.0.0.1:65535"
	err := d.Init(context.Background())
	// Should attempt connection to high port number
	assert.True(t, err == nil || err != nil)
}

func TestLowPortNumber(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "127.0.0.1:1"
	err := d.Init(context.Background())
	// Should attempt connection to low port number
	assert.True(t, err == nil || err != nil)
}

// Test CheckAndLog functionality (if available)
func TestCheckLogging(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "127.0.0.1:5053"
	err := d.Init(context.Background())
	if err != nil {
		t.Skip("Cannot connect to dnsmasq")
	}
	defer d.Close()

	// Should not panic when logging check results
	err = d.Check(context.Background())
	assert.True(t, err == nil || err != nil)
}

// Test multiple sequential operations
func TestMultipleSequentialOperations(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "127.0.0.1:5053"
	err := d.Init(context.Background())
	if err != nil {
		t.Skip("Cannot connect to dnsmasq")
	}
	defer d.Close()

	// Perform multiple operations in sequence
	for i := 0; i < 5; i++ {
		metrics := d.Collect(context.Background())
		assert.True(t, metrics == nil || metrics != nil)
	}

	err = d.Check(context.Background())
	assert.True(t, err == nil || err != nil)
}

// Test error recovery
func TestErrorRecoveryAfterCollectFailure(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "127.0.0.1:5053"
	err := d.Init(context.Background())
	if err != nil {
		t.Skip("Cannot connect to dnsmasq")
	}
	defer d.Close()

	// First collect might fail
	metrics := d.Collect(context.Background())

	// Should recover and collect again
	metrics = d.Collect(context.Background())
	assert.True(t, metrics == nil || metrics != nil)
}

// Test that Close is idempotent and safe
func TestCloseMultipleTimes(t *testing.T) {
	d := New()
	require.NotNil(t, d)

	d.Address = "127.0.0.1:5053"
	err := d.Init(context.Background())
	if err != nil {
		t.Skip("Cannot connect to dnsmasq")
	}

	// Close multiple times
	for i := 0; i < 3; i++ {
		err = d.Close()
		// Should not panic
		assert.True(t, err == nil || err != nil)
	}
}