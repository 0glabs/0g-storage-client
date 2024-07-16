package node

import (
	"github.com/openweb3/go-rpc-provider/interfaces"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
)

type rpcClient struct {
	url      string
	provider interfaces.Provider
}

func newRpcClient(url string, option ...providers.Option) (*rpcClient, error) {
	var opt providers.Option
	if len(option) > 0 {
		opt = option[0]
	}

	provider, err := providers.NewProviderWithOption(url, opt)
	if err != nil {
		return nil, err
	}

	return &rpcClient{url, provider}, nil
}

// URL Get the RPC server URL the client connected to.
func (c *rpcClient) URL() string {
	return c.url
}

// Close close the underlying RPC client.
func (c *rpcClient) Close() {
	c.provider.Close()
}
