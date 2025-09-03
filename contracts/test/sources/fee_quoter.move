// ================================================================
//          THIS IS A TEST CONTRACT FOR THE FEE QUOTER
// ================================================================f

module test::fee_quoter;

use std::bcs;
use std::string::{Self, String};
use std::type_name;

use sui::address;
use sui::event;
use sui::table;

// Event structures from the actual fee_quoter contract
public struct FeeTokenAdded has copy, drop {
    fee_token: address,
}

public struct FeeTokenRemoved has copy, drop {
    fee_token: address,
}

public struct TokenTransferFeeConfigAdded has copy, drop {
    dest_chain_selector: u64,
    token: address,
    token_transfer_fee_config: TokenTransferFeeConfig,
}

public struct TokenTransferFeeConfigRemoved has copy, drop {
    dest_chain_selector: u64,
    token: address,
}

public struct UsdPerTokenUpdated has copy, drop {
    token: address,
    usd_per_token: u256,
    timestamp: u64,
}

public struct UsdPerUnitGasUpdated has copy, drop {
    dest_chain_selector: u64,
    usd_per_unit_gas: u256,
    timestamp: u64,
}

public struct DestChainAdded has copy, drop {
    dest_chain_selector: u64,
    dest_chain_config: DestChainConfig,
}

public struct DestChainConfigUpdated has copy, drop {
    dest_chain_selector: u64,
    dest_chain_config: DestChainConfig,
}

public struct PremiumMultiplierWeiPerEthUpdated has copy, drop {
    token: address,
    premium_multiplier_wei_per_eth: u64,
}

// Supporting structures needed for the events
public struct TokenTransferFeeConfig has copy, drop, store {
    min_fee_usd_cents: u32,
    max_fee_usd_cents: u32,
    deci_bps: u16,
    dest_gas_overhead: u32,
    dest_bytes_overhead: u32,
    is_enabled: bool,
}

public struct DestChainConfig has copy, drop, store {
    is_enabled: bool,
    max_number_of_tokens_per_msg: u16,
    max_data_bytes: u32,
    max_per_msg_gas_limit: u32,
    dest_gas_overhead: u32,
    dest_gas_per_payload_byte_base: u8,
    dest_gas_per_payload_byte_high: u8,
    dest_gas_per_payload_byte_threshold: u16,
    dest_data_availability_overhead_gas: u32,
    dest_gas_per_data_availability_byte: u16,
    dest_data_availability_multiplier_bps: u16,
    chain_family_selector: vector<u8>,
    enforce_out_of_order: bool,
    default_token_fee_usd_cents: u16,
    default_token_dest_gas_overhead: u32,
    default_tx_gas_limit: u32,
    gas_multiplier_wei_per_eth: u64,
    gas_price_staleness_threshold: u32,
    network_fee_usd_cents: u32
}

// Simple methods to emit events for testing purposes

/// Emit a FeeTokenAdded event
public fun emit_fee_token_added_event(fee_token: address) {
    event::emit(FeeTokenAdded { fee_token });
}

/// Emit a FeeTokenRemoved event
public fun emit_fee_token_removed_event(fee_token: address) {
    event::emit(FeeTokenRemoved { fee_token });
}

/// Emit a TokenTransferFeeConfigAdded event
public fun emit_token_transfer_fee_config_added_event(
    dest_chain_selector: u64,
    token: address,
    token_transfer_fee_config: TokenTransferFeeConfig
) {
    event::emit(TokenTransferFeeConfigAdded {
        dest_chain_selector,
        token,
        token_transfer_fee_config
    });
}

/// Emit a TokenTransferFeeConfigRemoved event
public fun emit_token_transfer_fee_config_removed_event(
    dest_chain_selector: u64,
    token: address
) {
    event::emit(TokenTransferFeeConfigRemoved {
        dest_chain_selector,
        token
    });
}

/// Emit a UsdPerTokenUpdated event
public fun emit_usd_per_token_updated_event(
    token: address,
    usd_per_token: u256,
    timestamp: u64
) {
    event::emit(UsdPerTokenUpdated {
        token,
        usd_per_token,
        timestamp
    });
}

