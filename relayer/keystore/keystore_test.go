//go:build unit

package keystore_test

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/keystore"
)

func TestNewSuiKeystore(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)

	tests := []struct {
		name        string
		keyPath     string
		expectedErr bool
	}{
		{
			name:        "Default parameters",
			keyPath:     "",
			expectedErr: false,
		},
		{
			name:        "Custom keystore path",
			keyPath:     "/tmp/test-keystore",
			expectedErr: false,
		},
		{
			name:        "Custom signer type",
			keyPath:     "",
			expectedErr: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ks, err := keystore.NewSuiKeystore(log, tc.keyPath)

			if tc.expectedErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, ks)

				// Validate default values are set
				if tc.keyPath == "" {
					assert.Equal(t, keystore.SuiDefaultKeystorePath, ks.KeyStorePath())
				} else {
					assert.Equal(t, tc.keyPath, ks.KeyStorePath())
				}
			}
		})
	}
}

func TestGetSignerFromAddress(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)

	// Create a temporary keystore file
	tempDir := t.TempDir()
	tempKeystorePath := filepath.Join(tempDir, "sui.keystore")

	// Generate a key pair
	_, privKey, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)

	// Create a signer and get the address
	suiSigner := signer.NewSigner(privKey.Seed())
	address := suiSigner.Address

	// Store the key in keystore format (base64 encoded)
	// First byte is the key scheme flag (0x00 for ed25519)
	keystoreKey := append([]byte{0x00}, privKey.Seed()...)
	encodedKey := base64.StdEncoding.EncodeToString(keystoreKey)
	keystoreKeys := []string{encodedKey}

	// Write to the temporary keystore file
	keystoreData, err := json.Marshal(keystoreKeys)
	require.NoError(t, err)
	err = os.WriteFile(tempKeystorePath, keystoreData, 0600)
	require.NoError(t, err)

	tests := []struct {
		name        string
		address     string
		expectedErr bool
	}{
		{
			name:        "Valid address",
			address:     address,
			expectedErr: false,
		},
		{
			name:        "Invalid address",
			address:     "0x123456",
			expectedErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ks, err := keystore.NewSuiKeystore(log, tempKeystorePath)
			require.NoError(t, err)

			privateKey, err := ks.GetPrivateKeyByAddress(tc.address)

			if tc.expectedErr {
				require.Error(t, err)
				assert.Nil(t, privateKey)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, privateKey)
			}
		})
	}
}

func TestNonexistentKeystorePath(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)

	nonexistentPath := "/nonexistent/path/to/keystore"
	ks, err := keystore.NewSuiKeystore(log, nonexistentPath)
	require.NoError(t, err)

	_, err = ks.GetPrivateKeyByAddress("0x123456")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read keystore file")
}

func TestInvalidKeystoreFormat(t *testing.T) {
	t.Parallel()
	log := logger.Test(t)

	// Create a temporary keystore file with invalid JSON
	tempDir := t.TempDir()
	tempKeystorePath := filepath.Join(tempDir, "invalid.keystore")

	err := os.WriteFile(tempKeystorePath, []byte("invalid json"), 0600)
	require.NoError(t, err)

	ks, err := keystore.NewSuiKeystore(log, tempKeystorePath)
	require.NoError(t, err)

	_, err = ks.GetPrivateKeyByAddress("0x123456")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse keystore JSON")
}
