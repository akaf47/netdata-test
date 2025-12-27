// SPDX-License-Identifier: GPL-3.0-or-later

package file

import (
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/confgroup"
	"github.com/stretchr/testify/assert"
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

// TestValidateConfigWithValidRegistryAndRead tests validation succeeds with valid registry and Read paths
func TestValidateConfigWithValidRegistryAndRead(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Read:     []string{"/etc/config/file.conf"},
		Watch:    []string{},
	}

	err := validateConfig(cfg)

	assert.NoError(t, err)
}

// TestValidateConfigWithValidRegistryAndWatch tests validation succeeds with valid registry and Watch paths
func TestValidateConfigWithValidRegistryAndWatch(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Read:     []string{},
		Watch:    []string{"/etc/config/file.conf"},
	}

	err := validateConfig(cfg)

	assert.NoError(t, err)
}

// TestValidateConfigWithValidRegistryAndBothReadAndWatch tests validation with both Read and Watch
func TestValidateConfigWithValidRegistryAndBothReadAndWatch(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Read:     []string{"/etc/config/read.conf"},
		Watch:    []string{"/etc/config/watch.conf"},
	}

	err := validateConfig(cfg)

	assert.NoError(t, err)
}

// TestValidateConfigWithEmptyRegistry tests validation fails with empty registry
func TestValidateConfigWithEmptyRegistry(t *testing.T) {
	cfg := Config{
		Registry: MockRegistry{},
		Read:     []string{"/etc/config/file.conf"},
		Watch:    []string{},
	}

	err := validateConfig(cfg)

	assert.Error(t, err)
	assert.Equal(t, "empty config registry", err.Error())
}

// TestValidateConfigWithNilRegistry tests validation fails with nil registry
func TestValidateConfigWithNilRegistry(t *testing.T) {
	cfg := Config{
		Registry: nil,
		Read:     []string{"/etc/config/file.conf"},
		Watch:    []string{},
	}

	err := validateConfig(cfg)

	assert.Error(t, err)
}

// TestValidateConfigWithEmptyReadAndWatch tests validation fails with no discoverers
func TestValidateConfigWithEmptyReadAndWatch(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Read:     []string{},
		Watch:    []string{},
	}

	err := validateConfig(cfg)

	assert.Error(t, err)
	assert.Equal(t, "discoverers not set", err.Error())
}

// TestValidateConfigWithMultipleReadPaths tests validation with multiple Read paths
func TestValidateConfigWithMultipleReadPaths(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Read:     []string{"/etc/config/file1.conf", "/etc/config/file2.conf", "/etc/config/file3.conf"},
		Watch:    []string{},
	}

	err := validateConfig(cfg)

	assert.NoError(t, err)
}

// TestValidateConfigWithMultipleWatchPaths tests validation with multiple Watch paths
func TestValidateConfigWithMultipleWatchPaths(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Read:     []string{},
		Watch:    []string{"/etc/config/watch1.conf", "/etc/config/watch2.conf", "/etc/config/watch3.conf"},
	}

	err := validateConfig(cfg)

	assert.NoError(t, err)
}

// TestValidateConfigWithLargeRegistry tests validation with many registry entries
func TestValidateConfigWithLargeRegistry(t *testing.T) {
	registry := make(MockRegistry)
	for i := 0; i < 1000; i++ {
		registry["service_"+string(rune(i))] = &confgroup.ServiceConfig{}
	}

	cfg := Config{
		Registry: registry,
		Read:     []string{"/etc/config/file.conf"},
		Watch:    []string{},
	}

	err := validateConfig(cfg)

	assert.NoError(t, err)
}

// TestValidateConfigEmptyPathStrings tests with empty path strings
func TestValidateConfigEmptyPathStrings(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Read:     []string{""},
		Watch:    []string{},
	}

	err := validateConfig(cfg)

	assert.NoError(t, err)
}

// TestValidateConfigWithSpacePathStrings tests with whitespace path strings
func TestValidateConfigWithSpacePathStrings(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Read:     []string{"   "},
		Watch:    []string{},
	}

	err := validateConfig(cfg)

	assert.NoError(t, err)
}

// TestValidateConfigBothEmptyAndNilRegistry tests edge case with both nil and empty checks
func TestValidateConfigBothEmptyAndNilRegistry(t *testing.T) {
	// When len(cfg.Registry) == 0, should return error before checking discoverers
	cfg := Config{
		Registry: MockRegistry{},
		Read:     []string{},
		Watch:    []string{},
	}

	err := validateConfig(cfg)

	assert.Error(t, err)
	assert.Equal(t, "empty config registry", err.Error())
}

