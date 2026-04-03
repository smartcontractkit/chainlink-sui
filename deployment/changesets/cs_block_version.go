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
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type BlockVersionConfig struct {
	SuiChainSelector uint64                `yaml:"suiChainSelector"`
	ModuleName       string                `yaml:"moduleName"`
	Version          uint8                 `yaml:"version"`
	TimelockConfig   *utils.TimelockConfig `yaml:"timelockConfig,omitempty"`
}

var _ cldf.ChangeSetV2[BlockVersionConfig] = BlockVersion{}

type BlockVersion struct{}

func (d BlockVersion) Apply(e cldf.Environment, config BlockVersionConfig) (cldf.ChangesetOutput, error) {
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	chainState := state[config.SuiChainSelector]
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

	report, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.BlockVersionOp, deps, ccipops.BlockVersionInput{
		CCIPPackageId:    chainState.CCIPAddress,
		StateObjectId:    chainState.CCIPObjectRef,
		OwnerCapObjectId: chainState.CCIPOwnerCapObjectId,
		ModuleName:       config.ModuleName,
		Version:          config.Version,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to block version for Sui chain %d: %w", config.SuiChainSelector, err)
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

func (d BlockVersion) VerifyPreconditions(e cldf.Environment, config BlockVersionConfig) error {
	return nil
}
