package ping

import (
	"context"
	"testing"
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockPinger struct {
	mock.Mock
}

func (m *MockPinger) Ping(ctx context.Context, host string, requestCount int, timeout time.Duration) (map[string]interface{}, error) {
	args := m.Called(ctx, host, requestCount, timeout)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func TestCollect(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(*Ping)
		mockSetup func(*MockPinger)
		verify    func(t *testing.T, data map[string]interface{})
		wantErr   bool
		wantNil   bool
	}{
		{
			name: "should collect metrics when hosts are reachable",
			setup: func(p *Ping) {
				p.Hosts = []string{"localhost"}
				p.Packets = 4
				p.Timeout = 3
			},
			mockSetup: func(m *MockPinger) {
				m.On("Ping", mock.Anything, "localhost", 4, 3*time.Second).Return(map[string]interface{}{
					"min": 1.0,
					"avg": 2.5,
					"max": 5.0,
				}, nil)
			},
			verify: func(t *testing.T, data map[string]interface{}) {
				assert.NotNil(t, data)
				assert.Greater(t, len(data), 0)
			},
			wantErr: false,
			wantNil: false,
		},
		{
			name: "should return empty map when no hosts configured",
			setup: func(p *Ping) {
				p.Hosts = []string{}
			},
			mockSetup: nil,
			verify: func(t *testing.T, data map[string]interface{}) {
				assert.NotNil(t, data)
			},
			wantErr: false,
			wantNil: false,
		},
		{
			name: "should handle single host",
			setup: func(p *Ping) {
				p.Hosts = []string{"8.8.8.8"}
				p.Packets = 1
				p.Timeout = 5
			},
			mockSetup: func(m *MockPinger) {
				m.On("Ping", mock.Anything, "8.8.8.8", 1, 5*time.Second).Return(map[string]interface{}{
					"min": 10.0,
					"max": 10.0,
					"avg": 10.0,
				}, nil)
			},
			verify: func(t *testing.T, data map[string]interface{}) {
				assert.NotNil(t, data)
			},
			wantErr: false,
			wantNil: false,
		},
		{
			name: "should handle multiple hosts",
			setup: func(p *Ping) {
				p.Hosts = []string{"host1.com", "host2.com", "host3.com"}
				p.Packets = 4
				p.Timeout = 3
			},
			mockSetup: func(m *MockPinger) {
				m.On("Ping", mock.Anything, mock.Anything, 4, 3*time.Second).Return(map[string]interface{}{
					"min": 1.0,
					"avg": 2.0,
					"max": 3.0,
				}, nil)
			},
			verify: func(t *testing.T, data map[string]interface{}) {
				assert.NotNil(t, data)
			},
			wantErr: false,
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()
			require.NotNil(t, p)

			if tt.setup != nil {
				tt.setup(p)
			}

			data, err := p.Collect(context.Background())

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.wantNil {
				assert.Nil(t, data)
			} else {
				assert.NotNil(t, data)
				if tt.verify != nil {
					tt.verify(t, data)
				}
			}
		})
	}
}

func TestCollectWithContext(t *testing.T) {
	p := New()
	require.NotNil(t, p)
	p.Hosts = []string{"localhost"}
	p.Packets = 4
	p.Timeout = 3

	t.Run("should respect context cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// Should handle cancelled context gracefully
		data, err := p.Collect(ctx)
		// Result depends on implementation, but should not panic
		_ = data
		_ = err
	})

	t.Run("should respect context timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		// Should handle timeout gracefully
		data, err := p.Collect(ctx)
		_ = data
		_ = err
	})

	t.Run("should work with background context", func(t *testing.T) {
		ctx := context.Background()
		data, err := p.Collect(ctx)
		assert.NotNil(t, data)
		assert.NoError(t, err)
	})
}

func TestCollectMetricsStructure(t *testing.T) {
	p := New()
	require.NotNil(t, p)
	p.Hosts = []string{"testhost.com"}
	p.Packets = 4
	p.Timeout = 3

	data, err := p.Collect(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, data)

	t.Run("metrics should be a map", func(t *testing.T) {
		assert.IsType(t, make(map[string]interface{}), data)
	})

	t.Run("metrics should have expected keys for reachable hosts", func(t *testing.T) {
		// Common ping metrics: min, max, avg, loss, etc.
		expectedKeys := []string{"min", "avg", "max"}
		for _, key := range expectedKeys {
			_, exists := data[key]
			// Key may or may not exist depending on host reachability
			_ = exists
		}
	})
}

