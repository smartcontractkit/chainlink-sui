package changesets

import (
	"fmt"

	"github.com/smartcontractkit/mcms"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
)

var _ cldf.ChangeSetV2[DeployMCMSConfig] = DeployMCMS{}

// DeployMCMSConfig wraps DeployMCMSSeqInput and adds the IsFastCurse flag.
// When IsFastCurse is true the fast_mcms package is published and all address-book
// entries are stored with the "fastcurse" label so LoadOnchainStatesui can distinguish
// the two MCMS instances deployed on the same chain.
type DeployMCMSConfig struct {
	mcmsops.DeployMCMSSeqInput
	IsFastCurse bool `yaml:"isFastCurse,omitempty"`
}

type DeployMCMS struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeployMCMS) Apply(e cldf.Environment, config DeployMCMSConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	ds := fdatastore.NewMemoryDataStore()
	seqReports := make([]cld_ops.Report[any, any], 0)

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

	seqInput := config.DeployMCMSSeqInput
	seqInput.FastMCMS = config.IsFastCurse

	mcmsReport, err := cld_ops.ExecuteSequence(e.OperationsBundle, mcmsops.DeployMCMSSequence, deps, seqInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy MCMS for Sui chain %d: %w", config.ChainSelector, err)
	}

	err = deployment.StoreMCMSInAddressBookAndDataStore(ab, ds.Addresses(), config.ChainSelector, mcmsReport.Output, deployment.MCMSInstanceFromFastCurseFlag(config.IsFastCurse))
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to store MCMS in address book for Sui chain %d: %w", config.ChainSelector, err)
	}

	proposals := []mcms.TimelockProposal{}
	if !seqInput.SkipOwnershipTransfer {
		proposals = append(proposals, mcmsReport.Output.AcceptOwnershipProposal)
	}

	return cldf.ChangesetOutput{
		AddressBook:           ab,
		DataStore:             ds,
		Reports:               seqReports,
		MCMSTimelockProposals: proposals,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeployMCMS) VerifyPreconditions(e cldf.Environment, config DeployMCMSConfig) error {
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return fmt.Errorf("load onchain state: %w", err)
	}
	chainState, ok := state[config.ChainSelector]
	if !ok {
		return nil
	}
	instance := deployment.MCMSInstanceFromFastCurseFlag(config.IsFastCurse)
	if chainState.HasMCMSInstance(instance) {
		return fmt.Errorf("%s MCMS is already recorded for chain %d", instance, config.ChainSelector)
	}
	return nil
}
