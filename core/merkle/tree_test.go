package merkle

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func createChunkData(i int) []byte {
	return []byte(fmt.Sprintf("chunk data - %v", i))
}

func createTreeByChunks(chunks int) *Tree {
	var builder TreeBuilder

	for i := 0; i < chunks; i++ {
		chunk := createChunkData(i)
		builder.Append(chunk)
	}

	return builder.Build()
}

func TestTreeRoot(t *testing.T) {
	assert.Equal(t, "0x2dea03c693750777940bcd0cc3f5d93543c075fa3b9a07b9fd86ec8fbaf6a8b2", createTreeByChunks(5).Root().Hex())
	assert.Equal(t, "0x318c92000aefba6ebf570a8a6daa57aa643f04350ffbe583999ddd9e24ceb147", createTreeByChunks(6).Root().Hex())
	assert.Equal(t, "0xca80116fb7fb8d6ef4a47e322f22e94ae8beb03e6fcbf8ab59c4d6f54fe42c4d", createTreeByChunks(7).Root().Hex())
}

func TestTreeProof(t *testing.T) {
	for numChunks := 1; numChunks <= 32; numChunks++ {
		tree := createTreeByChunks(numChunks)

		for i := 0; i < numChunks; i++ {
			proof := tree.ProofAt(i)
			assert.NoError(t, proof.Validate(tree.Root(), createChunkData(i), uint64(i), uint64(numChunks)))
		}
	}
}

// chunksPerSegment: 2^n
func calculateRootBySegments(chunks, chunksPerSegment int) common.Hash {
	var fileBuilder TreeBuilder

	for i := 0; i < chunks; i += chunksPerSegment {
		var segBuilder TreeBuilder

		for j := 0; j < chunksPerSegment; j++ {
			index := i + j
			if index >= chunks {
				break
			}

			segBuilder.Append(createChunkData(index))
		}

		segRoot := segBuilder.Build().Root()

		fileBuilder.AppendHash(segRoot)
	}

	return fileBuilder.Build().Root()
}

// Number of chunks in segment will not impact the merkle root
func TestRootBySegment(t *testing.T) {
	for chunks := 1; chunks <= 256; chunks++ {
		root1 := createTreeByChunks(chunks).Root()   // no segment
		root2 := calculateRootBySegments(chunks, 4)  // segment with 4 chunks
		root3 := calculateRootBySegments(chunks, 16) // segment with 16 chunks

		assert.Equal(t, root1, root2)
		assert.Equal(t, root2, root3)
	}
}
