package rpc

import (
	"net/http"
	"net/rpc"

	"github.com/ethereum/go-ethereum/node"
	"github.com/sirupsen/logrus"
)

// MustNewHandler creates a http.Handler for the specified RPC apis.
func MustNewHandler(apis map[string]interface{}) http.Handler {
	handler := rpc.NewServer()

	for namespace, impl := range apis {
		if err := handler.RegisterName(namespace, impl); err != nil {
			logrus.WithError(err).WithField("namespace", namespace).Fatal("Failed to register rpc service")
		}
	}

	// enable cors
	return node.NewHTTPHandlerStack(handler, []string{"*"}, []string{"*"}, []byte{})
}

// Start starts a HTTP service util shutdown.
func Start(endpoint string, handler http.Handler) {
	server := http.Server{
		Addr:    endpoint,
		Handler: handler,
	}

	server.ListenAndServe()
}

// MustServe starts RPC service until shutdown.
func MustServe(endpoint string, apis map[string]interface{}) {
	rpcHandler := MustNewHandler(apis)
	Start(endpoint, rpcHandler)
}
