package gateway

import (
	"context"
	"net/http"

	"github.com/0glabs/0g-storage-client/core"
	"github.com/0glabs/0g-storage-client/core/merkle"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type UploadSegmentRequest struct {
	Root            common.Hash  `json:"root"`            // file merkle root
	TxSeq           *uint64      `json:"txSeq"`           // Transaction sequence
	Data            []byte       `json:"data"`            // segment data
	Index           uint64       `json:"index"`           // segment index
	Proof           merkle.Proof `json:"proof"`           // segment merkle proof
	ExpectedReplica uint         `json:"expectedReplica"` // expected replica count, default 1
}

func (ctrl *RestController) uploadSegment(c *gin.Context) {
	var input UploadSegmentRequest

	// bind the `application/json` request
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, errors.WithMessagef(err, "Failed to bind input parameters").Error())
		return
	}

	// validate segment data
	if len(input.Data) == 0 {
		c.JSON(http.StatusBadRequest, "Segment data is empty")
		return
	}

	cid := Cid{
		Root:  input.Root.String(),
		TxSeq: input.TxSeq,
	}

	var fileInfo *node.FileInfo
	var selectedClients []*node.ZgsClient

	// select trusted storage nodes that have already synced the submitted event
	for _, client := range ctrl.nodeManager.TrustedClients() {
		info, err := getOverallFileInfo(c, []*node.ZgsClient{client}, cid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, errors.WithMessage(err, "Failed to retrieve file info").Error())
			return
		}
		if info != nil {
			selectedClients = append(selectedClients, client)
			fileInfo = info
		}
	}
	if len(selectedClients) == 0 || fileInfo == nil {
		c.JSON(http.StatusNotFound, "File not found")
		return
	}

	// validate merkle proof
	if err := validateMerkleProof(input, fileInfo); err != nil {
		c.JSON(http.StatusBadRequest, errors.WithMessagef(err, "Failed to validate merkle proof").Error())
		return
	}

	// upload the segment
	segment := node.SegmentWithProof{
		Root:     fileInfo.Tx.DataMerkleRoot,
		Data:     input.Data,
		Index:    input.Index,
		Proof:    input.Proof,
		FileSize: fileInfo.Tx.Size,
	}
	if err := uploadSegmentWithProof(c, selectedClients, segment, fileInfo, input.ExpectedReplica); err != nil {
		c.JSON(http.StatusInternalServerError, errors.WithMessage(err, "Failed to upload segment with proof").Error())
		return
	}

	c.JSON(http.StatusOK, "Segment upload ok")
}

// validateMerkleProof is a helper function to validate merkle proof for the upload request
func validateMerkleProof(req UploadSegmentRequest, fileInfo *node.FileInfo) error {
	fileSize := int64(fileInfo.Tx.Size)
	merkleRoot := fileInfo.Tx.DataMerkleRoot
	segmentRootHash, numSegmentsFlowPadded := core.PaddedSegmentRoot(req.Index, req.Data, fileSize)
	return req.Proof.ValidateHash(merkleRoot, segmentRootHash, req.Index, numSegmentsFlowPadded)
}

// uploadSegmentWithProof is a helper function to upload the segment with proof
func uploadSegmentWithProof(
	ctx context.Context, clients []*node.ZgsClient, segment node.SegmentWithProof, fileInfo *node.FileInfo, expectedReplica uint) error {

	if expectedReplica == 0 {
		expectedReplica = 1
	}

	opt := transfer.UploadOption{
		ExpectedReplica: expectedReplica,
	}
	fileSegements := transfer.FileSegmentsWithProof{
		Segments: []node.SegmentWithProof{segment},
		FileInfo: fileInfo,
	}
	return transfer.NewFileSegementUploader(clients).Upload(ctx, fileSegements, opt)
}
