package gateway

import (
	"encoding/json"
	"io"
	"mime/multipart"
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
	Root  string  `form:"root" json:"root"`   // Merkle root
	TxSeq *uint64 `form:"txSeq" json:"txSeq"` // Transaction sequence
	Index uint64  `form:"index" json:"index"` // Segment index
	Proof string  `form:"proof" json:"proof"` // Merkle proof (encoded as JSON string)

	// Data can be either byte array or nil depending on request type
	Data []byte `json:"data,omitempty"` // for `application/json` request
	// OR (mutually exclusive with Data)
	File *multipart.FileHeader `form:"data,omitempty"` // for `multipart/form-data` request
}

// GetData reads segment data from the request
func (req *UploadSegmentRequest) GetData() ([]byte, error) {
	if req.Data != nil {
		return req.Data, nil
	}
	if req.File == nil {
		return nil, errors.New("either `data` or `file` is required")
	}
	data, err := req.readFileData()
	if err != nil {
		return nil, errors.WithMessage(err, "failed to load file data")
	}
	req.Data = data
	return data, nil
}

// readFileData reads data from the uploaded file
func (req *UploadSegmentRequest) readFileData() ([]byte, error) {
	file, err := req.File.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return fileBytes, nil
}

func uploadSegment(c *gin.Context) {
	var input UploadSegmentRequest

	// bind the request (supports both `application/json` and `multipart/form-data`)
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, errors.WithMessagef(err, "Failed to bind input parameters").Error())
		return
	}

	// validate root and transaction sequence
	if input.TxSeq == nil && len(input.Root) == 0 {
		c.JSON(http.StatusBadRequest, "Either 'root' or 'txSeq' must be provided")
		return
	}

	// validate segment data
	if _, err := input.GetData(); err != nil {
		c.JSON(http.StatusBadRequest, errors.WithMessagef(err, "Failed to extract data").Error())
		return
	}

	if len(input.Data) == 0 {
		c.JSON(http.StatusBadRequest, "Segment data is empty")
		return
	}

	// retrieve and validate file info
	fileInfo, err := getFileInfo(c, common.HexToHash(input.Root), input.TxSeq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.WithMessage(err, "Failed to retrieve file info").Error())
		return
	}
	if fileInfo == nil {
		c.JSON(http.StatusNotFound, "File not found")
		return
	}

	// validate merkle proof
	var proof merkle.Proof
	if err := json.Unmarshal([]byte(input.Proof), &proof); err != nil {
		c.JSON(http.StatusBadRequest, "Failed to json unmarshal merkle proof")
		return
	}
	if err := validateMerkleProof(input, proof, fileInfo); err != nil {
		c.JSON(http.StatusBadRequest, errors.WithMessagef(err, "Failed to validate merkle proof").Error())
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
		c.JSON(http.StatusInternalServerError, errors.WithMessage(err, "Failed to upload segment with proof").Error())
		return
	}

	c.JSON(http.StatusOK, "Segment upload ok")
}

// validateMerkleProof is a helper function to validate merkle proof for the upload request
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
