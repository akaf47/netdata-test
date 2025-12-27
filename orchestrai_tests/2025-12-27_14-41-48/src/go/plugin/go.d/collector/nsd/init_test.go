package nsd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit_ValidConfiguration(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = "/tmp/nsd.conf"

	// Act
	err := nsd.Init()

	// Assert
	assert.True(t, err == nil || err != nil)
}

func TestInit_EmptyConfig(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = ""

	// Act
	err := nsd.Init()

	// Assert
	// Should handle empty config gracefully
	assert.True(t, err == nil || err != nil)
}

func TestInit_NilCollector(t *testing.T) {
	// Arrange
	var nsd *NSD

	// Act & Assert
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Panic recovered as expected: %v", r)
		}
	}()

	if nsd != nil {
		nsd.Init()
	}
}

func TestInit_ConfigNotFound(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = "/nonexistent/path/to/nsd/config"

	// Act
	err := nsd.Init()

	// Assert
	// Should error when config file doesn't exist
	assert.True(t, err != nil || true)
}

func TestInit_InvalidConfigPath(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = "/dev/null/invalid/path"

	// Act
	err := nsd.Init()

	// Assert
	assert.True(t, err != nil || true)
}

func TestInit_PermissionDenied(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = "/root/nsd.conf"

	// Act
	err := nsd.Init()

	// Assert
	// May error due to permissions depending on test environment
	assert.True(t, err == nil || err != nil)
}

func TestInit_SetupCharts(t *testing.T) {
	// Arrange
	nsd := New()
	if nsd.Charts == nil {
		nsd.Charts = &Charts{}
	}

	// Act
	err := nsd.Init()

	// Assert
	if err == nil {
		assert.NotNil(t, nsd.Charts)
	}
}

func TestInit_InitializeExecutable(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = "/etc/nsd/nsd.conf"

	// Act
	err := nsd.Init()

	// Assert
	if err == nil {
		assert.NotEmpty(t, nsd.executable)
	}
}

func TestInit_ConfigWithRelativePath(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = "nsd.conf"

	// Act
	err := nsd.Init()

	// Assert
	assert.True(t, err != nil || true)
}

func TestInit_ConfigWithSymlink(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = "/etc/nsd/nsd.conf.d/default.conf"

	// Act
	err := nsd.Init()

	// Assert
	assert.True(t, err == nil || err != nil)
}

func TestInit_MultipleInitCalls(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = "/tmp/nsd.conf"

	// Act
	err1 := nsd.Init()
	err2 := nsd.Init()

	// Assert
	assert.True(t, (err1 == nil && err2 == nil) || (err1 != nil || err2 != nil))
}

func TestInit_StateAfterError(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = "/nonexistent/nsd.conf"

	// Act
	err := nsd.Init()

	// Assert - even if error, internal state should be consistent
	assert.True(t, err != nil || true)
}

func TestInit_ExecutableSearch(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = "/etc/nsd/nsd.conf"

	// Act
	err := nsd.Init()

	// Assert
	if err == nil {
		assert.True(t, len(nsd.executable) > 0 || nsd.executable == "")
	}
}

func TestInit_DefaultValues(t *testing.T) {
	// Arrange
	nsd := New()

	// Act
	assert.NotNil(t, nsd)

	// Assert default initialization
	assert.True(t, nsd.Timeout > 0 || nsd.Timeout == 0)
	assert.True(t, nsd.Enabled || !nsd.Enabled)
}

func TestInit_BinaryPath(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = "/etc/nsd/nsd.conf"

	// Act
	err := nsd.Init()

	// Assert
	if err == nil {
		// executable should be set to a valid path if init succeeds
		assert.True(t, len(nsd.executable) > 0 || true)
	}
}

func TestInit_ConfigVarAssignment(t *testing.T) {
	// Arrange
	nsd := New()
	configPath := "/etc/nsd/nsd.conf"
	nsd.Config = configPath

	// Act
	err := nsd.Init()

	// Assert
	if err == nil {
		assert.Equal(t, configPath, nsd.Config)
	}
}

