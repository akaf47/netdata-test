package ceph

import (
	"fmt"
	"os"
	"testing"
)

func TestIsValidAuthCommand(t *testing.T) {
	tests := []struct {
		name    string
		command string
		want    bool
	}{
		{
			name:    "valid ceph auth command",
			command: "ceph",
			want:    true,
		},
		{
			name:    "valid cephx auth command",
			command: "cephx",
			want:    true,
		},
		{
			name:    "valid none auth command",
			command: "none",
			want:    true,
		},
		{
			name:    "invalid empty command",
			command: "",
			want:    false,
		},
		{
			name:    "invalid unknown command",
			command: "unknown",
			want:    false,
		},
		{
			name:    "case sensitive test lowercase",
			command: "ceph",
			want:    true,
		},
		{
			name:    "case sensitive test uppercase",
			command: "CEPH",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isValidAuthCommand(tt.command)
			if got != tt.want {
				t.Errorf("isValidAuthCommand(%q) = %v, want %v", tt.command, got, tt.want)
			}
		})
	}
}

func TestParseAuthCredentials(t *testing.T) {
	tests := []struct {
		name        string
		authStr     string
		expectName  string
		expectKey   string
		expectError bool
	}{
		{
			name:        "valid auth string format",
			authStr:     "client.admin=AQDo58ddJG==",
			expectName:  "client.admin",
			expectKey:   "AQDo58ddJG==",
			expectError: false,
		},
		{
			name:        "valid different client format",
			authStr:     "client.test=key123",
			expectName:  "client.test",
			expectKey:   "key123",
			expectError: false,
		},
		{
			name:        "empty auth string",
			authStr:     "",
			expectName:  "",
			expectKey:   "",
			expectError: true,
		},
		{
			name:        "malformed no equals sign",
			authStr:     "clientadminkey",
			expectName:  "",
			expectKey:   "",
			expectError: true,
		},
		{
			name:        "malformed empty name",
			authStr:     "=key",
			expectName:  "",
			expectKey:   "",
			expectError: true,
		},
		{
			name:        "malformed empty key",
			authStr:     "client.admin=",
			expectName:  "",
			expectKey:   "",
			expectError: true,
		},
		{
			name:        "multiple equals signs",
			authStr:     "client.admin=key=value",
			expectName:  "client.admin",
			expectKey:   "key=value",
			expectError: false,
		},
		{
			name:        "whitespace in name",
			authStr:     "client admin=key",
			expectName:  "client admin",
			expectKey:   "key",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			name, key, err := parseAuthCredentials(tt.authStr)
			if (err != nil) != tt.expectError {
				t.Errorf("parseAuthCredentials(%q) error = %v, want error %v", tt.authStr, err, tt.expectError)
			}
			if name != tt.expectName {
				t.Errorf("parseAuthCredentials(%q) name = %q, want %q", tt.authStr, name, tt.expectName)
			}
			if key != tt.expectKey {
				t.Errorf("parseAuthCredentials(%q) key = %q, want %q", tt.authStr, key, tt.expectKey)
			}
		})
	}
}

func TestValidateAuthKeyFormat(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		isValid   bool
		wantError bool
	}{
		{
			name:      "valid base64 key",
			key:       "AQDo58ddJG==",
			isValid:   true,
			wantError: false,
		},
		{
			name:      "valid hex key",
			key:       "0123456789abcdef",
			isValid:   true,
			wantError: false,
		},
		{
			name:      "empty key",
			key:       "",
			isValid:   false,
			wantError: true,
		},
		{
			name:      "short valid key",
			key:       "abc",
			isValid:   true,
			wantError: false,
		},
		{
			name:      "long valid key",
			key:       "0123456789abcdef0123456789abcdef0123456789abcdef",
			isValid:   true,
			wantError: false,
		},
		{
			name:      "whitespace only",
			key:       "   ",
			isValid:   false,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAuthKeyFormat(tt.key)
			if (err != nil) != tt.wantError {
				t.Errorf("validateAuthKeyFormat(%q) error = %v, want error %v", tt.key, err, tt.wantError)
			}
		})
	}
}

func TestGetAuthMethodFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		envVar   string
		envValue string
		expected string
		cleanup  bool
	}{
		{
			name:     "CEPH_AUTH env var set",
			envVar:   "CEPH_AUTH",
			envValue: "cephx",
			expected: "cephx",
			cleanup:  true,
		},
		{
			name:     "CEPH_AUTH not set",
			envVar:   "CEPH_AUTH",
			envValue: "",
			expected: "cephx",
			cleanup:  false,
		},
		{
			name:     "custom auth env var",
			envVar:   "CUSTOM_CEPH_AUTH",
			envValue: "none",
			expected: "none",
			cleanup:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.cleanup {
				os.Setenv(tt.envVar, tt.envValue)
				defer os.Unsetenv(tt.envVar)
			} else {
				os.Unsetenv(tt.envVar)
			}

			got := getAuthMethodFromEnv(tt.envVar, "cephx")
			if got != tt.expected {
				t.Errorf("getAuthMethodFromEnv(%q, %q) = %q, want %q", tt.envVar, "cephx", got, tt.expected)
			}
		})
	}
}

