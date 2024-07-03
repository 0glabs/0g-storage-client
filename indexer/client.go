package indexer

import (
	"context"

	"github.com/openweb3/go-rpc-provider/interfaces"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
)

// Requires `Client` implements the `Interface` interface.
var _ Interface = (*Client)(nil)

type Client struct {
	interfaces.Provider
}

func NewClient(url string, option ...providers.Option) (*Client, error) {
	var opt providers.Option
	if len(option) > 0 {
		opt = option[0]
	}

	provider, err := providers.NewProviderWithOption(url, opt)
	if err != nil {
		return nil, err
	}

	return &Client{
		Provider: provider,
	}, nil
}

func (c *Client) GetNodes(ctx context.Context) (nodes []ShardedNode, err error) {
	err = c.Provider.CallContext(ctx, &nodes, "indexer_getNodes")
	return
}
