# Deepbook V3 (Go) Status

Reference package: `/Users/mac/work/sui-sdks/ts-sdks/packages/deepbook-v3`

## Current Scope

Implemented Go package path: `deepbook_v3`

Implemented modules:

- `deepbook_v3/client.go`
- `deepbook_v3/index.go`
- `deepbook_v3/types/types.go`
- `deepbook_v3/utils/*`
- `deepbook_v3/transactions/*`
- `deepbook_v3/pyth/pyth.go`

## Implemented API Surface

Implemented and exported contract classes:

- `BalanceManagerContract`
- `DeepBookContract`
- `DeepBookAdminContract`
- `FlashLoanContract`
- `GovernanceContract`
- `MarginAdminContract`
- `MarginMaintainerContract`
- `MarginManagerContract`
- `MarginPoolContract`
- `MarginRegistryContract`
- `MarginLiquidationsContract`
- `PoolProxyContract`
- `MarginTPSLContract`

Implemented client-level query helpers:

- `CheckManagerBalance`
- `Whitelisted`
- `GetQuoteQuantityOut`
- `GetBaseQuantityOut`
- `GetQuantityOut`
- `MidPrice`
- `GetOrder`
- `GetOrders`
- `AccountOpenOrders`
- `VaultBalances`
- `GetPoolIDByAssets`
- `PoolTradeParams`
- `PoolBookParams`
- `Account`
- `LockedBalance`
- `GetPoolDeepPrice`
- `BalanceManagerReferralOwner`
- `GetPriceInfoObjectAge`
- `GetMarginAccountOrderDetails`
- `GetQuoteQuantityOutInputFee`
- `GetBaseQuantityOutInputFee`
- `GetQuantityOutInputFee`
- `GetBaseQuantityIn`
- `GetQuoteQuantityIn`
- `GetAccountOrderDetails`
- `GetOrderDeepRequired`
- `PoolTradeParamsNext`

## Tests Added

- `deepbook_v3/client_test.go`
- `deepbook_v3/client_methods_test.go`
- `deepbook_v3/transactions/contracts_test.go`
- `deepbook_v3/transactions/encoding_test.go`

Current test focus:

- Client dry-run result parsing for quantity/price/input-fee/order detail paths.
- Contract method target mapping for deepbook, governance, flash-loan, margin contracts.
- Multi-command path validation (e.g. margin manager account order details).
- Parameter encoding assertions for key methods (`u64`, `u128`, `bool`, `vector<u128>`).

Run:

```bash
GOCACHE=$(pwd)/.gocache go test ./deepbook_v3/...
```

## Is it 100% TS parity?

No.

Still missing for full parity:

- Full generated `src/contracts/*` object model parity from ts package.
- Full binary parser parity for all complex return types.
- Comprehensive method-by-method behavior parity tests against ts vectors.
- End-to-end network parity tests for every transaction path.
