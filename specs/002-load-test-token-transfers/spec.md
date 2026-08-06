# Feature Specification: Load Test Token Transfers

**Feature Branch**: `sui-load-tests`

**Created**: 2026-08-05

**Status**: Draft

**Input**: User description: "previously we've implemented standalone load tests for SUI in `/chainlink-sui/integration-tests/load`. Now we need to implement token transfers on these tests. You can use `sui-starter-kit` as a base and `core /chainlink/integration-tests/smoke/ccip` as example again."

## Clarifications

### Session 2026-08-05

- Q: Which token pool types should the load tests support in v1? → A: Start with CCIP-BnM token managed by `managed_token_pool` as the primary token, since it is the canonical test token available on testnet via the faucet. The pool module is `managed_token_pool` (pool kind: `managed`), not `burn_mint_token_pool`. LockRelease and other pool types can be added later if needed.
- Q: Which fee token should token-transfer load tests use? → A: Reuse the existing v1 decision: native SUI for Sui→EVM (via PTB gas budget) and native ETH for EVM→Sui. LINK fee token support is deferred, consistent with the message-only load tests.
- Q: Should token transfers also carry a message payload, or be token-only? → A: Support both modes via the run config: a "token-only" mode (empty message data, zero gas limit, empty receiver object IDs for Sui destination) and a "token + message" mode (non-empty data, non-zero gas limit, receiver package ID for Sui destination). This mirrors the sui-starter-kit's `ccipSendTokenRouter.ts` vs `ccipSendMsgAndTokenRouter.ts` split.
- Q: For EVM→Sui token transfers, does the destination need a registered CCIP receiver? → A: For token-only transfers to an EOA, no receiver is needed (empty receiver package ID, token receiver set in extraArgs). For token + message transfers to a Sui object, a registered dummy receiver is required, and its package ID must be supplied via the run config. The receiver state object ID and clock object ID are resolved at runtime from the receiver package via `getObjectFromPackage`, consistent with the sui-starter-kit approach.
- Q: Should agents run Go commands during implementation? → A: No. All Go commands will be run manually by the user. Agents must only create/edit source files. (Carried over from the prior feature.)
- Q: For Sui→EVM token transfers, how should the token coin pool pre-split coins — into exact per-message transfer amounts, or into equal-sized chunks? → A: Exact per-message amounts, configured via the run config TOML. Each pre-split coin equals exactly one transfer amount, so the `lock_or_burn` consumes the entire coin with no leftovers. This matches the sui-starter-kit's `splitCoin` pattern.
- Q: Should the token amount in the run config TOML be specified in base units (raw smallest units) or human-readable units? → A: Base units (raw smallest units, e.g., `1000000000` for 1 token with 9 decimals). This avoids floating-point parsing, matches how the existing SUI gas budget is configured (in MIST), and is consistent with the sui-starter-kit's `convertToBaseUnit` pattern. Operators convert manually using the token's decimals.

## User Scenarios & Testing

### User Story 1 - Run a Sui→EVM token-only load test (Priority: P1)

An operator wants to send a batch of CCIP token transfers from a Sui chain (testnet/mainnet) to an EVM destination chain (e.g., Sepolia, Base Sepolia). The operator configures the run via the existing four-layer config system, adding token-specific parameters to the run config TOML: the token's coin metadata object ID, the amount per message, and the EVM receiver address. The operator runs a Go test file that sends N token transfers sequentially, logs each message ID and transaction hash, and saves the results to a file.

**Why this priority**: Sui→EVM is the primary direction for CCIP load testing, and token-only transfers are the simplest token path (no destination receiver contract needed). This extends the existing message-only Sui→EVM load test with the token pool `lock_or_burn` call before `ccip_send`.

**Independent Test**: Can be fully tested by running the Go test against Sui testnet with a funded wallet holding CCIP-BnM tokens, valid `.env.testnet`, `addresses-testnet.json`, `networks-testnet.yaml`, and a run config TOML with token parameters. The test sends at least one token transfer and verifies that a message ID and transaction hash are logged and written to the results file.

**Acceptance Scenarios**:

