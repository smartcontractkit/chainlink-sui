package lanes

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-ccip/deployment/utils/sequences"
	mcmstypes "github.com/smartcontractkit/mcms/types"

	sui_ops "github.com/smartcontractkit/chainlink-sui/deployment/ops"
	"github.com/smartcontractkit/chainlink-sui/relayer/signer"
)

type stubSuiSigner struct{}

func (stubSuiSigner) Sign(_ []byte) ([]string, error) { return nil, nil }

func (stubSuiSigner) GetAddress() (string, error) { return "0x1", nil }

var _ signer.SuiSigner = stubSuiSigner{}

func TestAppendMCMSBatchOpFromCall(t *testing.T) {
	t.Parallel()

	const chainSelector uint64 = 9762610643973837292

	call := sui_ops.TransactionCall{
		PackageID:  "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		Module:     "offramp",
		Function:   "apply_source_chain_config_updates",
		Data:       []byte{0x01},
		StateObjID: "0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
	}

	t.Run("appends batch op when signer is nil", func(t *testing.T) {
		var out sequences.OnChainOutput
		err := appendMCMSBatchOpFromCall(&out, chainSelector, call, sui_ops.OpTxDeps{})
		require.NoError(t, err)
		require.Len(t, out.BatchOps, 1)
		require.Equal(t, mcmstypes.ChainSelector(chainSelector), out.BatchOps[0].ChainSelector)
		require.Len(t, out.BatchOps[0].Transactions, 1)
	})

	t.Run("skips batch op when signer is set", func(t *testing.T) {
		var out sequences.OnChainOutput
		err := appendMCMSBatchOpFromCall(&out, chainSelector, call, sui_ops.OpTxDeps{
			Signer: stubSuiSigner{},
		})
		require.NoError(t, err)
		require.Empty(t, out.BatchOps)
	})

	t.Run("skips empty package id", func(t *testing.T) {
		var out sequences.OnChainOutput
		err := appendMCMSBatchOpFromCall(&out, chainSelector, sui_ops.TransactionCall{}, sui_ops.OpTxDeps{})
		require.NoError(t, err)
		require.Empty(t, out.BatchOps)
	})
}
