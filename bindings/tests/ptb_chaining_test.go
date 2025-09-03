//go:build integration

package tests

import (
	"context"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_counter "github.com/smartcontractkit/chainlink-sui/bindings/generated/test/counter"
	module_generics "github.com/smartcontractkit/chainlink-sui/bindings/generated/test/generics"
	testpackage "github.com/smartcontractkit/chainlink-sui/bindings/packages/test"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
)

func TestPTBChaining(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	signer, client := testenv.SetupEnvironment(t)

	testPackage, tx, err := testpackage.PublishTest(ctx, &bind.CallOpts{
		Signer:           signer,
		WaitForExecution: true,
		GasBudget:        &DEFAULT_GAS_BUDGET,
	}, client)
	require.NoError(t, err)
	require.NotNil(t, testPackage)
	require.NotNil(t, tx)

	opts := &bind.CallOpts{
		Signer:           signer,
		WaitForExecution: true,
		GasBudget:        &DEFAULT_GAS_BUDGET,
	}

	t.Run("Counter PTB chaining", func(t *testing.T) {
		contract := testPackage.Counter()
		counterContract := contract.(*module_counter.CounterContract)

		initTx, err := contract.Initialize(ctx, opts)
		require.NoError(t, err)
		require.NotNil(t, initTx)

		objId, sharedVersion, err := FindCreatedObject(initTx.ObjectChanges, "counter::Counter")
		require.NoError(t, err)
		require.NotNil(t, sharedVersion, "Counter should be shared")

		counterObj := bind.Object{
			Id:                   objId,
			InitialSharedVersion: sharedVersion,
		}

		ptb := transaction.NewTransaction()

		signerAddr, err := signer.GetAddress()
		require.NoError(t, err)
		ptb.SetSender(models.SuiAddress(signerAddr))
		ptb.SetGasBudget(30_000_000)

		incrementBy10, err := counterContract.Encoder().IncrementBy(counterObj, 10)
		require.NoError(t, err)
		_, err = counterContract.AppendPTB(ctx, opts, ptb, incrementBy10)
		require.NoError(t, err)

		incrementMult, err := counterContract.Encoder().IncrementMult(counterObj, 5, 6)
		require.NoError(t, err)
		_, err = counterContract.AppendPTB(ctx, opts, ptb, incrementMult)
		require.NoError(t, err)

		ptbTx, err := bind.ExecutePTB(ctx, opts, client, ptb)
		require.NoError(t, err)
		require.NotNil(t, ptbTx)
		require.Equal(t, "success", ptbTx.Effects.Status.Status)

		count, err := contract.DevInspect().GetCount(ctx, opts, counterObj)
		require.NoError(t, err)
		require.Equal(t, uint64(40), count)
	})

	t.Run("Generics PTB chaining with objects", func(t *testing.T) {
		genericsInterface := testPackage.Generics()
		contract, ok := genericsInterface.(*module_generics.GenericsContract)
		require.True(t, ok, "Failed to cast to GenericsContract")

		signerAddr, err := signer.GetAddress()
		require.NoError(t, err)

		ptb := transaction.NewTransaction()
		ptb.SetSender(models.SuiAddress(signerAddr))
		ptb.SetGasBudget(30_000_000)

		createToken1, err := contract.Encoder().CreateSuiToken()
		require.NoError(t, err)
		token1Result, err := contract.AppendPTB(ctx, opts, ptb, createToken1)
		require.NoError(t, err)

		createToken2, err := contract.Encoder().CreateSuiToken()
		require.NoError(t, err)
		token2Result, err := contract.AppendPTB(ctx, opts, ptb, createToken2)
		require.NoError(t, err)

		depositEncoded, err := contract.Encoder().DepositWithArgs(
			[]string{"0x2::sui::SUI"},
			*token1Result,
			*token2Result,
		)
		require.NoError(t, err)
		_, err = contract.AppendPTB(ctx, opts, ptb, depositEncoded)
		require.NoError(t, err)

		ptb.TransferObjects(
			[]transaction.Argument{*token1Result},
			ptb.Pure(signerAddr),
		)

		tx, err := bind.ExecutePTB(ctx, opts, client, ptb)
		require.NoError(t, err)
		require.Equal(t, "success", tx.Effects.Status.Status)
	})

	t.Run("Verify WithArgument methods exist", func(t *testing.T) {
		counterContract := testPackage.Counter().(*module_counter.CounterContract)
		encoder := counterContract.Encoder()

		result := uint16(0)
		dummyArg := transaction.Argument{Result: &result}

		encoded, err := encoder.IncrementByOneWithArgs(dummyArg)
		require.NoError(t, err)
		require.NotNil(t, encoded)
		require.Equal(t, "increment_by_one", encoded.Function)
		require.Len(t, encoded.CallArgs, 1)
		require.True(t, encoded.CallArgs[0].IsArgument())
	})
}
