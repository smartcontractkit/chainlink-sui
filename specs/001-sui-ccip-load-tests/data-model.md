# Data Model: Sui CCIP Load Tests

## Entities

### LoadTestConfig

The top-level configuration for a load test run, assembled from the three config layers.

| Field | Type | Source | Description |
|-------|------|--------|-------------|
| `EnvName` | `string` | CLI flag `--env` | Environment name (e.g., `testnet`, `mainnet`) |
| `SourceChainSelector` | `uint64` | CLI flag or auto-detect | Source chain selector |
| `DestChainSelector` | `uint64` | CLI flag or auto-detect | Destination chain selector |
| `MessageCount` | `int` | CLI flag or env var | Number of messages to send (default: 1) |
| `MessageData` | `[]byte` | CLI flag or env var | Arbitrary message payload |
| `SuiRPCURL` | `string` | `.env` file | Sui JSON-RPC endpoint |
| `SuiPrivateKey` | `string` | `.env` file | Sui private key (bech32 format) |
| `EVMPrivateKey` | `string` | `.env` file | EVM private key (hex format) |
| `SuiAddresses` | `AddressBook` | `addresses.json` | Sui contract addresses by chain selector |
| `EVMAddresses` | `[]FlatAddress` | `addresses.json` | EVM contract addresses |
| `EVMNetworks` | `[]NetworkConfig` | YAML file | EVM network configurations |

### AddressBook

A mapping from chain selector to contract addresses, loaded from the nested `addresses.json` format (deployment pipeline output).

```
map[chainSelector]map[address]TypeAndVersion
```

### TypeAndVersion

| Field | Type | Description |
|-------|------|-------------|
| `Type` | `string` | Contract type (e.g., `SuiRouter`, `SuiOnRamp`, `SuiCCIP`) |
| `Version` | `string` | Semantic version (e.g., `1.0.0`) |
| `Labels` | `map[string]struct{}` | Optional labels (e.g., `{"LINK": {}}`) |

### FlatAddress

An EVM address entry from the flat `addresses.json` format.

| Field | Type | Description |
|-------|------|-------------|
| `Address` | `string` | Contract address (hex with 0x prefix) |
| `ChainSelector` | `uint64` | Chain selector for this address |
| `Type` | `string` | Contract type (e.g., `Router`, `OnRamp`, `CommitStore`) |
| `Version` | `string` | Version string |
| `Qualifier` | `string` | Optional qualifier/label |
| `Labels` | `[]string` | Optional labels |

### NetworkConfig

An EVM network configuration entry from the YAML file.

| Field | Type | Description |
|-------|------|-------------|
| `Type` | `string` | Network type (e.g., `testnet`, `mainnet`) |
| `ChainSelector` | `uint64` | Chain selector |
| `RPCs` | `[]RPCConfig` | List of RPC endpoints |
| `Metadata` | `map[string]any` | Optional metadata |

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
| `Timestamp` | `string` | ISO 8601 timestamp of when the message was sent |
| `Success` | `bool` | Whether the send succeeded |
| `Error` | `string` | Error message if failed (empty if success) |
| `SequenceNumber` | `string` | CCIP sequence number (if available) |

### RunResults

The output file for a load test run.

| Field | Type | Description |
|-------|------|-------------|
| `EnvName` | `string` | Environment name |
| `SourceChainSelector` | `uint64` | Source chain selector |
| `DestChainSelector` | `uint64` | Destination chain selector |
| `TotalMessages` | `int` | Total messages attempted |
| `SuccessfulMessages` | `int` | Messages sent successfully |
| `FailedMessages` | `int` | Messages that failed after retries |
| `RunStarted` | `string` | ISO 8601 timestamp of run start |
| `RunEnded` | `string` | ISO 8601 timestamp of run end |
| `Messages` | `[]SentMessage` | List of all sent messages |

## Validation Rules

1. **Chain selector consistency**: The source chain selector in the config must match a chain selector present in `addresses.json` for Sui contract types.
2. **Address completeness**: For Sui→EVM, the `addresses.json` must contain `SuiRouter`, `SuiOnRamp`, `SuiCCIP`, and `SuiLinkTokenCoinMetadataID` entries for the source chain selector.
3. **Network config match**: The destination chain selector must match a `chain_selector` in the YAML network config.
4. **Key format**: Sui private keys must be in bech32 format (`suiprivkey1...`). EVM private keys must be hex (with or without `0x` prefix).
5. **Message count**: Must be >= 1.

## State Transitions

```
Config Loaded → [For each message 1..N] → Send Message → Extract Event → Record Result
                                                    ↓ (on failure, retry up to 3x)
                                              Retry or Record Failure
                                                    ↓
                                              Next Message or Done
                                                    ↓
                                              Write Results File
```
