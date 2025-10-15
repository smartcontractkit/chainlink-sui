package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	offrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"
	onrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_onramp"
	routerops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_router"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	opregistry "github.com/smartcontractkit/chainlink-sui/deployment/ops/registry"
	"github.com/smartcontractkit/mcms"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
)

type AcceptOwnershipCCIPConfig struct {
	SuiChainSelector uint64
}

var _ cldf.ChangeSetV2[AcceptOwnershipCCIPConfig] = AcceptOwnershipCCIP{}

// AcceptOwnershipCCIP deploys Sui chain packages and modules
type AcceptOwnershipCCIP struct{}

// Apply implements deployment.ChangeSetV2.
func (d AcceptOwnershipCCIP) Apply(e cldf.Environment, config AcceptOwnershipCCIPConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]operations.Report[any, any], 0)

	suiChain := e.BlockChains.SuiChains()[config.SuiChainSelector]
	signer := suiChain.Signer

	deps := sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: signer,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(500_000_000)
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
	}

	// in case the registry is not loaded with all operations. Needed to build accept ownership proposals
	ops := make([]*cld_ops.Operation[any, any, any], len(opregistry.AllOperations))
	for i := range opregistry.AllOperations {
		ops[i] = &opregistry.AllOperations[i]
		cld_ops.RegisterOperation(e.OperationsBundle.OperationRegistry, &opregistry.AllOperations[i])
	}

	suiState, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	state := suiState[config.SuiChainSelector]

	// Generate the proposal to accept the ownership of the deployed contracts
	proposalInput := mcmsops.ProposalGenerateInput{
		Defs: []operations.Definition{
			ccipops.AcceptOwnershipStateObjectOp.Def(),
			routerops.AcceptOwnershipOp.Def(),
			onrampops.AcceptOwnershipOnRampOp.Def(),
			offrampops.AcceptOwnershipOffRampOp.Def(),
		},
		Inputs: []any{
			ccipops.AcceptOwnershipStateObjectInput{
				CCIPPackageId:         state.CCIPAddress,
				CCIPObjectRefObjectId: state.CCIPObjectRef,
			},
			routerops.AcceptOwnershipInput{
				RouterPackageId:     state.CCIPRouterAddress,
				RouterStateObjectId: state.CCIPRouterStateObjectID,
			},
			onrampops.AcceptOwnershipOnRampInput{
				OnRampPackageId: state.OnRampAddress,
				CCIPObjectRefId: state.CCIPObjectRef,
				StateObjectId:   state.OnRampStateObjectId,
			},
			offrampops.AcceptOwnershipOffRampInput{
				OffRampPackageId:     state.OffRampAddress,
				OffRampRefObjectId:   state.CCIPObjectRef,
				OffRampStateObjectId: state.OffRampStateObjectId,
			},
		},
		// MCMS related
		MmcsPackageID:  state.MCMSPackageID,
		McmsStateObjID: state.MCMSStateObjectID,
		TimelockObjID:  state.MCMSTimelockObjectID,
		AccountObjID:   state.MCMSAccountStateObjectID,
		RegistryObjID:  state.MCMSRegistryObjectID,

		// Proposal
		Role: suisdk.TimelockRoleProposer,

		ChainSelector: config.SuiChainSelector,
	}

	acceptOwnershipProposalReport, err := cld_ops.ExecuteSequence(e.OperationsBundle, mcmsops.MCMSDynamicProposalGenerateSeq, deps, proposalInput)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		AddressBook:           ab,
		Reports:               seqReports,
		MCMSTimelockProposals: []mcms.TimelockProposal{acceptOwnershipProposalReport.Output},
	}, nil
}

// TODO
// VerifyPreconditions imsplements deployment.ChangeSetV2.
func (d AcceptOwnershipCCIP) VerifyPreconditions(e cldf.Environment, config AcceptOwnershipCCIPConfig) error {
	return nil
}
