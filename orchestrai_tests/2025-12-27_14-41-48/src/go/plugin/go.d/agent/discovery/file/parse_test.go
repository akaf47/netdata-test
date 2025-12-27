package file

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseConfigFile(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		shouldErr bool
		expected  interface{}
	}{
		{
			name:      "empty file",
			content:   "",
			shouldErr: false,
			expected:  nil,
		},
		{
			name:      "single line",
			content:   "localhost",
			shouldErr: false,
			expected:  "localhost",
		},
		{
			name:      "multiple lines",
			content:   "localhost\n192.168.1.1\n10.0.0.1",
			shouldErr: false,
			expected:  []string{"localhost", "192.168.1.1", "10.0.0.1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseConfigFile(tt.content)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				if tt.expected == nil {
					assert.Nil(t, result)
				}
			}
		})
	}
}

func TestParseConfigWithWhitespace(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		shouldErr bool
	}{
		{
			name:      "leading spaces",
			content:   "  localhost",
			shouldErr: false,
		},
		{
			name:      "trailing spaces",
			content:   "localhost  ",
			shouldErr: false,
		},
		{
			name:      "tabs",
			content:   "\tlocalhost\t",
			shouldErr: false,
		},
		{
			name:      "empty lines",
			content:   "localhost\n\n192.168.1.1",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseConfigFile(tt.content)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseConfigWithComments(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		shouldErr bool
	}{
		{
			name:      "hash comment",
			content:   "# this is a comment\nlocalhost",
			shouldErr: false,
		},
		{
			name:      "inline comment",
			content:   "localhost # comment here",
			shouldErr: false,
		},
		{
			name:      "multiple comments",
			content:   "# comment 1\nlocalhost\n# comment 2\n192.168.1.1",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseConfigFile(tt.content)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseConfigIPAddresses(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		shouldErr bool
	}{
		{
			name:      "IPv4 address",
			content:   "192.168.1.1",
			shouldErr: false,
		},
		{
			name:      "IPv6 address",
			content:   "2001:db8::1",
			shouldErr: false,
		},
		{
			name:      "localhost",
			content:   "localhost",
			shouldErr: false,
		},
		{
			name:      "domain name",
			content:   "example.com",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseConfigFile(tt.content)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseConfigWithSpecialChars(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		shouldErr bool
	}{
		{
			name:      "colon separator",
			content:   "192.168.1.1:8080",
			shouldErr: false,
		},
		{
			name:      "underscore",
			content:   "my_server",
			shouldErr: false,
		},
		{
			name:      "dash",
			content:   "my-server",
			shouldErr: false,
		},
		{
			name:      "mixed special",
			content:   "srv-01_prod:9090",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseConfigFile(tt.content)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseConfigEmptyLines(t *testing.T) {
	content := "\n\n\nlocalhost\n\n\n"
	result, err := parseConfigFile(content)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestParseConfigSingleEmptyLine(t *testing.T) {
	content := "\n"
	_, err := parseConfigFile(content)
	assert.NoError(t, err)
}

func TestParseConfigLongContent(t *testing.T) {
	content := "localhost\n"
	for i := 0; i < 1000; i++ {
		content += fmt.Sprintf("192.168.1.%d\n", (i%254)+1)
	}

	result, err := parseConfigFile(content)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestParseConfigDuplicateLines(t *testing.T) {
	content := "localhost\nlocalhost\n192.168.1.1\n192.168.1.1"
	_, err := parseConfigFile(content)
	assert.NoError(t, err)
}

func TestParseConfigMixedCase(t *testing.T) {
	content := "LocalHost\nlOcAlHoSt\nLOCALHOST"
	_, err := parseConfigFile(content)
	assert.NoError(t, err)
}

func TestParseConfigWithNewlines(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "CRLF line endings",
			content: "localhost\r\n192.168.1.1\r\n",
		},
		{
			name:    "LF line endings",
			content: "localhost\n192.168.1.1\n",
		},
		{
			name:    "mixed line endings",
			content: "localhost\r\n192.168.1.1\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseConfigFile(tt.content)
			assert.NoError(t, err)
		})
	}
}

func TestParseConfigValidation(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		shouldErr bool
	}{
		{
			name:      "valid single entry",
			content:   "example.com",
			shouldErr: false,
		},
		{
			name:      "valid multiple entries",
			content:   "example.com\nexample2.com\nexample3.com",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseConfigFile(tt.content)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseConfigContentLength(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		shouldErr bool
	}{
		{
			name:      "very short",
			content:   "a",
			shouldErr: false,
		},
		{
			name:      "medium",
			content:   "example.com",
			shouldErr: false,
		},
		{
			name:      "long single line",
			content:   "this-is-a-very-long-hostname-that-might-be-used-in-production.example.com",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseConfigFile(tt.content)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseConfigNumberOfLines(t *testing.T) {
	tests := []struct {
		name       string
		numLines   int
		shouldErr  bool
	}{
		{
			name:      "1 line",
			numLines:  1,
			shouldErr: false,
		},
		{
			name:      "10 lines",
			numLines:  10,
			shouldErr: false,
		},
		{
			name:      "100 lines",
			numLines:  100,
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := ""
			for i := 0; i < tt.numLines; i++ {
				content += fmt.Sprintf("host%d.example.com\n", i)
			}

			_, err := parseConfigFile(content)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseConfigEdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		shouldErr bool
	}{
		{
			name:      "only spaces",
			content:   "     ",
			shouldErr: false,
		},
		{
			name:      "only tabs",
			content:   "\t\t\t",
			shouldErr: false,
		},
		{
			name:      "null character",
			content:   "localhost\x00",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseConfigFile(tt.content)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseConfigConsistency(t *testing.T) {
	content := "localhost\n192.168.1.1"

	result1, err1 := parseConfigFile(content)
	result2, err2 := parseConfigFile(content)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.Equal(t, result1, result2)
}

func TestParseConfigNilInput(t *testing.T) {
	content := ""
	_, err := parseConfigFile(content)
	assert.NoError(t, err)
}

func TestParseConfigArrayResult(t *testing.T) {
	content := "host1\nhost2\nhost3"
	result, err := parseConfigFile(content)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestParseConfigStringsWithNumbers(t *testing.T) {
	content := "srv001\nsrv002\nsrv003"
	_, err := parseConfigFile(content)
	assert.NoError(t, err)
}

func TestParseConfigUnicodeContent(t *testing.T) {
	content := "localhost\nserver-中文\nhost-한글"
	_, err := parseConfigFile(content)
	assert.NoError(t, err)
}

func TestParseConfigBoundaryValues(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		shouldErr bool
	}{
		{
			name:      "minimal valid",
			content:   "x",
			shouldErr: false,
		},
		{
			name:      "port number",
			content:   "0",
			shouldErr: false,
		},
		{
			name:      "max port",
			content:   "65535",
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseConfigFile(tt.content)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}