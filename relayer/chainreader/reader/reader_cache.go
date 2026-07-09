package reader

import (
	"context"
	"time"

	suirpcv2 "github.com/block-vision/sui-go-sdk/pb/sui/rpc/v2"
	cache "github.com/patrickmn/go-cache"
	"golang.org/x/sync/singleflight"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
)

// CacheConfig tunes the Cache. Zero-value durations fall back to the values in
// DefaultCacheConfig via withDefaults(), so a partially-populated config is always valid.
type CacheConfig struct {
	// ObjectCacheEnabled caches version-stable object reference metadata (owner/version/digest). The Sui
	// read hot path resolves the same shared CCIP config objects (the CCIPObjectRef, OffRamp state, fee
	// quoter, ...) on every read; caching their refs removes that redundant GetObject fan-out so it does
	// not hit the node on every read across every config-poll cycle.
	ObjectCacheEnabled bool
	// ObjectTTL is how long object metadata is retained. Only SHARED/IMMUTABLE objects are cached (their
	// refs are version-stable), so this can be minutes without serving a stale ref.
	ObjectTTL time.Duration
	// ReadCacheEnabled caches decoded read-call (devInspect) results keyed by read identifier + params.
	// Config reads (OffRamp OCR config, source-chain configs, static/dynamic config, price seq number)
	// change rarely, so caching them for a short ReadTTL cuts the redundant devInspect fan-out that
	// otherwise floods the node every config-poll cycle. Combined with StaleReadTTL below, it also makes
	// reads resilient to transient RPC cancellation (the EVM->Sui commit blocker).
	ReadCacheEnabled bool
	// ReadTTL bounds how long a cached read result is served as fresh before it is re-fetched.
	ReadTTL time.Duration
	// StaleReadTTL bounds how long the last successfully-read value may be served as a fallback when a
	// fresh read fails (see GetReadResults). This turns a transient read failure — e.g. a context
	// cancellation while a slow config poll is in flight — into "serve the last known-good config"
	// instead of "return a zero/empty config", which is what makes the Sui commit plugin reject every
	// EVM->Sui report. Only consulted when ReadCacheEnabled is true; a zero value falls back to the
	// default via withDefaults.
	StaleReadTTL time.Duration
	// CleanupInterval is how often expired entries are purged from the underlying caches.
	CleanupInterval time.Duration
}

// DefaultCacheConfig returns safe defaults: object caching ON (high value, no staleness risk for
// version-stable objects) and read-call caching ON with a short freshness TTL plus a bounded
// serve-stale fallback (config reads change rarely, and serving a slightly-stale config is far safer
// than serving a zero config on a transient read failure).
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		ObjectCacheEnabled: true,
		ObjectTTL:          5 * time.Minute,
		ReadCacheEnabled:   false,
		ReadTTL:            10 * time.Second,
		StaleReadTTL:       5 * time.Minute,
		CleanupInterval:    1 * time.Minute,
	}
}

func (c CacheConfig) withDefaults() CacheConfig {
	d := DefaultCacheConfig()
	if c.ObjectTTL <= 0 {
		c.ObjectTTL = d.ObjectTTL
	}
	if c.ReadTTL <= 0 {
		c.ReadTTL = d.ReadTTL
	}
	if c.StaleReadTTL <= 0 {
		c.StaleReadTTL = d.StaleReadTTL
	}
	if c.CleanupInterval <= 0 {
		c.CleanupInterval = d.CleanupInterval
	}
	return c
}

// Cache de-duplicates and caches the read-path RPCs that dominate Sui CCIP config polling:
//   - object reference metadata (GetObject) for shared/immutable objects, and
//   - optionally, decoded read-call (devInspect) results.
//
// Concurrent identical loads are collapsed via singleflight — this is what relieves the cold-start storm
// where every node polls config simultaneously with empty caches — and successful loads are retained for
// a configurable TTL. The cache is process-local, which is sufficient because the redundancy it removes
// lives within a single relayer instance (each node re-reads the same shared objects on every poll); it
// does not need to be shared across nodes.
//
// All methods are safe to call on a nil *Cache: they degrade to invoking the loader directly, so
// callers can hold an optional cache without nil checks.
type Cache struct {
	lggr logger.Logger
	cfg  CacheConfig

	objectCache *cache.Cache
	objectGroup singleflight.Group

	readCache *cache.Cache
	readGroup singleflight.Group
	// staleReadCache retains the last successfully-read value per key for StaleReadTTL, so it can be
	// served as a fallback when a fresh read fails. It intentionally outlives readCache entries.
	staleReadCache *cache.Cache
}

