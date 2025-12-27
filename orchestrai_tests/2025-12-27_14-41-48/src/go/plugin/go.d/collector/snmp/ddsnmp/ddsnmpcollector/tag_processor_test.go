package ddsnmpcollector

import (
	"testing"
)

// TestProcessTags tests the tag processing functionality
func TestProcessTags(t *testing.T) {
	tests := []struct {
		name      string
		input     map[string]string
		wantErr   bool
		assertion func(t *testing.T, result map[string]string)
	}{
		{
			name:  "empty tags",
			input: map[string]string{},
			assertion: func(t *testing.T, result map[string]string) {
				if len(result) != 0 {
					t.Errorf("expected empty result, got %v", result)
				}
			},
		},
		{
			name: "simple tags",
			input: map[string]string{
				"tag1": "value1",
				"tag2": "value2",
			},
			assertion: func(t *testing.T, result map[string]string) {
				if result["tag1"] != "value1" {
					t.Errorf("expected tag1=value1, got %v", result["tag1"])
				}
				if result["tag2"] != "value2" {
					t.Errorf("expected tag2=value2, got %v", result["tag2"])
				}
			},
		},
		{
			name: "tags with special characters",
			input: map[string]string{
				"tag-name": "value-with-dash",
				"tag_name": "value_with_underscore",
			},
			assertion: func(t *testing.T, result map[string]string) {
				if result["tag-name"] != "value-with-dash" {
					t.Errorf("expected tag-name=value-with-dash, got %v", result["tag-name"])
				}
				if result["tag_name"] != "value_with_underscore" {
					t.Errorf("expected tag_name=value_with_underscore, got %v", result["tag_name"])
				}
			},
		},
		{
			name: "nil input",
			input: nil,
			assertion: func(t *testing.T, result map[string]string) {
				// Should handle gracefully
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessTags(tt.input)
			tt.assertion(t, result)
		})
	}
}

// TestValidateTag tests tag validation
func TestValidateTag(t *testing.T) {
	tests := []struct {
		name  string
		key   string
		value string
		want  bool
	}{
		{
			name:  "valid tag",
			key:   "host",
			value: "localhost",
			want:  true,
		},
		{
			name:  "empty key",
			key:   "",
			value: "value",
			want:  false,
		},
		{
			name:  "empty value",
			key:   "key",
			value: "",
			want:  false,
		},
		{
			name:  "both empty",
			key:   "",
			value: "",
			want:  false,
		},
		{
			name:  "key with spaces",
			key:   "key with spaces",
			value: "value",
			want:  false,
		},
		{
			name:  "value with spaces",
			key:   "key",
			value: "value with spaces",
			want:  true,
		},
		{
			name:  "numeric key",
			key:   "123",
			value: "value",
			want:  true,
		},
		{
			name:  "special characters in key",
			key:   "key@invalid",
			value: "value",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateTag(tt.key, tt.value)
			if got != tt.want {
				t.Errorf("ValidateTag(%q, %q) = %v, want %v", tt.key, tt.value, got, tt.want)
			}
		})
	}
}

