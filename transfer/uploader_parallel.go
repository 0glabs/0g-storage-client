package transfer

import (
	"github.com/0glabs/0g-storage-client/common/parallel"
	"github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/core/merkle"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type UploadTask struct {
	clientIndex int
	segIndex    uint64
	numShard    uint64
}

type SegmentUploader struct {
	data     core.IterableData
	tree     *merkle.Tree
	clients  []*node.ZeroGStorageClient
	tasks    []*UploadTask
	taskSize uint
}

var _ parallel.Interface = (*SegmentUploader)(nil)

// ParallelCollect implements parallel.Interface.
func (uploader *SegmentUploader) ParallelCollect(result *parallel.Result) error {
	return nil
}

// ParallelDo implements parallel.Interface.
func (uploader *SegmentUploader) ParallelDo(routine int, task int) (interface{}, error) {
	numChunks := uploader.data.NumChunks()
	numSegments := uploader.data.NumSegments()
	uploadTask := uploader.tasks[task]
	segIndex := uploadTask.segIndex
	startSegIndex := segIndex
	segments := make([]node.SegmentWithProof, 0)
	for i := 0; i < int(uploader.taskSize); i++ {
		// check segment index
		startIndex := segIndex * core.DefaultSegmentMaxChunks
		allDataUploaded := false
		if startIndex >= numChunks {
			// file real data already uploaded
			break
		}
		// get segment
		segment, err := core.ReadAt(uploader.data, core.DefaultSegmentSize, int64(segIndex*core.DefaultSegmentSize), uploader.data.PaddedSize())
		if err != nil {
			return nil, err
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
		segments = append(segments, segWithProof)
		if allDataUploaded {
			break
		}
		segIndex += uploadTask.numShard
	}
	if _, err := uploader.clients[uploadTask.clientIndex].UploadSegments(segments); err != nil && !isDuplicateError(err.Error()) {
		return nil, errors.WithMessage(err, "Failed to upload segment")
	}

	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		logrus.WithFields(logrus.Fields{
			"total":          numSegments,
			"from_seg_index": startSegIndex,
			"to_seg_index":   segIndex,
			"step":           uploadTask.numShard,
			"root":           core.SegmentRoot(segments[0].Data),
		}).Debug("Segments uploaded")
	}
	return nil, nil
}
