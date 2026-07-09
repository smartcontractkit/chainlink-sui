package reader

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

func ownerKindPtr(k suirpcv2.Owner_OwnerKind) *suirpcv2.Owner_OwnerKind { return &k }

func objWithOwner(kind suirpcv2.Owner_OwnerKind) *suirpcv2.Object {
	return &suirpcv2.Object{Owner: &suirpcv2.Owner{Kind: ownerKindPtr(kind)}}
}

// shared/immutable objects have version-stable refs and should be fetched once then served from cache.
func TestReaderCache_ObjectMetadata_CachesVersionStable(t *testing.T) {
	t.Parallel()
	for _, kind := range []suirpcv2.Owner_OwnerKind{suirpcv2.Owner_SHARED, suirpcv2.Owner_IMMUTABLE} {
		rc := NewCache(logger.Test(t), CacheConfig{ObjectCacheEnabled: true})
		var calls int32
		loader := func(context.Context) (*suirpcv2.Object, error) {
			atomic.AddInt32(&calls, 1)
			return objWithOwner(kind), nil
		}
		for i := 0; i < 5; i++ {
			obj, err := rc.GetObjectMetadata(context.Background(), "0xabc", loader)
			require.NoError(t, err)
			require.NotNil(t, obj)
		}
		require.Equal(t, int32(1), atomic.LoadInt32(&calls), "kind %v should be fetched once and cached", kind)
	}
}

// address-owned objects bump their version on mutation, so they must never be cached.
func TestReaderCache_ObjectMetadata_DoesNotCacheAddressOwned(t *testing.T) {
	t.Parallel()
	rc := NewCache(logger.Test(t), CacheConfig{ObjectCacheEnabled: true})
	var calls int32
	loader := func(context.Context) (*suirpcv2.Object, error) {
		atomic.AddInt32(&calls, 1)
		return objWithOwner(suirpcv2.Owner_ADDRESS), nil
	}
	for i := 0; i < 3; i++ {
		_, err := rc.GetObjectMetadata(context.Background(), "0xowned", loader)
		require.NoError(t, err)
	}
	require.Equal(t, int32(3), atomic.LoadInt32(&calls), "address-owned object must be re-fetched every call")
}

// concurrent reads of the same object collapse onto one in-flight load (the cold-start-storm relief).
func TestReaderCache_ObjectMetadata_SingleflightCollapsesConcurrent(t *testing.T) {
	t.Parallel()
	rc := NewCache(logger.Test(t), CacheConfig{ObjectCacheEnabled: true})
	var calls int32
	release := make(chan struct{})
	loader := func(context.Context) (*suirpcv2.Object, error) {
		atomic.AddInt32(&calls, 1)
		<-release // hold the in-flight load so every goroutine piles onto the same flight
		return objWithOwner(suirpcv2.Owner_SHARED), nil
	}

	var wg sync.WaitGroup
	const n = 20
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			_, _ = rc.GetObjectMetadata(context.Background(), "0xshared", loader)
		}()
	}
	time.Sleep(50 * time.Millisecond) // let the goroutines queue on the singleflight key
	close(release)
	wg.Wait()

	require.Equal(t, int32(1), atomic.LoadInt32(&calls), "concurrent reads of the same object should collapse to one load")
}

// a disabled object cache (and a nil cache) must pass every call straight through to the loader.
func TestReaderCache_ObjectMetadata_DisabledAndNilPassthrough(t *testing.T) {
	t.Parallel()
	loader := func(c *int32) func(context.Context) (*suirpcv2.Object, error) {
		return func(context.Context) (*suirpcv2.Object, error) {
			atomic.AddInt32(c, 1)
			return objWithOwner(suirpcv2.Owner_SHARED), nil
		}
	}

	disabled := NewCache(logger.Test(t), CacheConfig{ObjectCacheEnabled: false})
	var dCalls int32
	for i := 0; i < 3; i++ {
		_, err := disabled.GetObjectMetadata(context.Background(), "0xabc", loader(&dCalls))
		require.NoError(t, err)
	}
	require.Equal(t, int32(3), atomic.LoadInt32(&dCalls), "disabled object cache should pass through")

	var nilCache *Cache
	var nCalls int32
	obj, err := nilCache.GetObjectMetadata(context.Background(), "0xabc", loader(&nCalls))
	require.NoError(t, err)
	require.NotNil(t, obj)
	require.Equal(t, int32(1), atomic.LoadInt32(&nCalls), "nil cache should pass through")
}

