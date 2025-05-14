/// This module is responsible for storage and retrieval of fee token and token transfer
/// information and pricing.
module ccip::fee_quoter;

use std::string::{Self, String};
use sui::clock;
use sui::event;
use sui::table;

use ccip::eth_abi;
use ccip::state_object::{Self, CCIPObjectRef, OwnerCap};

const FEE_QUOTER_STATE_NAME: vector<u8> = b"FeeQuoterState";
const CHAIN_FAMILY_SELECTOR_EVM: vector<u8> = x"2812d52c";
const CHAIN_FAMILY_SELECTOR_SVM: vector<u8> = x"1e10bdc4";
const CHAIN_FAMILY_SELECTOR_APTOS: vector<u8> = x"ac77ffec";
const CHAIN_FAMILY_SELECTOR_SUI: vector<u8> = x"c4e05953";
const EVM_PRECOMPILE_SPACE: u256 = 1024;
const SVM_EXTRA_ARGS_V1_TAG: vector<u8> = x"1f3b3aba";
const GENERIC_EXTRA_ARGS_V2_TAG: vector<u8> = x"181dcf10";

const GAS_PRICE_BITS: u8 = 112;

const MESSAGE_FIXED_BYTES: u64 = 32 * 15;
const MESSAGE_FIXED_BYTES_PER_TOKEN: u64 = 32 * (4 + (3 + 2));

const CCIP_LOCK_OR_BURN_V1_RET_BYTES: u32 = 32;

const MAX_U64: u256 = 18446744073709551615;
const MAX_U160: u256 = 1461501637330902918203684832716283019655932542975;
const MAX_U256: u256 =
    115792089237316195423570985008687907853269984665640564039457584007913129639935;
const VAL_1E5: u256 = 100_000;
const VAL_1E14: u256 = 100_000_000_000_000;
const VAL_1E16: u256 = 10_000_000_000_000_000;
const VAL_1E18: u256 = 1_000_000_000_000_000_000;

public struct FeeQuoterState has key, store {
    id: UID,
    max_fee_juels_per_msg: u64,
    link_token: address,
    token_price_staleness_threshold: u64,
    fee_tokens: vector<address>,
    usd_per_unit_gas_by_dest_chain: table::Table<u64, TimestampedPrice>,
    // TODO: we need to know the token address for common tokens like USDC, SUI, etc
    usd_per_token: table::Table<address, TimestampedPrice>,
    dest_chain_configs: table::Table<u64, DestChainConfig>,
    // dest chain selector -> local token -> TokenTransferFeeConfig
    token_transfer_fee_configs: table::Table<u64, table::Table<address, TokenTransferFeeConfig>>,
    // TODO: update calculations - should this be octa per apt?
    premium_multiplier_wei_per_eth: table::Table<address, u64>,
}

public struct FeeQuoterCap has key, store {
    id: UID,
}

public struct StaticConfig has drop {
    max_fee_juels_per_msg: u64,
    link_token: address,
    token_price_staleness_threshold: u64
}

public struct DestChainConfig has store, drop, copy {
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
    // TODO: should this be octa per apt?
    gas_multiplier_wei_per_eth: u64,
    gas_price_staleness_threshold: u32,
    network_fee_usd_cents: u32
}

public struct TokenTransferFeeConfig has store, drop, copy {
    min_fee_usd_cents: u32,
    max_fee_usd_cents: u32,
    deci_bps: u16,
    dest_gas_overhead: u32,
    dest_bytes_overhead: u32,
    is_enabled: bool
}

public struct TimestampedPrice has store, drop, copy {
    value: u256,
    timestamp: u64
}

public struct FeeTokenAdded has copy, drop {
    fee_token: address
}

public struct FeeTokenRemoved has copy, drop {
    fee_token: address
}

public struct TokenTransferFeeConfigAdded has copy, drop {
    dest_chain_selector: u64,
    token: address,
    token_transfer_fee_config: TokenTransferFeeConfig
}

public struct TokenTransferFeeConfigRemoved has copy, drop {
    dest_chain_selector: u64,
    token: address
}

public struct UsdPerTokenUpdated has copy, drop {
    token: address,
    usd_per_token: u256,
    timestamp: u64
}

public struct UsdPerUnitGasUpdated has copy, drop {
    dest_chain_selector: u64,
    usd_per_unit_gas: u256,
    timestamp: u64
}

public struct DestChainAdded has copy, drop {
    dest_chain_selector: u64,
    dest_chain_config: DestChainConfig
}

public struct DestChainConfigUpdated has copy, drop {
    dest_chain_selector: u64,
    dest_chain_config: DestChainConfig
}

public struct PremiumMultiplierWeiPerEthUpdated has copy, drop {
    token: address,
    premium_multiplier_wei_per_eth: u64
}

