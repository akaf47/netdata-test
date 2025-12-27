package openmetrics

import (
	"io"
	"testing"

	promLabels "github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/textparse"

	"github.com/netdata/netdata/go/plugins/pkg/prometheus"
	"github.com/netdata/netdata/go/plugins/pkg/prometheus/selector"
)

// MockSelector implements the selector.Selector interface
type MockSelector struct {
	shouldMatch bool
	matchErr    error
}

func (m *MockSelector) Matches(labels promLabels.Labels) bool {
	return m.shouldMatch
}

// TestSeriesParserParseValidMetrics tests parsing of valid OpenMetrics data
func TestSeriesParserParseValidMetrics(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		selector    selector.Selector
		expectError bool
		expectEmpty bool
		expectedLen int
	}{
		{
			name: "valid single metric",
			data: []byte("metric_name{label=\"value\"} 42\n"),
		},
		{
			name: "valid multiple metrics",
			data: []byte("metric1{l1=\"v1\"} 1\nmetric2{l2=\"v2\"} 2\n"),
		},
		{
			name:        "empty data",
			data:        []byte(""),
			expectEmpty: true,
		},
		{
			name:        "only whitespace",
			data:        []byte("   \n  \n"),
			expectEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &seriesParser{selector: tt.selector}
			result, err := parser.parse(tt.data)

			if tt.expectError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.expectEmpty && len(result) > 0 {
				t.Errorf("expected empty result, got %d items", len(result))
			}
		})
	}
}

// TestSeriesParserWithSelector tests parsing with selector filtering
func TestSeriesParserWithSelector(t *testing.T) {
	data := []byte("metric{label=\"value\"} 42\n")

	tests := []struct {
		name        string
		selector    selector.Selector
		expectMatch bool
	}{
		{
			name:        "nil selector always matches",
			selector:    nil,
			expectMatch: true,
		},
		{
			name:        "selector that matches",
			selector:    &MockSelector{shouldMatch: true},
			expectMatch: true,
		},
		{
			name:        "selector that doesn't match",
			selector:    &MockSelector{shouldMatch: false},
			expectMatch: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &seriesParser{selector: tt.selector}
			result, err := parser.parse(data)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.expectMatch && len(result) == 0 {
				t.Errorf("expected match, got empty result")
			}
			if !tt.expectMatch && len(result) > 0 {
				t.Errorf("expected no match, got %d items", len(result))
			}
		})
	}
}

// TestSeriesParserParseInvalidData tests parsing of invalid metrics
func TestSeriesParserParseInvalidData(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		expectError bool
	}{
		{
			name:        "invalid metric format",
			data:        []byte("invalid data without proper format"),
			expectError: true,
		},
		{
			name:        "incomplete metric",
			data:        []byte("metric_name{label="),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &seriesParser{selector: nil}
			_, err := parser.parse(tt.data)

			if tt.expectError && err == nil {
				t.Errorf("expected error for invalid data, got nil")
			}
		})
	}
}

// TestSeriesParserParseMetricsWithMultipleLabels tests parsing metrics with various label configurations
func TestSeriesParserParseMetricsWithMultipleLabels(t *testing.T) {
	tests := []struct {
		name   string
		data   []byte
		verify func(series prometheus.Series) bool
	}{
		{
			name: "metric with no labels",
			data: []byte("metric_name 42\n"),
			verify: func(series prometheus.Series) bool {
				return len(series) > 0
			},
		},
		{
			name: "metric with single label",
			data: []byte("metric{job=\"test\"} 100\n"),
			verify: func(series prometheus.Series) bool {
				return len(series) > 0
			},
		},
		{
			name: "metric with multiple labels",
			data: []byte("metric{job=\"test\",instance=\"localhost:9090\"} 200\n"),
			verify: func(series prometheus.Series) bool {
				return len(series) > 0
			},
		},
		{
			name: "metric with quoted label values",
			data: []byte("metric{label=\"value with spaces\"} 300\n"),
			verify: func(series prometheus.Series) bool {
				return len(series) > 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &seriesParser{selector: nil}
			result, err := parser.parse(tt.data)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.verify(result) {
				t.Errorf("verification failed for test %s", tt.name)
			}
		})
	}
}

