package gateway

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/0glabs/0g-storage-client/node"
	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

var clients []*node.ZgsClient
var maxDownloadFileSize uint64

func downloadFile(c *gin.Context) {
	var input struct {
		Name  string `form:"name" json:"name"`
		Root  string `form:"root" json:"root"`
		TxSeq uint64 `form:"txSeq" json:"txSeq"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Failed to bind input parameters, %v", err.Error()))
		return
	}

	var fileInfo *node.FileInfo
	var err error

	if len(input.Root) == 0 {
		if fileInfo, err = clients[0].GetFileInfoByTxSeq(context.Background(), input.TxSeq); err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to get file info by tx seq, %v", err.Error()))
			return
		}
	} else {
		if fileInfo, err = clients[0].GetFileInfo(context.Background(), common.HexToHash(input.Root)); err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to get file info by root, %v", err.Error()))
			return
		}
	}

	if fileInfo == nil {
		c.JSON(http.StatusNotFound, "File not found")
		return
	}

	if !fileInfo.Finalized {
		c.JSON(http.StatusBadRequest, "File not finalized yet")
		return
	}

	if fileInfo.Tx.Size > maxDownloadFileSize {
		errMsg := fmt.Sprintf("Requested file size too large, actual = %v, max = %v", fileInfo.Tx.Size, maxDownloadFileSize)
		c.JSON(http.StatusBadRequest, errMsg)
		return
	}

	downloader, err := transfer.NewDownloader(clients)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to create downloader, %v", err.Error()))
		return
	}

	root := fileInfo.Tx.DataMerkleRoot.Hex()
	tmpfile := path.Join(os.TempDir(), fmt.Sprintf("zgs_indexer_download_%v", root))
	defer os.Remove(tmpfile)

	if err = downloader.Download(context.Background(), fileInfo.Tx.DataMerkleRoot.Hex(), tmpfile, true); err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to download file, %v", err.Error()))
		return
	}

	if len(input.Name) == 0 {
		c.FileAttachment(tmpfile, root)
	} else {
		c.FileAttachment(tmpfile, input.Name)
	}
}
