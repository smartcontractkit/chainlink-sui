# Quickstart: Load Test Token Transfers

**Feature**: 002-load-test-token-transfers
**Date**: 2026-08-05

## Prerequisites

- Go 1.26.2+
- Funded wallets:
  - **Sui wallet**: Sui testnet/mainnet with sufficient SUI (for gas/fees) and CCIP-BnM tokens (for transfers)
  - **EVM wallet**: Sepolia (or other testnet) with sufficient ETH (for gas/fees) and CCIP-BnM ERC-20 tokens (for transfers)
- Contract addresses from a deployment run (see [feature 001 quickstart](../../001-sui-ccip-load-tests/quickstart.md) for base setup)
- CCIP-BnM token registered in the CCIP token admin registry on the source chain

## Configuration

### 1. Secrets (`.env.testnet`)

Same as message-only — no new fields needed:

```bash
SUI_PRIVATE_KEY=suiprivkey1...
EVM_PRIVATE_KEY=0x...
```

### 2. Addresses (`addresses-testnet.json`)

Same as message-only. The token pool addresses are resolved at runtime from the CCIP token admin registry — no additional entries needed.

### 3. Networks (`networks-testnet.yaml`)

Same as message-only.

### 4. Run Config (`runs/<name>.toml`)

Extended with optional `[token]` and `[sui_receiver]` sections.

#### Sui→EVM Token-Only

```toml
[run]
env = "testnet"
source_chain_selector = "9762610643973837292"
dest_chain_selector = "10344971235874465080"
message_count = 10
message_data = ""

[receiver]
address = "0x320f4dE5B482a8732fdEEf358e96D785a3C0De30"

[gas]
sui_gas_budget = 1000000000
evm_gas_limit = 100000
evm_callback_gas_limit = 0

[token]
coin_metadata_id = "0x331ce2ba0901fec09d863f0d4162ae29bae2898922e345f3e4cd356363ce3c1b"
amount = 1000000000
mode = "token_only"
```

#### Sui→EVM Token + Message

```toml
[run]
env = "testnet"
source_chain_selector = "9762610643973837292"
dest_chain_selector = "10344971235874465080"
message_count = 5
message_data = "hello with token"

[receiver]
address = "0x320f4dE5B482a8732fdEEf358e96D785a3C0De30"

[gas]
sui_gas_budget = 1000000000
evm_gas_limit = 100000
evm_callback_gas_limit = 200000

[token]
coin_metadata_id = "0x331ce2ba0901fec09d863f0d4162ae29bae2898922e345f3e4cd356363ce3c1b"
amount = 1000000000
mode = "token_and_message"
```

#### EVM→Sui Token-Only

```toml
[run]
env = "testnet"
source_chain_selector = "10344971235874465080"
dest_chain_selector = "9762610643973837292"
message_count = 10
message_data = ""

[receiver]
address = "0x0000000000000000000000000000000000000000000000000000000000000000"

[gas]
sui_gas_budget = 1000000000
evm_gas_limit = 100000
evm_callback_gas_limit = 0

[token]
token_address = "0x88A2d74F47a237a62e7A51cdDa67270CE381555e"
amount = 1000000000000000000
mode = "token_only"
```

#### EVM→Sui Token + Message

```toml
[run]
env = "testnet"
source_chain_selector = "10344971235874465080"
dest_chain_selector = "9762610643973837292"
message_count = 5
message_data = "hello with token"

[receiver]
address = "0x1810288cfa4414848aad6d0d11afcb95836255d195c114d902fd933e5965511e"

[gas]
sui_gas_budget = 1000000000
evm_gas_limit = 100000
evm_callback_gas_limit = 200000

[token]
token_address = "0x88A2d74F47a237a62e7A51cdDa67270CE381555e"
amount = 1000000000000000000
mode = "token_and_message"

[sui_receiver]
package_id = "0x1810288cfa4414848aad6d0d11afcb95836255d195c114d902fd933e5965511e"
```

## Running Tests

### Sui→EVM Token Transfer

```bash
cd integration-tests/load
go test -run TestSui2EVM -v --run-name sui-to-evm-token-run
```

### EVM→Sui Token Transfer

```bash
cd integration-tests/load
go test -run TestEVM2Sui -v --run-name evm-to-sui-token-run
```

## Validation Scenarios

### Scenario 1: Sui→EVM Token-Only — Single Transfer

