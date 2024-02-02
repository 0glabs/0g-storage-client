package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/zero-gravity-labs/zerog-storage-client/common/parallel"
	"github.com/zero-gravity-labs/zerog-storage-client/core/merkle"
)

type TreeBuilderInitializer struct {
	data    IterableData
	offset  int64
	batch   int64
	builder *merkle.TreeBuilder
}

var _ parallel.Interface = (*TreeBuilderInitializer)(nil)

// ParallelCollect implements parallel.Interface.
func (t *TreeBuilderInitializer) ParallelCollect(result *parallel.Result) error {
	t.builder.AppendHash(result.Value.(common.Hash))
	return nil
}

// ParallelDo implements parallel.Interface.
func (t *TreeBuilderInitializer) ParallelDo(routine int, task int) (interface{}, error) {
	offset := t.offset + int64(task)*t.batch
	buf, err := ReadAt(t.data, int(t.batch), offset, t.data.PaddedSize())
	if err != nil {
		return nil, err
	}

	hash := SegmentRoot(buf)
	return hash, nil
}
