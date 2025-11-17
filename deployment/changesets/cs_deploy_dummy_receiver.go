package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

type DeployDummyReceiverConfig struct {
	SuiChainSelector uint64
	McmsOwner        string
}

var _ cldf.ChangeSetV2[DeployDummyReceiverConfig] = DeployDummyReceiver{}

// DeployAptosChain deploys Aptos chain packages and modules
type DeployDummyReceiver struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeployDummyReceiver) Apply(e cldf.Environment, config DeployDummyReceiverConfig) (cldf.ChangesetOutput, error) {
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

	// Run DummyReceiver Operation
	DeployDummyReceiverOp, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.DeployCCIPDummyReceiverOp, deps, ccipops.DeployDummyReceiverInput{
		CCIPPackageId: state[config.SuiChainSelector].CCIPAddress,
		McmsPackageId: state[config.SuiChainSelector].MCMSPackageID,
		McmsOwner:     config.McmsOwner,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy dummy receiver for Sui chain %d: %w", config.SuiChainSelector, err)
	}

	// register receiver
	seqReports = append(seqReports, []operations.Report[any, any]{DeployDummyReceiverOp.ToGenericReport()}...)

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeployDummyReceiver) VerifyPreconditions(e cldf.Environment, config DeployDummyReceiverConfig) error {
	return nil
}
