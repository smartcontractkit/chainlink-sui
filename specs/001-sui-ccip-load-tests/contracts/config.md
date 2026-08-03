# Config Package Contract

**Package**: `github.com/smartcontractkit/chainlink-sui/integration-tests/load/config`

## Imports

```go
import (
    cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)
```

## Public Types

```go
// RunConfig is loaded from a TOML run config file (Layer 4).
type RunConfig struct {
    Run     RunParams     `toml:"run"`
    Receiver ReceiverParams `toml:"receiver"`
    Gas     GasParams     `toml:"gas"`
}

type RunParams struct {
    Env                string `toml:"env"`
    SourceChainSelector uint64 `toml:"source_chain_selector"`
    DestChainSelector   uint64 `toml:"dest_chain_selector"`
    MessageCount        int    `toml:"message_count"`
    MessageData         string `toml:"message_data"`
}

type ReceiverParams struct {
    Address string `toml:"address"`
}

type GasParams struct {
    SuiGasBudget uint64 `toml:"sui_gas_budget"`
    EvmGasLimit  uint64 `toml:"evm_gas_limit"`
}

// LoadTestConfig is the fully assembled configuration for a load test run.
type LoadTestConfig struct {
    RunName             string   // from TOML filename (without .toml)
    EnvName             string   // from RunConfig.Run.Env
    SourceChainSelector uint64
    DestChainSelector   uint64
    MessageCount        int
    MessageData         []byte
    ReceiverAddress     string   // hex-encoded receiver address
    SuiGasBudget        uint64
    EvmGasLimit         uint64

    // Secrets (from .env)
    SuiPrivateKey string
    EVMPrivateKey string

    // Address book (from addresses.json, uses cldf.AddressBook)
    AddressBook cldf.AddressBook

    // Network configs (from YAML)
    Networks []NetworkConfig
}

// NetworkConfig is a unified chain network entry from the YAML config.
// Same format for EVM and Sui chains.
type NetworkConfig struct {
    Type          string      `yaml:"type"`
    ChainSelector uint64      `yaml:"chain_selector"`
    RPCs          []RPCConfig `yaml:"rpcs"`
}

// RPCConfig describes an RPC endpoint.
type RPCConfig struct {
    RPCName string `yaml:"rpc_name"`
    HTTPURL string `yaml:"http_url"`
    WSURL   string `yaml:"ws_url"`
}

// SentMessage represents a single sent CCIP message.
type SentMessage struct {
    MessageID           string `json:"message_id"`
    TransactionHash     string `json:"transaction_hash"`
    SourceChainSelector uint64 `json:"source_chain_selector"`
    DestChainSelector   uint64 `json:"dest_chain_selector"`
    Timestamp           string `json:"timestamp"`
    Success             bool   `json:"success"`
    Error               string `json:"error,omitempty"`
    SequenceNumber      string `json:"sequence_number,omitempty"`
}

// RunResults is the output of a load test run.
type RunResults struct {
    RunName             string         `json:"run_name"`
    EnvName             string         `json:"env_name"`
    SourceChainSelector uint64         `json:"source_chain_selector"`
    DestChainSelector   uint64         `json:"dest_chain_selector"`
    TotalMessages       int            `json:"total_messages"`
    SuccessfulMessages  int            `json:"successful_messages"`
    FailedMessages      int            `json:"failed_messages"`
    RunStarted          string         `json:"run_started"`
    RunEnded            string         `json:"run_ended"`
    Messages            []SentMessage  `json:"messages"`
}
```

## Public Functions

```go
// LoadRunConfig loads a TOML run config file.
func LoadRunConfig(path string) (*RunConfig, error)

// LoadEnvConfig loads secrets from a .env file.
func LoadEnvConfig(envName string) (suiPrivKey string, evmPrivKey string, err error)

// LoadAddressBook loads contract addresses from a JSON file using cldf.
func LoadAddressBook(path string) (cldf.AddressBook, error)

// LoadNetworks loads chain network configs from a YAML file.
func LoadNetworks(path string) ([]NetworkConfig, error)

// LoadFullConfig loads all four config layers for a run.
//  1. Loads runs/<runName>.toml for run parameters
//  2. Loads .env.<env> for secrets
//  3. Loads addresses-<env>.json for contract addresses
//  4. Loads networks-<env>.yaml for chain RPC endpoints
func LoadFullConfig(runName string) (*LoadTestConfig, error)

// ResolveAddressesForChain extracts addresses for a given chain selector
// from the cldf address book, keyed by contract type.
func ResolveAddressesForChain(book cldf.AddressBook, chainSelector uint64) (map[string]string, error)

// SaveResults writes run results to a file.
// File: results/<runName>-<env>-<timestamp>.txt
func SaveResults(results *RunResults) error
```

## File Naming Convention

| Layer | File Pattern | Required |
|-------|-------------|----------|
| 4 — Run config | `runs/<name>.toml` | Yes |
| 1 — Secrets | `.env.<env>` | Yes |
| 2 — Addresses | `addresses-<env>.json` | Yes |
| 3 — Networks | `networks-<env>.yaml` | Yes |
| Results | `results/<name>-<env>-<timestamp>.txt` | Created at runtime |