const E_ALREADY_INITIALIZED: u64 = 1;
const E_OUT_OF_BOUND: u64 = 2;
const E_UNKNOWN_DEST_CHAIN_SELECTOR: u64 = 3;
const E_UNKNOWN_TOKEN: u64 = 4;
const E_DEST_CHAIN_NOT_ENABLED: u64 = 5;
const E_TOKEN_UPDATE_MISMATCH: u64 = 6;
const E_GAS_UPDATE_MISMATCH: u64 = 7;
const E_TOKEN_TRANSFER_FEE_CONFIG_MISMATCH: u64 = 8;
const E_FEE_TOKEN_NOT_SUPPORTED: u64 = 9;
const E_TOKEN_NOT_SUPPORTED: u64 = 10;
const E_UNKNOWN_CHAIN_FAMILY_SELECTOR: u64 = 11;
const E_STALE_GAS_PRICE: u64 = 12;
const E_MESSAGE_TOO_LARGE: u64 = 13;
const E_UNSUPPORTED_NUMBER_OF_TOKENS: u64 = 14;
const E_INVALID_EVM_ADDRESS: u64 = 15;
const E_INVALID_SVM_ADDRESS: u64 = 16;
const E_FEE_TOKEN_COST_TOO_HIGH: u64 = 17;
const E_MESSAGE_GAS_LIMIT_TOO_HIGH: u64 = 18;
const E_EXTRA_ARG_OUT_OF_ORDER_EXECUTION_MUST_BE_TRUE: u64 = 19;
const E_INVALID_EXTRA_ARGS_TAG: u64 = 20;
const E_INVALID_EXTRA_ARGS_DATA: u64 = 21;
const E_INVALID_TOKEN_RECEIVER: u64 = 22;
const E_MESSAGE_COMPUTE_UNIT_LIMIT_TOO_HIGH: u64 = 23;
const E_MESSAGE_FEE_TOO_HIGH: u64 = 24;
const E_SOURCE_TOKEN_DATA_TOO_LARGE: u64 = 25;
const E_INVALID_DEST_CHAIN_SELECTOR: u64 = 26;
const E_INVALID_GAS_LIMIT: u64 = 27;
const E_INVALID_CHAIN_FAMILY_SELECTOR: u64 = 28;
const E_TO_TOKEN_AMOUNT_TOO_LARGE: u64 = 29;

public fun type_and_version(): String {
    string::utf8(b"FeeQuoter 1.6.0")
}

#[allow(lint(self_transfer))]
public fun initialize(
    ref: &mut CCIPObjectRef,
    _: &OwnerCap,
    max_fee_juels_per_msg: u64,
    link_token: address, // can pass in the LINK metadata object to make sure this is a valid token address
    token_price_staleness_threshold: u64,
    fee_tokens: vector<address>,
    ctx: &mut TxContext
) {
    assert!(
        !state_object::contains(ref, FEE_QUOTER_STATE_NAME),
        E_ALREADY_INITIALIZED
    );

    let state = FeeQuoterState {
        id: object::new(ctx),
        max_fee_juels_per_msg,
        link_token,
        token_price_staleness_threshold,
        fee_tokens,
        usd_per_unit_gas_by_dest_chain: table::new(ctx),
        usd_per_token: table::new(ctx),
        dest_chain_configs: table::new(ctx),
        token_transfer_fee_configs: table::new(ctx),
        premium_multiplier_wei_per_eth: table::new(ctx),
    };
    let fee_quoter_cap = FeeQuoterCap {
        id: object::new(ctx),
    };
    transfer::transfer(
        fee_quoter_cap,
        ctx.sender()
    );
    state_object::add(ref, FEE_QUOTER_STATE_NAME, state, ctx);
}

public fun apply_fee_token_updates(
    ref: &mut CCIPObjectRef,
    _: &OwnerCap,
    fee_tokens_to_remove: vector<address>,
    fee_tokens_to_add: vector<address>,
) {
    let state = state_object::borrow_mut<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);

    // Remove tokens
    vector::do_ref!(
        &fee_tokens_to_remove,
        |fee_token| {
            let fee_token = *fee_token;
            let (found, index) = vector::index_of(&state.fee_tokens, &fee_token);
            if (found) {
                state.fee_tokens.remove(index);
                event::emit(FeeTokenRemoved { fee_token });
            };
        }
    );

    // Add new tokens
    vector::do_ref!(
        &fee_tokens_to_add,
        |fee_token| {
            let fee_token = *fee_token;
            let (found, _) = vector::index_of(&state.fee_tokens, &fee_token);
            if (!found) {
                state.fee_tokens.push_back(fee_token);
                event::emit(FeeTokenAdded { fee_token });
            };
        }
    );
}

// Note that unlike EVM, this only allows changes for a single dest chain selector at a time.
public fun apply_token_transfer_fee_config_updates(
    ref: &mut CCIPObjectRef,
    _: &OwnerCap,
    dest_chain_selector: u64,
    add_tokens: vector<address>,
    add_min_fee_usd_cents: vector<u32>,
    add_max_fee_usd_cents: vector<u32>,
    add_deci_bps: vector<u16>,
    add_dest_gas_overhead: vector<u32>,
    add_dest_bytes_overhead: vector<u32>,
    add_is_enabled: vector<bool>,
    remove_tokens: vector<address>,
    ctx: &mut TxContext
) {
    let state = state_object::borrow_mut<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);

    if (!table::contains(
        &state.token_transfer_fee_configs, dest_chain_selector
    )) {
        table::add(
            &mut state.token_transfer_fee_configs,
            dest_chain_selector,
            table::new(ctx)
        );
    };
    let token_transfer_fee_configs =
        table::borrow_mut(
            &mut state.token_transfer_fee_configs, dest_chain_selector
        );

    let add_tokens_len = add_tokens.length();
    assert!(
        add_tokens_len == add_min_fee_usd_cents.length(),
        E_TOKEN_TRANSFER_FEE_CONFIG_MISMATCH
    );
    assert!(
        add_tokens_len == add_max_fee_usd_cents.length(),
        E_TOKEN_TRANSFER_FEE_CONFIG_MISMATCH
    );
    assert!(
        add_tokens_len == add_deci_bps.length(),
        E_TOKEN_TRANSFER_FEE_CONFIG_MISMATCH
    );
    assert!(
        add_tokens_len == add_dest_gas_overhead.length(),
        E_TOKEN_TRANSFER_FEE_CONFIG_MISMATCH
    );
    assert!(
        add_tokens_len == add_dest_bytes_overhead.length(),
        E_TOKEN_TRANSFER_FEE_CONFIG_MISMATCH
    );
    assert!(
        add_tokens_len == add_is_enabled.length(),
        E_TOKEN_TRANSFER_FEE_CONFIG_MISMATCH
    );

    let mut i = 0;
    while (i < add_tokens_len) {
        let token = add_tokens[i];
        let min_fee_usd_cents = add_min_fee_usd_cents[i];
        let max_fee_usd_cents = add_max_fee_usd_cents[i];
        let deci_bps = add_deci_bps[i];
        let dest_gas_overhead = add_dest_gas_overhead[i];
        let dest_bytes_overhead = add_dest_bytes_overhead[i];
        let is_enabled = add_is_enabled[i];

        let token_transfer_fee_config = TokenTransferFeeConfig {
            min_fee_usd_cents,
            max_fee_usd_cents,
            deci_bps,
            dest_gas_overhead,
            dest_bytes_overhead,
            is_enabled
        };

        table::add(
            token_transfer_fee_configs, token, token_transfer_fee_config
        );

        event::emit(
            TokenTransferFeeConfigAdded {
                dest_chain_selector,
                token,
                token_transfer_fee_config
            }
        );

        i = i + 1;
    };

    vector::do_ref!(
        &remove_tokens,
        |token| {
            let token = *token;
            if (token_transfer_fee_configs.contains(token)) {
                token_transfer_fee_configs.remove(token);

                event::emit(
                    TokenTransferFeeConfigRemoved { dest_chain_selector, token }
                );
            }
        }
    );
}

