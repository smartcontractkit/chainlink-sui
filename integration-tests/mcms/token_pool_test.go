//go:build integration

package mcms

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_token_admin_registry "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/token_admin_registry"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	ccipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	lockreleasetokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_lock_release_token_pool"
	managedtokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_managed_token_pool"
	tokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_token_pool"
	linkops "github.com/smartcontractkit/chainlink-sui/deployment/ops/link"
	managedtokenops "github.com/smartcontractkit/chainlink-sui/deployment/ops/managed_token"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	ownershipops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ownership"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

type TokenPoolTestSuite struct {
	MCMSTestSuite

	// managed token
	managedTokenLinkPackageId string
	managedTokenLinkObjects   linkops.DeployLinkObjects
	managedTokenPackageId     string
	managedTokenObjects       managedtokenops.DeployManagedTokenObjects

	// managed token pool
	managedTokenPoolPackageId string
	managedTokenPoolObjects   managedtokenpoolops.DeployManagedTokenPoolObjects

	// lnr
	lnrPackageId      string
	lnrObjects        lockreleasetokenpoolops.DeployLockReleaseTokenPoolObjects
	lnrTokenPackageId string
	lnrTokenObjects   linkops.DeployLinkObjects

	// bnm
	bnmPackageId string
	bnmObjects   burnminttokenpoolops.DeployBurnMintTokenPoolObjects
}

