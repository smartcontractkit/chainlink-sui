package keystore

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

var SuiDefaultKeystorePath = os.Getenv("HOME") + "/.sui/sui_config/sui.keystore"

const keyWithSchemeLength = 33

type Keystore interface {
	GetPrivateKeyFromAddress(address string) (ed25519.PrivateKey, error)
}

type SuiKeystore struct {
	logger       logger.Logger
	keyStorePath string
}

func NewSuiKeystore(log logger.Logger, keyStorePath string) (SuiKeystore, error) {
	if keyStorePath == "" {
		keyStorePath = SuiDefaultKeystorePath
	}

	return SuiKeystore{
		logger:       log,
		keyStorePath: keyStorePath,
	}, nil
}

func (s SuiKeystore) GetPrivateKeyFromAddress(address string) (ed25519.PrivateKey, error) {
	data, err := os.ReadFile(s.keyStorePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read keystore file: %w", err)
	}

	// Parse keystore JSON (Base64-encoded private keys)
	var keys []string
	if err := json.Unmarshal(data, &keys); err != nil {
		return nil, fmt.Errorf("failed to parse keystore JSON: %w", err)
	}

	for _, encodedKey := range keys {
		// Decode Base64 private key
		privateKeyBytes, err := base64.StdEncoding.DecodeString(encodedKey)
		if err != nil {
			continue
		}

		// Trim the key if it's 33 bytes (includes a key scheme identifier)
		if len(privateKeyBytes) == keyWithSchemeLength {
			privateKeyBytes = privateKeyBytes[1:] // Remove the first byte (scheme flag)
		}

		// Create signer and check if it matches the target Sui address
		signerAccount := signer.NewSigner(privateKeyBytes)
		if signerAccount.Address == address {
			ed25519PrivateKey := ed25519.NewKeyFromSeed(privateKeyBytes)
			return ed25519PrivateKey, nil
		}
	}

	return nil, fmt.Errorf("no matching private key found for address %s", address)
}
