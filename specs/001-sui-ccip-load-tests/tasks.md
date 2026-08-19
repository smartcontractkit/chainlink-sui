---

description: "Task list for Sui CCIP Load Tests feature implementation"
---

# Tasks: Sui CCIP Load Tests

**Input**: Design documents from `specs/001-sui-ccip-load-tests/`

**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: No test tasks — these are manual-run load tests, not automated CI tests.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

All source code lives under `integration-tests/load/` in the `chainlink-sui` repository.

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Create directory structure and sample config files

- [x] T001 Create directory structure: `config/`, `sui/`, `evm/`, `runs/`, `results/` under `integration-tests/load/`
- [x] T002 Create sample `.env.testnet` with `SUI_PRIVATE_KEY` and `EVM_PRIVATE_KEY` placeholders in `integration-tests/load/.env.testnet`
- [x] T003 Create sample run config `runs/my-first-sui-to-evm-run.toml` with `[run]`, `[receiver]`, `[gas]` sections in `integration-tests/load/runs/`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Config Package

- [x] T004 [P] Create config types (`RunConfig`, `LoadTestConfig`, `NetworkConfig`, `RPCConfig`, `SentMessage`, `RunResults`) in `integration-tests/load/config/types.go`
- [x] T005 [P] Implement TOML run config parser (`LoadRunConfig`) in `integration-tests/load/config/runconfig.go` — parses `runs/<name>.toml` using `github.com/pelletier/go-toml/v2`
- [x] T006 [P] Implement `.env` loader (`LoadEnvConfig`) in `integration-tests/load/config/config.go` — loads secrets using `github.com/joho/godotenv`
- [x] T007 [P] Implement YAML network config parser (`LoadNetworks`) in `integration-tests/load/config/config.go` — parses `networks-<env>.yaml` using `gopkg.in/yaml.v3`
- [x] T008 [P] Implement address book loader (`LoadAddressBook`) in `integration-tests/load/config/config.go` — loads `addresses-<env>.json` using `cldf.NewMemoryAddressBookFromFile()`
- [x] T009 Implement `LoadFullConfig(runName)` in `integration-tests/load/config/config.go` — assembles all 4 layers (TOML → .env → addresses.json → YAML), resolves RPC URLs for source/dest chains from network config, resolves contract addresses from address book

### Sui Client & Signer

- [x] T010 [P] Implement Sui client + signer setup in `integration-tests/load/sui/client.go`:
  - `NewSuiClient(rpcURL)` — creates `client.SuiPTBClient` from RPC URL
  - `NewSuiSigner(bech32PrivKey)` — decodes bech32 `suiprivkey1...` to Ed25519 seed, creates `relayer/signer.NewPrivateKeySigner`, returns `bindutils.SuiSigner` + sender address

### EVM Client & Signer

- [x] T011 [P] Implement EVM client + signer setup in `integration-tests/load/evm/client.go`:
  - `NewEVMClient(rpcURL)` — creates `*ethclient.Client`
  - `NewEVMSigner(privateKeyHex, chainID)` — creates `*bind.TransactOpts` from hex private key

### Extra Args Encoding

- [x] T012 [P] Implement `MakeBCSEVMExtraArgsV2` in `integration-tests/load/sui/extras.go` — BCS-encodes GenericExtraArgsV2 (tag `0x181dcf10` + 32-byte LE gasLimit + 1-byte bool) for Sui→EVM messages
- [x] T013 [P] Implement `SerializeClientSUIExtraArgsV1` in `integration-tests/load/evm/extras.go` — ABI-encodes SuiExtraArgsV1 (tag `0x21ea4ca9` + ABI-packed `ClientSuiExtraArgsV1`) using `message_hasher.MessageHasherABI` for EVM→Sui messages

**Checkpoint**: Foundation ready — user story implementation can now begin

---

## Phase 3: User Story 1 — Run a Sui→EVM message-only load test (Priority: P1) 🎯 MVP