func (s *TokenPoolTestSuite) SetupSuite() {
	s.MCMSTestSuite.SetupSuite()

	// Deploy another link token to wrap into managed token
	// Deploy LINK

	// build another reporter, common reporter is getting the report from the prev deployment...
	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(
		s.T().Context,
		logger.Test(s.T()),
		reporter,
	)
	linkManagedTokenReport, err := cld_ops.ExecuteOperation(bundle, linkops.DeployLINKOp, s.deps, cld_ops.EmptyInput{})
	require.NoError(s.T(), err, "failed to deploy LINK token")
	s.managedTokenLinkPackageId = linkManagedTokenReport.Output.PackageId
	s.managedTokenLinkObjects = linkManagedTokenReport.Output.Objects

	// TODO: Deploy Managed token
	managedTokenReport, err := cld_ops.ExecuteSequence(s.bundle, managedtokenops.DeployAndInitManagedTokenSequence, s.deps, managedtokenops.DeployAndInitManagedTokenInput{
		ManagedTokenDeployInput: managedtokenops.ManagedTokenDeployInput{
			MCMSAddress:      s.mcmsPackageID,
			MCMSOwnerAddress: s.mcmsPackageID, // mcms is the owner
		},
		CoinObjectTypeArg:   fmt.Sprintf("%s::link::LINK", s.managedTokenLinkPackageId),
		TreasuryCapObjectId: s.managedTokenLinkObjects.TreasuryCapObjectId,
		// configure_new_minter
		MinterAddress: s.mcmsOwnerAddress,
		Allowance:     1000000,
		IsUnlimited:   false,
	})
	s.Require().NoError(err, "failed to deploy managed token")
	s.managedTokenPackageId = managedTokenReport.Output.ManagedTokenPackageId
	s.managedTokenObjects = managedTokenReport.Output.Objects

	// Deploy another link token for lnr token pool
	reporter = cld_ops.NewMemoryReporter()
	bundle = cld_ops.NewBundle(
		s.T().Context,
		logger.Test(s.T()),
		reporter,
	)

	linkReport, err := cld_ops.ExecuteOperation(bundle, linkops.DeployLINKOp, s.deps, cld_ops.EmptyInput{})
	require.NoError(s.T(), err, "failed to deploy LINK token")
	s.lnrTokenPackageId = linkReport.Output.PackageId
	s.lnrTokenObjects = linkReport.Output.Objects

	linkTokenType := fmt.Sprintf("%s::link::LINK", s.linkPackageId)
	secondLinkTokenType := fmt.Sprintf("%s::link::LINK", linkReport.Output.PackageId)
	// Deploy a token pool of each class
	deployInput := tokenpoolops.DeployAndInitAllTokenPoolsInput{
		SuiChainSelector: uint64(s.chainSelector),
		TokenPoolTypes:   []deployment.TokenPoolType{deployment.TokenPoolTypeBurnMint, deployment.TokenPoolTypeLockRelease, deployment.TokenPoolTypeManaged},
		LockReleaseTPInput: lockreleasetokenpoolops.DeployAndInitLockReleaseTokenPoolInput{
			LockReleaseTokenPoolDeployInput: lockreleasetokenpoolops.LockReleaseTokenPoolDeployInput{
				CCIPPackageId:    s.ccipPackageId,
				MCMSAddress:      s.mcmsPackageID,
				MCMSOwnerAddress: s.mcmsPackageID, // mcms is the owner
			},
			CoinObjectTypeArg:      secondLinkTokenType,
			CCIPObjectRefObjectId:  s.ccipObjects.CCIPObjectRefObjectId,
			CoinMetadataObjectId:   linkReport.Output.Objects.CoinMetadataObjectId,
			TreasuryCapObjectId:    linkReport.Output.Objects.TreasuryCapObjectId,
			TokenPoolAdministrator: s.mcmsPackageID,
			Rebalancer:             "0x5555666677778888999900001111222233334444",
			// apply chain updates
			RemoteChainSelectorsToRemove: []uint64{}, // Empty - no chains to remove from new token pool
			RemoteChainSelectorsToAdd:    []uint64{4, 5, 6},
			RemotePoolAddressesToAdd:     [][]string{{"0x1111111111111111111111111111111111111111"}, {"0x2222222222222222222222222222222222222222"}, {"0x3333333333333333333333333333333333333333"}}, // Must match number of chains
			RemoteTokenAddressesToAdd:    []string{"0x4444444444444444444444444444444444444444", "0x5555555555555555555555555555555555555555", "0x6666666666666666666666666666666666666666"},         // Must match number of chains
			// set chain rate limiter configs
			RemoteChainSelectors: []uint64{7, 8, 9},
			OutboundIsEnableds:   []bool{true, false, true},
			OutboundCapacities:   []uint64{1000000, 2000000, 3000000},
			OutboundRates:        []uint64{100, 200, 300},
			InboundIsEnableds:    []bool{false, true, false},
			InboundCapacities:    []uint64{500000, 1500000, 2500000},
			InboundRates:         []uint64{50, 150, 250},
		},
		BurnMintTpInput: burnminttokenpoolops.DeployAndInitBurnMintTokenPoolInput{
			BurnMintTokenPoolDeployInput: burnminttokenpoolops.BurnMintTokenPoolDeployInput{
				CCIPPackageId:    s.ccipPackageId,
				MCMSAddress:      s.mcmsPackageID,
				MCMSOwnerAddress: s.mcmsPackageID, // mcms is the owner
			},
			CoinObjectTypeArg:      linkTokenType,
			CCIPObjectRefObjectId:  s.ccipObjects.CCIPObjectRefObjectId,
			CoinMetadataObjectId:   s.linkObjects.CoinMetadataObjectId,
			TreasuryCapObjectId:    s.linkObjects.TreasuryCapObjectId,
			TokenPoolAdministrator: s.mcmsPackageID,

			// apply chain updates
			RemoteChainSelectorsToRemove: []uint64{}, // Empty - no chains to remove from new token pool
			RemoteChainSelectorsToAdd:    []uint64{4, 5, 6},
			RemotePoolAddressesToAdd:     [][]string{{"0x1111111111111111111111111111111111111111"}, {"0x2222222222222222222222222222222222222222"}, {"0x3333333333333333333333333333333333333333"}}, // Must match number of chains
			RemoteTokenAddressesToAdd:    []string{"0x4444444444444444444444444444444444444444", "0x5555555555555555555555555555555555555555", "0x6666666666666666666666666666666666666666"},         // Must match number of chains
			// set chain rate limiter configs
			RemoteChainSelectors: []uint64{7, 8, 9},
			OutboundIsEnableds:   []bool{true, false, true},
			OutboundCapacities:   []uint64{1000000, 2000000, 3000000},
			OutboundRates:        []uint64{100, 200, 300},
			InboundIsEnableds:    []bool{false, true, false},
			InboundCapacities:    []uint64{500000, 1500000, 2500000},
			InboundRates:         []uint64{50, 150, 250},
		},
		ManagedTPInput: managedtokenpoolops.DeployAndInitManagedTokenPoolInput{
			// deploy
			CCIPPackageId:         s.ccipPackageId,
			ManagedTokenPackageId: s.managedTokenPackageId,
			MCMSAddress:           s.mcmsPackageID,
			MCMSOwnerAddress:      s.mcmsPackageID, // mcms is the owner
			// initialize
			CoinObjectTypeArg:         fmt.Sprintf("%s::link::LINK", s.managedTokenLinkPackageId),
			CCIPObjectRefObjectId:     s.ccipObjects.CCIPObjectRefObjectId,
			ManagedTokenStateObjectId: s.managedTokenObjects.StateObjectId,
			ManagedTokenOwnerCapId:    s.managedTokenObjects.OwnerCapObjectId,
			CoinMetadataObjectId:      s.managedTokenLinkObjects.CoinMetadataObjectId,
			MintCapObjectId:           s.managedTokenObjects.MinterCapObjectId,
			TokenPoolAdministrator:    s.mcmsPackageID,
			// apply chain updates
			RemoteChainSelectorsToRemove: []uint64{}, // Empty - no chains to remove from new token pool
			RemoteChainSelectorsToAdd:    []uint64{4, 5, 6},
			RemotePoolAddressesToAdd:     [][]string{{"0x1111111111111111111111111111111111111111"}, {"0x2222222222222222222222222222222222222222"}, {"0x3333333333333333333333333333333333333333"}}, // Must match number of chains
			RemoteTokenAddressesToAdd:    []string{"0x4444444444444444444444444444444444444444", "0x5555555555555555555555555555555555555555", "0x6666666666666666666666666666666666666666"},         // Must match number of chains
			// set chain rate limiter configs
			RemoteChainSelectors: []uint64{7, 8, 9},
			OutboundIsEnableds:   []bool{true, false, true},
			OutboundCapacities:   []uint64{1000000, 2000000, 3000000},
			OutboundRates:        []uint64{100, 200, 300},
			InboundIsEnableds:    []bool{false, true, false},
			InboundCapacities:    []uint64{500000, 1500000, 2500000},
			InboundRates:         []uint64{50, 150, 250},
		},
	}

	deploymentReport, err := cld_ops.ExecuteSequence(s.bundle, tokenpoolops.DeployAndInitAllTokenPoolsSequence, s.deps, deployInput)
	s.Require().NoError(err, "failed to deploy and initialize token pools")

	s.bnmPackageId = deploymentReport.Output.DeployBurnMintTokenPoolOutput.BurnMintTPPackageID
	s.bnmObjects = deploymentReport.Output.DeployBurnMintTokenPoolOutput.Objects

	s.lnrPackageId = deploymentReport.Output.DeployLockReleaseTokenPoolOutput.LockReleaseTPPackageID
	s.lnrObjects = deploymentReport.Output.DeployLockReleaseTokenPoolOutput.Objects

	s.managedTokenPoolPackageId = deploymentReport.Output.DeployManagedTokenPoolOutput.ManagedTPPackageId
	s.managedTokenPoolObjects = deploymentReport.Output.DeployManagedTokenPoolOutput.Objects
}

