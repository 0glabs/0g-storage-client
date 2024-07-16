package indexer

import (
	"context"

	"github.com/0glabs/0g-storage-client/common/shard"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/pkg/errors"
)

// Requires `indexerApi` implements the `Interface` interface.
var _ Interface = (*IndexerApi)(nil)

// IndexerApi indexer service configuration
type IndexerApi struct {
	Namespace string
	nodes     []*node.Client
}

// NewIndexerApi creates indexer service configuration
func NewIndexerApi(nodes []*node.Client) *IndexerApi {
	return &IndexerApi{"indexer", nodes}
}

// GetNodes return storage node list
func (api *IndexerApi) GetNodes(ctx context.Context) ([]shard.ShardedNode, error) {
	var result []shard.ShardedNode

	for _, v := range api.nodes {
		config, err := v.ZeroGStorage().GetShardConfig(ctx)
		if err != nil {
			return nil, errors.WithMessage(err, "Failed to query shard config from storage node")
		}
		if config.IsValid() {
			result = append(result, shard.ShardedNode{
				URL:    v.URL(),
				Config: config,
			})
		}

	}

	return result, nil
}
