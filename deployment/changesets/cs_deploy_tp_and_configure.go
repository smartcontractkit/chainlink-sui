package changesets

import (
	"fmt"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
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
	// ReplaceExisting allows this changeset to take datastore keys that are already recorded,
	// as redeploying a pool for a token does. Without it, an occupied key is an error raised
	// before anything is deployed.
	ReplaceExisting bool `yaml:"replaceExisting"`
}

// ConnectSuiToEVM connects sui chain with EVM
type DeployTPAndConfigure struct{}

var _ cldf.ChangeSetV2[DeployTPAndConfigureConfig] = DeployTPAndConfigure{}

// Apply implements deployment.ChangeSetV2.
func (d DeployTPAndConfigure) Apply(e cldf.Environment, config DeployTPAndConfigureConfig) (cldf.ChangesetOutput, error) {
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

	fastMcmsPackageID := state[config.SuiChainSelector].FastCurseMCMSPackageID
	if fastMcmsPackageID == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf(
			"fast MCMS package not deployed for Sui chain %d; run DeploySuiChain first",
			config.SuiChainSelector,
		)
	}

	// Populate state information for each token pool type
	for _, tokenPoolType := range config.TokenPoolTypes {
		switch tokenPoolType {
		case deployment.TokenPoolTypeBurnMint:
			config.BurnMintTpInput.CCIPPackageId = state[config.SuiChainSelector].CCIPAddress
			config.BurnMintTpInput.MCMSAddress = state[config.SuiChainSelector].MCMSPackageID
			config.BurnMintTpInput.FastMcmsAddress = fastMcmsPackageID
			// TODO: MCMSOwner address should come state
			config.BurnMintTpInput.MCMSOwnerAddress = deployerAddr
			config.BurnMintTpInput.CCIPObjectRefObjectId = state[config.SuiChainSelector].CCIPObjectRef
			config.BurnMintTpInput.TokenPoolAdministrator = deployerAddr
		case deployment.TokenPoolTypeLockRelease:
			config.LockReleaseTPInput.CCIPPackageId = state[config.SuiChainSelector].CCIPAddress
			config.LockReleaseTPInput.MCMSAddress = state[config.SuiChainSelector].MCMSPackageID
			config.LockReleaseTPInput.FastMcmsAddress = fastMcmsPackageID
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
			config.ManagedTPInput.FastMcmsAddress = fastMcmsPackageID
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
			err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.SuiChainSelector, tokenPoolReport.Output.BurnMintTPPackageID, typeAndVersionBurnMintTokenPool, deployment.TokenQualifier(tokenPoolReport.Output.DeployBurnMintTokenPoolOutput.TokenSymbol)) // token-scoped; the label keeps the display name, the key uses the symbol form
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save BnMTokenPool address %s for Sui chain %d: %w", tokenPoolReport.Output.BurnMintTPPackageID, config.SuiChainSelector, err)
			}

			// save BnM Pool State to the addressBook
			typeAndVersionBurnMintTokenPoolState := cldf.NewTypeAndVersion(deployment.SuiBnMTokenPoolStateType, deployment.Version1_0_0)
			typeAndVersionBurnMintTokenPoolState.AddLabel(tokenPoolReport.Output.DeployBurnMintTokenPoolOutput.TokenSymbol)
			err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.SuiChainSelector, tokenPoolReport.Output.DeployBurnMintTokenPoolOutput.Objects.StateObjectId, typeAndVersionBurnMintTokenPoolState, deployment.TokenQualifier(tokenPoolReport.Output.DeployBurnMintTokenPoolOutput.TokenSymbol)) // token-scoped; the label keeps the display name, the key uses the symbol form
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save BnMTokenPoolState address %s for Sui chain %d: %w", tokenPoolReport.Output.DeployBurnMintTokenPoolOutput.Objects.StateObjectId, config.SuiChainSelector, err)
			}

			// save BnM Pool OwnerId to the addressBook
			typeAndVersionBurnMintTokenPoolOwnerId := cldf.NewTypeAndVersion(deployment.SuiBnMTokenPoolOwnerCapObjectIDType, deployment.Version1_0_0)
			typeAndVersionBurnMintTokenPoolOwnerId.AddLabel(tokenPoolReport.Output.DeployBurnMintTokenPoolOutput.TokenSymbol)
			err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.SuiChainSelector, tokenPoolReport.Output.DeployBurnMintTokenPoolOutput.Objects.OwnerCapObjectId, typeAndVersionBurnMintTokenPoolOwnerId, deployment.TokenQualifier(tokenPoolReport.Output.DeployBurnMintTokenPoolOutput.TokenSymbol)) // token-scoped; the label keeps the display name, the key uses the symbol form
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save BnMTokenPoolOwnerCapId address %s for Sui chain %d: %w", tokenPoolReport.Output.DeployBurnMintTokenPoolOutput.Objects.OwnerCapObjectId, config.SuiChainSelector, err)
			}

		case deployment.TokenPoolTypeLockRelease:
			// save LnR Pool to the addressbook
			typeAndVersionLnRTokenPool := cldf.NewTypeAndVersion(deployment.SuiLnRTokenPoolType, deployment.Version1_0_0)
			typeAndVersionLnRTokenPool.AddLabel(tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.TokenSymbol)
			err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.SuiChainSelector, tokenPoolReport.Output.LockReleaseTPPackageID, typeAndVersionLnRTokenPool, deployment.TokenQualifier(tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.TokenSymbol)) // token-scoped; the label keeps the display name, the key uses the symbol form
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LnRTokenPool address %s for Sui chain %d: %w", tokenPoolReport.Output.LockReleaseTPPackageID, config.SuiChainSelector, err)
			}

			// save LnR Pool State to the addressBook
			typeAndVersionLnRTokenPoolState := cldf.NewTypeAndVersion(deployment.SuiLnRTokenPoolStateType, deployment.Version1_0_0)
			typeAndVersionLnRTokenPoolState.AddLabel(tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.TokenSymbol)
			err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.SuiChainSelector, tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.Objects.StateObjectId, typeAndVersionLnRTokenPoolState, deployment.TokenQualifier(tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.TokenSymbol)) // token-scoped; the label keeps the display name, the key uses the symbol form
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LnRTokenPoolState address %s for Sui chain %d: %w", tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.Objects.StateObjectId, config.SuiChainSelector, err)
			}

			// save LnR Pool OwnerId to the addressBook
			typeAndVersionLnRTokenPoolOwnerId := cldf.NewTypeAndVersion(deployment.SuiLnRTokenPoolOwnerCapObjectIDType, deployment.Version1_0_0)
			typeAndVersionLnRTokenPoolOwnerId.AddLabel(tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.TokenSymbol)
			err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.SuiChainSelector, tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.Objects.OwnerCapObjectId, typeAndVersionLnRTokenPoolOwnerId, deployment.TokenQualifier(tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.TokenSymbol)) // token-scoped; the label keeps the display name, the key uses the symbol form
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LnRTokenPoolOwnerCapId address %s for Sui chain %d: %w", tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.Objects.OwnerCapObjectId, config.SuiChainSelector, err)
			}

			// save LnR Pool RebalancerCapId to the addressBook
			typeAndVersionLnRTokenPoolRebalancerCapId := cldf.NewTypeAndVersion(deployment.SuiLnRTokenPoolRebalancerCapIDType, deployment.Version1_0_0)
			typeAndVersionLnRTokenPoolRebalancerCapId.AddLabel(tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.TokenSymbol)
			err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.SuiChainSelector, tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.Objects.RebalancerCapObjectId, typeAndVersionLnRTokenPoolRebalancerCapId, deployment.TokenQualifier(tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.TokenSymbol)) // token-scoped; the label keeps the display name, the key uses the symbol form
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save LnRTokenPoolRebalancerCapId address %s for Sui chain %d: %w", tokenPoolReport.Output.DeployLockReleaseTokenPoolOutput.Objects.RebalancerCapObjectId, config.SuiChainSelector, err)
			}
		case deployment.TokenPoolTypeManaged:
			// save Managed Pool to the addressbook
			typeAndVersionManagedTokenPool := cldf.NewTypeAndVersion(deployment.SuiManagedTokenPoolType, deployment.Version1_0_0)
			typeAndVersionManagedTokenPool.AddLabel(tokenPoolReport.Output.DeployManagedTokenPoolOutput.TokenSymbol)
			err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.SuiChainSelector, tokenPoolReport.Output.ManagedTPPackageId, typeAndVersionManagedTokenPool, deployment.TokenQualifier(tokenPoolReport.Output.DeployManagedTokenPoolOutput.TokenSymbol)) // token-scoped; the label keeps the display name, the key uses the symbol form
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save ManagedTokenPool address %s for Sui chain %d: %w", tokenPoolReport.Output.ManagedTPPackageId, config.SuiChainSelector, err)
			}

			// save Managed Pool State to the addressBook
			typeAndVersionManagedTokenPoolState := cldf.NewTypeAndVersion(deployment.SuiManagedTokenPoolStateType, deployment.Version1_0_0)
			typeAndVersionManagedTokenPoolState.AddLabel(tokenPoolReport.Output.DeployManagedTokenPoolOutput.TokenSymbol)
			err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.SuiChainSelector, tokenPoolReport.Output.DeployManagedTokenPoolOutput.Objects.StateObjectId, typeAndVersionManagedTokenPoolState, deployment.TokenQualifier(tokenPoolReport.Output.DeployManagedTokenPoolOutput.TokenSymbol)) // token-scoped; the label keeps the display name, the key uses the symbol form
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save ManagedTokenPoolState address %s for Sui chain %d: %w", tokenPoolReport.Output.DeployManagedTokenPoolOutput.Objects.StateObjectId, config.SuiChainSelector, err)
			}

			// save Managed Pool OwnerId to the addressBook
			typeAndVersionManagedTokenPoolOwnerId := cldf.NewTypeAndVersion(deployment.SuiManagedTokenPoolOwnerCapObjectIDType, deployment.Version1_0_0)
			typeAndVersionManagedTokenPoolOwnerId.AddLabel(tokenPoolReport.Output.DeployManagedTokenPoolOutput.TokenSymbol)
			err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.SuiChainSelector, tokenPoolReport.Output.DeployManagedTokenPoolOutput.Objects.OwnerCapObjectId, typeAndVersionManagedTokenPoolOwnerId, deployment.TokenQualifier(tokenPoolReport.Output.DeployManagedTokenPoolOutput.TokenSymbol)) // token-scoped; the label keeps the display name, the key uses the symbol form
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save ManagedTokenPoolOwnerCapId address %s for Sui chain %d: %w", tokenPoolReport.Output.DeployManagedTokenPoolOutput.Objects.OwnerCapObjectId, config.SuiChainSelector, err)
			}
		}
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		DataStore:   ds,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeployTPAndConfigure) VerifyPreconditions(e cldf.Environment, config DeployTPAndConfigureConfig) error {
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return err
	}
	if state[config.SuiChainSelector].FastCurseMCMSPackageID == "" {
		return fmt.Errorf(
			"fast MCMS package not deployed for Sui chain %d; run DeploySuiChain first",
			config.SuiChainSelector,
		)
	}

	return deployment.ValidateNoDatastoreConflicts(e, config.SuiChainSelector, config.ReplaceExisting, func() ([]deployment.PlannedRef, error) {
		return plannedTokenPoolRefs(e, config)
	})
}

