package ccipops

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_upgrade_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/upgrade_registry"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

const testModuleName = "fee_quoter"
const testFunctionName = "apply_fee_token_updates"

func upgradeRegistryTestBundle(t *testing.T) cld_ops.Bundle {
	t.Helper()
	return cld_ops.NewBundle(
		func() context.Context { return t.Context() },
		logger.Test(t),
		cld_ops.NewMemoryReporter(),
	)
}

func TestBlockVersionOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		upgradeRegistryTestBundle(t),
		BlockVersionOp,
		sui_ops.OpTxDeps{},
		BlockVersionInput{
			PackageId:        testCCIPPackageID,
			StateObjectId:    testStateObjectID,
			OwnerCapObjectId: testOwnerCapID,
			ModuleName:       testModuleName,
			Version:          1,
		},
	)
	require.NoError(t, err)

	contract, err := module_upgrade_registry.NewUpgradeRegistry(testCCIPPackageID, nil)
	require.NoError(t, err)
	encodedCall, err := contract.Encoder().BlockVersion(
		bind.Object{Id: testStateObjectID},
		bind.Object{Id: testOwnerCapID},
		testModuleName,
		byte(1),
	)
	require.NoError(t, err)
	expected, err := sui_ops.ToTransactionCall(encodedCall, testStateObjectID)
	require.NoError(t, err)
	require.Equal(t, expected.Data, report.Output.Call.Data)
}

func TestUnblockVersionOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		upgradeRegistryTestBundle(t),
		UnblockVersionOp,
		sui_ops.OpTxDeps{},
		UnblockVersionInput{
			PackageId:        testCCIPPackageID,
			StateObjectId:    testStateObjectID,
			OwnerCapObjectId: testOwnerCapID,
			ModuleName:       testModuleName,
			Version:          1,
		},
	)
	require.NoError(t, err)

	contract, err := module_upgrade_registry.NewUpgradeRegistry(testCCIPPackageID, nil)
	require.NoError(t, err)
	encodedCall, err := contract.Encoder().UnblockVersion(
		bind.Object{Id: testStateObjectID},
		bind.Object{Id: testOwnerCapID},
		testModuleName,
		byte(1),
	)
	require.NoError(t, err)
	expected, err := sui_ops.ToTransactionCall(encodedCall, testStateObjectID)
	require.NoError(t, err)
	require.Equal(t, expected.Data, report.Output.Call.Data)
}

func TestBlockFunctionOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		upgradeRegistryTestBundle(t),
		BlockFunctionOp,
		sui_ops.OpTxDeps{},
		BlockFunctionInput{
			PackageId:        testCCIPPackageID,
			StateObjectId:    testStateObjectID,
			OwnerCapObjectId: testOwnerCapID,
			ModuleName:       testModuleName,
			FunctionName:     testFunctionName,
			Version:          1,
		},
	)
	require.NoError(t, err)

	contract, err := module_upgrade_registry.NewUpgradeRegistry(testCCIPPackageID, nil)
	require.NoError(t, err)
	encodedCall, err := contract.Encoder().BlockFunction(
		bind.Object{Id: testStateObjectID},
		bind.Object{Id: testOwnerCapID},
		testModuleName,
		testFunctionName,
		byte(1),
	)
	require.NoError(t, err)
	expected, err := sui_ops.ToTransactionCall(encodedCall, testStateObjectID)
	require.NoError(t, err)
	require.Equal(t, expected.Data, report.Output.Call.Data)
}

func TestUnblockFunctionOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		upgradeRegistryTestBundle(t),
		UnblockFunctionOp,
		sui_ops.OpTxDeps{},
		UnblockFunctionInput{
			PackageId:        testCCIPPackageID,
			StateObjectId:    testStateObjectID,
			OwnerCapObjectId: testOwnerCapID,
			ModuleName:       testModuleName,
			FunctionName:     testFunctionName,
			Version:          1,
		},
	)
	require.NoError(t, err)

	contract, err := module_upgrade_registry.NewUpgradeRegistry(testCCIPPackageID, nil)
	require.NoError(t, err)
	encodedCall, err := contract.Encoder().UnblockFunction(
		bind.Object{Id: testStateObjectID},
		bind.Object{Id: testOwnerCapID},
		testModuleName,
		testFunctionName,
		byte(1),
	)
	require.NoError(t, err)
	expected, err := sui_ops.ToTransactionCall(encodedCall, testStateObjectID)
	require.NoError(t, err)
	require.Equal(t, expected.Data, report.Output.Call.Data)
}

func TestUpgradeRegistryInitializeOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		upgradeRegistryTestBundle(t),
		UpgradeRegistryInitializeOp,
		sui_ops.OpTxDeps{},
		InitUpgradeRegistryInput{
			CCIPPackageId:    testCCIPPackageID,
			StateObjectId:    testStateObjectID,
			OwnerCapObjectId: testOwnerCapID,
		},
	)
	require.NoError(t, err)

	contract, err := module_upgrade_registry.NewUpgradeRegistry(testCCIPPackageID, nil)
	require.NoError(t, err)
	encodedCall, err := contract.Encoder().Initialize(
		bind.Object{Id: testStateObjectID},
		bind.Object{Id: testOwnerCapID},
	)
	require.NoError(t, err)
	expected, err := sui_ops.ToTransactionCall(encodedCall, testStateObjectID)
	require.NoError(t, err)
	require.Equal(t, expected.Data, report.Output.Call.Data)
}
