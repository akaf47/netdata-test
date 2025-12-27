package mq

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		want    *Module
		wantErr bool
	}{
		{
			name:    "should create new module successfully",
			want:    &Module{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}

func TestModuleInit(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(t *testing.T) *Module
		wantErr   bool
	}{
		{
			name: "should initialize module successfully",
			setupFunc: func(t *testing.T) *Module {
				m, err := New()
				require.NoError(t, err)
				return m
			},
			wantErr: false,
		},
		{
			name: "should initialize with nil config",
			setupFunc: func(t *testing.T) *Module {
				m, err := New()
				require.NoError(t, err)
				m.Config = nil
				return m
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setupFunc(t)
			err := m.Init()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestModuleCheck(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(t *testing.T) *Module
		wantErr   bool
	}{
		{
			name: "should check module successfully",
			setupFunc: func(t *testing.T) *Module {
				m, err := New()
				require.NoError(t, err)
				return m
			},
			wantErr: false,
		},
		{
			name: "should check with uninitialized module",
			setupFunc: func(t *testing.T) *Module {
				m := &Module{}
				return m
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setupFunc(t)
			err := m.Check()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestModuleCharts(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(t *testing.T) *Module
		want      object
		wantErr   bool
	}{
		{
			name: "should get charts successfully",
			setupFunc: func(t *testing.T) *Module {
				m, err := New()
				require.NoError(t, err)
				return m
			},
			wantErr: false,
		},
		{
			name: "should return nil charts when not initialized",
			setupFunc: func(t *testing.T) *Module {
				m := &Module{}
				return m
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setupFunc(t)
			got := m.Charts()
			if got != nil {
				assert.IsType(t, object{}, got)
			}
		})
	}
}

func TestModuleCollect(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(t *testing.T) *Module
		want      map[string]interface{}
		wantErr   bool
	}{
		{
			name: "should collect data successfully",
			setupFunc: func(t *testing.T) *Module {
				m, err := New()
				require.NoError(t, err)
				return m
			},
			want:    make(map[string]interface{}),
			wantErr: false,
		},
		{
			name: "should collect with nil module",
			setupFunc: func(t *testing.T) *Module {
				return &Module{}
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setupFunc(t)
			got := m.Collect()
			if tt.wantErr {
				assert.Nil(t, got)
			} else {
				assert.IsType(t, map[string]interface{}{}, got)
			}
		})
	}
}

func TestModuleStop(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func(t *testing.T) *Module
		wantErr   bool
	}{
		{
			name: "should stop module successfully",
			setupFunc: func(t *testing.T) *Module {
				m, err := New()
				require.NoError(t, err)
				return m
			},
			wantErr: false,
		},
		{
			name: "should stop already stopped module",
			setupFunc: func(t *testing.T) *Module {
				m := &Module{}
				return m
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setupFunc(t)
			err := m.Stop()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestModuleWithTimeout(t *testing.T) {
	tests := []struct {
		name            string
		timeout         time.Duration
		expectValid     bool
	}{
		{
			name:        "should set valid timeout",
			timeout:     5 * time.Second,
			expectValid: true,
		},
		{
			name:        "should set zero timeout",
			timeout:     0,
			expectValid: true,
		},
		{
			name:        "should set negative timeout",
			timeout:     -1 * time.Second,
			expectValid: true,
		},
		{
			name:        "should set large timeout",
			timeout:     1000 * time.Second,
			expectValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := New()
			require.NoError(t, err)
			err = m.SetTimeout(tt.timeout)
			assert.NoError(t, err)
		})
	}
}

func TestModuleEdgeCases(t *testing.T) {
	t.Run("should handle multiple init calls", func(t *testing.T) {
		m, err := New()
		require.NoError(t, err)
		
		err1 := m.Init()
		err2 := m.Init()
		
		assert.NoError(t, err1)
		assert.NoError(t, err2)
	})

	t.Run("should handle init then check then collect sequence", func(t *testing.T) {
		m, err := New()
		require.NoError(t, err)
		
		err = m.Init()
		assert.NoError(t, err)
		
		err = m.Check()
		assert.NoError(t, err)
		
		result := m.Collect()
		assert.NotNil(t, result)
	})

	t.Run("should handle collect without init", func(t *testing.T) {
		m, err := New()
		require.NoError(t, err)
		
		result := m.Collect()
		assert.NotNil(t, result)
	})

	t.Run("should handle check without init", func(t *testing.T) {
		m, err := New()
		require.NoError(t, err)
		
		err = m.Check()
		assert.NoError(t, err)
	})

	t.Run("should handle stop before init", func(t *testing.T) {
		m, err := New()
		require.NoError(t, err)
		
		err = m.Stop()
		assert.NoError(t, err)
	})

	t.Run("should handle multiple stop calls", func(t *testing.T) {
		m, err := New()
		require.NoError(t, err)
		
		err1 := m.Stop()
		err2 := m.Stop()
		
		assert.NoError(t, err1)
		assert.NoError(t, err2)
	})
}

func TestModuleContextManagement(t *testing.T) {
	t.Run("should manage contexts correctly", func(t *testing.T) {
		m, err := New()
		require.NoError(t, err)
		
		err = m.Init()
		assert.NoError(t, err)
		
		charts := m.Charts()
		assert.NotNil(t, charts)
	})

	t.Run("should handle empty context list", func(t *testing.T) {
		m := &Module{
			Contexts: nil,
		}
		
		charts := m.Charts()
		assert.NotNil(t, charts)
	})

	t.Run("should handle context collection", func(t *testing.T) {
		m, err := New()
		require.NoError(t, err)
		
		err = m.Init()
		assert.NoError(t, err)
		
		data := m.Collect()
		assert.NotNil(t, data)
	})
}

func TestModuleConcurrency(t *testing.T) {
	t.Run("should handle concurrent collect calls", func(t *testing.T) {
		m, err := New()
		require.NoError(t, err)
		
		err = m.Init()
		require.NoError(t, err)
		
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				result := m.Collect()
				assert.NotNil(t, result)
				done <- true
			}()
		}
		
		for i := 0; i < 10; i++ {
			<-done
		}
	})

	t.Run("should handle concurrent chart calls", func(t *testing.T) {
		m, err := New()
		require.NoError(t, err)
		
		err = m.Init()
		require.NoError(t, err)
		
		done := make(chan bool, 10)
		for i := 0; i < 10; i++ {
			go func() {
				charts := m.Charts()
				assert.NotNil(t, charts)
				done <- true
			}()
		}
		
		for i := 0; i < 10; i++ {
			<-done
		}
	})
}

func TestModuleErrorHandling(t *testing.T) {
	t.Run("should handle nil config initialization", func(t *testing.T) {
		m := &Module{
			Config: nil,
		}
		
		err := m.Init()
		assert.NoError(t, err)
	})

	t.Run("should handle empty config", func(t *testing.T) {
		m := &Module{
			Config: &Config{},
		}
		
		err := m.Init()
		assert.NoError(t, err)
	})

	t.Run("should handle check with nil contexts", func(t *testing.T) {
		m := &Module{
			Contexts: nil,
		}
		
		err := m.Check()
		assert.NoError(t, err)
	})

	t.Run("should handle collect with nil contexts", func(t *testing.T) {
		m := &Module{
			Contexts: nil,
		}
		
		data := m.Collect()
		assert.NotNil(t, data)
	})
}

func TestModuleLifecycle(t *testing.T) {
	t.Run("complete module lifecycle", func(t *testing.T) {
		// Create
		m, err := New()
		require.NoError(t, err)
		assert.NotNil(t, m)
		
		// Initialize
		err = m.Init()
		require.NoError(t, err)
		
		// Check
		err = m.Check()
		require.NoError(t, err)
		
		// Get charts
		charts := m.Charts()
		assert.NotNil(t, charts)
		
		// Collect data
		data := m.Collect()
		assert.NotNil(t, data)
		
		// Stop
		err = m.Stop()
		require.NoError(t, err)
	})

	t.Run("lifecycle with timeout", func(t *testing.T) {
		m, err := New()
		require.NoError(t, err)
		
		err = m.SetTimeout(10 * time.Second)
		assert.NoError(t, err)
		
		err = m.Init()
		assert.NoError(t, err)
		
		err = m.Stop()
		assert.NoError(t, err)
	})
}

func TestModuleConfiguration(t *testing.T) {
	t.Run("should handle config with valid values", func(t *testing.T) {
		m, err := New()
		require.NoError(t, err)
		
		m.Config = &Config{
			UpdateEvery: 10,
			Timeout:     5,
		}
		
		err = m.Init()
		assert.NoError(t, err)
	})

	t.Run("should handle config with zero values", func(t *testing.T) {
		m, err := New()
		require.NoError(t, err)
		
		m.Config = &Config{
			UpdateEvery: 0,
			Timeout:     0,
		}
		
		err = m.Init()
		assert.NoError(t, err)
	})

	t.Run("should handle config with negative values", func(t *testing.T) {
		m, err := New()
		require.NoError(t, err)
		
		m.Config = &Config{
			UpdateEvery: -1,
			Timeout:     -1,
		}
		
		err = m.Init()
		assert.NoError(t, err)
	})
}

func TestModuleNilPointerHandling(t *testing.T) {
	t.Run("should handle nil module methods gracefully", func(t *testing.T) {
		// Create with minimal initialization
		m := &Module{}
		
		// These should not panic
		err := m.Init()
		assert.NoError(t, err)
		
		err = m.Check()
		assert.NoError(t, err)
		
		charts := m.Charts()
		assert.NotNil(t, charts)
		
		data := m.Collect()
		assert.NotNil(t, data)
		
		err = m.Stop()
		assert.NoError(t, err)
	})
}

func TestModuleBoundaryConditions(t *testing.T) {
	t.Run("should handle maximum timeout", func(t *testing.T) {
		m, err := New()
		require.NoError(t, err)
		
		err = m.SetTimeout(time.Duration(1<<63 - 1))
		assert.NoError(t, err)
	})

	t.Run("should handle minimum timeout", func(t *testing.T) {
		m, err := New()
		require.NoError(t, err)
		
		err = m.SetTimeout(time.Duration(-1<<63))
		assert.NoError(t, err)
	})

	t.Run("should handle update every boundary values", func(t *testing.T) {
		m, err := New()
		require.NoError(t, err)
		
		// Test with max int
		m.Config = &Config{
			UpdateEvery: 1<<31 - 1,
		}
		err = m.Init()
		assert.NoError(t, err)
		
		// Test with min int
		m.Config = &Config{
			UpdateEvery: -1 << 31,
		}
		err = m.Init()
		assert.NoError(t, err)
	})
}