func TestNewAuthConfig(t *testing.T) {
	tests := []struct {
		name        string
		authType    string
		clientName  string
		key         string
		expectError bool
	}{
		{
			name:        "valid cephx auth config",
			authType:    "cephx",
			clientName:  "client.admin",
			key:         "AQDo58ddJG==",
			expectError: false,
		},
		{
			name:        "valid none auth config",
			authType:    "none",
			clientName:  "",
			key:         "",
			expectError: false,
		},
		{
			name:        "empty auth type",
			authType:    "",
			clientName:  "client.admin",
			key:         "key",
			expectError: true,
		},
		{
			name:        "invalid auth type",
			authType:    "invalid",
			clientName:  "client.admin",
			key:         "key",
			expectError: true,
		},
		{
			name:        "cephx missing client name",
			authType:    "cephx",
			clientName:  "",
			key:         "key",
			expectError: true,
		},
		{
			name:        "cephx missing key",
			authType:    "cephx",
			clientName:  "client.admin",
			key:         "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := NewAuthConfig(tt.authType, tt.clientName, tt.key)
			if (err != nil) != tt.expectError {
				t.Errorf("NewAuthConfig(%q, %q, %q) error = %v, want error %v",
					tt.authType, tt.clientName, tt.key, err, tt.expectError)
			}
			if !tt.expectError && cfg == nil {
				t.Errorf("NewAuthConfig(%q, %q, %q) returned nil config, want non-nil",
					tt.authType, tt.clientName, tt.key)
			}
			if !tt.expectError && cfg != nil {
				if cfg.AuthType != tt.authType {
					t.Errorf("NewAuthConfig() AuthType = %q, want %q", cfg.AuthType, tt.authType)
				}
				if cfg.ClientName != tt.clientName {
					t.Errorf("NewAuthConfig() ClientName = %q, want %q", cfg.ClientName, tt.clientName)
				}
				if cfg.Key != tt.key {
					t.Errorf("NewAuthConfig() Key = %q, want %q", cfg.Key, tt.key)
				}
			}
		})
	}
}

func TestAuthConfig_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		cfg     AuthConfig
		isValid bool
	}{
		{
			name: "valid cephx config",
			cfg: AuthConfig{
				AuthType:   "cephx",
				ClientName: "client.admin",
				Key:        "AQDo58ddJG==",
			},
			isValid: true,
		},
		{
			name: "valid none config",
			cfg: AuthConfig{
				AuthType:   "none",
				ClientName: "",
				Key:        "",
			},
			isValid: true,
		},
		{
			name: "invalid empty auth type",
			cfg: AuthConfig{
				AuthType:   "",
				ClientName: "client.admin",
				Key:        "key",
			},
			isValid: false,
		},
		{
			name: "invalid unknown auth type",
			cfg: AuthConfig{
				AuthType:   "unknown",
				ClientName: "client.admin",
				Key:        "key",
			},
			isValid: false,
		},
		{
			name: "cephx missing client name",
			cfg: AuthConfig{
				AuthType:   "cephx",
				ClientName: "",
				Key:        "key",
			},
			isValid: false,
		},
		{
			name: "cephx missing key",
			cfg: AuthConfig{
				AuthType:   "cephx",
				ClientName: "client.admin",
				Key:        "",
			},
			isValid: false,
		},
		{
			name: "none with client name should still be valid",
			cfg: AuthConfig{
				AuthType:   "none",
				ClientName: "client.admin",
				Key:        "",
			},
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.IsValid()
			if got != tt.isValid {
				t.Errorf("AuthConfig.IsValid() = %v, want %v", got, tt.isValid)
			}
		})
	}
}

func TestAuthConfig_String(t *testing.T) {
	tests := []struct {
		name    string
		cfg     AuthConfig
		wantSub string
	}{
		{
			name: "cephx config string representation",
			cfg: AuthConfig{
				AuthType:   "cephx",
				ClientName: "client.admin",
				Key:        "secretkey",
			},
			wantSub: "client.admin",
		},
		{
			name: "none config string representation",
			cfg: AuthConfig{
				AuthType:   "none",
				ClientName: "",
				Key:        "",
			},
			wantSub: "none",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.String()
			if len(got) == 0 {
				t.Errorf("AuthConfig.String() returned empty string")
			}
			// Verify key is not exposed in string representation
			if tt.cfg.Key != "" && got == tt.cfg.Key {
				t.Errorf("AuthConfig.String() leaked secret key")
			}
		})
	}
}

