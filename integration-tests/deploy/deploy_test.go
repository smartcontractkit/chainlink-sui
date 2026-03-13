//go:build integration

package deploy

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"testing"

	"github.com/stretchr/testify/suite"

	cselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_rmn_remote "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/rmn_remote"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	"github.com/smartcontractkit/chainlink-sui/deployment/changesets"
	burnminttokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_burn_mint_token_pool"
	managedtokenpoolops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_managed_token_pool"
	managedtokenops "github.com/smartcontractkit/chainlink-sui/deployment/ops/managed_token"
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
	s.DeployLinkBurnMintTokenPool()
	// Phase 6: Deploy CCIP BnM Token
	s.DeployBnMToken()
	// Phase 7: Deploy Managed Token
	s.DeployManagedToken()
	// Phase 8: Deploy Managed Token Faucet
	s.DeployManagedTokenFaucet()
	// Phase 9: Configure Deployer as Minter
	s.ConfigureDeployerAsMinter()
	// Phase 10: Deploy Managed Token Pool
	s.DeployManagedTokenPool()

	// Phase 11: Test curse/uncurse RMN subjects
	// this test is placed here since we have a full deployment to test against
	curseCfg := changesets.CurseUncurseChainsConfig{
		SuiChainSelector:   SuiChainSelector,
		OperationType:      string(changesets.CurseOperationType),
		IsGlobalCurse:      false,
		DestChainSelectors: []uint64{EVMChainSelector},
	}
	curseOut, err := changesets.CurseUncurseChains{}.Apply(s.env, curseCfg)
	s.Require().NoError(err, "failed to curse RMN subjects")
	s.Require().Len(curseOut.Reports, 1, "expected single curse report")

	s.assertRMNCurseSubjects(EVMChainSelector, true)

	uncurseCfg := curseCfg
	uncurseCfg.OperationType = string(changesets.UncurseOperationType)
	uncurseOut, err := changesets.CurseUncurseChains{}.Apply(s.env, uncurseCfg)
	s.Require().NoError(err, "failed to uncurse RMN subjects")
	s.Require().Len(uncurseOut.Reports, 1, "expected single uncurse report")

	s.assertRMNCurseSubjects(EVMChainSelector, false)

	// Load view and check deployments
	states, err := deployment.LoadOnchainStatesui(s.env)
	state := states[cselectors.SUI_LOCALNET.Selector]
	s.Require().NoError(err, "failed to load on-chain state")
	actualView, err := state.GenerateView(&s.env, cselectors.SUI_LOCALNET.Selector, "sui_localnet")
	s.Require().NoError(err, "failed to generate on-chain view")

	owner, err := s.signer.GetAddress()
	s.Require().NoError(err, "failed to get signer address")

	expectedView := buildExpectedSuiChainView(s, state, owner)

	s.Require().Equal(expectedView, actualView)
}

