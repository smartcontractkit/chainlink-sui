package changesets

import (
	"errors"
	"fmt"

	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	rmn_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/rmn"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type CurseUncurseOperationType string

const (
	CurseOperationType   CurseUncurseOperationType = "curse"
	UncurseOperationType CurseUncurseOperationType = "uncurse"
)

type CurseUncurseChainsConfig struct {
	SuiChainSelector   uint64                `yaml:"suiChainSelector"`
	OperationType      string                `yaml:"operationType"`
	IsGlobalCurse      bool                  `yaml:"isGlobalCurse"`
	DestChainSelectors []uint64              `yaml:"destChainSelectors"`
	TimelockConfig     *utils.TimelockConfig `yaml:"timelockConfig,omitempty"`
	// IsFastCurse selects the CurserCap-based fast-curse path.
	// When TimelockConfig is set, the generated proposal targets the fastcurse MCMS instance.
	// When TimelockConfig is nil, the curse is executed directly using CurserCapObjectId (or the registered cap in chain state).
	IsFastCurse bool `yaml:"isFastCurse,omitempty"`
	// CurserCapObjectId is the CurserCap object ID registered in the fast MCMS Registry.
	// Optional if chain state already records the cap; otherwise required when IsFastCurse is true.
	CurserCapObjectId string `yaml:"curserCapObjectId,omitempty"`
}

var _ cldf.ChangeSetV2[CurseUncurseChainsConfig] = CurseUncurseChains{}

type CurseUncurseChains struct{}

func (c CurseUncurseChains) VerifyPreconditions(e cldf.Environment, cfg CurseUncurseChainsConfig) error {
	if cfg.OperationType != string(CurseOperationType) && cfg.OperationType != string(UncurseOperationType) {
		return fmt.Errorf("invalid operation type %s", cfg.OperationType)
	}
	if cfg.IsFastCurse && cfg.OperationType == string(UncurseOperationType) {
		return errors.New("uncurse via fastcurse MCMS is not supported; use slow MCMS")
	}
	if cfg.IsFastCurse {
		if _, err := resolveCurserCapObjectID(e, cfg.SuiChainSelector, cfg.CurserCapObjectId); err != nil {
			return err
		}
	}
	if cfg.IsGlobalCurse {
		if len(cfg.DestChainSelectors) > 0 {
			return errors.New("global curse config must not include destination selectors")
		}
		return nil
	}
	if len(cfg.DestChainSelectors) == 0 {
		return errors.New("no destination chain selectors provided")
	}
	return nil
}

func (c CurseUncurseChains) Apply(e cldf.Environment, cfg CurseUncurseChainsConfig) (cldf.ChangesetOutput, error) {
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	chainState, ok := state[cfg.SuiChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("no Sui chain state for selector %d", cfg.SuiChainSelector)
	}
	if chainState.CCIPAddress == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf("missing CCIP package address for chain %d", cfg.SuiChainSelector)
	}
	if chainState.CCIPObjectRef == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf("missing CCIP object ref for chain %d", cfg.SuiChainSelector)
	}
	if chainState.CCIPOwnerCapObjectId == "" {
		return cldf.ChangesetOutput{}, fmt.Errorf("missing CCIP owner cap object id for chain %d", cfg.SuiChainSelector)
	}

	suiChain, ok := e.BlockChains.SuiChains()[cfg.SuiChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("no Sui chain client for selector %d", cfg.SuiChainSelector)
	}

	subjects, err := buildCurseSubjects(cfg)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	deps := sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: suiChain.Signer,
		GetCallOpts: func() *bind.CallOpts {
			gasBudget := uint64(400_000_000)
			return &bind.CallOpts{WaitForExecution: true, GasBudget: &gasBudget}
		},
		SuiRPC: suiChain.URL,
	}

	// If timelock proposal is to be generated, disable signer in deps
	if cfg.TimelockConfig != nil {
		deps.Signer = nil
	}

	input := rmn_ops.CurseUncurseChainInput{
		CCIPPackageId:    chainState.CCIPAddress,
		StateObjectId:    chainState.CCIPObjectRef,
		OwnerCapObjectId: chainState.CCIPOwnerCapObjectId,
		Subjects:         subjects,
	}

	var genericReport operations.Report[any, any]
	if cfg.OperationType == string(UncurseOperationType) {
		report, execErr := operations.ExecuteOperation(e.OperationsBundle, rmn_ops.UncurseChainOp, deps, input)
		if execErr != nil {
			return cldf.ChangesetOutput{}, execErr
		}
		genericReport = report.ToGenericReport()
	} else if cfg.IsFastCurse {
		curserCapObjectID, err := resolveCurserCapObjectID(e, cfg.SuiChainSelector, cfg.CurserCapObjectId)
		if err != nil {
			return cldf.ChangesetOutput{}, err
		}
		fastInput := rmn_ops.CurseWithCurserCapInput{
			CCIPPackageId:     chainState.CCIPAddress,
			StateObjectId:     chainState.CCIPObjectRef,
			CurserCapObjectId: curserCapObjectID,
			Subjects:          subjects,
		}
		report, execErr := operations.ExecuteOperation(e.OperationsBundle, rmn_ops.CurseWithCurserCapOp, deps, fastInput)
		if execErr != nil {
			return cldf.ChangesetOutput{}, execErr
		}
		genericReport = report.ToGenericReport()
	} else {
		report, execErr := operations.ExecuteOperation(e.OperationsBundle, rmn_ops.CurseChainOp, deps, input)
		if execErr != nil {
			return cldf.ChangesetOutput{}, execErr
		}
		genericReport = report.ToGenericReport()
	}

	mcmsProposal := mcms.TimelockProposal{}
	if cfg.TimelockConfig != nil {
		defs := []cld_ops.Definition{genericReport.Def}
		inputs := []any{genericReport.Input}

		mcmsState := state[cfg.SuiChainSelector].MCMSState(cfg.IsFastCurse)
		mcmsConfig := mcmsops.ProposalGenerateInput{
			ChainSelector:      cfg.SuiChainSelector,
			Defs:               defs,
			Inputs:             inputs,
			MmcsPackageID:      mcmsState.PackageID,
			McmsStateObjID:     mcmsState.StateObjectID,
			TimelockObjID:      mcmsState.TimelockObjectID,
			AccountObjID:       mcmsState.AccountStateObjectID,
			RegistryObjID:      mcmsState.RegistryObjectID,
			DeployerStateObjID: mcmsState.DeployerStateObjectID,
			TimelockConfig:     *cfg.TimelockConfig,
		}

		result, err := cld_ops.ExecuteSequence(e.OperationsBundle, mcmsops.MCMSDynamicProposalGenerateSeq, deps, mcmsConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal: %w", err)
		}
		mcmsProposal = result.Output
	}

	return cldf.ChangesetOutput{
		Reports:               []operations.Report[any, any]{genericReport},
		MCMSTimelockProposals: []mcms.TimelockProposal{mcmsProposal},
	}, nil
}

