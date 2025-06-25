#[allow(implicit_const_copy)]
#[test_only]
module ccip::fee_quoter_tests;

use std::bcs;
use std::string;
use sui::test_scenario::{Self, Scenario};
use sui::clock;

use ccip::fee_quoter::{Self, FeeQuoterState};
use ccip::state_object::{Self, OwnerCap, CCIPObjectRef};

// === Constants ===

const CHAIN_FAMILY_SELECTOR_EVM: vector<u8> = x"2812d52c";
const CHAIN_FAMILY_SELECTOR_SVM: vector<u8> = x"1e10bdc4";

// Test addresses
const MOCK_ADDRESS_1: address = @0x1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b;
const MOCK_ADDRESS_2: address = @0x000000000000000000000000F4030086522a5bEEa4988F8cA5B36dbC97BeE88c;
const MOCK_ADDRESS_3: address = @0x8a7b6c5d4e3f2a1b0c9d8e7f6a5b4c3d2e1f0a9b8c7d6e5f4a3b2c1d0e9f8a7;
const MOCK_ADDRESS_4: address = @0x3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d;
const MOCK_ADDRESS_5: address = @0xd1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2;

// Fee and gas constants
const ONE_E_18: u256 = 1_000_000_000_000_000_000;
const DEFAULT_MAX_FEE_JUELS: u256 = 200;
const DEFAULT_TOKEN_PRICE_STALENESS_THRESHOLD: u64 = 1000;
const DEFAULT_GAS_PRICE: u256 = 1_000_000_000_000;
const DEFAULT_TOKEN_PRICE: u256 = 150_000_000_000;

// Test scenario constants
const ADMIN_ADDRESS: address = @0x1;

// === Helper Functions ===

fun create_test_scenario(): Scenario {
    test_scenario::begin(ADMIN_ADDRESS)
}

fun setup_ccip_environment(): (Scenario, OwnerCap, CCIPObjectRef) {
    let mut scenario = create_test_scenario();
    let ctx = scenario.ctx();

    state_object::test_init(ctx);
    
    // Advance to next transaction to retrieve the created objects
    scenario.next_tx(ADMIN_ADDRESS);
    
    // Retrieve the OwnerCap that was transferred to the sender
    let owner_cap = scenario.take_from_sender<OwnerCap>();
    
    // Retrieve the shared CCIPObjectRef
    let ref = scenario.take_shared<CCIPObjectRef>();
    
    (scenario, owner_cap, ref)
}

fun initialize_fee_quoter(ref: &mut CCIPObjectRef, owner_cap: &OwnerCap, ctx: &mut TxContext) {
    fee_quoter::initialize(
        ref,
        owner_cap,
        DEFAULT_MAX_FEE_JUELS * ONE_E_18, // 200 link
        MOCK_ADDRESS_1,
        DEFAULT_TOKEN_PRICE_STALENESS_THRESHOLD,
        vector[
            MOCK_ADDRESS_1,
            MOCK_ADDRESS_2,
            MOCK_ADDRESS_3
        ],
        ctx
    );
}

fun setup_basic_dest_chain_config(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    dest_chain_selector: u64,
    chain_family_selector: vector<u8>,
    enforce_out_of_order: bool
) {
    fee_quoter::apply_dest_chain_config_updates(
        ref,
        owner_cap,
        dest_chain_selector,
        true, // is_enabled
        1000, // max_number_of_tokens_per_msg
        30_000, // max_data_bytes
        30_000_000, // max_per_msg_gas_limit
        250_000, // dest_gas_overhead
        16, // dest_gas_per_payload_byte_base
        200, // dest_gas_per_payload_byte_high
        300, // dest_gas_per_payload_byte_threshold
        400000, // dest_data_availability_overhead_gas
        500, // dest_gas_per_data_availability_byte
        600, // dest_data_availability_multiplier_bps
        chain_family_selector,
        enforce_out_of_order,
        50, // default_token_fee_usd_cents
        90_000, // default_token_dest_gas_overhead
        200_000, // default_tx_gas_limit
        ONE_E_18 as u64, // gas_multiplier_wei_per_eth
        1_000_000, // gas_price_staleness_threshold
        50 // network_fee_usd_cents
    );
}

fun setup_token_transfer_fee_configs(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    dest_chain_selector: u64,
    ctx: &mut TxContext
) {
    fee_quoter::apply_token_transfer_fee_config_updates(
        ref,
        owner_cap,
        dest_chain_selector,
        vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2],
        vector[100, 200],
        vector[3000, 4000],
        vector[500, 600],
        vector[700, 800],
        vector[900, 1000],
        vector[true, false],
        vector[],
        ctx,
    );
}

fun setup_price_updates(
    ref: &mut CCIPObjectRef,
    fee_quoter_cap: &fee_quoter::FeeQuoterCap,
    clock: &clock::Clock
) {
    fee_quoter::update_prices(
        ref,
        fee_quoter_cap,
        clock,
        vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2], // source_tokens
        vector[DEFAULT_TOKEN_PRICE * ONE_E_18, DEFAULT_TOKEN_PRICE * ONE_E_18], // source_usd_per_token
        vector[100, 1000], // gas_dest_chain_selectors
        vector[DEFAULT_GAS_PRICE, DEFAULT_GAS_PRICE] // gas_usd_per_unit_gas
    );
}

fun cleanup_test_scenario(scenario: Scenario, owner_cap: OwnerCap, ref: CCIPObjectRef) {
    test_scenario::return_to_sender(&scenario, owner_cap);
    test_scenario::return_shared(ref);
    test_scenario::end(scenario);
}

// === Basic Initialization Tests ===

