package dbdriver

import (
	"errors"
	"os"
	"testing"
)

// TestNewConfig tests the NewConfig constructor
func TestNewConfig(t *testing.T) {
	tests := []struct {
		name      string
		database  string
		user      string
		password  string
		port      int
		hostname  string
		expect    *Config
		expectErr bool
	}{
		{
			name:      "valid config",
			database:  "testdb",
			user:      "testuser",
			password:  "testpass",
			port:      50000,
			hostname:  "localhost",
			expect: &Config{
				Database: "testdb",
				User:     "testuser",
				Password: "testpass",
				Port:     50000,
				Hostname: "localhost",
			},
			expectErr: false,
		},
		{
			name:      "config with empty database",
			database:  "",
			user:      "testuser",
			password:  "testpass",
			port:      50000,
			hostname:  "localhost",
			expectErr: true,
		},
		{
			name:      "config with empty user",
			database:  "testdb",
			user:      "",
			password:  "testpass",
			port:      50000,
			hostname:  "localhost",
			expectErr: true,
		},
		{
			name:      "config with empty password",
			database:  "testdb",
			user:      "testuser",
			password:  "",
			port:      50000,
			hostname:  "localhost",
			expectErr: true,
		},
		{
			name:      "config with invalid port zero",
			database:  "testdb",
			user:      "testuser",
			password:  "testpass",
			port:      0,
			hostname:  "localhost",
			expectErr: true,
		},
		{
			name:      "config with invalid port negative",
			database:  "testdb",
			user:      "testuser",
			password:  "testpass",
			port:      -1,
			hostname:  "localhost",
			expectErr: true,
		},
		{
			name:      "config with port too high",
			database:  "testdb",
			user:      "testuser",
			password:  "testpass",
			port:      65536,
			hostname:  "localhost",
			expectErr: true,
		},
		{
			name:      "config with empty hostname",
			database:  "testdb",
			user:      "testuser",
			password:  "testpass",
			port:      50000,
			hostname:  "",
			expectErr: true,
		},
		{
			name:      "config with min valid port",
			database:  "testdb",
			user:      "testuser",
			password:  "testpass",
			port:      1,
			hostname:  "localhost",
			expectErr: false,
		},
		{
			name:      "config with max valid port",
			database:  "testdb",
			user:      "testuser",
			password:  "testpass",
			port:      65535,
			hostname:  "localhost",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := NewConfig(tt.database, tt.user, tt.password, tt.port, tt.hostname)

			if (err != nil) != tt.expectErr {
				t.Errorf("NewConfig() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			if !tt.expectErr && cfg == nil {
				t.Errorf("NewConfig() returned nil config when error was not expected")
				return
			}

			if !tt.expectErr && cfg != nil {
				if cfg.Database != tt.expect.Database {
					t.Errorf("Database mismatch: got %s, want %s", cfg.Database, tt.expect.Database)
				}
				if cfg.User != tt.expect.User {
					t.Errorf("User mismatch: got %s, want %s", cfg.User, tt.expect.User)
				}
				if cfg.Password != tt.expect.Password {
					t.Errorf("Password mismatch: got %s, want %s", cfg.Password, tt.expect.Password)
				}
				if cfg.Port != tt.expect.Port {
					t.Errorf("Port mismatch: got %d, want %d", cfg.Port, tt.expect.Port)
				}
				if cfg.Hostname != tt.expect.Hostname {
					t.Errorf("Hostname mismatch: got %s, want %s", cfg.Hostname, tt.expect.Hostname)
				}
			}
		})
	}
}