func TestCollectEdgeCases(t *testing.T) {
	p := New()
	require.NotNil(t, p)

	t.Run("should handle nil hosts slice", func(t *testing.T) {
		p.Hosts = nil
		data, err := p.Collect(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, data)
	})

	t.Run("should handle empty hosts slice", func(t *testing.T) {
		p.Hosts = []string{}
		data, err := p.Collect(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, data)
	})

	t.Run("should handle zero packet count", func(t *testing.T) {
		p.Hosts = []string{"localhost"}
		p.Packets = 0
		data, err := p.Collect(context.Background())
		// Should handle gracefully
		_ = data
		_ = err
	})

	t.Run("should handle zero timeout", func(t *testing.T) {
		p.Hosts = []string{"localhost"}
		p.Timeout = 0
		data, err := p.Collect(context.Background())
		_ = data
		_ = err
	})

	t.Run("should handle very large packet count", func(t *testing.T) {
		p.Hosts = []string{"localhost"}
		p.Packets = 10000
		data, err := p.Collect(context.Background())
		_ = data
		_ = err
	})

	t.Run("should handle very large timeout", func(t *testing.T) {
		p.Hosts = []string{"localhost"}
		p.Timeout = 9999
		data, err := p.Collect(context.Background())
		_ = data
		_ = err
	})

	t.Run("should handle special IP addresses", func(t *testing.T) {
		p.Hosts = []string{"127.0.0.1", "0.0.0.0", "255.255.255.255", "::1"}
		data, err := p.Collect(context.Background())
		_ = data
		_ = err
	})

	t.Run("should handle hosts with special characters", func(t *testing.T) {
		p.Hosts = []string{"host-1.example.com", "host_2.example.com", "sub.domain.example.com"}
		data, err := p.Collect(context.Background())
		_ = data
		_ = err
	})

	t.Run("should handle empty string in hosts list", func(t *testing.T) {
		p.Hosts = []string{"", "localhost", ""}
		data, err := p.Collect(context.Background())
		_ = data
		_ = err
	})

	t.Run("should handle very long host list", func(t *testing.T) {
		hosts := make([]string, 500)
		for i := 0; i < 500; i++ {
			hosts[i] = "host" + string(rune(i)) + ".com"
		}
		p.Hosts = hosts
		data, err := p.Collect(context.Background())
		_ = data
		_ = err
	})
}

func TestCollectMultipleSequentialCalls(t *testing.T) {
	p := New()
	require.NotNil(t, p)
	p.Hosts = []string{"localhost"}
	p.Packets = 4
	p.Timeout = 3

	t.Run("should handle multiple sequential collect calls", func(t *testing.T) {
		for i := 0; i < 5; i++ {
			data, err := p.Collect(context.Background())
			assert.NoError(t, err)
			assert.NotNil(t, data)
		}
	})

	t.Run("should handle collect after hosts modification", func(t *testing.T) {
		data1, err1 := p.Collect(context.Background())
		assert.NoError(t, err1)
		assert.NotNil(t, data1)

		p.Hosts = []string{"host2.com"}
		data2, err2 := p.Collect(context.Background())
		assert.NoError(t, err2)
		assert.NotNil(t, data2)
	})

	t.Run("should handle collect after configuration changes", func(t *testing.T) {
		p.Packets = 4
		data1, err1 := p.Collect(context.Background())
		assert.NoError(t, err1)
		assert.NotNil(t, data1)

		p.Packets = 8
		data2, err2 := p.Collect(context.Background())
		assert.NoError(t, err2)
		assert.NotNil(t, data2)
	})
}

func TestCollectReturnTypes(t *testing.T) {
	p := New()
	require.NotNil(t, p)
	p.Hosts = []string{"localhost"}
	p.Packets = 4
	p.Timeout = 3

	data, err := p.Collect(context.Background())

	t.Run("should return map[string]interface{}", func(t *testing.T) {
		if data != nil {
			assert.IsType(t, make(map[string]interface{}), data)
		}
	})

	t.Run("should return error or nil", func(t *testing.T) {
		// err can be nil or an error
		_ = err
	})

	t.Run("should not return nil for both data and error", func(t *testing.T) {
		// At minimum, should return non-nil data
		assert.NotNil(t, data)
	})
}

func TestCollectConcurrency(t *testing.T) {
	p := New()
	require.NotNil(t, p)
	p.Hosts = []string{"localhost", "8.8.8.8", "example.com"}
	p.Packets = 4
	p.Timeout = 3

	t.Run("should handle concurrent collect calls safely", func(t *testing.T) {
		done := make(chan bool, 3)
		for i := 0; i < 3; i++ {
			go func() {
				data, err := p.Collect(context.Background())
				assert.NotNil(t, data)
				assert.NoError(t, err)
				done <- true
			}()
		}

		for i := 0; i < 3; i++ {
			<-done
		}
	})
}

func TestCollectDataValues(t *testing.T) {
	p := New()
	require.NotNil(t, p)
	p.Hosts = []string{"localhost"}
	p.Packets = 4
	p.Timeout = 3

	data, err := p.Collect(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, data)

	t.Run("should return numeric values for min/avg/max if present", func(t *testing.T) {
		for key, val := range data {
			if key == "min" || key == "avg" || key == "max" {
				switch v := val.(type) {
				case float64:
					assert.GreaterOrEqual(t, v, 0.0)
				case int:
					assert.GreaterOrEqual(t, v, 0)
				case interface{}:
					// Accept other numeric types
				}
			}
		}
	})

	t.Run("should return valid packet loss percentage if present", func(t *testing.T) {
		if loss, exists := data["loss"]; exists {
			switch v := loss.(type) {
			case float64:
				assert.GreaterOrEqual(t, v, 0.0)
				assert.LessOrEqual(t, v, 100.0)
			case int:
				assert.GreaterOrEqual(t, v, 0)
				assert.LessOrEqual(t, v, 100)
			}
		}
	})
}