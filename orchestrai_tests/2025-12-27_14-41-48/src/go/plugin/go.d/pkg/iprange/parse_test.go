package iprange

import (
	"testing"
)

// TestParseIPRange tests the main IP range parsing function
func TestParseIPRange(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		errMsg  string
	}{
		// Valid IPv4 single IP
		{
			name:    "valid single ipv4",
			input:   "192.168.1.1",
			wantErr: false,
		},
		// Valid IPv4 CIDR
		{
			name:    "valid ipv4 cidr",
			input:   "192.168.1.0/24",
			wantErr: false,
		},
		// Valid IPv4 range with hyphen
		{
			name:    "valid ipv4 range hyphen",
			input:   "192.168.1.1-192.168.1.255",
			wantErr: false,
		},
		// Valid IPv6 single IP
		{
			name:    "valid single ipv6",
			input:   "2001:db8::1",
			wantErr: false,
		},
		// Valid IPv6 CIDR
		{
			name:    "valid ipv6 cidr",
			input:   "2001:db8::/32",
			wantErr: false,
		},
		// Valid IPv6 range with hyphen
		{
			name:    "valid ipv6 range hyphen",
			input:   "2001:db8::1-2001:db8::255",
			wantErr: false,
		},
		// Empty string
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		// Invalid IPv4 - malformed
		{
			name:    "invalid ipv4 malformed",
			input:   "192.168.1.999",
			wantErr: true,
		},
		// Invalid IPv4 - incomplete
		{
			name:    "invalid ipv4 incomplete",
			input:   "192.168.1",
			wantErr: true,
		},
		// Invalid CIDR - bad prefix
		{
			name:    "invalid cidr bad prefix",
			input:   "192.168.1.0/33",
			wantErr: true,
		},
		// Invalid CIDR - bad IP
		{
			name:    "invalid cidr bad ip",
			input:   "999.999.999.999/24",
			wantErr: true,
		},
		// Invalid range - start > end
		{
			name:    "invalid range start greater than end",
			input:   "192.168.1.255-192.168.1.1",
			wantErr: true,
		},
		// Invalid range - bad IP format
		{
			name:    "invalid range bad ip format",
			input:   "192.168.1-192.168.2",
			wantErr: true,
		},
		// Whitespace edge case
		{
			name:    "whitespace only",
			input:   "   ",
			wantErr: true,
		},
		// Localhost
		{
			name:    "localhost",
			input:   "127.0.0.1",
			wantErr: false,
		},
		// All zeros
		{
			name:    "all zeros ipv4",
			input:   "0.0.0.0",
			wantErr: false,
		},
		// All ones
		{
			name:    "broadcast ipv4",
			input:   "255.255.255.255",
			wantErr: false,
		},
		// IPv6 loopback
		{
			name:    "ipv6 loopback",
			input:   "::1",
			wantErr: false,
		},
		// IPv6 all zeros
		{
			name:    "ipv6 all zeros",
			input:   "::",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseIPRange(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIPRange() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result == nil {
				t.Errorf("ParseIPRange() expected non-nil result for valid input")
			}
		})
	}
}

// TestParseIPRangeMultiple tests parsing multiple IP ranges
func TestParseIPRangeMultiple(t *testing.T) {
	tests := []struct {
		name      string
		inputs    []string
		wantErr   bool
		wantCount int
	}{
		{
			name:      "single range",
			inputs:    []string{"192.168.1.0/24"},
			wantErr:   false,
			wantCount: 1,
		},
		{
			name:      "multiple ranges",
			inputs:    []string{"192.168.1.0/24", "10.0.0.0/8", "172.16.0.0/12"},
			wantErr:   false,
			wantCount: 3,
		},
		{
			name:      "empty list",
			inputs:    []string{},
			wantErr:   false,
			wantCount: 0,
		},
		{
			name:      "mixed valid and invalid",
			inputs:    []string{"192.168.1.0/24", "invalid"},
			wantErr:   true,
			wantCount: 1,
		},
		{
			name:      "nil slice",
			inputs:    nil,
			wantErr:   false,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := ParseIPRangeMultiple(tt.inputs)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseIPRangeMultiple() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(results) != tt.wantCount {
				t.Errorf("ParseIPRangeMultiple() got %d results, want %d", len(results), tt.wantCount)
			}
		})
	}
}

