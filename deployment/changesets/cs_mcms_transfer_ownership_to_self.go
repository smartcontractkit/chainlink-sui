package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
)

var _ cldf.ChangeSetV2[TransferOwnershipMcmsToSelfInput] = TransferOwnershipMcmsToSelf{}

// TransferOwnershipMcmsToSelf calls mcms_account::transfer_ownership_to_self with the
// deployer-held MCMS OwnerCap. This must run before generating/executing the accept + execute
// MCMS ownership proposals so register_entrypoint can store OwnerCap on the registry.
type TransferOwnershipMcmsToSelf struct{}

type TransferOwnershipMcmsToSelfInput struct {
	ChainSelector uint64 `json:"chainSelector" yaml:"chainSelector"`
	IsFastCurse   bool   `json:"isFastCurse,omitempty" yaml:"isFastCurse,omitempty"`
}

func (d TransferOwnershipMcmsToSelf) Apply(e cldf.Environment, config TransferOwnershipMcmsToSelfInput) (cldf.ChangesetOutput, error) {
	suiState, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	mcmsFields := suiState[config.ChainSelector].MCMSState(config.IsFastCurse)
	if mcmsFields.PackageID == "" || mcmsFields.AccountOwnerCapObjectID == "" || mcmsFields.AccountStateObjectID == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf("missing MCMS package, account OwnerCap, or AccountState for chain %d (isFastCurse=%v)",
			config.ChainSelector, config.IsFastCurse)
	}

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

	_, err = cld_ops.ExecuteOperation(e.OperationsBundle, mcmsops.MCMSTransferOwnershipOp, deps, mcmsops.MCMSTransferOwnershipInput{
		McmsPackageID:   mcmsFields.PackageID,
		OwnerCap:        mcmsFields.AccountOwnerCapObjectID,
		AccountObjectID: mcmsFields.AccountStateObjectID,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("transfer MCMS ownership to self: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}

func (d TransferOwnershipMcmsToSelf) VerifyPreconditions(e cldf.Environment, config TransferOwnershipMcmsToSelfInput) error {
	return nil
}
