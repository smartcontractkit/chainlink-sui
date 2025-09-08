package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

var _ cldf.ChangeSetV2[ApplyFeeTokenUpdateConfig] = ApplyFeeToken{}

// ApplyFeeToken applies fee token updates
type ApplyFeeToken struct{}

// Apply implements deployment.ChangeSetV2.
func (d ApplyFeeToken) Apply(e cldf.Environment, config ApplyFeeTokenUpdateConfig) (cldf.ChangesetOutput, error) {
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

	// Run applyFeeTokenUpdate Operation
	applyFeeTokenUpdate, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.FeeQuoterApplyFeeTokenUpdatesOp, deps,
		ccipops.FeeQuoterApplyFeeTokenUpdatesInput{
			CCIPPackageId:     config.CCIPPackageId,
			StateObjectId:     config.StateObjectId,
			OwnerCapObjectId:  config.OwnerCapObjectId,
			FeeTokensToRemove: config.FeeTokensToRemove,
			FeeTokensToAdd:    config.FeeTokensToAdd,
		})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to applyFeeTokenUpdate for Sui chain %d: %w", config.ChainSelector, err)
	}

	seqReports = append(seqReports, applyFeeTokenUpdate.ToGenericReport())

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d ApplyFeeToken) VerifyPreconditions(e cldf.Environment, config ApplyFeeTokenUpdateConfig) error {
	return nil
}
