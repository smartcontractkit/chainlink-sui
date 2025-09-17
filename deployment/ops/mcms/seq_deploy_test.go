//go:build integration

package mcmsops

import (
	"context"
	"slices"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/bindings/tests/testenv"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
)

// generateSortedSigners creates a sorted array of test signers similar to Test_Sui_SetConfig
func generateSortedSigners(count int) []common.Address {
	signers := make([]common.Address, count)
	for i := range signers {
		key, _ := crypto.GenerateKey()
		signers[i] = crypto.PubkeyToAddress(key.PublicKey)
	}
	// Sort signers alphabetically by address
	slices.SortFunc(signers, func(a, b common.Address) int {
		return strings.Compare(strings.ToLower(a.Hex()), strings.ToLower(b.Hex()))
	})
	return signers
}

func TestDeployMCMSSeq(t *testing.T) {
	t.Parallel()

	signer, client := testenv.SetupEnvironment(t)

	deps := sui_ops.OpTxDeps{
		Client: client,
		Signer: signer,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(300_000_000)
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
	}

	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(
		context.Background,
		logger.Test(t),
		reporter,
	)

	// Generate sorted signers for testing (similar to Test_Sui_SetConfig)
	signers := generateSortedSigners(30)

	// Create bypasser config following the Test_Sui_SetConfig pattern
	bypasserConfig := &types.Config{
		Quorum: 2,
		Signers: []common.Address{
			signers[0],
			signers[1],
			signers[2],
		},
		GroupSigners: []types.Config{
			{
				Quorum: 4,
				Signers: []common.Address{
					signers[3],
					signers[4],
					signers[5],
					signers[6],
					signers[7],
				},
				GroupSigners: []types.Config{
					{
						Quorum: 1,
						Signers: []common.Address{
							signers[8],
							signers[9],
						},
						GroupSigners: []types.Config{},
					},
				},
			},
			{
				Quorum: 3,
				Signers: []common.Address{
					signers[10],
					signers[11],
					signers[12],
					signers[13],
				},
				GroupSigners: []types.Config{},
			},
		},
	}

	// Create canceller config
	cancellerConfig := &types.Config{
		Quorum: 1,
		Signers: []common.Address{
			signers[14],
			signers[15],
		},
		GroupSigners: []types.Config{
			{
				Quorum: 2,
				Signers: []common.Address{
					signers[16],
					signers[17],
					signers[18],
					signers[19],
				},
				GroupSigners: []types.Config{},
			},
		},
	}

	// Create proposer config
	proposerConfig := &types.Config{
		Quorum:  2,
		Signers: []common.Address{},
		GroupSigners: []types.Config{
			{
				Quorum: 2,
				Signers: []common.Address{
					signers[20],
					signers[21],
					signers[22],
					signers[23],
				},
				GroupSigners: []types.Config{},
			}, {
				Quorum: 2,
				Signers: []common.Address{
					signers[24],
					signers[25],
					signers[26],
					signers[27],
				},
				GroupSigners: []types.Config{},
			}, {
				Quorum: 1,
				Signers: []common.Address{
					signers[28],
					signers[29],
				},
				GroupSigners: []types.Config{},
			},
		},
	}

	report, err := cld_ops.ExecuteSequence(bundle, DeployMCMSSequence, deps, DeployMCMSSeqInput{
		ChainSelector: cselectors.SUI_TESTNET.Selector,
		Bypasser:      bypasserConfig,
		Proposer:      proposerConfig,
		Canceller:     cancellerConfig,
	})
	require.NoError(t, err, "failed to execute MCMS deploy sequence")

	objects := report.Output.Objects
	require.NotEmpty(t, objects.McmsMultisigStateObjectId, "MCMS Multisig State Object ID should not be empty")
	require.NotEmpty(t, objects.TimelockObjectId, "MCMS Timelock Object ID should not be empty")
	require.NotEmpty(t, objects.McmsDeployerObjectId, "MCMS Deployer Object ID should not be empty")
	require.NotEmpty(t, objects.McmsRegistryObjectId, "MCMS Registry Object ID should not be empty")
	require.NotEmpty(t, objects.McmsAccountStateObjectId, "MCMS Account State Object ID should not be empty")
	require.NotEmpty(t, objects.McmsAccountOwnerCapObjectId, "MCMS Account Owner Cap Object ID should not be empty")
	require.NotEmpty(t, report.Output.Digest, "Transaction digest should not be empty")
	require.NotEmpty(t, report.Output.PackageId, "Package ID should not be empty")
}
