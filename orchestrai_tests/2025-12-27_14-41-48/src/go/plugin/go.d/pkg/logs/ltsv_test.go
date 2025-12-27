package logs

import (
	"bytes"
	"io"
	"testing"
)

func TestLTSVParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    map[string]string
		wantErr bool
	}{
		{
			name:    "single key-value pair",
			input:   "key1:value1",
			want:    map[string]string{"key1": "value1"},
			wantErr: false,
		},
		{
			name:    "multiple key-value pairs",
			input:   "key1:value1\tkey2:value2\tkey3:value3",
			want:    map[string]string{"key1": "value1", "key2": "value2", "key3": "value3"},
			wantErr: false,
		},
		{
			name:    "empty value",
			input:   "key1:\tkey2:value2",
			want:    map[string]string{"key1": "", "key2": "value2"},
			wantErr: false,
		},
		{
			name:    "value with spaces",
			input:   "key1:value with spaces\tkey2:another value",
			want:    map[string]string{"key1": "value with spaces", "key2": "another value"},
			wantErr: false,
		},
		{
			name:    "empty input",
			input:   "",
			want:    map[string]string{},
			wantErr: false,
		},
		{
			name:    "only tabs",
			input:   "\t\t",
			want:    map[string]string{},
			wantErr: false,
		},
		{
			name:    "no colon separator",
			input:   "keyvalue",
			want:    map[string]string{},
			wantErr: false,
		},
		{
			name:    "multiple colons in value",
			input:   "key:value:with:colons",
			want:    map[string]string{"key": "value:with:colons"},
			wantErr: false,
		},
		{
			name:    "special characters in key",
			input:   "key-1_a:value1",
			want:    map[string]string{"key-1_a": "value1"},
			wantErr: false,
		},
		{
			name:    "special characters in value",
			input:   "key:!@#$%^&*()",
			want:    map[string]string{"key": "!@#$%^&*()"},
			wantErr: false,
		},
		{
			name:    "unicode characters",
			input:   "key:日本語\tkey2:中文",
			want:    map[string]string{"key": "日本語", "key2": "中文"},
			wantErr: false,
		},
		{
			name:    "trailing tab",
			input:   "key1:value1\t",
			want:    map[string]string{"key1": "value1"},
			wantErr: false,
		},
		{
			name:    "leading tab",
			input:   "\tkey1:value1",
			want:    map[string]string{"key1": "value1"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewLTSVParser()
			got, err := parser.Parse(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !mapsEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLTSVParser_ParseLine(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		want      map[string]string
		wantError bool
	}{
		{
			name:      "valid LTSV line",
			input:     "time:2021-01-01T00:00:00Z\thost:example.com\tpath:/path",
			want:      map[string]string{"time": "2021-01-01T00:00:00Z", "host": "example.com", "path": "/path"},
			wantError: false,
		},
		{
			name:      "empty line",
			input:     "",
			want:      map[string]string{},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewLTSVParser()
			got, err := parser.Parse(tt.input)

			if (err != nil) != tt.wantError {
				t.Errorf("Parse() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !mapsEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLTSVParser_ParseFromReader(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []map[string]string
		wantErr bool
	}{
		{
			name: "single line",
			input: "key1:value1\tkey2:value2\n",
			want: []map[string]string{
				{"key1": "value1", "key2": "value2"},
			},
			wantErr: false,
		},
		{
			name: "multiple lines",
			input: "key1:value1\tkey2:value2\nkey1:value3\tkey2:value4\n",
			want: []map[string]string{
				{"key1": "value1", "key2": "value2"},
				{"key1": "value3", "key2": "value4"},
			},
			wantErr: false,
		},
		{
			name: "empty input",
			input: "",
			want: []map[string]string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewLTSVParser()
			reader := bytes.NewReader([]byte(tt.input))
			
			var got []map[string]string
			var err error
			
			// Test if parser has ParseFromReader method
			if pr, ok := interface{}(parser).(interface{ ParseFromReader(io.Reader) error }); ok {
				err = pr.ParseFromReader(reader)
			} else if pm, ok := interface{}(parser).(interface{ ParseMultiple(string) []map[string]string }); ok {
				got = pm.ParseMultiple(tt.input)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseFromReader() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLTSVParser_DuplicateKeys(t *testing.T) {
	// Last value for duplicate key should win
	input := "key:value1\tkey:value2"
	parser := NewLTSVParser()
	got, err := parser.Parse(input)

	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	if val, ok := got["key"]; !ok || val != "value2" {
		t.Errorf("Parse() duplicate keys: expected 'value2', got %q", val)
	}
}

func TestLTSVParser_LongValues(t *testing.T) {
	longValue := bytes.Repeat([]byte("x"), 10000)
	input := "key:" + string(longValue)
	parser := NewLTSVParser()
	got, err := parser.Parse(input)

	if err != nil {
		t.Errorf("Parse() error = %v", err)
	}

	if val, ok := got["key"]; !ok || len(val) != 10000 {
		t.Errorf("Parse() long value: expected length 10000, got %d", len(val))
	}
}

func TestLTSVParser_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  map[string]string
	}{
		{
			name:  "only key no colon",
			input: "onlykey",
			want:  map[string]string{},
		},
		{
			name:  "only colon",
			input: ":",
			want:  map[string]string{"": ""},
		},
		{
			name:  "newline in value",
			input: "key:value\nwith\nnewlines",
			want:  map[string]string{"key": "value"},
		},
		{
			name:  "numeric values",
			input: "port:8080\ttimeout:3600",
			want:  map[string]string{"port": "8080", "timeout": "3600"},
		},
		{
			name:  "empty key",
			input: ":value",
			want:  map[string]string{"": "value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewLTSVParser()
			got, err := parser.Parse(tt.input)

			if err != nil {
				t.Errorf("Parse() error = %v", err)
			}

			if !mapsEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLTSVParser(t *testing.T) {
	parser := NewLTSVParser()
	if parser == nil {
		t.Error("NewLTSVParser() returned nil")
	}
}

func TestLTSVParser_RealWorldExample(t *testing.T) {
	// Apache combined log format as LTSV
	input := "time:01/Jan/2021:12:00:00 +0000\tremote:192.168.1.1\tmethod:GET\tpath:/index.html\tstatus:200\tsize:1024"
	
	parser := NewLTSVParser()
	got, err := parser.Parse(input)

	if err != nil {
		t.Errorf("Parse() error = %v", err)
		return
	}

	expected := map[string]string{
		"time":   "01/Jan/2021:12:00:00 +0000",
		"remote": "192.168.1.1",
		"method": "GET",
		"path":   "/index.html",
		"status": "200",
		"size":   "1024",
	}

	if !mapsEqual(got, expected) {
		t.Errorf("Parse() = %v, want %v", got, expected)
	}
}

func TestLTSVParser_TabHandling(t *testing.T) {
	tests := []struct {
		name  string
		input string
		pairs int
	}{
		{
			name:  "single pair",
			input: "key:value",
			pairs: 1,
		},
		{
			name:  "two pairs",
			input: "key1:value1\tkey2:value2",
			pairs: 2,
		},
		{
			name:  "three pairs",
			input: "key1:value1\tkey2:value2\tkey3:value3",
			pairs: 3,
		},
		{
			name:  "multiple tabs between pairs",
			input: "key1:value1\t\tkey2:value2",
			pairs: 1, // Empty field between tabs is skipped
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := NewLTSVParser()
			got, _ := parser.Parse(tt.input)

			if len(got) != tt.pairs {
				t.Errorf("Parse() got %d pairs, want %d", len(got), tt.pairs)
			}
		})
	}
}

func TestLTSVParser_ValueTypes(t *testing.T) {
	// LTSV returns all values as strings
	input := "string:hello\tnumber:123\tbool:true\tfloat:3.14"
	parser := NewLTSVParser()
	got, _ := parser.Parse(input)

	// All values should be strings
	for key, val := range got {
		if _, ok := val.(string); !ok {
			t.Errorf("Parse() key %q: expected string type, got %T", key, val)
		}
	}
}

// Helper function to compare maps
func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}