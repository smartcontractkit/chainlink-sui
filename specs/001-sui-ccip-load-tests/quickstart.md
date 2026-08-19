# Quickstart: Sui CCIP Load Tests

## Prerequisites

- Go 1.26.2+
- Access to a deployed CCIP environment (testnet/staging/mainnet)
- Funded wallets (Sui and EVM) with sufficient tokens for fees
- Contract addresses from a deployment run (`addresses-<env>.json`)
- Network configuration (`networks-<env>.yaml`)

## Setup

### 1. Create `.env.testnet`

```bash
# integration-tests/load/.env.testnet
SUI_PRIVATE_KEY=suiprivkey1...
EVM_PRIVATE_KEY=0x...
```

### 2. Copy `addresses-testnet.json`

Copy the deployment pipeline output to `integration-tests/load/addresses-testnet.json`. Uses the standard `cldf.AddressBook` format (flat array of `{address, chainSelector, type, version}` entries).

### 3. Create `networks-testnet.yaml`

```yaml
# integration-tests/load/networks-testnet.yaml
networks:
  - type: testnet
    chain_selector: 9762610643973837292  # Sui testnet
    rpcs:
      - rpc_name: Public
        http_url: https://fullnode.testnet.sui.io:443
  - type: testnet
    chain_selector: 16015286601757825753  # Sepolia
    rpcs:
      - rpc_name: CLL Proxy
        http_url: https://rpcs.cldev.sh/16015286601757825753
```

### 4. Create a run config TOML

```toml
# integration-tests/load/runs/my-first-sui-to-evm-run.toml
[run]
env = "testnet"
source_chain_selector = 9762610643973837292
dest_chain_selector = 16015286601757825753
message_count = 10
message_data = "hello from sui load test"

[receiver]
address = "0x..."  # 32-byte EVM receiver (left-padded with zeros)

[gas]
sui_gas_budget = 10000000000
evm_gas_limit = 200000
```

## Running Tests

### Sui → EVM

```bash
cd integration-tests/load

go test -run TestSui2EVM -v --run-name my-first-sui-to-evm-run
```

### EVM → Sui

```bash
cd integration-tests/load

go test -run TestEVM2Sui -v --run-name my-first-evm-to-sui-run
```

### Config Validation Only (no messages sent)

```bash
cd integration-tests/load

go test -run TestConfig -v --run-name my-first-sui-to-evm-run
```

## Expected Output

### Console Logs

```
=== RUN   TestSui2EVM
[2026-08-03T12:00:00Z] Loading run config: runs/my-first-sui-to-evm-run.toml
[2026-08-03T12:00:01Z] Sending message 1/10...
[2026-08-03T12:00:05Z]   ✓ Message sent | ID: 0xabc... | TX: 0xdef... | Seq: 42
[2026-08-03T12:00:06Z] Sending message 2/10...
[2026-08-03T12:00:10Z]   ✓ Message sent | ID: 0x123... | TX: 0x456... | Seq: 43
...
[2026-08-03T12:01:00Z] Run complete: 10/10 successful, 0 failed
[2026-08-03T12:01:00Z] Results saved to: results/my-first-sui-to-evm-run-testnet-20260803T120000.txt
```

### Results File (`results/my-first-sui-to-evm-run-testnet-20260803T120000.txt`)

```json
{
  "run_name": "my-first-sui-to-evm-run",
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
