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
	return ccip_ops.FeeQuoterApplyDestChainConfigUpdatesInput{
		DestChainSelector:                 destChainSelector,
		IsEnabled:                         cfg.IsEnabled,
		MaxNumberOfTokensPerMsg:           cfg.MaxNumberOfTokensPerMsg,
		MaxDataBytes:                      cfg.MaxDataBytes,
		MaxPerMsgGasLimit:                 cfg.MaxPerMsgGasLimit,
		DestGasOverhead:                   cfg.DestGasOverhead,
		DestGasPerPayloadByteBase:         cfg.DestGasPerPayloadByteBase,
		DestGasPerPayloadByteHigh:         cfg.DestGasPerPayloadByteHigh,
		DestGasPerPayloadByteThreshold:    cfg.DestGasPerPayloadByteThreshold,
		DestDataAvailabilityOverheadGas:   cfg.DestDataAvailabilityOverheadGas,
		DestGasPerDataAvailabilityByte:    cfg.DestGasPerDataAvailabilityByte,
		DestDataAvailabilityMultiplierBps: cfg.DestDataAvailabilityMultiplierBps,
		ChainFamilySelector:               binary.BigEndian.AppendUint32(nil, cfg.ChainFamilySelector),
		EnforceOutOfOrder:                 cfg.EnforceOutOfOrder,
		DefaultTokenFeeUsdCents:           cfg.DefaultTokenFeeUSDCents,
		DefaultTokenDestGasOverhead:       cfg.DefaultTokenDestGasOverhead,
		DefaultTxGasLimit:                 cfg.DefaultTxGasLimit,
		GasMultiplierWeiPerEth:            cfg.GasMultiplierWeiPerEth,
		GasPriceStalenessThreshold:        cfg.GasPriceStalenessThreshold,
		NetworkFeeUsdCents:                cfg.NetworkFeeUSDCents,
	}
}
