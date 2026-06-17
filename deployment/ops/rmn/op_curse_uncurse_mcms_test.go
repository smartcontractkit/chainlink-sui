package rmn

import (
	"testing"

	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_rmn_remote "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/rmn_remote"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/ops/mcmstest"
)

func TestRmnDualModeOps_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	contract, err := module_rmn_remote.NewRmnRemote(mcmstest.PackageID, nil)
	require.NoError(t, err)
	subject := testSubject(0xab)

	t.Run("curse_chain", func(t *testing.T) {
		t.Parallel()
		input := CurseUncurseChainInput{
			CCIPPackageId:    mcmstest.PackageID,
			StateObjectId:    mcmstest.StateObjectID,
			OwnerCapObjectId: mcmstest.OwnerCapID,
			Subjects:         [][]byte{subject},
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), CurseChainOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := contract.Encoder().CurseMultiple(bind.Object{Id: input.StateObjectId}, bind.Object{Id: input.OwnerCapObjectId}, input.Subjects)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.StateObjectId, nil)
	})

	t.Run("uncurse_chain", func(t *testing.T) {
		t.Parallel()
		input := CurseUncurseChainInput{
			CCIPPackageId:    mcmstest.PackageID,
			StateObjectId:    mcmstest.StateObjectID,
			OwnerCapObjectId: mcmstest.OwnerCapID,
			Subjects:         [][]byte{subject},
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), UncurseChainOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := contract.Encoder().UncurseMultiple(bind.Object{Id: input.StateObjectId}, bind.Object{Id: input.OwnerCapObjectId}, input.Subjects)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.StateObjectId, nil)
	})
}

func TestMcmsCreateCurserCapAndTransferOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		mcmstest.Bundle(t),
		McmsCreateCurserCapAndTransferOp,
		sui_ops.OpTxDeps{},
		McmsCreateCurserCapAndTransferInput{
			CCIPPackageId:        mcmstest.PackageID,
			StateObjectId:        mcmstest.StateObjectID,
			SlowOwnerCapObjectId: mcmstest.OwnerCapID,
			RecipientAddress:     mcmstest.Recipient,
		},
	)
	require.NoError(t, err)

	contract, err := module_rmn_remote.NewRmnRemote(mcmstest.PackageID, nil)
	require.NoError(t, err)
	encoded, err := contract.Encoder().CreateCurserCapAndTransfer(bind.Object{Id: mcmstest.StateObjectID}, bind.Object{Id: mcmstest.OwnerCapID}, mcmstest.Recipient)
	require.NoError(t, err)
	mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, mcmstest.StateObjectID, nil)
}

func TestMcmsInitializeAllowedCurserCapsOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		mcmstest.Bundle(t),
		McmsInitializeAllowedCurserCapsOp,
		sui_ops.OpTxDeps{},
		McmsInitializeAllowedCurserCapsInput{
			CCIPPackageId:        mcmstest.PackageID,
			StateObjectId:        mcmstest.StateObjectID,
			SlowOwnerCapObjectId: mcmstest.OwnerCapID,
			InitialCurserCapIds:  []string{mcmstest.RegistryID},
		},
	)
	require.NoError(t, err)

	contract, err := module_rmn_remote.NewRmnRemote(mcmstest.PackageID, nil)
	require.NoError(t, err)
	encoded, err := contract.Encoder().InitializeAllowedCurserCaps(bind.Object{Id: mcmstest.StateObjectID}, bind.Object{Id: mcmstest.OwnerCapID}, []string{mcmstest.RegistryID})
	require.NoError(t, err)
	mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, mcmstest.StateObjectID, nil)
}

func TestMcmsRegisterCurserCapIdsOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		mcmstest.Bundle(t),
		McmsRegisterCurserCapIdsOp,
		sui_ops.OpTxDeps{},
		McmsRegisterCurserCapIdsInput{
			CCIPPackageId:        mcmstest.PackageID,
			StateObjectId:        mcmstest.StateObjectID,
			SlowOwnerCapObjectId: mcmstest.OwnerCapID,
			CurserCapObjectIds:   []string{mcmstest.RegistryID},
		},
	)
	require.NoError(t, err)

	contract, err := module_rmn_remote.NewRmnRemote(mcmstest.PackageID, nil)
	require.NoError(t, err)
	encoded, err := contract.Encoder().RegisterCurserCapIds(bind.Object{Id: mcmstest.StateObjectID}, bind.Object{Id: mcmstest.OwnerCapID}, []string{mcmstest.RegistryID})
	require.NoError(t, err)
	mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, mcmstest.StateObjectID, nil)
}

func TestMcmsDeregisterCurserCapIdsOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		mcmstest.Bundle(t),
		McmsDeregisterCurserCapIdsOp,
		sui_ops.OpTxDeps{},
		McmsDeregisterCurserCapIdsInput{
			CCIPPackageId:        mcmstest.PackageID,
			StateObjectId:        mcmstest.StateObjectID,
			SlowOwnerCapObjectId: mcmstest.OwnerCapID,
			CurserCapObjectIds:   []string{mcmstest.RegistryID},
		},
	)
	require.NoError(t, err)

	contract, err := module_rmn_remote.NewRmnRemote(mcmstest.PackageID, nil)
	require.NoError(t, err)
	encoded, err := contract.Encoder().DeregisterCurserCapIds(bind.Object{Id: mcmstest.StateObjectID}, bind.Object{Id: mcmstest.OwnerCapID}, []string{mcmstest.RegistryID})
	require.NoError(t, err)
	mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, mcmstest.StateObjectID, nil)
}

func TestCreateCurserCapOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		mcmstest.Bundle(t),
		CreateCurserCapOp,
		sui_ops.OpTxDeps{},
		CreateCurserCapInput{
			CCIPPackageId:    mcmstest.PackageID,
			StateObjectId:    mcmstest.StateObjectID,
			OwnerCapObjectId: mcmstest.OwnerCapID,
		},
	)
	require.NoError(t, err)

	contract, err := module_rmn_remote.NewRmnRemote(mcmstest.PackageID, nil)
	require.NoError(t, err)
	encoded, err := contract.Encoder().CreateCurserCap(bind.Object{Id: mcmstest.StateObjectID}, bind.Object{Id: mcmstest.OwnerCapID})
	require.NoError(t, err)
	mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, mcmstest.StateObjectID, nil)
}

func TestCurseWithCurserCapOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	subject := testSubject(0xab)
	report, err := cld_ops.ExecuteOperation(
		mcmstest.Bundle(t),
		CurseWithCurserCapOp,
		sui_ops.OpTxDeps{},
		CurseWithCurserCapInput{
			CCIPPackageId:     mcmstest.PackageID,
			StateObjectId:     mcmstest.StateObjectID,
			CurserCapObjectId: mcmstest.RegistryID,
			Subjects:          [][]byte{subject},
		},
	)
	require.NoError(t, err)

	contract, err := module_rmn_remote.NewRmnRemote(mcmstest.PackageID, nil)
	require.NoError(t, err)
	encoded, err := contract.Encoder().CurseMultipleWithCurserCap(bind.Object{Id: mcmstest.StateObjectID}, bind.Object{Id: mcmstest.RegistryID}, [][]byte{subject})
	require.NoError(t, err)
	mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, mcmstest.StateObjectID, nil)
}
