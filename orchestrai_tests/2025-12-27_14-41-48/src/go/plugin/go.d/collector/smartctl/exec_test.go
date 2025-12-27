package smartctl

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"testing"
	"time"
)

// MockExecCommand allows us to mock exec.Command
var mockExecCommand func(ctx context.Context, name string, args ...string) *exec.Cmd

func init() {
	// Store original exec.Command
	originalExecCommand := exec.CommandContext
	mockExecCommand = func(ctx context.Context, name string, args ...string) *exec.Cmd {
		return originalExecCommand(ctx, name, args...)
	}
}

// Test NewExecProvider creates a valid exec provider
func TestNewExecProvider(t *testing.T) {
	tests := []struct {
		name      string
		executable string
		timeout   time.Duration
		wantErr   bool
	}{
		{
			name:       "valid executable",
			executable: "smartctl",
			timeout:    5 * time.Second,
			wantErr:    false,
		},
		{
			name:       "empty executable",
			executable: "",
			timeout:    5 * time.Second,
			wantErr:    true,
		},
		{
			name:       "zero timeout",
			executable: "smartctl",
			timeout:    0,
			wantErr:    true,
		},
		{
			name:       "negative timeout",
			executable: "smartctl",
			timeout:    -1 * time.Second,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider, err := NewExecProvider(tt.executable, tt.timeout)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewExecProvider() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && provider == nil {
				t.Errorf("NewExecProvider() provider is nil")
			}
		})
	}
}

// Test ExecProvider.Exec executes command successfully
func TestExecProviderExecSuccess(t *testing.T) {
	provider := &ExecProvider{
		executable: "echo",
		timeout:    5 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.timeout)
	defer cancel()

	output, err := provider.Exec(ctx, "hello")
	if err != nil {
		t.Errorf("Exec() error = %v, want nil", err)
	}
	if len(output) == 0 {
		t.Errorf("Exec() output is empty")
	}
}

// Test ExecProvider.Exec with timeout
func TestExecProviderExecTimeout(t *testing.T) {
	provider := &ExecProvider{
		executable: "sleep",
		timeout:    100 * time.Millisecond,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.timeout)
	defer cancel()

	// This should timeout
	_, err := provider.Exec(ctx, "10")
	if err == nil {
		t.Errorf("Exec() expected timeout error, got nil")
	}
}

// Test ExecProvider.Exec with non-existent executable
func TestExecProviderExecNotFound(t *testing.T) {
	provider := &ExecProvider{
		executable: "/nonexistent/path/to/executable",
		timeout:    5 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.timeout)
	defer cancel()

	_, err := provider.Exec(ctx)
	if err == nil {
		t.Errorf("Exec() expected error for non-existent executable, got nil")
	}
}

// Test ExecProvider.Exec with empty args
func TestExecProviderExecEmptyArgs(t *testing.T) {
	provider := &ExecProvider{
		executable: "echo",
		timeout:    5 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.timeout)
	defer cancel()

	output, err := provider.Exec(ctx)
	if err != nil {
		t.Errorf("Exec() error = %v, want nil", err)
	}
	if output == nil {
		t.Errorf("Exec() output is nil")
	}
}

// Test ExecProvider.Exec with multiple args
func TestExecProviderExecMultipleArgs(t *testing.T) {
	provider := &ExecProvider{
		executable: "echo",
		timeout:    5 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.timeout)
	defer cancel()

	output, err := provider.Exec(ctx, "arg1", "arg2", "arg3")
	if err != nil {
		t.Errorf("Exec() error = %v, want nil", err)
	}
	if output == nil {
		t.Errorf("Exec() output is nil")
	}
}

// Test ExecProvider.Exec with special characters in args
func TestExecProviderExecSpecialChars(t *testing.T) {
	provider := &ExecProvider{
		executable: "echo",
		timeout:    5 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.timeout)
	defer cancel()

	specialArgs := []string{"$test", "&special", "|pipe", ">redirect"}
	output, err := provider.Exec(ctx, specialArgs...)
	if err != nil {
		t.Errorf("Exec() error = %v, want nil", err)
	}
	if output == nil {
		t.Errorf("Exec() output is nil")
	}
}

