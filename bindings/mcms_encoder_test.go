//go:build unit

package mcmsencoder

import (
	"fmt"
	"testing"

	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/block-vision/sui-go-sdk/transaction"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
)

// Helper function to create a 32-byte address from a hex string
func createAddress(addr string) []byte {
	// Ensure address has "0x" prefix
	if len(addr) < 2 || addr[:2] != "0x" {
		addr = "0x" + addr
	}
	// Pad to 32 bytes (64 hex chars + 2 for "0x")
	for len(addr) < 66 {
		addr = "0x0" + addr[2:]
	}

	bytes := make([]byte, 32)
	for i := 0; i < 32; i++ {
		var b byte
		for j := 0; j < 2; j++ {
			c := addr[2+i*2+j]
			if c >= '0' && c <= '9' {
				b = b*16 + (c - '0')
			} else if c >= 'a' && c <= 'f' {
				b = b*16 + (c - 'a' + 10)
			} else if c >= 'A' && c <= 'F' {
				b = b*16 + (c - 'A' + 10)
			}
		}
		bytes[i] = b
	}
	return bytes
}

// Helper function to serialize a single address for BCS
func serializeAddress(addr string) []byte {
	s := &bcs.Serializer{}
	s.FixedBytes(createAddress(addr))
	return s.ToBytes()
}

// Helper function to serialize multiple addresses for BCS
func serializeAddresses(addrs ...string) []byte {
	s := &bcs.Serializer{}
	for _, addr := range addrs {
		s.FixedBytes(createAddress(addr))
	}
	return s.ToBytes()
}

// Helper function to extract object ID from an EncodedCallArgument
func extractObjectID(arg *bind.EncodedCallArgument) (string, error) {
	if arg == nil {
		return "", fmt.Errorf("nil argument")
	}

	if arg.CallArg == nil {
		return "", fmt.Errorf("nil CallArg")
	}

	if arg.CallArg.UnresolvedObject != nil {
		// Convert bytes to hex string
		objBytes := arg.CallArg.UnresolvedObject.ObjectId
		return toHexString(objBytes[:]), nil
	}

	if arg.CallArg.Object != nil {
		if arg.CallArg.Object.ImmOrOwnedObject != nil {
			return toHexString(arg.CallArg.Object.ImmOrOwnedObject.ObjectId[:]), nil
		}
		if arg.CallArg.Object.SharedObject != nil {
			return toHexString(arg.CallArg.Object.SharedObject.ObjectId[:]), nil
		}
	}

	return "", fmt.Errorf("no object ID found in argument")
}

func TestEncodeEntryPointArg_FeeQuoter(t *testing.T) {
	encoder := &CCIPEntrypointArgEncoder{
		registryObjID: "0x1234567890123456789012345678901234567890123456789012345678901234",
	}

	target := "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	stateObjID := "0x9999999999999999999999999999999999999999999999999999999999999999"
	executingCallbackParams := &transaction.Argument{}

	t.Run("update_prices_with_owner_cap", func(t *testing.T) {
		// For this function, ccipRef must match stateObjID
		data := serializeAddress(stateObjID)

		result, err := encoder.EncodeEntryPointArg(
			executingCallbackParams,
			target,
			"fee_quoter",
			"update_prices_with_owner_cap",
			stateObjID,
			data,
			[]string{"0x1::sui::SUI"},
		)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "fee_quoter", result.Module.ModuleName)
		assert.Equal(t, "mcms_update_prices_with_owner_cap", result.Function)

		// Verify deserialization - the ccipRef should be extracted from data and match stateObjID
		require.Len(t, result.CallArgs, 4, "Expected 4 arguments: ccipRef, registry, clock, executingCallbackParams")

		// Verify the ccipRef was deserialized correctly and matches stateObjID
		ccipRefFromResult, err := extractObjectID(result.CallArgs[0])
		require.NoError(t, err, "Failed to extract ccipRef object ID")
		assert.Equal(t, stateObjID, ccipRefFromResult, "CcipRef should match stateObjID (from BCS data)")
	})

	t.Run("fallback_case", func(t *testing.T) {
		// Test an unhandled function that falls back to default encoding
		data := []byte{}

		result, err := encoder.EncodeEntryPointArg(
			executingCallbackParams,
			target,
			"fee_quoter",
			"apply_fee_token_updates",
			stateObjID,
			data,
			[]string{"0x1::sui::SUI"},
		)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "fee_quoter", result.Module.ModuleName)
		assert.Equal(t, "mcms_apply_fee_token_updates", result.Function)
	})
}

