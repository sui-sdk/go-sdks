package transactions

import (
	"math"

	"github.com/sui-sdks/go-sdks/deepbook_v3/types"
	"github.com/sui-sdks/go-sdks/deepbook_v3/utils"
	stx "github.com/sui-sdks/go-sdks/sui/transactions"
)

type MarginAdminContract struct{ config *utils.DeepBookConfig }
type MarginMaintainerContract struct{ config *utils.DeepBookConfig }
type MarginManagerContract struct{ config *utils.DeepBookConfig }
type MarginPoolContract struct{ config *utils.DeepBookConfig }
type MarginRegistryContract struct{ config *utils.DeepBookConfig }
type MarginLiquidationsContract struct{ config *utils.DeepBookConfig }
type PoolProxyContract struct{ config *utils.DeepBookConfig }
type MarginTPSLContract struct{ config *utils.DeepBookConfig }

func NewMarginAdminContract(config *utils.DeepBookConfig) *MarginAdminContract {
	return &MarginAdminContract{config: config}
}
func NewMarginMaintainerContract(config *utils.DeepBookConfig) *MarginMaintainerContract {
	return &MarginMaintainerContract{config: config}
}
func NewMarginManagerContract(config *utils.DeepBookConfig) *MarginManagerContract {
	return &MarginManagerContract{config: config}
}
func NewMarginPoolContract(config *utils.DeepBookConfig) *MarginPoolContract {
	return &MarginPoolContract{config: config}
}
func NewMarginRegistryContract(config *utils.DeepBookConfig) *MarginRegistryContract {
	return &MarginRegistryContract{config: config}
}
func NewMarginLiquidationsContract(config *utils.DeepBookConfig) *MarginLiquidationsContract {
	return &MarginLiquidationsContract{config: config}
}
func NewPoolProxyContract(config *utils.DeepBookConfig) *PoolProxyContract {
	return &PoolProxyContract{config: config}
}
func NewMarginTPSLContract(config *utils.DeepBookConfig) *MarginTPSLContract {
	return &MarginTPSLContract{config: config}
}

func (c *MarginManagerContract) managerTypes(managerKey string) (types.Coin, types.Coin, types.MarginManager, types.Pool) {
	manager := c.config.GetMarginManager(managerKey)
	pool := c.config.GetPool(manager.PoolKey)
	return c.config.GetCoin(pool.BaseCoin), c.config.GetCoin(pool.QuoteCoin), manager, pool
}

