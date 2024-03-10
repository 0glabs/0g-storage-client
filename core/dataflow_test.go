package core

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFileAndInMemoryData(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	size := DefaultSegmentSize*10 + 10

	data := make([]byte, size)
	n, err := r.Read(data)
	assert.NoError(t, err)
	assert.Equal(t, n, len(data))

	tmpFile, err := os.CreateTemp("", "0g-storage-client-*")
	assert.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write(data)
	assert.NoError(t, err)

	file, err := Open(tmpFile.Name())
	assert.NoError(t, err)

	fileTree, err := MerkleTree(file)
	assert.NoError(t, err)

	inMem, _ := NewDataInMemory(data)
	inMemTree, err := MerkleTree(inMem)
	assert.NoError(t, err)

	assert.Equal(t, fileTree.Root(), inMemTree.Root())
}
