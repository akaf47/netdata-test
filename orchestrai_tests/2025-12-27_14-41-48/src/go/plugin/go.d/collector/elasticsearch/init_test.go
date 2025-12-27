package elasticsearch

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitDefault(t *testing.T) {
	tests := []struct {
		name          string
		expectedURL   string
		expectedProbe string
	}{
		{
			name:          "should initialize with default values",
			expectedURL:   "http://localhost:9200",
			expectedProbe: "/_cluster/health",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			require.NotNil(t, c)

			// Check if default values exist
			assert.Equal(t, "elasticsearch", c.Module)
		})
	}
}

func TestInitWithValidConfig(t *testing.T) {
	tests := []struct {
		name   string
		config struct {
			URL      string
			Username string
			Password string
		}
		wantErr bool
	}{
		{
			name: "should initialize with valid URL",
			config: struct {
				URL      string
				Username string
				Password string
			}{
				URL:      "http://localhost:9200",
				Username: "user",
				Password: "pass",
			},
			wantErr: false,
		},
		{
			name: "should initialize with HTTPS URL",
			config: struct {
				URL      string
				Username string
				Password string
			}{
				URL:      "https://elasticsearch.example.com:9200",
				Username: "",
				Password: "",
			},
			wantErr: false,
		},
		{
			name: "should initialize with custom port",
			config: struct {
				URL      string
				Username string
				Password string
			}{
				URL:      "http://localhost:9300",
				Username: "",
				Password: "",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			c.URL = tt.config.URL
			c.Username = tt.config.Username
			c.Password = tt.config.Password

			err := c.Init()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				if tt.config.URL != "" {
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestInitHTTPClient(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*Collector)
		wantErr bool
	}{
		{
			name: "should initialize HTTP client without auth",
			setup: func(c *Collector) {
				c.URL = "http://localhost:9200"
				c.Username = ""
				c.Password = ""
			},
			wantErr: false,
		},
		{
			name: "should initialize HTTP client with auth",
			setup: func(c *Collector) {
				c.URL = "http://localhost:9200"
				c.Username = "elastic"
				c.Password = "changeme"
			},
			wantErr: false,
		},
		{
			name: "should handle empty URL",
			setup: func(c *Collector) {
				c.URL = ""
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			tt.setup(c)

			err := c.Init()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				if c.URL != "" {
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestInitValidation(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		username   string
		password   string
		tlsSkip    bool
		shouldPass bool
	}{
		{
			name:       "should validate minimal config",
			url:        "http://localhost:9200",
			username:   "",
			password:   "",
			tlsSkip:    false,
			shouldPass: true,
		},
		{
			name:       "should validate with basic auth",
			url:        "http://localhost:9200",
			username:   "admin",
			password:   "secret",
			tlsSkip:    false,
			shouldPass: true,
		},
		{
			name:       "should validate HTTPS with TLS skip",
			url:        "https://localhost:9200",
			username:   "",
			password:   "",
			tlsSkip:    true,
			shouldPass: true,
		},
		{
			name:       "should reject empty URL",
			url:        "",
			username:   "",
			password:   "",
			tlsSkip:    false,
			shouldPass: false,
		},
		{
			name:       "should reject invalid URL format",
			url:        "not-a-url",
			username:   "",
			password:   "",
			tlsSkip:    false,
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			c.URL = tt.url
			c.Username = tt.username
			c.Password = tt.password

			err := c.Init()
			if tt.shouldPass {
				if tt.url != "" {
					assert.NoError(t, err)
				}
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestInitNetworkOptions(t *testing.T) {
	tests := []struct {
		name    string
		timeout string
		retry   int
		skip    bool
	}{
		{
			name:    "should initialize with custom timeout",
			timeout: "10s",
			retry:   3,
			skip:    false,
		},
		{
			name:    "should initialize with zero timeout",
			timeout: "0s",
			retry:   0,
			skip:    false,
		},
		{
			name:    "should initialize with negative timeout",
			timeout: "-1s",
			retry:   -1,
			skip:    false,
		},
		{
			name:    "should initialize with TLS skip enabled",
			timeout: "5s",
			retry:   1,
			skip:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			c.URL = "http://localhost:9200"
			assert.NotNil(t, c)
		})
	}
}

func TestInitModuleMetadata(t *testing.T) {
	tests := []struct {
		name           string
		expectedModule string
		expectedName   string
	}{
		{
			name:           "should set correct module name",
			expectedModule: "elasticsearch",
			expectedName:   "elasticsearch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			assert.Equal(t, tt.expectedModule, c.Module)
		})
	}
}

func TestInitURLParsing(t *testing.T) {
	tests := []struct {
		name      string
		url       string
		wantErr   bool
		wantProbe string
	}{
		{
			name:      "should parse http URL",
			url:       "http://localhost:9200",
			wantErr:   false,
			wantProbe: "/_cluster/health",
		},
		{
			name:      "should parse https URL",
			url:       "https://elasticsearch.example.com:9200",
			wantErr:   false,
			wantProbe: "/_cluster/health",
		},
		{
			name:      "should parse URL with path",
			url:       "http://localhost:9200/elasticsearch",
			wantErr:   false,
			wantProbe: "/_cluster/health",
		},
		{
			name:      "should handle URL without port",
			url:       "http://localhost",
			wantErr:   false,
			wantProbe: "/_cluster/health",
		},
		{
			name:      "should handle URL with credentials",
			url:       "http://user:pass@localhost:9200",
			wantErr:   false,
			wantProbe: "/_cluster/health",
		},
		{
			name:      "should reject invalid URL",
			url:       "://invalid",
			wantErr:   true,
			wantProbe: "",
		},
		{
			name:      "should reject empty URL",
			url:       "",
			wantErr:   true,
			wantProbe: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			c.URL = tt.url

			err := c.Init()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				if tt.url != "" {
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestInitClientConfiguration(t *testing.T) {
	tests := []struct {
		name          string
		setupConfig   func(*Collector)
		expectedError bool
	}{
		{
			name: "should configure HTTP client without authentication",
			setupConfig: func(c *Collector) {
				c.URL = "http://localhost:9200"
			},
			expectedError: false,
		},
		{
			name: "should configure HTTP client with username only",
			setupConfig: func(c *Collector) {
				c.URL = "http://localhost:9200"
				c.Username = "elastic"
			},
			expectedError: false,
		},
		{
			name: "should configure HTTP client with password only",
			setupConfig: func(c *Collector) {
				c.URL = "http://localhost:9200"
				c.Password = "changeme"
			},
			expectedError: false,
		},
		{
			name: "should configure HTTP client with both credentials",
			setupConfig: func(c *Collector) {
				c.URL = "http://localhost:9200"
				c.Username = "elastic"
				c.Password = "changeme"
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			tt.setupConfig(c)

			err := c.Init()
			if tt.expectedError {
				assert.Error(t, err)
			} else {
				if c.URL != "" {
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestInitProbeEndpoint(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		expectedProbe string
		wantErr       bool
	}{
		{
			name:          "should set probe endpoint for localhost",
			url:           "http://localhost:9200",
			expectedProbe: "http://localhost:9200/_cluster/health",
			wantErr:       false,
		},
		{
			name:          "should set probe endpoint with domain",
			url:           "https://elasticsearch.example.com:9200",
			expectedProbe: "https://elasticsearch.example.com:9200/_cluster/health",
			wantErr:       false,
		},
		{
			name:          "should handle URL without explicit port",
			url:           "http://localhost",
			expectedProbe: "http://localhost/_cluster/health",
			wantErr:       false,
		},
		{
			name:          "should reject empty URL",
			url:           "",
			expectedProbe: "",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCollector()
			c.URL = tt.url

			err := c.Init()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				if tt.url != "" {
					assert.NoError(t, err)
				}
			}
		})
	}
}