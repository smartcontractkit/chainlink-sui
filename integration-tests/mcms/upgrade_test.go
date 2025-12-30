//go:build integration

package mcms

import (
	"testing"

	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_user "github.com/smartcontractkit/chainlink-sui/bindings/generated/mcms/mcms_user"
	mcmsuser "github.com/smartcontractkit/chainlink-sui/bindings/packages/mcms/mcms_user"
	"github.com/smartcontractkit/chainlink-sui/contracts"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type UpgradeTestSuite struct {
	MCMSTestSuite

	// MCMSUser contract (v1)
	mcmsUserPackageId      string
	userDataObjectId       string
	userOwnerCapObjectId   string
	userUpgradeCapObjectId string
}

func (s *UpgradeTestSuite) SetupSuite() {
	s.MCMSTestSuite.SetupSuite()

	signerAddress, err := s.signer.GetAddress()
	s.Require().NoError(err, "getting signer address")
	// Deploy MCMSUser v1 contract
	mcmsUserPackage, tx, err := mcmsuser.PublishMCMSUser(s.T().Context(), s.deps.GetCallOpts(), s.client, s.mcmsPackageID, signerAddress, s.deps.SuiRPC)
	s.Require().NoError(err, "Failed to publish MCMS user package")

	require.NotEmpty(s.T(), tx.Digest, "transaction digest should not be empty")

	s.mcmsUserPackageId = mcmsUserPackage.Address()
	s.T().Logf("✅ Deployed MCMSUser v1 package: %s", s.mcmsUserPackageId)

	// Find created objects using the bind helper functions
	userDataObjectId, err := bind.FindObjectIdFromPublishTx(*tx, "mcms_user", "UserData")
	require.NoError(s.T(), err, "finding UserData object")

	ownerCapObjectId, err := bind.FindObjectIdFromPublishTx(*tx, "ownable", "OwnerCap")
	require.NoError(s.T(), err, "finding OwnerCap object")

	upgradeCapObjectId, err := bind.FindObjectIdFromPublishTx(*tx, "package", "UpgradeCap")
	require.NoError(s.T(), err, "finding UpgradeCap object")

	s.userDataObjectId = userDataObjectId
	s.userOwnerCapObjectId = ownerCapObjectId
	s.userUpgradeCapObjectId = upgradeCapObjectId

	s.T().Logf("✅ MCMSUser UserData object: %s", s.userDataObjectId)
	s.T().Logf("✅ MCMSUser OwnerCap object: %s", s.userOwnerCapObjectId)
	s.T().Logf("✅ MCMSUser UpgradeCap object: %s", s.userUpgradeCapObjectId)
	// Register MCMSUser entrypoint with MCMS
	userContract, err := module_user.NewMcmsUser(s.mcmsUserPackageId, s.client)
	require.NoError(s.T(), err, "creating MCMSUser contract")

	_, err = userContract.RegisterMcmsEntrypoint(
		s.T().Context(),
		s.deps.GetCallOpts(),
		bind.Object{Id: s.userOwnerCapObjectId},
		bind.Object{Id: s.registryObj},
		bind.Object{Id: s.userDataObjectId},
	)
	require.NoError(s.T(), err, "registering MCMS entrypoint")
	s.T().Logf("✅ Registered MCMSUser entrypoint with MCMS")

	// Register UpgradeCap with MCMS deployer
	_, err = userContract.RegisterUpgradeCap(
		s.T().Context(),
		s.deps.GetCallOpts(),
		bind.Object{Id: s.deployerStateObj},
		bind.Object{Id: s.userUpgradeCapObjectId},
		bind.Object{Id: s.registryObj},
	)
	require.NoError(s.T(), err, "registering UpgradeCap with MCMS")
	s.T().Logf("✅ Registered UpgradeCap with MCMS deployer")
}

func (s *UpgradeTestSuite) Test_Upgrade_MCMS_User() {
	s.T().Run("Verify initial version", func(t *testing.T) {
		s.VerifyMCMSUserVersion("MCMSUser 1.0.0")
	})

	s.T().Run("Upgrade MCMSUser through MCMS", func(t *testing.T) {
		s.RunUpgradeMCMSUserProposal()
	})

	s.T().Run("Verify upgraded version", func(t *testing.T) {
		s.VerifyMCMSUserVersion("MCMSUser 2.0.0")
	})
}

func (s *UpgradeTestSuite) RunUpgradeMCMSUserProposal() {
	signerAddress, err := s.signer.GetAddress()
	s.Require().NoError(err, "getting signer address")

	// 1. Build upgrade input
	input := mcmsops.UpgradeCCIPInput{
		// Package related
		PackageName:     contracts.MCMSUserV2,
		TargetPackageId: s.mcmsUserPackageId,
		NamedAddresses: map[string]string{
			"mcms":                      s.mcmsPackageID,
			"mcms_owner":                signerAddress,
			"original_mcms_user_v2_pkg": s.mcmsUserPackageId,
			"mcms_test":                 s.mcmsUserPackageId,
		},

		ChainSelector: uint64(s.chainSelector),
		// MCMS related
		MmcsPackageID:      s.mcmsPackageID,
		McmsStateObjID:     s.mcmsObj,
		RegistryObjID:      s.registryObj,
		TimelockObjID:      s.timelockObj,
		AccountObjID:       s.accountObj,
		DeployerStateObjID: s.deployerStateObj,
		OwnerCapObjID:      s.ownerCapObj,

		// Timelock related
		TimelockConfig: utils.TimelockConfig{
			MCMSAction:   types.TimelockActionBypass,
			MinDelay:     0,
			OverrideRoot: false,
		},
	}

	// 2. Execute operation to generate upgrade proposal
	upgradeReport, err := cld_ops.ExecuteOperation(s.bundle, mcmsops.UpgradeCCIPOp, s.deps, input)
	s.Require().NoError(err, "executing upgrade operation")

	timelockProposal := upgradeReport.Output

	s.T().Logf("✅ Generated upgrade proposal: %s", timelockProposal.Description)

	// 3. Execute the upgrade proposal through MCMS
	s.ExecuteProposalE2e(&timelockProposal, s.bypasserConfig, 0)

	s.T().Logf("✅ Successfully upgraded MCMSUser package from %s", s.mcmsUserPackageId)
}

func (s *UpgradeTestSuite) VerifyMCMSUserVersion(expectedVersion string) {
	// Create contract instance (will use the current package version)
	userContract, err := module_user.NewMcmsUser(s.mcmsUserPackageId, s.client)
	require.NoError(s.T(), err, "creating MCMSUser contract")

	// Call type_and_version function
	version, err := userContract.DevInspect().TypeAndVersion(s.T().Context(), s.deps.GetCallOpts())
	require.NoError(s.T(), err, "getting type and version")

	s.T().Logf("✅ MCMSUser version: %s", version)
	require.Equal(s.T(), expectedVersion, version, "version should match expected")
}
