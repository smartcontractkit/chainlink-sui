package keystore

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	suiSigner "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

var SuiDefaultKeystorePath = os.Getenv("HOME") + "/.sui/sui_config/sui.keystore"

const keyWithSchemeLength = 33

// SignerType is an integer-based enum for different signer implementations.
type SignerType int

const (
	// PrivateKeySigner represents a signer that uses a local private key.
	PrivateKeySigner SignerType = iota
	// UnknownSigner represents an unknown signer type.
	UnknownSigner SignerType = -1
)

type Keystore interface {
	GetSignerFromAddress(address string) (suiSigner.SuiSigner, error)
	KeyStorePath() string
	SignerType() SignerType
}

type SuiKeystore struct {
	logger       logger.Logger
	keyStorePath string
	signerType   SignerType
}

func NewSuiKeystore(log logger.Logger, keyStorePath string, signerType SignerType) (SuiKeystore, error) {
	if keyStorePath == "" {
		keyStorePath = SuiDefaultKeystorePath
	}

	return SuiKeystore{
		logger:       log,
		keyStorePath: keyStorePath,
		signerType:   signerType,
	}, nil
}

func (s SuiKeystore) GetSignerFromAddress(address string) (suiSigner.SuiSigner, error) {
	switch s.signerType {
	case PrivateKeySigner:
		return s.buildPrivateKeySigner(address)
	case UnknownSigner:
		return nil, fmt.Errorf("unknown signer type: %d", s.signerType)
	default:
		return nil, fmt.Errorf("unsupported signer type: %d", s.signerType)
	}
}

func (s SuiKeystore) buildPrivateKeySigner(address string) (suiSigner.SuiSigner, error) {
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
			return suiSigner.NewPrivateKeySigner(ed25519PrivateKey), nil
		}
	}

	return nil, fmt.Errorf("no private key found for address: %s", address)
}

func (s SuiKeystore) KeyStorePath() string {
	return s.keyStorePath
}

func (s SuiKeystore) SignerType() SignerType {
	return s.signerType
}
