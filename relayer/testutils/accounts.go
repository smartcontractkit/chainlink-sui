package testutils

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/fardream/go-bcs/bcs"
	"github.com/pattonkan/sui-go/sui/suiptb"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/stretchr/testify/require"

	suiAlt "github.com/pattonkan/sui-go/sui"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

// LoadAccountFromEnv loads a test account from environment variables
func LoadAccountFromEnv(t *testing.T, log logger.Logger) (ed25519.PrivateKey, ed25519.PublicKey, string) {
	t.Helper()
	// First try to load from private key
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex != "" {
		privateKey, err := hex.DecodeString(privateKeyHex)
		if err != nil {
			t.Fatal(fmt.Errorf("invalid PRIVATE_KEY format: %w", err))
		}

		if len(privateKey) != ed25519.PrivateKeySize {
			t.Fatal(fmt.Errorf("invalid PRIVATE_KEY length, expected %d got %d", ed25519.PrivateKeySize, len(privateKey)))
		}

		publicKey := privateKey[32:]
		address := DeriveAddressFromPublicKey(publicKey)

		log.Debugw("Loaded account from PRIVATE_KEY", "address", address)

		return privateKey, publicKey, address
	}

	// Then try to load from address
	address := os.Getenv("ADDRESS")
	if address != "" {
		log.Debugw("Only ADDRESS provided, can't use for signing", "address", address)
		return nil, nil, address
	}

	return nil, nil, ""
}

func GetAccountAndKeyFromSui(t *testing.T, lgr logger.Logger) string {
	t.Helper()
	cmd := exec.Command("sui", "client", "active-address")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Failed to get active address: %s", string(output))

	// Extract the account address using regex to handle any potential formatting
	re := regexp.MustCompile(`0x[a-fA-F0-9]+`)
	matches := re.FindString(string(output))
	if matches == "" {
		require.Fail(t, "Failed to extract account address from output: %s", string(output))
	}
	accountAddress := matches
	lgr.Info("Active address: ", accountAddress)

	return accountAddress
}

// GenerateAccountKeyPair Generates a public/private keypair with the ed25519 signature algorithm, then derives the address from the public key.
// Returns (private key, public key, address, error).
func GenerateAccountKeyPair(t *testing.T, log logger.Logger) (ed25519.PrivateKey, ed25519.PublicKey, string, error) {
	t.Helper()

	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err, "Failed to generate new account")

	// Generate Sui address from public key
	accountAddress := DeriveAddressFromPublicKey(publicKey)

	log.Debugw("Created account", "publicKey", hex.EncodeToString([]byte(publicKey)), "accountAddress", accountAddress)

	return privateKey, publicKey, accountAddress, nil
}

// DeriveAddressFromPublicKey derives a Sui address from an ed25519 public key
func DeriveAddressFromPublicKey(publicKey ed25519.PublicKey) string {
	return "0x" + hex.EncodeToString(publicKey)
}

func DrainAccountCoins(ctx context.Context, lgr logger.Logger, signerInstance *signer.SuiSigner, cli client.SuiClient, suiCoins []models.CoinData, receiver string) error {
	addr, err := (*signerInstance).GetAddress()
	if err != nil {
		return fmt.Errorf("failed to get address: %w", err)
	}
	senderAddress, _ := suiAlt.AddressFromHex(addr)
	receiverAddresss, _ := suiAlt.AddressFromHex(receiver)

	coins := make([]*suiAlt.ObjectRef, 0)

	for _, coin := range suiCoins {
		objectId := suiAlt.MustObjectIdFromHex(coin.CoinObjectId)

		version, _ := strconv.ParseUint(coin.Version, 10, 64)
		digest := suiAlt.MustNewDigest(coin.Digest)

		coinObject := &suiAlt.ObjectRef{
			ObjectId: objectId,
			Version:  version,
			Digest:   digest,
		}
		coins = append(coins, coinObject)
	}

	ptb := suiptb.NewTransactionDataTransactionBuilder()
	_ = ptb.PayAllSui(
		receiverAddresss,
	)
	pt := ptb.Finish()

	tx := suiptb.NewTransactionData(
		senderAddress,
		pt,
		coins,
		suiclient.DefaultGasBudget,
		suiclient.DefaultGasPrice,
	)
	txBytesBCS, err := bcs.Marshal(tx)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}

	_, err = cli.SignAndSendTransaction(ctx, base64.StdEncoding.EncodeToString(txBytesBCS), signerInstance, client.WaitForLocalExecution)
	if err != nil {
		return fmt.Errorf("failed to sign and send transaction: %w", err)
	}

	return nil
}

// NewTestKeystore creates a new test keystore
func NewTestKeystore(t *testing.T) *TestKeystore {
	t.Helper()
	return &TestKeystore{t: t, keys: map[string]ed25519.PrivateKey{}}
}

// TestKeystore is a simple keystore for testing
type TestKeystore struct {
	t    *testing.T
	keys map[string]ed25519.PrivateKey
}

// AddKey adds a private key to the keystore
func (k *TestKeystore) AddKey(key ed25519.PrivateKey) {
	// Derive address from private key (in Sui, address is derived from public key)
	publicKey := key.Public().(ed25519.PublicKey)
	address := DeriveAddressFromPublicKey(publicKey)
	k.keys[address] = key
}

// Get returns a private key by address
func (k *TestKeystore) Get(address string) (ed25519.PrivateKey, bool) {
	key, ok := k.keys[address]
	return key, ok
}
