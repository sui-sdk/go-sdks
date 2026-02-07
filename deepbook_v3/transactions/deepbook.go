package transactions

import (
	"math"

	"github.com/sui-sdks/go-sdks/deepbook_v3/types"
	"github.com/sui-sdks/go-sdks/deepbook_v3/utils"
	stx "github.com/sui-sdks/go-sdks/sui/transactions"
)

type DeepBookContract struct {
	config         *utils.DeepBookConfig
	balanceManager *BalanceManagerContract
}

func NewDeepBookContract(config *utils.DeepBookConfig, balanceManager *BalanceManagerContract) *DeepBookContract {
	return &DeepBookContract{config: config, balanceManager: balanceManager}
}

func (c *DeepBookContract) poolTypes(poolKey string) (types.Coin, types.Coin, types.Pool) {
	pool := c.config.GetPool(poolKey)
	return c.config.GetCoin(pool.BaseCoin), c.config.GetCoin(pool.QuoteCoin), pool
}

func (c *DeepBookContract) poolTarget(fn string) string {
	return c.config.DeepbookPackageID + "::pool::" + fn
}

func (c *DeepBookContract) PlaceLimitOrder(tx *stx.Transaction, params types.PlaceLimitOrderParams) stx.Argument {
	if params.Expiration == 0 {
		params.Expiration = utils.MaxTimestamp
	}
	base, quote, pool := c.poolTypes(params.PoolKey)
	manager := c.config.GetBalanceManager(params.BalanceManagerKey)
	price := uint64(math.Round((params.Price * utils.FloatScalar * quote.Scalar) / base.Scalar))
	qty := uint64(math.Round(params.Quantity * base.Scalar))
	proof := c.balanceManager.GenerateProof(tx, params.BalanceManagerKey)
	tx.SetGasBudgetIfNotSet(utils.GasBudget)
	return tx.MoveCall(c.poolTarget("place_limit_order"), []stx.Argument{
		tx.Object(pool.Address),
		tx.Object(manager.Address),
		proof,
		pureU64(tx, parseU64(params.ClientOrderID)),
		pureU8(tx, uint8(params.OrderType)),
		pureU8(tx, uint8(params.SelfMatchingOption)),
		pureU64(tx, price),
		pureU64(tx, qty),
		pureBool(tx, params.IsBid),
		pureBool(tx, params.PayWithDeep),
		pureU64(tx, params.Expiration),
		tx.Object("0x6"),
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) PlaceMarketOrder(tx *stx.Transaction, params types.PlaceMarketOrderParams) stx.Argument {
	base, quote, pool := c.poolTypes(params.PoolKey)
	manager := c.config.GetBalanceManager(params.BalanceManagerKey)
	qty := uint64(math.Round(params.Quantity * base.Scalar))
	proof := c.balanceManager.GenerateProof(tx, params.BalanceManagerKey)
	tx.SetGasBudgetIfNotSet(utils.GasBudget)
	return tx.MoveCall(c.poolTarget("place_market_order"), []stx.Argument{
		tx.Object(pool.Address),
		tx.Object(manager.Address),
		proof,
		pureU64(tx, parseU64(params.ClientOrderID)),
		pureU8(tx, uint8(params.SelfMatchingOption)),
		pureU64(tx, qty),
		pureBool(tx, params.IsBid),
		pureBool(tx, params.PayWithDeep),
		tx.Object("0x6"),
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) ModifyOrder(tx *stx.Transaction, poolKey, balanceManagerKey, orderID string, newQuantity float64) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	manager := c.config.GetBalanceManager(balanceManagerKey)
	qty := uint64(math.Round(newQuantity * base.Scalar))
	proof := c.balanceManager.GenerateProof(tx, balanceManagerKey)
	return tx.MoveCall(c.poolTarget("modify_order"), []stx.Argument{
		tx.Object(pool.Address), tx.Object(manager.Address), proof, pureU128String(tx, orderID), pureU64(tx, qty), tx.Object("0x6"),
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) CancelOrder(tx *stx.Transaction, poolKey, balanceManagerKey, orderID string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	manager := c.config.GetBalanceManager(balanceManagerKey)
	proof := c.balanceManager.GenerateProof(tx, balanceManagerKey)
	tx.SetGasBudgetIfNotSet(utils.GasBudget)
	return tx.MoveCall(c.poolTarget("cancel_order"), []stx.Argument{
		tx.Object(pool.Address), tx.Object(manager.Address), proof, pureU128String(tx, orderID), tx.Object("0x6"),
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) CancelOrders(tx *stx.Transaction, poolKey, balanceManagerKey string, orderIDs []string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	manager := c.config.GetBalanceManager(balanceManagerKey)
	proof := c.balanceManager.GenerateProof(tx, balanceManagerKey)
	tx.SetGasBudgetIfNotSet(utils.GasBudget)
	return tx.MoveCall(c.poolTarget("cancel_orders"), []stx.Argument{
		tx.Object(pool.Address), tx.Object(manager.Address), proof, pureVecU128(tx, orderIDs), tx.Object("0x6"),
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) CancelAllOrders(tx *stx.Transaction, poolKey, balanceManagerKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	manager := c.config.GetBalanceManager(balanceManagerKey)
	proof := c.balanceManager.GenerateProof(tx, balanceManagerKey)
	tx.SetGasBudgetIfNotSet(utils.GasBudget)
	return tx.MoveCall(c.poolTarget("cancel_all_orders"), []stx.Argument{
		tx.Object(pool.Address), tx.Object(manager.Address), proof, tx.Object("0x6"),
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) WithdrawSettledAmounts(tx *stx.Transaction, poolKey, balanceManagerKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	manager := c.config.GetBalanceManager(balanceManagerKey)
	proof := c.balanceManager.GenerateProof(tx, balanceManagerKey)
	return tx.MoveCall(c.poolTarget("withdraw_settled_amounts"), []stx.Argument{
		tx.Object(pool.Address), tx.Object(manager.Address), proof,
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) AddDeepPricePoint(tx *stx.Transaction, targetPoolKey, referencePoolKey string) stx.Argument {
	base, quote, target := c.poolTypes(targetPoolKey)
	ref := c.config.GetPool(referencePoolKey)
	return tx.MoveCall(c.poolTarget("add_deep_price_point"), []stx.Argument{
		tx.Object(target.Address), tx.Object(ref.Address), tx.Object(c.config.RegistryID), tx.Object("0x6"),
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) ClaimRebates(tx *stx.Transaction, poolKey, balanceManagerKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	manager := c.config.GetBalanceManager(balanceManagerKey)
	proof := c.balanceManager.GenerateProof(tx, balanceManagerKey)
	return tx.MoveCall(c.poolTarget("claim_rebates"), []stx.Argument{
		tx.Object(pool.Address), tx.Object(manager.Address), proof,
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) MintReferral(tx *stx.Transaction, poolKey string, multiplier float64) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	m := uint64(math.Round(multiplier * utils.FloatScalar))
	return tx.MoveCall(c.poolTarget("mint_referral"), []stx.Argument{
		tx.Object(pool.Address), pureU64(tx, m),
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) ClaimPoolReferralRewards(tx *stx.Transaction, poolKey, referral string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	return tx.MoveCall(c.poolTarget("claim_pool_referral_rewards"), []stx.Argument{
		tx.Object(pool.Address), tx.Object(referral),
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) UpdatePoolAllowedVersions(tx *stx.Transaction, poolKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	return tx.MoveCall(c.poolTarget("update_pool_allowed_versions"), []stx.Argument{
		tx.Object(pool.Address), tx.Object(c.config.RegistryID),
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) GetOrder(tx *stx.Transaction, poolKey, orderID string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	return tx.MoveCall(c.poolTarget("get_order"), []stx.Argument{tx.Object(pool.Address), pureU128String(tx, orderID)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) GetOrders(tx *stx.Transaction, poolKey string, orderIDs []string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	return tx.MoveCall(c.poolTarget("get_orders"), []stx.Argument{tx.Object(pool.Address), pureVecU128(tx, orderIDs)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) BurnDeep(tx *stx.Transaction, poolKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	return tx.MoveCall(c.poolTarget("burn_deep"), []stx.Argument{tx.Object(pool.Address), tx.Object(c.config.DeepTreasuryID)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) MidPrice(tx *stx.Transaction, poolKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	return tx.MoveCall(c.poolTarget("mid_price"), []stx.Argument{tx.Object(pool.Address)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) Whitelisted(tx *stx.Transaction, poolKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	return tx.MoveCall(c.poolTarget("whitelisted"), []stx.Argument{tx.Object(pool.Address)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) GetQuoteQuantityOut(tx *stx.Transaction, poolKey string, baseQuantity float64) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	q := uint64(math.Round(baseQuantity * base.Scalar))
	return tx.MoveCall(c.poolTarget("get_quote_quantity_out"), []stx.Argument{tx.Object(pool.Address), pureU64(tx, q), tx.Object("0x6")}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) GetBaseQuantityOut(tx *stx.Transaction, poolKey string, quoteQuantity float64) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	q := uint64(math.Round(quoteQuantity * quote.Scalar))
	return tx.MoveCall(c.poolTarget("get_base_quantity_out"), []stx.Argument{tx.Object(pool.Address), pureU64(tx, q), tx.Object("0x6")}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) GetQuantityOut(tx *stx.Transaction, poolKey string, baseQuantity, quoteQuantity float64) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	b := uint64(math.Round(baseQuantity * base.Scalar))
	q := uint64(math.Round(quoteQuantity * quote.Scalar))
	return tx.MoveCall(c.poolTarget("get_quantity_out"), []stx.Argument{tx.Object(pool.Address), pureU64(tx, b), pureU64(tx, q), tx.Object("0x6")}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) AccountOpenOrders(tx *stx.Transaction, poolKey, managerKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	manager := c.config.GetBalanceManager(managerKey)
	return tx.MoveCall(c.poolTarget("account_open_orders"), []stx.Argument{tx.Object(pool.Address), tx.Object(manager.Address)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) GetLevel2Range(tx *stx.Transaction, poolKey string, priceLow, priceHigh float64, isBid bool) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	pLow := uint64(math.Round((priceLow * utils.FloatScalar * quote.Scalar) / base.Scalar))
	pHigh := uint64(math.Round((priceHigh * utils.FloatScalar * quote.Scalar) / base.Scalar))
	return tx.MoveCall(c.poolTarget("get_level2_range"), []stx.Argument{
		tx.Object(pool.Address), pureU64(tx, pLow), pureU64(tx, pHigh), pureBool(tx, isBid),
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) GetLevel2TicksFromMid(tx *stx.Transaction, poolKey string, tickFromMid uint64) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	return tx.MoveCall(c.poolTarget("get_level2_ticks_from_mid"), []stx.Argument{
		tx.Object(pool.Address), pureU64(tx, tickFromMid),
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) VaultBalances(tx *stx.Transaction, poolKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	return tx.MoveCall(c.poolTarget("vault_balances"), []stx.Argument{tx.Object(pool.Address)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) GetPoolIDByAssets(tx *stx.Transaction, baseType, quoteType string) stx.Argument {
	return tx.MoveCall(c.poolTarget("get_pool_id_by_asset"), []stx.Argument{}, []string{baseType, quoteType})
}

func (c *DeepBookContract) PoolTradeParams(tx *stx.Transaction, poolKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	return tx.MoveCall(c.poolTarget("pool_trade_params"), []stx.Argument{tx.Object(pool.Address)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) PoolBookParams(tx *stx.Transaction, poolKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	return tx.MoveCall(c.poolTarget("pool_book_params"), []stx.Argument{tx.Object(pool.Address)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) Account(tx *stx.Transaction, poolKey, managerKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	manager := c.config.GetBalanceManager(managerKey)
	return tx.MoveCall(c.poolTarget("account"), []stx.Argument{tx.Object(pool.Address), tx.Object(manager.Address)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) LockedBalance(tx *stx.Transaction, poolKey, managerKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	manager := c.config.GetBalanceManager(managerKey)
	return tx.MoveCall(c.poolTarget("locked_balance"), []stx.Argument{tx.Object(pool.Address), tx.Object(manager.Address)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) GetPoolDeepPrice(tx *stx.Transaction, poolKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	return tx.MoveCall(c.poolTarget("get_order_deep_price"), []stx.Argument{tx.Object(pool.Address)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) GetBalanceManagerIDs(tx *stx.Transaction, owner string) stx.Argument {
	return tx.MoveCall(c.config.DeepbookPackageID+"::registry::get_balance_manager_ids", []stx.Argument{tx.PureBytes([]byte(owner)), tx.Object(c.config.RegistryID)}, nil)
}

func (c *DeepBookContract) GetPoolReferralBalances(tx *stx.Transaction, poolKey, referral string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	return tx.MoveCall(c.poolTarget("get_pool_referral_balances"), []stx.Argument{tx.Object(pool.Address), tx.Object(referral)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) PoolReferralMultiplier(tx *stx.Transaction, poolKey, referral string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	return tx.MoveCall(c.poolTarget("pool_referral_multiplier"), []stx.Argument{tx.Object(pool.Address), tx.Object(referral)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) StablePool(tx *stx.Transaction, poolKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	return tx.MoveCall(c.poolTarget("stable_pool"), []stx.Argument{tx.Object(pool.Address)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) RegisteredPool(tx *stx.Transaction, poolKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	return tx.MoveCall(c.poolTarget("registered_pool"), []stx.Argument{tx.Object(pool.Address)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) GetQuoteQuantityOutInputFee(tx *stx.Transaction, poolKey string, baseQuantity float64) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	q := uint64(math.Round(baseQuantity * base.Scalar))
	return tx.MoveCall(c.poolTarget("get_quote_quantity_out_input_fee"), []stx.Argument{tx.Object(pool.Address), pureU64(tx, q), tx.Object("0x6")}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) GetBaseQuantityOutInputFee(tx *stx.Transaction, poolKey string, quoteQuantity float64) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	q := uint64(math.Round(quoteQuantity * quote.Scalar))
	return tx.MoveCall(c.poolTarget("get_base_quantity_out_input_fee"), []stx.Argument{tx.Object(pool.Address), pureU64(tx, q), tx.Object("0x6")}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) GetQuantityOutInputFee(tx *stx.Transaction, poolKey string, baseQuantity, quoteQuantity float64) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	b := uint64(math.Round(baseQuantity * base.Scalar))
	q := uint64(math.Round(quoteQuantity * quote.Scalar))
	return tx.MoveCall(c.poolTarget("get_quantity_out_input_fee"), []stx.Argument{tx.Object(pool.Address), pureU64(tx, b), pureU64(tx, q), tx.Object("0x6")}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) GetBaseQuantityIn(tx *stx.Transaction, poolKey string, targetQuoteQuantity float64, payWithDeep bool) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	q := uint64(math.Round(targetQuoteQuantity * quote.Scalar))
	return tx.MoveCall(c.poolTarget("get_base_quantity_in"), []stx.Argument{tx.Object(pool.Address), pureU64(tx, q), pureBool(tx, payWithDeep), tx.Object("0x6")}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) GetQuoteQuantityIn(tx *stx.Transaction, poolKey string, targetBaseQuantity float64, payWithDeep bool) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	q := uint64(math.Round(targetBaseQuantity * base.Scalar))
	return tx.MoveCall(c.poolTarget("get_quote_quantity_in"), []stx.Argument{tx.Object(pool.Address), pureU64(tx, q), pureBool(tx, payWithDeep), tx.Object("0x6")}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) GetAccountOrderDetails(tx *stx.Transaction, poolKey, managerKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	manager := c.config.GetBalanceManager(managerKey)
	return tx.MoveCall(c.poolTarget("get_account_order_details"), []stx.Argument{tx.Object(pool.Address), tx.Object(manager.Address)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) GetOrderDeepRequired(tx *stx.Transaction, poolKey string, baseQuantity, price float64) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	q := uint64(math.Round(baseQuantity * base.Scalar))
	p := uint64(math.Round((price * utils.FloatScalar * quote.Scalar) / base.Scalar))
	return tx.MoveCall(c.poolTarget("get_order_deep_required"), []stx.Argument{tx.Object(pool.Address), pureU64(tx, q), pureU64(tx, p)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) AccountExists(tx *stx.Transaction, poolKey, managerKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	manager := c.config.GetBalanceManager(managerKey)
	return tx.MoveCall(c.poolTarget("account_exists"), []stx.Argument{tx.Object(pool.Address), tx.Object(manager.Address)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) PoolTradeParamsNext(tx *stx.Transaction, poolKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	return tx.MoveCall(c.poolTarget("pool_trade_params_next"), []stx.Argument{tx.Object(pool.Address)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) Quorum(tx *stx.Transaction, poolKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	return tx.MoveCall(c.poolTarget("quorum"), []stx.Argument{tx.Object(pool.Address)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) PoolID(tx *stx.Transaction, poolKey string) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	return tx.MoveCall(c.poolTarget("id"), []stx.Argument{tx.Object(pool.Address)}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) CanPlaceLimitOrder(tx *stx.Transaction, params types.CanPlaceLimitOrderParams) stx.Argument {
	base, quote, pool := c.poolTypes(params.PoolKey)
	manager := c.config.GetBalanceManager(params.BalanceManagerKey)
	p := uint64(math.Round((params.Price * utils.FloatScalar * quote.Scalar) / base.Scalar))
	q := uint64(math.Round(params.Quantity * base.Scalar))
	return tx.MoveCall(c.poolTarget("can_place_limit_order"), []stx.Argument{
		tx.Object(pool.Address), tx.Object(manager.Address), pureU64(tx, p), pureU64(tx, q), pureBool(tx, params.IsBid), pureBool(tx, params.PayWithDeep), pureU64(tx, params.ExpireTimestamp), tx.Object("0x6"),
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) CanPlaceMarketOrder(tx *stx.Transaction, params types.CanPlaceMarketOrderParams) stx.Argument {
	base, quote, pool := c.poolTypes(params.PoolKey)
	manager := c.config.GetBalanceManager(params.BalanceManagerKey)
	q := uint64(math.Round(params.Quantity * base.Scalar))
	return tx.MoveCall(c.poolTarget("can_place_market_order"), []stx.Argument{
		tx.Object(pool.Address), tx.Object(manager.Address), pureU64(tx, q), pureBool(tx, params.IsBid), pureBool(tx, params.PayWithDeep), tx.Object("0x6"),
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookContract) CheckMarketOrderParams(tx *stx.Transaction, poolKey string, quantity float64) stx.Argument {
	base, quote, pool := c.poolTypes(poolKey)
	q := uint64(math.Round(quantity * base.Scalar))
	return tx.MoveCall(c.poolTarget("check_market_order_params"), []stx.Argument{tx.Object(pool.Address), pureU64(tx, q)}, []string{base.Type, quote.Type})
}
