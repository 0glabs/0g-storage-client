package kv

import (
	"context"
	"math"

	zg_common "github.com/0glabs/0g-storage-client/common"
	"github.com/0glabs/0g-storage-client/contract"
	"github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

func (c *Client) GetValue(ctx context.Context, streamId common.Hash, key []byte, version ...uint64) (val *node.Value, err error) {
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
		seg, err = c.node.KV().GetValue(ctx, streamId, key, uint64(len(val.Data)), maxQuerySize, val.Version)
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
func (c *Client) Get(ctx context.Context, streamId common.Hash, key []byte, startIndex, length uint64, version ...uint64) (val *node.Value, err error) {
	return c.node.KV().GetValue(ctx, streamId, key, startIndex, length, version...)
}

func (c *Client) GetNext(ctx context.Context, streamId common.Hash, key []byte, startIndex, length uint64, inclusive bool, version ...uint64) (val *node.KeyValue, err error) {
	return c.node.KV().GetNext(ctx, streamId, key, startIndex, length, inclusive, version...)
}

func (c *Client) GetPrev(ctx context.Context, streamId common.Hash, key []byte, startIndex, length uint64, inclusive bool, version ...uint64) (val *node.KeyValue, err error) {
	return c.node.KV().GetPrev(ctx, streamId, key, startIndex, length, inclusive, version...)
}

func (c *Client) GetFirst(ctx context.Context, streamId common.Hash, startIndex, length uint64, version ...uint64) (val *node.KeyValue, err error) {
	return c.node.KV().GetFirst(ctx, streamId, startIndex, length, version...)
}

func (c *Client) GetLast(ctx context.Context, streamId common.Hash, startIndex, length uint64, version ...uint64) (val *node.KeyValue, err error) {
	return c.node.KV().GetLast(ctx, streamId, startIndex, length, version...)
}

func (c *Client) GetTransactionResult(ctx context.Context, txSeq uint64) (result string, err error) {
	return c.node.KV().GetTransactionResult(ctx, txSeq)
}

func (c *Client) GetHoldingStreamIds(ctx context.Context) (streamIds []common.Hash, err error) {
	return c.node.KV().GetHoldingStreamIds(ctx)
}

func (c *Client) HasWritePermission(ctx context.Context, account common.Address, streamId common.Hash, key []byte, version ...uint64) (hasPermission bool, err error) {
	return c.node.KV().HasWritePermission(ctx, account, streamId, key, version...)
}

func (c *Client) IsAdmin(ctx context.Context, account common.Address, streamId common.Hash, version ...uint64) (isAdmin bool, err error) {
	return c.node.KV().IsAdmin(ctx, account, streamId, version...)
}

func (c *Client) IsSpecialKey(ctx context.Context, streamId common.Hash, key []byte, version ...uint64) (isSpecialKey bool, err error) {
	return c.node.KV().IsSpecialKey(ctx, streamId, key, version...)
}

func (c *Client) IsWriterOfKey(ctx context.Context, account common.Address, streamId common.Hash, key []byte, version ...uint64) (isWriter bool, err error) {
	return c.node.KV().IsWriterOfKey(ctx, account, streamId, key, version...)
}

func (c *Client) IsWriterOfStream(ctx context.Context, account common.Address, streamId common.Hash, version ...uint64) (isWriter bool, err error) {
	return c.node.KV().IsWriterOfStream(ctx, account, streamId, version...)
}

// Batcher returns a Batcher instance for kv operations in batch.
func (c *Client) Batcher(opts ...zg_common.LogOption) *Batcher {
	return newBatcher(math.MaxUint64, c, opts...)
}

type Batcher struct {
	*StreamDataBuilder
	client *Client
	logger *logrus.Logger
}

func newBatcher(version uint64, client *Client, opts ...zg_common.LogOption) *Batcher {
	return &Batcher{
		StreamDataBuilder: NewStreamDataBuilder(version),
		client:            client,
		logger:            zg_common.NewLogger(opts...),
	}
}

// Exec submit the kv operations to ZeroGStorage network in batch.
//
// Note, this is a time consuming operation, e.g. several seconds or even longer.
// When it comes to a time sentitive context, it should be executed in a separate go-routine.
func (b *Batcher) Exec(ctx context.Context) error {
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
	uploader, err := transfer.NewUploader(b.client.flow, []*node.Client{b.client.node}, zg_common.LogOption{Logger: b.logger})
	if err != nil {
		return err
	}
	opt := transfer.UploadOption{
		Tags:  b.BuildTags(),
		Force: true,
	}
	if err = uploader.Upload(ctx, data, opt); err != nil {
		return errors.WithMessagef(err, "Failed to upload data")
	}
	return nil
}
