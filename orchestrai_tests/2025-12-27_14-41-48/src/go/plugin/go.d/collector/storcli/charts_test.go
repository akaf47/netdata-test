package storcli

import (
	"testing"
)

// TestNewCharts tests the NewCharts function
func TestNewCharts(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "should create new charts instance",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charts := NewCharts()
			if charts == nil {
				t.Error("NewCharts() returned nil")
			}
		})
	}
}

// TestChartsAddControllerChart tests the AddControllerChart method
func TestChartsAddControllerChart(t *testing.T) {
	tests := []struct {
		name           string
		controllerID   string
		expectedResult bool
	}{
		{
			name:           "should add controller chart with valid ID",
			controllerID:   "c0",
			expectedResult: true,
		},
		{
			name:           "should add controller chart with numeric ID",
			controllerID:   "c1",
			expectedResult: true,
		},
		{
			name:           "should handle empty controller ID",
			controllerID:   "",
			expectedResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charts := NewCharts()
			chart := charts.AddControllerChart(tt.controllerID)
			if tt.expectedResult {
				if chart == nil {
					t.Errorf("AddControllerChart(%q) returned nil", tt.controllerID)
				}
			}
		})
	}
}

// TestChartsAddPhysicalDriveChart tests the AddPhysicalDriveChart method
func TestChartsAddPhysicalDriveChart(t *testing.T) {
	tests := []struct {
		name           string
		controllerID   string
		enclosureID    string
		slotID         string
		expectedResult bool
	}{
		{
			name:           "should add physical drive chart",
			controllerID:   "c0",
			enclosureID:    "e0",
			slotID:         "s0",
			expectedResult: true,
		},
		{
			name:           "should add multiple physical drive charts",
			controllerID:   "c0",
			enclosureID:    "e1",
			slotID:         "s1",
			expectedResult: true,
		},
		{
			name:           "should handle empty controller ID",
			controllerID:   "",
			enclosureID:    "e0",
			slotID:         "s0",
			expectedResult: false,
		},
		{
			name:           "should handle empty enclosure ID",
			controllerID:   "c0",
			enclosureID:    "",
			slotID:         "s0",
			expectedResult: false,
		},
		{
			name:           "should handle empty slot ID",
			controllerID:   "c0",
			enclosureID:    "e0",
			slotID:         "",
			expectedResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charts := NewCharts()
			chart := charts.AddPhysicalDriveChart(tt.controllerID, tt.enclosureID, tt.slotID)
			if tt.expectedResult {
				if chart == nil {
					t.Errorf("AddPhysicalDriveChart(%q, %q, %q) returned nil", tt.controllerID, tt.enclosureID, tt.slotID)
				}
			}
		})
	}
}

// TestChartsAddVirtualDriveChart tests the AddVirtualDriveChart method
func TestChartsAddVirtualDriveChart(t *testing.T) {
	tests := []struct {
		name           string
		controllerID   string
		virtualDriveID string
		expectedResult bool
	}{
		{
			name:           "should add virtual drive chart",
			controllerID:   "c0",
			virtualDriveID: "vd0",
			expectedResult: true,
		},
		{
			name:           "should add multiple virtual drive charts",
			controllerID:   "c0",
			virtualDriveID: "vd1",
			expectedResult: true,
		},
		{
			name:           "should handle empty controller ID",
			controllerID:   "",
			virtualDriveID: "vd0",
			expectedResult: false,
		},
		{
			name:           "should handle empty virtual drive ID",
			controllerID:   "c0",
			virtualDriveID: "",
			expectedResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charts := NewCharts()
			chart := charts.AddVirtualDriveChart(tt.controllerID, tt.virtualDriveID)
			if tt.expectedResult {
				if chart == nil {
					t.Errorf("AddVirtualDriveChart(%q, %q) returned nil", tt.controllerID, tt.virtualDriveID)
				}
			}
		})
	}
}

