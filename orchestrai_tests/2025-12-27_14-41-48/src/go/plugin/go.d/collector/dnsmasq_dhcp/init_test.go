package dnsmasq_dhcp

import (
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDNSMasqDHCP(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "should create new instance with default values",
			test: func(t *testing.T) {
				d := New()
				assert.NotNil(t, d)
				assert.NotNil(t, d.Config)
				assert.NotNil(t, d.Client)
			},
		},
		{
			name: "should initialize with valid module base",
			test: func(t *testing.T) {
				d := New()
				require.NotNil(t, d)
				// Verify it implements the necessary interfaces
				assert.NotNil(t, d.FullModule)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestDNSMasqDHCP_Init(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*DNSMasqDHCP)
		wantErr bool
		test    func(t *testing.T, d *DNSMasqDHCP)
	}{
		{
			name: "should initialize successfully with valid config",
			setup: func(d *DNSMasqDHCP) {
				d.Config.LeaseFile = "/tmp/valid/lease/file"
			},
			wantErr: false,
			test: func(t *testing.T, d *DNSMasqDHCP) {
				assert.NotNil(t, d.Client)
			},
		},
		{
			name: "should fail when lease file is empty",
			setup: func(d *DNSMasqDHCP) {
				d.Config.LeaseFile = ""
			},
			wantErr: true,
			test: func(t *testing.T, d *DNSMasqDHCP) {
				// Verify that proper error handling occurred
			},
		},
		{
			name: "should fail when lease file path is invalid",
			setup: func(d *DNSMasqDHCP) {
				d.Config.LeaseFile = "\x00invalid\x00path"
			},
			wantErr: true,
			test: func(t *testing.T, d *DNSMasqDHCP) {
				// Verify that proper error handling occurred
			},
		},
		{
			name: "should handle nil Config gracefully",
			setup: func(d *DNSMasqDHCP) {
				d.Config = nil
			},
			wantErr: true,
			test: func(t *testing.T, d *DNSMasqDHCP) {
				// Verify that proper error handling occurred
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			if tt.setup != nil {
				tt.setup(d)
			}

			err := d.Init()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.test != nil {
				tt.test(t, d)
			}
		})
	}
}

func TestDNSMasqDHCP_Check(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*DNSMasqDHCP)
		wantErr bool
		test    func(t *testing.T, d *DNSMasqDHCP)
	}{
		{
			name: "should return no error on successful check",
			setup: func(d *DNSMasqDHCP) {
				d.Config.LeaseFile = "/tmp/valid/file"
			},
			wantErr: false,
			test: func(t *testing.T, d *DNSMasqDHCP) {
				assert.NotNil(t, d)
			},
		},
		{
			name: "should fail check when client is nil",
			setup: func(d *DNSMasqDHCP) {
				d.Client = nil
			},
			wantErr: true,
			test: func(t *testing.T, d *DNSMasqDHCP) {
				assert.Nil(t, d.Client)
			},
		},
		{
			name: "should fail check when unable to read lease file",
			setup: func(d *DNSMasqDHCP) {
				d.Config.LeaseFile = "/nonexistent/path/that/does/not/exist"
			},
			wantErr: true,
			test: func(t *testing.T, d *DNSMasqDHCP) {
				assert.NotNil(t, d)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			if tt.setup != nil {
				tt.setup(d)
			}

			err := d.Check()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.test != nil {
				tt.test(t, d)
			}
		})
	}
}

func TestDNSMasqDHCP_Charts(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "should return charts module",
			test: func(t *testing.T) {
				d := New()
				charts := d.Charts()
				assert.NotNil(t, charts)
			},
		},
		{
			name: "should have valid chart module structure",
			test: func(t *testing.T) {
				d := New()
				charts := d.Charts()
				assert.NotNil(t, charts)
				assert.IsType(t, (*module.Charts)(nil), charts)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestDNSMasqDHCP_Collect(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*DNSMasqDHCP)
		wantErr bool
		test    func(t *testing.T, d *DNSMasqDHCP)
	}{
		{
			name: "should collect metrics successfully",
			setup: func(d *DNSMasqDHCP) {
				d.Config.LeaseFile = "/tmp/dnsmasq.leases"
			},
			wantErr: false,
			test: func(t *testing.T, d *DNSMasqDHCP) {
				assert.NotNil(t, d)
			},
		},
		{
			name: "should fail when client is nil",
			setup: func(d *DNSMasqDHCP) {
				d.Client = nil
			},
			wantErr: true,
			test: func(t *testing.T, d *DNSMasqDHCP) {
				assert.Nil(t, d.Client)
			},
		},
		{
			name: "should fail when unable to read lease file",
			setup: func(d *DNSMasqDHCP) {
				d.Config.LeaseFile = "/nonexistent/dnsmasq.leases"
			},
			wantErr: true,
			test: func(t *testing.T, d *DNSMasqDHCP) {
				assert.NotNil(t, d)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			if tt.setup != nil {
				tt.setup(d)
			}

			result := d.Collect()
			if tt.wantErr {
				assert.Nil(t, result)
			} else {
				// Collect should return data even if it's empty
				assert.NotNil(t, result)
			}

			if tt.test != nil {
				tt.test(t, d)
			}
		})
	}
}

func TestDNSMasqDHCP_Cleanup(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "should execute cleanup without error",
			test: func(t *testing.T) {
				d := New()
				// Should not panic or error
				d.Cleanup()
			},
		},
		{
			name: "should be callable multiple times",
			test: func(t *testing.T) {
				d := New()
				d.Cleanup()
				d.Cleanup()
				// Should not panic on second call
			},
		},
		{
			name: "should handle nil client gracefully",
			test: func(t *testing.T) {
				d := New()
				d.Client = nil
				d.Cleanup()
				// Should not panic when client is nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestDNSMasqDHCP_InitialLogic(t *testing.T) {
	t.Run("should create collector with proper module interface", func(t *testing.T) {
		d := New()
		require.NotNil(t, d)
		
		// Test that it has the FullModule interface
		assert.NotNil(t, d.FullModule)
		
		// Test Config is properly initialized
		assert.NotNil(t, d.Config)
	})
}

func TestDNSMasqDHCP_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		test  func(t *testing.T)
	}{
		{
			name: "should handle empty lease file path",
			test: func(t *testing.T) {
				d := New()
				d.Config.LeaseFile = ""
				err := d.Init()
				assert.Error(t, err)
			},
		},
		{
			name: "should handle very long lease file path",
			test: func(t *testing.T) {
				d := New()
				d.Config.LeaseFile = "/path/" + string(make([]byte, 4096)) + "/file"
				// Should handle without crashing
				d.Init()
			},
		},
		{
			name: "should handle special characters in path",
			test: func(t *testing.T) {
				d := New()
				d.Config.LeaseFile = "/path/with spaces/and-special_chars/file.leases"
				// Should attempt to process
				_ = d.Init()
			},
		},
		{
			name: "should handle repeated initialization",
			test: func(t *testing.T) {
				d := New()
				d.Config.LeaseFile = "/tmp/test"
				d.Init()
				d.Init() // Second init should not crash
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}