#[test]
public fun test_initialize() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    let _state = state_object::borrow<FeeQuoterState>(&ref);

    let fee_tokens = fee_quoter::get_fee_tokens(&ref);
    assert!(fee_tokens == vector[
        MOCK_ADDRESS_1,
        MOCK_ADDRESS_2,
        MOCK_ADDRESS_3
    ]);

    cleanup_test_scenario(scenario, owner_cap, ref);
}

// === Fee Token Management Tests ===

#[test]
public fun test_apply_fee_token_updates() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    fee_quoter::apply_fee_token_updates(
        &mut ref,
        &owner_cap,
        vector[
            MOCK_ADDRESS_1,
            MOCK_ADDRESS_2
        ],
        vector[
            MOCK_ADDRESS_4,
            MOCK_ADDRESS_5
        ],
    );

    let fee_tokens = fee_quoter::get_fee_tokens(&ref);
    assert!(fee_tokens == vector[
        MOCK_ADDRESS_3,
        MOCK_ADDRESS_4,
        MOCK_ADDRESS_5
    ]);

    cleanup_test_scenario(scenario, owner_cap, ref);
}

// === Token Transfer Fee Configuration Tests ===

#[test]
public fun test_apply_token_transfer_fee_config_updates() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    setup_token_transfer_fee_configs(&mut ref, &owner_cap, 10, ctx);

    // Successful get means the config is created
    let _config1 = fee_quoter::get_token_transfer_fee_config(&ref, 10, MOCK_ADDRESS_1);
    let _config2 = fee_quoter::get_token_transfer_fee_config(&ref, 10, MOCK_ADDRESS_2);

    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
public fun test_apply_token_transfer_fee_config_updates_remove_token() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    setup_token_transfer_fee_configs(&mut ref, &owner_cap, 10, ctx);

    fee_quoter::apply_token_transfer_fee_config_updates(
        &mut ref,
        &owner_cap,
        10, // dest_chain_selector
        vector[], // source_tokens
        vector[], // min_fee_usd_cents
        vector[], // max_fee_usd_cents
        vector[], // dest_gas_overhead
        vector[], // dest_bytes_overhead
        vector[], // deci_bps
        vector[], // is_enabled
        vector[MOCK_ADDRESS_1], // remove MOCK_ADDRESS_1
        ctx,
    );

    let cfg = fee_quoter::get_token_transfer_fee_config(&ref, 10, MOCK_ADDRESS_1);
    let (min_fee_usd_cents, max_fee_usd_cents, deci_bps, dest_gas_overhead, dest_bytes_overhead, is_enabled)
        = fee_quoter::get_token_transfer_fee_config_fields(cfg);
    assert!(min_fee_usd_cents == 0);
    assert!(max_fee_usd_cents == 0);
    assert!(deci_bps == 0);
    assert!(dest_gas_overhead == 0);
    assert!(dest_bytes_overhead == 0);
    assert!(!is_enabled);

    cleanup_test_scenario(scenario, owner_cap, ref);
}

// === Premium Multiplier Tests ===

#[test]
public fun test_apply_premium_multiplier_wei_per_eth_updates() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    fee_quoter::apply_premium_multiplier_wei_per_eth_updates(
        &mut ref,
        &owner_cap,
        vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2], // source_tokens
        // 900_000_000_000_000_000 = 90%
        vector[900_000_000_000_000_000, 200_000_000_000_000_000], // premium_multiplier_wei_per_eth
    );

    assert!(fee_quoter::get_premium_multiplier_wei_per_eth(&ref, MOCK_ADDRESS_1) == 900000000000000000);
    assert!(fee_quoter::get_premium_multiplier_wei_per_eth(&ref, MOCK_ADDRESS_2) == 200000000000000000);

    cleanup_test_scenario(scenario, owner_cap, ref);
}

// === Price Update Tests ===

#[test]
public fun test_update_prices() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);
    let fee_quoter_cap = fee_quoter::create_fee_quoter_cap(ctx);

    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    
    setup_price_updates(&mut ref, &fee_quoter_cap, &clock);

    // Prices are successfully updated if we can find the config for the dest chain selector / token address
    let _timestamp_price = fee_quoter::get_dest_chain_gas_price(&ref, 100);
    let _token_price = fee_quoter::get_token_price(&ref, MOCK_ADDRESS_1);

    fee_quoter::destroy_fee_quoter_cap(fee_quoter_cap);
    clock::destroy_for_testing(clock);
    cleanup_test_scenario(scenario, owner_cap, ref);
}

// === Destination Chain Configuration Tests ===

#[test]
public fun test_apply_dest_chain_config_updates() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    setup_basic_dest_chain_config(&mut ref, &owner_cap, 100, CHAIN_FAMILY_SELECTOR_EVM, true);

    let _config = fee_quoter::get_dest_chain_config(&ref, 100);

    cleanup_test_scenario(scenario, owner_cap, ref);
}

// === Message Processing Tests ===

#[allow(implicit_const_copy)]
#[test]
public fun test_process_message_args_evm() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    setup_basic_dest_chain_config(&mut ref, &owner_cap, 100, CHAIN_FAMILY_SELECTOR_EVM, true);
    setup_token_transfer_fee_configs(&mut ref, &owner_cap, 100, ctx);

    let evm_extra_args = x"181dcf10a1a910000000000000000000000000000000000000000000000000000000000001";

    let (
        msg_fee_juels,
        is_out_of_order_execution,
        converted_extra_args,
        dest_exec_data_per_token
    ) = fee_quoter::process_message_args(
        &ref,
        100, // dest_chain_selector
        MOCK_ADDRESS_1, // fee_token
        1000, // fee_token_amount
        evm_extra_args, // extra_args
        vector[MOCK_ADDRESS_1], // source_token_addresses
        vector[
            bcs::to_bytes(&MOCK_ADDRESS_2)
        ], // dest_token_addresses
        vector[
            bcs::to_bytes(&MOCK_ADDRESS_3)
        ] // dest_pool_datas
    );

    assert!(msg_fee_juels == 10000000000000);
    assert!(is_out_of_order_execution == true);
    assert!(converted_extra_args == evm_extra_args);
    // This is the dest gas overhead. hex(02bc) = 700
    assert!(
        dest_exec_data_per_token == vector[x"bc020000"]
    );

    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[allow(implicit_const_copy)]
