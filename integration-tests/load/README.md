# CCIP Load Tests

Lightweight Go load test framework for sending CCIP messages and token transfers between Sui and EVM chains against remote environments (testnet/staging/mainnet).

**No Chainlink core dependency.** No confirmation/metrics features in v1.

## Architecture

```
SUI → EVM (message-only):  Build PTB → create_token_transfer_params → ccip_send → Execute → Extract CCIPMessageSent
SUI → EVM (token):         Build PTB → create_token_transfer_params → managed_token_pool.lock_or_burn → ccip_send → Execute → Extract CCIPMessageSent
EVM → SUI (message-only):  GetFee → Router.ccipSend (native ETH) → Extract CCIPMessageSent from receipt
EVM → SUI (token):         Approve Router → GetFee → Router.ccipSend (native ETH + TokenAmounts + SuiExtraArgsV1) → Extract CCIPMessageSent from receipt
```

## Prerequisites

- Go 1.26.2+
- Funded wallets (Sui and EVM) with sufficient tokens for fees
- Contract addresses from a deployment run
- Network configuration for all chains

## Configuration (4 Layers)

| Layer | File | Content |
|-------|------|---------|
| 1 | `.env.<env>` | Secrets: `SUI_PRIVATE_KEY`, `EVM_PRIVATE_KEY`, optional `SUI_GRPC_TOKEN` |
| 2 | `addresses-<env>.json` | Contract addresses (cldf AddressBook format) |
| 3 | `networks-<env>.yaml` | Chain RPC endpoints (unified EVM + Sui) |
| 4 | `runs/<name>.toml` | Run-specific parameters |

### Layer 1: `.env.testnet`

```bash
SUI_PRIVATE_KEY=suiprivkey1...
EVM_PRIVATE_KEY=0x...
SUI_GRPC_TOKEN=optional-provider-grpc-token
```

### Layer 2: `addresses-testnet.json`

Flat array of `{address, chainSelector, type, version}` entries. Copy from deployment pipeline output.

### Layer 3: `networks-testnet.yaml`

```yaml
networks:
  - type: testnet
    chain_selector: 9762610643973837292  # Sui testnet
    rpcs:
      - rpc_name: Public
        http_url: https://fullnode.testnet.sui.io:443
        # optional for providers that require explicit gRPC endpoint host:port
        grpc_target: sui-testnet.g.alchemy.com:443
  - type: testnet
    chain_selector: 16015286601757825753  # Sepolia
    rpcs:
      - rpc_name: CLL Proxy
        http_url: https://rpcs.cldev.sh/16015286601757825753
```

### Layer 4: `runs/my-first-sui-to-evm-run.toml`

```toml
[run]
env = "testnet"
source_chain_selector = 9762610643973837292
dest_chain_selector = 16015286601757825753
message_count = 10
message_data = "hello from sui load test"

[receiver]
address = "0x0000000000000000000000000000000000000000"

[gas]
sui_gas_budget = 10000000000
evm_gas_limit = 200000
```

### Token Transfer Run Config

Add an optional `[token]` section to enable token transfers. For Sui source chains, use `coin_metadata_id`. For EVM source chains, use `token_address`.

```toml
[token]
coin_metadata_id = "0x331ce2ba0901fec09d863f0d4162ae29bae2898922e345f3e4cd356363ce3c1b"
amount = 1000000000
mode = "token_only"  # or "token_and_message"
```

For EVM→Sui `token_and_message` mode, add a `[sui_receiver]` section with the receiver package ID:

