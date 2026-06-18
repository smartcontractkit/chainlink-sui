package reader

import (
	"context"
	"time"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	cache "github.com/patrickmn/go-cache"
	"golang.org/x/sync/singleflight"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// ReaderCacheConfig tunes the ReaderCache. Zero-value durations fall back to the values in
// DefaultReaderCacheConfig via withDefaults(), so a partially-populated config is always valid.
type ReaderCacheConfig struct {
	// ObjectCacheEnabled caches version-stable object reference metadata (owner/version/digest). The Sui
	// read hot path resolves the same shared CCIP config objects (the CCIPObjectRef, OffRamp state, fee
	// quoter, ...) on every read; caching their refs removes that redundant GetObject fan-out so it does
	// not hit the node on every read across every config-poll cycle.
	ObjectCacheEnabled bool
	// ObjectTTL is how long object metadata is retained. Only SHARED/IMMUTABLE objects are cached (their
	// refs are version-stable), so this can be minutes without serving a stale ref.
	ObjectTTL time.Duration
	// ReadCacheEnabled caches decoded read-call (devInspect) results keyed by read identifier + params.
	// It is OFF by default: config reads change rarely, but caching them trades a little staleness for
	// fewer node round-trips, so it is opt-in and should be enabled only with a short ReadTTL.
	ReadCacheEnabled bool
	// ReadTTL bounds how long a cached read result may be served before it is re-fetched.
	ReadTTL time.Duration
	// CleanupInterval is how often expired entries are purged from both underlying caches.
	CleanupInterval time.Duration
}

// DefaultReaderCacheConfig returns safe defaults: object caching ON (high value, no staleness risk for
// version-stable objects) and read-call caching OFF (opt-in, since it can briefly mask config changes).
func DefaultReaderCacheConfig() ReaderCacheConfig {
	return ReaderCacheConfig{
		ObjectCacheEnabled: true,
		ObjectTTL:          5 * time.Minute,
		ReadCacheEnabled:   false,
		ReadTTL:            2 * time.Second,
		CleanupInterval:    1 * time.Minute,
	}
}

func (c ReaderCacheConfig) withDefaults() ReaderCacheConfig {
	d := DefaultReaderCacheConfig()
	if c.ObjectTTL <= 0 {
		c.ObjectTTL = d.ObjectTTL
	}
	if c.ReadTTL <= 0 {
		c.ReadTTL = d.ReadTTL
	}
	if c.CleanupInterval <= 0 {
		c.CleanupInterval = d.CleanupInterval
	}
	return c
}

// ReaderCache de-duplicates and caches the read-path RPCs that dominate Sui CCIP config polling:
//   - object reference metadata (GetObject) for shared/immutable objects, and
//   - optionally, decoded read-call (devInspect) results.
//
// Concurrent identical loads are collapsed via singleflight — this is what relieves the cold-start storm
// where every node polls config simultaneously with empty caches — and successful loads are retained for
// a configurable TTL. The cache is process-local, which is sufficient because the redundancy it removes
// lives within a single relayer instance (each node re-reads the same shared objects on every poll); it
// does not need to be shared across nodes.
//
// All methods are safe to call on a nil *ReaderCache: they degrade to invoking the loader directly, so
// callers can hold an optional cache without nil checks.
type ReaderCache struct {
	lggr logger.Logger
	cfg  ReaderCacheConfig

	objectCache *cache.Cache
	objectGroup singleflight.Group

	readCache *cache.Cache
	readGroup singleflight.Group
}

// NewReaderCache builds a ReaderCache from cfg (missing durations are defaulted).
func NewReaderCache(lggr logger.Logger, cfg ReaderCacheConfig) *ReaderCache {
	cfg = cfg.withDefaults()
	return &ReaderCache{
		lggr:        logger.Named(lggr, "ReaderCache"),
		cfg:         cfg,
		objectCache: cache.New(cfg.ObjectTTL, cfg.CleanupInterval),
		readCache:   cache.New(cfg.ReadTTL, cfg.CleanupInterval),
	}
}

// GetObjectMetadata returns the object's reference metadata, serving it from cache when possible and
// otherwise invoking loader exactly once across concurrent callers for the same objectID. Only
// version-stable objects (SHARED/IMMUTABLE) are cached; address-owned objects, whose version changes on
// mutation, always go to loader so a stale ref is never served. This method satisfies the
// client.ObjectMetadataCache interface used by the gRPC client.
func (rc *ReaderCache) GetObjectMetadata(
	ctx context.Context,
	objectID string,
	loader func(context.Context) (*suirpcv2.Object, error),
) (*suirpcv2.Object, error) {
	if rc == nil || !rc.cfg.ObjectCacheEnabled {
		return loader(ctx)
	}

	if obj, ok := getTyped[*suirpcv2.Object](rc.objectCache, objectID); ok {
		return obj, nil
	}

	v, err, _ := rc.objectGroup.Do(objectID, func() (any, error) {
		// A prior flight may have populated the cache while this call was queued.
		if obj, ok := getTyped[*suirpcv2.Object](rc.objectCache, objectID); ok {
			return obj, nil
		}
		obj, err := loader(ctx)
		if err != nil {
			return nil, err
		}
		if isVersionStable(obj) {
			rc.objectCache.Set(objectID, obj, cache.DefaultExpiration)
		}
		return obj, nil
	})
	if err != nil {
		return nil, err
	}

	obj, _ := v.(*suirpcv2.Object)
	return obj, nil
}

// GetReadResults returns the decoded results of a read call, serving them from cache when enabled and
// otherwise invoking loader exactly once across concurrent callers for the same key. The cached value is
// the raw []any decoded result; callers decode it into their own return value on each call, so no shared
// mutable state escapes. Disabled by default — see ReaderCacheConfig.ReadCacheEnabled.
func (rc *ReaderCache) GetReadResults(
	ctx context.Context,
	key string,
	loader func(context.Context) ([]any, error),
) ([]any, error) {
	if rc == nil || !rc.cfg.ReadCacheEnabled {
		return loader(ctx)
	}

	if res, ok := getTyped[[]any](rc.readCache, key); ok {
		return res, nil
	}

	v, err, _ := rc.readGroup.Do(key, func() (any, error) {
		if res, ok := getTyped[[]any](rc.readCache, key); ok {
			return res, nil
		}
		res, err := loader(ctx)
		if err != nil {
			return nil, err
		}
		rc.readCache.Set(key, res, cache.DefaultExpiration)
		return res, nil
	})
	if err != nil {
		return nil, err
	}

	res, _ := v.([]any)
	return res, nil
}

// getTyped fetches key from c and asserts it to T, returning ok=false on miss or type mismatch.
func getTyped[T any](c *cache.Cache, key string) (T, bool) {
	var zero T
	if cached, found := c.Get(key); found {
		if typed, ok := cached.(T); ok {
			return typed, true
		}
	}
	return zero, false
}

// isVersionStable reports whether an object's reference is immutable for the object's lifetime and is
// therefore safe to cache. Shared objects expose an immutable InitialSharedVersion and immutable objects
// never change; address-owned objects bump their version on every mutation and must not be cached.
func isVersionStable(obj *suirpcv2.Object) bool {
	if obj == nil || obj.GetOwner() == nil {
		return false
	}
	switch obj.GetOwner().GetKind() {
	case suirpcv2.Owner_IMMUTABLE, suirpcv2.Owner_SHARED:
		return true
	default:
		return false
	}
}
