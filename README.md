# Go Sui SDKs (`go-sdks`)

Go implementation of selected packages from `ts-sdks`, with API naming adapted to idiomatic Go:

- `bcs` (binary canonical serialization)
- `sui/jsonrpc`
- `sui/grpc`
- `sui/graphql`
- `sui/faucet`
- `sui/cryptography`
- `sui/keypairs/ed25519`
- `sui/keypairs/secp256k1`
- `sui/keypairs/secp256r1`
- `sui/multisig`
- `sui/zklogin`
- `sui/verify`
- `sui/transactions`
- `walrus`
- `seal`
- `deepbook_v3`

This repository currently focuses on core usable capabilities first, then parity expansion.

## Quick Start

```bash
go test ./...
```

```go
package main

import (
	"context"
	"fmt"

	"github.com/sui-sdks/go-sdks/sui/jsonrpc"
)

func main() {
	c, _ := jsonrpc.NewClient(jsonrpc.ClientOptions{
		Network: "testnet",
	})

	version, _ := c.GetRPCAPIVersion(context.Background())
	fmt.Println(version)
}
```

## Package Overview

### `bcs`

- ULEB encode/decode
- Reader/Writer primitives (`u8/u16/u32/u64/u128/u256`)
- Primitive BCS types (`bool`, `string`, `bytes`, `byteVector`)
- Composite types (`vector`, `fixedArray`, `option`, `tuple`, `struct`, `enum`, `map`, `lazy`)
- Encodings (`hex`, `base64`, `base58`)

### `sui/jsonrpc`

- HTTP JSON-RPC transport
- Client constructor by `url` or network fullnode default
- Generic `Call` method
- Common methods: `GetCoins`, `GetAllCoins`, `GetBalance`, `GetAllBalances`,
  `GetCoinMetadata`, `GetTotalSupply`, `GetObject`, `MultiGetObjects`,
  `GetTransactionBlock`, `ExecuteTransactionBlock`, `GetReferenceGasPrice`,
  `QueryTransactionBlocks`, `GetRPCAPIVersion`

### `sui/graphql`

- HTTP GraphQL client
- `Query` and named `Execute` methods

### `sui/faucet`

- `RequestSuiFromFaucetV2`
- `GetFaucetHost`
- rate-limit error type (`FaucetRateLimitError`)

### `sui/cryptography` + `sui/keypairs/*`

- signature scheme constants + flags
- intent message domain separation
- serialized signature encode/decode
- public key base APIs
- keypair interfaces and helpers
- keypairs:
  - `ed25519`
  - `secp256k1` (current implementation uses stdlib ECDSA compatibility mode)
  - `secp256r1`

### `sui/transactions`

- transaction input helpers (`Inputs`)
- transaction command builders (`TransactionCommands`)
- transaction builder (`Transaction`)
- argument helpers (`Arguments`)
- build/serialize/restore flows (JSON/base64 JSON)
- transaction resolve plugin pipeline
- core resolver (input/object/gas resolution baseline)
- executors:
  - caching executor
  - serial executor
  - parallel executor
  - serial/parallel queue primitives

### `sui/grpc`

- `SuiGrpcClient`/`GrpcCoreClient` style API surface
- current implementation uses unified core transport and method mapping, exposed under grpc package path

### `sui/multisig`

- multisig public-key model
- multisig signature serialization/parsing
- multisig signer wrapper

### `sui/zklogin`

- JWT decode helpers
- nonce/randomness helpers
- zklogin signature encode/decode
- zklogin address helpers

### `sui/verify`

- verify helpers for raw signature / personal message / transaction intent

### `walrus`

- package constants for mainnet/testnet IDs
- storage node client for metadata/status/sliver/confirmation endpoints
- high-level client methods:
  - `ReadBlob` (baseline implementation)
  - `GetBlobMetadata`, `GetBlobStatus`, `GetSliver`
  - `WriteBlobMetadata`, `WriteSliver`

### `seal`

- session key generation/export/import
- key server metadata and fetch flow helpers
- Shamir split/combine
- DEM implementations:
  - `AesGcm256`
  - `Hmac256Ctr`
- encrypt/decrypt core flow
- key derivation and supporting utils/errors

### `deepbook_v3`

- `DeepBookClient` + config/options
- transaction contracts:
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
- pyth client adapter (`SuiPythClient`)
- mainnet/testnet constants and package IDs

## Tests

Current tests cover:

- `bcs`: ULEB/base58/struct-enum-vector-map roundtrip
- `sui/jsonrpc`: transport + client method path
- `sui/graphql`: query + named execute
- `sui/faucet`: success + 429 handling
- `sui/cryptography`: key encode/decode + signature serialization
- `sui/keypairs/*`: sign/verify for ed25519/secp256k1/secp256r1
- `sui/transactions`: build + serialize + restore
- `sui/transactions`: resolver + executor flows (caching/serial/parallel)
- `sui/grpc`: grpc package surface + core method coverage
- `sui/multisig`: serialize/parse/sign baseline flow
- `sui/zklogin`: jwt/nonce/signature/address helper flow
- `sui/verify`: verification helper flow
- `walrus`: read/write storage-node interaction
- `seal`: shamir + encrypt/decrypt + session key export/import
- `deepbook_v3`: client queries + contract target mapping coverage

Run:

```bash
go test ./...
```

## Parity Status (Important)

This is **not** yet a 100% full parity port of `ts-sdks`.

Detailed status is in:

- `docs/IMPLEMENTATION_STATUS.md`

In short:

- `bcs`: high functional coverage, but not full generic/type-level TS parity.
- `sui`: partial (jsonrpc/grpc/graphql/faucet + cryptography/keypairs + transactions builder/resolver/executor + zklogin/multisig/verify are implemented; full protobuf grpc parity and some advanced TS semantics are still pending).
- `walrus`: partial (storage-node + baseline client; no full on-chain flow/wasm encoding pipeline).
- `seal`: partial-to-medium (core crypto flow present; not full TS protocol compatibility/committee behavior).
- `deepbook_v3`: substantial API surface implemented, but not yet 100% TS parity (especially generated `contracts/*` BCS type model, advanced simulation parsers, and full parity vectors).

## Design Notes

- Go-first API naming and signatures.
- Practical runtime behavior first; strict TS type-system parity is intentionally not attempted.
- Additional parity should be added incrementally with tests mirroring TS package behavior.
