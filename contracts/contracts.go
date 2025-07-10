package contracts

import (
	"embed"
	"path/filepath"
)

//go:embed ccip mcms test
var Embed embed.FS

type Package string

const (
	// CCIP
	CCIP                 = Package("ccip")
	CCIPDummyReceiver    = Package("ccip_dummy_receiver")
	CCIPOfframp          = Package("ccip_offramp")
	CCIPOnramp           = Package("ccip_onramp")
	CCIPRouter           = Package("ccip_router")
	CCIPTokenPool        = Package("ccip_token_pool")
	LINKToken            = Package("link_token")
	LockReleaseTokenPool = Package("lock_release_token_pool")
	BurnMintTokenPool    = Package("burn_mint_token_pool")
	ManagedTokenPool     = Package("managed_token_pool")
	ManagedToken         = Package("managed_token")
	// MCMS
	MCMS = Package("mcms")
	// Other
	Test = Package("test")
)

// Contracts maps packages to their respective root directories within Embed
var Contracts map[Package]string = map[Package]string{
	// CCIP
	CCIP:                 filepath.Join("ccip", "ccip"),
	CCIPDummyReceiver:    filepath.Join("ccip", "ccip_dummy_receiver"),
	CCIPOfframp:          filepath.Join("ccip", "ccip_offramp"),
	CCIPOnramp:           filepath.Join("ccip", "ccip_onramp"),
	CCIPRouter:           filepath.Join("ccip", "ccip_router"),
	CCIPTokenPool:        filepath.Join("ccip", "ccip_token_pools", "token_pool"),
	LockReleaseTokenPool: filepath.Join("ccip", "ccip_token_pools", "lock_release_token_pool"),
	BurnMintTokenPool:    filepath.Join("ccip", "ccip_token_pools", "burn_mint_token_pool"),
	ManagedTokenPool:     filepath.Join("ccip", "ccip_token_pools", "managed_token_pool"),
	LINKToken:            filepath.Join("ccip", "link_token"),
	ManagedToken:         filepath.Join("ccip", "managed_token"),
	// MCMS
	MCMS: filepath.Join("mcms", "mcms"),
	// Other
	Test: filepath.Join("test"),
}
