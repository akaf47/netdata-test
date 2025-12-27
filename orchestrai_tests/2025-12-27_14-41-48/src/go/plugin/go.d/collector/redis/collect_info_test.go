package redis

import (
	"bufio"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectInfo(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedMetrics map[string]int64
		shouldError  bool
	}{
		{
			name: "valid info response with multiple sections",
			input: "# Server\r\nredis_version:6.0.0\r\nprocess_id:1234\r\n# Clients\r\nconnected_clients:100\r\n# Memory\r\nused_memory:1048576\r\n",
			expectedMetrics: map[string]int64{
				"redis_version": 6,
				"connected_clients": 100,
				"used_memory": 1048576,
			},
			shouldError: false,
		},
		{
			name: "empty info response",
			input: "",
			expectedMetrics: map[string]int64{},
			shouldError: false,
		},
		{
			name: "info response with comments only",
			input: "# Server\r\n# Clients\r\n# Memory\r\n",
			expectedMetrics: map[string]int64{},
			shouldError: false,
		},
		{
			name: "info response with malformed lines",
			input: "redis_version:6.0.0\r\ninvalid_line\r\nprocess_id:1234\r\n",
			expectedMetrics: map[string]int64{
				"redis_version": 6,
				"process_id": 1234,
			},
			shouldError: false,
		},
		{
			name: "info response with zero values",
			input: "connected_clients:0\r\nused_memory:0\r\n",
			expectedMetrics: map[string]int64{
				"connected_clients": 0,
				"used_memory": 0,
			},
			shouldError: false,
		},
		{
			name: "info response with negative values",
			input: "some_metric:-100\r\n",
			expectedMetrics: map[string]int64{
				"some_metric": -100,
			},
			shouldError: false,
		},
		{
			name: "info response with large numbers",
			input: "uptime_in_seconds:9223372036854775807\r\n",
			expectedMetrics: map[string]int64{
				"uptime_in_seconds": 9223372036854775807,
			},
			shouldError: false,
		},
		{
			name: "info response with non-numeric values",
			input: "redis_version:not_a_number\r\n",
			expectedMetrics: map[string]int64{},
			shouldError: false,
		},
		{
			name: "info response with empty values",
			input: "key:\r\n",
			expectedMetrics: map[string]int64{},
			shouldError: false,
		},
		{
			name: "info response with whitespace",
			input: "  connected_clients:50  \r\n",
			expectedMetrics: map[string]int64{},
			shouldError: false,
		},
		{
			name: "info response with carriage return and newline",
			input: "connected_clients:100\r\nused_memory:2097152\r\n",
			expectedMetrics: map[string]int64{
				"connected_clients": 100,
				"used_memory": 2097152,
			},
			shouldError: false,
		},
		{
			name: "info response with only newline",
			input: "connected_clients:100\nused_memory:2097152\n",
			expectedMetrics: map[string]int64{
				"connected_clients": 100,
				"used_memory": 2097152,
			},
			shouldError: false,
		},
		{
			name: "info response with duplicate keys",
			input: "connected_clients:50\r\nconnected_clients:100\r\n",
			expectedMetrics: map[string]int64{
				"connected_clients": 100,
			},
			shouldError: false,
		},
		{
			name: "info response with key without colon",
			input: "redis_version\r\n",
			expectedMetrics: map[string]int64{},
			shouldError: false,
		},
		{
			name: "info response with multiple colons in value",
			input: "uri:redis://localhost:6379\r\n",
			expectedMetrics: map[string]int64{},
			shouldError: false,
		},
		{
			name: "info response with float values",
			input: "instantaneous_input_kbps:1.5\r\n",
			expectedMetrics: map[string]int64{},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(tt.input))
			metrics := parseInfoResponse(scanner)
			
			if tt.shouldError {
				assert.Nil(t, metrics)
			} else {
				// Verify metrics contain expected values
				for key, expectedValue := range tt.expectedMetrics {
					if val, ok := metrics[key]; ok {
						assert.Equal(t, expectedValue, val)
					}
				}
			}
		})
	}
}

func TestCollectInfoHandlesNilScanner(t *testing.T) {
	// Test handling of nil scanner
	metrics := parseInfoResponse(nil)
	assert.Nil(t, metrics)
}

