package changesets

import (
	"fmt"
	"strings"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/mcms"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type BlockFunctionConfig struct {
	SuiChainSelector uint64                `yaml:"suiChainSelector"`
	PackageId        string                `yaml:"packageId"`
	LatestPackageId  string                `yaml:"latestPackageId,omitempty"`
	ModuleName       string                `yaml:"moduleName"`
	FunctionName     string                `yaml:"functionName"`
	Version          uint8                 `yaml:"version"`
	TimelockConfig   *utils.TimelockConfig `yaml:"timelockConfig,omitempty"`
}

var _ cldf.ChangeSetV2[BlockFunctionConfig] = BlockFunction{}

type BlockFunction struct{}

func (d BlockFunction) Apply(e cldf.Environment, config BlockFunctionConfig) (cldf.ChangesetOutput, error) {
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	chainState := state[config.SuiChainSelector]

	// When the upgraded package ID is not provided, fall back to the latest CCIP package
	// ID recorded in the address book so the call executes against the upgraded bytecode.
	if config.LatestPackageId == "" {
		config.LatestPackageId = chainState.LatestCCIPPackageID
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

	report, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.BlockFunctionOp, deps, ccipops.BlockFunctionInput{
		PackageId:        config.PackageId,
		LatestPackageId:  config.LatestPackageId,
		StateObjectId:    chainState.CCIPObjectRef,
		OwnerCapObjectId: chainState.CCIPOwnerCapObjectId,
		ModuleName:       config.ModuleName,
		FunctionName:     config.FunctionName,
		Version:          config.Version,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to block function for Sui chain %d: %w", config.SuiChainSelector, err)
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

func (d BlockFunction) VerifyPreconditions(e cldf.Environment, config BlockFunctionConfig) error {
	if strings.TrimSpace(config.PackageId) == "" {
		return fmt.Errorf("packageId is required")
	}
	return nil
}
