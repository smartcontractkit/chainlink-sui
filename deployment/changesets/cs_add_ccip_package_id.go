package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	offrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_offramp"
	onrampops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_onramp"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type AddPackageIdTarget string

const (
	AddPackageIdTargetCCIP    AddPackageIdTarget = "ccip"
	AddPackageIdTargetOnRamp  AddPackageIdTarget = "onramp"
	AddPackageIdTargetOffRamp AddPackageIdTarget = "offramp"
)

type AddCCIPPackageIdConfig struct {
	SuiChainSelector       uint64                `yaml:"suiChainSelector"`
	Target                 AddPackageIdTarget    `yaml:"target"`
	CCIPPackageId          string                `yaml:"ccipPackageId"`                 // original CCIP package ID
	LatestCCIPPackageId    string                `yaml:"latestCcipPackageId,omitempty"` // optional: upgraded CCIP binary
	OnRampPackageId        string                `yaml:"onRampPackageId,omitempty"`
	LatestOnRampPackageId  string                `yaml:"latestOnRampPackageId,omitempty"` // optional: upgraded OnRamp binary
	OffRampPackageId       string                `yaml:"offRampPackageId,omitempty"`
	LatestOffRampPackageId string                `yaml:"latestOffRampPackageId,omitempty"` // optional: upgraded OffRamp binary
	PackageId              string                `yaml:"packageId"`
	TimelockConfig         *utils.TimelockConfig `yaml:"timelockConfig,omitempty"`
}

var _ cldf.ChangeSetV2[AddCCIPPackageIdConfig] = AddCCIPPackageId{}

type AddCCIPPackageId struct{}

func (d AddCCIPPackageId) Apply(e cldf.Environment, config AddCCIPPackageIdConfig) (cldf.ChangesetOutput, error) {
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	chainState := state[config.SuiChainSelector]

	// When the upgraded package IDs are not provided in the config, fall back to the
	// latest package IDs recorded in the address book (addresses.json) so operations
	// execute against the upgraded bytecode without explicit YAML input.
	if config.LatestCCIPPackageId == "" {
		config.LatestCCIPPackageId = chainState.LatestCCIPPackageID
	}
	if config.LatestOnRampPackageId == "" {
		config.LatestOnRampPackageId = chainState.LatestOnRampPackageID
	}
	if config.LatestOffRampPackageId == "" {
		config.LatestOffRampPackageId = chainState.LatestOffRampPackageID
	}

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

	if config.TimelockConfig != nil {
		deps.Signer = nil
	}

	report, err := executeAddPackageId(e, config, chainState, deps)
	if err != nil {
		return cldf.ChangesetOutput{}, err
	}

	if config.TimelockConfig != nil {
		mcmsConfig := mcmsops.ProposalGenerateInput{
			ChainSelector:      config.SuiChainSelector,
			Defs:               []operations.Definition{report.Def},
			Inputs:             []any{report.Input},
			MmcsPackageID:      chainState.MCMSPackageID,
			McmsStateObjID:     chainState.MCMSStateObjectID,
			TimelockObjID:      chainState.MCMSTimelockObjectID,
			AccountObjID:       chainState.MCMSAccountStateObjectID,
			RegistryObjID:      chainState.MCMSRegistryObjectID,
			DeployerStateObjID: chainState.MCMSDeployerStateObjectID,
			TimelockConfig:     *config.TimelockConfig,
		}
		result, err := operations.ExecuteSequence(e.OperationsBundle, mcmsops.MCMSDynamicProposalGenerateSeq, deps, mcmsConfig)
		if err != nil {
			return cldf.ChangesetOutput{}, fmt.Errorf("failed to generate MCMS proposal: %w", err)
		}
		return cldf.ChangesetOutput{
			MCMSTimelockProposals: []mcms.TimelockProposal{result.Output},
		}, nil
	}

	return cldf.ChangesetOutput{
		Reports: []operations.Report[any, any]{report.ToGenericReport()},
	}, nil
}

