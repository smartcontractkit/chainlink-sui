package lanes

import (
	"encoding/binary"
	"fmt"

	laneapi "github.com/smartcontractkit/chainlink-ccip/deployment/lanes"

	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	ccip_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops/ccip"
)

// TranslateDestChainConfig maps a product-level FeeQuoterDestChainConfig to the Sui
// FeeQuoter op input for configuring a remote destination on the Sui chain.
// CCIP package/state object IDs are filled by the caller from address book state.
func TranslateDestChainConfig(
	cfg laneapi.FeeQuoterDestChainConfig,
	destChainSelector uint64,
) ccip_ops.FeeQuoterApplyDestChainConfigUpdatesInput {
	var v1 laneapi.FeeQuoterV1Params
	if cfg.V1Params != nil {
		v1 = *cfg.V1Params
	}
	return ccip_ops.FeeQuoterApplyDestChainConfigUpdatesInput{
		DestChainSelector:                 destChainSelector,
		IsEnabled:                         cfg.IsEnabled,
		MaxNumberOfTokensPerMsg:           v1.MaxNumberOfTokensPerMsg,
		MaxDataBytes:                      cfg.MaxDataBytes,
		MaxPerMsgGasLimit:                 cfg.MaxPerMsgGasLimit,
		DestGasOverhead:                   cfg.DestGasOverhead,
		DestGasPerPayloadByteBase:         cfg.DestGasPerPayloadByteBase,
		DestGasPerPayloadByteHigh:         v1.DestGasPerPayloadByteHigh,
		DestGasPerPayloadByteThreshold:    v1.DestGasPerPayloadByteThreshold,
		DestDataAvailabilityOverheadGas:   v1.DestDataAvailabilityOverheadGas,
		DestGasPerDataAvailabilityByte:    v1.DestGasPerDataAvailabilityByte,
		DestDataAvailabilityMultiplierBps: v1.DestDataAvailabilityMultiplierBps,
		ChainFamilySelector:               binary.BigEndian.AppendUint32(nil, cfg.ChainFamilySelector),
		EnforceOutOfOrder:                 v1.EnforceOutOfOrder,
		DefaultTokenFeeUsdCents:           cfg.DefaultTokenFeeUSDCents,
		DefaultTokenDestGasOverhead:       cfg.DefaultTokenDestGasOverhead,
		DefaultTxGasLimit:                 cfg.DefaultTxGasLimit,
		GasMultiplierWeiPerEth:            v1.GasMultiplierWeiPerEth,
		GasPriceStalenessThreshold:        v1.GasPriceStalenessThreshold,
		NetworkFeeUsdCents:                uint32(cfg.NetworkFeeUSDCents),
	}
}

// TranslateDestChainConfigFromMove maps an on-chain (Move) FeeQuoter DestChainConfig onto
// the chain-agnostic lane config. It is the inverse of TranslateDestChainConfig and is
// used by read-modify-write flows (e.g. surgical dest-chain-config updates) that must
// round-trip the full 21-field config because the Move setter replaces it wholesale.
// V1Params is always populated (nil-tolerant on the encode side); the on-chain config has
// no OverrideExistingConfig/V2Params counterparts, so those stay unset.
func TranslateDestChainConfigFromMove(
	cfg module_fee_quoter.DestChainConfig,
) (laneapi.FeeQuoterDestChainConfig, error) {
	if len(cfg.ChainFamilySelector) != 4 {
		return laneapi.FeeQuoterDestChainConfig{}, fmt.Errorf(
			"chain family selector must be 4 bytes, got %d", len(cfg.ChainFamilySelector))
	}
	// NetworkFeeUsdCents is u32 on-chain but uint16 in the chain-agnostic config; the
	// narrowing is safe for all realistic fee values (max ~655 USD per message).
	return laneapi.FeeQuoterDestChainConfig{
		IsEnabled:                   cfg.IsEnabled,
		MaxDataBytes:                cfg.MaxDataBytes,
		MaxPerMsgGasLimit:           cfg.MaxPerMsgGasLimit,
		DestGasOverhead:             cfg.DestGasOverhead,
		DestGasPerPayloadByteBase:   cfg.DestGasPerPayloadByteBase,
		ChainFamilySelector:         binary.BigEndian.Uint32(cfg.ChainFamilySelector),
		DefaultTokenFeeUSDCents:     cfg.DefaultTokenFeeUsdCents,
		DefaultTokenDestGasOverhead: cfg.DefaultTokenDestGasOverhead,
		DefaultTxGasLimit:           cfg.DefaultTxGasLimit,
		NetworkFeeUSDCents:          uint16(cfg.NetworkFeeUsdCents),
		V1Params: &laneapi.FeeQuoterV1Params{
			MaxNumberOfTokensPerMsg:           cfg.MaxNumberOfTokensPerMsg,
			DestGasPerPayloadByteHigh:         cfg.DestGasPerPayloadByteHigh,
			DestGasPerPayloadByteThreshold:    cfg.DestGasPerPayloadByteThreshold,
			DestDataAvailabilityOverheadGas:   cfg.DestDataAvailabilityOverheadGas,
			DestGasPerDataAvailabilityByte:    cfg.DestGasPerDataAvailabilityByte,
			DestDataAvailabilityMultiplierBps: cfg.DestDataAvailabilityMultiplierBps,
			EnforceOutOfOrder:                 cfg.EnforceOutOfOrder,
			GasMultiplierWeiPerEth:            cfg.GasMultiplierWeiPerEth,
			GasPriceStalenessThreshold:        cfg.GasPriceStalenessThreshold,
		},
	}, nil
}
