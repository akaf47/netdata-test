package tlscfg

import (
	"crypto/tls"
	"testing"
)

// TestNewConfig tests the creation of a new TLS configuration
func TestNewConfig(t *testing.T) {
	tests := []struct {
		name    string
		want    *Config
		wantErr bool
	}{
		{
			name:    "should create new config successfully",
			want:    &Config{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == nil && !tt.wantErr {
				t.Errorf("NewConfig() = nil, want non-nil Config")
			}
		})
	}
}

// TestConfigSetCertFile tests setting certificate file
func TestConfigSetCertFile(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		path    string
		wantErr bool
	}{
		{
			name:    "should set certificate file with valid path",
			cfg:     &Config{},
			path:    "/etc/ssl/certs/server.crt",
			wantErr: false,
		},
		{
			name:    "should handle empty path",
			cfg:     &Config{},
			path:    "",
			wantErr: false,
		},
		{
			name:    "should handle nonexistent file",
			cfg:     &Config{},
			path:    "/nonexistent/path/file.crt",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.SetCertFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetCertFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestConfigSetKeyFile tests setting key file
func TestConfigSetKeyFile(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		path    string
		wantErr bool
	}{
		{
			name:    "should set key file with valid path",
			cfg:     &Config{},
			path:    "/etc/ssl/private/server.key",
			wantErr: false,
		},
		{
			name:    "should handle empty path",
			cfg:     &Config{},
			path:    "",
			wantErr: false,
		},
		{
			name:    "should handle nonexistent file",
			cfg:     &Config{},
			path:    "/nonexistent/path/file.key",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.SetKeyFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetKeyFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestConfigSetCAFile tests setting CA file
func TestConfigSetCAFile(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		path    string
		wantErr bool
	}{
		{
			name:    "should set CA file with valid path",
			cfg:     &Config{},
			path:    "/etc/ssl/certs/ca.crt",
			wantErr: false,
		},
		{
			name:    "should handle empty path",
			cfg:     &Config{},
			path:    "",
			wantErr: false,
		},
		{
			name:    "should handle nonexistent file",
			cfg:     &Config{},
			path:    "/nonexistent/path/ca.crt",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.SetCAFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetCAFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestConfigSetMinVersion tests setting minimum TLS version
func TestConfigSetMinVersion(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		version uint16
		wantErr bool
	}{
		{
			name:    "should set TLS 1.0 as min version",
			cfg:     &Config{},
			version: tls.VersionTLS10,
			wantErr: false,
		},
		{
			name:    "should set TLS 1.1 as min version",
			cfg:     &Config{},
			version: tls.VersionTLS11,
			wantErr: false,
		},
		{
			name:    "should set TLS 1.2 as min version",
			cfg:     &Config{},
			version: tls.VersionTLS12,
			wantErr: false,
		},
		{
			name:    "should set TLS 1.3 as min version",
			cfg:     &Config{},
			version: tls.VersionTLS13,
			wantErr: false,
		},
		{
			name:    "should handle invalid version",
			cfg:     &Config{},
			version: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.SetMinVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetMinVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestConfigSetMaxVersion tests setting maximum TLS version
func TestConfigSetMaxVersion(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		version uint16
		wantErr bool
	}{
		{
			name:    "should set TLS 1.2 as max version",
			cfg:     &Config{},
			version: tls.VersionTLS12,
			wantErr: false,
		},
		{
			name:    "should set TLS 1.3 as max version",
			cfg:     &Config{},
			version: tls.VersionTLS13,
			wantErr: false,
		},
		{
			name:    "should handle invalid version",
			cfg:     &Config{},
			version: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.SetMaxVersion(tt.version)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetMaxVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestConfigSetInsecureSkipVerify tests setting InsecureSkipVerify flag
func TestConfigSetInsecureSkipVerify(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		value   bool
		wantErr bool
	}{
		{
			name:    "should set InsecureSkipVerify to true",
			cfg:     &Config{},
			value:   true,
			wantErr: false,
		},
		{
			name:    "should set InsecureSkipVerify to false",
			cfg:     &Config{},
			value:   false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.SetInsecureSkipVerify(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetInsecureSkipVerify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestConfigBuild tests building the TLS configuration
func TestConfigBuild(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		want    *tls.Config
		wantErr bool
	}{
		{
			name:    "should build config successfully with minimal settings",
			cfg:     &Config{},
			want:    &tls.Config{},
			wantErr: false,
		},
		{
			name:    "should build config with InsecureSkipVerify",
			cfg:     &Config{InsecureSkipVerify: true},
			want:    &tls.Config{InsecureSkipVerify: true},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.cfg.Build()
			if (err != nil) != tt.wantErr {
				t.Errorf("Build() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && got == nil {
				t.Errorf("Build() = nil, want *tls.Config")
			}
		})
	}
}

// TestConfigGetCertFile tests getting the certificate file path
func TestConfigGetCertFile(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Config
		want string
	}{
		{
			name: "should return empty string when not set",
			cfg:  &Config{},
			want: "",
		},
		{
			name: "should return set certificate file",
			cfg:  &Config{CertFile: "/etc/ssl/certs/server.crt"},
			want: "/etc/ssl/certs/server.crt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.GetCertFile()
			if got != tt.want {
				t.Errorf("GetCertFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestConfigGetKeyFile tests getting the key file path
func TestConfigGetKeyFile(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Config
		want string
	}{
		{
			name: "should return empty string when not set",
			cfg:  &Config{},
			want: "",
		},
		{
			name: "should return set key file",
			cfg:  &Config{KeyFile: "/etc/ssl/private/server.key"},
			want: "/etc/ssl/private/server.key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.GetKeyFile()
			if got != tt.want {
				t.Errorf("GetKeyFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestConfigGetCAFile tests getting the CA file path
func TestConfigGetCAFile(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Config
		want string
	}{
		{
			name: "should return empty string when not set",
			cfg:  &Config{},
			want: "",
		},
		{
			name: "should return set CA file",
			cfg:  &Config{CAFile: "/etc/ssl/certs/ca.crt"},
			want: "/etc/ssl/certs/ca.crt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.GetCAFile()
			if got != tt.want {
				t.Errorf("GetCAFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestConfigGetMinVersion tests getting the minimum TLS version
func TestConfigGetMinVersion(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Config
		want uint16
	}{
		{
			name: "should return default value when not set",
			cfg:  &Config{},
			want: 0,
		},
		{
			name: "should return set min version",
			cfg:  &Config{MinVersion: tls.VersionTLS12},
			want: tls.VersionTLS12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.GetMinVersion()
			if got != tt.want {
				t.Errorf("GetMinVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestConfigGetMaxVersion tests getting the maximum TLS version
func TestConfigGetMaxVersion(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Config
		want uint16
	}{
		{
			name: "should return default value when not set",
			cfg:  &Config{},
			want: 0,
		},
		{
			name: "should return set max version",
			cfg:  &Config{MaxVersion: tls.VersionTLS13},
			want: tls.VersionTLS13,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.GetMaxVersion()
			if got != tt.want {
				t.Errorf("GetMaxVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestConfigGetInsecureSkipVerify tests getting the InsecureSkipVerify flag
func TestConfigGetInsecureSkipVerify(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Config
		want bool
	}{
		{
			name: "should return false when not set",
			cfg:  &Config{},
			want: false,
		},
		{
			name: "should return true when set",
			cfg:  &Config{InsecureSkipVerify: true},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.GetInsecureSkipVerify()
			if got != tt.want {
				t.Errorf("GetInsecureSkipVerify() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestConfigValidate tests configuration validation
func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
	}{
		{
			name:    "should validate empty config",
			cfg:     &Config{},
			wantErr: false,
		},
		{
			name:    "should validate config with cert and key",
			cfg:     &Config{CertFile: "/etc/ssl/certs/server.crt", KeyFile: "/etc/ssl/private/server.key"},
			wantErr: false,
		},
		{
			name:    "should fail validation with only cert file",
			cfg:     &Config{CertFile: "/etc/ssl/certs/server.crt"},
			wantErr: true,
		},
		{
			name:    "should fail validation with only key file",
			cfg:     &Config{KeyFile: "/etc/ssl/private/server.key"},
			wantErr: true,
		},
		{
			name:    "should fail validation with invalid min version",
			cfg:     &Config{MinVersion: 0},
			wantErr: true,
		},
		{
			name:    "should fail validation when min > max",
			cfg:     &Config{MinVersion: tls.VersionTLS13, MaxVersion: tls.VersionTLS12},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestConfigReset tests resetting the configuration
func TestConfigReset(t *testing.T) {
	t.Run("should reset all config values", func(t *testing.T) {
		cfg := &Config{
			CertFile:             "/etc/ssl/certs/server.crt",
			KeyFile:              "/etc/ssl/private/server.key",
			CAFile:               "/etc/ssl/certs/ca.crt",
			MinVersion:           tls.VersionTLS12,
			MaxVersion:           tls.VersionTLS13,
			InsecureSkipVerify:   true,
		}

		cfg.Reset()

		if cfg.CertFile != "" {
			t.Errorf("Reset() CertFile = %v, want empty", cfg.CertFile)
		}
		if cfg.KeyFile != "" {
			t.Errorf("Reset() KeyFile = %v, want empty", cfg.KeyFile)
		}
		if cfg.CAFile != "" {
			t.Errorf("Reset() CAFile = %v, want empty", cfg.CAFile)
		}
		if cfg.MinVersion != 0 {
			t.Errorf("Reset() MinVersion = %v, want 0", cfg.MinVersion)
		}
		if cfg.MaxVersion != 0 {
			t.Errorf("Reset() MaxVersion = %v, want 0", cfg.MaxVersion)
		}
		if cfg.InsecureSkipVerify != false {
			t.Errorf("Reset() InsecureSkipVerify = %v, want false", cfg.InsecureSkipVerify)
		}
	})
}

// TestConfigClone tests cloning the configuration
func TestConfigClone(t *testing.T) {
	t.Run("should clone config correctly", func(t *testing.T) {
		original := &Config{
			CertFile:           "/etc/ssl/certs/server.crt",
			KeyFile:            "/etc/ssl/private/server.key",
			CAFile:             "/etc/ssl/certs/ca.crt",
			MinVersion:         tls.VersionTLS12,
			MaxVersion:         tls.VersionTLS13,
			InsecureSkipVerify: true,
		}

		cloned := original.Clone()

		if cloned.CertFile != original.CertFile {
			t.Errorf("Clone() CertFile = %v, want %v", cloned.CertFile, original.CertFile)
		}
		if cloned.KeyFile != original.KeyFile {
			t.Errorf("Clone() KeyFile = %v, want %v", cloned.KeyFile, original.KeyFile)
		}
		if cloned.CAFile != original.CAFile {
			t.Errorf("Clone() CAFile = %v, want %v", cloned.CAFile, original.CAFile)
		}
		if cloned.MinVersion != original.MinVersion {
			t.Errorf("Clone() MinVersion = %v, want %v", cloned.MinVersion, original.MinVersion)
		}
		if cloned.MaxVersion != original.MaxVersion {
			t.Errorf("Clone() MaxVersion = %v, want %v", cloned.MaxVersion, original.MaxVersion)
		}
		if cloned.InsecureSkipVerify != original.InsecureSkipVerify {
			t.Errorf("Clone() InsecureSkipVerify = %v, want %v", cloned.InsecureSkipVerify, original.InsecureSkipVerify)
		}
	})

	t.Run("should clone independently from original", func(t *testing.T) {
		original := &Config{
			CertFile: "/etc/ssl/certs/server.crt",
		}

		cloned := original.Clone()
		cloned.CertFile = "/etc/ssl/certs/other.crt"

		if original.CertFile == cloned.CertFile {
			t.Errorf("Clone() is not independent - modifications affect original")
		}
	})
}

// TestConfigString tests the string representation
func TestConfigString(t *testing.T) {
	t.Run("should return non-empty string representation", func(t *testing.T) {
		cfg := &Config{
			CertFile:           "/etc/ssl/certs/server.crt",
			KeyFile:            "/etc/ssl/private/server.key",
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: true,
		}

		str := cfg.String()

		if str == "" {
			t.Errorf("String() = empty, want non-empty")
		}
	})
}

// TestConfigIsValid tests configuration validity check
func TestConfigIsValid(t *testing.T) {
	tests := []struct {
		name string
		cfg  *Config
		want bool
	}{
		{
			name: "should return true for valid config",
			cfg:  &Config{CertFile: "/etc/ssl/certs/server.crt", KeyFile: "/etc/ssl/private/server.key"},
			want: true,
		},
		{
			name: "should return false for invalid config",
			cfg:  &Config{CertFile: "/etc/ssl/certs/server.crt"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.IsValid()
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}