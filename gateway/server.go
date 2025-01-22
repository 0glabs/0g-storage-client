package gateway

import (
	"net/http"

	"github.com/0glabs/0g-storage-client/common/api"
	"github.com/0glabs/0g-storage-client/node"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var allClients []*node.ZgsClient

func MustServeLocal(nodes []*node.ZgsClient) {
	if len(nodes) == 0 {
		logrus.Fatal("storage nodes not configured")
	}

	allClients = nodes

	if err := api.Serve("127.0.0.1:6789", routes); err != http.ErrServerClosed {
		logrus.WithError(err).Fatal("Failed to serve API")
	}
}

func routes(router *gin.Engine) {
	localApi := router.Group("/local")
	localApi.GET("/nodes", api.Wrap(listNodes))
	localApi.GET("/file", api.Wrap(getLocalFileInfo))
	localApi.GET("/status", api.Wrap(getFileStatus))
	localApi.POST("/upload", api.Wrap(uploadLocalFile))
	localApi.POST("/download", api.Wrap(downloadFileLocal))
}