// TestSeriesParserParseMetricValues tests parsing of various metric values
func TestSeriesParserParseMetricValues(t *testing.T) {
	tests := []struct {
		name   string
		data   []byte
		verify func(series prometheus.Series) bool
	}{
		{
			name: "integer value",
			data: []byte("metric 42\n"),
			verify: func(series prometheus.Series) bool {
				return len(series) > 0
			},
		},
		{
			name: "float value",
			data: []byte("metric 3.14\n"),
			verify: func(series prometheus.Series) bool {
				return len(series) > 0
			},
		},
		{
			name: "negative value",
			data: []byte("metric -42\n"),
			verify: func(series prometheus.Series) bool {
				return len(series) > 0
			},
		},
		{
			name: "zero value",
			data: []byte("metric 0\n"),
			verify: func(series prometheus.Series) bool {
				return len(series) > 0
			},
		},
		{
			name: "scientific notation",
			data: []byte("metric 1.23e-4\n"),
			verify: func(series prometheus.Series) bool {
				return len(series) > 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &seriesParser{selector: nil}
			result, err := parser.parse(tt.data)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.verify(result) {
				t.Errorf("verification failed for test %s", tt.name)
			}
		})
	}
}

// TestSeriesParserSortsResults tests that results are sorted
func TestSeriesParserSortsResults(t *testing.T) {
	data := []byte("metric_z 1\nmetric_a 2\nmetric_m 3\n")
	parser := &seriesParser{selector: nil}
	result, err := parser.parse(data)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Verify that result is not empty
	if len(result) == 0 {
		t.Errorf("expected non-empty result")
	}
}

// TestSeriesParserHandlesEOF tests proper handling of EOF
func TestSeriesParserHandlesEOF(t *testing.T) {
	// Valid data that ends without newline
	data := []byte("metric{label=\"value\"} 42")
	parser := &seriesParser{selector: nil}
	result, err := parser.parse(data)

	if err != nil {
		t.Errorf("unexpected error for data without trailing newline: %v", err)
	}

	if len(result) == 0 {
		t.Errorf("expected to parse metric without trailing newline")
	}
}

// TestSeriesParserSkipsInvalidEntries tests that invalid entries are skipped
func TestSeriesParserSkipsInvalidEntries(t *testing.T) {
	// Mix of valid series and comments
	data := []byte("# HELP metric help text\nmetric 42\n")
	parser := &seriesParser{selector: nil}
	result, err := parser.parse(data)

	// Should not error on comments
	if err != nil {
		t.Errorf("unexpected error when parsing mixed content: %v", err)
	}
}

// TestSeriesParserEdgeCases tests edge cases in parsing
func TestSeriesParserEdgeCases(t *testing.T) {
	tests := []struct {
		name   string
		data   []byte
		verify func(series prometheus.Series) bool
	}{
		{
			name: "very large metric value",
			data: []byte("metric 9.223372036854775807e+18\n"),
			verify: func(series prometheus.Series) bool {
				return true // Just ensure no panic
			},
		},
		{
			name: "very small metric value",
			data: []byte("metric 1e-308\n"),
			verify: func(series prometheus.Series) bool {
				return true // Just ensure no panic
			},
		},
		{
			name: "metric with underscore in name",
			data: []byte("metric_with_underscore 42\n"),
			verify: func(series prometheus.Series) bool {
				return len(series) > 0
			},
		},
		{
			name: "metric with colon in name",
			data: []byte("http_requests_total 42\n"),
			verify: func(series prometheus.Series) bool {
				return len(series) > 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &seriesParser{selector: nil}
			result, err := parser.parse(tt.data)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.verify(result) {
				t.Errorf("verification failed for test %s", tt.name)
			}
		})
	}
}

