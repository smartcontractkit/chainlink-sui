package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	managedtokenops "github.com/smartcontractkit/chainlink-sui/deployment/ops/managed_token"
)

type DeployManagedTokenConfig struct {
	managedtokenops.DeployAndInitManagedTokenInput
	ChainSelector uint64 `yaml:"chainSelector"`
}

var _ cldf.ChangeSetV2[DeployManagedTokenConfig] = DeployManagedToken{}

// DeployAptosChain deploys Sui chain packages and modules
type DeployManagedToken struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeployManagedToken) Apply(e cldf.Environment, config DeployManagedTokenConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]operations.Report[any, any], 0)

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
	}

	// Run DeployManagedToken Operation
	managedTokenReport, err := operations.ExecuteSequence(e.OperationsBundle, managedtokenops.DeployAndInitManagedTokenSequence, deps, config.DeployAndInitManagedTokenInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy ManagedToken for Sui chain %d: %w", config.ChainSelector, err)
	}

	// save ManagedToken address to the addressbook
	typeAndVersionManagedToken := cldf.NewTypeAndVersion(deployment.SuiManagedTokenType, deployment.Version1_0_0)
	err = ab.Save(config.ChainSelector, managedTokenReport.Output.ManagedTokenPackageId, typeAndVersionManagedToken)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save ManagedToken address %s for Sui chain %d: %w", managedTokenReport.Output.ManagedTokenPackageId, config.ChainSelector, err)
	}

	// save ManagedTokenOwnerCapObjectID address to the addressbook
	typeAndVersionOwnerCapObjectID := cldf.NewTypeAndVersion(deployment.SuiManagedTokenOwnerCapObjectID, deployment.Version1_0_0)
	err = ab.Save(config.ChainSelector, managedTokenReport.Output.Objects.OwnerCapObjectId, typeAndVersionOwnerCapObjectID)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save ManagedToken OwnerCapObjectId address %s for Sui chain %d: %w", managedTokenReport.Output.Objects.OwnerCapObjectId, config.ChainSelector, err)
	}

	// save ManagedTokenMinterCapID address to the addressbook
	typeAndVersionMinterCapID := cldf.NewTypeAndVersion(deployment.SuiManagedTokenMinterCapID, deployment.Version1_0_0)
	err = ab.Save(config.ChainSelector, managedTokenReport.Output.Objects.MinterCapObjectId, typeAndVersionMinterCapID)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save ManagedToken MinterCapObjectId address %s for Sui chain %d: %w", managedTokenReport.Output.Objects.MinterCapObjectId, config.ChainSelector, err)
	}

	// save ManagedTokenStateObjectID address to the addressbook
	typeAndVersionStateObjectID := cldf.NewTypeAndVersion(deployment.SuiManagedTokenStateObjectID, deployment.Version1_0_0)
	err = ab.Save(config.ChainSelector, managedTokenReport.Output.Objects.StateObjectId, typeAndVersionStateObjectID)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save ManagedToken StateObjectId address %s for Sui chain %d: %w", managedTokenReport.Output.Objects.StateObjectId, config.ChainSelector, err)
	}

	// save PublisherObjectId address to the addressbook
	typeAndVersionPublisherObjectId := cldf.NewTypeAndVersion(deployment.SuiManagedTokenPublisherObjectId, deployment.Version1_0_0)
	err = ab.Save(config.ChainSelector, managedTokenReport.Output.Objects.PublisherObjectId, typeAndVersionPublisherObjectId)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save ManagedToken PublisherObjectId address %s for Sui chain %d: %w", managedTokenReport.Output.Objects.PublisherObjectId, config.ChainSelector, err)
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeployManagedToken) VerifyPreconditions(e cldf.Environment, config DeployManagedTokenConfig) error {
	return nil
}
