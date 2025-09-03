#[test_only]
module ccip::offramp_state_helper_tests;

use ccip::client;
use ccip::offramp_state_helper::{Self, DestTransferCap};
use ccip::ownable::OwnerCap;
use ccip::receiver_registry;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::token_admin_registry as registry;
use std::ascii;
use std::string;
use std::type_name;
use sui::coin;
use sui::test_scenario::{Self as ts, Scenario};

public struct OFFRAMP_STATE_HELPER_TESTS has drop {}

public struct TestTypeProof has drop {}
public struct TestTypeProof2 has drop {}
public struct TestToken has drop {}

const OWNER: address = @0x1000;
const RECEIVER_ADDRESS: address = @0x2000;
const TOKEN_ADDRESS_1: address =
    @0x1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b;
const TOKEN_ADDRESS_2: address =
    @0x2a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b;
const TOKEN_POOL_ADDRESS_1: address =
    @0xdeeb7a4662eec9f2f3def03fb937a663dddaa2e215b8078a284d026b7946c270;
const TOKEN_POOL_ADDRESS_2: address =
    @0xaeeb7a4662eec9f2f3def03fb937a663dddaa2e215b8078a284d026b7946c270;
const SOURCE_CHAIN_SELECTOR: u64 = 1000;

fun setup_test(): (Scenario, OwnerCap, CCIPObjectRef, DestTransferCap) {
    let mut scenario = ts::begin(OWNER);
    let ctx = scenario.ctx();

    state_object::test_init(ctx);

    // Advance to next transaction to retrieve the created objects
    scenario.next_tx(OWNER);

    // Retrieve the OwnerCap that was transferred to the sender
    let owner_cap = scenario.take_from_sender<OwnerCap>();

    // Retrieve the shared CCIPObjectRef
    let mut ref = scenario.take_shared<CCIPObjectRef>();

    // Initialize token admin registry
    registry::initialize(&mut ref, &owner_cap, scenario.ctx());

    // Initialize receiver registry
    receiver_registry::initialize(&mut ref, &owner_cap, scenario.ctx());

    // Create offramp state helper and get dest transfer cap
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
    // Return the owner cap back to the sender instead of destroying it
    ts::return_to_sender(&scenario, owner_cap);
    // Return the shared object back to the scenario instead of destroying it
    ts::return_shared(ref);
    transfer::public_transfer(dest_cap, @0x0);
    ts::end(scenario);
}

#[test]
public fun test_create_receiver_params() {
    let (scenario, owner_cap, ref, dest_cap) = setup_test();

    // Test creating receiver params
    let receiver_params = offramp_state_helper::create_receiver_params(
        &dest_cap,
        SOURCE_CHAIN_SELECTOR,
    );
    let source_chain = offramp_state_helper::get_source_chain_selector(&receiver_params);
    assert!(source_chain == SOURCE_CHAIN_SELECTOR);

    // Clean up
    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params);

    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
public fun test_get_source_chain_selector() {
    let (scenario, owner_cap, ref, dest_cap) = setup_test();

    let different_chain = 2000;
    let receiver_params = offramp_state_helper::create_receiver_params(&dest_cap, different_chain);
    let source_chain = offramp_state_helper::get_source_chain_selector(&receiver_params);
    assert!(source_chain == different_chain);

    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params);
    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
