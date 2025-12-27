package netdataapi

import (
	"testing"
	"time"
)

func TestNewOptions(t *testing.T) {
	opts := NewOptions()
	if opts == nil {
		t.Errorf("NewOptions() returned nil")
	}
}

func TestOptionsWithTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		valid   bool
	}{
		{
			name:    "positive timeout",
			timeout: 5 * time.Second,
			valid:   true,
		},
		{
			name:    "zero timeout",
			timeout: 0,
			valid:   false,
		},
		{
			name:    "negative timeout",
			timeout: -1 * time.Second,
			valid:   false,
		},
		{
			name:    "large timeout",
			timeout: 1 * time.Hour,
			valid:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewOptions()
			err := opts.WithTimeout(tt.timeout)
			if (err != nil) != !tt.valid {
				t.Errorf("WithTimeout() error = %v, valid = %v", err, tt.valid)
			}
			if tt.valid && opts.Timeout != tt.timeout {
				t.Errorf("WithTimeout() timeout = %v, want %v", opts.Timeout, tt.timeout)
			}
		})
	}
}

func TestOptionsWithRetries(t *testing.T) {
	tests := []struct {
		name    string
		retries int
		valid   bool
	}{
		{
			name:    "zero retries",
			retries: 0,
			valid:   true,
		},
		{
			name:    "positive retries",
			retries: 3,
			valid:   true,
		},
		{
			name:    "negative retries",
			retries: -1,
			valid:   false,
		},
		{
			name:    "large retries",
			retries: 100,
			valid:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewOptions()
			err := opts.WithRetries(tt.retries)
			if (err != nil) != !tt.valid {
				t.Errorf("WithRetries() error = %v, valid = %v", err, tt.valid)
			}
			if tt.valid && opts.Retries != tt.retries {
				t.Errorf("WithRetries() retries = %v, want %v", opts.Retries, tt.retries)
			}
		})
	}
}

func TestOptionsWithMaxConnections(t *testing.T) {
	tests := []struct {
		name       string
		maxConns   int
		valid      bool
	}{
		{
			name:       "positive connections",
			maxConns:   10,
			valid:      true,
		},
		{
			name:       "zero connections",
			maxConns:   0,
			valid:      false,
		},
		{
			name:       "negative connections",
			maxConns:   -5,
			valid:      false,
		},
		{
			name:       "large connections",
			maxConns:   10000,
			valid:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewOptions()
			err := opts.WithMaxConnections(tt.maxConns)
			if (err != nil) != !tt.valid {
				t.Errorf("WithMaxConnections() error = %v, valid = %v", err, tt.valid)
			}
			if tt.valid && opts.MaxConnections != tt.maxConns {
				t.Errorf("WithMaxConnections() maxConns = %v, want %v", opts.MaxConnections, tt.maxConns)
			}
		})
	}
}

func TestOptionsWithHeaderValue(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
		valid bool
	}{
		{
			name:  "normal header",
			key:   "Authorization",
			value: "Bearer token",
			valid: true,
		},
		{
			name:  "empty key",
			key:   "",
			value: "value",
			valid: false,
		},
		{
			name:  "empty value",
			key:   "X-Custom",
			value: "",
			valid: false,
		},
		{
			name:  "both empty",
			key:   "",
			value: "",
			valid: false,
		},
		{
			name:  "multiple headers",
			key:   "X-Request-ID",
			value: "12345",
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewOptions()
			err := opts.WithHeaderValue(tt.key, tt.value)
			if (err != nil) != !tt.valid {
				t.Errorf("WithHeaderValue() error = %v, valid = %v", err, tt.valid)
			}
			if tt.valid {
				if headerVal, exists := opts.Headers[tt.key]; !exists || headerVal != tt.value {
					t.Errorf("WithHeaderValue() header not set correctly: key=%s, value=%s", tt.key, tt.value)
				}
			}
		})
	}
}

func TestOptionsWithMultipleHeaders(t *testing.T) {
	opts := NewOptions()
	
	err := opts.WithHeaderValue("Authorization", "Bearer token1")
	if err != nil {
		t.Fatalf("WithHeaderValue() first call error = %v", err)
	}
	
	err = opts.WithHeaderValue("X-Custom", "value1")
	if err != nil {
		t.Fatalf("WithHeaderValue() second call error = %v", err)
	}

	if len(opts.Headers) != 2 {
		t.Errorf("expected 2 headers, got %d", len(opts.Headers))
	}

	if opts.Headers["Authorization"] != "Bearer token1" {
		t.Errorf("Authorization header mismatch")
	}
	if opts.Headers["X-Custom"] != "value1" {
		t.Errorf("X-Custom header mismatch")
	}
}

func TestOptionsWithInsecureSkipVerify(t *testing.T) {
	tests := []struct {
		name       string
		skipVerify bool
	}{
		{
			name:       "skip verify true",
			skipVerify: true,
		},
		{
			name:       "skip verify false",
			skipVerify: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewOptions()
			opts.WithInsecureSkipVerify(tt.skipVerify)
			if opts.InsecureSkipVerify != tt.skipVerify {
				t.Errorf("WithInsecureSkipVerify() skipVerify = %v, want %v", opts.InsecureSkipVerify, tt.skipVerify)
			}
		})
	}
}

func TestOptionsChaining(t *testing.T) {
	opts := NewOptions()
	
	err := opts.
		WithTimeout(10*time.Second).
		WithRetries(3).
		WithMaxConnections(20)
	
	if err != nil {
		t.Fatalf("chained options error = %v", err)
	}

	if opts.Timeout != 10*time.Second {
		t.Errorf("timeout not set correctly after chaining")
	}
	if opts.Retries != 3 {
		t.Errorf("retries not set correctly after chaining")
	}
	if opts.MaxConnections != 20 {
		t.Errorf("maxConnections not set correctly after chaining")
	}
}