// TestValidateConfigOrderOfValidation tests registry check happens first
func TestValidateConfigOrderOfValidation(t *testing.T) {
	// When both registry is empty and discoverers are not set, registry error should be returned first
	cfg := Config{
		Registry: MockRegistry{},
		Read:     []string{},
		Watch:    []string{},
	}

	err := validateConfig(cfg)

	assert.Error(t, err)
	assert.Equal(t, "empty config registry", err.Error())
}

// TestValidateConfigWithSingleReadPath tests single Read path
func TestValidateConfigWithSingleReadPath(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Read:     []string{"/etc/config/file.conf"},
		Watch:    []string{},
	}

	err := validateConfig(cfg)

	assert.NoError(t, err)
}

// TestValidateConfigWithSingleWatchPath tests single Watch path
func TestValidateConfigWithSingleWatchPath(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Read:     []string{},
		Watch:    []string{"/etc/config/file.conf"},
	}

	err := validateConfig(cfg)

	assert.NoError(t, err)
}

// TestValidateConfigWithValidRegistryButEmptyDiscoverers tests discoverer validation
func TestValidateConfigWithValidRegistryButEmptyDiscoverers(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
		"service2": &confgroup.ServiceConfig{},
		"service3": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Read:     []string{},
		Watch:    []string{},
	}

	err := validateConfig(cfg)

	assert.Error(t, err)
	assert.Equal(t, "discoverers not set", err.Error())
}

// TestValidateConfigWithOneReadAndMultipleWatch tests mixed path types
func TestValidateConfigWithOneReadAndMultipleWatch(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Read:     []string{"/etc/config/read.conf"},
		Watch:    []string{"/etc/config/watch1.conf", "/etc/config/watch2.conf"},
	}

	err := validateConfig(cfg)

	assert.NoError(t, err)
}

// TestValidateConfigWithMultipleReadAndOneWatch tests mixed path types
func TestValidateConfigWithMultipleReadAndOneWatch(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Read:     []string{"/etc/config/read1.conf", "/etc/config/read2.conf"},
		Watch:    []string{"/etc/config/watch.conf"},
	}

	err := validateConfig(cfg)

	assert.NoError(t, err)
}

// TestValidateConfigPathsWithSpecialCharacters tests paths with special characters
func TestValidateConfigPathsWithSpecialCharacters(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Read:     []string{"/etc/config/file-name_2.conf", "/etc/config/file@test.conf"},
		Watch:    []string{},
	}

	err := validateConfig(cfg)

	assert.NoError(t, err)
}

// TestValidateConfigWithUnicodePathStrings tests paths with unicode characters
func TestValidateConfigWithUnicodePathStrings(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	cfg := Config{
		Registry: registry,
		Read:     []string{"/etc/config/файл.conf"},
		Watch:    []string{},
	}

	err := validateConfig(cfg)

	assert.NoError(t, err)
}

// TestValidateConfigLengthOfReadArray tests len(cfg.Read) check correctly
func TestValidateConfigLengthOfReadArray(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	// Exactly one item in Read should pass
	cfg := Config{
		Registry: registry,
		Read:     []string{"path"},
		Watch:    []string{},
	}

	err := validateConfig(cfg)

	assert.NoError(t, err)
}

// TestValidateConfigLengthOfWatchArray tests len(cfg.Watch) check correctly
func TestValidateConfigLengthOfWatchArray(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	// Exactly one item in Watch should pass
	cfg := Config{
		Registry: registry,
		Read:     []string{},
		Watch:    []string{"path"},
	}

	err := validateConfig(cfg)

	assert.NoError(t, err)
}

// TestValidateConfigBoundaryConditionReadAndWatch tests sum of Read and Watch
func TestValidateConfigBoundaryConditionReadAndWatch(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	// len(cfg.Read)+len(cfg.Watch) = 2, should pass
	cfg := Config{
		Registry: registry,
		Read:     []string{"path1"},
		Watch:    []string{"path2"},
	}

	err := validateConfig(cfg)

	assert.NoError(t, err)
}

// TestValidateConfigNilSlices tests with nil slices instead of empty
func TestValidateConfigNilSlices(t *testing.T) {
	registry := MockRegistry{
		"service1": &confgroup.ServiceConfig{},
	}
	// nil slices have len() == 0
	cfg := Config{
		Registry: registry,
		Read:     nil,
		Watch:    nil,
	}

	err := validateConfig(cfg)

	assert.Error(t, err)
	assert.Equal(t, "discoverers not set", err.Error())
}