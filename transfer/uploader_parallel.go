package transfer

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/zero-gravity-labs/zerog-storage-client/common/parallel"
	"github.com/zero-gravity-labs/zerog-storage-client/core"
	"github.com/zero-gravity-labs/zerog-storage-client/core/merkle"
	"github.com/zero-gravity-labs/zerog-storage-client/node"
)

type SegmentUploader struct {
	data     core.IterableData
	tree     *merkle.Tree
	clients  []*node.ZeroGStorageClient
	offset   int64
	disperse bool
	taskSize uint
}

var _ parallel.Interface = (*SegmentUploader)(nil)

// ParallelCollect implements parallel.Interface.
func (uploader *SegmentUploader) ParallelCollect(result *parallel.Result) error {
	return nil
}

// ParallelDo implements parallel.Interface.
func (uploader *SegmentUploader) ParallelDo(routine int, task int) (interface{}, error) {
	offset := uploader.offset + int64(task)*int64(uploader.taskSize)*core.DefaultSegmentSize
	numChunks := uploader.data.NumChunks()
	numSegments := uploader.data.NumSegments()
	segIndex := uint64(offset / core.DefaultSegmentSize)
	startSegIndex := segIndex
	segments := make([]node.SegmentWithProof, 0)
	for i := 0; i < int(uploader.taskSize); i++ {
		// get segment
		segment, err := core.ReadAt(uploader.data, core.DefaultSegmentSize, offset, uploader.data.PaddedSize())
		if err != nil {
			return nil, err
		}
		startIndex := segIndex * core.DefaultSegmentMaxChunks
		allDataUploaded := false
		if startIndex >= numChunks {
			// file real data already uploaded
			break
		} else if startIndex+uint64(len(segment))/core.DefaultChunkSize >= numChunks {
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
		segIndex++
		offset += core.DefaultSegmentSize
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
				"clientIndex":    clientIndex,
			}).Debug("Uploading segment to node..")
			if _, err := uploader.clients[clientIndex].UploadSegments(segments); err != nil && !isDuplicateError(err.Error()) {
				logrus.WithFields(logrus.Fields{
					"total":          numSegments,
					"from_seg_index": startSegIndex,
					"to_seg_index":   segIndex,
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
			"root":           core.SegmentRoot(segments[0].Data),
		}).Debug("Segments uploaded")
	}
	return nil, nil
}
