package changesets

import (
	"testing"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/contracts"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
)

func TestMCMSProposalUpgradePackage_VerifyPreconditions_RequiresFastMCMSForCCIP(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	ab := cldf.NewMemoryAddressBook()
	require.NoError(t, deployment.StoreMCMSInAddressBook(ab, selector, mcmsops.DeployMCMSSeqOutput{
		PackageId: "0xslow_pkg",
		Objects: mcmsops.DeployMCMSObjects{
			McmsMultisigStateObjectId:   "0xslow_state",
			McmsRegistryObjectId:        "0xslow_registry",
			McmsAccountStateObjectId:    "0xslow_account",
			McmsAccountOwnerCapObjectId: "0xslow_owner_cap",
			TimelockObjectId:            "0xslow_timelock",
			McmsDeployerStateObjectId:   "0xslow_deployer",
		},
	}, deployment.MCMSInstanceSlow))

	cs := MCMSProposalUpgradePackage{}
	err := cs.VerifyPreconditions(dualMCMSEnv(t, ab, selector), UpgradePackageConfig{
		UpgradeCCIPInput: mcmsops.UpgradeCCIPInput{
			ChainSelector:   selector,
			PackageName:     contracts.CCIP,
			TargetPackageId: "0xccip_genesis",
		},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "fast_mcms")
}

func TestMCMSProposalUpgradePackage_VerifyPreconditions_SucceedsWithFastMCMSInAddressBook(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	ab := cldf.NewMemoryAddressBook()
	require.NoError(t, deployment.StoreMCMSInAddressBook(ab, selector, mcmsops.DeployMCMSSeqOutput{
		PackageId: "0xslow_pkg",
		Objects: mcmsops.DeployMCMSObjects{
			McmsMultisigStateObjectId:   "0xslow_state",
			McmsRegistryObjectId:        "0xslow_registry",
			McmsAccountStateObjectId:    "0xslow_account",
			McmsAccountOwnerCapObjectId: "0xslow_owner_cap",
			TimelockObjectId:            "0xslow_timelock",
			McmsDeployerStateObjectId:   "0xslow_deployer",
		},
	}, deployment.MCMSInstanceSlow))
	require.NoError(t, deployment.StoreMCMSInAddressBook(ab, selector, mcmsops.DeployMCMSSeqOutput{
		PackageId: "0xfast_pkg",
		Objects: mcmsops.DeployMCMSObjects{
			McmsMultisigStateObjectId:   "0xfast_state",
			McmsRegistryObjectId:        "0xfast_registry",
			McmsAccountStateObjectId:    "0xfast_account",
			McmsAccountOwnerCapObjectId: "0xfast_owner_cap",
			TimelockObjectId:            "0xfast_timelock",
			McmsDeployerStateObjectId:   "0xfast_deployer",
		},
	}, deployment.MCMSInstanceFastCurse))
	require.NoError(t, ab.Save(selector, "0xccip_genesis", cldf.NewTypeAndVersion(deployment.SuiCCIPType, deployment.Version1_0_0)))

	cs := MCMSProposalUpgradePackage{}
	err := cs.VerifyPreconditions(dualMCMSEnv(t, ab, selector), UpgradePackageConfig{
		UpgradeCCIPInput: mcmsops.UpgradeCCIPInput{
			ChainSelector:   selector,
			PackageName:     contracts.CCIP,
			TargetPackageId: "0xccip_genesis",
		},
	})
	require.NoError(t, err)
}

func TestMCMSProposalUpgradePackage_VerifyPreconditions_ExplicitFastMCMSOverride(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	ab := cldf.NewMemoryAddressBook()
	require.NoError(t, deployment.StoreMCMSInAddressBook(ab, selector, mcmsops.DeployMCMSSeqOutput{
		PackageId: "0xslow_pkg",
		Objects: mcmsops.DeployMCMSObjects{
			McmsMultisigStateObjectId:   "0xslow_state",
			McmsRegistryObjectId:        "0xslow_registry",
			McmsAccountStateObjectId:    "0xslow_account",
			McmsAccountOwnerCapObjectId: "0xslow_owner_cap",
			TimelockObjectId:            "0xslow_timelock",
			McmsDeployerStateObjectId:   "0xslow_deployer",
		},
	}, deployment.MCMSInstanceSlow))
	require.NoError(t, ab.Save(selector, "0xccip_genesis", cldf.NewTypeAndVersion(deployment.SuiCCIPType, deployment.Version1_0_0)))

	cs := MCMSProposalUpgradePackage{}
	err := cs.VerifyPreconditions(dualMCMSEnv(t, ab, selector), UpgradePackageConfig{
		UpgradeCCIPInput: mcmsops.UpgradeCCIPInput{
			ChainSelector:   selector,
			PackageName:     contracts.CCIP,
			TargetPackageId: "0xccip_genesis",
			NamedAddresses: map[string]string{
				"fast_mcms": "0xexplicit_fast",
			},
		},
	})
	require.NoError(t, err)
}

func TestMCMSProposalUpgradePackage_VerifyPreconditions_LINKDoesNotRequireFastMCMS(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	ab := cldf.NewMemoryAddressBook()
	require.NoError(t, deployment.StoreMCMSInAddressBook(ab, selector, mcmsops.DeployMCMSSeqOutput{
		PackageId: "0xslow_pkg",
		Objects: mcmsops.DeployMCMSObjects{
			McmsMultisigStateObjectId:   "0xslow_state",
			McmsRegistryObjectId:        "0xslow_registry",
			McmsAccountStateObjectId:    "0xslow_account",
			McmsAccountOwnerCapObjectId: "0xslow_owner_cap",
			TimelockObjectId:            "0xslow_timelock",
			McmsDeployerStateObjectId:   "0xslow_deployer",
		},
	}, deployment.MCMSInstanceSlow))

	cs := MCMSProposalUpgradePackage{}
	err := cs.VerifyPreconditions(dualMCMSEnv(t, ab, selector), UpgradePackageConfig{
		UpgradeCCIPInput: mcmsops.UpgradeCCIPInput{
			ChainSelector:   selector,
			PackageName:     contracts.LINK,
			TargetPackageId: "0xlink_pkg",
		},
	})
	require.NoError(t, err)
}

func TestBackfillUpgradePackageConfig_FillsMCMSAndNamedAddresses(t *testing.T) {
	t.Parallel()

	chainState := deployment.CCIPChainState{
		MCMSPackageID:               "0xslow_pkg",
		MCMSStateObjectID:           "0xslow_state",
		MCMSRegistryObjectID:        "0xslow_registry",
		MCMSAccountStateObjectID:    "0xslow_account",
		MCMSAccountOwnerCapObjectID: "0xslow_owner_cap",
		MCMSTimelockObjectID:        "0xslow_timelock",
		MCMSDeployerStateObjectID:   "0xslow_deployer",
		FastCurseMCMSPackageID:      "0xfast_pkg",
		CCIPAddress:                 "0xccip_genesis",
		LatestCCIPPackageID:         "0xccip_v2",
	}

	cfg := UpgradePackageConfig{
		UpgradeCCIPInput: mcmsops.UpgradeCCIPInput{
			PackageName: contracts.CCIP,
		},
		IsFastCurse: false,
	}

	backfillUpgradePackageConfig(&cfg, chainState)

	require.Equal(t, "0xslow_deployer", cfg.DeployerStateObjID)
	require.Equal(t, "0xslow_owner_cap", cfg.OwnerCapObjID)
	require.Equal(t, "0xslow_pkg", cfg.NamedAddresses["mcms"])
	require.Equal(t, "0xfast_pkg", cfg.NamedAddresses["fast_mcms"])
	require.Equal(t, "0xccip_v2", cfg.NamedAddresses["ccip"])
	require.Equal(t, "0xccip_genesis", cfg.NamedAddresses["original_ccip_pkg"])
}

func TestCCIPChainState_EffectivePackageIDs(t *testing.T) {
	t.Parallel()

	state := deployment.CCIPChainState{
		CCIPAddress:            "0xccip_genesis",
		LatestCCIPPackageID:    "0xccip_v2",
		OnRampAddress:          "0xonramp_genesis",
		LatestOnRampPackageID:  "0xonramp_v2",
		OffRampAddress:         "0xofframp_genesis",
		LatestOffRampPackageID: "0xofframp_v2",
	}
	require.Equal(t, "0xccip_v2", state.EffectiveCCIPPackageID())
	require.Equal(t, "0xonramp_v2", state.EffectiveOnRampPackageID())
	require.Equal(t, "0xofframp_v2", state.EffectiveOffRampPackageID())

	emptyLatest := deployment.CCIPChainState{
		CCIPAddress:    "0xccip_genesis",
		OnRampAddress:  "0xonramp_genesis",
		OffRampAddress: "0xofframp_genesis",
	}
	require.Equal(t, "0xccip_genesis", emptyLatest.EffectiveCCIPPackageID())
	require.Equal(t, "0xonramp_genesis", emptyLatest.EffectiveOnRampPackageID())
	require.Equal(t, "0xofframp_genesis", emptyLatest.EffectiveOffRampPackageID())
}

func TestMCMSProposalUpgradePackage_VerifyPreconditions_RequiresSuiChainClient(t *testing.T) {
	t.Parallel()

	cs := MCMSProposalUpgradePackage{}
	err := cs.VerifyPreconditions(cldf.Environment{
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{}),
	}, UpgradePackageConfig{
		UpgradeCCIPInput: mcmsops.UpgradeCCIPInput{
			ChainSelector:   cselectors.SUI_TESTNET.Selector,
			PackageName:     contracts.LINK,
			TargetPackageId: "0xlink",
		},
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "no Sui chain client")
}
