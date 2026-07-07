package mcmsops

import (
	"testing"

	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
)

func TestAssertBypassAllowsCall(t *testing.T) {
	t.Parallel()

	forbidden := []string{
		"mcms_timelock_update_min_delay",
		"mcms_timelock_block_function",
		"mcms_timelock_unblock_function",
		"mcms_timelock_cancel",
		"mcms_set_config",
	}

	t.Run("bypass rejects each forbidden mcms admin function", func(t *testing.T) {
		t.Parallel()
		for _, fn := range forbidden {
			fn := fn
			t.Run(fn, func(t *testing.T) {
				t.Parallel()
				err := assertBypassAllowsCall(types.TimelockActionBypass, sui_ops.TransactionCall{
					Module:   mcmsModuleName,
					Function: fn,
				})
				require.Error(t, err)
				require.Contains(t, err.Error(), "F30 defense")
				require.Contains(t, err.Error(), fn)
			})
		}
	})

	t.Run("bypass allows non-mcms calls", func(t *testing.T) {
		t.Parallel()
		err := assertBypassAllowsCall(types.TimelockActionBypass, sui_ops.TransactionCall{
			Module:   "state_object",
			Function: "add_package_id",
		})
		require.NoError(t, err)
	})

	t.Run("bypass allows unrelated mcms functions", func(t *testing.T) {
		t.Parallel()
		err := assertBypassAllowsCall(types.TimelockActionBypass, sui_ops.TransactionCall{
			Module:   mcmsModuleName,
			Function: "some_future_read_only_helper",
		})
		require.NoError(t, err)
	})

	t.Run("schedule allows forbidden functions", func(t *testing.T) {
		t.Parallel()
		for _, fn := range forbidden {
			err := assertBypassAllowsCall(types.TimelockActionSchedule, sui_ops.TransactionCall{
				Module:   mcmsModuleName,
				Function: fn,
			})
			require.NoError(t, err, "schedule path must remain the intended route for %s", fn)
		}
	})

	t.Run("cancel allows forbidden functions", func(t *testing.T) {
		t.Parallel()
		err := assertBypassAllowsCall(types.TimelockActionCancel, sui_ops.TransactionCall{
			Module:   mcmsModuleName,
			Function: "mcms_timelock_cancel",
		})
		require.NoError(t, err)
	})
}
