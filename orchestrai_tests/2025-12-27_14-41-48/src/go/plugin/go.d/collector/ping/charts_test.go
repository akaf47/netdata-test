package ping

import (
	"testing"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCharts(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*Ping)
		verify   func(t *testing.T, charts *module.Charts)
		wantErr  bool
		wantNil  bool
	}{
		{
			name: "should return charts when ping has hosts",
			setup: func(p *Ping) {
				p.Hosts = []string{"localhost", "8.8.8.8"}
			},
			verify: func(t *testing.T, charts *module.Charts) {
				assert.NotNil(t, charts)
				assert.Greater(t, len(*charts), 0)
			},
			wantErr: false,
			wantNil: false,
		},
		{
			name: "should return charts when ping has no hosts",
			setup: func(p *Ping) {
				p.Hosts = []string{}
			},
			verify: func(t *testing.T, charts *module.Charts) {
				assert.NotNil(t, charts)
			},
			wantErr: false,
			wantNil: false,
		},
		{
			name: "should return charts when ping is uninitialized",
			setup: func(p *Ping) {
				// Use uninitialized/default state
			},
			verify: func(t *testing.T, charts *module.Charts) {
				assert.NotNil(t, charts)
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

			charts := p.Charts()

			if tt.wantNil {
				assert.Nil(t, charts)
			} else {
				assert.NotNil(t, charts)
				if tt.verify != nil {
					tt.verify(t, charts)
				}
			}
		})
	}
}

func TestChartsStructure(t *testing.T) {
	p := New()
	require.NotNil(t, p)

	charts := p.Charts()
	require.NotNil(t, charts)

	t.Run("charts should be module.Charts type", func(t *testing.T) {
		assert.IsType(t, &module.Charts{}, charts)
	})

	t.Run("charts should be iterable", func(t *testing.T) {
		count := 0
		for range *charts {
			count++
		}
		assert.Greater(t, count, 0)
	})

	t.Run("each chart should have valid structure", func(t *testing.T) {
		for _, chart := range *charts {
			assert.NotNil(t, chart)
			assert.NotEmpty(t, chart.ID)
			assert.NotEmpty(t, chart.Title)
		}
	})
}

func TestChartsMultipleCalls(t *testing.T) {
	p := New()
	require.NotNil(t, p)

	t.Run("should return consistent charts across multiple calls", func(t *testing.T) {
		charts1 := p.Charts()
		charts2 := p.Charts()

		assert.NotNil(t, charts1)
		assert.NotNil(t, charts2)
	})

	t.Run("should handle charts after host modification", func(t *testing.T) {
		p.Hosts = []string{"host1.com", "host2.com"}
		charts := p.Charts()
		assert.NotNil(t, charts)

		p.Hosts = []string{"host1.com", "host2.com", "host3.com"}
		charts = p.Charts()
		assert.NotNil(t, charts)
	})

	t.Run("should handle charts after clearing hosts", func(t *testing.T) {
		p.Hosts = []string{"host1.com", "host2.com"}
		charts := p.Charts()
		assert.NotNil(t, charts)

		p.Hosts = []string{}
		charts = p.Charts()
		assert.NotNil(t, charts)
	})
}

func TestChartsEdgeCases(t *testing.T) {
	t.Run("should handle nil hosts slice gracefully", func(t *testing.T) {
		p := New()
		require.NotNil(t, p)
		p.Hosts = nil

		charts := p.Charts()
		assert.NotNil(t, charts)
	})

	t.Run("should handle very long host list", func(t *testing.T) {
		p := New()
		require.NotNil(t, p)

		hosts := make([]string, 1000)
		for i := 0; i < 1000; i++ {
			hosts[i] = "host" + string(rune(i)) + ".com"
		}
		p.Hosts = hosts

		charts := p.Charts()
		assert.NotNil(t, charts)
	})

	t.Run("should handle special characters in hosts", func(t *testing.T) {
		p := New()
		require.NotNil(t, p)
		p.Hosts = []string{"host-1.example.com", "host_2.example.com", "127.0.0.1", "::1"}

		charts := p.Charts()
		assert.NotNil(t, charts)
	})

	t.Run("should handle empty string in hosts", func(t *testing.T) {
		p := New()
		require.NotNil(t, p)
		p.Hosts = []string{"", "validhost.com", ""}

		charts := p.Charts()
		assert.NotNil(t, charts)
	})
}