public fun apply_dest_chain_config_updates(
    ref: &mut CCIPObjectRef,
    _: &OwnerCap,
    dest_chain_selector: u64,
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
    network_fee_usd_cents: u32,
) {
    let state = state_object::borrow_mut<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);

    assert!(
        dest_chain_selector != 0,
        E_INVALID_DEST_CHAIN_SELECTOR
    );
    assert!(
        default_tx_gas_limit != 0 && default_tx_gas_limit <= max_per_msg_gas_limit,
        E_INVALID_GAS_LIMIT
    );

    assert!(
        chain_family_selector == CHAIN_FAMILY_SELECTOR_EVM
            || chain_family_selector == CHAIN_FAMILY_SELECTOR_SVM
            || chain_family_selector == CHAIN_FAMILY_SELECTOR_APTOS
            || chain_family_selector == CHAIN_FAMILY_SELECTOR_SUI,
        E_INVALID_CHAIN_FAMILY_SELECTOR
    );

    let dest_chain_config = DestChainConfig {
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
    };

    if (table::contains(&state.dest_chain_configs, dest_chain_selector)) {
        let dest_chain_config_ref =
            table::borrow_mut(
                &mut state.dest_chain_configs, dest_chain_selector
            );
        *dest_chain_config_ref = dest_chain_config;
        event::emit(DestChainConfigUpdated { dest_chain_selector, dest_chain_config });
    } else {
        table::add(
            &mut state.dest_chain_configs, dest_chain_selector, dest_chain_config
        );
        event::emit(DestChainAdded { dest_chain_selector, dest_chain_config });
    }
}

public fun apply_premium_multiplier_wei_per_eth_updates(
    ref: &mut CCIPObjectRef,
    _: &OwnerCap,
    tokens: vector<address>,
    premium_multiplier_wei_per_eth: vector<u64>,
) {
    let state = state_object::borrow_mut<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);

    vector::zip_do_ref!(
        &tokens,
        &premium_multiplier_wei_per_eth,
        |token, premium_multiplier_wei_per_eth| {
            let token: address = *token;
            let premium_multiplier_wei_per_eth: u64 = *premium_multiplier_wei_per_eth;

            if (state.premium_multiplier_wei_per_eth.contains(token)) {
                state.premium_multiplier_wei_per_eth.remove(token);
            };
            state.premium_multiplier_wei_per_eth.add(token, premium_multiplier_wei_per_eth);

            event::emit(
                PremiumMultiplierWeiPerEthUpdated {
                    token,
                    premium_multiplier_wei_per_eth
                }
            );
        }
    );
}

public fun get_static_config(ref: &CCIPObjectRef): StaticConfig {
    let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
    StaticConfig {
        max_fee_juels_per_msg: state.max_fee_juels_per_msg,
        link_token: state.link_token,
        token_price_staleness_threshold: state.token_price_staleness_threshold
    }
}

public fun get_token_transfer_fee_config(
    ref: &CCIPObjectRef,
    dest_chain_selector: u64,
    token: address
): TokenTransferFeeConfig {
    let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
    *get_token_transfer_fee_config_internal(
        state, dest_chain_selector, token
    )
}

public fun get_token_price(ref: &CCIPObjectRef, token: address): TimestampedPrice {
    let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
    get_token_price_internal(state, token)
}

public fun get_token_prices(
    ref: &CCIPObjectRef,
    tokens: vector<address>
): (vector<TimestampedPrice>) {
    let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
    tokens.map_ref!(|token| get_token_price_internal(state, *token))
}

public fun get_dest_chain_gas_price(
    ref: &CCIPObjectRef,
    dest_chain_selector: u64
): TimestampedPrice {
    let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
    get_dest_chain_gas_price_internal(state, dest_chain_selector)
}

public fun get_token_and_gas_prices(
    ref: &CCIPObjectRef, clock: &clock::Clock, token: address, dest_chain_selector: u64
): (u256, u256) {
    let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
    let dest_chain_config = get_dest_chain_config_internal(
        state, dest_chain_selector
    );
    assert!(
        dest_chain_config.is_enabled,
        E_DEST_CHAIN_NOT_ENABLED
    );
    let token_price = get_token_price_internal(state, token);
    let gas_price_value =
        get_validated_gas_price_internal(
            state, clock, dest_chain_config, dest_chain_selector
        );
    (token_price.value, gas_price_value)
}