#[test]
public fun test_process_message_args_svm() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    setup_basic_dest_chain_config(&mut ref, &owner_cap, 100, CHAIN_FAMILY_SELECTOR_SVM, true);

    fee_quoter::apply_token_transfer_fee_config_updates(
        &mut ref,
        &owner_cap,
        100, // dest_chain_selector
        vector[MOCK_ADDRESS_1, MOCK_ADDRESS_4], // source_tokens
        vector[100, 200], // source_usd_per_token
        vector[3000, 4000], // gas_dest_chain_selectors
        vector[500, 600], // gas_usd_per_unit_gas
        vector[700, 800], // dest_gas_overhead
        vector[900, 1000], // dest_bytes_overhead
        vector[true, false], // is_enabled
        vector[], // dest_chain_selectors
        ctx,
    );

    let svm_extra_args = x"1f3b3aba65000000660000000000000001201234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef02202234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdea203234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdeb";

    let (
        msg_fee_juels,
        is_out_of_order_execution,
        converted_extra_args,
        dest_exec_data_per_token
    ) = fee_quoter::process_message_args(
        &ref,
        100, // dest_chain_selector
        MOCK_ADDRESS_1, // fee_token
        1000, // fee_token_amount
        svm_extra_args, // extra_args
        vector[MOCK_ADDRESS_1], // source_token_addresses
        vector[
            bcs::to_bytes(&MOCK_ADDRESS_4)
        ], // dest_token_addresses
        vector[
            bcs::to_bytes(&MOCK_ADDRESS_3)
        ] // dest_pool_datas
    );

    assert!(msg_fee_juels == 10000000000000);
    assert!(is_out_of_order_execution == true);
    assert!(converted_extra_args == svm_extra_args);
    // This is the dest gas overhead. hex(02bc) = 700
    assert!(
        dest_exec_data_per_token == vector[x"bc020000"]
    );

    cleanup_test_scenario(scenario, owner_cap, ref);
}

// === Fee Validation Tests ===

#[test]
public fun test_get_validated_fee() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    let fee_quoter_cap = fee_quoter::create_fee_quoter_cap(ctx);

    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    
    setup_price_updates(&mut ref, &fee_quoter_cap, &clock);

    fee_quoter::apply_dest_chain_config_updates(
        &mut ref,
        &owner_cap,
        100, // dest_chain_selector
        true, // is_enabled
        1, // max_number_of_tokens_per_msg
        30_000, // max_data_bytes
        30_000_000, // max_per_msg_gas_limit
        250_000, // dest_gas_overhead
        16, // dest_gas_per_payload_byte_base
        0, // dest_gas_per_payload_byte_high
        0, // dest_gas_per_payload_byte_threshold
        0, // dest_data_availability_overhead_gas
        0, // dest_gas_per_data_availability_byte
        0, // dest_data_availability_multiplier_bps
        CHAIN_FAMILY_SELECTOR_EVM, // chain_family_selector
        false, // enforce_out_of_order
        50, // default_token_fee_usd_cents
        90_000, // default_token_dest_gas_overhead
        200_000, // default_tx_gas_limit
        ONE_E_18 as u64, // gas_multiplier_wei_per_eth
        1_000_000, // gas_price_staleness_threshold
        50 // network_fee_usd_cents
    );

    fee_quoter::apply_token_transfer_fee_config_updates(
        &mut ref,
        &owner_cap,
        100, // dest_chain_selector
        vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2], // add_tokens
        vector[1, 2], // add_min_fee_usd_cents
        vector[30, 40], // add_max_fee_usd_cents
        vector[50, 60], // add_deci_bps
        vector[700, 800], // add_dest_gas_overhead
        vector[90, 100], // add_dest_bytes_overhead
        vector[true, false], // add_is_enabled
        vector[], // remove_tokens
        ctx,
    );

    fee_quoter::apply_premium_multiplier_wei_per_eth_updates(
        &mut ref,
        &owner_cap,
        vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2], // source_tokens
        // 900_000_000_000_000_000 = 90%
        vector[900_000_000_000_000_000, 900_000_000_000_000_000] // premium_multiplier_wei_per_eth
    );

    // The gas limit is hex(503412) = 5256210
    let evm_extra_args = x"181dcf10123450000000000000000000000000000000000000000000000000000000000001";

    let val = fee_quoter::get_validated_fee(
        &ref,
        &clock,
        100,
        x"000000000000000000000000f4030086522a5beea4988f8ca5b36dbc97bee88c", // receiver
        b"456abc", // data
        vector[MOCK_ADDRESS_1], // token_addresses
        vector[100], // token_amounts
        MOCK_ADDRESS_1, // fee_token
        evm_extra_args, // extra_args
    );

    assert!(val == 36772733);

    fee_quoter::destroy_fee_quoter_cap(fee_quoter_cap);
    clock::destroy_for_testing(clock);
    cleanup_test_scenario(scenario, owner_cap, ref);
}

// === Uncovered Function Tests ===