func (s *TokenPoolTestSuite) Test_Token_Pool_MCMS() {
	s.T().Run("Transfer ownership of CCIP to MCMS", func(t *testing.T) {
		s.RunOwnershipCCIPTransfer()
	})

	s.T().Run("Transfer ownership of token pools to MCMS", func(t *testing.T) {
		RunOwnershipTokenPoolProposal(s)
	})

	s.T().Run("Run Lock and Release TP config ops through MCMS", func(t *testing.T) {
		RunLnRConfigOpsTokenPoolProposal(s)
		RunTransferAdminTokenPoolProposal(s, s.lnrTokenObjects.CoinMetadataObjectId)
	})

	s.T().Run("Run Burn and Mint TP config ops through MCMS", func(t *testing.T) {
		RunBnMConfigOpsTokenPoolProposal(s)
		RunTransferAdminTokenPoolProposal(s, s.linkObjects.CoinMetadataObjectId)
	})

	s.T().Run("Run Managed Token TP config ops through MCMS", func(t *testing.T) {
		RunManagedConfigOpsTokenPoolProposal(s)
		RunTransferAdminTokenPoolProposal(s, s.managedTokenLinkObjects.CoinMetadataObjectId)
	})

	s.T().Run("Unregister Token Pool through MCMS", func(t *testing.T) {
		RunUnregisterLnRTokenPoolProposal(s)
	})
}

func RunOwnershipTokenPoolProposal(s *TokenPoolTestSuite) {
	// 1. Generate proposal to accept the ownership from MCMS
	proposalInput := mcmsops.ProposalGenerateInput{
		Defs: []cld_ops.Definition{
			lockreleasetokenpoolops.AcceptOwnershipLockReleaseTokenPoolOp.Def(),
			burnminttokenpoolops.AcceptOwnershipBurnMintTokenPoolOp.Def(),
			managedtokenpoolops.AcceptOwnershipManagedTokenPoolOp.Def(),
		},
		Inputs: []any{
			lockreleasetokenpoolops.AcceptOwnershipLockReleaseTokenPoolInput{
				LockReleaseTokenPoolPackageId: s.lnrPackageId,
				TypeArgs:                      []string{fmt.Sprintf("%s::link::LINK", s.lnrTokenPackageId)},
				StateObjectId:                 s.lnrObjects.StateObjectId,
			},
			burnminttokenpoolops.AcceptOwnershipBurnMintTokenPoolInput{
				BurnMintTokenPoolPackageId: s.bnmPackageId,
				TypeArgs:                   []string{fmt.Sprintf("%s::link::LINK", s.linkPackageId)},
				StateObjectId:              s.bnmObjects.StateObjectId,
			},
			managedtokenpoolops.AcceptOwnershipManagedTokenPoolInput{
				ManagedTokenPoolPackageId: s.managedTokenPoolPackageId,
				TypeArgs:                  []string{fmt.Sprintf("%s::link::LINK", s.managedTokenLinkPackageId)},
				StateObjectId:             s.managedTokenPoolObjects.StateObjectId,
			},
		},

		// MCMS related
		MmcsPackageID:      s.mcmsPackageID,
		McmsStateObjID:     s.mcmsObj,
		TimelockObjID:      s.timelockObj,
		AccountObjID:       s.accountObj,
		RegistryObjID:      s.registryObj,
		DeployerStateObjID: s.deployerStateObj,
		ChainSelector:      uint64(s.chainSelector),

		// Proposal
		TimelockConfig: utils.TimelockConfig{
			MCMSAction:   types.TimelockActionBypass,
			MinDelay:     0,
			OverrideRoot: false,
		},
	}

	acceptOwnershipProposalReport, err := cld_ops.ExecuteSequence(s.bundle, mcmsops.MCMSDynamicProposalGenerateSeq, s.deps, proposalInput)
	s.Require().NoError(err, "executing ownership acceptance proposal sequence")

	timelockProposal := acceptOwnershipProposalReport.Output

	s.ExecuteProposalE2e(&timelockProposal, s.bypasserConfig, 0)

	// 2. Execute transfer ownership from original owner
	input := ownershipops.ExecuteOwnershipTransferToMcmsSeqInput{
		BurnMintTokenPool: &burnminttokenpoolops.ExecuteOwnershipTransferToMcmsBurnMintTokenPoolInput{
			BurnMintTokenPoolPackageId: s.bnmPackageId,
			TypeArgs:                   []string{fmt.Sprintf("%s::link::LINK", s.linkPackageId)},
			OwnerCapObjectId:           s.bnmObjects.OwnerCapObjectId,
			StateObjectId:              s.bnmObjects.StateObjectId,
			RegistryObjectId:           s.registryObj,
			To:                         s.mcmsPackageID,
		},
		LockReleaseTokenPool: &lockreleasetokenpoolops.ExecuteOwnershipTransferToMcmsLockReleaseTokenPoolInput{
			LockReleaseTokenPoolPackageId: s.lnrPackageId,
			TypeArgs:                      []string{fmt.Sprintf("%s::link::LINK", s.lnrTokenPackageId)},
			OwnerCapObjectId:              s.lnrObjects.OwnerCapObjectId,
			StateObjectId:                 s.lnrObjects.StateObjectId,
			RegistryObjectId:              s.registryObj,
			To:                            s.mcmsPackageID,
		},
		ManagedTokenPool: &managedtokenpoolops.ExecuteOwnershipTransferToMcmsManagedTokenPoolInput{
			ManagedTokenPoolPackageId: s.managedTokenPoolPackageId,
			TypeArgs:                  []string{fmt.Sprintf("%s::link::LINK", s.managedTokenLinkPackageId)},
			OwnerCapObjectId:          s.managedTokenPoolObjects.OwnerCapObjectId,
			StateObjectId:             s.managedTokenPoolObjects.StateObjectId,
			RegistryObjectId:          s.registryObj,
			To:                        s.mcmsPackageID,
		},
	}

	executeOwnershipReport, err := cld_ops.ExecuteSequence(s.bundle, ownershipops.ExecuteOwnershipTransferToMcmsSequence, s.deps, input)
	s.Require().NoError(err, "executing ownership transfer to MCMS sequence") // ownership transfer is checked inside the op

	s.Require().NotNil(executeOwnershipReport.Output.Results[ownershipops.ContractTypeBurnMintTokenPool], "burn mint token pool ownership transfer tx is nil")
	s.Require().NotNil(executeOwnershipReport.Output.Results[ownershipops.ContractTypeLockReleaseTokenPool], "lock release token pool ownership transfer tx is nil")
	s.Require().NotNil(executeOwnershipReport.Output.Results[ownershipops.ContractTypeManagedTokenPool], "managed token pool ownership transfer tx is nil")
}

