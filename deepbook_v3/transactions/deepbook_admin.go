package transactions

import (
	"math"

	"github.com/sui-sdks/go-sdks/deepbook_v3/types"
	"github.com/sui-sdks/go-sdks/deepbook_v3/utils"
	stx "github.com/sui-sdks/go-sdks/sui/transactions"
)

type DeepBookAdminContract struct{ config *utils.DeepBookConfig }

func NewDeepBookAdminContract(config *utils.DeepBookConfig) *DeepBookAdminContract {
	return &DeepBookAdminContract{config: config}
}

func (c *DeepBookAdminContract) adminCap() string {
	if c.config.AdminCap == "" {
		panic(&utils.ConfigurationError{DeepBookError: utils.DeepBookError{Msg: utils.ErrorMessages.AdminCapNotSet}})
	}
	return c.config.AdminCap
}

func (c *DeepBookAdminContract) CreatePoolAdmin(tx *stx.Transaction, params types.CreatePoolAdminParams) stx.Argument {
	base := c.config.GetCoin(params.BaseCoinKey)
	quote := c.config.GetCoin(params.QuoteCoinKey)
	adjustedTickSize := uint64(math.Round((params.TickSize * utils.FloatScalar * quote.Scalar) / base.Scalar))
	adjustedLotSize := uint64(math.Round(params.LotSize * base.Scalar))
	adjustedMinSize := uint64(math.Round(params.MinSize * base.Scalar))
	return tx.MoveCall(c.config.DeepbookPackageID+"::pool::create_pool_admin", []stx.Argument{
		tx.Object(c.config.RegistryID),
		pureU64(tx, adjustedTickSize),
		pureU64(tx, adjustedLotSize),
		pureU64(tx, adjustedMinSize),
		pureBool(tx, params.Whitelisted),
		pureBool(tx, params.StablePool),
		tx.Object(c.adminCap()),
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookAdminContract) UnregisterPoolAdmin(tx *stx.Transaction, poolKey string) stx.Argument {
	pool := c.config.GetPool(poolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	return tx.MoveCall(c.config.DeepbookPackageID+"::pool::unregister_pool_admin", []stx.Argument{
		tx.Object(pool.Address), tx.Object(c.adminCap()),
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookAdminContract) UpdateAllowedVersions(tx *stx.Transaction, poolKey string) stx.Argument {
	pool := c.config.GetPool(poolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	return tx.MoveCall(c.config.DeepbookPackageID+"::pool::update_allowed_versions", []stx.Argument{
		tx.Object(pool.Address), tx.Object(c.config.RegistryID), tx.Object(c.adminCap()),
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookAdminContract) EnableVersion(tx *stx.Transaction, version uint64) stx.Argument {
	return tx.MoveCall(c.config.DeepbookPackageID+"::registry::enable_version", []stx.Argument{
		tx.Object(c.config.RegistryID), pureU64(tx, version), tx.Object(c.adminCap()),
	}, nil)
}

func (c *DeepBookAdminContract) DisableVersion(tx *stx.Transaction, version uint64) stx.Argument {
	return tx.MoveCall(c.config.DeepbookPackageID+"::registry::disable_version", []stx.Argument{
		tx.Object(c.config.RegistryID), pureU64(tx, version), tx.Object(c.adminCap()),
	}, nil)
}

func (c *DeepBookAdminContract) SetTreasuryAddress(tx *stx.Transaction, treasuryAddress string) stx.Argument {
	return tx.MoveCall(c.config.DeepbookPackageID+"::registry::set_treasury_address", []stx.Argument{
		tx.Object(c.config.RegistryID), tx.PureBytes([]byte(treasuryAddress)), tx.Object(c.adminCap()),
	}, nil)
}

func (c *DeepBookAdminContract) AddStableCoin(tx *stx.Transaction, stableCoinKey string) stx.Argument {
	coin := c.config.GetCoin(stableCoinKey)
	return tx.MoveCall(c.config.DeepbookPackageID+"::registry::add_stablecoin", []stx.Argument{
		tx.Object(c.config.RegistryID), tx.Object(c.adminCap()),
	}, []string{coin.Type})
}

func (c *DeepBookAdminContract) RemoveStableCoin(tx *stx.Transaction, stableCoinKey string) stx.Argument {
	coin := c.config.GetCoin(stableCoinKey)
	return tx.MoveCall(c.config.DeepbookPackageID+"::registry::remove_stablecoin", []stx.Argument{
		tx.Object(c.config.RegistryID), tx.Object(c.adminCap()),
	}, []string{coin.Type})
}

func (c *DeepBookAdminContract) AdjustTickSize(tx *stx.Transaction, poolKey string, newTickSize float64) stx.Argument {
	pool := c.config.GetPool(poolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	tick := uint64(math.Round((newTickSize * utils.FloatScalar * quote.Scalar) / base.Scalar))
	return tx.MoveCall(c.config.DeepbookPackageID+"::pool::adjust_tick_size_admin", []stx.Argument{
		tx.Object(pool.Address), pureU64(tx, tick), tx.Object(c.adminCap()),
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookAdminContract) InitBalanceManagerMap(tx *stx.Transaction) stx.Argument {
	return tx.MoveCall(c.config.DeepbookPackageID+"::registry::init_balance_manager_map", []stx.Argument{
		tx.Object(c.config.RegistryID), tx.Object(c.adminCap()),
	}, nil)
}

func (c *DeepBookAdminContract) SetEwmaParams(tx *stx.Transaction, poolKey string, params types.SetEwmaParams) stx.Argument {
	pool := c.config.GetPool(poolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	alpha := uint64(math.Round(params.Alpha * utils.FloatScalar))
	zScoreThreshold := uint64(math.Round(params.ZScoreThreshold * utils.FloatScalar))
	additionalTakerFee := uint64(math.Round(params.AdditionalTakerFee * utils.FloatScalar))
	return tx.MoveCall(c.config.DeepbookPackageID+"::pool::set_ewma_params", []stx.Argument{
		tx.Object(pool.Address), pureU64(tx, alpha), pureU64(tx, zScoreThreshold), pureU64(tx, additionalTakerFee), tx.Object(c.adminCap()),
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookAdminContract) EnableEwmaState(tx *stx.Transaction, poolKey string, enable bool) stx.Argument {
	pool := c.config.GetPool(poolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	return tx.MoveCall(c.config.DeepbookPackageID+"::pool::enable_ewma_state", []stx.Argument{
		tx.Object(pool.Address), pureBool(tx, enable), tx.Object(c.adminCap()),
	}, []string{base.Type, quote.Type})
}

func (c *DeepBookAdminContract) AuthorizeMarginApp(tx *stx.Transaction) stx.Argument {
	return tx.MoveCall(c.config.DeepbookPackageID+"::registry::authorize_app", []stx.Argument{
		tx.Object(c.config.RegistryID), tx.Object(c.adminCap()),
	}, nil)
}

func (c *DeepBookAdminContract) DeauthorizeMarginApp(tx *stx.Transaction) stx.Argument {
	return tx.MoveCall(c.config.DeepbookPackageID+"::registry::deauthorize_app", []stx.Argument{
		tx.Object(c.config.RegistryID), tx.Object(c.adminCap()),
	}, nil)
}
