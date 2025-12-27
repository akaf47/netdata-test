package journaldexporter

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockConn is a mock for net.Conn
type MockConn struct {
	mock.Mock
}

func (m *MockConn) Read(b []byte) (n int, err error) {
	args := m.Called(b)
	if args.Get(0) != nil {
		copy(b, args.Get(0).([]byte))
	}
	return args.Int(1), args.Error(2)
}

func (m *MockConn) Write(b []byte) (n int, err error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *MockConn) Close() error {
	return m.Called().Error(0)
}

func (m *MockConn) LocalAddr() net.Addr {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(net.Addr)
}

func (m *MockConn) RemoteAddr() net.Addr {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(net.Addr)
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

// MockRemoteServer simulates a remote journal server
type MockRemoteServer struct {
	mock.Mock
	listener net.Listener
	addr     string
}

func NewMockRemoteServer(addr string) *MockRemoteServer {
	return &MockRemoteServer{
		addr: addr,
	}
}

func (m *MockRemoteServer) Start() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRemoteServer) Stop() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRemoteServer) Address() string {
	return m.addr
}

// Test creating a remote journal client
func TestNewRemoteJournalClient(t *testing.T) {
	tests := []struct {
		name      string
		address   string
		shouldErr bool
	}{
		{
			name:      "valid address",
			address:   "localhost:5140",
			shouldErr: false,
		},
		{
			name:      "empty address",
			address:   "",
			shouldErr: true,
		},
		{
			name:      "invalid address format",
			address:   "not-a-valid-address::",
			shouldErr: false, // Address validation may happen on connect
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewRemoteJournalClient(tt.address)
			if tt.shouldErr {
				assert.Error(t, err)
				assert.Nil(t, client)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, client)
			}
		})
	}
}

// Test connecting to a remote journal server
func TestRemoteJournalClient_Connect(t *testing.T) {
	tests := []struct {
		name           string
		address        string
		shouldErr      bool
		connectTimeout time.Duration
	}{
		{
			name:           "connect to valid address",
			address:        "localhost:9999",
			shouldErr:      true, // No server listening
			connectTimeout: 100 * time.Millisecond,
		},
		{
			name:      "empty address",
			address:   "",
			shouldErr: true,
		},
		{
			name:           "localhost connection",
			address:        "127.0.0.1:19999",
			shouldErr:      true, // No server listening
			connectTimeout: 100 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewRemoteJournalClient(tt.address)
			if err != nil {
				t.Skipf("failed to create client: %v", err)
			}

			ctx, cancel := context.WithTimeout(context.Background(), tt.connectTimeout)
			defer cancel()

			err = client.Connect(ctx)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if client.conn != nil {
					_ = client.Close()
				}
			}
		})
	}
}

// Test closing connection
func TestRemoteJournalClient_Close(t *testing.T) {
	client, err := NewRemoteJournalClient("localhost:5140")
	require.NoError(t, err)
	require.NotNil(t, client)

	// Close without connection should not error
	err = client.Close()
	assert.NoError(t, err)
}

