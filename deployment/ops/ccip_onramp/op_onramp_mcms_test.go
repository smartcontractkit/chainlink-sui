package onrampops

import (
	"testing"

	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_onramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/onramp"
	module_ownable "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/ownable"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/ops/mcmstest"
)

func TestOnRampDualModeOps_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	onramp, err := module_onramp.NewOnramp(mcmstest.PackageID, nil)
	require.NoError(t, err)
	ownable, err := module_ownable.NewOwnable(mcmstest.PackageID, nil)
	require.NoError(t, err)

	t.Run("apply_dest_chain_config_updates", func(t *testing.T) {
		t.Parallel()
		input := ApplyDestChainConfigureOnRampInput{
			OnRampPackageId:           mcmstest.PackageID,
			CCIPObjectRefId:           mcmstest.StateObjectID,
			OwnerCapObjectId:          mcmstest.OwnerCapID,
			StateObjectId:             mcmstest.CoinMetadata,
			DestChainSelector:         []uint64{mcmstest.DestChainSel},
			DestChainAllowListEnabled: []bool{true},
			DestChainRouters:          []string{mcmstest.Recipient},
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), ApplyDestChainConfigUpdateOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := onramp.Encoder().ApplyDestChainConfigUpdates(bind.Object{Id: input.CCIPObjectRefId}, bind.Object{Id: input.StateObjectId}, bind.Object{Id: input.OwnerCapObjectId}, input.DestChainSelector, input.DestChainAllowListEnabled, input.DestChainRouters)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.StateObjectId, nil)
	})

	t.Run("apply_allowlist_updates", func(t *testing.T) {
		t.Parallel()
		input := ApplyAllowListUpdatesInput{
			OnRampPackageId:               mcmstest.PackageID,
			CCIPObjectRefId:               mcmstest.StateObjectID,
			OwnerCapObjectId:              mcmstest.OwnerCapID,
			StateObjectId:                 mcmstest.CoinMetadata,
			DestChainSelector:             []uint64{mcmstest.DestChainSel},
			DestChainAllowListEnabled:     []bool{true},
			DestChainAddAllowedSenders:    [][]string{{mcmstest.Recipient}},
			DestChainRemoveAllowedSenders: [][]string{},
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), ApplyAllowListUpdateOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := onramp.Encoder().ApplyAllowlistUpdates(bind.Object{Id: input.CCIPObjectRefId}, bind.Object{Id: input.StateObjectId}, bind.Object{Id: input.OwnerCapObjectId}, input.DestChainSelector, input.DestChainAllowListEnabled, input.DestChainAddAllowedSenders, input.DestChainRemoveAllowedSenders)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.StateObjectId, nil)
	})

	t.Run("set_dynamic_config", func(t *testing.T) {
		t.Parallel()
		input := SetDynamicConfigInput{
			OnRampPackageId:  mcmstest.PackageID,
			CCIPObjectRefId:  mcmstest.StateObjectID,
			StateObjectId:    mcmstest.CoinMetadata,
			OwnerCapObjectId: mcmstest.OwnerCapID,
			FeeAggregator:    mcmstest.Recipient,
			AllowListAdmin:   mcmstest.Recipient,
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), SetDynamicConfigOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := onramp.Encoder().SetDynamicConfig(bind.Object{Id: input.CCIPObjectRefId}, bind.Object{Id: input.StateObjectId}, bind.Object{Id: input.OwnerCapObjectId}, input.FeeAggregator, input.AllowListAdmin)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.StateObjectId, nil)
	})

	t.Run("withdraw_fee_tokens", func(t *testing.T) {
		t.Parallel()
		typeArgs := []string{mcmstest.CoinTypeArg}
		input := WithdrawFeeTokensInput{
			OnRampPackageId:    mcmstest.PackageID,
			CCIPObjectRefId:    mcmstest.StateObjectID,
			StateObjectId:      mcmstest.CoinMetadata,
			OwnerCapObjectId:   mcmstest.OwnerCapID,
			FeeTokenMetadataId: mcmstest.RegistryID,
			TypeArg:            mcmstest.CoinTypeArg,
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), WithdrawFeeTokensOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := onramp.Encoder().WithdrawFeeTokens(typeArgs, bind.Object{Id: input.CCIPObjectRefId}, bind.Object{Id: input.StateObjectId}, bind.Object{Id: input.OwnerCapObjectId}, bind.Object{Id: input.FeeTokenMetadataId})
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.StateObjectId, typeArgs)
	})

	t.Run("add_package_id", func(t *testing.T) {
		t.Parallel()
		input := AddPackageIdInput{
			PackageId:        mcmstest.PackageID,
			StateObjectId:    mcmstest.CoinMetadata,
			OwnerCapObjectId: mcmstest.OwnerCapID,
			NewPackageId:     mcmstest.Recipient,
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), AddPackageIdOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := onramp.Encoder().AddPackageId(bind.Object{Id: input.StateObjectId}, bind.Object{Id: input.OwnerCapObjectId}, input.NewPackageId)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.StateObjectId, nil)
	})

	t.Run("remove_package_id", func(t *testing.T) {
		t.Parallel()
		input := RemovePackageIdOnRampInput{
			OnRampPackageId:  mcmstest.PackageID,
			StateObjectId:    mcmstest.CoinMetadata,
			OwnerCapObjectId: mcmstest.OwnerCapID,
			PackageId:        mcmstest.Recipient,
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), RemovePackageIdOnRampOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := onramp.Encoder().RemovePackageId(bind.Object{Id: input.StateObjectId}, bind.Object{Id: input.OwnerCapObjectId}, input.PackageId)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.StateObjectId, nil)
	})

	t.Run("accept_ownership", func(t *testing.T) {
		t.Parallel()
		input := AcceptOwnershipOnRampInput{
			OnRampPackageId: mcmstest.PackageID,
			CCIPObjectRefId: mcmstest.StateObjectID,
			StateObjectId:   mcmstest.CoinMetadata,
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), AcceptOwnershipOnRampOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := ownable.Encoder().AcceptOwnershipWithArgs(bind.Object{Id: input.StateObjectId})
		require.NoError(t, err)
		encoded.Module.ModuleName = "onramp"
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.CCIPObjectRefId, nil)
	})
}
