package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	lockreleasetokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_lock_release_token_pool"
	managedtokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_managed_token_pool"
	tokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_token_pool"
	coin_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/coin"
)

type DeployTPAndConfigureConfig struct {
	SuiChainSelector   uint64
	TokenPoolTypes     []deployment.TokenPoolType
	ManagedTPInput     managedtokenpoolops.DeployAndInitManagedTokenPoolInput
	LockReleaseTPInput lockreleasetokenpoolops.DeployAndInitLockReleaseTokenPoolInput
	BurnMintTpInput    burnminttokenpoolops.DeployAndInitBurnMintTokenPoolInput
}

// ConnectSuiToEVM connects sui chain with EVM
type DeployTPAndConfigure struct{}

var _ cldf.ChangeSetV2[DeployTPAndConfigureConfig] = DeployTPAndConfigure{}

// Apply implements deployment.ChangeSetV2.
func (d DeployTPAndConfigure) Apply(e cldf.Environment, config DeployTPAndConfigureConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
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

	// Populate state information for each token pool type
	for _, tokenPoolType := range config.TokenPoolTypes {
		switch tokenPoolType {
		case deployment.TokenPoolTypeBurnMint:
			config.BurnMintTpInput.CCIPPackageId = state[config.SuiChainSelector].CCIPAddress
			config.BurnMintTpInput.MCMSAddress = state[config.SuiChainSelector].MCMSPackageID
			// TODO: MCMSOwner address should come state
			config.BurnMintTpInput.MCMSOwnerAddress = deployerAddr
			config.BurnMintTpInput.CCIPObjectRefObjectId = state[config.SuiChainSelector].CCIPObjectRef
			config.BurnMintTpInput.TokenPoolAdministrator = deployerAddr
		case deployment.TokenPoolTypeLockRelease:
			config.LockReleaseTPInput.CCIPPackageId = state[config.SuiChainSelector].CCIPAddress
			config.LockReleaseTPInput.MCMSAddress = state[config.SuiChainSelector].MCMSPackageID
			config.LockReleaseTPInput.MCMSOwnerAddress = deployerAddr
			config.LockReleaseTPInput.CCIPObjectRefObjectId = state[config.SuiChainSelector].CCIPObjectRef
			config.LockReleaseTPInput.TokenPoolAdministrator = deployerAddr
		case deployment.TokenPoolTypeManaged:
			symbolReport, err := cld_ops.ExecuteOperation(e.OperationsBundle, coin_ops.GetCoinSymbolOp, deps, config.ManagedTPInput.CoinObjectTypeArg)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to get coin symbol: %w", err)
			}
			managedTokenState, ok := state[config.SuiChainSelector].ManagedTokens[symbolReport.Output.Symbol]
			if !ok {
				return cldf.ChangesetOutput{}, fmt.Errorf("managed token not found for coin object type arg: %s with symbol: %s", config.ManagedTPInput.CoinObjectTypeArg, symbolReport.Output.Symbol)
			}
			config.ManagedTPInput.CCIPPackageId = state[config.SuiChainSelector].CCIPAddress
			config.ManagedTPInput.ManagedTokenPackageId = managedTokenState.PackageID
			config.ManagedTPInput.MCMSAddress = state[config.SuiChainSelector].MCMSPackageID
			config.ManagedTPInput.MCMSOwnerAddress = deployerAddr
			config.ManagedTPInput.CCIPObjectRefObjectId = state[config.SuiChainSelector].CCIPObjectRef
			config.ManagedTPInput.TokenPoolAdministrator = deployerAddr
		}
	}

	// Execute the unified token pool deployment sequence
	tokenPoolInput := tokenpoolops.DeployAndInitAllTokenPoolsInput{
		SuiChainSelector:   config.SuiChainSelector,
		TokenPoolTypes:     config.TokenPoolTypes,
		ManagedTPInput:     config.ManagedTPInput,
		LockReleaseTPInput: config.LockReleaseTPInput,
		BurnMintTpInput:    config.BurnMintTpInput,
	}

	tokenPoolReport, err := operations.ExecuteSequence(e.OperationsBundle, tokenpoolops.DeployAndInitAllTokenPoolsSequence, deps, tokenPoolInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy token pools: %w", err)
	}

	// Save addresses to the address book based on what was deployed
	for _, tokenPoolType := range config.TokenPoolTypes {
		switch tokenPoolType {
		case deployment.TokenPoolTypeBurnMint:
			// save BnM Pool to the addressbook
			typeAndVersionBurnMintTokenPool := cldf.NewTypeAndVersion(deployment.SuiBnMTokenPoolType, deployment.Version1_0_0)
			typeAndVersionBurnMintTokenPool.AddLabel(tokenPoolReport.Output.DeployBurnMintTokenPoolOutput.TokenSymbol)
			err = ab.Save(config.SuiChainSelector, tokenPoolReport.Output.BurnMintTPPackageID, typeAndVersionBurnMintTokenPool)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save BnMTokenPool address %s for Sui chain %d: %w", tokenPoolReport.Output.BurnMintTPPackageID, config.SuiChainSelector, err)
			}

			// save BnM Pool State to the addressBook
			typeAndVersionBurnMintTokenPoolState := cldf.NewTypeAndVersion(deployment.SuiBnMTokenPoolStateType, deployment.Version1_0_0)
			typeAndVersionBurnMintTokenPoolState.AddLabel(tokenPoolReport.Output.DeployBurnMintTokenPoolOutput.TokenSymbol)
			err = ab.Save(config.SuiChainSelector, tokenPoolReport.Output.DeployBurnMintTokenPoolOutput.Objects.StateObjectId, typeAndVersionBurnMintTokenPoolState)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save BnMTokenPoolState address %s for Sui chain %d: %w", tokenPoolReport.Output.DeployBurnMintTokenPoolOutput.Objects.StateObjectId, config.SuiChainSelector, err)
			}

			// save BnM Pool OwnerId to the addressBook
			typeAndVersionBurnMintTokenPoolOwnerId := cldf.NewTypeAndVersion(deployment.SuiBnMTokenPoolOwnerIDType, deployment.Version1_0_0)
			typeAndVersionBurnMintTokenPoolOwnerId.AddLabel(tokenPoolReport.Output.DeployBurnMintTokenPoolOutput.TokenSymbol)
			err = ab.Save(config.SuiChainSelector, tokenPoolReport.Output.DeployBurnMintTokenPoolOutput.Objects.OwnerCapObjectId, typeAndVersionBurnMintTokenPoolOwnerId)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save BnMTokenPoolOwnerCapId address %s for Sui chain %d: %w", tokenPoolReport.Output.DeployBurnMintTokenPoolOutput.Objects.OwnerCapObjectId, config.SuiChainSelector, err)
			}

		case deployment.TokenPoolTypeLockRelease:
			// save LnR Pool to the addressbook
			typeAndVersionLnRTokenPool := cldf.NewTypeAndVersion(deployment.SuiLnRTokenPoolType, deployment.Version1_0_0)
			typeAndVersionLnRTokenPool.AddLabel(tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.TokenSymbol)
			err = ab.Save(config.SuiChainSelector, tokenPoolReport.Output.LockReleaseTPPackageID, typeAndVersionLnRTokenPool)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LnRTokenPool address %s for Sui chain %d: %w", tokenPoolReport.Output.LockReleaseTPPackageID, config.SuiChainSelector, err)
			}

			// save LnR Pool State to the addressBook
			typeAndVersionLnRTokenPoolState := cldf.NewTypeAndVersion(deployment.SuiLnRTokenPoolStateType, deployment.Version1_0_0)
			typeAndVersionLnRTokenPoolState.AddLabel(tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.TokenSymbol)
			err = ab.Save(config.SuiChainSelector, tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.Objects.StateObjectId, typeAndVersionLnRTokenPoolState)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LnRTokenPoolState address %s for Sui chain %d: %w", tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.Objects.StateObjectId, config.SuiChainSelector, err)
			}

			// save LnR Pool OwnerId to the addressBook
			typeAndVersionLnRTokenPoolOwnerId := cldf.NewTypeAndVersion(deployment.SuiLnRTokenPoolOwnerIDType, deployment.Version1_0_0)
			typeAndVersionLnRTokenPoolOwnerId.AddLabel(tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.TokenSymbol)
			err = ab.Save(config.SuiChainSelector, tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.Objects.OwnerCapObjectId, typeAndVersionLnRTokenPoolOwnerId)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LnRTokenPoolOwnerCapId address %s for Sui chain %d: %w", tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.Objects.OwnerCapObjectId, config.SuiChainSelector, err)
			}

			// save LnR Pool RebalancerCapId to the addressBook
			typeAndVersionLnRTokenPoolRebalancerCapId := cldf.NewTypeAndVersion(deployment.SuiLnRTokenPoolRebalancerCapIDType, deployment.Version1_0_0)
			typeAndVersionLnRTokenPoolRebalancerCapId.AddLabel(tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.TokenSymbol)
			err = ab.Save(config.SuiChainSelector, tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.Objects.RebalancerCapObjectId, typeAndVersionLnRTokenPoolRebalancerCapId)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LnRTokenPoolRebalancerCapId address %s for Sui chain %d: %w", tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.Objects.RebalancerCapObjectId, config.SuiChainSelector, err)
			}
		case deployment.TokenPoolTypeManaged:
			// save Managed Pool to the addressbook
			typeAndVersionManagedTokenPool := cldf.NewTypeAndVersion(deployment.SuiManagedTokenPoolType, deployment.Version1_0_0)
			typeAndVersionManagedTokenPool.AddLabel(tokenPoolReport.Output.DeployManagedTokenPoolOutput.TokenSymbol)
			err = ab.Save(config.SuiChainSelector, tokenPoolReport.Output.ManagedTPPackageId, typeAndVersionManagedTokenPool)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save ManagedTokenPool address %s for Sui chain %d: %w", tokenPoolReport.Output.ManagedTPPackageId, config.SuiChainSelector, err)
			}

			// save Managed Pool State to the addressBook
			typeAndVersionManagedTokenPoolState := cldf.NewTypeAndVersion(deployment.SuiManagedTokenPoolStateType, deployment.Version1_0_0)
			typeAndVersionManagedTokenPoolState.AddLabel(tokenPoolReport.Output.DeployManagedTokenPoolOutput.TokenSymbol)
			err = ab.Save(config.SuiChainSelector, tokenPoolReport.Output.DeployManagedTokenPoolOutput.Objects.StateObjectId, typeAndVersionManagedTokenPoolState)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save ManagedTokenPoolState address %s for Sui chain %d: %w", tokenPoolReport.Output.DeployManagedTokenPoolOutput.Objects.StateObjectId, config.SuiChainSelector, err)
			}

			// save Managed Pool OwnerId to the addressBook
			typeAndVersionManagedTokenPoolOwnerId := cldf.NewTypeAndVersion(deployment.SuiManagedTokenPoolOwnerIDType, deployment.Version1_0_0)
			typeAndVersionManagedTokenPoolOwnerId.AddLabel(tokenPoolReport.Output.DeployManagedTokenPoolOutput.TokenSymbol)
			err = ab.Save(config.SuiChainSelector, tokenPoolReport.Output.DeployManagedTokenPoolOutput.Objects.OwnerCapObjectId, typeAndVersionManagedTokenPoolOwnerId)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save ManagedTokenPoolOwnerCapId address %s for Sui chain %d: %w", tokenPoolReport.Output.DeployManagedTokenPoolOutput.Objects.OwnerCapObjectId, config.SuiChainSelector, err)
			}
		}
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeployTPAndConfigure) VerifyPreconditions(e cldf.Environment, config DeployTPAndConfigureConfig) error {
	return nil
}
