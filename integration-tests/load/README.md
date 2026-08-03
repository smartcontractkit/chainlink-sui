# CCIP Cross-Chain Messaging: Go Load Test Guide

This document describes the necessary steps to send **arbitrary messages** (no tokens) between
Sui and EVM chains in Go, targeting **real deployed environments** (staging, testnet, mainnet).
It assumes you already have all required on-chain addresses.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        SUI → EVM                                │
│                                                                 │
│  Go script                                                      │
│    │                                                            │
│    ├─ 1. Build PTB (Programmable Transaction Block)             │
│    │      a. create_token_transfer_params (empty, no tokens)    │
│    │      b. ccip_send (onramp)                                 │
│    ├─ 2. Sign & execute PTB via Sui fullnode                    │
│    └─ 3. Extract CCIPMessageSent event → sequence number        │
│                                                                 │
│  DON (off-chain) handles commit + execution automatically       │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                        EVM → SUI                                │
│                                                                 │
│  Go script                                                      │
│    │                                                            │
│    ├─ 1. Call Router.getFee(destChainSelector, message)         │
│    ├─ 2. If LINK fee: approve(Router, feeAmount)                │
│    ├─ 3. Call Router.ccipSend(destChainSelector, message)       │
│    └─ 4. Extract CCIPMessageSent event from OnRamp logs         │
│                                                                 │
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

## Step-by-Step: SUI → EVM (Message Only)

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

Even for message-only transfers, you must call this with an empty token receiver:

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

// Encode create_token_transfer_params with empty receiver
encoded, _ := helperContract.EncodeCallArgsWithGenerics(
    "create_token_transfer_params",
    []string{},          // typeArgs
    []string{},          // typeParams
    []string{"vector<u8>"}, // paramTypes
    []any{make([]byte, 32)}, // zeroed 32-byte receiver
    nil,                 // returnTypes
)

tokenParamsResult, _ := helperContract.AppendPTB(ctx, callOpts, ptb, encoded)
```

#### 5b. `ccip_send` (onramp)

```go
// Bind the onramp contract
onRampContract, _ := suiBind.NewBoundContract(
    latestOnRampPkg, "ccip_onramp", "onramp", ptbClient,
)

// Get the function signature from normalized module
normalizedModule, _ := ptbClient.GetNormalizedModule(ctx, latestOnRampPkg, "onramp")
fnSig := normalizedModule.ExposedFunctions["ccip_send"].(map[string]any)
paramTypes, _ := suiofframp_helper.DecodeParameters(logger, fnSig, "parameters")

// Encode ccip_send
encoded, _ := onRampContract.EncodeCallArgsWithGenerics(
    "ccip_send",
    []string{linkTokenPkgID + "::link::LINK"}, // fee token type
    []string{},
    paramTypes,
    []any{
        suiBind.Object{Id: ccipObjectRef},       // ccip state
        suiBind.Object{Id: onRampState},          // onramp state
        suiBind.Object{Id: "0x6"},                // clock
        destChainSelector,                        // uint64
        receiver,                                 // []byte (32-byte EVM address)
        data,                                     // []byte (message payload)
        tokenParamsResult,                        // from step 5a
        suiBind.Object{Id: linkCoinMetadataID},   // fee token metadata
        suiBind.Object{Id: feeTokenCoinObject},   // fee token coin
        extraArgs,                                // []byte
    },
    nil,
)

onRampContract.AppendPTB(ctx, callOpts, ptb, encoded)
```

**Reusable from:** `chainlink/deployment/ccip/changeset/testhelpers/test_sui_helpers.go` —
`SendSuiCCIPRequest()` does exactly this (with additional token-pool logic).

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

- **LINK fees:** Split a LINK coin to cover `fee + 20% buffer`, pass as `feeToken` arg.
- **Native SUI fees:** Use `ptb.SetGasBudget(budget)` — the fee is deducted from the gas coin.

The fee can be queried via `fee_quoter.get_fee` (devInspect) before building the PTB.

---

## Step-by-Step: EVM → SUI (Message Only)

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
    TokenAmounts: []router.ClientEVMTokenAmount{}, // empty = no tokens
    FeeToken:     feeTokenAddress,  // LINK address or 0x0 for native
    ExtraArgs:    extraArgsBytes,   // SuiExtraArgsV1
}
```

