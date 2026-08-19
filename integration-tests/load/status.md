# CCIP Load Tests — Status

Status of the Sui CCIP load test framework (`integration-tests/load/`).
Lightweight Go framework (built on `chainlink-testing-framework/wasp`) for sending
CCIP messages and token transfers between Sui and EVM chains. No Chainlink core dependency.

---

## Functional

### 1. SUI → EVM arbitrary messaging — ✅ done (WASP)
Message-only sends from Sui to EVM run in parallel using WASP with N wallets
(one generator per wallet). PTB path: `create_token_transfer_params` → `ccip_send`.
Fees paid in native SUI via the PTB gas budget. See `sui2evm_gun.go`, `sui/sender.go`.

### 2. EVM → SUI arbitrary messaging — ✅ done (WASP)
Message-only sends from EVM to Sui run in parallel using WASP with N wallets.
Path: `Router.getFee` → `Router.ccipSend` (native ETH). Fees paid in native ETH with a
20% buffer. See `evm2sui_gun.go`, `evm/sender.go`.

### 3. SUI ↔ EVM token transfers — ⚠️ sequential sends only
Token transfers (`token_only` and `token_and_message` modes) are implemented but run
**sequentially** (single signer, one message at a time) via `runSui2EVMSequential` /
`runEVM2SuiSequential`. The WASP parallel path is only wired for message-only sends.
- Sui→EVM: `create_token_transfer_params` → `managed_token_pool.lock_or_burn` → `ccip_send`.
- EVM→Sui: ERC-20 `approve` → `Router.ccipSend` with `TokenAmounts` + `SuiExtraArgsV1`.
- Token pool config resolved at runtime from the CCIP token admin registry.

### 4. SUI ↔ EVM PTT (programmable token transfer) — ⚠️ sequential, untested
`token_and_message` mode (programmable token transfer) is implemented but **untested**.
It requires a Sui receiver package implementing `ccip_receive` to be deployed
(`[sui_receiver]` section in the run TOML). The receiver state object ID is resolved at
runtime via `ResolveReceiverState`. No receiver has been deployed/validated yet.

### 5. Verification + metrics — ⚠️ done SUI->EVM using scripts
`scripts/check_commit_execute_on_dest_test.go` reads a results file and queries the
destination OffRamp for `ExecutionStateChanged` / `CommitReportAccepted` events, then
computes Sent→Commit→Execute latencies (avg, median, p95, p99). Helper scripts:
`merge_sui_coins`, `query_sui_coins`.

---

## Non-functional

### 1. Deterministic wallet generation + parallel sends
Wallets are generated deterministically from a `WALLET_SEED` via HKDF-SHA256
(Ed25519 for Sui, secp256k1 for EVM). N wallets send messages in parallel, one WASP
generator per wallet, each owning its wallet exclusively (no nonce/coin contention).

### 2. SUI coin splitting
SUI coins are pre-split into a pool of gas/fee coins before sending
(`sui/coin_pool.go`). Each message consumes two distinct coins (gas + fee) to avoid the
"same mutable coin twice in one PTB" constraint. Token coins are also pre-split into
exact per-message amounts.

### 3. SUI coin recovery after tests
Remaining funds are swept back to the main signer via `SweepSuiWallets` /
`SweepEVMWallets` (registered as `t.Cleanup`). Best-effort: wallets with insufficient
balance for gas are skipped.

---  

## TODOs

1. **Token transfers over WASP** — parallelize token transfers (currently sequential-only).
2. **Retries on send (exponential backoff)** — currently no retries; failed sends are logged and counted.
3. **LINK as fee token** — deferred until LINK is obtainable on Sui; currently native SUI/ETH only.
4. **Support for other pool types** — only `managed_token_pool` lock/burn is implemented;
   `lock_release` and `burn_mint` pool kinds are recognized but not yet handled.
5. **Re-organize code** — cleanup/refactor of the load test package structure.