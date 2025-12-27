package snmputils

import (
	"bytes"
	"net"
	"testing"
)

// Test ParseOID - parses OID strings into numeric arrays
func TestParseOID(t *testing.T) {
	tests := []struct {
		name      string
		oid       string
		expected  []int
		shouldErr bool
	}{
		{
			name:      "valid simple OID",
			oid:       "1.3.6.1",
			expected:  []int{1, 3, 6, 1},
			shouldErr: false,
		},
		{
			name:      "valid complex OID",
			oid:       "1.3.6.1.2.1.1.1.0",
			expected:  []int{1, 3, 6, 1, 2, 1, 1, 1, 0},
			shouldErr: false,
		},
		{
			name:      "single number OID",
			oid:       "1",
			expected:  []int{1},
			shouldErr: false,
		},
		{
			name:      "OID with leading dot",
			oid:       ".1.3.6",
			expected:  []int{1, 3, 6},
			shouldErr: false,
		},
		{
			name:      "empty OID",
			oid:       "",
			expected:  nil,
			shouldErr: true,
		},
		{
			name:      "invalid non-numeric OID",
			oid:       "1.3.abc.1",
			expected:  nil,
			shouldErr: true,
		},
		{
			name:      "OID with spaces",
			oid:       "1 3 6 1",
			expected:  nil,
			shouldErr: true,
		},
		{
			name:      "negative number in OID",
			oid:       "1.-3.6.1",
			expected:  nil,
			shouldErr: true,
		},
		{
			name:      "very large OID numbers",
			oid:       "1.3.6.1.2147483647",
			expected:  []int{1, 3, 6, 1, 2147483647},
			shouldErr: false,
		},
		{
			name:      "OID with trailing dot",
			oid:       "1.3.6.",
			expected:  nil,
			shouldErr: true,
		},
		{
			name:      "OID with consecutive dots",
			oid:       "1..3.6",
			expected:  nil,
			shouldErr: true,
		},
		{
			name:      "zero in OID",
			oid:       "1.0.6.1",
			expected:  []int{1, 0, 6, 1},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseOID(tt.oid)

			if (err != nil) != tt.shouldErr {
				t.Fatalf("expected error: %v, got: %v", tt.shouldErr, err != nil)
			}

			if !tt.shouldErr {
				if len(result) != len(tt.expected) {
					t.Fatalf("expected length %d, got %d", len(tt.expected), len(result))
				}
				for i, v := range result {
					if v != tt.expected[i] {
						t.Errorf("at index %d: expected %d, got %d", i, tt.expected[i], v)
					}
				}
			}
		})
	}
}

