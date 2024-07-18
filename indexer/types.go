package indexer

import (
	"context"

	"github.com/0glabs/0g-storage-client/common/shard"
)

type ShardedNodes struct {
	Trusted    []*shard.ShardedNode `json:"trusted"`
	Discovered []*shard.ShardedNode `json:"discovered"`
}

type Interface interface {
	GetShardedNodes(ctx context.Context) (ShardedNodes, error)

	GetNodeLocations(ctx context.Context) (map[string]*IPLocation, error)
}