// TestNewConfigFromEnv tests creating config from environment variables
func TestNewConfigFromEnv(t *testing.T) {
	tests := []struct {
		name      string
		env       map[string]string
		expect    *Config
		expectErr bool
	}{
		{
			name: "valid env config",
			env: map[string]string{
				"DB_DATABASE": "testdb",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_PORT":     "50000",
				"DB_HOSTNAME": "localhost",
			},
			expect: &Config{
				Database: "testdb",
				User:     "testuser",
				Password: "testpass",
				Port:     50000,
				Hostname: "localhost",
			},
			expectErr: false,
		},
		{
			name: "missing database env var",
			env: map[string]string{
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_PORT":     "50000",
				"DB_HOSTNAME": "localhost",
			},
			expectErr: true,
		},
		{
			name: "missing user env var",
			env: map[string]string{
				"DB_DATABASE": "testdb",
				"DB_PASSWORD": "testpass",
				"DB_PORT":     "50000",
				"DB_HOSTNAME": "localhost",
			},
			expectErr: true,
		},
		{
			name: "invalid port format",
			env: map[string]string{
				"DB_DATABASE": "testdb",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_PORT":     "invalid",
				"DB_HOSTNAME": "localhost",
			},
			expectErr: true,
		},
		{
			name: "port out of range",
			env: map[string]string{
				"DB_DATABASE": "testdb",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_PORT":     "99999",
				"DB_HOSTNAME": "localhost",
			},
			expectErr: true,
		},
		{
			name: "empty env vars",
			env: map[string]string{
				"DB_DATABASE": "",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_PORT":     "50000",
				"DB_HOSTNAME": "localhost",
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			for key, value := range tt.env {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			cfg, err := NewConfigFromEnv()

			if (err != nil) != tt.expectErr {
				t.Errorf("NewConfigFromEnv() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			if !tt.expectErr && cfg == nil {
				t.Errorf("NewConfigFromEnv() returned nil config when error was not expected")
				return
			}

			if !tt.expectErr && cfg != nil {
				if cfg.Database != tt.expect.Database {
					t.Errorf("Database mismatch: got %s, want %s", cfg.Database, tt.expect.Database)
				}
				if cfg.User != tt.expect.User {
					t.Errorf("User mismatch: got %s, want %s", cfg.User, tt.expect.User)
				}
				if cfg.Password != tt.expect.Password {
					t.Errorf("Password mismatch: got %s, want %s", cfg.Password, tt.expect.Password)
				}
				if cfg.Port != tt.expect.Port {
					t.Errorf("Port mismatch: got %d, want %d", cfg.Port, tt.expect.Port)
				}
				if cfg.Hostname != tt.expect.Hostname {
					t.Errorf("Hostname mismatch: got %s, want %s", cfg.Hostname, tt.expect.Hostname)
				}
			}
		})
	}
}

// TestConfigValidate tests the validation of config
func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		expectErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Database: "testdb",
				User:     "testuser",
				Password: "testpass",
				Port:     50000,
				Hostname: "localhost",
			},
			expectErr: false,
		},
		{
			name: "empty database",
			config: &Config{
				Database: "",
				User:     "testuser",
				Password: "testpass",
				Port:     50000,
				Hostname: "localhost",
			},
			expectErr: true,
		},
		{
			name: "empty user",
			config: &Config{
				Database: "testdb",
				User:     "",
				Password: "testpass",
				Port:     50000,
				Hostname: "localhost",
			},
			expectErr: true,
		},
		{
			name: "empty password",
			config: &Config{
				Database: "testdb",
				User:     "testuser",
				Password: "",
				Port:     50000,
				Hostname: "localhost",
			},
			expectErr: true,
		},
		{
			name: "invalid port",
			config: &Config{
				Database: "testdb",
				User:     "testuser",
				Password: "testpass",
				Port:     0,
				Hostname: "localhost",
			},
			expectErr: true,
		},
		{
			name: "empty hostname",
			config: &Config{
				Database: "testdb",
				User:     "testuser",
				Password: "testpass",
				Port:     50000,
				Hostname: "",
			},
			expectErr: true,
		},
		{
			name: "nil config",
			config: nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var cfg *Config
			if tt.config != nil {
				cfg = tt.config
			}

			err := ValidateConfig(cfg)

			if (err != nil) != tt.expectErr {
				t.Errorf("ValidateConfig() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

// TestConfigString tests the string representation of config
func TestConfigString(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		expect string
	}{
		{
			name: "config with all fields",
			config: &Config{
				Database: "testdb",
				User:     "testuser",
				Password: "testpass",
				Port:     50000,
				Hostname: "localhost",
			},
			expect: "Config{Database: testdb, User: testuser, Port: 50000, Hostname: localhost}",
		},
		{
			name: "config with empty fields",
			config: &Config{
				Database: "",
				User:     "",
				Password: "",
				Port:     0,
				Hostname: "",
			},
			expect: "Config{Database: , User: , Port: 0, Hostname: }",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.String()

			if result != tt.expect {
				t.Errorf("String() got %s, want %s", result, tt.expect)
			}
		})
	}
}

// TestConfigClone tests cloning a config
func TestConfigClone(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
	}{
		{
			name: "clone valid config",
			config: &Config{
				Database: "testdb",
				User:     "testuser",
				Password: "testpass",
				Port:     50000,
				Hostname: "localhost",
			},
		},
		{
			name: "clone empty config",
			config: &Config{
				Database: "",
				User:     "",
				Password: "",
				Port:     0,
				Hostname: "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cloned := tt.config.Clone()

			if cloned == nil {
				t.Errorf("Clone() returned nil")
				return
			}

			if cloned == tt.config {
				t.Errorf("Clone() returned the same object, expected a copy")
				return
			}

			if cloned.Database != tt.config.Database {
				t.Errorf("Database mismatch: got %s, want %s", cloned.Database, tt.config.Database)
			}
			if cloned.User != tt.config.User {
				t.Errorf("User mismatch: got %s, want %s", cloned.User, tt.config.User)
			}
			if cloned.Password != tt.config.Password {
				t.Errorf("Password mismatch: got %s, want %s", cloned.Password, tt.config.Password)
			}
			if cloned.Port != tt.config.Port {
				t.Errorf("Port mismatch: got %d, want %d", cloned.Port, tt.config.Port)
			}
			if cloned.Hostname != tt.config.Hostname {
				t.Errorf("Hostname mismatch: got %s, want %s", cloned.Hostname, tt.config.Hostname)
			}
		})
	}
}