```toml
[sui_receiver]
package_id = "0x1810288cfa4414848aad6d0d11afcb95836255d195c114d902fd933e5965511e"
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

### Token Transfer Examples

```bash
cd integration-tests/load
go test -run TestSui2EVM -v --run-name sui-to-evm-token-only-run
go test -run TestEVM2Sui -v --run-name evm-to-sui-token-only-run
go test -run TestSui2EVM -v --run-name sui-to-evm-token-and-message-run
go test -run TestEVM2Sui -v --run-name evm-to-sui-token-and-message-run
```

## Results

Results are saved to `results/<runName>-<env>-<timestamp>.txt` in JSON format:

```json
{
  "run_name": "my-first-sui-to-evm-run",
  "env_name": "testnet",
  "total_messages": 10,
  "successful_messages": 10,
  "failed_messages": 0,
  "messages": [
    {
      "message_id": "0xabc...",
      "transaction_hash": "0xdef...",
      "success": true,
      "sequence_number": "42",
      "token_amount": "1000000000",
      "token_identifier": "0x331ce2ba0901fec09d863f0d4162ae29bae2898922e345f3e4cd356363ce3c1b"
    }
  ]
}
```

Token fields are only populated when a `[token]` section is present in the run config.

## Fee Handling

- **Sui→EVM**: Fees paid via PTB gas budget (native SUI). No LINK coin needed.
- **EVM→Sui**: Fees paid in native ETH. 20% buffer added to estimated fee.
- **Token transfers**: Token amounts are specified in base units (raw smallest units). No decimal conversion is performed at runtime.

## Constraints

- Messages sent sequentially (one at a time)
- No retries in v1 (will be added later)
- No confirmation waiting (DON handles execution)
- No Chainlink core imports

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        SUI → EVM                                │
│                                                                 │
│  Go script                                                      │
│    │                                                            │
│    ├─ 1. Build PTB (Programmable Transaction Block)             │
│    │      a. create_token_transfer_params                       │
│    │      b. managed_token_pool.lock_or_burn (token transfers)  │
│    │      c. ccip_send (onramp)                                 │
│    ├─ 2. Sign & execute PTB via Sui fullnode                    │
│    └─ 3. Extract CCIPMessageSent event → sequence number        │
│                                                                 │
│  Token pool config resolved at runtime from token admin registry│
│  Token coins pre-split into exact per-message amounts           │
│  DON (off-chain) handles commit + execution automatically       │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                        EVM → SUI                                │
│                                                                 │
│  Go script                                                      │
│    │                                                            │
│    ├─ 1. Approve Router for total ERC-20 amount (once)        │
│    ├─ 2. Call Router.getFee(destChainSelector, message)         │
│    ├─ 3. Call Router.ccipSend(destChainSelector, message)       │
│    │      with TokenAmounts and SuiExtraArgsV1                  │
│    └─ 4. Extract CCIPMessageSent event from OnRamp logs         │
│                                                                 │
│  SuiExtraArgsV1 tokenReceiver = EOA for token-only              │
│  SuiExtraArgsV1 tokenReceiver = receiverState for token+message │
│  DON (off-chain) handles commit + execution automatically       │
└─────────────────────────────────────────────────────────────────┘
```

---

## Prerequisites: Addresses & Configuration

You need the following addresses for each lane. Example values are from the
`sui-starter-kit/helperConfig.ts`.

### Sui-side addresses

| Field | Testnet Example | Description |
|---|---|---|
| `SuiRPC` | `https://fullnode.testnet.sui.io` | Sui JSON-RPC endpoint |
| `CCIPPackageID` | `0x5ef4b4...bdd6bf` | Original CCIP package (for object lookups) |
| `CCIPLatestPackageID` | (derived at runtime) | Latest CCIP package (for PTB calls) |
| `OnRampPackageID` | `0x01a0a2...a205a4` | Original OnRamp package (for object lookups) |
| `OnRampLatestPackageID` | (derived at runtime) | Latest OnRamp package (for PTB calls) |
| `CCIPObjectRef` | (derived at runtime) | `CCIPObjectRef` shared object ID |
| `OnRampState` | (derived at runtime) | `OnRampState` shared object ID |
| `LinkCoinMetadataID` | `0x7b1f2e...efc8dc` | LINK token metadata (for fee payment) |
| `LinkTokenPackageID` | (from metadata) | LINK token package (e.g. `0x...::link::LINK`) |
| `SuiPrivateKey` | (env) | Ed25519 private key (32-byte seed) |

### EVM-side addresses

| Field | Sepolia Example | Description |
|---|---|---|
| `EVM_RPC` | `https://sepolia...` | EVM JSON-RPC endpoint |
| `RouterAddress` | `0x0BF3dE...3A59` | CCIP Router contract |
| `OnRampAddress` | `0x23a508...9DeE` | CCIP OnRamp (for event filtering) |
| `LinkTokenAddress` | `0x779877...4789` | LINK ERC-677 token |
| `ChainSelector` | `16015286601757825753` | Destination chain selector |
| `EVMPrivateKey` | (env) | EVM private key |

