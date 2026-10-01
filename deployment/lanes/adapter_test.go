package lanes_test

import (
	"encoding/binary"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	laneapi "github.com/smartcontractkit/chainlink-ccip/deployment/lanes"

	module_fee_quoter "github.com/smartcontractkit/chainlink-sui/bindings/generated/ccip/ccip/fee_quoter"
	"github.com/smartcontractkit/chainlink-sui/deployment/lanes"
)

func TestSuiAdapter_GetFeeQuoterDestChainConfig(t *testing.T) {
	t.Parallel()

	cfg := (&lanes.SuiAdapter{}).GetFeeQuoterDestChainConfig()

	require.True(t, cfg.IsEnabled)
	require.Equal(t, uint32(16_000), cfg.MaxDataBytes)
	require.Equal(t, uint32(3_000_000), cfg.MaxPerMsgGasLimit)
	require.Equal(t, uint32(300_000), cfg.DestGasOverhead)
	require.Equal(t, uint8(16), cfg.DestGasPerPayloadByteBase)
	require.Equal(t, uint32(0xc4e05953), cfg.ChainFamilySelector)
	require.Equal(t, uint16(25), cfg.DefaultTokenFeeUSDCents)
	require.Equal(t, uint32(90_000), cfg.DefaultTokenDestGasOverhead)
	require.Equal(t, uint32(200_000), cfg.DefaultTxGasLimit)
	require.Equal(t, uint16(10), cfg.NetworkFeeUSDCents)
	require.NotNil(t, cfg.V1Params)
	require.Equal(t, uint16(10), cfg.V1Params.MaxNumberOfTokensPerMsg)
	require.True(t, cfg.V1Params.EnforceOutOfOrder)
	require.Equal(t, uint64(11e17), cfg.V1Params.GasMultiplierWeiPerEth)
}

func TestSuiAdapter_GetDefaultGasPrice(t *testing.T) {
	t.Parallel()

	require.Equal(t, big.NewInt(15e11), (&lanes.SuiAdapter{}).GetDefaultGasPrice())
}

func TestSuiAdapter_GetChainFamilySelector(t *testing.T) {
	t.Parallel()

	selector := (&lanes.SuiAdapter{}).GetChainFamilySelector()
	require.Equal(t, [4]byte{0xc4, 0xe0, 0x59, 0x53}, selector)
	require.Equal(t, uint32(0xc4e05953), binary.BigEndian.Uint32(selector[:]))
}

func TestTranslateDestChainConfig(t *testing.T) {
	t.Parallel()

	const destSel = uint64(945045181441419236)
	cfg := laneapi.FeeQuoterDestChainConfig{
		IsEnabled:                   true,
		MaxDataBytes:                30_000,
		MaxPerMsgGasLimit:           3_000_000,
		DestGasOverhead:             300_000,
		DestGasPerPayloadByteBase:   16,
		ChainFamilySelector:         0x2812d52c,
		DefaultTokenFeeUSDCents:     25,
		DefaultTokenDestGasOverhead: 90_000,
		DefaultTxGasLimit:           200_000,
		NetworkFeeUSDCents:          10,
		V1Params: &laneapi.FeeQuoterV1Params{
			MaxNumberOfTokensPerMsg:           1,
			DestGasPerPayloadByteHigh:         40,
			DestGasPerPayloadByteThreshold:    3_000,
			DestDataAvailabilityOverheadGas:   100,
			DestGasPerDataAvailabilityByte:    16,
			DestDataAvailabilityMultiplierBps: 1,
			EnforceOutOfOrder:                 true,
			GasMultiplierWeiPerEth:            1e18,
			GasPriceStalenessThreshold:        1_000_000,
		},
	}

	got := lanes.TranslateDestChainConfig(cfg, destSel)

	require.Equal(t, destSel, got.DestChainSelector)
	require.True(t, got.IsEnabled)
	require.Equal(t, uint16(1), got.MaxNumberOfTokensPerMsg)
	require.Equal(t, uint32(30_000), got.MaxDataBytes)
	require.Equal(t, byte(16), got.DestGasPerPayloadByteBase)
	require.Equal(t, []byte{0x28, 0x12, 0xd5, 0x2c}, got.ChainFamilySelector)
	require.Equal(t, uint64(1e18), got.GasMultiplierWeiPerEth)
	require.Equal(t, uint32(10), got.NetworkFeeUsdCents)
}

