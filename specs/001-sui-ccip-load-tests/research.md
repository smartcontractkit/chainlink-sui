# Research: Sui CCIP Load Tests

## Unknown 1: EVM Router/OnRamp Binding Availability

**Decision**: Use `chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router` for the Router contract. For OnRamp event filtering, use the Router's `CCIPMessageSent` event or raw log parsing — the `v1_6_3/onramp` package does not exist in the gobindings tree.

**Rationale**: The `v1_2_0/router` package provides `RouterCCIPMessageSent` event filtering directly on the Router contract, which is sufficient for extracting message IDs. No need for a separate OnRamp binding.

**Alternatives considered**:
- `v1_6_3/onramp` — does not exist in the dependency tree
- Raw log parsing with go-ethereum — possible but unnecessary when Router events are available

**Action**: Add `github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0` as a dependency.

---

## Unknown 2: SuiExtraArgsV1 Encoding (EVM→Sui)

**Decision**: Replicate the encoding using `message_hasher.MessageHasherABI` and go-ethereum's ABI packer. No Chainlink core dependency needed.

**Rationale**: The `message_hasher` package (from `chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/message_hasher`) exports `MessageHasherABI` which contains the `encodeSUIExtraArgsV1` method ABI. We can use `go-ethereum`'s `accounts/abi` package to pack the inputs with the 4-byte tag prefix `0x21ea4ca9`.

**Implementation**:
```go
import (
    "math/big"
    "github.com/ethereum/go-ethereum/accounts/abi"
    "github.com/ethereum/go-ethereum/common/hexutil"
    "github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/message_hasher"
)

var suiExtraArgsABI = func() abi.ABI {
    a, _ := abi.JSON(strings.NewReader(message_hasher.MessageHasherABI))
    return a
}()

func SerializeClientSUIExtraArgsV1(data message_hasher.ClientSuiExtraArgsV1) ([]byte, error) {
    tag := hexutil.MustDecode("0x21ea4ca9")
    v, err := suiExtraArgsABI.Methods["encodeSUIExtraArgsV1"].Inputs.Pack(data)
    if err != nil {
        return nil, err
    }
    return append(tag, v...), nil
}
```

**Alternatives considered**:
- Importing `chainlink/v2/core/capabilities/ccip/ccipevm` — blocked by no-Chainlink-core constraint
- Hardcoding the ABI encoding — fragile and error-prone

---

## Unknown 3: GenericExtraArgsV2 Encoding (Sui→EVM)

**Decision**: Implement as documented in the existing README — simple BCS encoding with 4-byte tag `0x181dcf10`.

**Rationale**: The encoding is trivial: 4-byte tag + 32-byte LE gasLimit + 1-byte bool. No external dependency needed beyond `go-ethereum`'s `hexutil` (already in go.mod).

**Implementation**: Already documented in `integration-tests/load/README.md`.

---

## Unknown 4: Sui Client Approach

**Decision**: Use the raw `sui.ISuiAPI` (HTTP-only) from `block-vision/sui-go-sdk` for simplicity. Construct PTBs using `transaction.NewTransaction()` and `ptb.MoveCall()`. Sign with `block-vision/sui-go-sdk/signer`. Submit via `SuiExecuteTransactionBlock`.

**Rationale**: The `PTBClient` from `chainlink-sui/relayer/client` requires gRPC configuration which adds complexity and is unnecessary for simple message sending. The raw SDK provides all needed functionality via HTTP.

**Key types**:
- `sui.ISuiAPI` — JSON-RPC client interface
- `transaction.Transaction` — PTB builder
- `signer.Signer` — Transaction signer (from `block-vision/sui-go-sdk/signer`, NOT `chainlink-sui/relayer/signer`)
- `models.SuiExecuteTransactionBlockRequest` — Transaction submission

**For argument encoding**: Use `transaction.Pure()` for BCS-encoded values and `transaction.Object()` for object references. For complex type encoding (generics), construct `transaction.TypeTag` values manually.

**Alternatives considered**:
- `PTBClient` with gRPC — richer API but requires gRPC endpoint and more complex setup
- `BindingsClient` via `sui.Client` wrapper — would need adapter code, not worth it

---

## Unknown 5: Addresses.json Format

**Decision**: Support both formats — detect at parse time.

**Rationale**: The deployment pipeline produces the nested format (`map[chainSelector]map[address]TypeAndVersion`) for Sui addresses. The existing `addresses-testnet.json` uses a flat array format for EVM addresses. Both need to be supported.

**Detection**: Check if the JSON root is an object (nested format) or array (flat format).

**Types**:
```go
// Nested format (Sui addresses from deployment)
type AddressBook map[string]map[string]TypeAndVersion

type TypeAndVersion struct {
    Type    string              `json:"Type"`
    Version string              `json:"Version"`
    Labels  map[string]struct{} `json:"Labels,omitempty"`
}

// Flat format (EVM addresses)
type FlatAddress struct {
    Address       string   `json:"address"`
    ChainSelector uint64   `json:"chainSelector"`
    Type          string   `json:"type"`
    Version       string   `json:"version"`
    Qualifier     string   `json:"qualifier,omitempty"`
    Labels        []string `json:"labels,omitempty"`
}
```

---

## Unknown 6: Sui Signer Availability

**Decision**: Use `github.com/smartcontractkit/chainlink-sui/relayer/signer.NewPrivateKeySigner` — it's importable (via `replace` directive) and production-ready.

**Rationale**: The `integration-tests/go.mod` has `replace github.com/smartcontractkit/chainlink-sui => ../`, making the `relayer/signer` package available. The `NewPrivateKeySigner` function takes an `ed25519.PrivateKey` and implements the `SuiSigner` interface with proper intent-scope prefix and Blake2b hashing.

**Bech32 decoding**: Use `github.com/btcsuite/btcutil/bech32` (already indirect dep) to decode `suiprivkey1...` format. The decoded payload is `flag(1B) || seed(32B)` where flag `0x00` indicates Ed25519.

**Alternatives considered**:
- `bindings/utils/test_signer.go` — gated behind `//go:build integration`, not usable
- Manual signing implementation — unnecessary when `relayer/signer` is available
