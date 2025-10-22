package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	lockreleasetokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_lock_release_token_pool"
	managedtokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_managed_token_pool"
	tokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_token_pool"
)

type DeployTPAndConfigureConfig struct {
	SuiChainSelector   uint64
	TokenPoolTypes     []string
	ManagedTPInput     managedtokenpoolops.SeqDeployAndInitManagedTokenPoolInput
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
	}

	// Populate state information for each token pool type
	for _, tokenPoolType := range config.TokenPoolTypes {
		switch tokenPoolType {
		case "bnm":
			config.BurnMintTpInput.CCIPPackageId = state[config.SuiChainSelector].CCIPAddress
			config.BurnMintTpInput.MCMSAddress = state[config.SuiChainSelector].MCMSPackageID
			// TODO: MCMSOwner address should come state
			config.BurnMintTpInput.MCMSOwnerAddress = deployerAddr
			config.BurnMintTpInput.CCIPObjectRefObjectId = state[config.SuiChainSelector].CCIPObjectRef
			config.BurnMintTpInput.TokenPoolAdministrator = deployerAddr
		case "lnr":
			config.LockReleaseTPInput.CCIPPackageId = state[config.SuiChainSelector].CCIPAddress
			config.LockReleaseTPInput.MCMSAddress = state[config.SuiChainSelector].MCMSPackageID
			config.LockReleaseTPInput.MCMSOwnerAddress = deployerAddr
			config.LockReleaseTPInput.CCIPObjectRefObjectId = state[config.SuiChainSelector].CCIPObjectRef
			config.LockReleaseTPInput.TokenPoolAdministrator = deployerAddr
		case "managed":
			config.ManagedTPInput.CCIPPackageId = state[config.SuiChainSelector].CCIPAddress
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
		case "bnm":
			// save BnM Pool to the addressbook
			typeAndVersionBurnMintTokenPool := cldf.NewTypeAndVersion(deployment.SuiBnMTokenPoolType, deployment.Version1_0_0)
			err = ab.Save(config.SuiChainSelector, tokenPoolReport.Output.BurnMintTPPackageID, typeAndVersionBurnMintTokenPool)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save BnMTokenPool address %s for Sui chain %d: %w", tokenPoolReport.Output.BurnMintTPPackageID, config.SuiChainSelector, err)
			}

			// save BnM Pool State to the addressBook
			typeAndVersionBurnMintTokenPoolState := cldf.NewTypeAndVersion(deployment.SuiBnMTokenPoolStateType, deployment.Version1_0_0)
			err = ab.Save(config.SuiChainSelector, tokenPoolReport.Output.DeployBurnMintTokenPoolOutput.Objects.StateObjectId, typeAndVersionBurnMintTokenPoolState)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save BnMTokenPoolState address %s for Sui chain %d: %w", tokenPoolReport.Output.DeployBurnMintTokenPoolOutput.Objects.StateObjectId, config.SuiChainSelector, err)
			}

			// save BnM Pool OwnerId to the addressBook
			typeAndVersionBurnMintTokenPoolOwnerId := cldf.NewTypeAndVersion(deployment.SuiBnMTokenPoolOwnerIDType, deployment.Version1_0_0)
			err = ab.Save(config.SuiChainSelector, tokenPoolReport.Output.DeployBurnMintTokenPoolOutput.Objects.OwnerCapObjectId, typeAndVersionBurnMintTokenPoolOwnerId)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save BnMTokenPoolOwnerCapId address %s for Sui chain %d: %w", tokenPoolReport.Output.DeployBurnMintTokenPoolOutput.Objects.OwnerCapObjectId, config.SuiChainSelector, err)
			}
		}
		// Note: Address book saving for "lnr" and "managed" token pools can be added here if needed
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
