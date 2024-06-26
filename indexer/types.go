package indexer

import (
	"math/rand"

	"github.com/0glabs/0g-storage-client/node"
)

type ShardedNode struct {
	URL    string
	Config node.ShardConfig
}

type Interface interface {
	GetNodes() ([]ShardedNode, error)
}

func Select(nodes []ShardedNode, segmentIndex uint64, replica int) []ShardedNode {
	var matched []ShardedNode

	for _, v := range nodes {
		if v.Config.HasSegment(segmentIndex) {
			matched = append(matched, v)
		}
	}

	numMatched := len(matched)
	if numMatched == 0 {
		return nil
	}

	perm := rand.Perm(numMatched)
	result := make([]ShardedNode, numMatched)
	for i := 0; i < numMatched; i++ {
		result[i] = matched[perm[i]]
	}

	if replica < numMatched {
		result = result[:replica]
	}

	return result
}
