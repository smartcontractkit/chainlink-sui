package lanes

import (
	"github.com/Masterminds/semver/v3"

	laneapi "github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"
)

// DisableRemoteChain is a no-op on Sui until remote-chain disable is implemented on-chain.
var DisableRemoteChain = cldf_ops.NewSequence(
	"DisableRemoteChain",
	semver.MustParse("1.6.0"),
	"No-op disable remote chain sequence for Sui",
	func(_ cldf_ops.Bundle, _ cldf_chain.BlockChains, _ laneapi.DisableRemoteChainInput) (sequences.OnChainOutput, error) {
		return sequences.OnChainOutput{}, nil
	},
)

func (a *SuiAdapter) DisableRemoteChain() *cldf_ops.Sequence[laneapi.DisableRemoteChainInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return DisableRemoteChain
}
