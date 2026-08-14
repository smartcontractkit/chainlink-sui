package changesets

import (
	"fmt"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccip_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

// SetTokenTransferFeesConfig sets per-token, per-destination-chain transfer fee config on the
// Sui CCIP FeeQuoter. Token object IDs are coin-metadata object IDs, matching the fee-token
// identifier convention. Each Add* slice must be the same length as AddTokens.
type SetTokenTransferFeesConfig struct {
	SuiChainSelector     uint64   `yaml:"suiChainSelector"`
	DestChainSelector    uint64   `yaml:"destChainSelector"`
	AddTokens            []string `yaml:"addTokens"`
	AddMinFeeUsdCents    []uint32 `yaml:"addMinFeeUsdCents"`
	AddMaxFeeUsdCents    []uint32 `yaml:"addMaxFeeUsdCents"`
	AddDeciBps           []uint16 `yaml:"addDeciBps"`
	AddDestGasOverhead   []uint32 `yaml:"addDestGasOverhead"`
	AddDestBytesOverhead []uint32 `yaml:"addDestBytesOverhead"`
	AddIsEnabled         []bool   `yaml:"addIsEnabled"`
	RemoveTokens         []string `yaml:"removeTokens"`
}

// SetTokenTransferFees applies token transfer fee config updates to the Sui CCIP FeeQuoter.
type SetTokenTransferFees struct{}

var _ cldf.ChangeSetV2[SetTokenTransferFeesConfig] = SetTokenTransferFees{}

// Apply implements deployment.ChangeSetV2.
func (d SetTokenTransferFees) Apply(e cldf.Environment, config SetTokenTransferFeesConfig) (cldf.ChangesetOutput, error) {
	ab := cldf.NewMemoryAddressBook()
	ds := fdatastore.NewMemoryDataStore()
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

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

	report, err := operations.ExecuteOperation(e.OperationsBundle, ccip_ops.FeeQuoterApplyTokenTransferFeeConfigUpdatesOp, deps, ccip_ops.FeeQuoterApplyTokenTransferFeeConfigUpdatesInput{
		CCIPPackageId:        state[suiChain.Selector].CCIPAddress,
		StateObjectId:        state[suiChain.Selector].CCIPObjectRef,
		OwnerCapObjectId:     state[suiChain.Selector].CCIPOwnerCapObjectId,
		DestChainSelector:    config.DestChainSelector,
		AddTokens:            config.AddTokens,
		AddMinFeeUsdCents:    config.AddMinFeeUsdCents,
		AddMaxFeeUsdCents:    config.AddMaxFeeUsdCents,
		AddDeciBps:           config.AddDeciBps,
		AddDestGasOverhead:   config.AddDestGasOverhead,
		AddDestBytesOverhead: config.AddDestBytesOverhead,
		AddIsEnabled:         config.AddIsEnabled,
		RemoveTokens:         config.RemoveTokens,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to apply token transfer fee config updates for Sui chain %d dest %d: %w", config.SuiChainSelector, config.DestChainSelector, err)
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		DataStore:   ds,
		Reports:     []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d SetTokenTransferFees) VerifyPreconditions(e cldf.Environment, config SetTokenTransferFeesConfig) error {
	return nil
}