---

## Step-by-Step: SUI → EVM (Message Only or Token Transfer)

Token transfers reuse the same PTB path with an additional `managed_token_pool.lock_or_burn` call between `create_token_transfer_params` and `ccip_send`. The token pool config is resolved at runtime from the CCIP token admin registry using the token's coin metadata ID.

### 1. Create the Sui Client

Use the `block-vision/sui-go-sdk` JSON-RPC client directly (no gRPC needed for sending):

```go
import (
    "net/http"
    "github.com/block-vision/sui-go-sdk/sui"
)

suiClient := sui.NewSuiClientWithCustomClient(rpcURL, &http.Client{Timeout: 30 * time.Second})
```

**Reusable from:** `chainlink-sui/relayer/client` — `PTBClient` wraps this and adds gRPC,
caching, rate-limiting. For a simple load script, the raw `sui.ISuiAPI` is sufficient.

### 2. Create the Sui Signer

Implement the `bindutils.SuiSigner` interface (from `chainlink-sui/bindings/utils/signer.go`):

```go
type SuiSigner interface {
    Sign(message []byte) ([]string, error)
    GetAddress() (string, error)
}
```

**Reusable from:** `chainlink-sui/bindings/utils/test_signer.go` — `TestPrivateKeySigner`
implements this with Ed25519 + Blake2b hashing + intent-scope prefix.

### 3. Resolve Latest Package IDs

CCIP packages may be upgraded. Always resolve the latest package ID before building a PTB:

```go
// Use suiClient.GetLatestPackageId or the PTBClient equivalent
latestCcipPkg, _ := getLatestPackageID(ctx, suiClient, ccipPackageID, "ccip")
latestOnRampPkg, _ := getLatestPackageID(ctx, suiClient, onRampPackageID, "ccip_onramp")
```

**Reusable from:** `chainlink-sui/relayer/client` — `BindingsClient.GetLatestPackageId()`.

### 4. Resolve Shared Object IDs

```go
// Get CCIPObjectRef from the CCIP package
ccipObjectRef, _ := getObjectFromPackage(ctx, suiClient, ccipPackageID, "CCIPObjectRef")

// Get OnRampState from the OnRamp package
onRampState, _ := getObjectFromPackage(ctx, suiClient, onRampPackageID, "OnRampState")
```

These are done via `suiClient.SuiGetOwnedObjects` filtered by struct type.

### 5. Build the PTB

The PTB has **two Move calls**:

#### 5a. `create_token_transfer_params` (onramp_state_helper)

Even for message-only transfers, you must call this with the normalized receiver bytes:

```go
import (
    suitx "github.com/block-vision/sui-go-sdk/transaction"
    suiBind "github.com/smartcontractkit/chainlink-sui/bindings/bind"
)

ptb := suitx.NewTransaction()

// Bind the onramp_state_helper contract
helperContract, _ := suiBind.NewBoundContract(
    latestCcipPkg, "ccip", "onramp_state_helper", ptbClient,
)

// Encode create_token_transfer_params with receiver bytes
encoded, _ := helperContract.EncodeCallArgsWithGenerics(
    "create_token_transfer_params",
    []string{},          // typeArgs
    []string{},          // typeParams
    []string{"vector<u8>"}, // paramTypes
    []any{receiverBytes}, // normalized 32-byte receiver
    nil,                 // returnTypes
)

tokenParamsResult, _ := helperContract.AppendPTB(ctx, callOpts, ptb, encoded)
```

#### 5b. `managed_token_pool.lock_or_burn` (token transfers only)

```go
poolContract, _ := managedtokenpool.NewManagedTokenPool(latestPoolPkgID, ptbClient)

encoded, _ := poolContract.Encoder().LockOrBurn(
    []string{coinType},
    suiBind.Object{Id: ccipObjectRef},
    tokenParamsResult,
    suiBind.Object{Id: tokenCoinID},
    destChainSelector,
    suiBind.Object{Id: "0x6"},
    suiBind.Object{Id: denyListObjectID},
    suiBind.Object{Id: tokenStateObjectID},
    suiBind.Object{Id: poolStateObjectID},
)

poolContract.Bound().AppendPTB(ctx, callOpts, ptb, encoded)
```

