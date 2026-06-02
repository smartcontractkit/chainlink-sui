package rmn

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_rmn_remote "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/rmn_remote"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

func testBundle(t *testing.T) cld_ops.Bundle {
	t.Helper()
	return cld_ops.NewBundle(
		func() context.Context { return t.Context() },
		logger.Test(t),
		cld_ops.NewMemoryReporter(),
	)
}

func TestSerializeMcmsObjectAddrs(t *testing.T) {
	t.Parallel()

	data, err := SerializeMcmsObjectAddrs(
		"0x1",
		"0x2",
		"0x3",
	)
	require.NoError(t, err)
	require.Len(t, data, 96)
}

func TestSerializeSubjects(t *testing.T) {
	t.Parallel()

	subject := make([]byte, 16)
	subject[15] = 1

	data, err := SerializeSubjects([][]byte{subject, subject})
	require.NoError(t, err)
	require.NotEmpty(t, data)
}

func TestMcmsMintAndRegisterCurserCapOp_EncodesProposalLeaf(t *testing.T) {
	t.Parallel()

	bundle := testBundle(t)
	deps := sui_ops.OpTxDeps{Signer: nil}

	const (
		ccipPkg  = "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		ref      = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
		ownerCap = "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
		fastReg  = "0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"
	)

	report, err := cld_ops.ExecuteOperation(
		bundle,
		McmsMintAndRegisterCurserCapOp,
		deps,
		McmsMintAndRegisterCurserCapInput{
			CCIPPackageId:        ccipPkg,
			StateObjectId:        ref,
			SlowOwnerCapObjectId: ownerCap,
			FastRegistryObjectId: fastReg,
		},
	)
	require.NoError(t, err)
	require.Equal(t, "rmn_remote", report.Output.Call.Module)
	require.Equal(t, "mint_and_register_curser_cap", report.Output.Call.Function)
	require.Equal(t, ref, report.Output.Call.StateObjID)
	require.Len(t, report.Output.Call.Data, 96)
}

func TestMcmsMintAndRegisterCurserCapOp_RejectsDirectSigner(t *testing.T) {
	t.Parallel()

	bundle := testBundle(t)
	_, err := cld_ops.ExecuteOperation(
		bundle,
		McmsMintAndRegisterCurserCapOp,
		sui_ops.OpTxDeps{Signer: stubSigner{}},
		McmsMintAndRegisterCurserCapInput{
			CCIPPackageId:        "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			StateObjectId:        "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			SlowOwnerCapObjectId: "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
			FastRegistryObjectId: "0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd",
		},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "slow MCMS proposal")
}

type stubSigner struct{}

func (stubSigner) Sign([]byte) ([]string, error) { return nil, nil }
func (stubSigner) GetAddress() (string, error)   { return "0x1", nil }

func TestCurseWithCurserCapOp_EncodesProposalLeaf(t *testing.T) {
	t.Parallel()

	bundle := testBundle(t)
	deps := sui_ops.OpTxDeps{Signer: nil}

	const (
		ccipPkg   = "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		ref       = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
		curserCap = "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
	)

	subject := make([]byte, 16)
	subject[15] = 0xab

	report, err := cld_ops.ExecuteOperation(
		bundle,
		CurseWithCurserCapOp,
		deps,
		CurseWithCurserCapInput{
			CCIPPackageId:     ccipPkg,
			StateObjectId:     ref,
			CurserCapObjectId: curserCap,
			Subjects:          [][]byte{subject},
		},
	)
	require.NoError(t, err)
	require.Equal(t, "curse_multiple_with_curser_cap", report.Output.Call.Function)
	require.Greater(t, len(report.Output.Call.Data), 64)
}

func TestCurseWithCurserCapOp_ProposalDataMatchesManualBCS(t *testing.T) {
	t.Parallel()

	const (
		ccipPkg   = "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		ref       = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
		curserCap = "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
	)

	subject := make([]byte, 16)
	subject[15] = 0xab
	subjects := [][]byte{subject}

	contract, err := module_rmn_remote.NewRmnRemote(ccipPkg, nil)
	require.NoError(t, err)

	encodedCall, err := contract.Encoder().CurseMultipleWithCurserCap(
		bind.Object{Id: ref},
		bind.Object{Id: curserCap},
		subjects,
	)
	require.NoError(t, err)

	call, err := sui_ops.ToTransactionCall(encodedCall, ref)
	require.NoError(t, err)

	addrData, err := SerializeMcmsObjectAddrs(ref, curserCap)
	require.NoError(t, err)
	subjectData, err := SerializeSubjects(subjects)
	require.NoError(t, err)

	require.Equal(t, append(addrData, subjectData...), call.Data)
}

func TestCurseWithCurserCapOp_RequiresSubjects(t *testing.T) {
	t.Parallel()

	bundle := testBundle(t)
	_, err := cld_ops.ExecuteOperation(
		bundle,
		CurseWithCurserCapOp,
		sui_ops.OpTxDeps{},
		CurseWithCurserCapInput{
			CCIPPackageId:     "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			StateObjectId:     "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			CurserCapObjectId: "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
		},
	)
	require.Error(t, err)
}

func TestCreateCurserCapOp_EncodesProposalLeaf(t *testing.T) {
	t.Parallel()

	bundle := testBundle(t)
	deps := sui_ops.OpTxDeps{Signer: nil}

	const (
		ccipPkg  = "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		ref      = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
		ownerCap = "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
	)

	report, err := cld_ops.ExecuteOperation(
		bundle,
		CreateCurserCapOp,
		deps,
		CreateCurserCapInput{
			CCIPPackageId:    ccipPkg,
			StateObjectId:    ref,
			OwnerCapObjectId: ownerCap,
		},
	)
	require.NoError(t, err)
	require.Equal(t, "create_curser_cap", report.Output.Call.Function)
	require.Len(t, report.Output.Call.Data, 64)
}
