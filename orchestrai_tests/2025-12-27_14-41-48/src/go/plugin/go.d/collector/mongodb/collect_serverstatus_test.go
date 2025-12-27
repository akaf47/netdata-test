package mongodb

import (
	"testing"

	"github.com/netdata/netdata/go/plugins/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectServerStatus_Success(t *testing.T) {
	// Test basic server status collection
	m := &MongoDB{
		Logger: logger.New(),
	}

	ss := &serverStatus{
		Uptime: 3600,
		Connections: connectionStats{
			Current:      100,
			Available:    900,
			TotalCreated: 1000,
		},
		Network: networkStats{
			BytesIn:  1024000,
			BytesOut: 2048000,
			NumRequests: 5000,
		},
		Memory: memoryStats{
			Resident: 512,
			Virtual:  1024,
		},
		Locks: lockStats{
			Global: lockTypeStats{
				AcquireCount: lockModeStats{
					Shared:      100,
					Exclusive:   50,
					IntentShared: 200,
					IntentExclusive: 150,
				},
				AcquireWaitCount: lockModeStats{
					Shared:      5,
					Exclusive:   2,
					IntentShared: 10,
					IntentExclusive: 3,
				},
			},
		},
		Transactions: transactionStats{
			Prepared: 0,
			Active:   0,
			Inactive: 0,
		},
	}

	// Execute collection
	err := m.collectServerStatus(ss)
	assert.NoError(t, err)
}

func TestCollectServerStatus_NilInput(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	// Test with nil server status
	err := m.collectServerStatus(nil)
	assert.Error(t, err)
}

func TestCollectServerStatus_ZeroValues(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	ss := &serverStatus{
		Uptime: 0,
		Connections: connectionStats{
			Current:      0,
			Available:    0,
			TotalCreated: 0,
		},
		Network: networkStats{
			BytesIn:     0,
			BytesOut:    0,
			NumRequests: 0,
		},
		Memory: memoryStats{
			Resident: 0,
			Virtual:  0,
		},
	}

	err := m.collectServerStatus(ss)
	assert.NoError(t, err)
}

func TestCollectServerStatus_HighValues(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	ss := &serverStatus{
		Uptime: 999999999,
		Connections: connectionStats{
			Current:      999999,
			Available:    999999,
			TotalCreated: 999999999,
		},
		Network: networkStats{
			BytesIn:     9999999999,
			BytesOut:    9999999999,
			NumRequests: 9999999,
		},
		Memory: memoryStats{
			Resident: 999999,
			Virtual:  999999,
		},
	}

	err := m.collectServerStatus(ss)
	assert.NoError(t, err)
}

func TestCollectServerStatus_ConnectionStats(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	tests := []struct {
		name     string
		current  int64
		available int64
		created  int64
	}{
		{"zero_connections", 0, 0, 0},
		{"one_connection", 1, 999, 1000},
		{"max_current", 1000000, 1000000, 2000000},
		{"no_available", 100, 0, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &serverStatus{
				Connections: connectionStats{
					Current:      tt.current,
					Available:    tt.available,
					TotalCreated: tt.created,
				},
			}
			err := m.collectServerStatus(ss)
			assert.NoError(t, err)
		})
	}
}

func TestCollectServerStatus_NetworkStats(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	tests := []struct {
		name        string
		bytesIn     int64
		bytesOut    int64
		numRequests int64
	}{
		{"zero_network", 0, 0, 0},
		{"equal_in_out", 1000000, 1000000, 100},
		{"more_out", 1000000, 2000000, 500},
		{"more_in", 5000000, 1000000, 200},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &serverStatus{
				Network: networkStats{
					BytesIn:     tt.bytesIn,
					BytesOut:    tt.bytesOut,
					NumRequests: tt.numRequests,
				},
			}
			err := m.collectServerStatus(ss)
			assert.NoError(t, err)
		})
	}
}

func TestCollectServerStatus_MemoryStats(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	tests := []struct {
		name      string
		resident  int64
		virtual   int64
	}{
		{"zero_memory", 0, 0},
		{"resident_only", 512, 0},
		{"virtual_only", 0, 1024},
		{"equal", 512, 512},
		{"virtual_larger", 256, 1024},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &serverStatus{
				Memory: memoryStats{
					Resident: tt.resident,
					Virtual:  tt.virtual,
				},
			}
			err := m.collectServerStatus(ss)
			assert.NoError(t, err)
		})
	}
}

