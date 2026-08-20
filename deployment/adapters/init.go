package adapters

import (
	"github.com/Masterminds/semver/v3"
	chainsel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/deployment/fastcurse"
	tokensapi "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/changesets"
)

func init() {
	curseRegistry := fastcurse.GetCurseRegistry()
	curseRegistry.RegisterNewCurse(fastcurse.CurseRegistryInput{
		CursingFamily:       chainsel.FamilySui,
		CursingVersion:      semver.MustParse("1.6.0"),
		CurseAdapter:        NewCurseAdapter(),
		CurseSubjectAdapter: NewCurseAdapter(),
	})

	mcmsRegistry := changesets.GetRegistry()
	mcmsRegistry.RegisterMCMSReader(chainsel.FamilySui, &MCMSReader{})

	tokenRegistry := tokensapi.GetTokenAdapterRegistry()
	tokenRegistry.RegisterTokenRefResolver(chainsel.FamilySui, &SuiTokenAdapter{})
	tokenRegistry.RegisterTokenAdminRegistryReader(chainsel.FamilySui, &SuiTokenAdminRegistryReader{})
	// Sui CCIP ships on the 1.6.0 release line, so generic token changesets dispatch the Sui
	// adapter under family version 1.6.0, matching the curse adapter above and the Solana
	// sibling. Sui Move pool packages are themselves versioned 1.0.0 and existing pool refs are
	// saved under that version, so ResolveAdapter paths that key off the stored ref version need
	// a 1.0.0 entry too. The version key only selects dispatch; the same stateless adapter serves both.
	tokenRegistry.RegisterTokenAdapter(chainsel.FamilySui, semver.MustParse("1.6.0"), &SuiTokenAdapter{})
	tokenRegistry.RegisterTokenAdapter(chainsel.FamilySui, semver.MustParse("1.0.0"), &SuiTokenAdapter{})
}
