package phpfpm

import (
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/pkg/web"
)

func TestInit(t *testing.T) {
	tests := []struct {
		name    string
		phpfpm  *PHPFpm
		wantErr bool
		setup   func(*PHPFpm)
	}{
		{
			name:    "should initialize successfully with default config",
			phpfpm:  New(),
			wantErr: false,
			setup: func(p *PHPFpm) {
				p.URL = "http://localhost:9000/status"
			},
		},
		{
			name:    "should initialize with empty URL and fail",
			phpfpm:  New(),
			wantErr: true,
			setup: func(p *PHPFpm) {
				p.URL = ""
			},
		},
		{
			name:    "should initialize with invalid URL and fail",
			phpfpm:  New(),
			wantErr: true,
			setup: func(p *PHPFpm) {
				p.URL = "not a valid url"
			},
		},
		{
			name:    "should initialize with unix socket URL",
			phpfpm:  New(),
			wantErr: false,
			setup: func(p *PHPFpm) {
				p.URL = "http+unix:///run/php-fpm.sock/status"
			},
		},
		{
			name:    "should handle socket file URL",
			phpfpm:  New(),
			wantErr: false,
			setup: func(p *PHPFpm) {
				p.URL = "unix:///run/php-fpm.sock"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt.phpfpm)
			}
			err := tt.phpfpm.Init()
			if (err != nil) != tt.wantErr {
				t.Errorf("PHPFpm.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitWithHTTPClient(t *testing.T) {
	tests := []struct {
		name       string
		httpClient *web.Client
		wantErr    bool
	}{
		{
			name:       "should handle nil HTTP client",
			httpClient: nil,
			wantErr:    true,
		},
		{
			name: "should handle valid HTTP client",
			httpClient: &web.Client{
				Client: web.NewHTTPClient(),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.URL = "http://localhost:9000/status"
			p.HTTPClient = tt.httpClient

			err := p.Init()
			if (err != nil) != tt.wantErr {
				t.Errorf("PHPFpm.Init() with HTTPClient error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitWithCustomTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout int
		wantErr bool
	}{
		{
			name:    "should handle positive timeout",
			timeout: 5,
			wantErr: false,
		},
		{
			name:    "should handle zero timeout",
			timeout: 0,
			wantErr: true,
		},
		{
			name:    "should handle negative timeout",
			timeout: -1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.URL = "http://localhost:9000/status"
			p.Timeout = tt.timeout

			err := p.Init()
			if (err != nil) != tt.wantErr {
				t.Errorf("PHPFpm.Init() with timeout %d error = %v, wantErr %v", tt.timeout, err, tt.wantErr)
			}
		})
	}
}

func TestInitMultiplePools(t *testing.T) {
	tests := []struct {
		name    string
		pools   []string
		wantErr bool
	}{
		{
			name:    "should handle single pool",
			pools:   []string{"www"},
			wantErr: false,
		},
		{
			name:    "should handle multiple pools",
			pools:   []string{"www", "api", "worker"},
			wantErr: false,
		},
		{
			name:    "should handle empty pools list",
			pools:   []string{},
			wantErr: true,
		},
		{
			name:    "should handle nil pools",
			pools:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.URL = "http://localhost:9000/status"
			p.Pools = tt.pools

			err := p.Init()
			if (err != nil) != tt.wantErr {
				t.Errorf("PHPFpm.Init() with pools %v error = %v, wantErr %v", tt.pools, err, tt.wantErr)
			}
		})
	}
}

func TestInitAddressValidation(t *testing.T) {
	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		{
			name:    "should validate localhost address",
			address: "localhost:9000",
			wantErr: false,
		},
		{
			name:    "should validate IP address",
			address: "127.0.0.1:9000",
			wantErr: false,
		},
		{
			name:    "should validate default port 9000",
			address: "127.0.0.1:9000",
			wantErr: false,
		},
		{
			name:    "should handle missing port",
			address: "127.0.0.1:",
			wantErr: true,
		},
		{
			name:    "should handle invalid port",
			address: "127.0.0.1:99999",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.URL = "http://" + tt.address + "/status"

			err := p.Init()
			if (err != nil) != tt.wantErr {
				t.Errorf("PHPFpm.Init() with address %s error = %v, wantErr %v", tt.address, err, tt.wantErr)
			}
		})
	}
}

func TestInitTLSConfiguration(t *testing.T) {
	tests := []struct {
		name       string
		tlsEnabled bool
		tlsVerify  bool
		certFile   string
		keyFile    string
		caFile     string
		wantErr    bool
	}{
		{
			name:       "should initialize with TLS enabled",
			tlsEnabled: true,
			tlsVerify:  true,
			wantErr:    false,
		},
		{
			name:       "should initialize with TLS disabled",
			tlsEnabled: false,
			tlsVerify:  false,
			wantErr:    false,
		},
		{
			name:       "should handle missing cert file with TLS",
			tlsEnabled: true,
			certFile:   "",
			wantErr:    true,
		},
		{
			name:       "should initialize without TLS verification",
			tlsEnabled: true,
			tlsVerify:  false,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.URL = "https://localhost:9443/status"
			p.HTTPClient = &web.Client{
				TLSConfig: &web.TLSConfig{
					InsecureSkipVerify: !tt.tlsVerify,
					CertFile:           tt.certFile,
					KeyFile:            tt.keyFile,
					CAFile:             tt.caFile,
				},
			}

			err := p.Init()
			if (err != nil) != tt.wantErr {
				t.Errorf("PHPFpm.Init() with TLS config error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInitConcurrentRequests(t *testing.T) {
	tests := []struct {
		name              string
		maxConcurrentReqs int
		wantErr           bool
	}{
		{
			name:              "should handle positive concurrent requests",
			maxConcurrentReqs: 10,
			wantErr:           false,
		},
		{
			name:              "should handle zero concurrent requests",
			maxConcurrentReqs: 0,
			wantErr:           true,
		},
		{
			name:              "should handle negative concurrent requests",
			maxConcurrentReqs: -1,
			wantErr:           true,
		},
		{
			name:              "should handle large concurrent requests",
			maxConcurrentReqs: 1000,
			wantErr:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.URL = "http://localhost:9000/status"
			p.MaxConcurrentRequests = tt.maxConcurrentReqs

			err := p.Init()
			if (err != nil) != tt.wantErr {
				t.Errorf("PHPFpm.Init() with maxConcurrentReqs %d error = %v, wantErr %v", tt.maxConcurrentReqs, err, tt.wantErr)
			}
		})
	}
}

func TestInitPathValidation(t *testing.T) {
	tests := []struct {
		name    string
		urlPath string
		wantErr bool
	}{
		{
			name:    "should handle standard status path",
			urlPath: "/status",
			wantErr: false,
		},
		{
			name:    "should handle json path",
			urlPath: "/status?json",
			wantErr: false,
		},
		{
			name:    "should handle full url status",
			urlPath: "/status?full",
			wantErr: false,
		},
		{
			name:    "should handle empty path",
			urlPath: "",
			wantErr: true,
		},
		{
			name:    "should handle missing slash prefix",
			urlPath: "status",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			p.URL = "http://localhost:9000" + tt.urlPath

			err := p.Init()
			if (err != nil) != tt.wantErr {
				t.Errorf("PHPFpm.Init() with path %s error = %v, wantErr %v", tt.urlPath, err, tt.wantErr)
			}
		})
	}
}

func TestInitEdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*PHPFpm)
		wantErr bool
	}{
		{
			name: "should handle whitespace in URL",
			setup: func(p *PHPFpm) {
				p.URL = "  http://localhost:9000/status  "
			},
			wantErr: true,
		},
		{
			name: "should handle special characters in URL",
			setup: func(p *PHPFpm) {
				p.URL = "http://localhost:9000/status?param=value&other=123"
			},
			wantErr: false,
		},
		{
			name: "should handle port 80",
			setup: func(p *PHPFpm) {
				p.URL = "http://localhost:80/status"
			},
			wantErr: false,
		},
		{
			name: "should handle port 443 with https",
			setup: func(p *PHPFpm) {
				p.URL = "https://localhost:443/status"
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			tt.setup(p)

			err := p.Init()
			if (err != nil) != tt.wantErr {
				t.Errorf("PHPFpm.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}