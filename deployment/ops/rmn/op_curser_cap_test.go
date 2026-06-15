package rmn

import (
	"context"
	"testing"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_rmn_remote "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/rmn_remote"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

const (
	testCCIPPackageID = "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	testStateObjectID = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	testOwnerCapID    = "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"
	testCurserCapID   = "0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"
	testFastRegistry  = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
)

func testBundle(t *testing.T) cld_ops.Bundle {
	t.Helper()
	return cld_ops.NewBundle(
		func() context.Context { return t.Context() },
		logger.Test(t),
		cld_ops.NewMemoryReporter(),
	)
}

func testSubject(b byte) []byte {
	subject := make([]byte, 16)
	subject[15] = b
	return subject
}

func TestSerializeMcmsObjectAddrs(t *testing.T) {
	t.Parallel()

	data, err := SerializeMcmsObjectAddrs(testStateObjectID, testOwnerCapID, testFastRegistry)
	require.NoError(t, err)
	require.Len(t, data, 96)
}

func TestSerializeMcmsObjectAddrsWithAddressVector(t *testing.T) {
	t.Parallel()

	data, err := SerializeMcmsObjectAddrsWithAddressVector(
		[]string{testStateObjectID, testOwnerCapID},
		[]string{testCurserCapID},
	)
	require.NoError(t, err)
	require.Len(t, data, 64+1+32) // two pinned addrs + uleb128(1) + one address
}

func TestSerializeMcmsObjectAddrsWithBool(t *testing.T) {
	t.Parallel()

	data, err := SerializeMcmsObjectAddrsWithBool(
		[]string{testStateObjectID, testOwnerCapID},
		true,
	)
	require.NoError(t, err)
	require.Len(t, data, 65) // two pinned addrs + bool
	require.Equal(t, byte(1), data[len(data)-1])
}

func TestSerializeMcmsObjectAddrs_PadsShortIDs(t *testing.T) {
	t.Parallel()

	data, err := SerializeMcmsObjectAddrs("0x1", "0x2", "0x3")
	require.NoError(t, err)
	require.Len(t, data, 96)
}

func TestSerializeMcmsObjectAddrs_InvalidObjectID(t *testing.T) {
	t.Parallel()

	_, err := SerializeMcmsObjectAddrs("not-an-object-id")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid")
}

func TestSerializeMcmsObjectAddrs_InvalidHex(t *testing.T) {
	t.Parallel()

	_, err := SerializeMcmsObjectAddrs("0xzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid hex")
}

func TestSerializeSubjects(t *testing.T) {
	t.Parallel()

	subject := testSubject(1)
	data, err := SerializeSubjects([][]byte{subject, subject})
	require.NoError(t, err)
	require.NotEmpty(t, data)
}

func TestSerializeSubjects_Empty(t *testing.T) {
	t.Parallel()

	data, err := SerializeSubjects(nil)
	require.NoError(t, err)
	require.Equal(t, []byte{0}, data)
}

func TestSerializeSubjects_Multiple(t *testing.T) {
	t.Parallel()

	subjects := [][]byte{testSubject(1), testSubject(2), testSubject(3)}
	data, err := SerializeSubjects(subjects)
	require.NoError(t, err)
	require.Len(t, data, 1+3*(1+16)) // uleb128 count + 3 * (uleb128 len + 16-byte subject)
}

func TestFindCurserCapObjectIDFromTx(t *testing.T) {
	t.Parallel()

	tx := models.SuiTransactionBlockResponse{
		ObjectChanges: []models.ObjectChange{
			{
				Type:       "created",
				ObjectType: testCCIPPackageID + "::rmn_remote::CurserCap",
				ObjectId:   testCurserCapID,
			},
		},
	}

	found, err := bind.FindObjectIdFromPublishTx(tx, "rmn_remote", "CurserCap")
	require.NoError(t, err)
	require.Equal(t, testCurserCapID, found)
}

func TestFindCurserCapObjectIDFromTx_NotFound(t *testing.T) {
	t.Parallel()

	tx := models.SuiTransactionBlockResponse{
		ObjectChanges: []models.ObjectChange{
			{
				Type:       "created",
				ObjectType: testCCIPPackageID + "::rmn_remote::OwnerCap",
				ObjectId:   testOwnerCapID,
			},
		},
	}

	_, err := bind.FindObjectIdFromPublishTx(tx, "rmn_remote", "CurserCap")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")
}

func TestMcmsMintAndRegisterCurserCapOp_EncodesProposalLeaf(t *testing.T) {
	t.Parallel()

	bundle := testBundle(t)
	deps := sui_ops.OpTxDeps{Signer: nil}

	report, err := cld_ops.ExecuteOperation(
		bundle,
		McmsMintAndRegisterCurserCapOp,
		deps,
		McmsMintAndRegisterCurserCapInput{
			CCIPPackageId:        testCCIPPackageID,
			StateObjectId:        testStateObjectID,
			SlowOwnerCapObjectId: testOwnerCapID,
			FastRegistryObjectId: testFastRegistry,
		},
	)
	require.NoError(t, err)
	require.Equal(t, "rmn_remote", report.Output.Call.Module)
	require.Equal(t, "mint_and_register_curser_cap", report.Output.Call.Function)
	require.Equal(t, testStateObjectID, report.Output.Call.StateObjID)
	require.Len(t, report.Output.Call.Data, 96)
}

func TestMcmsMintAndRegisterCurserCapOp_ProposalDataMatchesManualBCS(t *testing.T) {
	t.Parallel()

	bundle := testBundle(t)
	report, err := cld_ops.ExecuteOperation(
		bundle,
		McmsMintAndRegisterCurserCapOp,
		sui_ops.OpTxDeps{Signer: nil},
		McmsMintAndRegisterCurserCapInput{
			CCIPPackageId:        testCCIPPackageID,
			StateObjectId:        testStateObjectID,
			SlowOwnerCapObjectId: testOwnerCapID,
			FastRegistryObjectId: testFastRegistry,
		},
	)
	require.NoError(t, err)

	expected, err := SerializeMcmsObjectAddrs(testStateObjectID, testOwnerCapID, testFastRegistry)
	require.NoError(t, err)
	require.Equal(t, expected, report.Output.Call.Data)
}

func TestMcmsCreateCurserCapAndTransferOp_EncodesProposalLeaf(t *testing.T) {
	t.Parallel()

	bundle := testBundle(t)
	report, err := cld_ops.ExecuteOperation(
		bundle,
		McmsCreateCurserCapAndTransferOp,
		sui_ops.OpTxDeps{Signer: nil},
		McmsCreateCurserCapAndTransferInput{
			CCIPPackageId:        testCCIPPackageID,
			StateObjectId:        testStateObjectID,
			SlowOwnerCapObjectId: testOwnerCapID,
			RecipientAddress:     "0xf45ca00000000000000000000000000000000000000000000000000000000000",
		},
	)
	require.NoError(t, err)
	require.Equal(t, "create_curser_cap_and_transfer", report.Output.Call.Function)
	require.Len(t, report.Output.Call.Data, 96)
}

func TestMcmsCreateCurserCapAndTransferOp_RequiresRecipient(t *testing.T) {
	t.Parallel()

	bundle := testBundle(t)
	_, err := cld_ops.ExecuteOperation(
		bundle,
		McmsCreateCurserCapAndTransferOp,
		sui_ops.OpTxDeps{Signer: nil},
		McmsCreateCurserCapAndTransferInput{
			CCIPPackageId:        testCCIPPackageID,
			StateObjectId:        testStateObjectID,
			SlowOwnerCapObjectId: testOwnerCapID,
		},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "recipient address")
}

func TestMcmsRegisterCurserCapOp_EncodesProposalLeaf(t *testing.T) {
	t.Parallel()

	bundle := testBundle(t)
	report, err := cld_ops.ExecuteOperation(
		bundle,
		McmsRegisterCurserCapOp,
		sui_ops.OpTxDeps{Signer: nil},
		McmsRegisterCurserCapInput{
			CCIPPackageId:        testCCIPPackageID,
			StateObjectId:        testStateObjectID,
			SlowOwnerCapObjectId: testOwnerCapID,
			FastRegistryObjectId: testFastRegistry,
			CurserCapObjectId:    testCurserCapID,
		},
	)
	require.NoError(t, err)
	require.Equal(t, "rmn_remote", report.Output.Call.Module)
	require.Equal(t, "register_curser_cap", report.Output.Call.Function)
	require.Equal(t, testStateObjectID, report.Output.Call.StateObjID)
	require.Len(t, report.Output.Call.Data, 128)
}

func TestMcmsRegisterCurserCapOp_ProposalDataMatchesManualBCS(t *testing.T) {
	t.Parallel()

	bundle := testBundle(t)
	report, err := cld_ops.ExecuteOperation(
		bundle,
		McmsRegisterCurserCapOp,
		sui_ops.OpTxDeps{Signer: nil},
		McmsRegisterCurserCapInput{
			CCIPPackageId:        testCCIPPackageID,
			StateObjectId:        testStateObjectID,
			SlowOwnerCapObjectId: testOwnerCapID,
			FastRegistryObjectId: testFastRegistry,
			CurserCapObjectId:    testCurserCapID,
		},
	)
	require.NoError(t, err)

	expected, err := SerializeMcmsObjectAddrs(testStateObjectID, testOwnerCapID, testFastRegistry, testCurserCapID)
	require.NoError(t, err)
	require.Equal(t, expected, report.Output.Call.Data)
}

func TestMcmsRegisterCurserCapOp_RequiresCurserCapObjectId(t *testing.T) {
	t.Parallel()

	bundle := testBundle(t)
	_, err := cld_ops.ExecuteOperation(
		bundle,
		McmsRegisterCurserCapOp,
		sui_ops.OpTxDeps{Signer: nil},
		McmsRegisterCurserCapInput{
			CCIPPackageId:        testCCIPPackageID,
			StateObjectId:        testStateObjectID,
			SlowOwnerCapObjectId: testOwnerCapID,
			FastRegistryObjectId: testFastRegistry,
		},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "curser cap object id")
}

func TestMcmsRegisterCurserCapOp_RejectsDirectSigner(t *testing.T) {
	t.Parallel()

	bundle := testBundle(t)
	_, err := cld_ops.ExecuteOperation(
		bundle,
		McmsRegisterCurserCapOp,
		sui_ops.OpTxDeps{Signer: stubSigner{}},
		McmsRegisterCurserCapInput{
			CCIPPackageId:        testCCIPPackageID,
			StateObjectId:        testStateObjectID,
			SlowOwnerCapObjectId: testOwnerCapID,
			FastRegistryObjectId: testFastRegistry,
			CurserCapObjectId:    testCurserCapID,
		},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "slow MCMS proposal")
}

func TestMcmsMintAndRegisterCurserCapOp_RejectsDirectSigner(t *testing.T) {
	t.Parallel()

	bundle := testBundle(t)
	_, err := cld_ops.ExecuteOperation(
		bundle,
		McmsMintAndRegisterCurserCapOp,
		sui_ops.OpTxDeps{Signer: stubSigner{}},
		McmsMintAndRegisterCurserCapInput{
			CCIPPackageId:        testCCIPPackageID,
			StateObjectId:        testStateObjectID,
			SlowOwnerCapObjectId: testOwnerCapID,
			FastRegistryObjectId: testFastRegistry,
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

	subject := testSubject(0xab)

	report, err := cld_ops.ExecuteOperation(
		bundle,
		CurseWithCurserCapOp,
		deps,
		CurseWithCurserCapInput{
			CCIPPackageId:     testCCIPPackageID,
			StateObjectId:     testStateObjectID,
			CurserCapObjectId: testCurserCapID,
			Subjects:          [][]byte{subject},
		},
	)
	require.NoError(t, err)
	require.Equal(t, "curse_multiple_with_curser_cap", report.Output.Call.Function)
	require.Greater(t, len(report.Output.Call.Data), 64)
	require.Empty(t, report.Output.Digest)
}

func TestCurseWithCurserCapOp_ProposalDataMatchesManualBCS(t *testing.T) {
	t.Parallel()

	subject := testSubject(0xab)
	subjects := [][]byte{subject}

	contract, err := module_rmn_remote.NewRmnRemote(testCCIPPackageID, nil)
	require.NoError(t, err)

	encodedCall, err := contract.Encoder().CurseMultipleWithCurserCap(
		bind.Object{Id: testStateObjectID},
		bind.Object{Id: testCurserCapID},
		subjects,
	)
	require.NoError(t, err)

	call, err := sui_ops.ToTransactionCall(encodedCall, testStateObjectID)
	require.NoError(t, err)

	addrData, err := SerializeMcmsObjectAddrs(testStateObjectID, testCurserCapID)
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
			CCIPPackageId:     testCCIPPackageID,
			StateObjectId:     testStateObjectID,
			CurserCapObjectId: testCurserCapID,
		},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "at least one subject")
}

func TestCreateCurserCapOp_EncodesProposalLeaf(t *testing.T) {
	t.Parallel()

	bundle := testBundle(t)
	deps := sui_ops.OpTxDeps{Signer: nil}

	report, err := cld_ops.ExecuteOperation(
		bundle,
		CreateCurserCapOp,
		deps,
		CreateCurserCapInput{
			CCIPPackageId:    testCCIPPackageID,
			StateObjectId:    testStateObjectID,
			OwnerCapObjectId: testOwnerCapID,
		},
	)
	require.NoError(t, err)
	require.Equal(t, "create_curser_cap", report.Output.Call.Function)
	require.Equal(t, testCCIPPackageID, report.Output.PackageId)
	require.Len(t, report.Output.Call.Data, 64)
	require.Empty(t, report.Output.Objects.CurserCapObjectId)
	require.Empty(t, report.Output.Digest)
}

func TestCreateCurserCapOp_ProposalDataMatchesManualBCS(t *testing.T) {
	t.Parallel()

	bundle := testBundle(t)
	report, err := cld_ops.ExecuteOperation(
		bundle,
		CreateCurserCapOp,
		sui_ops.OpTxDeps{Signer: nil},
		CreateCurserCapInput{
			CCIPPackageId:    testCCIPPackageID,
			StateObjectId:    testStateObjectID,
			OwnerCapObjectId: testOwnerCapID,
		},
	)
	require.NoError(t, err)

	expected, err := SerializeMcmsObjectAddrs(testStateObjectID, testOwnerCapID)
	require.NoError(t, err)
	require.Equal(t, expected, report.Output.Call.Data)
}

func TestMcmsInitializeAllowedCurserCapsOp_EncodesProposalLeaf(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		testBundle(t),
		McmsInitializeAllowedCurserCapsOp,
		sui_ops.OpTxDeps{},
		McmsInitializeAllowedCurserCapsInput{
			CCIPPackageId:        testCCIPPackageID,
			StateObjectId:        testStateObjectID,
			SlowOwnerCapObjectId: testOwnerCapID,
			InitialCurserCapIds:  []string{testCurserCapID},
		},
	)
	require.NoError(t, err)
	require.Equal(t, "initialize_allowed_curser_caps", report.Output.Call.Function)

	expected, err := SerializeMcmsObjectAddrsWithAddressVector(
		[]string{testStateObjectID, testOwnerCapID},
		[]string{testCurserCapID},
	)
	require.NoError(t, err)
	require.Equal(t, expected, report.Output.Call.Data)
}

