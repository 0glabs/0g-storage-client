package util

import (
	"net/http"

	"github.com/ethereum/go-ethereum/node"
	"github.com/openweb3/go-rpc-provider"
	"github.com/sirupsen/logrus"
)

// MustNewRPCHandler creates a http.Handler for the specified RPC apis.
func MustNewRPCHandler(apis map[string]interface{}) http.Handler {
	handler := rpc.NewServer()

	for namespace, impl := range apis {
		if err := handler.RegisterName(namespace, impl); err != nil {
			logrus.WithError(err).WithField("namespace", namespace).Fatal("Failed to register rpc service")
		}
	}

	// enable cors
	return node.NewHTTPHandlerStack(handler, []string{"*"}, []string{"*"}, []byte{})
}

// MustServe starts a HTTP service util shutdown.
func MustServe(endpoint string, handler http.Handler) {
	server := http.Server{
		Addr:    endpoint,
		Handler: handler,
	}

	server.ListenAndServe()
}

// MustServeRPC starts RPC service until shutdown.
func MustServeRPC(endpoint string, apis map[string]interface{}) {
	rpcHandler := MustNewRPCHandler(apis)
	MustServe(endpoint, rpcHandler)
}
