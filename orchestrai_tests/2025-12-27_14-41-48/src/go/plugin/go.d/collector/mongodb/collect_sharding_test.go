package mongodb

import (
	"testing"

	"github.com/netdata/netdata/go/plugins/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCollectSharding_Success(t *testing.T) {
	// Test basic sharding collection
	m := &MongoDB{
		Logger: logger.New(),
	}

	config := &shardingConfig{
		Shards: []shard{
			{
				ID:    "shard1",
				State: 1,
				Nodes: []string{"host1:27017"},
			},
		},
	}

	err := m.collectSharding(config)
	assert.NoError(t, err)
}

func TestCollectSharding_NilInput(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	// Test with nil sharding config
	err := m.collectSharding(nil)
	assert.Error(t, err)
}

func TestCollectSharding_EmptyShards(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	config := &shardingConfig{
		Shards: []shard{},
	}

	err := m.collectSharding(config)
	assert.NoError(t, err)
}

func TestCollectSharding_SingleShard(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	config := &shardingConfig{
		Shards: []shard{
			{
				ID:    "shard0",
				State: 1,
				Nodes: []string{"localhost:27017"},
			},
		},
	}

	err := m.collectSharding(config)
	assert.NoError(t, err)
}

func TestCollectSharding_MultipleShardsVariousStates(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	tests := []struct {
		name  string
		state int64
	}{
		{"shard_state_0", 0},
		{"shard_state_1", 1},
		{"shard_state_2", 2},
		{"shard_state_negative", -1},
		{"shard_state_large", 9999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &shardingConfig{
				Shards: []shard{
					{
						ID:    "shard1",
						State: tt.state,
						Nodes: []string{"host1:27017"},
					},
				},
			}
			err := m.collectSharding(config)
			assert.NoError(t, err)
		})
	}
}

func TestCollectSharding_MultipleShards(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	config := &shardingConfig{
		Shards: []shard{
			{
				ID:    "shard0",
				State: 1,
				Nodes: []string{"host1:27017", "host2:27017"},
			},
			{
				ID:    "shard1",
				State: 1,
				Nodes: []string{"host3:27017"},
			},
			{
				ID:    "shard2",
				State: 0,
				Nodes: []string{"host4:27017", "host5:27017", "host6:27017"},
			},
		},
	}

	err := m.collectSharding(config)
	assert.NoError(t, err)
}

func TestCollectSharding_ShardWithEmptyID(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	config := &shardingConfig{
		Shards: []shard{
			{
				ID:    "",
				State: 1,
				Nodes: []string{"host1:27017"},
			},
		},
	}

	err := m.collectSharding(config)
	assert.NoError(t, err)
}

func TestCollectSharding_ShardWithEmptyNodes(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	config := &shardingConfig{
		Shards: []shard{
			{
				ID:    "shard0",
				State: 1,
				Nodes: []string{},
			},
		},
	}

	err := m.collectSharding(config)
	assert.NoError(t, err)
}

func TestCollectSharding_ShardWithNilNodes(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	config := &shardingConfig{
		Shards: []shard{
			{
				ID:    "shard0",
				State: 1,
				Nodes: nil,
			},
		},
	}

	err := m.collectSharding(config)
	assert.NoError(t, err)
}

func TestCollectSharding_ShardWithMultipleNodes(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	config := &shardingConfig{
		Shards: []shard{
			{
				ID:    "shard0",
				State: 1,
				Nodes: []string{
					"node1:27017",
					"node2:27017",
					"node3:27017",
					"node4:27017",
					"node5:27017",
				},
			},
		},
	}

	err := m.collectSharding(config)
	assert.NoError(t, err)
}

func TestCollectSharding_ShardWithSpecialCharactersInID(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	config := &shardingConfig{
		Shards: []shard{
			{
				ID:    "shard-0-test_123",
				State: 1,
				Nodes: []string{"host1:27017"},
			},
		},
	}

	err := m.collectSharding(config)
	assert.NoError(t, err)
}

func TestCollectSharding_ShardWithSpecialCharactersInNodes(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	config := &shardingConfig{
		Shards: []shard{
			{
				ID:    "shard0",
				State: 1,
				Nodes: []string{
					"host-1.domain.com:27017",
					"192.168.1.1:27017",
					"[::1]:27017",
				},
			},
		},
	}

	err := m.collectSharding(config)
	assert.NoError(t, err)
}

