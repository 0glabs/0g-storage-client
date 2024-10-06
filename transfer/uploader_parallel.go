package transfer

import (
	"context"
	"time"

	"github.com/0glabs/0g-storage-client/common/parallel"
	"github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/core/merkle"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type uploadTask struct {
	clientIndex int
	segIndex    uint64
	numShard    uint64
}

type segmentUploader struct {
	data     core.IterableData
	tree     *merkle.Tree
	clients  []*node.ZgsClient
	tasks    []*uploadTask
	taskSize uint
	aligned  bool
	logger   *logrus.Logger
}

var _ parallel.Interface = (*segmentUploader)(nil)

// ParallelCollect implements parallel.Interface.
func (uploader *segmentUploader) ParallelCollect(result *parallel.Result) error {
	return nil
}

func (uploader *segmentUploader) getSegment(segIndex uint64) (bool, *node.SegmentWithProof, error) {
	numChunks := uploader.data.NumChunks()
	// check segment index
	startIndex := segIndex * core.DefaultSegmentMaxChunks
	allDataUploaded := false
	if startIndex >= numChunks {
		// file real data already uploaded
		return true, nil, nil
	}
	// get segment
	segment, err := core.ReadAt(uploader.data, core.DefaultSegmentSize, int64(segIndex*core.DefaultSegmentSize), uploader.data.PaddedSize())
	if err != nil {
		return false, nil, err
	}
	if startIndex+uint64(len(segment))/core.DefaultChunkSize >= numChunks {
		// last segment has real data
		expectedLen := core.DefaultChunkSize * int(numChunks-startIndex)
		segment = segment[:expectedLen]
		allDataUploaded = true
	}
	// fill proof
	proof := uploader.tree.ProofAt(int(segIndex))
	segWithProof := node.SegmentWithProof{
		Root:     uploader.tree.Root(),
		Data:     segment,
		Index:    segIndex,
		Proof:    proof,
		FileSize: uint64(uploader.data.Size()),
	}
	return allDataUploaded, &segWithProof, nil
}

// ParallelDo implements parallel.Interface.
func (uploader *segmentUploader) ParallelDo(ctx context.Context, routine int, task int) (interface{}, error) {
	numSegments := uploader.data.NumSegments()
	uploadTask := uploader.tasks[task]
	segIndex := uploadTask.segIndex
	startSegIndex := segIndex
	segments := make([]node.SegmentWithProof, 0)
	// upload segment (segIndex-1) in case the file start entry position is not aligned with SegmentSize in flow
	if segIndex > 0 && !uploader.aligned {
		_, segWithProof, err := uploader.getSegment(segIndex - 1)
		if err != nil {
			return nil, err
		}
		if segWithProof != nil {
			segments = append(segments, *segWithProof)
		}
	}
	for i := 0; i < int(uploader.taskSize); i++ {
		allDataUploaded, segWithProof, err := uploader.getSegment(segIndex)
		if err != nil {
			return nil, err
		}
		if segWithProof != nil {
			segments = append(segments, *segWithProof)
		}
		if allDataUploaded {
			break
		}
		segIndex += uploadTask.numShard
	}

	for i := 0; i < tooManyDataRetries; i++ {
		_, err := uploader.clients[uploadTask.clientIndex].UploadSegments(ctx, segments)
		if err == nil || isDuplicateError(err.Error()) {
			break
		}

		if isTooManyDataError(err.Error()) && i < tooManyDataRetries-1 {
			time.Sleep(10 * time.Second)
			continue
		}

		return nil, errors.WithMessage(err, "Failed to upload segment")
	}

	if uploader.logger.IsLevelEnabled(logrus.DebugLevel) {
		uploader.logger.WithFields(logrus.Fields{
			"total":          numSegments,
			"from_seg_index": startSegIndex,
			"to_seg_index":   segIndex,
			"step":           uploadTask.numShard,
			"root":           core.SegmentRoot(segments[0].Data),
			"to_node":        uploader.clients[uploadTask.clientIndex].URL(),
		}).Debug("Segments uploaded")
	}
	return nil, nil
}
