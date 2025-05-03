package client

import (
	"encoding/base64"
	"fmt"

	"github.com/pattonkan/sui-go/suisigner"

	"github.com/smartcontractkit/chainlink-sui/relayer/codec"
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

func ToSuiSignatures(signatures []string) ([]*suisigner.Signature, error) {
	suiSignatures := make([]*suisigner.Signature, 0)
	for _, sig := range signatures {
		bytes, err := toSuiSignature(sig)
		if err != nil {
			return nil, err
		}
		suiSignatures = append(suiSignatures, bytes)
	}

	return suiSignatures, nil
}

func toSuiSignature(sig string) (*suisigner.Signature, error) {
	decoded, err := codec.DecodeBase64(sig)
	if err != nil {
		return nil, err
	}
	const lengthSig = 97

	if len(decoded) != lengthSig {
		return nil, fmt.Errorf("invalid signature length: expected 97, got %d", len(decoded))
	}
	var sigArray [lengthSig]byte
	copy(sigArray[:], decoded)

	return &suisigner.Signature{
		Ed25519SuiSignature: &suisigner.Ed25519SuiSignature{
			Signature: sigArray,
		},
	}, nil
}
