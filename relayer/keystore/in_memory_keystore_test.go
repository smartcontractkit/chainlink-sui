//go:build unit

package keystore_test

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"testing"

	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/keystore"
)

// Helper function to generate a base64-encoded key
func generateBase64Key(t *testing.T) string {
	t.Helper()
	_, privateKey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)
	return base64.StdEncoding.EncodeToString(privateKey.Seed())
}

// Helper function to generate a base64-encoded key with scheme identifier
func generateBase64KeyWithScheme(t *testing.T) string {
	t.Helper()
	_, privateKey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)
	// Add scheme identifier (0x00 for ed25519)
	keyWithScheme := append([]byte{0x00}, privateKey.Seed()...)
	return base64.StdEncoding.EncodeToString(keyWithScheme)
}

func TestNewInMemoryKeystore(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)

	ks := keystore.NewInMemoryKeystore(log)
	require.NotNil(t, ks)

	ctx := context.Background()
	accounts, err := ks.Accounts(ctx)
	require.NoError(t, err)
	assert.Empty(t, accounts)
}

func TestInMemoryKeystore_AddKey(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)
	ctx := context.Background()

	ks := keystore.NewInMemoryKeystore(log)

	// Generate a base64-encoded key
	base64Key := generateBase64Key(t)

	// Add the key
	address, err := ks.AddKey(base64Key)
	require.NoError(t, err)
	assert.NotEmpty(t, address)

	// Verify the address is in the accounts list
	accounts, err := ks.Accounts(ctx)
	require.NoError(t, err)
	assert.Len(t, accounts, 1)
	assert.Contains(t, accounts, address)
}

func TestInMemoryKeystore_AddKeyWithScheme(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)
	ctx := context.Background()

	ks := keystore.NewInMemoryKeystore(log)

	// Generate a base64-encoded key with scheme identifier
	base64Key := generateBase64KeyWithScheme(t)

	// Add the key
	address, err := ks.AddKey(base64Key)
	require.NoError(t, err)
	assert.NotEmpty(t, address)

	// Verify the address is in the accounts list
	accounts, err := ks.Accounts(ctx)
	require.NoError(t, err)
	assert.Len(t, accounts, 1)
	assert.Contains(t, accounts, address)
}

func TestInMemoryKeystore_AddKey_InvalidBase64(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)

	ks := keystore.NewInMemoryKeystore(log)

	// Try to add invalid base64
	_, err := ks.AddKey("invalid-base64!")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode base64 key")
}

func TestInMemoryKeystore_AddKey_InvalidKeyLength(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)

	ks := keystore.NewInMemoryKeystore(log)

	// Try to add key with invalid length
	shortKey := base64.StdEncoding.EncodeToString([]byte("short"))
	_, err := ks.AddKey(shortKey)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid key length")
}

func TestInMemoryKeystore_Sign(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)
	ctx := context.Background()

	ks := keystore.NewInMemoryKeystore(log)

	// Generate and add a base64-encoded key
	base64Key := generateBase64Key(t)
	address, err := ks.AddKey(base64Key)
	require.NoError(t, err)

	// Test signing with data
	testData := []byte("test message to sign")
	signature, err := ks.Sign(ctx, address, testData)
	require.NoError(t, err)
	assert.NotEmpty(t, signature)

	// Test signing with nil data (no-op to check account existence)
	signature, err = ks.Sign(ctx, address, nil)
	require.NoError(t, err)
	assert.Empty(t, signature)

	// Test signing with non-existent address
	_, err = ks.Sign(ctx, "nonexistent", testData)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no private key found")
}

func TestInMemoryKeystore_Accounts(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)
	ctx := context.Background()

	ks := keystore.NewInMemoryKeystore(log)

	// Initially should be empty
	accounts, err := ks.Accounts(ctx)
	require.NoError(t, err)
	assert.Empty(t, accounts)

	// Add multiple keys
	addresses := make([]string, 3)
	for i := 0; i < 3; i++ {
		base64Key := generateBase64Key(t)
		address, err := ks.AddKey(base64Key)
		require.NoError(t, err)
		addresses[i] = address
	}

	// Should have all addresses
	accounts, err = ks.Accounts(ctx)
	require.NoError(t, err)
	assert.Len(t, accounts, 3)

	// Check all addresses are present
	for _, addr := range addresses {
		assert.Contains(t, accounts, addr)
	}
}

func TestInMemoryKeystore_ConcurrentAccess(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)
	ctx := context.Background()

	ks := keystore.NewInMemoryKeystore(log)

	// Test concurrent access
	const numRoutines = 10
	const numKeysPerRoutine = 5

	done := make(chan bool, numRoutines)

	// Launch multiple goroutines adding keys
	for i := 0; i < numRoutines; i++ {
		go func() {
			defer func() { done <- true }()

			for j := 0; j < numKeysPerRoutine; j++ {
				base64Key := generateBase64Key(t)
				_, err := ks.AddKey(base64Key)
				require.NoError(t, err)
			}
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < numRoutines; i++ {
		<-done
	}

	// Verify all keys were added
	accounts, err := ks.Accounts(ctx)
	require.NoError(t, err)
	assert.Len(t, accounts, numRoutines*numKeysPerRoutine)
}

func TestInMemoryKeystore_AddressDerivation(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)

	ks := keystore.NewInMemoryKeystore(log)

	// Generate a test key
	_, privateKey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	// Encode as base64
	base64Key := base64.StdEncoding.EncodeToString(privateKey.Seed())

	// Add the key and get the derived address
	address, err := ks.AddKey(base64Key)
	require.NoError(t, err)

	// Verify the address matches what we'd expect from the signer
	expectedSigner := signer.NewSigner(privateKey.Seed())
	assert.Equal(t, expectedSigner.Address, address)

	// Test that we can sign with this address
	ctx := context.Background()
	testData := []byte("test")
	signature, err := ks.Sign(ctx, address, testData)
	require.NoError(t, err)
	assert.NotEmpty(t, signature)
}

func TestInMemoryKeystore_Interface(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)

	// Verify InMemoryKeystore implements Keystore interface
	var _ keystore.Keystore = keystore.NewInMemoryKeystore(log)
}
