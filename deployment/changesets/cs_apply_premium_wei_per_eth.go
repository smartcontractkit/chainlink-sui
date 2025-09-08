package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

var _ cldf.ChangeSetV2[ApplyPremiumMultiplierWeiPerEthConfig] = ApplyPremiumMultiplierWeiPerEth{}

// ApplyPremiumMultiplierWeiPerEth applies premium multiplier wei per eth updates
type ApplyPremiumMultiplierWeiPerEth struct{}

// Apply implements deployment.ChangeSetV2.
func (d ApplyPremiumMultiplierWeiPerEth) Apply(e cldf.Environment, config ApplyPremiumMultiplierWeiPerEthConfig) (cldf.ChangesetOutput, error) {
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

	// Run applyPremiumMultiplierWeiPerEthUpdate Operation
	applyPremiumUpdate, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp, deps,
		ccipops.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput{
			CCIPPackageId:              config.CCIPPackageId,
			StateObjectId:              config.StateObjectId,
			OwnerCapObjectId:           config.OwnerCapObjectId,
			Tokens:                     config.Tokens,
			PremiumMultiplierWeiPerEth: config.PremiumMultiplierWeiPerEth,
		})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to applyPremiumMultiplierWeiPerEthUpdate for Sui chain %d: %w", config.ChainSelector, err)
	}

	seqReports = append(seqReports, applyPremiumUpdate.ToGenericReport())

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d ApplyPremiumMultiplierWeiPerEth) VerifyPreconditions(e cldf.Environment, config ApplyPremiumMultiplierWeiPerEthConfig) error {
	return nil
}
