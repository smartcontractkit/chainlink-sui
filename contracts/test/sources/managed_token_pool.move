// ================================================================
//          THIS IS A TEST CONTRACT FOR THE MANAGED TOKEN POOL
// ================================================================

module test::managed_token_pool;

use std::string::{Self, String};
use std::type_name::{Self, TypeName};
use std::option::{Self, Option};

use sui::address;
use sui::event;
use sui::object::{Self, UID};
use sui::table::{Self, Table};

// Event structures from the actual managed_token_pool contract
// Note: The actual contract uses token_pool::emit_locked_or_burned and token_pool::emit_released_or_minted
// These events are defined in the ccip_token_pool module, so we'll create test versions here

public struct TokenLockedOrBurned has copy, drop {
    amount: u64,
    remote_chain_selector: u64,
    token: address,
}

public struct TokenReleasedOrMinted has copy, drop {
    receiver: address,
    amount: u64,
    remote_chain_selector: u64,
}

// Supporting structures needed for the events
public struct ManagedTokenPoolState<phantom T> has key {
    id: UID,
    token_pool_state: TokenPoolState,
    mint_cap: MintCap<T>,
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

public struct MintCap<phantom T> has key, store {
    id: UID,
}

public struct OwnableState has store {
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

/// Emit a TokenLockedOrBurned event
public fun emit_token_locked_or_burned_event(
    amount: u64,
    remote_chain_selector: u64,
    token: address
) {
    event::emit(TokenLockedOrBurned {
        amount,
        remote_chain_selector,
        token
    });
}

/// Emit a TokenReleasedOrMinted event
public fun emit_token_released_or_minted_event(
    receiver: address,
    amount: u64,
    remote_chain_selector: u64
) {
    event::emit(TokenReleasedOrMinted {
        receiver,
        amount,
        remote_chain_selector
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
        true,        // inbound_is_enabled
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

/// Create test lock_or_burn_params
public fun create_test_lock_or_burn_params(): vector<address> {
    vector[@0x6, @0x403, @0x1, @0x2]  // Clock, DenyList, ManagedTokenState, ManagedTokenPoolState
}

/// Create test release_or_mint_params
public fun create_test_release_or_mint_params(): vector<address> {
    vector[@0x6, @0x403, @0x1, @0x2]  // Clock, DenyList, ManagedTokenState, ManagedTokenPoolState
}
