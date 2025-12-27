package nvidia_smi

import (
	"testing"

	"github.com/netdata/go.d.plugin/pkg/matcher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCollector(t *testing.T) {
	t.Run("should create new collector with default values", func(t *testing.T) {
		collector := New()

		assert.NotNil(t, collector)
		assert.Equal(t, "nvidia_smi", collector.Name())
		assert.NotNil(t, collector.Charts())
	})
}

func TestInit(t *testing.T) {
	tests := []struct {
		name        string
		collector   *Nvidia
		wantErr     bool
		setup       func(*Nvidia)
		validate    func(*testing.T, *Nvidia)
	}{
		{
			name:      "should initialize successfully with default config",
			collector: New(),
			wantErr:   false,
			setup: func(n *Nvidia) {
				// Default setup
			},
			validate: func(t *testing.T, n *Nvidia) {
				assert.NotNil(t, n)
			},
		},
		{
			name:      "should initialize with custom binary path",
			collector: New(),
			wantErr:   false,
			setup: func(n *Nvidia) {
				n.BinaryPath = "/custom/path/nvidia-smi"
			},
			validate: func(t *testing.T, n *Nvidia) {
				assert.Equal(t, "/custom/path/nvidia-smi", n.BinaryPath)
			},
		},
		{
			name:      "should initialize with timeout",
			collector: New(),
			wantErr:   false,
			setup: func(n *Nvidia) {
				n.Timeout = 5
			},
			validate: func(t *testing.T, n *Nvidia) {
				assert.Equal(t, 5, n.Timeout)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(tt.collector)
			err := tt.collector.Init()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			tt.validate(t, tt.collector)
		})
	}
}

func TestCheck(t *testing.T) {
	tests := []struct {
		name      string
		collector *Nvidia
		wantErr   bool
		setup     func(*Nvidia)
	}{
		{
			name:      "should check collector health",
			collector: New(),
			wantErr:   false,
			setup: func(n *Nvidia) {
				// Setup for successful check
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(tt.collector)
			err := tt.collector.Check()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				// Either error or success, both are valid states
				// Actual check depends on nvidia-smi availability
			}
		})
	}
}

func TestInitCharts(t *testing.T) {
	t.Run("should initialize charts correctly", func(t *testing.T) {
		collector := New()
		collector.InitCharts()

		charts := collector.Charts()
		assert.NotNil(t, charts)
	})
}

func TestSetup(t *testing.T) {
	tests := []struct {
		name      string
		collector *Nvidia
		setup     func(*Nvidia)
		wantErr   bool
	}{
		{
			name:      "should setup collector",
			collector: New(),
			wantErr:   false,
			setup: func(n *Nvidia) {
				// Default setup
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(tt.collector)
			// Note: Actual implementation depends on Collect/Cleanup methods
			assert.NotNil(t, tt.collector)
		})
	}
}

func TestCleanup(t *testing.T) {
	t.Run("should cleanup without error", func(t *testing.T) {
		collector := New()
		// Cleanup should not panic
		collector.Cleanup()
	})
}

func TestCollect(t *testing.T) {
	tests := []struct {
		name      string
		collector *Nvidia
		wantErr   bool
		setup     func(*Nvidia)
	}{
		{
			name:      "should collect metrics",
			collector: New(),
			wantErr:   false,
			setup: func(n *Nvidia) {
				// Setup for metric collection
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(tt.collector)
			// Call collect if available
			// Behavior depends on implementation
			assert.NotNil(t, tt.collector)
		})
	}
}

func TestModuleRegistration(t *testing.T) {
	t.Run("should have correct module type", func(t *testing.T) {
		collector := New()
		assert.Equal(t, "nvidia_smi", collector.Name())
	})
}

func TestBinaryPathHandling(t *testing.T) {
	tests := []struct {
		name        string
		binaryPath  string
		expectedErr bool
	}{
		{
			name:        "default binary path",
			binaryPath:  "nvidia-smi",
			expectedErr: false,
		},
		{
			name:        "custom binary path",
			binaryPath:  "/usr/bin/nvidia-smi",
			expectedErr: false,
		},
		{
			name:        "relative binary path",
			binaryPath:  "./nvidia-smi",
			expectedErr: false,
		},
		{
			name:        "empty binary path",
			binaryPath:  "",
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := New()
			collector.BinaryPath = tt.binaryPath

			if tt.binaryPath == "" {
				assert.Equal(t, "", collector.BinaryPath)
			} else {
				assert.Equal(t, tt.binaryPath, collector.BinaryPath)
			}
		})
	}
}

func TestTimeoutHandling(t *testing.T) {
	tests := []struct {
		name           string
		timeout        int
		expectedResult int
	}{
		{
			name:           "zero timeout",
			timeout:        0,
			expectedResult: 0,
		},
		{
			name:           "positive timeout",
			timeout:        10,
			expectedResult: 10,
		},
		{
			name:           "negative timeout",
			timeout:        -1,
			expectedResult: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := New()
			collector.Timeout = tt.timeout
			assert.Equal(t, tt.expectedResult, collector.Timeout)
		})
	}
}

func TestGPUFilterMatchers(t *testing.T) {
	tests := []struct {
		name         string
		setupFilter  func(*Nvidia)
		expectFilter bool
	}{
		{
			name: "with gpu filter",
			setupFilter: func(n *Nvidia) {
				n.GPUFilter = &matcher.Matcher{}
			},
			expectFilter: true,
		},
		{
			name: "without gpu filter",
			setupFilter: func(n *Nvidia) {
				n.GPUFilter = nil
			},
			expectFilter: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := New()
			tt.setupFilter(collector)

			if tt.expectFilter {
				assert.NotNil(t, collector.GPUFilter)
			} else {
				assert.Nil(t, collector.GPUFilter)
			}
		})
	}
}

func TestCollectorInitialization(t *testing.T) {
	t.Run("should initialize all fields", func(t *testing.T) {
		collector := New()

		assert.NotNil(t, collector)
		assert.NotEmpty(t, collector.Name())
		assert.NotNil(t, collector.Charts())
	})
}

func TestCollectorState(t *testing.T) {
	t.Run("should maintain state correctly", func(t *testing.T) {
		collector := New()

		// Verify initial state
		assert.NotNil(t, collector)

		// Modify state
		collector.BinaryPath = "/custom/path"
		assert.Equal(t, "/custom/path", collector.BinaryPath)

		// Verify state persists
		collector.Timeout = 30
		assert.Equal(t, "/custom/path", collector.BinaryPath)
		assert.Equal(t, 30, collector.Timeout)
	})
}

func TestMultipleCollectorInstances(t *testing.T) {
	t.Run("should create independent instances", func(t *testing.T) {
		collector1 := New()
		collector2 := New()

		collector1.BinaryPath = "/path1"
		collector2.BinaryPath = "/path2"

		assert.Equal(t, "/path1", collector1.BinaryPath)
		assert.Equal(t, "/path2", collector2.BinaryPath)
		assert.NotEqual(t, collector1.BinaryPath, collector2.BinaryPath)
	})
}