package isc_dhcpd

import (
	"testing"
)

// TestParseLeases tests the main lease parsing function
func TestParseLeases(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedCount int
		shouldError   bool
	}{
		{
			name:          "should parse empty file",
			input:         "",
			expectedCount: 0,
			shouldError:   false,
		},
		{
			name: "should parse single lease",
			input: `lease 192.168.1.100 {
  starts 1 2023/01/01 12:00:00;
  ends 1 2023/01/01 13:00:00;
  binding state active;
  hardware ethernet aa:bb:cc:dd:ee:ff;
}`,
			expectedCount: 1,
			shouldError:   false,
		},
		{
			name: "should parse multiple leases",
			input: `lease 192.168.1.100 {
  starts 1 2023/01/01 12:00:00;
  ends 1 2023/01/01 13:00:00;
  binding state active;
}
lease 192.168.1.101 {
  starts 1 2023/01/01 12:00:00;
  ends 1 2023/01/01 13:00:00;
  binding state released;
}`,
			expectedCount: 2,
			shouldError:   false,
		},
		{
			name:          "should handle whitespace",
			input:         "   \n  \t  \n  ",
			expectedCount: 0,
			shouldError:   false,
		},
		{
			name: "should parse lease with all fields",
			input: `lease 192.168.1.50 {
  starts 1 2023/01/01 00:00:00;
  ends 1 2023/01/01 23:59:59;
  binding state active;
  hardware ethernet 11:22:33:44:55:66;
  uid "device-uid";
  set vendor-string "vendor-name";
}`,
			expectedCount: 1,
			shouldError:   false,
		},
		{
			name: "should handle malformed lease gracefully",
			input: `lease 192.168.1.100 {
  invalid line here
}`,
			expectedCount: 1,
			shouldError:   false,
		},
		{
			name:          "should handle files with comments",
			input:         "# This is a comment\nlease 192.168.1.100 {\n}\n# Another comment",
			expectedCount: 1,
			shouldError:   false,
		},
		{
			name:          "should parse lease with missing semicolons",
			input:         "lease 192.168.1.100 {\nhardware ethernet aa:bb:cc:dd:ee:ff\n}",
			expectedCount: 1,
			shouldError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := NewDhcpdCollector()
			leases, err := dc.ParseLeases(tt.input)

			if (err != nil) != tt.shouldError {
				t.Errorf("ParseLeases() error = %v, shouldError %v", err, tt.shouldError)
			}

			if len(leases) != tt.expectedCount {
				t.Errorf("ParseLeases() got %d leases, expected %d", len(leases), tt.expectedCount)
			}
		})
	}
}

// TestParseLeasesStates tests parsing different binding states
func TestParseLeasesStates(t *testing.T) {
	tests := []struct {
		name     string
		state    string
		validate func(*Lease) bool
	}{
		{
			name:  "should parse active state",
			state: "active",
			validate: func(l *Lease) bool {
				return l.State == "active"
			},
		},
		{
			name:  "should parse released state",
			state: "released",
			validate: func(l *Lease) bool {
				return l.State == "released"
			},
		},
		{
			name:  "should parse abandoned state",
			state: "abandoned",
			validate: func(l *Lease) bool {
				return l.State == "abandoned"
			},
		},
		{
			name:  "should parse free state",
			state: "free",
			validate: func(l *Lease) bool {
				return l.State == "free"
			},
		},
		{
			name:  "should parse reserved state",
			state: "reserved",
			validate: func(l *Lease) bool {
				return l.State == "reserved"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := `lease 192.168.1.100 {
  binding state ` + tt.state + `;
}`
			dc := NewDhcpdCollector()
			leases, err := dc.ParseLeases(input)

			if err != nil {
				t.Errorf("ParseLeases() error = %v", err)
				return
			}

			if len(leases) == 0 {
				t.Error("ParseLeases() returned empty leases")
				return
			}

			if !tt.validate(leases[0]) {
				t.Errorf("State validation failed for %s", tt.state)
			}
		})
	}
}

// TestParseIPAddress tests IP address parsing
func TestParseIPAddress(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectedIP string
	}{
		{
			name:       "should parse IPv4 address",
			input:      "lease 192.168.1.100 {",
			expectedIP: "192.168.1.100",
		},
		{
			name:       "should parse different IPv4",
			input:      "lease 10.0.0.1 {",
			expectedIP: "10.0.0.1",
		},
		{
			name:       "should parse edge case IPv4",
			input:      "lease 255.255.255.255 {",
			expectedIP: "255.255.255.255",
		},
		{
			name:       "should parse zero IP",
			input:      "lease 0.0.0.0 {",
			expectedIP: "0.0.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := tt.input + "\nbinding state active;\n}"
			dc := NewDhcpdCollector()
			leases, err := dc.ParseLeases(input)

			if err != nil {
				t.Errorf("ParseLeases() error = %v", err)
				return
			}

			if len(leases) == 0 {
				t.Error("ParseLeases() returned empty leases")
				return
			}

			if leases[0].IP != tt.expectedIP {
				t.Errorf("Got IP %s, expected %s", leases[0].IP, tt.expectedIP)
			}
		})
	}
}