func plannedTokenPoolRefs(e cldf.Environment, config DeployTPAndConfigureConfig) ([]deployment.PlannedRef, error) {
	var planned []deployment.PlannedRef
	for _, tokenPoolType := range config.TokenPoolTypes {
		var (
			coinObjectTypeArg string
			types             []cldf.ContractType
		)
		switch tokenPoolType {
		case deployment.TokenPoolTypeBurnMint:
			coinObjectTypeArg = config.BurnMintTpInput.CoinObjectTypeArg
			types = []cldf.ContractType{
				deployment.SuiBnMTokenPoolType,
				deployment.SuiBnMTokenPoolStateType,
				deployment.SuiBnMTokenPoolOwnerIDType,
			}
		case deployment.TokenPoolTypeLockRelease:
			coinObjectTypeArg = config.LockReleaseTPInput.CoinObjectTypeArg
			types = []cldf.ContractType{
				deployment.SuiLnRTokenPoolType,
				deployment.SuiLnRTokenPoolStateType,
				deployment.SuiLnRTokenPoolOwnerIDType,
				deployment.SuiLnRTokenPoolRebalancerCapIDType,
			}
		case deployment.TokenPoolTypeManaged:
			coinObjectTypeArg = config.ManagedTPInput.CoinObjectTypeArg
			types = []cldf.ContractType{
				deployment.SuiManagedTokenPoolType,
				deployment.SuiManagedTokenPoolStateType,
				deployment.SuiManagedTokenPoolOwnerIDType,
			}
		default:
			return nil, fmt.Errorf("unknown token pool type %v", tokenPoolType)
		}

		symbol, err := coinSymbol(e, config.SuiChainSelector, coinObjectTypeArg)
		if err != nil {
			return nil, err
		}
		qualifier := deployment.TokenQualifier(symbol)
		for _, t := range types {
			planned = append(planned, deployment.PlannedRef{Type: t, Qualifier: qualifier})
		}
	}

	return planned, nil
}
