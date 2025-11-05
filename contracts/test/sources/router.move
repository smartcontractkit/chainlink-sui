module test::router;

use std::string::{Self, String};
use std::type_name::{Self, TypeName};

use sui::address;
use sui::event;
use sui::object::{Self, UID};

// Event structures from the actual router contract
public struct OnRampSet has copy, drop {
    dest_chain_selector: u64,
    on_ramp_info: OnRampInfo,
}

// Supporting structures needed for the events
public struct RouterState has key {
    id: UID,
    on_ramp_infos: vector<OnRampInfo>,
}

public struct OnRampInfo has copy, store, drop {
    onramp_address: address,
    onramp_version: vector<u8>,
}

// Simple methods to emit events for testing purposes

/// Emit an OnRampSet event
public fun emit_on_ramp_set_event(
    dest_chain_selector: u64,
    on_ramp_info: OnRampInfo
) {
    event::emit(OnRampSet {
        dest_chain_selector,
        on_ramp_info
    });
}

// Helper functions to create test structures

/// Create a test OnRampInfo
public fun create_test_on_ramp_info(
    onramp_address: address,
    onramp_version: vector<u8>
): OnRampInfo {
    OnRampInfo {
        onramp_address,
        onramp_version
    }
}

/// Create a default test OnRampInfo with reasonable values
public fun create_default_test_on_ramp_info(): OnRampInfo {
    create_test_on_ramp_info(
        @0x1,                           // onramp_address
        vector[1, 6, 0]                 // onramp_version (1.6.0)
    )
}

/// Create a test OnRampInfo for a specific version
public fun create_test_on_ramp_info_with_version(
    onramp_address: address,
    major: u8,
    minor: u8,
    patch: u8
): OnRampInfo {
    create_test_on_ramp_info(
        onramp_address,
        vector[major, minor, patch]
    )
}

/// Create a test OnRampInfo for a specific address
public fun get_mock_onramp_address(): address {
    return @0x1
}

/// Create test chain selectors
public fun create_test_chain_selectors(): vector<u64> {
    vector[1, 137, 42161, 10]  // Ethereum, Polygon, Arbitrum, Optimism
}

/// Create test onramp addresses
public fun create_test_onramp_addresses(): vector<address> {
    vector[@0x1, @0x2, @0x3, @0x4]
}

/// Create test onramp versions
public fun create_test_onramp_versions(): vector<vector<u8>> {
    vector[
        vector[1, 6, 0],  // 1.6.0
        vector[1, 5, 0],  // 1.5.0
        vector[1, 4, 0],  // 1.4.0
        vector[1, 3, 0]   // 1.3.0
    ]
}
