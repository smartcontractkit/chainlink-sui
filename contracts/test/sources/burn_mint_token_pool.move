// ================================================================
//          THIS IS A TEST CONTRACT FOR THE BURN MINT TOKEN POOL
// ================================================================

module test::burn_mint_token_pool;

use std::string::{Self, String};
use std::type_name::{Self, TypeName};
use std::option::{Self, Option};

use sui::address;
use sui::event;
use sui::object::{Self, UID};
use sui::table::{Self, Table};

// Event structures from the actual burn_mint_token_pool and token_pool contracts
public struct LockedOrBurned has copy, drop {
    remote_chain_selector: u64,
    local_token: address,
    amount: u64,
}

public struct ReleasedOrMinted has copy, drop {
    remote_chain_selector: u64,
    local_token: address,
    recipient: address,
    amount: u64,
}

public struct RemotePoolAdded has copy, drop {
    remote_chain_selector: u64,
    remote_pool_address: vector<u8>,
}

public struct RemotePoolRemoved has copy, drop {
    remote_chain_selector: u64,
    remote_pool_address: vector<u8>,
}

public struct ChainAdded has copy, drop {
    remote_chain_selector: u64,
    remote_token_address: vector<u8>,
}

public struct ChainRemoved has copy, drop {
    remote_chain_selector: u64,
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
public struct BurnMintTokenPoolState<phantom T> has key, store {
    id: UID,
    token_pool_state: TokenPoolState,
    treasury_cap: TreasuryCap<T>,
    ownable_state: OwnableState,
}

public struct TokenPoolState has key, store {
    id: UID,
    token: address,
    local_decimals: u8,
    remote_chain_selectors: vector<u64>,
    remote_pools: Table<u64, vector<vector<u8>>>,
    remote_tokens: Table<u64, vector<u8>>,
    allowlist_enabled: bool,
    allowlist: vector<address>,
    rate_limiters: Table<u64, RateLimiter>,
}

public struct TreasuryCap<phantom T> has key, store {
    id: UID,
}

public struct OwnableState has key, store {
    id: UID,
    owner: address,
    pending_owner: Option<address>,
    pending_transfer: Option<TransferRequest>,
}

public struct TransferRequest has store, drop {
    from: address,
    to: address,
    accepted: Option<bool>,
}

public struct RateLimiter has store, drop {
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
    amount: u64,
    remote_chain_selector: u64,
) {
    let local_token = @0x1;
    event::emit(LockedOrBurned {
        remote_chain_selector,
        local_token,
        amount
    });
}

/// Emit a ReleasedOrMinted event
public fun emit_released_or_minted_event(
    amount: u64,
    remote_chain_selector: u64,
) {
    let local_token = @0x1;
    let recipient = @0x2;
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
) {
    let remote_pool_address = x"1234567890abcdef1234567890abcdef1234567890abcdef";
    event::emit(RemotePoolAdded {
        remote_chain_selector,
        remote_pool_address
    });
}

/// Emit a RemotePoolRemoved event
public fun emit_remote_pool_removed_event(
    remote_chain_selector: u64,
) {
    let remote_pool_address = x"1234567890abcdef1234567890abcdef1234567890abcdef";
    event::emit(RemotePoolRemoved {
        remote_chain_selector,
        remote_pool_address
    });
}

/// Emit a ChainAdded event
public fun emit_chain_added_event(
    remote_chain_selector: u64,
) {
    let remote_token_address = x"2234567890abcdef1234567890abcdef1234567890abcdef";
    event::emit(ChainAdded {
        remote_chain_selector,
        remote_token_address
    });
}

/// Emit a ChainRemoved event
public fun emit_chain_removed_event(
    remote_chain_selector: u64,
) {
    event::emit(ChainRemoved {
        remote_chain_selector
    });
}

/// Emit a LiquidityAdded event
public fun emit_liquidity_added_event(
    amount: u64,
) {
    let local_token = @0x1;
    let provider = @0x3;
    event::emit(LiquidityAdded {
        local_token,
        provider,
        amount
    });
}

/// Emit a LiquidityRemoved event
public fun emit_liquidity_removed_event(
    amount: u64,
) {
    let local_token = @0x1;
    let provider = @0x3;
    event::emit(LiquidityRemoved {
        local_token,
        provider,
        amount
    });
}

/// Emit a RebalancerSet event
public fun emit_rebalancer_set_event() {
    let local_token = @0x1;
    let previous_rebalancer = @0x4;
    let rebalancer = @0x5;
    event::emit(RebalancerSet {
        local_token,
        previous_rebalancer,
        rebalancer
    });
}

// Helper functions to create test structures

/// Create a test RateLimiter
public fun create_test_rate_limiter(
    outbound_is_enabled: bool,
    outbound_capacity: u64,
    outbound_rate: u64,
    inbound_is_enabled: bool,
    inbound_capacity: u64,
    inbound_rate: u64
): RateLimiter {
    RateLimiter {
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

/// Create a default test RateLimiter with reasonable values
public fun create_default_test_rate_limiter(): RateLimiter {
    create_test_rate_limiter(
        true,       // outbound_is_enabled
        1000000,    // outbound_capacity
        1000,       // outbound_rate
        true,       // inbound_is_enabled
        1000000,    // inbound_capacity
        1000        // inbound_rate
    )
}

/// Create a test TransferRequest
public fun create_test_transfer_request(
    from: address,
    to: address,
    accepted: Option<bool>
): TransferRequest {
    TransferRequest {
        from,
        to,
        accepted
    }
}

/// Create a test TransferRequest that is pending
public fun create_pending_transfer_request(
    from: address,
    to: address
): TransferRequest {
    create_test_transfer_request(from, to, option::none<bool>())
}

/// Create a test TransferRequest that is accepted
public fun create_accepted_transfer_request(
    from: address,
    to: address
): TransferRequest {
    create_test_transfer_request(from, to, option::some(true))
}

/// Create a test TransferRequest that is rejected
public fun create_rejected_transfer_request(
    from: address,
    to: address
): TransferRequest {
    create_test_transfer_request(from, to, option::some(false))
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

/// Create test token addresses
public fun create_test_token_addresses(): vector<address> {
    vector[@0x7, @0x8, @0x9, @0xa]
}

/// Create test amounts for testing
public fun create_test_amounts(): vector<u64> {
    vector[1000, 5000, 10000, 50000]
}