func RunLnRConfigOpsTokenPoolProposal(s *TokenPoolTestSuite) {
	lnrProposalInput := mcmsops.ProposalGenerateInput{
		Defs: []cld_ops.Definition{
			// lnr config ops
			lockreleasetokenpoolops.LockReleaseTokenPoolSetAllowlistEnabledOp.Def(),
			lockreleasetokenpoolops.LockReleaseTokenPoolApplyAllowlistUpdatesOp.Def(),
			lockreleasetokenpoolops.LockReleaseTokenPoolApplyChainUpdatesOp.Def(),
			lockreleasetokenpoolops.LockReleaseTokenPoolAddRemotePoolOp.Def(),
			lockreleasetokenpoolops.LockReleaseTokenPoolRemoveRemotePoolOp.Def(),
			lockreleasetokenpoolops.LockReleaseTokenPoolSetChainRateLimiterOp.Def(),
		},
		Inputs: []any{
			// lnr config ops
			lockreleasetokenpoolops.LockReleaseTokenPoolSetAllowlistEnabledInput{
				LockReleaseTokenPoolPackageId: s.lnrPackageId,
				CoinObjectTypeArg:             fmt.Sprintf("%s::link::LINK", s.lnrTokenPackageId),
				StateObjectId:                 s.lnrObjects.StateObjectId,
				OwnerCap:                      s.lnrObjects.OwnerCapObjectId,
				Enabled:                       true,
			},
			lockreleasetokenpoolops.LockReleaseTokenPoolApplyAllowlistUpdatesInput{
				LockReleaseTokenPoolPackageId: s.lnrPackageId,
				CoinObjectTypeArg:             fmt.Sprintf("%s::link::LINK", s.lnrTokenPackageId),
				StateObjectId:                 s.lnrObjects.StateObjectId,
				OwnerCap:                      s.lnrObjects.OwnerCapObjectId,
				Removes:                       []string{},
				Adds:                          []string{"0x1111111111111111111111111111111111111111"},
			},
			lockreleasetokenpoolops.LockReleaseTokenPoolApplyChainUpdatesInput{
				LockReleasePackageId:         s.lnrPackageId,
				CoinObjectTypeArg:            fmt.Sprintf("%s::link::LINK", s.lnrTokenPackageId),
				StateObjectId:                s.lnrObjects.StateObjectId,
				OwnerCap:                     s.lnrObjects.OwnerCapObjectId,
				RemoteChainSelectorsToRemove: []uint64{},
				RemoteChainSelectorsToAdd:    []uint64{421}, // Use a different chain selector to avoid conflicts with initial setup
				RemotePoolAddressesToAdd:     [][]string{{"0x2222222222222222222222222222222222222222"}},
				RemoteTokenAddressesToAdd:    []string{"0x3333333333333333333333333333333333333333"},
			},
			lockreleasetokenpoolops.LockReleaseTokenPoolAddRemotePoolInput{
				LockReleaseTokenPoolPackageId: s.lnrPackageId,
				CoinObjectTypeArg:             fmt.Sprintf("%s::link::LINK", s.lnrTokenPackageId),
				StateObjectId:                 s.lnrObjects.StateObjectId,
				OwnerCap:                      s.lnrObjects.OwnerCapObjectId,
				RemoteChainSelector:           421, // Use the same chain we added in ApplyChainUpdates
				RemotePoolAddress:             "0x4444444444444444444444444444444444444444",
			},
			lockreleasetokenpoolops.LockReleaseTokenPoolRemoveRemotePoolInput{
				LockReleaseTokenPoolPackageId: s.lnrPackageId,
				CoinObjectTypeArg:             fmt.Sprintf("%s::link::LINK", s.lnrTokenPackageId),
				StateObjectId:                 s.lnrObjects.StateObjectId,
				OwnerCap:                      s.lnrObjects.OwnerCapObjectId,
				RemoteChainSelector:           421,                                          // Remove from the same chain we just added to
				RemotePoolAddress:             "0x4444444444444444444444444444444444444444", // Remove the same pool we just added
			},
			lockreleasetokenpoolops.LockReleaseTokenPoolSetChainRateLimiterInput{
				LockReleasePackageId: s.lnrPackageId,
				CoinObjectTypeArg:    fmt.Sprintf("%s::link::LINK", s.lnrTokenPackageId),
				StateObjectId:        s.lnrObjects.StateObjectId,
				OwnerCap:             s.lnrObjects.OwnerCapObjectId,
				RemoteChainSelectors: []uint64{420},
				OutboundIsEnableds:   []bool{true},
				OutboundCapacities:   []uint64{1000000},
				OutboundRates:        []uint64{100000},
				InboundIsEnableds:    []bool{true},
				InboundCapacities:    []uint64{2000000},
				InboundRates:         []uint64{200000},
			},
		},

		// MCMS related
		MmcsPackageID:      s.mcmsPackageID,
		McmsStateObjID:     s.mcmsObj,
		TimelockObjID:      s.timelockObj,
		AccountObjID:       s.accountObj,
		RegistryObjID:      s.registryObj,
		DeployerStateObjID: s.deployerStateObj,

		ChainSelector: uint64(s.chainSelector),
		// Proposal
		TimelockConfig: utils.TimelockConfig{
			MCMSAction:   types.TimelockActionBypass,
			MinDelay:     0,
			OverrideRoot: false,
		},
	}

	// Execute LNR proposal
	proposalReport, err := cld_ops.ExecuteSequence(s.bundle, mcmsops.MCMSDynamicProposalGenerateSeq, s.deps, lnrProposalInput)
	s.Require().NoError(err, "executing ownership acceptance proposal sequence")

	timelockProposal := proposalReport.Output

	s.ExecuteProposalE2e(&timelockProposal, s.bypasserConfig, 0)
}

