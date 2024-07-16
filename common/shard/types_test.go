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
	selected, found := Select(shardedNodes, 2)
	assert.Equal(t, found, true)
	fmt.Println(selected)
	assert.Equal(t, len(selected), 5)
	assert.DeepEqual(t, selected[0], makeShardNode(1, 0))
	assert.DeepEqual(t, selected[1], makeShardNode(2, 0))
	assert.DeepEqual(t, selected[2], makeShardNode(4, 3))
	assert.DeepEqual(t, selected[3], makeShardNode(8, 1))
	assert.DeepEqual(t, selected[4], makeShardNode(8, 5))
	selected, found = Select(shardedNodes, 3)
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
	_, found = Select(shardedNodes, 4)
	assert.Equal(t, found, false)
}
