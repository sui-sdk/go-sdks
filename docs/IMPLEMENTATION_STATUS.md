# Implementation Status vs `ts-sdks`

Reference used: `/Users/mac/work/sui-sdks/ts-sdks`

## Summary

Question: “Is this 100% implementation of TS SDK features?”

Answer: **No, not yet.**

Current repository provides a substantial initial baseline, but does not fully cover all TS packages/submodules/protocol features.

## Status Matrix

### 1) `@mysten/bcs`

Implemented:

- reader/writer
- uleb
- primitive/composite type builders
- serialized wrapper and transforms
- base58/base64/hex helpers

Gaps:

- TS compile-time generic inference equivalents are naturally absent in Go
- exhaustive edge-case behavior parity not yet proven against TS vectors

Status: **~80-90% functional parity** for common runtime use.

### 2) `@mysten/sui`

Implemented:

- `jsonRpc` client + transport + common method set
- `grpc` package surface + core client mapping
- `graphql` client
- `faucet` helper
- `cryptography` base module (intent/signature/public key/keypair helpers)
- keypairs:
  - `ed25519`
  - `secp256k1`
  - `secp256r1`
- `transactions` builder/input/command helpers
- `transactions` resolve plugin + core resolver
- `transactions` executors (caching/serial/parallel)
- `multisig` baseline model/signer
- `zklogin` helper module (jwt/nonce/address/signature)
- `verify` helper module

Major missing modules (TS has many):

- full protobuf-native `grpc` service/type parity
- strict transaction wire-level parity (current serializer/executor is baseline-compatible, not full TS internal parity)
- passkey keypair
- rich typed RPC schemas and converters

Status: **partial (~70-80%)** of total TS `sui` surface.

Notes:

- `secp256k1` currently runs in a stdlib-compatible ECDSA mode to keep zero external dependencies.
- For strict secp256k1 compatibility with TS/noble vectors, a dedicated secp256k1 implementation is still required.
- `grpc` module currently maps through unified core method calls; generated protobuf types/services parity is still pending.

### 3) `@mysten/walrus`

Implemented:

- constants and error types
- storage-node API client
- basic high-level client read/write wrappers

Missing / incomplete:

- full blob encode/decode wasm pipeline
- committee/shard assignment logic parity
- register/upload/certify full transaction flow
- quilt/file abstractions and upload relay complete flow
- contract object decoding parity

Status: **partial (~25-40%)**.

### 4) `@mysten/seal`

Implemented:

- core errors/types
- session keys
- shamir
- DEM ciphers
- encrypt/decrypt core data flow
- key-server fetch helper primitives

Missing / incompatible areas:

- full BLS12-381 + IBE protocol parity with TS impl
- exact wire/protocol/BSC compatibility with TS for every edge path
- committee key-server and verification semantics parity
- all TS specialized error mapping and advanced consistency checks

Status: **partial-to-medium (~40-60%)** depending on scenario.

### 5) `@mysten/deepbook-v3`

Implemented:

- Go package `deepbook_v3` with client/config/types/utils
- major transaction contract classes under `deepbook_v3/transactions`
- dry-run based query helpers in client
- testnet/mainnet constants and package IDs

Missing / incomplete:

- full `src/contracts/*` generated type parity
- exact BCS structure decoders for all complex return values
- full method-by-method conformance vectors against TS implementation

Status: **partial-to-high (~65-80%)** API coverage, not full behavioral parity yet.

## What “Full Parity” Would Require Next

1. Sui package completion (`grpc`, `transactions`, `keypairs`, `zklogin`, `multisig`, `verify`).
2. Walrus full flow (encode/register/upload/certify + quilt/file abstractions).
3. Seal cryptographic parity against TS reference vectors and protocol behavior.
4. Cross-language conformance tests:
   - serialize in TS, parse in Go; and vice versa
   - encrypt/decrypt vectors across TS/Go
   - same network responses through both SDKs produce same outputs
5. Add CI with regression tests for every module.

## Testing Status in This Repo

Added unit tests for:

- `bcs`
- `sui/jsonrpc`
- `sui/graphql`
- `sui/faucet`
- `walrus`
- `seal`

Run with:

```bash
go test ./...
```
