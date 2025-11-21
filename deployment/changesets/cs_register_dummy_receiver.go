package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

type RegisterDummyReceiverConfig struct {
	SuiChainSelector       uint64
	OwnerCapObjectId       string
	CCIPObjectRefObjectId  string
	DummyReceiverPackageId string
}

var _ cldf.ChangeSetV2[RegisterDummyReceiverConfig] = RegisterDummyReceiver{}

type RegisterDummyReceiver struct{}

// Apply implements deployment.ChangeSetV2.
func (d RegisterDummyReceiver) Apply(e cldf.Environment, config RegisterDummyReceiverConfig) (cldf.ChangesetOutput, error) {
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

	// Run RegisterReceiver Operation
	RegisterReceiverOp, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.RegisterDummyReceiverOp, deps, ccipops.RegisterDummyReceiverInput{
		OwnerCapObjectId:       config.OwnerCapObjectId,
		CCIPObjectRefObjectId:  config.CCIPObjectRefObjectId,
		DummyReceiverPackageId: config.DummyReceiverPackageId,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to register receiver for Sui chain %d: %w", config.SuiChainSelector, err)
	}

	seqReports = append(seqReports, []operations.Report[any, any]{RegisterReceiverOp.ToGenericReport()}...)

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d RegisterDummyReceiver) VerifyPreconditions(e cldf.Environment, config RegisterDummyReceiverConfig) error {
	return nil
}
