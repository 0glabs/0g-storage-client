package kv

import (
	"math"

	"github.com/0glabs/0g-storage-client/contract"
	"github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

// Client is used for users to communicate with server for kv operations.
type Client struct {
	node *node.Client
	flow *contract.FlowContract
}

// NewClient creates a new client for kv operations.
//
// Generally, you could refer to the `upload` function in `cmd/upload.go` file
// for how to create storage node client and flow contract client.
func NewClient(node *node.Client, flow *contract.FlowContract) *Client {
	return &Client{
		node: node,
		flow: flow,
	}
}

func (c *Client) NewIterator(streamId common.Hash, version ...uint64) *Iterator {
	var v uint64
	v = math.MaxUint64
	if len(version) > 0 {
		v = version[0]
	}
	return &Iterator{
		client:      c,
		streamId:    streamId,
		version:     v,
		currentPair: nil,
	}
}

func (c *Client) GetValue(streamId common.Hash, key []byte, version ...uint64) (val *node.Value, err error) {
	var v uint64
	v = math.MaxUint64
	if len(version) > 0 {
		v = version[0]
	}
	val = &node.Value{
		Version: v,
		Data:    make([]byte, 0),
		Size:    0,
	}
	for {
		var seg *node.Value
		seg, err = c.node.KV().GetValue(streamId, key, uint64(len(val.Data)), maxQuerySize, val.Version)
		if err != nil {
			return
		}
		if val.Version == math.MaxUint64 {
			val.Version = seg.Version
		} else if val.Version != seg.Version {
			val.Version = seg.Version
			val.Data = make([]byte, 0)
		}
		val.Size = seg.Size
		val.Data = append(val.Data, seg.Data...)
		if uint64(len(val.Data)) == val.Size {
			return
		}
	}
}

// Get returns paginated value for the specified stream key and offset.
func (c *Client) Get(streamId common.Hash, key []byte, startIndex, length uint64, version ...uint64) (val *node.Value, err error) {
	return c.node.KV().GetValue(streamId, key, startIndex, length, version...)
}

func (c *Client) GetNext(streamId common.Hash, key []byte, startIndex, length uint64, inclusive bool, version ...uint64) (val *node.KeyValue, err error) {
	return c.node.KV().GetNext(streamId, key, startIndex, length, inclusive, version...)
}

func (c *Client) GetPrev(streamId common.Hash, key []byte, startIndex, length uint64, inclusive bool, version ...uint64) (val *node.KeyValue, err error) {
	return c.node.KV().GetPrev(streamId, key, startIndex, length, inclusive, version...)
}

func (c *Client) GetFirst(streamId common.Hash, startIndex, length uint64, version ...uint64) (val *node.KeyValue, err error) {
	return c.node.KV().GetFirst(streamId, startIndex, length, version...)
}

func (c *Client) GetLast(streamId common.Hash, startIndex, length uint64, version ...uint64) (val *node.KeyValue, err error) {
	return c.node.KV().GetLast(streamId, startIndex, length, version...)
}

func (c *Client) GetTransactionResult(txSeq uint64) (result string, err error) {
	return c.node.KV().GetTransactionResult(txSeq)
}

func (c *Client) GetHoldingStreamIds() (streamIds []common.Hash, err error) {
	return c.node.KV().GetHoldingStreamIds()
}

func (c *Client) HasWritePermission(account common.Address, streamId common.Hash, key []byte, version ...uint64) (hasPermission bool, err error) {
	return c.node.KV().HasWritePermission(account, streamId, key, version...)
}

func (c *Client) IsAdmin(account common.Address, streamId common.Hash, version ...uint64) (isAdmin bool, err error) {
	return c.node.KV().IsAdmin(account, streamId, version...)
}

func (c *Client) IsSpecialKey(streamId common.Hash, key []byte, version ...uint64) (isSpecialKey bool, err error) {
	return c.node.KV().IsSpecialKey(streamId, key, version...)
}

func (c *Client) IsWriterOfKey(account common.Address, streamId common.Hash, key []byte, version ...uint64) (isWriter bool, err error) {
	return c.node.KV().IsWriterOfKey(account, streamId, key, version...)
}

func (c *Client) IsWriterOfStream(account common.Address, streamId common.Hash, version ...uint64) (isWriter bool, err error) {
	return c.node.KV().IsWriterOfStream(account, streamId, version...)
}

// Batcher returns a Batcher instance for kv operations in batch.
func (c *Client) Batcher() *Batcher {
	return newBatcher(math.MaxUint64, c)
}

type Batcher struct {
	*StreamDataBuilder
	client *Client
}

func newBatcher(version uint64, client *Client) *Batcher {
	return &Batcher{
		StreamDataBuilder: NewStreamDataBuilder(version),
		client:            client,
	}
}

// Exec submit the kv operations to ZeroGStorage network in batch.
//
// Note, this is a time consuming operation, e.g. several seconds or even longer.
// When it comes to a time sentitive context, it should be executed in a separate go-routine.
func (b *Batcher) Exec() error {
	// build stream data
	streamData, err := b.Build()
	if err != nil {
		return errors.WithMessage(err, "Failed to build stream data")
	}

	encoded, err := streamData.Encode()
	if err != nil {
		return errors.WithMessage(err, "Failed to encode data")
	}
	data, err := core.NewDataInMemory(encoded)
	if err != nil {
		return err
	}

	// upload file
	uploader, err := transfer.NewUploader(b.client.flow, []*node.Client{b.client.node})
	if err != nil {
		return err
	}
	opt := transfer.UploadOption{
		Tags:  b.BuildTags(),
		Force: true,
	}
	if err = uploader.Upload(data, opt); err != nil {
		return errors.WithMessagef(err, "Failed to upload data")
	}
	return nil
}
