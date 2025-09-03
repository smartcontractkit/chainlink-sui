// ================================================================
//          THIS IS A TEST CONTRACT FOR THE TOKEN ADMIN REGISTRY
// ================================================================

module test::token_admin_registry;

use std::ascii;
use std::string::{Self, String};
use std::type_name;

use sui::address;
use sui::event;

// Event structures from the actual token_admin_registry contract
public struct PoolSet has copy, drop {
    coin_metadata_address: address,
    previous_pool_package_id: address,
    new_pool_package_id: address,
    // type proof of the new token pool
    token_pool_type_proof: ascii::String,
    lock_or_burn_params: vector<address>,
    release_or_mint_params: vector<address>,
}

public struct PoolRegistered has copy, drop {
    coin_metadata_address: address,
    token_pool_package_id: address,
    administrator: address,
    // type proof of the token pool
    token_pool_type_proof: ascii::String,
}

public struct PoolUnregistered has copy, drop {
    coin_metadata_address: address,
    previous_pool_address: address,
}

public struct AdministratorTransferRequested has copy, drop {
    coin_metadata_address: address,
    current_admin: address,
    new_admin: address
}

public struct AdministratorTransferred has copy, drop {
    coin_metadata_address: address,
    new_admin: address
}

// Supporting structures needed for the events
public struct TokenConfig has store, drop, copy {
    token_pool_package_id: address,
    token_pool_module: String,
    // the type of the token, this should be the full type name of the token, e.g. "0x2::token::Token<0x1::sui::SUI>"
    token_type: ascii::String,
    administrator: address,
    pending_administrator: address,
    // type proof of the token pool
    token_pool_type_proof: ascii::String,
    lock_or_burn_params: vector<address>,
    release_or_mint_params: vector<address>,
}

// Simple methods to emit events for testing purposes

/// Emit a PoolSet event
public fun emit_pool_set_event(
    coin_metadata_address: address,
    previous_pool_package_id: address,
    new_pool_package_id: address,
    token_pool_type_proof: ascii::String,
    lock_or_burn_params: vector<address>,
    release_or_mint_params: vector<address>
) {
    event::emit(PoolSet {
        coin_metadata_address,
        previous_pool_package_id,
        new_pool_package_id,
        token_pool_type_proof,
        lock_or_burn_params,
        release_or_mint_params
    });
}

/// Emit a PoolRegistered event
public fun emit_pool_registered_event(
    coin_metadata_address: address,
    token_pool_package_id: address,
    administrator: address,
    token_pool_type_proof: ascii::String
) {
    event::emit(PoolRegistered {
        coin_metadata_address,
        token_pool_package_id,
        administrator,
        token_pool_type_proof
    });
}

/// Emit a PoolUnregistered event
public fun emit_pool_unregistered_event(
    coin_metadata_address: address,
    previous_pool_address: address
) {
    event::emit(PoolUnregistered {
        coin_metadata_address,
        previous_pool_address
    });
}

/// Emit an AdministratorTransferRequested event
public fun emit_administrator_transfer_requested_event(
    coin_metadata_address: address,
    current_admin: address,
    new_admin: address
) {
    event::emit(AdministratorTransferRequested {
        coin_metadata_address,
        current_admin,
        new_admin
    });
}

/// Emit an AdministratorTransferred event
public fun emit_administrator_transferred_event(
    coin_metadata_address: address,
    new_admin: address
) {
    event::emit(AdministratorTransferred {
        coin_metadata_address,
        new_admin
    });
}

// Helper functions to create test structures

/// Create a test TokenConfig
public fun create_test_token_config(
    token_pool_package_id: address,
    token_pool_module: String,
    token_type: ascii::String,
    administrator: address,
    pending_administrator: address,
    token_pool_type_proof: ascii::String,
    lock_or_burn_params: vector<address>,
    release_or_mint_params: vector<address>
): TokenConfig {
    TokenConfig {
        token_pool_package_id,
        token_pool_module,
        token_type,
        administrator,
        pending_administrator,
        token_pool_type_proof,
        lock_or_burn_params,
        release_or_mint_params
    }
}

/// Create a default test TokenConfig with reasonable values
public fun create_default_test_token_config(): TokenConfig {
    create_test_token_config(
        @0x0,
        string::utf8(b"TestModule"),
        ascii::string(b"TestType"),
        @0x0,
        @0x0,
        ascii::string(b"TestProof"),
        vector[],
        vector[]
    )
}

public fun create_test_token_config_with_pool(
    token_pool_package_id: address,
    administrator: address
): TokenConfig {
    create_test_token_config(
        token_pool_package_id,
        string::utf8(b"TestModule"),
        ascii::string(b"TestType"),
        administrator,
        @0x0,
        ascii::string(b"TestProof"),
        vector[@0x1, @0x2],
        vector[@0x3, @0x4]
    )
}

/// Create test lock_or_burn_params
public fun create_test_lock_or_burn_params(): vector<address> {
    vector[@0x1, @0x2, @0x3]
}

/// Create test release_or_mint_params
public fun create_test_release_or_mint_params(): vector<address> {
    vector[@0x4, @0x5, @0x6]
}

/// Create a test ascii::String
public fun create_test_ascii_string(value: vector<u8>): ascii::String {
    ascii::string(value)
}

/// Create a test String
public fun create_test_string(value: vector<u8>): String {
    string::utf8(value)
}
