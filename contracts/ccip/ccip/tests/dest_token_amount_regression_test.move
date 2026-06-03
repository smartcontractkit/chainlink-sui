#[test_only]
module ccip::dest_token_amount_regression_test;

use ccip::client;
use ccip::offramp_state_helper::{Self, DestTransferCap};
use ccip::ownable::OwnerCap;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::token_admin_registry as registry;
use ccip::upgrade_registry;
use std::ascii;
use std::string;
use std::type_name;
use sui::test_scenario::{Self as ts, Scenario};

public struct TestTypeProof has drop {}

const OWNER: address = @0x1000;
const RECEIVER_ADDRESS: address = @0x2000;
const TOKEN_ADDRESS: address =
    @0x1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b;
const TOKEN_POOL_ADDRESS: address =
    @0xdeeb7a4662eec9f2f3def03fb937a663dddaa2e215b8078a284d026b7946c270;
const SOURCE_CHAIN_SELECTOR: u64 = 1000;
const SOURCE_AMOUNT_18_DEC: u256 = 1_000_000_000_000_000_000;
const EXPECTED_LOCAL_AMOUNT_9_DEC: u256 = 1_000_000_000;
const LOCAL_DECIMALS: u8 = 9;
// ABI-encoded u256(18) = remote decimals
const SOURCE_POOL_DATA_18_DEC: vector<u8> =
    x"0000000000000000000000000000000000000000000000000000000000000012";

fun setup_test(): (Scenario, OwnerCap, CCIPObjectRef, DestTransferCap) {
    let mut scenario = ts::begin(OWNER);
    let ctx = scenario.ctx();

    state_object::test_init(ctx);
    scenario.next_tx(OWNER);

    let owner_cap = scenario.take_from_sender<OwnerCap>();
    let mut ref = scenario.take_shared<CCIPObjectRef>();

    upgrade_registry::initialize(&mut ref, &owner_cap, scenario.ctx());
    registry::initialize(&mut ref, &owner_cap, scenario.ctx());
    registry::initialize_local_decimals(&mut ref, &owner_cap, scenario.ctx());

    offramp_state_helper::test_init(scenario.ctx());
    scenario.next_tx(OWNER);
    let dest_cap = scenario.take_from_sender<DestTransferCap>();

    (scenario, owner_cap, ref, dest_cap)
}

fun cleanup_test(
    scenario: Scenario,
    owner_cap: OwnerCap,
    ref: CCIPObjectRef,
    dest_cap: DestTransferCap,
) {
    ts::return_to_sender(&scenario, owner_cap);
    ts::return_shared(ref);
    transfer::public_transfer(dest_cap, @0x0);
    ts::end(scenario);
}

#[test]
public fun test_token_message_18_to_9_receiver_amount_equals_local() {
    let (mut scenario, owner_cap, mut ref, dest_cap) = setup_test();

    registry::register_pool_as_owner_v2(
        &owner_cap,
        &mut ref,
        TOKEN_ADDRESS,
        TOKEN_POOL_ADDRESS,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::with_defining_ids<TestTypeProof>()),
        vector<address>[],
        vector<address>[],
        LOCAL_DECIMALS,
        scenario.ctx(),
    );

    let mut receiver_params = offramp_state_helper::create_receiver_params(
        &dest_cap,
        SOURCE_CHAIN_SELECTOR,
    );

    offramp_state_helper::add_dest_token_transfer(
        &dest_cap,
        &mut receiver_params,
        RECEIVER_ADDRESS,
        SOURCE_CHAIN_SELECTOR,
        SOURCE_AMOUNT_18_DEC,
        TOKEN_ADDRESS,
        TOKEN_POOL_ADDRESS,
        b"source_pool_address",
        SOURCE_POOL_DATA_18_DEC,
        b"offchain_data",
    );

    let dest_token_amounts = client::new_dest_token_amounts(
        vector[TOKEN_ADDRESS],
        vector[SOURCE_AMOUNT_18_DEC],
    );
    let test_message = client::new_any2sui_message(
        b"message_id_32_bytes_long_test_msg",
        SOURCE_CHAIN_SELECTOR,
        b"sender_address",
        b"test_data",
        @0x5432,
        @0x12345,
        dest_token_amounts,
    );
    offramp_state_helper::populate_message(&dest_cap, &mut receiver_params, test_message);

    offramp_state_helper::complete_token_transfer(
        &ref,
        &mut receiver_params,
        TestTypeProof {},
    );

    let message = offramp_state_helper::extract_any2sui_message(&mut receiver_params);
    let (_, _, _, _, _, _, extracted_dest_token_amounts) = client::consume_any2sui_message(
        message,
        @0x5432,
    );

    assert!(client::get_amount(&extracted_dest_token_amounts[0]) == EXPECTED_LOCAL_AMOUNT_9_DEC);

    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params);
    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