func (c *MarginManagerContract) NewMarginManager(tx *stx.Transaction, poolKey string) stx.Argument {
	pool := c.config.GetPool(poolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	return tx.MoveCall(c.config.MarginPackageID+"::margin_manager::new", []stx.Argument{
		tx.Object(pool.Address), tx.Object(c.config.RegistryID), tx.Object(c.config.MarginRegistryID), tx.Object("0x6"),
	}, []string{base.Type, quote.Type})
}

func (c *MarginManagerContract) NewMarginManagerWithInitializer(tx *stx.Transaction, poolKey string) stx.Argument {
	pool := c.config.GetPool(poolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	return tx.MoveCall(c.config.MarginPackageID+"::margin_manager::new_with_initializer", []stx.Argument{
		tx.Object(pool.Address), tx.Object(c.config.RegistryID), tx.Object(c.config.MarginRegistryID), tx.Object("0x6"),
	}, []string{base.Type, quote.Type})
}

func (c *MarginManagerContract) DepositBase(tx *stx.Transaction, params types.DepositParams) stx.Argument {
	base, quote, manager, _ := c.managerTypes(params.ManagerKey)
	amount := uint64(math.Round(params.Amount * base.Scalar))
	return tx.MoveCall(c.config.MarginPackageID+"::margin_manager::deposit", []stx.Argument{
		tx.Object(manager.Address), pureU64(tx, amount), tx.Object(c.config.MarginRegistryID),
	}, []string{base.Type, quote.Type, base.Type})
}

func (c *MarginManagerContract) DepositQuote(tx *stx.Transaction, params types.DepositParams) stx.Argument {
	base, quote, manager, _ := c.managerTypes(params.ManagerKey)
	amount := uint64(math.Round(params.Amount * quote.Scalar))
	return tx.MoveCall(c.config.MarginPackageID+"::margin_manager::deposit", []stx.Argument{
		tx.Object(manager.Address), pureU64(tx, amount), tx.Object(c.config.MarginRegistryID),
	}, []string{base.Type, quote.Type, quote.Type})
}

func (c *MarginManagerContract) WithdrawBase(tx *stx.Transaction, managerKey string, amount float64) stx.Argument {
	base, quote, manager, _ := c.managerTypes(managerKey)
	input := uint64(math.Round(amount * base.Scalar))
	return tx.MoveCall(c.config.MarginPackageID+"::margin_manager::withdraw", []stx.Argument{
		tx.Object(manager.Address), pureU64(tx, input), tx.Object(c.config.MarginRegistryID),
	}, []string{base.Type, quote.Type, base.Type})
}

func (c *MarginManagerContract) WithdrawQuote(tx *stx.Transaction, managerKey string, amount float64) stx.Argument {
	base, quote, manager, _ := c.managerTypes(managerKey)
	input := uint64(math.Round(amount * quote.Scalar))
	return tx.MoveCall(c.config.MarginPackageID+"::margin_manager::withdraw", []stx.Argument{
		tx.Object(manager.Address), pureU64(tx, input), tx.Object(c.config.MarginRegistryID),
	}, []string{base.Type, quote.Type, quote.Type})
}

func (c *MarginManagerContract) BorrowBase(tx *stx.Transaction, managerKey string, amount float64) stx.Argument {
	base, quote, manager, _ := c.managerTypes(managerKey)
	input := uint64(math.Round(amount * base.Scalar))
	return tx.MoveCall(c.config.MarginPackageID+"::margin_manager::borrow_base", []stx.Argument{
		tx.Object(manager.Address), pureU64(tx, input), tx.Object(c.config.MarginRegistryID), tx.Object("0x6"),
	}, []string{base.Type, quote.Type})
}

func (c *MarginManagerContract) BorrowQuote(tx *stx.Transaction, managerKey string, amount float64) stx.Argument {
	base, quote, manager, _ := c.managerTypes(managerKey)
	input := uint64(math.Round(amount * quote.Scalar))
	return tx.MoveCall(c.config.MarginPackageID+"::margin_manager::borrow_quote", []stx.Argument{
		tx.Object(manager.Address), pureU64(tx, input), tx.Object(c.config.MarginRegistryID), tx.Object("0x6"),
	}, []string{base.Type, quote.Type})
}

func (c *MarginManagerContract) RepayBase(tx *stx.Transaction, managerKey string, amount float64) stx.Argument {
	base, quote, manager, _ := c.managerTypes(managerKey)
	input := uint64(math.Round(amount * base.Scalar))
	return tx.MoveCall(c.config.MarginPackageID+"::margin_manager::repay_base", []stx.Argument{
		tx.Object(manager.Address), pureU64(tx, input), tx.Object(c.config.MarginRegistryID),
	}, []string{base.Type, quote.Type})
}

func (c *MarginManagerContract) RepayQuote(tx *stx.Transaction, managerKey string, amount float64) stx.Argument {
	base, quote, manager, _ := c.managerTypes(managerKey)
	input := uint64(math.Round(amount * quote.Scalar))
	return tx.MoveCall(c.config.MarginPackageID+"::margin_manager::repay_quote", []stx.Argument{
		tx.Object(manager.Address), pureU64(tx, input), tx.Object(c.config.MarginRegistryID),
	}, []string{base.Type, quote.Type})
}

func (c *MarginManagerContract) SetMarginManagerReferral(tx *stx.Transaction, managerKey, referral string) stx.Argument {
	_, _, manager, _ := c.managerTypes(managerKey)
	return tx.MoveCall(c.config.MarginPackageID+"::margin_manager::set_margin_manager_referral", []stx.Argument{
		tx.Object(manager.Address), tx.Object(referral), tx.Object(c.config.MarginRegistryID),
	}, nil)
}

func (c *MarginManagerContract) UnsetMarginManagerReferral(tx *stx.Transaction, managerKey, poolKey string) stx.Argument {
	_, _, manager, pool := c.managerTypes(managerKey)
	_ = poolKey
	return tx.MoveCall(c.config.MarginPackageID+"::margin_manager::unset_margin_manager_referral", []stx.Argument{
		tx.Object(manager.Address), tx.Object(pool.Address), tx.Object(c.config.MarginRegistryID),
	}, nil)
}

func (c *MarginManagerContract) GetMarginAccountOrderDetails(tx *stx.Transaction, managerKey string) stx.Argument {
	base, quote, manager, pool := c.managerTypes(managerKey)
	bm := tx.MoveCall(c.config.MarginPackageID+"::margin_manager::balance_manager", []stx.Argument{
		tx.Object(manager.Address), tx.Object(c.config.MarginRegistryID),
	}, []string{base.Type, quote.Type})
	return tx.MoveCall(c.config.DeepbookPackageID+"::pool::get_account_order_details", []stx.Argument{
		tx.Object(pool.Address), bm,
	}, []string{base.Type, quote.Type})
}

func (c *MarginPoolContract) MintSupplierCap(tx *stx.Transaction) stx.Argument {
	return tx.MoveCall(c.config.MarginPackageID+"::margin_pool::mint_supplier_cap", []stx.Argument{tx.Object(c.config.MarginRegistryID), tx.Object("0x6")}, nil)
}

func (c *MarginPoolContract) GetID(tx *stx.Transaction, coinKey string) stx.Argument {
	coin := c.config.GetCoin(coinKey)
	return tx.MoveCall(c.config.MarginPackageID+"::margin_pool::id", []stx.Argument{tx.Object(c.config.MarginRegistryID)}, []string{coin.Type})
}

func (c *MarginPoolContract) TotalSupply(tx *stx.Transaction, coinKey string) stx.Argument {
	coin := c.config.GetCoin(coinKey)
	return tx.MoveCall(c.config.MarginPackageID+"::margin_pool::total_supply", []stx.Argument{tx.Object(c.config.MarginRegistryID)}, []string{coin.Type})
}

func (c *MarginPoolContract) SupplyShares(tx *stx.Transaction, coinKey string) stx.Argument {
	coin := c.config.GetCoin(coinKey)
	return tx.MoveCall(c.config.MarginPackageID+"::margin_pool::supply_shares", []stx.Argument{tx.Object(c.config.MarginRegistryID)}, []string{coin.Type})
}

func (c *MarginPoolContract) TotalBorrow(tx *stx.Transaction, coinKey string) stx.Argument {
	coin := c.config.GetCoin(coinKey)
	return tx.MoveCall(c.config.MarginPackageID+"::margin_pool::total_borrow", []stx.Argument{tx.Object(c.config.MarginRegistryID)}, []string{coin.Type})
}

func (c *MarginPoolContract) BorrowShares(tx *stx.Transaction, coinKey string) stx.Argument {
	coin := c.config.GetCoin(coinKey)
	return tx.MoveCall(c.config.MarginPackageID+"::margin_pool::borrow_shares", []stx.Argument{tx.Object(c.config.MarginRegistryID)}, []string{coin.Type})
}

func (c *MarginPoolContract) InterestRate(tx *stx.Transaction, coinKey string) stx.Argument {
	coin := c.config.GetCoin(coinKey)
	return tx.MoveCall(c.config.MarginPackageID+"::margin_pool::interest_rate", []stx.Argument{tx.Object(c.config.MarginRegistryID)}, []string{coin.Type})
}

func (c *MarginRegistryContract) PoolEnabled(tx *stx.Transaction, poolKey string) stx.Argument {
	pool := c.config.GetPool(poolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	return tx.MoveCall(c.config.MarginPackageID+"::margin_registry::pool_enabled", []stx.Argument{
		tx.Object(pool.Address), tx.Object(c.config.MarginRegistryID),
	}, []string{base.Type, quote.Type})
}

func (c *MarginRegistryContract) GetMarginPoolID(tx *stx.Transaction, coinKey string) stx.Argument {
	coin := c.config.GetCoin(coinKey)
	return tx.MoveCall(c.config.MarginPackageID+"::margin_registry::get_margin_pool_id", []stx.Argument{
		tx.Object(c.config.MarginRegistryID),
	}, []string{coin.Type})
}

func (c *MarginRegistryContract) GetMarginManagerIDs(tx *stx.Transaction, owner string) stx.Argument {
	return tx.MoveCall(c.config.MarginPackageID+"::margin_registry::get_margin_manager_ids", []stx.Argument{
		tx.PureBytes([]byte(owner)), tx.Object(c.config.MarginRegistryID),
	}, nil)
}

func (c *MarginLiquidationsContract) CreateLiquidationVault(tx *stx.Transaction, liquidationAdminCap string) stx.Argument {
	return tx.MoveCall(c.config.LiquidationPackageID+"::liquidation_vault::create_liquidation_vault", []stx.Argument{
		tx.Object(c.config.MarginRegistryID), tx.Object(liquidationAdminCap),
	}, nil)
}

func (c *MarginLiquidationsContract) Balance(tx *stx.Transaction, vaultID, coinKey string) stx.Argument {
	coin := c.config.GetCoin(coinKey)
	return tx.MoveCall(c.config.LiquidationPackageID+"::liquidation_vault::balance", []stx.Argument{
		tx.Object(vaultID),
	}, []string{coin.Type})
}

func (c *PoolProxyContract) PlaceLimitOrder(tx *stx.Transaction, params types.PlaceMarginLimitOrderParams) stx.Argument {
	manager := c.config.GetMarginManager(params.MarginManagerKey)
	pool := c.config.GetPool(manager.PoolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	price := uint64(math.Round((params.Price * utils.FloatScalar * quote.Scalar) / base.Scalar))
	quantity := uint64(math.Round(params.Quantity * base.Scalar))
	exp := params.Expiration
	if exp == 0 {
		exp = utils.MaxTimestamp
	}
	return tx.MoveCall(c.config.MarginPackageID+"::pool_proxy::place_limit_order", []stx.Argument{
		tx.Object(manager.Address),
		pureU64(tx, parseU64(params.ClientOrderID)),
		pureU8(tx, uint8(params.OrderType)),
		pureU8(tx, uint8(params.SelfMatchingOption)),
		pureU64(tx, price),
		pureU64(tx, quantity),
		pureBool(tx, params.IsBid),
		pureBool(tx, params.PayWithDeep),
		pureU64(tx, exp),
		tx.Object(c.config.MarginRegistryID),
		tx.Object("0x6"),
	}, []string{base.Type, quote.Type})
}

func (c *PoolProxyContract) PlaceMarketOrder(tx *stx.Transaction, params types.PlaceMarginMarketOrderParams) stx.Argument {
	manager := c.config.GetMarginManager(params.MarginManagerKey)
	pool := c.config.GetPool(manager.PoolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	quantity := uint64(math.Round(params.Quantity * base.Scalar))
	return tx.MoveCall(c.config.MarginPackageID+"::pool_proxy::place_market_order", []stx.Argument{
		tx.Object(manager.Address),
		pureU64(tx, parseU64(params.ClientOrderID)),
		pureU8(tx, uint8(params.SelfMatchingOption)),
		pureU64(tx, quantity),
		pureBool(tx, params.IsBid),
		pureBool(tx, params.PayWithDeep),
		tx.Object(c.config.MarginRegistryID),
		tx.Object("0x6"),
	}, []string{base.Type, quote.Type})
}

func (c *PoolProxyContract) CancelOrder(tx *stx.Transaction, marginManagerKey, orderID string) stx.Argument {
	manager := c.config.GetMarginManager(marginManagerKey)
	pool := c.config.GetPool(manager.PoolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	return tx.MoveCall(c.config.MarginPackageID+"::pool_proxy::cancel_order", []stx.Argument{
		tx.Object(manager.Address), pureU128String(tx, orderID), tx.Object(c.config.MarginRegistryID), tx.Object("0x6"),
	}, []string{base.Type, quote.Type})
}

func (c *PoolProxyContract) CancelOrders(tx *stx.Transaction, marginManagerKey string, orderIDs []string) stx.Argument {
	manager := c.config.GetMarginManager(marginManagerKey)
	pool := c.config.GetPool(manager.PoolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	return tx.MoveCall(c.config.MarginPackageID+"::pool_proxy::cancel_orders", []stx.Argument{
		tx.Object(manager.Address), pureVecU128(tx, orderIDs), tx.Object(c.config.MarginRegistryID), tx.Object("0x6"),
	}, []string{base.Type, quote.Type})
}

func (c *PoolProxyContract) CancelAllOrders(tx *stx.Transaction, marginManagerKey string) stx.Argument {
	manager := c.config.GetMarginManager(marginManagerKey)
	pool := c.config.GetPool(manager.PoolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	return tx.MoveCall(c.config.MarginPackageID+"::pool_proxy::cancel_all_orders", []stx.Argument{
		tx.Object(manager.Address), tx.Object(c.config.MarginRegistryID), tx.Object("0x6"),
	}, []string{base.Type, quote.Type})
}

func (c *PoolProxyContract) WithdrawSettledAmounts(tx *stx.Transaction, marginManagerKey string) stx.Argument {
	manager := c.config.GetMarginManager(marginManagerKey)
	pool := c.config.GetPool(manager.PoolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	return tx.MoveCall(c.config.MarginPackageID+"::pool_proxy::withdraw_settled_amounts", []stx.Argument{
		tx.Object(manager.Address), tx.Object(c.config.MarginRegistryID),
	}, []string{base.Type, quote.Type})
}

func (c *PoolProxyContract) Stake(tx *stx.Transaction, marginManagerKey string, stakeAmount float64) stx.Argument {
	manager := c.config.GetMarginManager(marginManagerKey)
	pool := c.config.GetPool(manager.PoolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	stake := uint64(math.Round(stakeAmount * utils.DeepScalar))
	return tx.MoveCall(c.config.MarginPackageID+"::pool_proxy::stake", []stx.Argument{
		tx.Object(manager.Address), pureU64(tx, stake), tx.Object(c.config.MarginRegistryID),
	}, []string{base.Type, quote.Type})
}

func (c *PoolProxyContract) Unstake(tx *stx.Transaction, marginManagerKey string) stx.Argument {
	manager := c.config.GetMarginManager(marginManagerKey)
	pool := c.config.GetPool(manager.PoolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	return tx.MoveCall(c.config.MarginPackageID+"::pool_proxy::unstake", []stx.Argument{
		tx.Object(manager.Address), tx.Object(c.config.MarginRegistryID),
	}, []string{base.Type, quote.Type})
}

func (c *PoolProxyContract) Vote(tx *stx.Transaction, marginManagerKey, proposalID string) stx.Argument {
	manager := c.config.GetMarginManager(marginManagerKey)
	pool := c.config.GetPool(manager.PoolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	return tx.MoveCall(c.config.MarginPackageID+"::pool_proxy::vote", []stx.Argument{
		tx.Object(manager.Address), pureU128String(tx, proposalID), tx.Object(c.config.MarginRegistryID),
	}, []string{base.Type, quote.Type})
}

func (c *PoolProxyContract) ClaimRebate(tx *stx.Transaction, marginManagerKey string) stx.Argument {
	manager := c.config.GetMarginManager(marginManagerKey)
	pool := c.config.GetPool(manager.PoolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	return tx.MoveCall(c.config.MarginPackageID+"::pool_proxy::claim_rebate", []stx.Argument{
		tx.Object(manager.Address), tx.Object(c.config.MarginRegistryID),
	}, []string{base.Type, quote.Type})
}

func (c *MarginTPSLContract) AddConditionalOrder(tx *stx.Transaction, params types.AddConditionalOrderParams) stx.Argument {
	manager := c.config.GetMarginManager(params.MarginManagerKey)
	pool := c.config.GetPool(manager.PoolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	triggerPrice := uint64(math.Round((params.TriggerPrice * utils.FloatScalar * quote.Scalar) / base.Scalar))
	return tx.MoveCall(c.config.MarginPackageID+"::margin_manager::add_conditional_order", []stx.Argument{
		tx.Object(manager.Address), pureU128String(tx, params.ConditionalOrderID), pureBool(tx, params.TriggerBelowPrice), pureU64(tx, triggerPrice), tx.Object(c.config.MarginRegistryID), tx.Object("0x6"),
	}, []string{base.Type, quote.Type})
}

func (c *MarginTPSLContract) CancelAllConditionalOrders(tx *stx.Transaction, marginManagerKey string) stx.Argument {
	manager := c.config.GetMarginManager(marginManagerKey)
	pool := c.config.GetPool(manager.PoolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	return tx.MoveCall(c.config.MarginPackageID+"::margin_manager::cancel_all_conditional_orders", []stx.Argument{
		tx.Object(manager.Address), tx.Object(c.config.MarginRegistryID),
	}, []string{base.Type, quote.Type})
}

func (c *MarginTPSLContract) ConditionalOrderIDs(tx *stx.Transaction, poolKey, marginManagerID string) stx.Argument {
	pool := c.config.GetPool(poolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	return tx.MoveCall(c.config.MarginPackageID+"::margin_manager::conditional_order_ids", []stx.Argument{
		tx.Object(marginManagerID), tx.Object(c.config.MarginRegistryID),
	}, []string{base.Type, quote.Type})
}

func (c *MarginTPSLContract) LowestTriggerAbovePrice(tx *stx.Transaction, poolKey, marginManagerID string) stx.Argument {
	pool := c.config.GetPool(poolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	return tx.MoveCall(c.config.MarginPackageID+"::margin_manager::lowest_trigger_above_price", []stx.Argument{
		tx.Object(marginManagerID), tx.Object(c.config.MarginRegistryID),
	}, []string{base.Type, quote.Type})
}

func (c *MarginTPSLContract) HighestTriggerBelowPrice(tx *stx.Transaction, poolKey, marginManagerID string) stx.Argument {
	pool := c.config.GetPool(poolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	return tx.MoveCall(c.config.MarginPackageID+"::margin_manager::highest_trigger_below_price", []stx.Argument{
		tx.Object(marginManagerID), tx.Object(c.config.MarginRegistryID),
	}, []string{base.Type, quote.Type})
}