public fun get_dest_chain_config(
    ref: &CCIPObjectRef, dest_chain_selector: u64
): DestChainConfig {
    let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
    *get_dest_chain_config_internal(state, dest_chain_selector)
}

fun get_dest_chain_config_internal(
    state: &FeeQuoterState, dest_chain_selector: u64
): &DestChainConfig {
    assert!(
        table::contains(&state.dest_chain_configs, dest_chain_selector),
        E_UNKNOWN_DEST_CHAIN_SELECTOR
    );
    table::borrow(&state.dest_chain_configs, dest_chain_selector)
}

fun get_dest_chain_gas_price_internal(
    state: &FeeQuoterState, dest_chain_selector: u64
): TimestampedPrice {
    assert!(
        table::contains(
            &state.usd_per_unit_gas_by_dest_chain, dest_chain_selector
        ),
        E_UNKNOWN_DEST_CHAIN_SELECTOR
    );
    *table::borrow(
        &state.usd_per_unit_gas_by_dest_chain, dest_chain_selector
    )
}

fun get_validated_gas_price_internal(
    state: &FeeQuoterState,
    clock: &clock::Clock,
    dest_chain_config: &DestChainConfig,
    dest_chain_selector: u64
): u256 {
    let gas_price = get_dest_chain_gas_price_internal(state, dest_chain_selector);
    if (dest_chain_config.gas_price_staleness_threshold > 0) {
        let time_passed_secs = clock::timestamp_ms(clock) / 1000 - gas_price.timestamp;
        assert!(
            time_passed_secs <= (dest_chain_config.gas_price_staleness_threshold as u64),
            E_STALE_GAS_PRICE
        );
    };
    gas_price.value
}

// Token prices can be stale. On EVM we have additional fallbacks to a price feed, if configured.
// Since these fallbacks don't exist on Sui, we simply return the price as is.
fun get_token_price_internal(
    state: &FeeQuoterState, token: address
): TimestampedPrice {
    assert!(
        table::contains(&state.usd_per_token, token),
        E_UNKNOWN_TOKEN
    );
    *table::borrow(&state.usd_per_token, token)
}

fun convert_token_amount_internal(
    state: &FeeQuoterState,
    from_token: address,
    from_token_amount: u64,
    to_token: address
): u64 {
    let from_token_price = get_token_price_internal(state, from_token);
    let to_token_price = get_token_price_internal(state, to_token);

    let to_token_amount =
        ((from_token_amount as u256) * from_token_price.value) / to_token_price.value;
    assert!(
        to_token_amount <= MAX_U64,
        E_TO_TOKEN_AMOUNT_TOO_LARGE
    );
    to_token_amount as u64
}

// this should only be called from offramp, hence gated by a special cap owned by offramp
public fun update_prices(
    ref: &mut CCIPObjectRef,
    _: &FeeQuoterCap,
    clock: &clock::Clock,
    source_tokens: vector<address>,
    source_usd_per_token: vector<u256>,
    gas_dest_chain_selectors: vector<u64>,
    gas_usd_per_unit_gas: vector<u256>,
) {
    assert!(
        source_tokens.length() == source_usd_per_token.length(),
        E_TOKEN_UPDATE_MISMATCH
    );
    assert!(
        gas_dest_chain_selectors.length() == gas_usd_per_unit_gas.length(),
        E_GAS_UPDATE_MISMATCH
    );

    let state = state_object::borrow_mut<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
    let timestamp = clock.timestamp_ms() / 1000;

    vector::zip_do_ref!(
        &source_tokens,
        &source_usd_per_token,
        |token, usd_per_token| {
            let timestamped_price = TimestampedPrice {
                value: *usd_per_token,
                timestamp
            };

            if (state.usd_per_token.contains(*token)) {
                let _old_value = state.usd_per_token.remove(*token);
            };
            state.usd_per_token.add(*token, timestamped_price);

            event::emit(
                UsdPerTokenUpdated {
                    token: *token,
                    usd_per_token: *usd_per_token,
                    timestamp
                }
            );
        }
    );

    vector::zip_do_ref!(
        &gas_dest_chain_selectors,
        &gas_usd_per_unit_gas,
        |dest_chain_selector, usd_per_unit_gas| {
            let timestamped_price = TimestampedPrice {
                value: *usd_per_unit_gas,
                timestamp
            };

            if (state.usd_per_unit_gas_by_dest_chain.contains(*dest_chain_selector)) {
                let _old_value = state.usd_per_unit_gas_by_dest_chain.remove(*dest_chain_selector);
            };
            state.usd_per_unit_gas_by_dest_chain.add(*dest_chain_selector, timestamped_price);

            event::emit(
                UsdPerUnitGasUpdated {
                    dest_chain_selector: *dest_chain_selector,
                    usd_per_unit_gas: *usd_per_unit_gas,
                    timestamp
                }
            );
        }
    );
}

