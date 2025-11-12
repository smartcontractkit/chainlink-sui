module managed_token_pool::rate_limiter;

use sui::clock::Clock;

public struct TokenBucket has drop, store {
    tokens: u64,
    last_updated: u64,
    is_enabled: bool,
    capacity: u64,
    rate: u64,
}

const TIME_CONVERSION_TO_SECONDS: u64 = 1000;

const ETokenMaxCapacityExceeded: u64 = 1;
const ETokenRateLimitReached: u64 = 2;

public(package) fun new(clock: &Clock, is_enabled: bool, capacity: u64, rate: u64): TokenBucket {
    TokenBucket {
        tokens: 0,
        last_updated: clock.timestamp_ms() / TIME_CONVERSION_TO_SECONDS,
        is_enabled,
        capacity,
        rate,
    }
}

public(package) fun get_current_token_bucket_state(
    clock: &Clock,
    state: &TokenBucket,
): TokenBucket {
    TokenBucket {
        tokens: calculate_refill(
            state,
            clock.timestamp_ms() / TIME_CONVERSION_TO_SECONDS - state.last_updated,
        ),
        last_updated: clock.timestamp_ms() / TIME_CONVERSION_TO_SECONDS,
        is_enabled: state.is_enabled,
        capacity: state.capacity,
        rate: state.rate,
    }
}

public(package) fun consume(clock: &Clock, bucket: &mut TokenBucket, requested_tokens: u64) {
    if (!bucket.is_enabled || requested_tokens == 0) { return };

    update_bucket(clock, bucket);

    assert!(requested_tokens <= bucket.capacity, ETokenMaxCapacityExceeded);

    assert!(requested_tokens <= bucket.tokens, ETokenRateLimitReached);

    bucket.tokens = bucket.tokens - requested_tokens;
}

/// We allow 0 rate and/or 0 capacity rate limits to effectively disable value transfer.
public(package) fun set_token_bucket_config(
    clock: &Clock,
    bucket: &mut TokenBucket,
    is_enabled: bool,
    capacity: u64,
    rate: u64,
) {
    update_bucket(clock, bucket);

    bucket.tokens = min(bucket.tokens, capacity);
    bucket.capacity = capacity;
    bucket.rate = rate;
    bucket.is_enabled = is_enabled;
}

fun update_bucket(clock: &Clock, bucket: &mut TokenBucket) {
    let time_now_seconds = clock.timestamp_ms() / TIME_CONVERSION_TO_SECONDS;
    let time_diff = time_now_seconds - bucket.last_updated;

    if (time_diff > 0) {
        bucket.tokens = calculate_refill(bucket, time_diff);
        bucket.last_updated = time_now_seconds;
    };
}

fun calculate_refill(bucket: &TokenBucket, time_diff: u64): u64 {
    if (bucket.tokens >= bucket.capacity) { return bucket.capacity };
    if (bucket.rate == 0) { return bucket.tokens };

    let remaining = bucket.capacity - bucket.tokens;
    // If time_diff * rate would exceed remaining, saturate to capacity
    if (time_diff > remaining / bucket.rate) {
        bucket.capacity
    } else {
        bucket.tokens + time_diff * bucket.rate
    }
}

fun min(a: u64, b: u64): u64 {
    if (a > b) b else a
}

public fun get_token_bucket_fields(bucket: &TokenBucket): (u64, u64, bool, u64, u64) {
    (bucket.tokens, bucket.last_updated, bucket.is_enabled, bucket.capacity, bucket.rate)
}

// ================================================================
// |                      Test-only functions                     |
// ================================================================

#[test_only]
public fun test_calculate_refill(bucket: &TokenBucket, time_diff: u64): u64 {
    calculate_refill(bucket, time_diff)
}

#[test_only]
public fun test_create_token_bucket(
    tokens: u64,
    last_updated: u64,
    is_enabled: bool,
    capacity: u64,
    rate: u64,
): TokenBucket {
    TokenBucket {
        tokens,
        last_updated,
        is_enabled,
        capacity,
        rate,
    }
}
