package indexer

import (
	"context"

	"github.com/0glabs/0g-storage-client/node"
	"github.com/pkg/errors"
)

// Requires `indexerApi` implements the `Interface` interface.
var _ Interface = (*IndexerApi)(nil)

type IndexerApi struct {
	Namespace string
	nodes     []*node.Client
}

func NewIndexerApi(nodes []*node.Client) *IndexerApi {
	return &IndexerApi{"indexer", nodes}
}

func (api *IndexerApi) GetNodes(ctx context.Context) ([]ShardedNode, error) {
	var result []ShardedNode

	for _, v := range api.nodes {
		config, err := v.ZeroGStorage().GetShardConfig(ctx)
		if err != nil {
			return nil, errors.WithMessage(err, "Failed to query shard config from storage node")
		}
		if config.IsValid() {
			result = append(result, ShardedNode{
				URL:    v.URL(),
				Config: config,
			})
		}

	}

	return result, nil
}