public fun test_token_only_path_unchanged() {
    let (mut scenario, owner_cap, mut ref, dest_cap) = setup_test();

    registry::register_pool_as_owner_v2(
        &owner_cap,
        &mut ref,
        TOKEN_ADDRESS,
        TOKEN_POOL_ADDRESS,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::with_defining_ids<TestTypeProof>()),
        vector<address>[],
        vector<address>[],
        LOCAL_DECIMALS,
        scenario.ctx(),
    );

    let mut receiver_params = offramp_state_helper::create_receiver_params(
        &dest_cap,
        SOURCE_CHAIN_SELECTOR,
    );

    offramp_state_helper::add_dest_token_transfer(
        &dest_cap,
        &mut receiver_params,
        RECEIVER_ADDRESS,
        SOURCE_CHAIN_SELECTOR,
        SOURCE_AMOUNT_18_DEC,
        TOKEN_ADDRESS,
        TOKEN_POOL_ADDRESS,
        b"source_pool_address",
        SOURCE_POOL_DATA_18_DEC,
        b"offchain_data",
    );

    offramp_state_helper::complete_token_transfer(
        &ref,
        &mut receiver_params,
        TestTypeProof {},
    );

    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params);
    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
public fun test_message_only_path_unchanged() {
    let (scenario, owner_cap, ref, dest_cap) = setup_test();

    let mut receiver_params = offramp_state_helper::create_receiver_params(
        &dest_cap,
        SOURCE_CHAIN_SELECTOR,
    );

    let test_message = client::new_any2sui_message(
        b"message_id_32_bytes_long_test_msg",
        SOURCE_CHAIN_SELECTOR,
        b"sender_address",
        b"test_data",
        @0x5432,
        @0x12345,
        vector[],
    );
    offramp_state_helper::populate_message(&dest_cap, &mut receiver_params, test_message);

    let message = offramp_state_helper::extract_any2sui_message(&mut receiver_params);
    client::consume_any2sui_message(message, @0x5432);

    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params);
    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
public fun test_same_decimals_no_conversion() {
    let (mut scenario, owner_cap, mut ref, dest_cap) = setup_test();

    registry::register_pool_as_owner_v2(
        &owner_cap,
        &mut ref,
        TOKEN_ADDRESS,
        TOKEN_POOL_ADDRESS,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::with_defining_ids<TestTypeProof>()),
        vector<address>[],
        vector<address>[],
        9,
        scenario.ctx(),
    );

    let mut receiver_params = offramp_state_helper::create_receiver_params(
        &dest_cap,
        SOURCE_CHAIN_SELECTOR,
    );

    let source_amount: u256 = 1_000_000_000;
    // ABI-encoded u256(9) = remote decimals same as local
    let source_pool_data_9_dec =
        x"0000000000000000000000000000000000000000000000000000000000000009";

    offramp_state_helper::add_dest_token_transfer(
        &dest_cap,
        &mut receiver_params,
        RECEIVER_ADDRESS,
        SOURCE_CHAIN_SELECTOR,
        source_amount,
        TOKEN_ADDRESS,
        TOKEN_POOL_ADDRESS,
        b"source_pool_address",
        source_pool_data_9_dec,
        b"offchain_data",
    );

    let dest_token_amounts = client::new_dest_token_amounts(
        vector[TOKEN_ADDRESS],
        vector[source_amount],
    );
    let test_message = client::new_any2sui_message(
        b"message_id_32_bytes_long_test_msg",
        SOURCE_CHAIN_SELECTOR,
        b"sender_address",
        b"test_data",
        @0x5432,
        @0x12345,
        dest_token_amounts,
    );
    offramp_state_helper::populate_message(&dest_cap, &mut receiver_params, test_message);

    offramp_state_helper::complete_token_transfer(
        &ref,
        &mut receiver_params,
        TestTypeProof {},
    );

    let message = offramp_state_helper::extract_any2sui_message(&mut receiver_params);
    let (_, _, _, _, _, _, extracted_dest_token_amounts) = client::consume_any2sui_message(
        message,
        @0x5432,
    );

    assert!(client::get_amount(&extracted_dest_token_amounts[0]) == source_amount);

    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params);
    cleanup_test(scenario, owner_cap, ref, dest_cap);
}
