package offramp

import (
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testCcipPackageId    = "0x00000000000000000000000000000000000000000000000000000000ccipccip"
	testOffRampPackageId = "0x000000000000000000000000000000000000000000000000000000000ff2a3f0"
	testCcipObjectRef    = "0x000000000000000000000000000000000000000000000000000000000bj3c7r3"
	testOffRampState     = "0x000000000000000000000000000000000000000000000000000000000ff2a357"
	testCcipOwnerCap     = "0x0000000000000000000000000000000000000000000000000000000000ca9ca9"
)

func testAddressMappings() *OffRampAddressMappings {
	return &OffRampAddressMappings{
		CcipPackageId:    testCcipPackageId,
		CcipObjectRef:    testCcipObjectRef,
		CcipOwnerCap:     testCcipOwnerCap,
		OffRampPackageId: testOffRampPackageId,
		OffRampState:     testOffRampState,
	}
}

func standardParams(ccipPkgId string) []any {
	return []any{
		map[string]any{"Vector": "U8"},
		map[string]any{
			"Reference": map[string]any{
				"Struct": map[string]any{
					"address":       ccipPkgId,
					"module":        "state_object",
					"name":          "CCIPObjectRef",
					"typeArguments": []any{},
				},
			},
		},
		map[string]any{
			"Struct": map[string]any{
				"address":       ccipPkgId,
				"module":        "client",
				"name":          "Any2SuiMessage",
				"typeArguments": []any{},
			},
		},
	}
}

func TestValidateReceiverCallbackSignature_StandardParams(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	params := standardParams(testCcipPackageId)
	funcSig := map[string]any{"parameters": params}
	decodedTypes := []string{"vector<u8>", "&object", "object_id"}

	err := ValidateReceiverCallbackSignature(lggr, funcSig, decodedTypes, testCcipPackageId, testOffRampPackageId)
	require.NoError(t, err)
}

func TestValidateReceiverCallbackSignature_LegitExtraParams(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	params := append(standardParams(testCcipPackageId),
		map[string]any{
			"Reference": map[string]any{
				"Struct": map[string]any{
					"address":       "0x2",
					"module":        "clock",
					"name":          "Clock",
					"typeArguments": []any{},
				},
			},
		},
		map[string]any{
			"MutableReference": map[string]any{
				"Struct": map[string]any{
					"address":       "0xdeadbeef",
					"module":        "my_receiver",
					"name":          "ReceiverState",
					"typeArguments": []any{},
				},
			},
		},
	)
	funcSig := map[string]any{"parameters": params}
	decodedTypes := []string{"vector<u8>", "&object", "object_id", "&object", "&mut object"}

	err := ValidateReceiverCallbackSignature(lggr, funcSig, decodedTypes, testCcipPackageId, testOffRampPackageId)
	require.NoError(t, err, "legitimate extra params (Clock + receiver's own state) should pass")
}

func TestValidateReceiverCallbackSignature_RejectsMutableCcipProtocolType(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	params := append(standardParams(testCcipPackageId),
		map[string]any{
			"MutableReference": map[string]any{
				"Struct": map[string]any{
					"address":       testCcipPackageId,
					"module":        "fee_quoter",
					"name":          "FeeQuoterState",
					"typeArguments": []any{},
				},
			},
		},
	)
	funcSig := map[string]any{"parameters": params}
	decodedTypes := []string{"vector<u8>", "&object", "object_id", "&mut object"}

	err := ValidateReceiverCallbackSignature(lggr, funcSig, decodedTypes, testCcipPackageId, testOffRampPackageId)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mutable reference to CCIP protocol type")
	assert.Contains(t, err.Error(), "FeeQuoterState")
}

func TestValidateReceiverCallbackSignature_RejectsMutableOnRampState(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	onrampPackageId := "0x0000000000000000000000000000000000000000000000000000000000012345"
	params := append(standardParams(testCcipPackageId),
		map[string]any{
			"MutableReference": map[string]any{
				"Struct": map[string]any{
					"address":       onrampPackageId,
					"module":        "onramp",
					"name":          "OnRampState",
					"typeArguments": []any{},
				},
			},
		},
	)
	funcSig := map[string]any{"parameters": params}
	decodedTypes := []string{"vector<u8>", "&object", "object_id", "&mut object"}

	err := ValidateReceiverCallbackSignature(lggr, funcSig, decodedTypes, testCcipPackageId, testOffRampPackageId)
	require.Error(t, err, "OnRampState should be caught by module name denylist even when package ID is unknown")
	assert.Contains(t, err.Error(), "denied protocol type")
	assert.Contains(t, err.Error(), "OnRampState")
}

func TestValidateReceiverCallbackSignature_RejectsMutableOffRampPackageType(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	params := append(standardParams(testCcipPackageId),
		map[string]any{
			"MutableReference": map[string]any{
				"Struct": map[string]any{
					"address":       testOffRampPackageId,
					"module":        "offramp",
					"name":          "OffRampState",
					"typeArguments": []any{},
				},
			},
		},
	)
	funcSig := map[string]any{"parameters": params}
	decodedTypes := []string{"vector<u8>", "&object", "object_id", "&mut object"}

	err := ValidateReceiverCallbackSignature(lggr, funcSig, decodedTypes, testCcipPackageId, testOffRampPackageId)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "CCIP protocol type")
}

func TestValidateReceiverCallbackSignature_TooFewParams(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	funcSig := map[string]any{"parameters": []any{map[string]any{"Vector": "U8"}}}
	decodedTypes := []string{"vector<u8>"}

	err := ValidateReceiverCallbackSignature(lggr, funcSig, decodedTypes, testCcipPackageId, testOffRampPackageId)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "expected at least 3 standard parameters")
}

