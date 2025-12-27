// SPDX-License-Identifier: GPL-3.0-or-later

package smartctl

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

// Mock executor for testing
type mockExecutor struct {
	mock.Mock
}

func (m *mockExecutor) scan(scanOpen bool) (*gjson.Result, error) {
	args := m.Called(scanOpen)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gjson.Result), args.Error(1)
}

func (m *mockExecutor) deviceInfo(name, typ, noCheckPowerMode string) (*gjson.Result, error) {
	args := m.Called(name, typ, noCheckPowerMode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*gjson.Result), args.Error(1)
}

// Mock logger
type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) Infof(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *mockLogger) Warningf(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *mockLogger) Debugf(msg string, args ...interface{}) {
	m.Called(msg, args)
}

func (m *mockLogger) Errorf(msg string, args ...interface{}) {
	m.Called(msg, args)
}

// Test scanDevice.key()
func TestScanDeviceKey(t *testing.T) {
	tests := []struct {
		name     string
		device   scanDevice
		expected string
	}{
		{
			name:     "normal device",
			device:   scanDevice{name: "/dev/sda", typ: "sat"},
			expected: "/dev/sda|sat",
		},
		{
			name:     "scsi device",
			device:   scanDevice{name: "/dev/sdb", typ: "scsi"},
			expected: "/dev/sdb|scsi",
		},
		{
			name:     "empty name",
			device:   scanDevice{name: "", typ: "sat"},
			expected: "|sat",
		},
		{
			name:     "empty type",
			device:   scanDevice{name: "/dev/sdc", typ: ""},
			expected: "/dev/sdc|",
		},
		{
			name:     "both empty",
			device:   scanDevice{name: "", typ: ""},
			expected: "|",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.device.key()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test scanDevice.shortName()
func TestScanDeviceShortName(t *testing.T) {
	tests := []struct {
		name     string
		device   scanDevice
		expected string
	}{
		{
			name:     "device with /dev/ prefix",
			device:   scanDevice{name: "/dev/sda"},
			expected: "sda",
		},
		{
			name:     "device with /dev/ prefix multiple segments",
			device:   scanDevice{name: "/dev/nvme0n1"},
			expected: "nvme0n1",
		},
		{
			name:     "device without /dev/ prefix",
			device:   scanDevice{name: "sdb"},
			expected: "sdb",
		},
		{
			name:     "empty device name",
			device:   scanDevice{name: ""},
			expected: "",
		},
		{
			name:     "only /dev/",
			device:   scanDevice{name: "/dev/"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.device.shortName()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test Collector.scanDevices() - success cases
func TestScanDevices_Success(t *testing.T) {
	tests := []struct {
		name                  string
		noCheckPowerMode      string
		scanResponseJSON      string
		deviceSelector        string
		extraDevices          []ExtraDevice
		expectedDeviceCount   int
		expectedDeviceNames   map[string]bool
		expectedExtraIncluded bool
	}{
		{
			name:             "single device sat",
			noCheckPowerMode: "auto",
			scanResponseJSON: `{
				"devices": [
					{
						"name": "/dev/sda",
						"info_name": "sda",
						"type": "sat"
					}
				]
			}`,
			deviceSelector:      ".*",
			expectedDeviceCount: 1,
			expectedDeviceNames: map[string]bool{"/dev/sda|sat": true},
		},
		{
			name:             "multiple devices mixed types",
			noCheckPowerMode: "auto",
			scanResponseJSON: `{
				"devices": [
					{
						"name": "/dev/sda",
						"info_name": "sda",
						"type": "sat"
					},
					{
						"name": "/dev/sdb",
						"info_name": "sdb",
						"type": "scsi"
					}
				]
			}`,
			deviceSelector:      ".*",
			expectedDeviceCount: 2,
			expectedDeviceNames: map[string]bool{"/dev/sda|sat": true, "/dev/sdb|scsi": true},
		},
		{
			name:             "device selector filtering",
			noCheckPowerMode: "auto",
			scanResponseJSON: `{
				"devices": [
					{
						"name": "/dev/sda",
						"info_name": "sda",
						"type": "sat"
					},
					{
						"name": "/dev/sdb",
						"info_name": "sdb",
						"type": "sat"
					}
				]
			}`,
			deviceSelector:      "sda",
			expectedDeviceCount: 1,
			expectedDeviceNames: map[string]bool{"/dev/sda|sat": true},
		},
		{
			name:             "extra devices added",
			noCheckPowerMode: "auto",
			scanResponseJSON: `{
				"devices": [
					{
						"name": "/dev/sda",
						"info_name": "sda",
						"type": "sat"
					}
				]
			}`,
			deviceSelector:        ".*",
			extraDevices:          []ExtraDevice{{Name: "/dev/sdc", Type: "sat"}},
			expectedDeviceCount:   2,
			expectedDeviceNames:   map[string]bool{"/dev/sda|sat": true, "/dev/sdc|sat": true},
			expectedExtraIncluded: true,
		},
		{
			name:             "duplicate extra device ignored",
			noCheckPowerMode: "auto",
			scanResponseJSON: `{
				"devices": [
					{
						"name": "/dev/sda",
						"info_name": "sda",
						"type": "sat"
					}
				]
			}`,
			deviceSelector:      ".*",
			extraDevices:        []ExtraDevice{{Name: "/dev/sda", Type: "sat"}},
			expectedDeviceCount: 1,
			expectedDeviceNames: map[string]bool{"/dev/sda|sat": true},
		},
		{
			name:             "noCheckPowerMode never uses scan-open",
			noCheckPowerMode: "never",
			scanResponseJSON: `{
				"devices": [
					{
						"name": "/dev/sda",
						"info_name": "sda",
						"type": "sat"
					}
				]
			}`,
			deviceSelector:      ".*",
			expectedDeviceCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExec := new(mockExecutor)
			mockLogger := new(mockLogger)

			result := gjson.Parse(tt.scanResponseJSON)
			mockExec.On("scan", tt.noCheckPowerMode == "never").Return(&result, nil)
			mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
			mockLogger.On("Warningf", mock.Anything, mock.Anything).Return()

			deviceSelectorRe, _ := regexp.Compile(tt.deviceSelector)

			collector := &Collector{
				exec:              mockExec,
				deviceSr:          deviceSelectorRe,
				NoCheckPowerMode:  tt.noCheckPowerMode,
				ExtraDevices:      tt.extraDevices,
				scannedDevices:    make(map[string]*scanDevice),
			}
			collector.Debugf = mockLogger.Debugf
			collector.Warningf = mockLogger.Warningf

			devices, err := collector.scanDevices()

			require.NoError(t, err)
			assert.Equal(t, tt.expectedDeviceCount, len(devices))
			for deviceKey := range tt.expectedDeviceNames {
				assert.Contains(t, devices, deviceKey)
			}

			if tt.expectedExtraIncluded {
				extra := devices["/dev/sdc|sat"]
				assert.NotNil(t, extra)
				assert.True(t, extra.extra)
			}
		})
	}
}

// Test Collector.scanDevices() - error cases
func TestScanDevices_Errors(t *testing.T) {
	tests := []struct {
		name              string
		scanError         error
		expectedErrorMsg  string
		noCheckPowerMode  string
	}{
		{
			name:              "scan exec error",
			scanError:         fmt.Errorf("permission denied"),
			expectedErrorMsg:  "failed to scan devices: permission denied",
			noCheckPowerMode:  "auto",
		},
		{
			name:              "scan exec error with never mode",
			scanError:         fmt.Errorf("command not found"),
			expectedErrorMsg:  "failed to scan devices: command not found",
			noCheckPowerMode:  "never",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExec := new(mockExecutor)
			mockLogger := new(mockLogger)

			mockExec.On("scan", mock.Anything).Return(nil, tt.scanError)
			mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()

			collector := &Collector{
				exec:              mockExec,
				deviceSr:          regexp.MustCompile(".*"),
				NoCheckPowerMode:  tt.noCheckPowerMode,
				scannedDevices:    make(map[string]*scanDevice),
			}
			collector.Debugf = mockLogger.Debugf

			devices, err := collector.scanDevices()

			assert.Nil(t, devices)
			assert.EqualError(t, err, tt.expectedErrorMsg)
		})
	}
}

// Test Collector.scanDevices() - no devices found
func TestScanDevices_NoDevices(t *testing.T) {
	mockExec := new(mockExecutor)
	mockLogger := new(mockLogger)

	result := gjson.Parse(`{"devices": []}`)
	mockExec.On("scan", false).Return(&result, nil)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()

	collector := &Collector{
		exec:             mockExec,
		deviceSr:         regexp.MustCompile(".*"),
		NoCheckPowerMode: "auto",
		scannedDevices:   make(map[string]*scanDevice),
	}
	collector.Debugf = mockLogger.Debugf

	devices, err := collector.scanDevices()

	assert.Nil(t, devices)
	assert.EqualError(t, err, "no devices found during scan")
}

// Test Collector.scanDevices() - missing required fields
func TestScanDevices_MissingFields(t *testing.T) {
	tests := []struct {
		name             string
		scanResponseJSON string
		expectedCount    int
	}{
		{
			name: "missing name field",
			scanResponseJSON: `{
				"devices": [
					{
						"name": "",
						"info_name": "sda",
						"type": "sat"
					}
				]
			}`,
			expectedCount: 0,
		},
		{
			name: "missing type field",
			scanResponseJSON: `{
				"devices": [
					{
						"name": "/dev/sda",
						"info_name": "sda",
						"type": ""
					}
				]
			}`,
			expectedCount: 0,
		},
		{
			name: "missing both name and type",
			scanResponseJSON: `{
				"devices": [
					{
						"name": "",
						"info_name": "sda",
						"type": ""
					}
				]
			}`,
			expectedCount: 0,
		},
		{
			name: "mixed valid and invalid devices",
			scanResponseJSON: `{
				"devices": [
					{
						"name": "/dev/sda",
						"info_name": "sda",
						"type": "sat"
					},
					{
						"name": "",
						"info_name": "sdb",
						"type": "sat"
					},
					{
						"name": "/dev/sdc",
						"info_name": "sdc",
						"type": ""
					}
				]
			}`,
			expectedCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExec := new(mockExecutor)
			mockLogger := new(mockLogger)

			result := gjson.Parse(tt.scanResponseJSON)
			mockExec.On("scan", false).Return(&result, nil)
			mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()
			mockLogger.On("Warningf", mock.Anything, mock.Anything).Return()

			collector := &Collector{
				exec:             mockExec,
				deviceSr:         regexp.MustCompile(".*"),
				NoCheckPowerMode: "auto",
				scannedDevices:   make(map[string]*scanDevice),
			}
			collector.Debugf = mockLogger.Debugf
			collector.Warningf = mockLogger.Warningf

			devices, err := collector.scanDevices()

			if tt.expectedCount == 0 {
				assert.Nil(t, devices)
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(devices))
			}
		})
	}
}

// Test Collector.handleGuessedScsiScannedDevice() - no type change
func TestHandleGuessedScsiScannedDevice_NoChange(t *testing.T) {
	tests := []struct {
		name         string
		deviceType   string
		deviceExists bool
	}{
		{
			name:         "device type not scsi",
			deviceType:   "sat",
			deviceExists: false,
		},
		{
			name:         "device already scanned",
			deviceType:   "scsi",
			deviceExists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockLogger := new(mockLogger)
			mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()

			dev := &scanDevice{name: "/dev/sda", typ: tt.deviceType}

			scannedDevices := make(map[string]*scanDevice)
			if tt.deviceExists {
				scannedDevices[dev.key()] = dev
			}

			collector := &Collector{
				scannedDevices: scannedDevices,
			}
			collector.Debugf = mockLogger.Debugf

			originalType := dev.typ
			collector.handleGuessedScsiScannedDevice(dev)

			assert.Equal(t, originalType, dev.typ)
		})
	}
}

// Test Collector.handleGuessedScsiScannedDevice() - sat device in scanned list
func TestHandleGuessedScsiScannedDevice_SatInScanned(t *testing.T) {
	mockExec := new(mockExecutor)
	mockLogger := new(mockLogger)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()

	dev := &scanDevice{name: "/dev/sda", typ: "scsi"}

	// SAT version already in scanned devices
	satDev := &scanDevice{name: "/dev/sda", typ: "sat"}

	scannedDevices := make(map[string]*scanDevice)
	scannedDevices[satDev.key()] = satDev

	collector := &Collector{
		scannedDevices: scannedDevices,
		exec:           mockExec,
	}
	collector.Debugf = mockLogger.Debugf

	collector.handleGuessedScsiScannedDevice(dev)

	assert.Equal(t, "sat", dev.typ)
	mockLogger.AssertCalled(t, "Debugf", "changing device '%s' type 'scsi' -> 'sat'", "/dev/sda")
}

// Test Collector.handleGuessedScsiScannedDevice() - device info call fails
func TestHandleGuessedScsiScannedDevice_DeviceInfoFails(t *testing.T) {
	mockExec := new(mockExecutor)
	mockLogger := new(mockLogger)

	dev := &scanDevice{name: "/dev/sda", typ: "scsi"}

	mockExec.On("deviceInfo", "/dev/sda", "sat", mock.Anything).Return(nil, fmt.Errorf("error"))
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()

	collector := &Collector{
		scannedDevices: make(map[string]*scanDevice),
		exec:           mockExec,
	}
	collector.Debugf = mockLogger.Debugf

	collector.handleGuessedScsiScannedDevice(dev)

	assert.Equal(t, "scsi", dev.typ)
}

// Test Collector.handleGuessedScsiScannedDevice() - device info returns nil
func TestHandleGuessedScsiScannedDevice_DeviceInfoNil(t *testing.T) {
	mockExec := new(mockExecutor)
	mockLogger := new(mockLogger)

	dev := &scanDevice{name: "/dev/sda", typ: "scsi"}

	mockExec.On("deviceInfo", "/dev/sda", "sat", mock.Anything).Return(nil, nil)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()

	collector := &Collector{
		scannedDevices: make(map[string]*scanDevice),
		exec:           mockExec,
	}
	collector.Debugf = mockLogger.Debugf

	collector.handleGuessedScsiScannedDevice(dev)

	assert.Equal(t, "scsi", dev.typ)
}

// Test Collector.handleGuessedScsiScannedDevice() - exit status with bits set
func TestHandleGuessedScsiScannedDevice_ExitStatusBits(t *testing.T) {
	tests := []struct {
		name       string
		exitStatus int
	}{
		{
			name:       "exit status bit 0",
			exitStatus: 1,
		},
		{
			name:       "exit status bit 1",
			exitStatus: 2,
		},
		{
			name:       "exit status bit 2",
			exitStatus: 4,
		},
		{
			name:       "exit status bits 0 and 1",
			exitStatus: 3,
		},
		{
			name:       "exit status all bits",
			exitStatus: 7,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExec := new(mockExecutor)
			mockLogger := new(mockLogger)

			dev := &scanDevice{name: "/dev/sda", typ: "scsi"}

			respJSON := fmt.Sprintf(`{
				"smartctl": {
					"exit_status": %d
				}
			}`, tt.exitStatus)
			result := gjson.Parse(respJSON)

			mockExec.On("deviceInfo", "/dev/sda", "sat", mock.Anything).Return(&result, nil)
			mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()

			collector := &Collector{
				scannedDevices: make(map[string]*scanDevice),
				exec:           mockExec,
			}
			collector.Debugf = mockLogger.Debugf

			collector.handleGuessedScsiScannedDevice(dev)

			assert.Equal(t, "scsi", dev.typ)
		})
	}
}

// Test Collector.handleGuessedScsiScannedDevice() - empty attribute table
func TestHandleGuessedScsiScannedDevice_EmptyAttributeTable(t *testing.T) {
	mockExec := new(mockExecutor)
	mockLogger := new(mockLogger)

	dev := &scanDevice{name: "/dev/sda", typ: "scsi"}

	respJSON := `{
		"smartctl": {
			"exit_status": 0
		},
		"ata_smart_attributes": {
			"table": []
		}
	}`
	result := gjson.Parse(respJSON)

	mockExec.On("deviceInfo", "/dev/sda", "sat", mock.Anything).Return(&result, nil)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()

	collector := &Collector{
		scannedDevices: make(map[string]*scanDevice),
		exec:           mockExec,
	}
	collector.Debugf = mockLogger.Debugf

	collector.handleGuessedScsiScannedDevice(dev)

	assert.Equal(t, "scsi", dev.typ)
}

// Test Collector.handleGuessedScsiScannedDevice() - no attribute table
func TestHandleGuessedScsiScannedDevice_NoAttributeTable(t *testing.T) {
	mockExec := new(mockExecutor)
	mockLogger := new(mockLogger)

	dev := &scanDevice{name: "/dev/sda", typ: "scsi"}

	respJSON := `{
		"smartctl": {
			"exit_status": 0
		}
	}`
	result := gjson.Parse(respJSON)

	mockExec.On("deviceInfo", "/dev/sda", "sat", mock.Anything).Return(&result, nil)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()

	collector := &Collector{
		scannedDevices: make(map[string]*scanDevice),
		exec:           mockExec,
	}
	collector.Debugf = mockLogger.Debugf

	collector.handleGuessedScsiScannedDevice(dev)

	assert.Equal(t, "scsi", dev.typ)
}

// Test Collector.handleGuessedScsiScannedDevice() - successful type change
func TestHandleGuessedScsiScannedDevice_SuccessfulChange(t *testing.T) {
	mockExec := new(mockExecutor)
	mockLogger := new(mockLogger)

	dev := &scanDevice{name: "/dev/sda", typ: "scsi"}

	respJSON := `{
		"smartctl": {
			"exit_status": 0
		},
		"ata_smart_attributes": {
			"table": [
				{
					"id": "1",
					"name": "Raw_Read_Error_Rate",
					"value": "100",
					"raw": {"value": "0", "string": "0"}
				}
			]
		}
	}`
	result := gjson.Parse(respJSON)

	mockExec.On("deviceInfo", "/dev/sda", "sat", mock.Anything).Return(&result, nil)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()

	collector := &Collector{
		scannedDevices: make(map[string]*scanDevice),
		exec:           mockExec,
	}
	collector.Debugf = mockLogger.Debugf

	collector.handleGuessedScsiScannedDevice(dev)

	assert.Equal(t, "sat", dev.typ)
	mockLogger.AssertCalled(t, "Debugf", "changing device '%s' type 'scsi' -> 'sat'", "/dev/sda")
}

// Test Collector.handleGuessedScsiScannedDevice() - multiple attributes
func TestHandleGuessedScsiScannedDevice_MultipleAttributes(t *testing.T) {
	mockExec := new(mockExecutor)
	mockLogger := new(mockLogger)

	dev := &scanDevice{name: "/dev/sdb", typ: "scsi"}

	respJSON := `{
		"smartctl": {
			"exit_status": 0
		},
		"ata_smart_attributes": {
			"table": [
				{
					"id": "1",
					"name": "Raw_Read_Error_Rate",
					"value": "100",
					"raw": {"value": "0", "string": "0"}
				},
				{
					"id": "2",
					"name": "Throughput_Performance",
					"value": "95",
					"raw": {"value": "1234", "string": "1234"}
				},
				{
					"id": "3",
					"name": "Spin_Up_Time",
					"value": "90",
					"raw": {"value": "5678", "string": "5678"}
				}
			]
		}
	}`
	result := gjson.Parse(respJSON)

	mockExec.On("deviceInfo", "/dev/sdb", "sat", mock.Anything).Return(&result, nil)
	mockLogger.On("Debugf", mock.Anything, mock.Anything).Return()

	collector := &Collector{
		scannedDevices: make(map[string]*scanDevice),
		exec:           mockExec,
	}
	collector.Debugf = mockLogger.Debugf

	collector.handleGuessedScsiScannedDevice(dev)

	assert.Equal(t, "sat", dev.typ)
	mockLogger.AssertCalled(t, "Debugf", "changing device '%s' type 'scsi' -> 'sat'", "/dev/sdb")
}

// Test Collector.hasScannedDevice()
func TestHasScannedDevice(t *testing.T) {
	tests := []struct {
		name           string
		deviceKey      string
		scannedDevices map[string]*scanDevice
		expected       bool
	}{
		{
			name:      "device exists",
			deviceKey: "/dev/sda|sat",
			scannedDevices: map[string]*scanDevice{
				"/dev/sda|sat": {name: "/dev/sda", typ: "sat"},
			},
			expected: true,
		},
		{
			name:      "device does not exist",
			deviceKey: "/dev/sdb|scsi",
			scannedDevices: map[string]*scanDevice{
				"/dev/sda|sat": {name: "/dev/sda", typ: "sat"},
			},
			expected: false,
		},
		{
			name:            "empty scanned devices",
			deviceKey:       "/dev/sda|sat",
			scannedDevices:  map[string]*scanDevice{},
			expected:        false,
		},
		{
			name:            "multiple devices, one matches",
			deviceKey:       "/dev/sdb|scsi",
			scannedDevices: map[string]*scanDevice{
				"/dev/sda|sat":   {name: "/dev/sda", typ: "sat"},
				"/dev/sdb|scsi":  {name: "/dev/sdb", typ: "scsi"},
				"/dev/sdc|scsi":  {name: "/dev/sdc", typ: "scsi"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			collector := &Collector{
				scannedDevices: tt.scannedDevices,
			}

			dev := &scanDevice{name: tt.deviceKey[:strings.LastIndex(tt.deviceKey, "|")], typ: tt.deviceKey[strings.LastIndex(tt.deviceKey, "|")+1:]}
			result := collector.hasScannedDevice(dev)

			assert.Equal(t, tt.expected, result)
		})
	}
}