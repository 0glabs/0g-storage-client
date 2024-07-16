package indexer

import (
	"context"
	"fmt"

	"github.com/0glabs/0g-storage-client/common"
	"github.com/0glabs/0g-storage-client/common/shard"
	"github.com/0glabs/0g-storage-client/contract"
	"github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/0glabs/0g-storage-client/transfer"
	eth_common "github.com/ethereum/go-ethereum/common"
	"github.com/openweb3/go-rpc-provider/interfaces"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
	"github.com/pkg/errors"
)

// Requires `Client` implements the `Interface` interface.
var _ Interface = (*Client)(nil)

// Client indexer client
type Client struct {
	interfaces.Provider
	option IndexerClientOption
}

// IndexerClientOption indexer client option
type IndexerClientOption struct {
	ProviderOption providers.Option
	LogOption      common.LogOption // log option when uploading data
}

// NewClient create new indexer client, url is indexer service url
func NewClient(url string, option ...IndexerClientOption) (*Client, error) {
	var opt IndexerClientOption
	if len(option) > 0 {
		opt = option[0]
	}

	provider, err := providers.NewProviderWithOption(url, opt.ProviderOption)
	if err != nil {
		return nil, err
	}

	return &Client{
		Provider: provider,
		option:   opt,
	}, nil
}

// GetNodes get node list from indexer service
func (c *Client) GetNodes(ctx context.Context) (nodes ShardedNodes, err error) {
	err = c.Provider.CallContext(ctx, &nodes, "indexer_getNodes")
	return
}

// SelectNodes get node list from indexer service and select a subset of it, which is sufficient to store expected number of replications.
func (c *Client) SelectNodes(ctx context.Context, expectedReplica uint) ([]*node.ZgsClient, error) {
	nodes, err := c.GetNodes(ctx)
	if err != nil {
		return nil, err
	}
	trusted, ok := shard.Select(nodes.Trusted, expectedReplica)
	if !ok {
		return nil, fmt.Errorf("cannot select a subset from the returned nodes that meets the replication requirement")
	}
	clients := make([]*node.ZgsClient, len(trusted))
	for i, shardedNode := range trusted {
		clients[i], err = node.NewZgsClient(shardedNode.URL, c.option.ProviderOption)
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("failed to initialize storage node client with %v", shardedNode.URL))
		}
	}
	return clients, nil
}

// NewUploaderFromIndexerNodes return an uploader with selected storage nodes from indexer service.
func (c *Client) NewUploaderFromIndexerNodes(ctx context.Context, flow *contract.FlowContract, expectedReplica uint) (*transfer.Uploader, error) {
	clients, err := c.SelectNodes(ctx, expectedReplica)
	if err != nil {
		return nil, err
	}
	return transfer.NewUploader(flow, clients, c.option.LogOption)
}

// Upload submit data to 0g storage contract, then transfer the data to the storage nodes selected from indexer service.
func (c *Client) Upload(ctx context.Context, flow *contract.FlowContract, data core.IterableData, option ...transfer.UploadOption) error {
	expectedReplica := uint(1)
	if len(option) > 0 {
		expectedReplica = max(expectedReplica, option[0].ExpectedReplica)
	}
	uploader, err := c.NewUploaderFromIndexerNodes(ctx, flow, expectedReplica)
	if err != nil {
		return err
	}
	return uploader.Upload(ctx, data, option...)
}

// BatchUpload submit multiple data to 0g storage contract batchly in single on-chain transaction, then transfer the data to the storage nodes selected from indexer service.
func (c *Client) BatchUpload(ctx context.Context, flow *contract.FlowContract, datas []core.IterableData, waitForLogEntry bool, option ...[]transfer.UploadOption) (eth_common.Hash, []eth_common.Hash, error) {
	expectedReplica := uint(1)
	if len(option) > 0 {
		for _, opt := range option[0] {
			expectedReplica = max(expectedReplica, opt.ExpectedReplica)
		}
	}
	uploader, err := c.NewUploaderFromIndexerNodes(ctx, flow, expectedReplica)
	if err != nil {
		return eth_common.Hash{}, nil, err
	}
	return uploader.BatchUpload(ctx, datas, waitForLogEntry, option...)
}
