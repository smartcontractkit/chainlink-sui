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
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	lockreleasetokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_lock_release_token_pool"
	managedtokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_managed_token_pool"
	tokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_token_pool"
	coin_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/coin"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type TPConfigureConfig struct {
	SuiChainSelector   uint64
	TokenPoolTypes     []string
	ManagedTPInput     managedtokenpoolops.ConfigureManagedTokenPoolInput
	LockReleaseTPInput lockreleasetokenpoolops.ConfigureLockReleaseTokenPoolInput
	BurnMintTpInput    burnminttokenpoolops.ConfigureBurnMintTokenPoolInput
	TimelockConfig     *utils.TimelockConfig
}

// ConnectSuiToEVM connects sui chain with EVM
type TPConfigure struct{}

var _ cldf.ChangeSetV2[TPConfigureConfig] = TPConfigure{}

// Apply implements deployment.ChangeSetV2.
func (d TPConfigure) Apply(e cldf.Environment, config TPConfigureConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	ds := fdatastore.NewMemoryDataStore()
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	seqReports := make([]operations.Report[any, any], 0)

	suiChains := e.BlockChains.SuiChains()
	suiChain := suiChains[config.SuiChainSelector]

	deployerAddr, err := suiChain.Signer.GetAddress()
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

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

	defs := []operations.Definition{}
	inputs := []any{}

	// Populate state information for each token pool type
	for _, tokenPoolType := range config.TokenPoolTypes {
		switch tokenPoolType {
		case "bnm":
			config.BurnMintTpInput.CCIPPackageId = state[config.SuiChainSelector].CCIPAddress
			config.BurnMintTpInput.MCMSAddress = state[config.SuiChainSelector].MCMSPackageID
			// TODO: MCMSOwner address should come state
			config.BurnMintTpInput.MCMSOwnerAddress = deployerAddr
		case "lnr":
			// Resolve the already-deployed pool's object ids from on-chain state by symbol,
			// mirroring AcceptOwnershipTokenPool. The pool is keyed by the coin symbol derived
			// from CoinObjectTypeArg, the same key used at deploy. The config ops are dual-mode:
			// EOA-direct when a Signer is set, MCMS-encoded operation when Signer is nil
			// (the TimelockConfig branch above). Remote/rate-limit data stay caller-supplied.
			symbolReport, err := operations.ExecuteOperation(e.OperationsBundle, coin_ops.GetCoinSymbolOp, deps, config.LockReleaseTPInput.CoinObjectTypeArg)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get coin symbol: %w", err)
			}
			pool, ok := state[config.SuiChainSelector].LnRTokenPools[symbolReport.Output.Symbol]
			if !ok {
				return cldf.ChangesetOutput{}, fmt.Errorf("lock release token pool not found for symbol %s", symbolReport.Output.Symbol)
			}
			config.LockReleaseTPInput.TokenPoolPkgID = pool.PackageID
			config.LockReleaseTPInput.StateObjectId = pool.StateObjectId
			config.LockReleaseTPInput.OwnerCap = pool.OwnerCapObjectId
		case "managed":
			config.ManagedTPInput.CCIPPackageId = state[config.SuiChainSelector].CCIPAddress
			config.ManagedTPInput.MCMSAddress = state[config.SuiChainSelector].MCMSPackageID
			// TODO: MCMSOwner address should come state
			config.ManagedTPInput.MCMSOwnerAddress = deployerAddr
		}
	}

	// Execute the unified token pool deployment sequence
	tokenPoolInput := tokenpoolops.ConfigureAllTokenPoolsInput{
		SuiChainSelector:   config.SuiChainSelector,
		TokenPoolTypes:     config.TokenPoolTypes,
		ManagedTPInput:     config.ManagedTPInput,
		LockReleaseTPInput: config.LockReleaseTPInput,
		BurnMintTpInput:    config.BurnMintTpInput,
	}

	report, err := operations.ExecuteSequence(e.OperationsBundle, tokenpoolops.ConfigureAllTokenPoolsSequence, deps, tokenPoolInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy token pools: %w", err)
	}

	for _, r := range report.Output.Reports {
		defs = append(defs, r.Def)
		inputs = append(inputs, r.Input)
	}

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
		AddressBook:           ab,
		DataStore:             ds,
		Reports:               seqReports,
		MCMSTimelockProposals: []mcms.TimelockProposal{mcmsProposal},
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d TPConfigure) VerifyPreconditions(e cldf.Environment, config TPConfigureConfig) error {
	return nil
}
