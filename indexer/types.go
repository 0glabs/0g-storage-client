package indexer

import (
	"sort"

	"github.com/0glabs/0g-storage-client/node"
)

type ShardedNode struct {
	URL    string
	Config node.ShardConfig
}

type Interface interface {
	GetNodes() ([]ShardedNode, error)
}

type shardSegmentTreeNode struct {
	childs   []*shardSegmentTreeNode
	numShard uint
	lazyTags uint
	replica  uint
}

func (node *shardSegmentTreeNode) pushdown() {
	if node.childs == nil {
		node.childs = make([]*shardSegmentTreeNode, 2)
		for i := 0; i < 2; i += 1 {
			node.childs[i] = &shardSegmentTreeNode{
				numShard: node.numShard << 1,
				replica:  0,
				lazyTags: 0,
			}
		}
	}
	for i := 0; i < 2; i += 1 {
		node.childs[i].replica += node.lazyTags
		node.childs[i].lazyTags += node.lazyTags
	}
	node.lazyTags = 0
}

// insert a shard if it contributes to the replica
func (node *shardSegmentTreeNode) insert(numShard uint, shardId uint, expectedReplica uint) bool {
	if node.replica >= expectedReplica {
		return false
	}
	if node.numShard == numShard {
		node.replica += 1
		node.lazyTags += 1
		return true
	}
	node.pushdown()
	inserted := node.childs[shardId%2].insert(numShard, shardId>>1, expectedReplica)
	node.replica = min(node.childs[0].replica, node.childs[1].replica)
	return inserted
}

// select a set of given sharded node and make the data is replicated at least expctedReplica times
// return the selected nodes and if selection is successful
func Select(nodes []ShardedNode, expectedReplica uint) ([]ShardedNode, bool) {
	selected := make([]ShardedNode, 0)
	if expectedReplica == 0 {
		return selected, true
	}
	// sort by shard size from large to small
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].Config.NumShard == nodes[j].Config.NumShard {
			return nodes[i].Config.ShardId < nodes[j].Config.ShardId
		}
		return nodes[i].Config.NumShard < nodes[j].Config.NumShard
	})
	// build segment tree to select proper nodes
	root := shardSegmentTreeNode{
		numShard: 1,
		replica:  0,
		lazyTags: 0,
	}

	for _, node := range nodes {
		if root.insert(uint(node.Config.NumShard), uint(node.Config.ShardId), expectedReplica) {
			selected = append(selected, node)
		}
		if root.replica >= expectedReplica {
			return selected, true
		}
	}
	return make([]ShardedNode, 0), false
}
