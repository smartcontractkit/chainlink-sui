package changesets

import (
	"encoding/binary"
	"errors"
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	rmn_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/rmn"
)

type CurseUncurseOperationType string

const (
	CurseOperationType   CurseUncurseOperationType = "curse"
	UncurseOperationType CurseUncurseOperationType = "uncurse"
)

type CurseUncurseChainsConfig struct {
	SuiChainSelector   uint64   `yaml:"suiChainSelector"`
	OperationType      string   `yaml:"operationType"`
	IsGlobalCurse      bool     `yaml:"isGlobalCurse"`
	DestChainSelectors []uint64 `yaml:"destChainSelectors"`
}

var _ cldf.ChangeSetV2[CurseUncurseChainsConfig] = CurseUncurseChains{}

type CurseUncurseChains struct{}

var globalCurseSubjectBytes = [16]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01}

func (c CurseUncurseChains) VerifyPreconditions(e cldf.Environment, cfg CurseUncurseChainsConfig) error {
	if cfg.OperationType != string(CurseOperationType) && cfg.OperationType != string(UncurseOperationType) {
		return fmt.Errorf("invalid operation type %s", cfg.OperationType)
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
	} else {
		report, execErr := operations.ExecuteOperation(e.OperationsBundle, rmn_ops.CurseChainOp, deps, input)
		if execErr != nil {
			return cldf.ChangesetOutput{}, execErr
		}
		genericReport = report.ToGenericReport()
	}

	return cldf.ChangesetOutput{Reports: []operations.Report[any, any]{genericReport}}, nil
}

func buildCurseSubjects(cfg CurseUncurseChainsConfig) ([][]byte, error) {
	if cfg.IsGlobalCurse {
		subject := make([]byte, len(globalCurseSubjectBytes))
		copy(subject, globalCurseSubjectBytes[:])
		return [][]byte{subject}, nil
	}
	subjects := make([][]byte, 0, len(cfg.DestChainSelectors))
	for _, selector := range cfg.DestChainSelectors {
		subjects = append(subjects, selectorToSubject(selector))
	}
	return subjects, nil
}

func selectorToSubject(selector uint64) []byte {
	subject := make([]byte, 16)
	binary.BigEndian.PutUint64(subject[8:], selector)
	return subject
}
