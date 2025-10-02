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
	// the below can be part of (DeployAndInitBurnMintTokenPoolSequence)
	// Initialize TP
	// ApplyChainUpdates
	// SetChainRateLimiterConfigs
	// Add remote TP

	for _, tokenPoolType := range config.TokenPoolTypes {
		if tokenPoolType == "bnm" {
			config.BurnMintTpInput.CCIPPackageId = state[config.SuiChainSelector].CCIPAddress
			config.BurnMintTpInput.CCIPTokenPoolPackageId = state[config.SuiChainSelector].TokenPoolAddress
			config.BurnMintTpInput.MCMSAddress = state[config.SuiChainSelector].MCMsAddress
			config.BurnMintTpInput.MCMSOwnerAddress = deployerAddr
			config.BurnMintTpInput.CCIPObjectRefObjectId = state[config.SuiChainSelector].CCIPObjectRef
			config.BurnMintTpInput.TokenPoolAdministrator = deployerAddr // check with felix if this is fine

			BnMTokenPoolSeqReport, err := operations.ExecuteSequence(e.OperationsBundle, burnminttokenpoolops.DeployAndInitBurnMintTokenPoolSequence, deps, config.BurnMintTpInput)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}

			// save BnM Pool to the addressbook
			typeAndVersionBurnMintTokenPool := cldf.NewTypeAndVersion(deployment.SuiBnMTokenPoolType, deployment.Version1_0_0)
			err = ab.Save(config.SuiChainSelector, BnMTokenPoolSeqReport.Output.BurnMintTPPackageID, typeAndVersionBurnMintTokenPool)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save BnMTokenPool address %s for Sui chain %d: %w", BnMTokenPoolSeqReport.Output.BurnMintTPPackageID, config.SuiChainSelector, err)
			}

			// save BnM Pool State to the addressBook
			typeAndVersionBurnMintTokenPoolState := cldf.NewTypeAndVersion(deployment.SuiBnMTokenPoolStateType, deployment.Version1_0_0)
			err = ab.Save(config.SuiChainSelector, BnMTokenPoolSeqReport.Output.Objects.StateObjectId, typeAndVersionBurnMintTokenPoolState)
			if err != nil {
				return cldf.ChangesetOutput{}, fmt.Errorf("failed to save BnMTokenPoolState address %s for Sui chain %d: %w", BnMTokenPoolSeqReport.Output.Objects.StateObjectId, config.SuiChainSelector, err)
			}
		}

		if tokenPoolType == "lnr" {
			config.LockReleaseTPInput.CCIPPackageId = state[config.SuiChainSelector].CCIPAddress
			config.LockReleaseTPInput.CCIPTokenPoolPackageId = state[config.SuiChainSelector].TokenPoolAddress
			config.LockReleaseTPInput.MCMSAddress = state[config.SuiChainSelector].MCMsAddress
			config.LockReleaseTPInput.MCMSOwnerAddress = deployerAddr
			config.LockReleaseTPInput.CCIPObjectRefObjectId = state[config.SuiChainSelector].CCIPObjectRef
			config.LockReleaseTPInput.TokenPoolAdministrator = deployerAddr

			_, err = operations.ExecuteSequence(e.OperationsBundle, lockreleasetokenpoolops.DeployAndInitLockReleaseTokenPoolSequence, deps, config.LockReleaseTPInput)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
		}

		if tokenPoolType == "managed" {
			config.ManagedTPInput.CCIPPackageId = state[config.SuiChainSelector].CCIPAddress
			config.ManagedTPInput.CCIPTokenPoolPackageId = state[config.SuiChainSelector].TokenPoolAddress
			config.ManagedTPInput.MCMSAddress = state[config.SuiChainSelector].MCMsAddress
			config.ManagedTPInput.MCMSOwnerAddress = deployerAddr
			config.ManagedTPInput.CCIPObjectRefObjectId = state[config.SuiChainSelector].CCIPObjectRef
			config.ManagedTPInput.TokenPoolAdministrator = deployerAddr

			_, err = operations.ExecuteSequence(e.OperationsBundle, managedtokenpoolops.DeployAndInitManagedTokenPoolSequence, deps, config.ManagedTPInput)
			if err != nil {
				return cldf.ChangesetOutput{}, err
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
