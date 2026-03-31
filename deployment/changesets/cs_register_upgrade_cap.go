package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
)

var _ cldf.ChangeSetV2[RegisterUpgradeCapConfig] = RegisterUpgradeCap{}

// PackageUpgradeCap contains the upgrade cap information for a package
type PackageUpgradeCap struct {
	PackageName  string
	PackageID    string
	UpgradeCapID string
}

// RegisterUpgradeCapConfig contains the configuration for registering upgrade capabilities
type RegisterUpgradeCapConfig struct {
	ChainSelector      uint64              `json:"chainSelector"`
	PackageUpgradeCaps []PackageUpgradeCap `json:"packageUpgradeCaps" validate:"required,min=1"`
	// Optional MCMS details - if empty, they will be auto-populated from address book
	MCMSPackageID         string `json:"mcmsPackageID,omitempty"`
	RegistryObjectID      string `json:"registryObjectID,omitempty"`
	DeployerStateObjectID string `json:"deployerStateObjectID,omitempty"`
}

type RegisterUpgradeCap struct{}

// Apply implements ChangeSetV2 - registers upgrade capabilities for the specified packages
func (r RegisterUpgradeCap) Apply(e cldf.Environment, config RegisterUpgradeCapConfig) (cldf.ChangesetOutput, error) {
	suiState, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	chainState, ok := suiState[config.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("no Sui chain state found for chain selector %d", config.ChainSelector)
	}

	// Get MCMS state from address book (using normal MCMS, not FastCurse)
	mcmsState := chainState.MCMSState(false)

	// Auto-populate MCMS details if not provided
	if config.MCMSPackageID == "" {
		config.MCMSPackageID = mcmsState.PackageID
	}
	if config.RegistryObjectID == "" {
		config.RegistryObjectID = mcmsState.RegistryObjectID
	}
	if config.DeployerStateObjectID == "" {
		config.DeployerStateObjectID = mcmsState.DeployerStateObjectID
	}

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
	seqReports := make([]cld_ops.Report[any, any], 0)

	// Register upgrade caps for each package
	for _, upgradeCap := range config.PackageUpgradeCaps {
		// Create the registration input
		regInput := mcmsops.RegisterUpgradeCapInput{
			MCMSPackageID:         config.MCMSPackageID,
			UpgradeCapObjectID:    upgradeCap.UpgradeCapID,
			RegistryObjectID:      config.RegistryObjectID,
			DeployerStateObjectID: config.DeployerStateObjectID,
			PackageName:           upgradeCap.PackageName,
		}

		// Execute the register upgrade cap operation
		result, err := cld_ops.ExecuteOperation(e.OperationsBundle, mcmsops.RegisterUpgradeCapOp, deps, regInput)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to register upgrade cap for package %s: %w", upgradeCap.PackageName, err)
		}

		seqReports = append(seqReports, cld_ops.Report[any, any]{
			Input:  regInput,
			Output: result.Output,
		})

		e.Logger.Infow("Successfully registered upgrade cap",
			"package", upgradeCap.PackageName,
			"packageID", upgradeCap.PackageID,
			"upgradeCapID", upgradeCap.UpgradeCapID,
		)
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements ChangeSetV2
func (r RegisterUpgradeCap) VerifyPreconditions(e cldf.Environment, config RegisterUpgradeCapConfig) error {
	// Verify that MCMS state exists
	suiState, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	chainState, ok := suiState[config.ChainSelector]
	if !ok {
		return fmt.Errorf("no Sui chain state found for chain selector %d", config.ChainSelector)
	}

	mcmsState := chainState.MCMSState(false)
	if mcmsState.PackageID == "" {
		return fmt.Errorf("MCMS package not found on chain %d", config.ChainSelector)
	}

	// Verify that all packages have upgrade caps provided
	if len(config.PackageUpgradeCaps) == 0 {
		return fmt.Errorf("at least one package upgrade cap must be provided")
	}

	for _, upgradeCap := range config.PackageUpgradeCaps {
		if upgradeCap.PackageName == "" || upgradeCap.PackageID == "" || upgradeCap.UpgradeCapID == "" {
			return fmt.Errorf("all package upgrade cap fields must be non-empty")
		}
	}

	return nil
}
