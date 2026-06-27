package offrampops

import (
	"testing"

	"github.com/stretchr/testify/require"

	cld_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	"github.com/smartcontractkit/chainlink-sui/bindings/bind"
	module_offramp "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_offramp/offramp"
	module_ownable "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip_onramp/ownable"
	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/ops/mcmstest"
)

func TestOffRampDualModeOps_ProposalDataMatchesBindingEncoder(t *testing.T) {
	t.Parallel()

	offramp, err := module_offramp.NewOfframp(mcmstest.PackageID, nil)
	require.NoError(t, err)
	ownable, err := module_ownable.NewOwnable(mcmstest.PackageID, nil)
	require.NoError(t, err)

	t.Run("set_ocr3_config", func(t *testing.T) {
		t.Parallel()
		input := SetOCR3ConfigInput{
			OffRampPackageId:               mcmstest.PackageID,
			CCIPObjectRefId:                mcmstest.StateObjectID,
			OffRampStateId:                 mcmstest.CoinMetadata,
			OwnerCapObjectId:               mcmstest.OwnerCapID,
			ConfigDigest:                   make([]byte, 32),
			OCRPluginType:                  0,
			BigF:                           1,
			IsSignatureVerificationEnabled: false,
			Signers:                        [][]byte{},
			Transmitters:                   []string{},
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), SetOCR3ConfigOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := offramp.Encoder().SetOcr3Config(bind.Object{Id: input.CCIPObjectRefId}, bind.Object{Id: input.OffRampStateId}, bind.Object{Id: input.OwnerCapObjectId}, input.ConfigDigest, input.OCRPluginType, input.BigF, input.IsSignatureVerificationEnabled, input.Signers, input.Transmitters)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.OffRampStateId, nil)
	})

	t.Run("apply_source_chain_config_updates", func(t *testing.T) {
		t.Parallel()
		input := ApplySourceChainConfigUpdateInput{
			CCIPObjectRef:                         mcmstest.StateObjectID,
			OffRampPackageId:                      mcmstest.PackageID,
			OffRampStateId:                        mcmstest.CoinMetadata,
			OwnerCapObjectId:                      mcmstest.OwnerCapID,
			SourceChainsSelectors:                 []uint64{mcmstest.DestChainSel},
			SourceChainsIsEnabled:                 []bool{true},
			SourceChainsIsRMNVerificationDisabled: []bool{false},
			SourceChainsOnRamp:                    [][]byte{{0x01}},
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), ApplySourceChainConfigUpdatesOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := offramp.Encoder().ApplySourceChainConfigUpdates(bind.Object{Id: input.CCIPObjectRef}, bind.Object{Id: input.OffRampStateId}, bind.Object{Id: input.OwnerCapObjectId}, input.SourceChainsSelectors, input.SourceChainsIsEnabled, input.SourceChainsIsRMNVerificationDisabled, input.SourceChainsOnRamp)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.OffRampStateId, nil)
	})

	t.Run("apply_source_chain_config_updates_with_latest_package_id", func(t *testing.T) {
		t.Parallel()
		const latestPackageID = "0xdddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd"
		input := ApplySourceChainConfigUpdateInput{
			CCIPObjectRef:                         mcmstest.StateObjectID,
			OffRampPackageId:                      mcmstest.PackageID,
			LatestPackageId:                       latestPackageID,
			OffRampStateId:                        mcmstest.CoinMetadata,
			OwnerCapObjectId:                      mcmstest.OwnerCapID,
			SourceChainsSelectors:                 []uint64{mcmstest.DestChainSel},
			SourceChainsIsEnabled:                 []bool{true},
			SourceChainsIsRMNVerificationDisabled: []bool{false},
			SourceChainsOnRamp:                    [][]byte{{0x01}},
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), ApplySourceChainConfigUpdatesOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		require.Equal(t, mcmstest.PackageID, report.Output.Call.PackageID)
		require.Equal(t, latestPackageID, report.Output.Call.LatestPackageID)

		latestOfframp, err := module_offramp.NewOfframp(latestPackageID, nil)
		require.NoError(t, err)
		encoded, err := latestOfframp.Encoder().ApplySourceChainConfigUpdates(bind.Object{Id: input.CCIPObjectRef}, bind.Object{Id: input.OffRampStateId}, bind.Object{Id: input.OwnerCapObjectId}, input.SourceChainsSelectors, input.SourceChainsIsEnabled, input.SourceChainsIsRMNVerificationDisabled, input.SourceChainsOnRamp)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.OffRampStateId, nil)
	})

	t.Run("add_package_id", func(t *testing.T) {
		t.Parallel()
		input := AddPackageIdOffRampInput{
			PackageId:        mcmstest.PackageID,
			StateObjectId:    mcmstest.CoinMetadata,
			OwnerCapObjectId: mcmstest.OwnerCapID,
			NewPackageId:     mcmstest.Recipient,
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), AddPackageIdOffRampOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := offramp.Encoder().AddPackageId(bind.Object{Id: input.StateObjectId}, bind.Object{Id: input.OwnerCapObjectId}, input.NewPackageId)
		require.NoError(t, err)
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.StateObjectId, nil)
	})

	t.Run("accept_ownership", func(t *testing.T) {
		t.Parallel()
		input := AcceptOwnershipOffRampInput{
			OffRampPackageId:     mcmstest.PackageID,
			OffRampRefObjectId:   mcmstest.StateObjectID,
			OffRampStateObjectId: mcmstest.CoinMetadata,
		}
		report, err := cld_ops.ExecuteOperation(mcmstest.Bundle(t), AcceptOwnershipOffRampOp, sui_ops.OpTxDeps{}, input)
		require.NoError(t, err)
		encoded, err := ownable.Encoder().AcceptOwnershipWithArgs(bind.Object{Id: input.OffRampStateObjectId})
		require.NoError(t, err)
		encoded.Module.ModuleName = "offramp"
		mcmstest.AssertProposalDataMatches(t, report.Output.Call.Data, encoded, input.OffRampRefObjectId, nil)
	})
}
