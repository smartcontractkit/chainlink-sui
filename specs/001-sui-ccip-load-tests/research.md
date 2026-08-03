# Research: Sui CCIP Load Tests

## Unknown 1: EVM Router Binding Availability

**Decision**: Use `chainlink-ccip/chains/evm/gobindings/generated/v1_2_0/router` for the Router contract. The Router binding provides `CCIPSend`, `GetFee`, and `FilterCCIPMessageSent` — everything needed for EVM→Sui message sending.

**Rationale**: The `v1_2_0/router` package is already an indirect dependency. It provides `RouterCCIPMessageSent` event filtering directly on the Router contract, which is sufficient for extracting message IDs. No separate OnRamp binding needed.

---

## Unknown 2: SuiExtraArgsV1 Encoding (EVM→Sui)

**Decision**: Replicate the encoding using `message_hasher.MessageHasherABI` and go-ethereum's ABI packer. The `message_hasher` package is already in the dependency tree.

**Rationale**: The `message_hasher` package (from `chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/message_hasher`) exports `MessageHasherABI` which contains the `encodeSUIExtraArgsV1` method ABI. Use `go-ethereum`'s `accounts/abi` to pack inputs with the 4-byte tag prefix `0x21ea4ca9`.

**Implementation**:
```go
import (
    "strings"
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
    if err != nil { return nil, err }
    return append(tag, v...), nil
}
```

---

## Unknown 3: GenericExtraArgsV2 Encoding (Sui→EVM)

**Decision**: Implement as documented in the existing README — simple BCS encoding with 4-byte tag `0x181dcf10`.

**Rationale**: The encoding is trivial: 4-byte tag + 32-byte LE gasLimit + 1-byte bool. No external dependency needed beyond `go-ethereum`'s `hexutil` (already in go.mod).

---

## Unknown 4: Sui Client & Signer Approach

**Decision**: Use `cldf` (chainlink-deployments-framework) helpers for Sui client and signer setup. The `cldf_chain/sui` package provides `NewSuiClient` and signer utilities that existing integration tests use.

**Rationale**: The existing integration tests in `integration-tests/deploy/` and `integration-tests/offramp/` use `testenv.SetupEnvironment` which returns a `client.SuiPTBClient` and `bindutils.SuiSigner`. For remote environments, we can use `cldf_sui.NewSuiClient(rpcURL)` to create a client and `relayer/signer.NewPrivateKeySigner` for signing. The `bindings/bind` package's `BoundContract`, `ExecutePTB`, and `CallOpts` work with these types.

**Key types**:
- `client.SuiPTBClient` — full PTB client (from `chainlink-sui/relayer/client`)
- `bindutils.SuiSigner` — signer interface (from `chainlink-sui/bindings/utils`)
- `relayer/signer.NewPrivateKeySigner(ed25519.PrivateKey)` — production signer
- `bind.CallOpts{Signer, GasBudget, WaitForExecution}` — PTB execution options
- `bind.NewBoundContract(packageID, packageName, moduleName, client)` — contract binding
- `bind.ExecutePTB(ctx, opts, client, ptb)` — PTB execution

**Bech32 decoding**: Use `github.com/btcsuite/btcutil/bech32` (already indirect dep) to decode `suiprivkey1...` format. The decoded payload is `flag(1B) || seed(32B)` where flag `0x00` indicates Ed25519.

---

## Unknown 5: Addresses.json Format

**Decision**: Use `cldf.AddressBook` (from `chainlink-deployments-framework/deployment`) directly. The `addresses-testnet.json` file is a flat array that can be loaded into a `cldf.AddressBook` via `cldf.NewMemoryAddressBookFromFile()`.

**Rationale**: The `cldf.AddressBook` interface provides `AddressesForChain(chainSelector)` which returns `map[string]cldf.TypeAndVersion`. This is the standard format used throughout the deployment framework. The flat array format in `addresses-testnet.json` is just a serialization format — `cldf` can load it.

**Key types**:
- `cldf.AddressBook` — interface with `Save(chainSelector, address, TypeAndVersion)` and `AddressesForChain(chainSelector)`
- `cldf.TypeAndVersion` — has `.Type`, `.Version`, `.Labels` fields
- `cldf.NewMemoryAddressBookFromFile(path)` — loads from JSON file

**Addresses needed for sending** (message-only, no tokens):
- For Sui→EVM: `SuiRouter` (router package ID), `SuiOnRamp` (onramp package ID), `SuiCCIP` (CCIP package ID for object refs), `SuiLinkTokenCoinMetadataID` (fee token metadata)
- For EVM→Sui: `Router` (EVM router address), `LinkToken` (fee token address)

---

## Unknown 6: Network Config Format

**Decision**: The `networks-testnet.yaml` already defines a unified format for all chains (EVM and Sui). Each entry has `type`, `chain_selector`, and `rpcs` (array of `rpc_name`, `http_url`, `ws_url`). No distinction between EVM and Sui — both are just chains with RPC endpoints.

**Rationale**: The existing YAML already includes Sui entries (e.g., `chain_selector: 9762610643973837292`). The config loader just needs to look up the chain by selector and use the first HTTP RPC URL.

**Key types**:
```go
type NetworkConfig struct {
    Type          string      `yaml:"type"`
    ChainSelector uint64      `yaml:"chain_selector"`
    RPCs          []RPCConfig `yaml:"rpcs"`
}

type RPCConfig struct {
    RPCName string `yaml:"rpc_name"`
    HTTPURL string `yaml:"http_url"`
    WSURL   string `yaml:"ws_url"`
}
```

---

## Unknown 7: Run Config (4th Layer)

**Decision**: Add a TOML run config file as the 4th config layer. Each run gets its own file in `runs/`. The filename becomes the results filename.

**Rationale**: Separates run-specific parameters (source/dest chain selectors, message count, message data, receiver address) from environment config. Makes it easy to define multiple runs against the same environment.

**Format** (`runs/my-first-sui-to-evm-run.toml`):
```toml
[run]
env = "testnet"
source_chain_selector = 9762610643973837292
dest_chain_selector = 16015286601757825753
message_count = 10
message_data = "hello from sui load test"

[receiver]
# For Sui→EVM: EVM address (20 bytes, hex)
# For EVM→Sui: Sui package ID (32 bytes, hex)
address = "0x..."

[gas]
# Sui PTB gas budget (for Sui→EVM)
sui_gas_budget = 10000000000
# EVM gas limit (for EVM→Sui)
evm_gas_limit = 200000
```

**Results file**: `results/my-first-sui-to-evm-run-testnet-20260803T120000.txt`

---

## Unknown 8: Dependencies

**Decision**: No new dependencies needed. All required packages are already in `go.mod` (direct or indirect).

**Already available**:
- `gopkg.in/yaml.v3` — indirect dep (used by cldf)
- `github.com/pelletier/go-toml/v2` — indirect dep (used by cldf)
- `github.com/joho/godotenv` — already in go.mod
- `github.com/btcsuite/btcutil/bech32` — indirect dep
- `github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_2_0` — indirect dep
- `github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/message_hasher` — indirect dep
- `github.com/smartcontractkit/chainlink-deployments-framework` — already direct dep
- `github.com/smartcontractkit/chainlink-sui/bindings/bind` — same module
- `github.com/smartcontractkit/chainlink-sui/bindings/generated` — same module
- `github.com/smartcontractkit/chainlink-sui/relayer/client` — same module
- `github.com/smartcontractkit/chainlink-sui/relayer/signer` — same module
