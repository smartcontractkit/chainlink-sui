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

	// Version must match the version Sui pool refs are saved under (deployment.Version1_0_0).
	v := semver.MustParse("1.0.0")
	tokenRegistry := tokensapi.GetTokenAdapterRegistry()
	tokenRegistry.RegisterTokenRefResolver(chainsel.FamilySui, &SuiTokenAdapter{})
	tokenRegistry.RegisterTokenAdapter(chainsel.FamilySui, v, &SuiTokenAdapter{})
	tokenRegistry.RegisterTokenAdminRegistryReader(chainsel.FamilySui, &SuiTokenAdminRegistryReader{})
}
