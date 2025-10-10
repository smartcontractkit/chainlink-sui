#[test_only]
#[allow(implicit_const_copy)]
module ccip::fee_quoter_mcms_test;

use ccip::fee_quoter;
use ccip::ownable::OwnerCap;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::upgrade_registry;
use mcms::mcms_account;
use mcms::mcms_deployer;
use mcms::mcms_registry::{Self, Registry};
use std::string;
use sui::bcs;
use sui::clock::{Self, Clock};
use sui::test_scenario::{Self as ts, Scenario};

const OWNER: address = @0x123;
const TOKEN_1: address = @0x1000;
const TOKEN_2: address = @0x2000;
const LINK_TOKEN: address = @0x4000;
const DEST_CHAIN_SELECTOR_1: u64 = 1;

const FEE_QUOTER_MODULE_NAME: vector<u8> = b"fee_quoter";

public struct Env {
    scenario: Scenario,
    ref: CCIPObjectRef,
    registry: Registry,
    clock: Clock,
}

fun setup(): Env {
    let mut scenario = ts::begin(OWNER);

    {
        let ctx = scenario.ctx();
        // Initialize MCMS components
        mcms_account::test_init(ctx);
        mcms_registry::test_init(ctx);
        mcms_deployer::test_init(ctx);

        // Initialize CCIP state object
        state_object::test_init(ctx);
    };

    scenario.next_tx(OWNER);

    let registry = ts::take_shared<Registry>(&scenario);
    let mut ref = ts::take_shared<CCIPObjectRef>(&scenario);
    let clock = clock::create_for_testing(scenario.ctx());

    // Initialize fee quoter
    let state_object_owner_cap = ts::take_from_sender<ccip::ownable::OwnerCap>(&scenario);
    upgrade_registry::initialize(&mut ref, &state_object_owner_cap, scenario.ctx());
    fee_quoter::initialize(
        &mut ref,
        &state_object_owner_cap,
        1000000000000000000000, // max_fee_juels_per_msg
        LINK_TOKEN,
        3600, // token_price_staleness_threshold
        vector[LINK_TOKEN], // fee_tokens
        scenario.ctx(),
    );
    ts::return_to_address(OWNER, state_object_owner_cap);

    scenario.next_tx(OWNER);

    Env {
        scenario,
        ref,
        registry,
        clock,
    }
}

fun tear_down(env: Env) {
    let Env { scenario, ref, registry, clock } = env;
    ts::return_shared(ref);
    ts::return_shared(registry);
    clock.destroy_for_testing();
    ts::end(scenario);
}

/// Helper function to transfer ownership to MCMS using the standard 3-step process
fun transfer_ownership_to_mcms(env: &mut Env, owner_cap: OwnerCap) {
    // Step 1: transfer_ownership to MCMS multisig address
    state_object::transfer_ownership(
        &mut env.ref,
        &owner_cap,
        mcms_registry::get_multisig_address(),
        env.scenario.ctx(),
    );

    // Step 2: accept the ownership transfer as the multisig address
    env.scenario.next_tx(mcms_registry::get_multisig_address());
    state_object::accept_ownership(&mut env.ref, env.scenario.ctx());

    // Step 3: register the OwnerCap with MCMS
    state_object::execute_ownership_transfer_to_mcms(
        &mut env.ref,
        owner_cap,
        &mut env.registry,
        @mcms,
        env.scenario.ctx(),
    );

    // Switch back to original test context
    env.scenario.next_tx(OWNER);
}

#[test]
public fun test_mcms_apply_fee_token_updates() {
    let mut env = setup();

    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

    // Prepare data for mcms_apply_fee_token_updates
    let mut data = vector::empty<u8>();
    data.append(bcs::to_bytes(&object::id_address(&env.ref)));
    data.append(bcs::to_bytes(&object::id_address(&owner_cap)));
    data.append(bcs::to_bytes(&vector<address>[])); // fee_tokens_to_remove
    data.append(bcs::to_bytes(&vector[TOKEN_1, TOKEN_2])); // fee_tokens_to_add

    transfer_ownership_to_mcms(&mut env, owner_cap);

    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(FEE_QUOTER_MODULE_NAME),
        string::utf8(b"apply_fee_token_updates"),
        data,
        x"0000000000000000000000000000000000000000000000000000000000000001",
        0,
        1,
    );

    fee_quoter::mcms_apply_fee_token_updates(
        &mut env.ref,
        &mut env.registry,
        params,
        env.scenario.ctx(),
    );

    // Verify fee tokens were added
    let fee_tokens = fee_quoter::get_fee_tokens(&env.ref);
    assert!(fee_tokens.contains(&LINK_TOKEN), 0); // Original token
    assert!(fee_tokens.contains(&TOKEN_1), 1);
    assert!(fee_tokens.contains(&TOKEN_2), 2);

    env.tear_down();
}

