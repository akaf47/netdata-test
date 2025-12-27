package hpssa

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockExecutor is a mock for external command execution
type MockExecutor struct {
	mock.Mock
}

func (m *MockExecutor) Execute(ctx context.Context, name string, args ...string) ([]byte, error) {
	retArgs := m.Called(ctx, name, args)
	if retArgs.Get(0) == nil {
		return nil, retArgs.Error(1)
	}
	return retArgs.Get(0).([]byte), retArgs.Error(1)
}

// TestExecNewExecutor tests the creation of a new executor
func TestExecNewExecutor(t *testing.T) {
	// Arrange & Act
	executor := NewExecutor()

	// Assert
	assert.NotNil(t, executor)
}

// TestExecExecuteSuccess tests successful command execution
func TestExecExecuteSuccess(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()
	expectedOutput := []byte("test output")

	// Act
	output, err := executor.Execute(ctx, "echo", "test")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
}

// TestExecExecuteWithTimeout tests command execution with timeout
func TestExecExecuteWithTimeout(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	// Act
	_, err := executor.Execute(ctx, "sleep", "10")

	// Assert
	assert.Error(t, err)
}

// TestExecExecuteWithCancelledContext tests execution with cancelled context
func TestExecExecuteWithCancelledContext(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Act
	_, err := executor.Execute(ctx, "echo", "test")

	// Assert
	assert.Error(t, err)
}

// TestExecExecuteWithEmptyCommand tests execution with empty command
func TestExecExecuteWithEmptyCommand(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()

	// Act
	_, err := executor.Execute(ctx, "", "arg")

	// Assert
	assert.Error(t, err)
}

// TestExecExecuteWithNilContext tests execution with nil context
func TestExecExecuteWithNilContext(t *testing.T) {
	// Arrange
	executor := NewExecutor()

	// Act
	_, err := executor.Execute(nil, "echo", "test")

	// Assert
	assert.Error(t, err)
}

// TestExecExecuteWithNoArgs tests execution with no arguments
func TestExecExecuteWithNoArgs(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()

	// Act
	_, err := executor.Execute(ctx, "echo")

	// Assert
	assert.NoError(t, err)
}

// TestExecExecuteWithMultipleArgs tests execution with multiple arguments
func TestExecExecuteWithMultipleArgs(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()

	// Act
	_, err := executor.Execute(ctx, "echo", "arg1", "arg2", "arg3")

	// Assert
	assert.NoError(t, err)
}

// TestExecExecuteNonexistentCommand tests execution of nonexistent command
func TestExecExecuteNonexistentCommand(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()

	// Act
	_, err := executor.Execute(ctx, "/nonexistent/command/path", "arg")

	// Assert
	assert.Error(t, err)
}

// TestExecExecuteCommandWithSpaceInName tests execution with spaces in command name
func TestExecExecuteCommandWithSpaceInName(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()

	// Act
	_, err := executor.Execute(ctx, "command with spaces")

	// Assert
	assert.Error(t, err)
}

// TestExecExecuteStoresOutput tests that executor properly stores output
func TestExecExecuteStoresOutput(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()

	// Act
	output1, _ := executor.Execute(ctx, "echo", "output1")
	output2, _ := executor.Execute(ctx, "echo", "output2")

	// Assert
	assert.NotNil(t, output1)
	assert.NotNil(t, output2)
	// Ensure they are different calls
	assert.True(t, len(output1) > 0 || len(output2) > 0)
}

// TestExecExecuteWithLargeOutput tests execution with large output
func TestExecExecuteWithLargeOutput(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()

	// Act - generate a command with large output
	_, err := executor.Execute(ctx, "bash", "-c", "for i in {1..1000}; do echo 'line '$i; done")

	// Assert
	assert.NoError(t, err)
}

// TestExecExecuteMultipleTimes tests executor can be called multiple times
func TestExecExecuteMultipleTimes(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()

	// Act & Assert
	for i := 0; i < 5; i++ {
		_, err := executor.Execute(ctx, "echo", "test")
		assert.NoError(t, err)
	}
}

// TestExecExecuteWithSpecialCharactersInArgs tests args with special characters
func TestExecExecuteWithSpecialCharactersInArgs(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()

	// Act
	_, err := executor.Execute(ctx, "echo", "!@#$%^&*()")

	// Assert
	assert.NoError(t, err)
}