1. **Given** a configured environment with a funded Sui wallet holding CCIP-BnM tokens, **When** the operator runs the Sui→EVM token load test with 1 message, **Then** the test successfully sends the token transfer, logs the message ID and transaction hash, and writes them to the results file.
2. **Given** a configured environment, **When** the operator runs the test with N=10 token transfers, **Then** all 10 transfers are sent successfully, each with a unique message ID and transaction hash recorded.
3. **Given** a Sui wallet with insufficient CCIP-BnM token balance for the requested amount, **When** the operator runs the test, **Then** the test fails with a clear error indicating insufficient token balance.
4. **Given** a run config that specifies a token coin metadata ID not registered in the CCIP token admin registry, **When** the operator runs the test, **Then** the test fails with a clear error indicating the token is not registered.

---

### User Story 2 - Run an EVM→Sui token-only load test (Priority: P1)

An operator wants to send a batch of CCIP token transfers from an EVM chain (e.g., Sepolia) to a Sui destination chain (testnet/mainnet). The operator configures the EVM side via `.env` (EVM private key), `addresses.json` (contract addresses including the ERC-20 token), network config YAML, and a run config TOML with token parameters: the ERC-20 token address, the amount per message, and the Sui token receiver address (an EOA). The operator runs a Go test file that sends N token transfers sequentially, handling ERC-20 allowance approvals, logging each message ID and transaction hash, and saving the results to a file.

**Why this priority**: The reverse direction (EVM→Sui) is equally important for comprehensive load testing. Token-only transfers to a Sui EOA are the simplest EVM→Sui token path and extend the existing message-only EVM→Sui load test with ERC-20 token amounts and the Sui-specific extra args (token receiver).

**Independent Test**: Can be fully tested by running the Go test against Sepolia testnet with a funded EVM wallet holding CCIP-BnM ERC-20 tokens, valid config files, and a run config TOML with token parameters. The test sends at least one token transfer to a Sui EOA and verifies the message ID and transaction hash are recorded.

**Acceptance Scenarios**:

1. **Given** a configured environment with a funded EVM wallet holding CCIP-BnM ERC-20 tokens, **When** the operator runs the EVM→Sui token load test with 1 message, **Then** the test successfully sends the token transfer, logs the message ID and transaction hash, and writes them to the results file.
2. **Given** a configured environment, **When** the operator runs the test with N=10 token transfers, **Then** all 10 transfers are sent successfully, each with a unique message ID and transaction hash recorded.
3. **Given** an EVM wallet with insufficient ERC-20 token balance, **When** the operator runs the test, **Then** the test fails with a clear error indicating insufficient token balance.
4. **Given** an EVM wallet with insufficient native ETH for fees, **When** the operator runs the test, **Then** the test fails with a clear error about insufficient funds for gas.

---

### User Story 3 - Run a Sui→EVM token + message load test (Priority: P2)

An operator wants to send a batch of CCIP messages that carry both a token transfer and an arbitrary data payload from a Sui chain to an EVM destination chain. This is the "programmable token transfer" path. The operator configures the run config TOML with both token parameters (coin metadata ID, amount) and message parameters (data payload, non-zero EVM callback gas limit). The operator runs a Go test file that sends N combined transfers sequentially and records the results.

**Why this priority**: Combined token + message transfers exercise the full programmable CCIP path and are a common real-world use case. This is secondary to token-only because it requires a destination receiver contract that can handle the message callback, adding configuration complexity.

**Independent Test**: Can be fully tested by running the Go test against Sui testnet with a funded wallet, valid config, and a run config TOML with both token and message parameters. The test sends at least one combined transfer to an EVM receiver contract and verifies the message ID and transaction hash are recorded.

**Acceptance Scenarios**:

1. **Given** a configured environment with a funded Sui wallet and an EVM receiver contract address, **When** the operator runs the Sui→EVM token + message load test with 1 message, **Then** the test successfully sends the combined transfer with the token and data payload, logs the message ID and transaction hash, and writes them to the results file.
2. **Given** a run config with a non-zero EVM callback gas limit and a valid receiver contract, **When** the operator runs the test with N=5 combined transfers, **Then** all 5 transfers are sent successfully with unique message IDs recorded.

---

### User Story 4 - Run an EVM→Sui token + message load test (Priority: P2)

An operator wants to send a batch of CCIP messages that carry both a token transfer and a data payload from an EVM chain to a Sui destination chain, targeting a registered Sui CCIP receiver object. The operator configures the run config TOML with token parameters (ERC-20 token address, amount), message parameters (data payload, gas limit), and the Sui receiver package ID. The receiver state object ID and clock object ID are resolved at runtime from the receiver package. The operator runs a Go test file that sends N combined transfers sequentially and records the results.

