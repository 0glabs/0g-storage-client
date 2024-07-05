package node

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
)

type KvClient struct {
	provider *providers.MiddlewarableProvider
}

func (c *KvClient) GetValue(ctx context.Context, streamId common.Hash, key []byte, startIndex, length uint64, version ...uint64) (val *Value, err error) {
	args := []interface{}{streamId, key, startIndex, length}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(ctx, &val, "kv_getValue", args...)
	return
}

func (c *KvClient) GetNext(ctx context.Context, streamId common.Hash, key []byte, startIndex, length uint64, inclusive bool, version ...uint64) (val *KeyValue, err error) {
	args := []interface{}{streamId, key, startIndex, length, inclusive}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(ctx, &val, "kv_getNext", args...)
	return
}

func (c *KvClient) GetPrev(ctx context.Context, streamId common.Hash, key []byte, startIndex, length uint64, inclusive bool, version ...uint64) (val *KeyValue, err error) {
	args := []interface{}{streamId, key, startIndex, length, inclusive}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(ctx, &val, "kv_getPrev", args...)
	return
}

func (c *KvClient) GetFirst(ctx context.Context, streamId common.Hash, startIndex, length uint64, version ...uint64) (val *KeyValue, err error) {
	args := []interface{}{streamId, startIndex, length}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(ctx, &val, "kv_getFirst", args...)
	return
}

func (c *KvClient) GetLast(ctx context.Context, streamId common.Hash, startIndex, length uint64, version ...uint64) (val *KeyValue, err error) {
	args := []interface{}{streamId, startIndex, length}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(ctx, &val, "kv_getLast", args...)
	return
}

func (c *KvClient) GetTransactionResult(ctx context.Context, txSeq uint64) (result string, err error) {
	err = c.provider.CallContext(ctx, &result, "kv_getTransactionResult", txSeq)
	return
}

func (c *KvClient) GetHoldingStreamIds(ctx context.Context) (streamIds []common.Hash, err error) {
	err = c.provider.CallContext(ctx, &streamIds, "kv_getHoldingStreamIds")
	return
}

func (c *KvClient) HasWritePermission(ctx context.Context, account common.Address, streamId common.Hash, key []byte, version ...uint64) (hasPermission bool, err error) {
	args := []interface{}{account, streamId, key}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(ctx, &hasPermission, "kv_hasWritePermission", args...)
	return
}

func (c *KvClient) IsAdmin(ctx context.Context, account common.Address, streamId common.Hash, version ...uint64) (isAdmin bool, err error) {
	args := []interface{}{account, streamId}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(ctx, &isAdmin, "kv_isAdmin", args...)
	return
}

func (c *KvClient) IsSpecialKey(ctx context.Context, streamId common.Hash, key []byte, version ...uint64) (isSpecialKey bool, err error) {
	args := []interface{}{streamId, key}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(ctx, &isSpecialKey, "kv_isSpecialKey", args...)
	return
}

func (c *KvClient) IsWriterOfKey(ctx context.Context, account common.Address, streamId common.Hash, key []byte, version ...uint64) (isWriter bool, err error) {
	args := []interface{}{account, streamId, key}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(ctx, &isWriter, "kv_isWriterOfKey", args...)
	return
}

func (c *KvClient) IsWriterOfStream(ctx context.Context, account common.Address, streamId common.Hash, version ...uint64) (isWriter bool, err error) {
	args := []interface{}{account, streamId}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(ctx, &isWriter, "kv_isWriterOfStream", args...)
	return
}
