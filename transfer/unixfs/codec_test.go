package unixfs_test

import (
	"reflect"
	"testing"

	"github.com/0glabs/0g-storage-client/transfer/unixfs"
	"github.com/stretchr/testify/assert"
)

var (
	// Example usage
	exampleFsNode = &unixfs.FsNode{
		Name: "root",
		Type: unixfs.Directory,
		Entries: []*unixfs.FsNode{
			{
				Name: "file1.txt",
				Type: unixfs.File,
				Size: 1234,
			},
			{
				Name: "file2.txt",
				Type: unixfs.File,
				Size: 5678,
			},
			{
				Name: "subdir",
				Type: unixfs.Directory,
				Entries: []*unixfs.FsNode{
					{
						Name: "file3.txt",
						Type: unixfs.File,
						Size: 91011,
					},
				},
			},
		},
	}
)

func TestCodec(t *testing.T) {

	// Encode the FsNode to bytes
	encoded, err := unixfs.EncodeFsNode(exampleFsNode)
	assert.NoError(t, err)

	// Decode the bytes back to an FsNode
	decodedNode, err := unixfs.DecodeFsNode(encoded)
	assert.NoError(t, err)

	assert.True(t, reflect.DeepEqual(exampleFsNode, decodedNode))
}

// BenchmarkEncodeFsNode benchmarks the performance of the EncodeFsNode function.
func BenchmarkEncodeFsNode(b *testing.B) {
	// Run the benchmark loop
	for i := 0; i < b.N; i++ {
		_, err := unixfs.EncodeFsNode(exampleFsNode)
		if err != nil {
			b.Fatalf("Error encoding FsNode: %v", err)
		}
	}
}

// BenchmarkDecodeFsNode benchmarks the performance of the DecodeFsNode function.
func BenchmarkDecodeFsNode(b *testing.B) {
	// Pre-encode the FsNode to use as input for decoding benchmark
	encodedData, err := unixfs.EncodeFsNode(exampleFsNode)
	if err != nil {
		b.Fatalf("Error encoding FsNode for decoding benchmark: %v", err)
	}

	// Run the benchmark loop
	for i := 0; i < b.N; i++ {
		_, err := unixfs.DecodeFsNode(encodedData)
		if err != nil {
			b.Fatalf("Error decoding FsNode: %v", err)
		}
	}
}
