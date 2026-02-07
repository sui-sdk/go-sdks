package deepbookv3

import (
	"github.com/sui-sdks/go-sdks/deepbook_v3/pyth"
	"github.com/sui-sdks/go-sdks/deepbook_v3/transactions"
	"github.com/sui-sdks/go-sdks/deepbook_v3/types"
	"github.com/sui-sdks/go-sdks/deepbook_v3/utils"
)

type DeepBookOptions = Options
type DeepBookClientOptions = ClientOptions

type DeepBookCompatibleClient = CompatibleClient

type DeepBookConfig = utils.DeepBookConfig

var (
	MainnetCoins       = utils.MainnetCoins
	TestnetCoins       = utils.TestnetCoins
	MainnetPools       = utils.MainnetPools
	TestnetPools       = utils.TestnetPools
	MainnetMarginPools = utils.MainnetMarginPools
	TestnetMarginPools = utils.TestnetMarginPools
	MainnetPackageIDs  = utils.MainnetPackageIDs
	TestnetPackageIDs  = utils.TestnetPackageIDs
	MainnetPythConfigs = utils.MainnetPythConfigs
	TestnetPythConfigs = utils.TestnetPythConfigs
	DeepScalar         = utils.DeepScalar
	FloatScalar        = utils.FloatScalar
	GasBudget          = utils.GasBudget
	MaxTimestamp       = utils.MaxTimestamp

	NewBalanceManagerContract   = transactions.NewBalanceManagerContract
	NewDeepBookContract         = transactions.NewDeepBookContract
	NewDeepBookAdminContract    = transactions.NewDeepBookAdminContract
	NewFlashLoanContract        = transactions.NewFlashLoanContract
	NewGovernanceContract       = transactions.NewGovernanceContract
	NewMarginAdminContract      = transactions.NewMarginAdminContract
	NewMarginMaintainerContract = transactions.NewMarginMaintainerContract
	NewMarginManagerContract    = transactions.NewMarginManagerContract
	NewMarginPoolContract       = transactions.NewMarginPoolContract
	NewPoolProxyContract        = transactions.NewPoolProxyContract
	NewMarginTPSLContract       = transactions.NewMarginTPSLContract
	NewSuiPythClient            = pyth.NewSuiPythClient
)

type (
	BalanceManager = types.BalanceManager
	Coin           = types.Coin
	Pool           = types.Pool
	MarginManager  = types.MarginManager
	MarginPool     = types.MarginPool
	Config         = types.Config
)

func Deepbook(opts DeepBookClientOptions) *Client {
	return NewClient(opts)
}
