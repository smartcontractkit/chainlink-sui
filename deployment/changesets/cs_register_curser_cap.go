package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/mcms"
	"github.com/smartcontractkit/mcms/types"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	rmn_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/rmn"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type RegisterCurserCapConfig struct {
	SuiChainSelector uint64                `yaml:"suiChainSelector"`
	TimelockConfig   *utils.TimelockConfig `yaml:"timelockConfig"`
}

var _ cldf.ChangeSetV2[RegisterCurserCapConfig] = RegisterCurserCap{}

type RegisterCurserCap struct{}

func (c RegisterCurserCap) VerifyPreconditions(e cldf.Environment, cfg RegisterCurserCapConfig) error {
	if cfg.TimelockConfig == nil {
		return fmt.Errorf("timelockConfig is required")
	}
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return fmt.Errorf("load onchain state: %w", err)
	}
	chainState, ok := state[cfg.SuiChainSelector]
	if !ok {
		return fmt.Errorf("no Sui chain state for selector %d", cfg.SuiChainSelector)
	}
	if !chainState.HasMCMSInstance(deployment.MCMSInstanceSlow) {
		return fmt.Errorf("slow MCMS must be deployed before registering CurserCap")
	}
	if !chainState.HasMCMSInstance(deployment.MCMSInstanceFastCurse) {
		return fmt.Errorf("fastcurse MCMS must be deployed before registering CurserCap")
	}
	return nil
}

func (c RegisterCurserCap) Apply(e cldf.Environment, cfg RegisterCurserCapConfig) (cldf.ChangesetOutput, error) {
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	chainState, ok := state[cfg.SuiChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("no Sui chain state for selector %d", cfg.SuiChainSelector)
	}

	slowMCMS := chainState.MCMSStateByInstance(deployment.MCMSInstanceSlow)
	fastMCMS := chainState.MCMSStateByInstance(deployment.MCMSInstanceFastCurse)
	if slowMCMS.RegistryObjectID == "" || fastMCMS.RegistryObjectID == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf("slow and fast MCMS registry object IDs are required")
	}

	suiChain, ok := e.BlockChains.SuiChains()[cfg.SuiChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("no Sui chain client for selector %d", cfg.SuiChainSelector)
	}

	deps := sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: nil,
		GetCallOpts: func() *bind.CallOpts {
			gasBudget := uint64(400_000_000)
			return &bind.CallOpts{WaitForExecution: true, GasBudget: &gasBudget}
		},
		SuiRPC: suiChain.URL,
	}

	input := rmn_ops.McmsMintAndRegisterCurserCapInput{
		CCIPPackageId:        chainState.EffectiveCCIPPackageID(),
		StateObjectId:        chainState.CCIPObjectRef,
		SlowOwnerCapObjectId: chainState.CCIPOwnerCapObjectId,
		FastRegistryObjectId: fastMCMS.RegistryObjectID,
	}

	report, err := cld_ops.ExecuteOperation(e.OperationsBundle, rmn_ops.McmsMintAndRegisterCurserCapOp, deps, input)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	mcmsConfig := mcmsops.ProposalGenerateInput{
		ChainSelector:      cfg.SuiChainSelector,
		Defs:               []cld_ops.Definition{report.Def},
		Inputs:             []any{report.Input},
		MmcsPackageID:      slowMCMS.PackageID,
		McmsStateObjID:     slowMCMS.StateObjectID,
		TimelockObjID:      slowMCMS.TimelockObjectID,
		AccountObjID:       slowMCMS.AccountStateObjectID,
		RegistryObjID:      slowMCMS.RegistryObjectID,
		DeployerStateObjID: slowMCMS.DeployerStateObjectID,
		TimelockConfig:     *cfg.TimelockConfig,
	}

	result, err := cld_ops.ExecuteSequence(e.OperationsBundle, mcmsops.MCMSDynamicProposalGenerateSeq, deps, mcmsConfig)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal: %w", err)
	}

	return cldf.ChangesetOutput{
		Reports:               []cld_ops.Report[any, any]{report.ToGenericReport()},
		MCMSTimelockProposals: []mcms.TimelockProposal{result.Output},
	}, nil
}

// LastSuccessfulReportTxDigest returns the transaction hash from the last successful
// MCMS timelock execution report.
func LastSuccessfulReportTxDigest(reports []types.TransactionResult) (string, error) {
	for i := len(reports) - 1; i >= 0; i-- {
		if reports[i].Hash != "" {
			return reports[i].Hash, nil
		}
	}
	return "", fmt.Errorf("no transaction hash in execution reports")
}