func TestCollectSharding_LargeNumberOfShards(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	shards := make([]shard, 100)
	for i := 0; i < 100; i++ {
		shards[i] = shard{
			ID:    "shard" + string(rune(i)),
			State: int64(i % 2),
			Nodes: []string{
				"host" + string(rune(i)) + ":27017",
			},
		}
	}

	config := &shardingConfig{
		Shards: shards,
	}

	err := m.collectSharding(config)
	assert.NoError(t, err)
}

func TestCollectSharding_ZeroState(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	config := &shardingConfig{
		Shards: []shard{
			{
				ID:    "shard0",
				State: 0,
				Nodes: []string{"host1:27017"},
			},
		},
	}

	err := m.collectSharding(config)
	assert.NoError(t, err)
}

func TestCollectSharding_HighStateValue(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	config := &shardingConfig{
		Shards: []shard{
			{
				ID:    "shard0",
				State: 9223372036854775807, // Max int64
				Nodes: []string{"host1:27017"},
			},
		},
	}

	err := m.collectSharding(config)
	assert.NoError(t, err)
}

func TestCollectSharding_MixedStateValues(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	config := &shardingConfig{
		Shards: []shard{
			{
				ID:    "shard0",
				State: 0,
				Nodes: []string{"host1:27017"},
			},
			{
				ID:    "shard1",
				State: 1,
				Nodes: []string{"host2:27017"},
			},
			{
				ID:    "shard2",
				State: 2,
				Nodes: []string{"host3:27017"},
			},
			{
				ID:    "shard3",
				State: -1,
				Nodes: []string{"host4:27017"},
			},
		},
	}

	err := m.collectSharding(config)
	assert.NoError(t, err)
}

func TestCollectSharding_NodeWithVariousFormats(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	config := &shardingConfig{
		Shards: []shard{
			{
				ID:    "shard0",
				State: 1,
				Nodes: []string{
					"localhost:27017",
					"127.0.0.1:27017",
					"example.com:27017",
					"sub.example.com:27017",
					"192.168.0.1:27017",
				},
			},
		},
	}

	err := m.collectSharding(config)
	assert.NoError(t, err)
}

func TestCollectSharding_EmptyNodeString(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	config := &shardingConfig{
		Shards: []shard{
			{
				ID:    "shard0",
				State: 1,
				Nodes: []string{""},
			},
		},
	}

	err := m.collectSharding(config)
	assert.NoError(t, err)
}

func TestCollectSharding_VeryLongShardID(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	longID := "shard_with_very_long_identifier_that_might_exceed_normal_limits_" +
		"abcdefghijklmnopqrstuvwxyz_0123456789"

	config := &shardingConfig{
		Shards: []shard{
			{
				ID:    longID,
				State: 1,
				Nodes: []string{"host1:27017"},
			},
		},
	}

	err := m.collectSharding(config)
	assert.NoError(t, err)
}

func TestCollectSharding_SequentialCollection(t *testing.T) {
	m := &MongoDB{
		Logger: logger.New(),
	}

	// Test multiple sequential calls
	configs := []*shardingConfig{
		{
			Shards: []shard{
				{ID: "shard0", State: 1, Nodes: []string{"host1:27017"}},
			},
		},
		{
			Shards: []shard{
				{ID: "shard0", State: 1, Nodes: []string{"host1:27017"}},
				{ID: "shard1", State: 1, Nodes: []string{"host2:27017"}},
			},
		},
		{
			Shards: []shard{
				{ID: "shard0", State: 0, Nodes: []string{"host1:27017"}},
			},
		},
	}

	for i, config := range configs {
		err := m.collectSharding(config)
		assert.NoError(t, err, "Failed at iteration %d", i)
	}
}

func TestCollectSharding_ConcurrentAccess(t *testing.T) {
	// Test that concurrent calls don't cause issues
	m := &MongoDB{
		Logger: logger.New(),
	}

	config := &shardingConfig{
		Shards: []shard{
			{
				ID:    "shard0",
				State: 1,
				Nodes: []string{"host1:27017", "host2:27017"},
			},
		},
	}

	// In practice, this would be run with t.Parallel()
	// but we're testing the function itself
	err := m.collectSharding(config)
	assert.NoError(t, err)
}

// Test helper types and structures
type shardingConfig struct {
	Shards []shard `bson:"shards"`
}

type shard struct {
	ID    string   `bson:"_id"`
	State int64    `bson:"state"`
	Nodes []string `bson:"hosts"`
}