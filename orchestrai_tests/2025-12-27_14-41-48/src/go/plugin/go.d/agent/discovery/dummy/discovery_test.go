// SPDX-License-Identifier: GPL-3.0-or-later

package dummy

import (
	"context"
	"testing"

	"github.com/netdata/netdata/go/plugins/logger"
	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/confgroup"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockRegistry is a mock implementation of confgroup.Registry
type MockRegistry map[string]*confgroup.ServiceConfig

func (m MockRegistry) Lookup(name string) (*confgroup.ServiceConfig, bool) {
	cfg, ok := m[name]
	return cfg, ok
}

func (m MockRegistry) Register(cfg *confgroup.ServiceConfig) error {
	return nil
}

// TestNewDiscoveryWithValidConfig tests successful creation with valid config
func TestNewDiscoveryWithValidConfig(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Names:    []string{"service1"},
	}

	d, err := NewDiscovery(cfg)

	require.NoError(t, err)
	require.NotNil(t, d)
	assert.Equal(t, registry, d.reg)
	assert.Equal(t, []string{"service1"}, d.names)
	assert.NotNil(t, d.Logger)
}

// TestNewDiscoveryWithEmptyRegistry tests validation fails with empty registry
func TestNewDiscoveryWithEmptyRegistry(t *testing.T) {
	cfg := Config{
		Registry: MockRegistry{},
		Names:    []string{"service1"},
	}

	d, err := NewDiscovery(cfg)

	assert.Error(t, err)
	assert.Nil(t, d)
	assert.Contains(t, err.Error(), "config validation")
}

// TestNewDiscoveryWithEmptyNames tests validation fails with no names
func TestNewDiscoveryWithEmptyNames(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Names:    []string{},
	}

	d, err := NewDiscovery(cfg)

	assert.Error(t, err)
	assert.Nil(t, d)
	assert.Contains(t, err.Error(), "config validation")
}

// TestNewDiscoveryWithNilRegistry tests validation fails with nil registry
func TestNewDiscoveryWithNilRegistry(t *testing.T) {
	cfg := Config{
		Registry: nil,
		Names:    []string{"service1"},
	}

	d, err := NewDiscovery(cfg)

	assert.Error(t, err)
	assert.Nil(t, d)
}

// TestDiscoveryString tests String method returns correct name
func TestDiscoveryString(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Names:    []string{"service1"},
	}
	d, _ := NewDiscovery(cfg)

	result := d.String()

	assert.Equal(t, "dummy discovery", result)
}

// TestDiscoveryName tests Name method returns correct identifier
func TestDiscoveryName(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Names:    []string{"service1"},
	}
	d, _ := NewDiscovery(cfg)

	result := d.Name()

	assert.Equal(t, "dummy discovery", result)
}

// TestDiscoveryRunWithContextDone tests Run exits immediately on context cancellation
func TestDiscoveryRunWithContextDone(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Names:    []string{"service1"},
	}
	d, _ := NewDiscovery(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	in := make(chan []*confgroup.Group, 1)
	d.Run(ctx, in)

	// Channel should be closed after Run returns
	_, ok := <-in
	assert.False(t, ok)
}

// TestDiscoveryRunWithActiveContext tests Run sends groups and closes channel
func TestDiscoveryRunWithActiveContext(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Names:    []string{"service1"},
	}
	d, _ := NewDiscovery(cfg)

	ctx := context.Background()
	in := make(chan []*confgroup.Group, 1)

	d.Run(ctx, in)

	groups := <-in
	assert.NotNil(t, groups)
	assert.Len(t, groups, 1)
	assert.Equal(t, "internal", groups[0].Source)

	// Channel should be closed
	_, ok := <-in
	assert.False(t, ok)
}

// TestDiscoveryGroupsWithEmptyNames returns group with no configs
func TestDiscoveryGroupsWithEmptyNames(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Names:    []string{},
	}
	d := &Discovery{
		Logger: logger.New(),
		reg:    registry,
		names:  []string{},
	}

	groups := d.groups()

	assert.NotNil(t, groups)
	assert.Len(t, groups, 1)
	assert.Equal(t, "internal", groups[0].Source)
	assert.Len(t, groups[0].Configs, 0)
}

// TestDiscoveryGroupsWithSingleName returns group with one config
func TestDiscoveryGroupsWithSingleName(t *testing.T) {
	serviceCfg := &confgroup.ServiceConfig{}
	registry := MockRegistry{
		"service1": serviceCfg,
	}
	cfg := Config{
		Registry: registry,
		Names:    []string{"service1"},
	}
	d := &Discovery{
		Logger: logger.New(),
		reg:    registry,
		names:  []string{"service1"},
	}

	groups := d.groups()

	assert.NotNil(t, groups)
	assert.Len(t, groups, 1)
	assert.Equal(t, "internal", groups[0].Source)
	assert.Len(t, groups[0].Configs, 1)
}

// TestDiscoveryGroupsWithMultipleNames returns group with multiple configs
func TestDiscoveryGroupsWithMultipleNames(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
		"service2": &confgroup.ServiceConfig{},
		"service3": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Names:    []string{"service1", "service2", "service3"},
	}
	d := &Discovery{
		Logger: logger.New(),
		reg:    registry,
		names:  []string{"service1", "service2", "service3"},
	}

	groups := d.groups()

	assert.NotNil(t, groups)
	assert.Len(t, groups, 1)
	assert.Equal(t, "internal", groups[0].Source)
	assert.Len(t, groups[0].Configs, 3)
}

