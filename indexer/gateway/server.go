package gateway

import (
	"net/http"

	"github.com/0glabs/0g-storage-client/common/util"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Endpoint string

	Nodes               []string // storage nodes to download files
	MaxDownloadFileSize uint64

	RPCHandler http.Handler // enable to provide both RPC and REST API service
}

func MustServeWithRPC(config Config) {
	if len(config.Nodes) == 0 {
		logrus.Fatal("Nodes not specified to start HTTP server")
	}

	// init global variables
	clients = node.MustNewZgsClients(config.Nodes)
	maxDownloadFileSize = config.MaxDownloadFileSize

	// init router
	router := newRouter()
	if config.RPCHandler != nil {
		router.POST("/", gin.WrapH(config.RPCHandler))
	}

	util.MustServe(config.Endpoint, router)
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

	return router
}

func middlewareCors() gin.HandlerFunc {
	conf := cors.DefaultConfig()
	conf.AllowMethods = append(conf.AllowMethods, "OPTIONS")
	conf.AllowHeaders = append(conf.AllowHeaders, "*")
	conf.AllowAllOrigins = true

	return cors.New(conf)
}
