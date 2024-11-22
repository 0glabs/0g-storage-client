package gateway

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/0glabs/0g-storage-client/node"
	"github.com/gin-gonic/gin"
)

func (ctrl *RestController) getFileStatus(c *gin.Context) {
	cidStr := strings.TrimSpace(c.Param("cid"))
	cid := NewCid(cidStr)

	fileInfo, err := ctrl.fetchFileInfo(c, cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve file info: %v", err))
		return
	}

	if fileInfo == nil {
		c.JSON(http.StatusNotFound, "File not found")
		return
	}

	c.JSON(http.StatusOK, fileInfo)
}

func (ctrl *RestController) getNodeStatus(c *gin.Context) {
	var finalStatus *node.Status
	for _, client := range ctrl.nodeManager.TrustedClients() {
		status, err := client.GetStatus(c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve node status: %v", err))
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

	c.JSON(http.StatusOK, finalStatus)
}