#[test]
public fun test_type_and_version() {
    let version = fee_quoter::type_and_version();
    assert!(version == string::utf8(b"FeeQuoter 1.6.0"));
}

#[test]
public fun test_get_timestamped_price_fields() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);
    let fee_quoter_cap = fee_quoter::create_fee_quoter_cap(ctx);

    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    
    setup_price_updates(&mut ref, &fee_quoter_cap, &clock);

    let timestamped_price = fee_quoter::get_token_price(&ref, MOCK_ADDRESS_1);
    let (value, timestamp) = fee_quoter::get_timestamped_price_fields(timestamped_price);
    
    assert!(value == DEFAULT_TOKEN_PRICE * ONE_E_18);
    assert!(timestamp == 20);

    fee_quoter::destroy_fee_quoter_cap(fee_quoter_cap);
    clock::destroy_for_testing(clock);
    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
public fun test_get_token_prices() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);
    let fee_quoter_cap = fee_quoter::create_fee_quoter_cap(ctx);

    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    
    setup_price_updates(&mut ref, &fee_quoter_cap, &clock);

    let prices = fee_quoter::get_token_prices(&ref, vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2]);
    assert!(prices.length() == 2);
    
    let (value1, _timestamp1) = fee_quoter::get_timestamped_price_fields(prices[0]);
    let (value2, _timestamp2) = fee_quoter::get_timestamped_price_fields(prices[1]);
    
    assert!(value1 == DEFAULT_TOKEN_PRICE * ONE_E_18);
    assert!(value2 == DEFAULT_TOKEN_PRICE * ONE_E_18);

    fee_quoter::destroy_fee_quoter_cap(fee_quoter_cap);
    clock::destroy_for_testing(clock);
    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
public fun test_get_token_and_gas_prices() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);
    let fee_quoter_cap = fee_quoter::create_fee_quoter_cap(ctx);

    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    
    setup_price_updates(&mut ref, &fee_quoter_cap, &clock);
    setup_basic_dest_chain_config(&mut ref, &owner_cap, 100, CHAIN_FAMILY_SELECTOR_EVM, false);

    let (token_price, gas_price) = fee_quoter::get_token_and_gas_prices(&ref, &clock, MOCK_ADDRESS_1, 100);
    
    assert!(token_price == DEFAULT_TOKEN_PRICE * ONE_E_18);
    assert!(gas_price == DEFAULT_GAS_PRICE);

    fee_quoter::destroy_fee_quoter_cap(fee_quoter_cap);
    clock::destroy_for_testing(clock);
    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
public fun test_convert_token_amount() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);
    let fee_quoter_cap = fee_quoter::create_fee_quoter_cap(ctx);

    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    
    // Set different prices for tokens
    fee_quoter::update_prices(
        &mut ref,
        &fee_quoter_cap,
        &clock,
        vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2], 
        vector[100 * ONE_E_18, 200 * ONE_E_18], // MOCK_ADDRESS_1: $100, MOCK_ADDRESS_2: $200
        vector[100], 
        vector[DEFAULT_GAS_PRICE] 
    );

    // Convert 100 units of MOCK_ADDRESS_1 ($100 each) to MOCK_ADDRESS_2 ($200 each)
    // Expected: 100 * $100 / $200 = 50 units
    let converted_amount = fee_quoter::convert_token_amount(&ref, MOCK_ADDRESS_1, 100, MOCK_ADDRESS_2);
    assert!(converted_amount == 50);

    fee_quoter::destroy_fee_quoter_cap(fee_quoter_cap);
    clock::destroy_for_testing(clock);
    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
public fun test_get_token_receiver() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    setup_basic_dest_chain_config(&mut ref, &owner_cap, 100, CHAIN_FAMILY_SELECTOR_EVM, false);

    let message_receiver = x"000000000000000000000000f4030086522a5beea4988f8ca5b36dbc97bee88c";
    let evm_extra_args = x"181dcf10a1a910000000000000000000000000000000000000000000000000000000000001";

    let token_receiver = fee_quoter::get_token_receiver(&ref, 100, evm_extra_args, message_receiver);
    assert!(token_receiver == message_receiver); // For EVM, token_receiver = message_receiver

    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
public fun test_get_dest_chain_config_fields() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    setup_basic_dest_chain_config(&mut ref, &owner_cap, 100, CHAIN_FAMILY_SELECTOR_EVM, true);

    let config = fee_quoter::get_dest_chain_config(&ref, 100);
    let (
        is_enabled,
        max_number_of_tokens_per_msg,
        max_data_bytes,
        _max_per_msg_gas_limit,
        _dest_gas_overhead,
        _dest_gas_per_payload_byte_base,
        _dest_gas_per_payload_byte_high,
        _dest_gas_per_payload_byte_threshold,
        _dest_data_availability_overhead_gas,
        _dest_gas_per_data_availability_byte,
        _dest_data_availability_multiplier_bps,
        chain_family_selector,
        enforce_out_of_order,
        _default_token_fee_usd_cents,
        _default_token_dest_gas_overhead,
        _default_tx_gas_limit,
        _gas_multiplier_wei_per_eth,
        _gas_price_staleness_threshold,
        _network_fee_usd_cents,
    ) = fee_quoter::get_dest_chain_config_fields(config);

    assert!(is_enabled == true);
    assert!(max_number_of_tokens_per_msg == 1000);
    assert!(max_data_bytes == 30000);
    assert!(chain_family_selector == CHAIN_FAMILY_SELECTOR_EVM);
    assert!(enforce_out_of_order == true);

    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