### 5. Construct ExtraArgs (SuiExtraArgsV1)

For EVM→Sui, use the **SuiExtraArgsV1** format:

```go
import (
    "math/big"
    "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/message_hasher"
    "github.com/smartcontractkit/chainlink/core/capabilities/ccip/ccipevm"
)

extraArgs, _ := ccipevm.SerializeClientSUIExtraArgsV1(message_hasher.ClientSuiExtraArgsV1{
    GasLimit:                 big.NewInt(100_000),
    AllowOutOfOrderExecution: true,
    TokenReceiver:            [32]byte{}, // zero for message-only
    ReceiverObjectIds: [][32]byte{
        clockObjectID,        // 0x0000...0006
        receiverStateObjID,   // CCIPReceiverState from the receiver package
    },
})
```

**Reusable from:** `chainlink/core/capabilities/ccip/ccipevm/msghasher.go` —
`SerializeClientSUIExtraArgsV1()`.

### 6. Get the Fee

```go
fee, _ := routerContract.GetFee(&bind.CallOpts{}, destChainSelector, msg)
```

### 7. Approve LINK (if paying with LINK)

```go
if msg.FeeToken != (common.Address{}) {
    // Add 20% buffer
    feeWithBuffer := new(big.Int).Add(fee, new(big.Int).Div(fee, big.NewInt(5)))

    linkToken, _ := erc677.NewBurnMintERC677(msg.FeeToken, ethClient)
    tx, _ := linkToken.Approve(auth, routerAddress, feeWithBuffer)
    // Wait for confirmation...
}
```

### 8. Send the Message

```go
// For LINK fees:
auth.Value = nil
tx, _ := routerContract.CcipSend(auth, destChainSelector, msg)

// For native fees:
auth.Value = feeWithBuffer
tx, _ := routerContract.CcipSend(auth, destChainSelector, msg)
```

**Reusable from:** `chainlink/deployment/ccip/changeset/testhelpers/test_helpers_solana_v0_1_0.go` —
`CCIPSendRequest()` and `retryCcipSendUntilNativeFeeIsSufficient()`.

### 9. Extract the CCIPMessageSent Event

```go
import (
    onramp "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/onramp"
)

receipt, _ := ethClient.TransactionReceipt(ctx, tx.Hash())
onRampContract, _ := onramp.NewOnRamp(onRampAddress, ethClient)

// Filter CCIPMessageSent events in the receipt block
it, _ := onRampContract.FilterCCIPMessageSent(&bind.FilterOpts{
    StartBlock: receipt.BlockNumber,
    EndBlock:   &receipt.BlockNumber,
}, []uint64{destChainSelector}, []uint64{})

for it.Next() {
    seqNum := it.Event.SequenceNumber
    messageID := it.Event.Message.Header.MessageId
}
```

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
| ExtraArgs serialization | `chainlink/core/capabilities/ccip/ccipevm` | `SerializeClientSUIExtraArgsV1()` |
| GenericExtraArgsV2 | `chainlink/deployment/ccip/changeset/testhelpers` | `MakeBCSEVMExtraArgsV2()`, `GenericExtraArgsV2Tag` |
| Sui Go SDK | `github.com/block-vision/sui-go-sdk/sui` | `ISuiAPI` (JSON-RPC), `NewSuiClientWithCustomClient` |
| Sui Go SDK tx | `github.com/block-vision/sui-go-sdk/transaction` | `NewTransaction()` (PTB builder) |

---

## Minimal Dependencies for a Standalone Script

If you want a script that depends **only** on `chainlink-sui` and `chainlink-ccip`
(not the full `chainlink` monorepo), you need to:

1. **Vendor or copy** `MakeBCSEVMExtraArgsV2` (~10 lines, just BCS-encodes a u256 + bool with a 4-byte tag).
2. **Vendor or copy** `SerializeClientSUIExtraArgsV1` (uses `message_hasher.ClientSuiExtraArgsV1` from `chainlink-ccip` — already available).
3. **Use `chainlink-sui/bindings/bind`** for PTB construction and execution.
4. **Use `chainlink-sui/bindings/utils`** for the `SuiSigner` interface and `TestPrivateKeySigner`.
5. **Use `chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router`** for EVM Router interactions.

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