func TestEncodeEntryPointArg_Offramp(t *testing.T) {
	encoder := &CCIPEntrypointArgEncoder{
		registryObjID: "0x1234567890123456789012345678901234567890123456789012345678901234",
	}

	target := "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	stateObjID := "0x9999999999999999999999999999999999999999999999999999999999999999"
	ccipRefID := "0x8888888888888888888888888888888888888888888888888888888888888888"
	executingCallbackParams := &transaction.Argument{}

	testCases := []string{
		"set_dynamic_config",
		"apply_source_chain_config_updates",
		"set_ocr3_config",
		"transfer_ownership",
		"execute_ownership_transfer",
	}

	for _, fn := range testCases {
		t.Run(fn, func(t *testing.T) {
			data := serializeAddress(ccipRefID)

			result, err := encoder.EncodeEntryPointArg(
				executingCallbackParams,
				target,
				"offramp",
				fn,
				stateObjID,
				data,
				[]string{"0x1::sui::SUI"},
			)

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "offramp", result.Module.ModuleName)
			assert.Equal(t, "mcms_"+fn, result.Function)

			// Verify deserialization - the ccipRef should be extracted from data and match stateObjID
			require.Len(t, result.CallArgs, 4, "Expected 4 arguments: ccipRef, state, registry, executingCallbackParams")

			// Verify the ccipRef was deserialized correctly and matches stateObjID
			ccipRefFromResult, err := extractObjectID(result.CallArgs[0])
			require.NoError(t, err, "Failed to extract ccipRef object ID")
			assert.Equal(t, ccipRefID, ccipRefFromResult, "CcipRef should match ccipRefID (from BCS data)")
		})
	}
}

func TestEncodeEntryPointArg_Onramp(t *testing.T) {
	encoder := &CCIPEntrypointArgEncoder{
		registryObjID: "0x1234567890123456789012345678901234567890123456789012345678901234",
	}

	target := "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	stateObjID := "0x9999999999999999999999999999999999999999999999999999999999999999"
	ccipRefID := "0x8888888888888888888888888888888888888888888888888888888888888888"
	executingCallbackParams := &transaction.Argument{}

	t.Run("withdraw_fee_tokens", func(t *testing.T) {
		ownerCapID := "0x7777777777777777777777777777777777777777777777777777777777777777"
		feeTokenMetadataID := "0x6666666666666666666666666666666666666666666666666666666666666666"

		data := serializeAddresses(ccipRefID, stateObjID, ownerCapID, feeTokenMetadataID)

		result, err := encoder.EncodeEntryPointArg(
			executingCallbackParams,
			target,
			"onramp",
			"withdraw_fee_tokens",
			stateObjID,
			data,
			[]string{"0x1::sui::SUI"},
		)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "onramp", result.Module.ModuleName)
		assert.Equal(t, "mcms_withdraw_fee_tokens", result.Function)
		assert.Len(t, result.TypeArgs, 1, "Expected 1 type argument")

		// Verify deserialization - check that the encoder correctly deserialized BCS data
		require.Len(t, result.CallArgs, 5, "Expected 5 arguments: ccipRef, state, registry, coinMetadata, executingCallbackParams")

		// Verify ccipRef was deserialized correctly (1st argument)
		ccipRefFromResult, err := extractObjectID(result.CallArgs[0])
		require.NoError(t, err, "Failed to extract ccipRef object ID")
		assert.Equal(t, toHexString(createAddress(ccipRefID)), ccipRefFromResult, "CcipRef should match deserialized value from data")

		// Verify state object matches
		stateFromResult, err := extractObjectID(result.CallArgs[1])
		require.NoError(t, err, "Failed to extract state object ID")
		assert.Equal(t, stateObjID, stateFromResult, "State should match the provided stateObjID")

		// Verify coinMetadata was deserialized correctly (4th argument, index 3)
		coinMetadataFromResult, err := extractObjectID(result.CallArgs[3])
		require.NoError(t, err, "Failed to extract coinMetadata object ID")
		assert.Equal(t, toHexString(createAddress(feeTokenMetadataID)), coinMetadataFromResult, "CoinMetadata should match deserialized value from data")
	})

	ccipTestCases := []string{
		"set_dynamic_config",
		"apply_dest_chain_config_updates",
		"apply_allowlist_updates",
		"transfer_ownership",
		"execute_ownership_transfer",
	}

	for _, fn := range ccipTestCases {
		t.Run(fn, func(t *testing.T) {
			data := serializeAddress(ccipRefID)

			result, err := encoder.EncodeEntryPointArg(
				executingCallbackParams,
				target,
				"onramp",
				fn,
				stateObjID,
				data,
				[]string{"0x1::sui::SUI"},
			)

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "onramp", result.Module.ModuleName)
			assert.Equal(t, "mcms_"+fn, result.Function)

			// Verify deserialization - the ccipRef should be extracted from data
			require.Len(t, result.CallArgs, 4, "Expected 4 arguments: ccipRef, state, registry, executingCallbackParams")

			// Verify the ccipRef was deserialized correctly
			ccipRefFromResult, err := extractObjectID(result.CallArgs[0])
			require.NoError(t, err, "Failed to extract ccipRef object ID")
			assert.Equal(t, ccipRefID, ccipRefFromResult, "CcipRef should match ccipRefID (from BCS data)")
		})
	}
}