func TestTranslateDestChainConfigFromMove(t *testing.T) {
	t.Parallel()

	const destSel = uint64(16015286601757825753)
	moveCfg := module_fee_quoter.DestChainConfig{
		IsEnabled:                         true,
		MaxNumberOfTokensPerMsg:           1,
		MaxDataBytes:                      30_000,
		MaxPerMsgGasLimit:                 3_000_000,
		DestGasOverhead:                   300_000,
		DestGasPerPayloadByteBase:         16,
		DestGasPerPayloadByteHigh:         40,
		DestGasPerPayloadByteThreshold:    3_000,
		DestDataAvailabilityOverheadGas:   100,
		DestGasPerDataAvailabilityByte:    16,
		DestDataAvailabilityMultiplierBps: 1,
		ChainFamilySelector:               []byte{0x28, 0x12, 0xd5, 0x2c},
		EnforceOutOfOrder:                 true,
		DefaultTokenFeeUsdCents:           25,
		DefaultTokenDestGasOverhead:       90_000,
		DefaultTxGasLimit:                 200_000,
		GasMultiplierWeiPerEth:            1e18,
		GasPriceStalenessThreshold:        1_000_000,
		NetworkFeeUsdCents:                10,
	}

	laneCfg, err := lanes.TranslateDestChainConfigFromMove(moveCfg)
	require.NoError(t, err)

	// Round-trip the whole 21-field config: Move -> lane config -> op input must
	// preserve every field, since the Move setter replaces the config wholesale.
	opInput := lanes.TranslateDestChainConfig(laneCfg, destSel)
	require.Equal(t, destSel, opInput.DestChainSelector)
	require.Equal(t, moveCfg.IsEnabled, opInput.IsEnabled)
	require.Equal(t, moveCfg.MaxNumberOfTokensPerMsg, opInput.MaxNumberOfTokensPerMsg)
	require.Equal(t, moveCfg.MaxDataBytes, opInput.MaxDataBytes)
	require.Equal(t, moveCfg.MaxPerMsgGasLimit, opInput.MaxPerMsgGasLimit)
	require.Equal(t, moveCfg.DestGasOverhead, opInput.DestGasOverhead)
	require.Equal(t, moveCfg.DestGasPerPayloadByteBase, opInput.DestGasPerPayloadByteBase)
	require.Equal(t, moveCfg.DestGasPerPayloadByteHigh, opInput.DestGasPerPayloadByteHigh)
	require.Equal(t, moveCfg.DestGasPerPayloadByteThreshold, opInput.DestGasPerPayloadByteThreshold)
	require.Equal(t, moveCfg.DestDataAvailabilityOverheadGas, opInput.DestDataAvailabilityOverheadGas)
	require.Equal(t, moveCfg.DestGasPerDataAvailabilityByte, opInput.DestGasPerDataAvailabilityByte)
	require.Equal(t, moveCfg.DestDataAvailabilityMultiplierBps, opInput.DestDataAvailabilityMultiplierBps)
	require.Equal(t, moveCfg.ChainFamilySelector, opInput.ChainFamilySelector)
	require.Equal(t, moveCfg.EnforceOutOfOrder, opInput.EnforceOutOfOrder)
	require.Equal(t, moveCfg.DefaultTokenFeeUsdCents, opInput.DefaultTokenFeeUsdCents)
	require.Equal(t, moveCfg.DefaultTokenDestGasOverhead, opInput.DefaultTokenDestGasOverhead)
	require.Equal(t, moveCfg.DefaultTxGasLimit, opInput.DefaultTxGasLimit)
	require.Equal(t, moveCfg.GasMultiplierWeiPerEth, opInput.GasMultiplierWeiPerEth)
	require.Equal(t, moveCfg.GasPriceStalenessThreshold, opInput.GasPriceStalenessThreshold)
	require.Equal(t, moveCfg.NetworkFeeUsdCents, opInput.NetworkFeeUsdCents)
}

func TestTranslateDestChainConfigFromMove_InvalidChainFamilySelector(t *testing.T) {
	t.Parallel()

	moveCfg := module_fee_quoter.DestChainConfig{ChainFamilySelector: []byte{0x28, 0x12}}
	_, err := lanes.TranslateDestChainConfigFromMove(moveCfg)
	require.ErrorContains(t, err, "chain family selector must be 4 bytes")
}

func TestDefaultLinkTokenTransferFees(t *testing.T) {
	t.Parallel()

	require.Equal(t, uint32(3000), lanes.DefaultLinkTokenTransferMinFeeUsdCents)
	require.Equal(t, uint32(30000), lanes.DefaultLinkTokenTransferMaxFeeUsdCents)
	require.Equal(t, uint16(1000), lanes.DefaultLinkTokenTransferDeciBps)
	require.Equal(t, uint64(900_000_000_000_000_000), lanes.DefaultLinkPremiumMultiplierWeiPerEth)
}