#[test]
public fun test_mcms_apply_dest_chain_config_updates() {
    let mut env = setup();

    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

    // Prepare data for mcms_apply_dest_chain_config_updates
    let mut data = vector::empty<u8>();
    data.append(bcs::to_bytes(&object::id_address(&env.ref)));
    data.append(bcs::to_bytes(&object::id_address(&owner_cap)));
    data.append(bcs::to_bytes(&DEST_CHAIN_SELECTOR_1)); // dest_chain_selector
    data.append(bcs::to_bytes(&true)); // is_enabled
    data.append(bcs::to_bytes(&(10 as u16))); // max_number_of_tokens_per_msg
    data.append(bcs::to_bytes(&(10000 as u32))); // max_data_bytes
    data.append(bcs::to_bytes(&(1000000 as u32))); // max_per_msg_gas_limit
    data.append(bcs::to_bytes(&(100000 as u32))); // dest_gas_overhead
    data.append(bcs::to_bytes(&(10 as u8))); // dest_gas_per_payload_byte_base
    data.append(bcs::to_bytes(&(20 as u8))); // dest_gas_per_payload_byte_high
    data.append(bcs::to_bytes(&(1000 as u16))); // dest_gas_per_payload_byte_threshold
    data.append(bcs::to_bytes(&(50000 as u32))); // dest_data_availability_overhead_gas
    data.append(bcs::to_bytes(&(100 as u16))); // dest_gas_per_data_availability_byte
    data.append(bcs::to_bytes(&(1000 as u16))); // dest_data_availability_multiplier_bps
    data.append(bcs::to_bytes(&x"2812d52c")); // chain_family_selector (EVM)
    data.append(bcs::to_bytes(&false)); // enforce_out_of_order
    data.append(bcs::to_bytes(&(50 as u16))); // default_token_fee_usd_cents
    data.append(bcs::to_bytes(&(50000 as u32))); // default_token_dest_gas_overhead
    data.append(bcs::to_bytes(&(500000 as u32))); // default_tx_gas_limit
    data.append(bcs::to_bytes(&(1000000000000000000 as u64))); // gas_multiplier_wei_per_eth
    data.append(bcs::to_bytes(&(3600 as u32))); // gas_price_staleness_threshold
    data.append(bcs::to_bytes(&(100 as u32))); // network_fee_usd_cents

    transfer_ownership_to_mcms(&mut env, owner_cap);

    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(FEE_QUOTER_MODULE_NAME),
        string::utf8(b"apply_dest_chain_config_updates"),
        data,
        x"0000000000000000000000000000000000000000000000000000000000000002",
        0,
        1,
    );

    fee_quoter::mcms_apply_dest_chain_config_updates(
        &mut env.ref,
        &mut env.registry,
        params,
        env.scenario.ctx(),
    );

    // Verify dest chain config was updated
    let dest_chain_config = fee_quoter::get_dest_chain_config(&env.ref, DEST_CHAIN_SELECTOR_1);
    let (
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
        network_fee_usd_cents,
    ) = fee_quoter::get_dest_chain_config_fields(dest_chain_config);

    assert!(is_enabled == true);
    assert!(max_number_of_tokens_per_msg == 10);
    assert!(max_data_bytes == 10000);
    assert!(max_per_msg_gas_limit == 1000000);
    assert!(dest_gas_overhead == 100000);
    assert!(dest_gas_per_payload_byte_base == 10);
    assert!(dest_gas_per_payload_byte_high == 20);
    assert!(dest_gas_per_payload_byte_threshold == 1000);
    assert!(dest_data_availability_overhead_gas == 50000);
    assert!(dest_gas_per_data_availability_byte == 100);
    assert!(dest_data_availability_multiplier_bps == 1000);
    assert!(chain_family_selector == x"2812d52c");
    assert!(enforce_out_of_order == false);
    assert!(default_token_fee_usd_cents == 50);
    assert!(default_token_dest_gas_overhead == 50000);
    assert!(default_tx_gas_limit == 500000);
    assert!(gas_multiplier_wei_per_eth == 1000000000000000000);
    assert!(gas_price_staleness_threshold == 3600);
    assert!(network_fee_usd_cents == 100);

    env.tear_down();
}

