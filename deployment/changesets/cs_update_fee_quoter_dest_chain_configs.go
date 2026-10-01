package changesets

import (
	"fmt"
	"maps"
	"slices"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/mcms"

	lanes "github.com/smartcontractkit/chainlink-ccip/deployment/lanes"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	suilanes "github.com/smartcontractkit/chainlink-sui/deployment/lanes"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccip_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

// DestChainConfigOverride expresses optional spot-replacements applied on top of a
// destination chain's CURRENT on-chain FeeQuoter config. Nil fields pass the on-chain
// value through unchanged. This keeps the update surgical: only the overridden fields
// change, everything else is re-applied byte-identically. (The Move setter
// apply_dest_chain_config_updates replaces the whole 21-field config per dest chain
// selector, so a partial write would zero the omitted fields.)
type DestChainConfigOverride struct {
	IsEnabled                   *bool   `yaml:"isEnabled,omitempty"`
	DestGasOverhead             *uint32 `yaml:"destGasOverhead,omitempty"`
	DefaultTokenDestGasOverhead *uint32 `yaml:"defaultTokenDestGasOverhead,omitempty"`
	MaxDataBytes                *uint32 `yaml:"maxDataBytes,omitempty"`
	MaxPerMsgGasLimit           *uint32 `yaml:"maxPerMsgGasLimit,omitempty"`
	DefaultTxGasLimit           *uint32 `yaml:"defaultTxGasLimit,omitempty"`
	DefaultTokenFeeUSDCents     *uint16 `yaml:"defaultTokenFeeUSDCents,omitempty"`
	NetworkFeeUSDCents          *uint16 `yaml:"networkFeeUsdCents,omitempty"`
}

// apply sets every non-nil field on cfg, leaving nil fields untouched.
func (o DestChainConfigOverride) apply(cfg *lanes.FeeQuoterDestChainConfig) {
	if o.IsEnabled != nil {
		cfg.IsEnabled = *o.IsEnabled
	}
	if o.DestGasOverhead != nil {
		cfg.DestGasOverhead = *o.DestGasOverhead
	}
	if o.DefaultTokenDestGasOverhead != nil {
		cfg.DefaultTokenDestGasOverhead = *o.DefaultTokenDestGasOverhead
	}
	if o.MaxDataBytes != nil {
		cfg.MaxDataBytes = *o.MaxDataBytes
	}
	if o.MaxPerMsgGasLimit != nil {
		cfg.MaxPerMsgGasLimit = *o.MaxPerMsgGasLimit
	}
	if o.DefaultTxGasLimit != nil {
		cfg.DefaultTxGasLimit = *o.DefaultTxGasLimit
	}
	if o.DefaultTokenFeeUSDCents != nil {
		cfg.DefaultTokenFeeUSDCents = *o.DefaultTokenFeeUSDCents
	}
	if o.NetworkFeeUSDCents != nil {
		cfg.NetworkFeeUSDCents = *o.NetworkFeeUSDCents
	}
}

// UpdateFeeQuoterDestChainConfigsConfig updates dest-chain configs on the Sui CCIP
// FeeQuoter for one or more destination chains, keyed by dest chain selector.
//
// Because the Move setter is all-or-nothing, the changeset reads each dest chain's
// current on-chain config via DevInspect and re-applies it with only the overridden
// fields changed — a dest chain that was never configured is an error (the read
// aborts with fee_quoter::EUnknownDestChainSelector), since there is no meaningful
// default to merge into for a remote chain.
type UpdateFeeQuoterDestChainConfigsConfig struct {
	SuiChainSelector uint64                             `yaml:"suiChainSelector"`
	Overrides        map[uint64]DestChainConfigOverride `yaml:"overrides"`
	// If non-nil, transactions are built as call payloads and wrapped into an MCMS
	// timelock proposal instead of being signed and broadcast directly.
	TimelockConfig *utils.TimelockConfig `yaml:"timelockConfig,omitempty"`
}

// UpdateFeeQuoterDestChainConfigs updates dest-chain configs on the Sui CCIP FeeQuoter.
type UpdateFeeQuoterDestChainConfigs struct{}

var _ cldf.ChangeSetV2[UpdateFeeQuoterDestChainConfigsConfig] = UpdateFeeQuoterDestChainConfigs{}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d UpdateFeeQuoterDestChainConfigs) VerifyPreconditions(e cldf.Environment, config UpdateFeeQuoterDestChainConfigsConfig) error {
	if len(config.Overrides) == 0 {
		return fmt.Errorf("UpdateFeeQuoterDestChainConfigs called with nothing to update: provide at least one entry in Overrides")
	}
	if _, ok := e.BlockChains.SuiChains()[config.SuiChainSelector]; !ok {
		return fmt.Errorf("chain with selector %d not found in Sui chains", config.SuiChainSelector)
	}
	for dst := range config.Overrides {
		if dst == config.SuiChainSelector {
			return fmt.Errorf("dest chain selector %d cannot be the source chain selector", dst)
		}
	}
	if config.TimelockConfig != nil {
		if err := config.TimelockConfig.Validate(); err != nil {
			return fmt.Errorf("invalid TimelockConfig: %w", err)
		}
	}
	return nil
}

