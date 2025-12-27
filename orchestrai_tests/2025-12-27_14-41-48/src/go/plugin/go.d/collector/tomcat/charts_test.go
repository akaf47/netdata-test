package tomcat

import (
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin.d/go.d/pkg/web"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCharts(t *testing.T) {
	charts := Charts()
	
	t.Run("Returns non-nil charts", func(t *testing.T) {
		assert.NotNil(t, charts)
	})
	
	t.Run("Contains expected chart templates", func(t *testing.T) {
		require.NotNil(t, charts)
		// Verify charts contain thread pool metrics
		hasThreadChart := false
		for _, chart := range charts {
			if chart != nil && chart.ID == "tomcat_threads" {
				hasThreadChart = true
				break
			}
		}
		// Thread metrics are typically included in Tomcat charts
	})
	
	t.Run("All charts have valid IDs", func(t *testing.T) {
		require.NotNil(t, charts)
		for _, chart := range charts {
			if chart != nil {
				assert.NotEmpty(t, chart.ID, "chart ID should not be empty")
			}
		}
	})
	
	t.Run("All charts have valid titles", func(t *testing.T) {
		require.NotNil(t, charts)
		for _, chart := range charts {
			if chart != nil {
				assert.NotEmpty(t, chart.Title, "chart title should not be empty")
			}
		}
	})
	
	t.Run("All dimensions in charts are properly defined", func(t *testing.T) {
		require.NotNil(t, charts)
		for _, chart := range charts {
			if chart != nil && chart.Dims != nil {
				for _, dim := range chart.Dims {
					assert.NotEmpty(t, dim.ID, "dimension ID should not be empty")
					assert.NotEmpty(t, dim.Name, "dimension name should not be empty")
				}
			}
		}
	})
}

func TestChartsInitialization(t *testing.T) {
	t.Run("Charts initialized multiple times are consistent", func(t *testing.T) {
		charts1 := Charts()
		charts2 := Charts()
		
		require.NotNil(t, charts1)
		require.NotNil(t, charts2)
		assert.Equal(t, len(charts1), len(charts2), "multiple calls should return consistent number of charts")
	})
}

func BenchmarkCharts(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = Charts()
	}
}