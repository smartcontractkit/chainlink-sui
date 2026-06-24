package lanes_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	laneapi "github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-sui/deployment/lanes"
	"github.com/smartcontractkit/chainlink-sui/deployment/ops/mcmstest"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

const (
	evmSepoliaSelector = uint64(16015286601757825753)
	evmRouterAddress   = "0xd3E190f381f06DC0d289590fd452C42Fa2DAC586"

	testCCIPPackageID    = "0xece742a763bddf1e36629fa06b605497e413241afd14f05e558e80eef4f64e95"
	testOffRampPackageID = "0x9438693fb18f5660aff9277240a2282be44dc01cdd7eed4e1d8de0591ad52c03"
	testCCIPObjectRef    = "0xbeace36c3c1e1f37c5806c4954140d15cbc7b0002ef7ccb490de26e82f5ec4ca"
	testCCIPOwnerCapID   = "0x874447a3ac6ae37bf545c6d71b0fa4a3d0af56a3f339faf4a4bab16ca4956ce7"
	testLinkCoinMetadata = "0x8afb916ec72b91d28f519539659ebb1200b1824ff1f8d4c8f433acbb03017f2f"
)

func testSuiChains() chain.BlockChains {
	return chain.NewBlockChains(map[uint64]chain.BlockChain{
		suiTestnetSelector: sui.Chain{},
	})
}

func destLegInput() laneapi.UpdateLanesInput {
	return laneapi.UpdateLanesInput{
		Source: &laneapi.ChainDefinition{
			Selector:               evmSepoliaSelector,
			OnRamp:                 make([]byte, 32),
			RMNVerificationEnabled: true,
		},
		Dest: &laneapi.ChainDefinition{
			Selector: suiTestnetSelector,
		},
	}
}

func sourceLegInput(destRouter []byte) laneapi.UpdateLanesInput {
	return laneapi.UpdateLanesInput{
		Source: &laneapi.ChainDefinition{
			Selector: suiTestnetSelector,
		},
		Dest: &laneapi.ChainDefinition{
			Selector: evmSepoliaSelector,
			Router:   destRouter,
			FeeQuoterDestChainConfig: laneapi.FeeQuoterDestChainConfig{
				IsEnabled:               true,
				MaxDataBytes:            30_000,
				MaxPerMsgGasLimit:       3_000_000,
				ChainFamilySelector:     0x2812d52c,
				DefaultTokenFeeUSDCents: 25,
				NetworkFeeUSDCents:      10,
			},
		},
	}
}

func TestConfigureLaneLegAsDest_MCMSBatchOp(t *testing.T) {
	env := testEnvWithAddressBook(t)
	input := destLegInput()

	var report cldf_ops.SequenceReport[laneapi.UpdateLanesInput, sequences.OnChainOutput]
	err := lanes.WithConnectChainsEnvironment(env, func() error {
		var execErr error
		report, execErr = cldf_ops.ExecuteSequence(
			mcmstest.Bundle(t),
			lanes.ConfigureLaneLegAsDest,
			testSuiChains(),
			input,
		)
		return execErr
	})
	require.NoError(t, err)
	require.Len(t, report.Output.BatchOps, 1)
	require.Equal(t, suiTestnetSelector, uint64(report.Output.BatchOps[0].ChainSelector))
	require.Len(t, report.Output.BatchOps[0].Transactions, 1)

	tx := report.Output.BatchOps[0].Transactions[0]
	require.NotEmpty(t, tx.Data)
	require.Equal(t, testOffRampPackageID, tx.To)
	latestPackageID, err := utils.TransactionLatestPackageID(tx)
	require.NoError(t, err)
	require.Empty(t, latestPackageID)
}

func TestConfigureLaneLegAsDest_MCMSBatchOp_WithLatestPackageID(t *testing.T) {
	env := testEnvWithAddressBook(t)
	const latestOffRampPackageID = "0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"
	input := destLegInput()

	var report cldf_ops.SequenceReport[laneapi.UpdateLanesInput, sequences.OnChainOutput]
	err := lanes.WithConnectChainsEnvironment(env, func() error {
		return lanes.WithSuiLatestPackageIDs(map[uint64]lanes.LatestPackageIDsConfig{
			suiTestnetSelector: {OffRamp: latestOffRampPackageID},
		}, func() error {
			var execErr error
			report, execErr = cldf_ops.ExecuteSequence(
				mcmstest.Bundle(t),
				lanes.ConfigureLaneLegAsDest,
				testSuiChains(),
				input,
			)
			return execErr
		})
	})
	require.NoError(t, err)
	require.Len(t, report.Output.BatchOps, 1)
	tx := report.Output.BatchOps[0].Transactions[0]
	require.Equal(t, testOffRampPackageID, tx.To)
	latestPackageIDFromTx, err := utils.TransactionLatestPackageID(tx)
	require.NoError(t, err)
	require.Equal(t, latestOffRampPackageID, latestPackageIDFromTx)
}

