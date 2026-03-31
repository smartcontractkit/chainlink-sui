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

var _ cldf.ChangeSetV2[RegisterCCIPUpgradeCapMcmsInput] = RegisterCCIPUpgradeCapMcms{}

// RegisterCCIPUpgradeCapMcms deposits the CCIP package UpgradeCap into mcms_deployer::DeployerState.
// Preconditions: CCIP package is registered in MCMS registry (after execute_ownership_transfer_to_mcms).
type RegisterCCIPUpgradeCapMcms struct{}

type RegisterCCIPUpgradeCapMcmsInput struct {
	ChainSelector uint64 `json:"chainSelector" yaml:"chainSelector"`
	IsFastCurse   bool   `json:"isFastCurse,omitempty" yaml:"isFastCurse,omitempty"`
}

func (d RegisterCCIPUpgradeCapMcms) Apply(e cldf.Environment, config RegisterCCIPUpgradeCapMcmsInput) (cldf.ChangesetOutput, error) {
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
			b := uint64(1_000_000_000)
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &b,
			}
		},
		SuiRPC: suiChain.URL,
	}

	_, err = cld_ops.ExecuteOperation(e.OperationsBundle, ccipops.RegisterCCIPUpgradeCapMcmsOp, deps, ccipops.RegisterCCIPUpgradeCapMcmsInput{
		CCIPPackageId:         state.CCIPAddress,
		UpgradeCapObjectId:    state.CCIPUpgradeCapObjectId,
		RegistryObjectId:      mcmsFields.RegistryObjectID,
		DeployerStateObjectId: mcmsFields.DeployerStateObjectID,
		McmsPackageId:         mcmsFields.PackageID,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("register CCIP upgrade cap with MCMS: %w", err)
	}

	return cldf.ChangesetOutput{}, nil
}

func (d RegisterCCIPUpgradeCapMcms) VerifyPreconditions(e cldf.Environment, config RegisterCCIPUpgradeCapMcmsInput) error {
	return nil
}
