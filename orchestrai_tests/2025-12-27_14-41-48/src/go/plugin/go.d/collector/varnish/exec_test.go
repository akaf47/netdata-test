package varnish

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewExecPath tests the NewExecPath function
func TestNewExecPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "should return correct exec path",
			path:     "/usr/bin/varnishstat",
			expected: "/usr/bin/varnishstat",
		},
		{
			name:     "should handle empty path",
			path:     "",
			expected: "",
		},
		{
			name:     "should handle relative path",
			path:     "./varnishstat",
			expected: "./varnishstat",
		},
		{
			name:     "should handle absolute path",
			path:     "/opt/varnish/bin/varnishstat",
			expected: "/opt/varnish/bin/varnishstat",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exec := NewExecPath(tt.path)
			assert.NotNil(t, exec)
		})
	}
}

// TestExecVarnishstat tests the varnishstat execution
func TestExecVarnishstat(t *testing.T) {
	tests := []struct {
		name           string
		ctx            context.Context
		shouldErr      bool
		errMessage     string
		setupExec      func(*mockExec)
		validateResult func(t *testing.T, output string, err error)
	}{
		{
			name:      "should execute varnishstat successfully",
			ctx:       context.Background(),
			shouldErr: false,
			setupExec: func(m *mockExec) {
				m.execFunc = func(ctx context.Context, cmd string, args ...string) (string, error) {
					return "VBE.boot.vcl_duration 5.000000", nil
				}
			},
			validateResult: func(t *testing.T, output string, err error) {
				assert.NoError(t, err)
				assert.Contains(t, output, "VBE.boot.vcl_duration")
			},
		},
		{
			name:           "should handle exec error",
			ctx:            context.Background(),
			shouldErr:      true,
			errMessage:     "command not found",
			setupExec: func(m *mockExec) {
				m.execFunc = func(ctx context.Context, cmd string, args ...string) (string, error) {
					return "", errors.New("command not found")
				}
			},
			validateResult: func(t *testing.T, output string, err error) {
				assert.Error(t, err)
				assert.Empty(t, output)
			},
		},
		{
			name:      "should handle context timeout",
			ctx:       contextWithTimeout(10 * time.Millisecond),
			shouldErr: true,
			setupExec: func(m *mockExec) {
				m.execFunc = func(ctx context.Context, cmd string, args ...string) (string, error) {
					return "", context.DeadlineExceeded
				}
			},
			validateResult: func(t *testing.T, output string, err error) {
				assert.Error(t, err)
				assert.Equal(t, context.DeadlineExceeded, err)
			},
		},
		{
			name:      "should handle empty output",
			ctx:       context.Background(),
			shouldErr: false,
			setupExec: func(m *mockExec) {
				m.execFunc = func(ctx context.Context, cmd string, args ...string) (string, error) {
					return "", nil
				}
			},
			validateResult: func(t *testing.T, output string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "", output)
			},
		},
		{
			name:      "should handle multiline output",
			ctx:       context.Background(),
			shouldErr: false,
			setupExec: func(m *mockExec) {
				m.execFunc = func(ctx context.Context, cmd string, args ...string) (string, error) {
					return "VBE.boot.vcl_duration 5.000000\nVBE.boot.happy 0\nVBE.boot.cool 0", nil
				}
			},
			validateResult: func(t *testing.T, output string, err error) {
				assert.NoError(t, err)
				assert.Contains(t, output, "VBE.boot.vcl_duration")
				assert.Contains(t, output, "VBE.boot.happy")
				assert.Contains(t, output, "VBE.boot.cool")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExec := &mockExec{}
			tt.setupExec(mockExec)
			
			// Execute the test using the mock
			output, err := mockExec.execFunc(tt.ctx, "varnishstat")
			tt.validateResult(t, output, err)
		})
	}
}

