package lanes

import (
	"github.com/Masterminds/semver/v3"

	chainsel "github.com/smartcontractkit/chain-selectors"

	laneapi "github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
)

func init() {
	laneapi.GetLaneAdapterRegistry().RegisterLaneAdapter(chainsel.FamilySui, semver.MustParse("1.6.0"), &SuiAdapter{})
}
