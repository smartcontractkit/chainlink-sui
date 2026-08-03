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
func LoadRunConfig(runName string) (*RunConfig, error) {
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
	if cfg.Run.MessageData == "" {
		return fmt.Errorf("run.message_data is required")
	}
	if cfg.Receiver.Address == "" {
		return fmt.Errorf("receiver.address is required")
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
