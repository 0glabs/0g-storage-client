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
	// get segment
	segment, err := core.ReadAt(uploader.data, core.DefaultSegmentSize, offset, uploader.data.PaddedSize())
	if err != nil {
		return nil, err
	}
	startIndex := segIndex * core.DefaultSegmentMaxChunks
	if startIndex+uint64(len(segment))/core.DefaultChunkSize >= numChunks {
		// last segment has real data
		expectedLen := core.DefaultChunkSize * int(numChunks-startIndex)
		segment = segment[:expectedLen]
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
	// upload
	if !uploader.disperse {
		if _, err := uploader.clients[0].UploadSegment(segWithProof); err != nil && !isDuplicateError(err.Error()) {
			return nil, errors.WithMessage(err, "Failed to upload segment")
		}
	} else {
		clientIndex := segIndex % uint64(len(uploader.clients))
		ok := false
		// retry
		for i := 0; i < len(uploader.clients); i++ {
			logrus.WithFields(logrus.Fields{
				"total":       numSegments,
				"index":       segIndex,
				"clientIndex": clientIndex,
			}).Debug("Uploading segment to node..")
			if _, err := uploader.clients[clientIndex].UploadSegment(segWithProof); err != nil && !isDuplicateError(err.Error()) {
				logrus.WithFields(logrus.Fields{
					"total":       numSegments,
					"index":       segIndex,
					"clientIndex": clientIndex,
					"error":       err,
				}).Warn("Failed to upload segment to node, try next node..")
				clientIndex = (clientIndex + 1) % uint64(len(uploader.clients))
			} else {
				ok = true
				break
			}
		}
		if !ok {
			if _, err := uploader.clients[clientIndex].UploadSegment(segWithProof); err != nil {
				return nil, errors.WithMessage(err, "Failed to upload segment")
			}
		}
	}

	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		chunkIndex := segIndex * core.DefaultSegmentMaxChunks
		logrus.WithFields(logrus.Fields{
			"total":      numSegments,
			"index":      segIndex,
			"chunkStart": chunkIndex,
			"chunkEnd":   chunkIndex + uint64(len(segment))/core.DefaultChunkSize,
			"root":       core.SegmentRoot(segment),
		}).Debug("Segment uploaded")
	}
	return nil, nil
}
