package changesets

import (
	"fmt"

	"github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
)

var _ cldf.ChangeSetV2[mcmsops.AcceptMCMSOwnershipSeqInput] = AcceptMCMSOwnership{}

type AcceptMCMSOwnership struct{}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (a AcceptMCMSOwnership) VerifyPreconditions(e cldf.Environment, config mcmsops.AcceptMCMSOwnershipSeqInput) error {
	return nil
}

// Apply implements deployment.ChangeSetV2.
func (a AcceptMCMSOwnership) Apply(e cldf.Environment, config mcmsops.AcceptMCMSOwnershipSeqInput) (cldf.ChangesetOutput, error) {
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
	acceptReport, err := cld_ops.ExecuteSequence(e.OperationsBundle, mcmsops.AcceptMCMSOwnershipSequence, deps, config)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate accept ownership proposal for Sui chain %d: %w", config.ChainSelector, err)
	}

	return cldf.ChangesetOutput{
		Reports:               []cld_ops.Report[any, any]{acceptReport.ToGenericReport()},
		MCMSTimelockProposals: []mcms.TimelockProposal{acceptReport.Output},
	}, nil
}