// TestParseWithHyphen tests parsing IP ranges with hyphen notation
func TestParseWithHyphen(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid hyphen range",
			input:   "192.168.1.1-192.168.1.100",
			wantErr: false,
		},
		{
			name:    "single hyphen",
			input:   "-",
			wantErr: true,
		},
		{
			name:    "hyphen with empty parts",
			input:   "-192.168.1.1",
			wantErr: true,
		},
		{
			name:    "ipv6 hyphen range",
			input:   "2001:db8::1-2001:db8::100",
			wantErr: false,
		},
		{
			name:    "hyphen with spaces",
			input:   "192.168.1.1 - 192.168.1.100",
			wantErr: true,
		},
		{
			name:    "multiple hyphens",
			input:   "192.168.1.1-192.168.1.50-192.168.1.100",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseWithHyphen(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseWithHyphen() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result == nil {
				t.Errorf("ParseWithHyphen() expected non-nil result for valid input")
			}
		})
	}
}

// TestParseWithCIDR tests parsing IP ranges with CIDR notation
func TestParseWithCIDR(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid ipv4 cidr",
			input:   "192.168.0.0/16",
			wantErr: false,
		},
		{
			name:    "valid ipv6 cidr",
			input:   "2001:db8::/32",
			wantErr: false,
		},
		{
			name:    "ipv4 host (/32)",
			input:   "192.168.1.1/32",
			wantErr: false,
		},
		{
			name:    "ipv6 host (/128)",
			input:   "2001:db8::1/128",
			wantErr: false,
		},
		{
			name:    "no prefix",
			input:   "192.168.1.0",
			wantErr: true,
		},
		{
			name:    "invalid prefix",
			input:   "192.168.1.0/abc",
			wantErr: true,
		},
		{
			name:    "prefix too large ipv4",
			input:   "192.168.1.0/33",
			wantErr: true,
		},
		{
			name:    "prefix too large ipv6",
			input:   "2001:db8::/129",
			wantErr: true,
		},
		{
			name:    "negative prefix",
			input:   "192.168.1.0/-1",
			wantErr: true,
		},
		{
			name:    "slash only",
			input:   "/24",
			wantErr: true,
		},
		{
			name:    "multiple slashes",
			input:   "192.168.1.0/24/8",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseWithCIDR(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseWithCIDR() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && result == nil {
				t.Errorf("ParseWithCIDR() expected non-nil result for valid input")
			}
		})
	}
}

