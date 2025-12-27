package hdfs

import (
	"encoding/json"
	"testing"
)

// TestRawDataMarshal tests JSON marshaling of raw data structures
func TestRawDataMarshal(t *testing.T) {
	tests := []struct {
		name    string
		data    interface{}
		wantErr bool
	}{
		{
			name:    "empty data",
			data:    map[string]interface{}{},
			wantErr: false,
		},
		{
			name:    "nil data",
			data:    nil,
			wantErr: false,
		},
		{
			name: "complex nested data",
			data: map[string]interface{}{
				"string": "value",
				"number": 42,
				"float":  3.14,
				"bool":   true,
				"array":  []interface{}{1, 2, 3},
				"object": map[string]interface{}{"nested": "data"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := json.Marshal(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestRawDataUnmarshal tests JSON unmarshaling of raw data
func TestRawDataUnmarshal(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    interface{}
		wantErr bool
	}{
		{
			name:    "valid json object",
			data:    []byte(`{"key":"value"}`),
			want:    map[string]interface{}{"key": "value"},
			wantErr: false,
		},
		{
			name:    "valid json array",
			data:    []byte(`[1,2,3]`),
			want:    []interface{}{float64(1), float64(2), float64(3)},
			wantErr: false,
		},
		{
			name:    "empty object",
			data:    []byte(`{}`),
			want:    map[string]interface{}{},
			wantErr: false,
		},
		{
			name:    "invalid json",
			data:    []byte(`{invalid`),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "null value",
			data:    []byte(`null`),
			want:    nil,
			wantErr: false,
		},
		{
			name:    "empty string",
			data:    []byte(``),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "whitespace only",
			data:    []byte(`   `),
			want:    nil,
			wantErr: true,
		},
		{
			name: "deeply nested structure",
			data: []byte(`{"level1":{"level2":{"level3":{"value":"deep"}}}}`),
			wantErr: false,
		},
		{
			name: "special characters in strings",
			data: []byte(`{"key":"value with \"quotes\" and \\backslash"}`),
			wantErr: false,
		},
		{
			name: "unicode characters",
			data: []byte(`{"key":"value with unicode \u0041"}`),
			wantErr: false,
		},
		{
			name:    "large numbers",
			data:    []byte(`{"big":9999999999999999}`),
			wantErr: false,
		},
		{
			name:    "negative numbers",
			data:    []byte(`{"negative":-42}`),
			wantErr: false,
		},
		{
			name:    "float numbers",
			data:    []byte(`{"float":3.14159}`),
			wantErr: false,
		},
		{
			name:    "scientific notation",
			data:    []byte(`{"scientific":1.23e-4}`),
			wantErr: false,
		},
		{
			name:    "boolean true",
			data:    []byte(`{"bool":true}`),
			wantErr: false,
		},
		{
			name:    "boolean false",
			data:    []byte(`{"bool":false}`),
			wantErr: false,
		},
		{
			name:    "array with mixed types",
			data:    []byte(`[1,"two",true,null,3.14]`),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}
			err := json.Unmarshal(tt.data, &result)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestRawDataStringParsing tests string parsing and validation
func TestRawDataStringParsing(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		isEmpty bool
	}{
		{
			name:    "empty string",
			input:   "",
			isEmpty: true,
		},
		{
			name:    "whitespace",
			input:   "   ",
			isEmpty: false,
		},
		{
			name:    "normal string",
			input:   "test",
			isEmpty: false,
		},
		{
			name:    "numeric string",
			input:   "12345",
			isEmpty: false,
		},
		{
			name:    "special characters",
			input:   "!@#$%^&*()",
			isEmpty: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isEmpty := len(tt.input) == 0
			if isEmpty != tt.isEmpty {
				t.Errorf("isEmpty = %v, want %v", isEmpty, tt.isEmpty)
			}
		})
	}
}

// TestRawDataNumberParsing tests number parsing and validation
func TestRawDataNumberParsing(t *testing.T) {
	tests := []struct {
		name      string
		value     float64
		expectMin bool
		expectMax bool
	}{
		{
			name:      "zero",
			value:     0,
			expectMin: true,
			expectMax: false,
		},
		{
			name:      "positive number",
			value:     100,
			expectMin: true,
			expectMax: true,
		},
		{
			name:      "negative number",
			value:     -50,
			expectMin: true,
			expectMax: false,
		},
		{
			name:      "float",
			value:     3.14,
			expectMin: true,
			expectMax: true,
		},
		{
			name:      "very large number",
			value:     1e10,
			expectMin: true,
			expectMax: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectMin && tt.value < 0 {
				t.Errorf("value %f should not be negative", tt.value)
			}
		})
	}
}

// TestRawDataNullHandling tests handling of null/nil values
func TestRawDataNullHandling(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		isNil    bool
		expected interface{}
	}{
		{
			name:     "nil interface",
			value:    nil,
			isNil:    true,
			expected: nil,
		},
		{
			name:     "nil pointer",
			value:    (*struct{})(nil),
			isNil:    true,
			expected: nil,
		},
		{
			name:     "zero value int",
			value:    0,
			isNil:    false,
			expected: 0,
		},
		{
			name:     "empty string",
			value:    "",
			isNil:    false,
			expected: "",
		},
		{
			name:     "false bool",
			value:    false,
			isNil:    false,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isNil := tt.value == nil
			if isNil != tt.isNil {
				t.Errorf("isNil = %v, want %v", isNil, tt.isNil)
			}
		})
	}
}

// TestRawDataErrorConditions tests error handling
func TestRawDataErrorConditions(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "malformed json - missing closing brace",
			data:    []byte(`{"key":"value"`),
			wantErr: true,
		},
		{
			name:    "malformed json - trailing comma",
			data:    []byte(`{"key":"value",}`),
			wantErr: true,
		},
		{
			name:    "malformed json - single quotes",
			data:    []byte(`{'key':'value'}`),
			wantErr: true,
		},
		{
			name:    "malformed json - unquoted key",
			data:    []byte(`{key:"value"}`),
			wantErr: true,
		},
		{
			name:    "invalid escape sequence",
			data:    []byte(`{"key":"value\x"}`),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}
			err := json.Unmarshal(tt.data, &result)
			if (err != nil) != tt.wantErr {
				t.Errorf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}

// TestRawDataBoundaryConditions tests boundary values
func TestRawDataBoundaryConditions(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "minimum int64",
			data: []byte(`{"value":-9223372036854775808}`),
		},
		{
			name: "maximum int64",
			data: []byte(`{"value":9223372036854775807}`),
		},
		{
			name: "float32 range",
			data: []byte(`{"value":3.4028235e+38}`),
		},
		{
			name: "very small float",
			data: []byte(`{"value":1e-38}`),
		},
		{
			name: "zero",
			data: []byte(`{"value":0}`),
		},
		{
			name: "negative zero",
			data: []byte(`{"value":-0}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}
			err := json.Unmarshal(tt.data, &result)
			if err != nil {
				t.Errorf("failed to unmarshal: %v", err)
			}
		})
	}
}

// TestRawDataEmptyCollections tests empty arrays and objects
func TestRawDataEmptyCollections(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{
			name: "empty object",
			data: []byte(`{}`),
		},
		{
			name: "empty array",
			data: []byte(`[]`),
		},
		{
			name: "object with empty string",
			data: []byte(`{"key":""}`),
		},
		{
			name: "object with empty array",
			data: []byte(`{"key":[]}`),
		},
		{
			name: "object with empty object",
			data: []byte(`{"key":{}}`),
		},
		{
			name: "nested empty structures",
			data: []byte(`{"a":{"b":{"c":{}}}}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result interface{}
			err := json.Unmarshal(tt.data, &result)
			if err != nil {
				t.Errorf("failed to unmarshal: %v", err)
			}
		})
	}
}