func TestEncodeEntryPointArg_Router(t *testing.T) {
	encoder := &CCIPEntrypointArgEncoder{
		registryObjID: "0x1234567890123456789012345678901234567890123456789012345678901234",
	}

	target := "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	stateObjID := "0x9999999999999999999999999999999999999999999999999999999999999999"
	executingCallbackParams := &transaction.Argument{}

	t.Run("accept_ownership", func(t *testing.T) {
		data := []byte{}

		result, err := encoder.EncodeEntryPointArg(
			executingCallbackParams,
			target,
			"router",
			"accept_ownership",
			stateObjID,
			data,
			[]string{"0x1::sui::SUI"},
		)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "router", result.Module.ModuleName)
		assert.Equal(t, "mcms_accept_ownership", result.Function)

		// Verify deserialization - the state should be extracted from data
		require.Len(t, result.CallArgs, 3, "Expected 3 arguments: state, registry, executingCallbackParams")

		// Verify the state was deserialized correctly
		stateFromResult, err := extractObjectID(result.CallArgs[0])
		require.NoError(t, err, "Failed to extract state object ID")
		assert.Equal(t, stateObjID, stateFromResult, "State should match stateObjID (from BCS data)")
	})
}

func TestEncodeEntryPointArg_BurnMintTokenPool(t *testing.T) {
	encoder := &CCIPEntrypointArgEncoder{
		registryObjID:      "0x1234567890123456789012345678901234567890123456789012345678901234",
		deployerStateObjID: "0x8888888888888888888888888888888888888888888888888888888888888888",
	}

	target := "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	stateObjID := "0x9999999999999999999999999999999999999999999999999999999999999999"
	executingCallbackParams := &transaction.Argument{}

	t.Run("accept_ownership", func(t *testing.T) {
		data := []byte{}

		result, err := encoder.EncodeEntryPointArg(
			executingCallbackParams,
			target,
			"burn_mint_token_pool",
			"accept_ownership",
			stateObjID,
			data,
			[]string{"0x1::sui::SUI"},
		)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "burn_mint_token_pool", result.Module.ModuleName)
		assert.Equal(t, "mcms_accept_ownership", result.Function)
	})

	typeArgTestCases := []string{
		"set_allowlist_enabled",
		"apply_allowlist_updates",
		"apply_chain_updates",
		"add_remote_pool",
		"remove_remote_pool",
		"transfer_ownership",
	}

	for _, fn := range typeArgTestCases {
		t.Run(fn, func(t *testing.T) {
			data := []byte{}

			result, err := encoder.EncodeEntryPointArg(
				executingCallbackParams,
				target,
				"burn_mint_token_pool",
				fn,
				stateObjID,
				data,
				[]string{"0x1::sui::SUI"},
			)

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "burn_mint_token_pool", result.Module.ModuleName)
			assert.Equal(t, "mcms_"+fn, result.Function)

			require.Len(t, result.TypeArgs, 1, "Expected 1 type argument")
			require.Len(t, result.CallArgs, 3, "Expected 3 arguments: state, registry, executingCallbackParams")

			// Verify the state was deserialized correctly
			stateFromResult, err := extractObjectID(result.CallArgs[0])
			require.NoError(t, err, "Failed to extract state object ID")
			assert.Equal(t, stateObjID, stateFromResult, "State should match stateObjID (from BCS data)")
		})
	}

	t.Run("execute_ownership_transfer", func(t *testing.T) {
		data := []byte{}

		result, err := encoder.EncodeEntryPointArg(
			executingCallbackParams,
			target,
			"burn_mint_token_pool",
			"execute_ownership_transfer",
			stateObjID,
			data,
			[]string{"0x1::sui::SUI"},
		)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "burn_mint_token_pool", result.Module.ModuleName)
		assert.Equal(t, "mcms_execute_ownership_transfer", result.Function)

		require.Len(t, result.TypeArgs, 1, "Expected 1 type argument")
		require.Len(t, result.CallArgs, 4, "Expected 4 arguments: state, registry, deployer_state, executingCallbackParams")

		// Verify the state was deserialized correctly
		stateFromResult, err := extractObjectID(result.CallArgs[0])
		require.NoError(t, err, "Failed to extract state object ID")
		assert.Equal(t, stateObjID, stateFromResult, "State should match stateObjID")

		// Verify the deployer state was set correctly
		deployerStateFromResult, err := extractObjectID(result.CallArgs[2])
		require.NoError(t, err, "Failed to extract deployer state object ID")
		assert.Equal(t, encoder.deployerStateObjID, deployerStateFromResult, "Deployer state should match deployerStateObjID")
	})

	rateLimiterTestCases := []string{
		"set_chain_rate_limiter_configs",
		"set_chain_rate_limiter_config",
	}

	for _, fn := range rateLimiterTestCases {
		t.Run(fn, func(t *testing.T) {
			data := []byte{}

			result, err := encoder.EncodeEntryPointArg(
				executingCallbackParams,
				target,
				"burn_mint_token_pool",
				fn,
				stateObjID,
				data,
				[]string{"0x1::sui::SUI"},
			)

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "burn_mint_token_pool", result.Module.ModuleName)
			assert.Equal(t, "mcms_"+fn, result.Function)

			require.Len(t, result.TypeArgs, 1, "Expected 1 type argument")
			require.Len(t, result.CallArgs, 4, "Expected 4 arguments: state, clock, registry, executingCallbackParams")

			// Verify the state was deserialized correctly
			stateFromResult, err := extractObjectID(result.CallArgs[0])
			require.NoError(t, err, "Failed to extract state object ID")
			assert.Equal(t, stateObjID, stateFromResult, "State should match stateObjID (from BCS data)")
		})
	}
}

