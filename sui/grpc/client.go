package grpc

import (
	"context"

	jsonrpc "github.com/sui-sdks/go-sdks/sui/jsonrpc"
)

type ClientOptions struct {
	Network string
	BaseURL string
	RPC     *jsonrpc.Client
}

type Client struct {
	network string
	rpc     *jsonrpc.Client
	Core    *CoreClient
}

func NewClient(opts ClientOptions) (*Client, error) {
	rpc := opts.RPC
	if rpc == nil {
		url := opts.BaseURL
		if url == "" {
			var err error
			url, err = GetGrpcFullnodeURL(opts.Network)
			if err != nil {
				return nil, err
			}
		}
		var err error
		rpc, err = jsonrpc.NewClient(jsonrpc.ClientOptions{Network: opts.Network, URL: url})
		if err != nil {
			return nil, err
		}
	}
	c := &Client{network: opts.Network, rpc: rpc}
	c.Core = NewCoreClient(CoreClientOptions{Client: c})
	return c, nil
}

func (c *Client) Network() string { return c.network }

func (c *Client) GetObjects(ctx context.Context, objectIDs []string, include map[string]any) (map[string]any, error) {
	return c.Core.GetObjects(ctx, objectIDs, include)
}
func (c *Client) GetObject(ctx context.Context, objectID string, include map[string]any) (map[string]any, error) {
	return c.Core.GetObject(ctx, objectID, include)
}
func (c *Client) ListCoins(ctx context.Context, owner, coinType string, cursor any, limit *int) (map[string]any, error) {
	return c.Core.ListCoins(ctx, owner, coinType, cursor, limit)
}
func (c *Client) ListOwnedObjects(ctx context.Context, owner string, filter map[string]any, cursor any, limit *int) (map[string]any, error) {
	return c.Core.ListOwnedObjects(ctx, owner, filter, cursor, limit)
}
func (c *Client) GetBalance(ctx context.Context, owner, coinType string) (map[string]any, error) {
	return c.Core.GetBalance(ctx, owner, coinType)
}
func (c *Client) ListBalances(ctx context.Context, owner string) ([]map[string]any, error) {
	return c.Core.ListBalances(ctx, owner)
}
func (c *Client) GetCoinMetadata(ctx context.Context, coinType string) (map[string]any, error) {
	return c.Core.GetCoinMetadata(ctx, coinType)
}
func (c *Client) GetTransaction(ctx context.Context, digest string, include map[string]any) (map[string]any, error) {
	return c.Core.GetTransaction(ctx, digest, include)
}
func (c *Client) ExecuteTransaction(ctx context.Context, txBytesBase64 string, signatures []string, include map[string]any, requestType string) (map[string]any, error) {
	return c.Core.ExecuteTransaction(ctx, txBytesBase64, signatures, include, requestType)
}
func (c *Client) SimulateTransaction(ctx context.Context, txBytesBase64 string) (map[string]any, error) {
	return c.Core.SimulateTransaction(ctx, txBytesBase64)
}
func (c *Client) GetReferenceGasPrice(ctx context.Context) (string, error) {
	return c.Core.GetReferenceGasPrice(ctx)
}
func (c *Client) ListDynamicFields(ctx context.Context, parentObjectID string, cursor any, limit *int) (map[string]any, error) {
	return c.Core.ListDynamicFields(ctx, parentObjectID, cursor, limit)
}
func (c *Client) GetDynamicField(ctx context.Context, parentObjectID string, name any) (map[string]any, error) {
	return c.Core.GetDynamicFieldObject(ctx, parentObjectID, name)
}
func (c *Client) GetMoveFunction(ctx context.Context, packageID, module, function string) (map[string]any, error) {
	return c.Core.GetMoveFunction(ctx, packageID, module, function)
}
func (c *Client) VerifyZkLoginSignature(ctx context.Context, signature string, bytes string, intentScope string, author string) (map[string]any, error) {
	return c.Core.VerifyZkLoginSignature(ctx, signature, bytes, intentScope, author)
}
func (c *Client) DefaultNameServiceName(ctx context.Context, address string) (string, error) {
	return c.Core.DefaultNameServiceName(ctx, address)
}