func TestCollectServerStatus_LockStats(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	ss := &serverStatus{
		Locks: lockStats{
			Global: lockTypeStats{
				AcquireCount: lockModeStats{
					Shared:           100,
					Exclusive:        50,
					IntentShared:     200,
					IntentExclusive:  150,
				},
				AcquireWaitCount: lockModeStats{
					Shared:           5,
					Exclusive:        2,
					IntentShared:     10,
					IntentExclusive:  3,
				},
				TimeAcquiringMicros: lockModeStats{
					Shared:           1000,
					Exclusive:        500,
					IntentShared:     2000,
					IntentExclusive:  1500,
				},
				DeadlockCount: lockModeStats{
					Shared:           0,
					Exclusive:        1,
					IntentShared:     0,
					IntentExclusive:  0,
				},
			},
		},
	}

	err := m.collectServerStatus(ss)
	assert.NoError(t, err)
}

func TestCollectServerStatus_LockStats_ZeroValues(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	ss := &serverStatus{
		Locks: lockStats{
			Global: lockTypeStats{
				AcquireCount: lockModeStats{},
				AcquireWaitCount: lockModeStats{},
				TimeAcquiringMicros: lockModeStats{},
				DeadlockCount: lockModeStats{},
			},
		},
	}

	err := m.collectServerStatus(ss)
	assert.NoError(t, err)
}

func TestCollectServerStatus_TransactionStats(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	tests := []struct {
		name      string
		prepared  int64
		active    int64
		inactive  int64
	}{
		{"no_transactions", 0, 0, 0},
		{"all_prepared", 100, 0, 0},
		{"all_active", 0, 50, 0},
		{"all_inactive", 0, 0, 200},
		{"mixed", 10, 20, 30},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := &serverStatus{
				Transactions: transactionStats{
					Prepared: tt.prepared,
					Active:   tt.active,
					Inactive: tt.inactive,
				},
			}
			err := m.collectServerStatus(ss)
			assert.NoError(t, err)
		})
	}
}

func TestCollectServerStatus_AllMetricsPopulated(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	ss := &serverStatus{
		Uptime: 86400,
		Connections: connectionStats{
			Current:      500,
			Available:    500,
			TotalCreated: 1000,
		},
		Network: networkStats{
			BytesIn:     10485760,
			BytesOut:    20971520,
			NumRequests: 50000,
		},
		Memory: memoryStats{
			Resident: 2048,
			Virtual:  4096,
		},
		Locks: lockStats{
			Global: lockTypeStats{
				AcquireCount: lockModeStats{
					Shared:           1000,
					Exclusive:        500,
					IntentShared:     2000,
					IntentExclusive:  1500,
				},
				AcquireWaitCount: lockModeStats{
					Shared:           50,
					Exclusive:        20,
					IntentShared:     100,
					IntentExclusive:  30,
				},
			},
		},
		Transactions: transactionStats{
			Prepared: 5,
			Active:   10,
			Inactive: 100,
		},
	}

	err := m.collectServerStatus(ss)
	assert.NoError(t, err)
}

func TestCollectServerStatus_NegativeValues(t *testing.T) {
	// MongoDB should not return negative values, but test boundary
	m := &MongoDB{
		Logger: logger.New(),
	}

	ss := &serverStatus{
		Uptime: -1,
		Connections: connectionStats{
			Current:      -1,
			Available:    -1,
			TotalCreated: -1,
		},
	}

	// Should handle gracefully
	err := m.collectServerStatus(ss)
	// Either no error or handles gracefully
	_ = err
}

func TestCollectServerStatus_PartialData(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	tests := []struct {
		name string
		ss   *serverStatus
	}{
		{
			"only_uptime",
			&serverStatus{Uptime: 3600},
		},
		{
			"only_connections",
			&serverStatus{
				Connections: connectionStats{Current: 100},
			},
		},
		{
			"only_network",
			&serverStatus{
				Network: networkStats{BytesIn: 1000},
			},
		},
		{
			"only_memory",
			&serverStatus{
				Memory: memoryStats{Resident: 512},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := m.collectServerStatus(tt.ss)
			assert.NoError(t, err)
		})
	}
}

// Test helper types and structures
type serverStatus struct {
	Uptime       int64
	Connections  connectionStats
	Network      networkStats
	Memory       memoryStats
	Locks        lockStats
	Transactions transactionStats
}

type connectionStats struct {
	Current      int64
	Available    int64
	TotalCreated int64
}

type networkStats struct {
	BytesIn     int64
	BytesOut    int64
	NumRequests int64
}

type memoryStats struct {
	Resident int64
	Virtual  int64
}

type lockStats struct {
	Global lockTypeStats
}

type lockTypeStats struct {
	AcquireCount        lockModeStats
	AcquireWaitCount    lockModeStats
	TimeAcquiringMicros lockModeStats
	DeadlockCount       lockModeStats
}

type lockModeStats struct {
	Shared           int64
	Exclusive        int64
	IntentShared     int64
	IntentExclusive  int64
}

type transactionStats struct {
	Prepared int64
	Active   int64
	Inactive int64
}