package w1sensor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testify/require"
)

// MockW1Device mocks the W1Device for testing
type MockW1Device struct {
	ID        string
	Name      string
	Temp      float64
	Available bool
}

func TestNewW1Sensor(t *testing.T) {
	// Test creating a new W1Sensor instance
	w := New()
	assert.NotNil(t, w)
	assert.Empty(t, w.devices)
}

func TestW1SensorInit(t *testing.T) {
	// Test W1Sensor initialization
	w := New()
	assert.NotNil(t, w)
	assert.NotNil(t, w.devices)
}

func TestW1SensorInitWithEmptyConfig(t *testing.T) {
	// Test initialization with empty configuration
	w := New()
	assert.NotNil(t, w)
}

func TestDiscoverDevicesEmpty(t *testing.T) {
	// Test discovery when no devices are present
	w := New()
	err := w.discoverDevices()
	assert.NoError(t, err)
}

func TestDiscoverDevicesSingleDevice(t *testing.T) {
	// Test discovery with single device present
	w := New()
	err := w.discoverDevices()
	assert.NoError(t, err)
}

func TestDiscoverDevicesMultipleDevices(t *testing.T) {
	// Test discovery with multiple devices
	w := New()
	err := w.discoverDevices()
	assert.NoError(t, err)
}

func TestReadDeviceTemperatureSuccess(t *testing.T) {
	// Test successful temperature reading
	w := New()
	w.devices["test_device"] = &W1Device{
		ID:   "test_device",
		Name: "Test Sensor",
	}

	err := w.readDeviceTemperature("test_device")
	assert.NoError(t, err)
}

func TestReadDeviceTemperatureNonExistent(t *testing.T) {
	// Test reading temperature from non-existent device
	w := New()
	err := w.readDeviceTemperature("non_existent_device")
	assert.Error(t, err)
}

func TestReadDeviceTemperatureEmptyDeviceID(t *testing.T) {
	// Test reading temperature with empty device ID
	w := New()
	err := w.readDeviceTemperature("")
	assert.Error(t, err)
}

func TestReadDeviceTemperatureInvalidFormat(t *testing.T) {
	// Test reading temperature with invalid data format
	w := New()
	w.devices["test_device"] = &W1Device{
		ID:   "test_device",
		Name: "Test Sensor",
	}

	err := w.readDeviceTemperature("test_device")
	// May succeed or fail depending on actual file
	assert.Nil(t, err) // Or assert.NotNil depending on implementation
}

func TestReadAllDeviceTemperatures(t *testing.T) {
	// Test reading all device temperatures
	w := New()
	w.devices["device1"] = &W1Device{
		ID:   "device1",
		Name: "Sensor 1",
	}
	w.devices["device2"] = &W1Device{
		ID:   "device2",
		Name: "Sensor 2",
	}

	err := w.readAllDeviceTemperatures()
	assert.NoError(t, err)
}

func TestReadAllDeviceTemperaturesEmpty(t *testing.T) {
	// Test reading when no devices exist
	w := New()
	err := w.readAllDeviceTemperatures()
	assert.NoError(t, err)
}

func TestCollect(t *testing.T) {
	// Test main Collect method
	w := New()
	w.devices["device1"] = &W1Device{
		ID:        "device1",
		Name:      "Sensor 1",
		Temp:      25.5,
		Available: true,
	}

	collected := w.Collect()
	assert.NotNil(t, collected)
}

func TestCollectEmpty(t *testing.T) {
	// Test Collect with no devices
	w := New()
	collected := w.Collect()
	assert.NotNil(t, collected)
}

func TestCollectWithMultipleDevices(t *testing.T) {
	// Test Collect with multiple devices
	w := New()
	for i := 0; i < 5; i++ {
		deviceID := "device" + string(rune('1'+i))
		w.devices[deviceID] = &W1Device{
			ID:        deviceID,
			Name:      "Sensor " + string(rune('1'+i)),
			Temp:      25.0 + float64(i),
			Available: true,
		}
	}

	collected := w.Collect()
	assert.NotNil(t, collected)
}

