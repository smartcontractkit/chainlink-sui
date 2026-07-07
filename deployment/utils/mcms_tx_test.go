package utils_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

func TestTransactionCallToMCMSTransaction_LatestPackageID(t *testing.T) {
	t.Parallel()

	const (
		originalPackageID = "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
		latestPackageID   = "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	)

	call := sui_ops.TransactionCall{
		PackageID:       originalPackageID,
		LatestPackageID: latestPackageID,
		Module:          "offramp",
		Function:        "apply_source_chain_config_updates",
		Data:            []byte{0x01},
		StateObjID:      "0xcccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
	}

	tx, err := utils.TransactionCallToMCMSTransaction(call)
	require.NoError(t, err)
	require.Equal(t, originalPackageID, tx.To)
	latestPackageIDFromTx, err := utils.TransactionLatestPackageID(tx)
	require.NoError(t, err)
	require.Equal(t, latestPackageID, latestPackageIDFromTx)
}

func TestTimelockConfig_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		cfg     utils.TimelockConfig
		wantErr string
	}{
		{
			name: "schedule with delay under cap is accepted",
			cfg: utils.TimelockConfig{
				MCMSAction: types.TimelockActionSchedule,
				MinDelay:   24 * time.Hour,
			},
		},
		{
			name: "schedule with delay equal to cap is accepted",
			cfg: utils.TimelockConfig{
				MCMSAction: types.TimelockActionSchedule,
				MinDelay:   utils.MaxTimelockScheduleDelay,
			},
		},
		{
			name: "schedule with delay above cap is rejected",
			cfg: utils.TimelockConfig{
				MCMSAction: types.TimelockActionSchedule,
				MinDelay:   utils.MaxTimelockScheduleDelay + time.Second,
			},
			wantErr: "exceeds MaxTimelockScheduleDelay",
		},
		{
			name: "schedule with negative delay is rejected",
			cfg: utils.TimelockConfig{
				MCMSAction: types.TimelockActionSchedule,
				MinDelay:   -time.Second,
			},
			wantErr: "must be non-negative",
		},
		{
			name: "bypass action skips delay validation",
			cfg: utils.TimelockConfig{
				MCMSAction: types.TimelockActionBypass,
				MinDelay:   utils.MaxTimelockScheduleDelay * 1000,
			},
		},
		{
			name: "cancel action skips delay validation",
			cfg: utils.TimelockConfig{
				MCMSAction: types.TimelockActionCancel,
				MinDelay:   -time.Hour,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.cfg.Validate()
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestAssertMinDelayWithinCap(t *testing.T) {
	t.Parallel()

	maxSeconds := uint64(utils.MaxTimelockScheduleDelay / time.Second)

	tests := []struct {
		name             string
		newMinDelayInSec uint64
		wantErr          string
	}{
		{
			name:             "zero is accepted (bootstrap window; F5 warned elsewhere)",
			newMinDelayInSec: 0,
		},
		{
			name:             "value under the cap is accepted",
			newMinDelayInSec: uint64(24 * time.Hour / time.Second),
		},
		{
			name:             "value at the cap is accepted",
			newMinDelayInSec: maxSeconds,
		},
		{
			name:             "value one second over the cap is rejected",
			newMinDelayInSec: maxSeconds + 1,
			wantErr:          "F8 defense",
		},
		{
			name:             "u64 max (the F8 attack value) is rejected",
			newMinDelayInSec: ^uint64(0),
			wantErr:          "exceeds MaxTimelockScheduleDelay",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := utils.AssertMinDelayWithinCap(tt.newMinDelayInSec)
			if tt.wantErr == "" {
				require.NoError(t, err)
				return
			}
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.wantErr)
		})
	}
}
