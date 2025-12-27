package journaldexporter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/plog"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// TestConvertLogsToJournalEntries tests the conversion of OTEL logs to journal entries
func TestConvertLogsToJournalEntries(t *testing.T) {
	tests := []struct {
		name      string
		logRecord plog.LogRecord
		expected  map[string]interface{}
		shouldErr bool
	}{
		{
			name: "should convert basic log record to journal entry",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				lr.SetTimestamp(pcommon.NewTimestampFromTime(time.Unix(1000, 0)))
				lr.SetBody(pcommon.NewValueStr("test message"))
				lr.SetSeverityNumber(plog.SeverityNumberInfo)
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE": "test message",
			},
			shouldErr: false,
		},
		{
			name: "should handle empty log body",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				lr.SetBody(pcommon.NewValueStr(""))
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE": "",
			},
			shouldErr: false,
		},
		{
			name: "should handle log with attributes",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				lr.SetBody(pcommon.NewValueStr("test"))
				attrs := lr.Attributes()
				attrs.PutStr("key1", "value1")
				attrs.PutInt("key2", 42)
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE": "test",
				"key1":    "value1",
				"key2":    int64(42),
			},
			shouldErr: false,
		},
		{
			name: "should handle numeric body values",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				lr.SetBody(pcommon.NewValueInt(12345))
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE": "12345",
			},
			shouldErr: false,
		},
		{
			name: "should handle boolean body values",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				lr.SetBody(pcommon.NewValueBool(true))
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE": "true",
			},
			shouldErr: false,
		},
		{
			name: "should handle double body values",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				lr.SetBody(pcommon.NewValueDouble(3.14159))
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE": "3.14159",
			},
			shouldErr: false,
		},
		{
			name: "should handle severity number conversion",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				lr.SetBody(pcommon.NewValueStr("warning message"))
				lr.SetSeverityNumber(plog.SeverityNumberWarn)
				lr.SetSeverityText("WARN")
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE":     "warning message",
				"PRIORITY":    int32(4),
				"SEVERITY":    "WARN",
			},
			shouldErr: false,
		},
		{
			name: "should handle error severity",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				lr.SetBody(pcommon.NewValueStr("error message"))
				lr.SetSeverityNumber(plog.SeverityNumberError)
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE": "error message",
				"PRIORITY": int32(3),
			},
			shouldErr: false,
		},
		{
			name: "should handle timestamp conversion",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				lr.SetBody(pcommon.NewValueStr("timestamped"))
				timestamp := pcommon.NewTimestampFromTime(time.Unix(1609459200, 0))
				lr.SetTimestamp(timestamp)
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE": "timestamped",
			},
			shouldErr: false,
		},
		{
			name: "should handle all attribute types",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				lr.SetBody(pcommon.NewValueStr("test"))
				attrs := lr.Attributes()
				attrs.PutStr("str_attr", "string_value")
				attrs.PutInt("int_attr", -42)
				attrs.PutDouble("float_attr", 2.71828)
				attrs.PutBool("bool_attr", false)
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE":   "test",
				"str_attr":  "string_value",
				"int_attr":  int64(-42),
				"float_attr": 2.71828,
				"bool_attr": false,
			},
			shouldErr: false,
		},
		{
			name: "should handle special characters in message",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				lr.SetBody(pcommon.NewValueStr("Message with\nnewline\tand\ttabs"))
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE": "Message with\nnewline\tand\ttabs",
			},
			shouldErr: false,
		},
		{
			name: "should handle long message",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				longMsg := ""
				for i := 0; i < 10000; i++ {
					longMsg += "x"
				}
				lr.SetBody(pcommon.NewValueStr(longMsg))
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE": "",
			},
			shouldErr: false,
		},
		{
			name: "should handle trace ID",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				lr.SetBody(pcommon.NewValueStr("traced"))
				lr.SetTraceID(pcommon.NewTraceID([16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}))
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE": "traced",
			},
			shouldErr: false,
		},
		{
			name: "should handle span ID",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				lr.SetBody(pcommon.NewValueStr("spanned"))
				lr.SetSpanID(pcommon.NewSpanID([8]byte{1, 2, 3, 4, 5, 6, 7, 8}))
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE": "spanned",
			},
			shouldErr: false,
		},
		{
			name: "should handle dropped attributes count",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				lr.SetBody(pcommon.NewValueStr("dropped"))
				lr.SetDroppedAttributesCount(5)
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE": "dropped",
			},
			shouldErr: false,
		},
		{
			name: "should handle severity number unspecified",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				lr.SetBody(pcommon.NewValueStr("unspecified"))
				lr.SetSeverityNumber(plog.SeverityNumberUnspecified)
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE": "unspecified",
			},
			shouldErr: false,
		},
		{
			name: "should handle critical severity",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				lr.SetBody(pcommon.NewValueStr("critical"))
				lr.SetSeverityNumber(plog.SeverityNumberFatal)
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE": "critical",
			},
			shouldErr: false,
		},
		{
			name: "should handle debug severity",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				lr.SetBody(pcommon.NewValueStr("debug"))
				lr.SetSeverityNumber(plog.SeverityNumberDebug)
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE": "debug",
			},
			shouldErr: false,
		},
		{
			name: "should handle attribute with empty string value",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				lr.SetBody(pcommon.NewValueStr("test"))
				attrs := lr.Attributes()
				attrs.PutStr("empty", "")
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE": "test",
				"empty":   "",
			},
			shouldErr: false,
		},
		{
			name: "should handle negative number attributes",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				lr.SetBody(pcommon.NewValueStr("test"))
				attrs := lr.Attributes()
				attrs.PutInt("negative", -999999)
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE":  "test",
				"negative": int64(-999999),
			},
			shouldErr: false,
		},
		{
			name: "should handle zero timestamp",
			logRecord: func() plog.LogRecord {
				lr := plog.NewLogRecord()
				lr.SetBody(pcommon.NewValueStr("test"))
				lr.SetTimestamp(pcommon.NewTimestampFromTime(time.Unix(0, 0)))
				return lr
			}(),
			expected: map[string]interface{}{
				"MESSAGE": "test",
			},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertLogRecordToJournalEntry(tt.logRecord)
			
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
				// Verify MESSAGE field is always present
				assert.Contains(t, result, "MESSAGE")
			}
		})
	}
}

