package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/mcms"
)

var _ cldf.ChangeSetV2[mcmsops.ProposalGenerateInput] = MCMSProposalGenerate{}

type MCMSProposalGenerate struct{}

func (d MCMSProposalGenerate) Apply(e cldf.Environment, config mcmsops.ProposalGenerateInput) (cldf.ChangesetOutput, error) {
	suiState, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	state := suiState[config.ChainSelector]

	// Get necessary MCMS state from onchain AB
	config.MmcsPackageID = state.MCMSPackageID
	config.McmsStateObjID = state.MCMSStateObjectID
	config.TimelockObjID = state.MCMSTimelockObjectID
	config.AccountObjID = state.MCMSAccountStateObjectID
	config.RegistryObjID = state.MCMSRegistryObjectID

	suiChains := e.BlockChains.SuiChains()

	suiChain := suiChains[config.ChainSelector]
	deps := sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: nil, // Signer is not needed since we are not executing any transactions
		GetCallOpts: func() *bind.CallOpts {
			return &bind.CallOpts{}
		},
	}
	result, err := cld_ops.ExecuteSequence(e.OperationsBundle, mcmsops.MCMSDynamicProposalGenerateSeq, deps, config)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute sequence: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{result.Output},
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d MCMSProposalGenerate) VerifyPreconditions(e cldf.Environment, config mcmsops.ProposalGenerateInput) error {
	return nil
}
