package ccipops

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFeeQuoterInputValidation(t *testing.T) {
	t.Run("InitFeeQuoterInput", func(t *testing.T) {
		input := InitFeeQuoterInput{
			CCIPPackageId:                 "0x123",
			StateObjectId:                 "0x456",
			OwnerCapObjectId:              "0x789",
			MaxFeeJuelsPerMsg:             "100000000",
			LinkTokenCoinMetadataObjectId: "0xabc",
			TokenPriceStalenessThreshold:  60,
			FeeTokens:                     []string{"0xdef"},
		}

		require.NotEmpty(t, input.CCIPPackageId)
		require.NotEmpty(t, input.StateObjectId)
		require.NotEmpty(t, input.OwnerCapObjectId)
		require.NotEmpty(t, input.MaxFeeJuelsPerMsg)
		require.NotEmpty(t, input.LinkTokenCoinMetadataObjectId)
		require.Greater(t, input.TokenPriceStalenessThreshold, uint64(0))
		require.NotEmpty(t, input.FeeTokens)
	})

	t.Run("FeeQuoterApplyFeeTokenUpdatesInput", func(t *testing.T) {
		input := FeeQuoterApplyFeeTokenUpdatesInput{
			CCIPPackageId:     "0x123",
			StateObjectId:     "0x456",
			OwnerCapObjectId:  "0x789",
			FeeTokensToRemove: []string{"0xdef"},
			FeeTokensToAdd:    []string{"0xabc"},
		}

		require.NotEmpty(t, input.CCIPPackageId)
		require.NotEmpty(t, input.StateObjectId)
		require.NotEmpty(t, input.OwnerCapObjectId)
		require.NotNil(t, input.FeeTokensToRemove)
		require.NotNil(t, input.FeeTokensToAdd)
	})

	t.Run("FeeQuoterApplyTokenTransferFeeConfigUpdatesInput", func(t *testing.T) {
		input := FeeQuoterApplyTokenTransferFeeConfigUpdatesInput{
			CCIPPackageId:        "0x123",
			StateObjectId:        "0x456",
			OwnerCapObjectId:     "0x789",
			DestChainSelector:    1,
			AddTokens:            []string{"0xdef"},
			AddMinFeeUsdCents:    []uint32{1000},
			AddMaxFeeUsdCents:    []uint32{10000},
			AddDeciBps:           []uint16{500},
			AddDestGasOverhead:   []uint32{1000000},
			AddDestBytesOverhead: []uint32{1000},
			AddIsEnabled:         []bool{true},
			RemoveTokens:         []string{},
		}

		require.NotEmpty(t, input.CCIPPackageId)
		require.NotEmpty(t, input.StateObjectId)
		require.NotEmpty(t, input.OwnerCapObjectId)
		require.Greater(t, input.DestChainSelector, uint64(0))
		require.Len(t, input.AddTokens, 1)
		require.Len(t, input.AddMinFeeUsdCents, 1)
		require.Len(t, input.AddMaxFeeUsdCents, 1)
		require.Len(t, input.AddDeciBps, 1)
		require.Len(t, input.AddDestGasOverhead, 1)
		require.Len(t, input.AddDestBytesOverhead, 1)
		require.Len(t, input.AddIsEnabled, 1)
		require.NotNil(t, input.RemoveTokens)
	})

	t.Run("FeeQuoterApplyDestChainConfigUpdatesInput", func(t *testing.T) {
		input := FeeQuoterApplyDestChainConfigUpdatesInput{
			CCIPPackageId:                     "0x123",
			StateObjectId:                     "0x456",
			OwnerCapObjectId:                  "0x789",
			DestChainSelector:                 1,
			IsEnabled:                         true,
			MaxNumberOfTokensPerMsg:           2,
			MaxDataBytes:                      2000,
			MaxPerMsgGasLimit:                 5000000,
			DestGasOverhead:                   1000000,
			DestGasPerPayloadByteBase:         byte(2),
			DestGasPerPayloadByteHigh:         byte(5),
			DestGasPerPayloadByteThreshold:    uint16(10),
			DestDataAvailabilityOverheadGas:   300000,
			DestGasPerDataAvailabilityByte:    4,
			DestDataAvailabilityMultiplierBps: 1,
			ChainFamilySelector:               []byte{0x28, 0x12, 0xd5, 0x2c},
			EnforceOutOfOrder:                 false,
			DefaultTokenFeeUsdCents:           3,
			DefaultTokenDestGasOverhead:       100000,
			DefaultTxGasLimit:                 500000,
			GasMultiplierWeiPerEth:            100,
			GasPriceStalenessThreshold:        300,
			NetworkFeeUsdCents:                5,
		}

		require.NotEmpty(t, input.CCIPPackageId)
		require.NotEmpty(t, input.StateObjectId)
		require.NotEmpty(t, input.OwnerCapObjectId)
		require.Greater(t, input.DestChainSelector, uint64(0))
		require.Greater(t, input.MaxNumberOfTokensPerMsg, uint16(0))
		require.Greater(t, input.MaxDataBytes, uint32(0))
		require.Greater(t, input.MaxPerMsgGasLimit, uint32(0))
		require.Greater(t, input.DestGasOverhead, uint32(0))
		require.Greater(t, input.DestGasPerPayloadByteThreshold, uint16(0))
		require.Greater(t, input.DestDataAvailabilityOverheadGas, uint32(0))
		require.Greater(t, input.DestGasPerDataAvailabilityByte, uint16(0))
		require.Greater(t, input.DestDataAvailabilityMultiplierBps, uint16(0))
		require.NotEmpty(t, input.ChainFamilySelector)
		require.Greater(t, input.DefaultTokenFeeUsdCents, uint16(0))
		require.Greater(t, input.DefaultTokenDestGasOverhead, uint32(0))
		require.Greater(t, input.DefaultTxGasLimit, uint32(0))
		require.Greater(t, input.GasMultiplierWeiPerEth, uint64(0))
		require.Greater(t, input.GasPriceStalenessThreshold, uint32(0))
		require.Greater(t, input.NetworkFeeUsdCents, uint32(0))
	})

	t.Run("FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput", func(t *testing.T) {
		input := FeeQuoterApplyPremiumMultiplierWeiPerEthUpdatesInput{
			CCIPPackageId:              "0x123",
			StateObjectId:              "0x456",
			OwnerCapObjectId:           "0x789",
			Tokens:                     []string{"0xdef"},
			PremiumMultiplierWeiPerEth: []uint64{1100000000000000000}, // 1.1e18
		}

		require.NotEmpty(t, input.CCIPPackageId)
		require.NotEmpty(t, input.StateObjectId)
		require.NotEmpty(t, input.OwnerCapObjectId)
		require.Len(t, input.Tokens, 1)
		require.Len(t, input.PremiumMultiplierWeiPerEth, 1)
		require.Equal(t, input.Tokens, []string{"0xdef"})
		require.Equal(t, input.PremiumMultiplierWeiPerEth, []uint64{1100000000000000000})
	})

	t.Run("FeeQuoterUpdateTokenPricesInput", func(t *testing.T) {
		input := FeeQuoterUpdateTokenPricesInput{
			CCIPPackageId:         "0x123",
			CCIPObjectRef:         "0x456",
			FeeQuoterCapId:        "0x789",
			SourceTokens:          []string{"0xdef"},
			SourceUsdPerToken:     []*big.Int{big.NewInt(1000000000000000000)},
			GasDestChainSelectors: []uint64{1},
			GasUsdPerUnitGas:      []*big.Int{big.NewInt(2000000000000000)},
		}

		require.NotEmpty(t, input.CCIPPackageId)
		require.NotEmpty(t, input.CCIPObjectRef)
		require.NotEmpty(t, input.FeeQuoterCapId)
		require.Len(t, input.SourceTokens, 1)
		require.Len(t, input.SourceUsdPerToken, 1)
		require.Len(t, input.GasDestChainSelectors, 1)
		require.Len(t, input.GasUsdPerUnitGas, 1)
		require.Equal(t, input.SourceTokens, []string{"0xdef"})
		require.Equal(t, input.SourceUsdPerToken[0].String(), "1000000000000000000")
		require.Equal(t, input.GasDestChainSelectors, []uint64{1})
		require.Equal(t, input.GasUsdPerUnitGas[0].String(), "2000000000000000")
	})

	t.Run("NewFeeQuoterCapInput", func(t *testing.T) {
		input := NewFeeQuoterCapInput{
			CCIPPackageId:    "0x123",
			CCIPObjectRef:    "0xabc",
			OwnerCapObjectId: "0x456",
		}

		require.NotEmpty(t, input.CCIPPackageId)
		require.NotEmpty(t, input.CCIPObjectRef)
		require.NotEmpty(t, input.OwnerCapObjectId)
	})

	t.Run("DestroyFeeQuoterCapInput", func(t *testing.T) {
		input := DestroyFeeQuoterCapInput{
			CCIPPackageId:        "0x123",
			CCIPObjectRef:        "0xabc",
			OwnerCapObjectId:     "0x456",
			FeeQuoterCapObjectId: "0x789",
		}

		require.NotEmpty(t, input.CCIPPackageId)
		require.NotEmpty(t, input.CCIPObjectRef)
		require.NotEmpty(t, input.OwnerCapObjectId)
		require.NotEmpty(t, input.FeeQuoterCapObjectId)
	})

	t.Run("FeeQuoterUpdatePricesWithOwnerCapInput", func(t *testing.T) {
		input := FeeQuoterUpdatePricesWithOwnerCapInput{
			CCIPPackageId:         "0x123",
			CCIPObjectRef:         "0x456",
			OwnerCapObjectId:      "0x789",
			SourceTokens:          []string{"0xdef"},
			SourceUsdPerToken:     []*big.Int{big.NewInt(1000000000000000000)},
			GasDestChainSelectors: []uint64{1},
			GasUsdPerUnitGas:      []*big.Int{big.NewInt(2000000000000000)},
		}

		require.NotEmpty(t, input.CCIPPackageId)
		require.NotEmpty(t, input.CCIPObjectRef)
		require.NotEmpty(t, input.OwnerCapObjectId)
		require.Len(t, input.SourceTokens, 1)
		require.Len(t, input.SourceUsdPerToken, 1)
		require.Len(t, input.GasDestChainSelectors, 1)
		require.Len(t, input.GasUsdPerUnitGas, 1)
		require.Equal(t, input.SourceTokens, []string{"0xdef"})
		require.Equal(t, input.SourceUsdPerToken[0].String(), "1000000000000000000")
		require.Equal(t, input.GasDestChainSelectors, []uint64{1})
		require.Equal(t, input.GasUsdPerUnitGas[0].String(), "2000000000000000")
	})
}

