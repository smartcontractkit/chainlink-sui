# Tasks: Load Test Token Transfers

**Input**: Design documents from `/specs/002-load-test-token-transfers/`

**Prerequisites**: plan.md, spec.md, research.md, data-model.md, quickstart.md

**Tests**: Tests are NOT requested for this feature. The load tests themselves are integration tests run manually against remote environments. No unit tests or contract tests are generated. This is an intentional decision: the feature IS the test — adding unit tests to a test framework would be redundant. All helper logic (BCS parsing, pool kind derivation, amount calculations) is validated by the integration tests themselves.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Extend the existing load test framework with shared token transfer infrastructure.

- [X] T001 [P] Extend `integration-tests/load/config/types.go` with `TokenTransferConfig`, `SuiReceiverConfig`, and extend `LoadTestConfig` and `SentMessage` with token fields
- [X] T002 [P] Extend `integration-tests/load/config/runconfig.go` with `[token]` and `[sui_receiver]` TOML sections and validation rules
- [X] T003 Extend `integration-tests/load/config/config.go` to populate token fields in `LoadFullConfig` and validate token mode requirements

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

- [X] T004 Implement `integration-tests/load/sui/token_registry.go` to resolve `TokenPoolConfig` at runtime from the CCIP token admin registry via `get_token_config_struct`
- [X] T005 Implement `integration-tests/load/sui/token_pool.go` to call `managed_token_pool.lock_or_burn` with the correct 8-argument signature (includes denyList and tokenState)
- [X] T006 Implement `integration-tests/load/sui/coin_pool.go` extension for `TokenCoinPool` — pre-split token coins into exact per-message amounts
- [X] T007 Implement `integration-tests/load/evm/token_sender.go` with ERC-20 allowance pre-approval (once at start, for total run amount) and `Router.ccipSend` with token amounts. Include fee query with token amounts in the fee calculation and 20% buffer (FR-009)
- [X] T008 Extend `integration-tests/load/sui/receiver.go` with `ResolveReceiverState` (getObjectFromPackage via devInspect) and `SuiObjectIdToBytes32` helper for converting Sui object IDs to 32-byte padded values. Note: `evm/extras.go` already supports `SuiExtraArgsV1` encoding with `TokenReceiver` and `ReceiverObjectIds` — no changes needed there

**Checkpoint**: Foundation ready — user story implementation can now begin in parallel.

---

## Phase 3: User Story 1 - Sui→EVM Token-Only Load Test (Priority: P1) 🎯 MVP

**Goal**: An operator can run a Sui→EVM token-only load test by configuring the run config TOML with `[token]` section and `mode = "token_only"`.

**Independent Test**: Run `go test -run TestSui2EVM -v --run-name sui-to-evm-token-only-run` against Sui testnet with a funded wallet holding CCIP-BnM tokens. Verify message IDs and transaction hashes are recorded in the results file.

### Implementation for User Story 1

- [X] T009 [US1] Extend `integration-tests/load/sui/sender.go` with `SendTokenMessage` function that builds the PTB: `create_token_transfer_params` → `managed_token_pool.lock_or_burn` → `onramp.ccip_send`. Include fee query with token amounts in the fee calculation and 20% buffer (FR-009)
- [X] T010 [US1] Extend `integration-tests/load/sui2evm_test.go` to detect token-only mode from run config, prepare `TokenCoinPool`, and call `SendTokenMessage` in the send loop. Populate `TokenAmount` and `TokenIdentifier` fields in each `SentMessage` and log token amount alongside message ID (FR-010, FR-011)
- [X] T011 [US1] Add example run config `integration-tests/load/runs/sui-to-evm-token-only-run.toml`

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently.

---

## Phase 4: User Story 2 - EVM→Sui Token-Only Load Test (Priority: P1)

**Goal**: An operator can run an EVM→Sui token-only load test by configuring the run config TOML with `[token]` section and `mode = "token_only"`.

**Independent Test**: Run `go test -run TestEVM2Sui -v --run-name evm-to-sui-token-only-run` against Sepolia testnet with a funded EVM wallet holding CCIP-BnM ERC-20 tokens. Verify message IDs and transaction hashes are recorded in the results file.

### Implementation for User Story 2

