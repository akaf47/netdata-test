package logger

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"testing"
)

// TestNewLogger tests logger initialization
func TestNewLogger(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "successful logger creation",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger()
			if logger == nil {
				t.Fatal("NewLogger returned nil")
			}
		})
	}
}

// TestLogInfo tests the Info logging function
func TestLogInfo(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
		wantErr  bool
	}{
		{
			name:     "log basic info message",
			message:  "test info message",
			expected: "test info message",
			wantErr:  false,
		},
		{
			name:     "log empty message",
			message:  "",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "log long message",
			message:  strings.Repeat("a", 10000),
			expected: strings.Repeat("a", 10000),
			wantErr:  false,
		},
		{
			name:     "log message with special characters",
			message:  "test\nwith\ttabs\rand\x00null",
			expected: "test\nwith\ttabs\rand\x00null",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger()
			if logger == nil {
				t.Fatal("NewLogger returned nil")
			}

			// Capture output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			logger.Info(tt.message)

			w.Close()
			os.Stdout = oldStdout

			out, _ := io.ReadAll(r)
			output := string(out)

			if tt.expected != "" && !strings.Contains(output, tt.expected) {
				t.Errorf("Info() output %q does not contain expected %q", output, tt.expected)
			}
		})
	}
}

// TestLogWarn tests the Warn logging function
func TestLogWarn(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
		wantErr  bool
	}{
		{
			name:     "log basic warn message",
			message:  "test warn message",
			expected: "test warn message",
			wantErr:  false,
		},
		{
			name:     "log empty warning",
			message:  "",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "log long warning",
			message:  strings.Repeat("w", 10000),
			expected: strings.Repeat("w", 10000),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger()
			if logger == nil {
				t.Fatal("NewLogger returned nil")
			}

			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			logger.Warn(tt.message)

			w.Close()
			os.Stderr = oldStderr

			out, _ := io.ReadAll(r)
			output := string(out)

			if tt.expected != "" && !strings.Contains(output, tt.expected) {
				t.Errorf("Warn() output %q does not contain expected %q", output, tt.expected)
			}
		})
	}
}

// TestLogError tests the Error logging function
func TestLogError(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		err      error
		expected string
		wantErr  bool
	}{
		{
			name:     "log error with message and error",
			message:  "test error",
			err:      errors.New("test error object"),
			expected: "test error",
			wantErr:  false,
		},
		{
			name:     "log error with nil error",
			message:  "test error message",
			err:      nil,
			expected: "test error message",
			wantErr:  false,
		},
		{
			name:     "log error with empty message",
			message:  "",
			err:      errors.New("actual error"),
			expected: "actual error",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger()
			if logger == nil {
				t.Fatal("NewLogger returned nil")
			}

			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			logger.Error(tt.message, tt.err)

			w.Close()
			os.Stderr = oldStderr

			out, _ := io.ReadAll(r)
			output := string(out)

			if tt.expected != "" && !strings.Contains(output, tt.expected) {
				t.Errorf("Error() output %q does not contain expected %q", output, tt.expected)
			}
		})
	}
}

// TestLogDebug tests the Debug logging function
func TestLogDebug(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
		wantErr  bool
	}{
		{
			name:     "log debug message",
			message:  "test debug",
			expected: "test debug",
			wantErr:  false,
		},
		{
			name:     "log empty debug",
			message:  "",
			expected: "",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger()
			if logger == nil {
				t.Fatal("NewLogger returned nil")
			}

			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			logger.Debug(tt.message)

			w.Close()
			os.Stdout = oldStdout

			out, _ := io.ReadAll(r)
			_ = string(out)
		})
	}
}

// TestLogPanic tests the Panic logging function
func TestLogPanic(t *testing.T) {
	tests := []struct {
		name    string
		message string
		wantErr bool
	}{
		{
			name:    "log panic message",
			message: "test panic",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger()
			if logger == nil {
				t.Fatal("NewLogger returned nil")
			}

			defer func() {
				if r := recover(); r == nil {
					t.Errorf("Panic() did not panic")
				}
			}()

			logger.Panic(tt.message)
		})
	}
}

// TestLogFatal tests the Fatal logging function
func TestLogFatal(t *testing.T) {
	tests := []struct {
		name    string
		message string
		wantErr bool
	}{
		{
			name:    "log fatal message",
			message: "test fatal",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger()
			if logger == nil {
				t.Fatal("NewLogger returned nil")
			}

			if os.Getenv("BE_FATAL") == "1" {
				logger.Fatal(tt.message)
				return
			}

			cmd := os.Getenv("TEST_FATAL")
			if cmd != "" {
				logger.Fatal(tt.message)
			}
		})
	}
}

