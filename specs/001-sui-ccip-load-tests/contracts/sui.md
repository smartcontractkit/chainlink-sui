# Sui Package Contract

**Package**: `github.com/smartcontractkit/chainlink-sui/integration-tests/load/sui`

## Public Functions

```go
// NewSuiClient creates a Sui JSON-RPC client from the given URL.
func NewSuiClient(rpcURL string) (sui.ISuiAPI, error)

// NewSuiSigner creates a Sui signer from a bech32-encoded private key.
// The key must be in suiprivkey1... format (Ed25519).
func NewSuiSigner(bech32PrivKey string) (*signer.Signer, string, error)

// SendMessage builds and executes a CCIP message from Sui to EVM.
// Returns the message ID, transaction digest, and sequence number.
func SendMessage(
    ctx context.Context,
    client sui.ISuiAPI,
    signer *signer.Signer,
    senderAddress string,
    suiAddrs config.SuiAddresses,
    destChainSelector uint64,
    receiver []byte,       // 32-byte EVM receiver address
    data []byte,           // message payload
    gasBudget uint64,      // PTB gas budget (covers fee)
) (messageID string, txDigest string, seqNum uint64, err error)

// MakeBCSEVMExtraArgsV2 constructs GenericExtraArgsV2 for Sui→EVM messages.
// Tag: 0x181dcf10
// BCS encoding: u256 gasLimit (32-byte LE) + bool allowOOO
func MakeBCSEVMExtraArgsV2(gasLimit *big.Int, allowOOO bool) []byte
```

## Key Implementation Details

### PTB Construction (Sui→EVM)

The PTB requires two Move calls:

1. **`create_token_transfer_params`** (on `ccip::onramp_state_helper`):
   - Even for message-only transfers, this must be called with an empty 32-byte receiver
   - Returns token params result used as input to `ccip_send`

2. **`ccip_send`** (on `ccip_onramp::onramp`):
   - Takes: CCIP state object, OnRamp state object, clock, dest chain selector, receiver, data, token params, fee token metadata, fee token coin, extra args
   - Fee token coin: use a SUI coin from the sender's wallet (gas coin)
   - Fee token metadata: use LINK coin metadata ID (required by contract even for native fee)

### Event Extraction

After PTB execution, scan `response.Events` for events where:
- `PackageId` matches the OnRamp package ID
- `Type` ends with `CCIPMessageSent`

Extract `sequence_number` and `message.header.message_id` from `ParsedJson`.

### Signing Flow

1. Decode bech32 private key → 32-byte seed
2. Create `ed25519.PrivateKey` from seed
3. Create `signer.Signer` from `block-vision/sui-go-sdk/signer`
4. Use `signer.Sign` for transaction signing
