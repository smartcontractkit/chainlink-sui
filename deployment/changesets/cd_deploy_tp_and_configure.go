package changesets

import (
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
			config.BurnMintTpInput.MCMSAddress = state[config.SuiChainSelector].MCMsAddress
			config.BurnMintTpInput.MCMSOwnerAddress = deployerAddr
			config.BurnMintTpInput.CCIPObjectRefObjectId = state[config.SuiChainSelector].CCIPObjectRef
			config.BurnMintTpInput.TokenPoolAdministrator = deployerAddr // check with felix if this is fine

			_, err = operations.ExecuteSequence(e.OperationsBundle, burnminttokenpoolops.DeployAndInitBurnMintTokenPoolSequence, deps, config.BurnMintTpInput)
			if err != nil {
				return cldf.ChangesetOutput{}, err
			}
		}

		if tokenPoolType == "lnr" {
			config.LockReleaseTPInput.CCIPPackageId = state[config.SuiChainSelector].CCIPAddress
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
		Reports: seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeployTPAndConfigure) VerifyPreconditions(e cldf.Environment, config DeployTPAndConfigureConfig) error {
	return nil
}
