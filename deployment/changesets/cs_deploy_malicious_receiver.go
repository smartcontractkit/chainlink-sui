package changesets

import (
	"fmt"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

type DeployMaliciousReceiverConfig struct {
	SuiChainSelector uint64
	McmsOwner        string
}

var _ cldf.ChangeSetV2[DeployMaliciousReceiverConfig] = DeployMaliciousReceiver{}

// DeployMaliciousReceiver deploys the TEST-only malicious receiver used by the
// transmitter-ownership guard E2E smoke test.
type DeployMaliciousReceiver struct{}

// Apply implements deployment.ChangeSetV2.
func (d DeployMaliciousReceiver) Apply(e cldf.Environment, config DeployMaliciousReceiverConfig) (cldf.ChangesetOutput, error) {
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	ab := cldf.NewMemoryAddressBook()
	ds := fdatastore.NewMemoryDataStore()
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

	deployMaliciousReceiverOp, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.DeployCCIPMaliciousReceiverOp, deps, ccipops.DeployMaliciousReceiverInput{
		CCIPPackageId:     state[config.SuiChainSelector].CCIPAddress,
		McmsPackageId:     state[config.SuiChainSelector].MCMSPackageID,
		FastMcmsPackageId: state[config.SuiChainSelector].FastCurseMCMSPackageID,
		McmsOwner:         config.McmsOwner,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to deploy malicious receiver for Sui chain %d: %w", config.SuiChainSelector, err)
	}

	seqReports = append(seqReports, []operations.Report[any, any]{deployMaliciousReceiverOp.ToGenericReport()}...)

	return cldf.ChangesetOutput{
		AddressBook: ab,
		DataStore:   ds,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d DeployMaliciousReceiver) VerifyPreconditions(e cldf.Environment, config DeployMaliciousReceiverConfig) error {
	return nil
}
