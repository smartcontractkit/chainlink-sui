package client

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"

	"golang.org/x/crypto/blake2b"
)

func GetAddressFromPublicKey(pubKey []byte) (string, error) {
	if len(pubKey) != ed25519.PublicKeySize {
		return "", fmt.Errorf("invalid public key length, expected %d got %d", ed25519.PublicKeySize, len(pubKey))
	}

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