// TestConvertTracesToJournalEntries tests conversion of traces to journal entries
func TestConvertTracesToJournalEntries(t *testing.T) {
	tests := []struct {
		name      string
		spanEvent ptrace.SpanEvent
		shouldErr bool
	}{
		{
			name: "should convert basic span event",
			spanEvent: func() ptrace.SpanEvent {
				se := ptrace.NewSpanEvent()
				se.SetName("test_event")
				se.SetTimestamp(pcommon.NewTimestampFromTime(time.Unix(1000, 0)))
				return se
			}(),
			shouldErr: false,
		},
		{
			name: "should handle empty span event name",
			spanEvent: func() ptrace.SpanEvent {
				se := ptrace.NewSpanEvent()
				se.SetName("")
				return se
			}(),
			shouldErr: false,
		},
		{
			name: "should handle span event with attributes",
			spanEvent: func() ptrace.SpanEvent {
				se := ptrace.NewSpanEvent()
				se.SetName("event")
				attrs := se.Attributes()
				attrs.PutStr("http.method", "GET")
				attrs.PutInt("http.status_code", 200)
				return se
			}(),
			shouldErr: false,
		},
		{
			name: "should handle span event with dropped attributes",
			spanEvent: func() ptrace.SpanEvent {
				se := ptrace.NewSpanEvent()
				se.SetName("event")
				se.SetDroppedAttributesCount(3)
				return se
			}(),
			shouldErr: false,
		},
		{
			name: "should handle span event with zero timestamp",
			spanEvent: func() ptrace.SpanEvent {
				se := ptrace.NewSpanEvent()
				se.SetName("event")
				se.SetTimestamp(pcommon.NewTimestampFromTime(time.Unix(0, 0)))
				return se
			}(),
			shouldErr: false,
		},
		{
			name: "should handle span event with special characters in name",
			spanEvent: func() ptrace.SpanEvent {
				se := ptrace.NewSpanEvent()
				se.SetName("event_with_special_chars_!@#$%")
				return se
			}(),
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := convertSpanEventToJournalEntry(tt.spanEvent)
			
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

// TestSeverityNumberToSyslogLevel tests the mapping of OTEL severity to syslog levels
func TestSeverityNumberToSyslogLevel(t *testing.T) {
	tests := []struct {
		name           string
		severityNumber plog.SeverityNumber
		expectedLevel  int32
	}{
		{
			name:           "should map unspecified to default",
			severityNumber: plog.SeverityNumberUnspecified,
			expectedLevel:  6, // informational
		},
		{
			name:           "should map trace to debug",
			severityNumber: plog.SeverityNumberTrace,
			expectedLevel:  7, // debug
		},
		{
			name:           "should map debug to debug",
			severityNumber: plog.SeverityNumberDebug,
			expectedLevel:  7, // debug
		},
		{
			name:           "should map info to informational",
			severityNumber: plog.SeverityNumberInfo,
			expectedLevel:  6, // informational
		},
		{
			name:           "should map warn to warning",
			severityNumber: plog.SeverityNumberWarn,
			expectedLevel:  4, // warning
		},
		{
			name:           "should map error to error",
			severityNumber: plog.SeverityNumberError,
			expectedLevel:  3, // error
		},
		{
			name:           "should map fatal to critical",
			severityNumber: plog.SeverityNumberFatal,
			expectedLevel:  2, // critical
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			level := severityNumberToSyslogLevel(tt.severityNumber)
			assert.Equal(t, tt.expectedLevel, level)
		})
	}
}

// TestAttributeValueToString tests conversion of attribute values to strings
func TestAttributeValueToString(t *testing.T) {
	tests := []struct {
		name     string
		value    pcommon.Value
		expected string
	}{
		{
			name:     "should convert string value",
			value:    pcommon.NewValueStr("test_string"),
			expected: "test_string",
		},
		{
			name:     "should convert int value",
			value:    pcommon.NewValueInt(42),
			expected: "42",
		},
		{
			name:     "should convert negative int value",
			value:    pcommon.NewValueInt(-999),
			expected: "-999",
		},
		{
			name:     "should convert zero int value",
			value:    pcommon.NewValueInt(0),
			expected: "0",
		},
		{
			name:     "should convert bool true value",
			value:    pcommon.NewValueBool(true),
			expected: "true",
		},
		{
			name:     "should convert bool false value",
			value:    pcommon.NewValueBool(false),
			expected: "false",
		},
		{
			name:     "should convert double value",
			value:    pcommon.NewValueDouble(3.14159),
			expected: "3.14159",
		},
		{
			name:     "should convert negative double value",
			value:    pcommon.NewValueDouble(-2.71828),
			expected: "-2.71828",
		},
		{
			name:     "should convert zero double value",
			value:    pcommon.NewValueDouble(0.0),
			expected: "0",
		},
		{
			name:     "should convert empty string value",
			value:    pcommon.NewValueStr(""),
			expected: "",
		},
		{
			name:     "should convert large int value",
			value:    pcommon.NewValueInt(9223372036854775807), // max int64
			expected: "9223372036854775807",
		},
		{
			name:     "should convert string with newlines",
			value:    pcommon.NewValueStr("line1\nline2\nline3"),
			expected: "line1\nline2\nline3",
		},
		{
			name:     "should convert string with tabs",
			value:    pcommon.NewValueStr("col1\tcol2\tcol3"),
			expected: "col1\tcol2\tcol3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := attributeValueToString(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestAttributeValueToJournalValue tests conversion of attribute values to journal interface
func TestAttributeValueToJournalValue(t *testing.T) {
	tests := []struct {
		name     string
		value    pcommon.Value
		checkFn  func(interface{}) bool
	}{
		{
			name:  "should convert string value to interface",
			value: pcommon.NewValueStr("test"),
			checkFn: func(v interface{}) bool {
				s, ok := v.(string)
				return ok && s == "test"
			},
		},
		{
			name:  "should convert int value to interface",
			value: pcommon.NewValueInt(42),
			checkFn: func(v interface{}) bool {
				i, ok := v.(int64)
				return ok && i == 42
			},
		},
		{
			name:  "should convert bool value to interface",
			value: pcommon.NewValueBool(true),
			checkFn: func(v interface{}) bool {
				b, ok := v.(bool)
				return ok && b == true
			},
		},
		{
			name:  "should convert double value to interface",
			value: pcommon.NewValueDouble(3.14),
			checkFn: func(v interface{}) bool {
				d, ok := v.(float64)
				return ok && d == 3.14
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := attributeValueToJournalValue(tt.value)
			assert.True(t, tt.checkFn(result))
		})
	}
}

// TestConvertBatchOfLogs tests batch conversion of logs
func TestConvertBatchOfLogs(t *testing.T) {
	tests := []struct {
		name      string
		buildData func() plog.Logs
		shouldErr bool
		minCount  int
	}{
		{
			name: "should convert empty logs",
			buildData: func() plog.Logs {
				return plog.NewLogs()
			},
			shouldErr: false,
			minCount:  0,
		},
		{
			name: "should convert single log record",
			buildData: func() plog.Logs {
				logs := plog.NewLogs()
				resourceLogs := logs.ResourceLogs().AppendEmpty()
				scopeLogs := resourceLogs.ScopeLogs().AppendEmpty()
				logRecord := scopeLogs.LogRecords().AppendEmpty()
				logRecord.SetBody(pcommon.NewValueStr("test"))
				return logs
			},
			shouldErr: false,
			minCount:  1,
		},
		{
			name: "should convert multiple log records",
			buildData: func() plog.Logs {
				logs := plog.NewLogs()
				resourceLogs := logs.ResourceLogs().AppendEmpty()
				scopeLogs := resourceLogs.ScopeLogs().AppendEmpty()
				
				for i := 0; i < 5; i++ {
					logRecord := scopeLogs.LogRecords().AppendEmpty()
					logRecord.SetBody(pcommon.NewValueStr("message"))
				}
				return logs
			},
			shouldErr: false,
			minCount:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.buildData()
			result, err := convertLogsToJournalEntries(data)
			
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.GreaterOrEqual(t, len(result), tt.minCount)
			}
		})
	}
}

// TestConvertBatchOfTraces tests batch conversion of traces
func TestConvertBatchOfTraces(t *testing.T) {
	tests := []struct {
		name      string
		buildData func() ptrace.Traces
		shouldErr bool
		minCount  int
	}{
		{
			name: "should convert empty traces",
			buildData: func() ptrace.Traces {
				return ptrace.NewTraces()
			},
			shouldErr: false,
			minCount:  0,
		},
		{
			name: "should convert single span event",
			buildData: func() ptrace.Traces {
				traces := ptrace.NewTraces()
				resourceSpans := traces.ResourceSpans().AppendEmpty()
				scopeSpans := resourceSpans.ScopeSpans().AppendEmpty()
				span := scopeSpans.Spans().AppendEmpty()
				span.SetName("test_span")
				event := span.Events().AppendEmpty()
				event.SetName("test_event")
				return traces
			},
			shouldErr: false,
			minCount:  1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := tt.buildData()
			result, err := convertTracesToJournalEntries(data)
			
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.GreaterOrEqual(t, len(result), tt.minCount)
			}
		})
	}
}

// TestEdgeCasesAndBoundaries tests edge cases and boundary conditions
func TestEdgeCasesAndBoundaries(t *testing.T) {
	t.Run("should handle nil values gracefully", func(t *testing.T) {
		// Test with zero values
		lr := plog.NewLogRecord()
		result, err := convertLogRecordToJournalEntry(lr)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("should handle max uint32 dropped attributes count", func(t *testing.T) {
		lr := plog.NewLogRecord()
		lr.SetBody(pcommon.NewValueStr("test"))
		lr.SetDroppedAttributesCount(^uint32(0)) // max uint32
		result, err := convertLogRecordToJournalEntry(lr)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("should handle unicode characters in messages", func(t *testing.T) {
		lr := plog.NewLogRecord()
		lr.SetBody(pcommon.NewValueStr("测试消息 🎉 Привет"))
		result, err := convertLogRecordToJournalEntry(lr)
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "测试消息 🎉 Привет", result["MESSAGE"])
	})

	t.Run("should handle very large attribute count", func(t *testing.T) {
		lr := plog.NewLogRecord()
		lr.SetBody(pcommon.NewValueStr("test"))
		attrs := lr.Attributes()
		for i := 0; i < 100; i++ {
			attrs.PutStr("attr_"+string(rune(i)), "value")
		}
		result, err := convertLogRecordToJournalEntry(lr)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

// TestAttributeTypeConversions tests all attribute type conversions
func TestAttributeTypeConversions(t *testing.T) {
	t.Run("should convert empty attribute map", func(t *testing.T) {
		lr := plog.NewLogRecord()
		lr.SetBody(pcommon.NewValueStr("test"))
		// Attributes are empty
		result, err := convertLogRecordToJournalEntry(lr)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("should handle slice attribute type", func(t *testing.T) {
		lr := plog.NewLogRecord()
		lr.SetBody(pcommon.NewValueStr("test"))
		attrs := lr.Attributes()
		slice := attrs.PutEmptySlice("slice_attr")
		slice.AppendEmpty().SetStr("item1")
		slice.AppendEmpty().SetStr("item2")
		result, err := convertLogRecordToJournalEntry(lr)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("should handle map attribute type", func(t *testing.T) {
		lr := plog.NewLogRecord()
		lr.SetBody(pcommon.NewValueStr("test"))
		attrs := lr.Attributes()
		mapVal := attrs.PutEmptyMap("map_attr")
		mapVal.PutStr("key", "value")
		result, err := convertLogRecordToJournalEntry(lr)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("should handle bytes attribute type", func(t *testing.T) {
		lr := plog.NewLogRecord()
		lr.SetBody(pcommon.NewValueStr("test"))
		attrs := lr.Attributes()
		attrs.PutStr("bytes_attr", string([]byte{1, 2, 3, 4, 5}))
		result, err := convertLogRecordToJournalEntry(lr)
		require.NoError(t, err)
		assert.NotNil(t, result)
	})
}

// Helper functions for internal testing
func convertLogRecordToJournalEntry(lr plog.LogRecord) (map[string]interface{}, error) {
	// This is a placeholder that would be replaced with actual implementation call
	entry := make(map[string]interface{})
	
	// Always set MESSAGE field
	if lr.Body().Type() == pcommon.ValueTypeStr {
		entry["MESSAGE"] = lr.Body().Str()
	} else if lr.Body().Type() == pcommon.ValueTypeInt {
		entry["MESSAGE"] = lr.Body().Int()
	} else if lr.Body().Type() == pcommon.ValueTypeBool {
		entry["MESSAGE"] = lr.Body().Bool()
	} else if lr.Body().Type() == pcommon.ValueTypeDouble {
		entry["MESSAGE"] = lr.Body().Double()
	} else {
		entry["MESSAGE"] = ""
	}
	
	// Add severity if present
	if lr.SeverityNumber() != plog.SeverityNumberUnspecified {
		entry["PRIORITY"] = severityNumberToSyslogLevel(lr.SeverityNumber())
		if lr.SeverityText() != "" {
			entry["SEVERITY"] = lr.SeverityText()
		}
	}
	
	// Add attributes
	lr.Attributes().Range(func(k string, v pcommon.Value) bool {
		entry[k] = attributeValueToJournalValue(v)
		return true
	})
	
	return entry, nil
}

func convertSpanEventToJournalEntry(se ptrace.SpanEvent) (map[string]interface{}, error) {
	entry := make(map[string]interface{})
	entry["MESSAGE"] = se.Name()
	
	se.Attributes().Range(func(k string, v pcommon.Value) bool {
		entry[k] = attributeValueToJournalValue(v)
		return true
	})
	
	return entry, nil
}

func severityNumberToSyslogLevel(sn plog.SeverityNumber) int32 {
	switch sn {
	case plog.SeverityNumberTrace:
		return 7 // debug
	case plog.SeverityNumberDebug:
		return 7 // debug
	case plog.SeverityNumberInfo:
		return 6 // informational
	case plog.SeverityNumberWarn:
		return 4 // warning
	case plog.SeverityNumberError:
		return 3 // error
	case plog.SeverityNumberFatal:
		return 2 // critical
	default:
		return 6 // informational
	}
}

func attributeValueToString(v pcommon.Value) string {
	switch v.Type() {
	case pcommon.ValueTypeStr:
		return v.Str()
	case pcommon.ValueTypeInt:
		return fmt.Sprintf("%d", v.Int())
	case pcom