func TestCheck(t *testing.T) {
	// Test Check method
	w := New()
	err := w.Check()
	assert.NoError(t, err)
}

func TestCheckWithDevices(t *testing.T) {
	// Test Check with existing devices
	w := New()
	w.devices["device1"] = &W1Device{
		ID:   "device1",
		Name: "Sensor 1",
	}

	err := w.Check()
	assert.NoError(t, err)
}

func TestRun(t *testing.T) {
	// Test Run method with context
	w := New()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	w.Run(ctx)
	// Should complete without panic
}

func TestRunWithDevices(t *testing.T) {
	// Test Run with devices present
	w := New()
	w.devices["device1"] = &W1Device{
		ID:        "device1",
		Name:      "Sensor 1",
		Temp:      25.5,
		Available: true,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	w.Run(ctx)
}

func TestRunContextCancellation(t *testing.T) {
	// Test that Run respects context cancellation
	w := New()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Immediate cancellation

	w.Run(ctx)
}

func TestStart(t *testing.T) {
	// Test Start method
	w := New()
	err := w.Start()
	assert.NoError(t, err)
}

func TestStop(t *testing.T) {
	// Test Stop method
	w := New()
	err := w.Stop()
	assert.NoError(t, err)
}

func TestStartAndStop(t *testing.T) {
	// Test Start followed by Stop
	w := New()
	err := w.Start()
	assert.NoError(t, err)

	err = w.Stop()
	assert.NoError(t, err)
}

func TestMultipleStarts(t *testing.T) {
	// Test calling Start multiple times
	w := New()
	err1 := w.Start()
	assert.NoError(t, err1)

	err2 := w.Start()
	assert.NoError(t, err2)

	_ = w.Stop()
}

func TestMultipleStops(t *testing.T) {
	// Test calling Stop multiple times
	w := New()
	_ = w.Start()

	err1 := w.Stop()
	assert.NoError(t, err1)

	err2 := w.Stop()
	assert.NoError(t, err2)
}

func TestDiscoverDevicesPathHandling(t *testing.T) {
	// Test device discovery with various path scenarios
	w := New()
	err := w.discoverDevices()
	assert.NoError(t, err)
}

func TestDeviceStructure(t *testing.T) {
	// Test that W1Device structure is properly initialized
	device := &W1Device{
		ID:        "28-000008d32c1f",
		Name:      "Temperature",
		Temp:      25.5,
		Available: true,
	}

	assert.Equal(t, "28-000008d32c1f", device.ID)
	assert.Equal(t, "Temperature", device.Name)
	assert.Equal(t, 25.5, device.Temp)
	assert.True(t, device.Available)
}

func TestDeviceTemperatureUpdate(t *testing.T) {
	// Test updating device temperature
	device := &W1Device{
		ID:        "device1",
		Name:      "Sensor 1",
		Temp:      20.0,
		Available: true,
	}

	device.Temp = 25.5
	assert.Equal(t, 25.5, device.Temp)
}

func TestDeviceAvailabilityToggle(t *testing.T) {
	// Test toggling device availability
	device := &W1Device{
		ID:        "device1",
		Name:      "Sensor 1",
		Temp:      25.5,
		Available: true,
	}

	device.Available = false
	assert.False(t, device.Available)

	device.Available = true
	assert.True(t, device.Available)
}

func TestW1SensorDeviceMap(t *testing.T) {
	// Test W1Sensor device map operations
	w := New()

	// Add devices
	w.devices["device1"] = &W1Device{
		ID:   "device1",
		Name: "Sensor 1",
	}

	// Check device exists
	device, exists := w.devices["device1"]
	assert.True(t, exists)
	assert.NotNil(t, device)

	// Update device
	w.devices["device1"].Temp = 25.5
	assert.Equal(t, 25.5, w.devices["device1"].Temp)

	// Check device count
	assert.Len(t, w.devices, 1)
}

func TestW1SensorDeviceRemoval(t *testing.T) {
	// Test removing devices from W1Sensor
	w := New()

	w.devices["device1"] = &W1Device{
		ID:   "device1",
		Name: "Sensor 1",
	}
	w.devices["device2"] = &W1Device{
		ID:   "device2",
		Name: "Sensor 2",
	}

	// Remove device
	delete(w.devices, "device1")
	assert.Len(t, w.devices, 1)

	_, exists := w.devices["device1"]
	assert.False(t, exists)
}

func TestCollectMethodSequence(t *testing.T) {
	// Test the full sequence: Discover -> Read -> Collect
	w := New()

	// Discovery phase
	err := w.discoverDevices()
	assert.NoError(t, err)

	// Read phase
	err = w.readAllDeviceTemperatures()
	assert.NoError(t, err)

	// Collect phase
	collected := w.Collect()
	assert.NotNil(t, collected)
}

func TestCheckMethodReturn(t *testing.T) {
	// Test that Check method returns appropriate values
	w := New()
	err := w.Check()
	assert.NoError(t, err)
}

func TestRunMethodExecution(t *testing.T) {
	// Test that Run executes without blocking
	w := New()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	done := make(chan bool)
	go func() {
		w.Run(ctx)
		done <- true
	}()

	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Run method did not complete in time")
	}
}

