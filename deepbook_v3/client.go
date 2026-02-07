package deepbookv3

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/sui-sdks/go-sdks/bcs"
	"github.com/sui-sdks/go-sdks/deepbook_v3/pyth"
	"github.com/sui-sdks/go-sdks/deepbook_v3/transactions"
	"github.com/sui-sdks/go-sdks/deepbook_v3/types"
	"github.com/sui-sdks/go-sdks/deepbook_v3/utils"
	stx "github.com/sui-sdks/go-sdks/sui/transactions"
	suiutils "github.com/sui-sdks/go-sdks/sui/utils"
)

type CompatibleClient interface {
	Call(ctx context.Context, method string, params []any, out any) error
	Network() string
}

type Options struct {
	Address             string
	BalanceManagers     map[string]types.BalanceManager
	MarginManagers      map[string]types.MarginManager
	Coins               utils.CoinMap
	Pools               utils.PoolMap
	AdminCap            string
	MarginAdminCap      string
	MarginMaintainerCap string
}

type ClientOptions struct {
	Client  CompatibleClient
	Network string
	Options
}

type Client struct {
	client             CompatibleClient
	config             *utils.DeepBookConfig
	Address            string
	BalanceManager     *transactions.BalanceManagerContract
	DeepBook           *transactions.DeepBookContract
	DeepBookAdmin      *transactions.DeepBookAdminContract
	FlashLoans         *transactions.FlashLoanContract
	Governance         *transactions.GovernanceContract
	MarginAdmin        *transactions.MarginAdminContract
	MarginMaintainer   *transactions.MarginMaintainerContract
	MarginPool         *transactions.MarginPoolContract
	MarginManager      *transactions.MarginManagerContract
	MarginRegistry     *transactions.MarginRegistryContract
	MarginLiquidations *transactions.MarginLiquidationsContract
	PoolProxy          *transactions.PoolProxyContract
	MarginTPSL         *transactions.MarginTPSLContract
}

func NewClient(opts ClientOptions) *Client {
	address := suiutils.NormalizeSuiAddress(opts.Address)
	config := utils.NewDeepBookConfig(utils.ConfigOptions{
		Address:             address,
		Network:             opts.Network,
		BalanceManagers:     opts.BalanceManagers,
		MarginManagers:      opts.MarginManagers,
		Coins:               opts.Coins,
		Pools:               opts.Pools,
		AdminCap:            opts.AdminCap,
		MarginAdminCap:      opts.MarginAdminCap,
		MarginMaintainerCap: opts.MarginMaintainerCap,
	})
	balanceManager := transactions.NewBalanceManagerContract(config)
	return &Client{
		client:             opts.Client,
		config:             config,
		Address:            address,
		BalanceManager:     balanceManager,
		DeepBook:           transactions.NewDeepBookContract(config, balanceManager),
		DeepBookAdmin:      transactions.NewDeepBookAdminContract(config),
		FlashLoans:         transactions.NewFlashLoanContract(config),
		Governance:         transactions.NewGovernanceContract(config, balanceManager),
		MarginAdmin:        transactions.NewMarginAdminContract(config),
		MarginMaintainer:   transactions.NewMarginMaintainerContract(config),
		MarginPool:         transactions.NewMarginPoolContract(config),
		MarginManager:      transactions.NewMarginManagerContract(config),
		MarginRegistry:     transactions.NewMarginRegistryContract(config),
		MarginLiquidations: transactions.NewMarginLiquidationsContract(config),
		PoolProxy:          transactions.NewPoolProxyContract(config),
		MarginTPSL:         transactions.NewMarginTPSLContract(config),
	}
}

func (c *Client) Config() *utils.DeepBookConfig { return c.config }

func (c *Client) CheckManagerBalance(ctx context.Context, managerKey, coinKey string) (map[string]any, error) {
	tx := stx.NewTransaction()
	c.BalanceManager.CheckManagerBalance(tx, managerKey, coinKey)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return nil, err
	}
	v, err := readU64(res, 0, 0)
	if err != nil {
		return nil, err
	}
	coin := c.config.GetCoin(coinKey)
	return map[string]any{
		"coinType": coin.Type,
		"balance":  float64(v) / coin.Scalar,
	}, nil
}

