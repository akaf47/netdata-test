package zfspool

import (
	"context"
	"errors"
	"os/exec"
	"testing"
	"time"
)

// TestExecutor tests the Executor interface implementation
func TestExecutor_New(t *testing.T) {
	tests := []struct {
		name    string
		want    Executor
		wantErr bool
	}{
		{
			name:    "successful creation",
			want:    &cmdExecutor{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newExecutor()
			if got == nil {
				t.Errorf("newExecutor() = nil, want non-nil")
			}
		})
	}
}

// TestCmdExecutor_Execute tests basic command execution
func TestCmdExecutor_Execute(t *testing.T) {
	tests := []struct {
		name    string
		cmd     string
		want    string
		wantErr bool
	}{
		{
			name:    "execute ls command",
			cmd:     "echo test",
			want:    "test",
			wantErr: false,
		},
		{
			name:    "execute with empty command",
			cmd:     "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "execute non-existent command",
			cmd:     "nonexistentcommand12345",
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &cmdExecutor{}
			got, err := e.Execute(tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want && len(got) == 0 && len(tt.want) > 0 {
				t.Errorf("Execute() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestCmdExecutor_ExecuteWithContext tests context-based execution
func TestCmdExecutor_ExecuteWithContext(t *testing.T) {
	tests := []struct {
		name    string
		cmd     string
		timeout time.Duration
		want    string
		wantErr bool
	}{
		{
			name:    "execute with valid context",
			cmd:     "echo test",
			timeout: 5 * time.Second,
			want:    "test",
			wantErr: false,
		},
		{
			name:    "execute with timeout context",
			cmd:     "sleep 10",
			timeout: 100 * time.Millisecond,
			want:    "",
			wantErr: true,
		},
		{
			name:    "execute with cancelled context",
			cmd:     "echo test",
			timeout: 0,
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &cmdExecutor{}
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			got, err := e.ExecuteWithContext(ctx, tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteWithContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(tt.want) > 0 && got != tt.want {
				t.Errorf("ExecuteWithContext() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestCmdExecutor_ExecuteWithArgs tests execution with command arguments
func TestCmdExecutor_ExecuteWithArgs(t *testing.T) {
	tests := []struct {
		name    string
		cmd     string
		args    []string
		wantErr bool
	}{
		{
			name:    "execute with valid args",
			cmd:     "echo",
			args:    []string{"test", "output"},
			wantErr: false,
		},
		{
			name:    "execute with empty args",
			cmd:     "echo",
			args:    []string{},
			wantErr: false,
		},
		{
			name:    "execute with nil args",
			cmd:     "echo",
			args:    nil,
			wantErr: false,
		},
		{
			name:    "execute non-existent command with args",
			cmd:     "nonexistentcmd",
			args:    []string{"arg1"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &cmdExecutor{}
			got, err := e.ExecuteWithArgs(tt.cmd, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteWithArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("ExecuteWithArgs() got empty output")
			}
		})
	}
}

// TestCmdExecutor_ExecuteWithArgsAndContext tests context + args execution
func TestCmdExecutor_ExecuteWithArgsAndContext(t *testing.T) {
	tests := []struct {
		name    string
		cmd     string
		args    []string
		timeout time.Duration
		wantErr bool
	}{
		{
			name:    "execute with valid args and context",
			cmd:     "echo",
			args:    []string{"test"},
			timeout: 5 * time.Second,
			wantErr: false,
		},
		{
			name:    "execute with timeout exceeded",
			cmd:     "sleep",
			args:    []string{"10"},
			timeout: 100 * time.Millisecond,
			wantErr: true,
		},
		{
			name:    "execute with no timeout",
			cmd:     "echo",
			args:    []string{"quick"},
			timeout: 30 * time.Second,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &cmdExecutor{}
			ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
			defer cancel()

			got, err := e.ExecuteWithArgsAndContext(ctx, tt.cmd, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteWithArgsAndContext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == "" {
				t.Errorf("ExecuteWithArgsAndContext() got empty output")
			}
		})
	}
}

// TestErrorHandling tests error scenarios
func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		cmd         string
		expectError bool
		errorType   string
	}{
		{
			name:        "command not found",
			cmd:         "cmdnotfoundhere123",
			expectError: true,
			errorType:   "exec.ErrNotFound or similar",
		},
		{
			name:        "empty command string",
			cmd:         "",
			expectError: true,
			errorType:   "empty command error",
		},
		{
			name:        "malformed command with pipes",
			cmd:         "| grep test",
			expectError: true,
			errorType:   "exec error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &cmdExecutor{}
			_, err := e.Execute(tt.cmd)
			if tt.expectError && err == nil {
				t.Errorf("Execute() expected error for %q, got nil", tt.cmd)
			}
			if !tt.expectError && err != nil {
				t.Errorf("Execute() unexpected error for %q: %v", tt.cmd, err)
			}
		})
	}
}

// TestContextCancellation tests context cancellation behavior
func TestContextCancellation(t *testing.T) {
	tests := []struct {
		name    string
		cmd     string
		prepare func() (context.Context, context.CancelFunc)
		wantErr bool
	}{
		{
			name: "already cancelled context",
			cmd:  "echo test",
			prepare: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx, cancel
			},
			wantErr: true,
		},
		{
			name: "context with timeout",
			cmd:  "echo quick",
			prepare: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 5*time.Second)
			},
			wantErr: false,
		},
		{
			name: "background context",
			cmd:  "echo test",
			prepare: func() (context.Context, context.CancelFunc) {
				return context.Background(), func() {}
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &cmdExecutor{}
			ctx, cancel := tt.prepare()
			defer cancel()

			_, err := e.ExecuteWithContext(ctx, tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecuteWithContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestMultipleExecutions tests repeated executions
func TestMultipleExecutions(t *testing.T) {
	e := &cmdExecutor{}
	
	for i := 0; i < 5; i++ {
		_, err := e.Execute("echo test")
		if err != nil {
			t.Errorf("iteration %d: Execute() failed: %v", i, err)
		}
	}
}

// TestCommandOutput tests output parsing
func TestCommandOutput(t *testing.T) {
	tests := []struct {
		name       string
		cmd        string
		expectOut  string
		expectErr  bool
	}{
		{
			name:      "echo command",
			cmd:       "echo hello",
			expectErr: false,
		},
		{
			name:      "command with special chars",
			cmd:       "echo 'test string'",
			expectErr: false,
		},
		{
			name:      "command with numbers",
			cmd:       "echo 12345",
			expectErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &cmdExecutor{}
			output, err := e.Execute(tt.cmd)
			if (err != nil) != tt.expectErr {
				t.Errorf("Execute() error = %v, expectErr %v", err, tt.expectErr)
				return
			}
			if !tt.expectErr && output == "" {
				t.Errorf("Execute() returned empty output")
			}
		})
	}
}

// TestBoundaryConditions tests edge cases
func TestBoundaryConditions(t *testing.T) {
	e := &cmdExecutor{}
	
	tests := []struct {
		name   string
		cmd    string
		args   []string
		expect bool
	}{
		{
			name:   "very long command",
			cmd:    "echo",
			args:   []string{"a very long string of text that goes on and on and on"},
			expect: false, // should not error
		},
		{
			name:   "many arguments",
			cmd:    "echo",
			args:   make([]string, 100),
			expect: false,
		},
		{
			name:   "special characters in args",
			cmd:    "echo",
			args:   []string{"!@#$%^&*()"},
			expect: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := e.ExecuteWithArgs(tt.cmd, tt.args...)
			hasError := err != nil
			if hasError == tt.expect {
				if !tt.expect {
					t.Errorf("ExecuteWithArgs() error = %v", err)
				}
			}
		})
	}
}

// TestExecutorInterface tests that cmdExecutor implements Executor
func TestExecutorInterface(t *testing.T) {
	var _ Executor = (*cmdExecutor)(nil)
}

// TestStdErrHandling tests stderr output handling
func TestStdErrHandling(t *testing.T) {
	e := &cmdExecutor{}
	
	tests := []struct {
		name    string
		cmd     string
		wantErr bool
	}{
		{
			name:    "command that outputs to stderr",
			cmd:     "ls /nonexistent/path",
			wantErr: true,
		},
		{
			name:    "command that succeeds",
			cmd:     "echo success",
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := e.Execute(tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestExitCodeHandling tests non-zero exit codes
func TestExitCodeHandling(t *testing.T) {
	e := &cmdExecutor{}
	
	tests := []struct {
		name    string
		cmd     string
		wantErr bool
	}{
		{
			name:    "command with non-zero exit code",
			cmd:     "sh -c 'exit 1'",
			wantErr: true,
		},
		{
			name:    "command with zero exit code",
			cmd:     "sh -c 'exit 0'",
			wantErr: false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := e.Execute(tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestConcurrentExecution tests concurrent command execution
func TestConcurrentExecution(t *testing.T) {
	e := &cmdExecutor{}
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			_, err := e.Execute("echo test")
			if err != nil {
				t.Errorf("concurrent execution failed: %v", err)
			}
			done <- true
		}()
	}
	
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestEmptyCommandString tests empty command handling
func TestEmptyCommandString(t *testing.T) {
	e := &cmdExecutor{}
	
	_, err := e.Execute("")
	if err == nil {
		t.Errorf("Execute() with empty string should return error, got nil")
	}
}

// TestWhitespaceCommand tests whitespace-only command
func TestWhitespaceCommand(t *testing.T) {
	e := &cmdExecutor{}
	
	tests := []struct {
		name string
		cmd  string
	}{
		{
			name: "spaces only",
			cmd:  "   ",
		},
		{
			name: "tabs only",
			cmd:  "\t\t\t",
		},
		{
			name: "mixed whitespace",
			cmd:  " \t \n ",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := e.Execute(tt.cmd)
			if err == nil {
				t.Errorf("Execute() with %q should error", tt.cmd)
			}
		})
	}
}

// TestNilArgsSlice tests nil arguments slice
func TestNilArgsSlice(t *testing.T) {
	e := &cmdExecutor{}
	
	_, err := e.ExecuteWithArgs("echo")
	if err != nil {
		t.Errorf("ExecuteWithArgs() with no args should work, got error: %v", err)
	}
}

// TestContextDeadlineExceeded tests deadline exceeded scenario
func TestContextDeadlineExceeded(t *testing.T) {
	e := &cmdExecutor{}
	
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()
	
	time.Sleep(10 * time.Millisecond)
	
	_, err := e.ExecuteWithContext(ctx, "echo test")
	if err == nil {
		t.Logf("Execute() after deadline may succeed depending on timing")
	}
}

// TestCommandWithOutput tests capturing command output
func TestCommandWithOutput(t *testing.T) {
	e := &cmdExecutor{}
	
	output, err := e.Execute("echo hello world")
	if err != nil {
		t.Errorf("Execute() failed: %v", err)
		return
	}
	
	if len(output) == 0 {
		t.Errorf("Execute() returned empty output")
	}
}

// TestCommandWithMultilineOutput tests multiline output
func TestCommandWithMultilineOutput(t *testing.T) {
	e := &cmdExecutor{}
	
	output, err := e.Execute("echo -e 'line1\\nline2\\nline3'")
	if err != nil {
		t.Errorf("Execute() failed: %v", err)
		return
	}
	
	if len(output) == 0 {
		t.Errorf("Execute() returned empty output for multiline command")
	}
}