//go:build integration

package tests

import (
	"context"
	"math/big"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_complex "github.com/smartcontractkit/chainlink-sui/bindings/generated/test/complex"
	module_counter "github.com/smartcontractkit/chainlink-sui/bindings/generated/test/counter"
	testpackage "github.com/smartcontractkit/chainlink-sui/bindings/packages/test"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
)

func TestBindingFeatures(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	signer, client := testenv.SetupEnvironment(t)

	opts := &bind.CallOpts{
		Signer:           signer,
		WaitForExecution: true,
		GasBudget:        &DEFAULT_GAS_BUDGET,
	}

	testPackage, tx, err := testpackage.PublishTest(ctx, opts, client)
	require.NoError(t, err)
	require.NotNil(t, testPackage)
	require.NotNil(t, tx)

	counterContract, ok := testPackage.Counter().(*module_counter.CounterContract)
	require.True(t, ok)
	complexContract, ok := testPackage.Complex().(*module_complex.ComplexContract)
	require.True(t, ok)

	t.Run("DevInspect return values", func(t *testing.T) {
		initTx, err := counterContract.Initialize(ctx, opts)
		require.NoError(t, err)

		objId, sharedVersion, err := FindCreatedObject(initTx.ObjectChanges, "::counter::Counter")
		require.NoError(t, err)
		counterObj := bind.Object{
			Id:                   objId,
			InitialSharedVersion: sharedVersion,
		}

		_, err = counterContract.IncrementBy(ctx, opts, counterObj, 25)
		require.NoError(t, err)

		t.Run("Typed single return value", func(t *testing.T) {
			count, err := counterContract.DevInspect().GetCount(ctx, opts, counterObj)
			require.NoError(t, err)
			require.Equal(t, uint64(25), count)

			bytes, err := counterContract.DevInspect().GetVectorOfU8(ctx, opts)
			require.NoError(t, err)
			require.Equal(t, []byte{1, 2, 3, 4}, bytes)

			addresses, err := counterContract.DevInspect().GetVectorOfAddresses(ctx, opts)
			require.NoError(t, err)
			require.Len(t, addresses, 4)

			expectedAddresses := []string{
				"0x0000000000000000000000000000000000000000000000000000000000000001",
				"0x0000000000000000000000000000000000000000000000000000000000000002",
				"0x0000000000000000000000000000000000000000000000000000000000000003",
				"0x0000000000000000000000000000000000000000000000000000000000000004",
			}
			for i, addr := range addresses {
				require.Equal(t, expectedAddresses[i], addr, "Address at index %d should match", i)
			}

			vectors, err := counterContract.DevInspect().GetVectorOfVectorsOfU8(ctx, opts)
			require.NoError(t, err)
			require.Len(t, vectors, 4)
			for i, vec := range vectors {
				require.Len(t, vec, 32, "Vector %d should be 32 bytes (address)", i)
			}
		})

		t.Run("Struct types", func(t *testing.T) {
			simpleResult, err := counterContract.DevInspect().GetSimpleResult(ctx, opts)
			require.NoError(t, err)
			require.Equal(t, uint64(42), simpleResult.Value)

			addressList, err := counterContract.DevInspect().GetAddressList(ctx, opts)
			require.NoError(t, err)
			require.Equal(t, uint64(4), addressList.Count)
			require.Len(t, addressList.Addresses, 4)

			nestedStruct, err := counterContract.DevInspect().GetNestedResultStruct(ctx, opts)
			require.NoError(t, err)
			require.True(t, nestedStruct.IsNested)
			require.Equal(t, uint64(42), nestedStruct.DoubleCount)
			require.NotEmpty(t, nestedStruct.NestedStruct.Addr)
		})

		t.Run("Different typed tuple", func(t *testing.T) {
			results, err := counterContract.DevInspect().GetTupleStruct(ctx, opts)
			require.NoError(t, err)
			require.Len(t, results, 4)

			val1, ok := results[0].(uint64)
			require.True(t, ok)
			require.Equal(t, uint64(42), val1)

			val2, ok := results[1].(string)
			require.True(t, ok)
			require.NotEmpty(t, val2)

			val3, ok := results[2].(bool)
			require.True(t, ok)
			require.True(t, val3)

			val4, ok := results[3].(module_counter.MultiNestedStruct)
			require.True(t, ok)
			require.True(t, val4.IsMultiNested)
		})

		t.Run("DevInspectInterface", func(t *testing.T) {
			devInspect := counterContract.DevInspect()

			count, err := devInspect.GetCount(ctx, opts, counterObj)
			require.NoError(t, err)
			require.GreaterOrEqual(t, count, uint64(0))

			result, err := devInspect.GetSimpleResult(ctx, opts)
			require.NoError(t, err)
			require.NotNil(t, result)
			require.Equal(t, uint64(42), result.Value)
		})
	})

	t.Run("PTB behavior", func(t *testing.T) {
		initTx, err := counterContract.Initialize(ctx, opts)
		require.NoError(t, err)

		objId, sharedVersion, err := FindCreatedObject(initTx.ObjectChanges, "::counter::Counter")
		require.NoError(t, err)
		counterObj := bind.Object{Id: objId, InitialSharedVersion: sharedVersion}

		t.Run("PTB execution", func(t *testing.T) {
			ptb := transaction.NewTransaction()
			encoder := counterContract.Encoder()

			encoded1, err := encoder.IncrementBy(counterObj, 10)
			require.NoError(t, err)
			_, err = counterContract.AppendPTB(ctx, opts, ptb, encoded1)
			require.NoError(t, err)

			encoded2, err := encoder.IncrementBy(counterObj, 20)
			require.NoError(t, err)
			_, err = counterContract.AppendPTB(ctx, opts, ptb, encoded2)
			require.NoError(t, err)

			signerAddr, err := signer.GetAddress()
			require.NoError(t, err)
			ptb.SetSender(models.SuiAddress(signerAddr))
			ptb.SetGasBudget(30_000_000)

			ptbTx, err := bind.ExecutePTB(ctx, opts, client, ptb)
			require.NoError(t, err)
			require.Equal(t, "success", ptbTx.Effects.Status.Status)

			count, err := counterContract.DevInspect().GetCount(ctx, opts, counterObj)
			require.NoError(t, err)
			require.Equal(t, uint64(30), count)
		})
	})

	t.Run("Type Conversions", func(t *testing.T) {
		t.Run("u128", func(t *testing.T) {
			u128Input := new(big.Int)
			u128Input.SetString("340282366920938463463374607431768211455", 10) // max u128

			u128Result, err := complexContract.DevInspect().CheckU128(ctx, opts, u128Input)
			require.NoError(t, err)
			require.Equal(t, 0, u128Input.Cmp(u128Result))
		})

		t.Run("u256", func(t *testing.T) {
			u256Input, _ := new(big.Int).SetString("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", 0)
			require.NotNil(t, u256Input)

			u256Result, err := complexContract.DevInspect().CheckU256(ctx, opts, u256Input)
			require.NoError(t, err)
			require.Equal(t, u256Input.String(), u256Result.String())
		})

		t.Run("Address Normalization", func(t *testing.T) {
			// addresses normalized to 66 chars
			shortAddress := "0x1"
			addresses := []string{"0x2", "0x3"}

			resultAddresses, err := complexContract.DevInspect().FlattenAddress(ctx, opts, shortAddress, addresses)
			require.NoError(t, err)
			require.Len(t, resultAddresses, 3)

			for _, addr := range resultAddresses {
				require.Len(t, addr, 66)
				require.Equal(t, "0x", addr[:2])
			}
		})

		t.Run("Vector encoding and decoding", func(t *testing.T) {
			input := [][]byte{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}

			flattened, err := complexContract.DevInspect().FlattenU8(ctx, opts, input)
			require.NoError(t, err)
			require.Equal(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, flattened)
		})
	})

	t.Run("Invalid object error handling", func(t *testing.T) {
		t.Run("InvalidObjectID", func(t *testing.T) {
			invalidObj := bind.Object{
				Id: "0xinvalid",
			}

			_, err := counterContract.Increment(ctx, opts, invalidObj)
			require.Error(t, err)
		})

		t.Run("MissingSharedVersion", func(t *testing.T) {
			invalidObj := bind.Object{
				Id: "0x0000000000000000000000000000000000000000000000000000000000000123",
			}

			_, err := counterContract.GetCount(ctx, opts, invalidObj)
			require.Error(t, err)
		})
	})

	t.Run("Nested struct decoding", func(t *testing.T) {
		droppableObj, err := complexContract.DevInspect().NewObject(ctx, opts,
			[]byte{1, 2, 3},
			42,
			"0x1",
			[]string{"0x2", "0x3"})
		require.NoError(t, err)
		require.Equal(t, []byte{1, 2, 3}, droppableObj.SomeId)
		require.Equal(t, uint64(42), droppableObj.SomeNumber)
		require.Len(t, droppableObj.SomeAddress, 66) // Normalized
		require.Len(t, droppableObj.SomeAddresses, 2)

		ocrConfig, err := counterContract.DevInspect().GetOcrConfig(ctx, opts)
		require.NoError(t, err)
		require.NotEmpty(t, ocrConfig.ConfigInfo.ConfigDigest)
		require.True(t, ocrConfig.ConfigInfo.IsSignatureVerificationEnabled)
		require.Len(t, ocrConfig.Transmitters, 1)
		require.Len(t, ocrConfig.Signers, 2)

		for i, transmitter := range ocrConfig.Transmitters {
			require.NotEmpty(t, transmitter)
			require.Len(t, transmitter, 66)
			t.Logf("  Transmitter %d: %s", i, transmitter)
		}

		for i, signer := range ocrConfig.Signers {
			require.NotEmpty(t, signer)
			t.Logf("  Signer %d: %x", i, signer)
		}
	})

	t.Run("DevInspect simulates modified state", func(t *testing.T) {
		initTx, err := counterContract.Initialize(ctx, opts)
		require.NoError(t, err)

		objId, sharedVersion, err := FindCreatedObject(initTx.ObjectChanges, "::counter::Counter")
		require.NoError(t, err)
		require.NotNil(t, sharedVersion)
		counterObj := bind.Object{
			Id:                   objId,
			InitialSharedVersion: sharedVersion,
		}

		initialCount, err := counterContract.DevInspect().GetCount(ctx, opts, counterObj)
		require.NoError(t, err)

		// Execute modifies state, no return values
		executeTx, err := counterContract.IncrementByOne(ctx, opts, counterObj)
		require.NoError(t, err)
		require.Equal(t, "success", executeTx.Effects.Status.Status)

		newValue, err := counterContract.DevInspect().IncrementByOne(ctx, opts, counterObj)
		require.NoError(t, err)
		require.Equal(t, initialCount+2, newValue)

		count, err := counterContract.DevInspect().GetCount(ctx, opts, counterObj)
		require.NoError(t, err)
		// Only incremented once by Execute
		require.Equal(t, initialCount+1, count)
	})
}