func TestEncodeEntryPointArg_LockReleaseTokenPool(t *testing.T) {
	encoder := &CCIPEntrypointArgEncoder{
		registryObjID:      "0x1234567890123456789012345678901234567890123456789012345678901234",
		deployerStateObjID: "0x8888888888888888888888888888888888888888888888888888888888888888",
	}

	target := "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	stateObjID := "0x9999999999999999999999999999999999999999999999999999999999999999"
	executingCallbackParams := &transaction.Argument{}

	t.Run("accept_ownership", func(t *testing.T) {
		data := []byte{}

		result, err := encoder.EncodeEntryPointArg(
			executingCallbackParams,
			target,
			"lock_release_token_pool",
			"accept_ownership",
			stateObjID,
			data,
			[]string{"0x1::sui::SUI"},
		)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "lock_release_token_pool", result.Module.ModuleName)
		assert.Equal(t, "mcms_accept_ownership", result.Function)
	})

	typeArgTestCases := []string{
		"set_rebalancer",
		"set_allowlist_enabled",
		"apply_allowlist_updates",
		"apply_chain_updates",
		"add_remote_pool",
		"remove_remote_pool",
		"transfer_ownership",
	}

	for _, fn := range typeArgTestCases {
		t.Run(fn, func(t *testing.T) {
			data := []byte{}

			result, err := encoder.EncodeEntryPointArg(
				executingCallbackParams,
				target,
				"lock_release_token_pool",
				fn,
				stateObjID,
				data,
				[]string{"0x1::sui::SUI"},
			)

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "lock_release_token_pool", result.Module.ModuleName)
			assert.Equal(t, "mcms_"+fn, result.Function)

			require.Len(t, result.TypeArgs, 1, "Expected 1 type argument")
			require.Len(t, result.CallArgs, 3, "Expected 3 arguments: state, registry, executingCallbackParams")

			// Verify the state was deserialized correctly
			stateFromResult, err := extractObjectID(result.CallArgs[0])
			require.NoError(t, err, "Failed to extract state object ID")
			assert.Equal(t, stateObjID, stateFromResult, "State should match stateObjID (from BCS data)")
		})
	}

	t.Run("execute_ownership_transfer", func(t *testing.T) {
		data := []byte{}

		result, err := encoder.EncodeEntryPointArg(
			executingCallbackParams,
			target,
			"lock_release_token_pool",
			"execute_ownership_transfer",
			stateObjID,
			data,
			[]string{"0x1::sui::SUI"},
		)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "lock_release_token_pool", result.Module.ModuleName)
		assert.Equal(t, "mcms_execute_ownership_transfer", result.Function)

		require.Len(t, result.TypeArgs, 1, "Expected 1 type argument")
		require.Len(t, result.CallArgs, 4, "Expected 4 arguments: state, registry, deployer_state, executingCallbackParams")

		// Verify the state was deserialized correctly
		stateFromResult, err := extractObjectID(result.CallArgs[0])
		require.NoError(t, err, "Failed to extract state object ID")
		assert.Equal(t, stateObjID, stateFromResult, "State should match stateObjID")

		// Verify the deployer state was set correctly
		deployerStateFromResult, err := extractObjectID(result.CallArgs[2])
		require.NoError(t, err, "Failed to extract deployer state object ID")
		assert.Equal(t, encoder.deployerStateObjID, deployerStateFromResult, "Deployer state should match deployerStateObjID")
	})

	rateLimiterTestCases := []string{
		"set_chain_rate_limiter_configs",
		"set_chain_rate_limiter_config",
	}

	for _, fn := range rateLimiterTestCases {
		t.Run(fn, func(t *testing.T) {
			data := []byte{}

			result, err := encoder.EncodeEntryPointArg(
				executingCallbackParams,
				target,
				"lock_release_token_pool",
				fn,
				stateObjID,
				data,
				[]string{"0x1::sui::SUI"},
			)

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "lock_release_token_pool", result.Module.ModuleName)
			assert.Equal(t, "mcms_"+fn, result.Function)

			require.Len(t, result.TypeArgs, 1, "Expected 1 type argument")
			require.Len(t, result.CallArgs, 4, "Expected 4 arguments: state, clock, registry, executingCallbackParams")

			// Verify the state was deserialized correctly
			stateFromResult, err := extractObjectID(result.CallArgs[0])
			require.NoError(t, err, "Failed to extract state object ID")
			assert.Equal(t, stateObjID, stateFromResult, "State should match stateObjID (from BCS data)")
		})
	}
}