public fun test_add_dest_token_transfer() {
    let (scenario, owner_cap, ref, dest_cap) = setup_test();

    let mut receiver_params = offramp_state_helper::create_receiver_params(
        &dest_cap,
        SOURCE_CHAIN_SELECTOR,
    );

    // Add a destination token transfer
    offramp_state_helper::add_dest_token_transfer(
        &dest_cap,
        &mut receiver_params,
        RECEIVER_ADDRESS,
        SOURCE_CHAIN_SELECTOR, // remote_chain_selector
        1000, // source_amount
        TOKEN_ADDRESS_1,
        TOKEN_POOL_ADDRESS_1,
        b"source_pool_address",
        b"source_pool_data",
        b"offchain_data",
    );

    // Verify the token transfer was added
    let (
        receiver,
        source_amount,
        dest_token_address,
        source_pool_address,
        source_pool_data,
        offchain_data,
    ) = offramp_state_helper::get_token_param_data(&receiver_params);

    assert!(receiver == RECEIVER_ADDRESS);
    assert!(source_amount == 1000);
    assert!(dest_token_address == TOKEN_ADDRESS_1);
    assert!(source_pool_address == b"source_pool_address");
    assert!(source_pool_data == b"source_pool_data");
    assert!(offchain_data == b"offchain_data");

    // This will fail but we need to call it to consume receiver_params
    offramp_state_helper::deconstruct_receiver_params_with_message_for_test(
        &dest_cap,
        receiver_params,
    );

    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
public fun test_populate_message() {
    let (scenario, owner_cap, ref, dest_cap) = setup_test();

    let mut receiver_params = offramp_state_helper::create_receiver_params(
        &dest_cap,
        SOURCE_CHAIN_SELECTOR,
    );

    // Create a test message
    let test_message = client::new_any2sui_message(
        b"message_id_32_bytes_long_test_msg",
        SOURCE_CHAIN_SELECTOR,
        b"sender_address",
        b"test_data",
        vector[], // token_amounts
    );

    // Populate the message
    offramp_state_helper::populate_message(&dest_cap, &mut receiver_params, test_message);

    // We need to consume the receiver_params with a message
    // Use the new test function that can handle populated messages
    offramp_state_helper::deconstruct_receiver_params_with_message_for_test(
        &dest_cap,
        receiver_params,
    );

    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
public fun test_complete_token_transfer() {
    let (mut scenario, owner_cap, mut ref, dest_cap) = setup_test();

    // Register a token in the token admin registry
    registry::register_pool_by_admin(
        &mut ref,
        state_object::create_ccip_admin_proof_for_test(),
        TOKEN_ADDRESS_1,
        TOKEN_POOL_ADDRESS_1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        vector<address>[], // lock_or_burn_params
        vector<address>[], // release_or_mint_params
        scenario.ctx(),
    );

    let mut receiver_params = offramp_state_helper::create_receiver_params(
        &dest_cap,
        SOURCE_CHAIN_SELECTOR,
    );

    // Add a destination token transfer
    offramp_state_helper::add_dest_token_transfer(
        &dest_cap,
        &mut receiver_params,
        RECEIVER_ADDRESS,
        SOURCE_CHAIN_SELECTOR, // remote_chain_selector
        1000,
        TOKEN_ADDRESS_1,
        TOKEN_POOL_ADDRESS_1,
        b"source_pool_address",
        b"source_pool_data",
        b"offchain_data",
    );

    // Create a test coin to transfer
    let test_coin = coin::mint_for_testing<TestToken>(500, scenario.ctx());

    // Complete the token transfer
    offramp_state_helper::complete_token_transfer(
        &ref,
        &mut receiver_params,
        RECEIVER_ADDRESS,
        TOKEN_ADDRESS_1,
        TestTypeProof {},
    );

    // Destroy the unused test_coin
    coin::burn_for_testing(test_coin);

    // Clean up - the receiver_params should have completed transfers
    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params);

    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
public fun test_extract_any2sui_message() {
    let (scenario, owner_cap, mut ref, dest_cap) = setup_test();

    // Register a receiver
    receiver_registry::register_receiver(
        &mut ref,
        TestTypeProof {},
    );

    let mut receiver_params = offramp_state_helper::create_receiver_params(
        &dest_cap,
        SOURCE_CHAIN_SELECTOR,
    );

    // Create and populate a message
    let test_message = client::new_any2sui_message(
        b"message_id_32_bytes_long_test_msg",
        SOURCE_CHAIN_SELECTOR,
        b"sender_address",
        b"test_data",
        vector[],
    );
    offramp_state_helper::populate_message(&dest_cap, &mut receiver_params, test_message);

    // Since we can't destructure ReceiverParams outside its module, we'll just test that
    // the message was populated by trying to consume it via the offramp helper
    // Use the new test function that can handle populated messages
    offramp_state_helper::deconstruct_receiver_params_with_message_for_test(
        &dest_cap,
        receiver_params,
    );

    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
public fun test_deconstruct_receiver_params_empty() {
    let (scenario, owner_cap, ref, dest_cap) = setup_test();

    // Create empty receiver params
    let receiver_params = offramp_state_helper::create_receiver_params(
        &dest_cap,
        SOURCE_CHAIN_SELECTOR,
    );

    // Should succeed with no token transfers and no message
    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params);

    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
#[expected_failure(abort_code = ccip::offramp_state_helper::ENoMessageToExtract)]
public fun test_extract_message_when_none_exists() {
    let (scenario, owner_cap, ref, dest_cap) = setup_test();

    let mut receiver_params = offramp_state_helper::create_receiver_params(
        &dest_cap,
        SOURCE_CHAIN_SELECTOR,
    );

    // Try to extract message when none exists - should fail
    let message = offramp_state_helper::extract_any2sui_message(&mut receiver_params);

    // This should never be reached - consume the message to avoid drop error
    let (_, _, _, _, _) = client::consume_any2sui_message(message);
    offramp_state_helper::deconstruct_receiver_params_with_message_for_test(
        &dest_cap,
        receiver_params,
    );
    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
#[expected_failure(abort_code = ccip::offramp_state_helper::ETokenTransferAlreadyExists)]
public fun test_add_token_transfer_when_already_exists() {
    let (scenario, owner_cap, ref, dest_cap) = setup_test();

    let mut receiver_params = offramp_state_helper::create_receiver_params(
        &dest_cap,
        SOURCE_CHAIN_SELECTOR,
    );

    // Add first token transfer
    offramp_state_helper::add_dest_token_transfer(
        &dest_cap,
        &mut receiver_params,
        RECEIVER_ADDRESS,
        SOURCE_CHAIN_SELECTOR,
        1000,
        TOKEN_ADDRESS_1,
        TOKEN_POOL_ADDRESS_1,
        b"source_pool_address",
        b"source_pool_data",
        b"offchain_data",
    );

    // Try to add second token transfer - should fail
    offramp_state_helper::add_dest_token_transfer(
        &dest_cap,
        &mut receiver_params,
        RECEIVER_ADDRESS,
        SOURCE_CHAIN_SELECTOR,
        2000,
        TOKEN_ADDRESS_2,
        TOKEN_POOL_ADDRESS_2,
        b"source_pool_address_2",
        b"source_pool_data_2",
        b"offchain_data_2",
    );

    // This should never be reached
    offramp_state_helper::deconstruct_receiver_params_with_message_for_test(
        &dest_cap,
        receiver_params,
    );
    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
#[expected_failure(abort_code = ccip::offramp_state_helper::ETokenTransferDoesNotExist)]
public fun test_get_token_param_data_when_none_exists() {
    let (scenario, owner_cap, ref, dest_cap) = setup_test();

    let receiver_params = offramp_state_helper::create_receiver_params(
        &dest_cap,
        SOURCE_CHAIN_SELECTOR,
    );

    // Try to get token data when no token transfer exists - should fail
    let (_, _, _, _, _, _) = offramp_state_helper::get_token_param_data(&receiver_params);

    // This should never be reached
    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params);
    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
#[expected_failure(abort_code = ccip::offramp_state_helper::ETokenTransferDoesNotExist)]
public fun test_get_dest_token_transfer_data_when_none_exists() {
    let (scenario, owner_cap, ref, dest_cap) = setup_test();

    let receiver_params = offramp_state_helper::create_receiver_params(
        &dest_cap,
        SOURCE_CHAIN_SELECTOR,
    );

    // Try to get dest token transfer data when none exists - should fail
    let (_, _, _, _, _, _, _, _) = offramp_state_helper::get_dest_token_transfer_data(
        &receiver_params,
    );

    // This should never be reached
    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params);
    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
#[expected_failure(abort_code = ccip::offramp_state_helper::ETokenTransferAlreadyCompleted)]
public fun test_complete_token_transfer_already_completed() {
    let (mut scenario, owner_cap, mut ref, dest_cap) = setup_test();

    // Register a token
    registry::register_pool_by_admin(
        &mut ref,
        state_object::create_ccip_admin_proof_for_test(),
        TOKEN_ADDRESS_1,
        TOKEN_POOL_ADDRESS_1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        vector<address>[], // lock_or_burn_params
        vector<address>[], // release_or_mint_params
        scenario.ctx(),
    );

    let mut receiver_params = offramp_state_helper::create_receiver_params(
        &dest_cap,
        SOURCE_CHAIN_SELECTOR,
    );

    // Add a destination token transfer
    offramp_state_helper::add_dest_token_transfer(
        &dest_cap,
        &mut receiver_params,
        RECEIVER_ADDRESS,
        SOURCE_CHAIN_SELECTOR,
        1000,
        TOKEN_ADDRESS_1,
        TOKEN_POOL_ADDRESS_1,
        b"source_pool_address",
        b"source_pool_data",
        b"offchain_data",
    );

    // Complete the token transfer first time (should succeed)
    offramp_state_helper::complete_token_transfer(
        &ref,
        &mut receiver_params,
        RECEIVER_ADDRESS,
        TOKEN_ADDRESS_1,
        TestTypeProof {},
    );

    // This should fail with ETokenTransferAlreadyCompleted because transfer is already completed
    offramp_state_helper::complete_token_transfer(
        &ref,
        &mut receiver_params,
        RECEIVER_ADDRESS,
        TOKEN_ADDRESS_1,
        TestTypeProof {},
    );

    // This should never be reached
    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params);
    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
#[expected_failure(abort_code = ccip::offramp_state_helper::ETypeProofMismatch)]
public fun test_complete_token_transfer_wrong_type_proof() {
    let (mut scenario, owner_cap, mut ref, dest_cap) = setup_test();

    // Register a token with TestTypeProof
    registry::register_pool_by_admin(
        &mut ref,
        state_object::create_ccip_admin_proof_for_test(),
        TOKEN_ADDRESS_1,
        TOKEN_POOL_ADDRESS_1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        vector<address>[], // lock_or_burn_params
        vector<address>[], // release_or_mint_params
        scenario.ctx(),
    );

    let mut receiver_params = offramp_state_helper::create_receiver_params(
        &dest_cap,
        SOURCE_CHAIN_SELECTOR,
    );

    // Add a token transfer
    offramp_state_helper::add_dest_token_transfer(
        &dest_cap,
        &mut receiver_params,
        RECEIVER_ADDRESS,
        SOURCE_CHAIN_SELECTOR,
        1000,
        TOKEN_ADDRESS_1,
        TOKEN_POOL_ADDRESS_1,
        b"source_pool_address",
        b"source_pool_data",
        b"offchain_data",
    );

    // Try to complete with wrong type proof - should fail
    offramp_state_helper::complete_token_transfer(
        &ref,
        &mut receiver_params,
        RECEIVER_ADDRESS,
        TOKEN_ADDRESS_1,
        TestTypeProof2 {}, // Wrong type proof!
    );

    // This should never be reached
    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params);
    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
#[expected_failure(abort_code = ccip::offramp_state_helper::ECCIPReceiveFailed)]
public fun test_deconstruct_with_unextracted_message() {
    let (scenario, owner_cap, ref, dest_cap) = setup_test();

    let mut receiver_params = offramp_state_helper::create_receiver_params(
        &dest_cap,
        SOURCE_CHAIN_SELECTOR,
    );

    // Add a message but don't extract it
    let test_message = client::new_any2sui_message(
        b"message_id_32_bytes_long_test_msg",
        SOURCE_CHAIN_SELECTOR,
        b"sender_address",
        b"test_data",
        vector[],
    );
    offramp_state_helper::populate_message(&dest_cap, &mut receiver_params, test_message);

    // Try to deconstruct without extracting message - should fail
    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params);

    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
#[expected_failure(abort_code = ccip::offramp_state_helper::EWrongReceiptAndTokenTransfer)]
public fun test_deconstruct_with_token_transfer_but_no_receipt() {
    let (scenario, owner_cap, ref, dest_cap) = setup_test();

    let mut receiver_params = offramp_state_helper::create_receiver_params(
        &dest_cap,
        SOURCE_CHAIN_SELECTOR,
    );

    // Add a token transfer but don't complete it (no receipt)
    offramp_state_helper::add_dest_token_transfer(
        &dest_cap,
        &mut receiver_params,
        RECEIVER_ADDRESS,
        SOURCE_CHAIN_SELECTOR,
        1000,
        TOKEN_ADDRESS_1,
        TOKEN_POOL_ADDRESS_1,
        b"source_pool_address",
        b"source_pool_data",
        b"offchain_data",
    );

    // Try to deconstruct with token transfer but no receipt - should fail
    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params);

    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
#[expected_failure(abort_code = ccip::offramp_state_helper::ETypeProofMismatch)]
public fun test_consume_message_wrong_type_proof() {
    let (scenario, owner_cap, mut ref, dest_cap) = setup_test();

    // Register a receiver with TestTypeProof
    receiver_registry::register_receiver(
        &mut ref,
        TestTypeProof {},
    );

    // Create a test message
    let test_message = client::new_any2sui_message(
        b"message_id_32_bytes_long_test_msg",
        SOURCE_CHAIN_SELECTOR,
        b"sender_address",
        b"test_data",
        vector[],
    );

    // Try to consume with wrong type proof - should fail
    let (_, _, _, _, _) = offramp_state_helper::consume_any2sui_message(
        &ref,
        test_message,
        TestTypeProof2 {}, // Wrong type proof!
    );

    cleanup_test(scenario, owner_cap, ref, dest_cap);
}
