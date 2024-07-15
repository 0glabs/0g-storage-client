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
	"github.com/sirupsen/logrus"
)

// Requires `Client` implements the `Interface` interface.
var _ Interface = (*Client)(nil)

type Client struct {
	interfaces.Provider
	providerOption providers.Option
}

func NewClient(url string, option ...providers.Option) (*Client, error) {
	var opt providers.Option
	if len(option) > 0 {
		opt = option[0]
	}

	provider, err := providers.NewProviderWithOption(url, opt)
	if err != nil {
		return nil, err
	}

	return &Client{
		Provider:       provider,
		providerOption: opt,
	}, nil
}

func (c *Client) GetNodes(ctx context.Context) (nodes []shard.ShardedNode, err error) {
	err = c.Provider.CallContext(ctx, &nodes, "indexer_getNodes")
	return
}

func (c *Client) NewUploaderFromIndexerNodes(ctx context.Context, flow *contract.FlowContract, expectedReplica uint) (*transfer.Uploader, error) {
	nodes, err := c.GetNodes(ctx)
	if err != nil {
		return nil, err
	}
	nodes, ok := shard.Select(nodes, expectedReplica)
	if !ok {
		return nil, fmt.Errorf("cannot select a subset from the returned nodes that meets the replication requirement")
	}
	clients := make([]*node.Client, len(nodes))
	for i, shardedNode := range nodes {
		clients[i], err = node.NewClient(shardedNode.URL, c.providerOption)
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("failed to initialize storage node client with %v", shardedNode.URL))
		}
	}
	logger := logrus.New()
	logger.Out = c.providerOption.Logger
	return transfer.NewUploader(flow, clients, common.LogOption{Logger: logger})
}

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