func (c *Client) Whitelisted(ctx context.Context, poolKey string) (bool, error) {
	tx := stx.NewTransaction()
	c.DeepBook.Whitelisted(tx, poolKey)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return false, err
	}
	b, err := extractReturnBCS(res, 0, 0)
	if err != nil {
		return false, err
	}
	return len(b) > 0 && b[0] == 1, nil
}

func (c *Client) GetQuoteQuantityOut(ctx context.Context, poolKey string, baseQuantity float64) (map[string]any, error) {
	tx := stx.NewTransaction()
	pool := c.config.GetPool(poolKey)
	baseScalar := c.config.GetCoin(pool.BaseCoin).Scalar
	quoteScalar := c.config.GetCoin(pool.QuoteCoin).Scalar
	c.DeepBook.GetQuoteQuantityOut(tx, poolKey, baseQuantity)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return nil, err
	}
	baseOut, _ := readU64(res, 0, 0)
	quoteOut, _ := readU64(res, 0, 1)
	deepRequired, _ := readU64(res, 0, 2)
	return map[string]any{
		"baseQuantity": baseQuantity,
		"baseOut":      float64(baseOut) / baseScalar,
		"quoteOut":     float64(quoteOut) / quoteScalar,
		"deepRequired": float64(deepRequired) / utils.DeepScalar,
	}, nil
}

func (c *Client) GetBaseQuantityOut(ctx context.Context, poolKey string, quoteQuantity float64) (map[string]any, error) {
	tx := stx.NewTransaction()
	pool := c.config.GetPool(poolKey)
	baseScalar := c.config.GetCoin(pool.BaseCoin).Scalar
	quoteScalar := c.config.GetCoin(pool.QuoteCoin).Scalar
	c.DeepBook.GetBaseQuantityOut(tx, poolKey, quoteQuantity)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return nil, err
	}
	baseOut, _ := readU64(res, 0, 0)
	quoteOut, _ := readU64(res, 0, 1)
	deepRequired, _ := readU64(res, 0, 2)
	return map[string]any{
		"quoteQuantity": quoteQuantity,
		"baseOut":       float64(baseOut) / baseScalar,
		"quoteOut":      float64(quoteOut) / quoteScalar,
		"deepRequired":  float64(deepRequired) / utils.DeepScalar,
	}, nil
}

func (c *Client) GetQuantityOut(ctx context.Context, poolKey string, baseQuantity, quoteQuantity float64) (map[string]any, error) {
	tx := stx.NewTransaction()
	pool := c.config.GetPool(poolKey)
	baseScalar := c.config.GetCoin(pool.BaseCoin).Scalar
	quoteScalar := c.config.GetCoin(pool.QuoteCoin).Scalar
	c.DeepBook.GetQuantityOut(tx, poolKey, baseQuantity, quoteQuantity)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return nil, err
	}
	baseOut, _ := readU64(res, 0, 0)
	quoteOut, _ := readU64(res, 0, 1)
	deepRequired, _ := readU64(res, 0, 2)
	return map[string]any{
		"baseQuantity":  baseQuantity,
		"quoteQuantity": quoteQuantity,
		"baseOut":       float64(baseOut) / baseScalar,
		"quoteOut":      float64(quoteOut) / quoteScalar,
		"deepRequired":  float64(deepRequired) / utils.DeepScalar,
	}, nil
}

func (c *Client) MidPrice(ctx context.Context, poolKey string) (float64, error) {
	tx := stx.NewTransaction()
	pool := c.config.GetPool(poolKey)
	base := c.config.GetCoin(pool.BaseCoin)
	quote := c.config.GetCoin(pool.QuoteCoin)
	c.DeepBook.MidPrice(tx, poolKey)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return 0, err
	}
	v, err := readU64(res, 0, 0)
	if err != nil {
		return 0, err
	}
	return float64(v) * base.Scalar / (utils.FloatScalar * quote.Scalar), nil
}

func (c *Client) GetOrder(ctx context.Context, poolKey, orderID string) (string, error) {
	tx := stx.NewTransaction()
	c.DeepBook.GetOrder(tx, poolKey, orderID)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return "", err
	}
	return readReturnBCSBase64(res, 0, 0)
}

