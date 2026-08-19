package changesets

import (
	"testing"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/stretchr/testify/require"
)

func TestAcceptOwnershipTokenPool_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	cs := AcceptOwnershipTokenPool{}

	// Missing chainSelector.
	err := cs.VerifyPreconditions(cldf.Environment{}, AcceptOwnershipTokenPoolConfig{
		ManagedTokenPoolTokenSymbol: "CCIP BnM",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "chainSelector is required")

	// No token pool selected.
	err = cs.VerifyPreconditions(cldf.Environment{}, AcceptOwnershipTokenPoolConfig{
		ChainSelector: 9762610643973837292,
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "at least one token pool symbol must be selected")

	// Valid: slow MCMS + a managed token pool.
	err = cs.VerifyPreconditions(cldf.Environment{}, AcceptOwnershipTokenPoolConfig{
		ChainSelector:               9762610643973837292,
		ManagedTokenPoolTokenSymbol: "CCIP BnM",
		TypeArg:                     "0xde9a44c43b1e5cf3bee4ae5d6c1aa53f2981513ab3354ebace4fba470f44f92a::ccip_burn_mint_token::CCIP_BURN_MINT_TOKEN",
	})
	require.NoError(t, err)

	// Valid: burn-mint pool instead.
	err = cs.VerifyPreconditions(cldf.Environment{}, AcceptOwnershipTokenPoolConfig{
		ChainSelector:                9762610643973837292,
		BurnMintTokenPoolTokenSymbol: "CCIP BnM",
	})
	require.NoError(t, err)
}
