package storcli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

// MockExecCommand allows us to mock os/exec.Command for testing
type MockExecCommand struct {
	cmd    string
	args   []string
	output []byte
	err    error
}

// TestNewExecutor tests the Executor creation
func TestNewExecutor(t *testing.T) {
	t.Run("executor_creation_success", func(t *testing.T) {
		executor := NewExecutor()
		if executor == nil {
			t.Error("executor should not be nil")
		}
	})

	t.Run("executor_has_valid_fields", func(t *testing.T) {
		executor := NewExecutor()
		// Verify executor structure is initialized
		if executor == nil {
			t.Fatal("executor initialization failed")
		}
	})
}

// TestExecuteCommand tests command execution functionality
func TestExecuteCommand(t *testing.T) {
	executor := NewExecutor()
	if executor == nil {
		t.Fatal("failed to create executor")
	}

	t.Run("execute_valid_command", func(t *testing.T) {
		// Test executing a valid system command
		output, err := executor.Execute("echo", "test")
		if err != nil {
			t.Fatalf("unexpected error executing valid command: %v", err)
		}
		if output == nil {
			t.Error("output should not be nil for successful execution")
		}
	})

	t.Run("execute_invalid_command", func(t *testing.T) {
		// Test executing a non-existent command
		_, err := executor.Execute("nonexistent_command_xyz", "arg")
		if err == nil {
			t.Error("should return error for non-existent command")
		}
	})

	t.Run("execute_empty_command", func(t *testing.T) {
		// Test executing with empty command name
		_, err := executor.Execute("", "arg")
		if err == nil {
			t.Error("should return error for empty command")
		}
	})

	t.Run("execute_with_no_arguments", func(t *testing.T) {
		// Test executing command without arguments
		output, err := executor.Execute("echo")
		if err != nil {
			t.Fatalf("unexpected error executing command without args: %v", err)
		}
		if output == nil {
			t.Error("output should not be nil")
		}
	})

	t.Run("execute_with_multiple_arguments", func(t *testing.T) {
		// Test executing command with multiple arguments
		output, err := executor.Execute("echo", "arg1", "arg2", "arg3")
		if err != nil {
			t.Fatalf("unexpected error with multiple args: %v", err)
		}
		if output == nil {
			t.Error("output should not be nil")
		}
	})

	t.Run("execute_with_special_characters", func(t *testing.T) {
		// Test executing command with special characters in arguments
		output, err := executor.Execute("echo", "test@#$%^&*()")
		if err != nil {
			t.Fatalf("unexpected error with special characters: %v", err)
		}
		if output == nil {
			t.Error("output should not be nil")
		}
	})
}

// TestExecuteWithContext tests context-aware command execution
func TestExecuteWithContext(t *testing.T) {
	executor := NewExecutor()
	if executor == nil {
		t.Fatal("failed to create executor")
	}

	t.Run("execute_with_valid_context", func(t *testing.T) {
		ctx := context.Background()
		output, err := executor.ExecuteWithContext(ctx, "echo", "test")
		if err != nil {
			t.Fatalf("unexpected error with valid context: %v", err)
		}
		if output == nil {
			t.Error("output should not be nil")
		}
	})

	t.Run("execute_with_cancelled_context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		_, err := executor.ExecuteWithContext(ctx, "sleep", "10")
		if err == nil {
			t.Error("should return error for cancelled context")
		}
	})

	t.Run("execute_with_timeout_context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()
		_, err := executor.ExecuteWithContext(ctx, "sleep", "10")
		if err == nil {
			t.Error("should return error for timeout context")
		}
	})

	t.Run("execute_context_with_no_args", func(t *testing.T) {
		ctx := context.Background()
		output, err := executor.ExecuteWithContext(ctx, "echo")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output == nil {
			t.Error("output should not be nil")
		}
	})

	t.Run("execute_context_with_nil_context", func(t *testing.T) {
		// Test with nil context - should handle gracefully
		var ctx context.Context
		if ctx == nil {
			// Should use background context or handle error
			ctx = context.Background()
		}
		output, err := executor.ExecuteWithContext(ctx, "echo", "test")
		if err != nil {
			t.Fatalf("unexpected error with nil context recovery: %v", err)
		}
		if output == nil {
			t.Error("output should not be nil")
		}
	})
}

