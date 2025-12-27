// SPDX-License-Identifier: GPL-3.0-or-later

//go:build linux

package systemdunits

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollector_validateConfig(t *testing.T) {
	tests := map[string]struct {
		config   Config
		wantFail bool
	}{
		"success when 'include' option is set": {
			config: Config{
				Include: []string{"*.service"},
			},
			wantFail: false,
		},
		"success when 'include' option has multiple patterns": {
			config: Config{
				Include: []string{"*.service", "*.socket", "*.timer"},
			},
			wantFail: false,
		},
		"success when 'include' option has single pattern": {
			config: Config{
				Include: []string{"*"},
			},
			wantFail: false,
		},
		"fails when 'include' option is empty": {
			config: Config{
				Include: []string{},
			},
			wantFail: true,
		},
		"fails when 'include' option is nil": {
			config: Config{
				Include: nil,
			},
			wantFail: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			collr := &Collector{
				Config: test.config,
			}

			err := collr.validateConfig()

			if test.wantFail {
				assert.Error(t, err)
				assert.Equal(t, "'include' option not set", err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCollector_initUnitSelector(t *testing.T) {
	tests := map[string]struct {
		config          Config
		shouldMatch     []string
		shouldNotMatch  []string
		wantErr         bool
	}{
		"success with single wildcard pattern": {
			config: Config{
				Include: []string{"*"},
			},
			shouldMatch:    []string{"test.service", "any.socket", "foo.timer", "bar"},
			shouldNotMatch: []string{},
			wantErr:        false,
		},
		"success with service pattern": {
			config: Config{
				Include: []string{"*.service"},
			},
			shouldMatch:    []string{"test.service", "foo.service"},
			shouldNotMatch: []string{"test.socket", "foo.timer"},
			wantErr:        false,
		},
		"success with multiple patterns": {
			config: Config{
				Include: []string{"*.service", "*.socket"},
			},
			shouldMatch:    []string{"test.service", "test.socket"},
			shouldNotMatch: []string{"test.timer", "test.mount"},
			wantErr:        false,
		},
		"success with complex pattern": {
			config: Config{
				Include: []string{"ssh*", "*-dbus.service"},
			},
			shouldMatch:    []string{"sshd.service", "org-dbus.service", "test-dbus.service"},
			shouldNotMatch: []string{"test.socket", "dbus.service"},
			wantErr:        false,
		},
		"fails with empty include": {
			config: Config{
				Include: []string{},
			},
			wantErr: false, // Note: According to code, empty Include returns matcher.TRUE()
		},
		"success with question mark pattern": {
			config: Config{
				Include: []string{"test?.service"},
			},
			shouldMatch:    []string{"test1.service", "testa.service"},
			shouldNotMatch: []string{"test.service", "test12.service"},
			wantErr:        false,
		},
		"success with character class pattern": {
			config: Config{
				Include: []string{"test[0-9].service"},
			},
			shouldMatch:    []string{"test0.service", "test5.service", "test9.service"},
			shouldNotMatch: []string{"testa.service", "test.service"},
			wantErr:        false,
		},
		"success with negated pattern": {
			config: Config{
				Include: []string{"systemd-*", "!systemd-*-pre*"},
			},
			shouldMatch:    []string{"systemd-test.service"},
			shouldNotMatch: []string{},
			wantErr:        false,
		},
		"success with bracket expansion": {
			config: Config{
				Include: []string{"test{1,2}.service"},
			},
			shouldMatch:    []string{"test{1,2}.service"},
			shouldNotMatch: []string{"test1.service"},
			wantErr:        false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			collr := &Collector{
				Config: test.config,
			}

			matcher, err := collr.initUnitSelector()

			if test.wantErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, matcher)

				// Test matches
				for _, shouldMatch := range test.shouldMatch {
					ok := matcher.Match(shouldMatch)
					assert.True(t, ok, "expected '%s' to match", shouldMatch)
				}

				// Test non-matches
				for _, shouldNotMatch := range test.shouldNotMatch {
					ok := matcher.Match(shouldNotMatch)
					assert.False(t, ok, "expected '%s' to not match", shouldNotMatch)
				}
			}
		})
	}
}

func TestCollector_initUnitSelector_EmptyInclude(t *testing.T) {
	// Special test for empty include case - should return matcher.TRUE()
	collr := &Collector{
		Config: Config{
			Include: []string{},
		},
	}

	matcher, err := collr.initUnitSelector()

	require.NoError(t, err)
	assert.NotNil(t, matcher)

	// matcher.TRUE() should match everything
	testCases := []string{"test.service", "any.socket", "anything", ""}
	for _, tc := range testCases {
		assert.True(t, matcher.Match(tc), "expected empty include to match '%s'", tc)
	}
}

func TestCollector_initUnitSelector_PatternJoining(t *testing.T) {
	// Test that patterns are properly joined with spaces
	collr := &Collector{
		Config: Config{
			Include: []string{"pattern1", "pattern2", "pattern3"},
		},
	}

	matcher, err := collr.initUnitSelector()

	require.NoError(t, err)
	assert.NotNil(t, matcher)

	// Each pattern should be matchable
	testCases := []string{"pattern1", "pattern2", "pattern3"}
	for _, tc := range testCases {
		assert.True(t, matcher.Match(tc), "expected matcher to match pattern: %s", tc)
	}
}

func TestCollector_initUnitSelector_EdgeCases(t *testing.T) {
	tests := map[string]struct {
		include []string
	}{
		"single dot pattern": {
			include: []string{"."},
		},
		"pattern with spaces": {
			include: []string{"test service"},
		},
		"pattern with special characters": {
			include: []string{"@", "$", "#"},
		},
		"very long pattern": {
			include: []string{"very-long-pattern-with-many-hyphens-and-dots.service.slice"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			collr := &Collector{
				Config: Config{
					Include: test.include,
				},
			}

			matcher, err := collr.initUnitSelector()

			assert.NoError(t, err)
			assert.NotNil(t, matcher)
		})
	}
}

func TestCollector_validateConfig_WithOtherConfigFields(t *testing.T) {
	// Test validateConfig is not affected by other config fields
	collr := &Collector{
		Config: Config{
			Include:               []string{"*.service"},
			Timeout:               nil,
			SkipTransient:         true,
			CollectUnitFiles:      true,
			IncludeUnitFiles:      []string{"*.timer"},
			CollectUnitFilesEvery: nil,
		},
	}

	err := collr.validateConfig()
	assert.NoError(t, err)
}

func TestCollector_initUnitSelector_Reliability(t *testing.T) {
	// Test that initUnitSelector can be called multiple times
	collr := &Collector{
		Config: Config{
			Include: []string{"*.service"},
		},
	}

	// Call multiple times
	for i := 0; i < 5; i++ {
		matcher, err := collr.initUnitSelector()
		assert.NoError(t, err)
		assert.NotNil(t, matcher)
	}
}