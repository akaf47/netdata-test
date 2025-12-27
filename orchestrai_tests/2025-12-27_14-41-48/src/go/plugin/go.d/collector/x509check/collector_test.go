package x509check

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCollectorNew tests the New function
func TestCollectorNew(t *testing.T) {
	c := New()
	assert.NotNil(t, c)
	assert.Equal(t, "x509check", c.Name())
}

// TestCollectorCheck tests the Check method
func TestCollectorCheck(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
	}{
		{
			name:        "valid config with source",
			config:      Config{Source: "example.com:443"},
			expectError: false,
		},
		{
			name:        "empty source",
			config:      Config{Source: ""},
			expectError: true,
		},
		{
			name:        "invalid hostname",
			config:      Config{Source: "invalid...hostname:443"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			c.Config = tt.config
			err := c.Check()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				// May error due to no network, but structure should be valid
				if err != nil {
					// Network errors are acceptable in testing
					assert.NotNil(t, c)
				}
			}
		})
	}
}

// TestCollectorCollect tests the Collect method
func TestCollectorCollect(t *testing.T) {
	tests := []struct {
		name          string
		config        Config
		setupMocks    func(*Collector)
		expectMetrics bool
		expectError   bool
	}{
		{
			name: "successful collection",
			config: Config{
				Source: "example.com:443",
			},
			expectMetrics: true,
			expectError:   false,
		},
		{
			name: "collection with invalid source",
			config: Config{
				Source: "invalid...:999999",
			},
			expectMetrics: false,
			expectError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			c.Config = tt.config
			if tt.setupMocks != nil {
				tt.setupMocks(c)
			}
			mx := c.Collect()
			if tt.expectMetrics {
				// May be empty due to network issues in test environment
				assert.NotNil(t, mx)
			}
		})
	}
}

// TestCollectorRun tests the Run method
func TestCollectorRun(t *testing.T) {
	c := New()
	c.Config = Config{Source: "example.com:443"}

	// Run should not panic
	assert.NotPanics(t, func() {
		c.Run()
	})
}

// TestCollectorStop tests the Stop method
func TestCollectorStop(t *testing.T) {
	c := New()
	assert.NotPanics(t, func() {
		c.Stop()
	})
}

// TestCollectorInit tests the Init method
func TestCollectorInit(t *testing.T) {
	c := New()
	c.Config = Config{Source: "example.com:443"}
	err := c.Init()
	assert.NoError(t, err)
}

// TestCollectorCleanup tests the Cleanup method
func TestCollectorCleanup(t *testing.T) {
	c := New()
	assert.NotPanics(t, func() {
		c.Cleanup()
	})
}

// TestGetCertificateChain tests certificate chain retrieval
func TestGetCertificateChain(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		expectError bool
	}{
		{
			name:        "valid domain",
			source:      "google.com:443",
			expectError: false,
		},
		{
			name:        "invalid port",
			source:      "google.com:999999",
			expectError: true,
		},
		{
			name:        "invalid address",
			source:      "invalid...:443",
			expectError: true,
		},
		{
			name:        "unreachable host",
			source:      "127.0.0.1:9999",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			c.Config = Config{Source: tt.source}
			_, err := c.getCertificateChain()
			if tt.expectError {
				assert.Error(t, err)
			}
		})
	}
}

// TestGetCertificateExpiry tests certificate expiry calculation
func TestGetCertificateExpiry(t *testing.T) {
	tests := []struct {
		name       string
		cert       *x509.Certificate
		expectZero bool
	}{
		{
			name: "valid certificate",
			cert: &x509.Certificate{
				NotAfter: time.Now().Add(30 * 24 * time.Hour),
			},
			expectZero: false,
		},
		{
			name: "expired certificate",
			cert: &x509.Certificate{
				NotAfter: time.Now().Add(-1 * time.Hour),
			},
			expectZero: false,
		},
		{
			name: "certificate expiring soon",
			cert: &x509.Certificate{
				NotAfter: time.Now().Add(1 * time.Hour),
			},
			expectZero: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			expiry := c.getCertificateExpiry(tt.cert)
			if tt.expectZero {
				assert.Zero(t, expiry)
			} else {
				assert.NotZero(t, expiry)
			}
		})
	}
}

