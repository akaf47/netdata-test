package iprange

import (
	"net"
	"testing"
)

// TestNewIterator tests creating a new iterator with valid IP range
func TestNewIterator_ValidIPv4Range(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("192.168.1.0/24")
	iter := NewIterator(ipnet)
	
	if iter == nil {
		t.Errorf("NewIterator() = nil, want iterator")
	}
}

// TestNewIterator tests creating a new iterator with IPv6 range
func TestNewIterator_ValidIPv6Range(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("2001:db8::/32")
	iter := NewIterator(ipnet)
	
	if iter == nil {
		t.Errorf("NewIterator() = nil, want iterator")
	}
}

// TestNewIterator tests creating iterator with nil network
func TestNewIterator_NilNetwork(t *testing.T) {
	iter := NewIterator(nil)
	
	if iter == nil {
		t.Errorf("NewIterator() = nil, want iterator even with nil network")
	}
}

// TestNext tests the Next method with valid range
func TestNext_IPv4SequentialValues(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("192.168.1.1/32")
	iter := NewIterator(ipnet)
	
	ip := iter.Next()
	if ip == nil {
		t.Errorf("Next() = nil, want IP address")
	}
	if ip.String() != "192.168.1.1" {
		t.Errorf("Next() = %s, want 192.168.1.1", ip.String())
	}
	
	ip2 := iter.Next()
	if ip2 != nil {
		t.Errorf("Next() = %v, want nil after exhausting range", ip2)
	}
}

// TestNext tests multiple iterations
func TestNext_MultipleIterations(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("10.0.0.0/30")
	iter := NewIterator(ipnet)
	
	expected := []string{"10.0.0.0", "10.0.0.1", "10.0.0.2", "10.0.0.3"}
	for _, exp := range expected {
		ip := iter.Next()
		if ip == nil {
			t.Errorf("Next() = nil, want %s", exp)
			break
		}
		if ip.String() != exp {
			t.Errorf("Next() = %s, want %s", ip.String(), exp)
		}
	}
	
	// Should return nil when exhausted
	if ip := iter.Next(); ip != nil {
		t.Errorf("Next() = %v, want nil when exhausted", ip)
	}
}

// TestNext tests with single host (/32)
func TestNext_SingleHost(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("172.16.0.1/32")
	iter := NewIterator(ipnet)
	
	ip := iter.Next()
	if ip == nil {
		t.Errorf("Next() = nil, want single IP")
	}
	if ip.String() != "172.16.0.1" {
		t.Errorf("Next() = %s, want 172.16.0.1", ip.String())
	}
	
	ip2 := iter.Next()
	if ip2 != nil {
		t.Errorf("Next() = %v, want nil after single IP", ip2)
	}
}

// TestNext tests with /24 range
func TestNext_Class24Range(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("192.168.0.0/24")
	iter := NewIterator(ipnet)
	
	count := 0
	for {
		ip := iter.Next()
		if ip == nil {
			break
		}
		count++
	}
	
	if count != 256 {
		t.Errorf("Iterated %d IPs, want 256 in /24 range", count)
	}
}

// TestNext tests with very large range (/16)
func TestNext_Class16Range(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("10.0.0.0/16")
	iter := NewIterator(ipnet)
	
	count := 0
	var lastIP net.IP
	for {
		ip := iter.Next()
		if ip == nil {
			break
		}
		lastIP = ip
		count++
		if count > 1000 { // Limit for test performance
			break
		}
	}
	
	if count != 1000 {
		t.Errorf("Iterated %d IPs, want 1000 before stopping", count)
	}
	if lastIP == nil {
		t.Errorf("lastIP = nil, want IP address")
	}
}

// TestNext tests IPv6 iteration
func TestNext_IPv6Iteration(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("2001:db8::0/127")
	iter := NewIterator(ipnet)
	
	expected := []string{"2001:db8::", "2001:db8::1"}
	for _, exp := range expected {
		ip := iter.Next()
		if ip == nil {
			t.Errorf("Next() = nil, want %s", exp)
			break
		}
		if ip.String() != exp {
			t.Errorf("Next() = %s, want %s", ip.String(), exp)
		}
	}
	
	if ip := iter.Next(); ip != nil {
		t.Errorf("Next() = %v, want nil when exhausted", ip)
	}
}

// TestNext tests with nil network
func TestNext_WithNilNetwork(t *testing.T) {
	iter := NewIterator(nil)
	
	ip := iter.Next()
	if ip != nil {
		t.Errorf("Next() = %v, want nil for nil network", ip)
	}
}

// TestNext tests sequential IPv4 calls
func TestNext_SequentialCalls(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("127.0.0.1/32")
	iter := NewIterator(ipnet)
	
	ip1 := iter.Next()
	ip2 := iter.Next()
	ip3 := iter.Next()
	
	if ip1 == nil {
		t.Errorf("First Next() = nil, want IP")
	} else if ip1.String() != "127.0.0.1" {
		t.Errorf("First Next() = %s, want 127.0.0.1", ip1.String())
	}
	
	if ip2 != nil {
		t.Errorf("Second Next() = %v, want nil", ip2)
	}
	
	if ip3 != nil {
		t.Errorf("Third Next() = %v, want nil", ip3)
	}
}

