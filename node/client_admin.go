package node

import (
	"context"

	"github.com/0glabs/0g-storage-client/common/rpc"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
	"github.com/sirupsen/logrus"
)

// AdminClient RPC Client connected to a 0g storage node's admin RPC endpoint.
type AdminClient struct {
	*rpcClient
}

// MustNewAdminClient initalize an admin client and panic on failure.
func MustNewAdminClient(url string, option ...providers.Option) *AdminClient {
	client, err := NewAdminClient(url, option...)
	if err != nil {
		logrus.WithError(err).WithField("url", url).Fatal("Failed to create admin client")
	}

	return client
}

// NewAdminClient initalize an admin client.
func NewAdminClient(url string, option ...providers.Option) (*AdminClient, error) {
	client, err := newRpcClient(url, option...)
	if err != nil {
		return nil, err
	}

	return &AdminClient{client}, nil
}

// FindFile Call find_file to update file location cache
func (c *AdminClient) FindFile(ctx context.Context, txSeq uint64) (int, error) {
	return rpc.CallContext[int](c, ctx, "admin_findFile", txSeq)
}

// Shutdown Call admin_shutdown to shutdown the node.
func (c *AdminClient) Shutdown(ctx context.Context) (int, error) {
	return rpc.CallContext[int](c, ctx, "admin_shutdown")
}

// StartSyncFile Call admin_startSyncFile to request synchronization of a file.
func (c *AdminClient) StartSyncFile(ctx context.Context, txSeq uint64) (int, error) {
	return rpc.CallContext[int](c, ctx, "admin_startSyncFile", txSeq)
}

// StartSyncChunks Call admin_startSyncChunks to request synchronization of specified chunks.
func (c *AdminClient) StartSyncChunks(ctx context.Context, txSeq, startIndex, endIndex uint64) (int, error) {
	return rpc.CallContext[int](c, ctx, "admin_startSyncChunks", txSeq, startIndex, endIndex)
}

// TerminateSync Call admin_terminateSync to terminate a file sync.
func (c *AdminClient) TerminateSync(ctx context.Context, txSeq uint64) (bool, error) {
	return rpc.CallContext[bool](c, ctx, "admin_terminateSync", txSeq)
}

// GetSyncStatus Call admin_getSyncStatus to retrieve the sync status of specified file.
func (c *AdminClient) GetSyncStatus(ctx context.Context, txSeq uint64) (string, error) {
	return rpc.CallContext[string](c, ctx, "admin_getSyncStatus", txSeq)
}

// GetSyncInfo Call admin_getSyncInfo to retrieve the sync status of specified file or all files.
func (c *AdminClient) GetSyncInfo(ctx context.Context, tx_seq ...uint64) (map[uint64]FileSyncInfo, error) {
	if len(tx_seq) > 0 {
		return rpc.CallContext[map[uint64]FileSyncInfo](c, ctx, "admin_getSyncInfo", tx_seq[0])
	}

	return rpc.CallContext[map[uint64]FileSyncInfo](c, ctx, "admin_getSyncInfo")
}

// GetNetworkInfo Call admin_getNetworkInfo to retrieve the network information.
func (c *AdminClient) GetNetworkInfo(ctx context.Context) (NetworkInfo, error) {
	return rpc.CallContext[NetworkInfo](c, ctx, "admin_getNetworkInfo")
}

// GetPeers Call admin_getPeers to retrieve all discovered network peers.
func (c *AdminClient) GetPeers(ctx context.Context) (map[string]*PeerInfo, error) {
	return rpc.CallContext[map[string]*PeerInfo](c, ctx, "admin_getPeers")
}

// GetFileLocation Get file location
func (c *AdminClient) GetFileLocation(ctx context.Context, txSeq uint64, allShards bool) ([]LocationInfo, error) {
	return rpc.CallContext[[]LocationInfo](c, ctx, "admin_getFileLocation", txSeq, allShards)
}
