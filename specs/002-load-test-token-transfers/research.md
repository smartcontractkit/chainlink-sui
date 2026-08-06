# Research: Load Test Token Transfers

**Feature**: 002-load-test-token-transfers
**Date**: 2026-08-05

## 1. Sui→EVM Token Transfer PTB Construction

### Decision
Extend the existing `SendMessage` PTB with a `lock_or_burn` call between `create_token_transfer_params` and `ccip_send`. The PTB sequence becomes:

```
create_token_transfer_params → lock_or_burn (token pool) → ccip_send (onramp)
```

### Rationale
This matches the sui-starter-kit's `buildTokenTransferPTB` pattern in `buildCcipSendPTB.ts`. The `create_token_transfer_params` result is passed to both `lock_or_burn` (as `tokenTransferParams`) and `ccip_send` (as the token params argument). The token pool's `lock_or_burn` consumes the token coin object, and `ccip_send` references the same `tokenTransferParams` to include the token in the CCIP message.

### Alternatives Considered
- **Separate transactions for lock_or_burn and ccip_send**: Rejected — would require two Sui transactions per message, doubling latency and gas costs. The PTB model is designed for atomic multi-call transactions.
- **Skip create_token_transfer_params for token-only**: Rejected — the onramp's `ccip_send` always expects a `TokenTransferParams` argument, even for message-only sends. The helper function is always called.