func (s *DeployTestSuite) DeployMCMS() {
	s.T().Log("Phase 1: Deploying MCMS...")

	out, err := changesets.DeployMCMS{}.Apply(s.env, changesets.DeployMCMSConfig{
		DeployMCMSSeqInput: mcmsops.DeployMCMSSeqInput{
			ChainSelector: SuiChainSelector,
			Bypasser:      GetMCMSConfig(1),
			Proposer:      GetMCMSConfig(1),
			Canceller:     GetMCMSConfig(2),
		},
		IsFastCurse: false,
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

func (s *DeployTestSuite) DeployLinkBurnMintTokenPool() {
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

	owner, err := s.signer.GetAddress()
	s.Require().NoError(err, "failed to get signer address")

	out, err := changesets.DeployCCIPBnMToken{}.Apply(s.env, changesets.DeployCCIPBnMTokenConfig{
		ChainSelector: SuiChainSelector,
		MintAmount:    1000,
		MintToAddress: owner,
	})

	s.Require().NoError(err, "failed to deploy CCIP BnM Token")
	err = s.env.ExistingAddresses.Merge(out.AddressBook)
	s.Require().NoError(err, "failed to merge CCIP BnM Token addresses")
}

func (s *DeployTestSuite) DeployManagedToken() {
	s.T().Log("Phase 7: Deploying Managed Token...")

	// Get the BnM token addresses that were just deployed
	addresses, err := s.env.ExistingAddresses.AddressesForChain(SuiChainSelector)
	s.Require().NoError(err, "failed to get addresses")

	var bnmPackageID, bnmTreasuryCapID string
	for addr, typeAndVersion := range addresses {
		if typeAndVersion.Type == deployment.SuiManagedTokenType {
			if _, exists := typeAndVersion.Labels[changesets.CCIPBnMSymbol]; exists {
				bnmPackageID = addr
			}
		}
		if typeAndVersion.Type == deployment.SuiManagedTokenTreasuryCapIDType {
			if _, exists := typeAndVersion.Labels[changesets.CCIPBnMSymbol]; exists {
				bnmTreasuryCapID = addr
			}
		}
	}
	s.Require().NotEmpty(bnmPackageID, "CCIP BnM token package ID not found")
	s.Require().NotEmpty(bnmTreasuryCapID, "CCIP BnM token treasury cap ID not found")

	// Construct coin type from package ID
	coinType := fmt.Sprintf("%s::ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN", bnmPackageID)

	out, err := changesets.DeployManagedToken{}.Apply(s.env, changesets.DeployManagedTokenConfig{
		ChainSelector: SuiChainSelector,
		DeployAndInitManagedTokenInput: managedtokenops.DeployAndInitManagedTokenInput{
			CoinObjectTypeArg:   coinType,
			TreasuryCapObjectId: bnmTreasuryCapID,
			DenyCapObjectId:     "",
			MinterAddress:       "",
			Allowance:           0,
			IsUnlimited:         true,
		},
	})

	s.Require().NoError(err, "failed to deploy managed token")
	err = s.env.ExistingAddresses.Merge(out.AddressBook)
	s.Require().NoError(err, "failed to merge managed token addresses")
}

func (s *DeployTestSuite) ConfigureDeployerAsMinter() {
	s.T().Log("Phase 9: Configuring Deployer as Minter...")

	// Get the managed token addresses
	addresses, err := s.env.ExistingAddresses.AddressesForChain(SuiChainSelector)
	s.Require().NoError(err, "failed to get addresses")

	var managedTokenPackageID, managedTokenStateID, managedTokenOwnerCapID string
	for addr, typeAndVersion := range addresses {
		if typeAndVersion.Type == deployment.SuiManagedTokenPackageIDType {
			if _, exists := typeAndVersion.Labels[changesets.CCIPBnMSymbol]; exists {
				managedTokenPackageID = addr
			}
		}
		if typeAndVersion.Type == deployment.SuiManagedTokenStateObjectID {
			if _, exists := typeAndVersion.Labels[changesets.CCIPBnMSymbol]; exists {
				managedTokenStateID = addr
			}
		}
		if typeAndVersion.Type == deployment.SuiManagedTokenOwnerCapObjectID {
			if _, exists := typeAndVersion.Labels[changesets.CCIPBnMSymbol]; exists {
				managedTokenOwnerCapID = addr
			}
		}
	}

	s.Require().NotEmpty(managedTokenPackageID, "Managed token package ID not found")
	s.Require().NotEmpty(managedTokenStateID, "Managed token state object ID not found")
	s.Require().NotEmpty(managedTokenOwnerCapID, "Managed token owner cap ID not found")

	// Get the BnM token package ID for coin type
	var bnmPackageID string
	for addr, typeAndVersion := range addresses {
		if typeAndVersion.Type == deployment.SuiManagedTokenType {
			if _, exists := typeAndVersion.Labels[changesets.CCIPBnMSymbol]; exists {
				bnmPackageID = addr
				break
			}
		}
	}
	s.Require().NotEmpty(bnmPackageID, "CCIP BnM token package ID not found")

	coinTypeArg := fmt.Sprintf("%s::ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN", bnmPackageID)

	out, err := changesets.ManagedTokenConfigureNewMinter{}.Apply(s.env, changesets.ManagedTokenConfigureNewMinterConfig{
		SuiChainSelector:      SuiChainSelector,
		StateObjectId:         managedTokenStateID,
		OwnerCapObjectId:      managedTokenOwnerCapID,
		ManagedTokenPackageId: managedTokenPackageID,
		CoinObjectTypeArg:     coinTypeArg,
		MinterAddress:         s.deployerAddr,
		Allowance:             0,
		IsUnlimited:           true,
		Source:                "configure_deployer_as_minter_" + changesets.CCIPBnMSymbol,
	})

	s.Require().NoError(err, "failed to configure deployer as minter")
	err = s.env.ExistingAddresses.Merge(out.AddressBook)
	s.Require().NoError(err, "failed to merge minter configuration addresses")
}

func (s *DeployTestSuite) DeployManagedTokenFaucet() {
	s.T().Log("Phase 8: Deploying Managed Token Faucet...")

	// Get the CCIP BnM token addresses (the underlying token)
	addresses, err := s.env.ExistingAddresses.AddressesForChain(SuiChainSelector)
	s.Require().NoError(err, "failed to get addresses")

	var bnmPackageID string
	for addr, typeAndVersion := range addresses {
		if typeAndVersion.Type == deployment.SuiManagedTokenType {
			if _, exists := typeAndVersion.Labels[changesets.CCIPBnMSymbol]; exists {
				bnmPackageID = addr
				break
			}
		}
	}
	s.Require().NotEmpty(bnmPackageID, "CCIP BnM token package ID not found")

	// Construct coin type from CCIP BnM token package ID (not the managed token)
	coinType := fmt.Sprintf("%s::ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN", bnmPackageID)

	out, err := changesets.DeployManagedTokenFaucet{}.Apply(s.env, changesets.DeployManagedTokenFaucetConfig{
		ChainSelector:   SuiChainSelector,
		TokenSymbol:     changesets.CCIPBnMSymbol,
		CoinType:        coinType,
		MintCapObjectId: "", // a mint cap will be issued if this is empty AND the deployer is the managed token owner
	})

	s.Require().NoError(err, "failed to deploy managed token faucet")
	err = s.env.ExistingAddresses.Merge(out.AddressBook)
	s.Require().NoError(err, "failed to merge managed token faucet addresses")

	// Print the mint cap object ID from DeployManagedTokenFaucet
	addresses, err = s.env.ExistingAddresses.AddressesForChain(SuiChainSelector)
	s.Require().NoError(err, "failed to get addresses")
	for addr, typeAndVersion := range addresses {
		if typeAndVersion.Type == deployment.SuiManagedTokenMinterCapID {
			if _, exists := typeAndVersion.Labels[changesets.CCIPBnMSymbol]; exists {
				s.T().Logf("DeployManagedTokenFaucet mint cap object ID: %s", addr)
			}
		}
	}
}

func (s *DeployTestSuite) DeployManagedTokenPool() {
	s.T().Log("Phase 10: Deploying Managed Token Pool...")

	// Get addresses from previous deployments
	addresses, err := s.env.ExistingAddresses.AddressesForChain(SuiChainSelector)
	s.Require().NoError(err, "failed to get addresses")

	var (
		managedTokenPackageID, managedTokenStateID, managedTokenOwnerCapID string
		bnmCoinMetadataID                                                  string
	)

	for addr, typeAndVersion := range addresses {
		if typeAndVersion.Type == deployment.SuiManagedTokenType {
			if _, exists := typeAndVersion.Labels[changesets.CCIPBnMSymbol]; exists {
				managedTokenPackageID = addr
			}
		}
		if typeAndVersion.Type == deployment.SuiManagedTokenStateObjectID {
			if _, exists := typeAndVersion.Labels[changesets.CCIPBnMSymbol]; exists {
				managedTokenStateID = addr
			}
		}
		if typeAndVersion.Type == deployment.SuiManagedTokenOwnerCapObjectID {
			if _, exists := typeAndVersion.Labels[changesets.CCIPBnMSymbol]; exists {
				managedTokenOwnerCapID = addr
			}
		}
		if typeAndVersion.Type == deployment.SuiManagedTokenCoinMetadataIDType {
			if _, exists := typeAndVersion.Labels[changesets.CCIPBnMSymbol]; exists {
				bnmCoinMetadataID = addr
			}
		}
	}

	s.Require().NotEmpty(managedTokenPackageID, "Managed token package ID not found")
	s.Require().NotEmpty(managedTokenStateID, "Managed token state object ID not found")
	s.Require().NotEmpty(managedTokenOwnerCapID, "Managed token owner cap ID not found")
	s.Require().NotEmpty(bnmCoinMetadataID, "CCIP BnM coin metadata ID not found")
	s.Require().NotEmpty(s.ccipPackageID, "CCIP package ID not found")
	s.Require().NotEmpty(s.ccipObjectRef, "CCIP object ref not found")
	s.Require().NotEmpty(s.mcmsPackageID, "MCMS package ID not found")

	// Find the unused mint cap ID (the one that wasn't consumed by the faucet)
	managedTokenMinterCapID, err := s.findUnusedManagedTokenMinterCapID()
	s.Require().NoError(err, "failed to find unused managed token minter cap ID")
	s.Require().NotEmpty(managedTokenMinterCapID, "Unused managed token minter cap ID not found")

	// Construct coin type from CCIP BnM token package ID
	bnmPackageID := ""
	for addr, typeAndVersion := range addresses {
		if typeAndVersion.Type == deployment.SuiManagedTokenType {
			if _, exists := typeAndVersion.Labels[changesets.CCIPBnMSymbol]; exists {
				bnmPackageID = addr
				break
			}
		}
	}
	s.Require().NotEmpty(bnmPackageID, "CCIP BnM token package ID not found")

	coinTypeArg := fmt.Sprintf("%s::ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN", bnmPackageID)

	tokenPoolOut, err := changesets.DeployTPAndConfigure{}.Apply(s.env, changesets.DeployTPAndConfigureConfig{
		SuiChainSelector: SuiChainSelector,
		TokenPoolTypes:   []deployment.TokenPoolType{deployment.TokenPoolTypeManaged},
		ManagedTPInput: managedtokenpoolops.DeployAndInitManagedTokenPoolInput{
			CCIPPackageId:             s.ccipPackageID,
			ManagedTokenPackageId:     managedTokenPackageID,
			MCMSAddress:               s.mcmsPackageID,
			MCMSOwnerAddress:          s.deployerAddr,
			CoinObjectTypeArg:         coinTypeArg,
			CCIPObjectRefObjectId:     s.ccipObjectRef,
			ManagedTokenStateObjectId: managedTokenStateID,
			ManagedTokenOwnerCapId:    managedTokenOwnerCapID,
			CoinMetadataObjectId:      bnmCoinMetadataID,
			MintCapObjectId:           managedTokenMinterCapID,
			TokenPoolAdministrator:    s.deployerAddr,
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

	s.Require().NoError(err, "failed to deploy managed token pool")
	err = s.env.ExistingAddresses.Merge(tokenPoolOut.AddressBook)
	s.Require().NoError(err, "failed to merge managed token pool addresses")
}

func (s *DeployTestSuite) assertRMNCurseSubjects(selector uint64, expectCursed bool) {
	s.T().Helper()

	s.Require().NotEmpty(s.ccipPackageID, "CCIP package ID not set for RMN assertions")
	s.Require().NotEmpty(s.ccipObjectRef, "CCIP object ref not set for RMN assertions")

	contract, err := module_rmn_remote.NewRmnRemote(s.ccipPackageID, s.client)
	s.Require().NoError(err, "failed to create RMN remote binding")

	callOpts := &bind.CallOpts{Signer: s.signer}
	subjects, err := contract.DevInspect().GetCursedSubjects(s.T().Context(), callOpts, bind.Object{Id: s.ccipObjectRef})
	s.Require().NoError(err, "failed to fetch cursed subjects")

	target := make([]byte, 16)
	binary.BigEndian.PutUint64(target[8:], selector)
	found := false
	for _, subj := range subjects {
		if bytes.Equal(subj, target) {
			found = true
			break
		}
	}

	if expectCursed {
		s.Require().True(found, "expected selector to be cursed")
		return
	}

	s.Require().False(found, "expected selector to be uncursed")
}
