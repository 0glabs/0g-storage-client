package merkle

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Binary merkle tree node
type node struct {
	parent *node
	left   *node
	right  *node
	hash   common.Hash
}

func newNode(hash common.Hash) *node {
	return &node{
		hash: hash,
	}
}

func newLeafNode(content []byte) *node {
	return &node{
		hash: crypto.Keccak256Hash(content),
	}
}

func newInteriorNode(left, right *node) *node {
	node := &node{
		left:  left,
		right: right,
		hash:  crypto.Keccak256Hash(left.hash.Bytes(), right.hash.Bytes()),
	}

	left.parent = node
	right.parent = node

	return node
}

func (n *node) isLeftSide() bool {
	return n.parent != nil && n.parent.left == n
}
