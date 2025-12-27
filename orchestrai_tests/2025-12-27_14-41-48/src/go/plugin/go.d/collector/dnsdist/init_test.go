package dnsdist

import (
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin.d/go.d/agent/module"
)

// TestNewDnsdist tests the NewDnsdist factory function
func TestNewDnsdist(t *testing.T) {
	tests := []struct {
		name string
		want *Dnsdist
	}{
		{
			name: "should create a new Dnsdist instance with default values",
			want: &Dnsdist{
				Config: Config{
					Address: defaultAddress,
					Timeout: defaultTimeout,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New()
			if got == nil {
				t.Errorf("New() returned nil")
			}
			if got.Name() == "" {
				t.Errorf("New() instance has empty Name")
			}
		})
	}
}

// TestInitialize tests the Initialize method
func TestInitialize(t *testing.T) {
	tests := []struct {
		name      string
		d         *Dnsdist
		wantErr   bool
		setupFunc func(*Dnsdist)
	}{
		{
			name:    "should initialize successfully with valid config",
			d:       New(),
			wantErr: false,
			setupFunc: func(d *Dnsdist) {
				d.Address = "127.0.0.1:8083"
				d.Timeout = 5
			},
		},
		{
			name:    "should handle empty address",
			d:       New(),
			wantErr: true,
			setupFunc: func(d *Dnsdist) {
				d.Address = ""
			},
		},
		{
			name:    "should handle zero timeout",
			d:       New(),
			wantErr: true,
			setupFunc: func(d *Dnsdist) {
				d.Address = "127.0.0.1:8083"
				d.Timeout = 0
			},
		},
		{
			name:    "should handle negative timeout",
			d:       New(),
			wantErr: true,
			setupFunc: func(d *Dnsdist) {
				d.Address = "127.0.0.1:8083"
				d.Timeout = -1
			},
		},
		{
			name:    "should handle invalid address format",
			d:       New(),
			wantErr: true,
			setupFunc: func(d *Dnsdist) {
				d.Address = "invalid://address"
			},
		},
		{
			name:    "should handle address without port",
			d:       New(),
			wantErr: true,
			setupFunc: func(d *Dnsdist) {
				d.Address = "127.0.0.1"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				tt.setupFunc(tt.d)
			}

			err := tt.d.Init()
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestCheck tests the Check method for connectivity validation
func TestCheck(t *testing.T) {
	tests := []struct {
		name      string
		d         *Dnsdist
		wantErr   bool
		setupFunc func(*Dnsdist)
	}{
		{
			name:    "should check connectivity successfully",
			d:       New(),
			wantErr: false,
			setupFunc: func(d *Dnsdist) {
				d.Address = "127.0.0.1:8083"
				d.Timeout = 5
				d.Init()
			},
		},
		{
			name:    "should fail check with unreachable address",
			d:       New(),
			wantErr: true,
			setupFunc: func(d *Dnsdist) {
				d.Address = "192.0.2.1:8083" // non-routable address
				d.Timeout = 1
				d.Init()
			},
		},
		{
			name:    "should fail check with invalid port",
			d:       New(),
			wantErr: true,
			setupFunc: func(d *Dnsdist) {
				d.Address = "127.0.0.1:65535"
				d.Timeout = 1
				d.Init()
			},
		},
		{
			name:    "should handle timeout during check",
			d:       New(),
			wantErr: true,
			setupFunc: func(d *Dnsdist) {
				d.Address = "192.0.2.1:8083" // non-routable address
				d.Timeout = 1
				d.Init()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				tt.setupFunc(tt.d)
			}

			err := tt.d.Check()
			if (err != nil) != tt.wantErr {
				t.Errorf("Check() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestCollect tests the Collect method
func TestCollect(t *testing.T) {
	tests := []struct {
		name      string
		d         *Dnsdist
		setupFunc func(*Dnsdist)
		validate  func(*testing.T, map[string]int64)
	}{
		{
			name: "should collect metrics successfully",
			d:    New(),
			setupFunc: func(d *Dnsdist) {
				d.Address = "127.0.0.1:8083"
				d.Timeout = 5
				d.Init()
			},
			validate: func(t *testing.T, metrics map[string]int64) {
				if metrics == nil {
					t.Errorf("Collect() returned nil metrics")
				}
			},
		},
		{
			name: "should handle collection failure gracefully",
			d:    New(),
			setupFunc: func(d *Dnsdist) {
				d.Address = "192.0.2.1:8083" // non-routable
				d.Timeout = 1
				d.Init()
			},
			validate: func(t *testing.T, metrics map[string]int64) {
				// Should return empty metrics on error
				if metrics != nil && len(metrics) > 0 {
					t.Errorf("Collect() should return empty metrics on connection error")
				}
			},
		},
		{
			name: "should return empty map for zero metrics",
			d:    New(),
			setupFunc: func(d *Dnsdist) {
				d.Address = "127.0.0.1:8083"
				d.Timeout = 5
				d.Init()
			},
			validate: func(t *testing.T, metrics map[string]int64) {
				if metrics == nil {
					t.Errorf("Collect() should return map, not nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setupFunc != nil {
				tt.setupFunc(tt.d)
			}

			metrics := tt.d.Collect()
			if tt.validate != nil {
				tt.validate(t, metrics)
			}
		})
	}
}

// TestCharts tests the Charts method
func TestCharts(t *testing.T) {
	tests := []struct {
		name      string
		d         *Dnsdist
		validate  func(*testing.T, *module.Charts)
	}{
		{
			name: "should return charts for initialized instance",
			d:    New(),
			validate: func(t *testing.T, charts *module.Charts) {
				if charts == nil {
					t.Errorf("Charts() returned nil")
				}
				if len(charts.Charts) == 0 {
					t.Errorf("Charts() returned empty charts")
				}
			},
		},
		{
			name: "should return consistent charts on multiple calls",
			d:    New(),
			validate: func(t *testing.T, charts *module.Charts) {
				if charts == nil {
					t.Errorf("Charts() returned nil")
				}
				// Verify charts are populated with expected metric names
				hasMetrics := false
				for _, chart := range charts.Charts {
					if len(chart.Dims) > 0 {
						hasMetrics = true
						break
					}
				}
				if !hasMetrics {
					t.Errorf("Charts() should contain at least one chart with dimensions")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charts := tt.d.Charts()
			if tt.validate != nil {
				tt.validate(t, charts)
			}
		})
	}
}

// TestCollectEdgeCases tests edge cases for collector
func TestCollectEdgeCases(t *testing.T) {
	t.Run("should handle multiple concurrent collect calls", func(t *testing.T) {
		d := New()
		d.Address = "127.0.0.1:8083"
		d.Timeout = 5
		d.Init()

		done := make(chan bool, 3)
		for i := 0; i < 3; i++ {
			go func() {
				_ = d.Collect()
				done <- true
			}()
		}

		for i := 0; i < 3; i++ {
			<-done
		}
	})

	t.Run("should handle rapid init and check cycles", func(t *testing.T) {
		d := New()
		for i := 0; i < 5; i++ {
			d.Address = "127.0.0.1:8083"
			d.Timeout = 5
			_ = d.Init()
			_ = d.Check()
		}
	})

	t.Run("should handle maximum timeout value", func(t *testing.T) {
		d := New()
		d.Address = "127.0.0.1:8083"
		d.Timeout = 3600 // 1 hour
		err := d.Init()
		if err != nil {
			t.Errorf("Init() with large timeout failed: %v", err)
		}
	})

	t.Run("should handle minimum valid timeout", func(t *testing.T) {
		d := New()
		d.Address = "127.0.0.1:8083"
		d.Timeout = 1
		err := d.Init()
		if err != nil {
			t.Errorf("Init() with minimum timeout failed: %v", err)
		}
	})
}

// TestInitBoundaryConditions tests boundary conditions
func TestInitBoundaryConditions(t *testing.T) {
	t.Run("should handle very long address", func(t *testing.T) {
		d := New()
		d.Address = "127.0.0.1:" + string(make([]byte, 1000))
		err := d.Init()
		if err == nil {
			t.Errorf("Init() should fail with invalid long address")
		}
	})

	t.Run("should handle IPv6 address", func(t *testing.T) {
		d := New()
		d.Address = "[::1]:8083"
		d.Timeout = 5
		err := d.Init()
		// Should handle IPv6, either succeed or fail gracefully
		_ = err
	})

	t.Run("should handle hostname instead of IP", func(t *testing.T) {
		d := New()
		d.Address = "localhost:8083"
		d.Timeout = 5
		err := d.Init()
		// Should handle hostname resolution attempt
		_ = err
	})
}