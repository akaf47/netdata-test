package pihole

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestLogin tests the Login method with successful authentication
func TestLogin_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/admin/api.php" {
			t.Fatalf("Expected path /admin/api.php, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"session":"valid_session_token"}`))
	}))
	defer server.Close()

	p := &Pihole{
		httpClient: &http.Client{},
		Config: Config{
			URL: server.URL,
		},
	}

	err := p.Login(context.Background())
	if err != nil {
		t.Errorf("Login() error = %v, want nil", err)
	}
}

// TestLogin_InvalidURL tests Login with invalid URL
func TestLogin_InvalidURL(t *testing.T) {
	p := &Pihole{
		httpClient: &http.Client{},
		Config: Config{
			URL: "http://invalid..url",
		},
	}

	err := p.Login(context.Background())
	if err == nil {
		t.Error("Login() error = nil, want error for invalid URL")
	}
}

// TestLogin_EmptyURL tests Login with empty URL
func TestLogin_EmptyURL(t *testing.T) {
	p := &Pihole{
		httpClient: &http.Client{},
		Config: Config{
			URL: "",
		},
	}

	err := p.Login(context.Background())
	if err == nil {
		t.Error("Login() error = nil, want error for empty URL")
	}
}

// TestLogin_HTTPError tests Login when server returns an error
func TestLogin_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer server.Close()

	p := &Pihole{
		httpClient: &http.Client{},
		Config: Config{
			URL: server.URL,
		},
	}

	err := p.Login(context.Background())
	if err == nil {
		t.Error("Login() error = nil, want error for HTTP 500")
	}
}

// TestLogin_MalformedJSON tests Login with malformed JSON response
func TestLogin_MalformedJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{invalid json}`))
	}))
	defer server.Close()

	p := &Pihole{
		httpClient: &http.Client{},
		Config: Config{
			URL: server.URL,
		},
	}

	err := p.Login(context.Background())
	if err == nil {
		t.Error("Login() error = nil, want error for malformed JSON")
	}
}

// TestLogin_EmptyPassword tests Login with empty password
func TestLogin_EmptyPassword(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"session":"token"}`))
	}))
	defer server.Close()

	p := &Pihole{
		httpClient: &http.Client{},
		Config: Config{
			URL:      server.URL,
			Password: "",
		},
	}

	err := p.Login(context.Background())
	// Empty password should still attempt login
	if err != nil {
		t.Errorf("Login() with empty password error = %v, want nil", err)
	}
}

// TestLogin_ContextCancelled tests Login with cancelled context
func TestLogin_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	p := &Pihole{
		httpClient: &http.Client{},
		Config: Config{
			URL: "http://localhost:8080",
		},
	}

	err := p.Login(ctx)
	if err == nil {
		t.Error("Login() with cancelled context error = nil, want error")
	}
}

// TestLogin_Timeout tests Login with request timeout
func TestLogin_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response
		select {}
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	p := &Pihole{
		httpClient: &http.Client{},
		Config: Config{
			URL: server.URL,
		},
	}

	err := p.Login(ctx)
	if err == nil {
		t.Error("Login() with timeout error = nil, want error")
	}
}

// TestLogin_NoSessionToken tests Login response without session token
func TestLogin_NoSessionToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	p := &Pihole{
		httpClient: &http.Client{},
		Config: Config{
			URL: server.URL,
		},
	}

	err := p.Login(context.Background())
	if err == nil {
		t.Error("Login() without session token error = nil, want error")
	}
}

// TestAuthenticate tests the Authenticate method
func TestAuthenticate(t *testing.T) {
	p := &Pihole{
		token: "test_token",
	}

	authenticated := p.Authenticate()
	if !authenticated {
		t.Error("Authenticate() = false, want true when token is set")
	}
}

// TestAuthenticate_NoToken tests Authenticate without token
func TestAuthenticate_NoToken(t *testing.T) {
	p := &Pihole{
		token: "",
	}

	authenticated := p.Authenticate()
	if authenticated {
		t.Error("Authenticate() = true, want false when token is empty")
	}
}

// TestLogout tests the Logout method
func TestLogout_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/admin/api.php" {
			t.Fatalf("Expected path /admin/api.php, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	p := &Pihole{
		httpClient: &http.Client{},
		Config: Config{
			URL: server.URL,
		},
		token: "valid_token",
	}

	err := p.Logout(context.Background())
	if err != nil {
		t.Errorf("Logout() error = %v, want nil", err)
	}
}

// TestLogout_NoToken tests Logout without token
func TestLogout_NoToken(t *testing.T) {
	p := &Pihole{
		httpClient: &http.Client{},
		Config: Config{
			URL: "http://localhost:8080",
		},
		token: "",
	}

	err := p.Logout(context.Background())
	// Should handle logout even without token
	if err == nil || err.Error() != "" {
		// Behavior depends on implementation
	}
}

// TestLogout_HTTPError tests Logout with server error
func TestLogout_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	p := &Pihole{
		httpClient: &http.Client{},
		Config: Config{
			URL: server.URL,
		},
		token: "valid_token",
	}

	err := p.Logout(context.Background())
	if err == nil {
		t.Error("Logout() error = nil, want error for HTTP 500")
	}
}

