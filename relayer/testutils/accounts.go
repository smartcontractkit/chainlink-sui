package testutils

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// LoadAccountFromEnv loads a test account from environment variables
func LoadAccountFromEnv(t *testing.T, log logger.Logger) (ed25519.PrivateKey, ed25519.PublicKey, string) {
	t.Helper()
	// First try to load from private key
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex != "" {
		privateKey, err := hex.DecodeString(privateKeyHex)
		if err != nil {
			t.Fatal(fmt.Errorf("invalid PRIVATE_KEY format: %w", err))
		}

		if len(privateKey) != ed25519.PrivateKeySize {
			t.Fatal(fmt.Errorf("invalid PRIVATE_KEY length, expected %d got %d", ed25519.PrivateKeySize, len(privateKey)))
		}

		publicKey := privateKey[32:]
		address := DeriveAddressFromPublicKey(publicKey)

		log.Debugw("Loaded account from PRIVATE_KEY", "address", address)

		return privateKey, publicKey, address
	}

	// Then try to load from address
	address := os.Getenv("ADDRESS")
	if address != "" {
		log.Debugw("Only ADDRESS provided, can't use for signing", "address", address)
		return nil, nil, address
	}

	return nil, nil, ""
}

// GenerateAccountKeyPair Generates a public/private keypair with the ed25519 signature algorithm, then derives the address from the public key.
// Returns (private key, public key, address, error).
func GenerateAccountKeyPair(t *testing.T, log logger.Logger) (ed25519.PrivateKey, ed25519.PublicKey, string, error) {
	t.Helper()

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err, "Failed to generate new account")

	// Generate Sui address from public key
	accountAddress := DeriveAddressFromPublicKey(publicKey)

	log.Debugw("Created account", "publicKey", hex.EncodeToString([]byte(publicKey)), "accountAddress", accountAddress)

	return privateKey, publicKey, accountAddress, nil
}

// DeriveAddressFromPublicKey derives a Sui address from an ed25519 public key
func DeriveAddressFromPublicKey(publicKey ed25519.PublicKey) string {
	return "0x" + hex.EncodeToString(publicKey)
}

// NewTestKeystore creates a new test keystore
func NewTestKeystore(t *testing.T) *TestKeystore {
	t.Helper()
	return &TestKeystore{t: t, keys: map[string]ed25519.PrivateKey{}}
}

// TestKeystore is a simple keystore for testing
type TestKeystore struct {
	t    *testing.T
	keys map[string]ed25519.PrivateKey
}

// AddKey adds a private key to the keystore
func (k *TestKeystore) AddKey(key ed25519.PrivateKey) {
	// Derive address from private key (in Sui, address is derived from public key)
	publicKey := key.Public().(ed25519.PublicKey)
	address := DeriveAddressFromPublicKey(publicKey)
	k.keys[address] = key
}

// Get returns a private key by address
func (k *TestKeystore) Get(address string) (ed25519.PrivateKey, bool) {
	key, ok := k.keys[address]
	return key, ok
}