func RunBnMConfigOpsTokenPoolProposal(s *TokenPoolTestSuite) {
	bnmProposalInput := mcmsops.ProposalGenerateInput{
		Defs: []cld_ops.Definition{
			// bnm config ops
			burnminttokenpoolops.BurnMintTokenPoolSetAllowlistEnabledOp.Def(),
			burnminttokenpoolops.BurnMintTokenPoolApplyAllowlistUpdatesOp.Def(),
			burnminttokenpoolops.BurnMintTokenPoolApplyChainUpdatesOp.Def(),
			burnminttokenpoolops.BurnMintTokenPoolAddRemotePoolOp.Def(),
			burnminttokenpoolops.BurnMintTokenPoolRemoveRemotePoolOp.Def(),
			burnminttokenpoolops.BurnMintTokenPoolSetChainRateLimiterOp.Def(),
		},
		Inputs: []any{
			// bnm config ops
			burnminttokenpoolops.BurnMintTokenPoolSetAllowlistEnabledInput{
				BurnMintPackageId: s.bnmPackageId,
				StateObjectId:     s.bnmObjects.StateObjectId,
				OwnerCap:          s.bnmObjects.OwnerCapObjectId,
				CoinObjectTypeArg: fmt.Sprintf("%s::link::LINK", s.linkPackageId),
				Enabled:           true,
			},
			burnminttokenpoolops.BurnMintTokenPoolApplyAllowlistUpdatesInput{
				BurnMintPackageId: s.bnmPackageId,
				StateObjectId:     s.bnmObjects.StateObjectId,
				OwnerCap:          s.bnmObjects.OwnerCapObjectId,
				CoinObjectTypeArg: fmt.Sprintf("%s::link::LINK", s.linkPackageId),
				Removes:           []string{},
				Adds:              []string{"0x1111111111111111111111111111111111111111"},
			},
			burnminttokenpoolops.BurnMintTokenPoolApplyChainUpdatesInput{
				BurnMintPackageId:            s.bnmPackageId,
				CoinObjectTypeArg:            fmt.Sprintf("%s::link::LINK", s.linkPackageId),
				StateObjectId:                s.bnmObjects.StateObjectId,
				OwnerCap:                     s.bnmObjects.OwnerCapObjectId,
				RemoteChainSelectorsToRemove: []uint64{},
				RemoteChainSelectorsToAdd:    []uint64{419},
				RemotePoolAddressesToAdd:     [][]string{{"0x8888888888888888888888888888888888888888"}}, // Use different pool address to avoid conflict
				RemoteTokenAddressesToAdd:    []string{"0x9999999999999999999999999999999999999999"},     // Use different token address too
			},
			burnminttokenpoolops.BurnMintTokenPoolAddRemotePoolInput{
				BurnMintTokenPoolPackageId: s.bnmPackageId,
				CoinObjectTypeArg:          fmt.Sprintf("%s::link::LINK", s.linkPackageId),
				StateObjectId:              s.bnmObjects.StateObjectId,
				OwnerCap:                   s.bnmObjects.OwnerCapObjectId,
				RemoteChainSelector:        5,                                            // Add to existing chain 5 that already has pools
				RemotePoolAddress:          "0x7777777777777777777777777777777777777777", // New pool address
			},
			burnminttokenpoolops.BurnMintTokenPoolRemoveRemotePoolInput{
				BurnMintPackageId:   s.bnmPackageId,
				StateObjectId:       s.bnmObjects.StateObjectId,
				OwnerCap:            s.bnmObjects.OwnerCapObjectId,
				CoinObjectTypeArg:   fmt.Sprintf("%s::link::LINK", s.linkPackageId),
				RemoteChainSelector: 5,                                            // Remove from existing chain 5
				RemotePoolAddress:   "0x2222222222222222222222222222222222222222", // Remove the pool that exists from initial setup
			},
			burnminttokenpoolops.BurnMintTokenPoolSetChainRateLimiterInput{
				BurnMintPackageId:    s.bnmPackageId,
				CoinObjectTypeArg:    fmt.Sprintf("%s::link::LINK", s.linkPackageId),
				StateObjectId:        s.bnmObjects.StateObjectId,
				OwnerCap:             s.bnmObjects.OwnerCapObjectId,
				RemoteChainSelectors: []uint64{419}, // Use the chain we added in ApplyChainUpdates
				OutboundIsEnableds:   []bool{true},
				OutboundCapacities:   []uint64{1000000},
				OutboundRates:        []uint64{100000},
				InboundIsEnableds:    []bool{true},
				InboundCapacities:    []uint64{2000000},
				InboundRates:         []uint64{200000},
			},
		},

		// MCMS related
		MmcsPackageID:      s.mcmsPackageID,
		McmsStateObjID:     s.mcmsObj,
		TimelockObjID:      s.timelockObj,
		AccountObjID:       s.accountObj,
		RegistryObjID:      s.registryObj,
		DeployerStateObjID: s.deployerStateObj,

		// Proposal
		TimelockConfig: utils.TimelockConfig{
			MCMSAction:   types.TimelockActionBypass,
			MinDelay:     0,
			OverrideRoot: false,
		},

		ChainSelector: uint64(s.chainSelector),
	}

	// Execute BNM proposal
	proposalReport, err := cld_ops.ExecuteSequence(s.bundle, mcmsops.MCMSDynamicProposalGenerateSeq, s.deps, bnmProposalInput)
	s.Require().NoError(err, "executing ownership acceptance proposal sequence")

	timelockProposal := proposalReport.Output

	s.ExecuteProposalE2e(&timelockProposal, s.bypasserConfig, 0)
}

