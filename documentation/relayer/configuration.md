# Relayer Configuration

The Chainlink SUI Relayer uses a comprehensive TOML-based configuration system that provides fine-grained control over all aspects of blockchain interaction. This guide covers the complete configuration structure, best practices, and advanced optimization strategies.

## Table of Contents

1. [Configuration Overview](#configuration-overview)
2. [Basic Configuration](#basic-configuration)
3. [Advanced Configuration](#advanced-configuration)
4. [Environment Variables](#environment-variables)
5. [Configuration Validation](#configuration-validation)
6. [Best Practices](#best-practices)
7. [Production Tuning](#production-tuning)
8. [Troubleshooting](#troubleshooting)

## Configuration Overview

The relayer configuration is structured around chains, with each chain having its own set of nodes, transaction management settings, and component-specific configurations.

### Configuration Hierarchy

```
TOMLConfig
├── ChainID
├── Enabled
├── NetworkName
├── NetworkNameFull
├── Nodes[]
├── TransactionManager
├── BalanceMonitor
├── ChainReader
└── ChainWriter
```

## Basic Configuration

### Minimal Configuration

Here's the minimum configuration required to run the relayer:

```toml
[[Chains]]
ChainID = '0x1'
Enabled = true
NetworkName = 'sui-mainnet'
NetworkNameFull = 'Sui Mainnet'

[[Chains.Nodes]]
Name = 'sui-mainnet-1'
URL = 'https://fullnode.mainnet.sui.io'
```

### Development Configuration

For local development and testing:

```toml
[[Chains]]
ChainID = '0x1'
Enabled = true
NetworkName = 'sui-localnet'
NetworkNameFull = 'Sui Local Network'

[[Chains.Nodes]]
Name = 'local-sui'
URL = 'http://localhost:9000'

# Development-friendly settings
[Chains.TransactionManager]
BroadcastChanSize = 100
ConfirmPollSecs = 1
DefaultMaxGasAmount = 10000000
MaxTxRetryAttempts = 3
TransactionTimeout = "30s"
MaxConcurrentRequests = 5
RequestType = "WaitForLocalExecution"

[Chains.BalanceMonitor]
Enabled = true
BalancePollPeriod = "30s"
```

## Advanced Configuration

### Complete Configuration Example

```toml
# Chain Configuration
[[Chains]]
ChainID = '0x1'
Enabled = true
NetworkName = 'sui-mainnet'
NetworkNameFull = 'Sui Mainnet'

# Node Configuration
[[Chains.Nodes]]
Name = 'sui-mainnet-primary'
URL = 'https://fullnode.mainnet.sui.io'

[[Chains.Nodes]]
Name = 'sui-mainnet-backup'
URL = 'https://sui-mainnet-rpc.nodereal.io'

# Transaction Manager Configuration
[Chains.TransactionManager]
BroadcastChanSize = 1000
ConfirmPollSecs = 2
DefaultMaxGasAmount = 200000000
MaxTxRetryAttempts = 5
TransactionTimeout = "60s"
MaxConcurrentRequests = 10
RequestType = "WaitForEffectsCert"

# Balance Monitor Configuration  
[Chains.BalanceMonitor]
Enabled = true
BalancePollPeriod = "60s"

# ChainPoller (checkpoint-based indexing)
[Chains.ChainPoller]
PollingIntervalSecs = 2
SyncTimeoutSecs = 60
ChannelBufferSize = 16
BackfillCheckpointCount = 100

# ChainWriter Configuration
[Chains.ChainWriter]
GasLimit = 200000000
MaxRetries = 3
TransactionTimeout = "45s"
```

## Configuration Sections

### Chain Settings

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `ChainID` | string | required | Unique identifier for the Sui chain |
| `Enabled` | bool | `true` | Whether this chain configuration is active |
| `NetworkName` | string | required | Short name for the network |
| `NetworkNameFull` | string | required | Full descriptive name for the network |

### Node Configuration

Each chain can have multiple nodes for redundancy and load balancing. Because the PTB Client now talks to the node over **gRPC**, each node may also specify a gRPC target and auth token:

```toml
[[Chains.Nodes]]
Name = 'primary-node'
URL = 'https://fullnode.mainnet.sui.io'
GrpcTarget = 'sui-mainnet-grpc.example.org:443'  # gRPC endpoint (host:port)
GrpcToken = 'your-api-key'                        # gRPC auth token
```

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `Name` | string | required | Unique name for the node |
| `URL` | string | required | JSON-RPC endpoint URL |
| `GrpcTarget` | string | - | gRPC endpoint (`host:port`); enables gRPC when set together with `GrpcToken` |
| `GrpcToken` | string | - | gRPC auth token |

> The gRPC client is only initialized when **both** `GrpcTarget` and `GrpcToken` are present. Without them the client runs in JSON-RPC-only mode. The checkpoint-based indexer (ChainPoller) relies on gRPC, so production deployments should configure these.

### Transaction Manager

The Transaction Manager handles transaction lifecycle, retries, and confirmation:

```toml
[Chains.TransactionManager]
BroadcastChanSize = 1000
ConfirmPollSecs = 2
DefaultMaxGasAmount = 200000000
MaxTxRetryAttempts = 5
TransactionTimeout = "60s"
MaxConcurrentRequests = 10
RequestType = "WaitForEffectsCert"
```

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `BroadcastChanSize` | uint64 | `4096` | Size of the transaction broadcast channel |
| `ConfirmPollSecs` | uint64 | `10` | Interval for polling transaction confirmations |
| `DefaultMaxGasAmount` | uint64 | `10000000` | Default gas limit for transactions |
| `MaxTxRetryAttempts` | uint64 | `5` | Maximum retry attempts for failed transactions |
| `TransactionTimeout` | string | `"120s"` | Timeout for individual transactions |
| `MaxConcurrentRequests` | uint64 | `500` | Maximum concurrent RPC requests |
| `RequestType` | string | `"WaitForEffectsCert"` | Transaction request type |
| `ReaperPollSecs` | uint64 | `10` | How often the reaper runs |
| `TransactionRetentionSecs` | uint64 | `10` | How long to retain finalized/failed transactions |

#### Request Types

- `WaitForEffectsCert`: Wait for transaction effects certificate (recommended for production)
- `WaitForLocalExecution`: Wait for local execution only (faster, suitable for testing)

### Balance Monitor

Monitors account balances and provides alerts when balances are low:

```toml
[Chains.BalanceMonitor]
Enabled = true
BalancePollPeriod = "60s"
```

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `Enabled` | bool | `true` | Enable balance monitoring |
| `BalancePollPeriod` | string | `"10s"` | Interval for checking balances |

### ChainPoller Configuration

The indexing pipeline is driven by a single **ChainPoller** that streams checkpoints over gRPC and fans their contents out to the Events and Transactions indexers (see [Event Indexing](./event-indexing.md)). It replaces the old per-selector / per-transmitter polling, so checkpoint fetching is configured here rather than under separate `EventsIndexer` / `TransactionsIndexer` polling settings:

```toml
[Chains.ChainPoller]
PollingIntervalSecs = 2
SyncTimeoutSecs = 60
ChannelBufferSize = 16
BackfillCheckpointCount = 100
ReplayCheckpointCount = 100         # checkpoints re-scanned per Replay request (0 = to latest)
# StartCheckpointSequence = 12345   # optional; overrides backfill when set
```

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `PollingIntervalSecs` | uint64 | `2` | How often live polling checks for new checkpoints |
| `SyncTimeoutSecs` | uint64 | `60` | Timeout for processing a single checkpoint |
| `ChannelBufferSize` | uint64 | `16` | Buffer size of the events/transactions channels |
| `BackfillCheckpointCount` | uint64 | `100` | Start at `latest - N` |
| `StartCheckpointSequence` | uint64 | _unset_ | Explicit starting checkpoint sequence (overrides backfill) |
| `ReplayCheckpointCount` | uint64 | `100` | Checkpoints re-scanned per `Replay` request, starting at the requested checkpoint; an explicit `0` re-scans to the latest checkpoint |

> **Note**: In the relayer's own TOML the chain block is keyed `[[Sui]]` with `[[Sui.Nodes]]`, `[Sui.ChainPoller]`, `[Sui.TransactionManager]`, etc. The `[[Chains]]` form used elsewhere in this document is illustrative.

### ChainWriter Configuration

Controls transaction submission behavior:

```toml
[Chains.ChainWriter]
GasLimit = 200000000
MaxRetries = 3
TransactionTimeout = "45s"
```

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `GasLimit` | int64 | `10000000` | Default gas limit for transactions |
| `MaxRetries` | int | `3` | Maximum retry attempts |
| `TransactionTimeout` | string | `"30s"` | Timeout for transaction submission |

## Environment Variables

The relayer supports environment variable overrides for sensitive configuration:

```bash
export CHAINLINK_SUI_CONFIG=/path/to/config.toml
export DATABASE_URL=postgresql://user:pass@localhost:5432/chainlink
export LOG_LEVEL=info
export KEYSTORE_PATH=/path/to/sui.keystore
```

### Required Environment Variables

| Variable | Description |
|----------|-------------|
| `CHAINLINK_SUI_CONFIG` | Path to the TOML configuration file |
| `DATABASE_URL` | PostgreSQL connection string |

### Optional Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `LOG_LEVEL` | `info` | Logging level (debug, info, warn, error) |
| `KEYSTORE_PATH` | `~/.sui/sui_config/sui.keystore` | Path to Sui keystore |
| `RPC_TIMEOUT` | `30s` | Default RPC timeout |

## Configuration Validation

The relayer validates configuration on startup and provides detailed error messages:

### Validation Rules

1. **Chain ID**: Must be a valid hex string or decimal number
2. **Node URLs**: Must be valid HTTP/HTTPS URLs
3. **Timeouts**: Must be valid duration strings (e.g., "30s", "5m")
4. **Numeric Values**: Must be positive integers within reasonable bounds

### Example Validation Error

```
invalid tron config: Chain[0].TransactionManager.TransactionTimeout: invalid duration format: "invalid-timeout"
```

## Best Practices

### Network-Specific Configurations

#### Mainnet Configuration

```toml
[[Chains]]
ChainID = '0x1'
NetworkName = 'sui-mainnet'
NetworkNameFull = 'Sui Mainnet'

# Use multiple reliable nodes
[[Chains.Nodes]]
Name = 'sui-mainnet-primary'
URL = 'https://fullnode.mainnet.sui.io'

[[Chains.Nodes]]
Name = 'sui-mainnet-backup'
URL = 'https://sui-mainnet-rpc.nodereal.io'

# Conservative settings for production
[Chains.TransactionManager]
BroadcastChanSize = 2000
ConfirmPollSecs = 3
DefaultMaxGasAmount = 300000000
MaxTxRetryAttempts = 7
TransactionTimeout = "120s"
MaxConcurrentRequests = 15
RequestType = "WaitForEffectsCert"

# Monitor balances frequently
[Chains.BalanceMonitor]
Enabled = true
BalancePollPeriod = "30s"
```

#### Testnet Configuration

```toml
[[Chains]]
ChainID = '0x2'
NetworkName = 'sui-testnet'
NetworkNameFull = 'Sui Testnet'

[[Chains.Nodes]]
Name = 'sui-testnet'
URL = 'https://fullnode.testnet.sui.io'

# Faster settings for testing
[Chains.TransactionManager]
BroadcastChanSize = 500
ConfirmPollSecs = 1
DefaultMaxGasAmount = 100000000
MaxTxRetryAttempts = 3
TransactionTimeout = "60s"
MaxConcurrentRequests = 5
RequestType = "WaitForLocalExecution"
```

### Security Best Practices

1. **Use Environment Variables**: Store sensitive data in environment variables
2. **Multiple Nodes**: Configure multiple RPC endpoints for redundancy
3. **Reasonable Timeouts**: Set appropriate timeouts to prevent hanging
4. **Resource Limits**: Configure appropriate channel sizes and concurrency limits

### Performance Optimization

1. **Channel Sizing**: Size broadcast channels based on expected transaction volume
2. **Concurrency Limits**: Balance concurrency with node rate limits
3. **Polling Intervals**: Optimize polling intervals based on finality requirements
4. **Gas Management**: Set appropriate gas limits for your use case

## Production Tuning

### High-Throughput Configuration

For applications requiring high transaction throughput:

```toml
[Chains.TransactionManager]
BroadcastChanSize = 5000
ConfirmPollSecs = 1
MaxConcurrentRequests = 25
TransactionTimeout = "180s"

[Chains.ChainPoller]
PollingIntervalSecs = 1
SyncTimeoutSecs = 60
```

### Low-Latency Configuration

For applications requiring fast transaction confirmation:

```toml
[Chains.TransactionManager]
ConfirmPollSecs = 1
RequestType = "WaitForLocalExecution"
TransactionTimeout = "30s"
MaxConcurrentRequests = 10

[Chains.ChainPoller]
PollingIntervalSecs = 1
```

### Resource-Constrained Configuration

For environments with limited resources:

```toml
[Chains.TransactionManager]
BroadcastChanSize = 100
MaxConcurrentRequests = 3
TransactionTimeout = "60s"

[Chains.ChainPoller]
PollingIntervalSecs = 5
SyncTimeoutSecs = 30
```

## Configuration Templates

### Multi-Chain Configuration

```toml
# Mainnet
[[Chains]]
ChainID = '0x1'
Enabled = true
NetworkName = 'sui-mainnet'
NetworkNameFull = 'Sui Mainnet'

[[Chains.Nodes]]
Name = 'mainnet-primary'
URL = 'https://fullnode.mainnet.sui.io'

[Chains.TransactionManager]
BroadcastChanSize = 2000
ConfirmPollSecs = 3
DefaultMaxGasAmount = 300000000

# Testnet
[[Chains]]
ChainID = '0x2'
Enabled = true
NetworkName = 'sui-testnet'
NetworkNameFull = 'Sui Testnet'

[[Chains.Nodes]]
Name = 'testnet-primary'
URL = 'https://fullnode.testnet.sui.io'

[Chains.TransactionManager]
BroadcastChanSize = 500
ConfirmPollSecs = 1
DefaultMaxGasAmount = 100000000
```

## Troubleshooting

### Common Configuration Issues

| Issue | Symptom | Solution |
|-------|---------|----------|
| **Invalid Chain ID** | `couldn't parse chain id` error | Use valid hex (0x1) or decimal format |
| **Connection Failures** | RPC timeout errors | Verify node URLs are accessible |
| **High Memory Usage** | Out of memory errors | Reduce BroadcastChanSize and MaxConcurrentRequests |
| **Slow Confirmations** | Transactions pending too long | Decrease ConfirmPollSecs, increase timeout |
| **Transaction Failures** | High failure rate | Increase MaxTxRetryAttempts and gas limits |

### Validation Commands

Check configuration validity:

```bash
# Dry run to validate configuration
./chainlink-sui --config config.toml --validate-only

# Check specific chain configuration
./chainlink-sui --config config.toml --chain-id 0x1 --validate
```

### Debug Configuration

Enable debug logging to troubleshoot configuration issues:

```toml
# Add to environment
LOG_LEVEL=debug

# Or in configuration
[Chains.Debug]
LogLevel = "debug"
EnableConfigDump = true
```

## Migration Guide

### Upgrading Configuration

When upgrading the relayer, configuration may need updates:

```bash
# Backup existing configuration
cp config.toml config.toml.backup

# Update configuration format
./chainlink-sui --migrate-config config.toml.backup > config.toml

# Validate new configuration
./chainlink-sui --config config.toml --validate-only
```

### Configuration Version History

- **v1.0**: Basic chain and node configuration
- **v1.1**: Added transaction manager settings
- **v1.2**: Added balance monitor configuration
- **v1.3**: Added ChainReader/ChainWriter specific settings

## Examples

### Complete Production Configuration

```toml
# Production-ready configuration for Sui Mainnet
[[Chains]]
ChainID = '0x1'
Enabled = true
NetworkName = 'sui-mainnet'
NetworkNameFull = 'Sui Mainnet'

# Primary RPC endpoints
[[Chains.Nodes]]
Name = 'sui-mainnet-primary'
URL = 'https://fullnode.mainnet.sui.io'

[[Chains.Nodes]]
Name = 'sui-mainnet-nodereal'
URL = 'https://sui-mainnet-rpc.nodereal.io'

[[Chains.Nodes]]
Name = 'sui-mainnet-ankr'
URL = 'https://rpc.ankr.com/sui'

# Transaction management
[Chains.TransactionManager]
BroadcastChanSize = 3000
ConfirmPollSecs = 2
DefaultMaxGasAmount = 250000000
MaxTxRetryAttempts = 6
TransactionTimeout = "90s"
MaxConcurrentRequests = 20
RequestType = "WaitForEffectsCert"

# Balance monitoring
[Chains.BalanceMonitor]
Enabled = true
BalancePollPeriod = "45s"

# Checkpoint-based indexing
[Chains.ChainPoller]
PollingIntervalSecs = 2
SyncTimeoutSecs = 60
ChannelBufferSize = 16
BackfillCheckpointCount = 100

# Transaction submission
[Chains.ChainWriter]
GasLimit = 250000000
MaxRetries = 4
TransactionTimeout = "60s"
```

This configuration provides a robust foundation for production deployments with appropriate redundancy, error handling, and performance characteristics. 