// TestLogSetLevel tests setting log level
func TestLogSetLevel(t *testing.T) {
	tests := []struct {
		name  string
		level string
	}{
		{
			name:  "set debug level",
			level: "debug",
		},
		{
			name:  "set info level",
			level: "info",
		},
		{
			name:  "set warn level",
			level: "warn",
		},
		{
			name:  "set error level",
			level: "error",
		},
		{
			name:  "set invalid level",
			level: "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger()
			if logger == nil {
				t.Fatal("NewLogger returned nil")
			}

			logger.SetLevel(tt.level)
		})
	}
}

// TestLogWithFields tests logging with additional fields
func TestLogWithFields(t *testing.T) {
	tests := []struct {
		name   string
		fields map[string]interface{}
	}{
		{
			name:   "log with empty fields",
			fields: map[string]interface{}{},
		},
		{
			name: "log with single field",
			fields: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name: "log with multiple fields",
			fields: map[string]interface{}{
				"string": "value",
				"int":    42,
				"bool":   true,
				"nil":    nil,
			},
		},
		{
			name: "log with complex field values",
			fields: map[string]interface{}{
				"slice": []int{1, 2, 3},
				"map":   map[string]string{"a": "b"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger()
			if logger == nil {
				t.Fatal("NewLogger returned nil")
			}

			logger.WithFields(tt.fields)
		})
	}
}

// TestLogConcurrency tests concurrent logging
func TestLogConcurrency(t *testing.T) {
	logger := NewLogger()
	if logger == nil {
		t.Fatal("NewLogger returned nil")
	}

	done := make(chan bool)

	for i := 0; i < 100; i++ {
		go func(index int) {
			logger.Info("Concurrent log message")
			logger.Warn("Concurrent warn message")
			logger.Error("Concurrent error message", errors.New("test error"))
			done <- true
		}(i)
	}

	for i := 0; i < 100; i++ {
		<-done
	}
}

// TestLogMultipleInstances tests multiple logger instances
func TestLogMultipleInstances(t *testing.T) {
	logger1 := NewLogger()
	logger2 := NewLogger()

	if logger1 == nil || logger2 == nil {
		t.Fatal("NewLogger returned nil")
	}

	logger1.Info("Logger 1")
	logger2.Info("Logger 2")
}

// TestLogInfof tests the Infof logging function with formatting
func TestLogInfof(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "format with single argument",
			format:   "Value: %v",
			args:     []interface{}{42},
			expected: "42",
		},
		{
			name:     "format with multiple arguments",
			format:   "%s %d %v",
			args:     []interface{}{"test", 123, true},
			expected: "test 123 true",
		},
		{
			name:     "format with no arguments",
			format:   "no args",
			args:     []interface{}{},
			expected: "no args",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger()
			if logger == nil {
				t.Fatal("NewLogger returned nil")
			}

			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			logger.Infof(tt.format, tt.args...)

			w.Close()
			os.Stdout = oldStdout

			out, _ := io.ReadAll(r)
			output := string(out)

			if !strings.Contains(output, tt.expected) {
				t.Errorf("Infof() output %q does not contain expected %q", output, tt.expected)
			}
		})
	}
}

// TestLogWarnf tests the Warnf logging function with formatting
func TestLogWarnf(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "warn format with single argument",
			format:   "Warning: %v",
			args:     []interface{}{999},
			expected: "999",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger()
			if logger == nil {
				t.Fatal("NewLogger returned nil")
			}

			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			logger.Warnf(tt.format, tt.args...)

			w.Close()
			os.Stderr = oldStderr

			out, _ := io.ReadAll(r)
			output := string(out)

			if !strings.Contains(output, tt.expected) {
				t.Errorf("Warnf() output %q does not contain expected %q", output, tt.expected)
			}
		})
	}
}

// TestLogErrorf tests the Errorf logging function with formatting
func TestLogErrorf(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "error format with single argument",
			format:   "Error: %v",
			args:     []interface{}{"failed"},
			expected: "failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger()
			if logger == nil {
				t.Fatal("NewLogger returned nil")
			}

			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			logger.Errorf(tt.format, tt.args...)

			w.Close()
			os.Stderr = oldStderr

			out, _ := io.ReadAll(r)
			output := string(out)

			if !strings.Contains(output, tt.expected) {
				t.Errorf("Errorf() output %q does not contain expected %q", output, tt.expected)
			}
		})
	}
}

// TestLogDebugf tests the Debugf logging function with formatting
func TestLogDebugf(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		args     []interface{}
		expected string
	}{
		{
			name:     "debug format with argument",
			format:   "Debug: %v",
			args:     []interface{}{"info"},
			expected: "info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := NewLogger()
			if logger == nil {
				t.Fatal("NewLogger returned nil")
			}

			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			logger.Debugf(tt.format, tt.args...)

			w.Close()
			os.Stdout = oldStdout

			out, _ := io.ReadAll(r)
			_ = string(out)
		})
	}
}