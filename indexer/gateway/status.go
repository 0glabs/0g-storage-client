package gateway

import (
	"strings"

	"github.com/0glabs/0g-storage-client/node"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

func (ctrl *RestController) getFileStatus(c *gin.Context) (interface{}, error) {
	cidStr := strings.TrimSpace(c.Param("cid"))
	cid := NewCid(cidStr)

	fileInfo, err := ctrl.fetchFileInfo(c, cid)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to retrieve file info")
	}

	if fileInfo == nil {
		return nil, ErrFileNotFound
	}

	return fileInfo, nil
}

func (ctrl *RestController) getNodeStatus(c *gin.Context) (interface{}, error) {
	var finalStatus *node.Status
	for _, client := range ctrl.nodeManager.TrustedClients() {
		status, err := client.GetStatus(c)
		if err != nil {
			return nil, errors.WithMessagef(err, "Failed to retrieve node status")
		}

		if finalStatus == nil {
			finalStatus = &status
			continue
		}

		if finalStatus.LogSyncHeight > status.LogSyncHeight {
			finalStatus.LogSyncHeight = status.LogSyncHeight
			finalStatus.LogSyncBlock = status.LogSyncBlock
		}

		finalStatus.ConnectedPeers = max(finalStatus.ConnectedPeers, status.ConnectedPeers)
		finalStatus.NextTxSeq = min(finalStatus.NextTxSeq, status.NextTxSeq)
	}

	return finalStatus, nil
}