public fun get_validated_fee(
    ref: &CCIPObjectRef,
    clock: &clock::Clock,
    dest_chain_selector: u64,
    receiver: vector<u8>,
    data: vector<u8>,
    local_token_addresses: vector<address>,
    local_token_amounts: vector<u64>,
    fee_token: address,
    extra_args: vector<u8>
): u64 {
    let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);

    let dest_chain_config = get_dest_chain_config_internal(
        state, dest_chain_selector
    );
    assert!(
        dest_chain_config.is_enabled,
        E_DEST_CHAIN_NOT_ENABLED
    );

    assert!(
        vector::contains(&state.fee_tokens, &fee_token),
        E_FEE_TOKEN_NOT_SUPPORTED
    );

    let chain_family_selector = dest_chain_config.chain_family_selector;

    let data_len = data.length();
    let tokens_len = local_token_addresses.length();
    validate_message(dest_chain_config, data_len, tokens_len);

    let gas_limit =
        if (chain_family_selector == CHAIN_FAMILY_SELECTOR_EVM) {
            validate_evm_address(receiver);
            resolve_generic_gas_limit(dest_chain_config, extra_args)
        } else if (chain_family_selector == CHAIN_FAMILY_SELECTOR_APTOS
            || chain_family_selector == CHAIN_FAMILY_SELECTOR_SUI) {
            validate_32byte_address(receiver, true);
            resolve_generic_gas_limit(dest_chain_config, extra_args)
        } else if (chain_family_selector == CHAIN_FAMILY_SELECTOR_SVM) {
            let require_valid_token_receiver = tokens_len > 0;
            let svm_gas_limit =
                resolve_svm_gas_limit(
                    dest_chain_config, extra_args, require_valid_token_receiver
                );
            let must_be_non_zero = svm_gas_limit > 0;
            validate_32byte_address(receiver, must_be_non_zero);
            svm_gas_limit
        } else {
            abort E_UNKNOWN_CHAIN_FAMILY_SELECTOR
        };

    // TODO: make sure this does not revert if no price
    let fee_token_price = get_token_price_internal(state, fee_token);
    let packed_gas_price =
        get_validated_gas_price_internal(
            state, clock, dest_chain_config, dest_chain_selector
        );

    // TODO: this should probably be premium_fee_usd_octa for sui?
    let (mut premium_fee_usd_wei, token_transfer_gas, token_transfer_bytes_overhead) =
    if (tokens_len > 0) {
        get_token_transfer_cost(
            state,
            dest_chain_config,
            dest_chain_selector,
            fee_token,
            fee_token_price,
            local_token_addresses,
            local_token_amounts
        )
    } else {
        ((dest_chain_config.network_fee_usd_cents as u256) * VAL_1E16, 0, 0)
    };
    let premium_multiplier =
        get_premium_multiplier_wei_per_eth_internal(state, fee_token);
    premium_fee_usd_wei = premium_fee_usd_wei * (premium_multiplier as u256); // Apply premium multiplier in wei/eth units

    let data_availability_cost_usd_36_decimals =
        if (dest_chain_config.dest_data_availability_multiplier_bps > 0) {
            // TODO: on EVM, the gas price is uint224 and the top 112 bits are used. here we're using a u256
            // and expecting that the extra top 22 bits are zeroes. update this and `gas_cost` below
            // if needed.
            let data_availability_gas_price = packed_gas_price >> GAS_PRICE_BITS;
            get_data_availability_cost(
                dest_chain_config,
                data_availability_gas_price,
                data_len,
                tokens_len,
                token_transfer_bytes_overhead
            )
        } else { 0 };

    let call_data_length: u256 =
        (data_len as u256) + (token_transfer_bytes_overhead as u256);
    let mut dest_call_data_cost =
        call_data_length
            * (dest_chain_config.dest_gas_per_payload_byte_base as u256);
    if (call_data_length
        > (dest_chain_config.dest_gas_per_payload_byte_threshold as u256)) {
        dest_call_data_cost = (
            dest_chain_config.dest_gas_per_payload_byte_base as u256
        ) * (dest_chain_config.dest_gas_per_payload_byte_threshold as u256)
            + (
            call_data_length
                - (dest_chain_config.dest_gas_per_payload_byte_threshold as u256)
        ) * (dest_chain_config.dest_gas_per_payload_byte_high as u256);
    };

    let total_dest_chain_gas =
        (dest_chain_config.dest_gas_overhead as u256) + (token_transfer_gas as u256)
            + dest_call_data_cost + gas_limit;

    let gas_cost = packed_gas_price & (MAX_U256 >> (255 - GAS_PRICE_BITS + 1));

    let total_cost_usd =
        (
            total_dest_chain_gas * gas_cost
                * (dest_chain_config.gas_multiplier_wei_per_eth as u256)
        ) + premium_fee_usd_wei + data_availability_cost_usd_36_decimals;

    let fee_token_cost = total_cost_usd / fee_token_price.value;

    // we need to convert back to a u64 which is what the fungible asset module uses for amounts.
    assert!(
        fee_token_cost <= MAX_U64, E_FEE_TOKEN_COST_TOO_HIGH
    );
    fee_token_cost as u64
}

fun validate_message(
    dest_chain_config: &DestChainConfig, data_len: u64, tokens_len: u64
) {
    assert!(
        data_len <= (dest_chain_config.max_data_bytes as u64),
        E_MESSAGE_TOO_LARGE
    );
    assert!(
        tokens_len <= (dest_chain_config.max_number_of_tokens_per_msg as u64),
        E_UNSUPPORTED_NUMBER_OF_TOKENS
    );
}

fun validate_evm_address(encoded_address: vector<u8>) {
    let encoded_address_len = encoded_address.length();
    // TODO: why this is checking against 32 bytes?
    assert!(
        encoded_address_len == 32, E_INVALID_EVM_ADDRESS
    );

    let encoded_address_uint = eth_abi::decode_u256_value(encoded_address);

    assert!(
        encoded_address_uint >= EVM_PRECOMPILE_SPACE,
        E_INVALID_EVM_ADDRESS
    );
    assert!(
        encoded_address_uint <= MAX_U160,
        E_INVALID_EVM_ADDRESS
    );
}

