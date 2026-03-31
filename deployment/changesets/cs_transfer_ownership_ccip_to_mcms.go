package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

var _ cldf.ChangeSetV2[TransferOwnershipCCIPToMcmsInput] = TransferOwnershipCCIPToMcms{}

// TransferOwnershipCCIPToMcms initiates CCIP state-object ownership transfer to the MCMS package
// (transfer_ownership on CCIPObjectRef). MCMS must then accept and execute the transfer.
type TransferOwnershipCCIPToMcms struct{}

type TransferOwnershipCCIPToMcmsInput struct {
	ChainSelector uint64 `json:"chainSelector" yaml:"chainSelector"`
	IsFastCurse   bool   `json:"isFastCurse,omitempty" yaml:"isFastCurse,omitempty"`
}

func (d TransferOwnershipCCIPToMcms) Apply(e cldf.Environment, config TransferOwnershipCCIPToMcmsInput) (cldf.ChangesetOutput, error) {
	suiState, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	state := suiState[config.ChainSelector]
	mcmsFields := state.MCMSState(config.IsFastCurse)

	suiChain := e.BlockChains.SuiChains()[config.ChainSelector]
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

	_, err = cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.TransferOwnershipStateObjectOp, deps, ccipops.TransferOwnershipStateObjectInput{
		CCIPPackageId:         state.CCIPAddress,
		CCIPObjectRefObjectId: state.CCIPObjectRef,
		OwnerCapObjectId:      state.CCIPOwnerCapObjectId,
		To:                    mcmsFields.PackageID,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("transfer ownership to MCMS: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}

func (d TransferOwnershipCCIPToMcms) VerifyPreconditions(e cldf.Environment, config TransferOwnershipCCIPToMcmsInput) error {
	return nil
}
