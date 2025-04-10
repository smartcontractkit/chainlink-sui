//go:build integration

package modulecomplex

import (
	"context"
	"fmt"
	"testing"

	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	rel "github.com/smartcontractkit/chainlink-sui/relayer/signer"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"

	"github.com/stretchr/testify/require"
)

func setupSuiTest(t *testing.T) (rel.SuiSigner, *suiclient.ClientImpl) {
	t.Helper()
	log := logger.Test(t)

	// Start the node and schedule cleanup.
	cmd, err := testutils.StartSuiNode(testutils.CLI)
	require.NoError(t, err)
	t.Cleanup(func() {
		if cmd.Process != nil {
			if err = cmd.Process.Kill(); err != nil {
				t.Logf("Failed to kill process: %v", err)
			}
		}
	})

	// Generate key pair and create signer.
	pk, _, _, err := testutils.GenerateAccountKeyPair(t, log)
	require.NoError(t, err)
	signer := rel.NewPrivateKeySigner(pk)

	// Create client.
	client := suiclient.NewClient("http://localhost:9000")

	// Fund account.
	signerAddress, err := signer.GetAddress()
	require.NoError(t, err)
	err = testutils.FundWithFaucet(log, "localnet", signerAddress)
	require.NoError(t, err)

	return signer, client
}

// helper to convert hex string to address string for readability.
func getAddress(t *testing.T, hexStr string) string {
	t.Helper()
	addr, err := sui.AddressFromHex(hexStr)
	require.NoError(t, err)

	return addr.String()
}

// nolint:tparallel
func TestComplex(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	signer, client := setupSuiTest(t)

	// Publish the complex contract.
	contract, tx, err := PublishComplex(ctx, bind.TxOpts{}, signer, *client)
	require.NoError(t, err)
	require.NotNil(t, contract)
	require.NotNil(t, tx)

	// Prepare addresses.
	addresses := []string{
		getAddress(t, "0x11234"),
		getAddress(t, "0x21234"),
		getAddress(t, "0x31234"),
		getAddress(t, "0x41234"),
	}

	// nolint:paralleltest
	t.Run("NewObject", func(t *testing.T) {
		id := []byte("0x1234")
		someNumber := uint64(1)
		newObjectTx, err := contract.NewObject(id, someNumber, addresses[3], addresses[:3]).Execute(ctx, bind.TxOpts{}, signer, *client)
		require.NoError(t, err)
		require.NotNil(t, newObjectTx)
	})

	// nolint:paralleltest
	t.Run("FlattenAddress", func(t *testing.T) {
		flattenTx, err := contract.FlattenAddress(addresses[0], addresses[:3]).Execute(ctx, bind.TxOpts{}, signer, *client)
		require.NoError(t, err)
		require.NotNil(t, flattenTx)

		inspection, err := contract.FlattenAddress(addresses[0], addresses[:3]).Inspect(ctx, bind.TxOpts{}, signer, *client)
		require.NoError(t, err)
		require.NotNil(t, inspection)

		bytes, err := getBytesFromResult(t, inspection)
		require.NoError(t, err)
		expectedBytes := []byte{4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 18, 52, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 18, 52, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 18, 52, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 18, 52}
		require.Equal(t, expectedBytes, bytes)
	})

	// nolint:paralleltest
	t.Run("FlattenU8", func(t *testing.T) {
		input := [][]uint8{{0, 1, 2}, {3, 4, 5}}
		flattenTx, err := contract.FlattenU8(input).Execute(ctx, bind.TxOpts{}, signer, *client)
		require.NoError(t, err)
		require.NotNil(t, flattenTx)

		inspection, err := contract.FlattenU8(input).Inspect(ctx, bind.TxOpts{}, signer, *client)
		require.NoError(t, err)
		require.NotNil(t, inspection)

		bytes, err := getBytesFromResult(t, inspection)
		require.NoError(t, err)
		expectedBytes := []byte{6, 0, 1, 2, 3, 4, 5}
		require.Equal(t, expectedBytes, bytes)
	})
}

func getBytesFromResult(t *testing.T, res *suiclient.DevInspectTransactionBlockResponse) ([]byte, error) {
	t.Helper()
	r := res.Results[0].ReturnValues[0]
	resultSlice, ok := r.([]any)
	require.True(t, ok)
	require.NotEmpty(t, resultSlice, "Result slice is empty")
	b := resultSlice[0]
	arr, ok := b.([]any)
	if !ok {
		return nil, fmt.Errorf("expected b to be a slice, got %T", b)
	}
	resultBytes := make([]byte, 0, len(arr))
	for _, item := range arr {
		// Assuming each element is numeric (float64) and represents a byte value.
		num, ok := item.(float64)
		if !ok {
			return nil, fmt.Errorf("unexpected element type %T in b", item)
		}
		resultBytes = append(resultBytes, byte(num))
	}

	return resultBytes, nil
}
