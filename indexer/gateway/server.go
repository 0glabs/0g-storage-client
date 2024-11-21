package gateway

import (
	"net/http"

	"github.com/0glabs/0g-storage-client/common/rpc"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Endpoint string // http endpoint

	Nodes               []string // storage nodes for file upload or download
	MaxDownloadFileSize uint64   // max download file size
	ExpectedReplica     uint     // expected upload replica number

	RPCHandler http.Handler // enable to provide both RPC and REST API service
}

func MustServeWithRPC(config Config) {
	if len(config.Nodes) == 0 {
		logrus.Fatal("Nodes not specified to start HTTP server")
	}

	// init global variables
	clients = node.MustNewZgsClients(config.Nodes)
	maxDownloadFileSize = config.MaxDownloadFileSize
	expectedReplica = config.ExpectedReplica

	// init router
	router := newRouter()
	if config.RPCHandler != nil {
		router.POST("/", gin.WrapH(config.RPCHandler))
	}

	rpc.Start(config.Endpoint, router)
}

func newRouter() *gin.Engine {
	router := gin.New()

	// middlewares
	router.Use(gin.Recovery())
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		router.Use(gin.Logger())
	}
	router.Use(middlewareCors())

	// handlers
	router.GET("/file", downloadFile)
	router.GET("/file/:cid/*filePath", downloadFileInFolder)
	router.GET("/status/:cid", getFileStatus)
	router.POST("/file/segment", uploadSegment)

	return router
}

func middlewareCors() gin.HandlerFunc {
	conf := cors.DefaultConfig()
	conf.AllowMethods = append(conf.AllowMethods, "OPTIONS")
	conf.AllowHeaders = append(conf.AllowHeaders, "*")
	conf.AllowAllOrigins = true

	return cors.New(conf)
}