func TestMcmsRegisterCurserCapIdsOp_EncodesProposalLeaf(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		testBundle(t),
		McmsRegisterCurserCapIdsOp,
		sui_ops.OpTxDeps{},
		McmsRegisterCurserCapIdsInput{
			CCIPPackageId:        testCCIPPackageID,
			StateObjectId:        testStateObjectID,
			SlowOwnerCapObjectId: testOwnerCapID,
			CurserCapObjectIds:   []string{testCurserCapID},
		},
	)
	require.NoError(t, err)
	require.Equal(t, "register_curser_cap_ids", report.Output.Call.Function)

	expected, err := SerializeMcmsObjectAddrsWithAddressVector(
		[]string{testStateObjectID, testOwnerCapID},
		[]string{testCurserCapID},
	)
	require.NoError(t, err)
	require.Equal(t, expected, report.Output.Call.Data)
}

func TestMcmsRegisterCurserCapIdsOp_RequiresCurserCapObjectIds(t *testing.T) {
	t.Parallel()

	_, err := cld_ops.ExecuteOperation(
		testBundle(t),
		McmsRegisterCurserCapIdsOp,
		sui_ops.OpTxDeps{},
		McmsRegisterCurserCapIdsInput{
			CCIPPackageId:        testCCIPPackageID,
			StateObjectId:        testStateObjectID,
			SlowOwnerCapObjectId: testOwnerCapID,
		},
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "at least one curser cap object id")
}

