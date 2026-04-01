package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

type AddCCIPPackageIdConfig struct {
	SuiChainSelector      uint64 `yaml:"suiChainSelector"`
	CCIPPackageId         string `yaml:"ccipPackageId"`
	CCIPObjectRefObjectId string `yaml:"ccipObjectRefObjectId"`
	OwnerCapObjectId      string `yaml:"ownerCapObjectId"`
	PackageId             string `yaml:"packageId"`
}

var _ cldf.ChangeSetV2[AddCCIPPackageIdConfig] = AddCCIPPackageId{}

type AddCCIPPackageId struct{}

func (d AddCCIPPackageId) Apply(e cldf.Environment, config AddCCIPPackageIdConfig) (cldf.ChangesetOutput, error) {
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

	addPackageIdReport, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.AddPackageIdStateObjectOp, deps, ccipops.AddPackageIdStateObjectInput{
		CCIPPackageId:         config.CCIPPackageId,
		CCIPObjectRefObjectId: config.CCIPObjectRefObjectId,
		OwnerCapObjectId:      config.OwnerCapObjectId,
		PackageId:             config.PackageId,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to add CCIP package ID for Sui chain %d: %w", config.SuiChainSelector, err)
	}

	seqReports = append(seqReports, addPackageIdReport.ToGenericReport())

	return cldf.ChangesetOutput{
		AddressBook: ab,
		Reports:     seqReports,
	}, nil
}

func (d AddCCIPPackageId) VerifyPreconditions(e cldf.Environment, config AddCCIPPackageIdConfig) error {
	return nil
}
