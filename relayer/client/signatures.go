package client

import (
	"encoding/base64"
)

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
