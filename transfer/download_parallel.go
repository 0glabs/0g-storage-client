package transfer

import (
	"fmt"
	"runtime"

	"github.com/0glabs/0g-storage-client/common/parallel"
	"github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/0glabs/0g-storage-client/transfer/download"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type SegmentDownloader struct {
	clients      []*node.Client
	shardConfigs []*node.ShardConfig
	file         *download.DownloadingFile

	withProof bool

	segmentOffset uint64
	numChunks     uint64
	numSegments   uint64
}

var _ parallel.Interface = (*SegmentDownloader)(nil)

func NewSegmentDownloader(clients []*node.Client, shardConfigs []*node.ShardConfig, file *download.DownloadingFile, withProof bool) (*SegmentDownloader, error) {
	offset := file.Metadata().Offset
	if offset%core.DefaultSegmentSize > 0 {
		return nil, errors.Errorf("Invalid data offset in downloading file %v", offset)
	}

	fileSize := file.Metadata().Size

	return &SegmentDownloader{
		clients:      clients,
		shardConfigs: shardConfigs,
		file:         file,

		withProof: withProof,

		segmentOffset: uint64(offset / core.DefaultSegmentSize),
		numChunks:     core.NumSplits(fileSize, core.DefaultChunkSize),
		numSegments:   core.NumSplits(fileSize, core.DefaultSegmentSize),
	}, nil
}

// Download downloads segments in parallel.
func (downloader *SegmentDownloader) Download() error {
	numTasks := downloader.numSegments - downloader.segmentOffset

	return parallel.Serial(downloader, int(numTasks), runtime.GOMAXPROCS(0), 0)
}

// ParallelDo implements the parallel.Interface interface.
func (downloader *SegmentDownloader) ParallelDo(routine, task int) (interface{}, error) {
	segmentIndex := downloader.segmentOffset + uint64(task)
	startIndex := segmentIndex * core.DefaultSegmentMaxChunks
	endIndex := startIndex + core.DefaultSegmentMaxChunks
	if endIndex > downloader.numChunks {
		endIndex = downloader.numChunks
	}

	root := downloader.file.Metadata().Root

	clientIndex := routine % len(downloader.shardConfigs)
	for segmentIndex%downloader.shardConfigs[clientIndex].NumShard != downloader.shardConfigs[clientIndex].ShardId {
		clientIndex = (clientIndex + 1) % len(downloader.shardConfigs)
		if clientIndex == routine%len(downloader.shardConfigs) {
			return nil, fmt.Errorf("no storage node holds segment with index %v", segmentIndex)
		}
	}

	var (
		segment []byte
		err     error
	)

	if downloader.withProof {
		segment, err = downloader.downloadWithProof(downloader.clients[clientIndex], root, startIndex, endIndex)
	} else {
		segment, err = downloader.clients[clientIndex].ZeroGStorage().DownloadSegment(root, startIndex, endIndex)
	}

	if err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"client index": clientIndex,
			"segment":      fmt.Sprintf("%v/%v", segmentIndex, downloader.numSegments),
			"chunks":       fmt.Sprintf("[%v, %v)", startIndex, endIndex),
		}).Error("Failed to download segment")
	} else if logrus.IsLevelEnabled(logrus.TraceLevel) {
		logrus.WithFields(logrus.Fields{
			"client index": clientIndex,
			"segment":      fmt.Sprintf("%v/%v", segmentIndex, downloader.numSegments),
			"chunks":       fmt.Sprintf("[%v, %v)", startIndex, endIndex),
		}).Trace("Succeeded to download segment")
	}

	// remove paddings for the last chunk
	if err == nil && segmentIndex == downloader.numSegments-1 {
		fileSize := downloader.file.Metadata().Size
		if lastChunkSize := fileSize % core.DefaultChunkSize; lastChunkSize > 0 {
			paddings := core.DefaultChunkSize - lastChunkSize
			segment = segment[0 : len(segment)-int(paddings)]
		}
	}

	return segment, err
}

// ParallelCollect implements the parallel.Interface interface.
func (downloader *SegmentDownloader) ParallelCollect(result *parallel.Result) error {
	return downloader.file.Write(result.Value.([]byte))
}

func (downloader *SegmentDownloader) downloadWithProof(client *node.Client, root common.Hash, startIndex, endIndex uint64) ([]byte, error) {
	segmentIndex := startIndex / core.DefaultSegmentMaxChunks

	segment, err := client.ZeroGStorage().DownloadSegmentWithProof(root, segmentIndex)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to download segment with proof from storage node")
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
