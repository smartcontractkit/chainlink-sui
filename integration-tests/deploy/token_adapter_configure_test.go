//go:build integration

package deploy

import (
	"encoding/hex"
	"fmt"
	"strings"

	fdatastore "github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	tokensapi "github.com/smartcontractkit/chainlink-ccip/deployment/tokens"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_burn_mint_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/burn_mint_token_pool"
	"github.com/smartcontractkit/chainlink-sui/deployment"
	"github.com/smartcontractkit/chainlink-sui/deployment/adapters"
)

// TestReconfigureTokenPoolForTransfers re-runs the adapter's ConfigureTokenForTransfers
// sequence against the LINK burn-mint pool already configured for the EVM remote in the
// deploy phases: an identical re-run must skip the chain-config call, and a new remote
// pool for the already-supported chain must be appended via add_remote_pool instead of
// aborting with ERemoteChainAlreadyExists.
func (s *DeployTestSuite) TestReconfigureTokenPoolForTransfers() {
	s.T().Log("Phase 12: Reconfigure token pool for transfers")

	coinType := fmt.Sprintf("%s::link::LINK", s.linkTokenPackageID)
	linkQualifier := deployment.TokenQualifier("LINK")

	var poolPkgID, stateObjID string
	for _, ref := range s.suiDatastoreRefs(SuiChainSelector) {
		if ref.Qualifier != linkQualifier {
			continue
		}
		switch ref.Type {
		case fdatastore.ContractType(deployment.SuiBnMTokenPoolType):
			poolPkgID = ref.Address
		case fdatastore.ContractType(deployment.SuiBnMTokenPoolStateType):
			stateObjID = ref.Address
		}
	}
	s.Require().NotEmpty(poolPkgID, "LINK burn-mint pool package not found")
	s.Require().NotEmpty(stateObjID, "LINK burn-mint pool state object not found")

	// The deploy phases registered the EVM remote chain with the raw token and pool
	// addresses, so the desired config must use the same encoding to be a match.
	remoteToken, err := hex.DecodeString(EVMTokenAddress)
	s.Require().NoError(err)
	remotePool, err := hex.DecodeString(strings.TrimPrefix(EVMPoolAddress, "0x"))
	s.Require().NoError(err)
	secondPool := append([]byte{}, remotePool...)
	secondPool[len(secondPool)-1] ^= 0x01

	baseInput := tokensapi.ConfigureTokenForTransfersInput{
		ChainSelector:    SuiChainSelector,
		TokenPoolAddress: poolPkgID,
		TokenRef: fdatastore.AddressRef{
			ChainSelector: SuiChainSelector,
			Address:       coinType,
			Qualifier:     linkQualifier,
		},
		PoolType:          string(deployment.SuiBnMTokenPoolType),
		ExistingDataStore: s.env.DataStore,
		RemoteChains: map[uint64]tokensapi.RemoteChainConfig[[]byte, string]{
			EVMChainSelector: {
				RemoteToken:    remoteToken,
				RemotePool:     remotePool,
				RemoteDecimals: 18,
			},
		},
	}
	seq := (&adapters.SuiTokenAdapter{}).ConfigureTokenForTransfersSequence()

	// Identical re-run: the remote chain is already configured with the same token and
	// pool, so no chain-config call is emitted and nothing is proposed.
	out, err := cldf_ops.ExecuteSequence(s.env.OperationsBundle, seq, s.env.BlockChains, baseInput)
	s.Require().NoError(err, "identical re-run must be a no-op")
	s.Require().Empty(out.Output.BatchOps)
	s.assertLinkPoolRemotePools(poolPkgID, stateObjID, coinType, [][]byte{remotePool})

	// A second remote pool for the already-supported chain must be appended via
	// add_remote_pool rather than aborting with ERemoteChainAlreadyExists.
	input := baseInput
	input.RemoteChains = map[uint64]tokensapi.RemoteChainConfig[[]byte, string]{
		EVMChainSelector: {
			RemoteToken:    remoteToken,
			RemotePool:     secondPool,
			RemoteDecimals: 18,
		},
	}
	out, err = cldf_ops.ExecuteSequence(s.env.OperationsBundle, seq, s.env.BlockChains, input)
	s.Require().NoError(err, "second remote pool must be appended")
	s.Require().Empty(out.Output.BatchOps)
	s.assertLinkPoolRemotePools(poolPkgID, stateObjID, coinType, [][]byte{remotePool, secondPool})
}

// assertLinkPoolRemotePools asserts the LINK burn-mint pool's stored remote pools for the
// EVM remote chain.
func (s *DeployTestSuite) assertLinkPoolRemotePools(poolPkgID, stateObjID, coinType string, want [][]byte) {
	s.T().Helper()
	contract, err := module_burn_mint_token_pool.NewBurnMintTokenPool(poolPkgID, s.client)
	s.Require().NoError(err)
	got, err := contract.DevInspect().GetRemotePools(s.T().Context(), &bind.CallOpts{Signer: s.signer},
		[]string{coinType}, bind.Object{Id: stateObjID}, EVMChainSelector)
	s.Require().NoError(err)
	s.Require().Equal(want, got)
}
