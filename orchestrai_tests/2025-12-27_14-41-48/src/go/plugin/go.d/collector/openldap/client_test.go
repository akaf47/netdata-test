package openldap

import (
	"crypto/tls"
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockLDAPConn is a mock LDAP connection
type MockLDAPConn struct {
	mock.Mock
}

func (m *MockLDAPConn) Start() error {
	return m.Called().Error(0)
}

func (m *MockLDAPConn) StartTLS(config *tls.Config) error {
	return m.Called(config).Error(0)
}

func (m *MockLDAPConn) Bind(username, password string) error {
	return m.Called(username, password).Error(0)
}

func (m *MockLDAPConn) Search(searchRequest *ldap.SearchRequest) (*ldap.SearchResult, error) {
	args := m.Called(searchRequest)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ldap.SearchResult), args.Error(1)
}

func (m *MockLDAPConn) Close() {
	m.Called()
}

// TestNewClient tests the client constructor
func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		host    string
		port    int
		wantErr bool
	}{
		{
			name: "should create client with valid host and port",
			host: "localhost",
			port: 389,
		},
		{
			name: "should create client with LDAPS port",
			host: "localhost",
			port: 636,
		},
		{
			name: "should create client with empty host",
			host: "",
			port: 389,
		},
		{
			name: "should create client with port 0",
			host: "localhost",
			port: 0,
		},
		{
			name: "should create client with high port number",
			host: "localhost",
			port: 65535,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.host, tt.port)
			assert.NotNil(t, client)
			assert.Equal(t, tt.host, client.Host)
			assert.Equal(t, tt.port, client.Port)
		})
	}
}

// TestClientConnect tests the Connect method
func TestClientConnect(t *testing.T) {
	tests := []struct {
		name        string
		host        string
		port        int
		useTLS      bool
		wantErr     bool
		errContains string
	}{
		{
			name:    "should fail to connect to invalid host",
			host:    "invalid.host.example.com",
			port:    389,
			wantErr: true,
		},
		{
			name:    "should fail to connect to invalid port",
			host:    "localhost",
			port:    0,
			wantErr: true,
		},
		{
			name:    "should fail to connect to unreachable host",
			host:    "192.0.2.1", // TEST-NET-1 - guaranteed to be unreachable
			port:    389,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.host, tt.port)
			client.UseTLS = tt.useTLS

			err := client.Connect()
			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.ErrorContains(t, err, tt.errContains)
				}
			}

			// Cleanup
			if client.conn != nil {
				client.Close()
			}
		})
	}
}

// TestClientConnectWithTLS tests Connect with TLS enabled
func TestClientConnectWithTLS(t *testing.T) {
	tests := []struct {
		name            string
		host            string
		port            int
		useTLS          bool
		skipTLSVerify   bool
		tlsConfig       *tls.Config
		wantErr         bool
	}{
		{
			name:          "should attempt TLS connection",
			host:          "localhost",
			port:          636,
			useTLS:        true,
			skipTLSVerify: true,
			wantErr:       true, // Will fail without actual server
		},
		{
			name:          "should attempt STARTTLS",
			host:          "localhost",
			port:          389,
			useTLS:        true,
			skipTLSVerify: true,
			wantErr:       true, // Will fail without actual server
		},
		{
			name:          "should use custom TLS config",
			host:          "localhost",
			port:          636,
			useTLS:        true,
			skipTLSVerify: false,
			tlsConfig:     &tls.Config{InsecureSkipVerify: true},
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.host, tt.port)
			client.UseTLS = tt.useTLS
			client.SkipTLSVerify = tt.skipTLSVerify
			if tt.tlsConfig != nil {
				client.TLSConfig = tt.tlsConfig
			}

			err := client.Connect()
			if tt.wantErr {
				assert.Error(t, err)
			}

			if client.conn != nil {
				client.Close()
			}
		})
	}
}

// TestClientClose tests the Close method
func TestClientClose(t *testing.T) {
	tests := []struct {
		name string
		conn interface{}
	}{
		{
			name: "should close non-nil connection",
			conn: &tls.Conn{},
		},
		{
			name: "should handle nil connection gracefully",
			conn: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient("localhost", 389)
			if tt.conn != nil {
				client.conn = tt.conn
			}
			// Should not panic
			client.Close()
		})
	}
}

