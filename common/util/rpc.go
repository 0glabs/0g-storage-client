package util

import (
	"net"
	"net/http"

	"github.com/openweb3/go-rpc-provider"
	"github.com/sirupsen/logrus"
)

// MustServeRPC starts RPC service until shutdown.
func MustServeRPC(endpoint string, apis map[string]interface{}) {
	handler := rpc.NewServer()

	for namespace, impl := range apis {
		if err := handler.RegisterName(namespace, impl); err != nil {
			logrus.WithError(err).WithField("namespace", namespace).Fatal("Failed to register rpc service")
		}
	}

	httpServer := http.Server{
		// "github.com/ethereum/go-ethereum/node"
		// Handler: node.NewHTTPHandlerStack(handler, []string{"*"}, []string{"*"}),
		Handler: handler,
	}

	listener, err := net.Listen("tcp", endpoint)
	if err != nil {
		logrus.WithError(err).WithField("endpoint", endpoint).Fatal("Failed to listen to endpoint")
	}

	err = httpServer.Serve(listener)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to serve http listener")
	}
}