func buildCurseSubjects(cfg CurseUncurseChainsConfig) ([][]byte, error) {
	if cfg.IsGlobalCurse {
		s := fastcurse.GlobalCurseSubject()
		return [][]byte{s[:]}, nil
	}
	subjects := make([][]byte, 0, len(cfg.DestChainSelectors))
	for _, selector := range cfg.DestChainSelectors {
		s := fastcurse.GenericSelectorToSubject(selector)
		subjects = append(subjects, s[:])
	}
	return subjects, nil
}

// resolveCurserCapObjectID returns the CurserCap object ID to use for a fast curse.
// When chain state already records the registered cap, cfg may omit the id or must match it.
func resolveCurserCapObjectID(e cldf.Environment, chainSelector uint64, cfgCurserCapObjectID string) (string, error) {
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return "", fmt.Errorf("failed to load Sui chain state: %w", err)
	}
	chainState, ok := state[chainSelector]
	if !ok {
		return "", fmt.Errorf("no Sui chain state for selector %d", chainSelector)
	}

	registered := chainState.CurserCapObjectId
	if cfgCurserCapObjectID != "" {
		if registered != "" && cfgCurserCapObjectID != registered {
			return "", fmt.Errorf(
				"curserCapObjectId %q does not match registered CurserCap %q",
				cfgCurserCapObjectID,
				registered,
			)
		}
		return cfgCurserCapObjectID, nil
	}
	if registered == "" {
		return "", errors.New("curserCapObjectId is required when isFastCurse is true and no CurserCap is recorded in chain state")
	}
	return registered, nil
}
