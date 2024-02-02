package core

import (
	"errors"
	"runtime"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
	"github.com/zero-gravity-labs/zerog-storage-client/common/parallel"
	"github.com/zero-gravity-labs/zerog-storage-client/core/merkle"
)

const (
	// DefaultChunkSize represents the default chunk size in bytes.
	DefaultChunkSize = 256

	// DefaultSegmentMaxChunks represents the default maximum number of chunks within a segment.
	DefaultSegmentMaxChunks = 1024

	// DefaultSegmentSize represents the default segment size in bytes.
	DefaultSegmentSize = DefaultChunkSize * DefaultSegmentMaxChunks
)

var (
	EmptyChunk     = make([]byte, DefaultChunkSize)
	EmptyChunkHash = crypto.Keccak256Hash(EmptyChunk)
)

type IterableData interface {
	NumChunks() uint64
	NumSegments() uint64
	Size() int64
	PaddedSize() uint64
	Iterate(offset int64, batch int64, flowPadding bool) Iterator
	Read(buf []byte, offset int64) (int, error)
}

func MerkleTree(data IterableData) (*merkle.Tree, error) {
	stageTimer := time.Now()
	var builder merkle.TreeBuilder
	initializer := &TreeBuilderInitializer{
		data:    data,
		builder: &builder,
	}

	err := parallel.Serial(initializer, NumSegmentsPadded(data), runtime.GOMAXPROCS(0), 0)
	if err != nil {
		return nil, err
	}

	logrus.WithField("duration", time.Since(stageTimer)).Info("create segment root took")

	stageTimer = time.Now()
	tree := builder.Build()
	logrus.WithField("duration", time.Since(stageTimer)).Info("build merkle tree took")

	return tree, nil
}

func NumSplits(total int64, unit int) uint64 {
	return uint64((total-1)/int64(unit) + 1)
}

func NumSegmentsPadded(data IterableData) int {
	return int((data.PaddedSize()-1)/DefaultSegmentSize + 1)
}

func SegmentRoot(chunks []byte, emptyChunksPadded ...uint64) common.Hash {
	var builder merkle.TreeBuilder

	// append chunks
	for offset, dataLen := 0, len(chunks); offset < dataLen; offset += DefaultChunkSize {
		chunk := chunks[offset : offset+DefaultChunkSize]
		builder.Append(chunk)
	}

	// append empty chunks
	if len(emptyChunksPadded) > 0 && emptyChunksPadded[0] > 0 {
		for i := uint64(0); i < emptyChunksPadded[0]; i++ {
			builder.AppendHash(EmptyChunkHash)
		}
	}

	if tree := builder.Build(); tree != nil {
		return tree.Root()
	}

	return common.Hash{}
}

func paddingZeros(buf []byte, startOffset int, length int) {
	for i := 0; i < length; i++ {
		buf[startOffset+i] = 0
	}
}

func ReadAt(data IterableData, readSize int, offset int64, paddedSize uint64) ([]byte, error) {
	// Reject invalid offset
	if offset < 0 || uint64(offset) >= paddedSize {
		return nil, errors.New("invalid offset")
	}

	var expectedBufSize int
	maxAvailableLength := paddedSize - uint64(offset)
	if maxAvailableLength >= uint64(readSize) {
		expectedBufSize = readSize
	} else {
		expectedBufSize = int(maxAvailableLength)
	}

	if offset >= data.Size() {
		return make([]byte, expectedBufSize), nil
	}

	buf := make([]byte, expectedBufSize)

	_, err := data.Read(buf, offset)
	if err != nil {
		return nil, err
	}

	return buf, nil
}
