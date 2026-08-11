package config

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
)

// LoadEnvConfig loads secrets from a .env.<env> file.
// Returns the Sui and EVM private keys and an optional wallet seed.
func LoadEnvConfig(envName string) (suiPrivKey string, evmPrivKey string, walletSeed string, err error) {
	path := fmt.Sprintf(".env.%s", envName)
	if err := godotenv.Load(path); err != nil {
		return "", "", "", fmt.Errorf("failed to load env file %s: %w", path, err)
	}

	suiPrivKey = os.Getenv("SUI_PRIVATE_KEY")
	if suiPrivKey == "" {
		return "", "", "", fmt.Errorf("SUI_PRIVATE_KEY not set in %s", path)
	}

	evmPrivKey = os.Getenv("EVM_PRIVATE_KEY")
	if evmPrivKey == "" {
		return "", "", "", fmt.Errorf("EVM_PRIVATE_KEY not set in %s", path)
	}

	walletSeed = os.Getenv("WALLET_SEED")

	return suiPrivKey, evmPrivKey, walletSeed, nil
}

// LoadAddressBook loads contract addresses from an addresses-<env>.json file
// using the cldf AddressBook format.
func LoadAddressBook(envName string) (cldf.AddressBook, error) {
	path := fmt.Sprintf("addresses-%s.json", envName)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read address book %s: %w", path, err)
	}

	// The file is a flat array of {address, chainSelector, type, version, ...} entries.
	// We parse it and build a cldf.AddressBook from the entries.
	var entries []struct {
		Address       string   `json:"address"`
		ChainSelector uint64   `json:"chainSelector"`
		Type          string   `json:"type"`
		Version       string   `json:"version"`
		Labels        []string `json:"labels,omitempty"`
	}
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, fmt.Errorf("failed to parse address book %s: %w", path, err)
	}

	// Build the address book map: chainSelector -> address -> TypeAndVersion
	addrsByChain := make(map[uint64]map[string]cldf.TypeAndVersion)
	for _, entry := range entries {
		if _, ok := addrsByChain[entry.ChainSelector]; !ok {
			addrsByChain[entry.ChainSelector] = make(map[string]cldf.TypeAndVersion)
		}
		ver, err := semver.NewVersion(entry.Version)
		if err != nil {
			return nil, fmt.Errorf("invalid version %q for address %s: %w", entry.Version, entry.Address, err)
		}
		tv := cldf.NewTypeAndVersion(cldf.ContractType(entry.Type), *ver)
		for _, label := range entry.Labels {
			tv.Labels.Add(label)
		}
		addrsByChain[entry.ChainSelector][entry.Address] = tv
	}

	return cldf.NewMemoryAddressBookFromMap(addrsByChain), nil
}

// LoadNetworks loads chain network configs from a networks-<env>.yaml file.
func LoadNetworks(envName string) ([]NetworkConfig, error) {
	path := fmt.Sprintf("networks-%s.yaml", envName)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read networks config %s: %w", path, err)
	}

	var wrapper struct {
		Networks []NetworkConfig `yaml:"networks"`
	}
	if err := yaml.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse networks config %s: %w", path, err)
	}

	if len(wrapper.Networks) == 0 {
		return nil, fmt.Errorf("no networks found in %s", path)
	}

	return wrapper.Networks, nil
}

// FindNetworkBySelector finds a network config by chain selector.
func FindNetworkBySelector(networks []NetworkConfig, selector uint64) (*NetworkConfig, error) {
	for i := range networks {
		if networks[i].ChainSelector == selector {
			return &networks[i], nil
		}
	}
	return nil, fmt.Errorf("chain selector %d not found in network config", selector)
}

// ResolveAddressesForChain extracts addresses for a given chain selector
// from the cldf address book, keyed by contract type.
func ResolveAddressesForChain(book cldf.AddressBook, chainSelector uint64) (map[string]string, error) {
	addresses, err := book.AddressesForChain(chainSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses for chain %d: %w", chainSelector, err)
	}

	result := make(map[string]string, len(addresses))
	for addr, tv := range addresses {
		result[string(tv.Type)] = addr
	}
	return result, nil
}

