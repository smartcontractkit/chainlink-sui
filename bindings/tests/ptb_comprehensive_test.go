//go:build integration

package tests

import (
	"context"
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

// TestProgrammableTransactionBlocks tests all PTB functionality including:
// - Basic PTB construction with multiple operations
// - PTB with view functions (DevInspect within PTB)
// - PTB with non-entry functions
// - PTB return value handling (Execute vs DevInspect)
// - Result chaining with WithArgument methods
func TestProgrammableTransactionBlocks(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	signer, client := testenv.SetupEnvironment(t)

	opts := &bind.CallOpts{
		Signer:           signer,
		WaitForExecution: true,
	}

	// Publish the test package
	testPackage, tx, err := testpackage.PublishTest(ctx, opts, client)
	require.NoError(t, err)
	require.NotNil(t, testPackage)
	require.NotNil(t, tx)

	t.Run("Basic PTB Operations", func(t *testing.T) {
		counterInterface := testPackage.Counter()
		counter, ok := counterInterface.(*module_counter.CounterContract)
		require.True(t, ok, "Failed to cast to CounterContract")

		// Initialize counter
		initTx, err := counter.Initialize(ctx, opts)
		require.NoError(t, err)
		require.Equal(t, "success", initTx.Effects.Status.Status)

		// Find the counter object
		objId, sharedVersion, err := FindCreatedObject(initTx.ObjectChanges, "counter::Counter")
		require.NoError(t, err)
		require.NotNil(t, sharedVersion, "Counter should be shared")

		counterObj := bind.Object{
			Id:                   objId,
			InitialSharedVersion: sharedVersion,
		}

		// Build a PTB with multiple operations
		ptb := transaction.NewTransaction()

		signerAddr, err := signer.GetAddress()
		require.NoError(t, err)
		ptb.SetSender(models.SuiAddress(signerAddr))
		ptb.SetGasBudget(30_000_000)

		encoder := counter.Encoder()

		// Add multiple operations to PTB
		encoded1, err := encoder.IncrementByOne(counterObj)
		require.NoError(t, err)
		_, err = counter.AppendPTB(ctx, opts, ptb, encoded1)
		require.NoError(t, err)

		encoded2, err := encoder.IncrementBy(counterObj, 5)
		require.NoError(t, err)
		_, err = counter.AppendPTB(ctx, opts, ptb, encoded2)
		require.NoError(t, err)

		encoded3, err := encoder.IncrementMult(counterObj, 2, 3)
		require.NoError(t, err)
		_, err = counter.AppendPTB(ctx, opts, ptb, encoded3)
		require.NoError(t, err)

		// Execute PTB (total: 1 + 5 + 6 = 12)
		ptbTx, err := bind.ExecutePTB(ctx, opts, client, ptb)
		require.NoError(t, err)
		require.Equal(t, "success", ptbTx.Effects.Status.Status)

		// Verify final count
		count, err := counter.DevInspect().GetCount(ctx, opts, counterObj)
		require.NoError(t, err)
		require.Equal(t, uint64(12), count)
	})

	t.Run("PTB Function Calls", func(t *testing.T) {
		counterInterface := testPackage.Counter()
		counter, ok := counterInterface.(*module_counter.CounterContract)
		require.True(t, ok, "Failed to cast to CounterContract")

		ptb := transaction.NewTransaction()
		signerAddr, err := signer.GetAddress()
		require.NoError(t, err)
		ptb.SetSender(models.SuiAddress(signerAddr))
		ptb.SetGasBudget(30_000_000)

		encoder := counter.Encoder()

		encodedCreate, err := encoder.Create()
		require.NoError(t, err)

		counterResult, err := counter.AppendPTB(ctx, opts, ptb, encodedCreate)
		require.NoError(t, err)
		require.NotNil(t, counterResult)

		encodedGetCount, err := encoder.GetCountWithArgs(*counterResult)
		require.NoError(t, err)

		countResult, err := counter.AppendPTB(ctx, opts, ptb, encodedGetCount)
		require.NoError(t, err)
		require.NotNil(t, countResult)

		// Transfer the created counter to complete the PTB
		ptb.TransferObjects(
			[]transaction.Argument{*counterResult},
			ptb.Pure(signerAddr),
		)

		// Execute PTB
		ptbTx, err := bind.ExecutePTB(ctx, opts, client, ptb)
		require.NoError(t, err)
		require.Equal(t, "success", ptbTx.Effects.Status.Status)
	})

	t.Run("PTB Return Values via DevInspect", func(t *testing.T) {
		counterInterface := testPackage.Counter()
		counter, ok := counterInterface.(*module_counter.CounterContract)
		require.True(t, ok, "Failed to cast to CounterContract")

		// Create and initialize a counter
		initTx, err := counter.Initialize(ctx, opts)
		require.NoError(t, err)

		objId, sharedVersion, err := FindCreatedObject(initTx.ObjectChanges, "counter::Counter")
		require.NoError(t, err)
		require.NotNil(t, sharedVersion, "Counter should be shared")

		counterObj := bind.Object{
			Id:                   objId,
			InitialSharedVersion: sharedVersion,
		}

		// Execute mode - no return values
		executeTx, err := counter.IncrementByOne(ctx, opts, counterObj)
		require.NoError(t, err)
		require.Equal(t, "success", executeTx.Effects.Status.Status)
		// Note: executeTx has no return values from the Move function

		// DevInspect mode - returns values
		newCount, err := counter.DevInspect().IncrementByOne(ctx, opts, counterObj)
		require.NoError(t, err)
		require.Equal(t, uint64(2), newCount) // 1 + 1 = 2
	})

	t.Run("Complex Types in PTB", func(t *testing.T) {
		complexInterface := testPackage.Complex()
		complexContract, ok := complexInterface.(*module_complex.ComplexContract)
		require.True(t, ok, "Failed to cast to ComplexContract")

		// Create complex object
		createTx, err := complexContract.NewObjectWithTransfer(ctx, opts, []byte("test-id"), uint64(100), "0x0000000000000000000000000000000000000000000000000000000000000001", []string{})
		require.NoError(t, err)
		require.Equal(t, "success", createTx.Effects.Status.Status)

		// Find the created object
		complexObjectId, initialSharedVersion, err := FindCreatedObject(createTx.ObjectChanges, "complex::SampleObject")
		require.NoError(t, err)

		// Build PTB with complex operations
		ptb := transaction.NewTransaction()
		signerAddr, err := signer.GetAddress()
		require.NoError(t, err)
		ptb.SetSender(models.SuiAddress(signerAddr))
		ptb.SetGasBudget(30_000_000)

		complexObj := bind.Object{
			Id:                   complexObjectId,
			InitialSharedVersion: initialSharedVersion,
		}
		encoder := complexContract.Encoder()

		// Check with object ref
		encoded1, err := encoder.CheckWithObjectRef(complexObj)
		require.NoError(t, err)
		_, err = complexContract.AppendPTB(ctx, opts, ptb, encoded1)
		require.NoError(t, err)

		// Check with mut object ref
		encoded2, err := encoder.CheckWithMutObjectRef(complexObj, uint64(200))
		require.NoError(t, err)
		_, err = complexContract.AppendPTB(ctx, opts, ptb, encoded2)
		require.NoError(t, err)

		// No need to transfer - object is shared

		// Execute PTB
		ptbTx, err := bind.ExecutePTB(ctx, opts, client, ptb)
		require.NoError(t, err)
		require.Equal(t, "success", ptbTx.Effects.Status.Status)
	})

	t.Run("PTB DevInspect for Multiple Return Values", func(t *testing.T) {
		counterInterface := testPackage.Counter()
		counter, ok := counterInterface.(*module_counter.CounterContract)
		require.True(t, ok, "Failed to cast to CounterContract")

		results, err := counter.DevInspect().GetTupleStruct(ctx, opts)
		require.NoError(t, err)
		require.Len(t, results, 4)

		_, ok = results[0].(uint64)
		require.True(t, ok, "First tuple element should be uint64")
		_, ok = results[1].(string)
		require.True(t, ok, "Second tuple element should be address")
		_, ok = results[2].(bool)
		require.True(t, ok, "Third tuple element should be bool")

		result, err := counter.DevInspect().GetSimpleResult(ctx, opts)
		require.NoError(t, err)
		require.Equal(t, uint64(42), result.Value)
	})
}
