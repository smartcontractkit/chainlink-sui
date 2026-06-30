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

type BackfillLocalDecimalsConfig struct {
	SuiChainSelector uint64                `yaml:"suiChainSelector"`
	CCIPPackageId    string                `yaml:"ccipPackageId,omitempty"`
	StateObjectId    string                `yaml:"stateObjectId,omitempty"`
	OwnerCapObjectId string                `yaml:"ownerCapObjectId,omitempty"`
	VerifyOnly       bool                  `yaml:"verifyOnly,omitempty"`
	TimelockConfig   *utils.TimelockConfig `yaml:"timelockConfig,omitempty"`
}

var _ cldf.ChangeSetV2[BackfillLocalDecimalsConfig] = BackfillLocalDecimals{}

type BackfillLocalDecimals struct{}

func (d BackfillLocalDecimals) Apply(e cldf.Environment, config BackfillLocalDecimalsConfig) (cldf.ChangesetOutput, error) {
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	chainState, ok := state[config.SuiChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("no onchain state for sui chain %d", config.SuiChainSelector)
	}

	suiChain := e.BlockChains.SuiChains()[config.SuiChainSelector]

	ccipPackageId := strings.TrimSpace(config.CCIPPackageId)
	if ccipPackageId == "" {
		ccipPackageId = strings.TrimSpace(chainState.EffectiveCCIPPackageID())
	}

	stateObjectId := strings.TrimSpace(config.StateObjectId)
	if stateObjectId == "" {
		stateObjectId = strings.TrimSpace(chainState.CCIPObjectRef)
	}

	ownerCapObjectId := strings.TrimSpace(config.OwnerCapObjectId)
	if ownerCapObjectId == "" {
		ownerCapObjectId = strings.TrimSpace(chainState.CCIPOwnerCapObjectId)
	}

	deps := sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: suiChain.Signer,
		GetCallOpts: func() *bind.CallOpts {
			gasBudget := uint64(400_000_000)
			return &bind.CallOpts{
				WaitForExecution: true,
				GasBudget:        &gasBudget,
			}
		},
		SuiRPC: suiChain.URL,
	}

	if config.TimelockConfig != nil {
		deps.Signer = nil
	}

	seqResult, err := operations.ExecuteSequence(
		e.OperationsBundle,
		ccipops.BackfillLocalDecimalsSequence,
		deps,
		ccipops.BackfillLocalDecimalsSeqInput{
			CCIPPackageId:         ccipPackageId,          // latest/effective head = reads + direct-exec binary
			OriginalCCIPPackageId: chainState.CCIPAddress, // original = MCMS on-chain identity (proposal target)
			StateObjectId:         stateObjectId,
			OwnerCapObjectId:      ownerCapObjectId,
			VerifyOnly:            config.VerifyOnly,
		},
	)
	if err != nil {
		reports := seqResult.Output.Reports
		return cldf.ChangesetOutput{Reports: reports}, fmt.Errorf(
			"failed to backfill local decimals for Sui chain %d: %w",
			config.SuiChainSelector,
			err,
		)
	}

	if config.TimelockConfig == nil || len(seqResult.Output.McmsDefs) == 0 {
		return cldf.ChangesetOutput{Reports: seqResult.Output.Reports}, nil
	}

	mcmsConfig := mcmsops.ProposalGenerateInput{
		ChainSelector:      config.SuiChainSelector,
		Defs:               seqResult.Output.McmsDefs,
		Inputs:             seqResult.Output.McmsInputs,
		MmcsPackageID:      chainState.MCMSPackageID,
		McmsStateObjID:     chainState.MCMSStateObjectID,
		TimelockObjID:      chainState.MCMSTimelockObjectID,
		AccountObjID:       chainState.MCMSAccountStateObjectID,
		RegistryObjID:      chainState.MCMSRegistryObjectID,
		DeployerStateObjID: chainState.MCMSDeployerStateObjectID,
		TimelockConfig:     *config.TimelockConfig,
	}
	proposalResult, err := operations.ExecuteSequence(e.OperationsBundle, mcmsops.MCMSDynamicProposalGenerateSeq, deps, mcmsConfig)
	if err != nil {
		return cldf.ChangesetOutput{Reports: seqResult.Output.Reports}, fmt.Errorf("failed to generate MCMS proposal: %w", err)
	}

	return cldf.ChangesetOutput{
		Reports:               seqResult.Output.Reports,
		MCMSTimelockProposals: []mcms.TimelockProposal{proposalResult.Output},
	}, nil
}

func (d BackfillLocalDecimals) VerifyPreconditions(e cldf.Environment, config BackfillLocalDecimalsConfig) error {
	if config.SuiChainSelector == 0 {
		return fmt.Errorf("sui chain selector is required")
	}

	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return err
	}

	chainState, ok := state[config.SuiChainSelector]
	if !ok {
		return fmt.Errorf("no onchain state for sui chain %d", config.SuiChainSelector)
	}

	ccipPackageId := strings.TrimSpace(config.CCIPPackageId)
	if ccipPackageId == "" {
		ccipPackageId = strings.TrimSpace(chainState.EffectiveCCIPPackageID())
	}
	if ccipPackageId == "" {
		return fmt.Errorf("ccip package id is required")
	}

	stateObjectId := strings.TrimSpace(config.StateObjectId)
	if stateObjectId == "" {
		stateObjectId = strings.TrimSpace(chainState.CCIPObjectRef)
	}
	if stateObjectId == "" {
		return fmt.Errorf("state object id is required")
	}

	if !config.VerifyOnly && config.TimelockConfig == nil {
		ownerCapObjectId := strings.TrimSpace(config.OwnerCapObjectId)
		if ownerCapObjectId == "" {
			ownerCapObjectId = strings.TrimSpace(chainState.CCIPOwnerCapObjectId)
		}
		if ownerCapObjectId == "" {
			return fmt.Errorf("owner cap object id is required when backfilling directly")
		}
	}

	if config.TimelockConfig != nil {
		if strings.TrimSpace(chainState.MCMSPackageID) == "" || strings.TrimSpace(chainState.MCMSRegistryObjectID) == "" {
			return fmt.Errorf("MCMS package ID and registry object ID must be present in chain state (selector %d)", config.SuiChainSelector)
		}
		if strings.TrimSpace(chainState.MCMSStateObjectID) == "" ||
			strings.TrimSpace(chainState.MCMSTimelockObjectID) == "" ||
			strings.TrimSpace(chainState.MCMSAccountStateObjectID) == "" ||
			strings.TrimSpace(chainState.MCMSDeployerStateObjectID) == "" {
			return fmt.Errorf("MCMS state, timelock, account, and deployer object IDs are required when timelockConfig is set (selector %d)", config.SuiChainSelector)
		}
	}

	return nil
}
