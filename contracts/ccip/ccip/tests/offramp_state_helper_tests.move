#[test_only]
module ccip::offramp_state_helper_tests;

use std::ascii;
use std::string;
use std::type_name;

use sui::coin;
use sui::test_scenario::{Self as ts, Scenario};

use ccip::client;
use ccip::offramp_state_helper::{Self, DestTransferCap};
use ccip::receiver_registry;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::ownable::OwnerCap;
use ccip::token_admin_registry as registry;

public struct OFFRAMP_STATE_HELPER_TESTS has drop {}

public struct TestTypeProof has drop {}
public struct TestTypeProof2 has drop {}
public struct TestToken has drop {}
public struct TestToken2 has drop {}

const OWNER: address = @0x1000;
const RECEIVER_ADDRESS: address = @0x2000;
const TOKEN_ADDRESS_1: address = @0x1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b;
const TOKEN_ADDRESS_2: address = @0x8a7b6c5d4e3f2a1b0c9d8e7f6a5b4c3d2e1f0a9b8c7d6e5f4a3b2c1d0e9f8a7;
const TOKEN_POOL_ADDRESS_1: address = @0xdeeb7a4662eec9f2f3def03fb937a663dddaa2e215b8078a284d026b7946c270;
const TOKEN_POOL_ADDRESS_2: address = @0xd8908c165dee785924e7421a0fd0418a19d5daeec395fd505a92a0fd3117e428;
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

fun cleanup_test(scenario: Scenario, owner_cap: OwnerCap, ref: CCIPObjectRef, dest_cap: DestTransferCap) {
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
    let receiver_params = offramp_state_helper::create_receiver_params(&dest_cap, SOURCE_CHAIN_SELECTOR);
    let source_chain = offramp_state_helper::get_source_chain_selector(&receiver_params);
    assert!(source_chain == SOURCE_CHAIN_SELECTOR);
    
    // Clean up
    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params, vector[]);
    
    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
public fun test_get_source_chain_selector() {
    let (scenario, owner_cap, ref, dest_cap) = setup_test();
    
    let different_chain = 2000;
    let receiver_params = offramp_state_helper::create_receiver_params(&dest_cap, different_chain);
    let source_chain = offramp_state_helper::get_source_chain_selector(&receiver_params);
    assert!(source_chain == different_chain);
    
    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params, vector[]);
    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
