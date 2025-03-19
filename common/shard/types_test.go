package shard

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
)

func makeShardNode(numShard uint, shardId uint) *ShardedNode {
	return &ShardedNode{Config: ShardConfig{
		NumShard: uint64(numShard),
		ShardId:  uint64(shardId),
	}}
}

func TestSelect(t *testing.T) {
	shardedNodes := []*ShardedNode{
		makeShardNode(4, 0),
		makeShardNode(4, 2),
		makeShardNode(4, 3),
		makeShardNode(1, 0),
		makeShardNode(2, 0),
		makeShardNode(8, 1),
		makeShardNode(8, 5),
		makeShardNode(16, 0),
		makeShardNode(16, 1),
		makeShardNode(16, 2),
		makeShardNode(16, 3),
		makeShardNode(16, 4),
		makeShardNode(16, 5),
		makeShardNode(16, 6),
		makeShardNode(16, 7),
		makeShardNode(16, 8),
		makeShardNode(16, 9),
		makeShardNode(16, 10),
		makeShardNode(16, 11),
		makeShardNode(16, 12),
		makeShardNode(16, 13),
		makeShardNode(16, 14),
		makeShardNode(16, 15),
	}
	selected, found := Select(shardedNodes, 2, "min")
	assert.Equal(t, found, true)
	fmt.Println(selected)
	assert.Equal(t, len(selected), 5)
	assert.DeepEqual(t, selected[0], makeShardNode(1, 0))
	assert.DeepEqual(t, selected[1], makeShardNode(2, 0))
	assert.DeepEqual(t, selected[2], makeShardNode(4, 3))
	assert.DeepEqual(t, selected[3], makeShardNode(8, 1))
	assert.DeepEqual(t, selected[4], makeShardNode(8, 5))
	selected, found = Select(shardedNodes, 3, "min")
	assert.Equal(t, found, true)
	assert.Equal(t, len(selected), 15)
	assert.DeepEqual(t, selected[0], makeShardNode(1, 0))
	assert.DeepEqual(t, selected[1], makeShardNode(2, 0))
	assert.DeepEqual(t, selected[2], makeShardNode(4, 0))
	assert.DeepEqual(t, selected[3], makeShardNode(4, 2))
	assert.DeepEqual(t, selected[4], makeShardNode(4, 3))
	assert.DeepEqual(t, selected[5], makeShardNode(8, 1))
	assert.DeepEqual(t, selected[6], makeShardNode(8, 5))
	assert.DeepEqual(t, selected[7], makeShardNode(16, 1))
	assert.DeepEqual(t, selected[8], makeShardNode(16, 3))
	assert.DeepEqual(t, selected[9], makeShardNode(16, 5))
	assert.DeepEqual(t, selected[10], makeShardNode(16, 7))
	assert.DeepEqual(t, selected[11], makeShardNode(16, 9))
	assert.DeepEqual(t, selected[12], makeShardNode(16, 11))
	assert.DeepEqual(t, selected[13], makeShardNode(16, 13))
	assert.DeepEqual(t, selected[14], makeShardNode(16, 15))
	_, found = Select(shardedNodes, 4, "max")
	assert.Equal(t, found, false)

	selected, found = Select(shardedNodes, 1, "16")
	assert.Equal(t, found, true)
	assert.Equal(t, len(selected), 16)
	for i := 0; i < 16; i++ {
		assert.DeepEqual(t, selected[i], makeShardNode(16, uint(i)))
	}

	selected, found = Select(shardedNodes, 2, "max")
	assert.Equal(t, found, true)
	assert.Equal(t, len(selected), 21)
	for i := 0; i < 16; i++ {
		assert.DeepEqual(t, selected[i], makeShardNode(16, uint(i)))
	}
	assert.DeepEqual(t, selected[16], makeShardNode(8, 1))
	assert.DeepEqual(t, selected[17], makeShardNode(8, 5))
	assert.DeepEqual(t, selected[18], makeShardNode(4, 0))
	assert.DeepEqual(t, selected[19], makeShardNode(4, 2))
	assert.DeepEqual(t, selected[20], makeShardNode(4, 3))
}

func TestNextSegmentIndex(t *testing.T) {
	tests := []struct {
		name              string
		shardConfig       ShardConfig
		startSegmentIndex uint64
		expectedResult    uint64
	}{
		{
			name:              "Single shard case (NumShard = 1)",
			shardConfig:       ShardConfig{NumShard: 1, ShardId: 0},
			startSegmentIndex: 5,
			expectedResult:    5, // Since it's the only shard, startSegmentIndex is returned as-is
		},
		{
			name:              "Multiple shards, ShardId 0",
			shardConfig:       ShardConfig{NumShard: 4, ShardId: 0},
			startSegmentIndex: 5,
			expectedResult:    8,
		},
		{
			name:              "Multiple shards, ShardId 1",
			shardConfig:       ShardConfig{NumShard: 4, ShardId: 1},
			startSegmentIndex: 5,
			expectedResult:    5,
		},
		{
			name:              "Multiple shards, ShardId 3",
			shardConfig:       ShardConfig{NumShard: 4, ShardId: 3},
			startSegmentIndex: 8,
			expectedResult:    11,
		},
		{
			name:              "Multiple shards, ShardId 2 with startSegmentIndex already covered",
			shardConfig:       ShardConfig{NumShard: 4, ShardId: 2},
			startSegmentIndex: 6,
			expectedResult:    6,
		},
		{
			name:              "Multiple shards, ShardId 0, already aligned",
			shardConfig:       ShardConfig{NumShard: 4, ShardId: 0},
			startSegmentIndex: 8,
			expectedResult:    8,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.shardConfig.NextSegmentIndex(test.startSegmentIndex)
			assert.Equal(t, test.expectedResult, result)
		})
	}
}
