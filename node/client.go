package node

import (
	"context"

	"github.com/0glabs/0g-storage-client/common/shard"
	"github.com/ethereum/go-ethereum/common"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
	"github.com/sirupsen/logrus"
)

// Client Common client connected to a url, can be derived to different specific clients.
type Client struct {
	url string
	*providers.MiddlewarableProvider

	zgs   *ZeroGStorageClient
	admin *AdminClient
	kv    *KvClient
}

// MustNewClient Initalize a client and panic on failure.
func MustNewClient(url string, option ...providers.Option) *Client {
	client, err := NewClient(url, option...)
	if err != nil {
		logrus.WithError(err).WithField("url", url).Fatal("Failed to connect to storage node")
	}

	return client
}

// NewClient Initalize a client.
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
		url:                   url,
		MiddlewarableProvider: provider,

		zgs:   &ZeroGStorageClient{provider},
		admin: &AdminClient{provider},
		kv:    &KvClient{provider},
	}, nil
}

// MustNewClients Initialize a list of clients and panic on failure.
func MustNewClients(urls []string, option ...providers.Option) []*Client {
	var clients []*Client

	for _, url := range urls {
		client := MustNewClient(url, option...)
		clients = append(clients, client)
	}

	return clients
}

// URL Get the RPC server URL the client connected to.
func (c *Client) URL() string {
	return c.url
}

// ZeroGStorage Get a ZeroGStorageClient connected to current client url.
func (c *Client) ZeroGStorage() *ZeroGStorageClient {
	return c.zgs
}

// Admin Get an AdminClient connected to current client url.
func (c *Client) Admin() *AdminClient {
	return c.admin
}

// KV Get a KVClient connected to current client url.
func (c *Client) KV() *KvClient {
	return c.kv
}

// ZeroGStorageClient RPC Client connected to a 0g storage node's RPC endpoint.
type ZeroGStorageClient struct {
	provider *providers.MiddlewarableProvider
}

// GetStatus Call zgs_getStatus RPC to get sync status of the node.
func (c *ZeroGStorageClient) GetStatus(ctx context.Context) (status Status, err error) {
	err = c.provider.CallContext(ctx, &status, "zgs_getStatus")
	return
}

// GetFileInfo Call zgs_getFileInfo RPC to get the information of a file by file data root from the node.
func (c *ZeroGStorageClient) GetFileInfo(ctx context.Context, root common.Hash) (file *FileInfo, err error) {
	err = c.provider.CallContext(ctx, &file, "zgs_getFileInfo", root)
	return
}

// GetFileInfoByTxSeq Call zgs_getFileInfoByTxSeq RPC to get the information of a file by file sequence id from the node.
func (c *ZeroGStorageClient) GetFileInfoByTxSeq(ctx context.Context, txSeq uint64) (file *FileInfo, err error) {
	err = c.provider.CallContext(ctx, &file, "zgs_getFileInfoByTxSeq", txSeq)
	return
}

// UploadSegment Call zgs_uploadSegment RPC to upload a segment to the node.
func (c *ZeroGStorageClient) UploadSegment(ctx context.Context, segment SegmentWithProof) (ret int, err error) {
	err = c.provider.CallContext(ctx, &ret, "zgs_uploadSegment", segment)
	return
}

// UploadSegments Call zgs_uploadSegments RPC to upload a slice of segments to the node.
func (c *ZeroGStorageClient) UploadSegments(ctx context.Context, segments []SegmentWithProof) (ret int, err error) {
	err = c.provider.CallContext(ctx, &ret, "zgs_uploadSegments", segments)
	return
}

// DownloadSegment Call zgs_downloadSegment RPC to download a segment from the node.
func (c *ZeroGStorageClient) DownloadSegment(ctx context.Context, root common.Hash, startIndex, endIndex uint64) (data []byte, err error) {
	err = c.provider.CallContext(ctx, &data, "zgs_downloadSegment", root, startIndex, endIndex)
	if len(data) == 0 {
		return nil, err
	}
	return
}

// DownloadSegmentWithProof Call zgs_downloadSegmentWithProof RPC to download a segment along with its merkle proof from the node.
func (c *ZeroGStorageClient) DownloadSegmentWithProof(ctx context.Context, root common.Hash, index uint64) (segment *SegmentWithProof, err error) {
	err = c.provider.CallContext(ctx, &segment, "zgs_downloadSegmentWithProof", root, index)
	return
}

// GetShardConfig Call zgs_getShardConfig RPC to get the current shard configuration of the node.
func (c *ZeroGStorageClient) GetShardConfig(ctx context.Context) (shardConfig shard.ShardConfig, err error) {
	err = c.provider.CallContext(ctx, &shardConfig, "zgs_getShardConfig")
	return
}

// AdminClient RPC Client connected to a 0g storage node's admin RPC endpoint.
type AdminClient struct {
	provider *providers.MiddlewarableProvider
}

// Shutdown Call admin_shutdown to shutdown the node.
func (c *AdminClient) Shutdown(ctx context.Context) (ret int, err error) {
	err = c.provider.CallContext(ctx, &ret, "admin_shutdown")
	return
}

// StartSyncFile Call admin_startSyncFile to request synchronization of a file.
func (c *AdminClient) StartSyncFile(ctx context.Context, txSeq uint64) (ret int, err error) {
	err = c.provider.CallContext(ctx, &ret, "admin_startSyncFile", txSeq)
	return
}

// GetSyncStatus Call admin_getSyncStatus to check synchronization status of a file.
func (c *AdminClient) GetSyncStatus(ctx context.Context, txSeq uint64) (status string, err error) {
	err = c.provider.CallContext(ctx, &status, "admin_getSyncStatus", txSeq)
	return
}
