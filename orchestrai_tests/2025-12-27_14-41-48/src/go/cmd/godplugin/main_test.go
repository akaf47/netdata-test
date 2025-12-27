package main

import (
	"os"
	"testing"

	"github.com/netdata/netdata/go/plugins/pkg/cli"
	"github.com/netdata/netdata/go/plugins/pkg/executable"
	"github.com/netdata/netdata/go/plugins/pkg/buildinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	t.Run("should unset TZ when it starts with colon", func(t *testing.T) {
		// Arrange
		os.Setenv("TZ", ":UTC")
		
		// Act
		init()
		
		// Assert
		val, exists := os.LookupEnv("TZ")
		assert.False(t, exists, "TZ should be unset")
		assert.Equal(t, "", val)
		
		// Cleanup
		os.Unsetenv("TZ")
	})

	t.Run("should not unset TZ when it does not start with colon", func(t *testing.T) {
		// Arrange
		os.Setenv("TZ", "America/New_York")
		
		// Act
		init()
		
		// Assert
		val, exists := os.LookupEnv("TZ")
		assert.True(t, exists, "TZ should still be set")
		assert.Equal(t, "America/New_York", val)
		
		// Cleanup
		os.Unsetenv("TZ")
	})

	t.Run("should handle empty TZ", func(t *testing.T) {
		// Arrange
		os.Unsetenv("TZ")
		
		// Act
		init()
		
		// Assert
		_, exists := os.LookupEnv("TZ")
		assert.False(t, exists)
	})

	t.Run("should handle TZ with multiple colons", func(t *testing.T) {
		// Arrange
		os.Setenv("TZ", ":Europe/London:extra")
		
		// Act
		init()
		
		// Assert
		_, exists := os.LookupEnv("TZ")
		assert.False(t, exists)
		
		// Cleanup
		os.Unsetenv("TZ")
	})
}

func TestParseCLI(t *testing.T) {
	t.Run("should parse valid CLI args", func(t *testing.T) {
		// Arrange
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"godplugin", "--version"}
		
		// Act
		opts := parseCLI()
		
		// Assert
		assert.NotNil(t, opts)
		assert.IsType(t, &cli.Option{}, opts)
	})

	t.Run("should return cli.Option struct", func(t *testing.T) {
		// Arrange
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"godplugin"}
		
		// Act
		opts := parseCLI()
		
		// Assert
		require.NotNil(t, opts)
		assert.IsType(t, &cli.Option{}, opts)
	})

	t.Run("should handle help flag", func(t *testing.T) {
		// Arrange
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"godplugin", "-h"}
		
		// Act & Assert
		// parseCLI should handle help internally and exit
		// We cannot directly test os.Exit behavior, but we test the logic path
		opts := parseCLI()
		assert.NotNil(t, opts)
	})

	t.Run("should handle empty arguments", func(t *testing.T) {
		// Arrange
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"godplugin"}
		
		// Act
		opts := parseCLI()
		
		// Assert
		assert.NotNil(t, opts)
	})
}

func TestParseCliEdgeCases(t *testing.T) {
	t.Run("should handle debug flag", func(t *testing.T) {
		// Arrange
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"godplugin", "--debug"}
		
		// Act
		opts := parseCLI()
		
		// Assert
		require.NotNil(t, opts)
		assert.Equal(t, true, opts.Debug)
	})

	t.Run("should handle module flag", func(t *testing.T) {
		// Arrange
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"godplugin", "--module", "test_module"}
		
		// Act
		opts := parseCLI()
		
		// Assert
		require.NotNil(t, opts)
		assert.Equal(t, "test_module", opts.Module)
	})

	t.Run("should handle job flag", func(t *testing.T) {
		// Arrange
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"godplugin", "--job", "test_job"}
		
		// Act
		opts := parseCLI()
		
		// Assert
		require.NotNil(t, opts)
		assert.Equal(t, "test_job", opts.Job)
	})

	t.Run("should handle update-every flag", func(t *testing.T) {
		// Arrange
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"godplugin", "--update-every", "30"}
		
		// Act
		opts := parseCLI()
		
		// Assert
		require.NotNil(t, opts)
		assert.Equal(t, int32(30), opts.UpdateEvery)
	})

	t.Run("should handle dump-summary flag", func(t *testing.T) {
		// Arrange
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"godplugin", "--dump-summary"}
		
		// Act
		opts := parseCLI()
		
		// Assert
		require.NotNil(t, opts)
	})

	t.Run("should handle version flag", func(t *testing.T) {
		// Arrange
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"godplugin", "--version"}
		
		// Act
		opts := parseCLI()
		
		// Assert
		require.NotNil(t, opts)
		assert.Equal(t, true, opts.Version)
	})

	t.Run("should handle multiple flags", func(t *testing.T) {
		// Arrange
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"godplugin", "--debug", "--module", "test", "--job", "job1"}
		
		// Act
		opts := parseCLI()
		
		// Assert
		require.NotNil(t, opts)
		assert.Equal(t, true, opts.Debug)
		assert.Equal(t, "test", opts.Module)
		assert.Equal(t, "job1", opts.Job)
	})
}

