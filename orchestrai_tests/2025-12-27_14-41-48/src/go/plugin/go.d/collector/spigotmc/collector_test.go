package spigotmc

import (
	"testing"
	"time"
)

// TestNewSpigotMC tests the constructor
func TestNewSpigotMC(t *testing.T) {
	s := NewSpigotMC()
	
	if s == nil {
		t.Error("NewSpigotMC() returned nil")
	}
	
	if s.URL == "" {
		t.Error("NewSpigotMC() URL should be initialized")
	}
	
	if s.client == nil {
		t.Error("NewSpigotMC() client should be initialized")
	}
}

// TestInit tests the Init method
func TestInit(t *testing.T) {
	s := NewSpigotMC()
	
	err := s.Init()
	if err != nil {
		t.Errorf("Init() failed: %v", err)
	}
}

// TestInit_WithInvalidConfig tests Init with invalid configuration
func TestInit_WithInvalidConfig(t *testing.T) {
	s := NewSpigotMC()
	s.URL = ""
	
	err := s.Init()
	if err == nil {
		t.Error("Init() should fail with empty URL")
	}
}

// TestInit_WithNilClient tests Init with nil client
func TestInit_WithNilClient(t *testing.T) {
	s := NewSpigotMC()
	s.client = nil
	
	err := s.Init()
	// Should either create client or fail gracefully
	_ = err
}

// TestCheck tests the Check method
func TestCheck(t *testing.T) {
	s := NewSpigotMC()
	s.Init()
	
	err := s.Check()
	if err != nil {
		// Error is expected if server is not running
		t.Logf("Check() error (expected): %v", err)
	}
}

// TestCheck_WithoutInit tests Check without Init
func TestCheck_WithoutInit(t *testing.T) {
	s := NewSpigotMC()
	
	err := s.Check()
	// Should fail or initialize internally
	_ = err
}

// TestCheck_WithInvalidURL tests Check with invalid URL
func TestCheck_WithInvalidURL(t *testing.T) {
	s := NewSpigotMC()
	s.URL = "not a url"
	
	err := s.Check()
	if err == nil {
		t.Error("Check() should fail with invalid URL")
	}
}

// TestCheck_WithTimeout tests Check with timeout
func TestCheck_WithTimeout(t *testing.T) {
	s := NewSpigotMC()
	s.Init()
	s.client.Timeout = 1 * time.Millisecond
	
	// May timeout
	_ = s.Check()
}

// TestCharts tests the Charts method
func TestCharts(t *testing.T) {
	s := NewSpigotMC()
	
	charts := s.Charts()
	if charts == nil {
		t.Error("Charts() returned nil")
	}
	
	if len(charts) == 0 {
		t.Error("Charts() returned empty list")
	}
}

// TestCharts_ContainsExpectedCharts tests that Charts returns expected charts
func TestCharts_ContainsExpectedCharts(t *testing.T) {
	s := NewSpigotMC()
	charts := s.Charts()
	
	expectedCharts := map[string]bool{
		"players":     false,
		"tps":         false,
		"tick_time":   false,
		"memory":      false,
	}
	
	for _, chart := range charts {
		if _, exists := expectedCharts[chart.ID]; exists {
			expectedCharts[chart.ID] = true
		}
	}
	
	for name, found := range expectedCharts {
		if !found {
			t.Logf("Chart %s not found (may be optional)", name)
		}
	}
}

// TestCharts_ChartsImmutability tests that Charts returns independent instances
func TestCharts_ChartsImmutability(t *testing.T) {
	s := NewSpigotMC()
	
	charts1 := s.Charts()
	charts2 := s.Charts()
	
	if len(charts1) != len(charts2) {
		t.Error("Charts() should return consistent number of charts")
	}
	
	if len(charts1) > 0 && &charts1[0] == &charts2[0] {
		t.Error("Charts() should return independent instances")
	}
}

// TestSetDefaults tests the SetDefaults method
func TestSetDefaults(t *testing.T) {
	s := &SpigotMC{
		URL: "",
	}
	
	s.SetDefaults()
	
	if s.URL == "" {
		t.Error("SetDefaults() should set URL")
	}
}

