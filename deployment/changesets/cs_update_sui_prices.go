package changesets

import (
	"fmt"
	"math/big"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

type UpdateSuiPriceConfig struct {
	ChainSelector       uint64
	DestChainSelector   []uint64
	CCIPPackageId       string
	CCIPObjectRef       string
	SourceTokenMetadata []string
	SourceUsdPerToken   []*big.Int
	GasUsdPerUnitGas    []*big.Int
}

var _ cldf.ChangeSetV2[UpdateSuiPriceConfig] = UpdateSuiFeeQuoterPrice{}

// DeployAptosChain deploys Aptos chain packages and modules
type UpdateSuiFeeQuoterPrice struct{}

// Apply implements deployment.ChangeSetV2.
func (d UpdateSuiFeeQuoterPrice) Apply(e cldf.Environment, config UpdateSuiPriceConfig) (cldf.ChangesetOutput, error) {

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

	// Run UpdateTokenPrices Operation
	updateTokenPrices, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.FeeQuoterUpdateTokenPricesOp, deps,
		ccipops.FeeQuoterUpdateTokenPricesInput{
			CCIPPackageId: config.CCIPPackageId,
			CCIPObjectRef: config.CCIPObjectRef,
			// FeeQuoterCapId:        feeQuoterCapId,
			SourceTokens:          config.SourceTokenMetadata,
			SourceUsdPerToken:     config.SourceUsdPerToken,
			GasDestChainSelectors: config.DestChainSelector,
			GasUsdPerUnitGas:      config.GasUsdPerUnitGas,
		})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to updatePrice for Sui chain %d: %w", config.ChainSelector, err)
	}

	seqReports = append(seqReports, updateTokenPrices.ToGenericReport())

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d UpdateSuiFeeQuoterPrice) VerifyPreconditions(e cldf.Environment, config UpdateSuiPriceConfig) error {
	return nil
}
