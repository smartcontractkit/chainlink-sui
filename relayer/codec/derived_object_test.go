package codec

import (
	"testing"

	"github.com/block-vision/sui-go-sdk/mystenbcs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type RouterOwnableKey bool
type RouterKey bool

func TestDeriveRouterOwnerCapId(t *testing.T) {
	t.Run("test RouterOwnableKey from Move test", func(t *testing.T) {
		routerObjectId := "0x34401905bebdf8c04f3cd5f04f442a39372c8dc321c29edfb4f9cb30b23ab96"
		packageId := "0x1001"

		key := RouterOwnableKey(false)
		keyBytes, err := mystenbcs.Marshal(key)
		require.NoError(t, err)

		result, err := DeriveDerivedObjectID(routerObjectId, packageId, "router", "RouterOwnableKey", keyBytes)
		require.NoError(t, err)

		// Move contract derived value from /contracts/ccip/ccip_router/tests/router_tests.move
		// Test: test_derive_address()
		assert.Equal(t, "0x6b91999dc9fdc7ff1490b40df428c23503c852e0843e6384b5889eca95cdbd7d", result,
			"Incorrect derived OwnerCap ID")
	})
}

func TestDeriveRouterStateId(t *testing.T) {
	t.Run("test RouterKey derivation", func(t *testing.T) {
		routerObjectId := "0x34401905bebdf8c04f3cd5f04f442a39372c8dc321c29edfb4f9cb30b23ab96"
		packageId := "0x1001"

		key := RouterKey(false)
		keyBytes, err := mystenbcs.Marshal(key)
		require.NoError(t, err)

		result, err := DeriveDerivedObjectID(routerObjectId, packageId, "router", "RouterKey", keyBytes)
		require.NoError(t, err)

		// RouterState is created with RouterKey() in router.move
		// Expected value from Move test: test_derive_address()
		assert.Equal(t, "0xbe237b4b91fd359e2ceaafcc2839bc316c7dd4275b4d3cd684739e0642274aa4", result,
			"Incorrect derived RouterState ID")
	})
}
