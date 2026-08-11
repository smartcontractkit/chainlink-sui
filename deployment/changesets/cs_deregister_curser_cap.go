package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	rmn_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/rmn"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type DeregisterCurserCapConfig struct {
	SuiChainSelector   uint64                `yaml:"suiChainSelector"`
	CurserCapObjectIds []string              `yaml:"curserCapObjectIds"`
	TimelockConfig     *utils.TimelockConfig `yaml:"timelockConfig"`
}

var _ cldf.ChangeSetV2[DeregisterCurserCapConfig] = DeregisterCurserCap{}

type DeregisterCurserCap struct{}

func (c DeregisterCurserCap) VerifyPreconditions(e cldf.Environment, cfg DeregisterCurserCapConfig) error {
	if cfg.TimelockConfig == nil {
		return fmt.Errorf("timelockConfig is required")
	}
	if len(cfg.CurserCapObjectIds) == 0 {
		return fmt.Errorf("curserCapObjectIds must not be empty")
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
		return fmt.Errorf("slow MCMS must be deployed before deregistering CurserCap IDs")
	}
	return nil
}

func (c DeregisterCurserCap) Apply(e cldf.Environment, cfg DeregisterCurserCapConfig) (cldf.ChangesetOutput, error) {
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	chainState, ok := state[cfg.SuiChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("no Sui chain state for selector %d", cfg.SuiChainSelector)
	}

	slowMCMS := chainState.MCMSStateByInstance(deployment.MCMSInstanceSlow)

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

	input := rmn_ops.McmsDeregisterCurserCapIdsInput{
		CCIPPackageId:        chainState.CCIPAddress, // original package = MCMS on-chain identity (proposal target)
		LatestCCIPPackageId:  chainState.LatestCCIPPackageID,
		StateObjectId:        chainState.CCIPObjectRef,
		SlowOwnerCapObjectId: chainState.CCIPOwnerCapObjectId,
		CurserCapObjectIds:   cfg.CurserCapObjectIds,
	}

	report, err := cld_ops.ExecuteOperation(e.OperationsBundle, rmn_ops.McmsDeregisterCurserCapIdsOp, deps, input)
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
