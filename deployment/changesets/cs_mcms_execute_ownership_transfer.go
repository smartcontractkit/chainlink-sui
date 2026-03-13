package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	lockreleasetokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_lock_release_token_pool"
	managedtokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_managed_token_pool"
	offrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"
	onrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_onramp"
	routerops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_router"
	usdctokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_usdc_token_pool"
	managedtokenops "github.com/smartcontractkit/chainlink-sui/deployment/ops/managed_token"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	ownershipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ownership"
)

var _ cldf.ChangeSetV2[MCMSExecuteTransferOwnershipInput] = MCMSExecuteTransferOwnership{}

type MCMSExecuteTransferOwnership struct{}

type MCMSExecuteTransferOwnershipInput struct {
	ChainSelector uint64 `json:"chainSelector" yaml:"chainSelector"`

	// IsFastCurse selects the fastcurse MCMS instance as the ownership target.
	// When false (default) the normal MCMS instance is used.
	IsFastCurse bool `json:"isFastCurse,omitempty" yaml:"isFastCurse,omitempty"`

	// Type of contracts to execute the transfer on
	MCMS                            bool   `json:"mcms,omitempty" yaml:"mcms,omitempty"`
	StateObject                     bool   `json:"state_object,omitempty" yaml:"state_object,omitempty"`
	OnRamp                          bool   `json:"onramp,omitempty" yaml:"onramp,omitempty"`
	OffRamp                         bool   `json:"offramp,omitempty" yaml:"offramp,omitempty"`
	Router                          bool   `json:"router,omitempty" yaml:"router,omitempty"`
	ManagedToken                    bool   `json:"managed_token,omitempty" yaml:"managed_token,omitempty"`
	UsdcTokenPool                   bool   `json:"usdc_token_pool,omitempty" yaml:"usdc_token_pool,omitempty"`
	BurnMintTokenPoolTokenSymbol    string `json:"burn_mint_token_pool,omitempty" yaml:"burn_mint_token_pool,omitempty"`
	LockReleaseTokenPoolTokenSymbol string `json:"lock_release_token_pool,omitempty" yaml:"lock_release_token_pool,omitempty"`
	ManagedTokenPoolTokenSymbol     string `json:"managed_token_pool,omitempty" yaml:"managed_token_pool,omitempty"`
	TypeArg                         string `json:"type_arg,omitempty" yaml:"type_arg,omitempty"`
}

