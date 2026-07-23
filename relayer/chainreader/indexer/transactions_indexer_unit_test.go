package indexer

import (
	"context"
	"errors"
	"testing"
	"time"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func TestParseMoveAbortFromExecutionError(t *testing.T) {
	t.Parallel()

	kind := suirpcv2.ExecutionError_MOVE_ABORT
	command := uint64(1)
	functionName := "finish_execute"
	pkg := "0591c76156cdbf86aa4bcd42398c369c9d20151a23d9c5803f2b21e41f93ec61"
	module := "offramp"
	function := uint32(5)
	instruction := uint32(0)
	abortCode := uint64(1)

	execErr := &suirpcv2.ExecutionError{
		Kind:    &kind,
		Command: &command,
		ErrorDetails: &suirpcv2.ExecutionError_Abort{
			Abort: &suirpcv2.MoveAbort{
				AbortCode: &abortCode,
				Location: &suirpcv2.MoveLocation{
					Package:      &pkg,
					Module:       &module,
					Function:     &function,
					Instruction:  &instruction,
					FunctionName: &functionName,
				},
			},
		},
	}

	indexer := &TransactionsIndexer{}
	moveAbort, err := indexer.parseMoveAbortFromExecutionError(execErr)
	require.NoError(t, err)
	require.NotNil(t, moveAbort)
	require.Equal(t, uint64(1), moveAbort.CommandIndex)
	require.Equal(t, uint64(1), moveAbort.AbortCode)
	require.Equal(t, pkg, moveAbort.Location.Module.Address)
	require.Equal(t, module, moveAbort.Location.Module.Name)
	require.Equal(t, "finish_execute", *moveAbort.Location.FunctionName)
}

func TestParseMoveAbortFromExecutionError_RejectsMissingAbort(t *testing.T) {
	t.Parallel()

	description := "MoveAbort(MoveLocation { module: ModuleId { address: abc, name: Identifier(\"offramp\") }, function: 1, instruction: 0, function_name: Some(\"finish_execute\") }, 1) in command 1"
	execErr := &suirpcv2.ExecutionError{
		Description: &description,
	}

	indexer := &TransactionsIndexer{}
	_, err := indexer.parseMoveAbortFromExecutionError(execErr)
	require.Error(t, err)
}

// newCacheTestIndexer builds a TransactionsIndexer whose transmitter fetch is driven
// by fetchFn, so the cache can be exercised without a database.
func newCacheTestIndexer(t *testing.T, fetchFn func(ctx context.Context) ([]string, error)) *TransactionsIndexer {
	t.Helper()
	return &TransactionsIndexer{
		logger:              logger.Test(t),
		fetchTransmittersFn: fetchFn,
	}
}

func TestGetTransmitters_CachesWithinTTL(t *testing.T) {
	t.Parallel()

	calls := 0
	fetch := func(ctx context.Context) ([]string, error) {
		calls++
		return []string{"tx1", "tx2"}, nil
	}
	indexer := newCacheTestIndexer(t, fetch)

	first, err := indexer.getTransmitters(context.Background())
	require.NoError(t, err)
	require.Equal(t, []string{"tx1", "tx2"}, first)
	require.Equal(t, 1, calls)

	second, err := indexer.getTransmitters(context.Background())
	require.NoError(t, err)
	require.Equal(t, []string{"tx1", "tx2"}, second)
	require.Equal(t, 1, calls, "second call within TTL must be served from cache without refetching")
}

func TestGetTransmitters_RefetchesAfterTTL(t *testing.T) {
	t.Parallel()

	calls := 0
	fetch := func(ctx context.Context) ([]string, error) {
		calls++
		if calls == 1 {
			return []string{"tx1"}, nil
		}
		return []string{"tx1", "tx2"}, nil
	}
	indexer := newCacheTestIndexer(t, fetch)

	first, err := indexer.getTransmitters(context.Background())
	require.NoError(t, err)
	require.Equal(t, []string{"tx1"}, first)
	require.Equal(t, 1, calls)

	// Force the cache entry to be stale.
	indexer.mu.Lock()
	indexer.transmittersCachedAt = time.Now().Add(-transmitterCacheTTL - time.Second)
	indexer.mu.Unlock()

	second, err := indexer.getTransmitters(context.Background())
	require.NoError(t, err)
	require.Equal(t, []string{"tx1", "tx2"}, second)
	require.Equal(t, 2, calls, "stale cache must trigger a refetch")
}

func TestGetTransmitters_ReturnsCopy(t *testing.T) {
	t.Parallel()

	fetch := func(ctx context.Context) ([]string, error) {
		return []string{"tx1", "tx2"}, nil
	}
	indexer := newCacheTestIndexer(t, fetch)

	// Cache miss: the returned slice is the fetch result, but the cache stores a
	// clone, so mutating the return must not corrupt the cached value.
	got, err := indexer.getTransmitters(context.Background())
	require.NoError(t, err)
	got[0] = "MUTATED"

	cached, err := indexer.getTransmitters(context.Background())
	require.NoError(t, err)
	require.Equal(t, []string{"tx1", "tx2"}, cached, "cache must hold a copy, not the mutated return")

	// Cache hit: the returned slice is a clone, so mutating it must not corrupt the
	// cached value either.
	cached[1] = "MUTATED"
	again, err := indexer.getTransmitters(context.Background())
	require.NoError(t, err)
	require.Equal(t, []string{"tx1", "tx2"}, again, "cache hit must return a copy")
}

func TestGetTransmitters_FetchErrorNotCached(t *testing.T) {
	t.Parallel()

	calls := 0
	fetch := func(ctx context.Context) ([]string, error) {
		calls++
		if calls == 1 {
			return nil, errors.New("boom")
		}
		return []string{"tx1"}, nil
	}
	indexer := newCacheTestIndexer(t, fetch)

	_, err := indexer.getTransmitters(context.Background())
	require.Error(t, err)

	indexer.mu.RLock()
	cached := indexer.transmittersCached
	indexer.mu.RUnlock()
	require.False(t, cached, "a failed fetch must not populate the cache")

	got, err := indexer.getTransmitters(context.Background())
	require.NoError(t, err)
	require.Equal(t, []string{"tx1"}, got)
	require.Equal(t, 2, calls, "failed first fetch must not short-circuit the second")
}
