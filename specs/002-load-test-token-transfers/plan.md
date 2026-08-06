# Implementation Plan: Load Test Token Transfers

**Branch**: `sui-load-tests` | **Date**: 2026-08-05 | **Spec**: [spec.md](./spec.md)

**Input**: Feature specification from `/specs/002-load-test-token-transfers/spec.md`

## Summary

Extend the existing standalone Sui CCIP load test framework (`integration-tests/load/`) with token transfer support in both directions (Sui→EVM and EVM→Sui), supporting token-only and token+message modes. The implementation adds token pool interaction (`lock_or_burn`) to the Sui→EVM PTB path, ERC-20 token amounts to the EVM→Sui Router path, token coin pre-splitting, ERC-20 allowance pre-approval, and runtime token pool config resolution from the CCIP token admin registry. The existing four-layer config system is extended with optional `[token]` and `[sui_receiver]` TOML sections, preserving backward compatibility with message-only runs.

## Technical Context

**Language/Version**: Go 1.26.2+

**Primary Dependencies**: `block-vision/sui-go-sdk` (Sui JSON-RPC), `ethereum/go-ethereum` (EVM RPC + ABI), `chainlink-ccip/chains/evm/gobindings` (EVM Router bindings), `chainlink-sui/bindings` (Sui OnRamp + token pool bindings), `chainlink-sui/relayer/client` (PTB client), `chainlink-deployments-framework` (address book), `chainlink-common/pkg/logger` (structured logging — NOT imported by load tests; `log/slog` used instead for standalone constraint)

**Storage**: Filesystem — results JSON files in `results/` directory

**Testing**: Go test framework (`go test`), run manually against remote testnet/staging/mainnet. No unit tests (these are integration tests by nature). Build tag: `//go:build integration`.