- [X] T012 [US2] Extend `integration-tests/load/evm/sender.go` with `SendTokenMessage` wrapper function that calls `token_sender.go` logic. The wrapper constructs `ClientEVM2AnyMessage` with `TokenAmounts` and correct `SuiExtraArgsV1` (empty receiver, tokenReceiver = Sui EOA). The actual ERC-20 approval and `Router.ccipSend` live in `token_sender.go` (T007)
- [X] T013 [US2] Extend `integration-tests/load/evm2sui_test.go` to detect token-only mode from run config, pre-approve Router for total token amount (via `token_sender.go`), and call `SendTokenMessage` in the send loop. Populate `TokenAmount` and `TokenIdentifier` fields in each `SentMessage` and log token amount alongside message ID (FR-010, FR-011)
- [X] T014 [US2] Add example run config `integration-tests/load/runs/evm-to-sui-token-only-run.toml`

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently.

---

## Phase 5: User Story 3 - Sui→EVM Token + Message Load Test (Priority: P2)

**Goal**: An operator can run a Sui→EVM combined token + message load test by configuring the run config TOML with `[token]` section and `mode = "token_and_message"`.

**Independent Test**: Run `go test -run TestSui2EVM -v --run-name sui-to-evm-token-and-message-run` against Sui testnet with a funded wallet and an EVM receiver contract. Verify combined transfers with non-zero gas limit are recorded.

### Implementation for User Story 3

- [X] T015 [US3] Extend `integration-tests/load/sui/sender.go` `SendTokenMessage` to include non-empty message data and non-zero `evm_callback_gas_limit` when mode is `token_and_message`
- [X] T016 [US3] Extend `integration-tests/load/sui2evm_test.go` to pass message data and gas limit for `token_and_message` mode. Ensure `TokenAmount` and `TokenIdentifier` are populated in `SentMessage` results (FR-010, FR-011)
- [X] T017 [US3] Add example run config `integration-tests/load/runs/sui-to-evm-token-and-message-run.toml`

**Checkpoint**: At this point, User Stories 1, 2, and 3 should all work independently.

---

## Phase 6: User Story 4 - EVM→Sui Token + Message Load Test (Priority: P2)

**Goal**: An operator can run an EVM→Sui combined token + message load test by configuring the run config TOML with `[token]` section, `mode = "token_and_message"`, and `[sui_receiver]` section.

**Independent Test**: Run `go test -run TestEVM2Sui -v --run-name evm-to-sui-token-and-message-run` against Sepolia testnet with a funded EVM wallet and a registered Sui dummy receiver. Verify combined transfers with receiver object IDs are recorded.

### Implementation for User Story 4

- [X] T018 [US4] Extend `integration-tests/load/evm/sender.go` `SendTokenMessage` to resolve receiver state object ID at runtime and construct `SuiExtraArgsV1` with `receiverObjectIds = [clock, receiverState]` and `tokenReceiver = receiverState` for `token_and_message` mode
- [X] T019 [US4] Extend `integration-tests/load/evm2sui_test.go` to require `[sui_receiver]` package ID when mode is `token_and_message`. Ensure `TokenAmount` and `TokenIdentifier` are populated in `SentMessage` results (FR-010, FR-011)
- [X] T020 [US4] Add example run config `integration-tests/load/runs/evm-to-sui-token-and-message-run.toml`

**Checkpoint**: At this point, all user stories should be independently functional.

---

## Phase 7: User Story 5 - Extend Four-Layer Config with Token Parameters (Priority: P2)

**Goal**: The run config TOML supports optional `[token]` and `[sui_receiver]` sections, and existing message-only configs continue to work unchanged.

**Independent Test**: Run `go test -run TestSui2EVM -v --run-name first-sui-to-evm-run` (existing message-only config) and verify it still works. Create a config with `[token]` and verify it loads correctly.

### Implementation for User Story 5

- [X] T021 [US5] Ensure `integration-tests/load/config/runconfig.go` validation treats `[token]` and `[sui_receiver]` as optional and fails gracefully with clear errors when required fields are missing
- [X] T022 [US5] Ensure `integration-tests/load/config/config.go` sets token fields to nil/empty defaults when sections are absent, preserving backward compatibility
- [X] T023 [US5] Update `integration-tests/load/README.md` to document the new `[token]` and `[sui_receiver]` TOML sections and token transfer modes