// TestIPRangeContains tests whether an IP is within a range
func TestIPRangeContains(t *testing.T) {
	tests := []struct {
		name    string
		rangeStr string
		ipStr   string
		want    bool
		wantErr bool
	}{
		// IPv4 CIDR tests
		{
			name:    "ipv4 in cidr range",
			rangeStr: "192.168.1.0/24",
			ipStr:   "192.168.1.100",
			want:    true,
			wantErr: false,
		},
		{
			name:    "ipv4 not in cidr range",
			rangeStr: "192.168.1.0/24",
			ipStr:   "192.168.2.100",
			want:    false,
			wantErr: false,
		},
		{
			name:    "ipv4 at network boundary",
			rangeStr: "192.168.1.0/24",
			ipStr:   "192.168.1.0",
			want:    true,
			wantErr: false,
		},
		{
			name:    "ipv4 at broadcast boundary",
			rangeStr: "192.168.1.0/24",
			ipStr:   "192.168.1.255",
			want:    true,
			wantErr: false,
		},
		// IPv6 CIDR tests
		{
			name:    "ipv6 in cidr range",
			rangeStr: "2001:db8::/32",
			ipStr:   "2001:db8::1",
			want:    true,
			wantErr: false,
		},
		{
			name:    "ipv6 not in cidr range",
			rangeStr: "2001:db8::/32",
			ipStr:   "2001:db9::1",
			want:    false,
			wantErr: false,
		},
		// Hyphen range tests
		{
			name:    "ipv4 in hyphen range",
			rangeStr: "192.168.1.1-192.168.1.100",
			ipStr:   "192.168.1.50",
			want:    true,
			wantErr: false,
		},
		{
			name:    "ipv4 not in hyphen range",
			rangeStr: "192.168.1.1-192.168.1.100",
			ipStr:   "192.168.1.101",
			want:    false,
			wantErr: false,
		},
		// Single IP tests
		{
			name:    "single ip match",
			rangeStr: "192.168.1.1",
			ipStr:   "192.168.1.1",
			want:    true,
			wantErr: false,
		},
		{
			name:    "single ip no match",
			rangeStr: "192.168.1.1",
			ipStr:   "192.168.1.2",
			want:    false,
			wantErr: false,
		},
		// Error cases
		{
			name:    "invalid range",
			rangeStr: "invalid",
			ipStr:   "192.168.1.1",
			want:    false,
			wantErr: true,
		},
		{
			name:    "invalid ip",
			rangeStr: "192.168.1.0/24",
			ipStr:   "invalid",
			want:    false,
			wantErr: true,
		},
		{
			name:    "empty range",
			rangeStr: "",
			ipStr:   "192.168.1.1",
			want:    false,
			wantErr: true,
		},
		{
			name:    "empty ip",
			rangeStr: "192.168.1.0/24",
			ipStr:   "",
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IPRangeContains(tt.rangeStr, tt.ipStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("IPRangeContains() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("IPRangeContains() got %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIPRangeOverlap tests whether two IP ranges overlap
func TestIPRangeOverlap(t *testing.T) {
	tests := []struct {
		name     string
		range1   string
		range2   string
		want     bool
		wantErr  bool
	}{
		// Overlapping ranges
		{
			name:    "ipv4 overlapping cidr",
			range1:  "192.168.0.0/16",
			range2:  "192.168.1.0/24",
			want:    true,
			wantErr: false,
		},
		// Non-overlapping ranges
		{
			name:    "ipv4 non-overlapping cidr",
			range1:  "192.168.1.0/24",
			range2:  "192.168.2.0/24",
			want:    false,
			wantErr: false,
		},
		// Same range
		{
			name:    "identical ranges",
			range1:  "192.168.1.0/24",
			range2:  "192.168.1.0/24",
			want:    true,
			wantErr: false,
		},
		// Partial overlap
		{
			name:    "partial hyphen overlap",
			range1:  "192.168.1.1-192.168.1.100",
			range2:  "192.168.1.50-192.168.1.150",
			want:    true,
			wantErr: false,
		},
		// IPv6 tests
		{
			name:    "ipv6 overlapping",
			range1:  "2001:db8::/32",
			range2:  "2001:db8:1::/48",
			want:    true,
			wantErr: false,
		},
		{
			name:    "ipv6 non-overlapping",
			range1:  "2001:db8::/32",
			range2:  "2001:db9::/32",
			want:    false,
			wantErr: false,
		},
		// Error cases
		{
			name:    "invalid range1",
			range1:  "invalid",
			range2:  "192.168.1.0/24",
			want:    false,
			wantErr: true,
		},
		{
			name:    "invalid range2",
			range1:  "192.168.1.0/24",
			range2:  "invalid",
			want:    false,
			wantErr: true,
		},
		{
			name:    "empty ranges",
			range1:  "",
			range2:  "",
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IPRangeOverlap(tt.range1, tt.range2)
			if (err != nil) != tt.wantErr {
				t.Errorf("IPRangeOverlap() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("IPRangeOverlap() got %v, want %v", got, tt.want)
			}
		})
	}
}

// TestValidateIPFormat tests IP validation
func TestValidateIPFormat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// Valid IPv4
		{
			name:    "valid ipv4",
			input:   "192.168.1.1",
			wantErr: false,
		},
		// Valid IPv6
		{
			name:    "valid ipv6",
			input:   "2001:db8::1",
			wantErr: false,
		},
		// Invalid formats
		{
			name:    "invalid ipv4 octet",
			input:   "192.168.1.999",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "whitespace",
			input:   "   ",
			wantErr: true,
		},
		{
			name:    "random string",
			input:   "not an ip",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIPFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIPFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}