public fun test_add_dest_token_transfer() {
    let (scenario, owner_cap, ref, dest_cap) = setup_test();
    
    let mut receiver_params = offramp_state_helper::create_receiver_params(&dest_cap, SOURCE_CHAIN_SELECTOR);
    
    // Add a destination token transfer
    offramp_state_helper::add_dest_token_transfer(
        &dest_cap,
        &mut receiver_params,
        RECEIVER_ADDRESS,
        1000, // source_amount
        TOKEN_ADDRESS_1,
        TOKEN_POOL_ADDRESS_1,
        b"source_pool_address",
        b"source_pool_data",
        b"offchain_data"
    );
    
    // Verify the token transfer was added
    let (receiver, source_amount, dest_token_address, source_pool_address, source_pool_data, offchain_data) = 
        offramp_state_helper::get_token_param_data(&receiver_params, 0);
    
    assert!(receiver == RECEIVER_ADDRESS);
    assert!(source_amount == 1000);
    assert!(dest_token_address == TOKEN_ADDRESS_1);
    assert!(source_pool_address == b"source_pool_address");
    assert!(source_pool_data == b"source_pool_data");
    assert!(offchain_data == b"offchain_data");
    
    // This will fail but we need to call it to consume receiver_params
    offramp_state_helper::deconstruct_receiver_params_for_test(&dest_cap, receiver_params);
    
    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
public fun test_add_multiple_dest_token_transfers() {
    let (scenario, owner_cap, ref, dest_cap) = setup_test();
    
    let mut receiver_params = offramp_state_helper::create_receiver_params(&dest_cap, SOURCE_CHAIN_SELECTOR);
    
    // Add first token transfer
    offramp_state_helper::add_dest_token_transfer(
        &dest_cap,
        &mut receiver_params,
        RECEIVER_ADDRESS,
        1000,
        TOKEN_ADDRESS_1,
        TOKEN_POOL_ADDRESS_1,
        b"source_pool_1",
        b"pool_data_1",
        b"offchain_1"
    );
    
    // Add second token transfer
    offramp_state_helper::add_dest_token_transfer(
        &dest_cap,
        &mut receiver_params,
        RECEIVER_ADDRESS,
        2000,
        TOKEN_ADDRESS_2,
        TOKEN_POOL_ADDRESS_2,
        b"source_pool_2",
        b"pool_data_2",
        b"offchain_2"
    );
    
    // Verify first transfer
    let (receiver1, source_amount1, dest_token_address1, source_pool_address1, source_pool_data1, offchain_data1) = 
        offramp_state_helper::get_token_param_data(&receiver_params, 0);
    
    assert!(receiver1 == RECEIVER_ADDRESS);
    assert!(source_amount1 == 1000);
    assert!(dest_token_address1 == TOKEN_ADDRESS_1);
    assert!(source_pool_address1 == b"source_pool_1");
    assert!(source_pool_data1 == b"pool_data_1");
    assert!(offchain_data1 == b"offchain_1");
    
    // Verify second transfer
    let (receiver2, source_amount2, dest_token_address2, source_pool_address2, source_pool_data2, offchain_data2) = 
        offramp_state_helper::get_token_param_data(&receiver_params, 1);
    
    assert!(receiver2 == RECEIVER_ADDRESS);
    assert!(source_amount2 == 2000);
    assert!(dest_token_address2 == TOKEN_ADDRESS_2);
    assert!(source_pool_address2 == b"source_pool_2");
    assert!(source_pool_data2 == b"pool_data_2");
    assert!(offchain_data2 == b"offchain_2");
    
    // We need to consume the receiver_params with incomplete transfers
    // This will fail but we need to call it to consume receiver_params
    offramp_state_helper::deconstruct_receiver_params_for_test(&dest_cap, receiver_params);
    
    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
#[expected_failure(abort_code = offramp_state_helper::EWrongIndexInReceiverParams)]
public fun test_get_token_param_data_wrong_index() {
    let (scenario, owner_cap, ref, dest_cap) = setup_test();
    
    let receiver_params = offramp_state_helper::create_receiver_params(&dest_cap, SOURCE_CHAIN_SELECTOR);
    
    // Try to access index 0 when no transfers have been added
    let (_receiver, _source_amount, _dest_token_address, _source_pool_address, _source_pool_data, _offchain_data) = 
        offramp_state_helper::get_token_param_data(&receiver_params, 0);
    
    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params, vector[]);
    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
public fun test_populate_message() {
    let (scenario, owner_cap, ref, dest_cap) = setup_test();
    
    let mut receiver_params = offramp_state_helper::create_receiver_params(&dest_cap, SOURCE_CHAIN_SELECTOR);
    
    // Create a test message
    let test_message = client::new_any2sui_message(
        b"message_id_32_bytes_long_test_msg",
        SOURCE_CHAIN_SELECTOR,
        b"sender_address",
        b"test_data",
        vector[] // token_amounts
    );
    
    // Populate the message
    offramp_state_helper::populate_message(&dest_cap, &mut receiver_params, test_message);
    
    // We need to consume the receiver_params with a message
    // This will fail but we need to call it to consume receiver_params
    offramp_state_helper::deconstruct_receiver_params_for_test(&dest_cap, receiver_params);
    
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
        vector[], // lock_or_burn_params
        vector[], // release_or_mint_params
        scenario.ctx(),
    );
    
    let mut receiver_params = offramp_state_helper::create_receiver_params(&dest_cap, SOURCE_CHAIN_SELECTOR);
    
    // Add a destination token transfer
    offramp_state_helper::add_dest_token_transfer(
        &dest_cap,
        &mut receiver_params,
        RECEIVER_ADDRESS,
        1000,
        TOKEN_ADDRESS_1,
        TOKEN_POOL_ADDRESS_1,
        b"source_pool_address",
        b"source_pool_data",
        b"offchain_data"
    );
    
    // Create a test coin to transfer
    let test_coin = coin::mint_for_testing<TestToken>(500, scenario.ctx());

    // let local_amount = coin::value(&test_coin);
    // Complete the token transfer
    let completed_transfer = offramp_state_helper::complete_token_transfer(
        &ref,
        &mut receiver_params,
        0, // index
        TestTypeProof {}
    );
    
    // Destroy the unused test_coin
    coin::burn_for_testing(test_coin);
    
    // Clean up - the receiver_params should have completed transfers
    let completed_transfers = vector[completed_transfer];
    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params, completed_transfers);
    
    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
public fun test_extract_any2sui_message() {
    let (scenario, owner_cap, mut ref, dest_cap) = setup_test();
    
    // Register a receiver
    receiver_registry::register_receiver(
        &mut ref,
        @0x123, // receiver_state_id
        vector[], // receiver_state_params
        TestTypeProof {}
    );
    
    let mut receiver_params = offramp_state_helper::create_receiver_params(&dest_cap, SOURCE_CHAIN_SELECTOR);
    
    // Create and populate a message
    let test_message = client::new_any2sui_message(
        b"message_id_32_bytes_long_test_msg",
        SOURCE_CHAIN_SELECTOR,
        b"sender_address",
        b"test_data",
        vector[]
    );
    offramp_state_helper::populate_message(&dest_cap, &mut receiver_params, test_message);
    
    // Extract the message
    let (extracted_message, receiver_params) = offramp_state_helper::extract_any2sui_message(
        &ref,
        receiver_params,
        TestTypeProof {}
    );
    
    // Verify message was extracted
    assert!(extracted_message.is_some());
    
    // Clean up
    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params, vector[]);
    
    cleanup_test(scenario, owner_cap, ref, dest_cap);
}

#[test]
public fun test_deconstruct_receiver_params_empty() {
    let (scenario, owner_cap, ref, dest_cap) = setup_test();
    
    // Create empty receiver params
    let receiver_params = offramp_state_helper::create_receiver_params(&dest_cap, SOURCE_CHAIN_SELECTOR);
    
    // Should succeed with no token transfers and no message
    offramp_state_helper::deconstruct_receiver_params(&dest_cap, receiver_params, vector[]);
    
    cleanup_test(scenario, owner_cap, ref, dest_cap);
} 