fun validate_32byte_address(
    encoded_address: vector<u8>, must_be_non_zero: bool
) {
    let encoded_address_len = encoded_address.length();
    assert!(
        encoded_address_len == 32, E_INVALID_SVM_ADDRESS
    );

    if (must_be_non_zero) {
        assert!(
            encoded_address.length()== 32,
            E_INVALID_SVM_ADDRESS
        );
        let encoded_address_uint = eth_abi::decode_u256_value(encoded_address);
        assert!(
            encoded_address_uint > 0,
            E_INVALID_SVM_ADDRESS
        );
    };
}

fun resolve_generic_gas_limit(
    dest_chain_config: &DestChainConfig, extra_args: vector<u8>
): u256 {
    let extra_args_len = extra_args.length();
    if (extra_args_len == 0) {
        dest_chain_config.default_tx_gas_limit as u256
    } else {
        let (gas_limit, allow_out_of_order_execution) =
            decode_generic_extra_args(extra_args);
        assert!(
            gas_limit <= (dest_chain_config.max_per_msg_gas_limit as u256),
            E_MESSAGE_GAS_LIMIT_TOO_HIGH
        );
        assert!(
            !dest_chain_config.enforce_out_of_order || allow_out_of_order_execution,
            E_EXTRA_ARG_OUT_OF_ORDER_EXECUTION_MUST_BE_TRUE
        );
        gas_limit
    }
}

fun resolve_svm_gas_limit(
    dest_chain_config: &DestChainConfig,
    extra_args: vector<u8>,
    require_valid_token_receiver: bool
): u256 {
    let extra_args_len = extra_args.length();
    assert!(extra_args_len > 0, E_INVALID_EXTRA_ARGS_DATA);
    let (
        compute_units,
        _account_is_writable_bitmap,
        allow_out_of_order_execution,
        token_receiver,
        _accounts
    ) = decode_svm_extra_args(extra_args);
    assert!(
        !dest_chain_config.enforce_out_of_order || allow_out_of_order_execution,
        E_EXTRA_ARG_OUT_OF_ORDER_EXECUTION_MUST_BE_TRUE
    );
    assert!(
        compute_units <= dest_chain_config.max_per_msg_gas_limit,
        E_MESSAGE_COMPUTE_UNIT_LIMIT_TOO_HIGH
    );
    if (require_valid_token_receiver) {
        assert!(
            token_receiver.length() == 32,
            E_INVALID_TOKEN_RECEIVER
        );
        let token_receiver_uint = eth_abi::decode_u256_value(token_receiver);
        assert!(
            token_receiver_uint > 0,
            E_INVALID_TOKEN_RECEIVER
        );
    };
    compute_units as u256
}

fun get_token_transfer_cost(
    state: &FeeQuoterState,
    dest_chain_config: &DestChainConfig,
    dest_chain_selector: u64,
    fee_token: address,
    fee_token_price: TimestampedPrice,
    local_token_addresses: vector<address>,
    local_token_amounts: vector<u64>
): (u256, u32, u32) {
    let mut token_transfer_fee_wei: u256 = 0;
    let mut token_transfer_gas: u32 = 0;
    let mut token_transfer_bytes_overhead: u32 = 0;

    vector::zip_do_ref!(
        &local_token_addresses,
        &local_token_amounts,
        |local_token_address, local_token_amount| {
            let local_token_address: address = *local_token_address;
            let local_token_amount: u64 = *local_token_amount;

            let transfer_fee_config =
                get_token_transfer_fee_config_internal(
                    state, dest_chain_selector, local_token_address
                );

            if (!transfer_fee_config.is_enabled) {
                token_transfer_fee_wei = token_transfer_fee_wei
                + ((dest_chain_config.default_token_fee_usd_cents as u256)
                * VAL_1E16);
                token_transfer_gas = token_transfer_gas
                + dest_chain_config.default_token_dest_gas_overhead;
                token_transfer_bytes_overhead = token_transfer_bytes_overhead
                + CCIP_LOCK_OR_BURN_V1_RET_BYTES;
            } else {
                let mut bps_fee_usd_wei = 0;
                if (transfer_fee_config.deci_bps > 0) {
                    let token_price =
                    if (local_token_address == fee_token) {
                        fee_token_price
                    } else {
                        get_token_price_internal(state, local_token_address)
                    };
                    let token_usd_value =
                        calc_usd_value_from_token_amount(
                            local_token_amount, token_price.value
                        );
                    bps_fee_usd_wei = (
                        token_usd_value * (transfer_fee_config.deci_bps as u256)
                    ) / VAL_1E5;
                };

                token_transfer_gas = token_transfer_gas + transfer_fee_config.dest_gas_overhead;
                token_transfer_bytes_overhead = token_transfer_bytes_overhead + transfer_fee_config.dest_bytes_overhead;

                let min_fee_usd_wei =
                    (transfer_fee_config.min_fee_usd_cents as u256) * VAL_1E16;
                let max_fee_usd_wei =
                    (transfer_fee_config.max_fee_usd_cents as u256) * VAL_1E16;
                let selected_fee_usd_wei =
                    if (bps_fee_usd_wei < min_fee_usd_wei) {
                        min_fee_usd_wei
                    } else if (bps_fee_usd_wei > max_fee_usd_wei) {
                        max_fee_usd_wei
                    } else {
                        bps_fee_usd_wei
                    };
                token_transfer_fee_wei = token_transfer_fee_wei+ selected_fee_usd_wei;
            }
        }
    );

    (token_transfer_fee_wei, token_transfer_gas, token_transfer_bytes_overhead)
}