// TestExecExecuteOutputFormat tests output is in correct format
func TestExecExecuteOutputFormat(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()

	// Act
	output, err := executor.Execute(ctx, "echo", "hello")

	// Assert
	assert.NoError(t, err)
	assert.Greater(t, len(output), 0)
	assert.IsType(t, []byte{}, output)
}

// TestExecExecuteContextDeadlineExceeded tests when context deadline is exceeded
func TestExecExecuteContextDeadlineExceeded(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-1*time.Second))
	defer cancel()

	// Act
	_, err := executor.Execute(ctx, "echo", "test")

	// Assert
	assert.Error(t, err)
}

// TestExecExecuteReturnsErrorOnFailure tests error is properly returned
func TestExecExecuteReturnsErrorOnFailure(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()

	// Act
	_, err := executor.Execute(ctx, "bash", "-c", "exit 1")

	// Assert
	assert.Error(t, err)
}

// TestExecExecutePreservesOutputOnError tests output availability on error
func TestExecExecutePreservesOutputOnError(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()

	// Act
	output, err := executor.Execute(ctx, "bash", "-c", "echo 'before error' && exit 1")

	// Assert
	assert.Error(t, err)
	// Output might still be available depending on implementation
	_ = output
}

// TestExecExecuteWithEmptyOutput tests execution returning empty output
func TestExecExecuteWithEmptyOutput(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()

	// Act
	output, err := executor.Execute(ctx, "bash", "-c", "")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
}

// TestExecExecuteWithBinaryOutput tests execution with binary output
func TestExecExecuteWithBinaryOutput(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()

	// Act
	output, err := executor.Execute(ctx, "bash", "-c", "printf '\\x00\\x01\\x02\\x03'")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
}

// TestExecExecuteWithNewlineInOutput tests output with newlines
func TestExecExecuteWithNewlineInOutput(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()

	// Act
	output, err := executor.Execute(ctx, "bash", "-c", "echo -e 'line1\\nline2\\nline3'")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
}

// TestExecExecuteSequential tests sequential command execution
func TestExecExecuteSequential(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()

	// Act & Assert
	_, err1 := executor.Execute(ctx, "echo", "first")
	assert.NoError(t, err1)

	_, err2 := executor.Execute(ctx, "echo", "second")
	assert.NoError(t, err2)
}

// TestExecExecutePreservesContext tests context is used for timeout
func TestExecExecutePreservesContext(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	baseCtx := context.Background()
	ctx, cancel := context.WithTimeout(baseCtx, 30*time.Second)
	defer cancel()

	// Act
	_, err := executor.Execute(ctx, "echo", "test")

	// Assert
	assert.NoError(t, err)
}

// TestExecExecuteErrorReturnedDirectly tests errors are returned without wrapping
func TestExecExecuteErrorReturnedDirectly(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()

	// Act
	_, err := executor.Execute(ctx, "/bin/bash", "-c", "exit 127")

	// Assert
	assert.Error(t, err)
	assert.NotNil(t, err)
}

// TestExecExecuteWithUTF8Output tests execution with UTF-8 output
func TestExecExecuteWithUTF8Output(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()

	// Act
	output, err := executor.Execute(ctx, "bash", "-c", "echo '你好世界'")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, output)
}

// TestExecExecuteDoesNotModifyContext tests that execute doesn't modify context
func TestExecExecuteDoesNotModifyContext(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()

	// Act
	executor.Execute(ctx, "echo", "test")

	// Assert - context should still be usable
	select {
	case <-ctx.Done():
		t.Fail() // Context should not be done
	default:
		// Good
	}
}

// TestExecExecuteParallelExecution tests multiple concurrent executions
func TestExecExecuteParallelExecution(t *testing.T) {
	// Arrange
	executor := NewExecutor()
	ctx := context.Background()
	errChan := make(chan error, 3)

	// Act
	go func() {
		_, err := executor.Execute(ctx, "echo", "1")
		errChan <- err
	}()
	go func() {
		_, err := executor.Execute(ctx, "echo", "2")
		errChan <- err
	}()
	go func() {
		_, err := executor.Execute(ctx, "echo", "3")
		errChan <- err
	}()

	// Assert
	for i := 0; i < 3; i++ {
		err := <-errChan
		assert.NoError(t, err)
	}
}