package dockerhost

import (
	"net"
	"testing"
)

// TestNewHostname tests the NewHostname function with valid input
func TestNewHostname_ValidInput(t *testing.T) {
	hostname := "192.168.1.1"
	result := NewHostname(hostname)
	
	if result != hostname {
		t.Errorf("NewHostname(%s) = %s, want %s", hostname, result, hostname)
	}
}

// TestNewHostname tests the NewHostname function with empty string
func TestNewHostname_EmptyString(t *testing.T) {
	hostname := ""
	result := NewHostname(hostname)
	
	if result != hostname {
		t.Errorf("NewHostname(%s) = %s, want empty string", hostname, result)
	}
}

// TestNewHostname tests the NewHostname function with special characters
func TestNewHostname_SpecialCharacters(t *testing.T) {
	hostname := "my-docker-host.local"
	result := NewHostname(hostname)
	
	if result != hostname {
		t.Errorf("NewHostname(%s) = %s, want %s", hostname, result, hostname)
	}
}

// TestNewHostname tests the NewHostname function with numeric hostname
func TestNewHostname_NumericHostname(t *testing.T) {
	hostname := "10.0.0.1"
	result := NewHostname(hostname)
	
	if result != hostname {
		t.Errorf("NewHostname(%s) = %s, want %s", hostname, result, hostname)
	}
}

// TestNewHostname tests the NewHostname function with localhost
func TestNewHostname_Localhost(t *testing.T) {
	hostname := "localhost"
	result := NewHostname(hostname)
	
	if result != hostname {
		t.Errorf("NewHostname(%s) = %s, want %s", hostname, result, hostname)
	}
}

// TestHostname tests the Hostname method
func TestHostname_ReturnsValue(t *testing.T) {
	hostname := "docker.example.com"
	h := NewHostname(hostname)
	result := h.Hostname()
	
	if result != hostname {
		t.Errorf("Hostname() = %s, want %s", result, hostname)
	}
}

// TestHostname tests the Hostname method with empty string
func TestHostname_EmptyString(t *testing.T) {
	h := NewHostname("")
	result := h.Hostname()
	
	if result != "" {
		t.Errorf("Hostname() = %s, want empty string", result)
	}
}

// TestString tests the String method
func TestString_ReturnsHostname(t *testing.T) {
	hostname := "test-host"
	h := NewHostname(hostname)
	result := h.String()
	
	if result != hostname {
		t.Errorf("String() = %s, want %s", result, hostname)
	}
}

// TestString tests the String method with IP address
func TestString_IPAddress(t *testing.T) {
	hostname := "172.17.0.1"
	h := NewHostname(hostname)
	result := h.String()
	
	if result != hostname {
		t.Errorf("String() = %s, want %s", result, hostname)
	}
}

// TestIsLocalhost tests the IsLocalhost method with localhost
func TestIsLocalhost_WithLocalhost(t *testing.T) {
	h := NewHostname("localhost")
	if !h.IsLocalhost() {
		t.Errorf("IsLocalhost() = false, want true for 'localhost'")
	}
}

// TestIsLocalhost tests the IsLocalhost method with non-localhost
func TestIsLocalhost_WithNonLocalhost(t *testing.T) {
	h := NewHostname("example.com")
	if h.IsLocalhost() {
		t.Errorf("IsLocalhost() = true, want false for 'example.com'")
	}
}

// TestIsLocalhost tests the IsLocalhost method with loopback IP
func TestIsLocalhost_WithLoopbackIP(t *testing.T) {
	h := NewHostname("127.0.0.1")
	if !h.IsLocalhost() {
		t.Errorf("IsLocalhost() = false, want true for '127.0.0.1'")
	}
}

// TestIsLocalhost tests the IsLocalhost method with loopback IPv6
func TestIsLocalhost_WithLoopbackIPv6(t *testing.T) {
	h := NewHostname("::1")
	if !h.IsLocalhost() {
		t.Errorf("IsLocalhost() = false, want true for '::1'")
	}
}

// TestIsLocalhost tests the IsLocalhost method with empty string
func TestIsLocalhost_WithEmptyString(t *testing.T) {
	h := NewHostname("")
	if h.IsLocalhost() {
		t.Errorf("IsLocalhost() = true, want false for empty string")
	}
}

