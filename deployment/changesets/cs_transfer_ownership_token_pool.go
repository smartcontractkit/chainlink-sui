package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	lockreleasetokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_lock_release_token_pool"
	managedtokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_managed_token_pool"
)

var _ cldf.ChangeSetV2[TransferOwnershipTokenPoolConfig] = TransferOwnershipTokenPool{}

// TransferOwnershipTokenPool initiates ownership transfer of one or more deployed
// Sui token pools (managed / burn-mint / lock-release) to the slow MCMS instance.
// This is step 1 of the 3-step Sui ownership-to-MCMS flow: it is deployer-signed
// and direct, setting the pending transfer (accepted=false) on each pool. The MCMS
// then accepts (step 2, AcceptOwnershipTokenPool) and the deployer finalizes via
// execute_ownership_transfer_to_mcms (step 3, MCMSExecuteTransferOwnership).
//
// Token-pool ownership always targets the slow MCMS instance, matching
// MCMSExecuteTransferOwnership which forbids fastcurse receiving CCIP ownership.
type TransferOwnershipTokenPool struct{}

type TransferOwnershipTokenPoolConfig struct {
	ChainSelector uint64 `json:"chainSelector" yaml:"chainSelector"`

	// Select the token pool(s) to transfer. TypeArg is applied to every selected
	// pool, matching MCMSExecuteTransferOwnership's single type-arg convention.
	ManagedTokenPoolTokenSymbol     string `json:"managed_token_pool,omitempty" yaml:"managed_token_pool,omitempty"`
	BurnMintTokenPoolTokenSymbol    string `json:"burn_mint_token_pool,omitempty" yaml:"burn_mint_token_pool,omitempty"`
	LockReleaseTokenPoolTokenSymbol string `json:"lock_release_token_pool,omitempty" yaml:"lock_release_token_pool,omitempty"`
	TypeArg                         string `json:"type_arg,omitempty" yaml:"type_arg,omitempty"`
}

// Apply implements deployment.ChangeSetV2.
func (d TransferOwnershipTokenPool) Apply(e cldf.Environment, config TransferOwnershipTokenPoolConfig) (cldf.ChangesetOutput, error) {
	suiState, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	state, ok := suiState[config.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("no Sui chain state for chain selector %d", config.ChainSelector)
	}

	// Ownership always transfers to the slow MCMS instance; its package id is the
	// MCMS multisig address (mcms_registry::get_multisig_address).
	mcmsFields := state.MCMSState(false)
	if mcmsFields.PackageID == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf("slow MCMS package id not found for chain selector %d", config.ChainSelector)
	}

	suiChain := e.BlockChains.SuiChains()[config.ChainSelector]
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

	reports := make([]cld_ops.Report[any, any], 0, 3)

	if config.ManagedTokenPoolTokenSymbol != "" {
		pool, ok := state.ManagedTokenPools[config.ManagedTokenPoolTokenSymbol]
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("managed token pool not found: %s", config.ManagedTokenPoolTokenSymbol)
		}
		report, err := cld_ops.ExecuteOperation(e.OperationsBundle, managedtokenpoolops.TransferOwnershipManagedTokenPoolOp, deps, managedtokenpoolops.TransferOwnershipManagedTokenPoolInput{
			ManagedTokenPoolPackageId: pool.PackageID,
			TypeArgs:                  []string{config.TypeArg},
			StateObjectId:             pool.StateObjectId,
			OwnerCapObjectId:          pool.OwnerCapObjectId,
			To:                        mcmsFields.PackageID,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to transfer managed token pool ownership: %w", err)
		}
		reports = append(reports, report.ToGenericReport())
	}

	if config.BurnMintTokenPoolTokenSymbol != "" {
		pool, ok := state.BnMTokenPools[config.BurnMintTokenPoolTokenSymbol]
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("burn mint token pool not found: %s", config.BurnMintTokenPoolTokenSymbol)
		}
		report, err := cld_ops.ExecuteOperation(e.OperationsBundle, burnminttokenpoolops.TransferOwnershipBurnMintTokenPoolOp, deps, burnminttokenpoolops.TransferOwnershipBurnMintTokenPoolInput{
			BurnMintTokenPoolPackageId: pool.PackageID,
			TypeArgs:                   []string{config.TypeArg},
			StateObjectId:              pool.StateObjectId,
			OwnerCapObjectId:           pool.OwnerCapObjectId,
			To:                         mcmsFields.PackageID,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to transfer burn mint token pool ownership: %w", err)
		}
		reports = append(reports, report.ToGenericReport())
	}

	if config.LockReleaseTokenPoolTokenSymbol != "" {
		pool, ok := state.LnRTokenPools[config.LockReleaseTokenPoolTokenSymbol]
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("lock release token pool not found: %s", config.LockReleaseTokenPoolTokenSymbol)
		}
		report, err := cld_ops.ExecuteOperation(e.OperationsBundle, lockreleasetokenpoolops.TransferOwnershipLockReleaseTokenPoolOp, deps, lockreleasetokenpoolops.TransferOwnershipLockReleaseTokenPoolInput{
			LockReleaseTokenPoolPackageId: pool.PackageID,
			TypeArgs:                      []string{config.TypeArg},
			StateObjectId:                 pool.StateObjectId,
			OwnerCapObjectId:              pool.OwnerCapObjectId,
			To:                            mcmsFields.PackageID,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to transfer lock release token pool ownership: %w", err)
		}
		reports = append(reports, report.ToGenericReport())
	}

	return cldf.ChangesetOutput{
		Reports: reports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d TransferOwnershipTokenPool) VerifyPreconditions(e cldf.Environment, config TransferOwnershipTokenPoolConfig) error {
	if config.ChainSelector == 0 {
		return fmt.Errorf("chainSelector is required")
	}
	if config.ManagedTokenPoolTokenSymbol == "" && config.BurnMintTokenPoolTokenSymbol == "" && config.LockReleaseTokenPoolTokenSymbol == "" {
		return fmt.Errorf("at least one token pool symbol must be selected for ownership transfer")
	}
	return nil
}
