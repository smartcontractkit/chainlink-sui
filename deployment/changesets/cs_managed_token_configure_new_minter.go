package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	coin_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/coin"
	managedtokenops "github.com/smartcontractkit/chainlink-sui/deployment/ops/managed_token"
)

type ManagedTokenConfigureNewMinterConfig struct {
	SuiChainSelector      uint64
	StateObjectId         string
	OwnerCapObjectId      string
	MinterAddress         string
	CoinObjectTypeArg     string
	ManagedTokenPackageId string
	Allowance             uint64
	IsUnlimited           bool
	Source                string
}

var _ cldf.ChangeSetV2[ManagedTokenConfigureNewMinterConfig] = ManagedTokenConfigureNewMinter{}

type ManagedTokenConfigureNewMinter struct{}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d ManagedTokenConfigureNewMinter) VerifyPreconditions(e cldf.Environment, config ManagedTokenConfigureNewMinterConfig) error {
	return nil
}

// Apply implements deployment.ChangeSetV2.
func (d ManagedTokenConfigureNewMinter) Apply(e cldf.Environment, config ManagedTokenConfigureNewMinterConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]operations.Report[any, any], 0)

	suiChain := e.BlockChains.SuiChains()[config.SuiChainSelector]

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

	configureNewMinterReport, err := operations.ExecuteOperation(e.OperationsBundle, managedtokenops.ManagedTokenConfigureNewMinterOp, deps, managedtokenops.ManagedTokenConfigureNewMinterInput{
		ManagedTokenPackageId: config.ManagedTokenPackageId,
		CoinObjectTypeArg:     config.CoinObjectTypeArg,
		StateObjectId:         config.StateObjectId,
		OwnerCapObjectId:      config.OwnerCapObjectId,
		MinterAddress:         config.MinterAddress,
		Allowance:             config.Allowance,
		IsUnlimited:           config.IsUnlimited,
		Source:                config.Source,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to configure new minter for managed token: %w", err)
	}

	symbolReport, err := cld_ops.ExecuteOperation(e.OperationsBundle, coin_ops.GetCoinSymbolOp, deps, config.CoinObjectTypeArg)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to get coin symbol: %w", err)
	}

	// save ManagedTokenMinterCapID address to the addressbook
	typeAndVersionMinterCapID := cldf.NewTypeAndVersion(deployment.SuiManagedTokenMinterCapID, deployment.Version1_0_0)
	typeAndVersionMinterCapID.AddLabel(symbolReport.Output.Symbol)
	err = ab.Save(config.SuiChainSelector, configureNewMinterReport.Output.Objects.MinterCapObjectId, typeAndVersionMinterCapID)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save ManagedToken MinterCapObjectId address %s for Sui chain %d: %w", configureNewMinterReport.Output.Objects.MinterCapObjectId, config.SuiChainSelector, err)
	}

	seqReports = append(seqReports, []operations.Report[any, any]{configureNewMinterReport.ToGenericReport()}...)

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}