// TestClientBind tests the Bind method
func TestClientBind(t *testing.T) {
	tests := []struct {
		name           string
		bindDN         string
		bindPassword   string
		expectBindCall bool
		bindErr        error
		wantErr        bool
	}{
		{
			name:           "should bind with valid DN and password",
			bindDN:         "cn=admin,dc=example,dc=com",
			bindPassword:   "password123",
			expectBindCall: true,
			bindErr:        nil,
			wantErr:        false,
		},
		{
			name:           "should bind with empty DN",
			bindDN:         "",
			bindPassword:   "password",
			expectBindCall: true,
			bindErr:        nil,
			wantErr:        false,
		},
		{
			name:           "should bind with empty password",
			bindDN:         "cn=admin,dc=example,dc=com",
			bindPassword:   "",
			expectBindCall: true,
			bindErr:        nil,
			wantErr:        false,
		},
		{
			name:           "should fail on bind error",
			bindDN:         "cn=admin,dc=example,dc=com",
			bindPassword:   "wrongpassword",
			expectBindCall: true,
			bindErr:        errors.New("invalid credentials"),
			wantErr:        true,
		},
		{
			name:           "should fail when connection is nil",
			bindDN:         "cn=admin,dc=example,dc=com",
			bindPassword:   "password",
			expectBindCall: false,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient("localhost", 389)

			if tt.expectBindCall {
				mockConn := new(MockLDAPConn)
				mockConn.On("Bind", tt.bindDN, tt.bindPassword).Return(tt.bindErr)
				client.conn = mockConn
			}

			err := client.Bind(tt.bindDN, tt.bindPassword)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestClientSearch tests the Search method
func TestClientSearch(t *testing.T) {
	tests := []struct {
		name            string
		searchBaseDN    string
		searchFilter    string
		searchAttrs     []string
		expectSearchErr error
		wantErr         bool
		expectedResult  *ldap.SearchResult
	}{
		{
			name:            "should search with valid parameters",
			searchBaseDN:    "dc=example,dc=com",
			searchFilter:    "(objectClass=*)",
			searchAttrs:     []string{"cn", "mail"},
			expectSearchErr: nil,
			wantErr:         false,
			expectedResult:  &ldap.SearchResult{Entries: []*ldap.Entry{}},
		},
		{
			name:            "should search with empty attributes",
			searchBaseDN:    "dc=example,dc=com",
			searchFilter:    "(objectClass=*)",
			searchAttrs:     []string{},
			expectSearchErr: nil,
			wantErr:         false,
			expectedResult:  &ldap.SearchResult{Entries: []*ldap.Entry{}},
		},
		{
			name:            "should handle search error",
			searchBaseDN:    "dc=example,dc=com",
			searchFilter:    "(objectClass=invalid)",
			searchAttrs:     []string{"cn"},
			expectSearchErr: errors.New("search failed"),
			wantErr:         true,
			expectedResult:  nil,
		},
		{
			name:            "should fail when connection is nil",
			searchBaseDN:    "dc=example,dc=com",
			searchFilter:    "(objectClass=*)",
			searchAttrs:     []string{"cn"},
			expectSearchErr: nil,
			wantErr:         true,
		},
		{
			name:            "should search with complex filter",
			searchBaseDN:    "ou=users,dc=example,dc=com",
			searchFilter:    "(&(objectClass=person)(cn=*admin*))",
			searchAttrs:     []string{"cn", "mail", "userAccountControl"},
			expectSearchErr: nil,
			wantErr:         false,
			expectedResult:  &ldap.SearchResult{Entries: []*ldap.Entry{}},
		},
		{
			name:            "should search with empty base DN",
			searchBaseDN:    "",
			searchFilter:    "(objectClass=*)",
			searchAttrs:     []string{"cn"},
			expectSearchErr: nil,
			wantErr:         false,
			expectedResult:  &ldap.SearchResult{Entries: []*ldap.Entry{}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient("localhost", 389)

			if tt.expectedResult != nil || tt.expectSearchErr != nil {
				mockConn := new(MockLDAPConn)
				mockConn.On("Search", mock.MatchedBy(func(sr *ldap.SearchRequest) bool {
					return sr.BaseDN == tt.searchBaseDN && sr.Filter == tt.searchFilter
				})).Return(tt.expectedResult, tt.expectSearchErr)
				client.conn = mockConn
			}

			result, err := client.Search(tt.searchBaseDN, tt.searchFilter, tt.searchAttrs)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

// TestClientIsConnected tests the IsConnected method
func TestClientIsConnected(t *testing.T) {
	tests := []struct {
		name      string
		connected bool
	}{
		{
			name:      "should return true when connected",
			connected: true,
		},
		{
			name:      "should return false when disconnected",
			connected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient("localhost", 389)
			if tt.connected {
				client.conn = &MockLDAPConn{}
			}

			assert.Equal(t, tt.connected, client.IsConnected())
		})
	}
}

// TestClientConcurrentOperations tests thread-safe operations
func TestClientConcurrentOperations(t *testing.T) {
	client := NewClient("localhost", 389)
	mockConn := new(MockLDAPConn)
	mockConn.On("Bind", mock.Anything, mock.Anything).Return(nil)
	mockConn.On("Search", mock.Anything).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
	client.conn = mockConn

	var wg sync.WaitGroup
	numGoroutines := 10

	// Test concurrent Bind operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = client.Bind("cn=user", "password")
		}()
	}

	// Test concurrent Search operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = client.Search("dc=example,dc=com", "(objectClass=*)", []string{"cn"})
		}()
	}

	// Test concurrent IsConnected checks
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = client.IsConnected()
		}()
	}

	wg.Wait()
	mockConn.AssertExpectations(t)
}

// TestClientSearchWithTimeout tests Search with timeout scenarios
func TestClientSearchWithTimeout(t *testing.T) {
	client := NewClient("localhost", 389)
	mockConn := new(MockLDAPConn)

	// Simulate timeout
	mockConn.On("Search", mock.Anything).Return(nil, errors.New("deadline exceeded"))
	client.conn = mockConn

	_, err := client.Search("dc=example,dc=com", "(objectClass=*)", []string{"cn"})
	assert.Error(t, err)
}

// TestClientSearchMultipleEntries tests search returning multiple entries
func TestClientSearchMultipleEntries(t *testing.T) {
	client := NewClient("localhost", 389)
	mockConn := new(MockLDAPConn)

	entries := []*ldap.Entry{
		{
			DN: "cn=user1,dc=example,dc=com",
			Attributes: []*ldap.EntryAttribute{
				{
					Name:   "cn",
					Values: []string{"user1"},
				},
			},
		},
		{
			DN: "cn=user2,dc=example,dc=com",
			Attributes: []*ldap.EntryAttribute{
				{
					Name:   "cn",
					Values: []string{"user2"},
				},
			},
		},
	}

	mockConn.On("Search", mock.Anything).Return(&ldap.SearchResult{Entries: entries}, nil)
	client.conn = mockConn

	result, err := client.Search("dc=example,dc=com", "(objectClass=*)", []string{"cn"})
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Entries, 2)
}

// TestClientReconnectAfterClose tests reconnection after closing
func TestClientReconnectAfterClose(t *testing.T) {
	client := NewClient("localhost", 389)
	mockConn := new(MockLDAPConn)
	mockConn.On("Close").Return()
	client.conn = mockConn

	client.Close()
	assert.False(t, client.IsConnected())

	// Attempting to reconnect should work
	client.conn = mockConn
	assert.True(t, client.IsConnected())
}

// TestClientBindAfterConnect tests the full connect and bind flow
func TestClientBindAfterConnect(t *testing.T) {
	client := NewClient("localhost", 389)
	mockConn := new(MockLDAPConn)
	mockConn.On("Bind", "cn=admin,dc=example,dc=com", "password").Return(nil)
	client.conn = mockConn

	err := client.Bind("cn=admin,dc=example,dc=com", "password")
	assert.NoError(t, err)
}

// TestClientSearchWithSpecialCharacters tests search with special LDAP characters
func TestClientSearchWithSpecialCharacters(t *testing.T) {
	client := NewClient("localhost", 389)
	mockConn := new(MockLDAPConn)

	// Test filter with special characters: * ( ) \ NUL
	specialFilter := "(cn=user\\2atest)"
	mockConn.On("Search", mock.MatchedBy(func(sr *ldap.SearchRequest) bool {
		return sr.Filter == specialFilter
	})).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
	client.conn = mockConn

	result, err := client.Search("dc=example,dc=com", specialFilter, []string{"cn"})
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

// TestClientMultipleSearches tests consecutive searches
func TestClientMultipleSearches(t *testing.T) {
	client := NewClient("localhost", 389)
	mockConn := new(MockLDAPConn)

	mockConn.On("Search", mock.Anything).Return(&ldap.SearchResult{Entries: []*ldap.Entry{}}, nil)
	client.conn = mockConn

	for i := 0; i < 5; i++ {
		_, err := client.Search("dc=example,dc=com", "(objectClass=*)", []string{"cn"})
		assert.NoError(t, err)
	}

	mockConn.AssertNumberOfCalls(t, "Search", 5)
}

// TestNewClientPortRange tests valid port ranges
func TestNewClientPortRange(t *testing.T) {
	testPorts := []int{1, 389, 636, 3268, 3269, 65535}

	for _, port := range testPorts {
		client := NewClient("localhost", port)
		assert.Equal(t, port, client.Port)
	}
}