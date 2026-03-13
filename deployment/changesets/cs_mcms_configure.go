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
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type ConfigureMCMSConfig struct {
	mcmsops.ConfigureMCMSSeqInput
	TimelockConfig *utils.TimelockConfig // If nil, configuration will be executed directly
	IsFastCurse    bool                  // If true, the fastcurse MCMS instance is configured
}

var _ cldf.ChangeSetV2[ConfigureMCMSConfig] = ConfigureMCMS{}

type ConfigureMCMS struct{}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (c ConfigureMCMS) VerifyPreconditions(e cldf.Environment, config ConfigureMCMSConfig) error {
	return nil
}

// Apply implements deployment.ChangeSetV2.
func (c ConfigureMCMS) Apply(e cldf.Environment, config ConfigureMCMSConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]cld_ops.Report[any, any], 0)

	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	suiChains := e.BlockChains.SuiChains()

	suiChain := suiChains[config.ChainSelector]

	mcmsState := state[config.ChainSelector].MCMSState(config.IsFastCurse)

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

	// If timelock proposal is to be generated, disable signer in deps
	if config.TimelockConfig != nil {
		deps.Signer = nil
	}

	// Run ConfigureMCMS Sequence
	configReport, err := cld_ops.ExecuteSequence(e.OperationsBundle, mcmsops.ConfigureMCMSSequence, deps, config.ConfigureMCMSSeqInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to configure MCMS for Sui chain %d: %w", config.ChainSelector, err)
	}

	seqReports = append(seqReports, configReport.ToGenericReport())

	mcmsProposal := mcms.TimelockProposal{}
	if config.TimelockConfig != nil {
		defs := []cld_ops.Definition{}
		inputs := []any{}

		for _, r := range configReport.Output.Reports {
			defs = append(defs, r.Def)
			inputs = append(inputs, r.Input)
		}

		mcmsConfig := mcmsops.ProposalGenerateInput{
			ChainSelector:      config.ChainSelector,
			Defs:               defs,
			Inputs:             inputs,
			MmcsPackageID:      mcmsState.PackageID,
			McmsStateObjID:     mcmsState.StateObjectID,
			TimelockObjID:      mcmsState.TimelockObjectID,
			AccountObjID:       mcmsState.AccountStateObjectID,
			RegistryObjID:      mcmsState.RegistryObjectID,
			DeployerStateObjID: mcmsState.DeployerStateObjectID,
			TimelockConfig:     *config.TimelockConfig,
		}

		result, err := cld_ops.ExecuteSequence(e.OperationsBundle, mcmsops.MCMSDynamicProposalGenerateSeq, deps, mcmsConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal: %w", err)
		}
		mcmsProposal = result.Output
	}

	return cldf.ChangesetOutput{
		AddressBook:           ab,
		Reports:               seqReports,
		MCMSTimelockProposals: []mcms.TimelockProposal{mcmsProposal},
	}, nil
}