func TestValidateReceiverCallbackSignature_TxContextSkipped(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	params := append(standardParams(testCcipPackageId),
		map[string]any{
			"MutableReference": map[string]any{
				"Struct": map[string]any{
					"address":       "0x2",
					"module":        "tx_context",
					"name":          "TxContext",
					"typeArguments": []any{},
				},
			},
		},
	)
	funcSig := map[string]any{"parameters": params}
	decodedTypes := []string{"vector<u8>", "&object", "object_id"}

	err := ValidateReceiverCallbackSignature(lggr, funcSig, decodedTypes, testCcipPackageId, testOffRampPackageId)
	require.NoError(t, err, "TxContext parameter should be skipped")
}

func TestValidateReceiverObjectIdCount_Matching(t *testing.T) {
	t.Parallel()

	decodedTypes := []string{"vector<u8>", "&object", "object_id", "&object", "&mut object"}
	err := ValidateReceiverObjectIdCount(decodedTypes, 2)
	require.NoError(t, err)
}

func TestValidateReceiverObjectIdCount_ExactlyStandard(t *testing.T) {
	t.Parallel()

	decodedTypes := []string{"vector<u8>", "&object", "object_id"}
	err := ValidateReceiverObjectIdCount(decodedTypes, 0)
	require.NoError(t, err)
}

func TestValidateReceiverObjectIdCount_Mismatch_TooMany(t *testing.T) {
	t.Parallel()

	decodedTypes := []string{"vector<u8>", "&object", "object_id"}
	err := ValidateReceiverObjectIdCount(decodedTypes, 2)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "declares 0 extra object parameters but receiverObjectIds contains 2")
}

func TestValidateReceiverObjectIdCount_Mismatch_TooFew(t *testing.T) {
	t.Parallel()

	decodedTypes := []string{"vector<u8>", "&object", "object_id", "&mut object", "&mut object"}
	err := ValidateReceiverObjectIdCount(decodedTypes, 1)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "declares 2 extra object parameters but receiverObjectIds contains 1")
}

func TestValidateReceiverObjectIds_Safe(t *testing.T) {
	t.Parallel()

	objectIds := []string{
		"0x0000000000000000000000000000000000000000000000000000000000aaaaaa",
		"0x0000000000000000000000000000000000000000000000000000000000bbbbbb",
	}
	err := ValidateReceiverObjectIds(objectIds, testAddressMappings())
	require.NoError(t, err)
}

func TestValidateReceiverObjectIds_RejectsCcipObjectRef(t *testing.T) {
	t.Parallel()

	objectIds := []string{testCcipObjectRef}
	err := ValidateReceiverObjectIds(objectIds, testAddressMappings())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "CCIPObjectRef")
}

func TestValidateReceiverObjectIds_RejectsOffRampState(t *testing.T) {
	t.Parallel()

	objectIds := []string{testOffRampState}
	err := ValidateReceiverObjectIds(objectIds, testAddressMappings())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "OffRampState")
}

func TestValidateReceiverObjectIds_RejectsCcipOwnerCap(t *testing.T) {
	t.Parallel()

	objectIds := []string{testCcipOwnerCap}
	err := ValidateReceiverObjectIds(objectIds, testAddressMappings())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "CcipOwnerCap")
}

func TestValidateReceiverObjectIds_EmptyList(t *testing.T) {
	t.Parallel()

	err := ValidateReceiverObjectIds(nil, testAddressMappings())
	require.NoError(t, err)
}

func TestIsDeniedProtocolPackage(t *testing.T) {
	t.Parallel()

	assert.True(t, isDeniedProtocolPackage(testCcipPackageId, testCcipPackageId, testOffRampPackageId))
	assert.True(t, isDeniedProtocolPackage(testOffRampPackageId, testCcipPackageId, testOffRampPackageId))
	assert.False(t, isDeniedProtocolPackage("0xdeadbeef", testCcipPackageId, testOffRampPackageId))
	assert.False(t, isDeniedProtocolPackage("", testCcipPackageId, testOffRampPackageId))
}

func TestIsDeniedProtocolModule(t *testing.T) {
	t.Parallel()

	assert.True(t, isDeniedProtocolModule("onramp", "OnRampState"))
	assert.True(t, isDeniedProtocolModule("offramp", "OffRampState"))
	assert.True(t, isDeniedProtocolModule("fee_quoter", "FeeQuoterState"))
	assert.True(t, isDeniedProtocolModule("state_object", "CCIPObjectRef"))
	assert.False(t, isDeniedProtocolModule("my_receiver", "ReceiverState"))
	assert.False(t, isDeniedProtocolModule("onramp", "SomeOtherType"))
	assert.False(t, isDeniedProtocolModule("clock", "Clock"))
}

func TestValidateReceiverCallbackSignature_ImmutableCcipRefAllowed(t *testing.T) {
	t.Parallel()
	lggr := logger.Test(t)

	// Immutable reference to a CCIP type as an extra param should be allowed
	// (read-only access is not dangerous in the same way mutable access is).
	params := append(standardParams(testCcipPackageId),
		map[string]any{
			"Reference": map[string]any{
				"Struct": map[string]any{
					"address":       testCcipPackageId,
					"module":        "state_object",
					"name":          "CCIPObjectRef",
					"typeArguments": []any{},
				},
			},
		},
	)
	funcSig := map[string]any{"parameters": params}
	decodedTypes := []string{"vector<u8>", "&object", "object_id", "&object"}

	err := ValidateReceiverCallbackSignature(lggr, funcSig, decodedTypes, testCcipPackageId, testOffRampPackageId)
	require.NoError(t, err, "immutable references are safe; only mutable references to protocol types are denied")
}