fun calc_usd_value_from_token_amount(
    token_amount: u64, token_price: u256
): u256 {
    (token_amount as u256) * (token_price as u256) / VAL_1E18
}

fun get_token_transfer_fee_config_internal(
    state: &FeeQuoterState, dest_chain_selector: u64, token: address
): &TokenTransferFeeConfig {
    assert!(
        table::contains(
            &state.token_transfer_fee_configs, dest_chain_selector
        ),
        E_UNKNOWN_DEST_CHAIN_SELECTOR
    );
    let dest_chain_fee_configs =
        table::borrow(&state.token_transfer_fee_configs, dest_chain_selector);
    assert!(
        table::contains(dest_chain_fee_configs, token),
        E_TOKEN_NOT_SUPPORTED
    );
    table::borrow(dest_chain_fee_configs, token)
}

public fun get_premium_multiplier_wei_per_eth(
    ref: &CCIPObjectRef,
    token: address
): u64 {
    let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
    get_premium_multiplier_wei_per_eth_internal(state, token)
}

fun get_premium_multiplier_wei_per_eth_internal(
    state: &FeeQuoterState, token: address
): u64 {
    assert!(
        table::contains(&state.premium_multiplier_wei_per_eth, token),
        E_UNKNOWN_TOKEN
    );
    *table::borrow(&state.premium_multiplier_wei_per_eth, token)
}

fun decode_generic_extra_args(extra_args: vector<u8>): (u256, bool) {
    let extra_args_len = extra_args.length();
    let args_tag = slice(&extra_args, 0, 4);
    let args_data = slice(&extra_args, 4, extra_args_len - 4);

    if (args_tag == GENERIC_EXTRA_ARGS_V2_TAG) {
        decode_generic_extra_args_v2(args_data)
    } else {
        abort E_INVALID_EXTRA_ARGS_TAG
    }
}

fun decode_generic_extra_args_v2(extra_args: vector<u8>): (u256, bool) {
    let mut stream = eth_abi::new_stream(extra_args);
    let gas_limit = eth_abi::decode_u256(&mut stream);
    let allow_out_of_order_execution = eth_abi::decode_bool(&mut stream);
    (gas_limit, allow_out_of_order_execution)
}

fun encode_generic_extra_args_v2(
    gas_limit: u256, allow_out_of_order_execution: bool
): vector<u8> {
    let mut extra_args = vector[];
    eth_abi::encode_selector(&mut extra_args, GENERIC_EXTRA_ARGS_V2_TAG);
    eth_abi::encode_u256(&mut extra_args, gas_limit);
    eth_abi::encode_bool(&mut extra_args, allow_out_of_order_execution);
    extra_args
}

fun decode_svm_extra_args(
    extra_args: vector<u8>
): (u32, u64, bool, vector<u8>, vector<vector<u8>>) {
    // TODO: we need extra validation here. if extra_args length is less than tag length + data length,
    let extra_args_len = extra_args.length();
    let args_tag = slice(&extra_args, 0, 4);
    assert!(
        args_tag == SVM_EXTRA_ARGS_V1_TAG,
        E_INVALID_EXTRA_ARGS_TAG
    );
    let args_data = slice(&extra_args, 4, extra_args_len - 4);
    decode_svm_extra_args_v1(args_data)
}

fun decode_svm_extra_args_v1(
    extra_args: vector<u8>
): (u32, u64, bool, vector<u8>, vector<vector<u8>>) {
    let mut stream = eth_abi::new_stream(extra_args);
    let compute_units = eth_abi::decode_u32(&mut stream);
    let account_is_writable_bitmap = eth_abi::decode_u64(&mut stream);
    let allow_out_of_order_execution = eth_abi::decode_bool(&mut stream);
    let token_receiver = eth_abi::decode_bytes32(&mut stream);
    let accounts =
        eth_abi::decode_vector!(
            &mut stream,
            |stream| { eth_abi::decode_bytes32(stream) }
        );
    (
        compute_units,
        account_is_writable_bitmap,
        allow_out_of_order_execution,
        token_receiver,
        accounts
    )
}

/// Returns a new vector containing `len` elements from `vec`
/// starting at index `start`. Panics if `start + len` exceeds the vector length.
fun slice<T: copy>(vec: &vector<T>, start: u64, len: u64): vector<T> {
    let vec_len = vec.length();
    // Ensure we have enough elements for the slice.
    assert!(start + len <= vec_len, E_OUT_OF_BOUND);
    let mut new_vec = vector::empty<T>();
    let mut i = start;
    while (i < start + len) {
        // Copy each element from the original vector into the new vector.
        new_vec.push_back(vec[i]);
        i = i + 1;
    };
    new_vec
}

fun get_data_availability_cost(
    dest_chain_config: &DestChainConfig,
    data_availability_gas_price: u256,
    data_len: u64,
    tokens_len: u64,
    total_transfer_bytes_overhead: u32
): u256 {
    let data_availability_length_bytes =
        MESSAGE_FIXED_BYTES + data_len + (tokens_len
            * MESSAGE_FIXED_BYTES_PER_TOKEN)
            + (total_transfer_bytes_overhead as u64);

    let data_availability_gas =
        ((data_availability_length_bytes as u256)
            * (dest_chain_config.dest_gas_per_data_availability_byte as u256)) + (
            dest_chain_config.dest_data_availability_overhead_gas as u256
        );

    data_availability_gas * data_availability_gas_price
        * (dest_chain_config.dest_data_availability_multiplier_bps as u256)
        * VAL_1E14
}

