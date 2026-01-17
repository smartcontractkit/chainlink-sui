package changesets

import (
	"fmt"

	"github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	onrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_onramp"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type WithdrawFeeTokensOnRampConfig struct {
	ChainSelector      uint64
	FeeTokenMetadataId string
	TypeArg            string
	TimelockConfig     *utils.TimelockConfig // If nil, execute directly; otherwise generate proposal
}

var _ cldf.ChangeSetV2[WithdrawFeeTokensOnRampConfig] = WithdrawFeeTokensOnRamp{}

type WithdrawFeeTokensOnRamp struct{}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (WithdrawFeeTokensOnRamp) VerifyPreconditions(e cldf.Environment, config WithdrawFeeTokensOnRampConfig) error {
	// Verify that the chain selector is valid
	if config.ChainSelector == 0 {
		return fmt.Errorf("chain selector is required")
	}

	// Verify that the fee token metadata id is valid
	if config.FeeTokenMetadataId == "" {
		return fmt.Errorf("fee token metadata id is required")
	}

	// Verify that the type arg is valid
	if config.TypeArg == "" {
		return fmt.Errorf("type arg is required")
	}

	return nil
}

// Apply implements deployment.ChangeSetV2.
func (WithdrawFeeTokensOnRamp) Apply(e cldf.Environment, config WithdrawFeeTokensOnRampConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]cld_ops.Report[any, any], 0)

	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
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

	// If timelock proposal is to be generated, disable signer in deps
	if config.TimelockConfig != nil {
		deps.Signer = nil
	}

	reportWithdraw, err := cld_ops.ExecuteOperation(e.OperationsBundle, onrampops.WithdrawFeeTokensOp, deps, onrampops.WithdrawFeeTokensInput{
		OnRampPackageId:    state[config.ChainSelector].OnRampAddress,
		CCIPObjectRefId:    state[config.ChainSelector].CCIPObjectRef,
		StateObjectId:      state[config.ChainSelector].OnRampStateObjectId,
		OwnerCapObjectId:   state[config.ChainSelector].OnRampOwnerCapObjectId,
		FeeTokenMetadataId: config.FeeTokenMetadataId,
		TypeArg:            config.TypeArg,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to withdraw fee tokens for Sui chain %d: %w", config.ChainSelector, err)
	}

	seqReports = append(seqReports, reportWithdraw.ToGenericReport())

	var timelockProposals []mcms.TimelockProposal
	if config.TimelockConfig != nil {
		defs := []cld_ops.Definition{reportWithdraw.Def}
		inputs := []any{reportWithdraw.Input}

		mcmsConfig := mcmsops.ProposalGenerateInput{
			ChainSelector:      config.ChainSelector,
			Defs:               defs,
			Inputs:             inputs,
			MmcsPackageID:      state[config.ChainSelector].MCMSPackageID,
			McmsStateObjID:     state[config.ChainSelector].MCMSStateObjectID,
			TimelockObjID:      state[config.ChainSelector].MCMSTimelockObjectID,
			AccountObjID:       state[config.ChainSelector].MCMSAccountStateObjectID,
			RegistryObjID:      state[config.ChainSelector].MCMSRegistryObjectID,
			DeployerStateObjID: state[config.ChainSelector].MCMSDeployerStateObjectID,
			TimelockConfig:     *config.TimelockConfig,
		}

		result, err := cld_ops.ExecuteSequence(e.OperationsBundle, mcmsops.MCMSDynamicProposalGenerateSeq, deps, mcmsConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal: %w", err)
		}
		timelockProposals = append(timelockProposals, result.Output)
	}

	return cldf.ChangesetOutput{
		AddressBook:           ab,
		Reports:               seqReports,
		MCMSTimelockProposals: timelockProposals,
	}, nil
}