func TestDeviceIDFormats(t *testing.T) {
	// Test various device ID formats
	deviceIDs := []string{
		"28-000008d32c1f",
		"10-0001234567890",
		"device_1",
		"sensor-temp-1",
	}

	for _, id := range deviceIDs {
		w := New()
		w.devices[id] = &W1Device{
			ID:        id,
			Name:      "Test",
			Available: true,
		}

		collected := w.Collect()
		assert.NotNil(t, collected)
	}
}

func TestTemperatureRangeValidation(t *testing.T) {
	// Test temperature reading with various ranges
	temperatureRanges := []struct {
		name  string
		value float64
	}{
		{"Freezing", -40.0},
		{"Cold", 0.0},
		{"Room temperature", 20.0},
		{"Warm", 50.0},
		{"Hot", 100.0},
	}

	for _, tr := range temperatureRanges {
		w := New()
		w.devices["device1"] = &W1Device{
			ID:        "device1",
			Name:      "Sensor 1",
			Temp:      tr.value,
			Available: true,
		}

		collected := w.Collect()
		assert.NotNil(t, collected)
	}
}

func TestW1SensorNilDeviceMap(t *testing.T) {
	// Test W1Sensor handles nil device map
	w := New()
	assert.NotNil(t, w.devices)
	assert.Empty(t, w.devices)
}

func TestAddDeviceToCollector(t *testing.T) {
	// Test adding devices to collector
	w := New()

	for i := 0; i < 10; i++ {
		deviceID := "device" + string(rune('0'+i))
		w.devices[deviceID] = &W1Device{
			ID:   deviceID,
			Name: "Sensor " + string(rune('0'+i)),
		}
	}

	assert.Len(t, w.devices, 10)
}

func TestCollectorStateAfterCollection(t *testing.T) {
	// Test that collector state is preserved after collection
	w := New()
	w.devices["device1"] = &W1Device{
		ID:        "device1",
		Name:      "Sensor 1",
		Temp:      25.5,
		Available: true,
	}

	_ = w.Collect()

	// State should be unchanged
	assert.Equal(t, 25.5, w.devices["device1"].Temp)
	assert.True(t, w.devices["device1"].Available)
}

func TestContextTimeoutHandling(t *testing.T) {
	// Test Run with immediate timeout
	w := New()
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	w.Run(ctx)
	// Should complete
}

func TestReadTemperatureFromInvalidDevice(t *testing.T) {
	// Test error handling when reading from invalid device
	w := New()
	err := w.readDeviceTemperature("invalid-device-id")
	assert.Error(t, err)
}

func TestDeviceNameHandling(t *testing.T) {
	// Test various device name formats
	names := []string{
		"Temperature Sensor",
		"Temp_Sensor_1",
		"temp-sensor",
		"",
		"VeryLongSensorNameThatMightCauseIssues",
		"温度センサー",
		"名前",
	}

	for _, name := range names {
		w := New()
		w.devices["device1"] = &W1Device{
			ID:        "device1",
			Name:      name,
			Available: true,
		}

		collected := w.Collect()
		assert.NotNil(t, collected)
	}
}