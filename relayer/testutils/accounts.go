package testutils

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"

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

func DrainAccountCoins(t *testing.T, ctx context.Context, lgr logger.Logger, accountAddress string, suiKeystore loop.Keystore, cli *client.PTBClient, suiCoins []models.CoinData, receiver string) error {
	// TODO: this should reuse transaction.go's GeneratePTBTransaction
	lgr.Infow("Draining account coins from account address", "accountAddress", accountAddress)

	accounts, err := suiKeystore.Accounts(ctx)
	require.NoError(t, err, "Failed to get accounts")
	if len(accounts) != 1 {
		return fmt.Errorf("expected 1 account, got %d", len(accounts))
	}

	publicKeyStr := accounts[0]
	publicKeyBytes, err := hex.DecodeString(publicKeyStr)
	if err != nil {
		return fmt.Errorf("failed to decode public key: %w", err)
	}

	expectedAddress, err := client.GetAddressFromPublicKey(publicKeyBytes)
	if err != nil {
		return fmt.Errorf("failed to get address from public key: %w", err)
	}

	if expectedAddress != accountAddress {
		return fmt.Errorf("expected account address %s, got %s", expectedAddress, accountAddress)
	}

	testKeystore, ok := suiKeystore.(*TestKeystore)
	require.True(t, ok, "Expected TestKeystore")

	txnSigner := testKeystore.GetSuiSigner(context.Background(), publicKeyStr)

	// Create new transaction
	tx := transaction.NewTransaction()

	// Convert coin data to object references for the transaction
	coinObjectRefs := make([]string, 0)
	for _, coin := range suiCoins {
		coinObjectRefs = append(coinObjectRefs, coin.CoinObjectId)
	}

	// Add PayAllSui command to transfer all coins to receiver
	err = cli.PayAllSui(ctx, receiver, coinObjectRefs, accountAddress)
	if err != nil {
		return fmt.Errorf("failed to add PayAllSui command: %w", err)
	}

	// Execute the transaction using the PTB client
	response, err := cli.FinishPTBAndSend(ctx, txnSigner, tx, client.WaitForLocalExecution)
	if err != nil {
		return fmt.Errorf("failed to execute drain transaction: %w", err)
	}

	lgr.Infow("Successfully drained account coins",
		"accountAddress", accountAddress,
		"receiver", receiver,
		"txDigest", response.TxDigest,
	)

	return nil
}