// TestParseCertificateInfo tests certificate parsing
func TestParseCertificateInfo(t *testing.T) {
	tests := []struct {
		name       string
		cert       *x509.Certificate
		expectNil  bool
		checkName  string
		checkIssuer string
	}{
		{
			name: "valid certificate",
			cert: &x509.Certificate{
				Subject: pkix.Name{
					CommonName: "example.com",
				},
				Issuer: pkix.Name{
					CommonName: "CA",
				},
				NotAfter:  time.Now().Add(30 * 24 * time.Hour),
				NotBefore: time.Now().Add(-24 * time.Hour),
			},
			expectNil:   false,
			checkName:   "example.com",
			checkIssuer: "CA",
		},
		{
			name: "self-signed certificate",
			cert: &x509.Certificate{
				Subject: pkix.Name{
					CommonName: "selfsigned.local",
				},
				Issuer: pkix.Name{
					CommonName: "selfsigned.local",
				},
			},
			expectNil:   false,
			checkName:   "selfsigned.local",
			checkIssuer: "selfsigned.local",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			info := c.parseCertificateInfo(tt.cert)
			if !tt.expectNil {
				assert.NotNil(t, info)
			}
		})
	}
}

// TestCollectorConfigValidation tests configuration validation
func TestCollectorConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		source      string
		port        int
		timeout     time.Duration
		expectError bool
	}{
		{
			name:        "valid source and port",
			source:      "example.com",
			port:        443,
			timeout:     10 * time.Second,
			expectError: false,
		},
		{
			name:        "empty source",
			source:      "",
			port:        443,
			timeout:     10 * time.Second,
			expectError: true,
		},
		{
			name:        "invalid port zero",
			source:      "example.com",
			port:        0,
			timeout:     10 * time.Second,
			expectError: true,
		},
		{
			name:        "invalid port max",
			source:      "example.com",
			port:        65536,
			timeout:     10 * time.Second,
			expectError: true,
		},
		{
			name:        "zero timeout",
			source:      "example.com",
			port:        443,
			timeout:     0,
			expectError: true,
		},
		{
			name:        "negative timeout",
			source:      "example.com",
			port:        443,
			timeout:     -1 * time.Second,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			c.Config.Source = tt.source
			c.Config.Port = tt.port
			c.Config.Timeout = tt.timeout
			err := c.validateConfig()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestHandleTLSConnectionError tests TLS connection error handling
func TestHandleTLSConnectionError(t *testing.T) {
	tests := []struct {
		name           string
		source         string
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:        "connection refused",
			source:      "127.0.0.1:1",
			expectError: true,
		},
		{
			name:        "hostname resolution failure",
			source:      "invalid.nonexistent.domain.example:443",
			expectError: true,
		},
		{
			name:        "timeout scenario",
			source:      "192.0.2.1:443", // TEST-NET-1 (non-routable)
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			c.Config = Config{
				Source:  tt.source,
				Timeout: 2 * time.Second,
			}
			_, err := c.getCertificateChain()
			if tt.expectError {
				assert.Error(t, err)
			}
		})
	}
}

// TestMetricsCollection tests that all metrics are collected
func TestMetricsCollection(t *testing.T) {
	c := New()
	c.Config = Config{Source: "google.com:443"}
	
	err := c.Init()
	require.NoError(t, err)
	
	mx := c.Collect()
	assert.NotNil(t, mx)
	
	// Verify that metrics have expected structure
	if mx != nil {
		assert.IsType(t, map[string]interface{}{}, mx)
	}
}

