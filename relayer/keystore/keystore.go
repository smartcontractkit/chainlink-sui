package keystore

import (
	"context"
	"crypto"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"golang.org/x/crypto/blake2b"
)

var SuiDefaultKeystorePath = os.Getenv("HOME") + "/.sui/sui_config/sui.keystore"

const keyWithSchemeLength = 33

type SuiKeystore struct {
	logger       logger.Logger
	keyStorePath string
}

func NewSuiKeystore(lggr logger.Logger, keyStorePath string) (SuiKeystore, error) {
	if keyStorePath == "" {
		keyStorePath = SuiDefaultKeystorePath
	}

	return SuiKeystore{
		logger:       lggr,
		keyStorePath: keyStorePath,
	}, nil
}

func (s SuiKeystore) Accounts(ctx context.Context) (accounts []string, err error) {
	return nil, errors.New("not implemented")
}

func (s SuiKeystore) Sign(ctx context.Context, account string, data []byte) (signed []byte, err error) {
	privateKey, err := s.GetPrivateKeyByAddress(account)
	if err != nil {
		return nil, err
	}
	signedPayload, err := sign_sui_message(data, privateKey)
	if err != nil {
		return nil, err
	}

	return signedPayload, nil
}

func (s SuiKeystore) GetPrivateKeyByAddress(address string) (ed25519.PrivateKey, error) {
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
			return ed25519.NewKeyFromSeed(privateKeyBytes), nil
		}
	}

	return nil, fmt.Errorf("no private key found for address: %s", address)
}

func (s SuiKeystore) KeyStorePath() string {
	return s.keyStorePath
}

type SigFlag byte

const (
	SigFlagEd25519 SigFlag = 0x00
)

// SerializeSuiSignature formats and serializes a signature for use with Sui transactions.
//
// This function follows the Sui transaction signature format specification:
// 1. A one-byte flag indicating the signature scheme (0x00 for Ed25519)
// 2. The raw signature bytes
// 3. The public key bytes
// These components are concatenated and then base64 encoded to produce the final signature string
// that can be submitted to the Sui network.
//
// Based on the implementation from block-vision's sui-go-sdk:
// https://github.com/block-vision/sui-go-sdk/blob/main/models/signature.go#L140
//
// Parameters:
//
//	signature - The raw signature bytes from the ed25519 Sign operation
//	pubKey    - The public key corresponding to the private key that produced the signature
//
// Returns:
//
//	A base64-encoded string containing the serialized signature ready for submission to Sui
func SerializeSuiSignature(signature, pubKey []byte) string {
	signatureLen := len(signature)
	pubKeyLen := len(pubKey)
	serializedSignature := make([]byte, 1+signatureLen+pubKeyLen)
	serializedSignature[0] = byte(SigFlagEd25519)
	copy(serializedSignature[1:], signature)
	copy(serializedSignature[1+signatureLen:], pubKey)

	return base64.StdEncoding.EncodeToString(serializedSignature)
}

var IntentBytes = []byte{0, 0, 0}

// Sign implements the SuiSigner interface for PrivateKeySigner.
// This is a port of the block vision's implementation.
// Check the full code [here](https://github.com/block-vision/sui-go-sdk/blob/main/models/signature.go#L117)
//
// This method follows the Sui signing protocol by:
// 1. Prepending the standard intent bytes (0,0,0) to the message
// 2. Computing the blake2b-256 hash of the intent+message
// 3. Signing the resulting hash with the ed25519 private key
//
// Parameters:
//
//	message - The raw message bytes to sign
//
// Returns:
//
//	The ed25519 signature bytes and any error encountered during signing
func sign_sui_message(message []byte, privateKey ed25519.PrivateKey) ([]byte, error) {
	messageWithIntent := messageWithIntent(message)
	digest := blake2b.Sum256(messageWithIntent)
	var noHash crypto.Hash
	sigBytes, err := privateKey.Sign(nil, digest[:], noHash)
	if err != nil {
		return nil, fmt.Errorf("failed to sign message: %w", err)
	}

	return sigBytes, nil
}

// TODO: add docstring and reference to block-vision implementation
func messageWithIntent(message []byte) []byte {
	intent := IntentBytes
	intentMessage := make([]byte, len(intent)+len(message))
	copy(intentMessage, intent)
	copy(intentMessage[len(intent):], message)

	return intentMessage
}
