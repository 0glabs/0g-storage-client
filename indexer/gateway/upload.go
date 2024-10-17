package gateway

import (
	"encoding/json"
	"net/http"

	"github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/core/merkle"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

var expectedReplica uint

type UploadSegmentRequest struct {
	Root  string          `form:"root" json:"root"`                      // Merkle root
	TxSeq *uint64         `form:"txSeq" json:"txSeq"`                    // Transaction sequence
	Data  []byte          `form:"data" json:"data" binding:"required"`   // Segment data
	Index uint64          `form:"index" json:"index" binding:"required"` // Segment index
	Proof json.RawMessage `form:"proof" json:"proof" binding:"required"` // Merkle proof (encoded as JSON string)
}

func uploadSegment(c *gin.Context) {
	var input UploadSegmentRequest

	// bind the request (supports both form and JSON)
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.WithMessagef(err, "Failed to bind input parameters").Error(),
		})
		return
	}

	if input.TxSeq == nil && len(input.Root) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Either 'root' or 'txSeq' must be provided"})
		return
	}

	// validate segment data length
	if len(input.Data) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Segment data is empty"})
		return
	}

	// retrieve and validate file info
	fileInfo, err := getFileInfo(c, common.HexToHash(input.Root), input.TxSeq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.WithMessage(err, "Failed to retrieve file info").Error(),
		})
		return
	}
	if fileInfo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// validate merkle proof
	var proof merkle.Proof
	if err := json.Unmarshal(input.Proof, &proof); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to json unmarshal merkle proof"})
		return
	}
	if err := validateMerkleProof(input, proof, fileInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errors.WithMessagef(err, "Failed to validate merkle proof").Error(),
		})
		return
	}

	// upload the segment
	segment := node.SegmentWithProof{
		Root:     fileInfo.Tx.DataMerkleRoot,
		Data:     input.Data,
		Index:    input.Index,
		Proof:    proof,
		FileSize: fileInfo.Tx.Size,
	}
	if err := uploadSegmentWithProof(c, segment, fileInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errors.WithMessage(err, "Failed to upload segment with proof").Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Segment upload ok"})
}

// validateMerkleProof is a helper function to validate merkle proof for the input request
func validateMerkleProof(req UploadSegmentRequest, proof merkle.Proof, fileInfo *node.FileInfo) error {
	fileSize := int64(fileInfo.Tx.Size)
	merkleRoot := fileInfo.Tx.DataMerkleRoot
	segmentRootHash, numSegmentsFlowPadded := core.PaddedSegmentRoot(req.Index, req.Data, fileSize)
	return proof.ValidateHash(merkleRoot, segmentRootHash, req.Index, numSegmentsFlowPadded)
}

// uploadSegmentWithProof is a helper function to upload the segment with proof
func uploadSegmentWithProof(c *gin.Context, segment node.SegmentWithProof, fileInfo *node.FileInfo) error {
	opt := transfer.UploadOption{
		ExpectedReplica: expectedReplica,
	}
	fileSegements := transfer.FileSegmentsWithProof{
		Segments: []node.SegmentWithProof{segment},
		FileInfo: fileInfo,
	}

	return transfer.NewFileSegementUploader(clients).Upload(c, fileSegements, opt)
}