// read-call caching reuses results within the TTL when enabled, and passes through when disabled.
func TestReaderCache_ReadResults(t *testing.T) {
	t.Parallel()

	enabled := NewCache(logger.Test(t), CacheConfig{ReadCacheEnabled: true, ReadTTL: time.Minute})
	var eCalls int32
	for i := 0; i < 4; i++ {
		res, err := enabled.GetReadResults(context.Background(), "key", func(context.Context) ([]any, error) {
			atomic.AddInt32(&eCalls, 1)
			return []any{"v"}, nil
		})
		require.NoError(t, err)
		require.Equal(t, []any{"v"}, res)
	}
	require.Equal(t, int32(1), atomic.LoadInt32(&eCalls), "enabled read cache should reuse the result")

	disabled := NewCache(logger.Test(t), CacheConfig{ReadCacheEnabled: false})
	var dCalls int32
	for i := 0; i < 3; i++ {
		_, err := disabled.GetReadResults(context.Background(), "key", func(context.Context) ([]any, error) {
			atomic.AddInt32(&dCalls, 1)
			return []any{"v"}, nil
		})
		require.NoError(t, err)
	}
	require.Equal(t, int32(3), atomic.LoadInt32(&dCalls), "disabled read cache should pass through")
}

// A cache hit must return an independent deep copy: callers mutate the decoded result in place (e.g.
// GetLatestValue renames the OffRamp OCR config's `big_f` field), so handing out the shared cached
// reference would corrupt the cached value and break every subsequent read of the same key. This is the
// regression that made the EVM->Sui config read return a broken/zero OCR config once caching was on.
func TestReaderCache_ReadResults_HitReturnsIndependentCopy(t *testing.T) {
	t.Parallel()

	rc := NewCache(logger.Test(t), CacheConfig{ReadCacheEnabled: true, ReadTTL: time.Minute})
	var calls int32
	loader := func(context.Context) ([]any, error) {
		atomic.AddInt32(&calls, 1)
		// Mirrors the shape of a Sui OCR config read (nested map with the big_f field).
		return []any{map[string]any{
			"config_info": map[string]any{"big_f": float64(1), "n": float64(2)},
			"signers":     []any{"a", "b"},
		}}, nil
	}

	// First read populates the cache, then mutates its own copy the way GetLatestValue does (rename
	// big_f -> f in place).
	first, err := rc.GetReadResults(context.Background(), "k", loader)
	require.NoError(t, err)
	ci := first[0].(map[string]any)["config_info"].(map[string]any)
	ci["f"] = ci["big_f"]
	delete(ci, "big_f")
	first[0].(map[string]any)["signers"].([]any)[0] = "MUTATED"

	// Second read is a cache hit (loader not called again) and must see the pristine original value,
	// unaffected by the first caller's in-place mutation.
	second, err := rc.GetReadResults(context.Background(), "k", loader)
	require.NoError(t, err)
	require.Equal(t, int32(1), atomic.LoadInt32(&calls), "second read must be served from cache")

	ci2 := second[0].(map[string]any)["config_info"].(map[string]any)
	require.Equal(t, float64(1), ci2["big_f"], "cached big_f must survive a prior caller's rename")
	_, renamed := ci2["f"]
	require.False(t, renamed, "prior caller's rename must not leak into the cached value")
	require.Equal(t, "a", second[0].(map[string]any)["signers"].([]any)[0],
		"prior caller's slice mutation must not leak into the cached value")
}

// A transient read failure after the fresh TTL has expired must be served from the last known-good
// value (serve-stale), not surfaced as an error/zero result. This is the direct fix for the EVM->Sui
// commit blocker where a cancelled config read produced a zero OCR config.
func TestReaderCache_ReadResults_ServesStaleOnTransientFailure(t *testing.T) {
	t.Parallel()

	rc := NewCache(logger.Test(t), CacheConfig{
		ReadCacheEnabled: true,
		ReadTTL:          20 * time.Millisecond, // short fresh window so we can force a re-fetch
		StaleReadTTL:     time.Minute,
	})

	good := []any{map[string]any{"config_digest": "digest-1"}}
	_, err := rc.GetReadResults(context.Background(), "k", func(context.Context) ([]any, error) {
		return good, nil
	})
	require.NoError(t, err)

	// Let the fresh entry expire so the next read must call the loader.
	time.Sleep(40 * time.Millisecond)

	// Loader now fails (simulating a context cancellation during a slow poll). Serve-stale must return
	// the last good value instead of the error.
	res, err := rc.GetReadResults(context.Background(), "k", func(context.Context) ([]any, error) {
		return nil, context.Canceled
	})
	require.NoError(t, err, "transient failure after TTL expiry must be served from the stale cache")
	require.Equal(t, "digest-1", res[0].(map[string]any)["config_digest"])

	// With no prior success there is nothing to serve stale, so the error propagates.
	_, err = rc.GetReadResults(context.Background(), "never-loaded", func(context.Context) ([]any, error) {
		return nil, context.Canceled
	})
	require.ErrorIs(t, err, context.Canceled, "a cold key with no stale entry must surface the read error")
}
