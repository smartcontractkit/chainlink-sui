package lanes

import (
	"encoding/binary"

	laneapi "github.com/smartcontractkit/chainlink-ccip/deployment/lanes"

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
