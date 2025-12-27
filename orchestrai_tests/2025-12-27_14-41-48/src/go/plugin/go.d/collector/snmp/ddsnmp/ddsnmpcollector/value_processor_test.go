package ddsnmpcollector

import (
	"fmt"
	"testing"
)

// Test ProcessValue function with various input types
func TestProcessValue_Integer(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "positive integer",
			input:    int64(42),
			expected: int64(42),
		},
		{
			name:     "negative integer",
			input:    int64(-42),
			expected: int64(-42),
		},
		{
			name:     "zero integer",
			input:    int64(0),
			expected: int64(0),
		},
		{
			name:     "max int64",
			input:    int64(9223372036854775807),
			expected: int64(9223372036854775807),
		},
		{
			name:     "min int64",
			input:    int64(-9223372036854775808),
			expected: int64(-9223372036854775808),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessValue(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessValue(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestProcessValue_Float(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "positive float",
			input:    float64(3.14),
			expected: float64(3.14),
		},
		{
			name:     "negative float",
			input:    float64(-3.14),
			expected: float64(-3.14),
		},
		{
			name:     "zero float",
			input:    float64(0.0),
			expected: float64(0.0),
		},
		{
			name:     "very small float",
			input:    float64(0.0001),
			expected: float64(0.0001),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessValue(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessValue(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestProcessValue_String(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "simple string",
			input:    "test",
			expected: "test",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "numeric string",
			input:    "123",
			expected: "123",
		},
		{
			name:     "string with special chars",
			input:    "test!@#$%",
			expected: "test!@#$%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessValue(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessValue(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestProcessValue_Nil(t *testing.T) {
	result := ProcessValue(nil)
	if result != nil {
		t.Errorf("ProcessValue(nil) = %v, want nil", result)
	}
}

func TestProcessValue_Bool(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "true boolean",
			input:    true,
			expected: true,
		},
		{
			name:     "false boolean",
			input:    false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessValue(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessValue(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestProcessValue_Byte(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "byte array",
			input:    []byte{1, 2, 3},
			expected: []byte{1, 2, 3},
		},
		{
			name:     "empty byte array",
			input:    []byte{},
			expected: []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessValue(tt.input)
			if fmt.Sprintf("%v", result) != fmt.Sprintf("%v", tt.expected) {
				t.Errorf("ProcessValue(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestProcessValue_Uint(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "uint value",
			input:    uint(42),
			expected: uint(42),
		},
		{
			name:     "uint zero",
			input:    uint(0),
			expected: uint(0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessValue(tt.input)
			if result != tt.expected {
				t.Errorf("ProcessValue(%v) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestProcessValue_WithStringConversion(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "convert int to string",
			input:    int64(42),
			expected: "42",
		},
		{
			name:     "convert float to string",
			input:    float64(3.14),
			expected: "3.14",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessValue(tt.input)
			resultStr := fmt.Sprintf("%v", result)
			if resultStr != tt.expected {
				t.Errorf("ProcessValue(%v) string = %v, want %v", tt.input, resultStr, tt.expected)
			}
		})
	}
}

// Test error handling for invalid types
func TestProcessValue_InvalidType(t *testing.T) {
	complexValue := map[string]interface{}{"key": "value"}
	result := ProcessValue(complexValue)
	if result == nil {
		t.Error("ProcessValue should return the input for unrecognized types")
	}
}

// Benchmark test for ProcessValue
func BenchmarkProcessValue(b *testing.B) {
	value := int64(42)
	for i := 0; i < b.N; i++ {
		ProcessValue(value)
	}
}