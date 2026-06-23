package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccip_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	ccip_offramp_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"
	ccip_onramp_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_onramp"
	ccip_router_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_router"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type ConnectSuiToEVMConfig struct {
	SuiChainSelector                                     uint64
	FeeQuoterApplyTokenTransferFeeConfigUpdatesInput     ccip_ops.FeeQuoterApplyTokenTransferFeeConfigUpdatesInput
	FeeQuoterApplyDestChainConfigUpdatesInput            ccip_ops.FeeQuoterApplyDestChainConfigUpdatesInput
	FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput ccip_ops.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput
	ApplyDestChainConfigureOnRampInput                   ccip_onramp_ops.ApplyDestChainConfigureOnRampInput
	ApplySourceChainConfigUpdateInput                    ccip_offramp_ops.ApplySourceChainConfigUpdateInput
	TimelockConfig                                       *utils.TimelockConfig // If nil, transactions will be executed
}

// ConnectSuiToEVM connects sui chain with EVM
type ConnectSuiToEVM struct{}

var _ cldf.ChangeSetV2[ConnectSuiToEVMConfig] = ConnectSuiToEVM{}

// Apply implements deployment.ChangeSetV2.
func (d ConnectSuiToEVM) Apply(e cldf.Environment, config ConnectSuiToEVMConfig) (cldf.ChangesetOutput, error) {
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	seqReports := make([]operations.Report[any, any], 0)

	suiChains := e.BlockChains.SuiChains()
	suiChain := suiChains[config.SuiChainSelector]

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

	if config.TimelockConfig != nil {
		deps.Signer = nil
	}

	defs := []operations.Definition{}
	inputs := []any{}

	// Configure FeeQuoter
	config.FeeQuoterApplyTokenTransferFeeConfigUpdatesInput.CCIPPackageId = state[config.SuiChainSelector].CCIPAddress
	config.FeeQuoterApplyTokenTransferFeeConfigUpdatesInput.StateObjectId = state[config.SuiChainSelector].CCIPObjectRef
	config.FeeQuoterApplyTokenTransferFeeConfigUpdatesInput.OwnerCapObjectId = state[config.SuiChainSelector].CCIPOwnerCapObjectId
	reportFeeQuoterApplyTokenTransferFeeConfigUpdatesOp, err := operations.ExecuteOperation(e.OperationsBundle, ccip_ops.FeeQuoterApplyTokenTransferFeeConfigUpdatesOp, deps, config.FeeQuoterApplyTokenTransferFeeConfigUpdatesInput)
	defs = append(defs, reportFeeQuoterApplyTokenTransferFeeConfigUpdatesOp.Def)
	inputs = append(inputs, reportFeeQuoterApplyTokenTransferFeeConfigUpdatesOp.Input)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to run FeeQuoterApplyTokenTransferFeeConfigUpdatesOp for Sui chain %d: %w", config.SuiChainSelector, err)
	}
	seqReports = append(seqReports, []operations.Report[any, any]{reportFeeQuoterApplyTokenTransferFeeConfigUpdatesOp.ToGenericReport()}...)

	config.FeeQuoterApplyDestChainConfigUpdatesInput.CCIPPackageId = state[config.SuiChainSelector].CCIPAddress
	config.FeeQuoterApplyDestChainConfigUpdatesInput.StateObjectId = state[config.SuiChainSelector].CCIPObjectRef
	config.FeeQuoterApplyDestChainConfigUpdatesInput.OwnerCapObjectId = state[config.SuiChainSelector].CCIPOwnerCapObjectId
	reportFeeQuoterApplyDestChainConfigUpdatesOp, err := operations.ExecuteOperation(e.OperationsBundle, ccip_ops.FeeQuoterApplyDestChainConfigUpdatesOp, deps, config.FeeQuoterApplyDestChainConfigUpdatesInput)
	defs = append(defs, reportFeeQuoterApplyDestChainConfigUpdatesOp.Def)
	inputs = append(inputs, reportFeeQuoterApplyDestChainConfigUpdatesOp.Input)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to run FeeQuoterApplyDestChainConfigUpdatesOp for Sui chain %d: %w", config.SuiChainSelector, err)
	}
	seqReports = append(seqReports, []operations.Report[any, any]{reportFeeQuoterApplyDestChainConfigUpdatesOp.ToGenericReport()}...)

	config.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput.CCIPPackageId = state[config.SuiChainSelector].CCIPAddress
	config.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput.StateObjectId = state[config.SuiChainSelector].CCIPObjectRef
	config.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput.OwnerCapObjectId = state[config.SuiChainSelector].CCIPOwnerCapObjectId
	reportFeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp, err := operations.ExecuteOperation(e.OperationsBundle, ccip_ops.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp, deps, config.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput)
	defs = append(defs, reportFeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp.Def)
	inputs = append(inputs, reportFeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp.Input)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to run FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp for Sui chain %d: %w", config.SuiChainSelector, err)
	}
	seqReports = append(seqReports, []operations.Report[any, any]{reportFeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp.ToGenericReport()}...)

	// Configure OnRamp
	config.ApplyDestChainConfigureOnRampInput.OnRampPackageId = state[config.SuiChainSelector].OnRampAddress
	config.ApplyDestChainConfigureOnRampInput.OwnerCapObjectId = state[config.SuiChainSelector].OnRampOwnerCapObjectId
	config.ApplyDestChainConfigureOnRampInput.StateObjectId = state[config.SuiChainSelector].OnRampStateObjectId
	config.ApplyDestChainConfigureOnRampInput.CCIPObjectRefId = state[config.SuiChainSelector].CCIPObjectRef
	reportApplyDestChainConfigUpdateOp, err := operations.ExecuteOperation(e.OperationsBundle, ccip_onramp_ops.ApplyDestChainConfigUpdateOp, deps, config.ApplyDestChainConfigureOnRampInput)
	defs = append(defs, reportApplyDestChainConfigUpdateOp.Def)
	inputs = append(inputs, reportApplyDestChainConfigUpdateOp.Input)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to run ApplyDestChainConfigUpdateOp for Sui chain %d: %w", config.SuiChainSelector, err)
	}
	seqReports = append(seqReports, []operations.Report[any, any]{reportApplyDestChainConfigUpdateOp.ToGenericReport()}...)

	// Configure OffRamp
	config.ApplySourceChainConfigUpdateInput.CCIPObjectRef = state[config.SuiChainSelector].CCIPObjectRef
	config.ApplySourceChainConfigUpdateInput.OffRampPackageId = state[config.SuiChainSelector].OffRampAddress
	config.ApplySourceChainConfigUpdateInput.OffRampStateId = state[config.SuiChainSelector].OffRampStateObjectId
	config.ApplySourceChainConfigUpdateInput.OwnerCapObjectId = state[config.SuiChainSelector].OffRampOwnerCapId
	reportApplySourceChainConfigUpdatesOp, err := operations.ExecuteOperation(e.OperationsBundle, ccip_offramp_ops.ApplySourceChainConfigUpdatesOp, deps, config.ApplySourceChainConfigUpdateInput)
	defs = append(defs, reportApplySourceChainConfigUpdatesOp.Def)
	inputs = append(inputs, reportApplySourceChainConfigUpdatesOp.Input)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to run ApplySourceChainConfigUpdatesOp for Sui chain %d: %w", config.SuiChainSelector, err)
	}
	seqReports = append(seqReports, []operations.Report[any, any]{reportApplySourceChainConfigUpdatesOp.ToGenericReport()}...)

	// Configure Router
	onrampAddresses := make([]string, len(config.ApplyDestChainConfigureOnRampInput.DestChainSelector))
	for i := range config.ApplyDestChainConfigureOnRampInput.DestChainSelector {
		onrampAddresses[i] = config.ApplyDestChainConfigureOnRampInput.OnRampPackageId
	}
	reportConfigureRouterOp, err := operations.ExecuteOperation(e.OperationsBundle, ccip_router_ops.SetOnRampsOp, deps, ccip_router_ops.SetOnRampsInput{
		RouterPackageId:     state[config.SuiChainSelector].CCIPRouterAddress,
		RouterStateObjectId: state[config.SuiChainSelector].CCIPRouterStateObjectID,
		OwnerCapObjectId:    state[config.SuiChainSelector].CCIPRouterOwnerCapObjectId,
		DestChainSelectors:  config.ApplyDestChainConfigureOnRampInput.DestChainSelector,
		OnRampAddresses:     onrampAddresses,
	})
	defs = append(defs, reportConfigureRouterOp.Def)
	inputs = append(inputs, reportConfigureRouterOp.Input)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to run ConfigureRouterOp for Sui chain %d: %w", config.SuiChainSelector, err)
	}
	seqReports = append(seqReports, []operations.Report[any, any]{reportConfigureRouterOp.ToGenericReport()}...)

	mcmsProposal := mcms.TimelockProposal{}
	if config.TimelockConfig != nil {
		mcmsConfig := mcmsops.ProposalGenerateInput{
			ChainSelector:      config.SuiChainSelector,
			Defs:               defs,
			Inputs:             inputs,
			MmcsPackageID:      state[config.SuiChainSelector].MCMSPackageID,
			McmsStateObjID:     state[config.SuiChainSelector].MCMSStateObjectID,
			TimelockObjID:      state[config.SuiChainSelector].MCMSTimelockObjectID,
			AccountObjID:       state[config.SuiChainSelector].MCMSAccountStateObjectID,
			RegistryObjID:      state[config.SuiChainSelector].MCMSRegistryObjectID,
			DeployerStateObjID: state[config.SuiChainSelector].MCMSDeployerStateObjectID,
			TimelockConfig:     *config.TimelockConfig,
		}
		result, err := operations.ExecuteSequence(e.OperationsBundle, mcmsops.MCMSDynamicProposalGenerateSeq, deps, mcmsConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute sequence: %w", err)
		}
		mcmsProposal = result.Output
	}

	return cldf.ChangesetOutput{
		Reports:               seqReports,
		MCMSTimelockProposals: []mcms.TimelockProposal{mcmsProposal},
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d ConnectSuiToEVM) VerifyPreconditions(e cldf.Environment, config ConnectSuiToEVMConfig) error {
	return nil
}
