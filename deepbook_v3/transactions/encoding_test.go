package transactions

import (
	"encoding/base64"
	"math"
	"math/big"
	"testing"

	"github.com/sui-sdks/go-sdks/bcs"
	"github.com/sui-sdks/go-sdks/deepbook_v3/types"
	"github.com/sui-sdks/go-sdks/deepbook_v3/utils"
	stx "github.com/sui-sdks/go-sdks/sui/transactions"
)

func findMoveCallByFunction(t *testing.T, tx *stx.Transaction, function string) map[string]any {
	t.Helper()
	for _, cmd := range tx.GetData().Commands {
		mv, ok := cmd["MoveCall"].(map[string]any)
		if !ok {
			continue
		}
		if fn, _ := mv["function"].(string); fn == function {
			return mv
		}
	}
	t.Fatalf("move call %q not found", function)
	return nil
}

func argsFromMoveCall(t *testing.T, mv map[string]any) []stx.Argument {
	t.Helper()
	raw, ok := mv["arguments"]
	if !ok {
		t.Fatalf("missing arguments")
	}
	switch v := raw.(type) {
	case []stx.Argument:
		return v
	case []any:
		out := make([]stx.Argument, 0, len(v))
		for _, it := range v {
			a, ok := it.(map[string]any)
			if !ok {
				t.Fatalf("invalid argument item")
			}
			out = append(out, stx.Argument(a))
		}
		return out
	default:
		t.Fatalf("unexpected arguments type %T", raw)
	}
	return nil
}

func pureBytesFromArg(t *testing.T, tx *stx.Transaction, arg stx.Argument) []byte {
	t.Helper()
	kind, _ := arg["$kind"].(string)
	if kind != "Input" {
		t.Fatalf("argument is not Input kind: %v", kind)
	}
	idx, ok := arg["Input"].(int)
	if !ok {
		t.Fatalf("input index is not int: %T", arg["Input"])
	}
	input := tx.GetData().Inputs[idx]
	if k, _ := input["$kind"].(string); k != "Pure" {
		t.Fatalf("input is not Pure: %v", k)
	}
	pure := input["Pure"].(map[string]any)
	b64 := pure["bytes"].(string)
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		t.Fatalf("decode pure bytes failed: %v", err)
	}
	return data
}

func readU64FromArg(t *testing.T, tx *stx.Transaction, arg stx.Argument) uint64 {
	t.Helper()
	r := bcs.NewReader(pureBytesFromArg(t, tx, arg))
	v, err := r.Read64()
	if err != nil {
		t.Fatalf("read u64 failed: %v", err)
	}
	return v
}

func readBoolFromArg(t *testing.T, tx *stx.Transaction, arg stx.Argument) bool {
	t.Helper()
	b := pureBytesFromArg(t, tx, arg)
	if len(b) != 1 {
		t.Fatalf("bool arg length must be 1, got %d", len(b))
	}
	return b[0] == 1
}

func readU128StringFromArg(t *testing.T, tx *stx.Transaction, arg stx.Argument) string {
	t.Helper()
	b := pureBytesFromArg(t, tx, arg)
	if len(b) != 16 {
		t.Fatalf("u128 arg length must be 16, got %d", len(b))
	}
	be := make([]byte, 16)
	for i := 0; i < 16; i++ {
		be[i] = b[15-i]
	}
	n := new(big.Int).SetBytes(be)
	return n.String()
}

func TestDeepBookPlaceLimitOrderEncoding(t *testing.T) {
	cfg := newTestConfig()
	bm := NewBalanceManagerContract(cfg)
	c := NewDeepBookContract(cfg, bm)
	tx := stx.NewTransaction()

	params := types.PlaceLimitOrderParams{
		PoolKey:            "DEEP_SUI",
		BalanceManagerKey:  "m1",
		ClientOrderID:      "42",
		Price:              1.25,
		Quantity:           2.5,
		IsBid:              true,
		OrderType:          types.OrderTypePostOnly,
		SelfMatchingOption: types.CancelMaker,
		PayWithDeep:        true,
	}
	c.PlaceLimitOrder(tx, params)

	mv := findMoveCallByFunction(t, tx, "place_limit_order")
	args := argsFromMoveCall(t, mv)
	if len(args) != 12 {
		t.Fatalf("expected 12 args, got %d", len(args))
	}

	pool := cfg.GetPool("DEEP_SUI")
	base := cfg.GetCoin(pool.BaseCoin)
	quote := cfg.GetCoin(pool.QuoteCoin)
	wantPrice := uint64(math.Round((params.Price * utils.FloatScalar * quote.Scalar) / base.Scalar))
	wantQty := uint64(math.Round(params.Quantity * base.Scalar))

	if got := readU64FromArg(t, tx, args[3]); got != 42 {
		t.Fatalf("clientOrderId mismatch: got %d", got)
	}
	if got := pureBytesFromArg(t, tx, args[4])[0]; got != byte(types.OrderTypePostOnly) {
		t.Fatalf("orderType mismatch: got %d", got)
	}
	if got := pureBytesFromArg(t, tx, args[5])[0]; got != byte(types.CancelMaker) {
		t.Fatalf("selfMatching mismatch: got %d", got)
	}
	if got := readU64FromArg(t, tx, args[6]); got != wantPrice {
		t.Fatalf("price mismatch: got %d want %d", got, wantPrice)
	}
	if got := readU64FromArg(t, tx, args[7]); got != wantQty {
		t.Fatalf("quantity mismatch: got %d want %d", got, wantQty)
	}
	if got := readBoolFromArg(t, tx, args[8]); !got {
		t.Fatalf("isBid mismatch")
	}
	if got := readBoolFromArg(t, tx, args[9]); !got {
		t.Fatalf("payWithDeep mismatch")
	}
	if got := readU64FromArg(t, tx, args[10]); got != utils.MaxTimestamp {
		t.Fatalf("expiration mismatch: got %d want %d", got, utils.MaxTimestamp)
	}
}

