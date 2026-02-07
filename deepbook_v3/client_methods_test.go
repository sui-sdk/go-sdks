package deepbookv3

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/sui-sdks/go-sdks/bcs"
	"github.com/sui-sdks/go-sdks/deepbook_v3/types"
	"github.com/sui-sdks/go-sdks/deepbook_v3/utils"
)

type deepbookMethodMockClient struct{}

func (m deepbookMethodMockClient) Network() string { return "testnet" }

func encodeU64(v uint64) string {
	w := bcs.NewWriter(nil)
	_ = w.Write64(v)
	return base64.StdEncoding.EncodeToString(w.ToBytes())
}

func moveFunctionFromTxBase64(txb64 string) string {
	raw, err := base64.StdEncoding.DecodeString(txb64)
	if err != nil {
		return ""
	}
	var payload map[string]any
	if err := json.Unmarshal(raw, &payload); err != nil {
		return ""
	}
	cmds, ok := payload["Commands"].([]any)
	if !ok || len(cmds) == 0 {
		return ""
	}
	cmd, ok := cmds[0].(map[string]any)
	if !ok {
		return ""
	}
	mv, ok := cmd["MoveCall"].(map[string]any)
	if !ok {
		return ""
	}
	fn, _ := mv["function"].(string)
	return fn
}

func (m deepbookMethodMockClient) Call(ctx context.Context, method string, params []any, out any) error {
	_ = ctx
	if method != "sui_dryRunTransactionBlock" {
		return nil
	}
	fn := ""
	if len(params) > 0 {
		if txb64, ok := params[0].(string); ok {
			fn = moveFunctionFromTxBase64(txb64)
		}
	}

	firstRet := []any{map[string]any{"bcs": encodeU64(100)}}
	if fn == "whitelisted" {
		firstRet = []any{map[string]any{"bcs": base64.StdEncoding.EncodeToString([]byte{1})}}
	}
	if fn == "get_quote_quantity_out" || fn == "get_base_quantity_out" || fn == "get_quantity_out" {
		firstRet = []any{
			map[string]any{"bcs": encodeU64(100)},
			map[string]any{"bcs": encodeU64(200)},
			map[string]any{"bcs": encodeU64(300)},
		}
	}

	if p, ok := out.(*map[string]any); ok {
		*p = map[string]any{
			"commandResults": []any{
				map[string]any{"returnValues": firstRet},
				map[string]any{"returnValues": []any{map[string]any{"bcs": encodeU64(777)}}},
			},
		}
	}
	return nil
}

func newMethodTestClient() *Client {
	return NewClient(ClientOptions{
		Client:  deepbookMethodMockClient{},
		Network: "testnet",
		Options: Options{
			Address: "0x1",
			BalanceManagers: map[string]types.BalanceManager{
				"m1": {Address: "0x2"},
			},
			MarginManagers: map[string]types.MarginManager{
				"mm1": {Address: "0x3", PoolKey: "DEEP_SUI"},
			},
		},
	})
}

func TestDeepBookClientQuantityAndPriceMethods(t *testing.T) {
	c := newMethodTestClient()
	ctx := context.Background()

	q1, err := c.GetQuoteQuantityOut(ctx, "DEEP_SUI", 1)
	if err != nil || q1["deepRequired"] == nil {
		t.Fatalf("GetQuoteQuantityOut failed: %v", err)
	}
	q2, err := c.GetBaseQuantityOut(ctx, "DEEP_SUI", 1)
	if err != nil || q2["deepRequired"] == nil {
		t.Fatalf("GetBaseQuantityOut failed: %v", err)
	}
	q3, err := c.GetQuantityOut(ctx, "DEEP_SUI", 1, 0)
	if err != nil || q3["deepRequired"] == nil {
		t.Fatalf("GetQuantityOut failed: %v", err)
	}

	mid, err := c.MidPrice(ctx, "DEEP_SUI")
	if err != nil || mid <= 0 {
		t.Fatalf("MidPrice failed: %v, value=%v", err, mid)
	}
}