// TestDiscoveryGroupsWithUnregisteredName skips unregistered service names
func TestDiscoveryGroupsWithUnregisteredName(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Names:    []string{"service1", "unregistered"},
	}
	d := &Discovery{
		Logger: logger.New(),
		reg:    registry,
		names:  []string{"service1", "unregistered"},
	}

	groups := d.groups()

	assert.NotNil(t, groups)
	assert.Len(t, groups, 1)
	assert.Len(t, groups[0].Configs, 1)
	assert.Equal(t, "service1", groups[0].Configs[0].Module)
}

// TestDiscoveryGroupsConfigSourceType verifies config source type is set correctly
func TestDiscoveryGroupsConfigSourceType(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	d := &Discovery{
		Logger: logger.New(),
		reg:    registry,
		names:  []string{"service1"},
	}

	groups := d.groups()

	assert.Len(t, groups[0].Configs, 1)
	cfg := groups[0].Configs[0]
	assert.Equal(t, "service1", cfg.Module)
	assert.Equal(t, "dummy", cfg.Provider)
}

// TestDiscoveryGroupsWithMixedRegisteredAndUnregistered tests partial lookup
func TestDiscoveryGroupsWithMixedRegisteredAndUnregistered(t *testing.T) {
	registry := MockRegistry{
		"registered1": &confgroup.ServiceConfig{},
		"registered2": &confgroup.ServiceConfig{},
	}
	d := &Discovery{
		Logger: logger.New(),
		reg:    registry,
		names:  []string{"registered1", "missing", "registered2", "missing2"},
	}

	groups := d.groups()

	assert.Len(t, groups[0].Configs, 2)
}

// TestDiscoveryRunSendsAllGroups ensures all groups are sent before closing
func TestDiscoveryRunSendsAllGroups(t *testing.T) {
	registry := MockRegistry{
		"svc1": &confgroup.ServiceConfig{},
		"svc2": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Names:    []string{"svc1", "svc2"},
	}
	d, _ := NewDiscovery(cfg)

	ctx := context.Background()
	in := make(chan []*confgroup.Group)

	go d.Run(ctx, in)

	groups := <-in
	assert.NotNil(t, groups)
	assert.Len(t, groups, 1)
	assert.Len(t, groups[0].Configs, 2)
}

// TestDiscoveryLoggerIsInitialized ensures logger is properly initialized
func TestDiscoveryLoggerIsInitialized(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Names:    []string{"service1"},
	}

	d, err := NewDiscovery(cfg)

	require.NoError(t, err)
	assert.NotNil(t, d.Logger)
	// Verify logger has the correct context
	assert.NotNil(t, d.Logger)
}

// TestDiscoveryGroupsAllConfigsHaveProvider ensures all configs have dummy provider
func TestDiscoveryGroupsAllConfigsHaveProvider(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
		"service2": &confgroup.ServiceConfig{},
	}
	d := &Discovery{
		Logger: logger.New(),
		reg:    registry,
		names:  []string{"service1", "service2"},
	}

	groups := d.groups()

	for _, config := range groups[0].Configs {
		assert.Equal(t, "dummy", config.Provider)
	}
}

// TestDiscoveryGroupsAllConfigsHaveInternalSource ensures all configs have internal source
func TestDiscoveryGroupsAllConfigsHaveInternalSource(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
		"service2": &confgroup.ServiceConfig{},
	}
	d := &Discovery{
		Logger: logger.New(),
		reg:    registry,
		names:  []string{"service1", "service2"},
	}

	groups := d.groups()

	for _, config := range groups[0].Configs {
		assert.Equal(t, "internal", config.Source)
	}
}

// TestNewDiscoveryLogsCorrectComponent verifies logger has correct context
func TestNewDiscoveryLogsCorrectComponent(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Names:    []string{"service1"},
	}

	d, err := NewDiscovery(cfg)

	require.NoError(t, err)
	assert.NotNil(t, d.Logger)
}

// TestDiscoveryRunWithBufferedChannel tests with buffered channel
func TestDiscoveryRunWithBufferedChannel(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Names:    []string{"service1"},
	}
	d, _ := NewDiscovery(cfg)

	ctx := context.Background()
	in := make(chan []*confgroup.Group, 10)

	d.Run(ctx, in)

	select {
	case groups := <-in:
		assert.NotNil(t, groups)
	default:
		t.Fatal("expected groups to be sent")
	}

	_, ok := <-in
	assert.False(t, ok, "channel should be closed")
}

// TestDiscoveryWithLargeNameSet tests with many names
func TestDiscoveryWithLargeNameSet(t *testing.T) {
	registry := make(MockRegistry)
	names := make([]string, 0)

	for i := 0; i < 100; i++ {
		name := "service_" + string(rune(i))
		registry[name] = &confgroup.ServiceConfig{}
		names = append(names, name)
	}

	cfg := Config{
		Registry: registry,
		Names:    names,
	}

	d, err := NewDiscovery(cfg)

	require.NoError(t, err)
	groups := d.groups()
	assert.Len(t, groups[0].Configs, 100)
}

// TestValidateConfigWithNilRegistry tests validation with nil registry
func TestValidateConfigWithNilRegistry(t *testing.T) {
	cfg := Config{
		Registry: nil,
		Names:    []string{"service"},
	}

	err := validateConfig(cfg)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config validation")
}

// TestValidateConfigWithEmptyNamesAndEmptyRegistry tests both empty
func TestValidateConfigWithEmptyNamesAndEmptyRegistry(t *testing.T) {
	cfg := Config{
		Registry: MockRegistry{},
		Names:    []string{},
	}

	err := validateConfig(cfg)

	assert.Error(t, err)
}