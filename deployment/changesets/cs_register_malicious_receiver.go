package changesets

import (
	"fmt"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

type RegisterMaliciousReceiverConfig struct {
	SuiChainSelector           uint64
	OwnerCapObjectId           string
	CCIPObjectRefObjectId      string
	MaliciousReceiverPackageId string
}

var _ cldf.ChangeSetV2[RegisterMaliciousReceiverConfig] = RegisterMaliciousReceiver{}

type RegisterMaliciousReceiver struct{}

// Apply implements deployment.ChangeSetV2.
func (d RegisterMaliciousReceiver) Apply(e cldf.Environment, config RegisterMaliciousReceiverConfig) (cldf.ChangesetOutput, error) {
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

	registerReceiverOp, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.RegisterMaliciousReceiverOp, deps, ccipops.RegisterMaliciousReceiverInput{
		OwnerCapObjectId:           config.OwnerCapObjectId,
		CCIPObjectRefObjectId:      config.CCIPObjectRefObjectId,
		MaliciousReceiverPackageId: config.MaliciousReceiverPackageId,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to register malicious receiver for Sui chain %d: %w", config.SuiChainSelector, err)
	}

	seqReports = append(seqReports, []operations.Report[any, any]{registerReceiverOp.ToGenericReport()}...)

	return cldf.ChangesetOutput{
		AddressBook: ab,
		DataStore:   ds,
		Reports:     seqReports,
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d RegisterMaliciousReceiver) VerifyPreconditions(e cldf.Environment, config RegisterMaliciousReceiverConfig) error {
	return nil
}
