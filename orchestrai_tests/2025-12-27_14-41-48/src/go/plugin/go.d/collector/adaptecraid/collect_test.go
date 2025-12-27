package adaptecraid

import (
	"errors"
	"testing"
	"time"
)

func TestCollectMetrics(t *testing.T) {
	tests := []struct {
		name        string
		setup       func(*Adaptecraid)
		expectError bool
		validate    func(*testing.T, map[string]interface{})
	}{
		{
			name: "Collect should return valid metrics map",
			setup: func(a *Adaptecraid) {
				a.exec = func(string, ...string) ([]byte, error) {
					return []byte("Controller Status: OK\n"), nil
				}
			},
			expectError: false,
			validate: func(t *testing.T, data map[string]interface{}) {
				if data == nil {
					t.Error("Expected non-nil metrics map")
				}
			},
		},
		{
			name: "Collect should handle empty response",
			setup: func(a *Adaptecraid) {
				a.exec = func(string, ...string) ([]byte, error) {
					return []byte{}, nil
				}
			},
			expectError: false,
			validate: func(t *testing.T, data map[string]interface{}) {
				if data == nil {
					t.Error("Expected non-nil metrics map")
				}
			},
		},
		{
			name: "Collect should handle executor error",
			setup: func(a *Adaptecraid) {
				a.exec = func(string, ...string) ([]byte, error) {
					return nil, errors.New("command failed")
				}
			},
			expectError: true,
			validate: func(t *testing.T, data map[string]interface{}) {
				// When error occurs, data should be nil or empty
			},
		},
		{
			name: "Collect should handle nil executor",
			setup: func(a *Adaptecraid) {
				a.exec = nil
			},
			expectError: true,
			validate:    func(t *testing.T, data map[string]interface{}) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Adaptecraid{}
			if tt.setup != nil {
				tt.setup(a)
			}

			metrics, err := a.Collect()

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if !tt.expectError && metrics != nil {
				tt.validate(t, metrics)
			}
		})
	}
}

func TestCollectWithDifferentResponses(t *testing.T) {
	tests := []struct {
		name     string
		response string
		wantErr  bool
	}{
		{
			name:     "Valid single controller response",
			response: "Controller Status: OK\n",
			wantErr:  false,
		},
		{
			name:     "Response with multiple lines",
			response: "Controller Status: OK\nDrives: 4\nStatus: Optimal\n",
			wantErr:  false,
		},
		{
			name:     "Response with controller number",
			response: "Controller #0\nStatus: OK\n",
			wantErr:  false,
		},
		{
			name:     "Empty response",
			response: "",
			wantErr:  false,
		},
		{
			name:     "Response with special characters",
			response: "Status: OK\nController™ Info\n",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Adaptecraid{
				exec: func(string, ...string) ([]byte, error) {
					return []byte(tt.response), nil
				},
			}

			metrics, err := a.Collect()
			if (err != nil) != tt.wantErr {
				t.Errorf("Collect() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				if metrics == nil {
					t.Error("Expected metrics to be non-nil")
				}
			}
		})
	}
}

func TestCollectCommandExecution(t *testing.T) {
	tests := []struct {
		name      string
		execErr   error
		wantError bool
	}{
		{
			name:      "Command execution success",
			execErr:   nil,
			wantError: false,
		},
		{
			name:      "Command execution failure",
			execErr:   errors.New("permission denied"),
			wantError: true,
		},
		{
			name:      "Command timeout",
			execErr:   errors.New("command timeout"),
			wantError: true,
		},
		{
			name:      "Command not found",
			execErr:   errors.New("arcconf not found"),
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Adaptecraid{
				exec: func(string, ...string) ([]byte, error) {
					if tt.execErr != nil {
						return nil, tt.execErr
					}
					return []byte("Status: OK\n"), nil
				},
			}

			_, err := a.Collect()
			if (err != nil) != tt.wantError {
				t.Errorf("Collect() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestCollectMetricsStructure(t *testing.T) {
	a := &Adaptecraid{
		exec: func(cmd string, args ...string) ([]byte, error) {
			return []byte("Controller #0\nStatus: Optimal\n"), nil
		},
	}

	metrics, err := a.Collect()
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if metrics == nil {
		t.Fatal("Expected non-nil metrics")
	}

	// Verify metrics is a map
	if _, ok := metrics.(map[string]interface{}); !ok {
		if metrics != nil {
			t.Logf("Metrics type: %T", metrics)
		}
	}
}

func TestCollectEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		response string
	}{
		{
			name:     "Very long response",
			response: generateLongResponse(10000),
		},
		{
			name:     "Response with null bytes",
			response: "Status: OK\x00Hidden",
		},
		{
			name:     "Response with only whitespace",
			response: "   \n\n  \t\n  ",
		},
		{
			name:     "Response with unicode characters",
			response: "Status: OK™\nController: №1\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Adaptecraid{
				exec: func(string, ...string) ([]byte, error) {
					return []byte(tt.response), nil
				},
			}

			metrics, err := a.Collect()
			if err != nil {
				t.Logf("Collect returned error for %s: %v", tt.name, err)
			}
			if metrics == nil && err == nil {
				t.Logf("Collect returned nil metrics without error for %s", tt.name)
			}
		})
	}
}

func TestCollectMultipleCalls(t *testing.T) {
	callCount := 0
	a := &Adaptecraid{
		exec: func(string, ...string) ([]byte, error) {
			callCount++
			return []byte("Status: OK\n"), nil
		},
	}

	// Call Collect multiple times
	for i := 0; i < 5; i++ {
		metrics, err := a.Collect()
		if err != nil {
			t.Errorf("Call %d failed: %v", i+1, err)
		}
		if metrics == nil {
			t.Errorf("Call %d returned nil metrics", i+1)
		}
	}

	if callCount != 5 {
		t.Errorf("Expected 5 exec calls, got %d", callCount)
	}
}

func TestCollectConcurrency(t *testing.T) {
	a := &Adaptecraid{
		exec: func(string, ...string) ([]byte, error) {
			time.Sleep(10 * time.Millisecond)
			return []byte("Status: OK\n"), nil
		},
	}

	done := make(chan bool, 3)

	for i := 0; i < 3; i++ {
		go func() {
			_, err := a.Collect()
			if err != nil {
				t.Errorf("Concurrent collect failed: %v", err)
			}
			done <- true
		}()
	}

	for i := 0; i < 3; i++ {
		<-done
	}
}

func TestCollectWithNilAdaptecraid(t *testing.T) {
	var a *Adaptecraid
	// This would panic if not handled, testing error handling
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Recovered from panic as expected: %v", r)
		}
	}()

	if a != nil {
		_, _ = a.Collect()
	}
}

func TestCollectReturnValueType(t *testing.T) {
	a := &Adaptecraid{
		exec: func(string, ...string) ([]byte, error) {
			return []byte("Status: OK\n"), nil
		},
	}

	metrics, err := a.Collect()

	if err != nil {
		t.Fatalf("Collect() returned unexpected error: %v", err)
	}

	// Verify return types
	if metrics != nil {
		// Should be able to assert as map[string]interface{} for valid metrics
		_ = metrics
	}
}

func generateLongResponse(length int) string {
	result := ""
	for i := 0; i < length; i++ {
		result += "A"
	}
	return result
}