// TestNext tests IPv6 single host
func TestNext_IPv6SingleHost(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("fe80::1/128")
	iter := NewIterator(ipnet)
	
	ip := iter.Next()
	if ip == nil {
		t.Errorf("Next() = nil, want IPv6 address")
	}
	if ip.String() != "fe80::1" {
		t.Errorf("Next() = %s, want fe80::1", ip.String())
	}
	
	ip2 := iter.Next()
	if ip2 != nil {
		t.Errorf("Next() = %v, want nil after single IPv6", ip2)
	}
}

// TestNext tests boundary values in IPv4
func TestNext_BoundaryIPv4Values(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("0.0.0.0/31")
	iter := NewIterator(ipnet)
	
	expected := []string{"0.0.0.0", "0.0.0.1"}
	for _, exp := range expected {
		ip := iter.Next()
		if ip == nil {
			t.Errorf("Next() = nil, want %s", exp)
			break
		}
		if ip.String() != exp {
			t.Errorf("Next() = %s, want %s", ip.String(), exp)
		}
	}
}

// TestNext tests max IPv4 boundary
func TestNext_MaxIPv4Boundary(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("255.255.255.254/31")
	iter := NewIterator(ipnet)
	
	expected := []string{"255.255.255.254", "255.255.255.255"}
	for _, exp := range expected {
		ip := iter.Next()
		if ip == nil {
			t.Errorf("Next() = nil, want %s", exp)
			break
		}
		if ip.String() != exp {
			t.Errorf("Next() = %s, want %s", ip.String(), exp)
		}
	}
	
	if ip := iter.Next(); ip != nil {
		t.Errorf("Next() = %v, want nil after max boundary", ip)
	}
}

// TestNext tests IPv4 with /25 range
func TestNext_IPv4_25Range(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("192.168.0.0/25")
	iter := NewIterator(ipnet)
	
	count := 0
	for {
		ip := iter.Next()
		if ip == nil {
			break
		}
		count++
	}
	
	if count != 128 {
		t.Errorf("Iterated %d IPs in /25, want 128", count)
	}
}

// TestNext tests state preservation between calls
func TestNext_StatePreservation(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("10.0.0.0/29")
	iter := NewIterator(ipnet)
	
	firstIP := iter.Next()
	secondIP := iter.Next()
	thirdIP := iter.Next()
	
	if firstIP.String() != "10.0.0.0" {
		t.Errorf("First IP = %s, want 10.0.0.0", firstIP.String())
	}
	if secondIP.String() != "10.0.0.1" {
		t.Errorf("Second IP = %s, want 10.0.0.1", secondIP.String())
	}
	if thirdIP.String() != "10.0.0.2" {
		t.Errorf("Third IP = %s, want 10.0.0.2", thirdIP.String())
	}
}

// TestNext tests iterator independence
func TestNext_MultipleIteratorIndependence(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("192.168.1.0/30")
	
	iter1 := NewIterator(ipnet)
	iter2 := NewIterator(ipnet)
	
	ip1_1 := iter1.Next()
	ip2_1 := iter2.Next()
	ip1_2 := iter1.Next()
	
	if ip1_1.String() != "192.168.1.0" {
		t.Errorf("iter1 first = %s, want 192.168.1.0", ip1_1.String())
	}
	if ip2_1.String() != "192.168.1.0" {
		t.Errorf("iter2 first = %s, want 192.168.1.0", ip2_1.String())
	}
	if ip1_2.String() != "192.168.1.1" {
		t.Errorf("iter1 second = %s, want 192.168.1.1", ip1_2.String())
	}
}

// TestNext tests IPv6 /120 range
func TestNext_IPv6_120Range(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("2001:db8::/120")
	iter := NewIterator(ipnet)
	
	count := 0
	for {
		ip := iter.Next()
		if ip == nil {
			break
		}
		count++
	}
	
	if count != 256 {
		t.Errorf("Iterated %d IPs in IPv6 /120, want 256", count)
	}
}

// TestNext tests first and last IP in range
func TestNext_FirstAndLastIP(t *testing.T) {
	_, ipnet, _ := net.ParseCIDR("10.0.0.0/30")
	iter := NewIterator(ipnet)
	
	var ips []net.IP
	for {
		ip := iter.Next()
		if ip == nil {
			break
		}
		ips = append(ips, ip)
	}
	
	if len(ips) != 4 {
		t.Errorf("Got %d IPs, want 4", len(ips))
	}
	
	if len(ips) > 0 && ips[0].String() != "10.0.0.0" {
		t.Errorf("First IP = %s, want 10.0.0.0", ips[0].String())
	}
	
	if len(ips) > 3 && ips[3].String() != "10.0.0.3" {
		t.Errorf("Last IP = %s, want 10.0.0.3", ips[3].String())
	}
}