**Why this priority**: The reverse programmable path completes bidirectional token + message coverage. It requires a registered Sui dummy receiver and correct SuiExtraArgsV1 encoding with receiver object IDs, making it the most configuration-heavy scenario.

**Independent Test**: Can be fully tested by running the Go test against Sepolia testnet with a funded EVM wallet, a registered Sui dummy receiver, valid config, and a run config TOML with all receiver parameters. The test sends at least one combined transfer and verifies the message ID and transaction hash are recorded.

**Acceptance Scenarios**:

1. **Given** a configured environment with a funded EVM wallet and a registered Sui dummy receiver, **When** the operator runs the EVM→Sui token + message load test with 1 message, **Then** the test successfully sends the combined transfer, logs the message ID and transaction hash, and writes them to the results file.
2. **Given** a run config with a Sui receiver package ID and a non-zero gas limit, **When** the operator runs the test with N=5 combined transfers, **Then** all 5 transfers are sent successfully with unique message IDs recorded.
3. **Given** a run config that omits the Sui receiver package ID for a token + message transfer, **When** the operator runs the test, **Then** the test fails with a clear error indicating the receiver package ID is required for programmable transfers.

---

### User Story 5 - Extend the four-layer config with token parameters (Priority: P2)

An operator needs to configure a token transfer load test run. The existing four-layer config system (`.env`, `addresses.json`, `networks.yaml`, run config TOML) is extended with token-specific fields in the run config TOML: a `[token]` section with the token type/identifier, amount per message, and (for EVM source) the ERC-20 token address; and a `[sui_receiver]` section (for EVM→Sui programmable transfers) with the receiver package ID. The receiver state object ID and clock object ID are resolved at runtime. The existing message-only run configs continue to work unchanged.

**Why this priority**: Configuration is a prerequisite for all token transfer test runs. This story ensures the config extension is backward-compatible and covers all token scenarios.

**Independent Test**: Can be tested by creating a run config TOML with token parameters and running a test that loads all four layers and prints the resolved configuration (including token fields) without sending any messages.

**Acceptance Scenarios**:

1. **Given** a run config TOML with a `[token]` section specifying a coin metadata ID and amount, **When** the config loader reads the file, **Then** the token parameters are available in the test configuration.
2. **Given** a run config TOML with a `[sui_receiver]` section specifying the receiver package ID, **When** the config loader reads the file, **Then** the Sui receiver parameters are available in the test configuration.
3. **Given** an existing message-only run config TOML without any `[token]` or `[sui_receiver]` sections, **When** the config loader reads the file, **Then** the test configuration loads successfully with token fields set to their zero/empty defaults, preserving backward compatibility.
4. **Given** a run config TOML with a `[token]` section but missing the amount field, **When** the config loader validates the file, **Then** it fails with a clear error indicating the token amount is required when a token is specified.

### Edge Cases

