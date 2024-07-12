package kv

import (
	"context"

	zg_common "github.com/0glabs/0g-storage-client/common"
	"github.com/0glabs/0g-storage-client/contract"
	"github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Batcher struct {
	*StreamDataBuilder
	clients []*node.Client
	flow    *contract.FlowContract
	logger  *logrus.Logger
}

func NewBatcher(version uint64, clients []*node.Client, flow *contract.FlowContract, opts ...zg_common.LogOption) *Batcher {
	return &Batcher{
		StreamDataBuilder: NewStreamDataBuilder(version),
		clients:           clients,
		flow:              flow,
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
	uploader, err := transfer.NewUploader(b.flow, b.clients, zg_common.LogOption{Logger: b.logger})
	if err != nil {
		return err
	}
	opt := transfer.UploadOption{
		Tags: b.BuildTags(),
	}
	if err = uploader.Upload(ctx, data, opt); err != nil {
		return errors.WithMessagef(err, "Failed to upload data")
	}
	return nil
}