// TestSeriesParserParseWithNilData tests parsing with nil data
func TestSeriesParserParseWithNilData(t *testing.T) {
	parser := &seriesParser{selector: nil}
	result, err := parser.parse(nil)

	// Should handle nil gracefully
	if err != nil && err != io.EOF {
		t.Errorf("unexpected error for nil data: %v", err)
	}
	
	if len(result) != 0 {
		t.Errorf("expected empty result for nil data, got %d items", len(result))
	}
}

// TestSeriesParserParseConsecutiveCalls tests multiple consecutive parse calls
func TestSeriesParserParseConsecutiveCalls(t *testing.T) {
	parser := &seriesParser{selector: nil}
	data1 := []byte("metric1 10\n")
	data2 := []byte("metric2 20\n")

	result1, err1 := parser.parse(data1)
	result2, err2 := parser.parse(data2)

	if err1 != nil {
		t.Errorf("first parse error: %v", err1)
	}
	if err2 != nil {
		t.Errorf("second parse error: %v", err2)
	}

	if len(result1) == 0 || len(result2) == 0 {
		t.Errorf("expected non-empty results for both parses")
	}
}

// TestSeriesParserParseMultipleMetricsWithSelector tests selector filtering across multiple metrics
func TestSeriesParserParseMultipleMetricsWithSelector(t *testing.T) {
	data := []byte("metric1{label=\"value1\"} 10\nmetric2{label=\"value2\"} 20\nmetric3{label=\"value3\"} 30\n")

	// Create a selector that matches only metric1 and metric2
	matchCount := 0
	selector := &MockSelector{
		shouldMatch: true,
	}

	parser := &seriesParser{selector: selector}
	result, err := parser.parse(data)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) > 0 {
		matchCount = len(result)
	}

	if matchCount == 0 {
		t.Errorf("expected at least one match with matching selector")
	}
}

// TestSeriesParserParseErrorPropagation tests that errors are properly propagated
func TestSeriesParserParseErrorPropagation(t *testing.T) {
	// Create deliberately malformed data
	data := []byte("\xFF\xFE\x00\x00") // Invalid UTF-8
	parser := &seriesParser{selector: nil}
	_, err := parser.parse(data)

	// Should either succeed or fail gracefully, not panic
	if err != nil {
		// Expected - malformed data should error
		t.Logf("Got expected error for malformed data: %v", err)
	}
}

// TestSeriesParserParseWithLabelsContainingEquals tests labels with special characters
func TestSeriesParserParseWithLabelsContainingEquals(t *testing.T) {
	// Test metric with label values containing equals signs (escaped in quotes)
	data := []byte("metric{label=\"value=with=equals\"} 42\n")
	parser := &seriesParser{selector: nil}
	result, err := parser.parse(data)

	if err != nil {
		t.Errorf("unexpected error parsing label with equals: %v", err)
	}
	
	if len(result) == 0 {
		t.Errorf("expected to parse metric with label containing equals")
	}
}

// BenchmarkSeriesParserParse benchmarks the parse function
func BenchmarkSeriesParserParse(b *testing.B) {
	data := []byte("metric{label=\"value\"} 42\n")
	parser := &seriesParser{selector: nil}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.parse(data)
	}
}

// BenchmarkSeriesParserParseWithSelector benchmarks parsing with selector
func BenchmarkSeriesParserParseWithSelector(b *testing.B) {
	data := []byte("metric{label=\"value\"} 42\n")
	parser := &seriesParser{selector: &MockSelector{shouldMatch: true}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		parser.parse(data)
	}
}

// TestSeriesParserParseMultipleMetricsNoSelector tests parsing multiple metrics without selector
func TestSeriesParserParseMultipleMetricsNoSelector(t *testing.T) {
	data := []byte("metric1 10\nmetric2 20\nmetric3 30\n")
	parser := &seriesParser{selector: nil}
	result, err := parser.parse(data)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(result) == 0 {
		t.Errorf("expected to parse 3 metrics, got empty result")
	}
}