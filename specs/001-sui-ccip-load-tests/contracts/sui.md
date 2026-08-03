# Sui Package Contract

**Package**: `github.com/smartcontractkit/chainlink-sui/integration-tests/load/sui`

## Imports

```go
import (
    "github.com/block-vision/sui-go-sdk/sui"
    "github.com/smartcontractkit/chainlink-sui/bindings/bind"
    "github.com/smartcontractkit/chainlink-sui/bindings/utils"  // bindutils.SuiSigner
    "github.com/smartcontractkit/chainlink-sui/relayer/client"
    "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)
```

## Public Functions

```go
// NewSuiClient creates a Sui PTB client from an RPC URL.
// Uses the same client type as existing integration tests (client.SuiPTBClient).
func NewSuiClient(rpcURL string) (client.SuiPTBClient, error)

// NewSuiSigner creates a Sui signer from a bech32-encoded private key.
// The key must be in suiprivkey1... format (Ed25519).
// Uses relayer/signer.NewPrivateKeySigner (production, no build tag).
func NewSuiSigner(bech32PrivKey string) (bindutils.SuiSigner, string, error)

// SendMessage builds and executes a CCIP message from Sui to EVM.
// Returns the message ID, transaction digest, and sequence number.
func SendMessage(
    ctx context.Context,
    ptbClient client.SuiPTBClient,
    signer bindutils.SuiSigner,
    senderAddress string,
    ccipPkgID string,       // SuiCCIP package ID
    onRampPkgID string,     // SuiOnRamp package ID
    ccipObjectRef string,   // CCIPObjectRef shared object ID
    onRampState string,     // OnRampState shared object ID
    linkCoinMetadataID string, // LINK coin metadata ID
    destChainSelector uint64,
    receiver []byte,        // 32-byte EVM receiver address
    data []byte,            // message payload
    gasBudget uint64,       // PTB gas budget (covers fee)
) (messageID string, txDigest string, seqNum uint64, err error)

// MakeBCSEVMExtraArgsV2 constructs GenericExtraArgsV2 for Sui→EVM messages.
// Tag: 0x181dcf10
// BCS encoding: u256 gasLimit (32-byte LE) + bool allowOOO
func MakeBCSEVMExtraArgsV2(gasLimit *big.Int, allowOOO bool) []byte
```

## Key Implementation Details

### PTB Construction (Sui→EVM)

Uses `bind.BoundContract` and `bind.ExecutePTB` (same pattern as existing integration tests):

1. Create `bind.NewBoundContract(packageID, "ccip", "onramp_state_helper", ptbClient)` for the helper
2. Create `bind.NewBoundContract(onRampPkgID, "ccip_onramp", "onramp", ptbClient)` for the onramp
3. Build PTB with `transaction.NewTransaction()`
4. Call `create_token_transfer_params` with empty 32-byte receiver (even for message-only)
5. Call `ccip_send` with the token params result, CCIP state, OnRamp state, clock, etc.
6. Execute with `bind.ExecutePTB(ctx, &bind.CallOpts{Signer: signer, GasBudget: &gasBudget, WaitForExecution: true}, ptbClient, ptb)`

### Event Extraction

After PTB execution, scan `response.Events` for events where:
- `PackageId` matches the OnRamp package ID
- `Type` ends with `CCIPMessageSent`

Extract `sequence_number` and `message.header.message_id` from `ParsedJson`.

### Signing Flow

1. Decode bech32 private key → 32-byte seed
2. Create `ed25519.PrivateKey` from seed
3. Create signer with `relayer/signer.NewPrivateKeySigner(privateKey)`
4. Returns `bindutils.SuiSigner` interface (compatible with `bind.CallOpts`)

### Latest Package ID Resolution

Use `ptbClient.GetLatestPackageId(ctx, packageID, moduleName)` to resolve upgraded packages before building PTBs.