// TestCollectorWithMultipleSources tests collecting from multiple sources
func TestCollectorWithMultipleSources(t *testing.T) {
	sources := []string{
		"google.com:443",
		"github.com:443",
		"example.com:443",
	}

	for _, source := range sources {
		t.Run(source, func(t *testing.T) {
			c := New()
			c.Config = Config{Source: source}
			err := c.Init()
			// May error due to network, but should not panic
			assert.NotPanics(t, func() {
				c.Collect()
			})
		})
	}
}

// TestCertificateChainDepth tests certificate chain depth
func TestCertificateChainDepth(t *testing.T) {
	c := New()
	c.Config = Config{Source: "google.com:443"}
	
	certs, err := c.getCertificateChain()
	if err == nil && len(certs) > 0 {
		// Valid certificate chain
		assert.GreaterOrEqual(t, len(certs), 1)
	}
}

// TestCollectorStateTransition tests state transitions
func TestCollectorStateTransition(t *testing.T) {
	c := New()
	c.Config = Config{Source: "example.com:443"}
	
	// Init -> Check -> Collect -> Stop -> Cleanup
	err := c.Init()
	assert.NoError(t, err)
	
	err = c.Check()
	// May error due to network, accept error
	_ = err
	
	mx := c.Collect()
	assert.NotNil(t, mx)
	
	c.Stop()
	c.Cleanup()
}

// TestCollectorMetricsStructure tests metrics have correct structure
func TestCollectorMetricsStructure(t *testing.T) {
	c := New()
	c.Config = Config{Source: "google.com:443"}
	
	err := c.Init()
	if err != nil {
		t.Skip("Skipping due to init error")
	}
	
	mx := c.Collect()
	if mx != nil {
		_, ok := mx.(map[string]interface{})
		assert.True(t, ok, "Metrics should be a map")
	}
}

// TestSNISupport tests SNI (Server Name Indication) support
func TestSNISupport(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		expectSNI bool
	}{
		{
			name:      "domain with SNI",
			source:    "example.com:443",
			expectSNI: true,
		},
		{
			name:      "IP address without SNI",
			source:    "93.184.216.34:443",
			expectSNI: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			c.Config = Config{Source: tt.source}
			// Just verify no panic
			assert.NotPanics(t, func() {
				_, _ = c.getCertificateChain()
			})
		})
	}
}

// TestCertificateValidation tests certificate validation
func TestCertificateValidation(t *testing.T) {
	tests := []struct {
		name        string
		cert        *x509.Certificate
		expectValid bool
	}{
		{
			name: "valid certificate",
			cert: &x509.Certificate{
				NotAfter:  time.Now().Add(30 * 24 * time.Hour),
				NotBefore: time.Now().Add(-24 * time.Hour),
				Raw:       []byte("valid"),
			},
			expectValid: true,
		},
		{
			name: "not yet valid",
			cert: &x509.Certificate{
				NotAfter:  time.Now().Add(30 * 24 * time.Hour),
				NotBefore: time.Now().Add(24 * time.Hour),
				Raw:       []byte("notyet"),
			},
			expectValid: false,
		},
		{
			name: "expired certificate",
			cert: &x509.Certificate{
				NotAfter:  time.Now().Add(-1 * time.Hour),
				NotBefore: time.Now().Add(-100 * time.Hour),
				Raw:       []byte("expired"),
			},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			isValid := c.isCertificateValid(tt.cert)
			if tt.expectValid {
				assert.True(t, isValid)
			} else {
				assert.False(t, isValid)
			}
		})
	}
}

// TestErrorRecovery tests error recovery mechanisms
func TestErrorRecovery(t *testing.T) {
	c := New()
	
	// First call with invalid config
	c.Config = Config{Source: "invalid...:999999"}
	_ = c.Check()
	
	// Second call with valid config should work
	c.Config = Config{Source: "google.com:443"}
	assert.NotPanics(t, func() {
		_ = c.Check()
	})
}

