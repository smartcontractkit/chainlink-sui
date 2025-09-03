//go:build integration

package tests

import (
	"context"
	"testing"

	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_generics "github.com/smartcontractkit/chainlink-sui/bindings/generated/test/generics"
	testpackage "github.com/smartcontractkit/chainlink-sui/bindings/packages/test"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
)

func TestGenerics(t *testing.T) {
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

	genericsInterface := testPackage.Generics()
	genericsContract, ok := genericsInterface.(*module_generics.GenericsContract)
	require.True(t, ok, "Failed to cast to GenericsContract")

	t.Run("Token<SUI> creation", func(t *testing.T) {
		tx, err := genericsContract.CreateAndTransferSuiToken(ctx, opts)
		require.NoError(t, err)
		require.Equal(t, "success", tx.Effects.Status.Status)

		tokenId, _, err := FindCreatedObject(tx.ObjectChanges, "generics::Token")
		require.NoError(t, err)
		t.Logf("Created Token<SUI> with ID: %s", tokenId)
	})

	t.Run("Token<SUI> balance check", func(t *testing.T) {
		tx, err := genericsContract.CreateAndTransferSuiToken(ctx, opts)
		require.NoError(t, err)
		require.Equal(t, "success", tx.Effects.Status.Status)

		tokenId, initialSharedVersion, err := FindCreatedObject(tx.ObjectChanges, "generics::Token")
		require.NoError(t, err)

		tokenObj := bind.Object{
			Id:                   tokenId,
			InitialSharedVersion: initialSharedVersion,
		}
		typeArgs := []string{"0x2::sui::SUI"}
		balance, err := genericsContract.DevInspect().Balance(ctx, opts, typeArgs, tokenObj)
		require.NoError(t, err)
		require.Equal(t, uint64(0), balance)
	})

	t.Run("Resolve generics", func(t *testing.T) {
		typeArgs := []string{"vector<Box<u64>>", "Pair<address,bool>"}

		resolver, err := bind.NewTypeResolver([]string{"T", "U"}, typeArgs)
		require.NoError(t, err)

		resolved := resolver.ResolveType("Box<T>")
		require.Equal(t, "Box<vector<Box<u64>>>", resolved)

		resolved = resolver.ResolveType("Pair<T,U>")
		require.Equal(t, "Pair<vector<Box<u64>>,Pair<address,bool>>", resolved)
	})

	t.Run("DevInspect with generic", func(t *testing.T) {
		tx, err := genericsContract.CreateAndTransferSuiToken(ctx, opts)
		require.NoError(t, err)
		require.Equal(t, "success", tx.Effects.Status.Status)

		// Find the token
		tokenId, _, err := FindCreatedObject(tx.ObjectChanges, "generics::Token")
		require.NoError(t, err)

		tokenObj := bind.Object{Id: tokenId}

		// Check balance using DevInspect
		typeArgs := []string{"0x2::sui::SUI"}
		balance, err := genericsContract.DevInspect().Balance(ctx, opts, typeArgs, tokenObj)
		require.NoError(t, err)
		require.Equal(t, uint64(0), balance)
	})

	t.Run("Execute generic functions", func(t *testing.T) {
		ptb := transaction.NewTransaction()

		typeArgs := []string{"0x2::sui::SUI"}
		encoded, err := genericsContract.Encoder().CreateAndTransferToken(typeArgs)
		require.NoError(t, err)

		_, err = genericsContract.AppendPTB(ctx, opts, ptb, encoded)
		require.NoError(t, err)

		tx, err := bind.ExecutePTB(ctx, opts, client, ptb)
		require.NoError(t, err)
		require.Equal(t, "success", tx.Effects.Status.Status)

		tokenId, _, err := FindCreatedObject(tx.ObjectChanges, "generics::Token")
		require.NoError(t, err)

		tokenObj := bind.Object{Id: tokenId}
		balanceTx, err := genericsContract.Balance(ctx, opts, []string{"0x2::sui::SUI"}, tokenObj)
		require.NoError(t, err)
		require.Equal(t, "success", balanceTx.Effects.Status.Status)
	})
}
