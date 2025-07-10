//go:build integration

package tests

import (
	"context"
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	testpackage "github.com/smartcontractkit/chainlink-sui/bindings/packages/test"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
)

func TestComplexModule(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	signer, client := testenv.SetupEnvironment(t)
	opts := &bind.CallOpts{
		Signer:           signer,
		WaitForExecution: true,
	}

	testPackage, tx, err := testpackage.PublishTest(ctx, opts, client)
	require.NoError(t, err)
	require.NotNil(t, testPackage)
	require.NotNil(t, tx)

	contract := testPackage.Complex()

	t.Run("NewObject", func(t *testing.T) {
		someId := []byte{1, 2, 3, 4, 5}
		someNumber := uint64(42)
		someAddress := "0x1"
		someAddresses := []string{"0x2", "0x3"}

		obj, err := contract.DevInspect().NewObject(ctx, opts, someId, someNumber, someAddress, someAddresses)
		require.NoError(t, err)
		require.NotNil(t, obj)
		require.Equal(t, someId, obj.SomeId)
		require.Equal(t, someNumber, obj.SomeNumber)
		// addresses normalized to 66 chars
		require.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000001", obj.SomeAddress)
		expectedAddresses := []string{
			"0x0000000000000000000000000000000000000000000000000000000000000002",
			"0x0000000000000000000000000000000000000000000000000000000000000003",
		}
		require.Equal(t, expectedAddresses, obj.SomeAddresses)
	})

	t.Run("FlattenAddress", func(t *testing.T) {
		someAddress := "0x1"
		someAddresses := []string{"0x2", "0x3", "0x4"}

		addresses, err := contract.DevInspect().FlattenAddress(ctx, opts, someAddress, someAddresses)
		require.NoError(t, err)
		require.Len(t, addresses, 4)

		for i, addr := range addresses {
			require.Len(t, addr, 66, "Address %d should be normalized", i)
			require.True(t, strings.HasPrefix(addr, "0x"))
		}
	})

	t.Run("FlattenU8", func(t *testing.T) {
		input := [][]byte{
			{1, 2, 3},
			{4, 5},
			{6, 7, 8, 9},
		}

		bytes, err := contract.DevInspect().FlattenU8(ctx, opts, input)
		require.NoError(t, err)

		expected := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}
		require.Equal(t, expected, bytes)
	})

	t.Run("CheckU128", func(t *testing.T) {
		testCases := []struct {
			name  string
			value string
		}{
			{"max u128", "340282366920938463463374607431768211455"},
			{"zero", "0"},
			{"mid value", "170141183460469231731687303715884105727"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				input := new(big.Int)
				input.SetString(tc.value, 10)

				output, err := contract.DevInspect().CheckU128(ctx, opts, input)
				require.NoError(t, err)
				require.Equal(t, 0, input.Cmp(output), "u128 round-trip failed for %s", tc.name)
			})
		}
	})

	t.Run("CheckU256", func(t *testing.T) {
		testCases := []struct {
			name  string
			value string
		}{
			{"max u64", "18446744073709551615"},
			{"large value", "0xFFFFFFFFFFFFFFFF"},
			{"max u256", "0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				input := new(big.Int)
				_, ok := input.SetString(tc.value, 0)
				require.True(t, ok, "Failed to parse %s", tc.value)

				output, err := contract.DevInspect().CheckU256(ctx, opts, input)
				require.NoError(t, err)
				require.Equal(t, 0, input.Cmp(output), "u256 round-trip failed for %s", tc.name)
			})
		}
	})

	t.Run("SharedObjectFunctions", func(t *testing.T) {
		someId := []byte("test-shared-object")
		someNumber := uint64(999)
		someAddress := "0xabc"
		someAddresses := []string{"0xdef", "0x123", "0x456"}

		createTx, err := contract.NewObjectWithTransfer(ctx, opts, someId, someNumber, someAddress, someAddresses)
		require.NoError(t, err)
		require.Equal(t, "success", createTx.Effects.Status.Status)

		objId, sharedVersion, err := FindCreatedObject(createTx.ObjectChanges, "::complex::SampleObject")
		require.NoError(t, err)
		sampleObj := bind.Object{
			Id:                   objId,
			InitialSharedVersion: sharedVersion,
		}
		require.NotEmpty(t, sampleObj.Id, "SampleObject not found")
		require.NotNil(t, sampleObj.InitialSharedVersion, "SampleObject should be shared")

		t.Run("CheckWithObjectRef", func(t *testing.T) {
			result, err := contract.DevInspect().CheckWithObjectRef(ctx, opts, sampleObj)
			require.NoError(t, err)
			require.Equal(t, someNumber, result)
		})

		t.Run("CheckWithMutObjectRef", func(t *testing.T) {
			newNumber := uint64(1234)

			result, err := contract.DevInspect().CheckWithMutObjectRef(ctx, opts, sampleObj, newNumber)
			require.NoError(t, err)
			require.Equal(t, newNumber, result)

			checkResult, err := contract.DevInspect().CheckWithObjectRef(ctx, opts, sampleObj)
			require.NoError(t, err)
			require.Equal(t, someNumber, checkResult, "State should not be modified by DevInspect")
		})
	})
}
