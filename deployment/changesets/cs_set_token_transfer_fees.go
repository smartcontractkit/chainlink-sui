package changesets

import (
	"fmt"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccip_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

// SetTokenTransferFeesConfig sets per-token, per-destination-chain transfer fee config on the
// Sui CCIP FeeQuoter. Token object IDs are coin-metadata object IDs, matching the fee-token
// identifier convention. Each Add* slice must be the same length as AddTokens.
type SetTokenTransferFeesConfig struct {
	SuiChainSelector     uint64   `yaml:"suiChainSelector"`
	DestChainSelector    uint64   `yaml:"destChainSelector"`
	AddTokens            []string `yaml:"addTokens"`
	AddMinFeeUsdCents    []uint32 `yaml:"addMinFeeUsdCents"`
	AddMaxFeeUsdCents    []uint32 `yaml:"addMaxFeeUsdCents"`
	AddDeciBps           []uint16 `yaml:"addDeciBps"`
	AddDestGasOverhead   []uint32 `yaml:"addDestGasOverhead"`
	AddDestBytesOverhead []uint32 `yaml:"addDestBytesOverhead"`
	AddIsEnabled         []bool   `yaml:"addIsEnabled"`
	RemoveTokens         []string `yaml:"removeTokens"`
	// If non-nil, the operation is built as a call payload and wrapped into an MCMS
	// timelock proposal instead of being signed and broadcast directly by the EOA.
	TimelockConfig *utils.TimelockConfig `yaml:"timelockConfig,omitempty"`
}

// SetTokenTransferFees applies token transfer fee config updates to the Sui CCIP FeeQuoter.
type SetTokenTransferFees struct{}

var _ cldf.ChangeSetV2[SetTokenTransferFeesConfig] = SetTokenTransferFees{}

// Apply implements deployment.ChangeSetV2.
func (d SetTokenTransferFees) Apply(e cldf.Environment, config SetTokenTransferFeesConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	ds := fdatastore.NewMemoryDataStore()
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	suiChain := e.BlockChains.SuiChains()[config.SuiChainSelector]

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

	// Nil-out the signer so the op only builds the call payload; the resulting batch is
	// wrapped into an MCMS proposal below.
	if config.TimelockConfig != nil {
		deps.Signer = nil
	}

	report, err := operations.ExecuteOperation(e.OperationsBundle, ccip_ops.FeeQuoterApplyTokenTransferFeeConfigUpdatesOp, deps, ccip_ops.FeeQuoterApplyTokenTransferFeeConfigUpdatesInput{
		CCIPPackageId:        state[suiChain.Selector].CCIPAddress,
		LatestPackageId:      state[suiChain.Selector].EffectiveCCIPPackageID(),
		StateObjectId:        state[suiChain.Selector].CCIPObjectRef,
		OwnerCapObjectId:     state[suiChain.Selector].CCIPOwnerCapObjectId,
		DestChainSelector:    config.DestChainSelector,
		AddTokens:            config.AddTokens,
		AddMinFeeUsdCents:    config.AddMinFeeUsdCents,
		AddMaxFeeUsdCents:    config.AddMaxFeeUsdCents,
		AddDeciBps:           config.AddDeciBps,
		AddDestGasOverhead:   config.AddDestGasOverhead,
		AddDestBytesOverhead: config.AddDestBytesOverhead,
		AddIsEnabled:         config.AddIsEnabled,
		RemoveTokens:         config.RemoveTokens,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to apply token transfer fee config updates for Sui chain %d dest %d: %w", config.SuiChainSelector, config.DestChainSelector, err)
	}

	mcmsProposal := mcms.TimelockProposal{}
	if config.TimelockConfig != nil {
		mcmsConfig := mcmsops.ProposalGenerateInput{
			ChainSelector:      config.SuiChainSelector,
			Defs:               []operations.Definition{report.Def},
			Inputs:             []any{report.Input},
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
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build MCMS proposal for token transfer fees: %w", err)
		}
		mcmsProposal = result.Output
	}

	return cldf.ChangesetOutput{
		AddressBook:           ab,
		DataStore:             ds,
		Reports:               []operations.Report[any, any]{report.ToGenericReport()},
		MCMSTimelockProposals: []mcms.TimelockProposal{mcmsProposal},
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d SetTokenTransferFees) VerifyPreconditions(e cldf.Environment, config SetTokenTransferFeesConfig) error {
	n := len(config.AddTokens)
	if len(config.AddMinFeeUsdCents) != n {
		return fmt.Errorf("AddTokens (%d) and AddMinFeeUsdCents (%d) must have the same length", n, len(config.AddMinFeeUsdCents))
	}
	if len(config.AddMaxFeeUsdCents) != n {
		return fmt.Errorf("AddTokens (%d) and AddMaxFeeUsdCents (%d) must have the same length", n, len(config.AddMaxFeeUsdCents))
	}
	if len(config.AddDeciBps) != n {
		return fmt.Errorf("AddTokens (%d) and AddDeciBps (%d) must have the same length", n, len(config.AddDeciBps))
	}
	if len(config.AddDestGasOverhead) != n {
		return fmt.Errorf("AddTokens (%d) and AddDestGasOverhead (%d) must have the same length", n, len(config.AddDestGasOverhead))
	}
	if len(config.AddDestBytesOverhead) != n {
		return fmt.Errorf("AddTokens (%d) and AddDestBytesOverhead (%d) must have the same length", n, len(config.AddDestBytesOverhead))
	}
	if len(config.AddIsEnabled) != n {
		return fmt.Errorf("AddTokens (%d) and AddIsEnabled (%d) must have the same length", n, len(config.AddIsEnabled))
	}
	return nil
}
