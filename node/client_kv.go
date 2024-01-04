package node

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
)

type KvClient struct {
	provider *providers.MiddlewarableProvider
}

func (c *KvClient) GetValue(streamId common.Hash, key []byte, startIndex, length uint64, version ...uint64) (val *Value, err error) {
	args := []interface{}{streamId, key, startIndex, length}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(context.Background(), &val, "kv_getValue", args...)
	return
}

func (c *KvClient) GetNext(streamId common.Hash, key []byte, startIndex, length uint64, inclusive bool, version ...uint64) (val *KeyValue, err error) {
	args := []interface{}{streamId, key, startIndex, length, inclusive}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(context.Background(), &val, "kv_getNext", args...)
	return
}

func (c *KvClient) GetPrev(streamId common.Hash, key []byte, startIndex, length uint64, inclusive bool, version ...uint64) (val *KeyValue, err error) {
	args := []interface{}{streamId, key, startIndex, length, inclusive}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(context.Background(), &val, "kv_getPrev", args...)
	return
}

func (c *KvClient) GetFirst(streamId common.Hash, startIndex, length uint64, version ...uint64) (val *KeyValue, err error) {
	args := []interface{}{streamId, startIndex, length}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(context.Background(), &val, "kv_getFirst", args...)
	return
}

func (c *KvClient) GetLast(streamId common.Hash, startIndex, length uint64, version ...uint64) (val *KeyValue, err error) {
	args := []interface{}{streamId, startIndex, length}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(context.Background(), &val, "kv_getLast", args...)
	return
}

func (c *KvClient) GetTransactionResult(txSeq uint64) (result string, err error) {
	err = c.provider.CallContext(context.Background(), &result, "kv_getTransactionResult", txSeq)
	return
}

func (c *KvClient) GetHoldingStreamIds() (streamIds []common.Hash, err error) {
	err = c.provider.CallContext(context.Background(), &streamIds, "kv_getHoldingStreamIds")
	return
}

func (c *KvClient) HasWritePermission(account common.Address, streamId common.Hash, key []byte, version ...uint64) (hasPermission bool, err error) {
	args := []interface{}{account, streamId, key}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(context.Background(), &hasPermission, "kv_hasWritePermission", args...)
	return
}

func (c *KvClient) IsAdmin(account common.Address, streamId common.Hash, version ...uint64) (isAdmin bool, err error) {
	args := []interface{}{account, streamId}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(context.Background(), &isAdmin, "kv_isAdmin", args...)
	return
}

func (c *KvClient) IsSpecialKey(streamId common.Hash, key []byte, version ...uint64) (isSpecialKey bool, err error) {
	args := []interface{}{streamId, key}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(context.Background(), &isSpecialKey, "kv_isSpecialKey", args...)
	return
}

func (c *KvClient) IsWriterOfKey(account common.Address, streamId common.Hash, key []byte, version ...uint64) (isWriter bool, err error) {
	args := []interface{}{account, streamId, key}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(context.Background(), &isWriter, "kv_isWriterOfKey", args...)
	return
}

func (c *KvClient) IsWriterOfStream(account common.Address, streamId common.Hash, version ...uint64) (isWriter bool, err error) {
	args := []interface{}{account, streamId}
	if len(version) > 0 {
		args = append(args, version[0])
	}
	err = c.provider.CallContext(context.Background(), &isWriter, "kv_isWriterOfStream", args...)
	return
}