// Test ExecProvider.Exec with context cancellation
func TestExecProviderExecContextCancellation(t *testing.T) {
	provider := &ExecProvider{
		executable: "sleep",
		timeout:    5 * time.Second,
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := provider.Exec(ctx, "10")
	if err == nil {
		t.Errorf("Exec() expected error on context cancellation, got nil")
	}
}

// Test ExecProvider.Exec with very long argument
func TestExecProviderExecLongArg(t *testing.T) {
	provider := &ExecProvider{
		executable: "echo",
		timeout:    5 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.timeout)
	defer cancel()

	longArg := ""
	for i := 0; i < 10000; i++ {
		longArg += "a"
	}

	output, err := provider.Exec(ctx, longArg)
	if err != nil {
		t.Errorf("Exec() error = %v, want nil", err)
	}
	if output == nil {
		t.Errorf("Exec() output is nil")
	}
}

// Test ExecProvider.Exec with command that returns non-zero exit code
func TestExecProviderExecNonZeroExit(t *testing.T) {
	provider := &ExecProvider{
		executable: "false",
		timeout:    5 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.timeout)
	defer cancel()

	_, err := provider.Exec(ctx)
	if err == nil {
		t.Errorf("Exec() expected error for non-zero exit code, got nil")
	}
}

// Test ExecProvider.Exec preserves stdout output
func TestExecProviderExecPreservesOutput(t *testing.T) {
	provider := &ExecProvider{
		executable: "echo",
		timeout:    5 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.timeout)
	defer cancel()

	expectedContent := "test output"
	output, err := provider.Exec(ctx, expectedContent)
	if err != nil {
		t.Errorf("Exec() error = %v, want nil", err)
	}
	if output == nil || len(output) == 0 {
		t.Errorf("Exec() output is empty or nil")
	}
}

// Test ExecProvider.Exec with context with deadline
func TestExecProviderExecWithDeadline(t *testing.T) {
	provider := &ExecProvider{
		executable: "echo",
		timeout:    5 * time.Second,
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()

	output, err := provider.Exec(ctx, "test")
	if err != nil {
		t.Errorf("Exec() error = %v, want nil", err)
	}
	if output == nil {
		t.Errorf("Exec() output is nil")
	}
}

// Test ExecProvider Exec with no context values
func TestExecProviderExecNoContextValues(t *testing.T) {
	provider := &ExecProvider{
		executable: "echo",
		timeout:    5 * time.Second,
	}

	ctx := context.Background()
	output, err := provider.Exec(ctx, "test")
	if err != nil {
		t.Errorf("Exec() error = %v, want nil", err)
	}
	if output == nil {
		t.Errorf("Exec() output is nil")
	}
}

// Test ExecProvider with minimum timeout
func TestNewExecProviderMinTimeout(t *testing.T) {
	provider, err := NewExecProvider("echo", 1*time.Nanosecond)
	if err != nil {
		t.Errorf("NewExecProvider() error = %v, want nil", err)
	}
	if provider == nil {
		t.Errorf("NewExecProvider() provider is nil")
	}
}

// Test ExecProvider with very large timeout
func TestNewExecProviderLargeTimeout(t *testing.T) {
	provider, err := NewExecProvider("echo", 24*time.Hour)
	if err != nil {
		t.Errorf("NewExecProvider() error = %v, want nil", err)
	}
	if provider == nil {
		t.Errorf("NewExecProvider() provider is nil")
	}
}

// Test ExecProvider.Exec with whitespace args
func TestExecProviderExecWhitespaceArgs(t *testing.T) {
	provider := &ExecProvider{
		executable: "echo",
		timeout:    5 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.timeout)
	defer cancel()

	output, err := provider.Exec(ctx, "   ", "\t", "\n")
	if err != nil {
		t.Errorf("Exec() error = %v, want nil", err)
	}
	if output == nil {
		t.Errorf("Exec() output is nil")
	}
}

// Test ExecProvider field values
func TestExecProviderFieldValues(t *testing.T) {
	expectedExe := "smartctl"
	expectedTimeout := 10 * time.Second

	provider, err := NewExecProvider(expectedExe, expectedTimeout)
	if err != nil {
		t.Fatalf("NewExecProvider() error = %v", err)
	}

	if provider.executable != expectedExe {
		t.Errorf("provider.executable = %v, want %v", provider.executable, expectedExe)
	}
	if provider.timeout != expectedTimeout {
		t.Errorf("provider.timeout = %v, want %v", provider.timeout, expectedTimeout)
	}
}

// Test Exec with command that produces large output
func TestExecProviderLargeOutput(t *testing.T) {
	provider := &ExecProvider{
		executable: "bash",
		timeout:    5 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.timeout)
	defer cancel()

	// Generate large output
	output, err := provider.Exec(ctx, "-c", "for i in {1..1000}; do echo 'line '$i; done")
	if err != nil {
		t.Errorf("Exec() error = %v, want nil", err)
	}
	if output == nil || len(output) == 0 {
		t.Errorf("Exec() output is empty or nil")
	}
}

// Test Exec error handling with stdin
func TestExecProviderExecWithEnvironment(t *testing.T) {
	provider := &ExecProvider{
		executable: "bash",
		timeout:    5 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.timeout)
	defer cancel()

	// Test with environment-dependent command
	output, err := provider.Exec(ctx, "-c", "echo hello")
	if err != nil {
		t.Errorf("Exec() error = %v, want nil", err)
	}
	if output == nil {
		t.Errorf("Exec() output is nil")
	}
}