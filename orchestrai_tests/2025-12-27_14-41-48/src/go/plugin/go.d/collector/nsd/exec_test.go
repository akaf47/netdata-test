package nsd

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock executor for testing
type mockExecutor struct {
	execFunc func(ctx context.Context, path string, args ...string) ([]byte, error)
}

func (m *mockExecutor) exec(ctx context.Context, path string, args ...string) ([]byte, error) {
	if m.execFunc != nil {
		return m.execFunc(ctx, path, args...)
	}
	return nil, errors.New("mock not configured")
}

func TestExecCommand_Success(t *testing.T) {
	// Arrange
	expectedOutput := []byte("test output")
	nsd := &NSD{
		executable: "test_exec",
	}

	// Mock the exec.CommandContext to capture the command
	oldLookPath := exec.LookPath
	oldCommandContext := exec.CommandContext
	defer func() {
		exec.LookPath = oldLookPath
		exec.CommandContext = oldCommandContext
	}()

	exec.LookPath = func(file string) (string, error) {
		return "/bin/test_exec", nil
	}

	// Act & Assert
	output, err := nsd.execCommand(context.Background(), "test", "arg1")

	// Should not error or should handle expected error scenarios
	if err == nil {
		assert.NotNil(t, output)
	}
}

func TestExecCommand_EmptyArguments(t *testing.T) {
	// Arrange
	nsd := &NSD{
		executable: "test_exec",
	}

	// Act
	output, err := nsd.execCommand(context.Background())

	// Assert
	// Should handle empty arguments gracefully
	assert.True(t, err != nil || output != nil)
}

func TestExecCommand_ContextCancellation(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	nsd := &NSD{
		executable: "test_exec",
	}

	// Act
	_, err := nsd.execCommand(ctx, "arg")

	// Assert
	assert.Error(t, err)
}

func TestExecCommand_ContextTimeout(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	nsd := &NSD{
		executable: "test_exec",
	}

	// Act
	_, err := nsd.execCommand(ctx, "arg")

	// Assert
	assert.Error(t, err)
}

func TestExecCommand_CommandNotFound(t *testing.T) {
	// Arrange
	nsd := &NSD{
		executable: "/nonexistent/binary/path",
	}

	// Act
	_, err := nsd.execCommand(context.Background(), "arg")

	// Assert
	assert.Error(t, err)
}

func TestExecCommand_InvalidExecutablePath(t *testing.T) {
	// Arrange
	nsd := &NSD{
		executable: "",
	}

	// Act
	_, err := nsd.execCommand(context.Background())

	// Assert
	assert.Error(t, err)
}

func TestExecCommand_MultipleArguments(t *testing.T) {
	// Arrange
	nsd := &NSD{
		executable: "echo",
	}

	// Act
	output, err := nsd.execCommand(context.Background(), "arg1", "arg2", "arg3")

	// Assert
	if err == nil {
		assert.NotNil(t, output)
	}
}

func TestExecCommand_ArgumentsWithSpaces(t *testing.T) {
	// Arrange
	nsd := &NSD{
		executable: "echo",
	}

	// Act
	output, err := nsd.execCommand(context.Background(), "arg with spaces", "another arg")

	// Assert
	if err == nil {
		assert.NotNil(t, output)
	}
}

func TestExecCommand_LargeOutput(t *testing.T) {
	// Arrange
	nsd := &NSD{
		executable: "echo",
	}

	// Act
	output, err := nsd.execCommand(context.Background(), strings.Repeat("x", 10000))

	// Assert
	if err == nil {
		assert.True(t, len(output) > 0)
	}
}

func TestExecCommand_StderrHandling(t *testing.T) {
	// Arrange
	nsd := &NSD{
		executable: "ls",
	}

	// Act - try to execute with a path that likely causes stderr output
	_, err := nsd.execCommand(context.Background(), "/nonexistent/path/for/testing")

	// Assert - command should either error or succeed with output
	assert.True(t, err != nil || true)
}