#[test]
public fun test_mcms_apply_token_transfer_fee_config_updates() {
    let mut env = setup();

    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

    // Prepare data for mcms_apply_token_transfer_fee_config_updates
    let mut data = vector::empty<u8>();
    data.append(bcs::to_bytes(&object::id_address(&env.ref)));
    data.append(bcs::to_bytes(&object::id_address(&owner_cap)));
    data.append(bcs::to_bytes(&DEST_CHAIN_SELECTOR_1)); // dest_chain_selector
    data.append(bcs::to_bytes(&vector[TOKEN_1, TOKEN_2])); // add_tokens
    data.append(bcs::to_bytes(&vector[25 as u32, 30 as u32])); // add_min_fee_usd_cents
    data.append(bcs::to_bytes(&vector[100 as u32, 150 as u32])); // add_max_fee_usd_cents
    data.append(bcs::to_bytes(&vector[50 as u16, 75 as u16])); // add_deci_bps
    data.append(bcs::to_bytes(&vector[25000 as u32, 30000 as u32])); // add_dest_gas_overhead
    data.append(bcs::to_bytes(&vector[32 as u32, 64 as u32])); // add_dest_bytes_overhead
    data.append(bcs::to_bytes(&vector[true, true])); // add_is_enabled
    data.append(bcs::to_bytes(&vector<address>[])); // remove_tokens

    transfer_ownership_to_mcms(&mut env, owner_cap);

    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(FEE_QUOTER_MODULE_NAME),
        string::utf8(b"apply_token_transfer_fee_config_updates"),
        data,
        x"0000000000000000000000000000000000000000000000000000000000000003",
        0,
        1,
    );

    fee_quoter::mcms_apply_token_transfer_fee_config_updates(
        &mut env.ref,
        &mut env.registry,
        params,
        env.scenario.ctx(),
    );

    // Verify token transfer fee configs were added
    let config_1 = fee_quoter::get_token_transfer_fee_config(
        &env.ref,
        DEST_CHAIN_SELECTOR_1,
        TOKEN_1,
    );
    let (
        min_fee_1,
        max_fee_1,
        deci_bps_1,
        dest_gas_overhead_1,
        dest_bytes_overhead_1,
        is_enabled_1,
    ) = fee_quoter::get_token_transfer_fee_config_fields(config_1);

    assert!(min_fee_1 == 25);
    assert!(max_fee_1 == 100);
    assert!(deci_bps_1 == 50);
    assert!(dest_gas_overhead_1 == 25000);
    assert!(dest_bytes_overhead_1 == 32);
    assert!(is_enabled_1 == true);

    let config_2 = fee_quoter::get_token_transfer_fee_config(
        &env.ref,
        DEST_CHAIN_SELECTOR_1,
        TOKEN_2,
    );
    let (
        min_fee_2,
        max_fee_2,
        deci_bps_2,
        dest_gas_overhead_2,
        dest_bytes_overhead_2,
        is_enabled_2,
    ) = fee_quoter::get_token_transfer_fee_config_fields(config_2);

    assert!(min_fee_2 == 30);
    assert!(max_fee_2 == 150);
    assert!(deci_bps_2 == 75);
    assert!(dest_gas_overhead_2 == 30000);
    assert!(dest_bytes_overhead_2 == 64);
    assert!(is_enabled_2 == true);

    env.tear_down();
}

