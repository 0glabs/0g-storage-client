package gateway

import (
	"net/http"

	"github.com/0glabs/0g-storage-client/common/api"
	"github.com/0glabs/0g-storage-client/indexer"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Endpoint            string       // http endpoint
	RPCHandler          http.Handler // enable to provide both RPC and REST API service
	MaxDownloadFileSize uint64       // max download file size
}

func MustServeWithRPC(nodeManager *indexer.NodeManager, locationCache *indexer.FileLocationCache, config Config) {
	controller := NewRestController(nodeManager, locationCache, config.MaxDownloadFileSize)

	api.Serve(config.Endpoint, func(router *gin.Engine) {
		router.GET("/file", controller.downloadFile)
		router.GET("/file/:cid/*filePath", controller.downloadFileInFolder)
		router.GET("/file/info/:cid", controller.getFileStatus)
		router.GET("/node/status", controller.getNodeStatus)
		router.POST("/file/segment", controller.uploadSegment)

		if config.RPCHandler != nil {
			router.POST("/", gin.WrapH(config.RPCHandler))
		}
	})
}
