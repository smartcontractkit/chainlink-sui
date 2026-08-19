package changesets

import (
	"fmt"

	"github.com/smartcontractkit/mcms"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	lockreleasetokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_lock_release_token_pool"
	managedtokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_managed_token_pool"
	ownershipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ownership"
	opregistry "github.com/smartcontractkit/chainlink-sui/deployment/ops/registry"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

var _ cldf.ChangeSetV2[AcceptOwnershipTokenPoolConfig] = AcceptOwnershipTokenPool{}

// AcceptOwnershipTokenPool generates an MCMS timelock proposal that accepts
// ownership of one or more deployed Sui token pools (managed / burn-mint /
// lock-release). It mirrors AcceptOwnershipCCIP, which only covers the CCIP
// core contracts Router/OnRamp/OffRamp/StateObject.
//
// The proposal must be signed and executed through the target MCMS instance's
// timelock before execute_ownership_transfer_to_mcms can run.
type AcceptOwnershipTokenPool struct{}

type AcceptOwnershipTokenPoolConfig struct {
	ChainSelector  uint64               `json:"chainSelector" yaml:"chainSelector"`
	IsFastCurse    bool                 `json:"isFastCurse,omitempty" yaml:"isFastCurse,omitempty"`
	TimelockConfig utils.TimelockConfig `json:"timelockConfig" yaml:"timelockConfig"`

	// Select the token pool(s) to accept. TypeArg is applied to every selected
	// pool, matching MCMSExecuteTransferOwnership's single type-arg convention.
	ManagedTokenPoolTokenSymbol     string `json:"managed_token_pool,omitempty" yaml:"managed_token_pool,omitempty"`
	BurnMintTokenPoolTokenSymbol    string `json:"burn_mint_token_pool,omitempty" yaml:"burn_mint_token_pool,omitempty"`
	LockReleaseTokenPoolTokenSymbol string `json:"lock_release_token_pool,omitempty" yaml:"lock_release_token_pool,omitempty"`
	TypeArg                         string `json:"type_arg,omitempty" yaml:"type_arg,omitempty"`
}

// Apply implements deployment.ChangeSetV2.
func (d AcceptOwnershipTokenPool) Apply(e cldf.Environment, config AcceptOwnershipTokenPoolConfig) (cldf.ChangesetOutput, error) {
	suiState, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	state, ok := suiState[config.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("no Sui chain state for chain selector %d", config.ChainSelector)
	}

	mcmsFields := state.MCMSState(config.IsFastCurse)

	suiChain := e.BlockChains.SuiChains()[config.ChainSelector]
	deps := sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: suiChain.Signer,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(1_000_000_000)
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
		SuiRPC: suiChain.URL,
	}

	// Ensure the registry holds every operation so the token-pool accept ops are retrievable.
	for i := range opregistry.AllOperations {
		cld_ops.RegisterOperation(e.OperationsBundle.OperationRegistry, opregistry.AllOperations[i])
	}

	seqInput := ownershipops.AcceptTokenPoolOwnershipInput{
		MCMSPackageId:          mcmsFields.PackageID,
		MCMSStateObjId:         mcmsFields.StateObjectID,
		MCMSTimelockObjId:      mcmsFields.TimelockObjectID,
		MCMSAccountObjId:       mcmsFields.AccountStateObjectID,
		MCMSRegistryObjId:      mcmsFields.RegistryObjectID,
		MCMSDeployerStateObjId: mcmsFields.DeployerStateObjectID,
		ChainSelector:          config.ChainSelector,
		TimelockConfig:         config.TimelockConfig,
	}

	if config.ManagedTokenPoolTokenSymbol != "" {
		pool, ok := state.ManagedTokenPools[config.ManagedTokenPoolTokenSymbol]
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("managed token pool not found: %s", config.ManagedTokenPoolTokenSymbol)
		}
		seqInput.ManagedTokenPool = &managedtokenpoolops.AcceptOwnershipManagedTokenPoolInput{
			ManagedTokenPoolPackageId: pool.PackageID,
			TypeArgs:                  []string{config.TypeArg},
			StateObjectId:             pool.StateObjectId,
		}
	}

	if config.BurnMintTokenPoolTokenSymbol != "" {
		pool, ok := state.BnMTokenPools[config.BurnMintTokenPoolTokenSymbol]
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("burn mint token pool not found: %s", config.BurnMintTokenPoolTokenSymbol)
		}
		seqInput.BurnMintTokenPool = &burnminttokenpoolops.AcceptOwnershipBurnMintTokenPoolInput{
			BurnMintTokenPoolPackageId: pool.PackageID,
			TypeArgs:                   []string{config.TypeArg},
			StateObjectId:              pool.StateObjectId,
		}
	}

	if config.LockReleaseTokenPoolTokenSymbol != "" {
		pool, ok := state.LnRTokenPools[config.LockReleaseTokenPoolTokenSymbol]
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("lock release token pool not found: %s", config.LockReleaseTokenPoolTokenSymbol)
		}
		seqInput.LockReleaseTokenPool = &lockreleasetokenpoolops.AcceptOwnershipLockReleaseTokenPoolInput{
			LockReleaseTokenPoolPackageId: pool.PackageID,
			TypeArgs:                      []string{config.TypeArg},
			StateObjectId:                 pool.StateObjectId,
		}
	}

	report, err := cld_ops.ExecuteSequence(e.OperationsBundle, ownershipops.AcceptTokenPoolOwnershipSeq, deps, seqInput)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{report.Output},
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d AcceptOwnershipTokenPool) VerifyPreconditions(e cldf.Environment, config AcceptOwnershipTokenPoolConfig) error {
	if config.ChainSelector == 0 {
		return fmt.Errorf("chainSelector is required")
	}
	// CCIP token-pool ownership must go to the slow MCMS instance, matching
	// MCMSExecuteTransferOwnership which forbids fastcurse receiving CCIP ownership.
	if config.IsFastCurse {
		return fmt.Errorf("fastcurse MCMS cannot receive CCIP ownership transfer; CCIP token-pool OwnerCap must remain with slow MCMS")
	}
	if config.ManagedTokenPoolTokenSymbol == "" && config.BurnMintTokenPoolTokenSymbol == "" && config.LockReleaseTokenPoolTokenSymbol == "" {
		return fmt.Errorf("at least one token pool symbol must be selected for accept ownership")
	}
	return nil
}
