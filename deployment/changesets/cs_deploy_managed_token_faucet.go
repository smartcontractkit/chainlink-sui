package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	managedtokenops "github.com/smartcontractkit/chainlink-sui/deployment/ops/managed_token"
	managedtokenfaucetops "github.com/smartcontractkit/chainlink-sui/deployment/ops/managed_token_faucet"
)

type DeployManagedTokenFaucetConfig struct {
	ChainSelector   uint64 `yaml:"chainSelector"`
	TokenSymbol     string `yaml:"tokenSymbol"`
	CoinType        string `yaml:"coinType"`
	MintCapObjectId string `yaml:"mintCapObjectId"`
}

var _ cldf.ChangeSetV2[DeployManagedTokenFaucetConfig] = DeployManagedTokenFaucet{}

type DeployManagedTokenFaucet struct{}

func (d DeployManagedTokenFaucet) Apply(e cldf.Environment, config DeployManagedTokenFaucetConfig) (cldf.ChangesetOutput, error) {
	if config.TokenSymbol == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf("tokenSymbol must be provided")
	}
	if config.CoinType == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf("coinType must be provided")
	}

	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]cld_ops.Report[any, any], 0)

	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	suiChains := e.BlockChains.SuiChains()
	suiChain, ok := suiChains[config.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("sui chain %d not found", config.ChainSelector)
	}

	chainState := state[config.ChainSelector]
	managedToken, ok := chainState.ManagedTokens[config.TokenSymbol]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("managed token for symbol %s not found on chain %d", config.TokenSymbol, config.ChainSelector)
	}
	if managedToken.PackageID == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf("managed token package id not found for symbol %s on chain %d", config.TokenSymbol, config.ChainSelector)
	}

	if chainState.MCMSPackageID == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf("mcms package id not found for chain %d", config.ChainSelector)
	}

	deployerAddr, err := suiChain.Signer.GetAddress()
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get deployer address: %w", err)
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

	deployReport, err := cld_ops.ExecuteOperation(e.OperationsBundle, managedtokenfaucetops.DeployManagedTokenFaucetOp, deps, managedtokenfaucetops.DeployManagedTokenFaucetInput{
		ManagedTokenPackageId: managedToken.PackageID,
		MCMSAddress:           chainState.MCMSPackageID,
		MCMSOwnerAddress:      deployerAddr,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy managed token faucet for chain %d: %w", config.ChainSelector, err)
	}

	mintCapObjectId := config.MintCapObjectId
	if mintCapObjectId == "" {
		// issue a minter cap for the deployer. this allows the deployer to initialize the faucet with this minter cap later.
		minterReport, err := cld_ops.ExecuteOperation(e.OperationsBundle, managedtokenops.ManagedTokenConfigureNewMinterOp, deps, managedtokenops.ManagedTokenConfigureNewMinterInput{
			ManagedTokenPackageId: managedToken.PackageID,
			CoinObjectTypeArg:     config.CoinType,
			StateObjectId:         managedToken.StateObjectId,
			OwnerCapObjectId:      managedToken.OwnerCapObjectId,
			MinterAddress:         deployerAddr,
			Allowance:             0,    // Unlimited allowance
			IsUnlimited:           true, // Set as unlimited minter
			Source:                "deploy_managed_token_faucet_" + config.TokenSymbol,
		})
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to configure deployer as minter for chain %d: %w", config.ChainSelector, err)
		}

		// Save the minter cap for the deployer
		typeAndVersionMinterCapID := cldf.NewTypeAndVersion(deployment.SuiManagedTokenMinterCapID, deployment.Version1_0_0)
		typeAndVersionMinterCapID.AddLabel(config.TokenSymbol)
		if err := ab.Save(config.ChainSelector, minterReport.Output.Objects.MinterCapObjectId, typeAndVersionMinterCapID); err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to save managed token minter cap id %s: %w", minterReport.Output.Objects.MinterCapObjectId, err)
		}
		mintCapObjectId = minterReport.Output.Objects.MinterCapObjectId
	}

	initReport, err := cld_ops.ExecuteOperation(e.OperationsBundle, managedtokenfaucetops.InitializeManagedTokenFaucetOp, deps, managedtokenfaucetops.InitializeManagedTokenFaucetInput{
		ManagedTokenFaucetPackageId: deployReport.Output.PackageId,
		CoinObjectTypeArg:           config.CoinType,
		MintCapObjectId:             mintCapObjectId,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to initialize managed token faucet for chain %d: %w", config.ChainSelector, err)
	}

	typeAndVersionPackageID := cldf.NewTypeAndVersion(deployment.SuiManagedTokenFaucetPackageIDType, deployment.Version1_0_0)
	typeAndVersionPackageID.AddLabel(config.TokenSymbol)
	if err := ab.Save(config.ChainSelector, deployReport.Output.PackageId, typeAndVersionPackageID); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save managed token faucet package id %s: %w", deployReport.Output.PackageId, err)
	}

	typeAndVersionUpgradeCapID := cldf.NewTypeAndVersion(deployment.SuiManagedTokenFaucetUpgradeCapObjectIDType, deployment.Version1_0_0)
	typeAndVersionUpgradeCapID.AddLabel(config.TokenSymbol)
	if err := ab.Save(config.ChainSelector, deployReport.Output.Objects.UpgradeCapObjectId, typeAndVersionUpgradeCapID); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save managed token faucet upgrade cap id %s: %w", deployReport.Output.Objects.UpgradeCapObjectId, err)
	}

	typeAndVersionStateID := cldf.NewTypeAndVersion(deployment.SuiManagedTokenFaucetStateObjectIDType, deployment.Version1_0_0)
	typeAndVersionStateID.AddLabel(config.TokenSymbol)
	if err := ab.Save(config.ChainSelector, initReport.Output.Objects.FaucetStateObjectId, typeAndVersionStateID); err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save managed token faucet state object id %s: %w", initReport.Output.Objects.FaucetStateObjectId, err)
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

func (d DeployManagedTokenFaucet) VerifyPreconditions(e cldf.Environment, config DeployManagedTokenFaucetConfig) error {
	return nil
}