// TestIsLocalhost tests the IsLocalhost method with other loopback variations
func TestIsLocalhost_WithLoopbackVariations(t *testing.T) {
	tests := []struct {
		hostname string
		want     bool
	}{
		{"127.0.0.2", true},
		{"127.255.255.255", true},
		{"::2", false},
		{"192.168.1.1", false},
	}

	for _, tt := range tests {
		t.Run(tt.hostname, func(t *testing.T) {
			h := NewHostname(tt.hostname)
			got := h.IsLocalhost()
			if got != tt.want {
				t.Errorf("IsLocalhost() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestIsIP tests the IsIP method with valid IP addresses
func TestIsIP_ValidIPv4(t *testing.T) {
	h := NewHostname("192.168.1.1")
	if !h.IsIP() {
		t.Errorf("IsIP() = false, want true for valid IPv4")
	}
}

// TestIsIP tests the IsIP method with valid IPv6
func TestIsIP_ValidIPv6(t *testing.T) {
	h := NewHostname("2001:db8::1")
	if !h.IsIP() {
		t.Errorf("IsIP() = false, want true for valid IPv6")
	}
}

// TestIsIP tests the IsIP method with hostname
func TestIsIP_WithHostname(t *testing.T) {
	h := NewHostname("example.com")
	if h.IsIP() {
		t.Errorf("IsIP() = true, want false for hostname")
	}
}

// TestIsIP tests the IsIP method with empty string
func TestIsIP_WithEmptyString(t *testing.T) {
	h := NewHostname("")
	if h.IsIP() {
		t.Errorf("IsIP() = true, want false for empty string")
	}
}

// TestIsIP tests the IsIP method with localhost
func TestIsIP_WithLocalhost(t *testing.T) {
	h := NewHostname("localhost")
	if h.IsIP() {
		t.Errorf("IsIP() = true, want false for 'localhost'")
	}
}

// TestIsIP tests the IsIP method with loopback address
func TestIsIP_WithLoopback(t *testing.T) {
	h := NewHostname("127.0.0.1")
	if !h.IsIP() {
		t.Errorf("IsIP() = false, want true for loopback IP")
	}
}

// TestIsIP tests with various IP formats
func TestIsIP_VariousFormats(t *testing.T) {
	tests := []struct {
		hostname string
		want     bool
	}{
		{"10.0.0.1", true},
		{"172.16.0.1", true},
		{"255.255.255.255", true},
		{"0.0.0.0", true},
		{"256.1.1.1", false},
		{"::ffff:192.0.2.1", true},
		{"fe80::", true},
		{"not.an.ip", false},
		{"123", false},
	}

	for _, tt := range tests {
		t.Run(tt.hostname, func(t *testing.T) {
			h := NewHostname(tt.hostname)
			got := h.IsIP()
			if got != tt.want {
				t.Errorf("IsIP() = %v, want %v for %s", got, tt.want, tt.hostname)
			}
		})
	}
}

// TestIP tests the IP method with valid IP address
func TestIP_ValidIPv4(t *testing.T) {
	hostname := "192.168.1.1"
	h := NewHostname(hostname)
	ip := h.IP()
	
	if ip == nil {
		t.Errorf("IP() = nil, want net.IP for valid IPv4")
	}
	if ip.String() != hostname {
		t.Errorf("IP().String() = %s, want %s", ip.String(), hostname)
	}
}

// TestIP tests the IP method with valid IPv6
func TestIP_ValidIPv6(t *testing.T) {
	hostname := "::1"
	h := NewHostname(hostname)
	ip := h.IP()
	
	if ip == nil {
		t.Errorf("IP() = nil, want net.IP for valid IPv6")
	}
	if ip.String() != hostname {
		t.Errorf("IP().String() = %s, want %s", ip.String(), hostname)
	}
}

// TestIP tests the IP method with hostname
func TestIP_WithHostname(t *testing.T) {
	h := NewHostname("example.com")
	ip := h.IP()
	
	if ip != nil {
		t.Errorf("IP() = %v, want nil for hostname", ip)
	}
}

// TestIP tests the IP method with empty string
func TestIP_WithEmptyString(t *testing.T) {
	h := NewHostname("")
	ip := h.IP()
	
	if ip != nil {
		t.Errorf("IP() = %v, want nil for empty string", ip)
	}
}

// TestIP tests the IP method with localhost
func TestIP_WithLocalhost(t *testing.T) {
	h := NewHostname("localhost")
	ip := h.IP()
	
	if ip != nil {
		t.Errorf("IP() = %v, want nil for 'localhost'", ip)
	}
}

// TestIP tests various IP formats
func TestIP_VariousFormats(t *testing.T) {
	tests := []struct {
		hostname string
		want     bool
		wantIP   string
	}{
		{"10.0.0.1", true, "10.0.0.1"},
		{"127.0.0.1", true, "127.0.0.1"},
		{"2001:db8::1", true, "2001:db8::1"},
		{"example.com", false, ""},
		{"", false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.hostname, func(t *testing.T) {
			h := NewHostname(tt.hostname)
			got := h.IP()
			if (got != nil) != tt.want {
				t.Errorf("IP() != nil = %v, want %v for %s", got != nil, tt.want, tt.hostname)
			}
			if tt.want && got.String() != tt.wantIP {
				t.Errorf("IP().String() = %s, want %s", got.String(), tt.wantIP)
			}
		})
	}
}

// TestEqual tests the Equal method with identical hostnames
func TestEqual_Identical(t *testing.T) {
	h1 := NewHostname("test-host")
	h2 := NewHostname("test-host")
	
	if !h1.Equal(h2) {
		t.Errorf("Equal() = false, want true for identical hostnames")
	}
}

// TestEqual tests the Equal method with different hostnames
func TestEqual_Different(t *testing.T) {
	h1 := NewHostname("host1")
	h2 := NewHostname("host2")
	
	if h1.Equal(h2) {
		t.Errorf("Equal() = true, want false for different hostnames")
	}
}

// TestEqual tests the Equal method with empty strings
func TestEqual_BothEmpty(t *testing.T) {
	h1 := NewHostname("")
	h2 := NewHostname("")
	
	if !h1.Equal(h2) {
		t.Errorf("Equal() = false, want true for both empty")
	}
}

// TestEqual tests the Equal method with case sensitivity
func TestEqual_CaseSensitive(t *testing.T) {
	h1 := NewHostname("TestHost")
	h2 := NewHostname("testhost")
	
	if h1.Equal(h2) {
		t.Errorf("Equal() = true, want false (case sensitive)")
	}
}

// TestEqual tests the Equal method with IP addresses
func TestEqual_IPAddresses(t *testing.T) {
	h1 := NewHostname("192.168.1.1")
	h2 := NewHostname("192.168.1.1")
	
	if !h1.Equal(h2) {
		t.Errorf("Equal() = false, want true for identical IPs")
	}
}