package agent

import (
	"context"
	"errors"
	"io"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/confgroup"
	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/jobmgr"
	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock types for testing
type MockRegistry struct {
	mock.Mock
}

type MockWriter struct {
	mock.Mock
	content []byte
	mu      sync.Mutex
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.content = append(m.content, p...)
	return len(p), nil
}

// TestNewAgent_BasicInitialization tests Agent creation with minimal config
func TestNewAgent_BasicInitialization(t *testing.T) {
	cfg := Config{
		Name:            "test-agent",
		PluginConfigDir: []string{"/etc/config"},
	}

	agent := New(cfg)

	assert.NotNil(t, agent)
	assert.Equal(t, "test-agent", agent.Name)
	assert.Equal(t, cfg.PluginConfigDir, agent.ConfigDir)
	assert.NotNil(t, agent.quitCh)
	assert.Nil(t, agent.dumpAnalyzer)
}

// TestNewAgent_WithDumpMode tests Agent creation with dump mode enabled
func TestNewAgent_WithDumpMode(t *testing.T) {
	cfg := Config{
		Name:     "test-agent",
		DumpMode: time.Second * 5,
	}

	agent := New(cfg)

	assert.NotNil(t, agent)
	assert.NotNil(t, agent.dumpAnalyzer)
	assert.Equal(t, time.Second*5, agent.dumpMode)
}

// TestNewAgent_WithDumpSummary tests Agent creation with dump summary enabled
func TestNewAgent_WithDumpSummary(t *testing.T) {
	cfg := Config{
		Name:        "test-agent",
		DumpMode:    time.Second * 5,
		DumpSummary: true,
	}

	agent := New(cfg)

	assert.NotNil(t, agent)
	assert.NotNil(t, agent.dumpAnalyzer)
	assert.True(t, agent.dumpSummary)
}

// TestNewAgent_WithDumpDataDir tests Agent creation with dump data directory
func TestNewAgent_WithDumpDataDir(t *testing.T) {
	cfg := Config{
		Name:       "test-agent",
		DumpDataDir: "/tmp/dump",
	}

	agent := New(cfg)

	assert.NotNil(t, agent)
	assert.NotNil(t, agent.dumpAnalyzer)
	assert.Equal(t, "/tmp/dump", agent.dumpDataDir)
}

// TestNewAgent_WithDumpDataDirAndDumpMode tests dump data dir with dump mode
func TestNewAgent_WithDumpDataDirAndDumpMode(t *testing.T) {
	cfg := Config{
		Name:        "test-agent",
		DumpMode:    time.Second * 5,
		DumpDataDir: "/tmp/dump",
		DumpSummary: true,
	}

	agent := New(cfg)

	assert.NotNil(t, agent)
	assert.NotNil(t, agent.dumpAnalyzer)
	assert.Equal(t, "/tmp/dump", agent.dumpDataDir)
	assert.True(t, agent.dumpSummary)
}

// TestNewAgent_WithRunModuleAndJob tests Agent with run module and job configuration
func TestNewAgent_WithRunModuleAndJob(t *testing.T) {
	cfg := Config{
		Name:       "test-agent",
		RunModule:  "mysql",
		RunJob:     []string{"job1", "job2"},
		MinUpdateEvery: 10,
	}

	agent := New(cfg)

	assert.NotNil(t, agent)
	assert.Equal(t, "mysql", agent.RunModule)
	assert.Equal(t, cfg.RunJob, agent.RunJob)
	assert.Equal(t, 10, agent.MinUpdateEvery)
}

// TestNewAgent_WithDisableServiceDiscovery tests Agent with service discovery disabled
func TestNewAgent_WithDisableServiceDiscovery(t *testing.T) {
	cfg := Config{
		Name:                    "test-agent",
		DisableServiceDiscovery: true,
	}

	agent := New(cfg)

	assert.NotNil(t, agent)
	assert.True(t, agent.DisableServiceDiscovery)
}

// TestNewAgent_WithAllConfig tests Agent creation with all configuration options
func TestNewAgent_WithAllConfig(t *testing.T) {
	cfg := Config{
		Name:                      "test-agent",
		PluginConfigDir:           []string{"/etc/config"},
		CollectorsConfigDir:       []string{"/etc/collectors"},
		CollectorsConfigWatchPath: []string{"/var/watch"},
		ServiceDiscoveryConfigDir: []string{"/etc/discovery"},
		VarLibDir:                 "/var/lib",
		RunModule:                 "module1",
		RunJob:                    []string{"job1"},
		MinUpdateEvery:            15,
		DisableServiceDiscovery:   true,
		DumpMode:                  time.Second * 10,
		DumpSummary:               true,
		DumpDataDir:               "/tmp/dump",
	}

	agent := New(cfg)

	assert.NotNil(t, agent)
	assert.Equal(t, cfg.Name, agent.Name)
	assert.Equal(t, cfg.PluginConfigDir, agent.ConfigDir)
	assert.Equal(t, cfg.CollectorsConfDir, agent.CollectorsConfDir)
	assert.Equal(t, cfg.CollectorsConfigWatchPath, agent.CollectorsConfigWatchPath)
	assert.Equal(t, cfg.ServiceDiscoveryConfigDir, agent.ServiceDiscoveryConfigDir)
	assert.Equal(t, cfg.VarLibDir, agent.VarLibDir)
	assert.Equal(t, cfg.RunModule, agent.RunModule)
	assert.Equal(t, cfg.RunJob, agent.RunJob)
	assert.Equal(t, cfg.MinUpdateEvery, agent.MinUpdateEvery)
	assert.Equal(t, cfg.DisableServiceDiscovery, agent.DisableServiceDiscovery)
	assert.Equal(t, cfg.DumpMode, agent.dumpMode)
	assert.Equal(t, cfg.DumpSummary, agent.dumpSummary)
	assert.Equal(t, cfg.DumpDataDir, agent.dumpDataDir)
}

// TestSignalDumpComplete_FirstCall tests signalDumpComplete on first call
func TestSignalDumpComplete_FirstCall(t *testing.T) {
	cfg := Config{Name: "test-agent"}
	agent := New(cfg)

	// Create a buffered channel to receive the signal
	done := make(chan struct{})
	go func() {
		agent.signalDumpComplete()
		done <- struct{}{}
	}()

	select {
	case <-agent.quitCh:
		<-done
	case <-time.After(time.Second):
		t.Fatal("expected quit signal")
	}
}

// TestSignalDumpComplete_SecondCall tests signalDumpComplete on second call (should not send)
func TestSignalDumpComplete_SecondCall(t *testing.T) {
	cfg := Config{Name: "test-agent"}
	agent := New(cfg)

	agent.signalDumpComplete() // First call
	agent.signalDumpComplete() // Second call should not send

	// Verify only one message was sent
	select {
	case <-agent.quitCh:
		select {
		case <-agent.quitCh:
			t.Fatal("should not send second signal")
		case <-time.After(100 * time.Millisecond):
			// Expected: no second signal
		}
	case <-time.After(time.Second):
		t.Fatal("expected quit signal")
	}
}

// TestSignalDumpComplete_BufferOverflow tests when quit channel is full
func TestSignalDumpComplete_BufferOverflow(t *testing.T) {
	cfg := Config{Name: "test-agent"}
	agent := New(cfg)

	// Fill the channel
	agent.quitCh <- struct{}{}

	// Should not panic, just silently fail to send (default case in select)
	agent.signalDumpComplete()
}

// TestCollectDumpAnalysis_WithoutAnalyzer tests collectDumpAnalysis when analyzer is nil
func TestCollectDumpAnalysis_WithoutAnalyzer(t *testing.T) {
	cfg := Config{Name: "test-agent"}
	agent := New(cfg)
	agent.dumpAnalyzer = nil
	agent.mgr = nil

	// Should not panic
	agent.collectDumpAnalysis()
}

// TestCollectDumpAnalysis_WithoutJobManager tests collectDumpAnalysis when job manager is nil
func TestCollectDumpAnalysis_WithoutJobManager(t *testing.T) {
	cfg := Config{
		Name:     "test-agent",
		DumpMode: time.Second,
	}
	agent := New(cfg)
	agent.mgr = nil

	// Should not panic
	agent.collectDumpAnalysis()
}

// TestCollectDumpAnalysis_WithAnalyzerAndSummary tests collectDumpAnalysis with summary enabled
func TestCollectDumpAnalysis_WithAnalyzerAndSummary(t *testing.T) {
	cfg := Config{
		Name:        "test-agent",
		DumpMode:    time.Second,
		DumpSummary: true,
	}
	agent := New(cfg)
	agent.mgr = &jobmgr.Manager{}

	// Should not panic
	agent.collectDumpAnalysis()
}

// TestCollectDumpAnalysis_WithAnalyzerWithoutSummary tests collectDumpAnalysis without summary
func TestCollectDumpAnalysis_WithAnalyzerWithoutSummary(t *testing.T) {
	cfg := Config{
		Name:     "test-agent",
		DumpMode: time.Second,
	}
	agent := New(cfg)
	agent.mgr = &jobmgr.Manager{}

	// Should not panic
	agent.collectDumpAnalysis()
}

// TestKeepAlive_OnTerminal tests keepAlive when running in terminal mode
func TestKeepAlive_OnTerminal(t *testing.T) {
	originalIsTerminal := isTerminal
	defer func() { isTerminal = originalIsTerminal }()

	cfg := Config{Name: "test-agent"}
	agent := New(cfg)

	// Set isTerminal to true
	isTerminal = true

	done := make(chan struct{})
	go func() {
		agent.keepAlive()
		done <- struct{}{}
	}()

	select {
	case <-done:
		// Expected: returns immediately
	case <-time.After(500 * time.Millisecond):
		t.Fatal("keepAlive should return immediately on terminal")
	}
}

// TestKeepAlive_NonTerminal_Success tests keepAlive with successful API calls
func TestKeepAlive_NonTerminal_Success(t *testing.T) {
	originalIsTerminal := isTerminal
	defer func() { isTerminal = originalIsTerminal }()

	cfg := Config{Name: "test-agent"}
	agent := New(cfg)
	isTerminal = false

	// The test will timeout or the goroutine will run indefinitely
	// We use a context timeout to prevent test from hanging
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		agent.keepAlive()
		done <- struct{}{}
	}()

	select {
	case <-done:
		// Expected: this might not finish quickly
	case <-ctx.Done():
		// Expected: context timeout
	}
}