func TestEncodeEntryPointArg_ManagedTokenPool(t *testing.T) {
	encoder := &CCIPEntrypointArgEncoder{
		registryObjID:      "0x1234567890123456789012345678901234567890123456789012345678901234",
		deployerStateObjID: "0x8888888888888888888888888888888888888888888888888888888888888888",
	}

	target := "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	stateObjID := "0x9999999999999999999999999999999999999999999999999999999999999999"
	executingCallbackParams := &transaction.Argument{}

	t.Run("accept_ownership", func(t *testing.T) {
		data := []byte{}

		result, err := encoder.EncodeEntryPointArg(
			executingCallbackParams,
			target,
			"managed_token_pool",
			"accept_ownership",
			stateObjID,
			data,
			[]string{"0x1::sui::SUI"},
		)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "managed_token_pool", result.Module.ModuleName)
		assert.Equal(t, "mcms_accept_ownership", result.Function)
	})

	typeArgTestCases := []string{
		"set_allowlist_enabled",
		"apply_allowlist_updates",
		"apply_chain_updates",
		"add_remote_pool",
		"remove_remote_pool",
		"transfer_ownership",
	}

	for _, fn := range typeArgTestCases {
		t.Run(fn, func(t *testing.T) {
			data := []byte{}

			result, err := encoder.EncodeEntryPointArg(
				executingCallbackParams,
				target,
				"managed_token_pool",
				fn,
				stateObjID,
				data,
				[]string{"0x1::sui::SUI"},
			)

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "managed_token_pool", result.Module.ModuleName)
			assert.Equal(t, "mcms_"+fn, result.Function)

			require.Len(t, result.TypeArgs, 1, "Expected 1 type argument")
			require.Len(t, result.CallArgs, 3, "Expected 3 arguments: state, registry, executingCallbackParams")

			// Verify the state was deserialized correctly
			stateFromResult, err := extractObjectID(result.CallArgs[0])
			require.NoError(t, err, "Failed to extract state object ID")
			assert.Equal(t, stateObjID, stateFromResult, "State should match stateObjID (from BCS data)")
		})
	}

	t.Run("execute_ownership_transfer", func(t *testing.T) {
		data := []byte{}

		result, err := encoder.EncodeEntryPointArg(
			executingCallbackParams,
			target,
			"managed_token_pool",
			"execute_ownership_transfer",
			stateObjID,
			data,
			[]string{"0x1::sui::SUI"},
		)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "managed_token_pool", result.Module.ModuleName)
		assert.Equal(t, "mcms_execute_ownership_transfer", result.Function)

		require.Len(t, result.TypeArgs, 1, "Expected 1 type argument")
		require.Len(t, result.CallArgs, 4, "Expected 4 arguments: state, registry, deployer_state, executingCallbackParams")

		// Verify the state was deserialized correctly
		stateFromResult, err := extractObjectID(result.CallArgs[0])
		require.NoError(t, err, "Failed to extract state object ID")
		assert.Equal(t, stateObjID, stateFromResult, "State should match stateObjID")

		// Verify the deployer state was set correctly
		deployerStateFromResult, err := extractObjectID(result.CallArgs[2])
		require.NoError(t, err, "Failed to extract deployer state object ID")
		assert.Equal(t, encoder.deployerStateObjID, deployerStateFromResult, "Deployer state should match deployerStateObjID")
	})

	rateLimiterTestCases := []string{
		"set_chain_rate_limiter_configs",
		"set_chain_rate_limiter_config",
	}

	for _, fn := range rateLimiterTestCases {
		t.Run(fn, func(t *testing.T) {
			data := []byte{}

			result, err := encoder.EncodeEntryPointArg(
				executingCallbackParams,
				target,
				"managed_token_pool",
				fn,
				stateObjID,
				data,
				[]string{"0x1::sui::SUI"},
			)

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "managed_token_pool", result.Module.ModuleName)
			assert.Equal(t, "mcms_"+fn, result.Function)

			require.Len(t, result.TypeArgs, 1, "Expected 1 type argument")
			require.Len(t, result.CallArgs, 4, "Expected 4 arguments: state, clock, registry, executingCallbackParams")

			// Verify the state was deserialized correctly
			stateFromResult, err := extractObjectID(result.CallArgs[0])
			require.NoError(t, err, "Failed to extract state object ID")
			assert.Equal(t, stateObjID, stateFromResult, "State should match stateObjID (from BCS data)")
		})
	}
}

