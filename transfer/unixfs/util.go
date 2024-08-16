package unixfs

import (
	"github.com/0glabs/0g-storage-client/core"
	"github.com/pkg/errors"
)

// GetMerkleTreeRootOfFile calculates and returns the Merkle tree root hash of a given file.
// It opens the file, generates the Merkle tree, and returns the root hash as a hexadecimal string.
func GetMerkleTreeRootOfFile(filename string) (string, error) {
	// Attempt to open the file
	file, err := core.Open(filename)
	if err != nil {
		return "", errors.WithMessagef(err, "failed to open file %s", filename)
	}
	defer file.Close()

	// Generate the Merkle tree from the file content
	tree, err := core.MerkleTree(file)
	if err != nil {
		return "", errors.WithMessagef(err, "failed to create Merkle tree for file %s", filename)
	}

	// Return the root hash of the Merkle tree as a hexadecimal string
	return tree.Root().Hex(), nil
}
