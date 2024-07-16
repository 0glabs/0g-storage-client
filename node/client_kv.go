package node

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
)

// KvClient RPC client connected to 0g kv node.
type KvClient struct {
	provider *providers.MiddlewarableProvider
}

// GetValue Call kv_getValue RPC to query the value of a key.
func (c *KvClient) GetValue(ctx context.Context, streamId common.Hash, key []byte, startIndex, length uint64, version ...uint64) (val *Value, err error) {
	args := []interface{}{streamId, key, startIndex, length}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(ctx, &val, "kv_getValue", args...)
	return
}

// GetNext Call kv_getNext RPC to query the next key of a given key.
func (c *KvClient) GetNext(ctx context.Context, streamId common.Hash, key []byte, startIndex, length uint64, inclusive bool, version ...uint64) (val *KeyValue, err error) {
	args := []interface{}{streamId, key, startIndex, length, inclusive}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(ctx, &val, "kv_getNext", args...)
	return
}

// GetPrev Call kv_getNext RPC to query the prev key of a given key.
func (c *KvClient) GetPrev(ctx context.Context, streamId common.Hash, key []byte, startIndex, length uint64, inclusive bool, version ...uint64) (val *KeyValue, err error) {
	args := []interface{}{streamId, key, startIndex, length, inclusive}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(ctx, &val, "kv_getPrev", args...)
	return
}

// GetFirst Call kv_getFirst RPC to query the first key.
func (c *KvClient) GetFirst(ctx context.Context, streamId common.Hash, startIndex, length uint64, version ...uint64) (val *KeyValue, err error) {
	args := []interface{}{streamId, startIndex, length}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(ctx, &val, "kv_getFirst", args...)
	return
}

// GetLast Call kv_getLast RPC to query the last key.
func (c *KvClient) GetLast(ctx context.Context, streamId common.Hash, startIndex, length uint64, version ...uint64) (val *KeyValue, err error) {
	args := []interface{}{streamId, startIndex, length}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(ctx, &val, "kv_getLast", args...)
	return
}

// GetTransactionResult Call kv_getTransactionResult RPC to query the kv replay status of a given file.
func (c *KvClient) GetTransactionResult(ctx context.Context, txSeq uint64) (result string, err error) {
	err = c.provider.CallContext(ctx, &result, "kv_getTransactionResult", txSeq)
	return
}

// GetHoldingStreamIds Call kv_getHoldingStreamIds RPC to query the stream ids monitered by the kv node.
func (c *KvClient) GetHoldingStreamIds(ctx context.Context) (streamIds []common.Hash, err error) {
	err = c.provider.CallContext(ctx, &streamIds, "kv_getHoldingStreamIds")
	return
}

// HasWritePermission Call kv_hasWritePermission RPC to check if the account is able to write the stream.
func (c *KvClient) HasWritePermission(ctx context.Context, account common.Address, streamId common.Hash, key []byte, version ...uint64) (hasPermission bool, err error) {
	args := []interface{}{account, streamId, key}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(ctx, &hasPermission, "kv_hasWritePermission", args...)
	return
}

// IsAdmin Call kv_isAdmin RPC to check if the account is the admin of the stream.
func (c *KvClient) IsAdmin(ctx context.Context, account common.Address, streamId common.Hash, version ...uint64) (isAdmin bool, err error) {
	args := []interface{}{account, streamId}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(ctx, &isAdmin, "kv_isAdmin", args...)
	return
}

// IsSpecialKey Call kv_isSpecialKey RPC to check if the key has unique access control.
func (c *KvClient) IsSpecialKey(ctx context.Context, streamId common.Hash, key []byte, version ...uint64) (isSpecialKey bool, err error) {
	args := []interface{}{streamId, key}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(ctx, &isSpecialKey, "kv_isSpecialKey", args...)
	return
}

// IsWriterOfKey Call kv_isWriterOfKey RPC to check if the account can write the special key.
func (c *KvClient) IsWriterOfKey(ctx context.Context, account common.Address, streamId common.Hash, key []byte, version ...uint64) (isWriter bool, err error) {
	args := []interface{}{account, streamId, key}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(ctx, &isWriter, "kv_isWriterOfKey", args...)
	return
}

// IsWriterOfStream Call kv_isWriterOfStream RPC to check if the account is the writer of the stream.
func (c *KvClient) IsWriterOfStream(ctx context.Context, account common.Address, streamId common.Hash, version ...uint64) (isWriter bool, err error) {
	args := []interface{}{account, streamId}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(ctx, &isWriter, "kv_isWriterOfStream", args...)
	return
}
