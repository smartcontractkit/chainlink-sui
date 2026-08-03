package config

import (
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// RunConfig is loaded from a TOML run config file (Layer 4).
type RunConfig struct {
	Run      RunParams      `toml:"run"`
	Receiver ReceiverParams `toml:"receiver"`
	Gas      GasParams      `toml:"gas"`
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
