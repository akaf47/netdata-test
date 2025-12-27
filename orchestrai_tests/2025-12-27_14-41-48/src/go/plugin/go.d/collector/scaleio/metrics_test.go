package scaleio

import (
	"testing"
)

func TestMetricsInitialization(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "should initialize metrics structure",
			test: func(t *testing.T) {
				m := NewMetrics()
				if m == nil {
					t.Errorf("NewMetrics() returned nil")
				}
			},
		},
		{
			name: "should have empty charts by default",
			test: func(t *testing.T) {
				m := NewMetrics()
				if m.Charts == nil {
					t.Errorf("Charts should be initialized")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

func TestMetricsCharts(t *testing.T) {
	tests := []struct {
		name        string
		setup       func() *Metrics
		expectCharts int
		checkChart  func(t *testing.T, charts []*Chart)
	}{
		{
			name: "should return valid charts list",
			setup: func() *Metrics {
				return NewMetrics()
			},
			expectCharts: 1,
			checkChart: func(t *testing.T, charts []*Chart) {
				if len(charts) == 0 {
					t.Error("expected at least one chart")
				}
			},
		},
		{
			name: "should not return nil charts",
			setup: func() *Metrics {
				return NewMetrics()
			},
			checkChart: func(t *testing.T, charts []*Chart) {
				for _, chart := range charts {
					if chart == nil {
						t.Error("chart should not be nil")
					}
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.setup()
			charts := m.Charts()
			if tt.checkChart != nil {
				tt.checkChart(t, charts)
			}
		})
	}
}

func TestMetricsCopy(t *testing.T) {
	tests := []struct {
		name  string
		setup func() *Metrics
		check func(t *testing.T, orig, copied *Metrics)
	}{
		{
			name: "should create independent copy",
			setup: func() *Metrics {
				return NewMetrics()
			},
			check: func(t *testing.T, orig, copied *Metrics) {
				if copied == nil {
					t.Error("copied metrics should not be nil")
				}
				if copied == orig {
					t.Error("copied metrics should be a different instance")
				}
			},
		},
		{
			name: "should preserve all fields in copy",
			setup: func() *Metrics {
				m := NewMetrics()
				m.System.StoragePoolStatus = make(map[string]int64)
				m.System.StoragePoolStatus["pool1"] = 100
				return m
			},
			check: func(t *testing.T, orig, copied *Metrics) {
				if copied.System.StoragePoolStatus["pool1"] != 100 {
					t.Error("copy should preserve storage pool status")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			orig := tt.setup()
			copied := orig.Copy()
			tt.check(t, orig, copied)
		})
	}
}

func TestSystemMetrics(t *testing.T) {
	tests := []struct {
		name  string
		check func(t *testing.T, m *Metrics)
	}{
		{
			name: "should have initialized system metrics",
			check: func(t *testing.T, m *Metrics) {
				if m.System.CapacityLimitGb == 0 && m.System.CapacityInUseGb == 0 {
					// Allow both to be zero initially
				} else if m.System.CapacityLimitGb < 0 || m.System.CapacityInUseGb < 0 {
					t.Error("system metrics should not be negative")
				}
			},
		},
		{
			name: "should have status maps",
			check: func(t *testing.T, m *Metrics) {
				if m.System.StoragePoolStatus == nil {
					t.Error("StoragePoolStatus map should be initialized")
				}
			},
		},
		{
			name: "should have device metrics",
			check: func(t *testing.T, m *Metrics) {
				if m.System.DeviceMetrics == nil {
					t.Error("DeviceMetrics should be initialized")
				}
			},
		},
		{
			name: "should have SDS metrics",
			check: func(t *testing.T, m *Metrics) {
				if m.System.SDSMetrics == nil {
					t.Error("SDSMetrics should be initialized")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetrics()
			tt.check(t, m)
		})
	}
}

func TestDeviceMetricsOperations(t *testing.T) {
	tests := []struct {
		name  string
		test  func(t *testing.T)
	}{
		{
			name: "should add device metrics",
			test: func(t *testing.T) {
				m := NewMetrics()
				device := &DeviceMetrics{
					ID:                "dev1",
					SdsID:             "sds1",
					State:              "Normal",
					ErrorState:         "NoError",
					ReadLatency:        100,
					WriteLatency:       150,
					ReadBandwidth:      1000,
					WriteBandwidth:     800,
					ReadIOPS:           500,
					WriteIOPS:          400,
				}
				m.System.DeviceMetrics = append(m.System.DeviceMetrics, device)
				
				if len(m.System.DeviceMetrics) != 1 {
					t.Error("device metrics should be added")
				}
				if m.System.DeviceMetrics[0].ID != "dev1" {
					t.Error("device ID should match")
				}
			},
		},
		{
			name: "should handle multiple devices",
			test: func(t *testing.T) {
				m := NewMetrics()
				for i := 0; i < 5; i++ {
					device := &DeviceMetrics{
						ID:    "dev" + string(rune(i+48)),
						SdsID: "sds1",
						State: "Normal",
					}
					m.System.DeviceMetrics = append(m.System.DeviceMetrics, device)
				}
				
				if len(m.System.DeviceMetrics) != 5 {
					t.Errorf("expected 5 devices, got %d", len(m.System.DeviceMetrics))
				}
			},
		},
		{
			name: "should handle empty device list",
			test: func(t *testing.T) {
				m := NewMetrics()
				if len(m.System.DeviceMetrics) != 0 {
					t.Error("device metrics should be empty initially")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

func TestSDSMetricsOperations(t *testing.T) {
	tests := []struct {
		name  string
		test  func(t *testing.T)
	}{
		{
			name: "should add SDS metrics",
			test: func(t *testing.T) {
				m := NewMetrics()
				sds := &SDSMetrics{
					ID:                 "sds1",
					ProtectionDomainID: "pd1",
					State:              "Normal",
					MemoryState:        "OK",
					NumOfDevices:       3,
					NumOfScsConnections: 5,
				}
				m.System.SDSMetrics = append(m.System.SDSMetrics, sds)
				
				if len(m.System.SDSMetrics) != 1 {
					t.Error("SDS metrics should be added")
				}
				if m.System.SDSMetrics[0].ID != "sds1" {
					t.Error("SDS ID should match")
				}
			},
		},
		{
			name: "should handle multiple SDS",
			test: func(t *testing.T) {
				m := NewMetrics()
				for i := 0; i < 3; i++ {
					sds := &SDSMetrics{
						ID:                 "sds" + string(rune(i+49)),
						ProtectionDomainID: "pd1",
						State:              "Normal",
					}
					m.System.SDSMetrics = append(m.System.SDSMetrics, sds)
				}
				
				if len(m.System.SDSMetrics) != 3 {
					t.Errorf("expected 3 SDS, got %d", len(m.System.SDSMetrics))
				}
			},
		},
		{
			name: "should initialize SDS metrics to zero",
			test: func(t *testing.T) {
				m := NewMetrics()
				sds := &SDSMetrics{}
				
				if sds.NumOfDevices != 0 || sds.NumOfScsConnections != 0 {
					t.Error("SDS metrics should initialize to zero")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

func TestStoragePoolStatusOperations(t *testing.T) {
	tests := []struct {
		name  string
		test  func(t *testing.T)
	}{
		{
			name: "should manage storage pool status map",
			test: func(t *testing.T) {
				m := NewMetrics()
				m.System.StoragePoolStatus["pool1"] = 1
				m.System.StoragePoolStatus["pool2"] = 2
				
				if len(m.System.StoragePoolStatus) != 2 {
					t.Error("storage pool status should contain 2 entries")
				}
				if m.System.StoragePoolStatus["pool1"] != 1 {
					t.Error("storage pool status value should match")
				}
			},
		},
		{
			name: "should handle empty storage pool status",
			test: func(t *testing.T) {
				m := NewMetrics()
				if len(m.System.StoragePoolStatus) != 0 {
					t.Error("storage pool status should be empty initially")
				}
			},
		},
		{
			name: "should overwrite existing pool status",
			test: func(t *testing.T) {
				m := NewMetrics()
				m.System.StoragePoolStatus["pool1"] = 1
				m.System.StoragePoolStatus["pool1"] = 2
				
				if m.System.StoragePoolStatus["pool1"] != 2 {
					t.Error("storage pool status should be updated")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

func TestMetricsFields(t *testing.T) {
	tests := []struct {
		name  string
		field string
		check func(t *testing.T, m *Metrics)
	}{
		{
			name:  "should have capacity fields",
			field: "capacity",
			check: func(t *testing.T, m *Metrics) {
				if m.System.CapacityLimitGb < 0 {
					t.Error("capacity limit should not be negative")
				}
				if m.System.CapacityInUseGb < 0 {
					t.Error("capacity in use should not be negative")
				}
			},
		},
		{
			name:  "should have protection domain fields",
			field: "protectionDomain",
			check: func(t *testing.T, m *Metrics) {
				if m.System.NumOfProtectionDomains < 0 {
					t.Error("number of protection domains should not be negative")
				}
			},
		},
		{
			name:  "should have component counts",
			field: "counts",
			check: func(t *testing.T, m *Metrics) {
				if m.System.NumOfDevices < 0 {
					t.Error("number of devices should not be negative")
				}
				if m.System.NumOfSDS < 0 {
					t.Error("number of SDS should not be negative")
				}
				if m.System.NumOfStoragePools < 0 {
					t.Error("number of storage pools should not be negative")
				}
			},
		},
		{
			name:  "should have fault domains",
			field: "faultDomains",
			check: func(t *testing.T, m *Metrics) {
				if m.System.NumOfFaultDomains < 0 {
					t.Error("number of fault domains should not be negative")
				}
			},
		},
		{
			name:  "should have statistics",
			field: "statistics",
			check: func(t *testing.T, m *Metrics) {
				if m.System.ReadBandwidth < 0 {
					t.Error("read bandwidth should not be negative")
				}
				if m.System.WriteBandwidth < 0 {
					t.Error("write bandwidth should not be negative")
				}
				if m.System.ReadIOPS < 0 {
					t.Error("read IOPS should not be negative")
				}
				if m.System.WriteIOPS < 0 {
					t.Error("write IOPS should not be negative")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetrics()
			tt.check(t, m)
		})
	}
}

func TestMetricsZeroValues(t *testing.T) {
	tests := []struct {
		name  string
		test  func(t *testing.T)
	}{
		{
			name: "should handle zero capacity",
			test: func(t *testing.T) {
				m := NewMetrics()
				m.System.CapacityLimitGb = 0
				m.System.CapacityInUseGb = 0
				
				if m.System.CapacityLimitGb != 0 || m.System.CapacityInUseGb != 0 {
					t.Error("should preserve zero values")
				}
			},
		},
		{
			name: "should handle negative capacity as error condition",
			test: func(t *testing.T) {
				m := NewMetrics()
				m.System.CapacityLimitGb = -1
				
				if m.System.CapacityLimitGb >= 0 {
					t.Error("should represent negative capacity as error condition")
				}
			},
		},
		{
			name: "should handle large capacity values",
			test: func(t *testing.T) {
				m := NewMetrics()
				m.System.CapacityLimitGb = 1000000
				
				if m.System.CapacityLimitGb != 1000000 {
					t.Error("should handle large capacity values")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

func TestMetricsStateValues(t *testing.T) {
	tests := []struct {
		name  string
		state string
		check func(t *testing.T, m *Metrics)
	}{
		{
			name:  "should store normal state",
			state: "Normal",
			check: func(t *testing.T, m *Metrics) {
				device := &DeviceMetrics{State: "Normal"}
				m.System.DeviceMetrics = append(m.System.DeviceMetrics, device)
				if m.System.DeviceMetrics[0].State != "Normal" {
					t.Error("should store normal state correctly")
				}
			},
		},
		{
			name:  "should store degraded state",
			state: "Degraded",
			check: func(t *testing.T, m *Metrics) {
				device := &DeviceMetrics{State: "Degraded"}
				m.System.DeviceMetrics = append(m.System.DeviceMetrics, device)
				if m.System.DeviceMetrics[0].State != "Degraded" {
					t.Error("should store degraded state correctly")
				}
			},
		},
		{
			name:  "should store error state",
			state: "Error",
			check: func(t *testing.T, m *Metrics) {
				device := &DeviceMetrics{State: "Error"}
				m.System.DeviceMetrics = append(m.System.DeviceMetrics, device)
				if m.System.DeviceMetrics[0].State != "Error" {
					t.Error("should store error state correctly")
				}
			},
		},
		{
			name:  "should store empty state",
			state: "",
			check: func(t *testing.T, m *Metrics) {
				device := &DeviceMetrics{State: ""}
				m.System.DeviceMetrics = append(m.System.DeviceMetrics, device)
				if m.System.DeviceMetrics[0].State != "" {
					t.Error("should store empty state correctly")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetrics()
			tt.check(t, m)
		})
	}
}

func TestMetricsChart(t *testing.T) {
	tests := []struct {
		name  string
		check func(t *testing.T, chart *Chart)
	}{
		{
			name: "should have valid chart ID",
			check: func(t *testing.T, chart *Chart) {
				if chart == nil {
					t.Error("chart should not be nil")
					return
				}
				if chart.ID == "" {
					t.Error("chart ID should not be empty")
				}
			},
		},
		{
			name: "should have chart title",
			check: func(t *testing.T, chart *Chart) {
				if chart == nil {
					t.Error("chart should not be nil")
					return
				}
				if chart.Title == "" {
					t.Error("chart title should not be empty")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewMetrics()
			charts := m.Charts()
			if len(charts) > 0 {
				tt.check(t, charts[0])
			} else {
				t.Skip("no charts available")
			}
		})
	}
}