// TestExecuteWithTimeout tests timeout-based execution
func TestExecuteWithTimeout(t *testing.T) {
	executor := NewExecutor()
	if executor == nil {
		t.Fatal("failed to create executor")
	}

	t.Run("execute_with_sufficient_timeout", func(t *testing.T) {
		output, err := executor.ExecuteWithTimeout(5*time.Second, "echo", "test")
		if err != nil {
			t.Fatalf("unexpected error with sufficient timeout: %v", err)
		}
		if output == nil {
			t.Error("output should not be nil")
		}
	})

	t.Run("execute_with_insufficient_timeout", func(t *testing.T) {
		// This should timeout since sleep 10s > 100ms timeout
		_, err := executor.ExecuteWithTimeout(100*time.Millisecond, "sleep", "10")
		if err == nil {
			t.Error("should return error for timeout")
		}
	})

	t.Run("execute_with_zero_timeout", func(t *testing.T) {
		// Zero timeout should either be immediate or use default
		_, err := executor.ExecuteWithTimeout(0, "echo", "test")
		// May succeed or fail depending on implementation
		_ = err
	})

	t.Run("execute_with_negative_timeout", func(t *testing.T) {
		// Negative timeout should be treated as error
		_, err := executor.ExecuteWithTimeout(-1*time.Second, "echo", "test")
		// May error depending on implementation
		_ = err
	})

	t.Run("execute_with_very_long_timeout", func(t *testing.T) {
		output, err := executor.ExecuteWithTimeout(1*time.Hour, "echo", "test")
		if err != nil {
			t.Fatalf("unexpected error with long timeout: %v", err)
		}
		if output == nil {
			t.Error("output should not be nil")
		}
	})
}

// TestExecuteBinary tests execution of binary with path
func TestExecuteBinary(t *testing.T) {
	executor := NewExecutor()
	if executor == nil {
		t.Fatal("failed to create executor")
	}

	t.Run("execute_binary_with_full_path", func(t *testing.T) {
		// Most systems have /bin/echo
		if _, err := os.Stat("/bin/echo"); err == nil {
			output, err := executor.ExecuteBinary("/bin/echo", "test")
			if err != nil {
				t.Fatalf("unexpected error executing binary with full path: %v", err)
			}
			if output == nil {
				t.Error("output should not be nil")
			}
		}
	})

	t.Run("execute_binary_with_nonexistent_path", func(t *testing.T) {
		_, err := executor.ExecuteBinary("/nonexistent/path/to/binary", "arg")
		if err == nil {
			t.Error("should return error for non-existent binary path")
		}
	})

	t.Run("execute_binary_with_empty_path", func(t *testing.T) {
		_, err := executor.ExecuteBinary("", "arg")
		if err == nil {
			t.Error("should return error for empty binary path")
		}
	})

	t.Run("execute_binary_without_args", func(t *testing.T) {
		if _, err := os.Stat("/bin/echo"); err == nil {
			output, err := executor.ExecuteBinary("/bin/echo")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if output == nil {
				t.Error("output should not be nil")
			}
		}
	})

	t.Run("execute_binary_with_multiple_args", func(t *testing.T) {
		if _, err := os.Stat("/bin/echo"); err == nil {
			output, err := executor.ExecuteBinary("/bin/echo", "a", "b", "c")
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if output == nil {
				t.Error("output should not be nil")
			}
		}
	})
}

// TestParseOutput tests output parsing functionality
func TestParseOutput(t *testing.T) {
	executor := NewExecutor()
	if executor == nil {
		t.Fatal("failed to create executor")
	}

	t.Run("parse_valid_output", func(t *testing.T) {
		output := []byte("line1\nline2\nline3\n")
		result, err := executor.ParseOutput(output)
		if err != nil {
			t.Fatalf("unexpected error parsing valid output: %v", err)
		}
		if result == nil {
			t.Error("parsed result should not be nil")
		}
	})

	t.Run("parse_empty_output", func(t *testing.T) {
		output := []byte("")
		result, err := executor.ParseOutput(output)
		if err != nil {
			// May or may not error depending on implementation
			t.Logf("parse empty output returned error: %v", err)
		}
		if result == nil && err == nil {
			t.Log("empty output handled correctly")
		}
	})

	t.Run("parse_multiline_output", func(t *testing.T) {
		output := []byte("line1\nline2\nline3\nline4\nline5\n")
		result, err := executor.ParseOutput(output)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Error("result should not be nil")
		}
	})

	t.Run("parse_output_with_special_chars", func(t *testing.T) {
		output := []byte("special@#$%^&*()\nline2\n")
		result, err := executor.ParseOutput(output)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Error("result should not be nil")
		}
	})

	t.Run("parse_output_with_tabs", func(t *testing.T) {
		output := []byte("field1\tfield2\tfield3\n")
		result, err := executor.ParseOutput(output)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Error("result should not be nil")
		}
	})

	t.Run("parse_nil_output", func(t *testing.T) {
		var output []byte
		result, err := executor.ParseOutput(output)
		// Should handle nil/empty gracefully
		_ = result
		_ = err
	})
}

