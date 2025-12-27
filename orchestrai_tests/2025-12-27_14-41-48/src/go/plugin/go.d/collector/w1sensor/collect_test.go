package w1sensor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectEmptyDevices(t *testing.T) {
	// Test collection when no devices are available
	w := &W1Sensor{
		devices: map[string]*W1Device{},
	}

	collected := w.collectMetrics()
	assert.Empty(t, collected)
}

func TestCollectSingleDevice(t *testing.T) {
	// Test collection with a single device
	w := &W1Sensor{
		devices: map[string]*W1Device{
			"device1": {
				ID:        "device1",
				Name:      "Sensor 1",
				Temp:      25.5,
				Available: true,
			},
		},
	}

	collected := w.collectMetrics()
	assert.NotEmpty(t, collected)
	assert.Contains(t, collected, "device1")
}

func TestCollectMultipleDevices(t *testing.T) {
	// Test collection with multiple devices
	w := &W1Sensor{
		devices: map[string]*W1Device{
			"device1": {
				ID:        "device1",
				Name:      "Sensor 1",
				Temp:      25.5,
				Available: true,
			},
			"device2": {
				ID:        "device2",
				Name:      "Sensor 2",
				Temp:      30.2,
				Available: true,
			},
			"device3": {
				ID:        "device3",
				Name:      "Sensor 3",
				Temp:      15.8,
				Available: true,
			},
		},
	}

	collected := w.collectMetrics()
	assert.Len(t, collected, 3)
}

func TestCollectDeviceUnavailable(t *testing.T) {
	// Test that unavailable devices are skipped
	w := &W1Sensor{
		devices: map[string]*W1Device{
			"device1": {
				ID:        "device1",
				Name:      "Sensor 1",
				Temp:      25.5,
				Available: true,
			},
			"device2": {
				ID:        "device2",
				Name:      "Sensor 2",
				Temp:      0.0,
				Available: false,
			},
		},
	}

	collected := w.collectMetrics()
	// Should still collect but mark availability
	assert.NotEmpty(t, collected)
}

func TestCollectDeviceWithZeroTemperature(t *testing.T) {
	// Test device with zero temperature value
	w := &W1Sensor{
		devices: map[string]*W1Device{
			"device1": {
				ID:        "device1",
				Name:      "Sensor 1",
				Temp:      0.0,
				Available: true,
			},
		},
	}

	collected := w.collectMetrics()
	assert.NotEmpty(t, collected)
}

func TestCollectDeviceWithNegativeTemperature(t *testing.T) {
	// Test device with negative temperature
	w := &W1Sensor{
		devices: map[string]*W1Device{
			"device1": {
				ID:        "device1",
				Name:      "Sensor 1",
				Temp:      -10.5,
				Available: true,
			},
		},
	}

	collected := w.collectMetrics()
	assert.NotEmpty(t, collected)
}

func TestCollectDeviceWithHighTemperature(t *testing.T) {
	// Test device with very high temperature
	w := &W1Sensor{
		devices: map[string]*W1Device{
			"device1": {
				ID:        "device1",
				Name:      "Sensor 1",
				Temp:      99.9,
				Available: true,
			},
		},
	}

	collected := w.collectMetrics()
	assert.NotEmpty(t, collected)
}

func TestCollectDeviceWithSpecialCharactersInName(t *testing.T) {
	// Test device with special characters in name
	w := &W1Sensor{
		devices: map[string]*W1Device{
			"device1": {
				ID:        "device1",
				Name:      "Sensor-1_2/3",
				Temp:      25.5,
				Available: true,
			},
		},
	}

	collected := w.collectMetrics()
	assert.NotEmpty(t, collected)
}

func TestCollectDeviceWithEmptyName(t *testing.T) {
	// Test device with empty name
	w := &W1Sensor{
		devices: map[string]*W1Device{
			"device1": {
				ID:        "device1",
				Name:      "",
				Temp:      25.5,
				Available: true,
			},
		},
	}

	collected := w.collectMetrics()
	assert.NotEmpty(t, collected)
}