**Reusable from:** `integration-tests/load/sui/token_pool.go` —
`AppendManagedTokenPoolLockOrBurn()`.

#### 5c. `ccip_send` (onramp)

```go
// Bind the onramp contract
onRampContract, _ := onramp.NewOnramp(latestOnRampPkg, ptbClient)

encoded, _ := onRampContract.Encoder().CcipSendWithArgs(
    []string{feeTokenType},           // fee token type
    suiBind.Object{Id: ccipObjectRef}, // ccip state
    suiBind.Object{Id: onRampState},   // onramp state
    suiBind.Object{Id: "0x6"},        // clock
    destChainSelector,                // uint64
    receiver,                         // []byte (32-byte EVM address)
    data,                             // []byte (message payload)
    tokenParamsResult,                // from step 5a
    suiBind.Object{Id: feeTokenMetadataID}, // fee token metadata
    suiBind.Object{Id: feeTokenCoinID},    // fee token coin
    extraArgs,                        // []byte
)

onRampContract.Bound().AppendPTB(ctx, callOpts, ptb, encoded)
```

**Reusable from:** `integration-tests/load/sui/sender.go` —
`SendMessage()` and `SendTokenMessage()`.

### 6. Construct ExtraArgs (GenericExtraArgsV2)

For Sui→EVM, use the **GenericExtraArgsV2** format:

```go
import (
    "math/big"
    "github.com/ethereum/go-ethereum/common/hexutil"
)

const GenericExtraArgsV2Tag = "0x181dcf10"

func MakeBCSEVMExtraArgsV2(gasLimit *big.Int, allowOOO bool) []byte {
    // BCS-encode: u256 gasLimit + bool allowOOO
    // gasLimit as 32-byte little-endian
    glBytes := make([]byte, 32)
    gasLimit.FillBytes(glBytes) // big-endian fill
    // Reverse to little-endian
    for i, j := 0, 31; i < j; i, j = i+1, j-1 {
        glBytes[i], glBytes[j] = glBytes[j], glBytes[i]
    }
    oooByte := byte(0)
    if allowOOO {
        oooByte = 1
    }
    payload := append(glBytes, oooByte)
    return append(hexutil.MustDecode(GenericExtraArgsV2Tag), payload...)
}
```

**Reusable from:** `chainlink/deployment/ccip/changeset/testhelpers/test_helpers_aptos.go` —
`MakeBCSEVMExtraArgsV2()`.

### 7. Execute the PTB

```go
resp, err := suiBind.ExecutePTB(ctx, callOpts, ptbClient, ptb)
```

**Reusable from:** `chainlink-sui/bindings/bind/call.go` — `ExecutePTB()`.

### 8. Extract the CCIPMessageSent Event

```go
for _, event := range resp.Events {
    if event.PackageId == latestOnRampPkg && strings.HasSuffix(event.Type, "CCIPMessageSent") {
        seqStr := event.ParsedJson["sequence_number"].(string)
        seq, _ := strconv.ParseUint(seqStr, 10, 64)
        // seq is your CCIP sequence number
    }
}
```

### 9. Fee Handling

- **Native SUI fees:** Use `ptb.SetGasBudget(budget)` — the fee is deducted from the gas coin.
- **Token transfer fees:** The fee is quoted via `onramp.get_fee` with token types and amounts included, then a 20% buffer is applied.

The fee can be queried via `onramp.get_fee` (devInspect) before building the PTB.

---

## Step-by-Step: EVM → SUI (Message Only or Token Transfer)

Token transfers populate `TokenAmounts` in `ClientEVM2AnyMessage` and set `TokenReceiver` in `SuiExtraArgsV1`. For token-only transfers to a Sui EOA, the receiver is `ZeroHash` and `tokenReceiver` is the EOA address. For token+message transfers to a Sui receiver object, `tokenReceiver` is the receiver state object ID and `ReceiverObjectIds` includes both clock and receiver state.

### 1. Create the EVM Client

