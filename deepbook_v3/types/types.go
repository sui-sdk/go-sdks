package types

import "github.com/sui-sdks/go-sdks/sui/transactions"

type BalanceManager struct {
	Address     string
	TradeCap    string
	DepositCap  string
	WithdrawCap string
}

type MarginManager struct {
	Address string
	PoolKey string
}

type Coin struct {
	Address           string
	Type              string
	Scalar            float64
	Feed              string
	CurrencyID        string
	PriceInfoObjectID string
}

type Pool struct {
	Address   string
	BaseCoin  string
	QuoteCoin string
}

type MarginPool struct {
	Address string
	Type    string
}

type OrderType uint8

const (
	OrderTypeNoRestriction OrderType = iota
	OrderTypeImmediateOrCancel
	OrderTypeFillOrKill
	OrderTypePostOnly
)

type SelfMatchingOptions uint8

const (
	SelfMatchingAllowed SelfMatchingOptions = iota
	CancelTaker
	CancelMaker
)

type PlaceLimitOrderParams struct {
	PoolKey            string
	BalanceManagerKey  string
	ClientOrderID      string
	Price              float64
	Quantity           float64
	IsBid              bool
	Expiration         uint64
	OrderType          OrderType
	SelfMatchingOption SelfMatchingOptions
	PayWithDeep        bool
}

type PlaceMarketOrderParams struct {
	PoolKey            string
	BalanceManagerKey  string
	ClientOrderID      string
	Quantity           float64
	IsBid              bool
	SelfMatchingOption SelfMatchingOptions
	PayWithDeep        bool
}

type CanPlaceLimitOrderParams struct {
	PoolKey           string
	BalanceManagerKey string
	Price             float64
	Quantity          float64
	IsBid             bool
	PayWithDeep       bool
	ExpireTimestamp   uint64
}

type CanPlaceMarketOrderParams struct {
	PoolKey           string
	BalanceManagerKey string
	Quantity          float64
	IsBid             bool
	PayWithDeep       bool
}

type PlaceMarginLimitOrderParams struct {
	PoolKey            string
	MarginManagerKey   string
	ClientOrderID      string
	Price              float64
	Quantity           float64
	IsBid              bool
	Expiration         uint64
	OrderType          OrderType
	SelfMatchingOption SelfMatchingOptions
	PayWithDeep        bool
}

type PlaceMarginMarketOrderParams struct {
	PoolKey            string
	MarginManagerKey   string
	ClientOrderID      string
	Quantity           float64
	IsBid              bool
	SelfMatchingOption SelfMatchingOptions
	PayWithDeep        bool
}

type SwapParams struct {
	PoolKey    string
	Amount     float64
	DeepAmount float64
	MinOut     float64
	DeepCoin   transactions.Argument
	BaseCoin   transactions.Argument
	QuoteCoin  transactions.Argument
}

type SwapWithManagerParams struct {
	PoolKey           string
	BalanceManagerKey string
	TradeCap          string
	DepositCap        string
	WithdrawCap       string
	Amount            float64
	MinOut            float64
	BaseCoin          transactions.Argument
	QuoteCoin         transactions.Argument
}

type ProposalParams struct {
	PoolKey           string
	BalanceManagerKey string
	TakerFee          float64
	MakerFee          float64
	StakeRequired     float64
}

type MarginProposalParams struct {
	TakerFee      float64
	MakerFee      float64
	StakeRequired float64
}

type CreatePoolAdminParams struct {
	BaseCoinKey  string
	QuoteCoinKey string
	TickSize     float64
	LotSize      float64
	MinSize      float64
	Whitelisted  bool
	StablePool   bool
}

type CreatePermissionlessPoolParams struct {
	BaseCoinKey  string
	QuoteCoinKey string
	TickSize     float64
	LotSize      float64
	MinSize      float64
	DeepCoin     transactions.Argument
}

type SetEwmaParams struct {
	Alpha              float64
	ZScoreThreshold    float64
	AdditionalTakerFee float64
}

type PoolConfigParams struct {
	MinWithdrawRiskRatio       float64
	MinBorrowRiskRatio         float64
	LiquidationRiskRatio       float64
	TargetLiquidationRiskRatio float64
	UserLiquidationReward      float64
	PoolLiquidationReward      float64
}

type MarginPoolConfigParams struct {
	SupplyCap                float64
	MaxUtilizationRate       float64
	ReferralSpread           float64
	MinBorrow                float64
	RateLimitCapacity        float64
	RateLimitRefillRatePerMs float64
	RateLimitEnabled         bool
}

type InterestConfigParams struct {
	BaseRate           float64
	BaseSlope          float64
	OptimalUtilization float64
	ExcessSlope        float64
}

type PendingLimitOrderParams struct {
	ClientOrderID      string
	OrderType          OrderType
	SelfMatchingOption SelfMatchingOptions
	Price              float64
	Quantity           float64
	IsBid              bool
	PayWithDeep        bool
	ExpireTimestamp    uint64
}

type PendingMarketOrderParams struct {
	ClientOrderID      string
	SelfMatchingOption SelfMatchingOptions
	Quantity           float64
	IsBid              bool
	PayWithDeep        bool
}

type AddConditionalOrderParams struct {
	MarginManagerKey   string
	ConditionalOrderID string
	TriggerBelowPrice  bool
	TriggerPrice       float64
	PendingLimitOrder  *PendingLimitOrderParams
	PendingMarketOrder *PendingMarketOrderParams
}

type DepositParams struct {
	ManagerKey string
	Amount     float64
	Coin       transactions.Argument
}

type DepositDuringInitParams struct {
	Manager  transactions.Argument
	PoolKey  string
	CoinType string
	Amount   float64
	Coin     transactions.Argument
}

type Config struct {
	DeepbookPackageID    string
	RegistryID           string
	DeepTreasuryID       string
	MarginPackageID      string
	MarginRegistryID     string
	LiquidationPackageID string
}