func TestMcmsSetCurserCapAllowlistEnabledOp_EncodesProposalLeaf(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		testBundle(t),
		McmsSetCurserCapAllowlistEnabledOp,
		sui_ops.OpTxDeps{},
		McmsSetCurserCapAllowlistEnabledInput{
			CCIPPackageId:        testCCIPPackageID,
			StateObjectId:        testStateObjectID,
			SlowOwnerCapObjectId: testOwnerCapID,
			Enabled:              true,
		},
	)
	require.NoError(t, err)
	require.Equal(t, "set_curser_cap_allowlist_enabled", report.Output.Call.Function)

	expected, err := SerializeMcmsObjectAddrsWithBool(
		[]string{testStateObjectID, testOwnerCapID},
		true,
	)
	require.NoError(t, err)
	require.Equal(t, expected, report.Output.Call.Data)
}

func TestMcmsDeregisterCurserCapIdsOp_EncodesProposalLeaf(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		testBundle(t),
		McmsDeregisterCurserCapIdsOp,
		sui_ops.OpTxDeps{},
		McmsDeregisterCurserCapIdsInput{
			CCIPPackageId:        testCCIPPackageID,
			StateObjectId:        testStateObjectID,
			SlowOwnerCapObjectId: testOwnerCapID,
			CurserCapObjectIds:   []string{testCurserCapID},
		},
	)
	require.NoError(t, err)
	require.Equal(t, "deregister_curser_cap_ids", report.Output.Call.Function)

	expected, err := SerializeMcmsObjectAddrsWithAddressVector(
		[]string{testStateObjectID, testOwnerCapID},
		[]string{testCurserCapID},
	)
	require.NoError(t, err)
	require.Equal(t, expected, report.Output.Call.Data)
}

