package grpc

import (
	"context"

	jsonrpc "github.com/sui-sdks/go-sdks/sui/jsonrpc"
)

// Transport abstracts the invocation layer for grpc package clients.
type Transport interface {
	Call(ctx context.Context, method string, params []any, out any) error
	Close() error
}

type jsonRPCTransport struct {
	rpc *jsonrpc.Client
}

func NewJSONRPCTransport(rpc *jsonrpc.Client) Transport {
	return &jsonRPCTransport{rpc: rpc}
}

func (t *jsonRPCTransport) Call(ctx context.Context, method string, params []any, out any) error {
	return t.rpc.Call(ctx, method, params, out)
}

func (t *jsonRPCTransport) Close() error { return nil }
