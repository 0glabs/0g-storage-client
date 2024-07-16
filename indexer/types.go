package indexer

import (
	"context"

	"github.com/0glabs/0g-storage-client/common/shard"
)

type ShardedNodes struct {
	Trusted    []*shard.ShardedNode
	Discovered []*shard.ShardedNode
}

type Interface interface {
	GetNodes(ctx context.Context) (ShardedNodes, error)
}