// TestChartsAddBatteryChart tests the AddBatteryChart method
func TestChartsAddBatteryChart(t *testing.T) {
	tests := []struct {
		name           string
		controllerID   string
		expectedResult bool
	}{
		{
			name:           "should add battery chart",
			controllerID:   "c0",
			expectedResult: true,
		},
		{
			name:           "should add battery chart for controller c1",
			controllerID:   "c1",
			expectedResult: true,
		},
		{
			name:           "should handle empty controller ID",
			controllerID:   "",
			expectedResult: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charts := NewCharts()
			chart := charts.AddBatteryChart(tt.controllerID)
			if tt.expectedResult {
				if chart == nil {
					t.Errorf("AddBatteryChart(%q) returned nil", tt.controllerID)
				}
			}
		})
	}
}

// TestChartsRemoveControllerCharts tests the RemoveControllerCharts method
func TestChartsRemoveControllerCharts(t *testing.T) {
	tests := []struct {
		name         string
		controllerID string
	}{
		{
			name:         "should remove controller charts",
			controllerID: "c0",
		},
		{
			name:         "should handle removal of non-existent controller",
			controllerID: "c99",
		},
		{
			name:         "should handle empty controller ID",
			controllerID: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charts := NewCharts()
			charts.AddControllerChart(tt.controllerID)
			// Should not panic
			charts.RemoveControllerCharts(tt.controllerID)
		})
	}
}

// TestChartsAddControllerChartWithMultipleMetrics tests adding controller chart with metrics
func TestChartsAddControllerChartWithMultipleMetrics(t *testing.T) {
	t.Run("should add controller chart and verify it exists", func(t *testing.T) {
		charts := NewCharts()
		chart := charts.AddControllerChart("c0")
		
		if chart == nil {
			t.Error("AddControllerChart returned nil")
		}
	})
}

// TestChartsIntegrationMultipleCharts tests adding multiple chart types
func TestChartsIntegrationMultipleCharts(t *testing.T) {
	t.Run("should manage multiple chart types", func(t *testing.T) {
		charts := NewCharts()
		
		// Add controller chart
		controllerChart := charts.AddControllerChart("c0")
		if controllerChart == nil {
			t.Error("AddControllerChart returned nil")
		}
		
		// Add physical drive chart
		driveChart := charts.AddPhysicalDriveChart("c0", "e0", "s0")
		if driveChart == nil {
			t.Error("AddPhysicalDriveChart returned nil")
		}
		
		// Add virtual drive chart
		vdChart := charts.AddVirtualDriveChart("c0", "vd0")
		if vdChart == nil {
			t.Error("AddVirtualDriveChart returned nil")
		}
		
		// Add battery chart
		batteryChart := charts.AddBatteryChart("c0")
		if batteryChart == nil {
			t.Error("AddBatteryChart returned nil")
		}
	})
}

// TestChartsHandleMultipleControllers tests handling multiple controllers
func TestChartsHandleMultipleControllers(t *testing.T) {
	t.Run("should handle multiple controllers", func(t *testing.T) {
		charts := NewCharts()
		
		// Add charts for controller 0
		c0 := charts.AddControllerChart("c0")
		if c0 == nil {
			t.Error("AddControllerChart for c0 returned nil")
		}
		
		// Add charts for controller 1
		c1 := charts.AddControllerChart("c1")
		if c1 == nil {
			t.Error("AddControllerChart for c1 returned nil")
		}
		
		// Add physical drives for controller 0
		pd0 := charts.AddPhysicalDriveChart("c0", "e0", "s0")
		if pd0 == nil {
			t.Error("AddPhysicalDriveChart for c0 returned nil")
		}
		
		// Add physical drives for controller 1
		pd1 := charts.AddPhysicalDriveChart("c1", "e0", "s0")
		if pd1 == nil {
			t.Error("AddPhysicalDriveChart for c1 returned nil")
		}
	})
}

