package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	linkops "github.com/smartcontractkit/chainlink-sui/deployment/ops/link"
)

type MintLinkTokenConfig struct {
	ChainSelector  uint64
	TokenPackageId string
	TreasuryCapId  string
	Amount         uint64
}

var _ cldf.ChangeSetV2[MintLinkTokenConfig] = MintLinkToken{}

type MintLinkToken struct{}

// Apply implements deployment.ChangeSetV2.
func (d MintLinkToken) Apply(e cldf.Environment, config MintLinkTokenConfig) (cldf.ChangesetOutput, error) {

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
		SuiRPC: suiChain.URL,
	}

	// Run MintLinkToken Operation
	mintLinkTokenReport, err := operations.ExecuteOperation(e.OperationsBundle, linkops.MintLinkOp, deps,
		linkops.MintLinkTokenInput{
			LinkTokenPackageId: config.TokenPackageId,
			TreasuryCapId:      config.TreasuryCapId,
			Amount:             config.Amount, // 1099999999999999984
		})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to Mint LinkToken for Sui chain %d: %w", config.ChainSelector, err)
	}

	seqReports = append(seqReports, mintLinkTokenReport.ToGenericReport())

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d MintLinkToken) VerifyPreconditions(e cldf.Environment, config MintLinkTokenConfig) error {
	return nil
}
