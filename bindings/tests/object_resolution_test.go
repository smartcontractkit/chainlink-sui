//go:build integration

package tests

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	testpackage "github.com/smartcontractkit/chainlink-sui/bindings/packages/test"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
)

func TestObjectResolution(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	signer, client := testenv.SetupEnvironment(t)

	opts := &bind.CallOpts{
		Signer:           signer,
		WaitForExecution: true,
	}

	testContract, _, err := testpackage.PublishTest(ctx, opts, client)
	require.NoError(t, err)

	t.Run("SharedObjectResolution", func(t *testing.T) {
		// Initialize creates a shared counter
		counter := testContract.Counter()
		createTx, err := counter.Initialize(ctx, opts)
		require.NoError(t, err)
		require.Equal(t, "success", createTx.Effects.Status.Status)

		// Find the counter object
		counterId, sharedVersion, err := FindCreatedObject(createTx.ObjectChanges, "counter::Counter")
		require.NoError(t, err)
		require.NotNil(t, sharedVersion, "Counter should be shared")
		require.Positive(t, *sharedVersion, "Counter should be shared")

		t.Run("DevInspectWithMinimalMetadata", func(t *testing.T) {
			// Create an object with just the ID (simulating what would create UnresolvedObject)
			counterObj := bind.Object{
				Id: counterId,
				// NOT setting InitialSharedVersion - this tests UnresolvedObject handling
			}

			// This previously failed with "invalid value: integer `3`" error
			// Now it works because we resolve the object
			count, err := counter.DevInspect().GetCount(ctx, opts, counterObj)
			require.NoError(t, err, "DevInspect should work even without full object metadata")
			require.Equal(t, uint64(0), count)
		})

		t.Run("ExecuteWithMinimalMetadata", func(t *testing.T) {
			// Create an object with just the ID
			counterObj := bind.Object{
				Id: counterId,
				// NOT setting InitialSharedVersion
			}

			// Execute should also work
			tx, err := counter.IncrementByOne(ctx, opts, counterObj)
			require.NoError(t, err)
			require.Equal(t, "success", tx.Effects.Status.Status)
		})

		t.Run("GetObjectHelper", func(t *testing.T) {
			// Use the helper to get a fully resolved object
			resolvedObj, err := bind.GetObject(ctx, client, counterId)
			require.NoError(t, err)
			require.NotNil(t, resolvedObj)
			require.Equal(t, counterId, resolvedObj.Id)
			require.NotNil(t, resolvedObj.InitialSharedVersion)
			require.Equal(t, *sharedVersion, *resolvedObj.InitialSharedVersion)

			count, err := counter.DevInspect().GetCount(ctx, opts, *resolvedObj)
			require.NoError(t, err)
			require.Equal(t, uint64(1), count) // Should be 1 after increment
		})

		t.Run("MultipleCallsWithCaching", func(t *testing.T) {
			// Create minimal object
			counterObj := bind.Object{Id: counterId}

			// Multiple calls should all work (and use cached resolution)
			for range 3 {
				count, err := counter.DevInspect().GetCount(ctx, opts, counterObj)
				require.NoError(t, err)
				require.Equal(t, uint64(1), count)
			}
		})
	})

	t.Run("SharedSampleObjectResolution", func(t *testing.T) {
		complexContract := testContract.Complex()

		// Create a shared SampleObject using NewObjectWithTransfer
		signerAddr, _ := signer.GetAddress()
		createTx, err := complexContract.NewObjectWithTransfer(ctx, opts,
			[]byte("test-id"), 42, signerAddr, []string{signerAddr})
		require.NoError(t, err)
		require.Equal(t, "success", createTx.Effects.Status.Status)

		objectId, sharedVersion, err := FindCreatedObject(createTx.ObjectChanges, "complex::SampleObject")
		require.NoError(t, err)
		require.NotNil(t, sharedVersion, "SampleObject should be shared")
		require.Positive(t, *sharedVersion, "SampleObject should be shared")

		t.Run("DevInspect with auto object resolution", func(t *testing.T) {
			obj := bind.Object{Id: objectId}

			result, err := complexContract.DevInspect().CheckWithObjectRef(ctx, opts, obj)
			require.NoError(t, err)
			require.Equal(t, uint64(42), result)
		})

		t.Run("GetObject", func(t *testing.T) {
			resolvedObj, err := bind.GetObject(ctx, client, objectId)
			require.NoError(t, err)
			require.NotNil(t, resolvedObj)
			require.Equal(t, objectId, resolvedObj.Id)
			require.NotNil(t, resolvedObj.InitialSharedVersion)
			require.Equal(t, *sharedVersion, *resolvedObj.InitialSharedVersion)
		})
	})

	t.Run("ObjectResolverCaching", func(t *testing.T) {
		resolver := bind.NewObjectResolver(client)

		// Deploy a fresh contract to get a test object
		_, initTx, err := testpackage.PublishTest(ctx, opts, client)
		require.NoError(t, err)

		// Find any created object from the transaction
		testObjectId := ""
		for _, change := range initTx.ObjectChanges {
			if change.Type == "created" && change.ObjectType != "" {
				testObjectId = change.ObjectId
				break
			}
		}
		require.NotEmpty(t, testObjectId, "No test object found")

		// First call - should fetch from blockchain
		obj1, err := resolver.GetObject(ctx, testObjectId)
		require.NoError(t, err)
		require.NotNil(t, obj1)

		// Second call - should use cache
		obj2, err := resolver.GetObject(ctx, testObjectId)
		require.NoError(t, err)
		require.NotNil(t, obj2)

		// Both should be identical
		require.Equal(t, obj1.Id, obj2.Id)
		require.Equal(t, obj1.InitialSharedVersion, obj2.InitialSharedVersion)

		// Clear cache
		resolver.ClearCache()

		// Next call should fetch from blockchain again
		obj3, err := resolver.GetObject(ctx, testObjectId)
		require.NoError(t, err)
		require.NotNil(t, obj3)
		require.Equal(t, obj1.Id, obj3.Id)
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		resolver := bind.NewObjectResolver(client)

		// Test with invalid format
		_, err := resolver.GetObject(ctx, "invalid-id")
		require.Error(t, err)
		// The error could be either from address validation or from fetching
		require.True(t,
			strings.Contains(err.Error(), "invalid object ID") ||
				strings.Contains(err.Error(), "failed to fetch object"),
			"Expected error about invalid object ID or fetch failure, got: %v", err)

		// Test with non-existent object
		_, err = resolver.GetObject(ctx, "0x0000000000000000000000000000000000000000000000000000000000000000")
		require.Error(t, err)
	})
}