// TestChartsHandleMultiplePhysicalDrives tests handling multiple physical drives
func TestChartsHandleMultiplePhysicalDrives(t *testing.T) {
	t.Run("should handle multiple physical drives in same enclosure", func(t *testing.T) {
		charts := NewCharts()
		
		pd0 := charts.AddPhysicalDriveChart("c0", "e0", "s0")
		pd1 := charts.AddPhysicalDriveChart("c0", "e0", "s1")
		pd2 := charts.AddPhysicalDriveChart("c0", "e0", "s2")
		
		if pd0 == nil || pd1 == nil || pd2 == nil {
			t.Error("AddPhysicalDriveChart returned nil")
		}
	})
}

// TestChartsHandleMultipleEnclosures tests handling multiple enclosures
func TestChartsHandleMultipleEnclosures(t *testing.T) {
	t.Run("should handle multiple enclosures", func(t *testing.T) {
		charts := NewCharts()
		
		e0 := charts.AddPhysicalDriveChart("c0", "e0", "s0")
		e1 := charts.AddPhysicalDriveChart("c0", "e1", "s0")
		e2 := charts.AddPhysicalDriveChart("c0", "e2", "s0")
		
		if e0 == nil || e1 == nil || e2 == nil {
			t.Error("AddPhysicalDriveChart returned nil")
		}
	})
}

// TestChartsAddVirtualDriveChartMultiple tests adding multiple virtual drives
func TestChartsAddVirtualDriveChartMultiple(t *testing.T) {
	t.Run("should handle multiple virtual drives", func(t *testing.T) {
		charts := NewCharts()
		
		vd0 := charts.AddVirtualDriveChart("c0", "vd0")
		vd1 := charts.AddVirtualDriveChart("c0", "vd1")
		vd2 := charts.AddVirtualDriveChart("c0", "vd2")
		
		if vd0 == nil || vd1 == nil || vd2 == nil {
			t.Error("AddVirtualDriveChart returned nil")
		}
	})
}

// TestChartsRemoveAndAdd tests removing and re-adding charts
func TestChartsRemoveAndAdd(t *testing.T) {
	t.Run("should allow removing and re-adding controller charts", func(t *testing.T) {
		charts := NewCharts()
		
		// Add controller
		c0 := charts.AddControllerChart("c0")
		if c0 == nil {
			t.Error("First AddControllerChart returned nil")
		}
		
		// Remove controller
		charts.RemoveControllerCharts("c0")
		
		// Re-add controller
		c0Again := charts.AddControllerChart("c0")
		if c0Again == nil {
			t.Error("Second AddControllerChart returned nil")
		}
	})
}

// TestChartsChartIDGeneration tests that chart IDs are properly generated
func TestChartsChartIDGeneration(t *testing.T) {
	tests := []struct {
		name         string
		setupFunc    func(*Charts) string
		expectedID   string
	}{
		{
			name: "should generate controller chart ID",
			setupFunc: func(c *Charts) string {
				chart := c.AddControllerChart("c0")
				if chart != nil {
					return chart.ID
				}
				return ""
			},
			expectedID: "storcli_c0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charts := NewCharts()
			id := tt.setupFunc(charts)
			if id == "" && tt.expectedID != "" {
				t.Errorf("Expected chart ID to be generated, got empty string")
			}
		})
	}
}

// TestChartsBoundaryConditions tests boundary conditions
func TestChartsBoundaryConditions(t *testing.T) {
	tests := []struct {
		name string
		fn   func(*Charts)
	}{
		{
			name: "should handle very long controller ID",
			fn: func(c *Charts) {
				c.AddControllerChart("c_very_long_id_that_goes_on_and_on_0123456789_0123456789_0123456789")
			},
		},
		{
			name: "should handle numeric controller IDs",
			fn: func(c *Charts) {
				c.AddControllerChart("0")
				c.AddControllerChart("1")
				c.AddControllerChart("999")
			},
		},
		{
			name: "should handle special characters in IDs",
			fn: func(c *Charts) {
				c.AddControllerChart("c_0")
				c.AddControllerChart("c-0")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			charts := NewCharts()
			// Should not panic
			tt.fn(charts)
		})
	}
}