func TestDeepBookClientInputFeeAndInMethods(t *testing.T) {
	c := newMethodTestClient()
	ctx := context.Background()

	if _, err := c.GetQuoteQuantityOutInputFee(ctx, "DEEP_SUI", 1); err != nil {
		t.Fatalf("GetQuoteQuantityOutInputFee failed: %v", err)
	}
	if _, err := c.GetBaseQuantityOutInputFee(ctx, "DEEP_SUI", 1); err != nil {
		t.Fatalf("GetBaseQuantityOutInputFee failed: %v", err)
	}
	if _, err := c.GetQuantityOutInputFee(ctx, "DEEP_SUI", 1, 0); err != nil {
		t.Fatalf("GetQuantityOutInputFee failed: %v", err)
	}
	if _, err := c.GetBaseQuantityIn(ctx, "DEEP_SUI", 1, true); err != nil {
		t.Fatalf("GetBaseQuantityIn failed: %v", err)
	}
	if _, err := c.GetQuoteQuantityIn(ctx, "DEEP_SUI", 1, false); err != nil {
		t.Fatalf("GetQuoteQuantityIn failed: %v", err)
	}
	if _, err := c.GetOrderDeepRequired(ctx, "DEEP_SUI", 1, 1); err != nil {
		t.Fatalf("GetOrderDeepRequired failed: %v", err)
	}
}

func TestDeepBookClientRawBCSMethods(t *testing.T) {
	c := newMethodTestClient()
	ctx := context.Background()

	cases := []struct {
		name string
		call func() (string, error)
	}{
		{"GetOrder", func() (string, error) { return c.GetOrder(ctx, "DEEP_SUI", "1") }},
		{"GetOrders", func() (string, error) { return c.GetOrders(ctx, "DEEP_SUI", []string{"1"}) }},
		{"AccountOpenOrders", func() (string, error) { return c.AccountOpenOrders(ctx, "DEEP_SUI", "m1") }},
		{"VaultBalances", func() (string, error) { return c.VaultBalances(ctx, "DEEP_SUI") }},
		{"GetPoolIDByAssets", func() (string, error) { return c.GetPoolIDByAssets(ctx, "0x2::sui::SUI", "0x3::coin::C") }},
		{"PoolTradeParams", func() (string, error) { return c.PoolTradeParams(ctx, "DEEP_SUI") }},
		{"PoolBookParams", func() (string, error) { return c.PoolBookParams(ctx, "DEEP_SUI") }},
		{"Account", func() (string, error) { return c.Account(ctx, "DEEP_SUI", "m1") }},
		{"LockedBalance", func() (string, error) { return c.LockedBalance(ctx, "DEEP_SUI", "m1") }},
		{"GetPoolDeepPrice", func() (string, error) { return c.GetPoolDeepPrice(ctx, "DEEP_SUI") }},
		{"BalanceManagerReferralOwner", func() (string, error) { return c.BalanceManagerReferralOwner(ctx, "0xaaa") }},
		{"GetMarginAccountOrderDetails", func() (string, error) { return c.GetMarginAccountOrderDetails(ctx, "mm1") }},
		{"GetAccountOrderDetails", func() (string, error) { return c.GetAccountOrderDetails(ctx, "DEEP_SUI", "m1") }},
		{"PoolTradeParamsNext", func() (string, error) { return c.PoolTradeParamsNext(ctx, "DEEP_SUI") }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v, err := tc.call()
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}
			if v == "" {
				t.Fatalf("%s returned empty bcs", tc.name)
			}
		})
	}
}

func TestDeepBookClientWhitelistedAndPriceInfoAge(t *testing.T) {
	c := newMethodTestClient()
	ctx := context.Background()

	w, err := c.Whitelisted(ctx, "DEEP_SUI")
	if err != nil || !w {
		t.Fatalf("Whitelisted failed: %v, value=%v", err, w)
	}

	age, err := c.GetPriceInfoObjectAge(ctx, "SUI")
	if err != nil {
		t.Fatalf("GetPriceInfoObjectAge failed: %v", err)
	}
	if age != -1 {
		t.Fatalf("expected -1 for coin without price info object, got %d", age)
	}

	custom := NewClient(ClientOptions{
		Client:  deepbookMethodMockClient{},
		Network: "testnet",
		Options: Options{
			Address: "0x1",
			Coins: utils.CoinMap{
				"SUI": {Address: "0x2", Type: "0x2::sui::SUI", Scalar: 1_000_000_000, PriceInfoObjectID: "0xabc"},
			},
			Pools: utils.PoolMap{
				"DEEP_SUI": {Address: "0x11", BaseCoin: "SUI", QuoteCoin: "SUI"},
			},
			BalanceManagers: map[string]types.BalanceManager{
				"m1": {Address: "0x2"},
			},
		},
	})
	age2, err := custom.GetPriceInfoObjectAge(ctx, "SUI")
	if err != nil {
		t.Fatalf("GetPriceInfoObjectAge(custom) failed: %v", err)
	}
	if age2 <= 0 {
		t.Fatalf("expected positive timestamp for price info object age, got %d", age2)
	}
}
