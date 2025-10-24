module test::onramp;

use std::string::{Self, String};
use std::type_name::{Self, TypeName};
use std::u256;

use sui::address;
use sui::event;
use sui::object::{Self, UID};

// Event structures from the actual onramp contract
public struct ConfigSet has copy, drop {
    static_config: StaticConfig,
    dynamic_config: DynamicConfig
}

public struct DestChainConfigSet has copy, drop {
    dest_chain_selector: u64,
    is_enabled: bool,
    sequence_number: u64,
    allowlist_enabled: bool
}

public struct CCIPMessageSent has copy, drop {
    dest_chain_selector: u64,
    sequence_number: u64,
    message: Sui2AnyRampMessage
}

public struct AllowlistSendersAdded has copy, drop {
    dest_chain_selector: u64,
    senders: vector<address>
}

public struct AllowlistSendersRemoved has copy, drop {
    dest_chain_selector: u64,
    senders: vector<address>
}

public struct FeeTokenWithdrawn has copy, drop {
    fee_aggregator: address,
    fee_token: address,
    amount: u64
}

// Supporting structures needed for the events
public struct OnRampState has key, store {
    id: UID,
    chain_selector: u64,
    fee_aggregator: address,
    allowlist_admin: address,
}

public struct StaticConfig has copy, drop {
    chain_selector: u64
}

public struct DynamicConfig has copy, drop {
    fee_aggregator: address,
    allowlist_admin: address
}

public struct DestChainConfig has store, drop {
    is_enabled: bool,
    sequence_number: u64,
    allowlist_enabled: bool,
    allowed_senders: vector<address>
}

public struct RampMessageHeader has store, drop, copy {
    message_id: vector<u8>,
    source_chain_selector: u64,
    dest_chain_selector: u64,
    sequence_number: u64,
    nonce: u64
}

public struct Sui2AnyRampMessage has store, drop, copy {
    header: RampMessageHeader,
    sender: address,
    data: vector<u8>,
    receiver: vector<u8>,
    extra_args: vector<u8>,
    fee_token: address,
    fee_token_amount: u64,
    fee_value_juels: u256,
    token_amounts: vector<Sui2AnyTokenTransfer>
}

public struct Sui2AnyTokenTransfer has store, drop, copy {
    source_pool_address: address,
    dest_token_address: vector<u8>,
    extra_data: vector<u8>,
    amount: u64,
    dest_exec_data: vector<u8>
}

// Simple methods to emit events for testing purposes

/// Emit a ConfigSet event
public fun emit_config_set_event(
    static_config: StaticConfig,
    dynamic_config: DynamicConfig
) {
    event::emit(ConfigSet {
        static_config,
        dynamic_config
    });
}

/// Emit a DestChainConfigSet event
public fun emit_dest_chain_config_set_event(
    dest_chain_selector: u64,
    is_enabled: bool,
    sequence_number: u64,
    allowlist_enabled: bool
) {
    event::emit(DestChainConfigSet {
        dest_chain_selector,
        is_enabled,
        sequence_number,
        allowlist_enabled
    });
}

/// Emit a CCIPMessageSent event
public fun emit_ccip_message_sent_event(
    dest_chain_selector: u64,
    sequence_number: u64,
    message: Sui2AnyRampMessage
) {
    event::emit(CCIPMessageSent {
        dest_chain_selector,
        sequence_number,
        message
    });
}

/// Emit a CCIPMessageSent event
public fun emit_sample_ccip_message_sent_event() {
    let header = RampMessageHeader {
        message_id: x"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
        source_chain_selector: 1,
        dest_chain_selector: 137,
        sequence_number: 1,
        nonce: 1,
    };
    
    let token_transfer = Sui2AnyTokenTransfer {
        source_pool_address: @0x123,
        dest_token_address: x"456789abcdef",
        extra_data: x"deadbeef",
        amount: 1000,
        dest_exec_data: x"cafebabe"
    };
    
    let message = Sui2AnyRampMessage {
        header,
        sender: @0x789,
        data: x"abcdef123456",
        receiver: x"abcdef123456",
        extra_args: x"abcdef123456",
        fee_token: @0xabc,
        fee_token_amount: 500,
        fee_value_juels: 1000000,
        token_amounts: vector[token_transfer]
    };
    
    event::emit(CCIPMessageSent {
        dest_chain_selector: 137,
        sequence_number: 1,
        message
    });
}

