package bind

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-sui/relayer/client"
)

type mockTxStatusClient struct {
	fn func(ctx context.Context, digest string) (client.TransactionResult, error)
}

func (m *mockTxStatusClient) GetTransactionStatus(ctx context.Context, digest string) (client.TransactionResult, error) {
	return m.fn(ctx, digest)
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
	mock := &mockTxStatusClient{
		fn: func(ctx context.Context, digest string) (client.TransactionResult, error) {
			calls++
			require.Equal(t, wantDigest, digest)
			if calls < 3 {
				return client.TransactionResult{}, errors.New("not indexed yet")
			}
			return client.TransactionResult{Status: "success"}, nil
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

	mock := &mockTxStatusClient{
		fn: func(ctx context.Context, digest string) (client.TransactionResult, error) {
			return client.TransactionResult{}, errors.New("still missing")
		},
	}

	err := WaitForTransactionIndexed(context.Background(), mock, "missingdigest")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrTxIndexingTimeout)
}

func TestWaitForTransactionIndexed_NonSuccessStatusKeepsPolling(t *testing.T) {
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
	mock := &mockTxStatusClient{
		fn: func(ctx context.Context, digest string) (client.TransactionResult, error) {
			calls++
			return client.TransactionResult{Status: "pending"}, nil
		},
	}

	err := WaitForTransactionIndexed(context.Background(), mock, "expected")
	require.Error(t, err)
	require.ErrorIs(t, err, ErrTxIndexingTimeout)
	require.GreaterOrEqual(t, calls, 2)
}
