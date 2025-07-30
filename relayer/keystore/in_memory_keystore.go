package keystore

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"
	"sync"

	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// Keystore interface defines the required methods for keystore implementations
type Keystore interface {
	Accounts(ctx context.Context) (accounts []string, err error)
	// Sign returns data signed by account.
	// nil data can be used as a no-op to check for account existence.
	Sign(ctx context.Context, account string, data []byte) (signed []byte, err error)
}

// InMemoryKeystore implements the Keystore interface using in-memory storage
type InMemoryKeystore struct {
	logger logger.Logger
	mu     sync.RWMutex
	keys   map[string]ed25519.PrivateKey // address -> private key mapping
}

// NewInMemoryKeystore creates a new in-memory keystore
func NewInMemoryKeystore(lggr logger.Logger) *InMemoryKeystore {
	lggr.Debugw("Creating in-memory keystore")

	return &InMemoryKeystore{
		logger: lggr,
		keys:   make(map[string]ed25519.PrivateKey),
	}
}

// AddKey adds a private key to the in-memory keystore
// The key should be base64-encoded. The address is derived from the private key's public key
func (s *InMemoryKeystore) AddKey(base64Key string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Decode Base64 private key
	privateKeyBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 key: %w", err)
	}

	// Trim the key if it's 33 bytes (includes a key scheme identifier)
	if len(privateKeyBytes) == keyWithSchemeLength {
		privateKeyBytes = privateKeyBytes[1:] // Remove the first byte (scheme flag)
	}

	// Validate key length
	if len(privateKeyBytes) != ed25519.SeedSize {
		return "", fmt.Errorf("invalid key length: expected %d bytes, got %d", ed25519.SeedSize, len(privateKeyBytes))
	}

	// Create private key from seed
	privateKey := ed25519.NewKeyFromSeed(privateKeyBytes)

	// Derive the address from the private key
	signerAccount := signer.NewSigner(privateKeyBytes)
	address := signerAccount.Address

	s.keys[address] = privateKey
	s.logger.Debugw("Added key to in-memory keystore", "address", address)

	return address, nil
}

// Accounts returns all account addresses stored in the keystore
func (s *InMemoryKeystore) Accounts(ctx context.Context) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	accounts := make([]string, 0, len(s.keys))
	for address := range s.keys {
		accounts = append(accounts, address)
	}

	return accounts, nil
}

// Sign signs the given data with the private key associated with the account
// If data is nil, this is used as a no-op to check for account existence
func (s *InMemoryKeystore) Sign(ctx context.Context, account string, data []byte) ([]byte, error) {
	s.mu.RLock()
	privateKey, exists := s.keys[account]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("no private key found for account: %s", account)
	}

	// If data is nil, return empty signature (no-op to check account existence)
	if data == nil {
		return []byte{}, nil
	}

	signedPayload, err := sign_sui_message(data, privateKey)
	if err != nil {
		return nil, err
	}

	return signedPayload, nil
}
