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

var _ cldf.ChangeSetV2[MCMSProposalGenerateConfig] = MCMSProposalGenerate{}

// MCMSProposalGenerateConfig wraps ProposalGenerateInput and adds IsFastCurse.
// When MCMS state fields in ProposalGenerateInput are left empty, they are
// auto-populated from the on-chain address book using the IsFastCurse flag.
type MCMSProposalGenerateConfig struct {
	mcmsops.ProposalGenerateInput
	IsFastCurse bool
}

type MCMSProposalGenerate struct{}

func (d MCMSProposalGenerate) Apply(e cldf.Environment, config MCMSProposalGenerateConfig) (cldf.ChangesetOutput, error) {
	suiState, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	mcmsState := suiState[config.ChainSelector].MCMSState(config.IsFastCurse)

	// Get necessary MCMS state from onchain AB
	if config.MmcsPackageID == "" || config.McmsStateObjID == "" || config.TimelockObjID == "" || config.AccountObjID == "" || config.RegistryObjID == "" {
		config.MmcsPackageID = mcmsState.PackageID
		config.McmsStateObjID = mcmsState.StateObjectID
		config.TimelockObjID = mcmsState.TimelockObjectID
		config.AccountObjID = mcmsState.AccountStateObjectID
		config.RegistryObjID = mcmsState.RegistryObjectID
	}

	suiChains := e.BlockChains.SuiChains()

	suiChain := suiChains[config.ChainSelector]
	deps := sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: nil, // Signer is not needed since we are not executing any transactions
		GetCallOpts: func() *bind.CallOpts {
			return &bind.CallOpts{}
		},
		SuiRPC: suiChain.URL,
	}
	result, err := cld_ops.ExecuteSequence(e.OperationsBundle, mcmsops.MCMSDynamicProposalGenerateSeq, deps, config.ProposalGenerateInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute sequence: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{result.Output},
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d MCMSProposalGenerate) VerifyPreconditions(e cldf.Environment, config MCMSProposalGenerateConfig) error {
	return nil
}
