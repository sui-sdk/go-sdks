package utils

import (
	"github.com/sui-sdks/go-sdks/deepbook_v3/types"
	suiutils "github.com/sui-sdks/go-sdks/sui/utils"
)

const (
	FloatScalar             = 1_000_000_000.0
	DeepScalar              = 1_000_000.0
	MaxTimestamp      uint64 = 1_844_674_407_370_955_161
	PriceInfoObjectMaxAgeMs = 15_000
	GasBudget         int64  = 250_000_000
	PoolCreationFeeDeep      = 500_000_000.0
)

type DeepBookConfig struct {
	Address             string
	Network             string
	BalanceManagers     map[string]types.BalanceManager
	MarginManagers      map[string]types.MarginManager
	Coins               CoinMap
	Pools               PoolMap
	MarginPools         MarginPoolMap
	AdminCap            string
	MarginAdminCap      string
	MarginMaintainerCap string

	DeepbookPackageID    string
	RegistryID           string
	DeepTreasuryID       string
	MarginPackageID      string
	MarginRegistryID     string
	LiquidationPackageID string
}

type ConfigOptions struct {
	Address             string
	Network             string
	BalanceManagers     map[string]types.BalanceManager
	MarginManagers      map[string]types.MarginManager
	Coins               CoinMap
	Pools               PoolMap
	MarginPools         MarginPoolMap
	AdminCap            string
	MarginAdminCap      string
	MarginMaintainerCap string
}

func NewDeepBookConfig(opts ConfigOptions) *DeepBookConfig {
	address := suiutils.NormalizeSuiAddress(opts.Address)

	packageIDs := MainnetPackageIDs
	coins := MainnetCoins
	pools := MainnetPools
	marginPools := MainnetMarginPools

	if opts.Network == "testnet" {
		packageIDs = TestnetPackageIDs
		coins = TestnetCoins
		pools = TestnetPools
		marginPools = TestnetMarginPools
	}

	if opts.Coins != nil {
		coins = opts.Coins
	}
	if opts.Pools != nil {
		pools = opts.Pools
	}
	if opts.MarginPools != nil {
		marginPools = opts.MarginPools
	}
	if opts.BalanceManagers == nil {
		opts.BalanceManagers = map[string]types.BalanceManager{}
	}
	if opts.MarginManagers == nil {
		opts.MarginManagers = map[string]types.MarginManager{}
	}

	return &DeepBookConfig{
		Address:             address,
		Network:             opts.Network,
		BalanceManagers:     opts.BalanceManagers,
		MarginManagers:      opts.MarginManagers,
		Coins:               coins,
		Pools:               pools,
		MarginPools:         marginPools,
		AdminCap:            opts.AdminCap,
		MarginAdminCap:      opts.MarginAdminCap,
		MarginMaintainerCap: opts.MarginMaintainerCap,
		DeepbookPackageID:   packageIDs.DeepbookPackageID,
		RegistryID:          packageIDs.RegistryID,
		DeepTreasuryID:      packageIDs.DeepTreasuryID,
		MarginPackageID:     packageIDs.MarginPackageID,
		MarginRegistryID:    packageIDs.MarginRegistryID,
		LiquidationPackageID: packageIDs.LiquidationPackageID,
	}
}

func (c *DeepBookConfig) GetCoin(key string) types.Coin {
	coin, ok := c.Coins[key]
	if !ok {
		panic(&ResourceNotFoundError{DeepBookError{Msg: ErrorMessages.CoinNotFound(key)}})
	}
	return coin
}

func (c *DeepBookConfig) GetPool(key string) types.Pool {
	pool, ok := c.Pools[key]
	if !ok {
		panic(&ResourceNotFoundError{DeepBookError{Msg: ErrorMessages.PoolNotFound(key)}})
	}
	return pool
}

func (c *DeepBookConfig) GetMarginPool(key string) types.MarginPool {
	pool, ok := c.MarginPools[key]
	if !ok {
		panic(&ResourceNotFoundError{DeepBookError{Msg: ErrorMessages.MarginPoolNotFound(key)}})
	}
	return pool
}

func (c *DeepBookConfig) GetBalanceManager(key string) types.BalanceManager {
	m, ok := c.BalanceManagers[key]
	if !ok {
		panic(&ResourceNotFoundError{DeepBookError{Msg: ErrorMessages.BalanceManagerNotFound(key)}})
	}
	return m
}

func (c *DeepBookConfig) GetMarginManager(key string) types.MarginManager {
	m, ok := c.MarginManagers[key]
	if !ok {
		panic(&ResourceNotFoundError{DeepBookError{Msg: ErrorMessages.MarginManagerNotFound(key)}})
	}
	return m
}