/// Emit an AllowlistSendersAdded event
public fun emit_allowlist_senders_added_event(
    dest_chain_selector: u64,
    senders: vector<address>
) {
    event::emit(AllowlistSendersAdded {
        dest_chain_selector,
        senders
    });
}

/// Emit an AllowlistSendersRemoved event
public fun emit_allowlist_senders_removed_event(
    dest_chain_selector: u64,
    senders: vector<address>
) {
    event::emit(AllowlistSendersRemoved {
        dest_chain_selector,
        senders
    });
}

/// Emit a FeeTokenWithdrawn event
public fun emit_fee_token_withdrawn_event(
    fee_aggregator: address,
    fee_token: address,
    amount: u64
) {
    event::emit(FeeTokenWithdrawn {
        fee_aggregator,
        fee_token,
        amount
    });
}

// Helper functions to create test structures

/// Create a test StaticConfig
public fun create_test_static_config(chain_selector: u64): StaticConfig {
    StaticConfig { chain_selector }
}

/// Create a test DynamicConfig
public fun create_test_dynamic_config(
    fee_aggregator: address,
    allowlist_admin: address
): DynamicConfig {
    DynamicConfig {
        fee_aggregator,
        allowlist_admin
    }
}

/// Create a default test DynamicConfig with reasonable values
public fun create_default_test_dynamic_config(): DynamicConfig {
    create_test_dynamic_config(
        @0x1,                           // fee_aggregator
        @0x2                            // allowlist_admin
    )
}

/// Create a test DestChainConfig
public fun create_test_dest_chain_config(
    is_enabled: bool,
    sequence_number: u64,
    allowlist_enabled: bool,
    allowed_senders: vector<address>
): DestChainConfig {
    DestChainConfig {
        is_enabled,
        sequence_number,
        allowlist_enabled,
        allowed_senders
    }
}

/// Create a default test DestChainConfig with reasonable values
public fun create_default_test_dest_chain_config(): DestChainConfig {
    create_test_dest_chain_config(
        true,                           // is_enabled
        1,                              // sequence_number
        false,                          // allowlist_enabled
        vector[]                        // allowed_senders
    )
}

/// Create a test RampMessageHeader
public fun create_test_ramp_message_header(
    message_id: vector<u8>,
    source_chain_selector: u64,
    dest_chain_selector: u64,
    sequence_number: u64,
    nonce: u64
): RampMessageHeader {
    RampMessageHeader {
        message_id,
        source_chain_selector,
        dest_chain_selector,
        sequence_number,
        nonce
    }
}

/// Create a test Sui2AnyTokenTransfer
public fun create_test_sui2any_token_transfer(
    source_pool_address: address,
    dest_token_address: vector<u8>,
    extra_data: vector<u8>,
    amount: u64,
    dest_exec_data: vector<u8>
): Sui2AnyTokenTransfer {
    Sui2AnyTokenTransfer {
        source_pool_address,
        dest_token_address,
        extra_data,
        amount,
        dest_exec_data
    }
}

/// Create a test Sui2AnyRampMessage
public fun create_test_sui2any_ramp_message(
    header: RampMessageHeader,
    sender: address,
    data: vector<u8>,
    receiver: vector<u8>,
    extra_args: vector<u8>,
    fee_token: address,
    fee_token_amount: u64,
    fee_value_juels: u256,
    token_amounts: vector<Sui2AnyTokenTransfer>
): Sui2AnyRampMessage {
    Sui2AnyRampMessage {
        header,
        sender,
        data,
        receiver,
        extra_args,
        fee_token,
        fee_token_amount,
        fee_value_juels,
        token_amounts
    }
}

/// Create test chain selectors
public fun create_test_chain_selectors(): vector<u64> {
    vector[1, 137, 42161, 10]  // Ethereum, Polygon, Arbitrum, Optimism
}

/// Create test message IDs
public fun create_test_message_ids(): vector<vector<u8>> {
    vector[
        x"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
        x"2234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdea"
    ]
}

/// Create test senders for allowlist operations
public fun create_test_senders(): vector<address> {
    vector[@0x3, @0x4, @0x5, @0x6]
}

/// Create test fee token addresses
public fun create_test_fee_tokens(): vector<address> {
    vector[@0x7, @0x8, @0x9, @0xa]
}

/// Create test amounts for testing
public fun create_test_amounts(): vector<u64> {
    vector[1000, 5000, 10000, 50000]
}
