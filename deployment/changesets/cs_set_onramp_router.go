package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccip_router_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_router"
)

type ConfigureRouterOnRampConfig struct {
	SuiChainSelector  uint64
	DestChainSelector []uint64
	OnRampPackageId   string
	McmsOwner         string
}

var _ cldf.ChangeSetV2[ConfigureRouterOnRampConfig] = ConfigureRouterOnRamp{}

type ConfigureRouterOnRamp struct{}

// Apply implements deployment.ChangeSetV2.
func (d ConfigureRouterOnRamp) Apply(e cldf.Environment, config ConfigureRouterOnRampConfig) (cldf.ChangesetOutput, error) {
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

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

	onrampAddresses := make([]string, len(config.DestChainSelector))
	for i := range config.DestChainSelector {
		onrampAddresses[i] = config.OnRampPackageId
	}

	reportConfigureRouterOp, err := operations.ExecuteOperation(e.OperationsBundle, ccip_router_ops.SetOnRampsOp, deps, ccip_router_ops.SetOnRampsInput{
		RouterPackageId:     state[config.SuiChainSelector].CCIPRouterAddress,
		RouterStateObjectId: state[config.SuiChainSelector].CCIPRouterStateObjectID,
		OwnerCapObjectId:    state[config.SuiChainSelector].CCIPRouterOwnerCapObjectId,
		DestChainSelectors:  config.DestChainSelector,
		OnRampAddresses:     onrampAddresses,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to configure router onRamp for Sui chain %d: %w", config.SuiChainSelector, err)
	}

	seqReports = append(seqReports, []operations.Report[any, any]{reportConfigureRouterOp.ToGenericReport()}...)

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d ConfigureRouterOnRamp) VerifyPreconditions(e cldf.Environment, config ConfigureRouterOnRampConfig) error {
	return nil
}