// TestConfigGetConnectionString tests connection string generation
func TestConfigGetConnectionString(t *testing.T) {
	tests := []struct {
		name   string
		config *Config
		want   string
	}{
		{
			name: "standard connection string",
			config: &Config{
				Database: "testdb",
				User:     "testuser",
				Password: "testpass",
				Port:     50000,
				Hostname: "localhost",
			},
			want: "Server=localhost:50000;Database=testdb;UID=testuser;PWD=testpass;",
		},
		{
			name: "connection string with special characters",
			config: &Config{
				Database: "test-db",
				User:     "test.user",
				Password: "test@pass",
				Port:     50000,
				Hostname: "db.example.com",
			},
			want: "Server=db.example.com:50000;Database=test-db;UID=test.user;PWD=test@pass;",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.GetConnectionString()
			if got != tt.want {
				t.Errorf("GetConnectionString() got %s, want %s", got, tt.want)
			}
		})
	}
}

// TestConfigEdgeCases tests edge cases and boundary conditions
func TestConfigEdgeCases(t *testing.T) {
	t.Run("very long database name", func(t *testing.T) {
		longName := ""
		for i := 0; i < 1000; i++ {
			longName += "a"
		}
		cfg, err := NewConfig(longName, "user", "pass", 50000, "localhost")
		if err != nil {
			t.Errorf("NewConfig() with long name failed: %v", err)
		}
		if cfg == nil {
			t.Errorf("NewConfig() returned nil")
		}
	})

	t.Run("special characters in password", func(t *testing.T) {
		specialPass := "p@ss!w#rd$%^&*()"
		cfg, err := NewConfig("db", "user", specialPass, 50000, "localhost")
		if err != nil {
			t.Errorf("NewConfig() with special chars failed: %v", err)
		}
		if cfg.Password != specialPass {
			t.Errorf("Password mismatch: got %s, want %s", cfg.Password, specialPass)
		}
	})

	t.Run("whitespace in values", func(t *testing.T) {
		cfg, err := NewConfig("db ", "user ", "pass ", 50000, "localhost ")
		if err == nil {
			t.Errorf("NewConfig() should reject whitespace-padded values")
		}
	})

	t.Run("unicode characters in database name", func(t *testing.T) {
		cfg, err := NewConfig("тестбд", "user", "pass", 50000, "localhost")
		if err != nil {
			t.Errorf("NewConfig() with unicode failed: %v", err)
		}
		if cfg.Database != "тестбд" {
			t.Errorf("Unicode database name not preserved")
		}
	})
}

// TestConfigMerge tests merging two configs
func TestConfigMerge(t *testing.T) {
	tests := []struct {
		name     string
		config1  *Config
		config2  *Config
		expectErr bool
	}{
		{
			name: "merge valid configs",
			config1: &Config{
				Database: "db1",
				User:     "user1",
				Password: "pass1",
				Port:     50000,
				Hostname: "host1",
			},
			config2: &Config{
				Database: "db2",
				User:     "user2",
				Password: "pass2",
				Port:     50001,
				Hostname: "host2",
			},
			expectErr: false,
		},
		{
			name: "merge nil configs",
			config1: nil,
			config2: nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config1.Merge(tt.config2)

			if (err != nil) != tt.expectErr {
				t.Errorf("Merge() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}