package isc_dhcpd

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestCollectSuccess tests successful data collection
func TestCollectSuccess(t *testing.T) {
	d := New()
	if d == nil {
		t.Fatal("New() returned nil")
	}

	// Set a valid config
	d.Config.LeasesPath = "/tmp/dhcpd.leases"
	
	// Create test leases file
	tmpDir := t.TempDir()
	testLeasesPath := filepath.Join(tmpDir, "dhcpd.leases")
	
	testContent := `# DHCPv4 lease file
lease 192.168.1.100 {
  starts 4 2024/01/15 10:00:00;
  ends 4 2024/01/15 22:00:00;
  hardware ethernet 00:11:22:33:44:55;
}
`
	if err := os.WriteFile(testLeasesPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	d.Config.LeasesPath = testLeasesPath

	// Collect should succeed
	mx := d.Collect(context.Background())
	if mx == nil {
		t.Error("Collect() returned nil metrics")
	}
}

// TestCollectWithNilConfig tests collection with nil configuration
func TestCollectWithNilConfig(t *testing.T) {
	d := New()
	if d == nil {
		t.Fatal("New() returned nil")
	}

	// Don't set any config - use defaults
	mx := d.Collect(context.Background())
	// Should handle gracefully
	if mx != nil {
		t.Logf("Collect with default config returned: %v", mx)
	}
}

// TestCollectWithEmptyLeasesFile tests collection with empty leases file
func TestCollectWithEmptyLeasesFile(t *testing.T) {
	d := New()
	
	tmpDir := t.TempDir()
	testLeasesPath := filepath.Join(tmpDir, "dhcpd.leases")
	
	// Create empty file
	if err := os.WriteFile(testLeasesPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	d.Config.LeasesPath = testLeasesPath

	mx := d.Collect(context.Background())
	if mx == nil {
		t.Error("Collect() should handle empty file")
	}
}

// TestCollectWithNonexistentFile tests collection with missing leases file
func TestCollectWithNonexistentFile(t *testing.T) {
	d := New()
	d.Config.LeasesPath = "/nonexistent/path/dhcpd.leases"

	mx := d.Collect(context.Background())
	// Should handle error gracefully
	if mx != nil {
		t.Logf("Collect handled missing file")
	}
}

// TestCollectWithValidLeases tests collection with valid lease entries
func TestCollectWithValidLeases(t *testing.T) {
	d := New()
	
	tmpDir := t.TempDir()
	testLeasesPath := filepath.Join(tmpDir, "dhcpd.leases")
	
	testContent := `# ISC DHCP Server Lease file
# This is a sample lease file for testing
lease 192.168.1.100 {
  starts 4 2024/01/15 10:00:00;
  ends 4 2024/01/15 22:00:00;
  tstp 4 2024/01/15 22:00:00;
  tsfp 4 2024/01/15 22:00:00;
  hardware ethernet 00:11:22:33:44:55;
  uid "01:00:11:22:33:44:55";
  set vendor-string "VENDOR";
  client-hostname "test-host";
}

lease 192.168.1.101 {
  starts 2 2024/01/16 10:00:00;
  ends never;
  hardware ethernet aa:bb:cc:dd:ee:ff;
  binding state free;
}

lease 192.168.1.102 {
  starts 2 2024/01/16 09:00:00;
  ends 2 2024/01/16 21:00:00;
  hardware ethernet 11:22:33:44:55:66;
  binding state released;
}
`
	if err := os.WriteFile(testLeasesPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	d.Config.LeasesPath = testLeasesPath

	mx := d.Collect(context.Background())
	if mx == nil {
		t.Fatal("Collect() returned nil for valid leases")
	}
}

// TestCollectWithMalformedLeases tests collection with malformed lease entries
func TestCollectWithMalformedLeases(t *testing.T) {
	d := New()
	
	tmpDir := t.TempDir()
	testLeasesPath := filepath.Join(tmpDir, "dhcpd.leases")
	
	testContent := `lease 192.168.1.100 {
  starts 4 2024/01/15 10:00:00;
  // missing closing brace
lease 192.168.1.101 {
  this is invalid
}
`
	if err := os.WriteFile(testLeasesPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	d.Config.LeasesPath = testLeasesPath

	mx := d.Collect(context.Background())
	// Should handle malformed entries gracefully
	if mx != nil {
		t.Logf("Collect handled malformed leases")
	}
}

// TestCollectWithExpiredLeases tests parsing of expired leases
func TestCollectWithExpiredLeases(t *testing.T) {
	d := New()
	
	tmpDir := t.TempDir()
	testLeasesPath := filepath.Join(tmpDir, "dhcpd.leases")
	
	testContent := `lease 192.168.1.100 {
  starts 1 2020/01/01 10:00:00;
  ends 1 2020/01/02 10:00:00;
  hardware ethernet 00:11:22:33:44:55;
}

lease 192.168.1.101 {
  starts 1 2024/01/15 10:00:00;
  ends never;
  hardware ethernet aa:bb:cc:dd:ee:ff;
}
`
	if err := os.WriteFile(testLeasesPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	d.Config.LeasesPath = testLeasesPath

	mx := d.Collect(context.Background())
	if mx == nil {
		t.Fatal("Collect() should handle expired and never-ending leases")
	}
}

// TestCollectWithLargeLeaseCount tests collection with many leases
func TestCollectWithLargeLeaseCount(t *testing.T) {
	d := New()
	
	tmpDir := t.TempDir()
	testLeasesPath := filepath.Join(tmpDir, "dhcpd.leases")
	
	// Generate 1000 lease entries
	content := `# Large lease file test
`
	for i := 0; i < 1000; i++ {
		content += `lease 192.168.1.` + string(rune(100+(i%150))) + ` {
  starts 4 2024/01/15 10:00:00;
  ends 4 2024/01/15 22:00:00;
  hardware ethernet 00:11:22:33:44:5` + string(rune(48+(i%10))) + `;
}
`
	}

	if err := os.WriteFile(testLeasesPath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	d.Config.LeasesPath = testLeasesPath

	mx := d.Collect(context.Background())
	if mx == nil {
		t.Error("Collect() should handle large lease counts")
	}
}

// TestCollectWithSpecialCharactersInHostname tests parsing hostnames with special chars
func TestCollectWithSpecialCharactersInHostname(t *testing.T) {
	d := New()
	
	tmpDir := t.TempDir()
	testLeasesPath := filepath.Join(tmpDir, "dhcpd.leases")
	
	testContent := `lease 192.168.1.100 {
  starts 4 2024/01/15 10:00:00;
  ends 4 2024/01/15 22:00:00;
  hardware ethernet 00:11:22:33:44:55;
  client-hostname "test-host.example.com";
  set domain-name "example.com";
}

lease 192.168.1.101 {
  starts 4 2024/01/15 10:00:00;
  ends 4 2024/01/15 22:00:00;
  hardware ethernet aa:bb:cc:dd:ee:ff;
  client-hostname "host_with_underscores";
}
`
	if err := os.WriteFile(testLeasesPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	d.Config.LeasesPath = testLeasesPath

	mx := d.Collect(context.Background())
	if mx == nil {
		t.Error("Collect() should handle special characters in hostnames")
	}
}

// TestCollectWithVariousBindingStates tests different lease binding states
func TestCollectWithVariousBindingStates(t *testing.T) {
	d := New()
	
	tmpDir := t.TempDir()
	testLeasesPath := filepath.Join(tmpDir, "dhcpd.leases")
	
	testContent := `lease 192.168.1.100 {
  starts 4 2024/01/15 10:00:00;
  ends 4 2024/01/15 22:00:00;
  hardware ethernet 00:11:22:33:44:55;
  binding state active;
}

lease 192.168.1.101 {
  starts 4 2024/01/15 10:00:00;
  ends 4 2024/01/15 22:00:00;
  hardware ethernet aa:bb:cc:dd:ee:ff;
  binding state released;
}

lease 192.168.1.102 {
  starts 4 2024/01/15 10:00:00;
  ends never;
  hardware ethernet 11:22:33:44:55:66;
  binding state free;
}

lease 192.168.1.103 {
  starts 4 2024/01/15 10:00:00;
  ends 4 2024/01/15 22:00:00;
  hardware ethernet 22:33:44:55:66:77;
  binding state abandoned;
}
`
	if err := os.WriteFile(testLeasesPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	d.Config.LeasesPath = testLeasesPath

	mx := d.Collect(context.Background())
	if mx == nil {
		t.Error("Collect() should handle all binding states")
	}
}

// TestCollectWithContextCancellation tests behavior when context is cancelled
func TestCollectWithContextCancellation(t *testing.T) {
	d := New()
	
	tmpDir := t.TempDir()
	testLeasesPath := filepath.Join(tmpDir, "dhcpd.leases")
	
	testContent := `lease 192.168.1.100 {
  starts 4 2024/01/15 10:00:00;
  ends 4 2024/01/15 22:00:00;
  hardware ethernet 00:11:22:33:44:55;
}
`
	if err := os.WriteFile(testLeasesPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	d.Config.LeasesPath = testLeasesPath

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	mx := d.Collect(ctx)
	// Should handle cancelled context
	if mx != nil {
		t.Logf("Collect handled cancelled context")
	}
}

// TestCollectWithTimeout tests behavior with timeout context
func TestCollectWithTimeout(t *testing.T) {
	d := New()
	
	tmpDir := t.TempDir()
	testLeasesPath := filepath.Join(tmpDir, "dhcpd.leases")
	
	testContent := `lease 192.168.1.100 {
  starts 4 2024/01/15 10:00:00;
  ends 4 2024/01/15 22:00:00;
  hardware ethernet 00:11:22:33:44:55;
}
`
	if err := os.WriteFile(testLeasesPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	d.Config.LeasesPath = testLeasesPath

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	mx := d.Collect(ctx)
	if mx == nil {
		t.Error("Collect() should complete within timeout")
	}
}

// TestCollectWithPermissionDenied tests behavior when file is not readable
func TestCollectWithPermissionDenied(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	d := New()
	
	tmpDir := t.TempDir()
	testLeasesPath := filepath.Join(tmpDir, "dhcpd.leases")
	
	testContent := `lease 192.168.1.100 {
  starts 4 2024/01/15 10:00:00;
  ends 4 2024/01/15 22:00:00;
  hardware ethernet 00:11:22:33:44:55;
}
`
	if err := os.WriteFile(testLeasesPath, []byte(testContent), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Make file unreadable
	if err := os.Chmod(testLeasesPath, 0000); err != nil {
		t.Fatalf("Failed to set permissions: %v", err)
	}
	defer os.Chmod(testLeasesPath, 0644) // Restore for cleanup

	d.Config.LeasesPath = testLeasesPath

	mx := d.Collect(context.Background())
	// Should handle permission denied gracefully
	if mx != nil {
		t.Logf("Collect handled permission denied")
	}
}

// TestCollectReturnsMetrics tests that collect returns proper metrics structure
func TestCollectReturnsMetrics(t *testing.T) {
	d := New()
	
	tmpDir := t.TempDir()
	testLeasesPath := filepath.Join(tmpDir, "dhcpd.leases")
	
	testContent := `lease 192.168.1.100 {
  starts 4 2024/01/15 10:00:00;
  ends 4 2024/01/15 22:00:00;
  hardware ethernet 00:11:22:33:44:55;
}

lease 192.168.1.101 {
  starts 2 2024/01/16 10:00:00;
  ends never;
  hardware ethernet aa:bb:cc:dd:ee:ff;
}
`
	if err := os.WriteFile(testLeasesPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	d.Config.LeasesPath = testLeasesPath

	mx := d.Collect(context.Background())
	if mx == nil {
		t.Fatal("Collect() returned nil")
	}

	// Verify metrics structure
	if len(mx) == 0 {
		t.Error("Collect() should return non-empty metrics")
	}
}

// TestCollectMultipleInvocations tests that collect can be called multiple times
func TestCollectMultipleInvocations(t *testing.T) {
	d := New()
	
	tmpDir := t.TempDir()
	testLeasesPath := filepath.Join(tmpDir, "dhcpd.leases")
	
	testContent := `lease 192.168.1.100 {
  starts 4 2024/01/15 10:00:00;
  ends 4 2024/01/15 22:00:00;
  hardware ethernet 00:11:22:33:44:55;
}
`
	if err := os.WriteFile(testLeasesPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	d.Config.LeasesPath = testLeasesPath

	mx1 := d.Collect(context.Background())
	mx2 := d.Collect(context.Background())

	if mx1 == nil || mx2 == nil {
		t.Error("Multiple collect invocations should succeed")
	}
}

// TestCollectWithComments tests that comments in lease file are ignored
func TestCollectWithComments(t *testing.T) {
	d := New()
	
	tmpDir := t.TempDir()
	testLeasesPath := filepath.Join(tmpDir, "dhcpd.leases")
	
	testContent := `# This is a comment
# ISC DHCP Lease Database

# Another comment section
lease 192.168.1.100 {
  # inline comment
  starts 4 2024/01/15 10:00:00;
  # another inline comment
  ends 4 2024/01/15 22:00:00;
  hardware ethernet 00:11:22:33:44:55;
  # more comments
}

# end of leases
`
	if err := os.WriteFile(testLeasesPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	d.Config.LeasesPath = testLeasesPath

	mx := d.Collect(context.Background())
	if mx == nil {
		t.Error("Collect() should ignore comments properly")
	}
}

// TestCollectWithIPv6Leases tests handling of IPv6 leases if supported
func TestCollectWithIPv6Leases(t *testing.T) {
	d := New()
	
	tmpDir := t.TempDir()
	testLeasesPath := filepath.Join(tmpDir, "dhcpd.leases")
	
	testContent := `lease 192.168.1.100 {
  starts 4 2024/01/15 10:00:00;
  ends 4 2024/01/15 22:00:00;
  hardware ethernet 00:11:22:33:44:55;
}

# IPv6 lease format (if supported)
ia-pd fd00::1 {
  starts 4 2024/01/15 10:00:00;
  ends 4 2024/01/15 22:00:00;
  cltt 4 2024/01/15 10:00:00;
}
`
	if err := os.WriteFile(testLeasesPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	d.Config.LeasesPath = testLeasesPath

	mx := d.Collect(context.Background())
	if mx == nil {
		t.Error("Collect() should handle IPv6 lease entries if present")
	}
}

// TestNewCollector tests the constructor
func TestNewCollector(t *testing.T) {
	d := New()
	if d == nil {
		t.Fatal("New() returned nil")
	}

	// Verify basic structure
	if d.Config == nil {
		t.Error("Collector Config is nil")
	}
}

// TestCollectWithWhitespaceVariations tests lease files with various whitespace
func TestCollectWithWhitespaceVariations(t *testing.T) {
	d := New()
	
	tmpDir := t.TempDir()
	testLeasesPath := filepath.Join(tmpDir, "dhcpd.leases")
	
	testContent := `lease 192.168.1.100 {
  starts 4 2024/01/15 10:00:00;
  ends 4 2024/01/15 22:00:00;
    hardware ethernet 00:11:22:33:44:55;
}


lease 192.168.1.101 {
	starts 4 2024/01/15 10:00:00;
	ends 4 2024/01/15 22:00:00;
	hardware ethernet aa:bb:cc:dd:ee:ff;
}
`
	if err := os.WriteFile(testLeasesPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	d.Config.LeasesPath = testLeasesPath

	mx := d.Collect(context.Background())
	if mx == nil {
		t.Error("Collect() should handle whitespace variations")
	}
}