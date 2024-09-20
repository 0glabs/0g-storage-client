package node

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
	"github.com/sirupsen/logrus"
)

// KvClient RPC client connected to 0g kv node.
type KvClient struct {
	*rpcClient
}

// MustNewKvClient initalize a kv client and panic on failure.
func MustNewKvClient(url string, option ...providers.Option) *KvClient {
	client, err := NewKvClient(url, option...)
	if err != nil {
		logrus.WithError(err).WithField("url", url).Fatal("Failed to create kv client")
	}

	return client
}

// NewKvClient initalize a kv client.
func NewKvClient(url string, option ...providers.Option) (*KvClient, error) {
	client, err := newRpcClient(url, option...)
	if err != nil {
		return nil, err
	}

	return &KvClient{client}, nil
}

// GetValue Call kv_getValue RPC to query the value of a key.
func (c *KvClient) GetValue(ctx context.Context, streamId common.Hash, key []byte, startIndex, length uint64, version ...uint64) (*Value, error) {
	args := []interface{}{streamId, key, startIndex, length}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	return providers.CallContext[*Value](c, ctx, "kv_getValue", args...)
}

// GetNext Call kv_getNext RPC to query the next key of a given key.
func (c *KvClient) GetNext(ctx context.Context, streamId common.Hash, key []byte, startIndex, length uint64, inclusive bool, version ...uint64) (*KeyValue, error) {
	args := []interface{}{streamId, key, startIndex, length, inclusive}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	return providers.CallContext[*KeyValue](c, ctx, "kv_getNext", args...)
}

// GetPrev Call kv_getNext RPC to query the prev key of a given key.
func (c *KvClient) GetPrev(ctx context.Context, streamId common.Hash, key []byte, startIndex, length uint64, inclusive bool, version ...uint64) (*KeyValue, error) {
	args := []interface{}{streamId, key, startIndex, length, inclusive}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	return providers.CallContext[*KeyValue](c, ctx, "kv_getPrev", args...)
}

// GetFirst Call kv_getFirst RPC to query the first key.
func (c *KvClient) GetFirst(ctx context.Context, streamId common.Hash, startIndex, length uint64, version ...uint64) (*KeyValue, error) {
	args := []interface{}{streamId, startIndex, length}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	return providers.CallContext[*KeyValue](c, ctx, "kv_getFirst", args...)
}

// GetLast Call kv_getLast RPC to query the last key.
func (c *KvClient) GetLast(ctx context.Context, streamId common.Hash, startIndex, length uint64, version ...uint64) (*KeyValue, error) {
	args := []interface{}{streamId, startIndex, length}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	return providers.CallContext[*KeyValue](c, ctx, "kv_getLast", args...)
}

// GetTransactionResult Call kv_getTransactionResult RPC to query the kv replay status of a given file.
func (c *KvClient) GetTransactionResult(ctx context.Context, txSeq uint64) (string, error) {
	return providers.CallContext[string](c, ctx, "kv_getTransactionResult", txSeq)
}

// GetHoldingStreamIds Call kv_getHoldingStreamIds RPC to query the stream ids monitered by the kv node.
func (c *KvClient) GetHoldingStreamIds(ctx context.Context) ([]common.Hash, error) {
	return providers.CallContext[[]common.Hash](c, ctx, "kv_getHoldingStreamIds")
}

// HasWritePermission Call kv_hasWritePermission RPC to check if the account is able to write the stream.
func (c *KvClient) HasWritePermission(ctx context.Context, account common.Address, streamId common.Hash, key []byte, version ...uint64) (bool, error) {
	args := []interface{}{account, streamId, key}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	return providers.CallContext[bool](c, ctx, "kv_hasWritePermission", args...)
}

// IsAdmin Call kv_isAdmin RPC to check if the account is the admin of the stream.
func (c *KvClient) IsAdmin(ctx context.Context, account common.Address, streamId common.Hash, version ...uint64) (bool, error) {
	args := []interface{}{account, streamId}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	return providers.CallContext[bool](c, ctx, "kv_isAdmin", args...)
}

// IsSpecialKey Call kv_isSpecialKey RPC to check if the key has unique access control.
func (c *KvClient) IsSpecialKey(ctx context.Context, streamId common.Hash, key []byte, version ...uint64) (bool, error) {
	args := []interface{}{streamId, key}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	return providers.CallContext[bool](c, ctx, "kv_isSpecialKey", args...)
}

// IsWriterOfKey Call kv_isWriterOfKey RPC to check if the account can write the special key.
func (c *KvClient) IsWriterOfKey(ctx context.Context, account common.Address, streamId common.Hash, key []byte, version ...uint64) (bool, error) {
	args := []interface{}{account, streamId, key}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	return providers.CallContext[bool](c, ctx, "kv_isWriterOfKey", args...)
}

// IsWriterOfStream Call kv_isWriterOfStream RPC to check if the account is the writer of the stream.
func (c *KvClient) IsWriterOfStream(ctx context.Context, account common.Address, streamId common.Hash, version ...uint64) (bool, error) {
	args := []interface{}{account, streamId}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	return providers.CallContext[bool](c, ctx, "kv_isWriterOfStream", args...)
}