**Prerequisites**: Sui wallet with ≥1 CCIP-BnM token (1e9 base units) and sufficient SUI for gas/fees.

1. Create a run config with `mode = "token_only"`, `message_count = 1`, `amount = 1000000000`.
2. Run: `go test -run TestSui2EVM -v --run-name <name>`
3. **Expected**: One message sent successfully. Results file contains one entry with `token_amount = "1000000000"` and `token_identifier` set to the coin metadata ID.

### Scenario 2: Sui→EVM Token-Only — Multiple Transfers

**Prerequisites**: Sui wallet with ≥10 CCIP-BnM tokens (10 × 1e9 base units).

1. Create a run config with `mode = "token_only"`, `message_count = 10`, `amount = 1000000000`.
2. Run: `go test -run TestSui2EVM -v --run-name <name>`
3. **Expected**: All 10 messages sent successfully. Results file contains 10 entries, each with unique message IDs.

### Scenario 3: Sui→EVM Token-Only — Insufficient Balance

**Prerequisites**: Sui wallet with <1 CCIP-BnM token.

1. Create a run config with `mode = "token_only"`, `message_count = 1`, `amount = 1000000000`.
2. Run: `go test -run TestSui2EVM -v --run-name <name>`
3. **Expected**: Test fails during coin pool preparation with error "insufficient token balance".

### Scenario 4: EVM→Sui Token-Only — Single Transfer

**Prerequisites**: EVM wallet with ≥1 CCIP-BnM ERC-20 token (1e18 base units) and sufficient ETH for gas/fees.

1. Create a run config with `mode = "token_only"`, `message_count = 1`, `amount = 1000000000000000000`.
2. Run: `go test -run TestEVM2Sui -v --run-name <name>`
3. **Expected**: One message sent successfully. ERC-20 approval happens once at start. Results file contains one entry.

### Scenario 5: EVM→Sui Token + Message — Single Transfer

**Prerequisites**: EVM wallet with ≥1 CCIP-BnM ERC-20 token, registered Sui dummy receiver.

1. Create a run config with `mode = "token_and_message"`, `message_count = 1`, `message_data = "hello"`, `evm_callback_gas_limit = 200000`, and `[sui_receiver]` with the receiver package ID.
2. Run: `go test -run TestEVM2Sui -v --run-name <name>`
3. **Expected**: One combined transfer sent successfully. SuiExtraArgsV1 includes `receiverObjectIds = [clock, receiverState]` and `tokenReceiver = receiverState`.

### Scenario 6: Backward Compatibility — Message-Only Run

**Prerequisites**: Existing message-only run config (no `[token]` section).

1. Run: `go test -run TestSui2EVM -v --run-name first-sui-to-evm-run`
2. **Expected**: Test runs in message-only mode, identical to pre-token-transfer behavior. Results file has no token fields.

### Scenario 7: Config Validation — Missing Amount

1. Create a run config with `[token]` section but no `amount` field.
2. Run: `go test -run TestSui2EVM -v --run-name <name>`
3. **Expected**: Test fails at config load with error "token.amount is required when [token] section is present".

### Scenario 8: Config Validation — Token+Message Without Data

1. Create a run config with `mode = "token_and_message"` but empty `message_data`.
2. Run: `go test -run TestSui2EVM -v --run-name <name>`
3. **Expected**: Test fails at config load with error "message_data is required for token_and_message mode".

## Expected Results Format

```json
{
  "run_name": "sui-to-evm-token-run",
  "env_name": "testnet",
  "source_chain_selector": 9762610643973837292,
  "dest_chain_selector": 10344971235874465080,
  "total_messages": 10,
  "successful_messages": 10,
  "failed_messages": 0,
  "run_started": "2026-08-05T17:00:00Z",
  "run_ended": "2026-08-05T17:05:00Z",
  "messages": [
    {
      "message_id": "0xabc...",
      "transaction_hash": "0xdef...",
      "source_chain_selector": 9762610643973837292,
      "dest_chain_selector": 10344971235874465080,
      "timestamp": "2026-08-05T17:00:05Z",
      "success": true,
      "sequence_number": "42",
      "token_amount": "1000000000",
      "token_identifier": "0x331ce2ba0901fec09d863f0d4162ae29bae2898922e345f3e4cd356363ce3c1b"
    }
  ]
}
```
