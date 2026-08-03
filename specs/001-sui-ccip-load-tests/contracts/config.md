# Config Package Contract

**Package**: `github.com/smartcontractkit/chainlink-sui/integration-tests/load/config`

## Public Types

```go
// LoadTestConfig is the assembled configuration for a load test run.
type LoadTestConfig struct {
    EnvName             string
    SourceChainSelector uint64
    DestChainSelector   uint64
    MessageCount        int
    MessageData         []byte
    SuiRPCURL           string
    SuiPrivateKey       string
    EVMPrivateKey       string
    SuiAddresses        AddressBook
    EVMAddresses        []FlatAddress
    EVMNetworks         []NetworkConfig
}

// AddressBook maps chain selector -> address -> TypeAndVersion.
type AddressBook map[string]map[string]TypeAndVersion

// TypeAndVersion describes a deployed contract.
type TypeAndVersion struct {
    Type    string              `json:"Type"`
    Version string              `json:"Version"`
    Labels  map[string]struct{} `json:"Labels,omitempty"`
}

// FlatAddress is an EVM address entry from the flat addresses.json format.
type FlatAddress struct {
    Address       string   `json:"address"`
    ChainSelector uint64   `json:"chainSelector"`
    Type          string   `json:"type"`
    Version       string   `json:"version"`
    Qualifier     string   `json:"qualifier,omitempty"`
    Labels        []string `json:"labels,omitempty"`
}

// NetworkConfig is an EVM network entry from the YAML config.
type NetworkConfig struct {
    Type          string         `yaml:"type"`
    ChainSelector uint64         `yaml:"chain_selector"`
    RPCs          []RPCConfig    `yaml:"rpcs"`
    Metadata      map[string]any `yaml:"metadata,omitempty"`
}

// RPCConfig describes an RPC endpoint.
type RPCConfig struct {
    RPCName string `yaml:"rpc_name"`
    HTTPURL string `yaml:"http_url"`
    WSURL   string `yaml:"ws_url"`
}
```

## Public Functions

```go
// LoadConfig loads the three-layer configuration for the given environment.
//  1. Loads .env.<envName> for secrets
//  2. Loads addresses.json for contract addresses
//  3. Loads networks.<envName>.yaml for EVM network config
// Returns an error if any required file is missing or malformed.
func LoadConfig(envName string, opts LoadOptions) (*LoadTestConfig, error)

// LoadOptions provides optional overrides for config loading.
type LoadOptions struct {
    SourceChainSelector *uint64 // override auto-detected source chain
    DestChainSelector   *uint64 // override auto-detected dest chain
    MessageCount        int     // override default message count
    MessageData         string  // override default message payload
    AddressesPath       string  // override default addresses.json path
    NetworksPath        string  // override default networks YAML path
    ResultsDir          string  // override default results directory
}

// ParseAddressesJSON parses an addresses.json file, detecting format (nested or flat).
func ParseAddressesJSON(data []byte) (AddressBook, []FlatAddress, error)

// ParseNetworksYAML parses a networks YAML file.
func ParseNetworksYAML(data []byte) ([]NetworkConfig, error)

// ResolveSuiAddresses extracts Sui contract addresses for a given chain selector
// from the address book.
func ResolveSuiAddresses(book AddressBook, chainSelector string) (SuiAddresses, error)

// SuiAddresses holds the resolved Sui contract addresses for one chain.
type SuiAddresses struct {
    CCIPPackageID       string
    RouterPackageID     string
    OnRampPackageID     string
    OffRampPackageID    string
    LinkCoinMetadataID  string
    CCIPObjectRef       string
    OnRampStateObjectID string
}
```

## File Naming Convention

| File | Pattern | Required |
|------|---------|----------|
| Secrets | `.env.<envName>` | Yes |
| Addresses | `addresses.json` or `addresses-<envName>.json` | Yes |
| Networks | `networks-<envName>.yaml` | Yes |
| Results | `results/<envName>-<timestamp>.json` | Created at runtime |
