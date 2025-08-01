package testutils

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
)

const (
	ACCOUNT_1_SEED = "0000000000000000000000000000000000000000000000000000000000000001"
	ACCOUNT_2_SEED = "0000000000000000000000000000000000000000000000000000000000000002"
	ACCOUNT_3_SEED = "0000000000000000000000000000000000000000000000000000000000000003"
)

const PUBKEY_LENGTH = 32
const ADDRESS_LENGTH = 20
const KEYPAIR_LENGTH = 64

type ECDSAKeyPair struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  ecdsa.PublicKey
	Address    []byte
}

// GenerateKeyPair generates a new Ethereum key pair for testing
func GenerateKeyPair() (*ecdsa.PrivateKey, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate Ethereum key pair: %w", err)
	}

	return privateKey, nil
}

// GenerateFromHexSeed generates a deterministic Ethereum key pair from a hex seed
func GenerateFromHexSeed(seed string) (*ecdsa.PrivateKey, error) {
	seedBytes, err := hex.DecodeString(seed)
	if err != nil {
		return nil, fmt.Errorf("invalid hex seed: %w", err)
	}

	if len(seedBytes) != PUBKEY_LENGTH {
		return nil, fmt.Errorf("invalid seed length: got %d, want 32", len(seedBytes))
	}

	privateKey, err := crypto.ToECDSA(seedBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key from seed: %w", err)
	}

	return privateKey, nil
}

// VerifySignatureRecovery verifies that a given signature and a signed hash can be used to recover the provided address
//
//nolint:all
func VerifySignatureRecovery(t *testing.T, signature []byte, signedHash []byte, address []byte) []byte {
	t.Helper()

	// 1. Normalize the v value exactly like the contract
	v := signature[64]
	if v == 27 {
		v = 0
	} else if v == 28 {
		v = 1
	} else if v > 35 {
		v = (v - 1) % 2
	}

	// Create new signature with normalized v
	sigWithNormalizedV := make([]byte, KEYPAIR_LENGTH+1)
	copy(sigWithNormalizedV[:KEYPAIR_LENGTH], signature[:KEYPAIR_LENGTH])
	sigWithNormalizedV[64] = v

	// 2. Recover public key
	pubKey, err := crypto.Ecrecover(signedHash, sigWithNormalizedV)
	t.Log("Recovered uncompressed", pubKey)
	require.NoError(t, err, "Failed to recover public key")

	// 3. Decompress public key (equivalent to contract's decompress_pubkey)
	// Note: Go's Ecrecover already returns uncompressed public key
	// We need to take bytes 1-64 (skipping the first byte which is the format)
	uncompressed64 := pubKey[1 : KEYPAIR_LENGTH+1]

	// 4. Hash the 64-byte uncompressed public key
	hash := crypto.Keccak256(uncompressed64)

	// 5. Take last 20 bytes (equivalent to contract's last 20 bytes of hash)
	recoveredAddr := hash[12:PUBKEY_LENGTH]

	// Debug output
	t.Logf("Original signature: %x", signature)
	t.Logf("Normalized v: %d", v)
	t.Logf("Signature with normalized v: %x", sigWithNormalizedV)
	t.Logf("Recovered public key: %x", pubKey)
	t.Logf("Uncompressed 64 bytes: %x", uncompressed64)
	t.Logf("Hash: %x", hash)
	t.Logf("Recovered address: %x", recoveredAddr)

	require.Equal(t, recoveredAddr, address, "Recovered address does not match signed hash")

	return recoveredAddr
}
