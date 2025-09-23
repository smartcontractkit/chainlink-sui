package events

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

// GenerateAccountKeyPair Generates a public/private keypair with the ed25519 signature algorithm, then derives the address from the public key.
// Returns (private key, public key, address, error).
func GenerateAccountKeyPair(lggr logger.Logger) (ed25519.PrivateKey, ed25519.PublicKey, string, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		lggr.Errorw("Failed to generate account", "error", err)
		return ed25519.PrivateKey{}, ed25519.PublicKey{}, "", err
	}

	accountAddress, err := client.GetAddressFromPublicKey([]byte(publicKey))
	if err != nil {
		lggr.Errorw("Failed to get address from public key", "error", err)
		return ed25519.PrivateKey{}, ed25519.PublicKey{}, "", err
	}

	lggr.Infow("Created account, publicKey: %s, accountAddress: %s", hex.EncodeToString([]byte(publicKey)), accountAddress)

	return privateKey, publicKey, accountAddress, nil
}

func GenerateAccount(lggr logger.Logger) (rel.SuiSigner, string, ed25519.PublicKey, error) {
	pk, pubKey, accountAddress, err := GenerateAccountKeyPair(lggr)

	if err != nil {
		lggr.Errorw("Failed to generate account", "error", err)
		return nil, "", ed25519.PublicKey{}, err
	}

	lggr.Infow("Created account, publicKey: %s, accountAddress: %s", hex.EncodeToString([]byte(pubKey)), accountAddress)

	signer := rel.NewPrivateKeySigner(pk)

	return signer, accountAddress, pubKey, nil
}
