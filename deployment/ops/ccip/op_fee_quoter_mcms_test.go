package ccipops

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

func feeQuoterTestBundle(t *testing.T) cld_ops.Bundle {
	t.Helper()
	return cld_ops.NewBundle(
		func() context.Context { return t.Context() },
		logger.Test(t),
		cld_ops.NewMemoryReporter(),
	)
}

func TestFeeQuoterApplyFeeTokenUpdatesOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		feeQuoterTestBundle(t),
		FeeQuoterApplyFeeTokenUpdatesOp,
		sui_ops.OpTxDeps{},
		FeeQuoterApplyFeeTokenUpdatesInput{
			CCIPPackageId:     testCCIPPackageID,
			StateObjectId:     testStateObjectID,
			OwnerCapObjectId:  testOwnerCapID,
			FeeTokensToRemove: []string{},
			FeeTokensToAdd:    []string{testCoinMetadata},
		},
	)
	require.NoError(t, err)

	contract, err := module_fee_quoter.NewFeeQuoter(testCCIPPackageID, nil)
	require.NoError(t, err)
	encodedCall, err := contract.Encoder().ApplyFeeTokenUpdates(
		bind.Object{Id: testStateObjectID},
		bind.Object{Id: testOwnerCapID},
		[]string{},
		[]string{testCoinMetadata},
	)
	require.NoError(t, err)
	expected, err := sui_ops.ToTransactionCall(encodedCall, testStateObjectID)
	require.NoError(t, err)
	require.Equal(t, expected.Data, report.Output.Call.Data)
}

func TestFeeQuoterApplyTokenTransferFeeConfigUpdatesOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	const destChainSelector uint64 = 16015286601757825753

	report, err := cld_ops.ExecuteOperation(
		feeQuoterTestBundle(t),
		FeeQuoterApplyTokenTransferFeeConfigUpdatesOp,
		sui_ops.OpTxDeps{},
		FeeQuoterApplyTokenTransferFeeConfigUpdatesInput{
			CCIPPackageId:        testCCIPPackageID,
			StateObjectId:        testStateObjectID,
			OwnerCapObjectId:     testOwnerCapID,
			DestChainSelector:    destChainSelector,
			AddTokens:            []string{testCoinMetadata},
			AddMinFeeUsdCents:    []uint32{1},
			AddMaxFeeUsdCents:    []uint32{2},
			AddDeciBps:           []uint16{3},
			AddDestGasOverhead:   []uint32{4},
			AddDestBytesOverhead: []uint32{5},
			AddIsEnabled:         []bool{true},
			RemoveTokens:         []string{},
		},
	)
	require.NoError(t, err)

	contract, err := module_fee_quoter.NewFeeQuoter(testCCIPPackageID, nil)
	require.NoError(t, err)
	encodedCall, err := contract.Encoder().ApplyTokenTransferFeeConfigUpdates(
		bind.Object{Id: testStateObjectID},
		bind.Object{Id: testOwnerCapID},
		destChainSelector,
		[]string{testCoinMetadata},
		[]uint32{1},
		[]uint32{2},
		[]uint16{3},
		[]uint32{4},
		[]uint32{5},
		[]bool{true},
		[]string{},
	)
	require.NoError(t, err)
	expected, err := sui_ops.ToTransactionCall(encodedCall, testStateObjectID)
	require.NoError(t, err)
	require.Equal(t, expected.Data, report.Output.Call.Data)
}

func TestFeeQuoterApplyDestChainConfigUpdatesOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	const destChainSelector uint64 = 16015286601757825753

	report, err := cld_ops.ExecuteOperation(
		feeQuoterTestBundle(t),
		FeeQuoterApplyDestChainConfigUpdatesOp,
		sui_ops.OpTxDeps{},
		FeeQuoterApplyDestChainConfigUpdatesInput{
			CCIPPackageId:                     testCCIPPackageID,
			StateObjectId:                     testStateObjectID,
			OwnerCapObjectId:                  testOwnerCapID,
			DestChainSelector:                 destChainSelector,
			IsEnabled:                         true,
			MaxNumberOfTokensPerMsg:           2,
			MaxDataBytes:                      100,
			MaxPerMsgGasLimit:                 200,
			DestGasOverhead:                   300,
			DestGasPerPayloadByteBase:         1,
			DestGasPerPayloadByteHigh:         2,
			DestGasPerPayloadByteThreshold:    3,
			DestDataAvailabilityOverheadGas:   400,
			DestGasPerDataAvailabilityByte:    5,
			DestDataAvailabilityMultiplierBps: 6,
			ChainFamilySelector:               []byte{0x28, 0x12, 0xd5, 0x2c},
			EnforceOutOfOrder:                 false,
			DefaultTokenFeeUsdCents:           7,
			DefaultTokenDestGasOverhead:       800,
			DefaultTxGasLimit:                 900,
			GasMultiplierWeiPerEth:            100,
			GasPriceStalenessThreshold:        200,
			NetworkFeeUsdCents:                10,
		},
	)
	require.NoError(t, err)

	contract, err := module_fee_quoter.NewFeeQuoter(testCCIPPackageID, nil)
	require.NoError(t, err)
	encodedCall, err := contract.Encoder().ApplyDestChainConfigUpdates(
		bind.Object{Id: testStateObjectID},
		bind.Object{Id: testOwnerCapID},
		destChainSelector,
		true,
		uint16(2),
		uint32(100),
		uint32(200),
		uint32(300),
		byte(1),
		byte(2),
		uint16(3),
		uint32(400),
		uint16(5),
		uint16(6),
		[]byte{0x28, 0x12, 0xd5, 0x2c},
		false,
		uint16(7),
		uint32(800),
		uint32(900),
		uint64(100),
		uint32(200),
		uint32(10),
	)
	require.NoError(t, err)
	expected, err := sui_ops.ToTransactionCall(encodedCall, testStateObjectID)
	require.NoError(t, err)
	require.Equal(t, expected.Data, report.Output.Call.Data)
}

func TestFeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	report, err := cld_ops.ExecuteOperation(
		feeQuoterTestBundle(t),
		FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp,
		sui_ops.OpTxDeps{},
		FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput{
			CCIPPackageId:              testCCIPPackageID,
			StateObjectId:              testStateObjectID,
			OwnerCapObjectId:           testOwnerCapID,
			Tokens:                     []string{testCoinMetadata},
			PremiumMultiplierWeiPerEth: []uint64{77},
		},
	)
	require.NoError(t, err)

	contract, err := module_fee_quoter.NewFeeQuoter(testCCIPPackageID, nil)
	require.NoError(t, err)
	encodedCall, err := contract.Encoder().ApplyPremiumMultiplierWeiPerEthUpdates(
		bind.Object{Id: testStateObjectID},
		bind.Object{Id: testOwnerCapID},
		[]string{testCoinMetadata},
		[]uint64{77},
	)
	require.NoError(t, err)
	expected, err := sui_ops.ToTransactionCall(encodedCall, testStateObjectID)
	require.NoError(t, err)
	require.Equal(t, expected.Data, report.Output.Call.Data)
}
