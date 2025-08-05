package testutils

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"gopkg.in/yaml.v3"
)

// NewTestKeystore creates a new test keystore
func NewTestKeystore(t *testing.T) *TestKeystore {
	t.Helper()
	return &TestKeystore{t: t, Keys: map[string]ed25519.PrivateKey{}}
}

// TestKeystore is a simple keystore for testing
type TestKeystore struct {
	t    *testing.T
	Keys map[string]ed25519.PrivateKey
}

var _ loop.Keystore = &TestKeystore{}

// AddKey adds a private key to the keystore
func (tk *TestKeystore) AddKey(key ed25519.PrivateKey) {
	publicKey := fmt.Sprintf("%064x", key.Public())
	tk.Keys[publicKey] = key
}

func (tk *TestKeystore) Sign(ctx context.Context, id string, hash []byte) ([]byte, error) {
	privateKey, ok := tk.Keys[id]
	if !ok {
		tk.t.Fatalf("No such key: %s", id)
	}

	// used to check if the account exists.
	if hash == nil {
		return nil, nil
	}

	return ed25519.Sign(privateKey, hash), nil
}

func (tk *TestKeystore) Accounts(ctx context.Context) ([]string, error) {
	accounts := make([]string, 0, len(tk.Keys))
	for id := range tk.Keys {
		accounts = append(accounts, id)
	}
	return accounts, nil
}

func (tk *TestKeystore) GetSuiSigner(ctx context.Context, publicKey string) *signer.Signer {
	privateKey, ok := tk.Keys[publicKey]
	if !ok {
		tk.t.Fatalf("No such key: %s", publicKey)
	}

	// Extract the 32-byte seed from the 64-byte ed25519 private key
	seed := privateKey.Seed()
	return signer.NewSigner(seed)
}

func GetAccountAndKeyFromSui(testKeystore *TestKeystore) (string, []byte) {
	keystorePath := filepath.Join(os.Getenv("HOME"), ".sui", "sui_config", "sui.keystore")
	signers := make([]*signer.Signer, 0)

	// Add accounts from CLI keystore if it exists
	if _, err := os.Stat(keystorePath); err == nil {
		signers, err = addSuiCLIAccountsToKeystore(testKeystore, keystorePath)
		if err != nil {
			testKeystore.t.Fatalf("Failed to add Sui CLI accounts: %v", err)
		}

		if len(signers) == 0 {
			testKeystore.t.Logf("No accounts found in CLI keystore")
		}
	} else {
		testKeystore.t.Logf("No accounts found in CLI keystore, generating new account")

		// Either keystore doesn't exist or no keys were loaded, generate a new one
		privateKey, _, _, err := GenerateAccountKeyPair(testKeystore.t)
		if err != nil {
			testKeystore.t.Fatalf("Failed to generate account key pair: %v", err)
		}

		testKeystore.AddKey(privateKey)
	}

	accounts, err := testKeystore.Accounts(context.Background())
	if err != nil {
		testKeystore.t.Fatalf("Failed to get accounts: %v", err)
	}
	if len(accounts) == 0 {
		testKeystore.t.Fatalf("No accounts found in keystore")
	}

	// Get the active address from the sui CLI
	activeAddress, err := GetActiveAddressFromSuiConfig()
	if err != nil {
		testKeystore.t.Fatalf("Failed to get active address from Sui config: %v", err)
	}
	if activeAddress == "" {
		testKeystore.t.Fatalf("No active address found in Sui config")
	}

	// Find the key associated with the `active-address` in the sui CLI
	publicKeyBytes := make([]byte, 32)
	for _, sgnr := range signers {
		if sgnr.Address == activeAddress {
			copy(publicKeyBytes, sgnr.PriKey.Public().(ed25519.PublicKey))
			testKeystore.t.Logf("Found active address: %s", activeAddress)
			break
		}
	}

	return activeAddress, publicKeyBytes
}

const keyWithSchemeLength = 33

func addSuiCLIAccountsToKeystore(testKeystore *TestKeystore, keystorePath string) ([]*signer.Signer, error) {
	data, err := os.ReadFile(keystorePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read keystore file: %w", err)
	}

	// Parse keystore JSON (Base64-encoded private keys)
	var keys []string
	if err := json.Unmarshal(data, &keys); err != nil {
		return nil, fmt.Errorf("failed to parse keystore JSON: %w", err)
	}

	var signers []*signer.Signer
	for _, encodedKey := range keys {
		// Decode Base64 private key
		privateKeyBytes, err := base64.StdEncoding.DecodeString(encodedKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decode private key: %w", err)
		}

		// Trim the key if it's 33 bytes (includes a key scheme identifier)
		if len(privateKeyBytes) == keyWithSchemeLength {
			privateKeyBytes = privateKeyBytes[1:] // Remove the first byte (scheme flag)
		}

		sgnr := signer.NewSigner(privateKeyBytes)

		signers = append(signers, sgnr)
		testKeystore.AddKey(sgnr.PriKey)
	}

	return signers, nil
}

// GetActiveAddressFromSuiConfig retrieves the active address from the Sui client configuration file
// located at ~/.sui/sui_config/client.yaml
func GetActiveAddressFromSuiConfig() (string, error) {
	configPath := filepath.Join(os.Getenv("HOME"), ".sui", "sui_config", "client.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("failed to read config file: %w", err)
	}

	var config struct {
		ActiveAddress string `yaml:"active_address"`
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return "", fmt.Errorf("failed to parse config file: %w", err)
	}

	return config.ActiveAddress, nil
}
