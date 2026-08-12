package changesets

import (
	"fmt"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
)

var _ cldf.ChangeSetV2[RegisterUpgradeCapConfig] = RegisterUpgradeCap{}

// RegisterUpgradeCapConfig contains the configuration for registering a package upgrade capability
type RegisterUpgradeCapConfig struct {
	ChainSelector uint64 `json:"chainSelector" validate:"required"`
	PackageName   string `json:"packageName" validate:"required"`
	PackageID     string `json:"packageID" validate:"required"`
	UpgradeCapID  string `json:"upgradeCapID" validate:"required"`
	// MCMS details
	MCMSPackageID         string `json:"mcmsPackageID" validate:"required"`
	RegistryObjectID      string `json:"registryObjectID" validate:"required"`
	DeployerStateObjectID string `json:"deployerStateObjectID" validate:"required"`
}

type RegisterUpgradeCap struct{}

// Apply implements ChangeSetV2 - registers upgrade capability for a package
func (r RegisterUpgradeCap) Apply(e cldf.Environment, config RegisterUpgradeCapConfig) (cldf.ChangesetOutput, error) {
	suiChains := e.BlockChains.SuiChains()
	suiChain, ok := suiChains[config.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("no Sui chain found for chain selector %d", config.ChainSelector)
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

	ab := cldf.NewMemoryAddressBook()
	ds := fdatastore.NewMemoryDataStore()
	seqReports := make([]cld_ops.Report[any, any], 0)

	// Create the registration input
	regInput := mcmsops.RegisterUpgradeCapInput{
		MCMSPackageID:         config.MCMSPackageID,
		UpgradeCapObjectID:    config.UpgradeCapID,
		RegistryObjectID:      config.RegistryObjectID,
		DeployerStateObjectID: config.DeployerStateObjectID,
		PackageName:           config.PackageName,
	}

	// Execute the register upgrade cap operation
	result, err := cld_ops.ExecuteOperation(e.OperationsBundle, mcmsops.RegisterUpgradeCapOp, deps, regInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to register upgrade cap for package %s: %w", config.PackageName, err)
	}

	seqReports = append(seqReports, cld_ops.Report[any, any]{
		Input:  regInput,
		Output: result.Output,
	})

	e.Logger.Infow("Successfully registered upgrade cap",
		"package", config.PackageName,
		"packageID", config.PackageID,
		"upgradeCapID", config.UpgradeCapID,
	)

	return cldf.ChangesetOutput{
		AddressBook: ab,
		DataStore:   ds,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements ChangeSetV2
func (r RegisterUpgradeCap) VerifyPreconditions(e cldf.Environment, config RegisterUpgradeCapConfig) error {
	// Verify that all required fields are provided
	if config.PackageName == "" || config.PackageID == "" || config.UpgradeCapID == "" {
		return fmt.Errorf("all package upgrade cap fields must be non-empty")
	}
	if config.MCMSPackageID == "" || config.RegistryObjectID == "" || config.DeployerStateObjectID == "" {
		return fmt.Errorf("all MCMS details must be provided")
	}

	return nil
}