```go
import (
    "github.com/ethereum/go-ethereum/ethclient"
)

ethClient, _ := ethclient.Dial(evmRPCURL)
```

### 2. Create the TransactOpts (Signer)

```go
import (
    "crypto/ecdsa"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
)

privateKey, _ := crypto.HexToECDSA(privKeyHex)
chainID, _ := ethClient.ChainID(ctx)
auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
```

### 3. Instantiate the Router Contract

```go
import (
    "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
)

routerContract, _ := router.NewRouter(routerAddress, ethClient)
```

**Reusable from:** `chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router`.

### 4. Construct the Message

```go
msg := router.ClientEVM2AnyMessage{
    Receiver:     receiverBytes,    // 32-byte Sui receiver package ID
    Data:         messageData,      // arbitrary bytes
    TokenAmounts: []router.ClientEVMTokenAmount{}, // empty = message-only
    FeeToken:     common.Address{}, // 0x0 for native ETH
    ExtraArgs:    extraArgsBytes,   // SuiExtraArgsV1
}
```

For token transfers, populate `TokenAmounts`:

```go
msg.TokenAmounts = []router.ClientEVMTokenAmount{
    {Token: tokenAddress, Amount: tokenAmount},
}
```

### 5. Construct ExtraArgs (SuiExtraArgsV1)

For EVM→Sui, use the **SuiExtraArgsV1** format:

```go
import (
    "math/big"
    "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/message_hasher"
)

extraArgs, _ := evm.SerializeClientSUIExtraArgsV1(message_hasher.ClientSuiExtraArgsV1{
    GasLimit:                 big.NewInt(100_000),
    AllowOutOfOrderExecution: true,
    TokenReceiver:            [32]byte{}, // zero for message-only
    ReceiverObjectIds: [][32]byte{
        clockObjectID,        // 0x0000...0006
        receiverStateObjID,   // CCIPReceiverState from the receiver package
    },
})
```

**Reusable from:** `integration-tests/load/evm/extras.go` —
`SerializeClientSUIExtraArgsV1()`.

For token-only transfers, set `TokenReceiver` to the Sui EOA address (as bytes32).
For token+message transfers, set `TokenReceiver` to the receiver state object ID and
include both clock and receiver state in `ReceiverObjectIds`.

### 6. Get the Fee

```go
fee, _ := routerContract.GetFee(&bind.CallOpts{}, destChainSelector, msg)
```

### 7. Approve ERC-20 Tokens (for token transfers)

```go
// Approve Router for total token amount once at start
totalAmount := new(big.Int).Mul(tokenAmount, big.NewInt(int64(messageCount)))
evmapproveRouterForTokens(ctx, ethClient, auth, tokenAddress, routerAddress, totalAmount)
```

The approval is skipped if the existing allowance is already sufficient.

### 8. Send the Message

```go
// For native ETH fees:
auth.Value = feeWithBuffer
tx, _ := routerContract.CcipSend(auth, destChainSelector, msg)
```

Token amounts are included in `msg.TokenAmounts`. The fee token remains native ETH
in v1.

**Reusable from:** `chainlink/deployment/ccip/changeset/testhelpers/test_helpers_solana_v0_1_0.go` —
`CCIPSendRequest()` and `retryCcipSendUntilNativeFeeIsSufficient()`.

### 9. Extract the CCIPMessageSent Event

```go
receipt, _ := ethClient.TransactionReceipt(ctx, tx.Hash())
messageID, seqNum, _ := evm.ExtractMessageIDFromReceipt(receipt, destChainSelector)
```

**Reusable from:** `integration-tests/load/evm/sender.go` —
`ExtractMessageIDFromReceipt()`.

---

## Key Reusable Packages

