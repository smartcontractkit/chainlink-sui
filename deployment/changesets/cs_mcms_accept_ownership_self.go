package changesets

import (
	"fmt"

	"github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
)

var _ cldf.ChangeSetV2[AcceptMCMSOwnershipConfig] = AcceptMCMSOwnership{}

// AcceptMCMSOwnershipConfig identifies an MCMS deployment whose ownership
// transfer proposal should be (re-)generated. When IsFastCurse is true the
// fastcurse MCMS instance is targeted; otherwise the normal instance is used.
type AcceptMCMSOwnershipConfig struct {
	ChainSelector uint64
	IsFastCurse   bool
}

type AcceptMCMSOwnership struct{}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (a AcceptMCMSOwnership) VerifyPreconditions(e cldf.Environment, config AcceptMCMSOwnershipConfig) error {
	return nil
}

// Apply implements deployment.ChangeSetV2.
func (a AcceptMCMSOwnership) Apply(e cldf.Environment, config AcceptMCMSOwnershipConfig) (cldf.ChangesetOutput, error) {
	suiState, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load sui onchain state: %w", err)
	}

	mcmsFields := suiState[config.ChainSelector].MCMSState(config.IsFastCurse)

	seqInput := mcmsops.AcceptMCMSOwnershipSeqInput{
		ChainSelector:             config.ChainSelector,
		PackageId:                 mcmsFields.PackageID,
		McmsMultisigStateObjectId: mcmsFields.StateObjectID,
		TimelockObjectId:          mcmsFields.TimelockObjectID,
		McmsAccountStateObjectId:  mcmsFields.AccountStateObjectID,
		McmsRegistryObjectId:      mcmsFields.RegistryObjectID,
		McmsDeployerStateObjectId: mcmsFields.DeployerStateObjectID,
	}
	suiChains := e.BlockChains.SuiChains()
	suiChain := suiChains[config.ChainSelector]
	deps := sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: suiChain.Signer,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(400_000_000)
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
		SuiRPC: suiChain.URL,
	}

	// Run AcceptMCMSOwnership Sequence
	acceptReport, err := cld_ops.ExecuteSequence(e.OperationsBundle, mcmsops.AcceptMCMSOwnershipSequence, deps, seqInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate accept ownership proposal for Sui chain %d: %w", config.ChainSelector, err)
	}

	return cldf.ChangesetOutput{
		Reports:               []cld_ops.Report[any, any]{acceptReport.ToGenericReport()},
		MCMSTimelockProposals: []mcms.TimelockProposal{acceptReport.Output},
	}, nil
}