**Checkpoint**: At this point, the config extension is complete and backward-compatible.

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories.

- [X] T024 [P] Update `integration-tests/load/README.md` with token transfer architecture, config examples, and validation scenarios from `quickstart.md`
- [X] T025 [P] Update `integration-tests/load/.env.example` if any new env vars are needed (none expected)
- [X] T026 Verify all new source files have `//go:build integration` build tag and follow project Go conventions. Verify no `chainlink/` (core repository) imports are present in any new or modified files (FR-013, SC-006)
- [ ] T027 Run `golangci-lint` on modified files (manually, per project workflow)
- [ ] T028 Validate `quickstart.md` scenarios by creating the 4 example run config TOML files and ensuring they load without errors

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately.
- **Foundational (Phase 2)**: Depends on Setup completion — BLOCKS all user stories.
- **User Stories (Phases 3-6)**: All depend on Foundational phase completion.
  - User stories can proceed in parallel (if staffed).
  - Or sequentially in priority order: US1 (P1) → US2 (P1) → US3 (P2) → US4 (P2).
- **Config Extension (Phase 7)**: Depends on Setup phase (T001-T003) only; can run in parallel with user stories but must be complete before final polish.
- **Polish (Phase 8)**: Depends on all user stories and config extension being complete.

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational phase. No dependencies on other stories.
- **User Story 2 (P1)**: Can start after Foundational phase. No dependencies on other stories.
- **User Story 3 (P2)**: Can start after User Story 1 (reuses `SendTokenMessage` infrastructure) or in parallel if the function is designed to handle both modes.
- **User Story 4 (P2)**: Can start after User Story 2 (reuses EVM token sender infrastructure) or in parallel if the function is designed to handle both modes.
- **User Story 5 (P2)**: Can start after Setup phase; config work is independent of user story implementation.

### Within Each User Story

- Config/models before test logic.
- Sender functions before test file integration.
- Example run configs after test logic is complete.
- Story complete before moving to next priority.

### Parallel Opportunities

- T001, T002, T003 (Setup) can run in parallel.
- T004, T005, T006, T007, T008 (Foundational) can run in parallel after Setup.
- T009/T010 (US1), T012/T013 (US2) can run in parallel after Foundational.
- T015/T016 (US3), T018/T019 (US4) can run in parallel after US1/US2 respectively.
- T021, T022, T023 (US5) can run in parallel with user stories after Setup.
- T024, T025, T026, T027, T028 (Polish) can run in parallel after all stories.

---

## Parallel Example: User Story 1

```bash
# Terminal 1: Implement config types
# Work on: integration-tests/load/config/types.go

# Terminal 2: Implement token registry
# Work on: integration-tests/load/sui/token_registry.go

# Terminal 3: Implement token pool lock_or_burn
# Work on: integration-tests/load/sui/token_pool.go

# Terminal 4: Implement token coin pool + receiver resolution
# Work on: integration-tests/load/sui/coin_pool.go, integration-tests/load/sui/receiver.go

# After all complete:
# Terminal 5: Integrate into sui2evm_test.go
# Work on: integration-tests/load/sui2evm_test.go
```

---

## Suggested MVP Scope

**MVP = User Story 1 only** (Sui→EVM token-only load test). This delivers the primary direction for CCIP load testing with the simplest token path. Once MVP is verified, implement User Story 2 (EVM→Sui token-only), then the P2 programmable token transfer stories.

---

## Implementation Strategy

1. **MVP First**: Implement Sui→EVM token-only (US1) end-to-end before touching other directions.
2. **Incremental Delivery**: Add EVM→Sui token-only (US2), then Sui→EVM token+message (US3), then EVM→Sui token+message (US4).
3. **Config First**: The `[token]` and `[sui_receiver]` config sections are implemented in Phase 1 so all stories can use them.
4. **Runtime Resolution**: Token pool config and Sui receiver state are always resolved at runtime — no static addresses in run configs.
5. **No Chainlink Core**: All imports stay within `chainlink-sui`, `chainlink-ccip/gobindings`, and external dependencies already in `go.mod`.
6. **Manual Execution**: All `go test`, `go build`, and lint commands are run manually by the operator; agents only create/edit source files.
