package transfer

import (
	"github.com/0glabs/0g-storage-client/common/parallel"
	"github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/core/merkle"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type SegmentUploader struct {
	data     core.IterableData
	tree     *merkle.Tree
	clients  []*node.ZeroGStorageClient
	offset   int64
	disperse bool
	taskSize uint
	numTasks int
}

var _ parallel.Interface = (*SegmentUploader)(nil)

// ParallelCollect implements parallel.Interface.
func (uploader *SegmentUploader) ParallelCollect(result *parallel.Result) error {
	return nil
}

// ParallelDo implements parallel.Interface.
func (uploader *SegmentUploader) ParallelDo(routine int, task int) (interface{}, error) {
	offset := uploader.offset + int64(task)*core.DefaultSegmentSize
	numChunks := uploader.data.NumChunks()
	numSegments := uploader.data.NumSegments()
	segIndex := uint64(offset / core.DefaultSegmentSize)
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
		segment, err := core.ReadAt(uploader.data, core.DefaultSegmentSize, offset, uploader.data.PaddedSize())
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
		segIndex += uint64(uploader.numTasks)
		offset += core.DefaultSegmentSize * int64(uploader.numTasks)
	}
	// upload
	if !uploader.disperse {
		if _, err := uploader.clients[0].UploadSegments(segments); err != nil && !isDuplicateError(err.Error()) {
			return nil, errors.WithMessage(err, "Failed to upload segment")
		}
	} else {
		clientIndex := task % (len(uploader.clients))
		ok := false
		// retry
		for i := 0; i < len(uploader.clients); i++ {
			logrus.WithFields(logrus.Fields{
				"total":          numSegments,
				"from_seg_index": startSegIndex,
				"to_seg_index":   segIndex,
				"step":           uploader.numTasks,
				"clientIndex":    clientIndex,
			}).Debug("Uploading segment to node..")
			if _, err := uploader.clients[clientIndex].UploadSegments(segments); err != nil && !isDuplicateError(err.Error()) {
				logrus.WithFields(logrus.Fields{
					"total":          numSegments,
					"from_seg_index": startSegIndex,
					"to_seg_index":   segIndex,
					"step":           uploader.numTasks,
					"clientIndex":    clientIndex,
					"error":          err,
				}).Warn("Failed to upload segment to node, try next node..")
				clientIndex = (clientIndex + 1) % (len(uploader.clients))
			} else {
				ok = true
				break
			}
		}
		if !ok {
			if _, err := uploader.clients[clientIndex].UploadSegments(segments); err != nil {
				return nil, errors.WithMessage(err, "Failed to upload segment")
			}
		}
	}

	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		logrus.WithFields(logrus.Fields{
			"total":          numSegments,
			"from_seg_index": startSegIndex,
			"to_seg_index":   segIndex,
			"step":           uploader.numTasks,
			"root":           core.SegmentRoot(segments[0].Data),
		}).Debug("Segments uploaded")
	}
	return nil, nil
}
