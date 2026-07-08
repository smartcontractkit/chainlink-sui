package changesets

import (
	"fmt"
	"math/big"

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

// SeedDestChainPricesConfig bootstraps FeeQuoter price entries on a Sui source
// chain. Use this after a dest chain has been connected (via ConnectSuiToEVM /
// ConfigureLaneLegAsSource) but before the DON's price pusher publishes its
// first update — otherwise `ccip::onramp::get_fee` for that dest chain aborts
// with `fee_quoter::EUnknownDestChainSelector` (code 3), because
// `usd_per_unit_gas_by_dest_chain` is empty for the newly-registered chain.
//
// New lanes wired through `ConfigureLaneLegAsSource` already seed a price
// automatically from `input.Dest.GasPrice`; this changeset is the manual
// escape hatch for lanes that were connected without a price (e.g. Avalanche
// Fuji as of 2026-07 on Sui testnet), or for topping up source-token prices
// after a token addition.
type SeedDestChainPricesConfig struct {
	SuiChainSelector uint64

	// Optional source-token prices — USD per 1e18 of the smallest denomination
	// (18 decimals). Leave nil / empty to only refresh gas prices. Length of
	// SourceTokens must match SourceUsdPerToken.
	SourceTokens      []string
	SourceUsdPerToken []*big.Int

	// Destination gas prices — USD per unit gas (18 decimals) keyed by dest
	// chain selector. Length of GasDestChainSelectors must match GasUsdPerUnitGas.
	// Passing at least one entry here is the whole point of this changeset —
	// leaving both empty is a no-op.
	GasDestChainSelectors []uint64
	GasUsdPerUnitGas      []*big.Int

	// If non-nil, transactions are batched into an MCMS timelock proposal
	// instead of being signed and broadcast directly.
	TimelockConfig *utils.TimelockConfig
}

// SeedDestChainPrices seeds initial FeeQuoter prices on Sui.
type SeedDestChainPrices struct{}

var _ cldf.ChangeSetV2[SeedDestChainPricesConfig] = SeedDestChainPrices{}

// VerifyPreconditions implements cldf.ChangeSetV2.
func (SeedDestChainPrices) VerifyPreconditions(_ cldf.Environment, config SeedDestChainPricesConfig) error {
	if len(config.SourceTokens) != len(config.SourceUsdPerToken) {
		return fmt.Errorf(
			"SourceTokens (%d) and SourceUsdPerToken (%d) must have the same length",
			len(config.SourceTokens), len(config.SourceUsdPerToken),
		)
	}
	if len(config.GasDestChainSelectors) != len(config.GasUsdPerUnitGas) {
		return fmt.Errorf(
			"GasDestChainSelectors (%d) and GasUsdPerUnitGas (%d) must have the same length",
			len(config.GasDestChainSelectors), len(config.GasUsdPerUnitGas),
		)
	}
	if len(config.SourceTokens) == 0 && len(config.GasDestChainSelectors) == 0 {
		return fmt.Errorf("SeedDestChainPrices called with nothing to update: provide SourceTokens or GasDestChainSelectors")
	}
	for i, price := range config.SourceUsdPerToken {
		if price == nil {
			return fmt.Errorf("SourceUsdPerToken[%d] is nil", i)
		}
	}
	for i, price := range config.GasUsdPerUnitGas {
		if price == nil {
			return fmt.Errorf("GasUsdPerUnitGas[%d] is nil", i)
		}
	}
	return nil
}

// Apply implements cldf.ChangeSetV2.
func (SeedDestChainPrices) Apply(e cldf.Environment, config SeedDestChainPricesConfig) (cldf.ChangesetOutput, error) {
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

	// Nil-out the signer so the op only builds the call payload — the resulting
	// batch is then wrapped into an MCMS proposal below.
	if config.TimelockConfig != nil {
		deps.Signer = nil
	}

	input := ccip_ops.FeeQuoterUpdatePricesWithOwnerCapInput{
		CCIPPackageId:         state[config.SuiChainSelector].CCIPAddress,
		CCIPObjectRef:         state[config.SuiChainSelector].CCIPObjectRef,
		OwnerCapObjectId:      state[config.SuiChainSelector].CCIPOwnerCapObjectId,
		SourceTokens:          config.SourceTokens,
		SourceUsdPerToken:     config.SourceUsdPerToken,
		GasDestChainSelectors: config.GasDestChainSelectors,
		GasUsdPerUnitGas:      config.GasUsdPerUnitGas,
	}

	report, err := operations.ExecuteOperation(
		e.OperationsBundle,
		ccip_ops.FeeQuoterUpdatePricesWithOwnerCapOp,
		deps,
		input,
	)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf(
			"failed to run FeeQuoterUpdatePricesWithOwnerCapOp on Sui chain %d: %w",
			config.SuiChainSelector, err,
		)
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
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to build MCMS proposal for price seeding: %w", err)
		}
		mcmsProposal = result.Output
	}

	return cldf.ChangesetOutput{
		Reports:               []operations.Report[any, any]{report.ToGenericReport()},
		MCMSTimelockProposals: []mcms.TimelockProposal{mcmsProposal},
	}, nil
}
