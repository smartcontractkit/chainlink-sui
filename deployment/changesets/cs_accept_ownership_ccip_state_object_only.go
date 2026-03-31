package changesets

import (
	"fmt"

	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/types"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ownershipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ownership"
	opregistry "github.com/smartcontractkit/chainlink-sui/deployment/ops/registry"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

var _ cldf.ChangeSetV2[AcceptOwnershipCCIPStateObjectOnlyConfig] = AcceptOwnershipCCIPStateObjectOnly{}

// AcceptOwnershipCCIPStateObjectOnly generates an MCMS timelock proposal that only schedules
// accept_ownership on the CCIP state object (not router, onramp, or offramp).
type AcceptOwnershipCCIPStateObjectOnly struct{}

type AcceptOwnershipCCIPStateObjectOnlyConfig struct {
	SuiChainSelector uint64 `json:"suiChainSelector" yaml:"suichainselector"`
}

func (d AcceptOwnershipCCIPStateObjectOnly) Apply(e cldf.Environment, config AcceptOwnershipCCIPStateObjectOnlyConfig) (cldf.ChangesetOutput, error) {
	suiChain := e.BlockChains.SuiChains()[config.SuiChainSelector]
	signer := suiChain.Signer

	deps := sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: signer,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(1_000_000_000)
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
		SuiRPC: suiChain.URL,
	}

	for i := range opregistry.AllOperations {
		cld_ops.RegisterOperation(e.OperationsBundle.OperationRegistry, opregistry.AllOperations[i])
	}

	suiState, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	state := suiState[config.SuiChainSelector]

	proposalInput := ownershipops.AcceptCCIPStateObjectOnlyInput{
		ChainSelector: config.SuiChainSelector,

		MCMSPackageId:          state.MCMSPackageID,
		MCMSStateObjId:         state.MCMSStateObjectID,
		MCMSTimelockObjId:      state.MCMSTimelockObjectID,
		MCMSAccountObjId:       state.MCMSAccountStateObjectID,
		MCMSRegistryObjId:      state.MCMSRegistryObjectID,
		MCMSDeployerStateObjId: state.MCMSDeployerStateObjectID,

		CCIPPackageId: state.CCIPAddress,
		CCIPObjectRef: state.CCIPObjectRef,

		TimelockConfig: utils.TimelockConfig{
			MCMSAction:   types.TimelockActionSchedule,
			MinDelay:     0,
			OverrideRoot: false,
		},
	}

	report, err := cld_ops.ExecuteSequence(e.OperationsBundle, ownershipops.AcceptCCIPStateObjectOwnershipOnlySeq, deps, proposalInput)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{report.Output},
	}, nil
}

func (d AcceptOwnershipCCIPStateObjectOnly) VerifyPreconditions(e cldf.Environment, config AcceptOwnershipCCIPStateObjectOnlyConfig) error {
	return nil
}