| Package | Import Path | What It Provides |
|---|---|---|
| Sui bindings | `chainlink-sui/bindings/bind` | `BoundContract`, `ExecutePTB`, `CallOpts`, `NewBoundContract` |
| Sui signer utils | `chainlink-sui/bindings/utils` | `SuiSigner` interface, `TestPrivateKeySigner` |
| Sui PTB client | `chainlink-sui/relayer/client` | `PTBClient`, `BindingsClient`, `PTBClientConfig` |
| Sui offramp helpers | `chainlink-sui/relayer/chainwriter/ptb/offramp` | `DecodeParameters` (for reading Move function signatures) |
| EVM Router | `chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router` | `Router`, `ClientEVM2AnyMessage` |
| EVM Message Hasher | `chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/message_hasher` | `ClientSuiExtraArgsV1`, `ClientGenericExtraArgsV2` |
| ExtraArgs serialization | `chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/message_hasher` | `SerializeClientSUIExtraArgsV1()` (vendored in `evm/extras.go`) |
| GenericExtraArgsV2 | `integration-tests/load/sui/extras.go` | `MakeBCSEVMExtraArgsV2()`, `GenericExtraArgsV2Tag` |
| Sui Go SDK | `github.com/block-vision/sui-go-sdk/sui` | `ISuiAPI` (JSON-RPC), `NewSuiClientWithCustomClient` |
| Sui Go SDK tx | `github.com/block-vision/sui-go-sdk/transaction` | `NewTransaction()` (PTB builder) |
| Token admin registry | `chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry` | `GetTokenConfigStruct()` |
| Managed token pool | `chainlink-sui/bindings/generated/ccip/ccip_token_pools/managed_token_pool` | `LockOrBurn()` |

---

## Minimal Dependencies for a Standalone Script

If you want a script that depends **only** on `chainlink-sui` and `chainlink-ccip`
(not the full `chainlink` monorepo), you need to:

1. **Vendor or copy** `MakeBCSEVMExtraArgsV2` (~10 lines, just BCS-encodes a u256 + bool with a 4-byte tag).
2. **Vendor or copy** `SerializeClientSUIExtraArgsV1` (uses `message_hasher.ClientSuiExtraArgsV1` from `chainlink-ccip` — already available).
3. **Use `chainlink-sui/bindings/bind`** for PTB construction and execution.
4. **Use `chainlink-sui/bindings/utils`** for the `SuiSigner` interface and `TestPrivateKeySigner`.
5. **Use `chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router`** for EVM Router interactions.
6. **Use `chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry`** to resolve token pool config at runtime.
7. **Use `chainlink-sui/bindings/generated/ccip/ccip_token_pools/managed_token_pool`** for `lock_or_burn` calls.

The `GenericExtraArgsV2Tag` constant is `"0x181dcf10"` (4 bytes).

---

## Environment Variables

```
# Sui
SUI_RPC_URL=https://fullnode.testnet.sui.io
SUI_PRIVATE_KEY=<32-byte hex-encoded Ed25519 seed>

# EVM
EVM_RPC_URL=https://sepolia.infura.io/v3/...
EVM_PRIVATE_KEY=<hex-encoded ECDSA private key>
```

No new environment variables are required for token transfers.

---

## Notes

- **No DON interaction needed:** The script only sends messages. The DON (off-chain
  CCIP network) automatically picks up `CCIPMessageSent` events, commits them on the
  destination chain, and executes them.
- **Wait for finality:** After sending, the message will be committed and executed
  asynchronously. You can poll the destination chain's OffRamp for execution state
  changes.
- **Fee buffer:** Always add ~20% to the quoted fee to account for gas price
  fluctuations between quote and execution.
- **Sui fullnode indexing lag:** After a Sui transaction, the fullnode may need a few
  hundred milliseconds to index the new objects. The `WaitForTransactionIndexed` helper
  in `chainlink-sui/bindings/bind` handles this.
- **Receiver requirements on Sui:** For EVM→Sui messages, the receiver must be a
  package that implements the `ccip_receive` interface. The `ReceiverObjectIds` in
  `SuiExtraArgsV1` must include the `Clock` (0x6) and the receiver's state object.
- **Token pool config:** The Sui token pool package ID, state object ID, deny list,
  and token state are resolved at runtime from the CCIP token admin registry using
  the token's coin metadata ID. No manual pool addresses are required in the run config.
- **Token amounts:** All token amounts in the run config are in base units (raw
  smallest units). Operators must convert from human-readable amounts using the token's
  decimals.
- **ERC-20 allowance:** For EVM→Sui token transfers, the Router is approved once at
  the start of the run for the total token amount. The approval is skipped if the
  existing allowance is already sufficient.
