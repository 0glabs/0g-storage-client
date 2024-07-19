package indexer

import (
	"context"

	"github.com/0glabs/0g-storage-client/common/shard"
	"github.com/pkg/errors"
)

// Requires `indexerApi` implements the `Interface` interface.
var _ Interface = (*IndexerApi)(nil)

// IndexerApi indexer service configuration
type IndexerApi struct {
	Namespace string
}

// NewIndexerApi creates indexer service configuration
func NewIndexerApi() *IndexerApi {
	return &IndexerApi{"indexer"}
}

// GetShardedNodes return storage node list
func (api *IndexerApi) GetShardedNodes(ctx context.Context) (ShardedNodes, error) {
	trusted, err := defaultNodeManager.Trusted()
	if err != nil {
		return ShardedNodes{}, errors.WithMessage(err, "Failed to retrieve trusted nodes")
	}

	return ShardedNodes{
		Trusted:    trusted,
		Discovered: defaultNodeManager.Discovered(),
	}, nil
}

// GetNodeLocations return IP locations of all nodes.
func (api *IndexerApi) GetNodeLocations(ctx context.Context) (map[string]*IPLocation, error) {
	return defaultIPLocationManager.All(), nil
}

// GetFileLocations return locations info of given file.
func (api *IndexerApi) GetFileLocations(ctx context.Context, txSeq uint64) (locations []*shard.ShardedNode, err error) {
	return defaultFileLocationCache.GetFileLocations(ctx, txSeq)
}