func TestConfigureLaneLegAsDest_ValidationErrors(t *testing.T) {
	env := testEnvWithAddressBook(t)
	chains := testSuiChains()

	t.Run("requires ConnectChains scope", func(t *testing.T) {
		_, err := cldf_ops.ExecuteSequence(
			mcmstest.Bundle(t),
			lanes.ConfigureLaneLegAsDest,
			chains,
			destLegInput(),
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "WithConnectChainsEnvironment")
	})

	t.Run("requires source OnRamp", func(t *testing.T) {
		input := destLegInput()
		input.Source.OnRamp = nil
		err := lanes.WithConnectChainsEnvironment(env, func() error {
			_, execErr := cldf_ops.ExecuteSequence(
				mcmstest.Bundle(t),
				lanes.ConfigureLaneLegAsDest,
				chains,
				input,
			)
			return execErr
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "source OnRamp address required")
	})
}

func TestConfigureLaneLegAsSource_ValidationErrors(t *testing.T) {
	chains := testSuiChains()
	router := common.HexToAddress(evmRouterAddress).Bytes()

	t.Run("requires ConnectChains scope", func(t *testing.T) {
		_, err := cldf_ops.ExecuteSequence(
			mcmstest.Bundle(t),
			lanes.ConfigureLaneLegAsSource,
			chains,
			sourceLegInput(router),
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "WithConnectChainsEnvironment")
	})

	t.Run("requires dest Router", func(t *testing.T) {
		_, err := cldf_ops.ExecuteSequence(
			mcmstest.Bundle(t),
			lanes.ConfigureLaneLegAsSource,
			chains,
			sourceLegInput(nil),
		)
		require.Error(t, err)
		require.Contains(t, err.Error(), "dest Router address required")
	})
}

func TestConfigureLaneLegAsSource_RouterAddressBytes(t *testing.T) {
	env := testEnvWithAddressBook(t)
	chains := testSuiChains()
	evmAddr := common.HexToAddress(evmRouterAddress)

	t.Run("20 byte EVM address", func(t *testing.T) {
		_, err := runSourceLeg(t, env, chains, sourceLegInput(evmAddr.Bytes()))
		require.NoError(t, err)
	})

	t.Run("32 byte left padded EVM address", func(t *testing.T) {
		padded := make([]byte, 32)
		copy(padded[12:], evmAddr.Bytes())
		_, err := runSourceLeg(t, env, chains, sourceLegInput(padded))
		require.NoError(t, err)
	})

	t.Run("invalid router length", func(t *testing.T) {
		_, err := runSourceLeg(t, env, chains, sourceLegInput([]byte{0x01, 0x02, 0x03}))
		require.Error(t, err)
		require.Contains(t, err.Error(), "unexpected address length")
	})
}

func TestConfigureLaneLegAsSource_MCMSBatchOp_FirstOpEncoderParity(t *testing.T) {
	env := testEnvWithAddressBook(t)
	chains := testSuiChains()
	input := sourceLegInput(common.HexToAddress(evmRouterAddress).Bytes())

	report, err := runSourceLeg(t, env, chains, input)
	require.NoError(t, err)
	require.Len(t, report.Output.BatchOps, 5)
	require.Equal(t, suiTestnetSelector, uint64(report.Output.BatchOps[0].ChainSelector))

	feeQuoter, err := module_fee_quoter.NewFeeQuoter(testCCIPPackageID, nil)
	require.NoError(t, err)
	encodedCall, err := feeQuoter.Encoder().ApplyTokenTransferFeeConfigUpdates(
		bind.Object{Id: testCCIPObjectRef},
		bind.Object{Id: testCCIPOwnerCapID},
		evmSepoliaSelector,
		[]string{testLinkCoinMetadata},
		[]uint32{lanes.DefaultLinkTokenTransferMinFeeUsdCents},
		[]uint32{lanes.DefaultLinkTokenTransferMaxFeeUsdCents},
		[]uint16{lanes.DefaultLinkTokenTransferDeciBps},
		[]uint32{lanes.DefaultLinkTokenTransferDestGasOverhead},
		[]uint32{lanes.DefaultLinkTokenTransferDestBytesOverhead},
		[]bool{true},
		[]string{},
	)
	require.NoError(t, err)
	mcmstest.AssertProposalDataMatches(
		t,
		report.Output.BatchOps[0].Transactions[0].Data,
		encodedCall,
		testCCIPObjectRef,
		nil,
	)
}

func runSourceLeg(
	t *testing.T,
	env cldf.Environment,
	chains chain.BlockChains,
	input laneapi.UpdateLanesInput,
) (cldf_ops.SequenceReport[laneapi.UpdateLanesInput, sequences.OnChainOutput], error) {
	t.Helper()

	var report cldf_ops.SequenceReport[laneapi.UpdateLanesInput, sequences.OnChainOutput]
	err := lanes.WithConnectChainsEnvironment(env, func() error {
		var execErr error
		report, execErr = cldf_ops.ExecuteSequence(
			mcmstest.Bundle(t),
			lanes.ConfigureLaneLegAsSource,
			chains,
			input,
		)
		return execErr
	})
	return report, err
}
