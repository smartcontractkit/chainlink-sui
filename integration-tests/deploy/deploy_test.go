//go:build integration

package deploy

import (
	"fmt"

	"testing"

	"github.com/stretchr/testify/suite"

	cselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-sui/deployment"
	"github.com/smartcontractkit/chainlink-sui/deployment/changesets"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	mcmsops "github.com/smartcontractkit/chainlink-sui/deployment/ops/mcms"
)

// TestDeployAndConfigureSuiChain tests the deployment and configuration of a Sui chain, lane and BnM TP
// using changesets and using views to verify the deployments.
func TestDeployAndConfigureSuiChain(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(DeployTestSuite))
}

func (s *DeployTestSuite) TestDeployAndConfigureSuiChain() {
	// "Phase 1: Deploy MCMS"
	s.DeployMCMS()
	// "Phase 2: Deploy LINK Token"
	s.DeployLink()
	// "Phase 3: Deploy Core CCIP Infrastructure"
	s.DeployCCIPCore()
	// "Phase 4: Connect Lanes"
	s.ConnectLanes()
	// "Phase 5: Deploy Token Pools"
	s.DeployTokenPools()
	// Phase 6: Deploy CCIP BnM Token
	s.DeployBnMToken()

	// Load view and check deployments
	states, err := deployment.LoadOnchainStatesui(s.env)
	state := states[cselectors.SUI_LOCALNET.Selector]
	s.Require().NoError(err, "failed to load on-chain state")
	actualView, err := state.GenerateView(&s.env, cselectors.SUI_LOCALNET.Selector, "sui_localnet")
	s.Require().NoError(err, "failed to generate on-chain view")

	owner, err := s.signer.GetAddress()
	s.Require().NoError(err, "failed to get signer address")

	expectedView := BuildExpectedSuiChainView(s, state, owner)

	s.Require().Equal(expectedView, actualView)
}

func (s *DeployTestSuite) DeployMCMS() {
	s.T().Log("Phase 1: Deploying MCMS...")

	out, err := changesets.DeployMCMS{}.Apply(s.env, mcmsops.DeployMCMSSeqInput{
		ChainSelector: SuiChainSelector,
		Bypasser:      GetMCMSConfig(1),
		Proposer:      GetMCMSConfig(1),
		Canceller:     GetMCMSConfig(2),
	})
	s.Require().NoError(err, "failed to deploy MCMS")

	err = s.env.ExistingAddresses.Merge(out.AddressBook)
	s.Require().NoError(err, "failed to merge MCMS addresses")
}

func (s *DeployTestSuite) DeployLink() {
	s.T().Log("Phase 2: Deploying LINK Token...")

	out, err := changesets.DeployLinkToken{}.Apply(s.env, changesets.DeployLinkTokenConfig{
		ChainSelector: SuiChainSelector,
	})
	s.Require().NoError(err, "failed to deploy LINK token")
	err = s.env.ExistingAddresses.Merge(out.AddressBook)
	s.Require().NoError(err, "failed to merge LINK token addresses")

	// Cache addresses
	addresses, err := s.env.ExistingAddresses.AddressesForChain(SuiChainSelector)
	s.Require().NoError(err, "failed to get addresses")
	for addr, typeAndVersion := range addresses {
		switch typeAndVersion.Type {
		case deployment.SuiLinkTokenType:
			s.linkTokenPackageID = addr
		case deployment.SuiLinkTokenObjectMetadataID:
			s.linkTokenMetadataID = addr
		case deployment.SuiLinkTokenTreasuryCapID:
			s.linkTokenTreasuryCapID = addr
		}
	}
}

func (s *DeployTestSuite) DeployCCIPCore() {
	s.T().Log("Phase 3: Deploying Core CCIP Infrastructure...")

	s.Require().NotEmpty(s.linkTokenMetadataID, "LINK CoinMetadata not found")

	out, err := changesets.DeploySuiChain{}.Apply(s.env, changesets.DeploySuiChainConfig{
		SuiChainSelector:              SuiChainSelector,
		DestChainSelector:             EVMChainSelector,
		DestChainOnRampAddressBytes:   DestChainOnRampAddressBytes,
		LinkTokenCoinMetadataObjectId: s.linkTokenMetadataID,
	})
	s.Require().NoError(err, "failed to deploy CCIP")
	err = s.env.ExistingAddresses.Merge(out.AddressBook)
	s.Require().NoError(err, "failed to merge CCIP addresses")

	// Cache addresses
	s.deployerAddr, err = s.signer.GetAddress()
	s.Require().NoError(err, "failed to get deployer address")
	addresses, err := s.env.ExistingAddresses.AddressesForChain(SuiChainSelector)
	s.Require().NoError(err, "failed to get addresses")
	for addr, typeAndVersion := range addresses {
		switch typeAndVersion.Type {
		case deployment.SuiCCIPType:
			s.ccipPackageID = addr
		case deployment.SuiCCIPObjectRefType:
			s.ccipObjectRef = addr
		case deployment.SuiMcmsPackageIDType:
			s.mcmsPackageID = addr
		}
	}
}

func (s *DeployTestSuite) ConnectLanes() {
	s.T().Log("Phase 4: Connecting Lanes...")

	_, err := changesets.ConnectSuiToEVM{}.Apply(s.env, changesets.ConnectSuiToEVMConfig{
		SuiChainSelector: SuiChainSelector,
		FeeQuoterApplyTokenTransferFeeConfigUpdatesInput:     TokenTransferFeeConfig,
		FeeQuoterApplyDestChainConfigUpdatesInput:            DestChainConfigUpdatesInput,
		FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput: PremiumMultiplierWeiPerEth,
		ApplyDestChainConfigureOnRampInput:                   DestChainConfigureOnRamp,
		ApplySourceChainConfigUpdateInput:                    SourceChainConfigUpdate,
	})
	s.Require().NoError(err, "failed to connect lanes")
}

func (s *DeployTestSuite) DeployTokenPools() {
	s.T().Log("Phase 5: Deploying Token Pools...")

	// Validate addresses are present
	s.Require().NotEmpty(s.linkTokenPackageID, "LINK token package ID not found")
	s.Require().NotEmpty(s.linkTokenMetadataID, "LINK token metadata ID not found")
	s.Require().NotEmpty(s.linkTokenTreasuryCapID, "LINK token treasury cap ID not found")
	s.Require().NotEmpty(s.ccipPackageID, "CCIP package ID not found")
	s.Require().NotEmpty(s.ccipObjectRef, "CCIP object ref not found")
	s.Require().NotEmpty(s.mcmsPackageID, "MCMS package ID not found")

	coinTypeArg := fmt.Sprintf("%s::link::LINK", s.linkTokenPackageID)

	tokenPoolOut, err := changesets.DeployTPAndConfigure{}.Apply(s.env, changesets.DeployTPAndConfigureConfig{
		SuiChainSelector: SuiChainSelector,
		TokenPoolTypes:   []deployment.TokenPoolType{deployment.TokenPoolTypeBurnMint},
		BurnMintTpInput: burnminttokenpoolops.DeployAndInitBurnMintTokenPoolInput{
			BurnMintTokenPoolDeployInput: burnminttokenpoolops.BurnMintTokenPoolDeployInput{
				CCIPPackageId:    s.ccipPackageID,
				MCMSAddress:      s.mcmsPackageID,
				MCMSOwnerAddress: s.deployerAddr,
			},
			CoinObjectTypeArg:      coinTypeArg,
			CCIPObjectRefObjectId:  s.ccipObjectRef,
			CoinMetadataObjectId:   s.linkTokenMetadataID,
			TreasuryCapObjectId:    s.linkTokenTreasuryCapID,
			TokenPoolAdministrator: s.deployerAddr,
			// Remote chain configuration
			RemoteChainSelectorsToRemove: []uint64{},
			RemoteChainSelectorsToAdd:    []uint64{EVMChainSelector},
			RemotePoolAddressesToAdd:     [][]string{{EVMPoolAddress}},
			RemoteTokenAddressesToAdd:    []string{fmt.Sprintf("0x%s", EVMTokenAddress)},
			// Rate limiter configs
			RemoteChainSelectors: []uint64{EVMChainSelector},
			OutboundIsEnableds:   []bool{false},
			OutboundCapacities:   []uint64{RateLimiterCapacity},
			OutboundRates:        []uint64{RateLimiterRate},
			InboundIsEnableds:    []bool{false},
			InboundCapacities:    []uint64{RateLimiterCapacity},
			InboundRates:         []uint64{RateLimiterRate},
		},
	})
	s.Require().NoError(err, "failed to deploy LINK token pool")

	err = s.env.ExistingAddresses.Merge(tokenPoolOut.AddressBook)
	s.Require().NoError(err, "failed to merge LINK token pool addresses")
}

func (s *DeployTestSuite) DeployBnMToken() {
	s.T().Log("Phase 6: Deploying CCIP BnM Token...")

	out, err := changesets.DeployCCIPBnMToken{}.Apply(s.env, changesets.DeployCCIPBnMTokenConfig{
		ChainSelector: SuiChainSelector,
		MintAmount:    1000,
	})
	s.Require().NoError(err, "failed to deploy CCIP BnM Token")
	err = s.env.ExistingAddresses.Merge(out.AddressBook)
	s.Require().NoError(err, "failed to merge CCIP BnM Token addresses")
}
