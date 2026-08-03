# Feature Specification: Sui CCIP Load Tests

**Feature Branch**: `sui-load-tests`

**Created**: 2026-08-03

**Status**: Draft

**Input**: User description: "I need to implement load tests under integration-tests/load using simple golang scripts (ie: test files to be run manually). The tests will run against remote envs only (staging/testnet/mainnet), no need to start any devenv. This test must not depend on Chainlink core (/chainlink). chainlink/integration-tests/load can be used as a base for message sending, but it requires implementing a bunch of support cause it reuses frameworks and libs. Our work here should be as minimal as possible, since sending messages is mostly calling the router on both ends. So the ts sdk /sui-starter-kit is a better example, since it's a standalone app. We won't add any confirmation/metrics features for now, just a framework that is able to send messages and log the messageIDs (maybe save a file with the run results - just the sends ids and transactions hashes). It's going to be configured using 3 layers: Layer 1: .env files — Secrets only; Layer 2: addresses.json — All contract addresses; Layer 3: Network config (optional YAML/embedded)."

## Clarifications

### Session 2026-08-03

- Q: How should the test handle sending multiple messages — sequentially or concurrently? → A: Sequential sends (one at a time, wait for each to complete before sending the next). Operators who need higher throughput can run multiple test process instances with different wallets.
- Q: When a single message send fails, what should the test do? → A: Retry up to 3 times (simple retry, no exponential backoff for now), then log the error, record the failure in results, and continue to the next message.
- Q: For EVM→Sui message fees, which fee token approach should the test support in v1? → A: Native ETH only for v1. Support for LINK and other fee tokens will be added later.
- Q: For Sui→EVM messages, how should fees be paid? → A: Native SUI gas budget only for v1. Set PTB gas budget with a buffer. Support for LINK and other fee tokens deferred.
- Q: For the Layer 3 network config (EVM chain parameters), should EVM chain configs be embedded as Go constants or loaded from a separate YAML file? → A: YAML file. Format: top-level `networks` array, each entry with `type`, `chain_selector`, `rpcs` (array of `rpc_name`, `http_url`, `ws_url`), and optional `metadata`.
- Q: How should the test determine which Sui chain to target — derived from env name or explicit flag? → A: Both. Env name (e.g., `--env testnet`) loads the corresponding `.env` and YAML files. Explicit source and destination chain selector flags override auto-detection, since both source and destination need to be configured.
- Q: Should agents run Go commands (go build, go test, go mod tidy) during implementation? → A: No. All Go commands will be run manually by the user. Agents must only create/edit source files.

## User Scenarios & Testing

### User Story 1 - Run a Sui→EVM message-only load test (Priority: P1)

An operator wants to send a batch of CCIP messages from a Sui chain (testnet/mainnet) to an EVM destination chain (e.g., Sepolia, Arbitrum Sepolia). The operator configures the environment via a `.env` file (Sui private key, Sui RPC URL), an `addresses.json` file (contract addresses from a deployment run), and a network config (EVM RPC URLs, chain selectors, router addresses). The operator runs a Go test file that sends N messages, logs each message ID and transaction hash, and saves the results to a file.

**Why this priority**: This is the primary direction for CCIP load testing — Sui as the source chain is the core use case. It covers the most critical path: constructing a Sui transaction that calls the OnRamp to send a CCIP message, and extracting the resulting message identifier.

**Independent Test**: Can be fully tested by running the Go test against Sui testnet with a valid `.env.testnet` and `addresses.json`. The test sends at least one message and verifies that a message ID and transaction hash are logged and written to the results file.

**Acceptance Scenarios**:

1. **Given** a configured `.env.testnet` file with `SUI_PRIVATE_KEY` and `SUI_RPC_URL`, **When** the operator runs the Sui→EVM load test with 1 message, **Then** the test successfully sends the message, logs the message ID and transaction hash, and writes them to the results file.
2. **Given** a configured environment with valid addresses, **When** the operator runs the test with N=10 messages, **Then** all 10 messages are sent successfully, each with a unique message ID and transaction hash recorded.
3. **Given** an invalid Sui private key in `.env`, **When** the operator runs the test, **Then** the test fails with a clear error message indicating the key is invalid.

---

### User Story 2 - Run an EVM→Sui message-only load test (Priority: P1)

An operator wants to send a batch of CCIP messages from an EVM chain (e.g., Sepolia) to a Sui destination chain (testnet/mainnet). The operator configures the EVM side via `.env` (EVM private key, EVM RPC URL), `addresses.json` (Sui contract addresses), and network config (EVM router address, chain selectors). The operator runs a Go test file that sends N messages, logs each message ID and transaction hash, and saves the results to a file.