func TestEncodeEntryPointArg_UsdcTokenPool(t *testing.T) {
	encoder := &CCIPEntrypointArgEncoder{
		registryObjID:      "0x1234567890123456789012345678901234567890123456789012345678901234",
		deployerStateObjID: "0x8888888888888888888888888888888888888888888888888888888888888888",
	}

	target := "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	stateObjID := "0x9999999999999999999999999999999999999999999999999999999999999999"
	executingCallbackParams := &transaction.Argument{}

	t.Run("accept_ownership", func(t *testing.T) {
		data := []byte{}

		result, err := encoder.EncodeEntryPointArg(
			executingCallbackParams,
			target,
			"usdc_token_pool",
			"accept_ownership",
			stateObjID,
			data,
			[]string{"0x1::sui::SUI"},
		)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "usdc_token_pool", result.Module.ModuleName)
		assert.Equal(t, "mcms_accept_ownership", result.Function)
	})

	typeArgTestCases := []string{
		"set_allowlist_enabled",
		"apply_allowlist_updates",
		"apply_chain_updates",
		"add_remote_pool",
		"remove_remote_pool",
		"transfer_ownership",
	}

	for _, fn := range typeArgTestCases {
		t.Run(fn, func(t *testing.T) {
			data := []byte{}

			result, err := encoder.EncodeEntryPointArg(
				executingCallbackParams,
				target,
				"usdc_token_pool",
				fn,
				stateObjID,
				data,
				[]string{"0x1::sui::SUI"},
			)

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "usdc_token_pool", result.Module.ModuleName)
			assert.Equal(t, "mcms_"+fn, result.Function)

			require.Len(t, result.TypeArgs, 1, "Expected 1 type argument")
			require.Len(t, result.CallArgs, 3, "Expected 3 arguments: state, registry, executingCallbackParams")

			// Verify the state was deserialized correctly
			stateFromResult, err := extractObjectID(result.CallArgs[0])
			require.NoError(t, err, "Failed to extract state object ID")
			assert.Equal(t, stateObjID, stateFromResult, "State should match stateObjID (from BCS data)")
		})
	}

	t.Run("execute_ownership_transfer", func(t *testing.T) {
		data := []byte{}

		result, err := encoder.EncodeEntryPointArg(
			executingCallbackParams,
			target,
			"usdc_token_pool",
			"execute_ownership_transfer",
			stateObjID,
			data,
			[]string{"0x1::sui::SUI"},
		)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "usdc_token_pool", result.Module.ModuleName)
		assert.Equal(t, "mcms_execute_ownership_transfer", result.Function)

		require.Len(t, result.TypeArgs, 1, "Expected 1 type argument")
		require.Len(t, result.CallArgs, 4, "Expected 4 arguments: state, registry, deployer_state, executingCallbackParams")

		// Verify the state was deserialized correctly
		stateFromResult, err := extractObjectID(result.CallArgs[0])
		require.NoError(t, err, "Failed to extract state object ID")
		assert.Equal(t, stateObjID, stateFromResult, "State should match stateObjID")

		// Verify the deployer state was set correctly
		deployerStateFromResult, err := extractObjectID(result.CallArgs[2])
		require.NoError(t, err, "Failed to extract deployer state object ID")
		assert.Equal(t, encoder.deployerStateObjID, deployerStateFromResult, "Deployer state should match deployerStateObjID")
	})

	rateLimiterTestCases := []string{
		"set_chain_rate_limiter_configs",
		"set_chain_rate_limiter_config",
	}

	for _, fn := range rateLimiterTestCases {
		t.Run(fn, func(t *testing.T) {
			data := []byte{}

			result, err := encoder.EncodeEntryPointArg(
				executingCallbackParams,
				target,
				"usdc_token_pool",
				fn,
				stateObjID,
				data,
				[]string{"0x1::sui::SUI"},
			)

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "usdc_token_pool", result.Module.ModuleName)
			assert.Equal(t, "mcms_"+fn, result.Function)

			require.Len(t, result.TypeArgs, 1, "Expected 1 type argument")
			require.Len(t, result.CallArgs, 4, "Expected 4 arguments: state, clock, registry, executingCallbackParams")

			// Verify the state was deserialized correctly
			stateFromResult, err := extractObjectID(result.CallArgs[0])
			require.NoError(t, err, "Failed to extract state object ID")
			assert.Equal(t, stateObjID, stateFromResult, "State should match stateObjID (from BCS data)")
		})
	}
}