// LoadFullConfig loads all four config layers for a run.
//
//  1. Loads runs/<runName>.toml for run parameters
//  2. Loads .env.<env> for secrets
//  3. Loads addresses-<env>.json for contract addresses
//  4. Loads networks-<env>.yaml for chain RPC endpoints
func LoadFullConfig(runName string) (*LoadTestConfig, error) {
	// Layer 4: Run config
	runCfg, err := LoadRunConfig(runName)
	if err != nil {
		return nil, fmt.Errorf("layer 4 (run config): %w", err)
	}

	// Canonicalize the run name (strip .toml extension) for results/logs.
	runName = strings.TrimSuffix(runName, ".toml")

	envName := runCfg.Run.Env

	// Layer 1: Secrets
	suiPrivKey, evmPrivKey, walletSeed, err := LoadEnvConfig(envName)
	if err != nil {
		return nil, fmt.Errorf("layer 1 (env): %w", err)
	}

	// Layer 2: Addresses
	addressBook, err := LoadAddressBook(envName)
	if err != nil {
		return nil, fmt.Errorf("layer 2 (addresses): %w", err)
	}

	// Layer 3: Networks
	networks, err := LoadNetworks(envName)
	if err != nil {
		return nil, fmt.Errorf("layer 3 (networks): %w", err)
	}

	// Parse chain selectors from string to uint64
	srcSelector, err := parseChainSelector(runCfg.Run.SourceChainSelector)
	if err != nil {
		return nil, fmt.Errorf("layer 4 (run config): %w", err)
	}
	destSelector, err := parseChainSelector(runCfg.Run.DestChainSelector)
	if err != nil {
		return nil, fmt.Errorf("layer 4 (run config): %w", err)
	}

	// Validate source and dest chains exist in network config
	if _, err := FindNetworkBySelector(networks, srcSelector); err != nil {
		return nil, fmt.Errorf("source chain %d: %w", srcSelector, err)
	}
	if _, err := FindNetworkBySelector(networks, destSelector); err != nil {
		return nil, fmt.Errorf("dest chain %d: %w", destSelector, err)
	}

	cfg := &LoadTestConfig{
		RunName:             runName,
		EnvName:             envName,
		SourceChainSelector: srcSelector,
		DestChainSelector:   destSelector,
		MessageCount:        runCfg.Run.MessageCount,
		MessageData:         []byte(runCfg.Run.MessageData),
		ReceiverAddress:     runCfg.Receiver.Address,
		SuiGasBudget:        runCfg.Gas.SuiGasBudget,
		EvmGasLimit:         runCfg.Gas.EvmGasLimit,
		EvmCallbackGasLimit: runCfg.Gas.EvmCallbackGasLimit,
		LoadRPS:             runCfg.Load.RPS,
		LoadWallets:         runCfg.Load.Wallets,
		SuiPrivateKey:       suiPrivKey,
		EVMPrivateKey:       evmPrivKey,
		WalletSeed:          walletSeed,
		AddressBook:         addressBook,
		Networks:            networks,
	}

	if runCfg.Token != nil {
		var tokenIdentifier string
		if runCfg.Token.CoinMetadataID != "" {
			tokenIdentifier = runCfg.Token.CoinMetadataID
		} else {
			tokenIdentifier = runCfg.Token.TokenAddress
		}
		cfg.TokenConfig = &TokenTransferConfig{
			TokenIdentifier: tokenIdentifier,
			Amount:          runCfg.Token.Amount,
			Mode:            runCfg.Token.Mode,
		}
	}

	if runCfg.SuiReceiver != nil {
		cfg.SuiReceiverConfig = &SuiReceiverConfig{
			PackageID: runCfg.SuiReceiver.PackageID,
		}
	}

	slog.Info("Config loaded successfully",
		"runName", cfg.RunName,
		"env", cfg.EnvName,
		"sourceChain", cfg.SourceChainSelector,
		"destChain", cfg.DestChainSelector,
		"messageCount", cfg.MessageCount,
		"tokenMode", cfg.TokenConfig,
		"suiReceiver", cfg.SuiReceiverConfig,
	)

	return cfg, nil
}

// SaveResults writes run results to a file.
// File: results/<runName>-<env>-<timestamp>.txt
func SaveResults(results *RunResults) error {
	if err := os.MkdirAll("results", 0755); err != nil {
		return fmt.Errorf("failed to create results directory: %w", err)
	}

	timestamp := time.Now().Format("20060102T150405")
	filename := fmt.Sprintf("%s-%s-%s.txt", results.RunName, results.EnvName, timestamp)
	path := filepath.Join("results", filename)

	data, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write results file: %w", err)
	}

	slog.Info("Results saved", "path", path)
	return nil
}
