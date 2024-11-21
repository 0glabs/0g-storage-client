package gateway

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/0glabs/0g-storage-client/node"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

func getFileStatus(c *gin.Context) {
	cid := c.Param("cid")

	var root common.Hash
	var txSeq *uint64

	if v, err := strconv.ParseUint(cid, 10, 64); err == nil { // TxnSeq is used as cid
		txSeq = &v
	} else {
		root = common.HexToHash(cid)
	}

	if txSeq == nil && (root == common.Hash{}) {
		c.JSON(http.StatusBadRequest, "Either 'root' or 'txSeq' must be provided")
		return
	}

	fileInfo, err := getFileInfo(c, root, txSeq)
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

func getNodeStatus(c *gin.Context) {
	if len(clients) == 0 {
		c.JSON(http.StatusInternalServerError, "no clients available")
		return
	}

	var finalStatus *node.Status
	for _, client := range clients {
		status, err := client.GetStatus(context.Background())
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
