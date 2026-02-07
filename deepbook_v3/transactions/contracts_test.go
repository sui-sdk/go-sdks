package transactions

import (
	"testing"

	"github.com/sui-sdks/go-sdks/deepbook_v3/types"
	"github.com/sui-sdks/go-sdks/deepbook_v3/utils"
	stx "github.com/sui-sdks/go-sdks/sui/transactions"
)

func newTestConfig() *utils.DeepBookConfig {
	return utils.NewDeepBookConfig(utils.ConfigOptions{
		Address: "0x1",
		Network: "testnet",
		BalanceManagers: map[string]types.BalanceManager{
			"m1": {Address: "0x2"},
		},
		MarginManagers: map[string]types.MarginManager{
			"mm1": {Address: "0x3", PoolKey: "DEEP_SUI"},
		},
	})
}

func firstFunction(tx *stx.Transaction) string {
	return commandFunction(tx, 0)
}

func commandFunction(tx *stx.Transaction, idx int) string {
	data := tx.GetData()
	if len(data.Commands) <= idx {
		return ""
	}
	mv := data.Commands[idx]["MoveCall"].(map[string]any)
	return mv["function"].(string)
}

func lastFunction(tx *stx.Transaction) string {
	data := tx.GetData()
	if len(data.Commands) == 0 {
		return ""
	}
	return commandFunction(tx, len(data.Commands)-1)
}

func TestBalanceManagerContract_MethodTargets(t *testing.T) {
	cfg := newTestConfig()
	contract := NewBalanceManagerContract(cfg)
	tx := stx.NewTransaction()
	contract.CheckManagerBalance(tx, "m1", "SUI")
	if got := firstFunction(tx); got != "balance" {
		t.Fatalf("expected balance, got %s", got)
	}
}

func TestDeepBookContract_MethodTargets(t *testing.T) {
	cfg := newTestConfig()
	bm := NewBalanceManagerContract(cfg)
	contract := NewDeepBookContract(cfg, bm)
	tx := stx.NewTransaction()
	contract.GetQuoteQuantityOut(tx, "DEEP_SUI", 1.0)
	if got := firstFunction(tx); got != "get_quote_quantity_out" {
		t.Fatalf("expected get_quote_quantity_out, got %s", got)
	}
}

func TestDeepBookAdminContract_MethodTargets(t *testing.T) {
	cfg := newTestConfig()
	cfg.AdminCap = "0x999"
	contract := NewDeepBookAdminContract(cfg)
	tx := stx.NewTransaction()
	contract.EnableVersion(tx, 1)
	if got := firstFunction(tx); got != "enable_version" {
		t.Fatalf("expected enable_version, got %s", got)
	}
}

func TestPoolProxyContract_MethodTargets(t *testing.T) {
	cfg := newTestConfig()
	contract := NewPoolProxyContract(cfg)
	tx := stx.NewTransaction()
	contract.CancelAllOrders(tx, "mm1")
	if got := firstFunction(tx); got != "cancel_all_orders" {
		t.Fatalf("expected cancel_all_orders, got %s", got)
	}
}

