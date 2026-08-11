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

var _ cldf.ChangeSetV2[InitMCMSOwnershipTransferConfig] = InitMCMSOwnershipTransfer{}

// InitMCMSOwnershipTransferConfig starts the MCMS self-ownership transfer flow by calling
// mcms_account::transfer_ownership_to_self. This is step 1 of 3; accept and execute follow
// via sui_accept_mcms_ownership_to_self and sui_execute_ownership_transfer.
type InitMCMSOwnershipTransferConfig struct {
	ChainSelector uint64 `json:"chainSelector" yaml:"chainSelector"`
	// IsFastCurse selects the fastcurse MCMS instance; otherwise the normal instance is used.
	IsFastCurse bool `json:"isFastCurse,omitempty" yaml:"isFastCurse,omitempty"`
}

type InitMCMSOwnershipTransfer struct{}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (i InitMCMSOwnershipTransfer) VerifyPreconditions(e cldf.Environment, config InitMCMSOwnershipTransferConfig) error {
	if config.ChainSelector == 0 {
		return fmt.Errorf("chainSelector is required")
	}
	return nil
}

// Apply implements deployment.ChangeSetV2.
func (i InitMCMSOwnershipTransfer) Apply(e cldf.Environment, config InitMCMSOwnershipTransferConfig) (cldf.ChangesetOutput, error) {
	suiState, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load sui onchain state: %w", err)
	}

	chainState, ok := suiState[config.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("no Sui chain state for chain selector %d", config.ChainSelector)
	}

	mcmsFields := chainState.MCMSState(config.IsFastCurse)
	if mcmsFields.PackageID == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf("MCMS package ID not found for chain selector %d (isFastCurse=%v)", config.ChainSelector, config.IsFastCurse)
	}
	if mcmsFields.AccountOwnerCapObjectID == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf("MCMS account owner cap not found for chain selector %d (isFastCurse=%v)", config.ChainSelector, config.IsFastCurse)
	}
	if mcmsFields.AccountStateObjectID == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf("MCMS account state object not found for chain selector %d (isFastCurse=%v)", config.ChainSelector, config.IsFastCurse)
	}

	suiChains := e.BlockChains.SuiChains()
	suiChain, ok := suiChains[config.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("no Sui chain found for chain selector %d", config.ChainSelector)
	}

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

	opInput := mcmsops.MCMSTransferOwnershipInput{
		McmsPackageID:   mcmsFields.PackageID,
		OwnerCap:        mcmsFields.AccountOwnerCapObjectID,
		AccountObjectID: mcmsFields.AccountStateObjectID,
	}

	report, err := cld_ops.ExecuteOperation(e.OperationsBundle, mcmsops.MCMSTransferOwnershipOp, deps, opInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to init MCMS ownership transfer for chain %d: %w", config.ChainSelector, err)
	}

	e.Logger.Infow("Initiated MCMS ownership transfer to self",
		"chainSelector", config.ChainSelector,
		"isFastCurse", config.IsFastCurse,
		"mcmsPackageID", mcmsFields.PackageID,
		"digest", report.Output.Digest,
	)

	return cldf.ChangesetOutput{
		Reports: []cld_ops.Report[any, any]{report.ToGenericReport()},
	}, nil
}