func TestCurserCapOpDefinitions(t *testing.T) {
	t.Parallel()

	require.Equal(t, "sui-ccip-rmn_remote-create_curser_cap", CreateCurserCapOp.Def().ID)
	require.Equal(t, "sui-ccip-rmn_remote-mcms_mint_and_register_curser_cap", McmsMintAndRegisterCurserCapOp.Def().ID)
	require.Equal(t, "sui-ccip-rmn_remote-mcms_create_curser_cap_and_transfer", McmsCreateCurserCapAndTransferOp.Def().ID)
	require.Equal(t, "sui-ccip-rmn_remote-mcms_register_curser_cap", McmsRegisterCurserCapOp.Def().ID)
	require.Equal(t, "sui-ccip-rmn_remote-mcms_initialize_allowed_curser_caps", McmsInitializeAllowedCurserCapsOp.Def().ID)
	require.Equal(t, "sui-ccip-rmn_remote-mcms_register_curser_cap_ids", McmsRegisterCurserCapIdsOp.Def().ID)
	require.Equal(t, "sui-ccip-rmn_remote-mcms_set_curser_cap_allowlist_enabled", McmsSetCurserCapAllowlistEnabledOp.Def().ID)
	require.Equal(t, "sui-ccip-rmn_remote-mcms_deregister_curser_cap_ids", McmsDeregisterCurserCapIdsOp.Def().ID)
	require.Equal(t, "sui-ccip-rmn_remote-curse_with_curser_cap", CurseWithCurserCapOp.Def().ID)
}
