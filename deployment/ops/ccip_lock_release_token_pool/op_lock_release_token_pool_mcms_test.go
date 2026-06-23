package lockreleasetokenpoolops

import (
	"testing"

	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_lock_release_token_pool "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_token_pools/lock_release_token_pool"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/ops/mcmstest"
)

func TestLockReleaseTokenPoolDualModeOps_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	typeArgs := []string{mcmstest.CoinTypeArg}
	contract, err := module_lock_release_token_pool.NewLockReleaseTokenPool(mcmstest.PackageID, nil)
	require.NoError(t, err)

	t.Run("apply_chain_updates", func(t *testing.T) {
		t.Parallel()
		input := LockReleaseTokenPoolApplyChainUpdatesInput{
			LockReleasePackageId:         mcmstest.PackageID,
			CoinObjectTypeArg:            mcmstest.CoinTypeArg,
			StateObjectId:                mcmstest.StateObjectID,
			OwnerCap:                     mcmstest.OwnerCapID,
			RemoteChainSelectorsToRemove: []uint64{},
			RemoteChainSelectorsToAdd:    []uint64{},
			RemotePoolAddressesToAdd:     [][]string{},
			RemoteTokenAddressesToAdd:    []string{},
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), LockReleaseTokenPoolApplyChainUpdatesOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := contract.Encoder().ApplyChainUpdates(typeArgs, bind.Object{Id: input.StateObjectId}, bind.Object{Id: input.OwnerCap}, input.RemoteChainSelectorsToRemove, input.RemoteChainSelectorsToAdd, [][][]byte{}, [][]byte{})
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.StateObjectId, typeArgs)
	})

	t.Run("set_chain_rate_limiter", func(t *testing.T) {
		t.Parallel()
		input := LockReleaseTokenPoolSetChainRateLimiterInput{
			LockReleasePackageId: mcmstest.PackageID,
			CoinObjectTypeArg:    mcmstest.CoinTypeArg,
			StateObjectId:        mcmstest.StateObjectID,
			OwnerCap:             mcmstest.OwnerCapID,
			RemoteChainSelectors: []uint64{mcmstest.DestChainSel},
			OutboundIsEnableds:   []bool{true},
			OutboundCapacities:   []uint64{1000},
			OutboundRates:        []uint64{100},
			InboundIsEnableds:    []bool{true},
			InboundCapacities:    []uint64{1000},
			InboundRates:         []uint64{100},
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), LockReleaseTokenPoolSetChainRateLimiterOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := contract.Encoder().SetChainRateLimiterConfigs(typeArgs, bind.Object{Id: input.StateObjectId}, bind.Object{Id: input.OwnerCap}, bind.Object{Id: "0x6"}, input.RemoteChainSelectors, input.OutboundIsEnableds, input.OutboundCapacities, input.OutboundRates, input.InboundIsEnableds, input.InboundCapacities, input.InboundRates)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.StateObjectId, typeArgs)
	})

	t.Run("add_remote_pool", func(t *testing.T) {
		t.Parallel()
		input := LockReleaseTokenPoolAddRemotePoolInput{
			LockReleaseTokenPoolPackageId: mcmstest.PackageID,
			CoinObjectTypeArg:             mcmstest.CoinTypeArg,
			StateObjectId:                 mcmstest.StateObjectID,
			OwnerCap:                      mcmstest.OwnerCapID,
			RemoteChainSelector:           mcmstest.DestChainSel,
			RemotePoolAddress:             "pool-address",
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), LockReleaseTokenPoolAddRemotePoolOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := contract.Encoder().AddRemotePool(typeArgs, bind.Object{Id: input.StateObjectId}, bind.Object{Id: input.OwnerCap}, input.RemoteChainSelector, []byte(input.RemotePoolAddress))
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.StateObjectId, typeArgs)
	})

	t.Run("set_allowlist_enabled", func(t *testing.T) {
		t.Parallel()
		input := LockReleaseTokenPoolSetAllowlistEnabledInput{
			LockReleaseTokenPoolPackageId: mcmstest.PackageID,
			StateObjectId:        mcmstest.StateObjectID,
			OwnerCap:             mcmstest.OwnerCapID,
			CoinObjectTypeArg:    mcmstest.CoinTypeArg,
			Enabled:              true,
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), LockReleaseTokenPoolSetAllowlistEnabledOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := contract.Encoder().SetAllowlistEnabled(typeArgs, bind.Object{Id: input.StateObjectId}, bind.Object{Id: input.OwnerCap}, input.Enabled)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.StateObjectId, typeArgs)
	})

	t.Run("apply_allowlist_updates", func(t *testing.T) {
		t.Parallel()
		input := LockReleaseTokenPoolApplyAllowlistUpdatesInput{
			LockReleaseTokenPoolPackageId: mcmstest.PackageID,
			StateObjectId:        mcmstest.StateObjectID,
			OwnerCap:             mcmstest.OwnerCapID,
			CoinObjectTypeArg:    mcmstest.CoinTypeArg,
			Removes:              []string{},
			Adds:                 []string{mcmstest.Recipient},
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), LockReleaseTokenPoolApplyAllowlistUpdatesOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := contract.Encoder().ApplyAllowlistUpdates(typeArgs, bind.Object{Id: input.StateObjectId}, bind.Object{Id: input.OwnerCap}, input.Removes, input.Adds)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.StateObjectId, typeArgs)
	})

	t.Run("remove_remote_pool", func(t *testing.T) {
		t.Parallel()
		poolAddr := []byte(mcmstest.CoinMetadata)
		input := LockReleaseTokenPoolRemoveRemotePoolInput{
			LockReleaseTokenPoolPackageId: mcmstest.PackageID,
			CoinObjectTypeArg:    mcmstest.CoinTypeArg,
			StateObjectId:        mcmstest.StateObjectID,
			OwnerCap:             mcmstest.OwnerCapID,
			RemoteChainSelector:  mcmstest.DestChainSel,
			RemotePoolAddress:    mcmstest.CoinMetadata,
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), LockReleaseTokenPoolRemoveRemotePoolOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := contract.Encoder().RemoveRemotePool(typeArgs, bind.Object{Id: input.StateObjectId}, bind.Object{Id: input.OwnerCap}, input.RemoteChainSelector, poolAddr)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.StateObjectId, typeArgs)
	})
}
