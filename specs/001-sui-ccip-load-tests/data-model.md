# Data Model: Sui CCIP Load Tests

## Four-Layer Configuration

```
Layer 1: .env.<env>          → Secrets (private keys)
Layer 2: addresses-<env>.json → Contract addresses (cldf AddressBook format)
Layer 3: networks-<env>.yaml  → Chain RPC endpoints (unified EVM + Sui)
Layer 4: runs/<name>.toml     → Run-specific parameters
```

## Entities

### RunConfig (Layer 4 — TOML file)

Loaded from `runs/<name>.toml`. The filename `<name>` is used for the results file.

| Field | Type | Source | Description |
|-------|------|--------|-------------|
| `Env` | `string` | TOML `[run].env` | Environment name (e.g., `testnet`, `mainnet`) |
| `SourceChainSelector` | `uint64` | TOML `[run].source_chain_selector` | Source chain selector |
| `DestChainSelector` | `uint64` | TOML `[run].dest_chain_selector` | Destination chain selector |
| `MessageCount` | `int` | TOML `[run].message_count` | Number of messages to send |
| `MessageData` | `string` | TOML `[run].message_data` | Message payload string |
| `ReceiverAddress` | `string` | TOML `[receiver].address` | Receiver address (EVM hex or Sui package ID) |
| `SuiGasBudget` | `uint64` | TOML `[gas].sui_gas_budget` | Sui PTB gas budget (Sui→EVM) |
| `EvmGasLimit` | `uint64` | TOML `[gas].evm_gas_limit` | EVM gas limit (EVM→Sui) |

### EnvConfig (Layer 1 — .env file)

| Field | Type | Description |
|-------|------|-------------|
| `SUI_PRIVATE_KEY` | `string` | Sui private key (bech32 `suiprivkey1...`) |
| `EVM_PRIVATE_KEY` | `string` | EVM private key (hex) |

### AddressBook (Layer 2 — addresses.json)

Uses `cldf.AddressBook` from `chainlink-deployments-framework/deployment`. Loaded via `cldf.NewMemoryAddressBookFromFile(path)`.

The `addresses-testnet.json` file is a flat array of `{address, chainSelector, type, version, qualifier, labels}` entries. `cldf` loads this into an `AddressBook` where `AddressesForChain(selector)` returns `map[string]cldf.TypeAndVersion`.

**Addresses needed for sending** (message-only, no tokens):
- For Sui→EVM: `SuiRouter` (router package ID), `SuiOnRamp` (onramp package ID), `SuiCCIP` (CCIP package ID), `SuiLinkTokenCoinMetadataID` (fee token metadata)
- For EVM→Sui: `Router` (EVM router address), `LinkToken` (fee token address)

### NetworkConfig (Layer 3 — YAML file)

Unified format for all chains (EVM and Sui). Each entry has `type`, `chain_selector`, and `rpcs`.

| Field | Type | Description |
|-------|------|-------------|
| `Type` | `string` | Network type (e.g., `testnet`, `mainnet`) |
| `ChainSelector` | `uint64` | Chain selector |
| `RPCs` | `[]RPCConfig` | List of RPC endpoints |

### RPCConfig

| Field | Type | Description |
|-------|------|-------------|
| `RPCName` | `string` | RPC provider name |
| `HTTPURL` | `string` | HTTP JSON-RPC URL |
| `WSURL` | `string` | WebSocket URL |

### SentMessage

A single sent CCIP message, recorded in the results file.

| Field | Type | Description |
|-------|------|-------------|
| `MessageID` | `string` | CCIP message ID (hex-encoded bytes32) |
| `TransactionHash` | `string` | Transaction hash (hex-encoded) |
| `SourceChainSelector` | `uint64` | Source chain selector |
| `DestChainSelector` | `uint64` | Destination chain selector |
| `Timestamp` | `string` | ISO 8601 timestamp |
| `Success` | `bool` | Whether the send succeeded |
| `Error` | `string` | Error message if failed (empty if success) |
| `SequenceNumber` | `string` | CCIP sequence number (if available) |

### RunResults

The output file for a load test run. Named `results/<run-name>-<env>-<timestamp>.txt`.

| Field | Type | Description |
|-------|------|-------------|
| `RunName` | `string` | Run config filename (without .toml) |
| `EnvName` | `string` | Environment name |
| `SourceChainSelector` | `uint64` | Source chain selector |
| `DestChainSelector` | `uint64` | Destination chain selector |
| `TotalMessages` | `int` | Total messages attempted |
| `SuccessfulMessages` | `int` | Messages sent successfully |
| `FailedMessages` | `int` | Messages that failed |
| `RunStarted` | `string` | ISO 8601 timestamp of run start |
| `RunEnded` | `string` | ISO 8601 timestamp of run end |
| `Messages` | `[]SentMessage` | List of all sent messages |

## Validation Rules

1. **Chain selector in address book**: The source chain selector must have entries in the address book.
2. **Address completeness**: For Sui→EVM, the address book must contain `SuiRouter` and `SuiOnRamp` for the source chain. For EVM→Sui, it must contain `Router` for the source chain.
3. **Network config match**: Both source and destination chain selectors must have entries in the YAML network config.
4. **Key format**: Sui private keys must be bech32 (`suiprivkey1...`). EVM private keys must be hex.
5. **Message count**: Must be >= 1.

## State Transitions

```
Load Run Config (TOML) → Load Env (.env) → Load Addresses (JSON) → Load Networks (YAML)
                                    ↓
[For each message 1..N] → Send Message → Extract Event → Record Result
                                    ↓ (on failure, no retry in v1)
                              Record Failure → Next Message
                                    ↓
                              Write Results File
```