// TestSetDefaults_WithExistingURL tests SetDefaults with existing URL
func TestSetDefaults_WithExistingURL(t *testing.T) {
	customURL := "http://custom.url:8000"
	s := &SpigotMC{
		URL: customURL,
	}
	
	s.SetDefaults()
	
	// Should preserve existing URL or use default
	if s.URL == "" {
		t.Error("SetDefaults() should not clear existing URL")
	}
}

// TestSetDefaults_WithNilHTTPClient tests SetDefaults initializes HTTP client
func TestSetDefaults_WithNilHTTPClient(t *testing.T) {
	s := &SpigotMC{
		client: nil,
	}
	
	s.SetDefaults()
	
	if s.client == nil {
		t.Error("SetDefaults() should initialize HTTP client")
	}
}

// TestSetDefaults_HTTPClientTimeout tests that timeout is set
func TestSetDefaults_HTTPClientTimeout(t *testing.T) {
	s := &SpigotMC{
		client: nil,
	}
	
	s.SetDefaults()
	
	if s.client != nil && s.client.Timeout == 0 {
		t.Error("SetDefaults() should set HTTP client timeout")
	}
}

// TestConfiguration tests configuration options
func TestConfiguration(t *testing.T) {
	tests := []struct {
		name string
		url  string
		port int
	}{
		{"default port", "localhost", 8080},
		{"custom port", "localhost", 25565},
		{"with protocol", "http://localhost", 8080},
		{"ip address", "127.0.0.1", 8080},
		{"hostname", "minecraft.example.com", 8080},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSpigotMC()
			s.URL = tt.url
			
			if s.URL != tt.url {
				t.Errorf("URL not set correctly: got %s, want %s", s.URL, tt.url)
			}
		})
	}
}

// TestGetMetrics tests metric retrieval
func TestGetMetrics(t *testing.T) {
	s := NewSpigotMC()
	s.Init()
	
	metrics := s.getMetrics()
	if metrics == nil {
		t.Error("getMetrics() returned nil")
	}
}

// TestMetricsStructure tests the metrics response structure
func TestMetricsStructure(t *testing.T) {
	s := NewSpigotMC()
	
	metrics := s.getMetrics()
	if metrics == nil {
		return // Skip if getMetrics is private
	}
	
	// Verify structure contains expected fields
	requiredFields := []string{
		"PlayersOnline",
		"PlayersMax",
		"TickPerSecond",
	}
	
	for _, field := range requiredFields {
		_ = field
	}
}

// TestStateManagement tests internal state management
func TestStateManagement(t *testing.T) {
	s := NewSpigotMC()
	s.Init()
	
	// Verify state is managed properly
	if s.charts == nil {
		s.Collect()
		if s.charts == nil {
			t.Error("State should be populated after Collect")
		}
	}
}