- What happens when the Sui wallet has enough tokens for some but not all messages in a multi-message run? (Partial run: send until balance is exhausted, record failures for the rest.)
- How does the system handle a token coin object that is too small (dust) to be split for a transfer? (Skip or merge coins before splitting.)
- What happens when the ERC-20 token allowance is reset or revoked mid-run on EVM? (The send fails; record the error and continue to the next message.)
- How does the system handle a Sui token coin that is consumed (locked/burned) by a previous message in the same PTB? (Each message must use a fresh token coin object, pre-split from the sender's token balance — similar to the existing SUI gas/fee coin pool.)
- What happens when the destination chain's token pool is cursed or rate-limited during a run? (The send fails at the source; record the error and continue.)
- How does the system handle a run config that specifies token transfer mode but the `addresses.json` does not contain the required token pool addresses? (Fail with a clear error listing the missing address types.)
- What happens when the EVM receiver address for a Sui→EVM token transfer is a contract that does not support `onTokenTransfer`? (The message is sent successfully from the source; destination execution failure is out of scope for v1 — no confirmation waiting.)

## Requirements

### Functional Requirements

- **FR-001**: The system MUST provide a Go test file (or extend the existing Sui→EVM test) that sends CCIP token transfers from Sui to EVM by calling the token pool's `lock_or_burn` followed by the OnRamp's `ccip_send` within a single Programmable Transaction Block.
- **FR-002**: The system MUST provide a Go test file (or extend the existing EVM→Sui test) that sends CCIP token transfers from EVM to Sui by calling the EVM Router's `ccipSend` with token amounts and the correct SuiExtraArgsV1 encoding (including the token receiver address).
- **FR-003**: The system MUST support token-only transfers (empty message data, zero destination gas limit, empty receiver object IDs for Sui destination) as a distinct mode selectable via the run config.
- **FR-004**: The system MUST support token + message transfers (non-empty data, non-zero destination gas limit, receiver package ID for Sui destination) as a distinct mode selectable via the run config.
- **FR-005**: The system MUST resolve the token's pool configuration (pool package ID, pool module, pool state object ID, pool kind) from the CCIP token admin registry at runtime, given the token's coin metadata object ID. These addresses are not static — they are stored on-chain in the token admin registry and fetched via `get_token_config_struct`, consistent with the sui-starter-kit approach. The address book only contains the CCIP package ID, not per-token pool configs.
- **FR-006**: The system MUST resolve the latest deployed package versions for the token pool before sending, consistent with the existing latest-package resolution for CCIP and OnRamp.
- **FR-007**: For Sui→EVM token transfers, the system MUST pre-split the sender's token coins into per-message coin objects before the run, where each coin equals exactly the transfer amount configured in the run config TOML. Each message consumes one token coin entirely via `lock_or_burn`, leaving no leftovers. This is analogous to the existing SUI gas/fee coin pool but with exact per-message amounts instead of equal-sized chunks.
- **FR-008**: For EVM→Sui token transfers, the system MUST handle ERC-20 token allowance by approving the Router to spend the total token amount (plus fee amount if paying fees in the same token) **once at the start of the run, before any messages are sent**. The system MUST skip the approval transaction if the existing allowance is already sufficient for the entire run. Per-message approval is not needed.
- **FR-009**: The system MUST query the CCIP fee before sending each token transfer, including the token amount in the fee calculation, and apply a buffer (20%) to the estimated fee.
- **FR-010**: The system MUST log each message ID, transaction hash, and token amount to stdout during the test run.
- **FR-011**: The system MUST save run results (message IDs, transaction hashes, token amounts) to a file in a machine-readable format after the test completes, extending the existing results format with token fields.
- **FR-012**: The system MUST send messages sequentially (one at a time, waiting for each to complete before sending the next), consistent with the existing message-only load tests.
- **FR-013**: The system MUST NOT depend on the Chainlink core repository for any imports, consistent with the existing load test constraint.
- **FR-014**: The run config TOML MUST be extended with an optional `[token]` section containing: the token identifier (using `coin_metadata_id` key for Sui source, or `token_address` key for EVM source — both map to a single `TokenIdentifier` field), the amount per message in base units (raw smallest units, e.g., `1000000000` for 1 token with 9 decimals), and the token transfer mode ("token_only" or "token_and_message"). When the `[token]` section is absent, the test runs in message-only mode (backward compatible).
- **FR-015**: The run config TOML MUST be extended with an optional `[sui_receiver]` section for EVM→Sui programmable transfers, containing: the receiver package ID. The receiver state object ID and clock object ID are resolved at runtime from the receiver package via `getObjectFromPackage`, consistent with the sui-starter-kit approach. When absent, token-only transfers to a Sui EOA use the token receiver from extra args with an empty receiver package ID.
- **FR-016**: The system MUST validate the run config at load time: if a `[token]` section is present, the token amount MUST be greater than zero; if the transfer mode is "token_and_message", the message data and destination gas limit MUST be non-empty/non-zero.
- **FR-017**: The system MUST support the `managed_token_pool` module (pool kind: `managed`) as the primary pool type for v1, which manages CCIP-BnM tokens. LockRelease and other pool types MAY be supported if the pool kind is auto-detected from the token admin registry, but are not required for v1.
- **FR-018**: The system MUST handle Sui token coins using the same two-coin-per-message discipline as the existing SUI gas/fee pool: each message uses a distinct gas coin, a distinct fee coin, and a distinct token coin. No coin object may be reused across messages in the same run.
- **FR-019**: The system MUST accept token amounts in base units (raw smallest units) in the run config TOML. No decimal conversion is performed at runtime — the operator is responsible for converting human-readable amounts to base units using the token's known decimals. This is consistent with how the existing SUI gas budget is configured in MIST.
- **FR-020**: The system MUST continue to pay fees in native SUI (Sui→EVM) and native ETH (EVM→Sui) for v1, consistent with the existing message-only load tests. LINK and other fee token support is deferred.

### Key Entities

- **TokenTransferConfig**: Represents the token-specific parameters for a load test run, loaded from the `[token]` section of the run config TOML. Contains the token identifier, amount per message, and transfer mode.
- **SuiReceiverConfig**: Represents the Sui receiver parameters for EVM→Sui programmable transfers, loaded from the `[sui_receiver]` section of the run config TOML. Contains the receiver package ID. The receiver state object ID and clock object ID are resolved at runtime from the receiver package.
- **TokenPoolConfig**: The runtime-resolved token pool configuration fetched from the CCIP token admin registry. Contains the pool package ID, pool module name, pool state object ID, pool kind (managed / lock_release / burn_mint), and the token's coin type.
- **TokenCoinPool**: A pre-split pool of token coin objects (analogous to the existing `SuiCoinPool` for SUI gas/fee coins) that provides one fresh token coin per message for Sui→EVM transfers.
- **SentMessage** (extended): The existing results entity, extended with token amount and token identifier fields to record what was transferred in each message.

## Success Criteria

### Measurable Outcomes

- **SC-001**: An operator can configure and run a Sui→EVM token-only load test with 50 token transfers in under 10 minutes, with all message IDs, transaction hashes, and token amounts recorded to a results file.
- **SC-002**: An operator can configure and run an EVM→Sui token-only load test with 50 token transfers in under 10 minutes, with all message IDs, transaction hashes, and token amounts recorded to a results file.
- **SC-003**: An operator can configure and run a Sui→EVM token + message load test with 20 combined transfers in under 10 minutes, with all message IDs and transaction hashes recorded to a results file.
- **SC-004**: An operator can configure and run an EVM→Sui token + message load test with 20 combined transfers in under 10 minutes, with all message IDs and transaction hashes recorded to a results file.
- **SC-005**: The extended run config TOML is backward-compatible: an existing message-only run config loads and runs without modification after the token transfer feature is added.
- **SC-006**: The token transfer load test framework compiles and runs without importing any package from the Chainlink core repository.
- **SC-007**: The token pool configuration is resolved at runtime from the CCIP token admin registry, so an operator does not need to manually specify pool package IDs or state object IDs in the run config.

## Assumptions

- The operator has access to a deployed CCIP environment (testnet/staging/mainnet) with all required contract addresses, including token pool deployments.
- The operator has funded wallets: Sui wallet with sufficient SUI (for gas/fees) and CCIP-BnM tokens (for transfers); EVM wallet with sufficient ETH (for gas/fees) and CCIP-BnM ERC-20 tokens (for transfers).
- The CCIP-BnM token is registered in the CCIP token admin registry on the source chain, with a deployed `managed_token_pool` (pool kind: `managed`) for v1.
- The Sui and EVM RPC endpoints are accessible from the machine running the tests.
- The `addresses-<env>.json` file includes the token pool addresses and coin metadata object IDs needed for token transfers, or these are resolvable at runtime from the CCIP objects.
- Token-only transfers to an EVM EOA do not require a destination receiver contract; the EVM address is left-padded to 32 bytes and used as the message receiver.
- Token-only transfers to a Sui EOA do not require a registered receiver; the receiver package ID is empty (zero bytes32) and the token receiver address is set in the SuiExtraArgsV1. This is consistent with the sui-starter-kit's `ccipSendTokenRouter.ts`, which uses `ethers.ZeroHash` for the receiver and sets `tokenReceiver` to the EOA address in extra args.
- Token + message transfers to a Sui object require a registered dummy receiver; its package ID is supplied via the run config. The receiver state object ID and clock object ID are resolved at runtime from the receiver package, consistent with the sui-starter-kit's `ccipSendMsgAndTokenRouter.ts` approach.
- No confirmation or metrics features (waiting for message execution on destination, checking token balances on destination) are included in v1, consistent with the message-only load tests.
- No retry logic in v1; a failed send is logged and the run continues to the next message, consistent with the existing message-only load tests.
- The existing `chainlink-sui/integration-tests/load` Go module and its dependencies are reused; no new external dependencies are needed beyond what is already in `go.mod`.
- All Go commands (go build, go test, go mod tidy, etc.) MUST be run manually by the operator. Agents MUST NOT execute any Go commands — they only create and edit source files.
- The Sui PTB constraint that a mutable coin object cannot appear twice in the same transaction (documented in repository memory) is respected: gas coin, fee coin, and token coin are all distinct objects per message.
