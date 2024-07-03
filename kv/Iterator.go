package kv

import (
	"context"
	"errors"

	"github.com/0glabs/0g-storage-client/node"
	"github.com/ethereum/go-ethereum/common"
)

var errIteratorInvalid = errors.New("iterator is invalid")

type Iterator struct {
	client      *Client
	streamId    common.Hash
	version     uint64
	currentPair *node.KeyValue
}

func (iter *Iterator) Valid() bool {
	return iter.currentPair != nil
}

func (iter *Iterator) KeyValue() *node.KeyValue {
	return iter.currentPair
}

func (iter *Iterator) move(ctx context.Context, kv *node.KeyValue) error {
	if kv == nil {
		iter.currentPair = nil
		return nil
	}
	value, err := iter.client.GetValue(ctx, iter.streamId, kv.Key, iter.version)
	if err != nil {
		return err
	}
	iter.currentPair = &node.KeyValue{
		Version: value.Version,
		Key:     kv.Key,
		Data:    value.Data,
		Size:    value.Size,
	}
	return nil
}

func (iter *Iterator) SeekBefore(ctx context.Context, key []byte) error {
	kv, err := iter.client.GetPrev(ctx, iter.streamId, key, 0, 0, true, iter.version)
	if err != nil {
		return err
	}
	return iter.move(ctx, kv)
}

func (iter *Iterator) SeekAfter(ctx context.Context, key []byte) error {
	kv, err := iter.client.GetNext(ctx, iter.streamId, key, 0, 0, true, iter.version)
	if err != nil {
		return err
	}
	return iter.move(ctx, kv)
}

func (iter *Iterator) SeekToFirst(ctx context.Context) error {
	kv, err := iter.client.GetFirst(ctx, iter.streamId, 0, 0, iter.version)
	if err != nil {
		return err
	}
	return iter.move(ctx, kv)
}

func (iter *Iterator) SeekToLast(ctx context.Context) error {
	kv, err := iter.client.GetLast(ctx, iter.streamId, 0, 0, iter.version)
	if err != nil {
		return err
	}
	return iter.move(ctx, kv)
}

func (iter *Iterator) Next(ctx context.Context) error {
	if !iter.Valid() {
		return errIteratorInvalid
	}
	kv, err := iter.client.GetNext(ctx, iter.streamId, iter.currentPair.Key, 0, 0, false, iter.version)
	if err != nil {
		return err
	}
	return iter.move(ctx, kv)
}

func (iter *Iterator) Prev(ctx context.Context) error {
	if !iter.Valid() {
		return errIteratorInvalid
	}
	kv, err := iter.client.GetPrev(ctx, iter.streamId, iter.currentPair.Key, 0, 0, false, iter.version)
	if err != nil {
		return err
	}
	return iter.move(ctx, kv)
}