func RunManagedConfigOpsTokenPoolProposal(s *TokenPoolTestSuite) {
	managedProposalInput := mcmsops.ProposalGenerateInput{
		Defs: []cld_ops.Definition{
			// managed config ops
			managedtokenpoolops.ManagedTokenPoolSetAllowlistEnabledOp.Def(),
			managedtokenpoolops.ManagedTokenPoolApplyAllowlistUpdatesOp.Def(),
			managedtokenpoolops.ManagedTokenPoolApplyChainUpdatesOp.Def(),
			managedtokenpoolops.ManagedTokenPoolAddRemotePoolOp.Def(),
			managedtokenpoolops.ManagedTokenPoolRemoveRemotePoolOp.Def(),
			managedtokenpoolops.ManagedTokenPoolSetChainRateLimiterOp.Def(),
		},
		Inputs: []any{
			managedtokenpoolops.ManagedTokenPoolSetAllowlistEnabledInput{
				ManagedTokenPoolPackageId: s.managedTokenPoolPackageId,
				CoinObjectTypeArg:         fmt.Sprintf("%s::link::LINK", s.managedTokenLinkPackageId),
				StateObjectId:             s.managedTokenPoolObjects.StateObjectId,
				OwnerCap:                  s.managedTokenPoolObjects.OwnerCapObjectId,
				Enabled:                   true,
			},
			managedtokenpoolops.ManagedTokenPoolApplyAllowlistUpdatesInput{
				ManagedTokenPoolPackageId: s.managedTokenPoolPackageId,
				CoinObjectTypeArg:         fmt.Sprintf("%s::link::LINK", s.managedTokenLinkPackageId),
				StateObjectId:             s.managedTokenPoolObjects.StateObjectId,
				OwnerCap:                  s.managedTokenPoolObjects.OwnerCapObjectId,
				Removes:                   []string{},
				Adds:                      []string{"0x1111111111111111111111111111111111111111"},
			},
			managedtokenpoolops.ManagedTokenPoolApplyChainUpdatesInput{
				ManagedTokenPoolPackageId:    s.managedTokenPoolPackageId,
				CoinObjectTypeArg:            fmt.Sprintf("%s::link::LINK", s.managedTokenLinkPackageId),
				StateObjectId:                s.managedTokenPoolObjects.StateObjectId,
				OwnerCap:                     s.managedTokenPoolObjects.OwnerCapObjectId,
				RemoteChainSelectorsToRemove: []uint64{},
				RemoteChainSelectorsToAdd:    []uint64{421},
				RemotePoolAddressesToAdd:     [][]string{{"0x8888888888888888888888888888888888888888"}}, // Use different pool address to avoid conflict
				RemoteTokenAddressesToAdd:    []string{"0x9999999999999999999999999999999999999999"},     // Use different token address too
			},
			managedtokenpoolops.ManagedTokenPoolAddRemotePoolInput{
				ManagedTokenPoolPackageId: s.managedTokenPoolPackageId,
				CoinObjectTypeArg:         fmt.Sprintf("%s::link::LINK", s.managedTokenLinkPackageId),
				StateObjectId:             s.managedTokenPoolObjects.StateObjectId,
				OwnerCap:                  s.managedTokenPoolObjects.OwnerCapObjectId,
				RemoteChainSelector:       421, // Use the same chain we added in ApplyChainUpdates
				RemotePoolAddress:         "0x4444444444444444444444444444444444444444",
			},
			managedtokenpoolops.ManagedTokenPoolRemoveRemotePoolInput{
				ManagedTokenPoolPackageId: s.managedTokenPoolPackageId,
				CoinObjectTypeArg:         fmt.Sprintf("%s::link::LINK", s.managedTokenLinkPackageId),
				StateObjectId:             s.managedTokenPoolObjects.StateObjectId,
				OwnerCap:                  s.managedTokenPoolObjects.OwnerCapObjectId,
				RemoteChainSelector:       421,                                          // Remove from the same chain we just added to
				RemotePoolAddress:         "0x4444444444444444444444444444444444444444", // Remove the same pool we just added
			},
			managedtokenpoolops.ManagedTokenPoolSetChainRateLimiterInput{
				ManagedTokenPoolPackageId: s.managedTokenPoolPackageId,
				CoinObjectTypeArg:         fmt.Sprintf("%s::link::LINK", s.managedTokenLinkPackageId),
				StateObjectId:             s.managedTokenPoolObjects.StateObjectId,
				OwnerCap:                  s.managedTokenPoolObjects.OwnerCapObjectId,
				RemoteChainSelectors:      []uint64{419}, // Use the same chain we added in ApplyChainUpdates
				OutboundIsEnableds:        []bool{true},
				OutboundCapacities:        []uint64{1000000},
				OutboundRates:             []uint64{100000},
				InboundIsEnableds:         []bool{true},
				InboundCapacities:         []uint64{2000000},
				InboundRates:              []uint64{200000},
			},
		},

		// MCMS related
		MmcsPackageID:      s.mcmsPackageID,
		McmsStateObjID:     s.mcmsObj,
		TimelockObjID:      s.timelockObj,
		AccountObjID:       s.accountObj,
		RegistryObjID:      s.registryObj,
		DeployerStateObjID: s.deployerStateObj,

		// Proposal
		TimelockConfig: utils.TimelockConfig{
			MCMSAction:   types.TimelockActionBypass,
			MinDelay:     0,
			OverrideRoot: false,
		},

		ChainSelector: uint64(s.chainSelector),
	}

	// Execute Managed proposal
	proposalReport, err := cld_ops.ExecuteSequence(s.bundle, mcmsops.MCMSDynamicProposalGenerateSeq, s.deps, managedProposalInput)
	s.Require().NoError(err, "executing ownership acceptance proposal sequence")

	timelockProposal := proposalReport.Output

	s.ExecuteProposalE2e(&timelockProposal, s.bypasserConfig, 0)
}

