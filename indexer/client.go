package indexer

import (
	"context"
	"fmt"
	"slices"
	"time"

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

// Client indexer client
type Client struct {
	interfaces.Provider
	option IndexerClientOption
	logger *logrus.Logger
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
		logger:   common.NewLogger(opt.LogOption),
	}, nil
}

// GetNodes get node list from indexer service
func (c *Client) GetShardedNodes(ctx context.Context) (nodes ShardedNodes, err error) {
	err = c.Provider.CallContext(ctx, &nodes, "indexer_getShardedNodes")
	return
}

// GetNodes return storage nodes with IP location information.
func (c *Client) GetNodeLocations(ctx context.Context) (locations map[string]*IPLocation, err error) {
	err = c.Provider.CallContext(ctx, &locations, "indexer_getNodeLocations")
	return
}

// GetFileLocations return locations info of given file.
func (c *Client) GetFileLocations(ctx context.Context, root string) (locations []*shard.ShardedNode, err error) {
	err = c.Provider.CallContext(ctx, &locations, "indexer_getFileLocations", root)
	return
}

// SelectNodes get node list from indexer service and select a subset of it, which is sufficient to store expected number of replications.
func (c *Client) SelectNodes(ctx context.Context, expectedReplica uint, dropped []string) ([]*node.ZgsClient, error) {
	allNodes, err := c.GetShardedNodes(ctx)
	if err != nil {
		return nil, err
	}
	// filter out nodes unable to connect
	nodes := make([]*shard.ShardedNode, 0)
	for _, shardedNode := range allNodes.Trusted {
		if slices.Contains(dropped, shardedNode.URL) {
			continue
		}
		client, err := node.NewZgsClient(shardedNode.URL, c.option.ProviderOption)
		if err != nil {
			c.logger.Debugf("failed to initialize client of node %v, dropped.", shardedNode.URL)
			continue
		}
		defer client.Close()
		start := time.Now()
		config, err := client.GetShardConfig(ctx)
		if err != nil || !config.IsValid() {
			c.logger.Debugf("failed to get shard config of node %v, dropped.", shardedNode.URL)
			continue
		}

		nodes = append(nodes, &shard.ShardedNode{
			URL:     shardedNode.URL,
			Config:  config,
			Latency: time.Since(start).Milliseconds(),
		})
	}
	// randomly select proper subset
	trusted, ok := shard.Select(nodes, expectedReplica, true)
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
func (c *Client) NewUploaderFromIndexerNodes(ctx context.Context, flow *contract.FlowContract, expectedReplica uint, dropped []string) (*transfer.Uploader, error) {
	clients, err := c.SelectNodes(ctx, expectedReplica, dropped)
	if err != nil {
		return nil, err
	}
	urls := make([]string, len(clients))
	for i, client := range clients {
		urls[i] = client.URL()
	}
	c.logger.Infof("get %v storage nodes from indexer: %v", len(urls), urls)
	return transfer.NewUploader(ctx, flow, clients, c.option.LogOption)
}

// Upload submit data to 0g storage contract, then transfer the data to the storage nodes selected from indexer service.
func (c *Client) Upload(ctx context.Context, flow *contract.FlowContract, data core.IterableData, option ...transfer.UploadOption) error {
	expectedReplica := uint(1)
	if len(option) > 0 {
		expectedReplica = max(expectedReplica, option[0].ExpectedReplica)
	}
	dropped := make([]string, 0)
	for {
		uploader, err := c.NewUploaderFromIndexerNodes(ctx, flow, expectedReplica, dropped)
		if err != nil {
			return err
		}
		err = uploader.Upload(ctx, data, option...)
		var rpcError *node.RPCError
		if errors.As(err, &rpcError) {
			dropped = append(dropped, rpcError.URL)
			c.logger.Infof("dropped problematic node %v and retry..", rpcError.URL)
		} else {
			return err
		}
	}
}

// BatchUpload submit multiple data to 0g storage contract batchly in single on-chain transaction, then transfer the data to the storage nodes selected from indexer service.
func (c *Client) BatchUpload(ctx context.Context, flow *contract.FlowContract, datas []core.IterableData, waitForLogEntry bool, option ...[]transfer.UploadOption) (eth_common.Hash, []eth_common.Hash, error) {
	expectedReplica := uint(1)
	if len(option) > 0 {
		for _, opt := range option[0] {
			expectedReplica = max(expectedReplica, opt.ExpectedReplica)
		}
	}
	dropped := make([]string, 0)
	for {
		uploader, err := c.NewUploaderFromIndexerNodes(ctx, flow, expectedReplica, dropped)
		if err != nil {
			return eth_common.Hash{}, nil, err
		}
		hash, roots, err := uploader.BatchUpload(ctx, datas, waitForLogEntry, option...)
		var rpcError *node.RPCError
		if errors.As(err, &rpcError) {
			dropped = append(dropped, rpcError.URL)
			c.logger.Infof("dropped problematic node %v and retry..", rpcError.URL)
		} else {
			return hash, roots, err
		}
	}
}

// Download download file by given data root
func (c *Client) Download(ctx context.Context, root, filename string, withProof bool) error {
	locations, err := c.GetFileLocations(ctx, root)
	if err != nil {
		return errors.WithMessage(err, "failed to get file locations")
	}
	clients := make([]*node.ZgsClient, 0)
	for _, location := range locations {
		client, err := node.NewZgsClient(location.URL, c.option.ProviderOption)
		if err != nil {
			c.logger.Debugf("failed to initialize client of node %v, dropped.", location.URL)
			continue
		}
		clients = append(clients, client)
	}
	if len(clients) == 0 {
		return fmt.Errorf("no node holding the file found, FindFile triggered, try again later")
	}
	downloader, err := transfer.NewDownloader(clients, c.option.LogOption)
	if err != nil {
		return err
	}

	return downloader.Download(ctx, root, filename, withProof)
}
