package signer

import (
	"crypto"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/hex"
	"log"

	"golang.org/x/crypto/blake2b"
)

type SigFlag byte

const (
	SigFlagEd25519 SigFlag = 0x00
)

// SuiSigner defines an interface for signing messages in the Sui blockchain format.
// This allow us to have multiple signing implementations, such as using a private key or a hardware wallet.
type SuiSigner interface {
	// Sign signs the given message and returns the serialized signature.
	Sign(message []byte) ([]string, error)

	// GetAddress returns the Sui address derived from the signer's public key
	GetAddress() (string, error)
}

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

type PrivateKeySigner struct {
	privateKey ed25519.PrivateKey
}

var IntentBytes = []byte{0, 0, 0}

func NewPrivateKeySigner(privateKey ed25519.PrivateKey) *PrivateKeySigner {
	return &PrivateKeySigner{
		privateKey: privateKey,
	}
}

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
func (s *PrivateKeySigner) Sign(message []byte) ([]string, error) {
	messageWithIntent := messageWithIntent(message)
	digest := blake2b.Sum256(messageWithIntent)
	var noHash crypto.Hash
	sigBytes, err := s.privateKey.Sign(nil, digest[:], noHash)
	if err != nil {
		log.Fatal(err)
	}
	pubKey := s.privateKey.Public().(ed25519.PublicKey)
	serializedSignature := SerializeSuiSignature(sigBytes, pubKey)

	return []string{serializedSignature}, nil
}

// GetAddress returns the Sui address derived from the signer's public key
func (s *PrivateKeySigner) GetAddress() (string, error) {
	pubKey := s.privateKey.Public().(ed25519.PublicKey)

	// Prepend the Ed25519 flag byte to the public key
	flaggedPubKey := make([]byte, 1+len(pubKey))
	flaggedPubKey[0] = byte(SigFlagEd25519)
	copy(flaggedPubKey[1:], pubKey)

	// Hash the flagged public key
	digest := blake2b.Sum256(flaggedPubKey)

	// Take the first 20 bytes of the hash as the address
	addressBytes := digest[:32]

	// Convert to hex string with "0x" prefix
	address := "0x" + hex.EncodeToString(addressBytes)

	return address, nil
}

// TODO: add docstring and reference to block-vision implementation
func messageWithIntent(message []byte) []byte {
	intent := IntentBytes
	intentMessage := make([]byte, len(intent)+len(message))
	copy(intentMessage, intent)
	copy(intentMessage[len(intent):], message)

	return intentMessage
}
