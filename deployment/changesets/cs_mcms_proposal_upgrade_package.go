package changesets

import (
	"fmt"
	"strings"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/mcms"
)

var _ cldf.ChangeSetV2[UpgradePackageConfig] = MCMSProposalUpgradePackage{}

// UpgradePackageConfig wraps UpgradeCCIPInput and adds IsFastCurse.
// When MCMS state fields in UpgradeCCIPInput are left empty, they are
// auto-populated from the on-chain address book using the IsFastCurse flag.
type UpgradePackageConfig struct {
	mcmsops.UpgradeCCIPInput
	IsFastCurse bool `yaml:"isFastCurse,omitempty"`
}

type MCMSProposalUpgradePackage struct{}

func (d MCMSProposalUpgradePackage) Apply(e cldf.Environment, config UpgradePackageConfig) (cldf.ChangesetOutput, error) {
	suiState, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	chainState, ok := suiState[config.ChainSelector]
	if !ok {
		return cldf.ChangesetOutput{}, fmt.Errorf("no Sui chain state for selector %d", config.ChainSelector)
	}

	backfillUpgradePackageConfig(&config, chainState)

	suiChains := e.BlockChains.SuiChains()

	suiChain := suiChains[config.ChainSelector]
	deps := sui_ops.OpTxDeps{
		Client: suiChain.Client,
		Signer: suiChain.Signer,
		GetCallOpts: func() *bind.CallOpts {
			return &bind.CallOpts{}
		},
		SuiRPC: suiChain.URL,
	}
	result, err := cld_ops.ExecuteOperation(e.OperationsBundle, mcmsops.UpgradeCCIPOp, deps, config.UpgradeCCIPInput)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to execute sequence: %w", err)
	}

	return cldf.ChangesetOutput{
		MCMSTimelockProposals: []mcms.TimelockProposal{result.Output},
	}, nil
}

func (d MCMSProposalUpgradePackage) VerifyPreconditions(e cldf.Environment, config UpgradePackageConfig) error {
	if config.ChainSelector == 0 {
		return fmt.Errorf("chainSelector is required")
	}
	if config.PackageName == "" {
		return fmt.Errorf("packageName is required")
	}
	if strings.TrimSpace(config.TargetPackageId) == "" {
		return fmt.Errorf("targetPackageId is required")
	}
	if _, ok := e.BlockChains.SuiChains()[config.ChainSelector]; !ok {
		return fmt.Errorf("no Sui chain client for selector %d", config.ChainSelector)
	}

	state, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return fmt.Errorf("load onchain state: %w", err)
	}
	chainState, ok := state[config.ChainSelector]
	if !ok {
		return fmt.Errorf("no Sui chain state for selector %d", config.ChainSelector)
	}

	mcmsState := chainState.MCMSState(config.IsFastCurse)
	if mcmsState.PackageID == "" || mcmsState.RegistryObjectID == "" {
		return fmt.Errorf("MCMS instance for isFastCurse=%v is not recorded for chain %d", config.IsFastCurse, config.ChainSelector)
	}

	if upgradePackageRequiresFastMCMS(config.PackageName) {
		fastMcms := ""
		if config.NamedAddresses != nil {
			fastMcms = config.NamedAddresses["fast_mcms"]
		}
		if fastMcms == "" {
			fastMcms = chainState.FastCurseMCMSPackageID
		}
		if fastMcms == "" {
			return fmt.Errorf(
				"fast_mcms package ID is required to upgrade %q on chain %d; deploy fast MCMS first or set namedAddresses.fast_mcms",
				config.PackageName, config.ChainSelector,
			)
		}
	}

	return nil
}

func backfillUpgradePackageConfig(config *UpgradePackageConfig, chainState deployment.CCIPChainState) {
	mcmsState := chainState.MCMSState(config.IsFastCurse)

	if config.MmcsPackageID == "" {
		config.MmcsPackageID = mcmsState.PackageID
	}
	if config.McmsStateObjID == "" {
		config.McmsStateObjID = mcmsState.StateObjectID
	}
	if config.TimelockObjID == "" {
		config.TimelockObjID = mcmsState.TimelockObjectID
	}
	if config.AccountObjID == "" {
		config.AccountObjID = mcmsState.AccountStateObjectID
	}
	if config.RegistryObjID == "" {
		config.RegistryObjID = mcmsState.RegistryObjectID
	}
	if config.DeployerStateObjID == "" {
		config.DeployerStateObjID = mcmsState.DeployerStateObjectID
	}
	if config.OwnerCapObjID == "" {
		config.OwnerCapObjID = mcmsState.AccountOwnerCapObjectID
	}

	if config.NamedAddresses == nil {
		config.NamedAddresses = map[string]string{}
	}
	if config.NamedAddresses["mcms"] == "" {
		config.NamedAddresses["mcms"] = chainState.MCMSPackageID
	}
	if config.NamedAddresses["fast_mcms"] == "" {
		config.NamedAddresses["fast_mcms"] = chainState.FastCurseMCMSPackageID
	}
	if config.NamedAddresses["ccip"] == "" {
		config.NamedAddresses["ccip"] = chainState.EffectiveCCIPPackageID()
	}
	if config.PackageName == contracts.CCIP && config.NamedAddresses["original_ccip_pkg"] == "" && chainState.CCIPAddress != "" {
		config.NamedAddresses["original_ccip_pkg"] = chainState.CCIPAddress
	}
}

func upgradePackageRequiresFastMCMS(pkg contracts.Package) bool {
	switch pkg {
	case contracts.CCIP,
		contracts.CCIPOnramp,
		contracts.CCIPOfframp,
		contracts.CCIPRouter,
		contracts.CCIPDummyReceiver,
		contracts.CCIPBrokenReceiver,
		contracts.LockReleaseTokenPool,
		contracts.BurnMintTokenPool,
		contracts.ManagedTokenPool:
		return true
	default:
		return false
	}
}