func (c *Client) GetOrders(ctx context.Context, poolKey string, orderIDs []string) (string, error) {
	tx := stx.NewTransaction()
	c.DeepBook.GetOrders(tx, poolKey, orderIDs)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return "", err
	}
	return readReturnBCSBase64(res, 0, 0)
}

func (c *Client) AccountOpenOrders(ctx context.Context, poolKey, managerKey string) (string, error) {
	tx := stx.NewTransaction()
	c.DeepBook.AccountOpenOrders(tx, poolKey, managerKey)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return "", err
	}
	return readReturnBCSBase64(res, 0, 0)
}

func (c *Client) VaultBalances(ctx context.Context, poolKey string) (string, error) {
	tx := stx.NewTransaction()
	c.DeepBook.VaultBalances(tx, poolKey)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return "", err
	}
	return readReturnBCSBase64(res, 0, 0)
}

func (c *Client) GetPoolIDByAssets(ctx context.Context, baseType, quoteType string) (string, error) {
	tx := stx.NewTransaction()
	c.DeepBook.GetPoolIDByAssets(tx, baseType, quoteType)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return "", err
	}
	return readReturnBCSBase64(res, 0, 0)
}

func (c *Client) PoolTradeParams(ctx context.Context, poolKey string) (string, error) {
	tx := stx.NewTransaction()
	c.DeepBook.PoolTradeParams(tx, poolKey)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return "", err
	}
	return readReturnBCSBase64(res, 0, 0)
}

func (c *Client) PoolBookParams(ctx context.Context, poolKey string) (string, error) {
	tx := stx.NewTransaction()
	c.DeepBook.PoolBookParams(tx, poolKey)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return "", err
	}
	return readReturnBCSBase64(res, 0, 0)
}

func (c *Client) Account(ctx context.Context, poolKey, managerKey string) (string, error) {
	tx := stx.NewTransaction()
	c.DeepBook.Account(tx, poolKey, managerKey)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return "", err
	}
	return readReturnBCSBase64(res, 0, 0)
}

func (c *Client) LockedBalance(ctx context.Context, poolKey, balanceManagerKey string) (string, error) {
	tx := stx.NewTransaction()
	c.DeepBook.LockedBalance(tx, poolKey, balanceManagerKey)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return "", err
	}
	return readReturnBCSBase64(res, 0, 0)
}

func (c *Client) GetPoolDeepPrice(ctx context.Context, poolKey string) (string, error) {
	tx := stx.NewTransaction()
	c.DeepBook.GetPoolDeepPrice(tx, poolKey)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return "", err
	}
	return readReturnBCSBase64(res, 0, 0)
}

func (c *Client) BalanceManagerReferralOwner(ctx context.Context, referral string) (string, error) {
	tx := stx.NewTransaction()
	c.BalanceManager.BalanceManagerReferralOwner(tx, referral)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return "", err
	}
	return readReturnBCSBase64(res, 0, 0)
}

func (c *Client) GetPriceInfoObjectAge(_ context.Context, coinKey string) (int64, error) {
	coin := c.config.GetCoin(coinKey)
	if coin.PriceInfoObjectID == "" {
		return -1, nil
	}
	return time.Now().UnixMilli(), nil
}

func (c *Client) GetMarginAccountOrderDetails(ctx context.Context, marginManagerKey string) (string, error) {
	tx := stx.NewTransaction()
	c.MarginManager.GetMarginAccountOrderDetails(tx, marginManagerKey)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return "", err
	}
	return readReturnBCSBase64(res, 1, 0)
}

func (c *Client) GetQuoteQuantityOutInputFee(ctx context.Context, poolKey string, baseQuantity float64) (map[string]any, error) {
	tx := stx.NewTransaction()
	c.DeepBook.GetQuoteQuantityOutInputFee(tx, poolKey, baseQuantity)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return nil, err
	}
	v, _ := readU64(res, 0, 0)
	return map[string]any{"baseQuantity": baseQuantity, "result": v}, nil
}

func (c *Client) GetBaseQuantityOutInputFee(ctx context.Context, poolKey string, quoteQuantity float64) (map[string]any, error) {
	tx := stx.NewTransaction()
	c.DeepBook.GetBaseQuantityOutInputFee(tx, poolKey, quoteQuantity)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return nil, err
	}
	v, _ := readU64(res, 0, 0)
	return map[string]any{"quoteQuantity": quoteQuantity, "result": v}, nil
}

