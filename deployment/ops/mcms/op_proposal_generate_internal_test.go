package mcmsops

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/smartcontractkit/mcms/types"
	"github.com/stretchr/testify/require"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/deployment/utils"
)

// updateMinDelayPayload builds a well-formed [32 addr | 8 u64 LE] payload
// matching mcms::deserialize_timelock_update_min_delay's expected input.
func updateMinDelayPayload(newMinDelaySeconds uint64) []byte {
	buf := make([]byte, updateMinDelayPayloadLen)
	// addr bytes stay zero; the Move-side validate_obj_addr asserts against
	// the actual timelock object at execute time. The value under test lives
	// in the trailing u64.
	binary.LittleEndian.PutUint64(buf[32:updateMinDelayPayloadLen], newMinDelaySeconds)
	return buf
}

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

func TestAssertUpdateMinDelayWithinCap(t *testing.T) {
	t.Parallel()

	maxSeconds := uint64(utils.MaxTimelockScheduleDelay / time.Second)

	t.Run("accepts value at the cap", func(t *testing.T) {
		t.Parallel()
		err := assertUpdateMinDelayWithinCap(sui_ops.TransactionCall{
			Module:   mcmsModuleName,
			Function: updateMinDelayFunctionName,
			Data:     updateMinDelayPayload(maxSeconds),
		})
		require.NoError(t, err)
	})

	t.Run("accepts value under the cap", func(t *testing.T) {
		t.Parallel()
		err := assertUpdateMinDelayWithinCap(sui_ops.TransactionCall{
			Module:   mcmsModuleName,
			Function: updateMinDelayFunctionName,
			Data:     updateMinDelayPayload(uint64(24 * time.Hour / time.Second)),
		})
		require.NoError(t, err)
	})

	t.Run("accepts zero (no-op update)", func(t *testing.T) {
		t.Parallel()
		err := assertUpdateMinDelayWithinCap(sui_ops.TransactionCall{
			Module:   mcmsModuleName,
			Function: updateMinDelayFunctionName,
			Data:     updateMinDelayPayload(0),
		})
		require.NoError(t, err)
	})

	t.Run("rejects value one second over the cap", func(t *testing.T) {
		t.Parallel()
		err := assertUpdateMinDelayWithinCap(sui_ops.TransactionCall{
			Module:   mcmsModuleName,
			Function: updateMinDelayFunctionName,
			Data:     updateMinDelayPayload(maxSeconds + 1),
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "F8 defense")
		require.Contains(t, err.Error(), "exceeds MaxTimelockScheduleDelay")
	})

	t.Run("rejects u64 max (the exact F8 attack value)", func(t *testing.T) {
		t.Parallel()
		err := assertUpdateMinDelayWithinCap(sui_ops.TransactionCall{
			Module:   mcmsModuleName,
			Function: updateMinDelayFunctionName,
			Data:     updateMinDelayPayload(^uint64(0)),
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "F8 defense")
	})

	t.Run("rejects malformed payload length", func(t *testing.T) {
		t.Parallel()
		err := assertUpdateMinDelayWithinCap(sui_ops.TransactionCall{
			Module:   mcmsModuleName,
			Function: updateMinDelayFunctionName,
			Data:     make([]byte, updateMinDelayPayloadLen-1),
		})
		require.Error(t, err)
		require.Contains(t, err.Error(), "must be 40 bytes")
	})

	t.Run("ignores non-mcms modules", func(t *testing.T) {
		t.Parallel()
		err := assertUpdateMinDelayWithinCap(sui_ops.TransactionCall{
			Module:   "state_object",
			Function: updateMinDelayFunctionName,
			Data:     updateMinDelayPayload(^uint64(0)),
		})
		require.NoError(t, err)
	})

	t.Run("ignores other mcms functions", func(t *testing.T) {
		t.Parallel()
		err := assertUpdateMinDelayWithinCap(sui_ops.TransactionCall{
			Module:   mcmsModuleName,
			Function: "mcms_timelock_block_function",
			Data:     make([]byte, 100),
		})
		require.NoError(t, err)
	})
}
