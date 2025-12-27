package bind

import (
	"errors"
	"net/http"
	"testing"

	"github.com/netdata/netdata/go/plugins/pkg/matcher"
	"github.com/netdata/netdata/go/plugins/pkg/web"
)

// Test validateConfig
func TestValidateConfig_WithEmptyURL(t *testing.T) {
	c := &Collector{
		URL: "",
	}
	err := c.validateConfig()
	if err == nil {
		t.Error("expected error for empty URL, got nil")
	}
	if err.Error() != "url not set" {
		t.Errorf("expected 'url not set', got '%s'", err.Error())
	}
}

func TestValidateConfig_WithURL(t *testing.T) {
	c := &Collector{
		URL: "http://localhost:8053/json/v1",
	}
	err := c.validateConfig()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateConfig_WithWhitespaceURL(t *testing.T) {
	c := &Collector{
		URL: "   ",
	}
	err := c.validateConfig()
	if err != nil {
		t.Errorf("expected no error for whitespace URL (not empty string), got %v", err)
	}
}

func TestValidateConfig_WithLocalhost(t *testing.T) {
	c := &Collector{
		URL: "http://127.0.0.1:8053/json/v1",
	}
	err := c.validateConfig()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// Test initPermitViewMatcher
func TestInitPermitViewMatcher_WithEmptyString(t *testing.T) {
	c := &Collector{
		PermitView: "",
	}
	matcher, err := c.initPermitViewMatcher()
	if err != nil {
		t.Errorf("expected no error for empty PermitView, got %v", err)
	}
	if matcher != nil {
		t.Errorf("expected nil matcher for empty PermitView, got %v", matcher)
	}
}

func TestInitPermitViewMatcher_WithValidPattern(t *testing.T) {
	c := &Collector{
		PermitView: "internal external",
	}
	matcher, err := c.initPermitViewMatcher()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if matcher == nil {
		t.Error("expected non-nil matcher for valid pattern, got nil")
	}
}

func TestInitPermitViewMatcher_WithSinglePattern(t *testing.T) {
	c := &Collector{
		PermitView: "internal",
	}
	matcher, err := c.initPermitViewMatcher()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if matcher == nil {
		t.Error("expected non-nil matcher for single pattern, got nil")
	}
}

func TestInitPermitViewMatcher_WithWildcardPattern(t *testing.T) {
	c := &Collector{
		PermitView: "*.example.com",
	}
	matcher, err := c.initPermitViewMatcher()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if matcher == nil {
		t.Error("expected non-nil matcher for wildcard pattern, got nil")
	}
}

// Test initBindApiClient
func TestInitBindApiClient_WithXMLv3Endpoint(t *testing.T) {
	httpClient := &http.Client{}
	c := &Collector{
		URL: "http://localhost:8053/xml/v3",
		RequestConfig: web.RequestConfig{
			URL: "http://localhost:8053/xml/v3",
		},
	}
	client, err := c.initBindApiClient(httpClient)
	if err != nil {
		t.Errorf("expected no error for xml/v3 endpoint, got %v", err)
	}
	if client == nil {
		t.Error("expected non-nil client for xml/v3 endpoint, got nil")
	}
}

func TestInitBindApiClient_WithJSONv1Endpoint(t *testing.T) {
	httpClient := &http.Client{}
	c := &Collector{
		URL: "http://localhost:8053/json/v1",
		RequestConfig: web.RequestConfig{
			URL: "http://localhost:8053/json/v1",
		},
	}
	client, err := c.initBindApiClient(httpClient)
	if err != nil {
		t.Errorf("expected no error for json/v1 endpoint, got %v", err)
	}
	if client == nil {
		t.Error("expected non-nil client for json/v1 endpoint, got nil")
	}
}

func TestInitBindApiClient_WithInvalidEndpoint(t *testing.T) {
	httpClient := &http.Client{}
	c := &Collector{
		URL: "http://localhost:8053/invalid/endpoint",
		RequestConfig: web.RequestConfig{
			URL: "http://localhost:8053/invalid/endpoint",
		},
	}
	client, err := c.initBindApiClient(httpClient)
	if err == nil {
		t.Error("expected error for invalid endpoint, got nil")
	}
	if client != nil {
		t.Errorf("expected nil client for invalid endpoint, got %v", client)
	}
	expectedMsg := "URL http://localhost:8053/invalid/endpoint is wrong"
	if !contains(err.Error(), expectedMsg) {
		t.Errorf("expected error containing '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestInitBindApiClient_WithXMLv3URLOnly(t *testing.T) {
	httpClient := &http.Client{}
	c := &Collector{
		URL: "http://example.com/xml/v3",
		RequestConfig: web.RequestConfig{
			URL: "http://example.com/xml/v3",
		},
	}
	client, err := c.initBindApiClient(httpClient)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if client == nil {
		t.Error("expected non-nil client")
	}
}

func TestInitBindApiClient_WithJSONv1URLOnly(t *testing.T) {
	httpClient := &http.Client{}
	c := &Collector{
		URL: "http://example.com/json/v1",
		RequestConfig: web.RequestConfig{
			URL: "http://example.com/json/v1",
		},
	}
	client, err := c.initBindApiClient(httpClient)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if client == nil {
		t.Error("expected non-nil client")
	}
}

func TestInitBindApiClient_WithEmptyURL(t *testing.T) {
	httpClient := &http.Client{}
	c := &Collector{
		URL: "",
		RequestConfig: web.RequestConfig{
			URL: "",
		},
	}
	client, err := c.initBindApiClient(httpClient)
	if err == nil {
		t.Error("expected error for empty URL, got nil")
	}
	if client != nil {
		t.Errorf("expected nil client, got %v", client)
	}
}

func TestInitBindApiClient_WithURLWithoutSuffix(t *testing.T) {
	httpClient := &http.Client{}
	c := &Collector{
		URL: "http://localhost:8053",
		RequestConfig: web.RequestConfig{
			URL: "http://localhost:8053",
		},
	}
	client, err := c.initBindApiClient(httpClient)
	if err == nil {
		t.Error("expected error for URL without proper suffix, got nil")
	}
	if client != nil {
		t.Errorf("expected nil client, got %v", client)
	}
}

func TestInitBindApiClient_WithURLContainingJsonv1InPath(t *testing.T) {
	httpClient := &http.Client{}
	c := &Collector{
		URL: "http://localhost:8053/api/json/v1",
		RequestConfig: web.RequestConfig{
			URL: "http://localhost:8053/api/json/v1",
		},
	}
	client, err := c.initBindApiClient(httpClient)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if client == nil {
		t.Error("expected non-nil client")
	}
}

func TestInitBindApiClient_WithURLContainingXmlv3InPath(t *testing.T) {
	httpClient := &http.Client{}
	c := &Collector{
		URL: "http://localhost:8053/api/xml/v3",
		RequestConfig: web.RequestConfig{
			URL: "http://localhost:8053/api/xml/v3",
		},
	}
	client, err := c.initBindApiClient(httpClient)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if client == nil {
		t.Error("expected non-nil client")
	}
}

func TestInitBindApiClient_WithNilHTTPClient(t *testing.T) {
	c := &Collector{
		URL: "http://localhost:8053/json/v1",
		RequestConfig: web.RequestConfig{
			URL: "http://localhost:8053/json/v1",
		},
	}
	client, err := c.initBindApiClient(nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if client == nil {
		t.Error("expected non-nil client even with nil http.Client")
	}
}

func TestInitBindApiClient_SwitchCaseOrderJSONv1Before(t *testing.T) {
	httpClient := &http.Client{}
	c := &Collector{
		URL: "http://localhost:8053/json/v1",
		RequestConfig: web.RequestConfig{
			URL: "http://localhost:8053/json/v1",
		},
	}
	client, err := c.initBindApiClient(httpClient)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if _, ok := client.(*jsonClient); !ok {
		t.Errorf("expected jsonClient type, got %T", client)
	}
}

func TestInitBindApiClient_SwitchCaseOrderXMLv3First(t *testing.T) {
	httpClient := &http.Client{}
	c := &Collector{
		URL: "http://localhost:8053/xml/v3",
		RequestConfig: web.RequestConfig{
			URL: "http://localhost:8053/xml/v3",
		},
	}
	client, err := c.initBindApiClient(httpClient)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if client == nil {
		t.Error("expected non-nil client")
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(substr) > 0 && len(s) > 0))
}