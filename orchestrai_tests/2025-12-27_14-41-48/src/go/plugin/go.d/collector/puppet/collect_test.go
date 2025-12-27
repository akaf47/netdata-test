package puppet

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectSummaryStats(t *testing.T) {
	type args struct {
		summaryStats map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should collect valid summary stats",
			args: args{
				summaryStats: map[string]interface{}{
					"total": float64(100),
					"time":  map[string]interface{}{"total": 5.5},
				},
			},
			wantErr: false,
		},
		{
			name: "should handle nil summary stats",
			args: args{
				summaryStats: nil,
			},
			wantErr: true,
		},
		{
			name: "should handle empty summary stats",
			args: args{
				summaryStats: map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name: "should handle missing time field",
			args: args{
				summaryStats: map[string]interface{}{
					"total": float64(100),
				},
			},
			wantErr: false,
		},
		{
			name: "should handle non-numeric total",
			args: args{
				summaryStats: map[string]interface{}{
					"total": "not a number",
				},
			},
			wantErr: false,
		},
		{
			name: "should handle negative total",
			args: args{
				summaryStats: map[string]interface{}{
					"total": float64(-1),
				},
			},
			wantErr: false,
		},
		{
			name: "should handle zero total",
			args: args{
				summaryStats: map[string]interface{}{
					"total": float64(0),
				},
			},
			wantErr: false,
		},
		{
			name: "should handle very large numbers",
			args: args{
				summaryStats: map[string]interface{}{
					"total": float64(999999999999),
					"time": map[string]interface{}{
						"total": 999999.999,
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Puppet{}
			err := p.collectSummaryStats(tt.args.summaryStats)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCollectResourceMetrics(t *testing.T) {
	type args struct {
		resources map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should collect valid resources",
			args: args{
				resources: map[string]interface{}{
					"changed": float64(5),
					"failed":  float64(0),
					"total":   float64(100),
				},
			},
			wantErr: false,
		},
		{
			name: "should handle nil resources",
			args: args{
				resources: nil,
			},
			wantErr: true,
		},
		{
			name: "should handle empty resources",
			args: args{
				resources: map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name: "should handle zero values",
			args: args{
				resources: map[string]interface{}{
					"changed": float64(0),
					"failed":  float64(0),
					"total":   float64(0),
				},
			},
			wantErr: false,
		},
		{
			name: "should handle negative values",
			args: args{
				resources: map[string]interface{}{
					"changed": float64(-1),
					"failed":  float64(-1),
					"total":   float64(-1),
				},
			},
			wantErr: false,
		},
		{
			name: "should handle non-numeric values",
			args: args{
				resources: map[string]interface{}{
					"changed": "not a number",
					"failed":  "invalid",
					"total":   "bad",
				},
			},
			wantErr: false,
		},
		{
			name: "should handle partial metrics",
			args: args{
				resources: map[string]interface{}{
					"changed": float64(5),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Puppet{}
			err := p.collectResourceMetrics(tt.args.resources)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCollectTimeMetrics(t *testing.T) {
	type args struct {
		times map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should collect valid time metrics",
			args: args{
				times: map[string]interface{}{
					"total":      5.5,
					"file":       1.2,
					"exec":       2.1,
					"filebucket": 0.5,
				},
			},
			wantErr: false,
		},
		{
			name: "should handle nil times",
			args: args{
				times: nil,
			},
			wantErr: true,
		},
		{
			name: "should handle empty times",
			args: args{
				times: map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name: "should handle zero time values",
			args: args{
				times: map[string]interface{}{
					"total": 0.0,
				},
			},
			wantErr: false,
		},
		{
			name: "should handle negative time values",
			args: args{
				times: map[string]interface{}{
					"total": -1.5,
				},
			},
			wantErr: false,
		},
		{
			name: "should handle non-numeric time values",
			args: args{
				times: map[string]interface{}{
					"total": "not a number",
				},
			},
			wantErr: false,
		},
		{
			name: "should handle very small time values",
			args: args{
				times: map[string]interface{}{
					"total": 0.0001,
				},
			},
			wantErr: false,
		},
		{
			name: "should handle very large time values",
			args: args{
				times: map[string]interface{}{
					"total": 99999.9999,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Puppet{}
			err := p.collectTimeMetrics(tt.args.times)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCollectChanges(t *testing.T) {
	type args struct {
		changes map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should collect valid changes",
			args: args{
				changes: map[string]interface{}{
					"total": float64(42),
				},
			},
			wantErr: false,
		},
		{
			name: "should handle nil changes",
			args: args{
				changes: nil,
			},
			wantErr: true,
		},
		{
			name: "should handle empty changes",
			args: args{
				changes: map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name: "should handle zero changes",
			args: args{
				changes: map[string]interface{}{
					"total": float64(0),
				},
			},
			wantErr: false,
		},
		{
			name: "should handle negative changes",
			args: args{
				changes: map[string]interface{}{
					"total": float64(-1),
				},
			},
			wantErr: false,
		},
		{
			name: "should handle non-numeric changes",
			args: args{
				changes: map[string]interface{}{
					"total": "invalid",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Puppet{}
			err := p.collectChanges(tt.args.changes)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCollectEvents(t *testing.T) {
	type args struct {
		events map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "should collect valid events",
			args: args{
				events: map[string]interface{}{
					"total":     float64(100),
					"success":   float64(95),
					"failure":   float64(5),
					"audit":     float64(0),
					"noop":      float64(0),
					"skipped":   float64(0),
					"timestamp": float64(1234567890),
				},
			},
			wantErr: false,
		},
		{
			name: "should handle nil events",
			args: args{
				events: nil,
			},
			wantErr: true,
		},
		{
			name: "should handle empty events",
			args: args{
				events: map[string]interface{}{},
			},
			wantErr: false,
		},
		{
			name: "should handle zero event counts",
			args: args{
				events: map[string]interface{}{
					"total":   float64(0),
					"success": float64(0),
				},
			},
			wantErr: false,
		},
		{
			name: "should handle negative event counts",
			args: args{
				events: map[string]interface{}{
					"total": float64(-1),
				},
			},
			wantErr: false,
		},
		{
			name: "should handle non-numeric event values",
			args: args{
				events: map[string]interface{}{
					"total": "not a number",
				},
			},
			wantErr: false,
		},
		{
			name: "should handle partial event data",
			args: args{
				events: map[string]interface{}{
					"total": float64(50),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Puppet{}
			err := p.collectEvents(tt.args.events)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCollectAllData(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		wantErr bool
	}{
		{
			name: "should collect complete valid data",
			data: map[string]interface{}{
				"version": "7.0.0",
				"transaction_uuid": "uuid-123",
				"time": map[string]interface{}{
					"total": 5.5,
				},
				"resources": map[string]interface{}{
					"total": float64(100),
				},
				"changes": map[string]interface{}{
					"total": float64(5),
				},
				"events": map[string]interface{}{
					"total": float64(100),
				},
			},
			wantErr: false,
		},
		{
			name: "should handle nil data",
			data: nil,
			wantErr: true,
		},
		{
			name: "should handle empty data",
			data: map[string]interface{}{},
			wantErr: false,
		},
		{
			name: "should handle missing optional fields",
			data: map[string]interface{}{
				"version": "7.0.0",
			},
			wantErr: false,
		},
		{
			name: "should handle invalid nested structures",
			data: map[string]interface{}{
				"time": "not a map",
				"resources": "invalid",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Puppet{}
			err := p.collectAllData(tt.data)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestParseFloat(t *testing.T) {
	tests := []struct {
		name    string
		val     interface{}
		want    float64
		wantErr bool
	}{
		{
			name:    "should parse float64",
			val:     float64(123.45),
			want:    123.45,
			wantErr: false,
		},
		{
			name:    "should parse float32",
			val:     float32(123.45),
			want:    123.45,
			wantErr: false,
		},
		{
			name:    "should parse int",
			val:     int(123),
			want:    123.0,
			wantErr: false,
		},
		{
			name:    "should parse int64",
			val:     int64(123),
			want:    123.0,
			wantErr: false,
		},
		{
			name:    "should parse zero",
			val:     float64(0),
			want:    0.0,
			wantErr: false,
		},
		{
			name:    "should parse negative",
			val:     float64(-123.45),
			want:    -123.45,
			wantErr: false,
		},
		{
			name:    "should fail on string",
			val:     "123.45",
			want:    0,
			wantErr: true,
		},
		{
			name:    "should fail on nil",
			val:     nil,
			want:    0,
			wantErr: true,
		},
		{
			name:    "should fail on bool",
			val:     true,
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseFloat(tt.val)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestGetNestedMap(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]interface{}
		key     string
		wantErr bool
	}{
		{
			name: "should get existing nested map",
			data: map[string]interface{}{
				"time": map[string]interface{}{
					"total": 5.5,
				},
			},
			key:     "time",
			wantErr: false,
		},
		{
			name: "should fail on nil data",
			data: nil,
			key:  "time",
			wantErr: true,
		},
		{
			name: "should fail on missing key",
			data: map[string]interface{}{
				"other": "value",
			},
			key:     "time",
			wantErr: true,
		},
		{
			name: "should fail when value is not a map",
			data: map[string]interface{}{
				"time": "not a map",
			},
			key:     "time",
			wantErr: true,
		},
		{
			name: "should fail when value is nil",
			data: map[string]interface{}{
				"time": nil,
			},
			key:     "time",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := getNestedMap(tt.data, tt.key)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetMetricValue(t *testing.T) {
	tests := []struct {
		name     string
		data     map[string]interface{}
		key      string
		expected float64
	}{
		{
			name: "should get existing metric",
			data: map[string]interface{}{
				"total": float64(100),
			},
			key:      "total",
			expected: 100.0,
		},
		{
			name: "should return 0 for missing metric",
			data: map[string]interface{}{
				"other": float64(50),
			},
			key:      "total",
			expected: 0.0,
		},
		{
			name: "should return 0 for nil data",
			data: nil,
			key:  "total",
			expected: 0.0,
		},
		{
			name: "should return 0 for empty data",
			data: map[string]interface{}{},
			key:  "total",
			expected: 0.0,
		},
		{
			name: "should handle non-numeric value",
			data: map[string]interface{}{
				"total": "not a number",
			},
			key:      "total",
			expected: 0.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getMetricValue(tt.data, tt.key)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestUpdateMetrics(t *testing.T) {
	tests := []struct {
		name string
		call func(p *Puppet)
	}{
		{
			name: "should update all metrics without error",
			call: func(p *Puppet) {
				p.updateMetrics()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Puppet{}
			require.NotPanics(t, func() {
				tt.call(p)
			})
		})
	}
}

func TestCollectWithoutError(t *testing.T) {
	tests := []struct {
		name string
		p    *Puppet
	}{
		{
			name: "should collect without error on empty puppet instance",
			p:    &Puppet{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.NotPanics(t, func() {
				tt.p.collectData()
			})
		})
	}
}

func TestCollectEdgeCases(t *testing.T) {
	t.Run("should handle large numbers in all metric types", func(t *testing.T) {
		p := &Puppet{}
		largeNum := float64(999999999999.999)
		
		err := p.collectTimeMetrics(map[string]interface{}{
			"total": largeNum,
		})
		assert.NoError(t, err)
		
		err = p.collectResourceMetrics(map[string]interface{}{
			"total": largeNum,
		})
		assert.NoError(t, err)
		
		err = p.collectEvents(map[string]interface{}{
			"total": largeNum,
		})
		assert.NoError(t, err)
	})

	t.Run("should handle boundary conditions", func(t *testing.T) {
		p := &Puppet{}
		
		// Test with exactly zero
		err := p.collectTimeMetrics(map[string]interface{}{"total": 0.0})
		assert.NoError(t, err)
		
		// Test with minimum positive
		err = p.collectTimeMetrics(map[string]interface{}{"total": 1e-9})
		assert.NoError(t, err)
	})
}