/// Emit a UsdPerUnitGasUpdated event
public fun emit_usd_per_unit_gas_updated_event(
    dest_chain_selector: u64,
    usd_per_unit_gas: u256,
    timestamp: u64
) {
    event::emit(UsdPerUnitGasUpdated {
        dest_chain_selector,
        usd_per_unit_gas,
        timestamp
    });
}

/// Emit a DestChainAdded event
public fun emit_dest_chain_added_event(
    dest_chain_selector: u64,
    dest_chain_config: DestChainConfig
) {
    event::emit(DestChainAdded {
        dest_chain_selector,
        dest_chain_config
    });
}

/// Emit a DestChainConfigUpdated event
public fun emit_dest_chain_config_updated_event(
    dest_chain_selector: u64,
    dest_chain_config: DestChainConfig
) {
    event::emit(DestChainConfigUpdated {
        dest_chain_selector,
        dest_chain_config
    });
}

/// Emit a PremiumMultiplierWeiPerEthUpdated event
public fun emit_premium_multiplier_wei_per_eth_updated_event(
    token: address,
    premium_multiplier_wei_per_eth: u64
) {
    event::emit(PremiumMultiplierWeiPerEthUpdated {
        token,
        premium_multiplier_wei_per_eth
    });
}

// Helper functions to create test structures

/// Create a test TokenTransferFeeConfig
public fun create_test_token_transfer_fee_config(
    min_fee_usd_cents: u32,
    max_fee_usd_cents: u32,
    deci_bps: u16,
    dest_gas_overhead: u32,
    dest_bytes_overhead: u32,
    is_enabled: bool
): TokenTransferFeeConfig {
    TokenTransferFeeConfig {
        min_fee_usd_cents,
        max_fee_usd_cents,
        deci_bps,
        dest_gas_overhead,
        dest_bytes_overhead,
        is_enabled
    }
}

/// Create a test DestChainConfig
public fun create_test_dest_chain_config(
    is_enabled: bool,
    max_number_of_tokens_per_msg: u16,
    max_data_bytes: u32,
    max_per_msg_gas_limit: u32,
    dest_gas_overhead: u32,
    dest_gas_per_payload_byte_base: u8,
    dest_gas_per_payload_byte_high: u8,
    dest_gas_per_payload_byte_threshold: u16,
    dest_data_availability_overhead_gas: u32,
    dest_gas_per_data_availability_byte: u16,
    dest_data_availability_multiplier_bps: u16,
    chain_family_selector: vector<u8>,
    enforce_out_of_order: bool,
    default_token_fee_usd_cents: u16,
    default_token_dest_gas_overhead: u32,
    default_tx_gas_limit: u32,
    gas_multiplier_wei_per_eth: u64,
    gas_price_staleness_threshold: u32,
    network_fee_usd_cents: u32
): DestChainConfig {
    DestChainConfig {
        is_enabled,
        max_number_of_tokens_per_msg,
        max_data_bytes,
        max_per_msg_gas_limit,
        dest_gas_overhead,
        dest_gas_per_payload_byte_base,
        dest_gas_per_payload_byte_high,
        dest_gas_per_payload_byte_threshold,
        dest_data_availability_overhead_gas,
        dest_gas_per_data_availability_byte,
        dest_data_availability_multiplier_bps,
        chain_family_selector,
        enforce_out_of_order,
        default_token_fee_usd_cents,
        default_token_dest_gas_overhead,
        default_tx_gas_limit,
        gas_multiplier_wei_per_eth,
        gas_price_staleness_threshold,
        network_fee_usd_cents
    }
}

public fun create_default_test_dest_chain_config(): DestChainConfig {
    create_test_dest_chain_config(
        true,
        10,
        10000,
        1000000,
        100000,
        10,
        20,
        1000,
        50000,
        100,
        1000,
        b"EVM",
        false,
        50,
        50000,
        500000,
        1000000000000000000,
        3600,
        100
    )
}

public fun create_default_test_token_transfer_fee_config(): TokenTransferFeeConfig {
    create_test_token_transfer_fee_config(
        25,
        100,
        50,
        25000,
        32,
        true
    )
}