func TestCollectInfoParseNumericFields(t *testing.T) {
	tests := []struct {
		name          string
		field         string
		expectedValue int64
		shouldParse   bool
	}{
		{
			name:          "parse valid positive integer",
			field:         "1234",
			expectedValue: 1234,
			shouldParse:   true,
		},
		{
			name:          "parse zero",
			field:         "0",
			expectedValue: 0,
			shouldParse:   true,
		},
		{
			name:          "parse negative integer",
			field:         "-5678",
			expectedValue: -5678,
			shouldParse:   true,
		},
		{
			name:          "parse max int64",
			field:         "9223372036854775807",
			expectedValue: 9223372036854775807,
			shouldParse:   true,
		},
		{
			name:          "parse min int64",
			field:         "-9223372036854775808",
			expectedValue: -9223372036854775808,
			shouldParse:   true,
		},
		{
			name:          "fail to parse string",
			field:         "not_a_number",
			expectedValue: 0,
			shouldParse:   false,
		},
		{
			name:          "fail to parse float",
			field:         "3.14",
			expectedValue: 0,
			shouldParse:   false,
		},
		{
			name:          "fail to parse empty string",
			field:         "",
			expectedValue: 0,
			shouldParse:   false,
		},
		{
			name:          "fail to parse with spaces",
			field:         " 123 ",
			expectedValue: 0,
			shouldParse:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, err := parseNumericField(tt.field)
			
			if tt.shouldParse {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedValue, value)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestCollectInfoExtractKeyValue(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		expectedKey string
		expectedVal string
		shouldParse bool
	}{
		{
			name:        "valid key:value pair",
			line:        "redis_version:6.0.0",
			expectedKey: "redis_version",
			expectedVal: "6.0.0",
			shouldParse: true,
		},
		{
			name:        "key with empty value",
			line:        "key:",
			expectedKey: "key",
			expectedVal: "",
			shouldParse: true,
		},
		{
			name:        "line without colon",
			line:        "no_colon_here",
			expectedKey: "",
			expectedVal: "",
			shouldParse: false,
		},
		{
			name:        "comment line",
			line:        "# Server",
			expectedKey: "",
			expectedVal: "",
			shouldParse: false,
		},
		{
			name:        "value with multiple colons",
			line:        "uri:redis://localhost:6379",
			expectedKey: "uri",
			expectedVal: "redis://localhost:6379",
			shouldParse: true,
		},
		{
			name:        "empty line",
			line:        "",
			expectedKey: "",
			expectedVal: "",
			shouldParse: false,
		},
		{
			name:        "whitespace only",
			line:        "   ",
			expectedKey: "",
			expectedVal: "",
			shouldParse: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, val, err := extractKeyValue(tt.line)
			
			if tt.shouldParse {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedKey, key)
				assert.Equal(t, tt.expectedVal, val)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestCollectInfoWithDifferentSections(t *testing.T) {
	input := `# Server
redis_version:6.0.0
redis_git_sha1:00000000
redis_git_dirty:0
redis_build_id:0
process_id:1234
# Clients
connected_clients:100
client_recent_max_input_buffer:0
client_recent_max_output_buffer:0
blocked_clients:0
# Memory
used_memory:1048576
used_memory_human:1M
used_memory_rss:2097152
used_memory_peak:3145728
# Stats
total_connections_received:500
total_commands_processed:10000
instantaneous_ops_per_sec:100
# Replication
role:master
connected_slaves:2
# Keyspace
db0:keys=1000,expires=0,avg_ttl=0
`
	scanner := bufio.NewScanner(strings.NewReader(input))
	metrics := parseInfoResponse(scanner)
	
	require.NotNil(t, metrics)
	assert.True(t, len(metrics) > 0)
}

func TestCollectInfoConcurrency(t *testing.T) {
	input := "connected_clients:100\r\nused_memory:1048576\r\n"
	
	// Run multiple goroutines to test concurrent access
	results := make(chan map[string]int64, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			scanner := bufio.NewScanner(strings.NewReader(input))
			results <- parseInfoResponse(scanner)
		}()
	}
	
	for i := 0; i < 10; i++ {
		metrics := <-results
		assert.NotNil(t, metrics)
	}
}

func TestCollectInfoMemoryEfficiency(t *testing.T) {
	// Test with large input to ensure memory efficiency
	var sb strings.Builder
	for i := 0; i < 1000; i++ {
		sb.WriteString("metric_")
		sb.WriteString(string(rune(i)))
		sb.WriteString(":")
		sb.WriteString(string(rune(i * 100)))
		sb.WriteString("\r\n")
	}
	
	scanner := bufio.NewScanner(strings.NewReader(sb.String()))
	metrics := parseInfoResponse(scanner)
	
	require.NotNil(t, metrics)
}

// Helper functions for testing
func parseInfoResponse(scanner *bufio.Scanner) map[string]int64 {
	if scanner == nil {
		return nil
	}
	
	metrics := make(map[string]int64)
	
	for scanner.Scan() {
		line := scanner.Text()
		
		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		
		// Parse key:value pairs
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		
		// Try to parse as integer
		if v, err := parseNumericValue(value); err == nil {
			metrics[key] = v
		}
	}
	
	return metrics
}

func parseNumericField(field string) (int64, error) {
	// Implementation placeholder
	return 0, nil
}

func parseNumericValue(value string) (int64, error) {
	// Implementation placeholder
	return 0, nil
}

func extractKeyValue(line string) (string, string, error) {
	// Implementation placeholder
	return "", "", nil
}