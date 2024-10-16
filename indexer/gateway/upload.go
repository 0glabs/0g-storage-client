package gateway

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"

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
	Root  string          `form:"root" json:"root" binding:"required"`   // File merkle root
	Data  []byte          `form:"data" json:"data" binding:"required"`   // Segment data (base64 encoded in JSON or raw bytes in form-data)
	Index uint64          `form:"index" json:"index" binding:"required"` // Segment index
	Proof json.RawMessage `form:"proof" json:"proof" binding:"required"` // Merkle proof (encoded as JSON string)
}

// decode base64 data for JSON requests
func (req *UploadSegmentRequest) decodeBase64Data() error {
	decodedData, err := base64.StdEncoding.DecodeString(string(req.Data))
	if err != nil {
		return err
	}

	req.Data = decodedData
	return nil
}

func uploadSegment(c *gin.Context) {
	var input UploadSegmentRequest

	// bind the request (supports both form and JSON)
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.WithMessagef(err, "Failed to bind input parameters")})
		return
	}

	// validate root hash
	if len(input.Root) != 66 || common.HexToHash(input.Root) == (common.Hash{}) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid root hash"})
		return
	}

	// handle base64 decoding for JSON content type
	if isJSON(c) {
		if err := input.decodeBase64Data(); err != nil {
			err = errors.WithMessage(err, "Failed to decode base64 data for json request")
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
	}

	// validate segment data length
	if len(input.Data) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Segment data is empty"})
		return
	}

	// retrieve and validate file info
	fileInfo, err := getFileInfo(c, input.Root, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errors.WithMessage(err, "Failed to retrieve file info")})
		return
	}
	if fileInfo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// validate merkle proof
	var proof merkle.Proof
	if err := json.Unmarshal([]byte(input.Proof), &proof); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to json unmarshal merkle proof"})
		return
	}
	if err := validateMerkleProof(input, proof, fileInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.WithMessagef(err, "Failed to validate merkle proof")})
		return
	}

	// upload the segment
	segment := node.SegmentWithProof{
		Root:     common.HexToHash(input.Root),
		Data:     input.Data,
		Index:    input.Index,
		Proof:    proof,
		FileSize: fileInfo.Tx.Size,
	}
	if err := uploadSegmentWithProof(c, segment); err != nil {
		err = errors.WithMessage(err, "Failed to upload segment with proof")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Segment upload ok"})
}

// isJSON is a helper function to determine if the request is JSON
func isJSON(c *gin.Context) bool {
	return strings.Contains(c.ContentType(), "application/json")
}

// validateMerkleProof is a helper function to validate merkle proof for the input request
func validateMerkleProof(req UploadSegmentRequest, proof merkle.Proof, fileInfo *node.FileInfo) error {
	fileSize := int64(fileInfo.Tx.Size)
	merkleRoot := common.HexToHash(req.Root)
	segmentRootHash, numSegmentsFlowPadded := core.PaddedSegmentRoot(req.Index, req.Data, fileSize)
	return proof.ValidateHash(merkleRoot, segmentRootHash, req.Index, numSegmentsFlowPadded)
}

// uploadSegmentWithProof is a helper function to upload the segment with proof
func uploadSegmentWithProof(c *gin.Context, segment node.SegmentWithProof) error {
	segments := []node.SegmentWithProof{segment}
	uploadOpt := transfer.UploadOption{
		ExpectedReplica: expectedReplica,
	}

	uploader := transfer.NewSegmentProofUploader(clients)
	return uploader.UploadSegments(c, segments, uploadOpt)
}
