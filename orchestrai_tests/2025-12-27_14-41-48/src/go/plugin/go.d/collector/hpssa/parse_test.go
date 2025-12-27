package hpssa

import (
	"fmt"
	"testing"
)

// Test uniqueKey methods for struct types
func TestHpssaControllerUniqueKey(t *testing.T) {
	tests := []struct {
		name     string
		cntrl    *hpssaController
		expected string
	}{
		{
			name:     "should generate unique key with model and slot",
			cntrl:    &hpssaController{model: "E208i", slot: "0"},
			expected: "E208i/0",
		},
		{
			name:     "should handle different model and slot",
			cntrl:    &hpssaController{model: "E208i-p", slot: "1"},
			expected: "E208i-p/1",
		},
		{
			name:     "should handle empty model",
			cntrl:    &hpssaController{model: "", slot: "0"},
			expected: "/0",
		},
		{
			name:     "should handle empty slot",
			cntrl:    &hpssaController{model: "E208i", slot: ""},
			expected: "E208i/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.cntrl.uniqueKey()
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestHpssaArrayUniqueKey(t *testing.T) {
	tests := []struct {
		name     string
		arr      *hpssaArray
		expected string
	}{
		{
			name:     "should generate unique key with controller model, slot, and array id",
			arr:      &hpssaArray{cntrl: &hpssaController{model: "E208i", slot: "0"}, id: "A"},
			expected: "E208i/0/A",
		},
		{
			name:     "should handle multiple arrays",
			arr:      &hpssaArray{cntrl: &hpssaController{model: "E208i", slot: "1"}, id: "B"},
			expected: "E208i/1/B",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.arr.uniqueKey()
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestHpssaLogicalDriveUniqueKey(t *testing.T) {
	tests := []struct {
		name     string
		ld       *hpssaLogicalDrive
		expected string
	}{
		{
			name: "should generate unique key with all components",
			ld: &hpssaLogicalDrive{
				cntrl: &hpssaController{model: "E208i", slot: "0"},
				arr:   &hpssaArray{id: "A"},
				id:    "1",
			},
			expected: "E208i/0/A/1",
		},
		{
			name: "should handle different array and drive ids",
			ld: &hpssaLogicalDrive{
				cntrl: &hpssaController{model: "E208i", slot: "1"},
				arr:   &hpssaArray{id: "B"},
				id:    "2",
			},
			expected: "E208i/1/B/2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.ld.uniqueKey()
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestHpssaPhysicalDriveUniqueKey(t *testing.T) {
	tests := []struct {
		name     string
		pd       *hpssaPhysicalDrive
		expected string
	}{
		{
			name: "should generate unique key with all components",
			pd: &hpssaPhysicalDrive{
				cntrl:    &hpssaController{model: "E208i", slot: "0"},
				arr:      &hpssaArray{id: "A"},
				ld:       &hpssaLogicalDrive{id: "1"},
				location: "1I:1:1",
			},
			expected: "E208i/0/A/1/1I:1:1",
		},
		{
			name: "should handle unassigned drive (no array)",
			pd: &hpssaPhysicalDrive{
				cntrl:    &hpssaController{model: "E208i", slot: "0"},
				arr:      nil,
				ld:       nil,
				location: "1I:1:2",
			},
			expected: "E208i/0/na/na/1I:1:2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pd.uniqueKey()
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestHpssaPhysicalDriveArrId(t *testing.T) {
	tests := []struct {
		name     string
		pd       *hpssaPhysicalDrive
		expected string
	}{
		{
			name:     "should return array id when array is not nil",
			pd:       &hpssaPhysicalDrive{arr: &hpssaArray{id: "A"}},
			expected: "A",
		},
		{
			name:     "should return 'na' when array is nil",
			pd:       &hpssaPhysicalDrive{arr: nil},
			expected: "na",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pd.arrId()
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestHpssaPhysicalDriveLdId(t *testing.T) {
	tests := []struct {
		name     string
		pd       *hpssaPhysicalDrive
		expected string
	}{
		{
			name:     "should return logical drive id when ld is not nil",
			pd:       &hpssaPhysicalDrive{ld: &hpssaLogicalDrive{id: "1"}},
			expected: "1",
		},
		{
			name:     "should return 'na' when ld is nil",
			pd:       &hpssaPhysicalDrive{ld: nil},
			expected: "na",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.pd.ldId()
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

// Test parsing functions
func TestGetColonSepValue(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected string
	}{
		{
			name:     "should extract value after colon",
			line:     "Serial Number: ABC123",
			expected: "ABC123",
		},
		{
			name:     "should trim whitespace",
			line:     "Serial Number:   ABC123   ",
			expected: "ABC123",
		},
		{
			name:     "should handle line with no colon",
			line:     "No colon here",
			expected: "",
		},
		{
			name:     "should handle colon at start",
			line:     ": value",
			expected: "value",
		},
		{
			name:     "should handle empty value after colon",
			line:     "Serial Number:",
			expected: "",
		},
		{
			name:     "should handle colon with spaces",
			line:     "Serial Number:    ",
			expected: "",
		},
		{
			name:     "should handle multiple colons",
			line:     "Key: value: extra",
			expected: "value: extra",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getColonSepValue(tt.line)
			if result != tt.expected {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestParseControllerLine(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		shouldErr bool
		validate  func(*hpssaController) bool
	}{
		{
			name:      "should parse HPE Smart Array line",
			line:      "HPE Smart Array E208i-p SR Gen10 in Slot 0",
			shouldErr: false,
			validate: func(c *hpssaController) bool {
				return c.model == "E208i-p" &&
					c.slot == "0" &&
					len(c.arrays) == 0 &&
					len(c.unassignedDrives) == 0
			},
		},
		{
			name:      "should parse Smart Array line without HPE prefix",
			line:      "Smart Array E208i in Slot 1",
			shouldErr: false,
			validate: func(c *hpssaController) bool {
				return c.model == "E208i" && c.slot == "1"
			},
		},
		{
			name:      "should error on malformed line with too few fields",
			line:      "Smart Array E208i",
			shouldErr: true,
			validate:  func(c *hpssaController) bool { return false },
		},
		{
			name:      "should error on empty parts",
			line:      "Smart Array",
			shouldErr: true,
			validate:  func(c *hpssaController) bool { return false },
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseControllerLine(tt.line)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if !tt.validate(result) {
					t.Errorf("validation failed for result: %+v", result)
				}
			}
		})
	}
}

func TestParseArrayLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected string
	}{
		{
			name:     "should parse array line with id",
			line:     "   Array: A",
			expected: "A",
		},
		{
			name:     "should handle different array ids",
			line:     "   Array: B",
			expected: "B",
		},
		{
			name:     "should handle array with numeric id",
			line:     "   Array: 1",
			expected: "1",
		},
		{
			name:     "should handle empty array id",
			line:     "   Array:",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseArrayLine(tt.line)
			if result.id != tt.expected {
				t.Errorf("got id %q, want %q", result.id, tt.expected)
			}
			if len(result.logicalDrives) != 0 {
				t.Errorf("expected empty logicalDrives map")
			}
		})
	}
}

func TestParseLogicalDriveLine(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		expected string
	}{
		{
			name:     "should parse logical drive line with id",
			line:     "      Logical Drive: 1",
			expected: "1",
		},
		{
			name:     "should handle different logical drive ids",
			line:     "      Logical Drive: 2",
			expected: "2",
		},
		{
			name:     "should handle empty logical drive id",
			line:     "      Logical Drive:",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseLogicalDriveLine(tt.line)
			if result.id != tt.expected {
				t.Errorf("got id %q, want %q", result.id, tt.expected)
			}
			if len(result.physicalDrives) != 0 {
				t.Errorf("expected empty physicalDrives map")
			}
		})
	}
}

func TestParsePhysicalDriveLine(t *testing.T) {
	tests := []struct {
		name      string
		line      string
		shouldErr bool
		expected  string
	}{
		{
			name:      "should parse physical drive line",
			line:      "         physicaldrive 1I:1:1",
			shouldErr: false,
			expected:  "1I:1:1",
		},
		{
			name:      "should handle different locations",
			line:      "         physicaldrive 2I:2:5",
			shouldErr: false,
			expected:  "2I:2:5",
		},
		{
			name:      "should error on malformed line",
			line:      "         physicaldrive",
			shouldErr: true,
			expected:  "",
		},
		{
			name:      "should error on too many parts",
			line:      "         physicaldrive 1I:1:1 extra",
			shouldErr: true,
			expected:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parsePhysicalDriveLine(tt.line)
			if tt.shouldErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if result.location != tt.expected {
					t.Errorf("got location %q, want %q", result.location, tt.expected)
				}
			}
		})
	}
}

func TestParseControllerSectionLine(t *testing.T) {
	tests := []struct {
		name             string
		line             string
		validateProperty func(*hpssaController) bool
	}{
		{
			name: "should parse Serial Number",
			line: "   Serial Number: ABC123",
			validateProperty: func(c *hpssaController) bool {
				return c.serialNumber == "ABC123"
			},
		},
		{
			name: "should parse Controller Status",
			line: "   Controller Status: OK",
			validateProperty: func(c *hpssaController) bool {
				return c.controllerStatus == "OK"
			},
		},
		{
			name: "should parse Cache Board Present",
			line: "   Cache Board Present: True",
			validateProperty: func(c *hpssaController) bool {
				return c.cacheBoardPresent == "True"
			},
		},
		{
			name: "should parse Cache Status",
			line: "   Cache Status: OK",
			validateProperty: func(c *hpssaController) bool {
				return c.cacheStatus == "OK"
			},
		},
		{
			name: "should parse Cache Ratio",
			line: "   Cache Ratio: 25% Read / 75% Write",
			validateProperty: func(c *hpssaController) bool {
				return c.cacheRatio == "25% Read / 75% Write"
			},
		},
		{
			name: "should parse Controller Temperature (C)",
			line: "   Controller Temperature (C): 45",
			validateProperty: func(c *hpssaController) bool {
				return c.controllerTemperatureC == "45"
			},
		},
		{
			name: "should parse Cache Module Temperature (C)",
			line: "   Cache Module Temperature (C): 50",
			validateProperty: func(c *hpssaController) bool {
				return c.cacheModuleTemperatureC == "50"
			},
		},
		{
			name: "should parse Number of Ports",
			line: "   Number of Ports: 8",
			validateProperty: func(c *hpssaController) bool {
				return c.numberOfPorts == "8"
			},
		},
		{
			name: "should parse Driver Name",
			line: "   Driver Name: hpsa",
			validateProperty: func(c *hpssaController) bool {
				return c.driverName == "hpsa"
			},
		},
		{
			name: "should parse Battery/Capacitor Count",
			line: "   Battery/Capacitor Count: 1",
			validateProperty: func(c *hpssaController) bool {
				return c.batteryCapacitorCount == "1"
			},
		},
		{
			name: "should parse Battery/Capacitor Status",
			line: "   Battery/Capacitor Status: OK",
			validateProperty: func(c *hpssaController) bool {
				return c.batteryCapacitorStatus == "OK"
			},
		},
		{
			name: "should ignore unknown property",
			line: "   Unknown Property: value",
			validateProperty: func(c *hpssaController) bool {
				return c.serialNumber == ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cntrl := &hpssaController{}
			parseControllerSectionLine(tt.line, cntrl)
			if !tt.validateProperty(cntrl) {
				t.Errorf("property validation failed for line: %s", tt.line)
			}
		})
	}
}

func TestParseArraySectionLine(t *testing.T) {
	tests := []struct {
		name             string
		line             string
		validateProperty func(*hpssaArray) bool
	}{
		{
			name: "should parse Interface Type",
			line: "      Interface Type: SAS",
			validateProperty: func(a *hpssaArray) bool {
				return a.interfaceType == "SAS"
			},
		},
		{
			name: "should parse Unused Space",
			line: "      Unused Space: 0 MB",
			validateProperty: func(a *hpssaArray) bool {
				return a.unusedSpace == "0 MB"
			},
		},
		{
			name: "should parse Used Space",
			line: "      Used Space: 1000 MB",
			validateProperty: func(a *hpssaArray) bool {
				return a.usedSpace == "1000 MB"
			},
		},
		{
			name: "should parse Status",
			line: "      Status: OK",
			validateProperty: func(a *hpssaArray) bool {
				return a.status == "OK"
			},
		},
		{
			name: "should parse Array Type",
			line: "      Array Type: Data",
			validateProperty: func(a *hpssaArray) bool {
				return a.arrayType == "Data"
			},
		},
		{
			name: "should ignore unknown property",
			line: "      Unknown Property: value",
			validateProperty: func(a *hpssaArray) bool {
				return a.interfaceType == ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arr := &hpssaArray{}
			parseArraySectionLine(tt.line, arr)
			if !tt.validateProperty(arr) {
				t.Errorf("property validation failed for line: %s", tt.line)
			}
		})
	}
}

func TestParseLogicalDriveSectionLine(t *testing.T) {
	tests := []struct {
		name             string
		line             string
		validateProperty func(*hpssaLogicalDrive) bool
	}{
		{
			name: "should parse Size",
			line: "         Size: 1000 MB",
			validateProperty: func(ld *hpssaLogicalDrive) bool {
				return ld.size == "1000 MB"
			},
		},
		{
			name: "should parse Status",
			line: "         Status: OK",
			validateProperty: func(ld *hpssaLogicalDrive) bool {
				return ld.status == "OK"
			},
		},
		{
			name: "should parse Disk Name",
			line: "         Disk Name: /dev/sda",
			validateProperty: func(ld *hpssaLogicalDrive) bool {
				return ld.diskName == "/dev/sda"
			},
		},
		{
			name: "should parse Unique Identifier",
			line: "         Unique Identifier: 12345",
			validateProperty: func(ld *hpssaLogicalDrive) bool {
				return ld.uniqueIdentifier == "12345"
			},
		},
		{
			name: "should parse Logical Drive Label",
			line: "         Logical Drive Label: MyDrive",
			validateProperty: func(ld *hpssaLogicalDrive) bool {
				return ld.logicalDriveLabel == "MyDrive"
			},
		},
		{
			name: "should parse Drive Type",
			line: "         Drive Type: Data",
			validateProperty: func(ld *hpssaLogicalDrive) bool {
				return ld.driveType == "Data"
			},
		},
		{
			name: "should ignore unknown property",
			line: "         Unknown Property: value",
			validateProperty: func(ld *hpssaLogicalDrive) bool {
				return ld.size == ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ld := &hpssaLogicalDrive{}
			parseLogicalDriveSectionLine(tt.line, ld)
			if !tt.validateProperty(ld) {
				t.Errorf("property validation failed for line: %s", tt.line)
			}
		})
	}
}

func TestParsePhysicalDriveSectionLine(t *testing.T) {
	tests := []struct {
		name             string
		line             string
		validateProperty func(*hpssaPhysicalDrive) bool
	}{
		{
			name: "should parse Status",
			line: "         Status: OK",
			validateProperty: func(pd *hpssaPhysicalDrive) bool {
				return pd.status == "OK"
			},
		},
		{
			name: "should parse Drive Type",
			line: "         Drive Type: HDD SAS",
			validateProperty: func(pd *hpssaPhysicalDrive) bool {
				return pd.driveType == "HDD SAS"
			},
		},
		{
			name: "should parse Interface Type",
			line: "         Interface Type: SAS",
			validateProperty: func(pd *hpssaPhysicalDrive) bool {
				return pd.interfaceType == "SAS"
			},
		},
		{
			name: "should parse Size",
			line: "         Size: 1000 GB",
			validateProperty: func(pd *hpssaPhysicalDrive) bool {
				return pd.size == "1000 GB"
			},
		},
		{
			name: "should parse Serial Number",
			line: "         Serial Number: SERNUM123",
			validateProperty: func(pd *hpssaPhysicalDrive) bool {
				return pd.serialNumber == "SERNUM123"
			},
		},
		{
			name: "should parse WWID",
			line: "         WWID: 5000c5a12345abcd",
			validateProperty: func(pd *hpssaPhysicalDrive) bool {
				return pd.wwid == "5000c5a12345abcd"
			},
		},
		{
			name: "should parse Model",
			line: "         Model: SEAGATE ST1000DM003",
			validateProperty: func(pd *hpssaPhysicalDrive) bool {
				return pd.model == "SEAGATE ST1000DM003"
			},
		},
		{
			name: "should parse Current Temperature (C)",
			line: "         Current Temperature (C): 35",
			validateProperty: func(pd *hpssaPhysicalDrive) bool {
				return pd.currentTemperatureC == "35"
			},
		},
		{
			name: "should ignore unknown property",
			line: "         Unknown Property: value",
			validateProperty: func(pd *hpssaPhysicalDrive) bool {
				return pd.status == ""
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pd := &hpssaPhysicalDrive{}
			parsePhysicalDriveSectionLine(tt.line, pd)
			if !tt.validateProperty(pd) {
				t.Errorf("property validation failed for line: %s", tt.line)
			}
		})
	}
}

func TestUpdateHpssaHierarchy(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() map[string]*hpssaController
		validate func(map[string]*hpssaController) bool
	}{
		{
			name: "should update hierarchy for unassigned drives",
			setup: func() map[string]*hpssaController {
				cntrl := &hpssaController{
					model:            "E208i",
					slot:             "0",
					arrays:           make(map[string]*hpssaArray),
					unassignedDrives: make(map[string]*hpssaPhysicalDrive),
				}
				pd := &hpssaPhysicalDrive{location: "1I:1:1"}
				cntrl.unassignedDrives["1I:1:1"] = pd

				return map[string]*hpssaController{"0": cntrl}
			},
			validate: func(controllers map[string]*hpssaController) bool {
				cntrl := controllers["0"]
				pd := cntrl.unassignedDrives["1I:1:1"]
				return pd.cntrl == cntrl
			},
		},
		{
			name: "should update hierarchy for arrays and logical drives",
			setup: func() map[string]*hpssaController {
				cntrl := &hpssaController{
					model:            "E208i",
					slot:             "0",
					arrays:           make(map[string]*hpssaArray),
					unassignedDrives: make(map[string]*hpssaPhysicalDrive),
				}
				arr := &hpssaArray{
					id:            "A",
					logicalDrives: make(map[string]*hpssaLogicalDrive),
				}
				ld := &hpssaLogicalDrive{
					id:             "1",
					physicalDrives: make(map[string]*hpssaPhysicalDrive),
				}
				pd := &hpssaPhysicalDrive{location: "1I:1:1"}

				arr.logicalDrives["A"] = ld
				ld.physicalDrives["1I:1:1"] = pd
				cntrl.arrays["A"] = arr

				return map[string]*hpssaController{"0": cntrl}
			},
			validate: func(controllers map[string]*hpssaController) bool {
				cntrl := controllers["0"]
				arr := cntrl.arrays["A"]
				ld := arr.logicalDrives["A"]
				pd := ld.physicalDrives["1I:1:1"]

				return arr.cntrl == cntrl &&
					ld.cntrl == cntrl &&
					ld.arr == arr &&
					pd.cntrl == cntrl &&
					pd.arr == arr &&
					pd.ld == ld
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			controllers := tt.setup()
			updateHpssaHierarchy(controllers)
			if !tt.validate(controllers) {
				t.Errorf("hierarchy validation failed")
			}
		})
	}
}

func TestParseSsacliControllersInfo(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		shouldErr bool
		validate  func(map[string]*hpssaController) bool
	}{
		{
			name: "should parse single controller with array and logical drive",
			data: []byte(`Smart Array E208i-p SR Gen10 in Slot 0
   Serial Number: PDVCC0BQM5J1W3
   Controller Status: OK
   Cache Board Present: True
   Cache Status: OK
   Number of Ports: 8
   
   Array: A
      Interface Type: SAS
      Status: OK
      
      Logical Drive: 1
         Size: 1000 MB
         Status: OK
         
         physicaldrive
         
         Status: OK
`),
			shouldErr: false,
			validate: func(controllers map[string]*hpssaController) bool {
				if len(controllers) != 1 {
					return false
				}
				cntrl, ok := controllers["0"]
				if !ok {
					return false
				}
				return cntrl.model == "E208i-p" &&
					cntrl.serialNumber == "PDVCC0BQM5J1W3" &&
					cntrl.controllerStatus == "OK" &&
					len(cntrl.arrays) == 1
			},
		},
		{
			name: "should parse unassigned drives",
			data: []byte(`Smart Array E208i in Slot 0
   Serial Number: ABC123
   
   Unassigned
   
   physicaldrive 1I:1:2
   
   Status: Failed
`),
			shouldErr: false,
			validate: func(controllers map[string]*hpssaController) bool {
				cntrl, ok := controllers["0"]
				if !ok {
					return false
				}
				pd, ok := cntrl.unassignedDrives["1I:1:2"]
				return ok && pd.status == "Failed" && pd.cntrl == cntrl
			},
		},
		{
			name: "should error when unassigned drive but controller is nil",
			data: []byte(`Unassigned

physicaldrive 1I:1:2

Status: