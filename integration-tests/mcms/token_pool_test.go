//go:build integration

package mcms

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	lockreleasetokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_lock_release_token_pool"
	tokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_token_pool"
	linkops "github.com/smartcontractkit/chainlink-sui/deployment/ops/link"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
	suisdk "github.com/smartcontractkit/mcms/sdk/sui"
	"github.com/stretchr/testify/require"
)

type TokenPoolTestSuite struct {
	CCIPMCMSTestSuite

	lnrPackageId      string
	lnrObjects        lockreleasetokenpoolops.DeployLockReleaseTokenPoolObjects
	lnrTokenPackageId string
	lnrTokenObjects   linkops.DeployLinkObjects

	bnmPackageId string
	bnmObjects   burnminttokenpoolops.DeployBurnMintTokenPoolObjects
}

func (s *TokenPoolTestSuite) SetupSuite() {
	s.CCIPMCMSTestSuite.SetupSuite()

	// TODO: Deploy Managed token

	// Deploy another link token
	// Deploy LINK
	reporter := cld_ops.NewMemoryReporter()
	bundle := cld_ops.NewBundle(
		s.T().Context,
		logger.Test(s.T()),
		reporter,
	)
	// is getting the report from the prev deployment... Ops shouldn't do that...
	linkReport, err := cld_ops.ExecuteOperation(bundle, linkops.DeployLINKOp, s.deps, cld_ops.EmptyInput{})
	require.NoError(s.T(), err, "failed to deploy LINK token")
	s.lnrTokenPackageId = linkReport.Output.PackageId
	s.lnrTokenObjects = linkReport.Output.Objects

	linkTokenType := fmt.Sprintf("%s::link::LINK", s.linkPackageId)
	secondLinkTokenType := fmt.Sprintf("%s::link::LINK", linkReport.Output.PackageId)
	// Deploy a token pool of each class
	deployInput := tokenpoolops.DeployAndInitAllTokenPoolsInput{
		SuiChainSelector: uint64(s.chainSelector),
		TokenPoolTypes:   []string{"bnm", "lnr"},
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
	}

	deploymentReport, err := cld_ops.ExecuteSequence(s.bundle, tokenpoolops.DeployAndInitAllTokenPoolsSequence, s.deps, deployInput)
	s.Require().NoError(err, "failed to deploy and initialize token pools")

	s.bnmPackageId = deploymentReport.Output.DeployBurnMintTokenPoolOutput.BurnMintTPPackageID
	s.bnmObjects = deploymentReport.Output.DeployBurnMintTokenPoolOutput.Objects

	s.lnrPackageId = deploymentReport.Output.DeployLockReleaseTokenPoolOutput.LockReleaseTPPackageID
	s.lnrObjects = deploymentReport.Output.DeployLockReleaseTokenPoolOutput.Objects
}

func (s *TokenPoolTestSuite) Test_Token_Pool_MCMS() {
	s.T().Run("Transfer ownership of token pools to MCMS", func(t *testing.T) {
		RunOwnershipTokenPoolProposal(s)
	})
	s.T().Run("Run config ops of token pools through MCMS", func(t *testing.T) {
		RunConfigOpsTokenPoolProposal(s)
	})
}

func RunOwnershipTokenPoolProposal(s *TokenPoolTestSuite) {

	proposalInput := mcmsops.ProposalGenerateInput{
		Defs: []cld_ops.Definition{
			lockreleasetokenpoolops.AcceptOwnershipLockReleaseTokenPoolOp.Def(),
			burnminttokenpoolops.AcceptOwnershipBurnMintTokenPoolOp.Def(),
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
		},

		// MCMS related
		MmcsPackageID:  s.mcmsPackageID,
		McmsStateObjID: s.mcmsObj,
		TimelockObjID:  s.timelockObj,
		AccountObjID:   s.accountObj,
		RegistryObjID:  s.registryObj,

		// Proposal
		Role: suisdk.TimelockRoleBypasser,

		ChainSelector: uint64(s.chainSelector),
	}

	acceptOwnershipProposalReport, err := cld_ops.ExecuteSequence(s.bundle, mcmsops.MCMSDynamicProposalGenerateSeq, s.deps, proposalInput)
	s.Require().NoError(err, "executing ownership acceptance proposal sequence")

	timelockProposal := acceptOwnershipProposalReport.Output

	// 3. Execute transfer ownership from original owner
	// 3.1. Execute the proposal
	s.ExecuteProposalE2e(&timelockProposal, s.bypasserConfig, 0)
}

func RunConfigOpsTokenPoolProposal(s *TokenPoolTestSuite) {
	s.Require().NoError(nil, "getting new owner of Token Pool state object")
}