func TestDeepBookCancelOrdersEncoding(t *testing.T) {
	cfg := newTestConfig()
	bm := NewBalanceManagerContract(cfg)
	c := NewDeepBookContract(cfg, bm)
	tx := stx.NewTransaction()

	c.CancelOrder(tx, "DEEP_SUI", "m1", "123456789")
	mvCancelOne := findMoveCallByFunction(t, tx, "cancel_order")
	argsOne := argsFromMoveCall(t, mvCancelOne)
	if got := readU128StringFromArg(t, tx, argsOne[3]); got != "123456789" {
		t.Fatalf("cancel_order u128 mismatch: got %s", got)
	}

	tx2 := stx.NewTransaction()
	c.CancelOrders(tx2, "DEEP_SUI", "m1", []string{"1", "340282366920938463463374607431768211455"})
	mv := findMoveCallByFunction(t, tx2, "cancel_orders")
	args := argsFromMoveCall(t, mv)
	raw := pureBytesFromArg(t, tx2, args[3])
	r := bcs.NewReader(raw)
	n, err := r.ReadULEB()
	if err != nil {
		t.Fatalf("read vec len failed: %v", err)
	}
	if n != 2 {
		t.Fatalf("vector length mismatch: got %d", n)
	}
	v1, err := r.ReadBytes(16)
	if err != nil {
		t.Fatalf("read first u128 failed: %v", err)
	}
	v2, err := r.ReadBytes(16)
	if err != nil {
		t.Fatalf("read second u128 failed: %v", err)
	}
	one := new(big.Int).SetUint64(1)
	max := new(big.Int)
	max.SetString("340282366920938463463374607431768211455", 10)
	gotOne := new(big.Int).SetBytes(reverse(v1))
	gotMax := new(big.Int).SetBytes(reverse(v2))
	if gotOne.Cmp(one) != 0 || gotMax.Cmp(max) != 0 {
		t.Fatalf("u128 vector values mismatch")
	}
}

func TestDeepBookCanPlaceLimitOrderEncoding(t *testing.T) {
	cfg := newTestConfig()
	bm := NewBalanceManagerContract(cfg)
	c := NewDeepBookContract(cfg, bm)
	tx := stx.NewTransaction()

	params := types.CanPlaceLimitOrderParams{
		PoolKey:           "DEEP_SUI",
		BalanceManagerKey: "m1",
		Price:             0.5,
		Quantity:          3,
		IsBid:             false,
		PayWithDeep:       true,
		ExpireTimestamp:   987654321,
	}
	c.CanPlaceLimitOrder(tx, params)

	mv := findMoveCallByFunction(t, tx, "can_place_limit_order")
	args := argsFromMoveCall(t, mv)
	if got := readBoolFromArg(t, tx, args[4]); got {
		t.Fatalf("isBid should be false")
	}
	if got := readBoolFromArg(t, tx, args[5]); !got {
		t.Fatalf("payWithDeep should be true")
	}
	if got := readU64FromArg(t, tx, args[6]); got != 987654321 {
		t.Fatalf("expire timestamp mismatch: %d", got)
	}
}

func TestPoolProxyPlaceLimitOrderDefaultExpiration(t *testing.T) {
	cfg := newTestConfig()
	c := NewPoolProxyContract(cfg)
	tx := stx.NewTransaction()

	c.PlaceLimitOrder(tx, types.PlaceMarginLimitOrderParams{
		PoolKey:          "DEEP_SUI",
		MarginManagerKey: "mm1",
		ClientOrderID:    "9",
		Price:            1.1,
		Quantity:         2.2,
		IsBid:            true,
		PayWithDeep:      false,
	})
	mv := findMoveCallByFunction(t, tx, "place_limit_order")
	args := argsFromMoveCall(t, mv)
	if got := readU64FromArg(t, tx, args[8]); got != utils.MaxTimestamp {
		t.Fatalf("default expiration mismatch: %d", got)
	}
	if got := readBoolFromArg(t, tx, args[7]); got {
		t.Fatalf("payWithDeep should be false")
	}
}

func TestFlashLoanBorrowEncoding(t *testing.T) {
	cfg := newTestConfig()
	c := NewFlashLoanContract(cfg)
	tx := stx.NewTransaction()

	c.BorrowBaseAsset(tx, "DEEP_SUI", 1.5)
	mv := findMoveCallByFunction(t, tx, "borrow_flashloan_base")
	args := argsFromMoveCall(t, mv)
	pool := cfg.GetPool("DEEP_SUI")
	base := cfg.GetCoin(pool.BaseCoin)
	want := uint64(math.Round(1.5 * base.Scalar))
	if got := readU64FromArg(t, tx, args[1]); got != want {
		t.Fatalf("flashloan base qty mismatch: got %d want %d", got, want)
	}
}

func reverse(in []byte) []byte {
	out := make([]byte, len(in))
	for i := 0; i < len(in); i++ {
		out[i] = in[len(in)-1-i]
	}
	return out
}