**Why this priority**: The reverse direction (EVM→Sui) is equally important for comprehensive load testing. It covers calling the EVM Router to send a CCIP message, handling fee estimation, and extracting the resulting message identifier from the transaction logs.

**Independent Test**: Can be fully tested by running the Go test against Sepolia testnet with a valid `.env.testnet` and `addresses.json`. The test sends at least one message and verifies that a message ID and transaction hash are logged and written to the results file.

**Acceptance Scenarios**:

1. **Given** a configured `.env.testnet` file with `EVM_PRIVATE_KEY` and `EVM_RPC_URL`, **When** the operator runs the EVM→Sui load test with 1 message, **Then** the test successfully sends the message, logs the message ID and transaction hash, and writes them to the results file.
2. **Given** a configured environment with valid addresses, **When** the operator runs the test with N=10 messages, **Then** all 10 messages are sent successfully, each with a unique message ID and transaction hash recorded.
3. **Given** insufficient LINK token balance for fees, **When** the operator runs the test, **Then** the test fails with a clear error about insufficient funds.

---

### User Story 3 - Configure a load test run via three-layer config (Priority: P2)

An operator needs to set up a load test for a new environment (e.g., staging, a different testnet, or mainnet). The operator creates a `.env.<env>` file with secrets, copies the `addresses.json` from the deployment pipeline, and creates or updates a network config YAML with chain-specific parameters. The operator then runs the test by specifying the environment name.

**Why this priority**: Configuration is a prerequisite for all test runs. This story ensures the three-layer config system works end-to-end before any test execution.

**Independent Test**: Can be tested by creating a `.env.testnet` file, a valid `addresses.json`, and a network config, then running a test that loads all three layers and prints the resolved configuration without sending any messages.

**Acceptance Scenarios**:

1. **Given** a `.env.testnet` file with `SUI_PRIVATE_KEY` and `SUI_RPC_URL`, **When** the config loader reads the file, **Then** the private key and RPC URL are available in the test configuration.
2. **Given** an `addresses.json` file in the deployment framework format, **When** the config loader reads the file, **Then** all contract addresses (SuiCCIP, SuiRouter, SuiOffRamp, SuiLinkTokenCoinMetadataID) are parsed correctly by chain selector.
3. **Given** a network config YAML with EVM chain parameters, **When** the config loader reads the file, **Then** EVM RPC URLs, chain selectors, and router addresses are available.
4. **Given** a missing `.env` file, **When** the config loader runs, **Then** it fails with a clear error indicating the missing file.

---

### User Story 4 - Save run results to a file (Priority: P2)

After a load test run completes (or during the run), the test framework writes the results to a file. Each result entry contains the message ID and transaction hash for each sent message. The file format is simple and machine-readable (e.g., JSON or CSV).

**Why this priority**: Persisting results enables post-run analysis, debugging, and sharing with other tools. Without this, the operator would need to manually capture output from logs.

**Independent Test**: Can be tested by running a test that sends 1+ messages and verifying that the results file is created with the correct content.

**Acceptance Scenarios**:

1. **Given** a completed load test run with 5 messages, **When** the test finishes, **Then** a results file exists with 5 entries, each containing a message ID and transaction hash.
2. **Given** a load test run that fails mid-way, **When** the test encounters an error, **Then** all successfully sent messages up to that point are already written to the results file.

### Edge Cases

- What happens when the Sui RPC endpoint is unreachable or returns errors?
- How does the system handle rate limiting from the Sui fullnode or EVM RPC?
- What happens when the LINK token balance is insufficient for EVM fees?
- How does the system handle duplicate transaction submissions (nonce management on EVM)?
- What happens when the `addresses.json` file is missing entries for a required contract type?
- How does the system behave when the destination chain selector in the config does not match any entry in `addresses.json`?

## Requirements

### Functional Requirements