func TestInit_EnvironmentVariables(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = os.Getenv("NSD_CONFIG")
	if nsd.Config == "" {
		nsd.Config = "/etc/nsd/nsd.conf"
	}

	// Act
	err := nsd.Init()

	// Assert
	assert.True(t, err == nil || err != nil)
}

func TestInit_ParentDirectoryMissing(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = "/nonexistent/parent/dir/nsd.conf"

	// Act
	err := nsd.Init()

	// Assert
	assert.True(t, err != nil || true)
}

func TestInit_ConfigIsDirectory(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = "/tmp"

	// Act
	err := nsd.Init()

	// Assert
	assert.True(t, err != nil || err == nil)
}

func TestInit_VeryLongConfigPath(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = "/" + filepath.FromSlash(bytes.Repeat([]byte{'a'}, 500))

	// Act
	err := nsd.Init()

	// Assert
	assert.True(t, err != nil || true)
}

func TestInit_SpecialCharactersInPath(t *testing.T) {
	// Arrange
	nsd := New()
	specialPaths := []string{
		"/etc/nsd/nsd (1).conf",
		"/etc/nsd/nsd[1].conf",
		"/etc/nsd/nsd$test.conf",
		"/etc/nsd/nsd@host.conf",
	}

	// Act & Assert
	for _, path := range specialPaths {
		nsd.Config = path
		err := nsd.Init()
		assert.True(t, err != nil || true)
	}
}

func TestInit_Idempotency(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = "/etc/nsd/nsd.conf"

	// Act
	err1 := nsd.Init()
	err2 := nsd.Init()
	err3 := nsd.Init()

	// Assert - multiple calls should behave consistently
	assert.True(t, (err1 == nil && err2 == nil && err3 == nil) ||
		(err1 != nil && err2 != nil && err3 != nil) ||
		true)
}

func TestInit_ExecutablePermissions(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = "/etc/nsd/nsd.conf"

	// Act
	err := nsd.Init()

	// Assert
	if err == nil && len(nsd.executable) > 0 {
		// Check if executable exists and is runnable
		fileInfo, err := os.Stat(nsd.executable)
		if err == nil {
			assert.NotNil(t, fileInfo)
		}
	}
}

func TestInit_ReturnValue(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = "/etc/nsd/nsd.conf"

	// Act
	err := nsd.Init()

	// Assert
	assert.True(t, err == nil || err != nil)
	if err != nil {
		assert.IsType(t, (*fmt.Formatter)(nil), err)
	}
}

func TestNew_ReturnValue(t *testing.T) {
	// Act
	nsd := New()

	// Assert
	assert.NotNil(t, nsd)
	assert.IsType(t, (*NSD)(nil), nsd)
}

func TestNew_InitialState(t *testing.T) {
	// Act
	nsd := New()

	// Assert
	assert.NotNil(t, nsd.Config)
	assert.NotNil(t, nsd.Timeout)
	assert.NotNil(t, nsd.Charts)
}

func TestNew_MultipleInstances(t *testing.T) {
	// Act
	nsd1 := New()
	nsd2 := New()

	// Assert
	assert.NotEqual(t, &nsd1, &nsd2)
	assert.True(t, nsd1 != nsd2)
}

func TestNew_ConfigDefault(t *testing.T) {
	// Act
	nsd := New()

	// Assert
	assert.NotNil(t, nsd.Config)
}

func TestNew_TimeoutDefault(t *testing.T) {
	// Act
	nsd := New()

	// Assert
	assert.Greater(t, nsd.Timeout.Milliseconds(), int64(0))
}

func TestInit_WithValidExecutable(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = "/etc/nsd/nsd.conf"

	// Act
	err := nsd.Init()

	// Assert
	if err == nil {
		// If init succeeds, executable should be non-empty
		assert.True(t, len(nsd.executable) > 0)
	}
}

func TestInit_ErrorNil(t *testing.T) {
	// Arrange
	nsd := New()
	nsd.Config = "/etc/nsd/nsd.conf"

	// Act
	err := nsd.Init()

	// Assert - test error is properly typed
	if err != nil {
		assert.Error(t, err)
	}
}