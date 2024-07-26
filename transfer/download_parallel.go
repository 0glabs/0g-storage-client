package transfer

import (
	"context"
	"fmt"
	"runtime"

	"github.com/0glabs/0g-storage-client/common/parallel"
	"github.com/0glabs/0g-storage-client/common/shard"
	"github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/0glabs/0g-storage-client/transfer/download"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type segmentDownloader struct {
	clients      []*node.ZgsClient
	shardConfigs []*shard.ShardConfig
	file         *download.DownloadingFile

	withProof bool

	segmentOffset uint64
	numChunks     uint64
	numSegments   uint64

	logger *logrus.Logger
}

var _ parallel.Interface = (*segmentDownloader)(nil)

func newSegmentDownloader(clients []*node.ZgsClient, shardConfigs []*shard.ShardConfig, file *download.DownloadingFile, withProof bool, logger *logrus.Logger) (*segmentDownloader, error) {
	offset := file.Metadata().Offset
	if offset%core.DefaultSegmentSize > 0 {
		return nil, errors.Errorf("Invalid data offset in downloading file %v", offset)
	}

	fileSize := file.Metadata().Size

	return &segmentDownloader{
		clients:      clients,
		shardConfigs: shardConfigs,
		file:         file,

		withProof: withProof,

		segmentOffset: uint64(offset / core.DefaultSegmentSize),
		numChunks:     core.NumSplits(fileSize, core.DefaultChunkSize),
		numSegments:   core.NumSplits(fileSize, core.DefaultSegmentSize),

		logger: logger,
	}, nil
}

// Download downloads segments in parallel.
func (downloader *segmentDownloader) Download(ctx context.Context) error {
	numTasks := downloader.numSegments - downloader.segmentOffset

	return parallel.Serial(ctx, downloader, int(numTasks), runtime.GOMAXPROCS(0), 0)
}

// ParallelDo implements the parallel.Interface interface.
func (downloader *segmentDownloader) ParallelDo(ctx context.Context, routine, task int) (interface{}, error) {
	segmentIndex := downloader.segmentOffset + uint64(task)
	startIndex := segmentIndex * core.DefaultSegmentMaxChunks
	endIndex := startIndex + core.DefaultSegmentMaxChunks
	if endIndex > downloader.numChunks {
		endIndex = downloader.numChunks
	}

	root := downloader.file.Metadata().Root

	var (
		segment []byte
		err     error
	)

	for i := 0; i < len(downloader.shardConfigs); i += 1 {
		nodeIndex := (routine + i) % len(downloader.shardConfigs)
		if segmentIndex%downloader.shardConfigs[nodeIndex].NumShard != downloader.shardConfigs[nodeIndex].ShardId {
			continue
		}
		// try download from current node
		if downloader.withProof {
			segment, err = downloader.downloadWithProof(ctx, downloader.clients[nodeIndex], root, startIndex, endIndex)
		} else {
			segment, err = downloader.clients[nodeIndex].DownloadSegment(ctx, root, startIndex, endIndex)
		}

		if err != nil {
			downloader.logger.WithError(err).WithFields(logrus.Fields{
				"node index": nodeIndex,
				"segment":    fmt.Sprintf("%v/%v", segmentIndex, downloader.numSegments),
				"chunks":     fmt.Sprintf("[%v, %v)", startIndex, endIndex),
			}).Error("Failed to download segment")
			continue
		}
		if segment == nil {
			downloader.logger.WithFields(logrus.Fields{
				"node index": nodeIndex,
				"segment":    fmt.Sprintf("%v/%v", segmentIndex, downloader.numSegments),
				"chunks":     fmt.Sprintf("[%v, %v)", startIndex, endIndex),
			}).Warn("segment not found")
			continue
		}
		if len(segment)%core.DefaultChunkSize != 0 {
			downloader.logger.WithFields(logrus.Fields{
				"node index": nodeIndex,
				"segment":    fmt.Sprintf("%v/%v", segmentIndex, downloader.numSegments),
				"chunks":     fmt.Sprintf("[%v, %v)", startIndex, endIndex),
			}).Warn("invalid segment length")
			continue
		}
		if downloader.logger.IsLevelEnabled(logrus.TraceLevel) {
			downloader.logger.WithFields(logrus.Fields{
				"node index": nodeIndex,
				"segment":    fmt.Sprintf("%v/%v", segmentIndex, downloader.numSegments),
				"chunks":     fmt.Sprintf("[%v, %v)", startIndex, endIndex),
			}).Trace("Succeeded to download segment")
		}

		// remove paddings for the last chunk
		if segmentIndex == downloader.numSegments-1 {
			fileSize := downloader.file.Metadata().Size
			if lastChunkSize := fileSize % core.DefaultChunkSize; lastChunkSize > 0 {
				paddings := core.DefaultChunkSize - lastChunkSize
				segment = segment[0 : len(segment)-int(paddings)]
			}
		}
		return segment, nil
	}
	return nil, fmt.Errorf("failed to download segment %v", segmentIndex)
}

// ParallelCollect implements the parallel.Interface interface.
func (downloader *segmentDownloader) ParallelCollect(result *parallel.Result) error {
	return downloader.file.Write(result.Value.([]byte))
}

func (downloader *segmentDownloader) downloadWithProof(ctx context.Context, client *node.ZgsClient, root common.Hash, startIndex, endIndex uint64) ([]byte, error) {
	segmentIndex := startIndex / core.DefaultSegmentMaxChunks

	segment, err := client.DownloadSegmentWithProof(ctx, root, segmentIndex)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to download segment with proof from storage node")
	}
	if segment == nil {
		return nil, nil
	}

	if expectedDataLen := (endIndex - startIndex) * core.DefaultChunkSize; int(expectedDataLen) != len(segment.Data) {
		return nil, errors.Errorf("Downloaded data length mismatch, expected = %v, actual = %v", expectedDataLen, len(segment.Data))
	}

	numChunksFlowPadded, _ := core.ComputePaddedSize(downloader.numChunks)
	numSegmentsFlowPadded := (numChunksFlowPadded-1)/core.DefaultSegmentMaxChunks + 1

	// pad empty chunks for the last segment to validate merkle proof
	var emptyChunksPadded uint64
	if numChunks := endIndex - startIndex; numChunks < core.DefaultSegmentMaxChunks {
		if segmentIndex < numSegmentsFlowPadded-1 || numChunksFlowPadded%core.DefaultSegmentMaxChunks == 0 {
			// pad empty chunks to a full segment
			emptyChunksPadded = core.DefaultSegmentMaxChunks - numChunks
		} else if lastSegmentChunks := numChunksFlowPadded % core.DefaultSegmentMaxChunks; numChunks < lastSegmentChunks {
			// pad for the last segment with flow padded empty chunks
			emptyChunksPadded = lastSegmentChunks - numChunks
		}
	}

	segmentRootHash := core.SegmentRoot(segment.Data, emptyChunksPadded)

	if err := segment.Proof.ValidateHash(root, segmentRootHash, segmentIndex, numSegmentsFlowPadded); err != nil {
		return nil, errors.WithMessage(err, "Failed to validate proof")
	}

	return segment.Data, nil
}
