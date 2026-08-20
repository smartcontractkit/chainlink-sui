package adapters

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	tokensapi "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	"github.com/stretchr/testify/require"
)

func TestInit_RegistersSuiCurseAndSubjectAdapters(t *testing.T) {
	t.Parallel()

	reg := fastcurse.GetCurseRegistry()

	// Pinned chainlink-ccip registers CurseSubjectAdapter under CursingFamily.
	_, ok := reg.GetCurseSubjectAdapter(chainsel.FamilySui)
	require.True(t, ok, "sui subject adapter must be registered under FamilySui")

	_, ok = reg.GetCurseAdapter(chainsel.FamilySui, semver.MustParse("1.6.0"))
	require.True(t, ok, "sui curse adapter must be registered for v1.6.0")
}

func TestInit_RegistersSuiTokenAdminRegistryReader(t *testing.T) {
	t.Parallel()

	reg := tokensapi.GetTokenAdapterRegistry()

	// Required by the generic token changesets whenever autoMigrateRemoteChains is enabled,
	// to read the pool registered for a token from the on-chain TokenAdminRegistry.
	_, ok := reg.GetTokenAdminRegistryReader(chainsel.FamilySui)
	require.True(t, ok, "sui token admin registry reader must be registered")
}