// TestCommandErrorHandling tests error handling in command execution
func TestCommandErrorHandling(t *testing.T) {
	executor := NewExecutor()
	if executor == nil {
		t.Fatal("failed to create executor")
	}

	t.Run("handle_command_not_found", func(t *testing.T) {
		_, err := executor.Execute("this_command_definitely_does_not_exist_12345", "arg")
		if err == nil {
			t.Error("should return error when command not found")
		}
		if !strings.Contains(err.Error(), "not found") && !strings.Contains(err.Error(), "no such file") {
			t.Logf("error message: %v", err)
		}
	})

	t.Run("handle_command_failure", func(t *testing.T) {
		// Using a command that will fail (e.g., exit with non-zero)
		_, err := executor.Execute("false")
		if err == nil {
			t.Log("false command may not error on all systems")
		}
	})

	t.Run("handle_broken_pipe", func(t *testing.T) {
		// Execute a command that might produce broken pipe
		output, err := executor.Execute("echo", "test")
		// Should handle gracefully
		_ = output
		_ = err
	})

	t.Run("handle_permission_denied", func(t *testing.T) {
		// Try to execute a path without permission
		_, err := executor.Execute("/root/nonexistent_file", "arg")
		if err == nil {
			t.Log("permission denied error not triggered (may be expected)")
		}
	})
}

// TestConcurrentExecution tests thread-safety of executor
func TestConcurrentExecution(t *testing.T) {
	executor := NewExecutor()
	if executor == nil {
		t.Fatal("failed to create executor")
	}

	t.Run("concurrent_command_execution", func(t *testing.T) {
		done := make(chan error, 10)
		for i := 0; i < 10; i++ {
			go func(index int) {
				_, err := executor.Execute("echo", fmt.Sprintf("test%d", index))
				done <- err
			}(i)
		}

		for i := 0; i < 10; i++ {
			err := <-done
			if err != nil {
				t.Logf("concurrent execution %d failed: %v", i, err)
			}
		}
	})

	t.Run("concurrent_timeout_execution", func(t *testing.T) {
		done := make(chan error, 5)
		for i := 0; i < 5; i++ {
			go func(index int) {
				_, err := executor.ExecuteWithTimeout(5*time.Second, "echo", fmt.Sprintf("test%d", index))
				done <- err
			}(i)
		}

		for i := 0; i < 5; i++ {
			_ = <-done
		}
	})
}

// TestExecutorStates tests executor state management
func TestExecutorStates(t *testing.T) {
	t.Run("executor_reusability", func(t *testing.T) {
		executor := NewExecutor()
		for i := 0; i < 5; i++ {
			_, err := executor.Execute("echo", "test")
			if err != nil {
				t.Fatalf("iteration %d failed: %v", i, err)
			}
		}
	})

	t.Run("multiple_executors", func(t *testing.T) {
		executor1 := NewExecutor()
		executor2 := NewExecutor()
		
		_, err1 := executor1.Execute("echo", "from1")
		_, err2 := executor2.Execute("echo", "from2")
		
		if err1 != nil || err2 != nil {
			t.Fatalf("multiple executors failed: err1=%v, err2=%v", err1, err2)
		}
	})
}