// Apply implements deployment.ChangeSetV2.
func (d UpdateFeeQuoterDestChainConfigs) Apply(e cldf.Environment, config UpdateFeeQuoterDestChainConfigsConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	ds := fdatastore.NewMemoryDataStore()
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}
	chainState := state[config.SuiChainSelector]

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

	// DevInspect reads target the upgraded package head (current state) and require a
	// non-nil signer, so they use suiChain.Signer directly rather than deps.Signer.
	readPkg := chainState.EffectiveCCIPPackageID()
	fq, err := module_fee_quoter.NewFeeQuoter(readPkg, suiChain.Client)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to instantiate FeeQuoter at %s on Sui chain %d: %w", readPkg, config.SuiChainSelector, err)
	}

	defs := make([]operations.Definition, 0, len(config.Overrides))
	inputs := make([]any, 0, len(config.Overrides))
	reports := make([]operations.Report[any, any], 0, len(config.Overrides))

	// Iterate dest selectors in sorted order so direct execution and any proposal built
	// from the calls are byte-deterministic across runs.
	for _, dst := range slices.Sorted(maps.Keys(config.Overrides)) {
		moveCfg, err := fq.DevInspect().GetDestChainConfig(
			e.GetContext(),
			&bind.CallOpts{Signer: suiChain.Signer},
			bind.Object{Id: chainState.CCIPObjectRef},
			dst,
		)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf(
				"failed to read current dest chain config for dst %d from FeeQuoter at %s on Sui chain %d (is the dest chain configured?): %w",
				dst, readPkg, config.SuiChainSelector, err)
		}
		laneCfg, err := suilanes.TranslateDestChainConfigFromMove(moveCfg)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to translate dest chain config for dst %d: %w", dst, err)
		}
		config.Overrides[dst].apply(&laneCfg)

		opInput := suilanes.TranslateDestChainConfig(laneCfg, dst)
		opInput.CCIPPackageId = chainState.CCIPAddress
		opInput.LatestPackageId = chainState.EffectiveCCIPPackageID()
		opInput.StateObjectId = chainState.CCIPObjectRef
		opInput.OwnerCapObjectId = chainState.CCIPOwnerCapObjectId

		report, err := operations.ExecuteOperation(e.OperationsBundle, ccip_ops.FeeQuoterApplyDestChainConfigUpdatesOp, deps, opInput)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to apply dest chain config update for Sui chain %d dst %d: %w", config.SuiChainSelector, dst, err)
		}
		defs = append(defs, report.Def)
		inputs = append(inputs, report.Input)
		reports = append(reports, report.ToGenericReport())
	}

	mcmsProposal := mcms.TimelockProposal{}
	if config.TimelockConfig != nil {
		mcmsConfig := mcmsops.ProposalGenerateInput{
			ChainSelector:      config.SuiChainSelector,
			Defs:               defs,
			Inputs:             inputs,
			MmcsPackageID:      chainState.MCMSPackageID,
			McmsStateObjID:     chainState.MCMSStateObjectID,
			TimelockObjID:      chainState.MCMSTimelockObjectID,
			AccountObjID:       chainState.MCMSAccountStateObjectID,
			RegistryObjID:      chainState.MCMSRegistryObjectID,
			DeployerStateObjID: chainState.MCMSDeployerStateObjectID,
			TimelockConfig:     *config.TimelockConfig,
		}
		result, err := operations.ExecuteSequence(e.OperationsBundle, mcmsops.MCMSDynamicProposalGenerateSeq, deps, mcmsConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build MCMS proposal for fee quoter dest chain config updates: %w", err)
		}
		mcmsProposal = result.Output
	}

	return cldf.ChangesetOutput{
		AddressBook:           ab,
		DataStore:             ds,
		Reports:               reports,
		MCMSTimelockProposals: []mcms.TimelockProposal{mcmsProposal},
	}, nil
}
