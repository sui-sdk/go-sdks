package deepbookv3

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/sui-sdks/go-sdks/deepbook_v3/types"
)

type mockClient struct{}

func (m mockClient) Network() string { return "testnet" }
func (m mockClient) Call(ctx context.Context, method string, params []any, out any) error {
	_ = ctx
	_ = params
	if p, ok := out.(*map[string]any); ok {
		switch method {
		case "sui_dryRunTransactionBlock":
			*p = map[string]any{
				"commandResults": []any{
					map[string]any{
						"returnValues": []any{
							map[string]any{"bcs": base64.StdEncoding.EncodeToString([]byte{100, 0, 0, 0, 0, 0, 0, 0})},
						},
					},
				},
			}
		}
	}
	return nil
}

func TestDeepBookClientCheckManagerBalance(t *testing.T) {
	client := NewClient(ClientOptions{
		Client:  mockClient{},
		Network: "testnet",
		Options: Options{
			Address: "0x1",
			BalanceManagers: map[string]types.BalanceManager{
				"m1": {Address: "0x2"},
			},
		},
	})

	res, err := client.CheckManagerBalance(context.Background(), "m1", "SUI")
	if err != nil {
		t.Fatalf("CheckManagerBalance failed: %v", err)
	}
	if res["coinType"].(string) == "" {
		t.Fatalf("coinType missing")
	}
}

func TestDeepBookClientWhitelisted(t *testing.T) {
	client := NewClient(ClientOptions{Client: mockClient{}, Network: "testnet", Options: Options{Address: "0x1"}})
	v, err := client.Whitelisted(context.Background(), "DEEP_SUI")
	if err != nil {
		t.Fatalf("Whitelisted failed: %v", err)
	}
	if v {
		t.Fatalf("expected false for return byte 100")
	}
}
