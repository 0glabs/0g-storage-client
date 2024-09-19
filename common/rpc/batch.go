package rpc

import (
	"context"

	"github.com/openweb3/go-rpc-provider"
	"github.com/openweb3/go-rpc-provider/interfaces"
)

type Request struct {
	Method string
	Args   []any
}

type Response[T any] struct {
	Data  T
	Error error
}

// BatchCall is a generic method to call RPC in batch.
func BatchCall[T any](provider interfaces.Provider, requests ...Request) ([]Response[T], error) {
	return BatchCallContext[T](provider, context.Background(), requests...)
}

// BatchCallContext is a generic method to call RPC with context in batch.
func BatchCallContext[T any](provider interfaces.Provider, ctx context.Context, requests ...Request) ([]Response[T], error) {
	batch := make([]rpc.BatchElem, 0, len(requests))
	responses := make([]Response[T], len(requests))

	for i, v := range requests {
		batch = append(batch, rpc.BatchElem{
			Method: v.Method,
			Args:   v.Args,
			Result: &responses[i].Data,
		})
	}

	if err := provider.BatchCallContext(ctx, batch); err != nil {
		return nil, err
	}

	for i, v := range batch {
		if v.Error != nil {
			responses[i].Error = v.Error
		}
	}

	return responses, nil
}
