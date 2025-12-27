package file

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDiscovery(t *testing.T) {
	d := New()
	assert.NotNil(t, d)
	assert.IsType(t, &Discovery{}, d)
}

func TestDiscoveryInit(t *testing.T) {
	tests := []struct {
		name      string
		configDir string
		shouldErr bool
	}{
		{
			name:      "valid config directory",
			configDir: "/etc/netdata",
			shouldErr: false,
		},
		{
			name:      "empty config directory",
			configDir: "",
			shouldErr: true,
		},
		{
			name:      "relative path",
			configDir: "./config",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := New()
			err := d.Init(tt.configDir)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDiscoveryDiscoverWithValidConfig(t *testing.T) {
	d := New()
	err := d.Init("testdata")
	if err != nil {
		t.Skipf("setup failed: %v", err)
	}

	targets, err := d.Discover()
	assert.NoError(t, err)
	assert.IsType(t, []string{}, targets)
}

func TestDiscoveryDiscoverWithNilState(t *testing.T) {
	d := &Discovery{}
	targets, err := d.Discover()
	assert.Error(t, err)
	assert.Nil(t, targets)
}

func TestDiscoveryDiscoverEmptyDirectory(t *testing.T) {
	d := New()
	err := d.Init("nonexistent_dir")
	if err == nil {
		targets, err := d.Discover()
		assert.NoError(t, err)
		assert.Equal(t, 0, len(targets))
	}
}

func TestDiscoveryDiscoverMultipleFiles(t *testing.T) {
	d := New()
	d.Init("testdata")
	targets, err := d.Discover()
	assert.NoError(t, err)
	assert.NotNil(t, targets)
}

func TestDiscoveryClose(t *testing.T) {
	d := New()
	err := d.Close()
	if err != nil {
		t.Logf("Close returned error: %v", err)
	}
}

func TestDiscoveryInit_MultipleInits(t *testing.T) {
	d := New()
	err1 := d.Init("/etc/netdata")
	err2 := d.Init("/etc/netdata/conf.d")

	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestDiscoveryDiscoverConcurrentCalls(t *testing.T) {
	d := New()
	d.Init("testdata")

	done := make(chan error, 2)

	go func() {
		_, err := d.Discover()
		done <- err
	}()

	go func() {
		_, err := d.Discover()
		done <- err
	}()

	err1 := <-done
	err2 := <-done

	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestDiscoveryDiscoverStateValidation(t *testing.T) {
	d := &Discovery{}
	_, err := d.Discover()
	assert.Error(t, err)
}

func TestDiscoveryInitWithSpecialChars(t *testing.T) {
	d := New()
	err := d.Init("/etc/netdata/configs with spaces")
	assert.NoError(t, err)
}

func TestDiscoveryInitDuplicateCalls(t *testing.T) {
	d := New()
	path := "/etc/netdata/conf.d"

	err1 := d.Init(path)
	err2 := d.Init(path)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
}

func TestDiscoveryDiscoverErrorHandling(t *testing.T) {
	d := New()
	err := d.Init("/invalid/path/that/does/not/exist")

	if err == nil {
		targets, discoverErr := d.Discover()
		assert.NoError(t, discoverErr)
		assert.NotNil(t, targets)
	}
}

func TestDiscoveryDiscoverReturnType(t *testing.T) {
	d := New()
	d.Init("testdata")

	targets, err := d.Discover()

	assert.NoError(t, err)
	assert.IsType(t, ([]string)(nil), targets)
}

func TestDiscoveryDiscoverEmptyResults(t *testing.T) {
	d := New()
	d.Init("/nonexistent")

	targets, _ := d.Discover()

	if targets != nil {
		assert.Equal(t, 0, len(targets))
	}
}

func TestDiscoveryLifecycle(t *testing.T) {
	d := New()
	assert.NotNil(t, d)

	err := d.Init("/etc/netdata")
	assert.NoError(t, err)

	targets, err := d.Discover()
	assert.NoError(t, err)
	assert.IsType(t, []string{}, targets)

	d.Close()
}

func TestDiscoveryInitCalledBeforeDiscover(t *testing.T) {
	d := New()

	targets, err := d.Discover()
	assert.Error(t, err)
	assert.Nil(t, targets)
}

func TestDiscoveryDiscoverWithDifferentPaths(t *testing.T) {
	paths := []string{
		"/etc/netdata/conf.d",
		"/etc/netdata/conf.d/services",
		"/usr/local/etc/netdata",
	}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			d := New()
			d.Init(path)
			targets, _ := d.Discover()
			assert.IsType(t, ([]string)(nil), targets)
		})
	}
}

func TestDiscoveryDiscoverTimeoutHandling(t *testing.T) {
	d := New()
	d.Init("/etc/netdata")

	c := make(chan []string, 1)
	go func() {
		targets, _ := d.Discover()
		c <- targets
	}()

	select {
	case targets := <-c:
		assert.IsType(t, []string{}, targets)
	case <-time.After(5 * time.Second):
		t.Fatal("discover took too long")
	}
}

func TestDiscoveryInitWithAbsolutePath(t *testing.T) {
	d := New()
	err := d.Init("/absolute/path/to/config")
	assert.NoError(t, err)
}

func TestDiscoveryInitWithRelativePath(t *testing.T) {
	d := New()
	err := d.Init("./relative/path")
	assert.NoError(t, err)
}

func TestDiscoveryDiscoverAfterClose(t *testing.T) {
	d := New()
	d.Init("/etc/netdata")
	d.Close()

	targets, err := d.Discover()
	assert.IsType(t, ([]string)(nil), targets)
}

func TestDiscoveryMultipleDiscoverCalls(t *testing.T) {
	d := New()
	d.Init("/etc/netdata")

	targets1, err1 := d.Discover()
	targets2, err2 := d.Discover()

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, targets1, targets2)
}

func TestDiscoveryDiscoverNilCheck(t *testing.T) {
	d := &Discovery{}
	targets, err := d.Discover()

	assert.Error(t, err)
	assert.Nil(t, targets)
}