// TestSanitizeTagValue tests tag value sanitization
func TestSanitizeTagValue(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "clean value",
			input: "clean_value",
			want:  "clean_value",
		},
		{
			name:  "value with leading spaces",
			input: "  value",
			want:  "value",
		},
		{
			name:  "value with trailing spaces",
			input: "value  ",
			want:  "value",
		},
		{
			name:  "value with spaces on both sides",
			input: "  value  ",
			want:  "value",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "only spaces",
			input: "   ",
			want:  "",
		},
		{
			name:  "value with internal spaces",
			input: "value with spaces",
			want:  "value with spaces",
		},
		{
			name:  "value with tabs",
			input: "\tvalue\t",
			want:  "value",
		},
		{
			name:  "value with newlines",
			input: "\nvalue\n",
			want:  "value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SanitizeTagValue(tt.input)
			if got != tt.want {
				t.Errorf("SanitizeTagValue(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// TestMergeTags tests tag merging functionality
func TestMergeTags(t *testing.T) {
	tests := []struct {
		name    string
		tags1   map[string]string
		tags2   map[string]string
		want    map[string]string
		wantErr bool
	}{
		{
			name:  "merge empty maps",
			tags1: map[string]string{},
			tags2: map[string]string{},
			want:  map[string]string{},
		},
		{
			name: "merge with first empty",
			tags1: map[string]string{},
			tags2: map[string]string{
				"key": "value",
			},
			want: map[string]string{
				"key": "value",
			},
		},
		{
			name: "merge with second empty",
			tags1: map[string]string{
				"key": "value",
			},
			tags2: map[string]string{},
			want: map[string]string{
				"key": "value",
			},
		},
		{
			name: "merge non-overlapping",
			tags1: map[string]string{
				"key1": "value1",
			},
			tags2: map[string]string{
				"key2": "value2",
			},
			want: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "merge with overlapping keys - second overrides",
			tags1: map[string]string{
				"key": "value1",
			},
			tags2: map[string]string{
				"key": "value2",
			},
			want: map[string]string{
				"key": "value2",
			},
		},
		{
			name: "merge with nil first",
			tags1: nil,
			tags2: map[string]string{
				"key": "value",
			},
			want: map[string]string{
				"key": "value",
			},
		},
		{
			name: "merge with nil second",
			tags1: map[string]string{
				"key": "value",
			},
			tags2: nil,
			want: map[string]string{
				"key": "value",
			},
		},
		{
			name:  "merge both nil",
			tags1: nil,
			tags2: nil,
			want:  map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MergeTags(tt.tags1, tt.tags2)
			if (err != nil) != tt.wantErr {
				t.Errorf("MergeTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !mapEqual(got, tt.want) {
				t.Errorf("MergeTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestFilterTags tests tag filtering functionality
func TestFilterTags(t *testing.T) {
	tests := []struct {
		name     string
		tags     map[string]string
		filter   []string
		want     map[string]string
		wantErr  bool
	}{
		{
			name:   "empty tags",
			tags:   map[string]string{},
			filter: []string{"key1"},
			want:   map[string]string{},
		},
		{
			name: "empty filter",
			tags: map[string]string{
				"key1": "value1",
			},
			filter: []string{},
			want:   map[string]string{},
		},
		{
			name: "filter matches keys",
			tags: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			filter: []string{"key1", "key3"},
			want: map[string]string{
				"key1": "value1",
				"key3": "value3",
			},
		},
		{
			name: "filter with no matches",
			tags: map[string]string{
				"key1": "value1",
			},
			filter: []string{"nonexistent"},
			want:   map[string]string{},
		},
		{
			name:   "nil tags",
			tags:   nil,
			filter: []string{"key1"},
			want:   map[string]string{},
		},
		{
			name: "nil filter",
			tags: map[string]string{
				"key1": "value1",
			},
			filter: nil,
			want:   map[string]string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FilterTags(tt.tags, tt.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("FilterTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !mapEqual(got, tt.want) {
				t.Errorf("FilterTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestExcludeTags tests tag exclusion functionality
func TestExcludeTags(t *testing.T) {
	tests := []struct {
		name     string
		tags     map[string]string
		exclude  []string
		want     map[string]string
		wantErr  bool
	}{
		{
			name:    "empty tags",
			tags:    map[string]string{},
			exclude: []string{"key1"},
			want:    map[string]string{},
		},
		{
			name: "empty exclude",
			tags: map[string]string{
				"key1": "value1",
			},
			exclude: []string{},
			want: map[string]string{
				"key1": "value1",
			},
		},
		{
			name: "exclude existing keys",
			tags: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
			exclude: []string{"key1", "key3"},
			want: map[string]string{
				"key2": "value2",
			},
		},
		{
			name: "exclude with no matches",
			tags: map[string]string{
				"key1": "value1",
			},
			exclude: []string{"nonexistent"},
			want: map[string]string{
				"key1": "value1",
			},
		},
		{
			name:    "nil tags",
			tags:    nil,
			exclude: []string{"key1"},
			want:    map[string]string{},
		},
		{
			name: "nil exclude",
			tags: map[string]string{
				"key1": "value1",
			},
			exclude: nil,
			want: map[string]string{
				"key1": "value1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExcludeTags(tt.tags, tt.exclude)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExcludeTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !mapEqual(got, tt.want) {
				t.Errorf("ExcludeTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestDuplicateTagKey tests handling of duplicate tag keys
func TestDuplicateTagKey(t *testing.T) {
	tests := []struct {
		name  string
		tags  map[string]string
		key   string
		want  bool
	}{
		{
			name: "key exists",
			tags: map[string]string{
				"key1": "value1",
			},
			key:  "key1",
			want: true,
		},
		{
			name: "key not exists",
			tags: map[string]string{
				"key1": "value1",
			},
			key:  "key2",
			want: false,
		},
		{
			name:  "empty tags",
			tags:  map[string]string{},
			key:   "key1",
			want:  false,
		},
		{
			name:  "nil tags",
			tags:  nil,
			key:   "key1",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DuplicateTagKey(tt.tags, tt.key)
			if got != tt.want {
				t.Errorf("DuplicateTagKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTransformTagKey tests tag key transformation
func TestTransformTagKey(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "lowercase letters",
			input: "Key",
			want:  "key",
		},
		{
			name:  "uppercase letters",
			input: "KEY",
			want:  "key",
		},
		{
			name:  "mixed case",
			input: "MyKey",
			want:  "mykey",
		},
		{
			name:  "with underscores",
			input: "My_Key",
			want:  "my_key",
		},
		{
			name:  "with dashes",
			input: "My-Key",
			want:  "my-key",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
		{
			name:  "single char",
			input: "A",
			want:  "a",
		},
		{
			name:  "with numbers",
			input: "Key123",
			want:  "key123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TransformTagKey(tt.input)
			if got != tt.want {
				t.Errorf("TransformTagKey(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

// Helper function to compare maps
func mapEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

// TestProcessTagsWithInvalidInput tests edge cases
func TestProcessTagsWithInvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input map[string]string
	}{
		{
			name: "large map",
			input: func() map[string]string {
				m := make(map[string]string)
				for i := 0; i < 1000; i++ {
					m[string(rune(i))] = "value"
				}
				return m
			}(),
		},
		{
			name: "unicode values",
			input: map[string]string{
				"key": "日本語",
			},
		},
		{
			name: "empty string values",
			input: map[string]string{
				"key": "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProcessTags(tt.input)
			if result == nil {
				t.Error("ProcessTags returned nil for valid input")
			}
		})
	}
}

// TestProcessTagsConcurrency tests thread safety
func TestProcessTagsConcurrency(t *testing.T) {
	tags := map[string]string{"key": "value"}
	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			_ = ProcessTags(tags)
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}
}