#[test]
public fun test_mcms_apply_premium_multiplier_wei_per_eth_updates() {
    let mut env = setup();

    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

    // Prepare data for mcms_apply_premium_multiplier_wei_per_eth_updates
    let mut data = vector::empty<u8>();
    data.append(bcs::to_bytes(&object::id_address(&env.ref)));
    data.append(bcs::to_bytes(&object::id_address(&owner_cap)));
    data.append(bcs::to_bytes(&vector[TOKEN_1, TOKEN_2])); // tokens
    data.append(bcs::to_bytes(&vector[1100000000000000000 as u64, 1200000000000000000 as u64])); // premium_multiplier_wei_per_eth (110%, 120%)

    transfer_ownership_to_mcms(&mut env, owner_cap);

    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(FEE_QUOTER_MODULE_NAME),
        string::utf8(b"apply_premium_multiplier_wei_per_eth_updates"),
        data,
        x"0000000000000000000000000000000000000000000000000000000000000004",
        0,
        1,
    );

    fee_quoter::mcms_apply_premium_multiplier_wei_per_eth_updates(
        &mut env.ref,
        &mut env.registry,
        params,
        env.scenario.ctx(),
    );

    // Verify premium multipliers were updated
    let multiplier_1 = fee_quoter::get_premium_multiplier_wei_per_eth(&env.ref, TOKEN_1);
    assert!(multiplier_1 == 1100000000000000000);

    let multiplier_2 = fee_quoter::get_premium_multiplier_wei_per_eth(&env.ref, TOKEN_2);
    assert!(multiplier_2 == 1200000000000000000);

    env.tear_down();
}

#[test]
public fun test_mcms_update_prices_with_owner_cap() {
    let mut env = setup();

    let owner_cap = ts::take_from_sender<OwnerCap>(&env.scenario);

    // Set up clock with current timestamp for staleness validation
    let current_timestamp = 1000000000; // 1 billion ms
    env.clock.set_for_testing(current_timestamp);

    // Prepare data for mcms_update_prices_with_owner_cap
    let mut data = vector::empty<u8>();
    data.append(bcs::to_bytes(&object::id_address(&env.ref)));
    data.append(bcs::to_bytes(&object::id_address(&owner_cap)));
    data.append(bcs::to_bytes(&object::id_address(&env.clock)));
    data.append(bcs::to_bytes(&vector[TOKEN_1, TOKEN_2])); // source_tokens
    data.append(bcs::to_bytes(&vector[2000000000000000000 as u256, 500000000000000 as u256])); // source_usd_per_token (2 ETH, 0.0005 ETH in wei)
    data.append(bcs::to_bytes(&vector[DEST_CHAIN_SELECTOR_1])); // gas_dest_chain_selectors
    data.append(bcs::to_bytes(&vector[30000000 as u256])); // gas_usd_per_unit_gas (0.03 USD per gas unit)

    transfer_ownership_to_mcms(&mut env, owner_cap);

    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(FEE_QUOTER_MODULE_NAME),
        string::utf8(b"update_prices_with_owner_cap"),
        data,
        x"0000000000000000000000000000000000000000000000000000000000000005",
        0,
        1,
    );

    {
        let clock = &env.clock;
        fee_quoter::mcms_update_prices_with_owner_cap(
            &mut env.ref,
            &mut env.registry,
            clock,
            params,
            env.scenario.ctx(),
        );
    };

    // Verify token prices were updated
    let price_1 = fee_quoter::get_token_price(&env.ref, TOKEN_1);
    let (token_price_1, timestamp_1) = fee_quoter::get_timestamped_price_fields(price_1);
    assert!(token_price_1 == (2000000000000000000 as u256));
    assert!(timestamp_1 == (current_timestamp / 1000)); // Timestamp is in seconds

    let price_2 = fee_quoter::get_token_price(&env.ref, TOKEN_2);
    let (token_price_2, timestamp_2) = fee_quoter::get_timestamped_price_fields(price_2);
    assert!(token_price_2 == (500000000000000 as u256));
    assert!(timestamp_2 == (current_timestamp / 1000)); // Timestamp is in seconds

    // Verify gas price was updated
    let gas_price = fee_quoter::get_dest_chain_gas_price(&env.ref, DEST_CHAIN_SELECTOR_1);
    let (usd_per_unit_gas, gas_timestamp) = fee_quoter::get_timestamped_price_fields(gas_price);
    assert!(usd_per_unit_gas == (30000000 as u256), 4);
    assert!(gas_timestamp == (current_timestamp / 1000), 5); // Timestamp is in seconds

    env.tear_down();
}
