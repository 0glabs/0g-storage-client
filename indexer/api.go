package indexer

import (
	"context"

	"github.com/pkg/errors"
)

// Requires `indexerApi` implements the `Interface` interface.
var _ Interface = (*IndexerApi)(nil)

// IndexerApi indexer service configuration
type IndexerApi struct {
	Namespace string
	manager   *NodeManager
}

// NewIndexerApi creates indexer service configuration
func NewIndexerApi(manager *NodeManager) *IndexerApi {
	return &IndexerApi{"indexer", manager}
}

// GetShardedNodes return storage node list
func (api *IndexerApi) GetShardedNodes(ctx context.Context) (ShardedNodes, error) {
	trusted, err := api.manager.Trusted()
	if err != nil {
		return ShardedNodes{}, errors.WithMessage(err, "Failed to retrieve trusted nodes")
	}

	return ShardedNodes{
		Trusted:    trusted,
		Discovered: api.manager.Discovered(),
	}, nil
}

// GetNodeLocations return IP locations of all nodes.
func (api *IndexerApi) GetNodeLocations(ctx context.Context) (map[string]*IPLocation, error) {
	return CachedLocations(), nil
}
