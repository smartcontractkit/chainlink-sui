package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

// LoadRunConfig loads a TOML run config file from runs/<name>.toml.
// The filename (without .toml) is used as the run name for results.
// Accepts runName with or without the .toml extension.
func LoadRunConfig(runName string) (*RunConfig, error) {
	runName = strings.TrimSuffix(runName, ".toml")
	path := filepath.Join("runs", runName+".toml")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read run config %s: %w", path, err)
	}

	var cfg RunConfig
	if err := toml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse run config %s: %w", path, err)
	}

	if err := validateRunConfig(&cfg); err != nil {
		return nil, fmt.Errorf("invalid run config %s: %w", path, err)
	}

	return &cfg, nil
}

func validateRunConfig(cfg *RunConfig) error {
	if cfg.Run.Env == "" {
		return fmt.Errorf("run.env is required")
	}
	if cfg.Run.SourceChainSelector == "" {
		return fmt.Errorf("run.source_chain_selector is required")
	}
	if cfg.Run.DestChainSelector == "" {
		return fmt.Errorf("run.dest_chain_selector is required")
	}
	if cfg.Run.MessageCount < 1 {
		return fmt.Errorf("run.message_count must be >= 1, got %d", cfg.Run.MessageCount)
	}
	if cfg.Receiver.Address == "" {
		return fmt.Errorf("receiver.address is required")
	}

	// message_data is optional for token_only mode; required for message-only and token_and_message
	if cfg.Token == nil && cfg.Run.MessageData == "" {
		return fmt.Errorf("run.message_data is required for message-only mode")
	}

	if cfg.Token != nil {
		if cfg.Token.Amount == 0 {
			return fmt.Errorf("token.amount is required when [token] section is present and must be > 0")
		}
		if cfg.Token.Mode != "token_only" && cfg.Token.Mode != "token_and_message" {
			return fmt.Errorf("token.mode must be either \"token_only\" or \"token_and_message\", got %q", cfg.Token.Mode)
		}
		if cfg.Token.CoinMetadataID == "" && cfg.Token.TokenAddress == "" {
			return fmt.Errorf("token identifier is required: set coin_metadata_id for Sui source or token_address for EVM source")
		}
		if cfg.Token.Mode == "token_and_message" {
			if cfg.Run.MessageData == "" {
				return fmt.Errorf("run.message_data is required for token_and_message mode")
			}
			if cfg.Gas.EvmCallbackGasLimit == 0 {
				return fmt.Errorf("gas.evm_callback_gas_limit is required for token_and_message mode and must be > 0")
			}
		}
	}

	if cfg.SuiReceiver != nil {
		if cfg.SuiReceiver.PackageID == "" {
			return fmt.Errorf("sui_receiver.package_id is required when [sui_receiver] section is present")
		}
	}

	return nil
}

// parseChainSelector parses a chain selector string (decimal) to uint64.
func parseChainSelector(s string) (uint64, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid chain selector %q: %w", s, err)
	}
	return v, nil
}

// RunNameFromPath extracts the run name from a TOML file path.
func RunNameFromPath(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, ".toml")
}
