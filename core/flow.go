package core

import (
	"math"
	"math/big"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zero-gravity-labs/zerog-storage-client/contract"
	"github.com/zero-gravity-labs/zerog-storage-client/core/merkle"
)

type Flow struct {
	data IterableData
	tags []byte
}

func NewFlow(data IterableData, tags []byte) *Flow {
	return &Flow{data, tags}
}

func (flow *Flow) CreateSubmission() (*contract.Submission, error) {
	// TODO(kevin): limit file size, e.g., 2^31
	submission := contract.Submission{
		Length: big.NewInt(flow.data.Size()),
		Tags:   flow.tags,
	}

	stageTimer := time.Now()
	var offset int64
	for _, chunks := range flow.splitNodes() {
		node, err := flow.createNode(offset, chunks)
		if err != nil {
			return nil, err
		}
		submission.Nodes = append(submission.Nodes, *node)
		offset += chunks * DefaultChunkSize
	}
	logrus.WithField("duration", time.Since(stageTimer)).Info("create submission nodes took")

	return &submission, nil
}

func nextPow2(input uint64) uint64 {
	x := input
	x -= 1
	x |= x >> 32
	x |= x >> 16
	x |= x >> 8
	x |= x >> 4
	x |= x >> 2
	x |= x >> 1
	x += 1
	return x
}

func ComputePaddedSize(chunks uint64) (uint64, uint64) {
	chunksNextPow2 := nextPow2(chunks)
	if chunksNextPow2 == chunks {
		return chunksNextPow2, chunksNextPow2
	}

	var minChunk uint64
	if chunksNextPow2 >= 16 {
		minChunk = chunksNextPow2 / 16
	} else {
		minChunk = 1
	}

	paddedChunks := ((chunks-1)/minChunk + 1) * minChunk
	return paddedChunks, chunksNextPow2
}

// e.g. 64, 32, 1 in chunks
func (flow *Flow) splitNodes() []int64 {
	var nodes []int64

	chunks := flow.data.NumChunks()
	paddedChunks, chunksNextPow2 := ComputePaddedSize(chunks)
	nextChunkSize := chunksNextPow2

	for paddedChunks > 0 {
		if paddedChunks >= nextChunkSize {
			paddedChunks -= nextChunkSize
			nodes = append(nodes, int64(nextChunkSize))
		}
		nextChunkSize /= 2
	}
	logrus.WithFields(logrus.Fields{
		"chunks":   chunks,
		"nodeSize": nodes,
	}).Debug("SplitNodes")

	return nodes
}

func (flow *Flow) createNode(offset, chunks int64) (*contract.SubmissionNode, error) {
	batch := chunks
	if chunks > DefaultSegmentMaxChunks {
		batch = DefaultSegmentMaxChunks
	}

	return flow.createSegmentNode(offset, DefaultChunkSize*batch, DefaultChunkSize*chunks)
}

func (flow *Flow) createSegmentNode(offset, batch, size int64) (*contract.SubmissionNode, error) {
	iter := flow.data.Iterate(offset, batch, true)
	var builder merkle.TreeBuilder

	for i := int64(0); i < size; {
		ok, err := iter.Next()
		if err != nil {
			return nil, err
		}

		// should always load data
		if !ok {
			logrus.WithFields(logrus.Fields{
				"offset": offset,
				"size":   size,
			}).Error("Not enough data to create submission node")
			break
		}

		segment := iter.Current()
		builder.AppendHash(SegmentRoot(segment))
		i += int64(len(segment))
	}

	numChunks := size / DefaultChunkSize
	height := int64(math.Log2(float64(numChunks)))

	return &contract.SubmissionNode{
		Root:   builder.Build().Root(),
		Height: big.NewInt(height),
	}, nil
}