**Goal**: Send CCIP messages from Sui to EVM by calling the Sui OnRamp contract, log message IDs and transaction hashes, save results.

**Independent Test**: Run `go test -run TestSui2EVM -v --run-name my-first-sui-to-evm-run` against Sui testnet with valid `.env.testnet`, `addresses-testnet.json`, and `networks-testnet.yaml`. Verify at least one message is sent and results file is created.

### Implementation for User Story 1

- [x] T014 [US1] Implement `SendMessage` in `integration-tests/load/sui/sender.go`:
  - Resolve latest package IDs via `ptbClient.GetLatestPackageId()`
  - Build PTB with `transaction.NewTransaction()`
  - Create `bind.NewBoundContract` for `ccip::onramp_state_helper` and `ccip_onramp::onramp`
  - Call `create_token_transfer_params` with empty 32-byte receiver
  - Call `ccip_send` with CCIP state object, OnRamp state, clock, dest chain selector, receiver, data, token params, LINK coin metadata, gas coin, extra args
  - Execute with `bind.ExecutePTB(ctx, opts, ptbClient, ptb)`
  - Extract `CCIPMessageSent` event from response events
  - Return message ID, transaction digest, sequence number
- [x] T015 [US1] Implement `TestSui2EVM` in `integration-tests/load/sui2evm_test.go`:
  - Load full config via `config.LoadFullConfig(runName)`
  - Create Sui client + signer
  - Resolve Sui addresses from address book (`SuiRouter`, `SuiOnRamp`, `SuiCCIP`, `SuiLinkTokenCoinMetadataID`)
  - Resolve EVM RPC URL from network config for dest chain
  - Loop N times: call `sui.SendMessage`, log result, append to results
  - Save results via `config.SaveResults`

**Checkpoint**: Sui→EVM message sending works end-to-end

---

## Phase 4: User Story 2 — Run an EVM→Sui message-only load test (Priority: P1)

**Goal**: Send CCIP messages from EVM to Sui by calling the EVM Router contract, log message IDs and transaction hashes, save results.

**Independent Test**: Run `go test -run TestEVM2Sui -v --run-name my-first-evm-to-sui-run` against Sepolia testnet with valid config. Verify at least one message is sent and results file is created.

### Implementation for User Story 2

- [x] T016 [US2] Implement `SendMessage` + `GetFee` + `ExtractMessageIDFromReceipt` in `integration-tests/load/evm/sender.go`:
  - `GetFee`: instantiate `router.NewRouter`, call `Router.GetFee(destChainSelector, msg)`
  - `SendMessage`: build `router.ClientEVM2AnyMessage` (empty `TokenAmounts`, zero `FeeToken` for native ETH), add 20% fee buffer to `auth.Value`, call `Router.CcipSend(auth, destChainSelector, msg)`
  - `ExtractMessageIDFromReceipt`: get receipt from tx hash, filter `RouterCCIPMessageSent` events, extract `Message.Header.MessageId` and `SequenceNumber`
- [x] T017 [US2] Implement `TestEVM2Sui` in `integration-tests/load/evm2sui_test.go`:
  - Load full config via `config.LoadFullConfig(runName)`
  - Create EVM client + signer (resolve chain ID from client)
  - Resolve EVM router address from address book
  - Resolve Sui RPC URL from network config (for reference, not needed for sending)
  - Construct `SuiExtraArgsV1` with clock object ID (`0x6`) and receiver state object
  - Loop N times: call `evm.SendMessage`, log result, append to results
  - Save results via `config.SaveResults`

**Checkpoint**: EVM→Sui message sending works end-to-end

---

## Phase 5: User Story 4 — Save run results to a file (Priority: P2)

**Goal**: Persist run results (message IDs, transaction hashes) to a machine-readable file.

**Independent Test**: Run any test with 1+ messages and verify the results file is created with correct content.