### Implementation Notes
- The `lock_or_burn` call uses the generated `ManagedTokenPoolContract.LockOrBurn` binding from `bindings/generated/ccip/ccip_token_pools/managed_token_pool/`.
- Parameters: `typeArgs` = coin type (e.g., `de9a44c43b1e5cf3bee4ae5d6c1aa53f2981513ab3354ebace4fba470f44f92a::ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN`), `ref` = CCIPObjectRef, `tokenTransferParams` = result from step 1, `c_` = token coin object, `remoteChainSelector` = dest chain, `clock` = `0x6`, `denyList` = deny list object ID, `tokenState` = token state object ID, `state` = pool state object ID.
- Note: `managed_token_pool.lock_or_burn` takes 8 arguments (includes denyList and tokenState), unlike `burn_mint_token_pool.lock_or_burn` which takes 6 arguments.
- The token pool package ID and state object ID are resolved at runtime from the token admin registry (see Research #3).

---

## 2. EVM→Sui Token Transfer via Router

### Decision
Extend the existing `SendMessage` with token amounts in the `ClientEVM2AnyMessage` struct and correct `SuiExtraArgsV1` encoding. The `ccipSend` call on the Router already supports `TokenAmounts` — the only change is populating that field and encoding the `tokenReceiver` in extra args.

### Rationale
The EVM Router's `ccipSend` accepts `ClientEVM2AnyMessage` which includes `TokenAmounts []ClientEVMTokenAmount`. The existing message-only code passes an empty slice. For token transfers, we populate it with the ERC-20 token address and amount. The `SuiExtraArgsV1` encoding already supports `TokenReceiver` — we set it to the Sui EOA address for token-only, or the receiver state object ID for token+message.

### Alternatives Considered
- **Separate Router function for token transfers**: Rejected — the Router's `ccipSend` is the unified entry point for all CCIP messages (message-only, token-only, and combined). No separate function exists.
- **Use LINK as fee token for EVM→Sui**: Deferred to future iteration (consistent with v1 scope).

### Implementation Notes
- Token-only mode: `receiver = ZeroHash` (32 zero bytes), `data = []byte{}`, `receiverObjectIds = []`, `tokenReceiver = Sui EOA address`.
- Token+message mode: `receiver = receiver package ID`, `data = message payload`, `receiverObjectIds = [clockObj, receiverState]`, `tokenReceiver = receiverState`.
- The `receiverState` object ID is resolved at runtime from the receiver package via `getObjectFromPackage` (devInspect call).
- ERC-20 allowance is approved once at the start of the run for the total amount (token amount × message count + fee buffer if paying fees in same token).

---

## 3. Runtime Token Pool Config Resolution

### Decision
Resolve the token pool configuration at runtime by calling the CCIP token admin registry's `get_token_config_struct` function via devInspect, then parse the BCS-encoded result to extract the pool package ID, module name, state object ID, and pool kind.

### Rationale
The token pool config is stored on-chain in the token admin registry, not in the address book. The sui-starter-kit resolves it at runtime via `getTokenConfig()` in `scripts/sui-helper/getTokenConfig.ts`. The address book only contains the CCIP package ID — per-token pool configs are dynamic and fetched on-chain. This avoids requiring operators to manually specify pool package IDs and state object IDs in the run config.

### Alternatives Considered
- **Static config in run TOML**: Rejected — would require operators to look up pool addresses manually, which is error-prone and fragile across deployments.
- **Address book extension**: Rejected — the address book format is standardized by the deployment framework and adding per-token pool entries would require changes to the deployment pipeline.

### Implementation Notes
- Use `bind.NewBoundContract` with the latest CCIP package to call `token_admin_registry::get_token_config_struct`.
- The function takes `CCIPObjectRef` and `coinMetadataId` as arguments.
- The return value is a BCS-encoded struct containing: `token_pool_package_id` (address), `token_pool_module` (string), `token_type` (string), `administrator` (address), `pending_administrator` (address), `token_pool_type_proof` (string), `lock_or_burn_params` (vector<address>), `release_or_mint_params` (vector<address>).
- Pool kind is determined from the module name: `managed_token_pool` → managed, `lock_release_token_pool` → lock_release, `burn_mint_token_pool` → burn_mint.
- Pool state object ID is extracted from `lock_or_burn_params[3]` for managed, or `lock_or_burn_params[1]` for lock_release/burn_mint.
- The pool package ID must also be resolved to its latest version, consistent with the existing CCIP/OnRamp latest-package resolution.

---

## 4. Token Coin Pre-Splitting (Sui→EVM)

### Decision
Pre-split the sender's token coins into exact per-message amounts before the run, using the same pattern as the existing `SuiCoinPool` but with exact amounts instead of equal-sized chunks. Each message consumes one token coin entirely via `lock_or_burn`.

### Rationale
The Sui PTB constraint (documented in repository memory) prohibits using the same mutable coin object twice in one transaction. Each message already uses distinct gas and fee coins — adding a distinct token coin per message follows the same discipline. Exact per-message amounts mean no leftovers, simplifying the coin management.

### Alternatives Considered
- **Equal-sized chunks**: Rejected — would leave dust leftovers that need to be merged or wasted, adding complexity.
- **Split on-demand per message**: Rejected — would add a split transaction before each send, doubling the transaction count and latency.

### Implementation Notes
- Create a `TokenCoinPool` struct analogous to `SuiCoinPool` but parameterized by coin type (not hardcoded to `0x2::sui::SUI`).
- The pool queries the sender's coins of the token type, finds the largest one, and splits it into N coins of exactly `amountPerMessage` each.
- Required coins = `messageCount` (one per message, no buffer needed since each coin is consumed entirely).
- The split uses `sui::pay::split_vec` via PTB, similar to the existing SUI coin split logic.

---

## 5. ERC-20 Allowance Pre-Approval (EVM→Sui)

### Decision
Approve the Router to spend the total token amount once at the start of the run, before any messages are sent. Skip the approval if the existing allowance is already sufficient.

### Rationale
This matches the sui-starter-kit's `approveToken` pattern in `ccipSendTokenRouter.ts`. Approving once for the total run amount avoids N approval transactions (one per message), reducing gas costs and latency. The check for existing allowance avoids unnecessary approval transactions on subsequent runs.

### Alternatives Considered
- **Per-message approval**: Rejected — wasteful in gas and adds latency.
- **Infinite approval**: Rejected — security best practice is to approve only what's needed.

### Implementation Notes
- Total approval amount = `tokenAmount × messageCount` (for token transfers) + fee buffer if paying fees in the same token (not applicable for v1 since fees are in native ETH).
- Use the standard ERC-20 `approve` function via `go-ethereum`'s `bind.TransactOpts`.
- Wait for the approval transaction to be mined before starting the send loop.

---

## 6. Config Extension: Backward Compatibility

### Decision
Add optional `[token]` and `[sui_receiver]` sections to the run config TOML. When absent, the test runs in message-only mode (existing behavior). When present, the test runs in token transfer mode.

### Rationale
This preserves backward compatibility with existing run configs. Operators who only need message-only tests are unaffected. The mode is determined by the presence of the `[token]` section, not a separate flag.

### Alternatives Considered
- **Separate test files for token transfers**: Rejected — would fragment the operator interface. A single test file with mode determined by config is simpler.
- **Separate `--mode` flag**: Rejected — the run config TOML is the single source of truth for run parameters. Adding a CLI flag would create ambiguity.

### Implementation Notes
- `[token]` section fields: `coin_metadata_id` (Sui source) or `token_address` (EVM source), `amount` (base units), `mode` ("token_only" or "token_and_message").
- `[sui_receiver]` section fields: `package_id` (receiver package ID for EVM→Sui programmable transfers).
- Validation: if `[token]` is present, `amount` must be > 0. If `mode` is "token_and_message", `message_data` and `evm_callback_gas_limit` must be non-empty/non-zero.

---

## 7. Results Extension

### Decision
Extend the `SentMessage` struct with `TokenAmount` and `TokenIdentifier` fields. The results JSON file includes these fields for token transfer runs.

### Rationale
Operators need to know what was transferred in each message for post-run analysis. The existing results format is extended with optional fields that are omitted (zero-value) for message-only runs.

### Implementation Notes
- New fields: `TokenAmount string` (base units as string to avoid JSON number precision issues), `TokenIdentifier string` (coin metadata ID or ERC-20 address).
- Backward compatible: existing message-only results files are still valid JSON (new fields are `omitempty`).
