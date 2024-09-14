package node

import (
	"github.com/openweb3/go-rpc-provider/interfaces"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
)

type rpcClient struct {
	interfaces.Provider
	url string
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

	return &rpcClient{provider, url}, nil
}

// URL Get the RPC server URL the client connected to.
func (c *rpcClient) URL() string {
	return c.url
}

func (c *rpcClient) wrapError(e error, method string) error {
	if e == nil {
		return nil
	}
	return &RPCError{
		Message: e.Error(),
		Method:  method,
		URL:     c.URL(),
	}
}
