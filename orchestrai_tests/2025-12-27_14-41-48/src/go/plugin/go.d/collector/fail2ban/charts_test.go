package fail2ban

import (
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/module"
)

func TestChartsInit(t *testing.T) {
	tests := []struct {
		name     string
		wantNil  bool
		validate func(t *testing.T, charts *module.Charts)
	}{
		{
			name:    "charts should not be nil",
			wantNil: false,
			validate: func(t *testing.T, charts *module.Charts) {
				if charts == nil {
					t.Error("expected charts to be not nil")
				}
			},
		},
		{
			name:    "charts should contain expected number of charts",
			wantNil: false,
			validate: func(t *testing.T, charts *module.Charts) {
				if len(*charts) == 0 {
					t.Error("expected charts to have items")
				}
			},
		},
		{
			name:    "charts should contain jails chart",
			wantNil: false,
			validate: func(t *testing.T, charts *module.Charts) {
				found := false
				for _, chart := range *charts {
					if chart.ID == "fail2ban_jails" {
						found = true
						break
					}
				}
				if !found {
					t.Error("expected jails chart to be present")
				}
			},
		},
		{
			name:    "jails chart should have banned dimension",
			wantNil: false,
			validate: func(t *testing.T, charts *module.Charts) {
				for _, chart := range *charts {
					if chart.ID == "fail2ban_jails" {
						if len(chart.Dims) == 0 {
							t.Error("expected jails chart to have dimensions")
						}
						break
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charts := Charts()
			if (charts == nil) != tt.wantNil {
				t.Errorf("Charts() returned nil=%v, want nil=%v", charts == nil, tt.wantNil)
			}
			if charts != nil {
				tt.validate(t, charts)
			}
		})
	}
}

func TestChartsStructure(t *testing.T) {
	tests := []struct {
		name     string
		validate func(t *testing.T, charts *module.Charts)
	}{
		{
			name: "each chart should have valid ID",
			validate: func(t *testing.T, charts *module.Charts) {
				for _, chart := range *charts {
					if chart.ID == "" {
						t.Error("chart ID should not be empty")
					}
				}
			},
		},
		{
			name: "each chart should have valid type",
			validate: func(t *testing.T, charts *module.Charts) {
				for _, chart := range *charts {
					if chart.Type == "" {
						t.Error("chart type should not be empty")
					}
				}
			},
		},
		{
			name: "each chart should have valid module",
			validate: func(t *testing.T, charts *module.Charts) {
				for _, chart := range *charts {
					if chart.Fqn == "" {
						t.Error("chart FQN should not be empty")
					}
				}
			},
		},
		{
			name: "dimensions should have valid names",
			validate: func(t *testing.T, charts *module.Charts) {
				for _, chart := range *charts {
					for _, dim := range chart.Dims {
						if dim.ID == "" {
							t.Error("dimension ID should not be empty")
						}
						if dim.Name == "" {
							t.Error("dimension name should not be empty")
						}
					}
				}
			},
		},
		{
			name: "charts should be properly initialized",
			validate: func(t *testing.T, charts *module.Charts) {
				for _, chart := range *charts {
					if chart.Dims == nil || len(chart.Dims) == 0 {
						t.Errorf("chart %s should have at least one dimension", chart.ID)
					}
				}
			},
		},
	}

	charts := Charts()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t, charts)
		})
	}
}

func TestChartsContent(t *testing.T) {
	charts := Charts()

	tests := []struct {
		name      string
		chartID   string
		wantFound bool
	}{
		{
			name:      "fail2ban_jails chart exists",
			chartID:   "fail2ban_jails",
			wantFound: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found := false
			for _, chart := range *charts {
				if chart.ID == tt.chartID {
					found = true
					break
				}
			}
			if found != tt.wantFound {
				t.Errorf("chart %s found=%v, want=%v", tt.chartID, found, tt.wantFound)
			}
		})
	}
}

func TestChartsDimensions(t *testing.T) {
	charts := Charts()

	tests := []struct {
		name           string
		chartID        string
		expectedDimLen int
		validate       func(t *testing.T, chart *module.Chart)
	}{
		{
			name:           "jails chart should have valid dimensions",
			chartID:        "fail2ban_jails",
			expectedDimLen: 1,
			validate: func(t *testing.T, chart *module.Chart) {
				if chart == nil {
					t.Error("expected chart to be not nil")
					return
				}
				if len(chart.Dims) != 1 {
					t.Errorf("expected 1 dimension, got %d", len(chart.Dims))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, chart := range *charts {
				if chart.ID == tt.chartID {
					tt.validate(t, chart)
					return
				}
			}
			t.Errorf("chart %s not found", tt.chartID)
		})
	}
}

func TestChartsAttributes(t *testing.T) {
	charts := Charts()

	tests := []struct {
		name     string
		validate func(t *testing.T, charts *module.Charts)
	}{
		{
			name: "all charts should have units",
			validate: func(t *testing.T, charts *module.Charts) {
				for _, chart := range *charts {
					if chart.Units == "" {
						t.Errorf("chart %s should have units", chart.ID)
					}
				}
			},
		},
		{
			name: "all dimensions should have algorithm",
			validate: func(t *testing.T, charts *module.Charts) {
				for _, chart := range *charts {
					for _, dim := range chart.Dims {
						if dim.Algo == "" {
							t.Errorf("dimension %s should have algorithm", dim.ID)
						}
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t, charts)
		})
	}
}

func TestChartsNilHandling(t *testing.T) {
	charts := Charts()

	if charts == nil {
		t.Error("Charts() should not return nil")
		return
	}

	// Verify we can iterate safely
	for _, chart := range *charts {
		if chart == nil {
			t.Error("chart should not be nil")
		}

		for _, dim := range chart.Dims {
			if dim == nil {
				t.Error("dimension should not be nil")
			}
		}
	}
}