public fun test_get_static_config() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    let static_config = fee_quoter::get_static_config(&ref);
    let (max_fee_juels_per_msg, link_token, token_price_staleness_threshold) = 
        fee_quoter::get_static_config_fields(static_config);

    assert!(max_fee_juels_per_msg == DEFAULT_MAX_FEE_JUELS * ONE_E_18);
    assert!(link_token == MOCK_ADDRESS_1);
    assert!(token_price_staleness_threshold == DEFAULT_TOKEN_PRICE_STALENESS_THRESHOLD);

    cleanup_test_scenario(scenario, owner_cap, ref);
}

// === Error Condition Tests ===

#[test]
#[expected_failure(abort_code = fee_quoter::EUnknownDestChainSelector)]
public fun test_get_dest_chain_config_unknown_chain() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    // Try to get config for a chain that was never configured
    let _config = fee_quoter::get_dest_chain_config(&ref, 999);

    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::EDestChainNotEnabled)]
public fun test_get_token_and_gas_prices_chain_not_enabled() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);
    let fee_quoter_cap = fee_quoter::create_fee_quoter_cap(ctx);

    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    
    setup_price_updates(&mut ref, &fee_quoter_cap, &clock);

    // Configure destination chain as DISABLED
    fee_quoter::apply_dest_chain_config_updates(
        &mut ref,
        &owner_cap,
        100, // dest_chain_selector
        false, // is_enabled = FALSE
        1000, // max_number_of_tokens_per_msg
        30_000, // max_data_bytes
        30_000_000, // max_per_msg_gas_limit
        250_000, // dest_gas_overhead
        16, // dest_gas_per_payload_byte_base
        200, // dest_gas_per_payload_byte_high
        300, // dest_gas_per_payload_byte_threshold
        400000, // dest_data_availability_overhead_gas
        500, // dest_gas_per_data_availability_byte
        600, // dest_data_availability_multiplier_bps
        CHAIN_FAMILY_SELECTOR_EVM,
        false, // enforce_out_of_order
        50, // default_token_fee_usd_cents
        90_000, // default_token_dest_gas_overhead
        200_000, // default_tx_gas_limit
        ONE_E_18 as u64, // gas_multiplier_wei_per_eth
        1_000_000, // gas_price_staleness_threshold
        50 // network_fee_usd_cents
    );

    // This should fail because the destination chain is disabled
    let (_token_price, _gas_price) = fee_quoter::get_token_and_gas_prices(&ref, &clock, MOCK_ADDRESS_1, 100);

    fee_quoter::destroy_fee_quoter_cap(fee_quoter_cap);
    clock::destroy_for_testing(clock);
    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::ETokenUpdateMismatch)]
public fun test_update_prices_token_update_mismatch() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);
    let fee_quoter_cap = fee_quoter::create_fee_quoter_cap(ctx);

    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    
    // This should fail because source_tokens has 2 elements but source_usd_per_token has 1 element
    fee_quoter::update_prices(
        &mut ref,
        &fee_quoter_cap,
        &clock,
        vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2], // source_tokens (2 elements)
        vector[DEFAULT_TOKEN_PRICE * ONE_E_18], // source_usd_per_token (1 element - MISMATCH!)
        vector[100], // gas_dest_chain_selectors
        vector[DEFAULT_GAS_PRICE] // gas_usd_per_unit_gas
    );

    fee_quoter::destroy_fee_quoter_cap(fee_quoter_cap);
    clock::destroy_for_testing(clock);
    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::EInvalidExtraArgsData)]
public fun test_get_validated_fee_invalid_extra_args_data_too_short() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    let fee_quoter_cap = fee_quoter::create_fee_quoter_cap(ctx);

    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    
    setup_price_updates(&mut ref, &fee_quoter_cap, &clock);
    setup_basic_dest_chain_config(&mut ref, &owner_cap, 100, CHAIN_FAMILY_SELECTOR_EVM, false);

    // Create extra args that are too short (less than 4 bytes for the tag)
    let invalid_extra_args = x"12"; // Only 1 byte, should be at least 4

    // This should fail because extra args data is too short
    let _val = fee_quoter::get_validated_fee(
        &ref,
        &clock,
        100, // dest_chain_selector
        x"000000000000000000000000f4030086522a5beea4988f8ca5b36dbc97bee88c", // receiver
        b"test", // data
        vector[], // token_addresses
        vector[], // token_amounts
        MOCK_ADDRESS_1, // fee_token
        invalid_extra_args, // extra_args too short
    );

    fee_quoter::destroy_fee_quoter_cap(fee_quoter_cap);
    clock::destroy_for_testing(clock);
    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::EInvalidTokenReceiver)]
public fun test_get_validated_fee_invalid_token_receiver_svm() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    let fee_quoter_cap = fee_quoter::create_fee_quoter_cap(ctx);

    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    
    setup_price_updates(&mut ref, &fee_quoter_cap, &clock);
    setup_basic_dest_chain_config(&mut ref, &owner_cap, 100, CHAIN_FAMILY_SELECTOR_SVM, false);

    // Create SVM extra args with ZERO token receiver (invalid for token transfers)
    // Using the client module to create properly formatted SVM extra args
    let svm_extra_args = {
        use ccip::client;
        client::encode_svm_extra_args_v1(
            100, // compute_units
            0, // account_is_writable_bitmap
            false, // allow_out_of_order_execution
            x"0000000000000000000000000000000000000000000000000000000000000000", // token_receiver = all zeros (INVALID)
            vector[] // accounts
        )
    };

    // This should fail because token_receiver is zero but we have token transfers
    let _val = fee_quoter::get_validated_fee(
        &ref,
        &clock,
        100, // dest_chain_selector
        x"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", // receiver (32 bytes for SVM)
        b"test", // data
        vector[MOCK_ADDRESS_1], // token_addresses (non-empty, requires valid token_receiver)
        vector[100], // token_amounts
        MOCK_ADDRESS_1, // fee_token
        svm_extra_args, // extra_args with zero token_receiver
    );

    fee_quoter::destroy_fee_quoter_cap(fee_quoter_cap);
    clock::destroy_for_testing(clock);
    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::EMessageFeeTooHigh)]
