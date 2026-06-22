package ccipops

import (
	"context"
	"testing"

	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	module_token_admin_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/ops/rmn"
)

const (
	testCCIPPackageID = "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	testStateObjectID = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	testOwnerCapID    = "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
	testCoinMetadata  = "0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"
	testNewAdmin      = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
)

func testBundle(t *testing.T) cld_ops.Bundle {
	t.Helper()
	return cld_ops.NewBundle(
		func() context.Context { return t.Context() },
		logger.Test(t),
		cld_ops.NewMemoryReporter(),
	)
}

func TestTokenAdminRegistryInitializeLocalDecimalsOp_EncodesProposalLeaf(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		testBundle(t),
		TokenAdminRegistryInitializeLocalDecimalsOp,
		sui_ops.OpTxDeps{},
		InitLocalDecimalsInput{
			CCIPPackageId:    testCCIPPackageID,
			StateObjectId:    testStateObjectID,
			OwnerCapObjectId: testOwnerCapID,
		},
	)
	require.NoError(t, err)
	require.Equal(t, "token_admin_registry", report.Output.Call.Module)
	require.Equal(t, "initialize_local_decimals", report.Output.Call.Function)
	require.Equal(t, testStateObjectID, report.Output.Call.StateObjID)
	require.Len(t, report.Output.Call.Data, 64)
	require.Empty(t, report.Output.Digest)
}

func TestTokenAdminRegistryInitializeLocalDecimalsOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		testBundle(t),
		TokenAdminRegistryInitializeLocalDecimalsOp,
		sui_ops.OpTxDeps{},
		InitLocalDecimalsInput{
			CCIPPackageId:    testCCIPPackageID,
			StateObjectId:    testStateObjectID,
			OwnerCapObjectId: testOwnerCapID,
		},
	)
	require.NoError(t, err)

	contract, err := module_token_admin_registry.NewTokenAdminRegistry(testCCIPPackageID, nil)
	require.NoError(t, err)
	encodedCall, err := contract.Encoder().InitializeLocalDecimals(
		bind.Object{Id: testStateObjectID},
		bind.Object{Id: testOwnerCapID},
	)
	require.NoError(t, err)
	expected, err := sui_ops.ToTransactionCall(encodedCall, testStateObjectID)
	require.NoError(t, err)
	require.Equal(t, expected.Data, report.Output.Call.Data)
}

func TestTokenAdminRegistryBackfillLocalDecimalsOp_EncodesProposalLeaf(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		testBundle(t),
		TokenAdminRegistryBackfillLocalDecimalsOp,
		sui_ops.OpTxDeps{},
		BackfillLocalDecimalsInput{
			CCIPPackageId:       testCCIPPackageID,
			StateObjectId:       testStateObjectID,
			OwnerCapObjectId:    testOwnerCapID,
			CoinMetadataAddress: testCoinMetadata,
			LocalDecimals:       9,
		},
	)
	require.NoError(t, err)
	require.Equal(t, "token_admin_registry", report.Output.Call.Module)
	require.Equal(t, "backfill_local_decimals", report.Output.Call.Function)
	require.Equal(t, testStateObjectID, report.Output.Call.StateObjID)
	require.Len(t, report.Output.Call.Data, 97)
	require.Empty(t, report.Output.Digest)
}

func TestTokenAdminRegistryBackfillLocalDecimalsOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		testBundle(t),
		TokenAdminRegistryBackfillLocalDecimalsOp,
		sui_ops.OpTxDeps{},
		BackfillLocalDecimalsInput{
			CCIPPackageId:       testCCIPPackageID,
			StateObjectId:       testStateObjectID,
			OwnerCapObjectId:    testOwnerCapID,
			CoinMetadataAddress: testCoinMetadata,
			LocalDecimals:       6,
		},
	)
	require.NoError(t, err)

	contract, err := module_token_admin_registry.NewTokenAdminRegistry(testCCIPPackageID, nil)
	require.NoError(t, err)
	encodedCall, err := contract.Encoder().BackfillLocalDecimals(
		bind.Object{Id: testOwnerCapID},
		bind.Object{Id: testStateObjectID},
		testCoinMetadata,
		byte(6),
	)
	require.NoError(t, err)
	expected, err := sui_ops.ToTransactionCall(encodedCall, testStateObjectID)
	require.NoError(t, err)
	require.Equal(t, expected.Data, report.Output.Call.Data)

	d := bcs.NewDeserializer(report.Output.Call.Data)
	require.Equal(t, testOwnerCapID[2:], bytesToHex(d.ReadFixedBytes(32)))
	require.Equal(t, testStateObjectID[2:], bytesToHex(d.ReadFixedBytes(32)))
	require.Equal(t, testCoinMetadata[2:], bytesToHex(d.ReadFixedBytes(32)))
	require.Equal(t, uint8(6), d.U8())
}

