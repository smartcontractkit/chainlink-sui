package changesets

import (
	"testing"
	"time"

	cselectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain"
	"github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	mcmstypes "github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	lanes "github.com/smartcontractkit/chainlink-ccip/deployment/lanes"

	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

func TestDestChainConfigOverride_Apply(t *testing.T) {
	t.Parallel()

	base := lanes.FeeQuoterDestChainConfig{
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
		V1Params: &lanes.FeeQuoterV1Params{
			MaxNumberOfTokensPerMsg:           5,
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

	// The Glamsterdam bump shape: only the two gas values are overridden; every other
	// field (incl. all of V1Params) must pass through untouched. `want` is a copy of
	// base with exactly those two fields changed, so full-struct equality proves no
	// other field moved.
	gasOverhead := uint32(500_000)
	tokenGasOverhead := uint32(270_000)
	want := base
	want.DestGasOverhead = gasOverhead
	want.DefaultTokenDestGasOverhead = tokenGasOverhead
	got := base
	DestChainConfigOverride{
		DestGasOverhead:             &gasOverhead,
		DefaultTokenDestGasOverhead: &tokenGasOverhead,
	}.apply(&got)

	require.Equal(t, want, got)

	// The full override set, one field at a time.
	isEnabled := false
	maxDataBytes := uint32(16_000)
	maxPerMsgGasLimit := uint32(4_000_000)
	defaultTxGasLimit := uint32(300_000)
	defaultTokenFeeUSDCents := uint16(50)
	networkFeeUSDCents := uint16(15)
	got = base
	DestChainConfigOverride{
		IsEnabled:                   &isEnabled,
		MaxDataBytes:                &maxDataBytes,
		MaxPerMsgGasLimit:           &maxPerMsgGasLimit,
		DefaultTxGasLimit:           &defaultTxGasLimit,
		DefaultTokenFeeUSDCents:     &defaultTokenFeeUSDCents,
		NetworkFeeUSDCents:          &networkFeeUSDCents,
		DestGasOverhead:             &gasOverhead,
		DefaultTokenDestGasOverhead: &tokenGasOverhead,
	}.apply(&got)
	require.Equal(t, isEnabled, got.IsEnabled)
	require.Equal(t, maxDataBytes, got.MaxDataBytes)
	require.Equal(t, maxPerMsgGasLimit, got.MaxPerMsgGasLimit)
	require.Equal(t, defaultTxGasLimit, got.DefaultTxGasLimit)
	require.Equal(t, defaultTokenFeeUSDCents, got.DefaultTokenFeeUSDCents)
	require.Equal(t, networkFeeUSDCents, got.NetworkFeeUSDCents)
	require.Equal(t, gasOverhead, got.DestGasOverhead)
	require.Equal(t, tokenGasOverhead, got.DefaultTokenDestGasOverhead)

	// An empty override is a no-op.
	got = base
	DestChainConfigOverride{}.apply(&got)
	require.Equal(t, base, got)
}

func TestUpdateFeeQuoterDestChainConfigs_VerifyPreconditions(t *testing.T) {
	t.Parallel()

	selector := cselectors.SUI_TESTNET.Selector
	const sepolia = uint64(16015286601757825753)

	env := cldf.Environment{
		BlockChains: chain.NewBlockChains(map[uint64]chain.BlockChain{
			selector: sui.Chain{},
		}),
	}
	cs := UpdateFeeQuoterDestChainConfigs{}

	t.Run("valid", func(t *testing.T) {
		t.Parallel()
		require.NoError(t, cs.VerifyPreconditions(env, UpdateFeeQuoterDestChainConfigsConfig{
			SuiChainSelector: selector,
			Overrides:        map[uint64]DestChainConfigOverride{sepolia: {}},
		}))
	})

	t.Run("empty overrides errors", func(t *testing.T) {
		t.Parallel()
		err := cs.VerifyPreconditions(env, UpdateFeeQuoterDestChainConfigsConfig{SuiChainSelector: selector})
		require.ErrorContains(t, err, "nothing to update")
	})

	t.Run("source chain not in env errors", func(t *testing.T) {
		t.Parallel()
		err := cs.VerifyPreconditions(env, UpdateFeeQuoterDestChainConfigsConfig{
			SuiChainSelector: 999,
			Overrides:        map[uint64]DestChainConfigOverride{sepolia: {}},
		})
		require.ErrorContains(t, err, "not found in Sui chains")
	})

	t.Run("dst equal to src errors", func(t *testing.T) {
		t.Parallel()
		err := cs.VerifyPreconditions(env, UpdateFeeQuoterDestChainConfigsConfig{
			SuiChainSelector: selector,
			Overrides:        map[uint64]DestChainConfigOverride{selector: {}},
		})
		require.ErrorContains(t, err, "cannot be the source chain selector")
	})

	t.Run("invalid timelock config errors", func(t *testing.T) {
		t.Parallel()
		err := cs.VerifyPreconditions(env, UpdateFeeQuoterDestChainConfigsConfig{
			SuiChainSelector: selector,
			Overrides:        map[uint64]DestChainConfigOverride{sepolia: {}},
			TimelockConfig: &utils.TimelockConfig{
				MCMSAction: mcmstypes.TimelockActionSchedule,
				MinDelay:   utils.MaxTimelockScheduleDelay + time.Second,
			},
		})
		require.ErrorContains(t, err, "invalid TimelockConfig")
	})
}
