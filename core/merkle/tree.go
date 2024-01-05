package merkle

import (
	"github.com/ethereum/go-ethereum/common"
)

// Tree represents a binary merkle tree, e.g. bitcoin-like merkle tree or complete BMT.
type Tree struct {
	root      *node // always not nil
	leafNodes []*node
}

func (tree *Tree) Root() common.Hash {
	return tree.root.hash
}

func (tree *Tree) ProofAt(i int) Proof {
	if i < 0 || i >= len(tree.leafNodes) {
		panic("index out of bound")
	}

	// only single root node
	if len(tree.leafNodes) == 1 {
		return Proof{
			Lemma: []common.Hash{tree.root.hash},
			Path:  []bool{},
		}
	}

	var proof Proof

	// append the target leaf node hash
	proof.Lemma = append(proof.Lemma, tree.leafNodes[i].hash)

	for current := tree.leafNodes[i]; current != tree.root; current = current.parent {
		if current.isLeftSide() {
			proof.Lemma = append(proof.Lemma, current.parent.right.hash)
			proof.Path = append(proof.Path, true)
		} else {
			proof.Lemma = append(proof.Lemma, current.parent.left.hash)
			proof.Path = append(proof.Path, false)
		}
	}

	// append the root node hash
	proof.Lemma = append(proof.Lemma, tree.root.hash)

	return proof
}