func TestEncodeEntryPointArg_ManagedToken(t *testing.T) {
	encoder := &CCIPEntrypointArgEncoder{
		registryObjID: "0x1234567890123456789012345678901234567890123456789012345678901234",
	}

	target := "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	stateObjID := "0x9999999999999999999999999999999999999999999999999999999999999999"
	executingCallbackParams := &transaction.Argument{}

	t.Run("accept_ownership", func(t *testing.T) {
		data := []byte{}

		result, err := encoder.EncodeEntryPointArg(
			executingCallbackParams,
			target,
			"managed_token",
			"accept_ownership",
			stateObjID,
			data,
			[]string{"0x1::sui::SUI"},
		)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "managed_token", result.Module.ModuleName)
		assert.Equal(t, "mcms_accept_ownership", result.Function)

		require.Len(t, result.TypeArgs, 1, "Expected 1 type argument")
		require.Len(t, result.CallArgs, 3, "Expected 3 arguments: state, registry, executingCallbackParams")

		// Verify the state was deserialized correctly
		stateFromResult, err := extractObjectID(result.CallArgs[0])
		require.NoError(t, err, "Failed to extract state object ID")
		assert.Equal(t, stateObjID, stateFromResult, "State should match stateObjID (from BCS data)")
	})

	t.Run("configure_new_minter", func(t *testing.T) {
		data := []byte{}

		result, err := encoder.EncodeEntryPointArg(
			executingCallbackParams,
			target,
			"managed_token",
			"configure_new_minter",
			stateObjID,
			data,
			[]string{"0x1::sui::SUI"},
		)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "managed_token", result.Module.ModuleName)
		assert.Equal(t, "mcms_configure_new_minter", result.Function)

		require.Len(t, result.TypeArgs, 1, "Expected 1 type argument")
		require.Len(t, result.CallArgs, 3, "Expected 3 arguments: state, registry, executingCallbackParams")

		// Verify the state was deserialized correctly
		stateFromResult, err := extractObjectID(result.CallArgs[0])
		require.NoError(t, err, "Failed to extract state object ID")
		assert.Equal(t, stateObjID, stateFromResult, "State should match stateObjID (from BCS data)")
	})

	denyListTestCases := []string{
		"increment_mint_allowance",
		"set_unlimited_mint_allowances",
		"blocklist",
		"unblocklist",
		"pause",
	}

	for _, fn := range denyListTestCases {
		t.Run(fn, func(t *testing.T) {
			ownerCapID := "0x7777777777777777777777777777777777777777777777777777777777777777"
			denyListID := "0x6666666666666666666666666666666666666666666666666666666666666666"

			data := serializeAddresses(stateObjID, ownerCapID, denyListID)

			result, err := encoder.EncodeEntryPointArg(
				executingCallbackParams,
				target,
				"managed_token",
				fn,
				stateObjID,
				data,
				[]string{"0x1::sui::SUI"},
			)

			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, "managed_token", result.Module.ModuleName)
			assert.Equal(t, "mcms_"+fn, result.Function)

			require.NotEmpty(t, result.CallArgs, "Expected call arguments to be populated")

			denyListFromResult, err := extractObjectID(result.CallArgs[2])
			require.NoError(t, err, "Failed to extract denyList object ID")
			assert.Equal(t, toHexString(createAddress(denyListID)), denyListFromResult, "DenyList should match deserialized value from data")
		})
	}
}

func TestEncodeEntryPointArg_UnknownModule(t *testing.T) {
	encoder := &CCIPEntrypointArgEncoder{
		registryObjID: "0x1234567890123456789012345678901234567890123456789012345678901234",
	}

	target := "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	stateObjID := "0x9999999999999999999999999999999999999999999999999999999999999999"
	executingCallbackParams := &transaction.Argument{}

	t.Run("fallback_to_default_encoder", func(t *testing.T) {
		data := []byte{}

		result, err := encoder.EncodeEntryPointArg(
			executingCallbackParams,
			target,
			"unknown_module",
			"unknown_function",
			stateObjID,
			data,
			[]string{"0x1::sui::SUI"},
		)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "unknown_module", result.Module.ModuleName)
		assert.Equal(t, "mcms_unknown_function", result.Function)
	})
}

