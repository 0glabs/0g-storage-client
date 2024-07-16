package shard

import "sort"

type ShardConfig struct {
	ShardId  uint64 `json:"shardId"`
	NumShard uint64 `json:"numShard"`
}

func (config *ShardConfig) HasSegment(segmentIndex uint64) bool {
	return config.NumShard < 2 || segmentIndex%config.NumShard == config.ShardId
}

func (config *ShardConfig) IsValid() bool {
	// NumShard should be larger than zero and be power of 2
	return config.NumShard > 0 && (config.NumShard&(config.NumShard-1) == 0) && config.ShardId < config.NumShard
}

type ShardedNode struct {
	URL    string      `json:"url"`
	Config ShardConfig `json:"config"`
	/// Latency RPC latency in milli seconds.
	Latency int64 `json:"latency"`
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
func Select(nodes []*ShardedNode, expectedReplica uint) ([]*ShardedNode, bool) {
	selected := make([]*ShardedNode, 0)
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
	return make([]*ShardedNode, 0), false
}

func CheckReplica(shardConfigs []*ShardConfig, expectedReplica uint) bool {
	shardedNodes := make([]*ShardedNode, len(shardConfigs))
	for i, shardConfig := range shardConfigs {
		shardedNodes[i] = &ShardedNode{Config: ShardConfig{
			NumShard: uint64(shardConfig.NumShard),
			ShardId:  uint64(shardConfig.ShardId),
		}}
	}
	_, ok := Select(shardedNodes, expectedReplica)
	return ok
}
