package deepbookv3

import (
	"context"
	"testing"
)

type emptyResponseClient struct{}

func (e emptyResponseClient) Network() string { return "testnet" }
func (e emptyResponseClient) Call(ctx context.Context, method string, params []any, out any) error {
	_ = ctx
	_ = method
	_ = params
	if p, ok := out.(*map[string]any); ok {
		*p = map[string]any{}
	}
	return nil
}

func TestClientHandlesMissingReturnValues(t *testing.T) {
	c := NewClient(ClientOptions{
		Client:  emptyResponseClient{},
		Network: "testnet",
		Options: Options{
			Address: "0x1",
		},
	})

	_, err := c.Whitelisted(context.Background(), "DEEP_SUI")
	if err == nil {
		t.Fatalf("expected error for missing commandResults")
	}

	_, err = c.GetOrder(context.Background(), "DEEP_SUI", "1")
	if err == nil {
		t.Fatalf("expected error for missing returnValues")
	}
}
