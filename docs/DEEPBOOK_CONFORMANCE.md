# Deepbook TS-Go Conformance Plan

## Current Environment Status

In the current workspace, TypeScript deepbook package cannot be executed directly because:

- `pnpm` is not installed in PATH.
- `corepack` is not available.
- `ts-sdks` workspace has no installed `node_modules`.

So this iteration focuses on Go-side deterministic encoding/decoding tests first.

## What is already covered in Go tests

- Move target mapping checks for major contracts.
- BCS argument encoding checks for:
  - `u64`
  - `u128`
  - `bool`
  - `vector<u128>`
- Client dry-run parsing happy-path and error-path checks.

## Next Step Once TS Runtime Is Available

1. Install and bootstrap TS workspace:
   - `pnpm install` in `/Users/mac/work/sui-sdks/ts-sdks`
2. Add a TS fixture generator script in `ts-sdks/packages/deepbook-v3`:
   - Build representative transactions.
   - Export normalized command payloads and selected pure-arg bytes as JSON.
3. Add Go conformance tests to load these JSON fixtures and assert:
   - same target function
   - same type arguments
   - same encoded pure argument bytes for key fields