**Target Platform**: macOS/Linux (operator's workstation)

**Project Type**: CLI test tool (Go test files with `--run-name` flag)

**Performance Goals**: Sequential sends, one at a time. 50 token transfers in under 10 minutes.

**Constraints**: No Chainlink core (`/chainlink`) imports. No new external dependencies beyond what is already in `go.mod`. No confirmation waiting (DON handles execution). No retries in v1. All Go commands run manually by operator.

**Scale/Scope**: 50-100 messages per run. 4 test files (Sui→EVM token, EVM→Sui token, Sui→EVM token+msg, EVM→Sui token+msg). ~5 new source files (token pool resolver, token coin pool, EVM token sender, config extensions).

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Standard Library First | ✅ PASS | No new dependencies. All required packages already in `go.mod`. Uses `log/slog` (stdlib) instead of `chainlink-common/pkg/logger` to satisfy the no-chainlink-core constraint. |
| II. Test Discipline | ✅ PASS | Load tests are integration tests (`//go:build integration`). Run manually against remote envs. No unit tests needed for the load test framework itself (it is the test). |
| III. Structured Logging | ✅ PASS | Uses `log/slog` (Go stdlib structured logging) with `slog.Info`, `slog.Error`, `slog.Debug` — consistent with existing load tests. The constitution mandates `chainlink-common/pkg/logger` for relayer components, but load tests are standalone and exempt from this requirement. |
| IV. Code Generation & Bindings | ✅ PASS | No new bindings needed. Existing `chainlink-sui/bindings/generated/` and `chainlink-ccip/chains/evm/gobindings/` are reused. |
| V. Configuration as Code | ✅ PASS | Run config uses TOML (`pelletier/go-toml/v2`). New `[token]` and `[sui_receiver]` sections are optional and validated. Secrets remain in `.env` files. |

**Gate Result**: ALL PASS — no violations.

### Post-Design Re-Evaluation (Phase 1 Complete)

| Principle | Status | Notes |
|-----------|--------|-------|
| I. Standard Library First | ✅ PASS | No new dependencies introduced. `log/slog` (stdlib) used for logging. All required packages (`go-ethereum`, `sui-go-sdk`, `chainlink-ccip/gobindings`, `chainlink-sui/bindings`, `pelletier/go-toml`) already in `go.mod`. |
| II. Test Discipline | ✅ PASS | Load tests are integration tests by nature (require remote RPC endpoints). The feature IS the test — adding unit tests to a test framework is intentionally excluded. New helper functions (BCS parsing, pool kind derivation) are simple enough to be validated by the integration tests themselves. Build tag `//go:build integration` applied. |
| III. Structured Logging | ✅ PASS | Uses `log/slog` with `slog.Info`, `slog.Error`, `slog.Debug` — consistent with existing load tests. The constitution's `chainlink-common/pkg/logger` requirement applies to relayer components, not standalone test tools. |
| IV. Code Generation & Bindings | ✅ PASS | No new bindings generated. Existing `ManagedTokenPoolContract`, `Onramp`, and `Router` bindings reused as-is. Token pool config resolved at runtime via devInspect (no new codegen needed). |
| V. Configuration as Code | ✅ PASS | New `[token]` and `[sui_receiver]` TOML sections are optional, validated at load time, and backward-compatible. Secrets remain in `.env`. No hard-coded addresses. |

**Post-Design Gate Result**: ALL PASS — no violations introduced by design decisions.

## Project Structure

### Documentation (this feature)

```text
specs/002-load-test-token-transfers/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (N/A — no external interfaces)
└── tasks.md             # Phase 2 output (NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
integration-tests/load/
├── config/
│   ├── config.go            # EXTEND: LoadFullConfig, add token/sui_receiver parsing
│   ├── runconfig.go          # EXTEND: RunConfig struct with [token] and [sui_receiver]
│   └── types.go              # EXTEND: LoadTestConfig, SentMessage, new types
├── evm/
│   ├── client.go             # (unchanged)
│   ├── extras.go             # (unchanged — already supports SuiExtraArgsV1 with TokenReceiver/ReceiverObjectIds)
│   ├── sender.go             # EXTEND: SendTokenMessage wrapper calling token_sender.go
│   └── token_sender.go       # NEW: ERC-20 allowance pre-approval + Router.ccipSend with token amounts
├── sui/
│   ├── client.go             # (unchanged)
│   ├── coin_pool.go          # EXTEND: TokenCoinPool (analogous to SuiCoinPool)
│   ├── extras.go             # (unchanged)
│   ├── receiver.go           # EXTEND: ResolveReceiverState (getObjectFromPackage) + SuiObjectIdToBytes32 helper
│   ├── sender.go             # EXTEND: SendTokenMessage
│   ├── token_pool.go         # NEW: managed_token_pool lock_or_burn call
│   └── token_registry.go     # NEW: Runtime token pool config resolution
├── evm2sui_test.go           # EXTEND: token-only and token+message modes
├── sui2evm_test.go           # EXTEND: token-only and token+message modes
├── runs/
│   ├── first-sui-to-evm-run.toml                     # (unchanged)
│   ├── first-evm-to-sui-run.toml                     # (unchanged)
│   ├── sui-to-evm-token-only-run.toml                # NEW: example token-only run config
│   ├── evm-to-sui-token-only-run.toml                # NEW: example token-only run config
│   ├── sui-to-evm-token-and-message-run.toml         # NEW: example token+message run config
│   └── evm-to-sui-token-and-message-run.toml         # NEW: example token+message run config
└── results/                 # (runtime output, not committed)
```

**Structure Decision**: Single project — extend the existing `integration-tests/load/` package. No new top-level directories. New files are added within the existing `config/`, `evm/`, and `sui/` sub-packages. The existing test files (`evm2sui_test.go`, `sui2evm_test.go`) are extended with token transfer modes rather than creating separate test files, keeping the operator interface simple (same `--run-name` flag, mode determined by run config).

## Complexity Tracking

> No violations — this section is intentionally empty.