func (c *Client) GetQuantityOutInputFee(ctx context.Context, poolKey string, baseQuantity, quoteQuantity float64) (map[string]any, error) {
	tx := stx.NewTransaction()
	c.DeepBook.GetQuantityOutInputFee(tx, poolKey, baseQuantity, quoteQuantity)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return nil, err
	}
	v, _ := readU64(res, 0, 0)
	return map[string]any{"baseQuantity": baseQuantity, "quoteQuantity": quoteQuantity, "result": v}, nil
}

func (c *Client) GetBaseQuantityIn(ctx context.Context, poolKey string, targetQuoteQuantity float64, payWithDeep bool) (uint64, error) {
	tx := stx.NewTransaction()
	c.DeepBook.GetBaseQuantityIn(tx, poolKey, targetQuoteQuantity, payWithDeep)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return 0, err
	}
	return readU64(res, 0, 0)
}

func (c *Client) GetQuoteQuantityIn(ctx context.Context, poolKey string, targetBaseQuantity float64, payWithDeep bool) (uint64, error) {
	tx := stx.NewTransaction()
	c.DeepBook.GetQuoteQuantityIn(tx, poolKey, targetBaseQuantity, payWithDeep)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return 0, err
	}
	return readU64(res, 0, 0)
}

func (c *Client) GetAccountOrderDetails(ctx context.Context, poolKey, managerKey string) (string, error) {
	tx := stx.NewTransaction()
	c.DeepBook.GetAccountOrderDetails(tx, poolKey, managerKey)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return "", err
	}
	return readReturnBCSBase64(res, 0, 0)
}

func (c *Client) GetOrderDeepRequired(ctx context.Context, poolKey string, baseQuantity, price float64) (uint64, error) {
	tx := stx.NewTransaction()
	c.DeepBook.GetOrderDeepRequired(tx, poolKey, baseQuantity, price)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return 0, err
	}
	return readU64(res, 0, 0)
}

func (c *Client) PoolTradeParamsNext(ctx context.Context, poolKey string) (string, error) {
	tx := stx.NewTransaction()
	c.DeepBook.PoolTradeParamsNext(tx, poolKey)
	res, err := c.simulate(ctx, tx)
	if err != nil {
		return "", err
	}
	return readReturnBCSBase64(res, 0, 0)
}

func (c *Client) GetPythClient(pythStateID, wormholeStateID string) *pyth.SuiPythClient {
	return pyth.NewSuiPythClient(c.client, pythStateID, wormholeStateID)
}

func (c *Client) simulate(ctx context.Context, tx *stx.Transaction) (map[string]any, error) {
	built, err := tx.BuildBase64()
	if err != nil {
		return nil, err
	}
	var out map[string]any
	err = c.client.Call(ctx, "sui_dryRunTransactionBlock", []any{built}, &out)
	return out, err
}

func readU64(res map[string]any, cmd, ret int) (uint64, error) {
	bytes, err := extractReturnBCS(res, cmd, ret)
	if err != nil {
		return 0, err
	}
	reader := bcs.NewReader(bytes)
	return reader.Read64()
}

func readReturnBCSBase64(res map[string]any, cmd, ret int) (string, error) {
	bytes, err := extractReturnBCS(res, cmd, ret)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func extractReturnBCS(res map[string]any, cmd, ret int) ([]byte, error) {
	commandResults, ok := res["commandResults"].([]any)
	if !ok || len(commandResults) <= cmd {
		return nil, fmt.Errorf("missing commandResults[%d]", cmd)
	}
	cr, ok := commandResults[cmd].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid command result")
	}
	returnValues, ok := cr["returnValues"].([]any)
	if !ok || len(returnValues) <= ret {
		return nil, fmt.Errorf("missing returnValues[%d]", ret)
	}
	rv, ok := returnValues[ret].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid return value")
	}
	b64, ok := rv["bcs"].(string)
	if !ok {
		return nil, fmt.Errorf("missing bcs")
	}
	return base64.StdEncoding.DecodeString(b64)
}