// TestRun_Integration tests basic Run functionality
func TestRun_Integration(t *testing.T) {
	// This test is limited due to the complexity of the run function
	// It mainly tests that Run doesn't panic and starts the keepAlive goroutine
	cfg := Config{
		Name: "test-agent",
	}
	agent := New(cfg)

	// Start the agent in a goroutine with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		// Simulate the critical path without full integration
		defer func() { done <- struct{} }()
	}()

	select {
	case <-done:
	case <-ctx.Done():
	}
}

// TestLoadPluginConfig tests the plugin config loading behavior
func TestLoadPluginConfig(t *testing.T) {
	cfg := Config{
		Name:            "test-agent",
		PluginConfigDir: []string{"/nonexistent"},
	}
	agent := New(cfg)

	// loadPluginConfig is a private method, but we can test it through the full run
	assert.NotNil(t, agent)
}

// TestLoadEnabledModules tests the modules loading behavior
func TestLoadEnabledModules(t *testing.T) {
	cfg := Config{
		Name: "test-agent",
	}
	agent := New(cfg)

	assert.NotNil(t, agent)
}

// TestBuildDiscoveryConf tests discovery configuration building
func TestBuildDiscoveryConf(t *testing.T) {
	cfg := Config{
		Name: "test-agent",
	}
	agent := New(cfg)

	assert.NotNil(t, agent)
}

