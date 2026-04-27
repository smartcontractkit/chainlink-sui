package bind

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/stretchr/testify/require"
)

type mockTxBlockGetter struct {
	fn func(ctx context.Context, req models.SuiGetTransactionBlockRequest) (models.SuiTransactionBlockResponse, error)
}

func (m *mockTxBlockGetter) SuiGetTransactionBlock(ctx context.Context, req models.SuiGetTransactionBlockRequest) (models.SuiTransactionBlockResponse, error) {
	return m.fn(ctx, req)
}

func TestWaitForTransactionIndexed_SucceedsAfterRetries(t *testing.T) {
	oldTimeout := WaitForTxIndexedTimeout
	oldInit := WaitForTxIndexedInitialBackoff
	oldMax := WaitForTxIndexedMaxBackoff
	t.Cleanup(func() {
		WaitForTxIndexedTimeout = oldTimeout
		WaitForTxIndexedInitialBackoff = oldInit
		WaitForTxIndexedMaxBackoff = oldMax
	})
	WaitForTxIndexedTimeout = 5 * time.Second
	WaitForTxIndexedInitialBackoff = time.Millisecond
	WaitForTxIndexedMaxBackoff = 2 * time.Millisecond

	const wantDigest = "wantdigest"
	var calls int
	mock := &mockTxBlockGetter{
		fn: func(ctx context.Context, req models.SuiGetTransactionBlockRequest) (models.SuiTransactionBlockResponse, error) {
			calls++
			if calls < 3 {
				return models.SuiTransactionBlockResponse{}, errors.New("not indexed yet")
			}
			require.Equal(t, wantDigest, req.Digest)
			return models.SuiTransactionBlockResponse{Digest: wantDigest}, nil
		},
	}

	err := WaitForTransactionIndexed(context.Background(), mock, wantDigest)
	require.NoError(t, err)
	require.Equal(t, 3, calls)
}

func TestWaitForTransactionIndexed_Timeout(t *testing.T) {
	oldTimeout := WaitForTxIndexedTimeout
	oldInit := WaitForTxIndexedInitialBackoff
	oldMax := WaitForTxIndexedMaxBackoff
	t.Cleanup(func() {
		WaitForTxIndexedTimeout = oldTimeout
		WaitForTxIndexedInitialBackoff = oldInit
		WaitForTxIndexedMaxBackoff = oldMax
	})
	WaitForTxIndexedTimeout = 50 * time.Millisecond
	WaitForTxIndexedInitialBackoff = 5 * time.Millisecond
	WaitForTxIndexedMaxBackoff = 5 * time.Millisecond

	mock := &mockTxBlockGetter{
		fn: func(ctx context.Context, req models.SuiGetTransactionBlockRequest) (models.SuiTransactionBlockResponse, error) {
			return models.SuiTransactionBlockResponse{}, errors.New("still missing")
		},
	}

	err := WaitForTransactionIndexed(context.Background(), mock, "missingdigest")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrTxIndexingTimeout)
}

func TestWaitForTransactionIndexed_WrongDigestKeepsPolling(t *testing.T) {
	oldTimeout := WaitForTxIndexedTimeout
	oldInit := WaitForTxIndexedInitialBackoff
	oldMax := WaitForTxIndexedMaxBackoff
	t.Cleanup(func() {
		WaitForTxIndexedTimeout = oldTimeout
		WaitForTxIndexedInitialBackoff = oldInit
		WaitForTxIndexedMaxBackoff = oldMax
	})
	WaitForTxIndexedTimeout = 80 * time.Millisecond
	WaitForTxIndexedInitialBackoff = 5 * time.Millisecond
	WaitForTxIndexedMaxBackoff = 5 * time.Millisecond

	var calls int
	mock := &mockTxBlockGetter{
		fn: func(ctx context.Context, req models.SuiGetTransactionBlockRequest) (models.SuiTransactionBlockResponse, error) {
			calls++
			// RPC returns a response but digest mismatch (should not happen in practice)
			return models.SuiTransactionBlockResponse{Digest: "other"}, nil
		},
	}

	err := WaitForTransactionIndexed(context.Background(), mock, "expected")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrTxIndexingTimeout)
	require.GreaterOrEqual(t, calls, 2)
}
