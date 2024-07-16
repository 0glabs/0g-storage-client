package indexer

import (
	"context"

	"github.com/0glabs/0g-storage-client/common/shard"
)

type Interface interface {
	GetNodes(ctx context.Context) ([]shard.ShardedNode, error)
}
