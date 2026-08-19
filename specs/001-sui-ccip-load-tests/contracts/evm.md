# EVM Package Contract

**Package**: `github.com/smartcontractkit/chainlink-sui/integration-tests/load/evm`

## Imports

```go
import (
    "github.com/ethereum/go-ethereum/ethclient"
    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    router "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router"
    "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/message_hasher"
)
```

## Public Functions

```go
// NewEVMClient creates an Ethereum JSON-RPC client from the given URL.
func NewEVMClient(rpcURL string) (*ethclient.Client, error)

// NewEVMSigner creates an EVM transactor from a hex-encoded private key.
func NewEVMSigner(privateKeyHex string, chainID *big.Int) (*bind.TransactOpts, error)

// SendMessage sends a CCIP message from EVM to Sui.
// Returns the message ID and transaction hash.
func SendMessage(
    ctx context.Context,
    client *ethclient.Client,
    auth *bind.TransactOpts,
    routerAddress common.Address,
    destChainSelector uint64,
    receiver [32]byte,       // 32-byte Sui receiver (package ID)
    data []byte,             // message payload
    extraArgs []byte,        // SuiExtraArgsV1 encoded
) (messageID string, txHash string, err error)

// GetFee queries the CCIP fee for a message.
func GetFee(
    ctx context.Context,
    client *ethclient.Client,
    routerAddress common.Address,
    destChainSelector uint64,
    receiver [32]byte,
    data []byte,
    extraArgs []byte,
) (*big.Int, error)

// SerializeClientSUIExtraArgsV1 encodes SuiExtraArgsV1 for EVM→Sui messages.
// Tag: 0x21ea4ca9
// Uses message_hasher.MessageHasherABI for ABI encoding.
func SerializeClientSUIExtraArgsV1(data message_hasher.ClientSuiExtraArgsV1) ([]byte, error)

// ExtractMessageIDFromReceipt extracts the CCIPMessageSent event from a transaction receipt.
func ExtractMessageIDFromReceipt(
    ctx context.Context,
    client *ethclient.Client,
    routerAddress common.Address,
    txHash common.Hash,
    destChainSelector uint64,
) (messageID string, seqNum uint64, err error)
```

## Key Implementation Details

### Message Construction (EVM→Sui)

```go
msg := router.ClientEVM2AnyMessage{
    Receiver:     receiver,       // 32-byte Sui package ID
    Data:         data,           // arbitrary bytes
    TokenAmounts: []router.ClientEVMTokenAmount{}, // empty for message-only
    FeeToken:     common.Address{}, // zero address = native ETH fee
    ExtraArgs:    extraArgs,      // SuiExtraArgsV1
}
```

### Fee Handling (Native ETH)

1. Call `Router.GetFee(destChainSelector, msg)` to get the fee
2. Add 20% buffer: `feeWithBuffer = fee + fee/5`
3. Set `auth.Value = feeWithBuffer`
4. Call `Router.CcipSend(auth, destChainSelector, msg)`

### Event Extraction

After getting the transaction receipt:
1. Instantiate `Router` contract binding at the router address
2. Call `RouterFilterCCIPMessageSent` with the receipt block range
3. Iterate events to find the one matching `destChainSelector`
4. Extract `Message.Header.MessageId` and `SequenceNumber`
