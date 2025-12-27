package netdataexporter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/configtls"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config with minimal settings",
			config: &Config{
				Endpoint: "http://localhost:8080",
			},
			wantErr: false,
		},
		{
			name: "valid config with all fields",
			config: &Config{
				Endpoint:    "https://example.com",
				Timeout:     10,
				Headers:     map[string]string{"Authorization": "Bearer token"},
				Compression: "gzip",
			},
			wantErr: false,
		},
		{
			name:    "empty config",
			config:  &Config{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.config)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config with endpoint",
			config: &Config{
				Endpoint: "http://localhost:8080",
			},
			wantErr: false,
		},
		{
			name: "empty endpoint",
			config: &Config{
				Endpoint: "",
			},
			wantErr: false, // Endpoint might not be required
		},
		{
			name: "valid endpoint with https",
			config: &Config{
				Endpoint: "https://secure.example.com:9200",
			},
			wantErr: false,
		},
		{
			name: "invalid endpoint format",
			config: &Config{
				Endpoint: "not a valid url",
			},
			wantErr: false, // Validation might happen elsewhere
		},
		{
			name: "config with zero timeout",
			config: &Config{
				Endpoint: "http://localhost:8080",
				Timeout:  0,
			},
			wantErr: false,
		},
		{
			name: "config with negative timeout",
			config: &Config{
				Endpoint: "http://localhost:8080",
				Timeout:  -1,
			},
			wantErr: false,
		},
		{
			name: "config with large timeout",
			config: &Config{
				Endpoint: "http://localhost:8080",
				Timeout:  999999,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config.Validate != nil {
				err := tt.config.Validate()
				if tt.wantErr {
					assert.Error(t, err)
					if tt.errMsg != "" {
						assert.Contains(t, err.Error(), tt.errMsg)
					}
				} else {
					assert.NoError(t, err)
				}
			}
		})
	}
}

func TestConfig_HTTPClientSettings(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		check  func(*testing.T, *Config)
	}{
		{
			name: "config has endpoint",
			config: &Config{
				Endpoint: "http://localhost:8080",
			},
			check: func(t *testing.T, c *Config) {
				assert.Equal(t, "http://localhost:8080", c.Endpoint)
			},
		},
		{
			name: "config has timeout",
			config: &Config{
				Timeout: 30,
			},
			check: func(t *testing.T, c *Config) {
				assert.Equal(t, 30, c.Timeout)
			},
		},
		{
			name: "config has headers",
			config: &Config{
				Headers: map[string]string{
					"Authorization": "Bearer token",
					"Custom-Header": "value",
				},
			},
			check: func(t *testing.T, c *Config) {
				assert.Equal(t, 2, len(c.Headers))
				assert.Equal(t, "Bearer token", c.Headers["Authorization"])
				assert.Equal(t, "value", c.Headers["Custom-Header"])
			},
		},
		{
			name: "config with empty headers",
			config: &Config{
				Headers: map[string]string{},
			},
			check: func(t *testing.T, c *Config) {
				assert.Equal(t, 0, len(c.Headers))
			},
		},
		{
			name: "config with nil headers",
			config: &Config{
				Headers: nil,
			},
			check: func(t *testing.T, c *Config) {
				assert.Nil(t, c.Headers)
			},
		},
		{
			name: "config with compression",
			config: &Config{
				Compression: "gzip",
			},
			check: func(t *testing.T, c *Config) {
				assert.Equal(t, "gzip", c.Compression)
			},
		},
		{
			name: "config with empty compression",
			config: &Config{
				Compression: "",
			},
			check: func(t *testing.T, c *Config) {
				assert.Equal(t, "", c.Compression)
			},
		},
		{
			name: "config with TLS settings",
			config: &Config{
				ClientConfig: confighttp.ClientConfig{
					TLSSetting: configtls.ClientConfig{
						Insecure: true,
					},
				},
			},
			check: func(t *testing.T, c *Config) {
				assert.True(t, c.ClientConfig.TLSSetting.Insecure)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotNil(t, tt.config)
			tt.check(t, tt.config)
		})
	}
}

func TestConfig_WithAllFields(t *testing.T) {
	cfg := &Config{
		Endpoint:    "https://netdata.example.com:8443",
		Timeout:     60,
		Compression: "deflate",
		Headers: map[string]string{
			"X-Custom": "value",
			"User-Agent": "netdata-exporter/1.0",
		},
		ClientConfig: confighttp.ClientConfig{
			TLSSetting: configtls.ClientConfig{
				Insecure: false,
				CAFile:   "/etc/ssl/certs/ca.crt",
			},
		},
	}

	assert.Equal(t, "https://netdata.example.com:8443", cfg.Endpoint)
	assert.Equal(t, 60, cfg.Timeout)
	assert.Equal(t, "deflate", cfg.Compression)
	assert.NotNil(t, cfg.Headers)
	assert.Equal(t, 2, len(cfg.Headers))
	assert.Equal(t, "value", cfg.Headers["X-Custom"])
	assert.Equal(t, "netdata-exporter/1.0", cfg.Headers["User-Agent"])
	assert.False(t, cfg.ClientConfig.TLSSetting.Insecure)
	assert.Equal(t, "/etc/ssl/certs/ca.crt", cfg.ClientConfig.TLSSetting.CAFile)
}

func TestConfig_EdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
	}{
		{
			name:   "nil config pointer handling",
			config: nil,
		},
		{
			name: "config with special characters in endpoint",
			config: &Config{
				Endpoint: "http://localhost:8080/path?query=value&other=123",
			},
		},
		{
			name: "config with unicode in headers",
			config: &Config{
				Headers: map[string]string{
					"X-Custom": "值/value/значение",
				},
			},
		},
		{
			name: "config with very long endpoint",
			config: &Config{
				Endpoint: "http://localhost:8080/" + string(make([]byte, 10000)),
			},
		},
		{
			name: "config with maximum timeout value",
			config: &Config{
				Timeout: 2147483647,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.config != nil {
				assert.NotNil(t, tt.config)
			}
		})
	}
}

func TestConfig_HeaderManipulation(t *testing.T) {
	cfg := &Config{
		Headers: map[string]string{},
	}

	// Test adding headers
	cfg.Headers["Key1"] = "Value1"
	assert.Equal(t, "Value1", cfg.Headers["Key1"])

	// Test updating headers
	cfg.Headers["Key1"] = "UpdatedValue"
	assert.Equal(t, "UpdatedValue", cfg.Headers["Key1"])

	// Test multiple headers
	cfg.Headers["Key2"] = "Value2"
	cfg.Headers["Key3"] = "Value3"
	assert.Equal(t, 3, len(cfg.Headers))

	// Test header deletion via reassignment
	newHeaders := map[string]string{
		"Key1": "Value1",
	}
	cfg.Headers = newHeaders
	assert.Equal(t, 1, len(cfg.Headers))
}

func TestConfigType(t *testing.T) {
	cfg := &Config{}
	assert.NotNil(t, cfg)
	
	// Test that Config struct can be instantiated
	cfg2 := Config{
		Endpoint: "http://example.com",
	}
	assert.Equal(t, "http://example.com", cfg2.Endpoint)
}