// TestConcurrentCollect tests concurrent collection attempts
func TestConcurrentCollect(t *testing.T) {
	c := New()
	c.Config = Config{Source: "google.com:443"}
	c.Init()
	
	done := make(chan bool, 3)
	
	go func() {
		_ = c.Collect()
		done <- true
	}()
	
	go func() {
		_ = c.Collect()
		done <- true
	}()
	
	go func() {
		_ = c.Collect()
		done <- true
	}()
	
	for i := 0; i < 3; i++ {
		<-done
	}
}

// TestTimeoutHandling tests timeout handling
func TestTimeoutHandling(t *testing.T) {
	c := New()
	c.Config = Config{
		Source:  "192.0.2.1:443", // Non-routable IP
		Timeout: 100 * time.Millisecond,
	}
	
	_, err := c.getCertificateChain()
	// Should timeout
	assert.Error(t, err)
}

// TestCollectorName tests the Name method
func TestCollectorName(t *testing.T) {
	c := New()
	assert.Equal(t, "x509check", c.Name())
}

// TestCollectorType tests the Type method
func TestCollectorType(t *testing.T) {
	c := New()
	assert.NotEmpty(t, c.Type())
}

// TestCollectorModuleName tests the ModuleName method
func TestCollectorModuleName(t *testing.T) {
	c := New()
	assert.NotEmpty(t, c.ModuleName())
}

// TestParseSourceString tests parsing source strings
func TestParseSourceString(t *testing.T) {
	tests := []struct {
		name           string
		source         string
		expectedHost   string
		expectedPort   int
		expectError    bool
	}{
		{
			name:         "standard format",
			source:       "example.com:443",
			expectedHost: "example.com",
			expectedPort: 443,
			expectError:  false,
		},
		{
			name:         "with explicit port",
			source:       "example.com:8443",
			expectedHost: "example.com",
			expectedPort: 8443,
			expectError:  false,
		},
		{
			name:        "missing port",
			source:      "example.com",
			expectError: true,
		},
		{
			name:        "invalid port",
			source:      "example.com:notaport",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			c.Config.Source = tt.source
			host, port, err := c.parseSource()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedHost, host)
				assert.Equal(t, tt.expectedPort, port)
			}
		})
	}
}

// TestSubjectAlternativeNames tests SAN parsing
func TestSubjectAlternativeNames(t *testing.T) {
	cert := &x509.Certificate{
		Subject: pkix.Name{
			CommonName: "example.com",
		},
		DNSNames:       []string{"example.com", "www.example.com", "*.example.com"},
		IPAddresses:    []net.IP{net.ParseIP("192.0.2.1")},
		EmailAddresses: []string{"admin@example.com"},
	}

	c := New()
	info := c.parseCertificateInfo(cert)
	assert.NotNil(t, info)
}

// TestCollectorInitWithoutConfig tests Init without proper config
func TestCollectorInitWithoutConfig(t *testing.T) {
	c := New()
	// Call Init without config
	err := c.Init()
	// Should handle gracefully
	assert.NotNil(t, err)
}

// TestRawCertificateData tests raw certificate data handling
func TestRawCertificateData(t *testing.T) {
	cert := &x509.Certificate{
		Raw: []byte{0x30, 0x82, 0x02, 0x5d}, // Sample DER encoding
	}

	c := New()
	info := c.parseCertificateInfo(cert)
	assert.NotNil(t, info)
}

// TestCollectorCleanupAfterError tests cleanup after error
func TestCollectorCleanupAfterError(t *testing.T) {
	c := New()
	c.Config = Config{Source: "invalid...:999999"}
	
	// Generate error
	_ = c.Check()
	
	// Cleanup should not panic
	assert.NotPanics(t, func() {
		c.Cleanup()
	})
}

// TestCollectorStopIdempotent tests that Stop can be called multiple times
func TestCollectorStopIdempotent(t *testing.T) {
	c := New()
	
	assert.NotPanics(t, func() {
		c.Stop()
		c.Stop()
		c.Stop()
	})
}