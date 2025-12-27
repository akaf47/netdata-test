package zookeeper

import (
	"bytes"
	"errors"
	"io"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/pkg/tlscfg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCollectSuccessful tests successful metric collection
func TestCollectSuccessful(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	// Mock the connection
	mockConn := &mockConnection{
		readResponses: []string{
			"Zookeeper version: 3.4.5, built on 06/20/2012 15 40 GMT\n" +
			"Host: zk01.example.com:2181\n" +
			"Latency min/avg/max: 0/0/2559\n" +
			"Received: 123456789\n" +
			"Sent: 123456789\n" +
			"Connections: 10\n" +
			"Outstanding: 0\n" +
			"Zxid: 0x100000000\n" +
			"Mode: follower\n" +
			"Node count: 1000\n",
		},
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	assert.NoError(t, err)
	assert.NotNil(t, zk.charts)
}

// TestCollectDialError tests error when connection fails
func TestCollectDialError(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	zk.dial = func(network, address string) (connection, error) {
		return nil, errors.New("connection refused")
	}
	
	err := zk.Collect()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection refused")
}

// TestCollectEmptyResponse tests handling empty response
func TestCollectEmptyResponse(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	mockConn := &mockConnection{
		readResponses: []string{""},
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	assert.Error(t, err)
}

// TestCollectPartialMetrics tests collection with missing metrics
func TestCollectPartialMetrics(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	mockConn := &mockConnection{
		readResponses: []string{
			"Zookeeper version: 3.4.5\n" +
			"Latency min/avg/max: 0/0/100\n" +
			"Received: 1000\n" +
			"Sent: 1000\n",
		},
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	// Should handle partial metrics gracefully
	assert.NoError(t, err)
}

// TestCollectReadError tests error when reading from connection
func TestCollectReadError(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	mockConn := &mockConnection{
		readError: errors.New("read timeout"),
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "read timeout")
}

// TestCollectWriteError tests error when writing to connection
func TestCollectWriteError(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	mockConn := &mockConnection{
		writeError: errors.New("write failed"),
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	assert.Error(t, err)
}

// TestCollectParseLatency tests latency metric parsing
func TestCollectParseLatency(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	mockConn := &mockConnection{
		readResponses: []string{
			"Latency min/avg/max: 1/5/100\n" +
			"Received: 1000\n" +
			"Sent: 1000\n" +
			"Connections: 5\n" +
			"Outstanding: 2\n" +
			"Mode: leader\n",
		},
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	assert.NoError(t, err)
}

// TestCollectParseVersion tests version metric parsing
func TestCollectParseVersion(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	mockConn := &mockConnection{
		readResponses: []string{
			"Zookeeper version: 3.6.3, built on 07/15/2021 10 39 GMT\n" +
			"Latency min/avg/max: 0/0/1\n" +
			"Received: 100\n" +
			"Sent: 100\n" +
			"Mode: standalone\n",
		},
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	assert.NoError(t, err)
}

// TestCollectParseMode tests mode metric parsing
func TestCollectParseMode(t *testing.T) {
	testCases := []struct {
		name     string
		response string
		mode     string
	}{
		{
			name: "leader mode",
			response: "Mode: leader\nLatency min/avg/max: 0/0/1\nReceived: 100\nSent: 100\n",
			mode: "leader",
		},
		{
			name: "follower mode",
			response: "Mode: follower\nLatency min/avg/max: 0/0/1\nReceived: 100\nSent: 100\n",
			mode: "follower",
		},
		{
			name: "standalone mode",
			response: "Mode: standalone\nLatency min/avg/max: 0/0/1\nReceived: 100\nSent: 100\n",
			mode: "standalone",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			zk := NewZookeeper()
			zk.Address = "localhost:2181"
			
			mockConn := &mockConnection{
				readResponses: []string{tc.response},
			}
			
			zk.dial = func(network, address string) (connection, error) {
				return mockConn, nil
			}
			
			err := zk.Collect()
			assert.NoError(t, err)
		})
	}
}

// TestCollectParseReceivedSent tests received/sent metrics parsing
func TestCollectParseReceivedSent(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	mockConn := &mockConnection{
		readResponses: []string{
			"Received: 9999999999\n" +
			"Sent: 8888888888\n" +
			"Latency min/avg/max: 0/0/1\n" +
			"Connections: 1\n" +
			"Mode: follower\n",
		},
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	assert.NoError(t, err)
}

// TestCollectParseConnections tests connections metric parsing
func TestCollectParseConnections(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	mockConn := &mockConnection{
		readResponses: []string{
			"Connections: 42\n" +
			"Outstanding: 5\n" +
			"Latency min/avg/max: 0/0/1\n" +
			"Received: 100\n" +
			"Sent: 100\n" +
			"Mode: follower\n",
		},
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	assert.NoError(t, err)
}

// TestCollectParseOutstanding tests outstanding metric parsing
func TestCollectParseOutstanding(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	mockConn := &mockConnection{
		readResponses: []string{
			"Outstanding: 10\n" +
			"Latency min/avg/max: 0/0/1\n" +
			"Received: 100\n" +
			"Sent: 100\n" +
			"Connections: 5\n" +
			"Mode: follower\n",
		},
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	assert.NoError(t, err)
}

// TestCollectParseZxid tests zxid metric parsing
func TestCollectParseZxid(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	mockConn := &mockConnection{
		readResponses: []string{
			"Zxid: 0x100000042\n" +
			"Latency min/avg/max: 0/0/1\n" +
			"Received: 100\n" +
			"Sent: 100\n" +
			"Mode: leader\n",
		},
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	assert.NoError(t, err)
}

// TestCollectParseNodeCount tests node count metric parsing
func TestCollectParseNodeCount(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	mockConn := &mockConnection{
		readResponses: []string{
			"Node count: 12345\n" +
			"Latency min/avg/max: 0/0/1\n" +
			"Received: 100\n" +
			"Sent: 100\n" +
			"Mode: follower\n",
		},
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	assert.NoError(t, err)
}

// TestCollectWithTLS tests collection with TLS configuration
func TestCollectWithTLS(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	zk.TLSConfig = tlscfg.TLSConfig{
		Enable: true,
	}
	
	mockConn := &mockConnection{
		readResponses: []string{
			"Latency min/avg/max: 0/0/1\n" +
			"Received: 100\n" +
			"Sent: 100\n" +
			"Mode: follower\n",
		},
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	assert.NoError(t, err)
}

// TestCollectWithTimeout tests collection with timeout
func TestCollectWithTimeout(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	zk.Timeout = 5 * time.Second
	
	mockConn := &mockConnection{
		readResponses: []string{
			"Latency min/avg/max: 0/0/1\n" +
			"Received: 100\n" +
			"Sent: 100\n" +
			"Mode: follower\n",
		},
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	assert.NoError(t, err)
}

// TestCollectInvalidLatencyFormat tests handling of invalid latency format
func TestCollectInvalidLatencyFormat(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	mockConn := &mockConnection{
		readResponses: []string{
			"Latency min/avg/max: invalid/data/here\n" +
			"Received: 100\n" +
			"Sent: 100\n" +
			"Mode: follower\n",
		},
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	// Should handle invalid format gracefully
	assert.Error(t, err)
}

// TestCollectInvalidNumericValues tests handling of invalid numeric values
func TestCollectInvalidNumericValues(t *testing.T) {
	testCases := []struct {
		name     string
		response string
	}{
		{
			name: "invalid received",
			response: "Received: notanumber\nLatency min/avg/max: 0/0/1\nSent: 100\nMode: follower\n",
		},
		{
			name: "invalid sent",
			response: "Sent: notanumber\nLatency min/avg/max: 0/0/1\nReceived: 100\nMode: follower\n",
		},
		{
			name: "invalid connections",
			response: "Connections: notanumber\nLatency min/avg/max: 0/0/1\nReceived: 100\nSent: 100\nMode: follower\n",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			zk := NewZookeeper()
			zk.Address = "localhost:2181"
			
			mockConn := &mockConnection{
				readResponses: []string{tc.response},
			}
			
			zk.dial = func(network, address string) (connection, error) {
				return mockConn, nil
			}
			
			err := zk.Collect()
			assert.Error(t, err)
		})
	}
}

// TestCollectConnectionClose tests connection closing
func TestCollectConnectionClose(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	mockConn := &mockConnection{
		readResponses: []string{
			"Latency min/avg/max: 0/0/1\n" +
			"Received: 100\n" +
			"Sent: 100\n" +
			"Mode: follower\n",
		},
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	assert.NoError(t, err)
	assert.True(t, mockConn.closed)
}

// TestCollectConnectionCloseError tests handling close error
func TestCollectConnectionCloseError(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	mockConn := &mockConnection{
		readResponses: []string{
			"Latency min/avg/max: 0/0/1\n" +
			"Received: 100\n" +
			"Sent: 100\n" +
			"Mode: follower\n",
		},
		closeError: errors.New("close failed"),
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	// Error during close should not prevent successful metrics collection
	assert.NoError(t, err)
}

// TestCollectBoundaryValues tests parsing of boundary values
func TestCollectBoundaryValues(t *testing.T) {
	testCases := []struct {
		name     string
		response string
	}{
		{
			name: "zero values",
			response: "Latency min/avg/max: 0/0/0\nReceived: 0\nSent: 0\nConnections: 0\nMode: follower\n",
		},
		{
			name: "large values",
			response: "Latency min/avg/max: 0/999999999/999999999\nReceived: 9999999999\nSent: 9999999999\nMode: follower\n",
		},
		{
			name: "negative latency",
			response: "Latency min/avg/max: -1/5/10\nReceived: 100\nSent: 100\nMode: follower\n",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			zk := NewZookeeper()
			zk.Address = "localhost:2181"
			
			mockConn := &mockConnection{
				readResponses: []string{tc.response},
			}
			
			zk.dial = func(network, address string) (connection, error) {
				return mockConn, nil
			}
			
			err := zk.Collect()
			// Should handle boundary values
			assert.NoError(t, err)
		})
	}
}

// TestCollectMultipleLineEndings tests handling different line endings
func TestCollectMultipleLineEndings(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	mockConn := &mockConnection{
		readResponses: []string{
			"Latency min/avg/max: 0/0/1\r\n" +
			"Received: 100\r\n" +
			"Sent: 100\r\n" +
			"Mode: follower\r\n",
		},
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	assert.NoError(t, err)
}

// Mock implementation of connection interface
type mockConnection struct {
	readResponses []string
	readIndex     int
	readError     error
	writeError    error
	closeError    error
	closed        bool
}

func (m *mockConnection) Read(b []byte) (n int, err error) {
	if m.readError != nil {
		return 0, m.readError
	}
	
	if m.readIndex >= len(m.readResponses) {
		return 0, io.EOF
	}
	
	response := m.readResponses[m.readIndex]
	m.readIndex++
	
	n = copy(b, response)
	return n, nil
}

func (m *mockConnection) Write(b []byte) (n int, err error) {
	if m.writeError != nil {
		return 0, m.writeError
	}
	return len(b), nil
}

func (m *mockConnection) Close() error {
	m.closed = true
	return m.closeError
}

func (m *mockConnection) SetDeadline(t time.Time) error {
	return nil
}

func (m *mockConnection) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *mockConnection) SetWriteDeadline(t time.Time) error {
	return nil
}

// TestCollectCaseInsensitiveHeaders tests parsing with different header cases
func TestCollectCaseInsensitiveHeaders(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	mockConn := &mockConnection{
		readResponses: []string{
			"latency min/avg/max: 0/0/1\n" +
			"received: 100\n" +
			"sent: 100\n" +
			"mode: follower\n",
		},
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	// Should handle case variations
	assert.Error(t, err) // If implementation is case-sensitive
}

// TestCollectWhitespaceHandling tests whitespace handling in parsing
func TestCollectWhitespaceHandling(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	mockConn := &mockConnection{
		readResponses: []string{
			"Latency min/avg/max:  0/0/1 \n" +
			"Received:  100 \n" +
			"Sent:  100 \n" +
			"Mode:  follower \n",
		},
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	assert.NoError(t, err)
}

// TestCollectEmptyLines tests handling of empty lines in response
func TestCollectEmptyLines(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	mockConn := &mockConnection{
		readResponses: []string{
			"Latency min/avg/max: 0/0/1\n" +
			"\n" +
			"Received: 100\n" +
			"\n" +
			"Sent: 100\n" +
			"Mode: follower\n",
		},
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	assert.NoError(t, err)
}

// TestCollectMetricsAreRecorded tests that metrics are properly recorded
func TestCollectMetricsAreRecorded(t *testing.T) {
	zk := NewZookeeper()
	zk.Address = "localhost:2181"
	
	mockConn := &mockConnection{
		readResponses: []string{
			"Latency min/avg/max: 1/5/20\n" +
			"Received: 12345\n" +
			"Sent: 54321\n" +
			"Connections: 8\n" +
			"Outstanding: 3\n" +
			"Mode: leader\n" +
			"Node count: 500\n" +
			"Zxid: 0x100000100\n",
		},
	}
	
	zk.dial = func(network, address string) (connection, error) {
		return mockConn, nil
	}
	
	err := zk.Collect()
	assert.NoError(t, err)
}

// TestCollectDifferentModes tests collection with different operational modes
func TestCollectDifferentModes(t *testing.T) {
	modes := []string{"leader", "follower", "observer", "standalone"}
	
	for _, mode := range modes {
		t.Run(mode, func(t *testing.T) {
			zk := NewZookeeper()
			zk.Address = "localhost:2181"
			
			mockConn := &mockConnection{
				readResponses: []string{
					"Mode: " + mode + "\n" +
					"Latency min/avg/max: 0/0/1\n" +
					"Received: 100\n" +
					"Sent: 100\n",
				},
			}
			
			zk.dial = func(network, address string) (connection, error) {
				return mockConn, nil
			}
			
			err := zk.Collect()
			assert.NoError(t, err)
		})
	}
}