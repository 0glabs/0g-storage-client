package rpc

import (
	"context"

	"github.com/openweb3/go-rpc-provider/interfaces"
	providers "github.com/openweb3/go-rpc-provider/provider_wrapper"
)

// Client is a base class of any RPC client.
type Client struct {
	*providers.MiddlewarableProvider
	url string
}

// NewClient creates a new client instance.
func NewClient(url string, option ...providers.Option) (*Client, error) {
	var opt providers.Option
	if len(option) > 0 {
		opt = option[0]
	}

	provider, err := providers.NewProviderWithOption(url, opt)
	if err != nil {
		return nil, err
	}

	return &Client{providers.NewMiddlewarableProvider(provider), url}, nil
}

// URL Get the RPC server URL the client connected to.
func (c *Client) URL() string {
	return c.url
}

// Call is a generic method to call RPC.
func Call[T any](provider interfaces.Provider, method string, args ...any) (result T, err error) {
	return CallContext[T](provider, context.Background(), method, args...)
}

// CallContext is a generic method to call RPC with context.
func CallContext[T any](provider interfaces.Provider, ctx context.Context, method string, args ...any) (result T, err error) {
	err = provider.CallContext(ctx, &result, method, args...)
	return
}