func TestEncodeEntryPointArg_ErrorCases(t *testing.T) {
	encoder := &CCIPEntrypointArgEncoder{
		registryObjID: "0x1234567890123456789012345678901234567890123456789012345678901234",
	}

	target := "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	stateObjID := "0x9999999999999999999999999999999999999999999999999999999999999999"
	wrongStateObjID := "0x8888888888888888888888888888888888888888888888888888888888888888"
	executingCallbackParams := &transaction.Argument{}

	t.Run("fee_quoter_update_prices_with_owner_cap_mismatch", func(t *testing.T) {
		// ccipRef does not match stateObjID - should error
		data := serializeAddress(wrongStateObjID)

		_, err := encoder.EncodeEntryPointArg(
			executingCallbackParams,
			target,
			"fee_quoter",
			"update_prices_with_owner_cap",
			stateObjID,
			data,
			[]string{"0x1::sui::SUI"},
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "does not match state object")
	})

	t.Run("onramp_withdraw_fee_tokens_state_mismatch", func(t *testing.T) {
		ccipRefID := "0x8888888888888888888888888888888888888888888888888888888888888888"
		ownerCapID := "0x7777777777777777777777777777777777777777777777777777777777777777"
		feeTokenMetadataID := "0x6666666666666666666666666666666666666666666666666666666666666666"

		// Second address is wrong state
		data := serializeAddresses(ccipRefID, wrongStateObjID, ownerCapID, feeTokenMetadataID)

		_, err := encoder.EncodeEntryPointArg(
			executingCallbackParams,
			target,
			"onramp",
			"withdraw_fee_tokens",
			stateObjID,
			data,
			[]string{"0x1::sui::SUI"},
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "does not match state object")
	})

	t.Run("managed_token_state_mismatch", func(t *testing.T) {
		ownerCapID := "0x7777777777777777777777777777777777777777777777777777777777777777"
		denyListID := "0x6666666666666666666666666666666666666666666666666666666666666666"

		// First address is wrong state
		data := serializeAddresses(wrongStateObjID, ownerCapID, denyListID)

		_, err := encoder.EncodeEntryPointArg(
			executingCallbackParams,
			target,
			"managed_token",
			"pause",
			stateObjID,
			data,
			[]string{"0x1::sui::SUI"},
		)

		require.Error(t, err)
		assert.Contains(t, err.Error(), "does not match state object")
	})
}

func TestBCSDeserialization(t *testing.T) {
	t.Run("deserialize_single_address", func(t *testing.T) {
		expectedAddr := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
		data := serializeAddress(expectedAddr)

		result := deserializeFirst32Bytes(data)

		assert.Equal(t, 32, len(result))
		assert.Equal(t, toHexString(result), toHexString(createAddress(expectedAddr)))
	})

	t.Run("deserialize_multiple_addresses", func(t *testing.T) {
		addr1 := "0x1111111111111111111111111111111111111111111111111111111111111111"
		addr2 := "0x2222222222222222222222222222222222222222222222222222222222222222"
		addr3 := "0x3333333333333333333333333333333333333333333333333333333333333333"

		data := serializeAddresses(addr1, addr2, addr3)

		// Deserialize each address sequentially
		deserializer := bcs.NewDeserializer(data)
		result1 := deserializer.ReadFixedBytes(SuiAddressLength)
		result2 := deserializer.ReadFixedBytes(SuiAddressLength)
		result3 := deserializer.ReadFixedBytes(SuiAddressLength)

		assert.Equal(t, toHexString(result1), toHexString(createAddress(addr1)))
		assert.Equal(t, toHexString(result2), toHexString(createAddress(addr2)))
		assert.Equal(t, toHexString(result3), toHexString(createAddress(addr3)))
	})

	t.Run("verify_address_format", func(t *testing.T) {
		// Test that addresses are correctly padded and formatted
		shortAddr := "0x123"
		data := serializeAddress(shortAddr)

		result := deserializeFirst32Bytes(data)
		resultHex := toHexString(result)

		// Should be 0x-prefixed with 64 hex chars (32 bytes)
		assert.True(t, len(resultHex) == 66, "Address should be 0x + 64 hex chars")
		assert.True(t, resultHex[:2] == "0x", "Address should start with 0x")
		assert.Contains(t, resultHex, "123", "Address should contain the original value")
	})
}

func TestOverrideCall(t *testing.T) {
	call := &bind.EncodedCall{
		Module: bind.ModuleInformation{
			ModuleName: "original_module",
		},
		Function: "original_function",
	}

	result := overrideCall(call, "new_module", "some_function")

	assert.Equal(t, "new_module", result.Module.ModuleName)
	assert.Equal(t, "mcms_some_function", result.Function)

	// Test with mcms_ prefix already present
	result2 := overrideCall(call, "another_module", "mcms_prefixed_function")

	assert.Equal(t, "another_module", result2.Module.ModuleName)
	assert.Equal(t, "mcms_prefixed_function", result2.Function)
}

func TestHelperFunctions(t *testing.T) {
	t.Run("deserializeFirst32Bytes", func(t *testing.T) {
		addr := createAddress("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
		result := deserializeFirst32Bytes(addr)

		assert.Equal(t, 32, len(result))
		assert.Equal(t, addr, result)
	})

	t.Run("toHexString", func(t *testing.T) {
		data := []byte{0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef}
		result := toHexString(data)

		assert.Equal(t, "0x1234567890abcdef", result)
	})
}

func TestNewCCIPEntrypointArgEncoder(t *testing.T) {
	registryObjID := "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
	deployerStateObjID := "0x8888888888888888888888888888888888888888888888888888888888888888"
	encoder := NewCCIPEntrypointArgEncoder(registryObjID, deployerStateObjID)

	assert.NotNil(t, encoder)
	assert.Equal(t, registryObjID, encoder.registryObjID)
	assert.Equal(t, deployerStateObjID, encoder.deployerStateObjID)
}