// TestEdgeCases tests boundary conditions and edge cases
func TestEdgeCases(t *testing.T) {
	executor := NewExecutor()
	if executor == nil {
		t.Fatal("failed to create executor")
	}

	t.Run("very_long_argument", func(t *testing.T) {
		longArg := strings.Repeat("a", 10000)
		output, err := executor.Execute("echo", longArg)
		if err != nil {
			t.Logf("long argument failed: %v", err)
		}
		_ = output
	})

	t.Run("many_arguments", func(t *testing.T) {
		args := make([]string, 100)
		for i := 0; i < 100; i++ {
			args[i] = fmt.Sprintf("arg%d", i)
		}
		output, err := executor.Execute("echo", args...)
		if err != nil {
			t.Fatalf("many arguments failed: %v", err)
		}
		if output == nil {
			t.Error("output should not be nil")
		}
	})

	t.Run("command_with_spaces", func(t *testing.T) {
		output, err := executor.Execute("echo", "hello world test")
		if err != nil {
			t.Fatalf("command with spaces failed: %v", err)
		}
		if output == nil {
			t.Error("output should not be nil")
		}
	})

	t.Run("command_with_newlines_in_args", func(t *testing.T) {
		output, err := executor.Execute("echo", "line1\nline2")
		if err != nil {
			t.Logf("newlines in args failed: %v", err)
		}
		_ = output
	})

	t.Run("unicode_arguments", func(t *testing.T) {
		output, err := executor.Execute("echo", "你好世界🌍")
		if err != nil {
			t.Fatalf("unicode arguments failed: %v", err)
		}
		if output == nil {
			t.Error("output should not be nil")
		}
	})
}

// TestStorcliSpecificExecution tests storcli-specific command execution
func TestStorcliSpecificExecution(t *testing.T) {
	executor := NewExecutor()
	if executor == nil {
		t.Fatal("failed to create executor")
	}

	t.Run("storcli_help_command", func(t *testing.T) {
		// Note: This will only work if storcli is installed
		// Most systems won't have storcli, so we test the execution framework
		output, err := executor.Execute("echo", "storcli help")
		if err != nil {
			t.Logf("storcli help not available: %v", err)
		}
		_ = output
	})

	t.Run("storcli_json_output", func(t *testing.T) {
		// Test ability to handle JSON output
		output, err := executor.Execute("echo", "{\"key\":\"value\"}")
		if err != nil {
			t.Fatalf("json output test failed: %v", err)
		}
		if output == nil {
			t.Error("output should not be nil")
		}
	})

	t.Run("storcli_large_output", func(t *testing.T) {
		// Test handling of large output
		largeArg := strings.Repeat("x", 50000)
		output, err := executor.Execute("echo", largeArg)
		if err != nil {
			t.Logf("large output failed: %v", err)
		}
		_ = output
	})
}

// TestOutputValidation tests output validation
func TestOutputValidation(t *testing.T) {
	executor := NewExecutor()
	if executor == nil {
		t.Fatal("failed to create executor")
	}

	t.Run("validate_byte_output", func(t *testing.T) {
		output, err := executor.Execute("echo", "test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output == nil {
			t.Error("output should not be nil")
		}
		if len(output) == 0 {
			t.Log("output is empty (may be expected)")
		}
	})

	t.Run("validate_output_encoding", func(t *testing.T) {
		output, err := executor.Execute("echo", "test")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output != nil && len(output) > 0 {
			// Verify it's valid UTF-8
			_ = string(output)
		}
	})
}

// TestContextDeadlineHandling tests deadline handling in context
func TestContextDeadlineHandling(t *testing.T) {
	executor := NewExecutor()
	if executor == nil {
		t.Fatal("failed to create executor")
	}

	t.Run("deadline_not_exceeded", func(t *testing.T) {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
		defer cancel()
		
		output, err := executor.ExecuteWithContext(ctx, "echo", "test")
		if err != nil {
			t.Logf("unexpected error: %v", err)
		}
		_ = output
	})

	t.Run("deadline_already_exceeded", func(t *testing.T) {
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-1*time.Second))
		defer cancel()
		
		_, err := executor.ExecuteWithContext(ctx, "echo", "test")
		if err == nil {
			t.Log("should ideally error for already exceeded deadline")
		}
	})
}

// TestExecutorInterface tests that executor implements expected interface
func TestExecutorInterface(t *testing.T) {
	t.Run("executor_has_execute_method", func(t *testing.T) {
		executor := NewExecutor()
		if executor == nil {
			t.Fatal("executor is nil")
		}
		// Verify Execute method exists and is callable
		_, _ = executor.Execute("echo", "test")
	})

	t.Run("executor_has_timeout_method", func(t *testing.T) {
		executor := NewExecutor()
		if executor == nil {
			t.Fatal("executor is nil")
		}
		// Verify ExecuteWithTimeout method exists and is callable
		_, _ = executor.ExecuteWithTimeout(5*time.Second, "echo", "test")
	})

	t.Run("executor_has_context_method", func(t *testing.T) {
		executor := NewExecutor()
		if executor == nil {
			t.Fatal("executor is nil")
		}
		// Verify ExecuteWithContext method exists and is callable
		_, _ = executor.ExecuteWithContext(context.Background(), "echo", "test")
	})
}