// Unit tests for core logic paths

func TestMainLogic(t *testing.T) {
	t.Run("should handle version output", func(t *testing.T) {
		// This tests the version branch in main
		// We verify that when opts.Version is true, the function should return early
		oldArgs := os.Args
		oldStdout := os.Stdout
		defer func() {
			os.Args = oldArgs
			os.Stdout = oldStdout
		}()
		os.Args = []string{"godplugin", "--version"}
		
		opts := parseCLI()
		assert.True(t, opts.Version)
	})
}

func TestInitTZVariations(t *testing.T) {
	tests := []struct {
		name        string
		tzValue     string
		shouldUnset bool
	}{
		{
			name:        "TZ with leading colon should be unset",
			tzValue:     ":UTC",
			shouldUnset: true,
		},
		{
			name:        "TZ without leading colon should not be unset",
			tzValue:     "UTC",
			shouldUnset: false,
		},
		{
			name:        "TZ with colon in middle should not be unset",
			tzValue:     "America:New_York",
			shouldUnset: false,
		},
		{
			name:        "TZ with path-like value starting with colon should be unset",
			tzValue:     ":/usr/share/zoneinfo/UTC",
			shouldUnset: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			os.Setenv("TZ", tt.tzValue)
			
			// Act
			init()
			
			// Assert
			_, exists := os.LookupEnv("TZ")
			if tt.shouldUnset {
				assert.False(t, exists, "TZ should be unset")
			} else {
				assert.True(t, exists, "TZ should still be set")
			}
			
			// Cleanup
			os.Unsetenv("TZ")
		})
	}
}

func TestParseCliReturnType(t *testing.T) {
	t.Run("parseCLI always returns non-nil cli.Option", func(t *testing.T) {
		// Arrange
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"godplugin"}
		
		// Act
		opts := parseCLI()
		
		// Assert
		assert.NotNil(t, opts)
		assert.IsType(t, (*cli.Option)(nil), opts)
	})
}

// Test boundary conditions and edge cases
func TestParseCliWithNoArgs(t *testing.T) {
	t.Run("should handle program name only", func(t *testing.T) {
		// Arrange
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"godplugin"}
		
		// Act
		opts := parseCLI()
		
		// Assert
		assert.NotNil(t, opts)
		assert.False(t, opts.Version)
		assert.False(t, opts.Debug)
	})
}

func TestInitBehavior(t *testing.T) {
	t.Run("init function should be called once on package load", func(t *testing.T) {
		// Arrange
		os.Unsetenv("TZ")
		os.Setenv("TZ", ":TEST")
		
		// Act
		init()
		
		// Assert
		_, exists := os.LookupEnv("TZ")
		assert.False(t, exists)
		
		// Cleanup
		os.Unsetenv("TZ")
	})
}

func TestParseCliIntegration(t *testing.T) {
	t.Run("should parse all supported flags", func(t *testing.T) {
		// Arrange
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{
			"godplugin",
			"--version",
			"--debug",
			"--module", "module1",
			"--job", "job1",
			"--update-every", "60",
		}
		
		// Act
		opts := parseCLI()
		
		// Assert
		require.NotNil(t, opts)
		assert.True(t, opts.Version)
		assert.True(t, opts.Debug)
		assert.Equal(t, "module1", opts.Module)
		assert.Equal(t, "job1", opts.Job)
		assert.Equal(t, int32(60), opts.UpdateEvery)
	})
}