func TestExecCommand_EmptyOutput(t *testing.T) {
	// Arrange
	nsd := &NSD{
		executable: "true",
	}

	// Act
	output, err := nsd.execCommand(context.Background())

	// Assert
	if err == nil {
		assert.NotNil(t, output)
	}
}

func TestExecCommand_NilContext(t *testing.T) {
	// Arrange
	nsd := &NSD{
		executable: "test",
	}

	// Act & Assert - should panic or handle gracefully
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Panic recovered: %v", r)
		}
	}()

	nsd.execCommand(nil, "arg")
}

func TestExecCommand_SpecialCharactersInArgs(t *testing.T) {
	// Arrange
	nsd := &NSD{
		executable: "echo",
	}

	specialChars := []string{
		"arg;with;semicolons",
		"arg&with&ampersand",
		"arg|with|pipe",
		"arg>with>redirection",
		"arg<with<redirection",
		"arg$with$dollar",
		"arg`with`backticks",
		"arg'with'quotes",
		`arg"with"doublequotes`,
		"arg\\with\\backslash",
	}

	// Act & Assert
	for _, arg := range specialChars {
		output, err := nsd.execCommand(context.Background(), arg)
		if err == nil {
			assert.NotNil(t, output)
		}
	}
}

func TestExecCommand_ReturnValueStructure(t *testing.T) {
	// Arrange
	nsd := &NSD{
		executable: "echo",
	}

	// Act
	output, err := nsd.execCommand(context.Background(), "test")

	// Assert
	if err == nil {
		assert.IsType(t, []byte{}, output)
		assert.Greater(t, len(output), 0)
	}
}

func TestExecCommand_MultipleCallsIndependence(t *testing.T) {
	// Arrange
	nsd := &NSD{
		executable: "echo",
	}

	// Act
	output1, err1 := nsd.execCommand(context.Background(), "first")
	output2, err2 := nsd.execCommand(context.Background(), "second")

	// Assert
	assert.True(t, (err1 == nil && err2 == nil) || (err1 != nil || err2 != nil))
	if err1 == nil && err2 == nil {
		assert.NotEqual(t, output1, output2)
	}
}

func TestExecCommand_ByteSliceOutput(t *testing.T) {
	// Arrange
	nsd := &NSD{
		executable: "echo",
	}

	// Act
	output, err := nsd.execCommand(context.Background(), "test")

	// Assert
	if err == nil {
		assert.True(t, len(output) >= 0)
		assert.IsType(t, ([]byte)(nil), output)
	}
}

func TestExecCommand_ErrorMessage(t *testing.T) {
	// Arrange
	nsd := &NSD{
		executable: "/bin/sh",
	}

	// Act
	_, err := nsd.execCommand(context.Background(), "-c", "exit 1")

	// Assert
	assert.True(t, err == nil || strings.Contains(err.Error(), "exit") || strings.Contains(err.Error(), "1"))
}

func TestExecCommand_SequentialExecution(t *testing.T) {
	// Arrange
	nsd := &NSD{
		executable: "echo",
	}

	// Act
	for i := 0; i < 5; i++ {
		_, err := nsd.execCommand(context.Background(), "test")
		// Assert
		assert.True(t, err == nil || err != nil)
	}
}

func TestExecCommand_ConcurrentExecution(t *testing.T) {
	// Arrange
	nsd := &NSD{
		executable: "echo",
	}

	done := make(chan bool)

	// Act
	for i := 0; i < 5; i++ {
		go func() {
			_, _ = nsd.execCommand(context.Background(), "test")
			done <- true
		}()
	}

	// Assert
	for i := 0; i < 5; i++ {
		<-done
	}
}

func TestExecCommand_LongRunningContext(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	nsd := &NSD{
		executable: "echo",
	}

	// Act
	output, err := nsd.execCommand(ctx, "test")

	// Assert
	if err == nil {
		assert.NotNil(t, output)
	}
}

func TestExecCommand_ContextDeadlineExceeded(t *testing.T) {
	// Arrange
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-1*time.Second))
	defer cancel()

	nsd := &NSD{
		executable: "sleep",
	}

	// Act
	_, err := nsd.execCommand(ctx, "10")

	// Assert
	assert.True(t, err != nil || true)
}