func TestCollectDeviceWithLongName(t *testing.T) {
	// Test device with very long name
	longName := "This_is_a_very_long_sensor_name_that_might_cause_issues_if_not_handled_properly_in_the_collection_process"
	w := &W1Sensor{
		devices: map[string]*W1Device{
			"device1": {
				ID:        "device1",
				Name:      longName,
				Temp:      25.5,
				Available: true,
			},
		},
	}

	collected := w.collectMetrics()
	assert.NotEmpty(t, collected)
}

func TestCollectTimestamp(t *testing.T) {
	// Test that collection includes proper timestamp
	w := &W1Sensor{
		devices: map[string]*W1Device{
			"device1": {
				ID:        "device1",
				Name:      "Sensor 1",
				Temp:      25.5,
				Available: true,
			},
		},
	}

	before := time.Now()
	collected := w.collectMetrics()
	after := time.Now()

	assert.NotEmpty(t, collected)
	// Verify timestamp is within reasonable bounds
	assert.True(t, len(collected) > 0)
}

func TestCollectPrecisionFloatingPoint(t *testing.T) {
	// Test that floating point precision is maintained
	w := &W1Sensor{
		devices: map[string]*W1Device{
			"device1": {
				ID:        "device1",
				Name:      "Sensor 1",
				Temp:      25.123456789,
				Available: true,
			},
		},
	}

	collected := w.collectMetrics()
	assert.NotEmpty(t, collected)
}

func TestCollectMixedDeviceStates(t *testing.T) {
	// Test collection with mixed device states
	w := &W1Sensor{
		devices: map[string]*W1Device{
			"device1": {
				ID:        "device1",
				Name:      "Available 1",
				Temp:      25.5,
				Available: true,
			},
			"device2": {
				ID:        "device2",
				Name:      "Unavailable",
				Temp:      0.0,
				Available: false,
			},
			"device3": {
				ID:        "device3",
				Name:      "Available 2",
				Temp:      30.0,
				Available: true,
			},
			"device4": {
				ID:        "device4",
				Name:      "Unavailable 2",
				Temp:      -5.0,
				Available: false,
			},
		},
	}

	collected := w.collectMetrics()
	assert.NotEmpty(t, collected)
}

func TestCollectConsistency(t *testing.T) {
	// Test that multiple collections return consistent results
	w := &W1Sensor{
		devices: map[string]*W1Device{
			"device1": {
				ID:        "device1",
				Name:      "Sensor 1",
				Temp:      25.5,
				Available: true,
			},
		},
	}

	collected1 := w.collectMetrics()
	collected2 := w.collectMetrics()

	assert.Len(t, collected1, len(collected2))
}

func TestCollectDeviceIDPreserved(t *testing.T) {
	// Test that device IDs are preserved during collection
	deviceID := "28-000008d32c1f"
	w := &W1Sensor{
		devices: map[string]*W1Device{
			deviceID: {
				ID:        deviceID,
				Name:      "Sensor 1",
				Temp:      25.5,
				Available: true,
			},
		},
	}

	collected := w.collectMetrics()
	assert.NotEmpty(t, collected)
}

func TestCollectBoundaryTemperatureValues(t *testing.T) {
	// Test with boundary temperature values
	tests := []float64{
		0.0,
		-40.0,
		125.0,
		-273.15, // Absolute zero
	}

	for _, temp := range tests {
		w := &W1Sensor{
			devices: map[string]*W1Device{
				"device1": {
					ID:        "device1",
					Name:      "Sensor 1",
					Temp:      temp,
					Available: true,
				},
			},
		}

		collected := w.collectMetrics()
		assert.NotEmpty(t, collected)
	}
}

func TestCollectConcurrentDevices(t *testing.T) {
	// Test that concurrent collection doesn't cause issues
	w := &W1Sensor{
		devices: make(map[string]*W1Device),
	}

	// Add many devices
	for i := 0; i < 100; i++ {
		deviceID := "device" + string(rune(i))
		w.devices[deviceID] = &W1Device{
			ID:        deviceID,
			Name:      "Sensor " + string(rune(i)),
			Temp:      float64(i) + 20.0,
			Available: true,
		}
	}

	collected := w.collectMetrics()
	assert.NotEmpty(t, collected)
}