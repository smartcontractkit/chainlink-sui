#[test_only]
module ccip::fee_quoter_tests;

use std::bcs;
use sui::test_scenario::{Self, Scenario};
use sui::clock;

use ccip::fee_quoter::{Self, FeeQuoterState};
use ccip::state_object::{Self, OwnerCap, CCIPObjectRef};

const CHAIN_FAMILY_SELECTOR_EVM: vector<u8> = x"2812d52c";
const CHAIN_FAMILY_SELECTOR_SVM: vector<u8> = x"1e10bdc4";
const MOCK_ADDRESS_1: address = @0x1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b;
const MOCK_ADDRESS_2: address = @0x000000000000000000000000F4030086522a5bEEa4988F8cA5B36dbC97BeE88c;
// EVM token address
const MOCK_ADDRESS_3: address = @0x8a7b6c5d4e3f2a1b0c9d8e7f6a5b4c3d2e1f0a9b8c7d6e5f4a3b2c1d0e9f8a7;
const MOCK_ADDRESS_4: address = @0x3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d;
const MOCK_ADDRESS_5: address = @0xd1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b3c4d5e6f7a8b9c0d1e2;

const ONE_E_18: u256 = 1_000_000_000_000_000_000;

fun set_up_test(): (Scenario, OwnerCap, CCIPObjectRef) {
    let mut scenario = test_scenario::begin(@0x1);
    let ctx = scenario.ctx();

    let (owner_cap, ref) = state_object::create(ctx);
    (scenario, owner_cap, ref)
}

fun initialize(ref: &mut CCIPObjectRef, owner_cap: &OwnerCap, ctx: &mut TxContext) {
    fee_quoter::initialize(
        ref,
        owner_cap,
        200 * ONE_E_18, // 200 link,
        MOCK_ADDRESS_1,
        1000,
        vector[
            MOCK_ADDRESS_1,
            MOCK_ADDRESS_2,
            MOCK_ADDRESS_3
        ],
        ctx
    );
}

fun tear_down_test(scenario: Scenario, owner_cap: OwnerCap, ref: CCIPObjectRef) {
    state_object::destroy_owner_cap(owner_cap);
    state_object::destroy_state_object(ref);
    test_scenario::end(scenario);
}

#[test]
public fun test_initialize() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();
    initialize(&mut ref, &owner_cap, ctx);

    let _state = state_object::borrow<FeeQuoterState>(&ref);

    let fee_tokens = fee_quoter::get_fee_tokens(&ref);
    assert!(fee_tokens == vector[
        MOCK_ADDRESS_1,
        MOCK_ADDRESS_2,
        MOCK_ADDRESS_3
    ]);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_apply_fee_token_updates() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();
    initialize(&mut ref, &owner_cap, ctx);

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

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_apply_token_transfer_fee_config_updates() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();
    initialize(&mut ref, &owner_cap, ctx);

    fee_quoter::apply_token_transfer_fee_config_updates(
        &mut ref,
        &owner_cap,
        10,
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

    // a successful get means the config is created.
    // we can verify the content of config but that requires an additional
    // function to expose fields within the config due to the fact that this
    // test is outside the module.
    let _config1 = fee_quoter::get_token_transfer_fee_config(&ref, 10, MOCK_ADDRESS_1);
    let _config2 = fee_quoter::get_token_transfer_fee_config(&ref, 10, MOCK_ADDRESS_2);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = fee_quoter::E_TOKEN_TRANSFER_FEE_CONFIG_MISMATCH)]