public fun test_process_message_args_message_fee_too_high() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    
    // Initialize with a very LOW max fee limit
    fee_quoter::initialize(
        &mut ref,
        &owner_cap,
        100, // max_fee_juels_per_msg = 100 juels (very low)
        MOCK_ADDRESS_1,
        DEFAULT_TOKEN_PRICE_STALENESS_THRESHOLD,
        vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2, MOCK_ADDRESS_3],
        ctx
    );

    setup_basic_dest_chain_config(&mut ref, &owner_cap, 100, CHAIN_FAMILY_SELECTOR_EVM, false);

    let evm_extra_args = x"181dcf10a1a910000000000000000000000000000000000000000000000000000000000001";

    // This should fail because fee_token_amount (1000) * multiplier exceeds max_fee_juels_per_msg (100)
    let (_msg_fee_juels, _is_out_of_order_execution, _converted_extra_args, _dest_exec_data_per_token) = 
        fee_quoter::process_message_args(
            &ref,
            100, // dest_chain_selector
            MOCK_ADDRESS_1, // fee_token
            1000, // fee_token_amount = 1000 (will be 10000000000000 juels, exceeds limit of 100)
            evm_extra_args, // extra_args
            vector[], // source_token_addresses
            vector[], // dest_token_addresses
            vector[] // dest_pool_datas
        );

    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::ESourceTokenDataTooLarge)]
public fun test_process_message_args_source_token_data_too_large() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    setup_basic_dest_chain_config(&mut ref, &owner_cap, 100, CHAIN_FAMILY_SELECTOR_EVM, false);

    // Configure token with very small dest_bytes_overhead
    fee_quoter::apply_token_transfer_fee_config_updates(
        &mut ref,
        &owner_cap,
        100, // dest_chain_selector
        vector[MOCK_ADDRESS_1], // add_tokens
        vector[100], // add_min_fee_usd_cents
        vector[3000], // add_max_fee_usd_cents
        vector[500], // add_deci_bps
        vector[700], // add_dest_gas_overhead
        vector[32], // add_dest_bytes_overhead = 32 (very small, equal to CCIP_LOCK_OR_BURN_V1_RET_BYTES)
        vector[true], // add_is_enabled
        vector[], // remove_tokens
        ctx,
    );

    let evm_extra_args = x"181dcf10a1a910000000000000000000000000000000000000000000000000000000000001";

    // Create pool data that's larger than the configured dest_bytes_overhead (32 bytes)
    let large_pool_data = x"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"; // 40 bytes > 32

    // This should fail because dest_pool_data (40 bytes) > dest_bytes_overhead (32 bytes)
    let (_msg_fee_juels, _is_out_of_order_execution, _converted_extra_args, _dest_exec_data_per_token) = 
        fee_quoter::process_message_args(
            &ref,
            100, // dest_chain_selector
            MOCK_ADDRESS_1, // fee_token
            1000, // fee_token_amount
            evm_extra_args, // extra_args
            vector[MOCK_ADDRESS_1], // source_token_addresses
            vector[bcs::to_bytes(&MOCK_ADDRESS_2)], // dest_token_addresses
            vector[large_pool_data] // dest_pool_datas (too large)
        );

    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::EInvalidDestChainSelector)]
public fun test_apply_dest_chain_config_invalid_dest_chain_selector_zero() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    // This should fail because dest_chain_selector is 0 (invalid)
    fee_quoter::apply_dest_chain_config_updates(
        &mut ref,
        &owner_cap,
        0, // dest_chain_selector = 0 (INVALID)
        true, // is_enabled
        1000, // max_number_of_tokens_per_msg
        30_000, // max_data_bytes
        30_000_000, // max_per_msg_gas_limit
        250_000, // dest_gas_overhead
        16, // dest_gas_per_payload_byte_base
        200, // dest_gas_per_payload_byte_high
        300, // dest_gas_per_payload_byte_threshold
        400000, // dest_data_availability_overhead_gas
        500, // dest_gas_per_data_availability_byte
        600, // dest_data_availability_multiplier_bps
        CHAIN_FAMILY_SELECTOR_EVM,
        false, // enforce_out_of_order
        50, // default_token_fee_usd_cents
        90_000, // default_token_dest_gas_overhead
        200_000, // default_tx_gas_limit
        ONE_E_18 as u64, // gas_multiplier_wei_per_eth
        1_000_000, // gas_price_staleness_threshold
        50 // network_fee_usd_cents
    );

    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::EInvalidExtraArgsData)]
public fun test_get_validated_fee_svm_empty_extra_args() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    let fee_quoter_cap = fee_quoter::create_fee_quoter_cap(ctx);

    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    
    setup_price_updates(&mut ref, &fee_quoter_cap, &clock);
    setup_basic_dest_chain_config(&mut ref, &owner_cap, 100, CHAIN_FAMILY_SELECTOR_SVM, false);

    // Create completely EMPTY extra args for SVM (invalid)
    let empty_extra_args = x""; // Empty extra args

    // This should fail because SVM requires extra args but none are provided
    let _val = fee_quoter::get_validated_fee(
        &ref,
        &clock,
        100, // dest_chain_selector
        x"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", // receiver (32 bytes for SVM)
        b"test", // data
        vector[], // token_addresses
        vector[], // token_amounts
        MOCK_ADDRESS_1, // fee_token
        empty_extra_args, // empty extra_args (invalid for SVM)
    );

    fee_quoter::destroy_fee_quoter_cap(fee_quoter_cap);
    clock::destroy_for_testing(clock);
    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::ETokenTransferFeeConfigMismatch)]
