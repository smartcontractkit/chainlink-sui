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
