package gateway

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/0glabs/0g-storage-client/transfer"
	"github.com/0glabs/0g-storage-client/transfer/dir"
	"github.com/gin-gonic/gin"
)

// downloadFile handles file downloads by root hash or transaction sequence.
func (ctrl *RestController) downloadFile(c *gin.Context) {
	var input struct {
		Cid
		Name string `form:"name" json:"name"`
	}

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("Failed to bind input parameters, %v", err.Error()))
		return
	}

	if input.TxSeq == nil && len(input.Root) == 0 {
		c.JSON(http.StatusBadRequest, "Either 'root' or 'txSeq' must be provided")
		return
	}

	ctrl.downloadAndServeFile(c, input.Cid, input.Name)
}

// downloadFileInFolder handles file downloads from a directory structure.
func (ctrl *RestController) downloadFileInFolder(c *gin.Context) {
	cidStr := strings.TrimSpace(c.Param("cid"))
	filePath := filepath.Clean(c.Param("filePath"))
	cid := NewCid(cidStr)

	clients, err := ctrl.getAvailableStorageNodes(c, cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to get available storage nodes, %v", err.Error()))
		return
	}

	fileInfo, err := getOverallFileInfo(c, clients, cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve file info: %v", err))
		return
	}

	if fileInfo == nil {
		c.JSON(http.StatusNotFound, "File not found")
		return
	}

	if fileInfo.Pruned {
		c.JSON(http.StatusBadRequest, "File already pruned")
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

	root := fileInfo.Tx.DataMerkleRoot

	ftree, err := transfer.BuildFileTree(c, downloader, root.Hex(), true)
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
		if fnode.Size > int64(ctrl.maxDownloadFileSize) {
			errMsg := fmt.Sprintf("Requested file size too large, actual = %v, max = %v", fnode.Size, ctrl.maxDownloadFileSize)
			c.JSON(http.StatusBadRequest, errMsg)
			return
		}

		ctrl.downloadAndServeFile(c, Cid{Root: fnode.Root}, fnode.Name)
	default:
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Unsupported file type: %v", fnode.Type))
		return
	}
}

// downloadAndServeFile downloads the file and serves it as an attachment.
func (ctrl *RestController) downloadAndServeFile(c *gin.Context, cid Cid, filename string) {
	clients, err := ctrl.getAvailableStorageNodes(c, cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to get available storage nodes, %v", err.Error()))
		return
	}

	fileInfo, err := getOverallFileInfo(c, clients, cid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve file info: %v", err))
		return
	}

	if fileInfo == nil {
		c.JSON(http.StatusNotFound, "File not found")
		return
	}

	if fileInfo.Pruned {
		c.JSON(http.StatusBadRequest, "File already pruned")
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

	root := fileInfo.Tx.DataMerkleRoot.Hex()
	tmpfile := filepath.Join(os.TempDir(), fmt.Sprintf("zgs_indexer_download_%v", root))
	defer os.Remove(tmpfile)

	if err := downloader.Download(c, root, tmpfile, true); err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Failed to download file: %v", err.Error()))
		return
	}

	if len(filename) == 0 {
		filename = root
	}

	c.FileAttachment(tmpfile, filename)
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