func executeAddPackageId(e cldf.Environment, config AddCCIPPackageIdConfig, chainState deployment.CCIPChainState, deps sui_ops.OpTxDeps) (operations.Report[any, any], error) {
	target := config.Target
	if target == "" {
		target = AddPackageIdTargetCCIP
	}

	switch target {
	case AddPackageIdTargetCCIP:
		r, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.AddPackageIdStateObjectOp, deps, ccipops.AddPackageIdStateObjectInput{
			PackageId:             config.CCIPPackageId,
			LatestPackageId:       config.LatestCCIPPackageId,
			CCIPObjectRefObjectId: chainState.CCIPObjectRef,
			OwnerCapObjectId:      chainState.CCIPOwnerCapObjectId,
			NewPackageId:          config.PackageId,
		})
		if err != nil {
			return operations.Report[any, any]{}, fmt.Errorf("failed to add package ID to CCIP state object for Sui chain %d: %w", config.SuiChainSelector, err)
		}
		return r.ToGenericReport(), nil

	case AddPackageIdTargetOnRamp:
		r, err := operations.ExecuteOperation(e.OperationsBundle, onrampops.AddPackageIdOp, deps, onrampops.AddPackageIdInput{
			PackageId:        config.OnRampPackageId,
			LatestPackageId:  config.LatestOnRampPackageId,
			StateObjectId:    chainState.OnRampStateObjectId,
			OwnerCapObjectId: chainState.OnRampOwnerCapObjectId,
			NewPackageId:     config.PackageId,
		})
		if err != nil {
			return operations.Report[any, any]{}, fmt.Errorf("failed to add package ID to OnRamp for Sui chain %d: %w", config.SuiChainSelector, err)
		}
		return r.ToGenericReport(), nil

	case AddPackageIdTargetOffRamp:
		r, err := operations.ExecuteOperation(e.OperationsBundle, offrampops.AddPackageIdOffRampOp, deps, offrampops.AddPackageIdOffRampInput{
			PackageId:        config.OffRampPackageId,
			LatestPackageId:  config.LatestOffRampPackageId,
			StateObjectId:    chainState.OffRampStateObjectId,
			OwnerCapObjectId: chainState.OffRampOwnerCapId,
			NewPackageId:     config.PackageId,
		})
		if err != nil {
			return operations.Report[any, any]{}, fmt.Errorf("failed to add package ID to OffRamp for Sui chain %d: %w", config.SuiChainSelector, err)
		}
		return r.ToGenericReport(), nil

	default:
		return operations.Report[any, any]{}, fmt.Errorf("unknown target %q: must be one of %q, %q, %q",
			config.Target, AddPackageIdTargetCCIP, AddPackageIdTargetOnRamp, AddPackageIdTargetOffRamp)
	}
}

func (d AddCCIPPackageId) VerifyPreconditions(e cldf.Environment, config AddCCIPPackageIdConfig) error {
	if config.Target != "" && config.Target != AddPackageIdTargetCCIP &&
		config.Target != AddPackageIdTargetOnRamp && config.Target != AddPackageIdTargetOffRamp {
		return fmt.Errorf("invalid target %q: must be one of %q, %q, %q",
			config.Target, AddPackageIdTargetCCIP, AddPackageIdTargetOnRamp, AddPackageIdTargetOffRamp)
	}

	target := config.Target
	if target == "" {
		target = AddPackageIdTargetCCIP
	}
	switch target {
	case AddPackageIdTargetOnRamp:
		if config.OnRampPackageId == "" {
			return fmt.Errorf("onRampPackageId is required when target is %q", AddPackageIdTargetOnRamp)
		}
	case AddPackageIdTargetOffRamp:
		if config.OffRampPackageId == "" {
			return fmt.Errorf("offRampPackageId is required when target is %q", AddPackageIdTargetOffRamp)
		}
	}
	return nil
}
