package grpc

import "context"

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
	return c.client.transport.Call(ctx, method, params, out)
}

func (c *CoreClient) GetObjects(ctx context.Context, objectIDs []string, include map[string]any) (map[string]any, error) {
	var results []map[string]any
	if err := c.Call(ctx, "sui_multiGetObjects", []any{objectIDs, include}, &results); err != nil {
		return nil, err
	}
	return map[string]any{"objects": results}, nil
}

func (c *CoreClient) GetObject(ctx context.Context, objectID string, include map[string]any) (map[string]any, error) {
	var obj map[string]any
	if err := c.Call(ctx, "sui_getObject", []any{objectID, include}, &obj); err != nil {
		return nil, err
	}
	return map[string]any{"object": obj}, nil
}

func (c *CoreClient) ListCoins(ctx context.Context, owner, coinType string, cursor any, limit *int) (map[string]any, error) {
	var out map[string]any
	err := c.Call(ctx, "suix_getCoins", []any{owner, emptyToNil(coinType), cursor, intOrNil(limit)}, &out)
	return out, err
}

func (c *CoreClient) ListOwnedObjects(ctx context.Context, owner string, filter map[string]any, cursor any, limit *int) (map[string]any, error) {
	var out map[string]any
	err := c.Call(ctx, "suix_getOwnedObjects", []any{owner, filter, cursor, intOrNil(limit)}, &out)
	return out, err
}

func (c *CoreClient) GetBalance(ctx context.Context, owner, coinType string) (map[string]any, error) {
	var out map[string]any
	err := c.Call(ctx, "suix_getBalance", []any{owner, emptyToNil(coinType)}, &out)
	return out, err
}

func (c *CoreClient) ListBalances(ctx context.Context, owner string) ([]map[string]any, error) {
	var out []map[string]any
	err := c.Call(ctx, "suix_getAllBalances", []any{owner}, &out)
	return out, err
}

func (c *CoreClient) GetCoinMetadata(ctx context.Context, coinType string) (map[string]any, error) {
	var out map[string]any
	err := c.Call(ctx, "suix_getCoinMetadata", []any{coinType}, &out)
	return out, err
}

func (c *CoreClient) GetTransaction(ctx context.Context, digest string, include map[string]any) (map[string]any, error) {
	var out map[string]any
	err := c.Call(ctx, "sui_getTransactionBlock", []any{digest, include}, &out)
	return out, err
}

func (c *CoreClient) ExecuteTransaction(ctx context.Context, txBytesBase64 string, signatures []string, include map[string]any, requestType string) (map[string]any, error) {
	var out map[string]any
	err := c.Call(ctx, "sui_executeTransactionBlock", []any{txBytesBase64, signatures, include, requestType}, &out)
	return out, err
}

func (c *CoreClient) SimulateTransaction(ctx context.Context, txBytesBase64 string) (map[string]any, error) {
	var out map[string]any
	err := c.Call(ctx, "sui_dryRunTransactionBlock", []any{txBytesBase64}, &out)
	return out, err
}

func (c *CoreClient) GetReferenceGasPrice(ctx context.Context) (string, error) {
	var out string
	err := c.Call(ctx, "suix_getReferenceGasPrice", []any{}, &out)
	return out, err
}

func (c *CoreClient) GetCurrentSystemState(ctx context.Context) (map[string]any, error) {
	var out map[string]any
	err := c.Call(ctx, "suix_getLatestSuiSystemState", []any{}, &out)
	return out, err
}

func (c *CoreClient) GetChainIdentifier(ctx context.Context) (string, error) {
	var out string
	err := c.Call(ctx, "sui_getChainIdentifier", []any{}, &out)
	return out, err
}

func (c *CoreClient) ListDynamicFields(ctx context.Context, parentObjectID string, cursor any, limit *int) (map[string]any, error) {
	var out map[string]any
	err := c.Call(ctx, "suix_getDynamicFields", []any{parentObjectID, cursor, intOrNil(limit)}, &out)
	return out, err
}

func (c *CoreClient) GetDynamicFieldObject(ctx context.Context, parentObjectID string, name any) (map[string]any, error) {
	var out map[string]any
	err := c.Call(ctx, "suix_getDynamicFieldObject", []any{parentObjectID, name}, &out)
	return out, err
}

func (c *CoreClient) VerifyZkLoginSignature(ctx context.Context, signature string, bytes string, intentScope string, author string) (map[string]any, error) {
	var out map[string]any
	err := c.Call(ctx, "sui_verifyZkLoginSignature", []any{signature, bytes, intentScope, author}, &out)
	return out, err
}

func (c *CoreClient) GetMoveFunction(ctx context.Context, packageID, module, function string) (map[string]any, error) {
	var out map[string]any
	err := c.Call(ctx, "sui_getMoveFunctionArgTypes", []any{packageID, module, function}, &out)
	if err != nil {
		return nil, err
	}
	return map[string]any{"function": out}, nil
}

func (c *CoreClient) DefaultNameServiceName(ctx context.Context, address string) (string, error) {
	var out []string
	err := c.Call(ctx, "suix_resolveNameServiceNames", []any{address, nil, 1}, &out)
	if err != nil {
		return "", err
	}
	if len(out) == 0 {
		return "", nil
	}
	return out[0], nil
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

func emptyToNil(v string) any {
	if v == "" {
		return nil
	}
	return v
}
