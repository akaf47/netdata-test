package varnish

import (
	"errors"
	"testing"

	"github.com/netdata/netdata/go/plugins/plugin/go.d/agent/module"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCollect tests the Collect method with successful data retrieval
func TestCollect_Success(t *testing.T) {
	v := &Varnish{
		client: &mockVarnishClient{
			stats: map[string]interface{}{
				"MAIN.uptime":           float64(1000),
				"MAIN.sess_conn":        float64(100),
				"MAIN.sess_drop":        float64(10),
				"MAIN.sess_fail":        float64(5),
				"MAIN.client_req":       float64(500),
				"MAIN.cache_hit":        float64(400),
				"MAIN.cache_miss":       float64(100),
				"MAIN.backend_conn":     float64(50),
				"MAIN.backend_unhealthy": float64(2),
				"MAIN.fetch_head":       float64(20),
				"MAIN.fetch_length":     float64(80),
				"MAIN.fetch_bad":        float64(1),
				"MAIN.fetch_failed":     float64(0),
				"MAIN.fetch_no_thread":  float64(1),
				"MAIN.pools":            float64(2),
				"MAIN.threads":          float64(4),
				"MAIN.threads_created":  float64(4),
				"MAIN.threads_failed":   float64(0),
				"MAIN.threads_limited":  float64(0),
				"MAIN.thread_queue_len": float64(0),
				"MAIN.n_object":         float64(100),
				"MAIN.n_vampireobject":  float64(0),
				"MAIN.n_objectcore":     float64(100),
				"MAIN.n_objecthead":     float64(100),
				"MAIN.n_waitinglist":    float64(5),
				"MAIN.n_backend":        float64(2),
				"MAIN.n_expired":        float64(10),
				"MAIN.n_lru_nuked":      float64(0),
				"MAIN.n_lru_saved":      float64(5),
			},
		},
	}

	mx := v.Collect()

	assert.NotNil(t, mx)
	assert.Greater(t, len(mx), 0)
}

// TestCollect_ClientError tests Collect when client returns error
func TestCollect_ClientError(t *testing.T) {
	v := &Varnish{
		client: &mockVarnishClient{
			err: errors.New("connection failed"),
		},
	}

	mx := v.Collect()

	assert.Nil(t, mx)
}

// TestCollect_EmptyStats tests Collect with empty stats
func TestCollect_EmptyStats(t *testing.T) {
	v := &Varnish{
		client: &mockVarnishClient{
			stats: map[string]interface{}{},
		},
	}

	mx := v.Collect()

	assert.Nil(t, mx)
}

// TestCollect_NilClient tests Collect with nil client
func TestCollect_NilClient(t *testing.T) {
	v := &Varnish{
		client: nil,
	}

	mx := v.Collect()

	assert.Nil(t, mx)
}

// TestCollect_PartialStats tests Collect with partial stats
func TestCollect_PartialStats(t *testing.T) {
	v := &Varnish{
		client: &mockVarnishClient{
			stats: map[string]interface{}{
				"MAIN.uptime":     float64(1000),
				"MAIN.sess_conn":  float64(100),
				"MAIN.cache_hit":  float64(400),
			},
		},
	}

	mx := v.Collect()

	assert.NotNil(t, mx)
}

// TestCollect_StatTypeConversion tests stats with unexpected types
func TestCollect_StatTypeConversion(t *testing.T) {
	v := &Varnish{
		client: &mockVarnishClient{
			stats: map[string]interface{}{
				"MAIN.uptime":   "invalid",
				"MAIN.sess_conn": int64(100),
			},
		},
	}

	mx := v.Collect()

	// Should handle gracefully
	assert.NotNil(t, mx)
}

// TestCollect_ZeroValues tests Collect with all zero values
func TestCollect_ZeroValues(t *testing.T) {
	v := &Varnish{
		client: &mockVarnishClient{
			stats: map[string]interface{}{
				"MAIN.uptime":           float64(0),
				"MAIN.sess_conn":        float64(0),
				"MAIN.sess_drop":        float64(0),
				"MAIN.sess_fail":        float64(0),
				"MAIN.client_req":       float64(0),
				"MAIN.cache_hit":        float64(0),
				"MAIN.cache_miss":       float64(0),
				"MAIN.backend_conn":     float64(0),
				"MAIN.backend_unhealthy": float64(0),
			},
		},
	}

	mx := v.Collect()

	assert.NotNil(t, mx)
}

// TestCollect_LargeValues tests Collect with large stat values
func TestCollect_LargeValues(t *testing.T) {
	v := &Varnish{
		client: &mockVarnishClient{
			stats: map[string]interface{}{
				"MAIN.uptime":     float64(999999999),
				"MAIN.sess_conn":  float64(1000000),
				"MAIN.cache_hit":  float64(10000000),
				"MAIN.client_req": float64(50000000),
			},
		},
	}

	mx := v.Collect()

	assert.NotNil(t, mx)
}

// TestCollect_NegativeValues tests Collect with negative values
func TestCollect_NegativeValues(t *testing.T) {
	v := &Varnish{
		client: &mockVarnishClient{
			stats: map[string]interface{}{
				"MAIN.uptime":    float64(-1),
				"MAIN.sess_conn": float64(-100),
				"MAIN.cache_hit": float64(-50),
			},
		},
	}

	mx := v.Collect()

	// Should handle gracefully
	assert.NotNil(t, mx)
}

// TestCollect_FloatPrecision tests Collect with float values
func TestCollect_FloatPrecision(t *testing.T) {
	v := &Varnish{
		client: &mockVarnishClient{
			stats: map[string]interface{}{
				"MAIN.uptime":     float64(1234.5678),
				"MAIN.sess_conn":  float64(100.1),
				"MAIN.cache_hit":  float64(400.9999),
			},
		},
	}

	mx := v.Collect()

	assert.NotNil(t, mx)
}

// Mock client for testing
type mockVarnishClient struct {
	stats map[string]interface{}
	err   error
}

func (m *mockVarnishClient) Stats() (map[string]interface{}, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.stats, nil
}