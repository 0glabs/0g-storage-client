package node

import (
	"context"

	"github.com/0glabs/0g-storage-client/common/shard"
	"github.com/ethereum/go-ethereum/common"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
	"github.com/sirupsen/logrus"
)

// ZgsClient RPC Client connected to a 0g storage node's zgs RPC endpoint.
type ZgsClient struct {
	*rpcClient
}

// MustNewZgsClient Initalize a zgs client and panic on failure.
func MustNewZgsClient(url string, option ...providers.Option) *ZgsClient {
	client, err := NewZgsClient(url, option...)
	if err != nil {
		logrus.WithError(err).WithField("url", url).Fatal("Failed to create zgs client")
	}

	return client
}

// NewZgsClient Initalize a zgs client.
func NewZgsClient(url string, option ...providers.Option) (*ZgsClient, error) {
	client, err := newRpcClient(url, option...)
	if err != nil {
		return nil, err
	}

	return &ZgsClient{client}, nil
}

// MustNewZgsClients Initialize a list of zgs clients and panic on failure.
func MustNewZgsClients(urls []string, option ...providers.Option) []*ZgsClient {
	var clients []*ZgsClient

	for _, url := range urls {
		client := MustNewZgsClient(url, option...)
		clients = append(clients, client)
	}

	return clients
}

// GetStatus Call zgs_getStatus RPC to get sync status of the node.
func (c *ZgsClient) GetStatus(ctx context.Context) (status Status, err error) {
	err = c.wrapError(c.provider.CallContext(ctx, &status, "zgs_getStatus"), "zgs_getStatus")
	return
}

// GetFileInfo Call zgs_getFileInfo RPC to get the information of a file by file data root from the node.
func (c *ZgsClient) GetFileInfo(ctx context.Context, root common.Hash) (file *FileInfo, err error) {
	err = c.wrapError(c.provider.CallContext(ctx, &file, "zgs_getFileInfo", root), "zgs_getFileInfo")
	return
}

// GetFileInfoByTxSeq Call zgs_getFileInfoByTxSeq RPC to get the information of a file by file sequence id from the node.
func (c *ZgsClient) GetFileInfoByTxSeq(ctx context.Context, txSeq uint64) (file *FileInfo, err error) {
	err = c.wrapError(c.provider.CallContext(ctx, &file, "zgs_getFileInfoByTxSeq", txSeq), "zgs_getFileInfoByTxSeq")
	return
}

// UploadSegment Call zgs_uploadSegment RPC to upload a segment to the node.
func (c *ZgsClient) UploadSegment(ctx context.Context, segment SegmentWithProof) (ret int, err error) {
	err = c.wrapError(c.provider.CallContext(ctx, &ret, "zgs_uploadSegment", segment), "zgs_uploadSegment")
	return
}

// UploadSegments Call zgs_uploadSegments RPC to upload a slice of segments to the node.
func (c *ZgsClient) UploadSegments(ctx context.Context, segments []SegmentWithProof) (ret int, err error) {
	err = c.wrapError(c.provider.CallContext(ctx, &ret, "zgs_uploadSegments", segments), "zgs_uploadSegments")
	return
}

// DownloadSegment Call zgs_downloadSegment RPC to download a segment from the node.
func (c *ZgsClient) DownloadSegment(ctx context.Context, root common.Hash, startIndex, endIndex uint64) (data []byte, err error) {
	err = c.wrapError(c.provider.CallContext(ctx, &data, "zgs_downloadSegment", root, startIndex, endIndex), "zgs_downloadSegment")
	if len(data) == 0 {
		return nil, err
	}
	return
}

// DownloadSegmentWithProof Call zgs_downloadSegmentWithProof RPC to download a segment along with its merkle proof from the node.
func (c *ZgsClient) DownloadSegmentWithProof(ctx context.Context, root common.Hash, index uint64) (segment *SegmentWithProof, err error) {
	err = c.wrapError(c.provider.CallContext(ctx, &segment, "zgs_downloadSegmentWithProof", root, index), "zgs_downloadSegmentWithProof")
	return
}

// GetShardConfig Call zgs_getShardConfig RPC to get the current shard configuration of the node.
func (c *ZgsClient) GetShardConfig(ctx context.Context) (shardConfig shard.ShardConfig, err error) {
	err = c.wrapError(c.provider.CallContext(ctx, &shardConfig, "zgs_getShardConfig"), "zgs_getShardConfig")
	return
}
