package gateway

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/0glabs/0g-storage-client/node"
	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/0glabs/0g-storage-client/transfer/dir"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

var clients []*node.ZgsClient
var maxDownloadFileSize uint64

// downloadFile handles file downloads by root hash or transaction sequence.
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

	fileInfo, err := getFileInfo(c, input.Root, input.TxSeq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve file info: %v", err))
		return
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

	if err := downloadAndServeFile(c, downloader, fileInfo.Tx.DataMerkleRoot.Hex(), input.Name); err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to download file: %v", err))
	}
}

// downloadFileInFolder handles file downloads from a directory structure.
func downloadFileInFolder(c *gin.Context) {
	cid := c.Param("cid")
	filePath := filepath.Clean(c.Param("filePath"))

	var root string
	var txSeq uint64

	if v, err := strconv.ParseUint(cid, 10, 64); err == nil { // TxnSeq is used as cid
		txSeq = v
	} else {
		root = cid
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

	if !fileInfo.Finalized {
		c.JSON(http.StatusBadRequest, "File not finalized yet")
		return
	}

	downloader, err := transfer.NewDownloader(clients)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to create downloader, %v", err.Error()))
		return
	}

	root = fileInfo.Tx.DataMerkleRoot.Hex()

	ftree, err := transfer.BuildFileTree(c, downloader, root, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to build file tree, %v", err.Error()))
		return
	}

	fnode, err := ftree.Locate(filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("File path not found: %v", err))
		return
	}

	switch fnode.Type {
	case dir.FileTypeDirectory:
		// Show the list of files in the directory.
		serveDirectoryListing(c, fnode)
	case dir.FileTypeSymbolic:
		// If the file type is symbolic (a symlink), return the metadata of the symbolic link itself
		// (i.e., information about the symlink, not the target it points to).
		// This prevents the server from following the symbolic link and returning the target file's content.
		c.JSON(http.StatusOK, fnode)
	case dir.FileTypeFile:
		if fnode.Size > int64(maxDownloadFileSize) {
			errMsg := fmt.Sprintf("Requested file size too large, actual = %v, max = %v", fnode.Size, maxDownloadFileSize)
			c.JSON(http.StatusBadRequest, errMsg)
			return
		}

		if err := downloadAndServeFile(c, downloader, fnode.Root, fnode.Name); err != nil {
			c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to download file: %v", err))
		}
	default:
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Unsupported file type: %v", fnode.Type))
		return
	}
}

// getFileInfo retrieves file info based on root or transaction sequence.
func getFileInfo(ctx context.Context, root string, txSeq uint64) (info *node.FileInfo, err error) {
	if len(clients) == 0 {
		return nil, errors.New("no clients available")
	}

	for _, client := range clients {
		if root != "" {
			info, err = client.GetFileInfo(ctx, common.HexToHash(root))
		} else {
			info, err = client.GetFileInfoByTxSeq(ctx, txSeq)
		}

		if err != nil {
			return nil, err
		}

		if info != nil {
			return info, nil
		}
	}

	return nil, nil
}

// downloadAndServeFile downloads the file and serves it as an attachment.
func downloadAndServeFile(c *gin.Context, downloader *transfer.Downloader, root, filename string) error {
	tmpfile := filepath.Join(os.TempDir(), fmt.Sprintf("zgs_indexer_download_%v", root))
	defer os.Remove(tmpfile)

	if err := downloader.Download(c, root, tmpfile, true); err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	if len(filename) == 0 {
		filename = root
	}

	c.FileAttachment(tmpfile, filename)
	return nil
}

// serveDirectoryListing serves the list of files in a directory.
func serveDirectoryListing(c *gin.Context, dirNode *dir.FsNode) {
	type DirListing struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Size uint64 `json:"size,omitempty"`
	}

	var dirListing []DirListing
	for _, entry := range dirNode.Entries {
		dirListing = append(dirListing, DirListing{
			Name: entry.Name,
			Type: string(entry.Type),
			Size: uint64(entry.Size),
		})
	}
	c.JSON(http.StatusOK, dirListing)
}
