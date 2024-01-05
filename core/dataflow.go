package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
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
	Iterate(offset int64, batch int64, flowPadding bool) Iterator
}

func MerkleTree(data IterableData) (*merkle.Tree, error) {
	iter := data.Iterate(0, DefaultSegmentSize, true)
	var builder merkle.TreeBuilder

	for {
		ok, err := iter.Next()
		if err != nil {
			return nil, err
		}

		if !ok {
			break
		}

		segRoot := SegmentRoot(iter.Current())

		builder.AppendHash(segRoot)
	}

	return builder.Build(), nil
}

func NumSplits(total int64, unit int) uint64 {
	return uint64((total-1)/int64(unit) + 1)
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
