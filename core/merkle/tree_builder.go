package merkle

import (
	"container/list"

	"github.com/ethereum/go-ethereum/common"
)

// TreeBuilder is used to build complete binary merkle tree.
type TreeBuilder struct {
	leafNodes []*node
}

func (builder *TreeBuilder) Append(content []byte) {
	node := newLeafNode(content)
	builder.leafNodes = append(builder.leafNodes, node)
}

func (builder *TreeBuilder) AppendHash(hash common.Hash) {
	node := newNode(hash)
	builder.leafNodes = append(builder.leafNodes, node)
}

func (builder *TreeBuilder) Build() *Tree {
	numLeafNodes := len(builder.leafNodes)
	if numLeafNodes == 0 {
		return nil
	}

	queue := list.New()

	for i := 0; i < numLeafNodes; i += 2 {
		// last single leaf node
		if i == numLeafNodes-1 {
			queue.PushBack(builder.leafNodes[i])
			continue
		}

		left, right := builder.leafNodes[i], builder.leafNodes[i+1]

		node := newInteriorNode(left, right)
		queue.PushBack(node)
	}

	for {
		numNodes := queue.Len()
		if numNodes <= 1 {
			break
		}

		for i := 0; i < numNodes/2; i++ {
			left := queue.Remove(queue.Front()).(*node)
			right := queue.Remove(queue.Front()).(*node)

			node := newInteriorNode(left, right)
			queue.PushBack(node)
		}

		// last single node
		if numNodes%2 > 0 {
			queue.MoveToBack(queue.Front())
		}
	}

	return &Tree{
		root:      queue.Front().Value.(*node),
		leafNodes: builder.leafNodes,
	}
}