func (d MCMSExecuteTransferOwnership) Apply(e cldf.Environment, config MCMSExecuteTransferOwnershipInput) (cldf.ChangesetOutput, error) {
	suiState, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	state := suiState[config.ChainSelector]
	mcmsFields := state.MCMSState(config.IsFastCurse)

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

	input := ownershipops.ExecuteOwnershipTransferToMcmsSeqInput{}

	// Populate the input fields
	if config.MCMS {
		input.MCMS = &mcmsops.MCMSExecuteTransferOwnershipInput{
			McmsPackageID:    mcmsFields.PackageID,
			OwnerCap:         mcmsFields.AccountOwnerCapObjectID,
			AccountObjectID:  mcmsFields.AccountStateObjectID,
			RegistryObjectID: mcmsFields.RegistryObjectID,
		}
	}

	if config.StateObject {
		input.StateObject = &ccipops.ExecuteOwnershipTransferToMcmsStateObjectInput{
			CCIPPackageId:         state.CCIPAddress,
			OwnerCapObjectId:      state.CCIPOwnerCapObjectId,
			CCIPObjectRefObjectId: state.CCIPObjectRef,
			RegistryObjectId:      mcmsFields.RegistryObjectID,
			To:                    mcmsFields.PackageID,
		}
	}

	if config.OnRamp {
		input.OnRamp = &onrampops.ExecuteOwnershipTransferToMcmsOnRampInput{
			OnRampPackageId:     state.OnRampAddress,
			OnRampRefObjectId:   state.CCIPObjectRef,
			OnRampStateObjectId: state.OnRampStateObjectId,
			OwnerCapObjectId:    state.OnRampOwnerCapObjectId,
			RegistryObjectId:    mcmsFields.RegistryObjectID,
			To:                  mcmsFields.PackageID,
		}
	}

	if config.OffRamp {
		input.OffRamp = &offrampops.ExecuteOwnershipTransferToMcmsOffRampInput{
			OffRampPackageId:     state.OffRampAddress,
			OffRampRefObjectId:   state.CCIPObjectRef,
			OffRampStateObjectId: state.OffRampStateObjectId,
			OwnerCapObjectId:     state.OffRampOwnerCapId,
			RegistryObjectId:     mcmsFields.RegistryObjectID,
			To:                   mcmsFields.PackageID,
		}
	}

	if config.Router {
		input.Router = &routerops.ExecuteOwnershipTransferToMcmsRouterInput{
			RouterPackageId:     state.CCIPRouterAddress,
			OwnerCapObjectId:    state.CCIPRouterOwnerCapObjectId,
			RouterStateObjectId: state.CCIPRouterStateObjectID,
			RegistryObjectId:    mcmsFields.RegistryObjectID,
			To:                  mcmsFields.PackageID,
		}
	}

	// TODO: Need typeargs support
	if config.ManagedToken {
		input.ManagedToken = &managedtokenops.ExecuteOwnershipTransferToMcmsManagedTokenInput{
			ManagedTokenPackageId: state.CCIPAddress,
			OwnerCapObjectId:      state.CCIPOwnerCapObjectId,
			RegistryObjectId:      mcmsFields.RegistryObjectID,
			To:                    mcmsFields.PackageID,
		}
	}

	if config.BurnMintTokenPoolTokenSymbol != "" {
		poolState, ok := state.BnMTokenPools[config.BurnMintTokenPoolTokenSymbol]
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("burn mint token pool not found: %s", config.BurnMintTokenPoolTokenSymbol)
		}
		input.BurnMintTokenPool = &burnminttokenpoolops.ExecuteOwnershipTransferToMcmsBurnMintTokenPoolInput{
			BurnMintTokenPoolPackageId: poolState.PackageID,
			TypeArgs:                   []string{config.TypeArg},
			StateObjectId:              poolState.StateObjectId,
			OwnerCapObjectId:           poolState.OwnerCapObjectId,
			RegistryObjectId:           mcmsFields.RegistryObjectID,
			To:                         mcmsFields.PackageID,
		}
	}

	if config.LockReleaseTokenPoolTokenSymbol != "" {
		poolState, ok := state.LnRTokenPools[config.LockReleaseTokenPoolTokenSymbol]
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("lock release token pool not found: %s", config.LockReleaseTokenPoolTokenSymbol)
		}
		input.LockReleaseTokenPool = &lockreleasetokenpoolops.ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolInput{
			LockReleaseTokenPoolPackageId: poolState.PackageID,
			TypeArgs:                      []string{config.TypeArg},
			StateObjectId:                 poolState.StateObjectId,
			OwnerCapObjectId:              poolState.OwnerCapObjectId,
			RegistryObjectId:              mcmsFields.RegistryObjectID,
			To:                            mcmsFields.PackageID,
		}
	}

	if config.ManagedTokenPoolTokenSymbol != "" {
		poolState, ok := state.ManagedTokenPools[config.ManagedTokenPoolTokenSymbol]
		if !ok {
			return cldf.ChangesetOutput{}, fmt.Errorf("managed token pool not found: %s", config.ManagedTokenPoolTokenSymbol)
		}
		input.ManagedTokenPool = &managedtokenpoolops.ExecuteOwnershipTransferToMcmsManagedTokenPoolInput{
			ManagedTokenPoolPackageId: poolState.PackageID,
			TypeArgs:                  []string{config.TypeArg},
			StateObjectId:             poolState.StateObjectId,
			OwnerCapObjectId:          poolState.OwnerCapObjectId,
			RegistryObjectId:          mcmsFields.RegistryObjectID,
			To:                        mcmsFields.PackageID,
		}
	}

	// TODO: not supported yet
	if config.UsdcTokenPool {
		input.UsdcTokenPool = &usdctokenpoolops.ExecuteOwnershipTransferToMcmsUsdcTokenPoolInput{}
		return cldf.ChangesetOutput{}, fmt.Errorf("usdc token pool ownership transfer not implemented yet")
	}

	// Execute the sequence
	_, err = cld_ops.ExecuteSequence(e.OperationsBundle, ownershipops.ExecuteOwnershipTransferToMcmsSequence, deps, input)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute sequence: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d MCMSExecuteTransferOwnership) VerifyPreconditions(e cldf.Environment, config MCMSExecuteTransferOwnershipInput) error {
	// Check that at least one contract type is selected
	if !config.MCMS && !config.StateObject && !config.OnRamp &&
		!config.OffRamp && !config.Router && !config.ManagedToken &&
		!config.UsdcTokenPool && config.LockReleaseTokenPoolTokenSymbol != "" &&
		config.ManagedTokenPoolTokenSymbol != "" && config.BurnMintTokenPoolTokenSymbol != "" {
		return fmt.Errorf("at least one contract type must be selected for ownership transfer")
	}
	return nil
}
