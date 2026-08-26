package changesets

import (
	"fmt"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
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

type UpgradeRegistryConfig struct {
	SuiChainSelector    uint64                `yaml:"suiChainSelector"`
	CCIPPackageId       string                `yaml:"ccipPackageId"`
	LatestCCIPPackageId string                `yaml:"latestCCIPPackageId,omitempty"`
	StateObjectId       string                `yaml:"stateObjectId"`
	OwnerCapObjectId    string                `yaml:"ownerCapObjectId"`
	TimelockConfig      *utils.TimelockConfig `yaml:"timelockConfig,omitempty"`
	// ReplaceExisting allows this changeset to take the datastore key that is already
	// recorded, as re-initializing an upgrade registry does. Without it, an occupied key is
	// an error raised before anything is deployed.
	ReplaceExisting bool `yaml:"replaceExisting"`
}

var _ cldf.ChangeSetV2[UpgradeRegistryConfig] = UpgradeRegistry{}

type UpgradeRegistry struct{}

// Apply implements deployment.ChangeSetV2.
// With timelockConfig set, the op runs without a signer and an MCMS timelock proposal is
// produced; the UpgradeRegistry object is only created when that proposal executes on-chain,
// so the address-book save is deferred to post-execution state generation. Without
// timelockConfig, the environment signer initializes the registry directly and the resulting
// object ID is saved to the address book here.
func (d UpgradeRegistry) Apply(e cldf.Environment, config UpgradeRegistryConfig) (cldf.ChangesetOutput, error) {
	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}
	chainState := state[config.SuiChainSelector]

	// When the upgraded package ID is not provided, fall back to the latest CCIP package
	// ID recorded in the address book so the call executes against the upgraded bytecode.
	if config.LatestCCIPPackageId == "" {
		config.LatestCCIPPackageId = chainState.LatestCCIPPackageID
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

	upgradeRegistryInitializeOp, err := operations.ExecuteOperation(e.OperationsBundle, ccipops.UpgradeRegistryInitializeOp, deps, ccipops.InitUpgradeRegistryInput{
		CCIPPackageId:       config.CCIPPackageId,
		LatestCCIPPackageId: config.LatestCCIPPackageId,
		StateObjectId:       config.StateObjectId,
		OwnerCapObjectId:    config.OwnerCapObjectId,
	})
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to initialize upgrade registry for Sui chain %d: %w", config.SuiChainSelector, err)
	}

	if config.TimelockConfig != nil {
		mcmsConfig := mcmsops.ProposalGenerateInput{
			ChainSelector:      config.SuiChainSelector,
			Defs:               []operations.Definition{upgradeRegistryInitializeOp.Def},
			Inputs:             []any{upgradeRegistryInitializeOp.Input},
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

	// EOA mode: save UpgradeRegistryObjectId to the addressbook now that it exists on-chain.
	ab := cldf.NewMemoryAddressBook()
	ds := fdatastore.NewMemoryDataStore()
	typeAndVersionUpgradeRegistryObjectId := cldf.NewTypeAndVersion(deployment.SuiUpgradeRegistryObjectId, deployment.Version1_0_0)
	err = deployment.SaveSuiAddress(ab, ds.Addresses(), config.SuiChainSelector, upgradeRegistryInitializeOp.Output.Objects.UpgradeRegistryObjectId, typeAndVersionUpgradeRegistryObjectId, deployment.ChainSingletonQualifier)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to save UpgradeRegistryInitializeOp address %s for Sui chain %d: %w", upgradeRegistryInitializeOp.Output.Objects.UpgradeRegistryObjectId, config.SuiChainSelector, err)
	}

	return cldf.ChangesetOutput{
		AddressBook: ab,
		DataStore:   ds,
		Reports:     []operations.Report[any, any]{upgradeRegistryInitializeOp.ToGenericReport()},
	}, nil
}

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d UpgradeRegistry) VerifyPreconditions(e cldf.Environment, config UpgradeRegistryConfig) error {
	return deployment.ValidateNoDatastoreConflicts(e, config.SuiChainSelector, config.ReplaceExisting,
		func() ([]deployment.PlannedRef, error) {
			// With a timelock config the registry object is created by the proposal's
			// execution and recorded by post-execution state generation; Apply writes nothing.
			if config.TimelockConfig != nil {
				return nil, nil
			}
			return []deployment.PlannedRef{
				{Type: deployment.SuiUpgradeRegistryObjectId, Qualifier: deployment.ChainSingletonQualifier},
			}, nil
		})
}