// Test JoinOID - joins numeric arrays back into OID strings
func TestJoinOID(t *testing.T) {
	tests := []struct {
		name     string
		oid      []int
		expected string
	}{
		{
			name:     "simple OID",
			oid:      []int{1, 3, 6, 1},
			expected: "1.3.6.1",
		},
		{
			name:     "complex OID",
			oid:      []int{1, 3, 6, 1, 2, 1, 1, 1, 0},
			expected: "1.3.6.1.2.1.1.1.0",
		},
		{
			name:     "single number OID",
			oid:      []int{1},
			expected: "1",
		},
		{
			name:     "empty OID",
			oid:      []int{},
			expected: "",
		},
		{
			name:     "OID with zero",
			oid:      []int{1, 0, 6, 1},
			expected: "1.0.6.1",
		},
		{
			name:     "large numbers",
			oid:      []int{1, 3, 6, 1, 2147483647},
			expected: "1.3.6.1.2147483647",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JoinOID(tt.oid)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// Test ParseOIDAndJoinOID round-trip
func TestParseOIDAndJoinOIDRoundTrip(t *testing.T) {
	tests := []string{
		"1.3.6.1",
		"1.3.6.1.2.1.1.1.0",
		"1",
		".1.3.6.1",
		"1.0.6.1",
		"1.2.3.4.5.6.7.8.9.10",
	}

	for _, original := range tests {
		t.Run(original, func(t *testing.T) {
			parsed, err := ParseOID(original)
			if err != nil {
				t.Fatalf("failed to parse OID: %v", err)
			}

			joined := JoinOID(parsed)
			expected := original
			if original[0] == '.' {
				expected = original[1:]
			}

			if joined != expected {
				t.Errorf("round-trip failed: %s -> %v -> %s", original, parsed, joined)
			}
		})
	}
}

// Test IsIPAddr
func TestIsIPAddr(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{
			name:     "valid IPv4",
			ip:       "192.168.1.1",
			expected: true,
		},
		{
			name:     "valid IPv4 localhost",
			ip:       "127.0.0.1",
			expected: true,
		},
		{
			name:     "valid IPv4 zero",
			ip:       "0.0.0.0",
			expected: true,
		},
		{
			name:     "valid IPv4 max",
			ip:       "255.255.255.255",
			expected: true,
		},
		{
			name:     "valid IPv6",
			ip:       "::1",
			expected: true,
		},
		{
			name:     "valid IPv6 full",
			ip:       "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			expected: true,
		},
		{
			name:     "valid IPv6 compressed",
			ip:       "2001:db8::1",
			expected: true,
		},
		{
			name:     "invalid IP",
			ip:       "256.256.256.256",
			expected: false,
		},
		{
			name:     "empty string",
			ip:       "",
			expected: false,
		},
		{
			name:     "hostname",
			ip:       "example.com",
			expected: false,
		},
		{
			name:     "partial IPv4",
			ip:       "192.168.1",
			expected: false,
		},
		{
			name:     "IPv4 with port",
			ip:       "192.168.1.1:8080",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsIPAddr(tt.ip)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test TargetAddr
func TestTargetAddr(t *testing.T) {
	tests := []struct {
		name        string
		ip          string
		port        string
		expected    string
		shouldMatch bool
	}{
		{
			name:        "IPv4 with port",
			ip:          "192.168.1.1",
			port:        "161",
			expected:    "192.168.1.1:161",
			shouldMatch: true,
		},
		{
			name:        "IPv4 localhost with port",
			ip:          "127.0.0.1",
			port:        "8080",
			expected:    "127.0.0.1:8080",
			shouldMatch: true,
		},
		{
			name:        "IPv6 with port",
			ip:          "::1",
			port:        "161",
			expected:    "[::1]:161",
			shouldMatch: true,
		},
		{
			name:        "IPv6 full with port",
			ip:          "2001:db8::1",
			port:        "161",
			expected:    "[2001:db8::1]:161",
			shouldMatch: true,
		},
		{
			name:        "empty IP",
			ip:          "",
			port:        "161",
			expected:    ":161",
			shouldMatch: true,
		},
		{
			name:        "empty port",
			ip:          "192.168.1.1",
			port:        "",
			expected:    "192.168.1.1:",
			shouldMatch: true,
		},
		{
			name:        "both empty",
			ip:          "",
			port:        "",
			expected:    ":",
			shouldMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TargetAddr(tt.ip, tt.port)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// Test IsIPv4
func TestIsIPv4(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{
			name:     "valid IPv4",
			ip:       "192.168.1.1",
			expected: true,
		},
		{
			name:     "IPv4 localhost",
			ip:       "127.0.0.1",
			expected: true,
		},
		{
			name:     "IPv4 zero",
			ip:       "0.0.0.0",
			expected: true,
		},
		{
			name:     "IPv4 max",
			ip:       "255.255.255.255",
			expected: true,
		},
		{
			name:     "IPv6",
			ip:       "::1",
			expected: false,
		},
		{
			name:     "invalid IP",
			ip:       "256.256.256.256",
			expected: false,
		},
		{
			name:     "empty string",
			ip:       "",
			expected: false,
		},
		{
			name:     "hostname",
			ip:       "example.com",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsIPv4(tt.ip)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test IsIPv6
func TestIsIPv6(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{
			name:     "IPv6 loopback",
			ip:       "::1",
			expected: true,
		},
		{
			name:     "IPv6 full",
			ip:       "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			expected: true,
		},
		{
			name:     "IPv6 compressed",
			ip:       "2001:db8::1",
			expected: true,
		},
		{
			name:     "IPv6 unspecified",
			ip:       "::",
			expected: true,
		},
		{
			name:     "IPv4",
			ip:       "192.168.1.1",
			expected: false,
		},
		{
			name:     "invalid IPv6",
			ip:       "gggg::1",
			expected: false,
		},
		{
			name:     "empty string",
			ip:       "",
			expected: false,
		},
		{
			name:     "hostname",
			ip:       "example.com",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsIPv6(tt.ip)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test ParseIP
func TestParseIP(t *testing.T) {
	tests := []struct {
		name          string
		ip            string
		expectedIP    net.IP
		expectedError bool
	}{
		{
			name:          "valid IPv4",
			ip:            "192.168.1.1",
			expectedIP:    net.ParseIP("192.168.1.1"),
			expectedError: false,
		},
		{
			name:          "valid IPv6",
			ip:            "::1",
			expectedIP:    net.ParseIP("::1"),
			expectedError: false,
		},
		{
			name:          "invalid IP",
			ip:            "not-an-ip",
			expectedIP:    nil,
			expectedError: true,
		},
		{
			name:          "empty string",
			ip:            "",
			expectedIP:    nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseIP(tt.ip)

			if (err != nil) != tt.expectedError {
				t.Fatalf("expected error: %v, got: %v", tt.expectedError, err != nil)
			}

			if !tt.expectedError {
				if !bytes.Equal(result, tt.expectedIP) {
					t.Errorf("expected %v, got %v", tt.expectedIP, result)
				}
			}
		})
	}
}

// Test ParsePort
func TestParsePort(t *testing.T) {
	tests := []struct {
		name          string
		port          string
		expected      uint16
		expectedError bool
	}{
		{
			name:          "valid port 161",
			port:          "161",
			expected:      161,
			expectedError: false,
		},
		{
			name:          "valid port 0",
			port:          "0",
			expected:      0,
			expectedError: false,
		},
		{
			name:          "valid port 65535",
			port:          "65535",
			expected:      65535,
			expectedError: false,
		},
		{
			name:          "port too high",
			port:          "65536",
			expected:      0,
			expectedError: true,
		},
		{
			name:          "negative port",
			port:          "-1",
			expected:      0,
			expectedError: true,
		},
		{
			name:          "non-numeric port",
			port:          "abc",
			expected:      0,
			expectedError: true,
		},
		{
			name:          "empty string",
			port:          "",
			expected:      0,
			expectedError: true,
		},
		{
			name:          "port with leading zeros",
			port:          "0161",
			expected:      161,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParsePort(tt.port)

			if (err != nil) != tt.expectedError {
				t.Fatalf("expected error: %v, got: %v", tt.expectedError, err != nil)
			}

			if !tt.expectedError && result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}

// Test IsHexString
func TestIsHexString(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected bool
	}{
		{
			name:     "valid hex lowercase",
			s:        "deadbeef",
			expected: true,
		},
		{
			name:     "valid hex uppercase",
			s:        "DEADBEEF",
			expected: true,
		},
		{
			name:     "valid hex mixed",
			s:        "DeAdBeEf",
			expected: true,
		},
		{
			name:     "valid hex with numbers",
			s:        "a1b2c3d4",
			expected: true,
		},
		{
			name:     "single hex char",
			s:        "f",
			expected: true,
		},
		{
			name:     "invalid hex char",
			s:        "gg",
			expected: false,
		},
		{
			name:     "empty string",
			s:        "",
			expected: false,
		},
		{
			name:     "hex with spaces",
			s:        "dead beef",
			expected: false,
		},
		{
			name:     "hex with special chars",
			s:        "dead-beef",
			expected: false,
		},
		{
			name:     "decimal numbers",
			s:        "123456",
			expected: true,
		},
		{
			name:     "single invalid char in hex",
			s:        "deadbeez",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsHexString(tt.s)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test IsNumString
func TestIsNumString(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected bool
	}{
		{
			name:     "single digit",
			s:        "5",
			expected: true,
		},
		{
			name:     "multiple digits",
			s:        "12345",
			expected: true,
		},
		{
			name:     "zero",
			s:        "0",
			expected: true,
		},
		{
			name:     "large number",
			s:        "9999999999",
			expected: true,
		},
		{
			name:     "negative number",
			s:        "-123",
			expected: false,
		},
		{
			name:     "with decimal",
			s:        "123.45",
			expected: false,
		},
		{
			name:     "with spaces",
			s:        "123 456",
			expected: false,
		},
		{
			name:     "with letters",
			s:        "123abc",
			expected: false,
		},
		{
			name:     "empty string",
			s:        "",
			expected: false,
		},
		{
			name:     "only spaces",
			s:        "   ",
			expected: false,
		},
		{
			name:     "hex numbers",
			s:        "abcdef",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsNumString(tt.s)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// Test Ipv6Brackets
func TestIpv6Brackets(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected string
	}{
		{
			name:     "IPv6 without brackets",
			ip:       "::1",
			expected: "[::1]",
		},
		{
			name:     "IPv6 full without brackets",
			ip:       "2001:db8::1",
			expected: "[2001:db8::1]",
		},
		{
			name:     "IPv6 already with brackets",
			ip:       "[::1]",
			expected: "[::1]",
		},
		{
			name:     "IPv4",
			ip:       "192.168.1.1",
			expected: "192.168.1.1",
		},
		{
			name:     "empty string",
			ip:       "",
			expected: "",
		},
		{
			name:     "IPv6 unspecified",
			ip:       "::",
			expected: "[::]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Ipv6Brackets(tt.ip)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// Test RemoveIpv6Brackets
func TestRemoveIpv6Brackets(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected string
	}{
		{
			name:     "IPv6 with brackets",
			ip:       "[::1]",
			expected: "::1",
		},
		{
			name:     "IPv6 full with brackets",
			ip:       "[2001:db8::1]",
			expected: "2001:db8::1",
		},
		{
			name:     "IPv6 without brackets",
			ip:       "::1",
			expected: "::1",
		},
		{
			name:     "IPv4",
			ip:       "192.168.1.1",
			expected: "192.168.1.1",
		},
		{
			name:     "empty string",
			ip:       "",
			expected: "",
		},
		{
			name:     "only opening bracket",
			ip:       "[::1",
			expected: "[::1",
		},
		{
			name:     "only closing bracket",
			ip:       "::1]",
			expected: "::1]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveIpv6Brackets(tt.ip)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// Test SplitHostPort
func TestSplitHostPort(t *testing.T) {
	tests := []struct {
		name          string
		hostport      string
		expectedHost  string
		expectedPort  string
		expectedError bool
	}{
		{
			name:          "IPv4 with port",
			hostport:      "192.168.1.1:161",
			expectedHost:  "192.168.1.1",
			expectedPort:  "161",
			expectedError: false,
		},
		{
			name:          "IPv6 with port",
			hostport:      "[::1]:161",
			expectedHost:  "::1",
			expectedPort:  "161",
			expectedError: false,
		},
		{
			name:          "IPv6 full with port",
			hostport:      "[2001:db8::1]:8080",
			expectedHost:  "2001:db8::1",
			expectedPort:  "8080",
			expectedError: false,
		},
		{
			name:          "hostname with port",
			hostport:      "example.com:80",
			expectedHost:  "example.com",
			expectedPort:  "80",
			expectedError: false,
		},
		{
			name:          "localhost with port",
			hostport:      "localhost:8080",
			expectedHost:  "localhost",
			expectedPort:  "8080",
			expectedError: false,
		},
		{
			name:          "no port",
			hostport:      "192.168.1.1",
			expectedHost:  "192.168.1.1",
			expectedPort:  "",
			expectedError: false,
		},
		{
			name:          "IPv6 without port",
			hostport:      "[::1]",
			expectedHost:  "::1",
			expectedPort:  "",
			expectedError: false,
		},
		{
			name:          "malformed IPv6 with port",
			hostport:      "::1:161",
			expectedHost:  "",
			expectedPort:  "",
			expectedError: true,
		},
		{
			name:          "empty string",
			hostport:      "",
			expectedHost:  "",
			expectedPort:  "",
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			host, port, err := SplitHostPort(tt.hostport)

			if (err != nil) != tt.expectedError {
				t.Fatalf("expected error: %v, got: %v", tt.expectedError, err != nil)
			}

			if !tt.expectedError {
				if host != tt.expectedHost {
					t.Errorf("host: expected %s, got %s", tt.expectedHost, host)
				}
				if port != tt.expectedPort {
					t.Errorf("port: expected %s, got %s", tt.expectedPort, port)
				}
			}
		})
	}
}