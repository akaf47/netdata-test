package dbdriver

import (
	"testing"
)

// TestParseDSN tests the DSN parsing functionality
func TestParseDSN(t *testing.T) {
	tests := []struct {
		name    string
		dsn     string
		want    *DSN
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid dsn with all parameters",
			dsn:  "hostname=localhost;port=50000;database=testdb;uid=user;pwd=password",
			want: &DSN{
				Hostname: "localhost",
				Port:     50000,
				Database: "testdb",
				UID:      "user",
				PWD:      "password",
			},
			wantErr: false,
		},
		{
			name: "valid dsn with minimal parameters",
			dsn:  "hostname=localhost;port=50000;database=testdb",
			want: &DSN{
				Hostname: "localhost",
				Port:     50000,
				Database: "testdb",
			},
			wantErr: false,
		},
		{
			name: "empty dsn string",
			dsn:  "",
			want: &DSN{},
			wantErr: false,
		},
		{
			name: "dsn with spaces in values",
			dsn:  "hostname=local host;port=50000;database=test db",
			want: &DSN{
				Hostname: "local host",
				Port:     50000,
				Database: "test db",
			},
			wantErr: false,
		},
		{
			name: "dsn with invalid port",
			dsn:  "hostname=localhost;port=invalid;database=testdb",
			want: nil,
			wantErr: true,
		},
		{
			name: "dsn with negative port",
			dsn:  "hostname=localhost;port=-1;database=testdb",
			want: nil,
			wantErr: true,
		},
		{
			name: "dsn with port zero",
			dsn:  "hostname=localhost;port=0;database=testdb",
			want: nil,
			wantErr: true,
		},
		{
			name: "dsn with port exceeding max",
			dsn:  "hostname=localhost;port=65536;database=testdb",
			want: nil,
			wantErr: true,
		},
		{
			name: "dsn with uppercase keys",
			dsn:  "HOSTNAME=localhost;PORT=50000;DATABASE=testdb",
			want: &DSN{
				Hostname: "localhost",
				Port:     50000,
				Database: "testdb",
			},
			wantErr: false,
		},
		{
			name: "dsn with mixed case keys",
			dsn:  "HostName=localhost;Port=50000;DataBase=testdb",
			want: &DSN{
				Hostname: "localhost",
				Port:     50000,
				Database: "testdb",
			},
			wantErr: false,
		},
		{
			name: "dsn with extra spaces",
			dsn:  "hostname = localhost ; port = 50000 ; database = testdb",
			want: &DSN{
				Hostname: "localhost",
				Port:     50000,
				Database: "testdb",
			},
			wantErr: false,
		},
		{
			name: "dsn with duplicate keys (last wins)",
			dsn:  "hostname=localhost1;hostname=localhost2;port=50000;database=testdb",
			want: &DSN{
				Hostname: "localhost2",
				Port:     50000,
				Database: "testdb",
			},
			wantErr: false,
		},
		{
			name: "dsn with empty value",
			dsn:  "hostname=;port=50000;database=testdb",
			want: &DSN{
				Hostname: "",
				Port:     50000,
				Database: "testdb",
			},
			wantErr: false,
		},
		{
			name: "dsn with special characters in password",
			dsn:  "hostname=localhost;port=50000;database=testdb;uid=user;pwd=p@ss!w0rd;extra=value",
			want: &DSN{
				Hostname: "localhost",
				Port:     50000,
				Database: "testdb",
				UID:      "user",
				PWD:      "p@ss!w0rd;extra=value",
			},
			wantErr: false,
		},
		{
			name: "dsn with missing equals sign",
			dsn:  "hostname:localhost;port:50000",
			want: nil,
			wantErr: true,
		},
		{
			name: "dsn with trailing semicolon",
			dsn:  "hostname=localhost;port=50000;database=testdb;",
			want: &DSN{
				Hostname: "localhost",
				Port:     50000,
				Database: "testdb",
			},
			wantErr: false,
		},
		{
			name: "dsn with leading semicolon",
			dsn:  ";hostname=localhost;port=50000;database=testdb",
			want: &DSN{
				Hostname: "localhost",
				Port:     50000,
				Database: "testdb",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDSN(tt.dsn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDSN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !dsnEqual(got, tt.want) {
				t.Errorf("ParseDSN() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

// TestValidateDSN tests the DSN validation
func TestValidateDSN(t *testing.T) {
	tests := []struct {
		name    string
		dsn     *DSN
		wantErr bool
	}{
		{
			name: "valid dsn",
			dsn: &DSN{
				Hostname: "localhost",
				Port:     50000,
				Database: "testdb",
			},
			wantErr: false,
		},
		{
			name: "nil dsn",
			dsn:  nil,
			wantErr: true,
		},
		{
			name: "dsn with empty hostname",
			dsn: &DSN{
				Hostname: "",
				Port:     50000,
				Database: "testdb",
			},
			wantErr: true,
		},
		{
			name: "dsn with empty database",
			dsn: &DSN{
				Hostname: "localhost",
				Port:     50000,
				Database: "",
			},
			wantErr: true,
		},
		{
			name: "dsn with invalid port",
			dsn: &DSN{
				Hostname: "localhost",
				Port:     -1,
				Database: "testdb",
			},
			wantErr: true,
		},
		{
			name: "dsn with port zero",
			dsn: &DSN{
				Hostname: "localhost",
				Port:     0,
				Database: "testdb",
			},
			wantErr: true,
		},
		{
			name: "dsn with port exceeding max",
			dsn: &DSN{
				Hostname: "localhost",
				Port:     65536,
				Database: "testdb",
			},
			wantErr: true,
		},
		{
			name: "dsn with valid max port",
			dsn: &DSN{
				Hostname: "localhost",
				Port:     65535,
				Database: "testdb",
			},
			wantErr: false,
		},
		{
			name: "dsn with valid min port",
			dsn: &DSN{
				Hostname: "localhost",
				Port:     1,
				Database: "testdb",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDSN(tt.dsn)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDSN() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestDSNString tests the string representation of DSN
func TestDSNString(t *testing.T) {
	tests := []struct {
		name string
		dsn  *DSN
		want string
	}{
		{
			name: "full dsn",
			dsn: &DSN{
				Hostname: "localhost",
				Port:     50000,
				Database: "testdb",
				UID:      "user",
				PWD:      "password",
			},
			want: "hostname=localhost;port=50000;database=testdb;uid=user;pwd=password",
		},
		{
			name: "dsn without credentials",
			dsn: &DSN{
				Hostname: "localhost",
				Port:     50000,
				Database: "testdb",
			},
			want: "hostname=localhost;port=50000;database=testdb",
		},
		{
			name: "empty dsn",
			dsn: &DSN{},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.dsn.String()
			if got != tt.want {
				t.Errorf("DSN.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDSNClone tests the cloning of DSN
func TestDSNClone(t *testing.T) {
	original := &DSN{
		Hostname: "localhost",
		Port:     50000,
		Database: "testdb",
		UID:      "user",
		PWD:      "password",
	}

	cloned := original.Clone()

	if !dsnEqual(cloned, original) {
		t.Errorf("Clone() = %+v, want %+v", cloned, original)
	}

	// Verify deep copy
	cloned.Hostname = "modified"
	if original.Hostname == "modified" {
		t.Error("Clone() did not create a deep copy")
	}
}

// TestDSNMerge tests merging two DSN structures
func TestDSNMerge(t *testing.T) {
	tests := []struct {
		name  string
		dsn1  *DSN
		dsn2  *DSN
		want  *DSN
	}{
		{
			name: "merge with non-empty fields",
			dsn1: &DSN{
				Hostname: "localhost",
				Port:     50000,
			},
			dsn2: &DSN{
				Hostname: "remotehost",
				Database: "testdb",
			},
			want: &DSN{
				Hostname: "remotehost",
				Port:     50000,
				Database: "testdb",
			},
		},
		{
			name: "merge with empty dsn2",
			dsn1: &DSN{
				Hostname: "localhost",
				Port:     50000,
				Database: "testdb",
			},
			dsn2: &DSN{},
			want: &DSN{
				Hostname: "localhost",
				Port:     50000,
				Database: "testdb",
			},
		},
		{
			name: "merge nil with valid",
			dsn1: nil,
			dsn2: &DSN{
				Hostname: "localhost",
				Port:     50000,
			},
			want: &DSN{
				Hostname: "localhost",
				Port:     50000,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeDSN(tt.dsn1, tt.dsn2)
			if !dsnEqual(got, tt.want) {
				t.Errorf("MergeDSN() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

// Helper function to compare DSN objects
func dsnEqual(a, b *DSN) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return a.Hostname == b.Hostname &&
		a.Port == b.Port &&
		a.Database == b.Database &&
		a.UID == b.UID &&
		a.PWD == b.PWD
}