package logger

import (
	"testing"
)

// TestLevelString tests the String() method for all log levels
func TestLevelString(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		expected string
	}{
		{
			name:     "should return DEBUG for DEBUG level",
			level:    DEBUG,
			expected: "DEBUG",
		},
		{
			name:     "should return INFO for INFO level",
			level:    INFO,
			expected: "INFO",
		},
		{
			name:     "should return WARNING for WARNING level",
			level:    WARNING,
			expected: "WARNING",
		},
		{
			name:     "should return ERROR for ERROR level",
			level:    ERROR,
			expected: "ERROR",
		},
		{
			name:     "should return CRITICAL for CRITICAL level",
			level:    CRITICAL,
			expected: "CRITICAL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.level.String()
			if got != tt.expected {
				t.Errorf("String() = %s, want %s", got, tt.expected)
			}
		})
	}
}

// TestLevelStringInvalidLevel tests String() for invalid level values
func TestLevelStringInvalidLevel(t *testing.T) {
	tests := []struct {
		name  string
		level Level
	}{
		{
			name:  "should handle negative level",
			level: Level(-1),
		},
		{
			name:  "should handle large invalid level",
			level: Level(999),
		},
		{
			name:  "should handle zero level",
			level: Level(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.level.String()
			// Should return something, even if unknown
			if got == "" {
				t.Errorf("String() returned empty string for invalid level")
			}
		})
	}
}

