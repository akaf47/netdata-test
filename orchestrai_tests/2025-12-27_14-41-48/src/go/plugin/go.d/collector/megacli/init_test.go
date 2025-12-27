// SPDX-License-Identifier: GPL-3.0-or-later

//go:build linux || freebsd || openbsd || netbsd || dragonfly

package megacli

import (
	"testing"
	"time"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/logs"
)

// TestInitMegaCliExecSuccess tests that initMegaCliExec successfully creates a megaCli executor
func TestInitMegaCliExecSuccess(t *testing.T) {
	// Arrange
	logger := logs.New()
	collector := &Collector{
		Timeout: Duration{Duration: 30 * time.Second},
		Logger:  logger,
	}

	// Act
	exec, err := collector.initMegaCliExec()

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if exec == nil {
		t.Fatal("expected non-nil megaCli executor, got nil")
	}
}

// TestInitMegaCliExecWithDifferentTimeouts tests initMegaCliExec with various timeout values
func TestInitMegaCliExecWithDifferentTimeouts(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
	}{
		{
			name:    "short timeout",
			timeout: 1 * time.Second,
		},
		{
			name:    "normal timeout",
			timeout: 30 * time.Second,
		},
		{
			name:    "long timeout",
			timeout: 5 * time.Minute,
		},
		{
			name:    "zero timeout",
			timeout: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			logger := logs.New()
			collector := &Collector{
				Timeout: Duration{Duration: tt.timeout},
				Logger:  logger,
			}

			// Act
			exec, err := collector.initMegaCliExec()

			// Assert
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if exec == nil {
				t.Fatal("expected non-nil executor")
			}
		})
	}
}

// TestInitMegaCliExecReturnsNoError verifies the return error is always nil
func TestInitMegaCliExecReturnsNoError(t *testing.T) {
	// Arrange
	logger := logs.New()
	collector := &Collector{
		Timeout: Duration{Duration: 10 * time.Second},
		Logger:  logger,
	}

	// Act
	_, err := collector.initMegaCliExec()

	// Assert - verify nil error is returned
	if err != nil {
		t.Fatalf("expected error to be nil, got %v", err)
	}
}

// TestInitMegaCliExecWithNilLogger tests initMegaCliExec when logger might be nil
func TestInitMegaCliExecWithNilLogger(t *testing.T) {
	// Arrange
	collector := &Collector{
		Timeout: Duration{Duration: 30 * time.Second},
		Logger:  nil,
	}

	// Act
	exec, err := collector.initMegaCliExec()

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if exec == nil {
		t.Fatal("expected non-nil executor even with nil logger")
	}
}

// TestInitMegaCliExecMultipleCalls tests that multiple calls produce independent executors
func TestInitMegaCliExecMultipleCalls(t *testing.T) {
	// Arrange
	logger := logs.New()
	collector := &Collector{
		Timeout: Duration{Duration: 30 * time.Second},
		Logger:  logger,
	}

	// Act
	exec1, err1 := collector.initMegaCliExec()
	exec2, err2 := collector.initMegaCliExec()

	// Assert
	if err1 != nil || err2 != nil {
		t.Fatalf("expected no errors, got %v and %v", err1, err2)
	}
	if exec1 == nil || exec2 == nil {
		t.Fatal("expected non-nil executors")
	}
	// Verify they are different instances
	if exec1 == exec2 {
		t.Fatal("expected different executor instances")
	}
}

// TestInitMegaCliExecCreatesCorrectType verifies the returned executor is of megaCli type
func TestInitMegaCliExecCreatesCorrectType(t *testing.T) {
	// Arrange
	logger := logs.New()
	collector := &Collector{
		Timeout: Duration{Duration: 30 * time.Second},
		Logger:  logger,
	}

	// Act
	exec, _ := collector.initMegaCliExec()

	// Assert - exec should be of megaCli interface type
	var _ megaCli = exec
}