// TestParseHardwareAddress tests MAC address parsing
func TestParseHardwareAddress(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedMAC  string
	}{
		{
			name:         "should parse standard MAC",
			input:        "hardware ethernet aa:bb:cc:dd:ee:ff;",
			expectedMAC:  "aa:bb:cc:dd:ee:ff",
		},
		{
			name:         "should parse uppercase MAC",
			input:        "hardware ethernet AA:BB:CC:DD:EE:FF;",
			expectedMAC:  "AA:BB:CC:DD:EE:FF",
		},
		{
			name:         "should parse zero MAC",
			input:        "hardware ethernet 00:00:00:00:00:00;",
			expectedMAC:  "00:00:00:00:00:00",
		},
		{
			name:         "should parse all F MAC",
			input:        "hardware ethernet ff:ff:ff:ff:ff:ff;",
			expectedMAC:  "ff:ff:ff:ff:ff:ff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := `lease 192.168.1.100 {
` + tt.input + `
binding state active;
}`
			dc := NewDhcpdCollector()
			leases, err := dc.ParseLeases(input)

			if err != nil {
				t.Errorf("ParseLeases() error = %v", err)
				return
			}

			if len(leases) == 0 {
				t.Error("ParseLeases() returned empty leases")
				return
			}

			if leases[0].MAC != tt.expectedMAC {
				t.Errorf("Got MAC %s, expected %s", leases[0].MAC, tt.expectedMAC)
			}
		})
	}
}

// TestParseTimestamps tests timestamp parsing
func TestParseTimestamps(t *testing.T) {
	tests := []struct {
		name    string
		starts  string
		ends    string
		hasDates bool
	}{
		{
			name:    "should parse valid timestamps",
			starts:  "1 2023/01/01 12:00:00",
			ends:    "1 2023/01/01 13:00:00",
			hasDates: true,
		},
		{
			name:    "should parse different dates",
			starts:  "1 2023/12/31 23:59:59",
			ends:    "2 2024/01/01 00:00:00",
			hasDates: true,
		},
		{
			name:    "should handle epoch times",
			starts:  "0 1970/01/01 00:00:00",
			ends:    "0 1970/01/01 00:00:01",
			hasDates: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := `lease 192.168.1.100 {
  starts ` + tt.starts + `;
  ends ` + tt.ends + `;
  binding state active;
}`
			dc := NewDhcpdCollector()
			leases, err := dc.ParseLeases(input)

			if err != nil {
				t.Errorf("ParseLeases() error = %v", err)
				return
			}

			if len(leases) == 0 {
				t.Error("ParseLeases() returned empty leases")
				return
			}

			if tt.hasDates {
				if leases[0].Starts == "" {
					t.Error("Starts timestamp not parsed")
				}
				if leases[0].Ends == "" {
					t.Error("Ends timestamp not parsed")
				}
			}
		})
	}
}

// TestParseEdgeCases tests edge cases in parsing
func TestParseEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		count   int
	}{
		{
			name:    "should handle lease without closing brace",
			input:   "lease 192.168.1.100 {",
			count:   1,
		},
		{
			name:    "should handle extra closing braces",
			input:   "lease 192.168.1.100 { binding state active; } }",
			count:   1,
		},
		{
			name:    "should handle nested braces",
			input:   "lease 192.168.1.100 { { binding state active; } }",
			count:   1,
		},
		{
			name:    "should handle very long MAC address line",
			input:   "lease 192.168.1.100 { hardware ethernet aa:bb:cc:dd:ee:ff; binding state active; }",
			count:   1,
		},
		{
			name:    "should handle multiple spaces",
			input:   "lease    192.168.1.100    {   binding state active;   }",
			count:   1,
		},
		{
			name:    "should handle tabs",
			input:   "lease\t192.168.1.100\t{\tbinding state active;\t}",
			count:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := NewDhcpdCollector()
			leases, err := dc.ParseLeases(tt.input)

			if err != nil {
				t.Errorf("ParseLeases() error = %v", err)
				return
			}

			if len(leases) != tt.count {
				t.Errorf("Got %d leases, expected %d", len(leases), tt.count)
			}
		})
	}
}

// TestParseLeasesRobustness tests robustness against various inputs
func TestParseLeasesRobustness(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "should handle nil input gracefully",
			input: "",
		},
		{
			name:  "should handle only whitespace",
			input: "   \n\t\n   ",
		},
		{
			name:  "should handle very large input",
			input: generateLargeInput(1000),
		},
		{
			name:  "should handle special characters in comments",
			input: "# Special chars: !@#$%^&*()\nlease 192.168.1.1 { binding state active; }",
		},
		{
			name:  "should handle unicode",
			input: "# Unicode: 你好\nlease 192.168.1.1 { binding state active; }",
		},
		{
			name:  "should handle mixed line endings",
			input: "lease 192.168.1.1 {\r\nbinding state active;\r\n}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := NewDhcpdCollector()
			_, err := dc.ParseLeases(tt.input)

			// Should not panic, error is acceptable
			if err != nil && tt.name != "should handle very large input" {
				// Log error but don't fail for robustness tests
				t.Logf("ParseLeases() error = %v", err)
			}
		})
	}
}

// TestParseLeaseFieldExtraction tests extraction of individual fields
func TestParseLeaseFieldExtraction(t *testing.T) {
	input := `lease 192.168.1.100 {
  starts 1 2023/01/01 12:00:00;
  ends 1 2023/01/01 13:00:00;
  binding state active;
  hardware ethernet aa:bb:cc:dd:ee:ff;
  uid "device-id";
  set vendor-string "vendor";
  client-hostname "hostname";
}`

	dc := NewDhcpdCollector()
	leases, err := dc.ParseLeases(input)

	if err != nil {
		t.Errorf("ParseLeases() error = %v", err)
		return
	}

	if len(leases) != 1 {
		t.Fatalf("Expected 1 lease, got %d", len(leases))
	}

	lease := leases[0]

	tests := []struct {
		name     string
		validate func() bool
	}{
		{"IP extracted", func() bool { return lease.IP == "192.168.1.100" }},
		{"MAC extracted", func() bool { return lease.MAC == "aa:bb:cc:dd:ee:ff" }},
		{"State extracted", func() bool { return lease.State == "active" }},
		{"Start time extracted", func() bool { return lease.Starts != "" }},
		{"End time extracted", func() bool { return lease.Ends != "" }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.validate() {
				t.Error("Field extraction failed")
			}
		})
	}
}

// TestParseConsecutiveLeases tests parsing many consecutive leases
func TestParseConsecutiveLeases(t *testing.T) {
	input := ""
	expectedCount := 100

	for i := 0; i < expectedCount; i++ {
		input += `lease 192.168.1.` + string(rune(i%256)) + ` {
  binding state active;
}
`
	}

	dc := NewDhcpdCollector()
	leases, err := dc.ParseLeases(input)

	if err != nil {
		t.Errorf("ParseLeases() error = %v", err)
		return
	}

	if len(leases) != expectedCount {
		t.Errorf("Got %d leases, expected %d", len(leases), expectedCount)
	}
}

// TestParseEmptyLeaseBlock tests parsing lease block with no properties
func TestParseEmptyLeaseBlock(t *testing.T) {
	input := "lease 192.168.1.100 { }"

	dc := NewDhcpdCollector()
	leases, err := dc.ParseLeases(input)

	if err != nil {
		t.Errorf("ParseLeases() error = %v", err)
		return
	}

	if len(leases) != 1 {
		t.Errorf("Expected 1 lease, got %d", len(leases))
		return
	}

	lease := leases[0]
	if lease.IP != "192.168.1.100" {
		t.Errorf("IP not parsed correctly: %s", lease.IP)
	}
}

// TestParseLeaseWithDifferentSeparators tests parsing with various whitespace
func TestParseLeaseWithDifferentSeparators(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "newlines and spaces",
			input: "lease 192.168.1.100 {\n  binding state active;\n}",
		},
		{
			name:  "tabs and spaces",
			input: "lease 192.168.1.100 {\t binding state active;\t}",
		},
		{
			name:  "minimal whitespace",
			input: "lease 192.168.1.100{binding state active;}",
		},
		{
			name:  "excessive whitespace",
			input: "lease   192.168.1.100   {   binding   state   active   ;   }",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := NewDhcpdCollector()
			leases, err := dc.ParseLeases(tt.input)

			if err != nil {
				t.Errorf("ParseLeases() error = %v", err)
				return
			}

			if len(leases) != 1 {
				t.Errorf("Expected 1 lease, got %d", len(leases))
			}
		})
	}
}

// Helper function to generate large input
func generateLargeInput(count int) string {
	input := ""
	for i := 0; i < count; i++ {
		input += "lease 192.168.1." + string(rune(i%256)) + " { binding state active; }\n"
	}
	return input
}