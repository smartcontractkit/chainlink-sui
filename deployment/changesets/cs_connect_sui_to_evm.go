package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccip_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	ccip_offramp_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"
	ccip_onramp_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_onramp"
	ccip_router_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_router"
)

type ConnectSuiToEVMConfig struct {
	SuiChainSelector                                     uint64
	FeeQuoterApplyTokenTransferFeeConfigUpdatesInput     ccip_ops.FeeQuoterApplyTokenTransferFeeConfigUpdatesInput
	FeeQuoterApplyDestChainConfigUpdatesInput            ccip_ops.FeeQuoterApplyDestChainConfigUpdatesInput
	FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput ccip_ops.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput
	ApplyDestChainConfigureOnRampInput                   ccip_onramp_ops.ApplyDestChainConfigureOnRampInput
	ApplySourceChainConfigUpdateInput                    ccip_offramp_ops.ApplySourceChainConfigUpdateInput
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

	// Configure FeeQuoter
	config.FeeQuoterApplyTokenTransferFeeConfigUpdatesInput.CCIPPackageId = state[config.SuiChainSelector].CCIPAddress
	config.FeeQuoterApplyTokenTransferFeeConfigUpdatesInput.StateObjectId = state[config.SuiChainSelector].CCIPObjectRef
	config.FeeQuoterApplyTokenTransferFeeConfigUpdatesInput.OwnerCapObjectId = state[config.SuiChainSelector].CCIPOwnerCapObjectId
	reportFeeQuoterApplyTokenTransferFeeConfigUpdatesOp, err := operations.ExecuteOperation(e.OperationsBundle, ccip_ops.FeeQuoterApplyTokenTransferFeeConfigUpdatesOp, deps, config.FeeQuoterApplyTokenTransferFeeConfigUpdatesInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to run FeeQuoterApplyTokenTransferFeeConfigUpdatesOp for Sui chain %d: %w", config.SuiChainSelector, err)
	}
	seqReports = append(seqReports, []operations.Report[any, any]{reportFeeQuoterApplyTokenTransferFeeConfigUpdatesOp.ToGenericReport()}...)

	config.FeeQuoterApplyDestChainConfigUpdatesInput.CCIPPackageId = state[config.SuiChainSelector].CCIPAddress
	config.FeeQuoterApplyDestChainConfigUpdatesInput.StateObjectId = state[config.SuiChainSelector].CCIPObjectRef
	config.FeeQuoterApplyDestChainConfigUpdatesInput.OwnerCapObjectId = state[config.SuiChainSelector].CCIPOwnerCapObjectId
	reportFeeQuoterApplyDestChainConfigUpdatesOp, err := operations.ExecuteOperation(e.OperationsBundle, ccip_ops.FeeQuoterApplyDestChainConfigUpdatesOp, deps, config.FeeQuoterApplyDestChainConfigUpdatesInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to run FeeQuoterApplyDestChainConfigUpdatesOp for Sui chain %d: %w", config.SuiChainSelector, err)
	}
	seqReports = append(seqReports, []operations.Report[any, any]{reportFeeQuoterApplyDestChainConfigUpdatesOp.ToGenericReport()}...)

	config.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput.CCIPPackageId = state[config.SuiChainSelector].CCIPAddress
	config.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput.StateObjectId = state[config.SuiChainSelector].CCIPObjectRef
	config.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput.OwnerCapObjectId = state[config.SuiChainSelector].CCIPOwnerCapObjectId
	reportFeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp, err := operations.ExecuteOperation(e.OperationsBundle, ccip_ops.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp, deps, config.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput)
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
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to run ConfigureRouterOp for Sui chain %d: %w", config.SuiChainSelector, err)
	}
	seqReports = append(seqReports, []operations.Report[any, any]{reportConfigureRouterOp.ToGenericReport()}...)

	return cldf.ChangesetOutput{
		Reports: seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d ConnectSuiToEVM) VerifyPreconditions(e cldf.Environment, config ConnectSuiToEVMConfig) error {
	return nil
}