func RunTransferAdminTokenPoolProposal(s *TokenPoolTestSuite, coinmetadataAddress string) {
	transferAdminProposalInput := mcmsops.ProposalGenerateInput{
		Defs: []cld_ops.Definition{
			ccipops.TokenAdminRegistryTransferAdminRoleOp.Def(),
		},
		Inputs: []any{
			ccipops.TransferAdminRoleInput{
				CCIPPackageId:       s.ccipPackageId,
				CCIPObjectRef:       s.ccipObjects.CCIPObjectRefObjectId,
				CoinMetadataAddress: coinmetadataAddress,
				NewAdmin:            s.mcmsOwnerAddress, // current signer address
			},
		},

		// MCMS related
		MmcsPackageID:      s.mcmsPackageID,
		McmsStateObjID:     s.mcmsObj,
		TimelockObjID:      s.timelockObj,
		AccountObjID:       s.accountObj,
		RegistryObjID:      s.registryObj,
		DeployerStateObjID: s.deployerStateObj,

		// Proposal
		TimelockConfig: utils.TimelockConfig{
			MCMSAction:   types.TimelockActionBypass,
			MinDelay:     0,
			OverrideRoot: false,
		},

		ChainSelector: uint64(s.chainSelector),
	}

	// Execute transfer admin proposal
	proposalReport, err := cld_ops.ExecuteSequence(s.bundle, mcmsops.MCMSDynamicProposalGenerateSeq, s.deps, transferAdminProposalInput)
	s.Require().NoError(err, "executing transfer admin proposal sequence")

	timelockProposal := proposalReport.Output

	s.ExecuteProposalE2e(&timelockProposal, s.bypasserConfig, 0)

	// Accept the admin role from signer
	contract, err := module_token_admin_registry.NewTokenAdminRegistry(s.ccipPackageId, s.deps.Client)
	s.Require().NoError(err, "creating token admin registry contract instance")

	acceptAdminInput := ccipops.AcceptAdminRoleInput{
		CCIPPackageId:       s.ccipPackageId,
		CoinMetadataAddress: coinmetadataAddress,
		CCIPObjectRef:       s.ccipObjects.CCIPObjectRefObjectId,
	}

	_, err = cld_ops.ExecuteOperation(s.bundle, ccipops.TokenAdminRegistryAcceptAdminRoleOp, s.deps, acceptAdminInput)
	s.Require().NoError(err, "executing accept admin role operation")

	// verify new admin is mcms owner
	isAdmin, err := contract.DevInspect().IsAdministrator(s.T().Context(), s.deps.GetCallOpts(), bind.Object{Id: s.ccipObjects.CCIPObjectRefObjectId}, coinmetadataAddress, s.mcmsOwnerAddress)
	s.Require().NoError(err, "checking if mcms owner is admin")
	s.Require().True(isAdmin, "mcms owner is not admin after transfer")
	s.T().Logf("MCMS owner %s is now admin of coin metadata %s", s.mcmsOwnerAddress, coinmetadataAddress)

	// Transfer again to MCMS
	transferAdminInput := ccipops.TransferAdminRoleInput{
		CCIPPackageId:       s.ccipPackageId,
		CCIPObjectRef:       s.ccipObjects.CCIPObjectRefObjectId,
		CoinMetadataAddress: coinmetadataAddress,
		NewAdmin:            s.mcmsPackageID, // tx signer
	}
	_, err = cld_ops.ExecuteOperation(s.bundle, ccipops.TokenAdminRegistryTransferAdminRoleOp, s.deps, transferAdminInput)
	s.Require().NoError(err, "executing transfer admin role operation")

	// Accept the admin role from MCMS
	acceptAdminProposalInput := mcmsops.ProposalGenerateInput{
		Defs: []cld_ops.Definition{
			ccipops.TokenAdminRegistryAcceptAdminRoleOp.Def(),
		},
		Inputs: []any{
			ccipops.AcceptAdminRoleInput{
				CCIPPackageId:       s.ccipPackageId,
				CCIPObjectRef:       s.ccipObjects.CCIPObjectRefObjectId,
				CoinMetadataAddress: coinmetadataAddress,
			},
		},

		// MCMS related
		MmcsPackageID:      s.mcmsPackageID,
		McmsStateObjID:     s.mcmsObj,
		TimelockObjID:      s.timelockObj,
		AccountObjID:       s.accountObj,
		RegistryObjID:      s.registryObj,
		DeployerStateObjID: s.deployerStateObj,

		// Proposal
		TimelockConfig: utils.TimelockConfig{
			MCMSAction:   types.TimelockActionBypass,
			MinDelay:     0,
			OverrideRoot: false,
		},

		ChainSelector: uint64(s.chainSelector),
	}

	// Execute transfer admin proposal
	proposalReport, err = cld_ops.ExecuteSequence(s.bundle, mcmsops.MCMSDynamicProposalGenerateSeq, s.deps, acceptAdminProposalInput)
	s.Require().NoError(err, "executing transfer admin proposal sequence")

	timelockProposal = proposalReport.Output

	s.ExecuteProposalE2e(&timelockProposal, s.bypasserConfig, 0)

	// verify new admin is mcms owner
	isAdmin, err = contract.DevInspect().IsAdministrator(s.T().Context(), s.deps.GetCallOpts(), bind.Object{Id: s.ccipObjects.CCIPObjectRefObjectId}, coinmetadataAddress, s.mcmsPackageID)
	s.Require().NoError(err, "checking if mcms owner is admin")
	s.Require().True(isAdmin, "mcms is not admin after transfer")
}