// NewCache builds a Cache from cfg (missing durations are defaulted).
func NewCache(lggr logger.Logger, cfg CacheConfig) *Cache {
	cfg = cfg.withDefaults()
	return &Cache{
		lggr:           logger.Named(lggr, "ReaderCache"),
		cfg:            cfg,
		objectCache:    cache.New(cfg.ObjectTTL, cfg.CleanupInterval),
		readCache:      cache.New(cfg.ReadTTL, cfg.CleanupInterval),
		staleReadCache: cache.New(cfg.StaleReadTTL, cfg.CleanupInterval),
	}
}

// GetObjectMetadata returns the object's reference metadata, serving it from cache when possible and
// otherwise invoking loader exactly once across concurrent callers for the same objectID. Only
// version-stable objects (SHARED/IMMUTABLE) are cached; address-owned objects, whose version changes on
// mutation, always go to loader so a stale ref is never served. This method satisfies the
// client.ObjectMetadataCache interface used by the gRPC client.
func (rc *Cache) GetObjectMetadata(
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
// otherwise invoking loader exactly once across concurrent callers for the same key. Disabled by
// default — see CacheConfig.ReadCacheEnabled.
//
// Every cache hit returns a deep copy of the cached value. This is required for correctness: callers
// (GetLatestValue -> prepareFunctionReadResult -> applyResultFieldRenames/MaybeRenameFields, and the
// NormalizeReturnValuesToHex path) mutate the decoded result in place — e.g. renaming the OffRamp OCR
// config's `big_f` field to `f`. Handing out the cached reference would let the first read mutate the
// shared value, so a later cache hit would see an already-transformed object (e.g. `big_f` missing) and
// fail. Copying on the way out keeps the cached value pristine and each caller's result independent.
func (rc *Cache) GetReadResults(
	ctx context.Context,
	key string,
	loader func(context.Context) ([]any, error),
) ([]any, error) {
	if rc == nil || !rc.cfg.ReadCacheEnabled {
		return loader(ctx)
	}

	if res, ok := getTyped[[]any](rc.readCache, key); ok {
		return deepCopyAnySlice(res), nil
	}

	v, err, _ := rc.readGroup.Do(key, func() (any, error) {
		if res, ok := getTyped[[]any](rc.readCache, key); ok {
			return res, nil
		}
		res, lErr := loader(ctx)
		if lErr != nil {
			// Serve-stale: a fresh read failed (commonly a context cancellation while a slow config
			// poll is in flight). Rather than surfacing a zero/empty result — which upstream turns into
			// a zero OCR config that blocks EVM->Sui commits — fall back to the last successfully-read
			// value if one is still within StaleReadTTL.
			if stale, ok := getTyped[[]any](rc.staleReadCache, key); ok {
				rc.lggr.Warnw("read failed; serving last-known-good cached result", "key", key, "err", lErr)
				return stale, nil
			}
			return nil, lErr
		}
		rc.readCache.Set(key, res, cache.DefaultExpiration)
		rc.staleReadCache.Set(key, res, cache.DefaultExpiration)
		return res, nil
	})
	if err != nil {
		return nil, err
	}

	res, _ := v.([]any)
	// Copy on the way out so concurrent singleflight sharers and later cache hits never mutate the
	// cached value (see doc comment).
	return deepCopyAnySlice(res), nil
}

// deepCopyAnySlice returns a structural deep copy of a decoded read result. The values originate from
// protobuf Struct.AsInterface, so they are only ever composed of map[string]any, []any and immutable
// scalars (string, float64, bool, nil); those containers are the ones downstream transforms mutate.
func deepCopyAnySlice(in []any) []any {
	if in == nil {
		return nil
	}
	out := make([]any, len(in))
	for i, v := range in {
		out[i] = deepCopyAny(v)
	}
	return out
}

func deepCopyAny(v any) any {
	switch t := v.(type) {
	case map[string]any:
		m := make(map[string]any, len(t))
		for k, val := range t {
			m[k] = deepCopyAny(val)
		}
		return m
	case []any:
		return deepCopyAnySlice(t)
	case []byte:
		cp := make([]byte, len(t))
		copy(cp, t)
		return cp
	default:
		// Scalars (string, float64, bool, nil, ...) are immutable; safe to share.
		return v
	}
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