- **FR-001**: The system MUST provide a Go test file that sends CCIP messages from Sui to EVM by calling the Sui OnRamp contract.
- **FR-002**: The system MUST provide a Go test file that sends CCIP messages from EVM to Sui by calling the EVM Router contract.
- **FR-003**: The system MUST load secrets (private keys, RPC URLs) from a `.env` file, with one `.env` file per environment (e.g., `.env.testnet`, `.env.mainnet`).
- **FR-004**: The system MUST load contract addresses from an `addresses.json` file that uses the same format as the deployment framework (keyed by chain selector, then by address, with `Type` and `Version` fields).
- **FR-005**: The system MUST load EVM chain network configuration from a YAML file. The YAML format MUST contain a top-level `networks` array, where each entry has: `type` (string), `chain_selector` (uint64), `rpcs` (array of objects with `rpc_name`, `http_url`, `ws_url`), and optional `metadata`.
- **FR-006**: The system MUST log each message ID and transaction hash to stdout during the test run.
- **FR-007**: The system MUST save run results (message IDs and transaction hashes) to a file in a machine-readable format (e.g., JSON) after the test completes.
- **FR-008**: The system MUST support configurable message count (N messages to send) via a command-line flag or environment variable. Messages MUST be sent sequentially (one at a time, waiting for each to complete before sending the next).
- **FR-016**: When a single message send fails, the system MUST retry up to 3 times (simple retry without exponential backoff). If all retries fail, the system MUST log the error, record the failure in the results file, and continue to the next message.
- **FR-009**: The system MUST support configurable message data payload via a command-line flag or environment variable.
- **FR-017**: The system MUST support an `--env` flag to select the environment (e.g., `testnet`, `mainnet`), which determines which `.env` and YAML config files are loaded. The system MUST also support explicit `--source-chain-selector` and `--dest-chain-selector` flags to override auto-detected chain selectors.
- **FR-010**: The system MUST NOT depend on the Chainlink core repository for any imports.
- **FR-011**: The system MUST handle Sui private keys in the standard Sui bech32 format by decoding them to the raw key material.
- **FR-012**: The system MUST handle EVM private keys in standard hex format.
- **FR-013**: The system MUST resolve the latest deployed package versions for Sui contracts before sending messages, to handle contract upgrades.
- **FR-014**: The system MUST construct the correct extra arguments format for each direction (Sui→EVM and EVM→Sui) as required by the CCIP protocol.
- **FR-015**: The system MUST query the CCIP fee before sending. For EVM→Sui messages, the system MUST pay fees in native ETH only (v1 scope). For Sui→EVM messages, the system MUST pay fees via the Sui PTB gas budget (native SUI) only (v1 scope). Support for LINK and other fee tokens is deferred to a future iteration.

### Key Entities

- **LoadTestRun**: Represents a single load test execution. Contains the environment name, message count, source and destination chain selectors, and a list of sent messages.
- **SentMessage**: Represents a single sent CCIP message. Contains the message ID, transaction hash, source chain selector, destination chain selector, and timestamp.
- **EnvironmentConfig**: The three-layer configuration for a test run. Composed of secrets (`.env`), contract addresses (`addresses.json`), and network parameters (YAML/constants).
- **AddressBook**: A mapping from chain selector to contract addresses, loaded from `addresses.json`. Each address has a `Type` (e.g., `SuiRouter`, `SuiOnRamp`) and `Version`.
- **RunResults**: The output of a test run. A JSON file containing an array of `SentMessage` entries, plus metadata (run timestamp, environment name, total messages sent).

## Success Criteria

### Measurable Outcomes

- **SC-001**: An operator can configure and run a Sui→EVM load test with 100 messages in under 5 minutes, with all message IDs and transaction hashes recorded to a results file.
- **SC-002**: An operator can configure and run an EVM→Sui load test with 100 messages in under 5 minutes, with all message IDs and transaction hashes recorded to a results file.
- **SC-003**: The three-layer configuration system (`.env` + `addresses.json` + network config) can be fully set up for a new environment in under 10 minutes by an operator who has the deployment addresses.
- **SC-004**: A load test run produces a results file that can be consumed by external tools (e.g., for analysis or reporting) without any transformation.
- **SC-005**: The load test framework compiles and runs without importing any package from the Chainlink core repository.

## Assumptions

- The operator has access to a deployed CCIP environment (testnet/staging/mainnet) with all required contract addresses.
- The operator has funded wallets (Sui and EVM) with sufficient tokens to pay CCIP fees.
- The Sui and EVM RPC endpoints are accessible from the machine running the tests.
- The `addresses.json` file is copied directly from the deployment pipeline output — no format translation is needed.
- Message-only transfers (no tokens) are the primary use case; token transfer support is out of scope for v1.
- EVM→Sui fees are paid in native ETH only in v1; LINK and other fee token support is deferred.
- No confirmation or metrics features (e.g., waiting for message execution on destination, Grafana dashboards) are included in v1.
- The existing `chainlink-sui/integration-tests` Go module will be used, with additional dependencies added as needed.
- Standard Go libraries for blockchain interactions (Sui and EVM clients) are already available in the module dependencies.
- A `.env` file loading mechanism will be added as a new dependency.
- All Go commands (go build, go test, go mod tidy, etc.) MUST be run manually by the operator. Agents MUST NOT execute any Go commands — they only create and edit source files.