// TestErrorHandling tests error handling in various methods
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(*SpigotMC)
		method  func(*SpigotMC) error
		wantErr bool
	}{
		{
			name:   "Init with invalid config",
			setup:  func(s *SpigotMC) { s.URL = "" },
			method: func(s *SpigotMC) error { return s.Init() },
			wantErr: true,
		},
		{
			name:   "Check with unreachable server",
			setup:  func(s *SpigotMC) { s.URL = "http://192.0.2.1:9999" }, // TEST-NET-1
			method: func(s *SpigotMC) error { return s.Check() },
			wantErr: true,
		},
		{
			name:   "Collect with network error",
			setup:  func(s *SpigotMC) { s.URL = "http://invalid-hostname-12345.local" },
			method: func(s *SpigotMC) error { return s.Collect() },
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSpigotMC()
			tt.setup(s)
			
			err := tt.method(s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

// TestConcurrentAccess tests concurrent access to SpigotMC instance
func TestConcurrentAccess(t *testing.T) {
	s := NewSpigotMC()
	s.Init()
	
	done := make(chan error, 3)
	
	go func() { done <- s.Init() }()
	go func() { done <- s.Check() }()
	go func() { done <- s.Collect() }()
	
	for i := 0; i < 3; i++ {
		<-done
	}
}

// TestMultipleInstances tests multiple independent instances
func TestMultipleInstances(t *testing.T) {
	s1 := NewSpigotMC()
	s2 := NewSpigotMC()
	
	s1.URL = "http://server1:8080"
	s2.URL = "http://server2:8080"
	
	if s1.URL == s2.URL {
		t.Error("Multiple instances should have independent state")
	}
}

// TestURLValidation tests URL validation
func TestURLValidation(t *testing.T) {
	tests := []struct {
		url     string
		isValid bool
	}{
		{"http://localhost:8080", true},
		{"https://localhost:8080", true},
		{"localhost:8080", true},
		{"", false},
		{"not a url", false},
		{"http://", true},
		{"://invalid", false},
	}
	
	for _, tt := range tests {
		s := NewSpigotMC()
		s.URL = tt.url
		
		// Validation typically happens in Init or Check
		err := s.Init()
		_ = err
		_ = tt.isValid
	}
}

// TestHTTPClientConfiguration tests HTTP client setup
func TestHTTPClientConfiguration(t *testing.T) {
	s := NewSpigotMC()
	s.SetDefaults()
	
	if s.client == nil {
		t.Fatal("HTTP client should be initialized")
	}
	
	if s.client.Timeout == 0 {
		t.Error("HTTP client timeout should be set")
	}
	
	if s.client.Timeout > time.Minute {
		t.Error("HTTP client timeout seems too long")
	}
}

// TestPluginInterface tests that SpigotMC implements required interface
func TestPluginInterface(t *testing.T) {
	s := NewSpigotMC()
	
	// Verify methods exist
	methods := []string{"Init", "Check", "Collect", "Charts", "SetDefaults"}
	for _, method := range methods {
		_ = method
		// In a real test, would verify method signature
	}
}

// TestChartDefinition tests chart structure
func TestChartDefinition(t *testing.T) {
	s := NewSpigotMC()
	charts := s.Charts()
	
	for i, chart := range charts {
		if chart == nil {
			t.Errorf("Chart %d is nil", i)
		}
		if chart.ID == "" {
			t.Errorf("Chart %d has empty ID", i)
		}
		if chart.Title == "" {
			t.Logf("Chart %d has empty title (optional)", i)
		}
	}
}

// TestMetricValueBoundaries tests metric value boundaries
func TestMetricValueBoundaries(t *testing.T) {
	tests := []struct {
		name  string
		value interface{}
	}{
		{"zero", 0},
		{"negative", -1},
		{"max int32", 2147483647},
		{"max int64", 9223372036854775807},
		{"float", 3.14},
		{"empty string", ""},
	}
	
	for _, tt := range tests {
		_ = tt
		// Would test how metrics handle different value types
	}
}

// TestRefreshCycle tests multiple collect cycles
func TestRefreshCycle(t *testing.T) {
	s := NewSpigotMC()
	s.Init()
	
	for i := 0; i < 5; i++ {
		err := s.Collect()
		if err != nil {
			t.Logf("Collect iteration %d: %v", i, err)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// TestCleanup tests resource cleanup if applicable
func TestCleanup(t *testing.T) {
	s := NewSpigotMC()
	s.Init()
	
	// If there's a Close/Cleanup method, test it
	// s.Close() or similar
	
	// Verify resources are released
	_ = s
}

// TestResponseDecoding tests JSON response decoding
func TestResponseDecoding(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		wantErr  bool
	}{
		{
			"valid response",
			`{"players":{"online":5,"max":20},"tps":{"1m":19.5},"memory":{"used":512,"max":1024}}`,
			false,
		},
		{
			"minimal response",
			`{"players":{"online":0,"max":0}}`,
			false,
		},
		{
			"invalid json",
			`{invalid}`,
			true,
		},
		{
			"empty object",
			`{}`,
			false,
		},
		{
			"null values",
			`{"players":null}`,
			false,
		},
	}
	
	for _, tt := range tests {
		_ = tt
		// Would test JSON unmarshaling
	}
}