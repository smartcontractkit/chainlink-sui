package contracts

import (
	"embed"
	"path/filepath"
)

//go:embed ccip mcms test test_secondary link
var Embed embed.FS

type Package string

const (
	// CCIP
	CCIP                   = Package("ccip")
	CCIPBrokenReceiver     = Package("ccip_broken_receiver")
	CCIPDummyReceiver      = Package("ccip_dummy_receiver")
	CCIPMaliciousReceiver  = Package("ccip_malicious_receiver")
	CCIPMaliciousTokenPool = Package("ccip_malicious_token_pool")
	CCIPOfframp            = Package("ccip_offramp")
	CCIPOnramp             = Package("ccip_onramp")
	CCIPRouter             = Package("ccip_router")
	LockReleaseTokenPool   = Package("lock_release_token_pool")
	BurnMintTokenPool      = Package("burn_mint_token_pool")
	ManagedTokenPool       = Package("managed_token_pool")
	USDCTokenPool          = Package("usdc_token_pool")
	ManagedToken           = Package("managed_token")
	ManagedTokenFaucet     = Package("managed_token_faucet")
	MockLinkToken          = Package("mock_link_token")
	MockEthToken           = Package("mock_eth_token")
	// LINK
	LINK    = Package("link")
	CCIPBnM = Package("ccip_burn_mint_token")
	CCIPLnR = Package("ccip_lock_release_token")
	// MCMS
	MCMS       = Package("mcms")
	FastMCMS   = Package("fast_mcms")
	MCMSUser   = Package("mcms_user")
	MCMSUserV2 = Package("mcms_user_v2")
	// Other
	Test          = Package("test")
	TestSecondary = Package("test_secondary")
)

// Contracts maps packages to their respective root directories within Embed
var Contracts map[Package]string = map[Package]string{
	// CCIP
	CCIP:                   filepath.Join("ccip", "ccip"),
	CCIPBrokenReceiver:     filepath.Join("ccip", "ccip_broken_receiver"),
	CCIPDummyReceiver:      filepath.Join("ccip", "ccip_dummy_receiver"),
	CCIPMaliciousReceiver:  filepath.Join("ccip", "ccip_malicious_receiver"),
	CCIPMaliciousTokenPool: filepath.Join("ccip", "ccip_malicious_token_pool"),
	CCIPBnM:                filepath.Join("ccip", "ccip_burn_mint_token"),
	CCIPLnR:                filepath.Join("ccip", "ccip_lock_release_token"),
	CCIPOfframp:            filepath.Join("ccip", "ccip_offramp"),
	CCIPOnramp:             filepath.Join("ccip", "ccip_onramp"),
	CCIPRouter:             filepath.Join("ccip", "ccip_router"),
	LockReleaseTokenPool:   filepath.Join("ccip", "ccip_token_pools", "lock_release_token_pool"),
	BurnMintTokenPool:      filepath.Join("ccip", "ccip_token_pools", "burn_mint_token_pool"),
	ManagedTokenPool:       filepath.Join("ccip", "ccip_token_pools", "managed_token_pool"),
	ManagedToken:           filepath.Join("ccip", "managed_token"),
	ManagedTokenFaucet:     filepath.Join("ccip", "managed_token_faucet"),
	MockLinkToken:          filepath.Join("ccip", "mock_link_token"),
	MockEthToken:           filepath.Join("ccip", "mock_eth_token"),
	// LINK
	LINK: filepath.Join("link"),
	// MCMS
	MCMS:       filepath.Join("mcms", "mcms"),
	FastMCMS:   filepath.Join("mcms", "fast_mcms"),
	MCMSUser:   filepath.Join("mcms", "mcms_test"),
	MCMSUserV2: filepath.Join("mcms", "mcms_test_v2"),
	// Other
	//nolint:gocritic // we need to handle these paths for tests
	Test: filepath.Join("test"),
	//nolint:gocritic // we need to handle these paths for tests
	TestSecondary: filepath.Join("test_secondary"),
}