func bytesToHex(b []byte) string {
	const hex = "0123456789abcdef"
	out := make([]byte, len(b)*2)
	for i, v := range b {
		out[i*2] = hex[v>>4]
		out[i*2+1] = hex[v&0x0f]
	}
	return string(out)
}

func TestTokenAdminRegistryUnregisterPoolOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		testBundle(t),
		TokenAdminRegistryUnregisterPoolOp,
		sui_ops.OpTxDeps{},
		UnregisterPoolInput{
			CCIPPackageId:       testCCIPPackageID,
			CCIPObjectRef:       testStateObjectID,
			OwnerCapObjectId:    testOwnerCapID,
			CoinMetadataAddress: testCoinMetadata,
		},
	)
	require.NoError(t, err)
	require.Equal(t, "token_admin_registry", report.Output.Call.Module)
	require.Equal(t, "unregister_pool", report.Output.Call.Function)
	require.Equal(t, testStateObjectID, report.Output.Call.StateObjID)

	expected, err := rmn.SerializeMcmsObjectAddrs(testOwnerCapID, testStateObjectID, testCoinMetadata)
	require.NoError(t, err)
	require.Equal(t, expected, report.Output.Call.Data)

	d := bcs.NewDeserializer(report.Output.Call.Data)
	require.Equal(t, testOwnerCapID[2:], bytesToHex(d.ReadFixedBytes(32)))
	require.Equal(t, testStateObjectID[2:], bytesToHex(d.ReadFixedBytes(32)))
	require.Equal(t, testCoinMetadata[2:], bytesToHex(d.ReadFixedBytes(32)))
}

func TestTokenAdminRegistryTransferAdminRoleOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		testBundle(t),
		TokenAdminRegistryTransferAdminRoleOp,
		sui_ops.OpTxDeps{},
		TransferAdminRoleInput{
			CCIPPackageId:       testCCIPPackageID,
			CCIPObjectRef:       testStateObjectID,
			CoinMetadataAddress: testCoinMetadata,
			NewAdmin:            testNewAdmin,
		},
	)
	require.NoError(t, err)

	contract, err := module_token_admin_registry.NewTokenAdminRegistry(testCCIPPackageID, nil)
	require.NoError(t, err)
	encodedCall, err := contract.Encoder().TransferAdminRole(
		bind.Object{Id: testStateObjectID},
		testCoinMetadata,
		testNewAdmin,
	)
	require.NoError(t, err)
	expected, err := sui_ops.ToTransactionCall(encodedCall, testStateObjectID)
	require.NoError(t, err)
	require.Equal(t, expected.Data, report.Output.Call.Data)
}

func TestTokenAdminRegistryAcceptAdminRoleOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		testBundle(t),
		TokenAdminRegistryAcceptAdminRoleOp,
		sui_ops.OpTxDeps{},
		AcceptAdminRoleInput{
			CCIPPackageId:       testCCIPPackageID,
			CCIPObjectRef:       testStateObjectID,
			CoinMetadataAddress: testCoinMetadata,
		},
	)
	require.NoError(t, err)

	contract, err := module_token_admin_registry.NewTokenAdminRegistry(testCCIPPackageID, nil)
	require.NoError(t, err)
	encodedCall, err := contract.Encoder().AcceptAdminRole(
		bind.Object{Id: testStateObjectID},
		testCoinMetadata,
	)
	require.NoError(t, err)
	expected, err := sui_ops.ToTransactionCall(encodedCall, testStateObjectID)
	require.NoError(t, err)
	require.Equal(t, expected.Data, report.Output.Call.Data)
}