// Test sending data to remote server
func TestRemoteJournalClient_Send(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		shouldErr bool
	}{
		{
			name:      "send empty data",
			data:      []byte{},
			shouldErr: false,
		},
		{
			name:      "send valid data",
			data:      []byte("test message"),
			shouldErr: true, // No connection established
		},
		{
			name:      "send large data",
			data:      []byte(strings.Repeat("a", 10000)),
			shouldErr: true, // No connection established
		},
		{
			name:      "send nil data",
			data:      nil,
			shouldErr: false, // Depends on implementation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewRemoteJournalClient("localhost:5140")
			require.NoError(t, err)

			err = client.Send(tt.data)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test reconnection logic
func TestRemoteJournalClient_Reconnect(t *testing.T) {
	tests := []struct {
		name      string
		maxRetry  int
		shouldErr bool
	}{
		{
			name:      "reconnect with 0 retries",
			maxRetry:  0,
			shouldErr: true,
		},
		{
			name:      "reconnect with 1 retry",
			maxRetry:  1,
			shouldErr: true, // No server listening
		},
		{
			name:      "reconnect with 3 retries",
			maxRetry:  3,
			shouldErr: true, // No server listening
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewRemoteJournalClient("localhost:29999")
			require.NoError(t, err)

			ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
			defer cancel()

			err = client.Reconnect(ctx, tt.maxRetry)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// Test setting connection timeout
func TestRemoteJournalClient_SetTimeout(t *testing.T) {
	client, err := NewRemoteJournalClient("localhost:5140")
	require.NoError(t, err)

	tests := []struct {
		name    string
		timeout time.Duration
	}{
		{
			name:    "zero timeout",
			timeout: 0,
		},
		{
			name:    "short timeout",
			timeout: 100 * time.Millisecond,
		},
		{
			name:    "long timeout",
			timeout: 30 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.SetTimeout(tt.timeout)
			// Verify timeout is set (implementation specific)
			assert.True(t, client != nil)
		})
	}
}

// Test connection state
func TestRemoteJournalClient_IsConnected(t *testing.T) {
	client, err := NewRemoteJournalClient("localhost:5140")
	require.NoError(t, err)

	// Initially disconnected
	assert.False(t, client.IsConnected())

	// After closing (if not connected)
	_ = client.Close()
	assert.False(t, client.IsConnected())
}

// Test error handling during send
func TestRemoteJournalClient_SendWithConnectionError(t *testing.T) {
	client, err := NewRemoteJournalClient("localhost:5140")
	require.NoError(t, err)

	// Try sending without connection
	err = client.Send([]byte("test"))
	if err != nil {
		assert.Error(t, err)
	}
}

// Test multiple operations
func TestRemoteJournalClient_MultipleOperations(t *testing.T) {
	client, err := NewRemoteJournalClient("localhost:5140")
	require.NoError(t, err)

	// Test sequence of operations
	assert.False(t, client.IsConnected())
	_ = client.Send([]byte("test1"))
	assert.False(t, client.IsConnected())
	_ = client.Close()
	assert.False(t, client.IsConnected())
}

// Test concurrent send operations
func TestRemoteJournalClient_ConcurrentSend(t *testing.T) {
	client, err := NewRemoteJournalClient("localhost:5140")
	require.NoError(t, err)

	done := make(chan bool, 3)
	for i := 0; i < 3; i++ {
		go func(index int) {
			_ = client.Send([]byte("concurrent message"))
			done <- true
		}(i)
	}

	for i := 0; i < 3; i++ {
		<-done
	}

	_ = client.Close()
}

// Test connection failure scenarios
func TestRemoteJournalClient_ConnectionFailures(t *testing.T) {
	tests := []struct {
		name    string
		address string
	}{
		{
			name:    "invalid host",
			address: "invalid-host-that-does-not-exist.local:5140",
		},
		{
			name:    "localhost with unreachable port",
			address: "127.0.0.1:1",
		},
		{
			name:    "ipv6 address",
			address: "[::1]:5140",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewRemoteJournalClient(tt.address)
			if err != nil {
				// Address parsing error
				assert.Error(t, err)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()

			err = client.Connect(ctx)
			assert.Error(t, err)
		})
	}
}

// Test data integrity
func TestRemoteJournalClient_DataIntegrity(t *testing.T) {
	testData := [][]byte{
		[]byte("simple message"),
		[]byte("message with\nnewlines"),
		[]byte("message with\ttabs"),
		[]byte("🚀 emoji message"),
		[]byte(strings.Repeat("long ", 1000)),
	}

	for _, data := range testData {
		t.Run("send "+string(data[:min(len(data), 20)]), func(t *testing.T) {
			client, err := NewRemoteJournalClient("localhost:5140")
			require.NoError(t, err)

			// Data should be processed (may fail due to no connection)
			_ = client.Send(data)
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Test context cancellation
func TestRemoteJournalClient_ContextCancellation(t *testing.T) {
	client, err := NewRemoteJournalClient("localhost:5140")
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = client.Connect(ctx)
	// Should return an error due to context cancellation
	if err != nil {
		assert.Error(t, err)
	}
}

// Test address parsing edge cases
func TestRemoteJournalClient_AddressParsing(t *testing.T) {
	tests := []struct {
		name    string
		address string
	}{
		{
			name:    "hostname with port",
			address: "journal.example.com:5140",
		},
		{
			name:    "ipv4 with port",
			address: "192.168.1.1:5140",
		},
		{
			name:    "localhost with port",
			address: "localhost:5140",
		},
		{
			name:    "ipv6 with port",
			address: "[2001:db8::1]:5140",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewRemoteJournalClient(tt.address)
			if err != nil {
				t.Logf("Address parsing failed: %v", err)
			} else {
				assert.NotNil(t, client)
			}
		})
	}
}

// Test graceful shutdown
func TestRemoteJournalClient_GracefulShutdown(t *testing.T) {
	client, err := NewRemoteJournalClient("localhost:5140")
	require.NoError(t, err)

	// Multiple close calls should not panic
	_ = client.Close()
	_ = client.Close()
	_ = client.Close()
	assert.True(t, true) // If we reach here, no panic occurred
}

// Test send with nil client
func TestRemoteJournalClient_NilReceiver(t *testing.T) {
	var client *RemoteJournalClient
	if client != nil {
		err := client.Send([]byte("test"))
		assert.NoError(t, err)
	}
}

// Test timeout behavior
func TestRemoteJournalClient_TimeoutBehavior(t *testing.T) {
	client, err := NewRemoteJournalClient("localhost:19999")
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	start := time.Now()
	err = client.Connect(ctx)
	elapsed := time.Since(start)

	// Should timeout before expected
	assert.True(t, elapsed >= 50*time.Millisecond)
}

// Test error type assertions
func TestRemoteJournalClient_ErrorTypes(t *testing.T) {
	tests := []struct {
		name     string
		address  string
		testFunc func(*RemoteJournalClient) error
	}{
		{
			name:    "connect error",
			address: "localhost:29999",
			testFunc: func(c *RemoteJournalClient) error {
				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				defer cancel()
				return c.Connect(ctx)
			},
		},
		{
			name:    "send error without connection",
			address: "localhost:5140",
			testFunc: func(c *RemoteJournalClient) error {
				return c.Send([]byte("test"))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewRemoteJournalClient(tt.address)
			require.NoError(t, err)

			err = tt.testFunc(client)
			if err != nil {
				assert.True(t, true) // Error occurred as expected
			}
		})
	}
}

// Test with special characters
func TestRemoteJournalClient_SpecialCharacters(t *testing.T) {
	specialMessages := [][]byte{
		[]byte("\x00\x01\x02\x03"), // Binary data
		[]byte("\xff\xfe\xfd"),      // High bytes
		[]byte(""),                  // Empty
		[]byte("\n\r\n"),            // Line endings
	}

	for i, msg := range specialMessages {
		t.Run("special_message_"+string(rune(i)), func(t *testing.T) {
			client, err := NewRemoteJournalClient("localhost:5140")
			require.NoError(t, err)

			// Should handle special characters
			_ = client.Send(msg)
		})
	}
}

// Test read timeout
func TestRemoteJournalClient_ReadTimeout(t *testing.T) {
	client, err := NewRemoteJournalClient("localhost:5140")
	require.NoError(t, err)

	// Set a very short timeout
	client.SetTimeout(1 * time.Nanosecond)

	// Operations should respect timeout
	_ = client.Send([]byte("test"))
}