func TestFeeQuoterObjects(t *testing.T) {
	t.Run("InitFeeQuoterObjects", func(t *testing.T) {
		objects := InitFeeQuoterObjects{
			FeeQuoterCapObjectId:   "0x123",
			FeeQuoterStateObjectId: "0x456",
		}

		require.NotEmpty(t, objects.FeeQuoterCapObjectId)
		require.NotEmpty(t, objects.FeeQuoterStateObjectId)
	})

	t.Run("NewFeeQuoterCapObjects", func(t *testing.T) {
		objects := NewFeeQuoterCapObjects{
			FeeQuoterCapObjectId: "0x123",
		}

		require.NotEmpty(t, objects.FeeQuoterCapObjectId)
	})

	t.Run("NoObjects", func(t *testing.T) {
		objects := NoObjects{}
		// No fields to test, just ensure it compiles
		require.NotNil(t, objects)
	})
}

func TestBigIntConversion(t *testing.T) {
	t.Run("MaxFeeJuelsConversion", func(t *testing.T) {
		maxFeeJuelsStr := "1000000000000000000" // 1e18
		maxFeeJuels, ok := new(big.Int).SetString(maxFeeJuelsStr, 10)
		require.True(t, ok, "should parse max fee juels string")
		require.Equal(t, maxFeeJuelsStr, maxFeeJuels.String())
	})

	t.Run("PriceConversion", func(t *testing.T) {
		priceStr := "1000000000000000000" // 1 USD in wei
		price, ok := new(big.Int).SetString(priceStr, 10)
		require.True(t, ok, "should parse price string")
		require.Equal(t, priceStr, price.String())
	})

	t.Run("GasPriceConversion", func(t *testing.T) {
		gasPriceStr := "2000000000000000" // 0.002 USD per unit gas
		gasPrice, ok := new(big.Int).SetString(gasPriceStr, 10)
		require.True(t, ok, "should parse gas price string")
		require.Equal(t, gasPriceStr, gasPrice.String())
	})
}
