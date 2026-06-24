package lanes_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	laneapi "github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/deployment/lanes"
	"github.com/smartcontractkit/chainlink-sui/deployment/ops/mcmstest"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

const evmSepoliaSelector = uint64(16015286601757825753)

func TestConfigureLaneLegAsDest_MCMSBatchOp(t *testing.T) {
	env := testEnvWithAddressBook(t)
	chains := chain.NewBlockChains(map[uint64]chain.BlockChain{
		suiTestnetSelector: sui.Chain{},
	})
	input := laneapi.UpdateLanesInput{
		Source: &laneapi.ChainDefinition{
			Selector:               evmSepoliaSelector,
			OnRamp:                 make([]byte, 32),
			RMNVerificationEnabled: true,
		},
		Dest: &laneapi.ChainDefinition{
			Selector: suiTestnetSelector,
		},
	}

	var report cldf_ops.SequenceReport[laneapi.UpdateLanesInput, sequences.OnChainOutput]
	err := lanes.WithConnectChainsEnvironment(env, func() error {
		var execErr error
		report, execErr = cldf_ops.ExecuteSequence(
			mcmstest.Bundle(t),
			lanes.ConfigureLaneLegAsDest,
			chains,
			input,
		)
		return execErr
	})
	require.NoError(t, err)
	require.Len(t, report.Output.BatchOps, 1)
	require.Equal(t, suiTestnetSelector, uint64(report.Output.BatchOps[0].ChainSelector))
	require.Len(t, report.Output.BatchOps[0].Transactions, 1)
	require.NotEmpty(t, report.Output.BatchOps[0].Transactions[0].Data)
}

func TestConfigureLaneLegAsDest_MCMSBatchOp_WithLatestPackageID(t *testing.T) {
	env := testEnvWithAddressBook(t)
	chains := chain.NewBlockChains(map[uint64]chain.BlockChain{
		suiTestnetSelector: sui.Chain{},
	})
	const latestOffRampPackageID = "0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"
	input := laneapi.UpdateLanesInput{
		Source: &laneapi.ChainDefinition{
			Selector:               evmSepoliaSelector,
			OnRamp:                 make([]byte, 32),
			RMNVerificationEnabled: true,
		},
		Dest: &laneapi.ChainDefinition{
			Selector: suiTestnetSelector,
		},
	}
	var report cldf_ops.SequenceReport[laneapi.UpdateLanesInput, sequences.OnChainOutput]
	err := lanes.WithConnectChainsEnvironment(env, func() error {
		return lanes.WithSuiLatestPackageIDs(map[uint64]lanes.LatestPackageIDsConfig{
			suiTestnetSelector: {OffRamp: latestOffRampPackageID},
		}, func() error {
			var execErr error
			report, execErr = cldf_ops.ExecuteSequence(
				mcmstest.Bundle(t),
				lanes.ConfigureLaneLegAsDest,
				chains,
				input,
			)
			return execErr
		})
	})
	require.NoError(t, err)
	require.Len(t, report.Output.BatchOps, 1)
	tx := report.Output.BatchOps[0].Transactions[0]
	latestPackageIDFromTx, err := utils.TransactionLatestPackageID(tx)
	require.NoError(t, err)
	require.Equal(t, latestOffRampPackageID, latestPackageIDFromTx)
}
