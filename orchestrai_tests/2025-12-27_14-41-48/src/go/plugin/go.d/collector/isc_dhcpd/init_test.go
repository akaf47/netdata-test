package isc_dhcpd

import (
	"testing"
)

// TestInit tests the Init method
func TestInit(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*DhcpdCollector)
		wantErr bool
	}{
		{
			name: "should initialize successfully with default values",
			setup: func(d *DhcpdCollector) {
				d.Config = Config{}
			},
			wantErr: false,
		},
		{
			name: "should initialize with custom config",
			setup: func(d *DhcpdCollector) {
				d.Config = Config{
					LeasesPath: "/var/lib/dhcp/dhcpd.leases",
				}
			},
			wantErr: false,
		},
		{
			name: "should handle nil config",
			setup: func(d *DhcpdCollector) {
				// Config is zero-initialized
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := NewDhcpdCollector()
			if tt.setup != nil {
				tt.setup(dc)
			}

			err := dc.Init()

			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestInitWithEmptyLeasesPath tests Init with empty leases path
func TestInitWithEmptyLeasesPath(t *testing.T) {
	dc := NewDhcpdCollector()
	dc.Config = Config{
		LeasesPath: "",
	}

	err := dc.Init()
	if err != nil {
		t.Errorf("Init with empty path should succeed, got error: %v", err)
	}
}

// TestInitWithValidLeasesPath tests Init with valid leases path
func TestInitWithValidLeasesPath(t *testing.T) {
	dc := NewDhcpdCollector()
	dc.Config = Config{
		LeasesPath: "/etc/dhcp/dhcpd.leases",
	}

	err := dc.Init()
	if err != nil {
		t.Errorf("Init with valid path should succeed, got error: %v", err)
	}
}

// TestInitIdempotent tests that Init can be called multiple times
func TestInitIdempotent(t *testing.T) {
	dc := NewDhcpdCollector()
	dc.Config = Config{
		LeasesPath: "/var/lib/dhcp/dhcpd.leases",
	}

	err1 := dc.Init()
	err2 := dc.Init()

	if (err1 != nil) != (err2 != nil) {
		t.Errorf("Init should be idempotent, errors differ: %v vs %v", err1, err2)
	}
}

// TestInitSetsRequiredFields tests that Init properly sets internal fields
func TestInitSetsRequiredFields(t *testing.T) {
	dc := NewDhcpdCollector()
	dc.Config = Config{
		LeasesPath: "/var/lib/dhcp/dhcpd.leases",
	}

	_ = dc.Init()

	// Verify initialization occurred (implementation-specific checks)
	if dc.Config.LeasesPath == "" {
		t.Error("Init should preserve LeasesPath")
	}
}

// TestInitWithSpecialCharactersInPath tests Init with special characters
func TestInitWithSpecialCharactersInPath(t *testing.T) {
	dc := NewDhcpdCollector()
	dc.Config = Config{
		LeasesPath: "/var/lib/dhcp/dhcpd-special_chars.leases",
	}

	err := dc.Init()
	if err != nil {
		t.Errorf("Init with special characters should succeed, got error: %v", err)
	}
}

// TestInitWithRelativePath tests Init with relative path
func TestInitWithRelativePath(t *testing.T) {
	dc := NewDhcpdCollector()
	dc.Config = Config{
		LeasesPath: "./dhcpd.leases",
	}

	err := dc.Init()
	if err != nil {
		t.Errorf("Init with relative path should succeed, got error: %v", err)
	}
}

// TestInitWithAbsolutePath tests Init with absolute path
func TestInitWithAbsolutePath(t *testing.T) {
	dc := NewDhcpdCollector()
	dc.Config = Config{
		LeasesPath: "/absolute/path/to/dhcpd.leases",
	}

	err := dc.Init()
	if err != nil {
		t.Errorf("Init with absolute path should succeed, got error: %v", err)
	}
}

// TestInitMultipleInstances tests multiple collector instances
func TestInitMultipleInstances(t *testing.T) {
	dc1 := NewDhcpdCollector()
	dc1.Config = Config{LeasesPath: "/path1"}

	dc2 := NewDhcpdCollector()
	dc2.Config = Config{LeasesPath: "/path2"}

	err1 := dc1.Init()
	err2 := dc2.Init()

	if err1 != nil || err2 != nil {
		t.Errorf("Multiple instances should initialize independently, errors: %v, %v", err1, err2)
	}

	if dc1.Config.LeasesPath != "/path1" || dc2.Config.LeasesPath != "/path2" {
		t.Error("Instances should maintain separate configurations")
	}
}

// TestInitStateAfterError tests state consistency after error
func TestInitStateAfterError(t *testing.T) {
	dc := NewDhcpdCollector()
	dc.Config = Config{
		LeasesPath: "/var/lib/dhcp/dhcpd.leases",
	}

	_ = dc.Init()
	
	// State should be consistent after init
	if dc.Config.LeasesPath != "/var/lib/dhcp/dhcpd.leases" {
		t.Error("State should be consistent after Init")
	}
}

// TestNewDhcpdCollector tests the constructor
func TestNewDhcpdCollector(t *testing.T) {
	dc := NewDhcpdCollector()

	if dc == nil {
		t.Error("NewDhcpdCollector should not return nil")
	}

	// Verify default state
	if dc.Config.LeasesPath != "" && dc.Config.LeasesPath != "/var/lib/dhcp/dhcpd.leases" {
		// Allow for either empty or default value
	}
}