public fun test_apply_token_transfer_fee_config_updates_config_mismatch() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    fee_quoter::apply_token_transfer_fee_config_updates(
        &mut ref,
        &owner_cap,
        10,
        vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2],
        vector[100, 200],
        vector[3000], // only one value - MISMATCH!
        vector[500, 600],
        vector[700, 800],
        vector[900, 1000],
        vector[true, false],
        vector[],
        ctx,
    );

    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::EGasUpdateMismatch)]
public fun test_update_prices_gas_update_mismatch() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);
    let fee_quoter_cap = fee_quoter::create_fee_quoter_cap(ctx);

    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    
    // This should fail because gas_dest_chain_selectors has 2 elements 
    // but gas_usd_per_unit_gas has only 1 element
    fee_quoter::update_prices(
        &mut ref,
        &fee_quoter_cap,
        &clock,
        vector[MOCK_ADDRESS_1], // source_tokens
        vector[DEFAULT_TOKEN_PRICE * ONE_E_18], // source_usd_per_token (matching length)
        vector[100, 1000], // gas_dest_chain_selectors (2 elements)
        vector[DEFAULT_GAS_PRICE] // gas_usd_per_unit_gas (1 element - MISMATCH!)
    );

    fee_quoter::destroy_fee_quoter_cap(fee_quoter_cap);
    clock::destroy_for_testing(clock);
    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::ETokenTransferFeeConfigMismatch)]
public fun test_apply_token_transfer_fee_config_updates_deci_bps_mismatch() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    // This should fail because add_tokens has 2 elements but add_deci_bps has only 1 element
    fee_quoter::apply_token_transfer_fee_config_updates(
        &mut ref,
        &owner_cap,
        10, // dest_chain_selector
        vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2], // add_tokens (2 elements)
        vector[100, 200], // add_min_fee_usd_cents (2 elements)
        vector[3000, 4000], // add_max_fee_usd_cents (2 elements)
        vector[500], // add_deci_bps (1 element - MISMATCH!)
        vector[700, 800], // add_dest_gas_overhead (2 elements)
        vector[900, 1000], // add_dest_bytes_overhead (2 elements)
        vector[true, false], // add_is_enabled (2 elements)
        vector[], // remove_tokens
        ctx,
    );

    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::EFeeTokenNotSupported)]
public fun test_get_validated_fee_unsupported_fee_token() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    let fee_quoter_cap = fee_quoter::create_fee_quoter_cap(ctx);

    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    
    setup_price_updates(&mut ref, &fee_quoter_cap, &clock);

    setup_basic_dest_chain_config(&mut ref, &owner_cap, 100, CHAIN_FAMILY_SELECTOR_EVM, false);

    let evm_extra_args = x"181dcf10a1a910000000000000000000000000000000000000000000000000000000000000";

    // This should fail because MOCK_ADDRESS_4 is not in the fee_tokens list
    // (initialized with MOCK_ADDRESS_1, MOCK_ADDRESS_2, MOCK_ADDRESS_3)
    let _val = fee_quoter::get_validated_fee(
        &ref,
        &clock,
        100, // dest_chain_selector
        x"000000000000000000000000f4030086522a5beea4988f8ca5b36dbc97bee88c", // receiver
        b"test", // data
        vector[], // token_addresses (empty for simplicity)
        vector[], // token_amounts (empty for simplicity)
        MOCK_ADDRESS_4, // fee_token (NOT SUPPORTED - this will trigger the error)
        evm_extra_args, // extra_args
    );

    fee_quoter::destroy_fee_quoter_cap(fee_quoter_cap);
    clock::destroy_for_testing(clock);
    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::EExtraArgOutOfOrderExecutionMustBeTrue)]
public fun test_get_validated_fee_out_of_order_execution_required() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    let fee_quoter_cap = fee_quoter::create_fee_quoter_cap(ctx);

    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    
    setup_price_updates(&mut ref, &fee_quoter_cap, &clock);

    // Configure destination chain to ENFORCE out-of-order execution
    setup_basic_dest_chain_config(&mut ref, &owner_cap, 100, CHAIN_FAMILY_SELECTOR_EVM, true);

    // Create extra args with allow_out_of_order_execution = FALSE
    // This will conflict with the chain config that enforces out-of-order execution
    let evm_extra_args = x"181dcf10a1a910000000000000000000000000000000000000000000000000000000000000"; // false for out-of-order

    // This should fail because the chain enforces out-of-order execution but extra args disable it
    let _val = fee_quoter::get_validated_fee(
        &ref,
        &clock,
        100, // dest_chain_selector
        x"000000000000000000000000f4030086522a5beea4988f8ca5b36dbc97bee88c", // receiver
        b"test", // data
        vector[], // token_addresses
        vector[], // token_amounts
        MOCK_ADDRESS_1, // fee_token
        evm_extra_args, // extra_args (out-of-order = false, but chain requires true)
    );

    fee_quoter::destroy_fee_quoter_cap(fee_quoter_cap);
    clock::destroy_for_testing(clock);
    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::EInvalidExtraArgsTag)]
