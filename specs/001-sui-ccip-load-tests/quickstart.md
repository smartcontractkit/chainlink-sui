# Quickstart: Sui CCIP Load Tests

## Prerequisites

- Go 1.26.2+
- Access to a deployed CCIP environment (testnet/staging/mainnet)
- Funded wallets (Sui and EVM) with sufficient tokens for fees
- Contract addresses from a deployment run (`addresses.json`)
- Network configuration for EVM chains

## Setup

### 1. Create `.env.testnet`

```bash
# integration-tests/load/.env.testnet
SUI_PRIVATE_KEY=suiprivkey1...
SUI_RPC_URL=https://fullnode.testnet.sui.io:443
EVM_PRIVATE_KEY=0x...
```

### 2. Copy `addresses.json`

Copy the deployment pipeline output to `integration-tests/load/addresses-testnet.json`.

The file supports two formats:
- **Nested format** (Sui addresses from deployment): `map[chainSelector]map[address]{Type, Version}`
- **Flat format** (EVM addresses): `[{address, chainSelector, type, version}]`

### 3. Create `networks-testnet.yaml`

```yaml
# integration-tests/load/networks-testnet.yaml
networks:
  - type: testnet
    chain_selector: 16015286601757825753  # Sepolia
    rpcs:
      - rpc_name: CLL Proxy
        http_url: https://rpcs.cldev.sh/16015286601757825753
        ws_url: wss://rpcs.cldev.sh/16015286601757825753
```

## Running Tests

### Sui → EVM

```bash
cd integration-tests/load

go test -run TestSui2EVM -v \
  --env testnet \
  --source-chain-selector 9762610643973837292 \
  --dest-chain-selector 16015286601757825753 \
  --count 10 \
  --data "hello from sui"
```

### EVM → Sui

```bash
cd integration-tests/load

go test -run TestEVM2Sui -v \
  --env testnet \
  --source-chain-selector 16015286601757825753 \
  --dest-chain-selector 9762610643973837292 \
  --count 10 \
  --data "hello from evm"
```

### Config Validation Only (no messages sent)

```bash
cd integration-tests/load

go test -run TestConfig -v --env testnet
```

## Expected Output

### Console Logs

```
=== RUN   TestSui2EVM
[2026-08-03T12:00:00Z] Loading config for environment: testnet
[2026-08-03T12:00:01Z] Sending message 1/10...
[2026-08-03T12:00:05Z]   ✓ Message sent | ID: 0xabc... | TX: 0xdef... | Seq: 42
[2026-08-03T12:00:06Z] Sending message 2/10...
[2026-08-03T12:00:10Z]   ✓ Message sent | ID: 0x123... | TX: 0x456... | Seq: 43
...
[2026-08-03T12:01:00Z] Run complete: 10/10 successful, 0 failed
[2026-08-03T12:01:00Z] Results saved to: results/testnet-20260803T120000.json
```

### Results File (`results/testnet-20260803T120000.json`)

```json
{
  "env_name": "testnet",
  "source_chain_selector": 9762610643973837292,
  "dest_chain_selector": 16015286601757825753,
  "total_messages": 10,
  "successful_messages": 10,
  "failed_messages": 0,
  "run_started": "2026-08-03T12:00:00Z",
  "run_ended": "2026-08-03T12:01:00Z",
  "messages": [
    {
      "message_id": "0xabc...",
      "transaction_hash": "0xdef...",
      "source_chain_selector": 9762610643973837292,
      "dest_chain_selector": 16015286601757825753,
      "timestamp": "2026-08-03T12:00:05Z",
      "success": true,
      "error": "",
      "sequence_number": "42"
    }
  ]
}
```

## Validation Scenarios

### Scenario 1: Single message Sui→EVM

1. Configure `.env.testnet`, `addresses-testnet.json`, `networks-testnet.yaml`
2. Run: `go test -run TestSui2EVM -v --env testnet --count 1`
3. **Expected**: 1 message sent, message ID and TX hash logged, results file created

### Scenario 2: Multiple messages EVM→Sui

1. Same config as above
2. Run: `go test -run TestEVM2Sui -v --env testnet --count 10`
3. **Expected**: 10 messages sent sequentially, all logged, results file with 10 entries

### Scenario 3: Config validation

1. Run: `go test -run TestConfig -v --env testnet`
2. **Expected**: Config loaded and printed, no messages sent, no error

### Scenario 4: Invalid private key

1. Set `SUI_PRIVATE_KEY=invalid` in `.env.testnet`
2. Run: `go test -run TestSui2EVM -v --env testnet --count 1`
3. **Expected**: Test fails with clear error about invalid key format

## Architecture Reference

- [Data Model](data-model.md) — Entity definitions and validation rules
- [Config Contract](contracts/config.md) — Config package types and functions
- [Sui Contract](contracts/sui.md) — Sui client, PTB construction, event extraction
- [EVM Contract](contracts/evm.md) — EVM client, message sending, event extraction