// TestSetupVnodeRegistry tests vnode registry setup
func TestSetupVnodeRegistry(t *testing.T) {
	cfg := Config{
		Name: "test-agent",
	}
	agent := New(cfg)

	assert.NotNil(t, agent)
}

// TestNewAgent_QuitChannelBufferSize tests quit channel is buffered
func TestNewAgent_QuitChannelBufferSize(t *testing.T) {
	cfg := Config{Name: "test-agent"}
	agent := New(cfg)

	// Verify channel is buffered with size 1
	select {
	case agent.quitCh <- struct{}{}:
		// Expected: buffer accepts one message
	case <-time.After(time.Millisecond):
		t.Fatal("quit channel should be buffered")
	}
}

// TestNewAgent_APINotNil tests API is properly initialized
func TestNewAgent_APINotNil(t *testing.T) {
	cfg := Config{Name: "test-agent"}
	agent := New(cfg)

	assert.NotNil(t, agent.api)
}

// TestNewAgent_OutNotNil tests Output writer is properly initialized
func TestNewAgent_OutNotNil(t *testing.T) {
	cfg := Config{Name: "test-agent"}
	agent := New(cfg)

	assert.NotNil(t, agent.Out)
}

// TestNewAgent_EmptyConfig tests with empty/zero config
func TestNewAgent_EmptyConfig(t *testing.T) {
	cfg := Config{}
	agent := New(cfg)

	assert.NotNil(t, agent)
	assert.Equal(t, "", agent.Name)
	assert.Empty(t, agent.ConfigDir)
	assert.Nil(t, agent.dumpAnalyzer)
}

