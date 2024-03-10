package kv

import (
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

func (iter *Iterator) move(kv *node.KeyValue) error {
	if kv == nil {
		iter.currentPair = nil
		return nil
	}
	value, err := iter.client.GetValue(iter.streamId, kv.Key, iter.version)
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

func (iter *Iterator) SeekBefore(key []byte) error {
	kv, err := iter.client.GetPrev(iter.streamId, key, 0, 0, true, iter.version)
	if err != nil {
		return err
	}
	return iter.move(kv)
}

func (iter *Iterator) SeekAfter(key []byte) error {
	kv, err := iter.client.GetNext(iter.streamId, key, 0, 0, true, iter.version)
	if err != nil {
		return err
	}
	return iter.move(kv)
}

func (iter *Iterator) SeekToFirst() error {
	kv, err := iter.client.GetFirst(iter.streamId, 0, 0, iter.version)
	if err != nil {
		return err
	}
	return iter.move(kv)
}

func (iter *Iterator) SeekToLast() error {
	kv, err := iter.client.GetLast(iter.streamId, 0, 0, iter.version)
	if err != nil {
		return err
	}
	return iter.move(kv)
}

func (iter *Iterator) Next() error {
	if !iter.Valid() {
		return errIteratorInvalid
	}
	kv, err := iter.client.GetNext(iter.streamId, iter.currentPair.Key, 0, 0, false, iter.version)
	if err != nil {
		return err
	}
	return iter.move(kv)
}

func (iter *Iterator) Prev() error {
	if !iter.Valid() {
		return errIteratorInvalid
	}
	kv, err := iter.client.GetPrev(iter.streamId, iter.currentPair.Key, 0, 0, false, iter.version)
	if err != nil {
		return err
	}
	return iter.move(kv)
}