public fun test_get_validated_fee_invalid_extra_args_tag() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    let fee_quoter_cap = fee_quoter::create_fee_quoter_cap(ctx);

    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    
    setup_price_updates(&mut ref, &fee_quoter_cap, &clock);
    setup_basic_dest_chain_config(&mut ref, &owner_cap, 100, CHAIN_FAMILY_SELECTOR_EVM, false);

    // Create extra args with INVALID tag (should be x"181dcf10" for EVM v2)
    let invalid_extra_args = x"deadbeefa1a910000000000000000000000000000000000000000000000000000000000000"; // invalid tag

    // This should fail because the extra args tag is invalid
    let _val = fee_quoter::get_validated_fee(
        &ref,
        &clock,
        100, // dest_chain_selector
        x"000000000000000000000000f4030086522a5beea4988f8ca5b36dbc97bee88c", // receiver
        b"test", // data
        vector[], // token_addresses
        vector[], // token_amounts
        MOCK_ADDRESS_1, // fee_token
        invalid_extra_args, // extra_args with invalid tag
    );

    fee_quoter::destroy_fee_quoter_cap(fee_quoter_cap);
    clock::destroy_for_testing(clock);
    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::EMessageComputeUnitLimitTooHigh)]
public fun test_get_validated_fee_compute_unit_limit_too_high() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    let fee_quoter_cap = fee_quoter::create_fee_quoter_cap(ctx);

    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    
    setup_price_updates(&mut ref, &fee_quoter_cap, &clock);

    // Configure SVM destination chain with LOW max_per_msg_gas_limit
    fee_quoter::apply_dest_chain_config_updates(
        &mut ref,
        &owner_cap,
        100, // dest_chain_selector
        true, // is_enabled
        1, // max_number_of_tokens_per_msg
        30_000, // max_data_bytes
        1000, // max_per_msg_gas_limit = 1000 (very low limit)
        250_000, // dest_gas_overhead
        16, // dest_gas_per_payload_byte_base
        0, // dest_gas_per_payload_byte_high
        0, // dest_gas_per_payload_byte_threshold
        0, // dest_data_availability_overhead_gas
        0, // dest_gas_per_data_availability_byte
        0, // dest_data_availability_multiplier_bps
        CHAIN_FAMILY_SELECTOR_SVM, // chain_family_selector = SVM
        false, // enforce_out_of_order
        50, // default_token_fee_usd_cents
        90_000, // default_token_dest_gas_overhead
        500, // default_tx_gas_limit (must be <= max_per_msg_gas_limit)
        ONE_E_18 as u64, // gas_multiplier_wei_per_eth
        1_000_000, // gas_price_staleness_threshold
        50 // network_fee_usd_cents
    );

    // Create SVM extra args with compute_units = 5000 (exceeds the limit of 1000)
    let svm_extra_args = x"1f3b3aba88130000660000000000000001201234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef00"; // compute_units = 5000 (0x1388)

    // This should fail because compute units (5000) exceed max_per_msg_gas_limit (1000)
    let _val = fee_quoter::get_validated_fee(
        &ref,
        &clock,
        100, // dest_chain_selector
        x"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef", // receiver (32 bytes for SVM)
        b"test", // data
        vector[], // token_addresses
        vector[], // token_amounts
        MOCK_ADDRESS_1, // fee_token
        svm_extra_args, // extra_args with high compute units
    );

    fee_quoter::destroy_fee_quoter_cap(fee_quoter_cap);
    clock::destroy_for_testing(clock);
    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::EInvalidChainFamilySelector)]
public fun test_apply_dest_chain_config_invalid_chain_family_selector() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    // This should fail because the chain_family_selector is invalid
    // Valid selectors are: EVM (x"2812d52c"), SVM (x"1e10bdc4"), Aptos (x"ac77ffec"), Sui (x"c4e05953")
    fee_quoter::apply_dest_chain_config_updates(
        &mut ref,
        &owner_cap,
        100, // dest_chain_selector
        true, // is_enabled
        1, // max_number_of_tokens_per_msg
        30_000, // max_data_bytes
        30_000_000, // max_per_msg_gas_limit
        250_000, // dest_gas_overhead
        16, // dest_gas_per_payload_byte_base
        0, // dest_gas_per_payload_byte_high
        0, // dest_gas_per_payload_byte_threshold
        0, // dest_data_availability_overhead_gas
        0, // dest_gas_per_data_availability_byte
        0, // dest_data_availability_multiplier_bps
        x"deadbeef", // chain_family_selector = INVALID selector
        false, // enforce_out_of_order
        50, // default_token_fee_usd_cents
        90_000, // default_token_dest_gas_overhead
        200_000, // default_tx_gas_limit
        ONE_E_18 as u64, // gas_multiplier_wei_per_eth
        1_000_000, // gas_price_staleness_threshold
        50 // network_fee_usd_cents
    );

    cleanup_test_scenario(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::EInvalidFeeRange)]
public fun test_apply_token_transfer_fee_config_invalid_fee_range() {
    let (mut scenario, owner_cap, mut ref) = setup_ccip_environment();
    let ctx = scenario.ctx();
    initialize_fee_quoter(&mut ref, &owner_cap, ctx);

    // This should fail because min_fee_usd_cents (5000) >= max_fee_usd_cents (3000)
    // The validation requires min_fee < max_fee
    fee_quoter::apply_token_transfer_fee_config_updates(
        &mut ref,
        &owner_cap,
        10, // dest_chain_selector
        vector[MOCK_ADDRESS_1], // add_tokens
        vector[5000], // add_min_fee_usd_cents = 5000
        vector[3000], // add_max_fee_usd_cents = 3000 (INVALID: min >= max)
        vector[500], // add_deci_bps
        vector[700], // add_dest_gas_overhead
        vector[900], // add_dest_bytes_overhead
        vector[true], // add_is_enabled
        vector[], // remove_tokens
        ctx,
    );

    cleanup_test_scenario(scenario, owner_cap, ref);
}
