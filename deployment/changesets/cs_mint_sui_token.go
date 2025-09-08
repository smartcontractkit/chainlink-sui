package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	linkops "github.com/smartcontractkit/chainlink-sui/deployment/ops/link"
)

var _ cldf.ChangeSetV2[MintSuiTokenConfig] = MintSuiToken{}

// MintSuiToken mints Sui tokens
type MintSuiToken struct{}

// Apply implements deployment.ChangeSetV2.
func (d MintSuiToken) Apply(e cldf.Environment, config MintSuiTokenConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	seqReports := make([]operations.Report[any, any], 0)

	suiChains := e.BlockChains.SuiChains()
	suiChain := suiChains[config.ChainSelector]
	suiSigner := suiChain.Signer

	deps := sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: suiSigner,
		GetCallOpts: func() *bind.CallOpts {
			b := uint64(400_000_000)
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
	}

	// Run MintSuiToken Operation
	mintLinkTokenReport, err := operations.ExecuteOperation(e.OperationsBundle, linkops.MintLinkOp, deps,
		linkops.MintLinkTokenInput{
			LinkTokenPackageId: config.TokenPackageId,
			TreasuryCapId:      config.TreasuryCapId,
			Amount:             config.Amount,
		})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to Mint SuiToken for Sui chain %d: %w", config.ChainSelector, err)
	}

	seqReports = append(seqReports, mintLinkTokenReport.ToGenericReport())

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d MintSuiToken) VerifyPreconditions(e cldf.Environment, config MintSuiTokenConfig) error {
	return nil
}
