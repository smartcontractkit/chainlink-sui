package testutils

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

// GenerateAccountKeyPair Generates a public/private keypair with the ed25519 signature algorithm, then derives the address from the public key.
// Returns (private key, public key, address, error).
func GenerateAccountKeyPair(t *testing.T) (ed25519.PrivateKey, ed25519.PublicKey, string, error) {
	t.Helper()

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err, "Failed to generate new account")

	accountAddress, err := client.GetAddressFromPublicKey([]byte(publicKey))
	require.NoError(t, err, "Failed to get address from public key")

	t.Logf("Created account, publicKey: %s, accountAddress: %s", hex.EncodeToString([]byte(publicKey)), accountAddress)

	return privateKey, publicKey, accountAddress, nil
}