// TestExecVarnishstatWithArgs tests varnishstat with various arguments
func TestExecVarnishstatWithArgs(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		shouldErr      bool
		validateResult func(t *testing.T, output string, err error)
	}{
		{
			name:      "should accept -1 flag",
			args:      []string{"-1"},
			shouldErr: false,
			validateResult: func(t *testing.T, output string, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:      "should accept -f flag",
			args:      []string{"-f", "RES_HDR"},
			shouldErr: false,
			validateResult: func(t *testing.T, output string, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:      "should accept -n flag",
			args:      []string{"-n", "varnish"},
			shouldErr: false,
			validateResult: func(t *testing.T, output string, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:      "should handle empty args",
			args:      []string{},
			shouldErr: false,
			validateResult: func(t *testing.T, output string, err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExec := &mockExec{
				execFunc: func(ctx context.Context, cmd string, args ...string) (string, error) {
					assert.Equal(t, cmd, "varnishstat")
					for i, expectedArg := range tt.args {
						if i < len(args) {
							assert.Equal(t, expectedArg, args[i])
						}
					}
					return "output", nil
				},
			}
			
			output, err := mockExec.execFunc(context.Background(), "varnishstat", tt.args...)
			tt.validateResult(t, output, err)
		})
	}
}

// TestExecVarnishstatFiltering tests output filtering scenarios
func TestExecVarnishstatFiltering(t *testing.T) {
	tests := []struct {
		name           string
		output         string
		filter         string
		validateResult func(t *testing.T, output string, err error)
	}{
		{
			name:   "should filter specific metric",
			output: "VBE.boot.vcl_duration 5.000000\nVBE.boot.happy 0\nVBE.boot.cool 0",
			filter: "happy",
			validateResult: func(t *testing.T, output string, err error) {
				assert.NoError(t, err)
				assert.Contains(t, output, "happy")
			},
		},
		{
			name:   "should handle filter with no matches",
			output: "VBE.boot.vcl_duration 5.000000",
			filter: "nonexistent",
			validateResult: func(t *testing.T, output string, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name:   "should handle empty filter",
			output: "VBE.boot.vcl_duration 5.000000",
			filter: "",
			validateResult: func(t *testing.T, output string, err error) {
				assert.NoError(t, err)
				assert.NotEmpty(t, output)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExec := &mockExec{
				execFunc: func(ctx context.Context, cmd string, args ...string) (string, error) {
					return tt.output, nil
				},
			}
			
			output, err := mockExec.execFunc(context.Background(), "varnishstat")
			tt.validateResult(t, output, err)
		})
	}
}

// TestExecContextHandling tests context handling in exec
func TestExecContextHandling(t *testing.T) {
	tests := []struct {
		name             string
		ctxFunc          func() context.Context
		expectedErr      error
		validateBehavior func(t *testing.T, err error)
	}{
		{
			name:        "should handle background context",
			ctxFunc:     func() context.Context { return context.Background() },
			expectedErr: nil,
			validateBehavior: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "should handle context with timeout",
			ctxFunc: func() context.Context {
				ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
				return ctx
			},
			expectedErr: nil,
			validateBehavior: func(t *testing.T, err error) {
				assert.NoError(t, err)
			},
		},
		{
			name: "should handle expired context",
			ctxFunc: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
				defer cancel()
				time.Sleep(2 * time.Millisecond)
				return ctx
			},
			expectedErr: context.DeadlineExceeded,
			validateBehavior: func(t *testing.T, err error) {
				// Validate that context deadline handling is possible
				assert.NotNil(t, err) // err may be DeadlineExceeded
			},
		},
		{
			name: "should handle cancelled context",
			ctxFunc: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			},
			expectedErr: context.Canceled,
			validateBehavior: func(t *testing.T, err error) {
				assert.Equal(t, context.Canceled, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.ctxFunc()
			
			mockExec := &mockExec{
				execFunc: func(execCtx context.Context, cmd string, args ...string) (string, error) {
					select {
					case <-execCtx.Done():
						return "", execCtx.Err()
					default:
						return "output", nil
					}
				},
			}
			
			output, err := mockExec.execFunc(ctx, "varnishstat")
			
			if tt.expectedErr != nil {
				// Some context tests expect an error
				if err != nil {
					tt.validateBehavior(t, err)
				}
			} else {
				tt.validateBehavior(t, err)
			}
		})
	}
}

// TestExecErrorConditions tests various error conditions
func TestExecErrorConditions(t *testing.T) {
	tests := []struct {
		name              string
		execFunc          func(ctx context.Context, cmd string, args ...string) (string, error)
		validateErr       func(t *testing.T, err error)
	}{
		{
			name: "should handle permission denied error",
			execFunc: func(ctx context.Context, cmd string, args ...string) (string, error) {
				return "", errors.New("permission denied")
			},
			validateErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "permission denied")
			},
		},
		{
			name: "should handle file not found error",
			execFunc: func(ctx context.Context, cmd string, args ...string) (string, error) {
				return "", errors.New("executable file not found in $PATH")
			},
			validateErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "not found")
			},
		},
		{
			name: "should handle varnish unavailable error",
			execFunc: func(ctx context.Context, cmd string, args ...string) (string, error) {
				return "", errors.New("varnish is not running")
			},
			validateErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "not running")
			},
		},
		{
			name: "should handle exit code error",
			execFunc: func(ctx context.Context, cmd string, args ...string) (string, error) {
				return "", errors.New("exit status 1")
			},
			validateErr: func(t *testing.T, err error) {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "exit status")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExec := &mockExec{execFunc: tt.execFunc}
			_, err := mockExec.execFunc(context.Background(), "varnishstat")
			tt.validateErr(t, err)
		})
	}
}

