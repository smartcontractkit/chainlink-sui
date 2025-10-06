package bind

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDeriveWithVectorU8Key(t *testing.T) {
	t.Run("vector<u8> key with RouterState and OwnerCap key", func(t *testing.T) {
		// These values are from /contracts/ccip/ccip_router/tests/router_tests.move
		// Test: `fun test_derive_address()`
		parentID := "0x34401905bebdf8c04f3cd5f04f442a39372c8dc321c29edfb4f9cb30b23ab96"
		expectedRouterStateID := "0xc2ab753588210ab5de22dca1caf6e6d18a0b514c28c1975655c5769117d6f9ef"
		expectedOwnerCapID := "0x8e574462de77f45ea5bf3e8c1da19bcf081d25796376367e311926a8f993177e"

		result1, err := DeriveObjectIDWithVectorU8Key(parentID, []byte("RouterState"))
		require.NoError(t, err)
		assert.Equal(t, expectedRouterStateID, result1)

		result2, err := DeriveObjectIDWithVectorU8Key(parentID, []byte("CCIP_OWNABLE"))
		require.NoError(t, err)
		assert.Equal(t, expectedOwnerCapID, result2)

		// Different keys should produce different addresses
		assert.NotEqual(t, result1, result2,
			"Different vector<u8> keys should produce different addresses")
	})
}