func TestDeepBookContract_MethodMatrix(t *testing.T) {
	cfg := newTestConfig()
	bm := NewBalanceManagerContract(cfg)
	c := NewDeepBookContract(cfg, bm)

	tests := []struct {
		name string
		call func(*stx.Transaction)
		want string
	}{
		{"CancelOrders", func(tx *stx.Transaction) { c.CancelOrders(tx, "DEEP_SUI", "m1", []string{"1", "2"}) }, "cancel_orders"},
		{"CancelAllOrders", func(tx *stx.Transaction) { c.CancelAllOrders(tx, "DEEP_SUI", "m1") }, "cancel_all_orders"},
		{"WithdrawSettledAmounts", func(tx *stx.Transaction) { c.WithdrawSettledAmounts(tx, "DEEP_SUI", "m1") }, "withdraw_settled_amounts"},
		{"GetOrder", func(tx *stx.Transaction) { c.GetOrder(tx, "DEEP_SUI", "1") }, "get_order"},
		{"GetOrders", func(tx *stx.Transaction) { c.GetOrders(tx, "DEEP_SUI", []string{"1", "2"}) }, "get_orders"},
		{"MidPrice", func(tx *stx.Transaction) { c.MidPrice(tx, "DEEP_SUI") }, "mid_price"},
		{"GetBaseQuantityOut", func(tx *stx.Transaction) { c.GetBaseQuantityOut(tx, "DEEP_SUI", 1) }, "get_base_quantity_out"},
		{"GetQuantityOutInputFee", func(tx *stx.Transaction) { c.GetQuantityOutInputFee(tx, "DEEP_SUI", 1, 0) }, "get_quantity_out_input_fee"},
		{"GetBaseQuantityIn", func(tx *stx.Transaction) { c.GetBaseQuantityIn(tx, "DEEP_SUI", 1, true) }, "get_base_quantity_in"},
		{"CanPlaceLimitOrder", func(tx *stx.Transaction) {
			c.CanPlaceLimitOrder(tx, types.CanPlaceLimitOrderParams{
				PoolKey: "DEEP_SUI", BalanceManagerKey: "m1", Price: 1, Quantity: 1, IsBid: true, PayWithDeep: true, ExpireTimestamp: 100,
			})
		}, "can_place_limit_order"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tx := stx.NewTransaction()
			tc.call(tx)
			if got := lastFunction(tx); got != tc.want {
				t.Fatalf("expected %s, got %s", tc.want, got)
			}
		})
	}
}

func TestGovernanceAndFlashLoanTargets(t *testing.T) {
	cfg := newTestConfig()
	bm := NewBalanceManagerContract(cfg)
	g := NewGovernanceContract(cfg, bm)
	f := NewFlashLoanContract(cfg)

	tx1 := stx.NewTransaction()
	g.Vote(tx1, "DEEP_SUI", "m1", "7")
	if got := lastFunction(tx1); got != "vote" {
		t.Fatalf("expected vote, got %s", got)
	}

	tx2 := stx.NewTransaction()
	f.ReturnBaseAsset(tx2, "DEEP_SUI", tx2.Object("0xa"), tx2.Object("0xb"))
	if got := firstFunction(tx2); got != "return_flashloan_base" {
		t.Fatalf("expected return_flashloan_base, got %s", got)
	}
}

func TestMarginContractsTargets(t *testing.T) {
	cfg := newTestConfig()

	mm := NewMarginManagerContract(cfg)
	tx1 := stx.NewTransaction()
	mm.GetMarginAccountOrderDetails(tx1, "mm1")
	if got := commandFunction(tx1, 0); got != "balance_manager" {
		t.Fatalf("expected first call balance_manager, got %s", got)
	}
	if got := commandFunction(tx1, 1); got != "get_account_order_details" {
		t.Fatalf("expected second call get_account_order_details, got %s", got)
	}

	mp := NewMarginPoolContract(cfg)
	tx2 := stx.NewTransaction()
	mp.InterestRate(tx2, "SUI")
	if got := firstFunction(tx2); got != "interest_rate" {
		t.Fatalf("expected interest_rate, got %s", got)
	}

	mr := NewMarginRegistryContract(cfg)
	tx3 := stx.NewTransaction()
	mr.GetMarginManagerIDs(tx3, "0x1")
	if got := firstFunction(tx3); got != "get_margin_manager_ids" {
		t.Fatalf("expected get_margin_manager_ids, got %s", got)
	}

	tpsl := NewMarginTPSLContract(cfg)
	tx4 := stx.NewTransaction()
	tpsl.CancelAllConditionalOrders(tx4, "mm1")
	if got := firstFunction(tx4); got != "cancel_all_conditional_orders" {
		t.Fatalf("expected cancel_all_conditional_orders, got %s", got)
	}
}
