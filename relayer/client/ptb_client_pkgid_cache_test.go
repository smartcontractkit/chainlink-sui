package client

import (
	"context"
	"testing"

	cache "github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/require"
)

// newPkgIDCacheTestClient builds a minimal PTBClient with only the cache wired. The cache-hit and
// invalidation paths under test never touch the network or the logger, so the other fields can be
// left zero-valued.
func newPkgIDCacheTestClient() *PTBClient {
	return &PTBClient{
		cache: cache.New(DefaultCacheExpiration, DefaultCacheCleanupInterval),
	}
}

func TestLoadModulePackageIdsInternal_CacheHitSkipsResolution(t *testing.T) {
	c := newPkgIDCacheTestClient()
	const pkg, module = "0xpkg", "offramp"
	want := []string{"0xv1", "0xv2"}

	// Pre-populate the cache; a hit must return without hitting loadModulePackageIdsUncached
	// (which would dial the RPC via the nil client and panic/error).
	c.cache.Set(packageIDCacheKey(pkg, module), append([]string(nil), want...), PackageIDCacheExpiration)

	got, err := c.loadModulePackageIdsInternal(context.Background(), pkg, module)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestLoadModulePackageIdsInternal_ReturnsCopy(t *testing.T) {
	c := newPkgIDCacheTestClient()
	const pkg, module = "0xpkg", "offramp"
	c.cache.Set(packageIDCacheKey(pkg, module), []string{"0xv1", "0xv2"}, PackageIDCacheExpiration)

	first, err := c.loadModulePackageIdsInternal(context.Background(), pkg, module)
	require.NoError(t, err)

	// Mutating the returned slice must not corrupt the cached value.
	first[0] = "0xMUTATED"

	second, err := c.loadModulePackageIdsInternal(context.Background(), pkg, module)
	require.NoError(t, err)
	require.Equal(t, []string{"0xv1", "0xv2"}, second, "cached slice must be isolated from caller mutations")
}

func TestInvalidatePackageIDCache(t *testing.T) {
	c := newPkgIDCacheTestClient()
	const pkg, module = "0xpkg", "offramp"
	key := packageIDCacheKey(pkg, module)
	c.cache.Set(key, []string{"0xv1"}, PackageIDCacheExpiration)

	_, found := c.cache.Get(key)
	require.True(t, found)

	c.InvalidatePackageIDCache(pkg, module)

	_, found = c.cache.Get(key)
	require.False(t, found, "InvalidatePackageIDCache must evict the cached entry")
}
