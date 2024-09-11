package node

import (
	"context"

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
func (c *AdminClient) FindFile(ctx context.Context, txSeq uint64) (ret int, err error) {
	err = c.wrapError(c.CallContext(ctx, &ret, "admin_findFile", txSeq), "admin_findFile")
	return
}

// Shutdown Call admin_shutdown to shutdown the node.
func (c *AdminClient) Shutdown(ctx context.Context) (ret int, err error) {
	err = c.wrapError(c.CallContext(ctx, &ret, "admin_shutdown"), "admin_shutdown")
	return
}

// StartSyncFile Call admin_startSyncFile to request synchronization of a file.
func (c *AdminClient) StartSyncFile(ctx context.Context, txSeq uint64) (ret int, err error) {
	err = c.wrapError(c.CallContext(ctx, &ret, "admin_startSyncFile", txSeq), "admin_startSyncFile")
	return
}

// StartSyncChunks Call admin_startSyncChunks to request synchronization of specified chunks.
func (c *AdminClient) StartSyncChunks(ctx context.Context, txSeq, startIndex, endIndex uint64) (ret int, err error) {
	err = c.wrapError(c.CallContext(ctx, &ret, "admin_startSyncChunks", txSeq, startIndex, endIndex), "admin_startSyncChunks")
	return
}

// TerminateSync Call admin_terminateSync to terminate a file sync.
func (c *AdminClient) TerminateSync(ctx context.Context, txSeq uint64) (terminated bool, err error) {
	err = c.wrapError(c.CallContext(ctx, &terminated, "admin_terminateSync", txSeq), "admin_terminateSync")
	return
}

// GetSyncStatus Call admin_getSyncStatus to retrieve the sync status of specified file.
func (c *AdminClient) GetSyncStatus(ctx context.Context, txSeq uint64) (status string, err error) {
	err = c.wrapError(c.CallContext(ctx, &status, "admin_getSyncStatus", txSeq), "admin_getSyncStatus")
	return
}

// GetSyncInfo Call admin_getSyncInfo to retrieve the sync status of specified file or all files.
func (c *AdminClient) GetSyncInfo(ctx context.Context, tx_seq ...uint64) (files map[uint64]FileSyncInfo, err error) {
	if len(tx_seq) > 0 {
		err = c.wrapError(c.CallContext(ctx, &files, "admin_getSyncInfo", tx_seq[0]), "admin_getSyncInfo")
	} else {
		err = c.wrapError(c.CallContext(ctx, &files, "admin_getSyncInfo"), "admin_getSyncInfo")
	}

	return
}

// GetNetworkInfo Call admin_getNetworkInfo to retrieve the network information.
func (c *AdminClient) GetNetworkInfo(ctx context.Context) (info NetworkInfo, err error) {
	err = c.wrapError(c.CallContext(ctx, &info, "admin_getNetworkInfo"), "admin_getNetworkInfo")
	return
}

// GetPeers Call admin_getPeers to retrieve all discovered network peers.
func (c *AdminClient) GetPeers(ctx context.Context) (peers map[string]*PeerInfo, err error) {
	err = c.wrapError(c.CallContext(ctx, &peers, "admin_getPeers"), "admin_getPeers")
	return
}

// getFileLocation Get file location
func (c *AdminClient) GetFileLocation(ctx context.Context, txSeq uint64, allShards bool) (locations []LocationInfo, err error) {
	err = c.wrapError(c.CallContext(ctx, &locations, "admin_getFileLocation", txSeq, allShards), "admin_getFileLocation")
	return
}
