package lanes

import (
	"fmt"
	"math/big"

	"github.com/Masterminds/semver/v3"

	laneapi "github.com/smartcontractkit/chainlink-ccip/deployment/lanes"
	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_ops "github.com/smartcontractkit/chainlink-deployments-framework/operations"

	ccip_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
	ccip_onramp_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_onramp"
	ccip_router_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip_router"
)

// ConfigureLaneLegAsSource wires FeeQuoter, OnRamp, and Router on Sui for outbound messages
// to a remote destination chain. Op order matches connect_sui_to_evm.
var ConfigureLaneLegAsSource = cldf_ops.NewSequence(
	"ConfigureLaneLegAsSource",
	semver.MustParse("1.6.0"),
	"Configures lane leg as source on CCIP 1.6.0 for Sui",
	func(b cldf_ops.Bundle, chains cldf_chain.BlockChains, input laneapi.UpdateLanesInput) (sequences.OnChainOutput, error) {
		if input.Source == nil || input.Dest == nil {
			return sequences.OnChainOutput{}, fmt.Errorf("ConfigureLaneLegAsSource requires Source and Dest chain definitions")
		}
		if len(input.Dest.Router) == 0 {
			return sequences.OnChainOutput{}, fmt.Errorf("dest Router address required to configure Sui as source for chain %d", input.Dest.Selector)
		}

		env, err := connectChainsEnvironment()
		if err != nil {
			return sequences.OnChainOutput{}, err
		}

		state, err := loadSourceChainState(env, input.Source.Selector)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}

		deps, err := opTxDepsForChain(chains, input.Source.Selector)
		if err != nil {
			return sequences.OnChainOutput{}, err
		}

		destRouter, err := remoteAddressBytesToHex(input.Dest.Router)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("encode dest router for chain %d: %w", input.Dest.Selector, err)
		}

		latestIDs := resolveLatestPackageIDs(input.Source.Selector)

		var out sequences.OnChainOutput

		tokenFeeReport, err := cldf_ops.ExecuteOperation(b, ccip_ops.FeeQuoterApplyTokenTransferFeeConfigUpdatesOp, deps, ccip_ops.FeeQuoterApplyTokenTransferFeeConfigUpdatesInput{
			CCIPPackageId:        state.CCIPAddress,
			LatestPackageId:      latestIDs.CCIP,
			StateObjectId:        state.CCIPObjectRef,
			OwnerCapObjectId:     state.CCIPOwnerCapObjectId,
			DestChainSelector:    input.Dest.Selector,
			AddTokens:            []string{state.LinkTokenCoinMetadataId},
			AddMinFeeUsdCents:    []uint32{DefaultLinkTokenTransferMinFeeUsdCents},
			AddMaxFeeUsdCents:    []uint32{DefaultLinkTokenTransferMaxFeeUsdCents},
			AddDeciBps:           []uint16{DefaultLinkTokenTransferDeciBps},
			AddDestGasOverhead:   []uint32{DefaultLinkTokenTransferDestGasOverhead},
			AddDestBytesOverhead: []uint32{DefaultLinkTokenTransferDestBytesOverhead},
			AddIsEnabled:         []bool{true},
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("apply token transfer fee config on Sui for dest %d: %w", input.Dest.Selector, err)
		}
		if err := appendMCMSBatchOpFromCall(&out, input.Source.Selector, tokenFeeReport.Output.Call, deps); err != nil {
			return sequences.OnChainOutput{}, err
		}

		destCfgInput := TranslateDestChainConfig(input.Dest.FeeQuoterDestChainConfig, input.Dest.Selector)
		destCfgInput.CCIPPackageId = state.CCIPAddress
		destCfgInput.LatestPackageId = latestIDs.CCIP
		destCfgInput.StateObjectId = state.CCIPObjectRef
		destCfgInput.OwnerCapObjectId = state.CCIPOwnerCapObjectId
		destCfgReport, err := cldf_ops.ExecuteOperation(b, ccip_ops.FeeQuoterApplyDestChainConfigUpdatesOp, deps, destCfgInput)
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("apply dest chain config on Sui FeeQuoter for dest %d: %w", input.Dest.Selector, err)
		}
		if err := appendMCMSBatchOpFromCall(&out, input.Source.Selector, destCfgReport.Output.Call, deps); err != nil {
			return sequences.OnChainOutput{}, err
		}

		premiumReport, err := cldf_ops.ExecuteOperation(b, ccip_ops.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesOp, deps, ccip_ops.FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput{
			CCIPPackageId:              state.CCIPAddress,
			LatestPackageId:            latestIDs.CCIP,
			StateObjectId:              state.CCIPObjectRef,
			OwnerCapObjectId:           state.CCIPOwnerCapObjectId,
			Tokens:                     []string{state.LinkTokenCoinMetadataId},
			PremiumMultiplierWeiPerEth: []uint64{DefaultLinkPremiumMultiplierWeiPerEth},
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("apply premium multiplier on Sui FeeQuoter for dest %d: %w", input.Dest.Selector, err)
		}
		if err := appendMCMSBatchOpFromCall(&out, input.Source.Selector, premiumReport.Output.Call, deps); err != nil {
			return sequences.OnChainOutput{}, err
		}

		onRampReport, err := cldf_ops.ExecuteOperation(b, ccip_onramp_ops.ApplyDestChainConfigUpdateOp, deps, ccip_onramp_ops.ApplyDestChainConfigureOnRampInput{
			OnRampPackageId:           state.OnRampAddress,
			LatestPackageId:           latestIDs.OnRamp,
			CCIPObjectRefId:           state.CCIPObjectRef,
			OwnerCapObjectId:          state.OnRampOwnerCapObjectId,
			StateObjectId:             state.OnRampStateObjectId,
			DestChainSelector:         []uint64{input.Dest.Selector},
			DestChainAllowListEnabled: []bool{input.Source.AllowListEnabled},
			DestChainRouters:          []string{destRouter},
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("apply dest chain config on Sui OnRamp for dest %d: %w", input.Dest.Selector, err)
		}
		if err := appendMCMSBatchOpFromCall(&out, input.Source.Selector, onRampReport.Output.Call, deps); err != nil {
			return sequences.OnChainOutput{}, err
		}

		routerReport, err := cldf_ops.ExecuteOperation(b, ccip_router_ops.SetOnRampsOp, deps, ccip_router_ops.SetOnRampsInput{
			RouterPackageId:     state.CCIPRouterAddress,
			LatestPackageId:     latestIDs.Router,
			RouterStateObjectId: state.CCIPRouterStateObjectID,
			OwnerCapObjectId:    state.CCIPRouterOwnerCapObjectId,
			DestChainSelectors:  []uint64{input.Dest.Selector},
			OnRampAddresses:     []string{state.OnRampAddress},
		})
		if err != nil {
			return sequences.OnChainOutput{}, fmt.Errorf("set on-ramps on Sui Router for dest %d: %w", input.Dest.Selector, err)
		}
		if err := appendMCMSBatchOpFromCall(&out, input.Source.Selector, routerReport.Output.Call, deps); err != nil {
			return sequences.OnChainOutput{}, err
		}

		// Seed the FeeQuoter's usd_per_unit_gas_by_dest_chain map for the new
		// dest chain. Without this the very next `get_fee` call for this dest
		// chain aborts with fee_quoter::EUnknownDestChainSelector (code 3) —
		// the DON's price pusher only starts publishing after a dest chain has
		// an initial entry. See connect_chains.go: input.Dest.GasPrice is
		// pre-populated by the framework (falls back to
		// adapter.GetDefaultGasPrice() if the caller didn't supply one).
		if input.Dest.GasPrice != nil {
			pricesReport, err := cldf_ops.ExecuteOperation(b, ccip_ops.FeeQuoterUpdatePricesWithOwnerCapOp, deps, ccip_ops.FeeQuoterUpdatePricesWithOwnerCapInput{
				CCIPPackageId:         state.CCIPAddress,
				LatestPackageId:       latestIDs.CCIP,
				CCIPObjectRef:         state.CCIPObjectRef,
				OwnerCapObjectId:      state.CCIPOwnerCapObjectId,
				GasDestChainSelectors: []uint64{input.Dest.Selector},
				GasUsdPerUnitGas:      []*big.Int{input.Dest.GasPrice},
			})
			if err != nil {
				return sequences.OnChainOutput{}, fmt.Errorf("seed initial gas price on Sui FeeQuoter for dest %d: %w", input.Dest.Selector, err)
			}
			if err := appendMCMSBatchOpFromCall(&out, input.Source.Selector, pricesReport.Output.Call, deps); err != nil {
				return sequences.OnChainOutput{}, err
			}
		}

		return out, nil
	},
)

func (a *SuiAdapter) ConfigureLaneLegAsSource() *cldf_ops.Sequence[laneapi.UpdateLanesInput, sequences.OnChainOutput, cldf_chain.BlockChains] {
	return ConfigureLaneLegAsSource
}
