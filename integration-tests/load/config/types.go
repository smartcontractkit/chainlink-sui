package config

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"gopkg.in/yaml.v3"
)

// RunConfig is loaded from a TOML run config file (Layer 4).
type RunConfig struct {
	Run         RunParams          `toml:"run"`
	Receiver    ReceiverParams     `toml:"receiver"`
	Gas         GasParams          `toml:"gas"`
	Load        *LoadParams        `toml:"load,omitempty"`
	Token       *TokenParams       `toml:"token,omitempty"`
	SuiReceiver *SuiReceiverParams `toml:"sui_receiver,omitempty"`
}

// RunParams holds the core run parameters.
// Chain selectors are strings in TOML to avoid int64 overflow (selectors can exceed 2^63-1).
type RunParams struct {
	Env                 string `toml:"env"`
	SourceChainSelector string `toml:"source_chain_selector"`
	DestChainSelector   string `toml:"dest_chain_selector"`
	MessageCount        int    `toml:"message_count"`
	MessageData         string `toml:"message_data"`
}

// ReceiverParams holds the receiver address.
type ReceiverParams struct {
	Address string `toml:"address"`
}

// GasParams holds gas settings for both chains.
type GasParams struct {
	SuiGasBudget        uint64 `toml:"sui_gas_budget"`
	EvmGasLimit         uint64 `toml:"evm_gas_limit"`
	EvmCallbackGasLimit uint64 `toml:"evm_callback_gas_limit"`
}

// LoadParams holds WASP load test settings.
// RPS is the target total requests per second across all wallets.
// Wallets is the number of parallel wallets (and generators) to use.
type LoadParams struct {
	RPS     int `toml:"rps"`
	Wallets int `toml:"wallets"`
}

// TokenParams holds token transfer parameters in the run config TOML.
// For Sui source chains, use coin_metadata_id. For EVM source chains, use token_address.
type TokenParams struct {
	CoinMetadataID string `toml:"coin_metadata_id,omitempty"`
	TokenAddress   string `toml:"token_address,omitempty"`
	Amount         uint64 `toml:"amount"`
	Mode           string `toml:"mode"`
}

// SuiReceiverParams holds the Sui receiver package ID for EVM→Sui programmable transfers.
type SuiReceiverParams struct {
	PackageID string `toml:"package_id"`
}

// LoadTestConfig is the fully assembled configuration for a load test run.
type LoadTestConfig struct {
	RunName             string // from TOML filename (without .toml)
	EnvName             string // from RunConfig.Run.Env
	SourceChainSelector uint64
	DestChainSelector   uint64
	MessageCount        int
	MessageData         []byte
	ReceiverAddress     string // hex-encoded receiver address
	SuiGasBudget        uint64
	EvmGasLimit         uint64
	EvmCallbackGasLimit uint64

	// Load test settings (from [load] section, defaults to rps=1, wallets=1)
	LoadRPS     int
	LoadWallets int

	// Secrets (from .env)
	SuiPrivateKey string
	EVMPrivateKey string
	WalletSeed    string // optional hex seed for deterministic load-test wallets
	SuiGrpcToken  string // optional gRPC auth token for Sui providers

	// Address book (from addresses.json, uses cldf.AddressBook)
	AddressBook cldf.AddressBook

	// Network configs (from YAML)
	Networks []NetworkConfig

	// Token transfer config (optional, nil for message-only runs)
	TokenConfig *TokenTransferConfig

	// Sui receiver config (optional, only for EVM→Sui programmable transfers)
	SuiReceiverConfig *SuiReceiverConfig
}

// TokenTransferConfig is the fully assembled token transfer configuration.
type TokenTransferConfig struct {
	TokenIdentifier string // coin metadata ID (Sui source) or ERC-20 address (EVM source)
	Amount          uint64 // token amount per message in base units
	Mode            string // "token_only" or "token_and_message"
}

// SuiReceiverConfig is the fully assembled Sui receiver configuration.
type SuiReceiverConfig struct {
	PackageID string // Sui receiver package ID
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
	RPCName    string `yaml:"rpc_name"`
	HTTPURL    string `yaml:"http_url"`
	WSURL      string `yaml:"ws_url"`
	GrpcTarget string `yaml:"grpc_target,omitempty"`
	GrpcToken  string `yaml:"grpc_token,omitempty"`
}

// UnmarshalYAML allows RPCs to be either a list of RPC configs or a single RPC config.
// The standard networks YAML uses both formats — some entries have `rpcs:` as a list
// (with `-` markers), others as a single map.
func (n *NetworkConfig) UnmarshalYAML(value *yaml.Node) error {
	type rawNetworkConfig struct {
		Type          string    `yaml:"type"`
		ChainSelector uint64    `yaml:"chain_selector"`
		RPCs          yaml.Node `yaml:"rpcs"`
	}
	var raw rawNetworkConfig
	if err := value.Decode(&raw); err != nil {
		return err
	}
	n.Type = raw.Type
	n.ChainSelector = raw.ChainSelector

	switch raw.RPCs.Kind {
	case yaml.SequenceNode:
		// List of RPC configs
		if err := raw.RPCs.Decode(&n.RPCs); err != nil {
			return err
		}
	case yaml.MappingNode:
		// Single RPC config (not a list)
		var single RPCConfig
		if err := raw.RPCs.Decode(&single); err != nil {
			return err
		}
		n.RPCs = []RPCConfig{single}
	default:
		n.RPCs = nil
	}
	return nil
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
	TokenAmount         string `json:"token_amount,omitempty"`
	TokenIdentifier     string `json:"token_identifier,omitempty"`
}

// RunResults is the output of a load test run.
type RunResults struct {
	RunName             string        `json:"run_name"`
	EnvName             string        `json:"env_name"`
	SourceChainSelector uint64        `json:"source_chain_selector"`
	DestChainSelector   uint64        `json:"dest_chain_selector"`
	TotalMessages       int           `json:"total_messages"`
	SuccessfulMessages  int           `json:"successful_messages"`
	FailedMessages      int           `json:"failed_messages"`
	RunStarted          string        `json:"run_started"`
	RunEnded            string        `json:"run_ended"`
	Messages            []SentMessage `json:"messages"`
}
