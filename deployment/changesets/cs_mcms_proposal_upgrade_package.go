package changesets

import (
	"fmt"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
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
	IsFastCurse bool
}

type MCMSProposalUpgradePackage struct{}

func (d MCMSProposalUpgradePackage) Apply(e cldf.Environment, config UpgradePackageConfig) (cldf.ChangesetOutput, error) {
	suiState, err := deployment.LoadOnchainStatesui(e)
	if err != nil {
		return cldf.ChangesetOutput{}, fmt.Errorf("failed to load onchain state: %w", err)
	}

	mcmsState := suiState[config.ChainSelector].MCMSState(config.IsFastCurse)

	// Get necessary MCMS state from onchain AB
	if config.MmcsPackageID == "" || config.McmsStateObjID == "" || config.TimelockObjID == "" || config.AccountObjID == "" || config.RegistryObjID == "" {
		config.MmcsPackageID = mcmsState.PackageID
		config.McmsStateObjID = mcmsState.StateObjectID
		config.TimelockObjID = mcmsState.TimelockObjectID
		config.AccountObjID = mcmsState.AccountStateObjectID
		config.RegistryObjID = mcmsState.RegistryObjectID
	}

	// Backfill dependency named addresses from on-chain state when the caller
	// omitted them. CCIP links both mcms and fast_mcms, and CompilePackage now
	// hard-errors if fast_mcms is missing when compiling against published CCIP.
	chainState := suiState[config.ChainSelector]
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
		config.NamedAddresses["ccip"] = chainState.CCIPAddress
	}

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

// VerifyPreconditions implements deployment.ChangeSetV2.
func (d MCMSProposalUpgradePackage) VerifyPreconditions(e cldf.Environment, config UpgradePackageConfig) error {
	return nil
}
