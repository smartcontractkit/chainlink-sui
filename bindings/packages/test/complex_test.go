//go:build integration

package test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/holiman/uint256"
	"github.com/pattonkan/sui-go/sui"
	"github.com/pattonkan/sui-go/suiclient"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"

	"github.com/stretchr/testify/require"
)

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
	testPackage, tx, err := PublishTest(ctx, bind.TxOpts{}, signer, *client)
	require.NoError(t, err)
	require.NotNil(t, testPackage)
	require.NotNil(t, tx)

	contract := testPackage.Complex()
	// Prepare addresses.
	addresses := []string{
		getAddress(t, "0x11234"),
		getAddress(t, "0x21234"),
		getAddress(t, "0x31234"),
		getAddress(t, "0x41234"),
	}

	//nolint:paralleltest
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

	//nolint:paralleltest
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

	//nolint:paralleltest
	t.Run("Check_u128", func(t *testing.T) {
		input := big.NewInt(2000)
		flattenTx, err := contract.CheckU128(input).Execute(ctx, bind.TxOpts{}, signer, *client)
		require.NoError(t, err)
		require.NotNil(t, flattenTx)

		inspection, err := contract.CheckU128(input).Inspect(ctx, bind.TxOpts{}, signer, *client)
		require.NoError(t, err)
		require.NotNil(t, inspection)
	})

	//nolint:paralleltest
	t.Run("Check_u256", func(t *testing.T) {
		input := uint256.NewInt(2000)
		flattenTx, err := contract.CheckU256(*input).Execute(ctx, bind.TxOpts{}, signer, *client)
		require.NoError(t, err)
		require.NotNil(t, flattenTx)

		inspection, err := contract.CheckU256(*input).Inspect(ctx, bind.TxOpts{}, signer, *client)
		require.NoError(t, err)
		require.NotNil(t, inspection)
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
