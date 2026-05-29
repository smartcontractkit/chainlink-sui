package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

type RegisterBrokenReceiverConfig struct {
	SuiChainSelector        uint64
	OwnerCapObjectId        string
	CCIPObjectRefObjectId   string
	BrokenReceiverPackageId string
}

var _ cldf.ChangeSetV2[RegisterBrokenReceiverConfig] = RegisterBrokenReceiver{}

type RegisterBrokenReceiver struct{}

// Apply implements deployment.ChangeSetV2.
func (d RegisterBrokenReceiver) Apply(e cldf.Environment, config RegisterBrokenReceiverConfig) (cldf.ChangesetOutput, error) {
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

	RegisterBrokenReceiverOp, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.RegisterBrokenReceiverOp, deps, ccipops.RegisterBrokenReceiverInput{
		OwnerCapObjectId:        config.OwnerCapObjectId,
		CCIPObjectRefObjectId:   config.CCIPObjectRefObjectId,
		BrokenReceiverPackageId: config.BrokenReceiverPackageId,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to register broken receiver for Sui chain %d: %w", config.SuiChainSelector, err)
	}

	seqReports = append(seqReports, []operations.Report[any, any]{RegisterBrokenReceiverOp.ToGenericReport()}...)

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d RegisterBrokenReceiver) VerifyPreconditions(e cldf.Environment, config RegisterBrokenReceiverConfig) error {
	return nil
}