func TestDefaultAuthConfig(t *testing.T) {
	cfg := DefaultAuthConfig()

	if cfg == nil {
		t.Fatal("DefaultAuthConfig() returned nil")
	}

	if !cfg.IsValid() {
		t.Errorf("DefaultAuthConfig() returned invalid config")
	}

	if cfg.AuthType != "none" {
		t.Errorf("DefaultAuthConfig() AuthType = %q, want 'none'", cfg.AuthType)
	}
}

func TestParseAuthFile(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expectCount int
		expectError bool
	}{
		{
			name:        "single entry auth file",
			content:     "client.admin=AQDo58ddJG==",
			expectCount: 1,
			expectError: false,
		},
		{
			name:        "multiple entries auth file",
			content:     "client.admin=AQDo58ddJG==\nclient.user=key123",
			expectCount: 2,
			expectError: false,
		},
		{
			name:        "empty auth file",
			content:     "",
			expectCount: 0,
			expectError: true,
		},
		{
			name:        "auth file with comments",
			content:     "# comment\nclient.admin=key\n# another comment",
			expectCount: 1,
			expectError: false,
		},
		{
			name:        "auth file with blank lines",
			content:     "client.admin=key\n\nclient.user=key2",
			expectCount: 2,
			expectError: false,
		},
		{
			name:        "malformed entry",
			content:     "malformed",
			expectCount: 0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpfile, err := os.CreateTemp("", "ceph_auth_*.txt")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpfile.Name())

			if _, err := tmpfile.WriteString(tt.content); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			tmpfile.Close()

			creds, err := parseAuthFile(tmpfile.Name())
			if (err != nil) != tt.expectError {
				t.Errorf("parseAuthFile(%q) error = %v, want error %v", tmpfile.Name(), err, tt.expectError)
			}
			if !tt.expectError && len(creds) != tt.expectCount {
				t.Errorf("parseAuthFile(%q) returned %d credentials, want %d", tmpfile.Name(), len(creds), tt.expectCount)
			}
		})
	}
}

func TestAuthConfig_Credentials(t *testing.T) {
	tests := []struct {
		name           string
		cfg            AuthConfig
		expectClientID string
		expectSecret   string
	}{
		{
			name: "cephx credentials",
			cfg: AuthConfig{
				AuthType:   "cephx",
				ClientName: "client.admin",
				Key:        "mysecret",
			},
			expectClientID: "client.admin",
			expectSecret:   "mysecret",
		},
		{
			name: "none auth credentials",
			cfg: AuthConfig{
				AuthType:   "none",
				ClientName: "",
				Key:        "",
			},
			expectClientID: "",
			expectSecret:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientID := tt.cfg.GetClientID()
			secret := tt.cfg.GetSecret()

			if clientID != tt.expectClientID {
				t.Errorf("GetClientID() = %q, want %q", clientID, tt.expectClientID)
			}
			if secret != tt.expectSecret {
				t.Errorf("GetSecret() = %q, want %q", secret, tt.expectSecret)
			}
		})
	}
}

func TestAuthConfigEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		cfg    AuthConfig
		testFn func(AuthConfig) error
	}{
		{
			name: "config with special characters in client name",
			cfg: AuthConfig{
				AuthType:   "cephx",
				ClientName: "client.admin@domain.com",
				Key:        "key",
			},
			testFn: func(cfg AuthConfig) error {
				if !cfg.IsValid() {
					return fmt.Errorf("config with special chars should be valid")
				}
				return nil
			},
		},
		{
			name: "config with very long key",
			cfg: AuthConfig{
				AuthType:   "cephx",
				ClientName: "client.admin",
				Key:        "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			},
			testFn: func(cfg AuthConfig) error {
				if !cfg.IsValid() {
					return fmt.Errorf("config with long key should be valid")
				}
				return nil
			},
		},
		{
			name: "config with unicode characters",
			cfg: AuthConfig{
				AuthType:   "cephx",
				ClientName: "client.测试",
				Key:        "key",
			},
			testFn: func(cfg AuthConfig) error {
				if !cfg.IsValid() {
					return fmt.Errorf("config with unicode should be valid")
				}
				return nil
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.testFn(tt.cfg); err != nil {
				t.Errorf("Edge case test failed: %v", err)
			}
		})
	}
}