func TestOptionsWithTLSClientCert(t *testing.T) {
	tests := []struct {
		name     string
		certPath string
		keyPath  string
		valid    bool
	}{
		{
			name:     "valid cert paths",
			certPath: "/path/to/cert.crt",
			keyPath:  "/path/to/key.key",
			valid:    true,
		},
		{
			name:     "empty cert path",
			certPath: "",
			keyPath:  "/path/to/key.key",
			valid:    false,
		},
		{
			name:     "empty key path",
			certPath: "/path/to/cert.crt",
			keyPath:  "",
			valid:    false,
		},
		{
			name:     "both empty",
			certPath: "",
			keyPath:  "",
			valid:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewOptions()
			err := opts.WithTLSClientCert(tt.certPath, tt.keyPath)
			if (err != nil) != !tt.valid {
				t.Errorf("WithTLSClientCert() error = %v, valid = %v", err, tt.valid)
			}
			if tt.valid {
				if opts.TLSCertPath != tt.certPath || opts.TLSKeyPath != tt.keyPath {
					t.Errorf("WithTLSClientCert() paths not set correctly")
				}
			}
		})
	}
}

func TestOptionsGetTimeout(t *testing.T) {
	tests := []struct {
		name           string
		setTimeoutFunc func(*Options)
		expected       time.Duration
	}{
		{
			name: "default timeout",
			setTimeoutFunc: func(opts *Options) {
				// Use default
			},
			expected: 30 * time.Second,
		},
		{
			name: "custom timeout",
			setTimeoutFunc: func(opts *Options) {
				opts.WithTimeout(15 * time.Second)
			},
			expected: 15 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewOptions()
			tt.setTimeoutFunc(opts)
			timeout := opts.GetTimeout()
			if timeout != tt.expected {
				t.Errorf("GetTimeout() = %v, want %v", timeout, tt.expected)
			}
		})
	}
}

func TestOptionsGetRetries(t *testing.T) {
	tests := []struct {
		name          string
		setRetriesFunc func(*Options)
		expected      int
	}{
		{
			name: "default retries",
			setRetriesFunc: func(opts *Options) {
				// Use default
			},
			expected: 3,
		},
		{
			name: "custom retries",
			setRetriesFunc: func(opts *Options) {
				opts.WithRetries(5)
			},
			expected: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewOptions()
			tt.setRetriesFunc(opts)
			retries := opts.GetRetries()
			if retries != tt.expected {
				t.Errorf("GetRetries() = %v, want %v", retries, tt.expected)
			}
		})
	}
}

func TestOptionsIsValid(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(*Options)
		valid     bool
	}{
		{
			name: "valid options",
			setupFunc: func(opts *Options) {
				opts.WithTimeout(5*time.Second)
				opts.WithRetries(3)
				opts.WithMaxConnections(10)
			},
			valid: true,
		},
		{
			name: "invalid timeout",
			setupFunc: func(opts *Options) {
				opts.Timeout = -1 * time.Second
			},
			valid: false,
		},
		{
			name: "invalid retries",
			setupFunc: func(opts *Options) {
				opts.Retries = -1
			},
			valid: false,
		},
		{
			name: "invalid max connections",
			setupFunc: func(opts *Options) {
				opts.MaxConnections = 0
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewOptions()
			tt.setupFunc(opts)
			isValid := opts.IsValid()
			if isValid != tt.valid {
				t.Errorf("IsValid() = %v, want %v", isValid, tt.valid)
			}
		})
	}
}

func TestOptionsWithProxy(t *testing.T) {
	tests := []struct {
		name     string
		proxyURL string
		valid    bool
	}{
		{
			name:     "valid proxy URL",
			proxyURL: "http://proxy.example.com:8080",
			valid:    true,
		},
		{
			name:     "empty proxy URL",
			proxyURL: "",
			valid:    false,
		},
		{
			name:     "https proxy",
			proxyURL: "https://secure-proxy.example.com:8443",
			valid:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewOptions()
			err := opts.WithProxy(tt.proxyURL)
			if (err != nil) != !tt.valid {
				t.Errorf("WithProxy() error = %v, valid = %v", err, tt.valid)
			}
			if tt.valid && opts.ProxyURL != tt.proxyURL {
				t.Errorf("WithProxy() proxyURL = %s, want %s", opts.ProxyURL, tt.proxyURL)
			}
		})
	}
}

func TestOptionsWithBasicAuth(t *testing.T) {
	tests := []struct {
		name     string
		username string
		password string
		valid    bool
	}{
		{
			name:     "valid credentials",
			username: "user",
			password: "pass",
			valid:    true,
		},
		{
			name:     "empty username",
			username: "",
			password: "pass",
			valid:    false,
		},
		{
			name:     "empty password",
			username: "user",
			password: "",
			valid:    false,
		},
		{
			name:     "both empty",
			username: "",
			password: "",
			valid:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := NewOptions()
			err := opts.WithBasicAuth(tt.username, tt.password)
			if (err != nil) != !tt.valid {
				t.Errorf("WithBasicAuth() error = %v, valid = %v", err, tt.valid)
			}
			if tt.valid {
				if opts.BasicAuthUsername != tt.username || opts.BasicAuthPassword != tt.password {
					t.Errorf("WithBasicAuth() credentials not set correctly")
				}
			}
		})
	}
}