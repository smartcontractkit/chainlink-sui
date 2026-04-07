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
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

// AddMCMSAllowedModulesConfig extends the MCMS registry allowlist for the CCIP package using
// mcms_registry::add_allowed_modules with proof type state_object::McmsCallback.
//
// With timelockConfig set, the op runs without a signer and an MCMS timelock proposal is produced.
// Without timelockConfig, the environment signer executes the transaction on-chain.
type AddMCMSAllowedModulesConfig struct {
	SuiChainSelector uint64                `yaml:"suiChainSelector"`
	CCIPPackageId    string                `yaml:"ccipPackageId,omitempty"`
	AllowedModules   []string              `yaml:"allowedModules"`
	TimelockConfig   *utils.TimelockConfig `yaml:"timelockConfig,omitempty"`
}

var _ cldf.ChangeSetV2[AddMCMSAllowedModulesConfig] = AddMCMSAllowedModules{}

type AddMCMSAllowedModules struct{}

func (AddMCMSAllowedModules) Apply(e cldf.Environment, config AddMCMSAllowedModulesConfig) (cldf.ChangesetOutput, error) {
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	chainState, ok := state[config.SuiChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("no Sui chain state for selector %d", config.SuiChainSelector)
	}

	suiChain := e.BlockChains.SuiChains()[config.SuiChainSelector]

	ccipPkg := strings.TrimSpace(config.CCIPPackageId)
	if ccipPkg == "" {
		ccipPkg = strings.TrimSpace(chainState.CCIPAddress)
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

	if config.TimelockConfig != nil {
		deps.Signer = nil
	}

	input := mcmsops.AddModulesMCMSInput{
		MCMSPackageId:     chainState.MCMSPackageID,
		MCMSRegistryObjId: chainState.MCMSRegistryObjectID,
		CCIPPackageId:     ccipPkg,
		AllowedModules:    config.AllowedModules,
	}

	report, err := operations.ExecuteOperation(e.OperationsBundle, mcmsops.AddModulesMCMSOp, deps, input)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("add_allowed_modules for chain %d: %w", config.SuiChainSelector, err)
	}

	if config.TimelockConfig != nil {
		mcmsConfig := mcmsops.ProposalGenerateInput{
			ChainSelector:      config.SuiChainSelector,
			Defs:               []operations.Definition{report.Def},
			Inputs:             []any{input},
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

func (AddMCMSAllowedModules) VerifyPreconditions(e cldf.Environment, config AddMCMSAllowedModulesConfig) error {
	if len(config.AllowedModules) == 0 {
		return fmt.Errorf("allowedModules must not be empty")
	}
	for i, m := range config.AllowedModules {
		if strings.TrimSpace(m) == "" {
			return fmt.Errorf("allowedModules[%d] is empty", i)
		}
	}

	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return fmt.Errorf("failed to load onchain state: %w", err)
	}

	chainState, ok := state[config.SuiChainSelector]
	if !ok {
		return fmt.Errorf("no Sui chain state for selector %d", config.SuiChainSelector)
	}

	ccipPkg := strings.TrimSpace(config.CCIPPackageId)
	if ccipPkg == "" {
		ccipPkg = strings.TrimSpace(chainState.CCIPAddress)
	}
	if ccipPkg == "" {
		return fmt.Errorf("ccipPackageId is required when chain state has no CCIPAddress (selector %d)", config.SuiChainSelector)
	}

	if strings.TrimSpace(chainState.MCMSPackageID) == "" || strings.TrimSpace(chainState.MCMSRegistryObjectID) == "" {
		return fmt.Errorf("MCMS package ID and registry object ID must be present in chain state (selector %d)", config.SuiChainSelector)
	}

	if config.TimelockConfig != nil {
		if strings.TrimSpace(chainState.MCMSStateObjectID) == "" ||
			strings.TrimSpace(chainState.MCMSTimelockObjectID) == "" ||
			strings.TrimSpace(chainState.MCMSAccountStateObjectID) == "" ||
			strings.TrimSpace(chainState.MCMSDeployerStateObjectID) == "" {
			return fmt.Errorf("MCMS state, timelock, account, and deployer object IDs are required when timelockConfig is set (selector %d)", config.SuiChainSelector)
		}
	}

	return nil
}
