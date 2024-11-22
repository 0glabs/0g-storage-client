package gateway

import (
	"net/http"

	"github.com/0glabs/0g-storage-client/common/rpc"
	"github.com/0glabs/0g-storage-client/indexer"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Endpoint            string       // http endpoint
	RPCHandler          http.Handler // enable to provide both RPC and REST API service
	MaxDownloadFileSize uint64       // max download file size
}

func MustServeWithRPC(nodeManager *indexer.NodeManager, locationCache *indexer.FileLocationCache, config Config) {
	controller := NewRestController(nodeManager, locationCache, config.MaxDownloadFileSize)

	// init router
	router := newRouter(controller)
	if config.RPCHandler != nil {
		router.POST("/", gin.WrapH(config.RPCHandler))
	}

	rpc.Start(config.Endpoint, router)
}

func newRouter(controller *RestController) *gin.Engine {
	router := gin.New()

	// middlewares
	router.Use(gin.Recovery())
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		router.Use(gin.Logger())
	}
	router.Use(middlewareCors())

	// handlers
	router.GET("/file", controller.downloadFile)
	router.GET("/file/:cid/*filePath", controller.downloadFileInFolder)
	router.GET("/file/info/:cid", controller.getFileStatus)
	router.GET("/node/status", controller.getNodeStatus)
	router.POST("/file/segment", controller.uploadSegment)

	return router
}

func middlewareCors() gin.HandlerFunc {
	conf := cors.DefaultConfig()
	conf.AllowMethods = append(conf.AllowMethods, "OPTIONS")
	conf.AllowHeaders = append(conf.AllowHeaders, "*")
	conf.AllowAllOrigins = true

	return cors.New(conf)
}