public fun test_apply_token_transfer_fee_config_updates_config_mismatch() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();
    initialize(&mut ref, &owner_cap, ctx);

    fee_quoter::apply_token_transfer_fee_config_updates(
        &mut ref,
        &owner_cap,
        10,
        vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2],
        vector[100, 200],
        vector[3000], // only one value
        vector[500, 600],
        vector[700, 800],
        vector[900, 1000],
        vector[true, false],
        vector[],
        ctx,
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_apply_token_transfer_fee_config_updates_remove_token() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();
    initialize(&mut ref, &owner_cap, ctx);

    fee_quoter::apply_token_transfer_fee_config_updates(
        &mut ref,
        &owner_cap,
        10,
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

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_apply_premium_multiplier_wei_per_eth_updates() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();
    initialize(&mut ref, &owner_cap, ctx);

    fee_quoter::apply_premium_multiplier_wei_per_eth_updates(
        &mut ref,
        &owner_cap,
        vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2], // source_tokens
        // 900_000_000_000_000_000 = 90%
        vector[900_000_000_000_000_000, 200_000_000_000_000_000], // premium_multiplier_wei_per_eth
    );

    assert!(fee_quoter::get_premium_multiplier_wei_per_eth(&ref, MOCK_ADDRESS_1) == 900000000000000000);
    assert!(fee_quoter::get_premium_multiplier_wei_per_eth(&ref, MOCK_ADDRESS_2) == 200000000000000000);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_update_prices() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();
    initialize(&mut ref, &owner_cap, ctx);
    let fee_quoter_cap = fee_quoter::create_fee_quoter_cap(ctx);

    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    fee_quoter::update_prices(
        &mut ref,
        &fee_quoter_cap,
        &clock,
        vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2], // source_tokens
        // 1e18 per 1e18 tokens. A token with 8 decimals that's worth $15 would be $15e10 * 1e18 = $15e28
        vector[150_000_000_000 * ONE_E_18, 150_000_000_000 * ONE_E_18], // source_usd_per_token
        vector[100, 1000], // gas_dest_chain_selectors
        vector[1_000_000_000_000, 1_000_000_000_000] // gas_usd_per_unit_gas
    );

    // prices are successfully updated if we can find the config for the dest chain selector / token address
    let _timestamp_price = fee_quoter::get_dest_chain_gas_price(&ref, 100);
    let _token_price = fee_quoter::get_token_price(&ref, MOCK_ADDRESS_1);

    fee_quoter::destroy_fee_quoter_cap(fee_quoter_cap);
    clock::destroy_for_testing(clock);
    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_apply_dest_chain_config_updates() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();
    initialize(&mut ref, &owner_cap, ctx);

    fee_quoter::apply_dest_chain_config_updates(
        &mut ref,
        &owner_cap,
        100, // dest_chain_selector
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
        CHAIN_FAMILY_SELECTOR_EVM, // chain_family_selector
        true, // enforce_out_of_order
        50, // default_token_fee_usd_cents
        90_000, // default_token_dest_gas_overhead
        200_000, // default_tx_gas_limit
        ONE_E_18 as u64, // gas_multiplier_wei_per_eth
        1_000_000, // gas_price_staleness_threshold
        50 // network_fee_usd_cents
    );

    let _config = fee_quoter::get_dest_chain_config(&ref, 100);

    tear_down_test(scenario, owner_cap, ref);
}

#[allow(implicit_const_copy)]
#[test]
public fun test_process_message_args_evm() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();
    initialize(&mut ref, &owner_cap, ctx);

    fee_quoter::apply_dest_chain_config_updates(
        &mut ref,
        &owner_cap,
        100, // dest_chain_selector
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
        CHAIN_FAMILY_SELECTOR_EVM, // chain_family_selector
        true, // enforce_out_of_order
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
        100,
        vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2], // source_tokens
        vector[100, 200], // source_usd_per_token
        vector[3000, 4000], // gas_dest_chain_selectors
        vector[500, 600], // gas_usd_per_unit_gas
        vector[700, 800], // dest_gas_overhead
        vector[900, 1000], // dest_bytes_overhead
        vector[true, false], // is_enabled
        vector[], // dest_chain_selectors
        ctx,
    );

    let evm_extra_args = x"181dcf10181dcf10181dcf10181dcf10181dcf10181dcf10181dcf10181dcf10181dcf100000000000000000000000000000000000000000000000000000000000000001";

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
    assert!(
        dest_exec_data_per_token == vector[x"00000000000000000000000000000000000000000000000000000000000002bc"]
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[allow(implicit_const_copy)]
#[test]
public fun test_process_message_args_svm() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();
    initialize(&mut ref, &owner_cap, ctx);

    fee_quoter::apply_dest_chain_config_updates(
        &mut ref,
        &owner_cap,
        100, // dest_chain_selector
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
        CHAIN_FAMILY_SELECTOR_SVM, // chain_family_selector
        true, // enforce_out_of_order
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

    let svm_extra_args = x"1f3b3aba00000000000000000000000000000000000000000000000000000000000dcf00000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000abc00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000111";

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
    assert!(
        dest_exec_data_per_token == vector[x"00000000000000000000000000000000000000000000000000000000000002bc"]
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_get_validated_fee() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();
    initialize(&mut ref, &owner_cap, ctx);

    let fee_quoter_cap = fee_quoter::create_fee_quoter_cap(ctx);

    let mut clock = clock::create_for_testing(ctx);
    clock::increment_for_testing(&mut clock, 20000);
    fee_quoter::update_prices(
        &mut ref,
        &fee_quoter_cap,
        &clock,
        vector[MOCK_ADDRESS_1, MOCK_ADDRESS_2], // source_tokens
        vector[150_000_000_000 * ONE_E_18, 150_000_000_000 * ONE_E_18], // source_usd_per_token
        vector[100, 1000], // gas_dest_chain_selectors
        vector[1_000_000_000_000, 1_000_000_000_000] // gas_usd_per_unit_gas
    );

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
        vector[9, 10], // add_dest_bytes_overhead
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

    let evm_extra_args = x"181dcf1000000000000000000000000000000000000000000000000000000000001dcf100000000000000000000000000000000000000000000000000000000000000001";

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

    assert!(val == 14755013);

    fee_quoter::destroy_fee_quoter_cap(fee_quoter_cap);
    clock::destroy_for_testing(clock);
    tear_down_test(scenario, owner_cap, ref);
}