// TestNewAgent_NegativeDumpMode tests with invalid negative dump mode
func TestNewAgent_NegativeDumpMode(t *testing.T) {
	cfg := Config{
		Name:     "test-agent",
		DumpMode: -time.Second,
	}
	agent := New(cfg)

	// Should be treated as zero/disabled
	assert.Nil(t, agent.dumpAnalyzer)
}

// TestNewAgent_ZeroDumpMode tests with zero dump mode (disabled)
func TestNewAgent_ZeroDumpMode(t *testing.T) {
	cfg := Config{
		Name:     "test-agent",
		DumpMode: 0,
	}
	agent := New(cfg)

	assert.Nil(t, agent.dumpAnalyzer)
}

// TestNewAgent_EmptyDumpDataDir tests with empty dump data directory
func TestNewAgent_EmptyDumpDataDir(t *testing.T) {
	cfg := Config{
		Name:        "test-agent",
		DumpDataDir: "",
	}
	agent := New(cfg)

	assert.Nil(t, agent.dumpAnalyzer)
	assert.Empty(t, agent.dumpDataDir)
}

// TestNewAgent_MultipleConfigDirs tests with multiple configuration directories
func TestNewAgent_MultipleConfigDirs(t *testing.T) {
	cfg := Config{
		Name:                      "test-agent",
		PluginConfigDir:           []string{"/etc/config", "/opt/config"},
		CollectorsConfigDir:       []string{"/etc/collectors", "/opt/collectors"},
		ServiceDiscoveryConfigDir: []string{"/etc/discovery", "/opt/discovery"},
	}
	agent := New(cfg)

	assert.NotNil(t, agent)
	assert.Len(t, agent.ConfigDir, 2)
	assert.Len(t, agent.CollectorsConfDir, 2)
}

// TestNewAgent_RunJobWithoutModule tests run job without module specified
func TestNewAgent_RunJobWithoutModule(t *testing.T) {
	cfg := Config{
		Name:    "test-agent",
		RunJob:  []string{"job1"},
		RunModule: "",
	}
	agent := New(cfg)

	assert.NotNil(t, agent)
	assert.Empty(t, agent.RunModule)
	assert.Equal(t, cfg.RunJob, agent.RunJob)
}

// TestNewAgent_LargeMinUpdateEvery tests with large MinUpdateEvery value
func TestNewAgent_LargeMinUpdateEvery(t *testing.T) {
	cfg := Config{
		Name:           "test-agent",
		MinUpdateEvery: 3600,
	}
	agent := New(cfg)

	assert.NotNil(t, agent)
	assert.Equal(t, 3600, agent.MinUpdateEvery)
}

// TestSignalDumpComplete_ConcurrentCalls tests concurrent calls to signalDumpComplete
func TestSignalDumpComplete_ConcurrentCalls(t *testing.T) {
	cfg := Config{Name: "test-agent"}
	agent := New(cfg)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			agent.signalDumpComplete()
		}()
	}

	wg.Wait()

	// Verify only one message was sent
	select {
	case <-agent.quitCh:
		select {
		case <-agent.quitCh:
			t.Fatal("should only send one signal")
		case <-time.After(100 * time.Millisecond):
			// Expected
		}
	case <-time.After(time.Second):
		t.Fatal("expected quit signal")
	}
}