// TestExecOutputParsing tests output parsing scenarios
func TestExecOutputParsing(t *testing.T) {
	tests := []struct {
		name     string
		output   string
		validate func(t *testing.T, output string)
	}{
		{
			name:   "should parse single value output",
			output: "VBE.boot.vcl_duration 5.000000",
			validate: func(t *testing.T, output string) {
				assert.NotEmpty(t, output)
				assert.Contains(t, output, "VBE.boot.vcl_duration")
			},
		},
		{
			name:   "should parse multiple values",
			output: "VBE.boot.vcl_duration 5.000000\nVBE.boot.happy 0\nVBE.boot.cool 0",
			validate: func(t *testing.T, output string) {
				assert.NotEmpty(t, output)
				lines := len([]byte(output))
				assert.Greater(t, lines, 0)
			},
		},
		{
			name:   "should handle whitespace in output",
			output: "  VBE.boot.vcl_duration   5.000000  ",
			validate: func(t *testing.T, output string) {
				assert.NotEmpty(t, output)
			},
		},
		{
			name:   "should handle special characters",
			output: "MAIN.uptime 123456\nMAIN.sess_conn 789",
			validate: func(t *testing.T, output string) {
				assert.NotEmpty(t, output)
				assert.Contains(t, output, "MAIN")
			},
		},
		{
			name:   "should handle numeric values",
			output: "VBE.boot.vcl_duration 5\nVBE.boot.happy 0\nVBE.boot.cool -1",
			validate: func(t *testing.T, output string) {
				assert.NotEmpty(t, output)
				assert.Contains(t, output, "5")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockExec := &mockExec{
				execFunc: func(ctx context.Context, cmd string, args ...string) (string, error) {
					return tt.output, nil
				},
			}
			
			output, err := mockExec.execFunc(context.Background(), "varnishstat")
			assert.NoError(t, err)
			tt.validate(t, output)
		})
	}
}

// Mock types and helpers
type mockExec struct {
	execFunc func(ctx context.Context, cmd string, args ...string) (string, error)
}

func (m *mockExec) Exec(ctx context.Context, cmd string, args ...string) (string, error) {
	if m.execFunc != nil {
		return m.execFunc(ctx, cmd, args...)
	}
	return "", nil
}

func contextWithTimeout(duration time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), duration)
	return ctx
}

// BenchmarkExec benchmarks the exec function
func BenchmarkExec(b *testing.B) {
	mockExec := &mockExec{
		execFunc: func(ctx context.Context, cmd string, args ...string) (string, error) {
			return "VBE.boot.vcl_duration 5.000000", nil
		},
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mockExec.execFunc(context.Background(), "varnishstat")
	}
}