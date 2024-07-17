package indexer

import (
	"context"

	"github.com/0glabs/0g-storage-client/common/shard"
)

type ShardedNodes struct {
	Trusted    []*shard.ShardedNode `json:"trusted"`
	Discovered []*shard.ShardedNode `json:"discovered"`
}

type NodeInfo struct {
	*shard.ShardedNode
	Location *IPLocation `json:"location"`
}

type Interface interface {
	GetShardedNodes(ctx context.Context) (ShardedNodes, error)

	GetNodes(ctx context.Context) ([]*NodeInfo, error)
}
