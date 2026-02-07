package grpc

import (
	"context"
	"fmt"

	jsonrpc "github.com/sui-sdks/go-sdks/sui/jsonrpc"
)

type CoreClientOptions struct {
	Client *Client
}

type CoreClient struct {
	client *Client
}

func NewCoreClient(opts CoreClientOptions) *CoreClient {
	return &CoreClient{client: opts.Client}
}

func (c *CoreClient) Call(ctx context.Context, method string, params []any, out any) error {
	return c.client.rpc.Call(ctx, method, params, out)
}

func (c *CoreClient) GetObjects(ctx context.Context, objectIDs []string, include map[string]any) (map[string]any, error) {
	results, err := c.client.rpc.MultiGetObjects(ctx, objectIDs, include)
	if err != nil {
		return nil, err
	}
	return map[string]any{"objects": results}, nil
}

func (c *CoreClient) GetObject(ctx context.Context, objectID string, include map[string]any) (map[string]any, error) {
	obj, err := c.client.rpc.GetObject(ctx, objectID, include)
	if err != nil {
		return nil, err
	}
	return map[string]any{"object": obj}, nil
}

func (c *CoreClient) ListCoins(ctx context.Context, owner, coinType string, cursor any, limit *int) (map[string]any, error) {
	return c.client.rpc.GetCoins(ctx, owner, coinType, cursor, limit)
}

func (c *CoreClient) ListOwnedObjects(ctx context.Context, owner string, filter map[string]any, cursor any, limit *int) (map[string]any, error) {
	var out map[string]any
	err := c.client.rpc.Call(ctx, "suix_getOwnedObjects", []any{owner, filter, cursor, intOrNil(limit)}, &out)
	return out, err
}

func (c *CoreClient) GetBalance(ctx context.Context, owner, coinType string) (map[string]any, error) {
	return c.client.rpc.GetBalance(ctx, owner, coinType)
}

func (c *CoreClient) ListBalances(ctx context.Context, owner string) ([]map[string]any, error) {
	return c.client.rpc.GetAllBalances(ctx, owner)
}

func (c *CoreClient) GetCoinMetadata(ctx context.Context, coinType string) (map[string]any, error) {
	return c.client.rpc.GetCoinMetadata(ctx, coinType)
}

func (c *CoreClient) GetTransaction(ctx context.Context, digest string, include map[string]any) (map[string]any, error) {
	return c.client.rpc.GetTransactionBlock(ctx, digest, include)
}

func (c *CoreClient) ExecuteTransaction(ctx context.Context, txBytesBase64 string, signatures []string, include map[string]any, requestType string) (map[string]any, error) {
	return c.client.rpc.ExecuteTransactionBlock(ctx, txBytesBase64, signatures, include, requestType)
}

func (c *CoreClient) SimulateTransaction(ctx context.Context, txBytesBase64 string) (map[string]any, error) {
	var out map[string]any
	err := c.client.rpc.Call(ctx, "sui_dryRunTransactionBlock", []any{txBytesBase64}, &out)
	return out, err
}

func (c *CoreClient) GetReferenceGasPrice(ctx context.Context) (string, error) {
	return c.client.rpc.GetReferenceGasPrice(ctx)
}

func (c *CoreClient) GetCurrentSystemState(ctx context.Context) (map[string]any, error) {
	var out map[string]any
	err := c.client.rpc.Call(ctx, "suix_getLatestSuiSystemState", []any{}, &out)
	return out, err
}

func (c *CoreClient) GetChainIdentifier(ctx context.Context) (string, error) {
	var out string
	err := c.client.rpc.Call(ctx, "sui_getChainIdentifier", []any{}, &out)
	return out, err
}

func (c *CoreClient) ListDynamicFields(ctx context.Context, parentObjectID string, cursor any, limit *int) (map[string]any, error) {
	var out map[string]any
	err := c.client.rpc.Call(ctx, "suix_getDynamicFields", []any{parentObjectID, cursor, intOrNil(limit)}, &out)
	return out, err
}

func (c *CoreClient) GetDynamicFieldObject(ctx context.Context, parentObjectID string, name any) (map[string]any, error) {
	var out map[string]any
	err := c.client.rpc.Call(ctx, "suix_getDynamicFieldObject", []any{parentObjectID, name}, &out)
	return out, err
}

func (c *CoreClient) VerifyZkLoginSignature(ctx context.Context, signature string, bytes string, intentScope string, author string) (map[string]any, error) {
	var out map[string]any
	err := c.client.rpc.Call(ctx, "sui_verifyZkLoginSignature", []any{signature, bytes, intentScope, author}, &out)
	return out, err
}

func (c *CoreClient) GetMoveFunction(ctx context.Context, packageID, module, function string) (map[string]any, error) {
	var out map[string]any
	err := c.client.rpc.Call(ctx, "sui_getMoveFunctionArgTypes", []any{packageID, module, function}, &out)
	if err != nil {
		return nil, err
	}
	return map[string]any{"function": out}, nil
}

func (c *CoreClient) DefaultNameServiceName(ctx context.Context, address string) (string, error) {
	var out string
	err := c.client.rpc.Call(ctx, "suix_resolveNameServiceNames", []any{address, nil, 1}, &out)
	if err != nil {
		return "", err
	}
	return out, nil
}

func (c *CoreClient) ResolveTransactionPlugin() any {
	return nil
}

func intOrNil(v *int) any {
	if v == nil {
		return nil
	}
	return *v
}

func mapErr(msg string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", msg, err)
}

var _ = jsonrpc.Client{}