func RunUnregisterLnRTokenPoolProposal(s *TokenPoolTestSuite) {
	// verify pool is registered before unregistering
	contract, err := module_token_admin_registry.NewTokenAdminRegistry(s.ccipPackageId, s.deps.Client)
	s.Require().NoError(err, "creating token admin registry contract instance")

	pool, err := contract.DevInspect().GetPool(s.T().Context(), s.deps.GetCallOpts(), bind.Object{Id: s.ccipObjects.CCIPObjectRefObjectId}, s.lnrTokenObjects.CoinMetadataObjectId)
	s.Require().NoError(err, "checking if token pool is registered")
	s.Require().NotEmpty(pool, "token pool is not registered before unregistering")
	s.Require().Equal(s.lnrPackageId, pool, "registered pool package ID does not match expected")

	// Generate unregister proposal
	unregisterProposalInput := mcmsops.ProposalGenerateInput{
		Defs: []cld_ops.Definition{
			ccipops.TokenAdminRegistryUnregisterPoolOp.Def(),
		},
		Inputs: []any{
			ccipops.UnregisterPoolInput{
				CCIPPackageId:       s.ccipPackageId,
				CCIPObjectRef:       s.ccipObjects.CCIPObjectRefObjectId,
				CoinMetadataAddress: s.lnrTokenObjects.CoinMetadataObjectId,
			},
		},

		// MCMS related
		MmcsPackageID:      s.mcmsPackageID,
		McmsStateObjID:     s.mcmsObj,
		TimelockObjID:      s.timelockObj,
		AccountObjID:       s.accountObj,
		RegistryObjID:      s.registryObj,
		DeployerStateObjID: s.deployerStateObj,

		// Proposal
		TimelockConfig: utils.TimelockConfig{
			MCMSAction:   types.TimelockActionBypass,
			MinDelay:     0,
			OverrideRoot: false,
		},

		ChainSelector: uint64(s.chainSelector),
	}

	// Execute unregister proposal
	proposalReport, err := cld_ops.ExecuteSequence(s.bundle, mcmsops.MCMSDynamicProposalGenerateSeq, s.deps, unregisterProposalInput)
	s.Require().NoError(err, "executing unregister token pool proposal sequence")

	timelockProposal := proposalReport.Output

	s.ExecuteProposalE2e(&timelockProposal, s.bypasserConfig, 0)

	// verify pool is unregistered
	pool, err = contract.DevInspect().GetPool(s.T().Context(), s.deps.GetCallOpts(), bind.Object{Id: s.ccipObjects.CCIPObjectRefObjectId}, s.lnrTokenObjects.CoinMetadataObjectId)
	s.Require().NoError(err, "checking if token pool is registered after unregistering")
	s.Require().Equal(pool, "0x0000000000000000000000000000000000000000000000000000000000000000", "token pool is still registered after unregistering")
}