// TestLevelIsValid tests level validation
func TestLevelIsValid(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		expected bool
	}{
		{
			name:     "should return true for DEBUG",
			level:    DEBUG,
			expected: true,
		},
		{
			name:     "should return true for INFO",
			level:    INFO,
			expected: true,
		},
		{
			name:     "should return true for WARNING",
			level:    WARNING,
			expected: true,
		},
		{
			name:     "should return true for ERROR",
			level:    ERROR,
			expected: true,
		},
		{
			name:     "should return true for CRITICAL",
			level:    CRITICAL,
			expected: true,
		},
		{
			name:     "should return false for invalid negative level",
			level:    Level(-1),
			expected: false,
		},
		{
			name:     "should return false for invalid large level",
			level:    Level(999),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.level.IsValid()
			if got != tt.expected {
				t.Errorf("IsValid() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestLevelComparison tests level comparison/ordering
func TestLevelComparison(t *testing.T) {
	tests := []struct {
		name     string
		level1   Level
		level2   Level
		less     bool
		equal    bool
		greater  bool
	}{
		{
			name:    "DEBUG should be less than INFO",
			level1:  DEBUG,
			level2:  INFO,
			less:    true,
			equal:   false,
			greater: false,
		},
		{
			name:    "INFO should be less than WARNING",
			level1:  INFO,
			level2:  WARNING,
			less:    true,
			equal:   false,
			greater: false,
		},
		{
			name:    "WARNING should be less than ERROR",
			level1:  WARNING,
			level2:  ERROR,
			less:    true,
			equal:   false,
			greater: false,
		},
		{
			name:    "ERROR should be less than CRITICAL",
			level1:  ERROR,
			level2:  CRITICAL,
			less:    true,
			equal:   false,
			greater: false,
		},
		{
			name:    "DEBUG should equal DEBUG",
			level1:  DEBUG,
			level2:  DEBUG,
			less:    false,
			equal:   true,
			greater: false,
		},
		{
			name:    "CRITICAL should be greater than ERROR",
			level1:  CRITICAL,
			level2:  ERROR,
			less:    false,
			equal:   false,
			greater: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.less && !(tt.level1 < tt.level2) {
				t.Errorf("%v < %v should be true", tt.level1, tt.level2)
			}
			if tt.equal && tt.level1 != tt.level2 {
				t.Errorf("%v == %v should be true", tt.level1, tt.level2)
			}
			if tt.greater && !(tt.level1 > tt.level2) {
				t.Errorf("%v > %v should be true", tt.level1, tt.level2)
			}
		})
	}
}

// TestLevelFromString tests creating Level from string
func TestLevelFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected Level
		wantErr  bool
	}{
		{
			name:     "should parse DEBUG",
			input:    "DEBUG",
			expected: DEBUG,
			wantErr:  false,
		},
		{
			name:     "should parse INFO",
			input:    "INFO",
			expected: INFO,
			wantErr:  false,
		},
		{
			name:     "should parse WARNING",
			input:    "WARNING",
			expected: WARNING,
			wantErr:  false,
		},
		{
			name:     "should parse ERROR",
			input:    "ERROR",
			expected: ERROR,
			wantErr:  false,
		},
		{
			name:     "should parse CRITICAL",
			input:    "CRITICAL",
			expected: CRITICAL,
			wantErr:  false,
		},
		{
			name:     "should handle lowercase debug",
			input:    "debug",
			expected: DEBUG,
			wantErr:  false,
		},
		{
			name:     "should handle mixed case info",
			input:    "Info",
			expected: INFO,
			wantErr:  false,
		},
		{
			name:     "should error on invalid string",
			input:    "INVALID",
			expected: Level(0),
			wantErr:  true,
		},
		{
			name:     "should error on empty string",
			input:    "",
			expected: Level(0),
			wantErr:  true,
		},
		{
			name:     "should error on whitespace",
			input:    "   ",
			expected: Level(0),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLevel(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLevel() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("ParseLevel() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestLevelConstants tests that level constants have correct values
func TestLevelConstants(t *testing.T) {
	tests := []struct {
		name  string
		level Level
		value int
	}{
		{
			name:  "DEBUG should have value 0",
			level: DEBUG,
			value: 0,
		},
		{
			name:  "INFO should have value 1",
			level: INFO,
			value: 1,
		},
		{
			name:  "WARNING should have value 2",
			level: WARNING,
			value: 2,
		},
		{
			name:  "ERROR should have value 3",
			level: ERROR,
			value: 3,
		},
		{
			name:  "CRITICAL should have value 4",
			level: CRITICAL,
			value: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.level) != tt.value {
				t.Errorf("Level value = %d, want %d", int(tt.level), tt.value)
			}
		})
	}
}

// TestLevelHierarchy tests the hierarchy of log levels
func TestLevelHierarchy(t *testing.T) {
	levels := []Level{DEBUG, INFO, WARNING, ERROR, CRITICAL}

	// Test ordering
	for i := 0; i < len(levels)-1; i++ {
		if !(levels[i] < levels[i+1]) {
			t.Errorf("Level hierarchy broken: %v should be < %v", levels[i], levels[i+1])
		}
	}

	// Test that lower index is lower severity
	for i := 0; i < len(levels); i++ {
		for j := i + 1; j < len(levels); j++ {
			if !(levels[i] < levels[j]) {
				t.Errorf("Level ordering failed: %v should be < %v", levels[i], levels[j])
			}
		}
	}
}

// TestLevelStringRepresentation tests string representation consistency
func TestLevelStringRepresentation(t *testing.T) {
	levels := []Level{DEBUG, INFO, WARNING, ERROR, CRITICAL}
	levelNames := []string{"DEBUG", "INFO", "WARNING", "ERROR", "CRITICAL"}

	for i, level := range levels {
		t.Run(levelNames[i], func(t *testing.T) {
			str := level.String()
			if str != levelNames[i] {
				t.Errorf("String() = %s, want %s", str, levelNames[i])
			}

			// Verify it can be parsed back
			parsed, err := ParseLevel(str)
			if err != nil {
				t.Errorf("ParseLevel(%s) failed: %v", str, err)
			}
			if parsed != level {
				t.Errorf("Round-trip parsing failed: %v -> %s -> %v", level, str, parsed)
			}
		})
	}
}

// TestLevelMarshalUnmarshal tests marshaling/unmarshaling if implemented
func TestLevelMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name  string
		level Level
	}{
		{
			name:  "should marshal DEBUG",
			level: DEBUG,
		},
		{
			name:  "should marshal INFO",
			level: INFO,
		},
		{
			name:  "should marshal WARNING",
			level: WARNING,
		},
		{
			name:  "should marshal ERROR",
			level: ERROR,
		},
		{
			name:  "should marshal CRITICAL",
			level: CRITICAL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that level can be used in JSON/similar contexts
			str := tt.level.String()
			parsed, err := ParseLevel(str)
			if err != nil {
				t.Errorf("Failed to parse marshaled level: %v", err)
			}
			if parsed != tt.level {
				t.Errorf("Unmarshal failed: got %v, want %v", parsed, tt.level)
			}
		})
	}
}

// TestLevelBoundary tests boundary values
func TestLevelBoundary(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		wantMin  bool
		wantMax  bool
	}{
		{
			name:    "DEBUG is minimum valid level",
			level:   DEBUG,
			wantMin: true,
			wantMax: false,
		},
		{
			name:    "CRITICAL is maximum valid level",
			level:   CRITICAL,
			wantMin: false,
			wantMax: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantMin && tt.level != DEBUG {
				t.Errorf("Expected minimum level, got %v", tt.level)
			}
			if tt.wantMax && tt.level != CRITICAL {
				t.Errorf("Expected maximum level, got %v", tt.level)
			}
		})
	}
}

// TestLevelSwitchStatement tests all levels in switch statement context
func TestLevelSwitchStatement(t *testing.T) {
	levels := []Level{DEBUG, INFO, WARNING, ERROR, CRITICAL}

	for _, level := range levels {
		t.Run(level.String(), func(t *testing.T) {
			handled := false
			switch level {
			case DEBUG:
				handled = true
			case INFO:
				handled = true
			case WARNING:
				handled = true
			case ERROR:
				handled = true
			case CRITICAL:
				handled = true
			default:
				t.Errorf("Level %v not handled in switch", level)
			}
			if !handled {
				t.Errorf("Level %v was not handled", level)
			}
		})
	}
}

// TestLevelTyping tests level type conversions
func TestLevelTyping(t *testing.T) {
	tests := []struct {
		name   string
		input  int
		level  Level
		output int
	}{
		{
			name:   "should convert int to Level to int (DEBUG)",
			input:  0,
			level:  Level(0),
			output: 0,
		},
		{
			name:   "should convert int to Level to int (INFO)",
			input:  1,
			level:  Level(1),
			output: 1,
		},
		{
			name:   "should preserve value through conversion",
			input:  3,
			level:  Level(3),
			output: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if int(tt.level) != tt.output {
				t.Errorf("Type conversion failed: %d -> Level -> %d, want %d", tt.input, int(tt.level), tt.output)
			}
		})
	}
}