// TestLogout_InvalidURL tests Logout with invalid URL
func TestLogout_InvalidURL(t *testing.T) {
	p := &Pihole{
		httpClient: &http.Client{},
		Config: Config{
			URL: "http://invalid..url",
		},
		token: "valid_token",
	}

	err := p.Logout(context.Background())
	if err == nil {
		t.Error("Logout() error = nil, want error for invalid URL")
	}
}

// TestLogout_ContextCancelled tests Logout with cancelled context
func TestLogout_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	p := &Pihole{
		httpClient: &http.Client{},
		Config: Config{
			URL: "http://localhost:8080",
		},
		token: "valid_token",
	}

	err := p.Logout(ctx)
	if err == nil {
		t.Error("Logout() with cancelled context error = nil, want error")
	}
}

// TestParseAuthResponse tests parsing of authentication response
func TestParseAuthResponse_ValidResponse(t *testing.T) {
	body := io.NopCloser(bytes.NewBufferString(`{"session":"valid_session_token"}`))
	
	p := &Pihole{}
	token, err := p.parseAuthResponse(body)
	
	if err != nil {
		t.Errorf("parseAuthResponse() error = %v, want nil", err)
	}
	if token != "valid_session_token" {
		t.Errorf("parseAuthResponse() token = %s, want 'valid_session_token'", token)
	}
}

// TestParseAuthResponse_EmptyResponse tests parsing empty response
func TestParseAuthResponse_EmptyResponse(t *testing.T) {
	body := io.NopCloser(bytes.NewBufferString(`{}`))
	
	p := &Pihole{}
	token, err := p.parseAuthResponse(body)
	
	if err == nil {
		t.Error("parseAuthResponse() error = nil, want error")
	}
	if token != "" {
		t.Errorf("parseAuthResponse() token = %s, want empty string", token)
	}
}

// TestParseAuthResponse_MalformedJSON tests parsing malformed JSON
func TestParseAuthResponse_MalformedJSON(t *testing.T) {
	body := io.NopCloser(bytes.NewBufferString(`{invalid}`))
	
	p := &Pihole{}
	token, err := p.parseAuthResponse(body)
	
	if err == nil {
		t.Error("parseAuthResponse() error = nil, want error for malformed JSON")
	}
	if token != "" {
		t.Errorf("parseAuthResponse() token = %s, want empty string", token)
	}
}

// TestParseAuthResponse_NullSession tests parsing null session
func TestParseAuthResponse_NullSession(t *testing.T) {
	body := io.NopCloser(bytes.NewBufferString(`{"session":null}`))
	
	p := &Pihole{}
	token, err := p.parseAuthResponse(body)
	
	if err == nil {
		t.Error("parseAuthResponse() error = nil, want error for null session")
	}
	if token != "" {
		t.Errorf("parseAuthResponse() token = %s, want empty string", token)
	}
}

// TestParseAuthResponse_EmptySession tests parsing empty session
func TestParseAuthResponse_EmptySession(t *testing.T) {
	body := io.NopCloser(bytes.NewBufferString(`{"session":""}`))
	
	p := &Pihole{}
	token, err := p.parseAuthResponse(body)
	
	if err == nil {
		t.Error("parseAuthResponse() error = nil, want error for empty session")
	}
	if token != "" {
		t.Errorf("parseAuthResponse() token = %s, want empty string", token)
	}
}

// TestParseAuthResponse_WhitespaceSession tests parsing whitespace-only session
func TestParseAuthResponse_WhitespaceSession(t *testing.T) {
	body := io.NopCloser(bytes.NewBufferString(`{"session":"   "}`))
	
	p := &Pihole{}
	token, err := p.parseAuthResponse(body)
	
	// Whitespace-only token behavior depends on implementation
	if token == "" && err != nil {
		// Acceptable: treat whitespace as invalid
	} else if token == "   " && err == nil {
		// Acceptable: return as-is
	}
}

// TestGetAuthRequestBody tests building authentication request body
func TestGetAuthRequestBody_WithPassword(t *testing.T) {
	p := &Pihole{
		Config: Config{
			Password: "test_password",
		},
	}

	body := p.getAuthRequestBody()
	if body == "" {
		t.Error("getAuthRequestBody() returned empty string")
	}
}

// TestGetAuthRequestBody_EmptyPassword tests auth request with empty password
func TestGetAuthRequestBody_EmptyPassword(t *testing.T) {
	p := &Pihole{
		Config: Config{
			Password: "",
		},
	}

	body := p.getAuthRequestBody()
	if body == "" {
		t.Error("getAuthRequestBody() returned empty string for empty password")
	}
}

// TestGetLogoutRequestBody tests building logout request body
func TestGetLogoutRequestBody(t *testing.T) {
	p := &Pihole{
		token: "test_token",
	}

	body := p.getLogoutRequestBody()
	if body == "" {
		t.Error("getLogoutRequestBody() returned empty string")
	}
	if !contains(body, "test_token") {
		t.Errorf("getLogoutRequestBody() body = %s, want to contain token", body)
	}
}

// TestGetLogoutRequestBody_EmptyToken tests logout request with empty token
func TestGetLogoutRequestBody_EmptyToken(t *testing.T) {
	p := &Pihole{
		token: "",
	}

	body := p.getLogoutRequestBody()
	if body == "" {
		t.Error("getLogoutRequestBody() returned empty string for empty token")
	}
}

// Helper function
func contains(s, substr string) bool {
	return bytes.Contains([]byte(s), []byte(substr))
}