### Implementation for User Story 4

- [x] T018 [US4] Implement `SaveResults` in `integration-tests/load/config/config.go`:
  - Generate filename: `results/<runName>-<env>-<timestamp>.txt`
  - Marshal `RunResults` to JSON with indentation
  - Write to file, creating `results/` directory if needed
  - Log the file path on completion
- [x] T019 [US4] Integrate results saving into `sui2evm_test.go` and `evm2sui_test.go`:
  - Track `RunStarted`/`RunEnded` timestamps
  - Append each `SentMessage` to results during the loop
  - Call `SaveResults` after the loop completes (or on panic/defer)

**Checkpoint**: All test runs produce results files

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, sample configs, and final cleanup

- [x] T020 Update `integration-tests/load/README.md` with:
  - Architecture overview (Sui→EVM and EVM→Sui flows)
  - Prerequisites (Go 1.26.2+, funded wallets, deployment addresses)
  - Setup instructions (4 config layers)
  - Run commands for both directions
  - Expected output format
  - No Chainlink core dependency note
- [x] T021 Create sample run configs in `integration-tests/load/runs/`:
  - `my-first-sui-to-evm-run.toml` — Sui testnet → Sepolia
  - `my-first-evm-to-sui-run.toml` — Sepolia → Sui testnet

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion — BLOCKS all user stories
- **US1 — Sui→EVM (Phase 3)**: Depends on Foundational (T009, T010, T012)
- **US2 — EVM→Sui (Phase 4)**: Depends on Foundational (T009, T011, T013)
- **US4 — Results (Phase 5)**: Depends on Foundational (T009) — can run in parallel with US1/US2
- **Polish (Phase 6)**: Depends on all user stories being complete

### User Story Dependencies

- **US1 (P1)**: Can start after Foundational — No dependencies on other stories
- **US2 (P1)**: Can start after Foundational — No dependencies on other stories
- **US4 (P2)**: Can start after Foundational — No dependencies on other stories

### Within Each Phase

- Models/types before services
- Services before test files
- Core implementation before integration

### Parallel Opportunities

- **Phase 2**: T004, T005, T006, T007, T008, T010, T011, T012, T013 can all run in parallel (different files, no cross-dependencies)
- **Phase 3+**: US1, US2, and US4 can all be implemented in parallel once Foundational is complete
- **Phase 6**: T020 and T021 can run in parallel

---

## Parallel Example: Phase 2 Foundational

```text
Week 1:
  Dev A: T004 (types.go)
  Dev B: T005 (runconfig.go)
  Dev C: T006 + T007 + T008 (config loaders)
  Dev D: T010 (sui/client.go) + T012 (sui/extras.go)
  Dev E: T011 (evm/client.go) + T013 (evm/extras.go)

Week 2 (after T009 merge):
  Dev A: T014 + T015 (US1 — Sui→EVM)
  Dev B: T016 + T017 (US2 — EVM→Sui)
  Dev C: T018 + T019 (US4 — Results)
```

## Implementation Strategy

### MVP Scope (Phase 1 + 2 + 3)

The minimum viable product is **US1 only** (Sui→EVM message sending). This requires:
- Phase 1 (Setup): directory structure + sample configs
- Phase 2 (Foundational): config package, Sui client/signer, GenericExtraArgsV2
- Phase 3 (US1): Sui sender + test file

Once US1 works, US2 (EVM→Sui) and US4 (Results) can be added incrementally.

### Incremental Delivery Order

1. Config types + loaders (T004–T009)
2. Sui client + signer (T010)
3. GenericExtraArgsV2 (T012)
4. Sui sender + test (T014–T015) → **MVP complete**
5. EVM client + signer (T011)
6. SuiExtraArgsV1 (T013)
7. EVM sender + test (T016–T017)
8. Results persistence (T018–T019)
9. README + sample configs (T020–T021)
