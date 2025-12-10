package changesets

import (
	"fmt"
	"math/big"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

type NewFeeTokenConfig struct {
	SuiChainSelector  uint64
	FeeTokensToRemove []string
	FeeTokensToAdd    []string // should be the objectID

	// update price
	SourceUsdPerToken []*big.Int

	// premium multiplier wei per eth
	PremiumMultiplierWeiPerEth []uint64
}

// ConnectSuiToEVM connects sui chain with EVM
type NewFeeToken struct{}

var _ cldf.ChangeSetV2[NewFeeTokenConfig] = NewFeeToken{}

// Apply implements deployment.ChangeSetV2.
func (d NewFeeToken) Apply(e cldf.Environment, config NewFeeTokenConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	seqReports := make([]operations.Report[any, any], 0)

	suiChains := e.BlockChains.SuiChains()
	suiChain := suiChains[config.SuiChainSelector]

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

	// Run ApplyFeeTokenUpdate Operation
	applyFeeTokenUpdateOP, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.FeeQuoterApplyFeeTokenUpdatesOp, deps, ccipops.FeeQuoterApplyFeeTokenUpdatesInput{
		CCIPPackageId:     state[suiChain.Selector].CCIPAddress,
		StateObjectId:     state[suiChain.Selector].CCIPObjectRef,
		OwnerCapObjectId:  state[suiChain.Selector].CCIPOwnerCapObjectId,
		FeeTokensToRemove: config.FeeTokensToRemove,
		FeeTokensToAdd:    config.FeeTokensToAdd,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to register receiver for Sui chain %d: %w", config.SuiChainSelector, err)
	}

	seqReports = append(seqReports, []operations.Report[any, any]{applyFeeTokenUpdateOP.ToGenericReport()}...)

	// Run ApplyPremiumMultiplier
	applyPremiumMultiplierOP, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp, deps, ccipops.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput{
		CCIPPackageId:              state[suiChain.Selector].CCIPAddress,
		StateObjectId:              state[suiChain.Selector].CCIPObjectRef,
		OwnerCapObjectId:           state[suiChain.Selector].CCIPOwnerCapObjectId,
		Tokens:                     config.FeeTokensToAdd,
		PremiumMultiplierWeiPerEth: config.PremiumMultiplierWeiPerEth,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to register receiver for Sui chain %d: %w", config.SuiChainSelector, err)
	}

	seqReports = append(seqReports, []operations.Report[any, any]{applyPremiumMultiplierOP.ToGenericReport()}...)

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d NewFeeToken) VerifyPreconditions(e cldf.Environment, config NewFeeTokenConfig) error {
	return nil
}
