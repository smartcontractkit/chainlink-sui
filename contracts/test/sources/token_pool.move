// ================================================================
//          THIS IS A TEST CONTRACT FOR THE TOKEN POOL
// ================================================================

module test::token_pool;

use std::string::{Self, String};
use std::type_name::{Self, TypeName};

use sui::address;
use sui::event;
use sui::object::{Self, UID};
use sui::vec_map::{Self, VecMap};

// Event structures from the actual token_pool contract
public struct LockedOrBurned has copy, drop {
    remote_chain_selector: u64,
    local_token: address,
    amount: u64
}

public struct ReleasedOrMinted has copy, drop {
    remote_chain_selector: u64,
    local_token: address,
    recipient: address,
    amount: u64
}

public struct RemotePoolAdded has copy, drop {
    remote_chain_selector: u64,
    remote_pool_address: vector<u8>
}

public struct RemotePoolRemoved has copy, drop {
    remote_chain_selector: u64,
    remote_pool_address: vector<u8>
}

public struct ChainAdded has copy, drop {
    remote_chain_selector: u64,
    remote_token_address: vector<u8>
}

public struct LiquidityAdded has copy, drop {
    local_token: address,
    provider: address,
    amount: u64,
}

public struct LiquidityRemoved has copy, drop {
    local_token: address,
    provider: address,
    amount: u64,
}

public struct RebalancerSet has copy, drop {
    local_token: address,
    previous_rebalancer: address,
    rebalancer: address,
}

// Supporting structures needed for the events
public struct TokenPoolState has store {
    allowlist_state: AllowlistState,
    coin_metadata: address,
    local_decimals: u8,
    remote_chain_configs: VecMap<u64, RemoteChainConfig>,
    rate_limiter_config: RateLimitState
}

public struct RemoteChainConfig has store, drop, copy {
    remote_token_address: vector<u8>,
    remote_pools: vector<vector<u8>>
}

public struct AllowlistState has key, store {
    id: UID,
    enabled: bool,
    allowlist: vector<address>
}

public struct RateLimitState has key, store {
    id: UID,
    chain_rate_limiters: VecMap<u64, ChainRateLimiter>
}

public struct ChainRateLimiter has store, drop {
    outbound_is_enabled: bool,
    outbound_capacity: u64,
    outbound_rate: u64,
    outbound_current: u64,
    outbound_last_reset: u64,
    inbound_is_enabled: bool,
    inbound_capacity: u64,
    inbound_rate: u64,
    inbound_current: u64,
    inbound_last_reset: u64,
}

// Simple methods to emit events for testing purposes

/// Emit a LockedOrBurned event
public fun emit_locked_or_burned_event(
    remote_chain_selector: u64,
    local_token: address,
    amount: u64
) {
    event::emit(LockedOrBurned {
        remote_chain_selector,
        local_token,
        amount
    });
}

/// Emit a ReleasedOrMinted event
public fun emit_released_or_minted_event(
    remote_chain_selector: u64,
    local_token: address,
    recipient: address,
    amount: u64
) {
    event::emit(ReleasedOrMinted {
        remote_chain_selector,
        local_token,
        recipient,
        amount
    });
}

/// Emit a RemotePoolAdded event
public fun emit_remote_pool_added_event(
    remote_chain_selector: u64,
    remote_pool_address: vector<u8>
) {
    event::emit(RemotePoolAdded {
        remote_chain_selector,
        remote_pool_address
    });
}

/// Emit a RemotePoolRemoved event
public fun emit_remote_pool_removed_event(
    remote_chain_selector: u64,
    remote_pool_address: vector<u8>
) {
    event::emit(RemotePoolRemoved {
        remote_chain_selector,
        remote_pool_address
    });
}

/// Emit a ChainAdded event
public fun emit_chain_added_event(
    remote_chain_selector: u64,
    remote_token_address: vector<u8>
) {
    event::emit(ChainAdded {
        remote_chain_selector,
        remote_token_address
    });
}

/// Emit a LiquidityAdded event
public fun emit_liquidity_added_event(
    local_token: address,
    provider: address,
    amount: u64
) {
    event::emit(LiquidityAdded {
        local_token,
        provider,
        amount
    });
}

/// Emit a LiquidityRemoved event
public fun emit_liquidity_removed_event(
    local_token: address,
    provider: address,
    amount: u64
) {
    event::emit(LiquidityRemoved {
        local_token,
        provider,
        amount
    });
}

/// Emit a RebalancerSet event
public fun emit_rebalancer_set_event(
    local_token: address,
    previous_rebalancer: address,
    rebalancer: address
) {
    event::emit(RebalancerSet {
        local_token,
        previous_rebalancer,
        rebalancer
    });
}

// Helper functions to create test structures

/// Create a test RemoteChainConfig
public fun create_test_remote_chain_config(
    remote_token_address: vector<u8>,
    remote_pools: vector<vector<u8>>
): RemoteChainConfig {
    RemoteChainConfig {
        remote_token_address,
        remote_pools
    }
}

/// Create a default test RemoteChainConfig with reasonable values
public fun create_default_test_remote_chain_config(): RemoteChainConfig {
    create_test_remote_chain_config(
        x"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",  // remote_token_address
        vector[
            x"2234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdea",
            x"3234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdeb"
        ]  // remote_pools
    )
}

/// Create a test ChainRateLimiter
public fun create_test_chain_rate_limiter(
    outbound_is_enabled: bool,
    outbound_capacity: u64,
    outbound_rate: u64,
    inbound_is_enabled: bool,
    inbound_capacity: u64,
    inbound_rate: u64
): ChainRateLimiter {
    ChainRateLimiter {
        outbound_is_enabled,
        outbound_capacity,
        outbound_rate,
        outbound_current: 0,
        outbound_last_reset: 0,
        inbound_is_enabled,
        inbound_capacity,
        inbound_rate,
        inbound_current: 0,
        inbound_last_reset: 0
    }
}

/// Create a default test ChainRateLimiter with reasonable values
public fun create_default_test_chain_rate_limiter(): ChainRateLimiter {
    create_test_chain_rate_limiter(
        true,       // outbound_is_enabled
        1000000,    // outbound_capacity
        1000,       // outbound_rate
        true,        // inbound_is_enabled
        1000000,    // inbound_capacity
        1000        // inbound_rate
    )
}

/// Create test remote pool addresses
public fun create_test_remote_pool_addresses(): vector<vector<u8>> {
    vector[
        x"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
        x"2234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdea"
    ]
}

/// Create test remote token addresses
public fun create_test_remote_token_address(): vector<u8> {
    x"3234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdeb"
}

/// Create test chain selectors
public fun create_test_chain_selectors(): vector<u64> {
    vector[1, 137, 42161, 10]  // Ethereum, Polygon, Arbitrum, Optimism
}

/// Create test allowlist addresses
public fun create_test_allowlist(): vector<address> {
    vector[@0x1, @0x2, @0x3, @0x4]
}

/// Create test amounts for testing
public fun create_test_amounts(): vector<u64> {
    vector[1000, 5000, 10000, 50000]
}

/// Create test providers for liquidity operations
public fun create_test_providers(): vector<address> {
    vector[@0x5, @0x6, @0x7, @0x8]
}

/// Create test rebalancer addresses
public fun create_test_rebalancers(): vector<address> {
    vector[@0x9, @0xa, @0xb, @0xc]
}