fun process_chain_family_selector(
    dest_chain_config: &DestChainConfig,
    is_message_with_token_transfers: bool,
    extra_args: vector<u8>
): (vector<u8>, bool) {
    let chain_family_selector = dest_chain_config.chain_family_selector;
    if (chain_family_selector == CHAIN_FAMILY_SELECTOR_EVM
        || chain_family_selector == CHAIN_FAMILY_SELECTOR_APTOS
        || chain_family_selector == CHAIN_FAMILY_SELECTOR_SUI) {
        let (gas_limit, allow_out_of_order_execution) =
            decode_generic_extra_args(extra_args);
        let extra_args_v2 =
            encode_generic_extra_args_v2(gas_limit, allow_out_of_order_execution);
        (extra_args_v2, allow_out_of_order_execution)
    } else if (chain_family_selector == CHAIN_FAMILY_SELECTOR_SVM) {
        let (
            _compute_units,
            _account_is_writable_bitmap,
            allow_out_of_order_execution,
            token_receiver,
            _accounts
        ) = decode_svm_extra_args(extra_args);
        if (is_message_with_token_transfers) {
            assert!(
                token_receiver.length() == 32,
                E_INVALID_TOKEN_RECEIVER
            );
            let token_receiver_uint = eth_abi::decode_u256_value(token_receiver);
            assert!(
                token_receiver_uint > 0,
                E_INVALID_TOKEN_RECEIVER
            );
        };
        (extra_args, allow_out_of_order_execution)
    } else {
        // TODO: add support for Sui?
        abort E_UNKNOWN_CHAIN_FAMILY_SELECTOR
    }
}

fun process_pool_return_data(
    dest_chain_config: &DestChainConfig,
    token_transfer_fee_config: &TokenTransferFeeConfig,
    dest_token_addresses: vector<vector<u8>>,
    dest_pool_datas: vector<vector<u8>>
): vector<vector<u8>> {
    let chain_family_selector = dest_chain_config.chain_family_selector;

    let tokens_len = dest_token_addresses.length();

    let mut dest_exec_data_per_token = vector[];
    let mut i = 0;
    while (i < tokens_len) {
        let dest_token_address = dest_token_addresses[i];
        let dest_pool_data = &dest_pool_datas[i];
        let dest_pool_data_len = dest_pool_data.length();
        if (dest_pool_data_len > (CCIP_LOCK_OR_BURN_V1_RET_BYTES as u64)) {
            assert!(
                dest_pool_data_len <= (token_transfer_fee_config.dest_bytes_overhead as u64),
                E_SOURCE_TOKEN_DATA_TOO_LARGE
            );
        };

        if (chain_family_selector == CHAIN_FAMILY_SELECTOR_EVM) {
            validate_evm_address(dest_token_address);
            // TODO: does it make sense to add SVM here since SVM is not going to be the dest chain?
        } else if (chain_family_selector == CHAIN_FAMILY_SELECTOR_SVM
            || chain_family_selector == CHAIN_FAMILY_SELECTOR_APTOS
            || chain_family_selector == CHAIN_FAMILY_SELECTOR_SUI) {
            validate_32byte_address(dest_token_address, /* must_be_non_zero= */ true);
        };

        let dest_gas_amount =
            if (token_transfer_fee_config.is_enabled) {
                token_transfer_fee_config.dest_gas_overhead
            } else {
                dest_chain_config.default_token_dest_gas_overhead
            };

        let mut dest_exec_data = vector[];
        eth_abi::encode_u32(&mut dest_exec_data, dest_gas_amount);
        dest_exec_data_per_token.push_back(dest_exec_data);

        i = i + 1;
    };

    dest_exec_data_per_token
}

public fun process_message_args(
    ref: &CCIPObjectRef,
    dest_chain_selector: u64,
    fee_token: address,
    fee_token_amount: u64,
    extra_args: vector<u8>,
    dest_token_addresses: vector<vector<u8>>,
    dest_pool_datas: vector<vector<u8>>
): (u64, bool, vector<u8>, vector<vector<u8>>) {
    let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
    let msg_fee_juels =
        if (fee_token == state.link_token) {
            fee_token_amount
        } else {
            convert_token_amount_internal(
                state,
                fee_token,
                fee_token_amount,
                state.link_token
            )
        };

    assert!(
        msg_fee_juels <= state.max_fee_juels_per_msg,
        E_MESSAGE_FEE_TOO_HIGH
    );

    let dest_chain_config = get_dest_chain_config_internal(
        state, dest_chain_selector
    );
    let token_transfer_fee_config =
        get_token_transfer_fee_config_internal(state, dest_chain_selector, fee_token);

    let (converted_extra_args, is_out_of_order_execution) =
        process_chain_family_selector(
            dest_chain_config,
            !vector::is_empty(&dest_token_addresses),
            extra_args
        );

    let dest_exec_data_per_token =
        process_pool_return_data(
            dest_chain_config,
            token_transfer_fee_config,
            dest_token_addresses,
            dest_pool_datas
        );

    (
        msg_fee_juels,
        is_out_of_order_execution,
        converted_extra_args,
        dest_exec_data_per_token
    )
}

#[test_only]
public fun get_fee_tokens(ref: &CCIPObjectRef): vector<address> {
    let state = state_object::borrow<FeeQuoterState>(ref, FEE_QUOTER_STATE_NAME);
    state.fee_tokens
}

#[test_only]
public fun create_fee_quoter_cap(ctx: &mut TxContext): FeeQuoterCap {
    FeeQuoterCap {
        id: object::new(ctx),
    }
}

#[test_only]
public fun destroy_fee_quoter_cap(cap: FeeQuoterCap) {
    let FeeQuoterCap { id } = cap;
    object::delete(id);
}
