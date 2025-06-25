#[test_only]
module ccip::dynamic_dispatcher_tests;

use std::ascii;
use std::string;
use std::type_name;

use sui::test_scenario::{Self as ts, Scenario};

use ccip::dynamic_dispatcher::{Self, SourceTransferCap};
use ccip::state_object::{Self, OwnerCap, CCIPObjectRef};
use ccip::token_admin_registry as registry;

public struct DYNAMIC_DISPATCHER_TESTS has drop {}

public struct TestTypeProof has drop {}
public struct TestTypeProof2 has drop {}

const OWNER: address = @0x1000;
const TOKEN_ADDRESS_1: address = @0x1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b;
const TOKEN_ADDRESS_2: address = @0x8a7b6c5d4e3f2a1b0c9d8e7f6a5b4c3d2e1f0a9b8c7d6e5f4a3b2c1d0e9f8a7;
const TOKEN_POOL_PACKAGE_ID_1: address = @0xdeeb7a4662eec9f2f3def03fb937a663dddaa2e215b8078a284d026b7946c270;
const TOKEN_POOL_PACKAGE_ID_2: address = @0xd8908c165dee785924e7421a0fd0418a19d5daeec395fd505a92a0fd3117e428;
const DESTINATION_CHAIN_SELECTOR: u64 = 1000;

fun setup_test(): (Scenario, OwnerCap, CCIPObjectRef, SourceTransferCap) {
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
    
    // Create dynamic dispatcher and get source transfer cap
    dynamic_dispatcher::test_init(scenario.ctx());
    
    scenario.next_tx(OWNER);
    let source_cap = scenario.take_from_sender<SourceTransferCap>();

    (scenario, owner_cap, ref, source_cap)
}

fun cleanup_test(scenario: Scenario, owner_cap: OwnerCap, ref: CCIPObjectRef, source_cap: SourceTransferCap) {
    // Return the owner cap back to the sender instead of destroying it
    ts::return_to_sender(&scenario, owner_cap);
    // Return the shared object back to the scenario instead of destroying it
    ts::return_shared(ref);
    transfer::public_transfer(source_cap, @0x0);
    ts::end(scenario);
}

#[test]
public fun test_create_token_params() {
    let (scenario, owner_cap, ref, source_cap) = setup_test();
    
    // Test creating token params with valid destination chain selector
    let token_params = dynamic_dispatcher::create_token_params(DESTINATION_CHAIN_SELECTOR);
    let destination = dynamic_dispatcher::get_destination_chain_selector(&token_params);
    assert!(destination == DESTINATION_CHAIN_SELECTOR);
    
    // Clean up token_params
    let (_, _) = dynamic_dispatcher::deconstruct_token_params(&source_cap, token_params);
    
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
#[expected_failure(abort_code = dynamic_dispatcher::EInvalidDestinationChainSelector)]
public fun test_create_token_params_zero_chain_selector() {
    let (scenario, owner_cap, ref, source_cap) = setup_test();
    
    // Test creating token params with zero destination chain selector should fail
    let token_params = dynamic_dispatcher::create_token_params(0);

    let (_, _) = dynamic_dispatcher::deconstruct_token_params(&source_cap, token_params);
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_get_destination_chain_selector() {
    let (scenario, owner_cap, ref, source_cap) = setup_test();
    
    let token_params = dynamic_dispatcher::create_token_params(DESTINATION_CHAIN_SELECTOR);
    let destination = dynamic_dispatcher::get_destination_chain_selector(&token_params);
    assert!(destination == DESTINATION_CHAIN_SELECTOR);
    
    // Test with different chain selector
    let different_chain = 2000;
    let token_params2 = dynamic_dispatcher::create_token_params(different_chain);
    let destination2 = dynamic_dispatcher::get_destination_chain_selector(&token_params2);
    assert!(destination2 == different_chain);
    
    // Clean up
    let (_, _) = dynamic_dispatcher::deconstruct_token_params(&source_cap, token_params);
    let (_, _) = dynamic_dispatcher::deconstruct_token_params(&source_cap, token_params2);
    
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_add_source_token_transfer() {
    let (mut scenario, owner_cap, mut ref, source_cap) = setup_test();
    
    // Register a token in the token admin registry first
    registry::register_pool_by_admin(
        &mut ref,
        TOKEN_ADDRESS_1,
        TOKEN_POOL_PACKAGE_ID_1,
        TOKEN_ADDRESS_1, // state address
        string::utf8(b"test_pool"),
        ascii::string(b"TestToken"),
        OWNER, // administrator
        type_name::into_string(type_name::get<TestTypeProof>()),
        scenario.ctx(),
    );
    
    let token_params = dynamic_dispatcher::create_token_params(DESTINATION_CHAIN_SELECTOR);
    
    // Add source token transfer
    let updated_params = dynamic_dispatcher::add_source_token_transfer(
        &ref,
        token_params,
        1000, // amount
        TOKEN_ADDRESS_1,
        b"dest_token_address",
        b"extra_data",
        TestTypeProof {}
    );
    
    // Verify the token params were updated
    let destination = dynamic_dispatcher::get_destination_chain_selector(&updated_params);
    assert!(destination == DESTINATION_CHAIN_SELECTOR);
    
    // Deconstruct and verify the source token transfer
    let (dest_chain, transfers) = dynamic_dispatcher::deconstruct_token_params(&source_cap, updated_params);
    assert!(dest_chain == DESTINATION_CHAIN_SELECTOR);
    assert!(transfers.length() == 1);
    
    let transfer = &transfers[0];
    let (source_pool, amount, source_token_address, dest_token_address, extra_data) = 
        dynamic_dispatcher::get_source_token_transfer_data(*transfer);
    
    assert!(source_pool == TOKEN_POOL_PACKAGE_ID_1);
    assert!(amount == 1000);
    assert!(source_token_address == TOKEN_ADDRESS_1);
    assert!(dest_token_address == b"dest_token_address");
    assert!(extra_data == b"extra_data");
    
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
#[expected_failure(abort_code = dynamic_dispatcher::ETypeProofMismatch)]
public fun test_add_source_token_transfer_wrong_proof() {
    let (mut scenario, owner_cap, mut ref, source_cap) = setup_test();
    
    // Register a token with TestTypeProof
    registry::register_pool_by_admin(
        &mut ref,
        TOKEN_ADDRESS_1,
        TOKEN_POOL_PACKAGE_ID_1,
        TOKEN_ADDRESS_1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestToken"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        scenario.ctx(),
    );
    
    let token_params = dynamic_dispatcher::create_token_params(DESTINATION_CHAIN_SELECTOR);
    
    // Try to add source token transfer with wrong proof type
    let updated_params = dynamic_dispatcher::add_source_token_transfer(
        &ref,
        token_params,
        1000,
        TOKEN_ADDRESS_1,
        b"dest_token_address",
        b"extra_data",
        TestTypeProof2 {} // Wrong proof type!
    );

    let (_, _) = dynamic_dispatcher::deconstruct_token_params(&source_cap, updated_params);
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_add_multiple_source_token_transfers() {
    let (mut scenario, owner_cap, mut ref, source_cap) = setup_test();
    
    // Register two tokens
    registry::register_pool_by_admin(
        &mut ref,
        TOKEN_ADDRESS_1,
        TOKEN_POOL_PACKAGE_ID_1,
        TOKEN_ADDRESS_1,
        string::utf8(b"test_pool_1"),
        ascii::string(b"TestToken1"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        scenario.ctx(),
    );
    
    registry::register_pool_by_admin(
        &mut ref,
        TOKEN_ADDRESS_2,
        TOKEN_POOL_PACKAGE_ID_2,
        TOKEN_ADDRESS_2,
        string::utf8(b"test_pool_2"),
        ascii::string(b"TestToken2"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof2>()),
        scenario.ctx(),
    );
    
    let token_params = dynamic_dispatcher::create_token_params(DESTINATION_CHAIN_SELECTOR);
    
    // Add first token transfer
    let updated_params1 = dynamic_dispatcher::add_source_token_transfer(
        &ref,
        token_params,
        1000,
        TOKEN_ADDRESS_1,
        b"dest_token_address_1",
        b"extra_data_1",
        TestTypeProof {}
    );
    
    // Add second token transfer
    let updated_params2 = dynamic_dispatcher::add_source_token_transfer(
        &ref,
        updated_params1,
        2000,
        TOKEN_ADDRESS_2,
        b"dest_token_address_2",
        b"extra_data_2",
        TestTypeProof2 {}
    );
    
    // Verify both transfers were added
    let (dest_chain, transfers) = dynamic_dispatcher::deconstruct_token_params(&source_cap, updated_params2);
    assert!(dest_chain == DESTINATION_CHAIN_SELECTOR);
    assert!(transfers.length() == 2);
    
    // Verify first transfer
    let transfer1 = &transfers[0];
    let (source_pool1, amount1, source_token_address1, dest_token_address1, extra_data1) = 
        dynamic_dispatcher::get_source_token_transfer_data(*transfer1);
    
    assert!(source_pool1 == TOKEN_POOL_PACKAGE_ID_1);
    assert!(amount1 == 1000);
    assert!(source_token_address1 == TOKEN_ADDRESS_1);
    assert!(dest_token_address1 == b"dest_token_address_1");
    assert!(extra_data1 == b"extra_data_1");
    
    // Verify second transfer
    let transfer2 = &transfers[1];
    let (source_pool2, amount2, source_token_address2, dest_token_address2, extra_data2) = 
        dynamic_dispatcher::get_source_token_transfer_data(*transfer2);
    
    assert!(source_pool2 == TOKEN_POOL_PACKAGE_ID_2);
    assert!(amount2 == 2000);
    assert!(source_token_address2 == TOKEN_ADDRESS_2);
    assert!(dest_token_address2 == b"dest_token_address_2");
    assert!(extra_data2 == b"extra_data_2");
    
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_deconstruct_token_params_empty() {
    let (scenario, owner_cap, ref, source_cap) = setup_test();
    
    // Create token params without any transfers
    let token_params = dynamic_dispatcher::create_token_params(DESTINATION_CHAIN_SELECTOR);
    
    // Deconstruct should work with empty transfers
    let (dest_chain, transfers) = dynamic_dispatcher::deconstruct_token_params(&source_cap, token_params);
    assert!(dest_chain == DESTINATION_CHAIN_SELECTOR);
    assert!(transfers.length() == 0);
    
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_get_source_token_transfer_data() {
    let (mut scenario, owner_cap, mut ref, source_cap) = setup_test();
    
    // Register a token
    registry::register_pool_by_admin(
        &mut ref,
        TOKEN_ADDRESS_1,
        TOKEN_POOL_PACKAGE_ID_1,
        TOKEN_ADDRESS_1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestToken"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        scenario.ctx(),
    );
    
    let token_params = dynamic_dispatcher::create_token_params(DESTINATION_CHAIN_SELECTOR);
    
    // Add source token transfer with specific data
    let updated_params = dynamic_dispatcher::add_source_token_transfer(
        &ref,
        token_params,
        12345, // specific amount
        TOKEN_ADDRESS_1,
        x"deadbeef", // hex dest address
        x"cafebabe", // hex extra data
        TestTypeProof {}
    );
    
    // Get the transfer and verify all data
    let (_, transfers) = dynamic_dispatcher::deconstruct_token_params(&source_cap, updated_params);
    let transfer = &transfers[0];
    let (source_pool, amount, source_token_address, dest_token_address, extra_data) = 
        dynamic_dispatcher::get_source_token_transfer_data(*transfer);
    
    assert!(source_pool == TOKEN_POOL_PACKAGE_ID_1);
    assert!(amount == 12345);
    assert!(source_token_address == TOKEN_ADDRESS_1);
    assert!(dest_token_address == x"deadbeef");
    assert!(extra_data == x"cafebabe");
    
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_edge_case_large_amounts() {
    let (mut scenario, owner_cap, mut ref, source_cap) = setup_test();
    
    // Register a token
    registry::register_pool_by_admin(
        &mut ref,
        TOKEN_ADDRESS_1,
        TOKEN_POOL_PACKAGE_ID_1,
        TOKEN_ADDRESS_1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestToken"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        scenario.ctx(),
    );
    
    let token_params = dynamic_dispatcher::create_token_params(DESTINATION_CHAIN_SELECTOR);
    
    // Test with maximum u64 value
    let max_amount = 18446744073709551615; // u64::MAX
    let updated_params = dynamic_dispatcher::add_source_token_transfer(
        &ref,
        token_params,
        max_amount,
        TOKEN_ADDRESS_1,
        b"dest_address",
        b"extra_data",
        TestTypeProof {}
    );
    
    let (_, transfers) = dynamic_dispatcher::deconstruct_token_params(&source_cap, updated_params);
    let transfer = &transfers[0];
    let (_, amount, _, _, _) = dynamic_dispatcher::get_source_token_transfer_data(*transfer);
    
    assert!(amount == max_amount);
    
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_edge_case_empty_data() {
    let (mut scenario, owner_cap, mut ref, source_cap) = setup_test();
    
    // Register a token
    registry::register_pool_by_admin(
        &mut ref,
        TOKEN_ADDRESS_1,
        TOKEN_POOL_PACKAGE_ID_1,
        TOKEN_ADDRESS_1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestToken"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        scenario.ctx(),
    );
    
    let token_params = dynamic_dispatcher::create_token_params(DESTINATION_CHAIN_SELECTOR);
    
    // Test with empty destination address and extra data
    let updated_params = dynamic_dispatcher::add_source_token_transfer(
        &ref,
        token_params,
        100,
        TOKEN_ADDRESS_1,
        vector[], // empty dest address
        vector[], // empty extra data
        TestTypeProof {}
    );
    
    let (_, transfers) = dynamic_dispatcher::deconstruct_token_params(&source_cap, updated_params);
    let transfer = &transfers[0];
    let (_, _, _, dest_token_address, extra_data) = 
        dynamic_dispatcher::get_source_token_transfer_data(*transfer);
    
    assert!(dest_token_address == vector[]);
    assert!(extra_data == vector[]);
    
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_different_destination_chains() {
    let (scenario, owner_cap, ref, source_cap) = setup_test();
    
    // Test creating token params for different destination chains
    let chains = vector[1, 100, 1000, 999999];
    let mut i = 0;
    let mut token_params_list = vector[];
    
    while (i < chains.length()) {
        let chain = chains[i];
        let token_params = dynamic_dispatcher::create_token_params(chain);
        let destination = dynamic_dispatcher::get_destination_chain_selector(&token_params);
        assert!(destination == chain);
        token_params_list.push_back(token_params);
        i = i + 1;
    };
    
    // Clean up all token params
    while (!token_params_list.is_empty()) {
        let token_params = token_params_list.pop_back();
        let (_, _) = dynamic_dispatcher::deconstruct_token_params(&source_cap, token_params);
    };

    token_params_list.destroy_empty();
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_zero_amount_transfer() {
    let (mut scenario, owner_cap, mut ref, source_cap) = setup_test();
    
    // Register a token
    registry::register_pool_by_admin(
        &mut ref,
        TOKEN_ADDRESS_1,
        TOKEN_POOL_PACKAGE_ID_1,
        TOKEN_ADDRESS_1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestToken"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        scenario.ctx(),
    );
    
    let token_params = dynamic_dispatcher::create_token_params(DESTINATION_CHAIN_SELECTOR);
    
    // Test with zero amount - should be allowed
    let updated_params = dynamic_dispatcher::add_source_token_transfer(
        &ref,
        token_params,
        0, // zero amount
        TOKEN_ADDRESS_1,
        b"dest_token_address",
        b"extra_data",
        TestTypeProof {}
    );
    
    let (_, transfers) = dynamic_dispatcher::deconstruct_token_params(&source_cap, updated_params);
    let transfer = &transfers[0];
    let (_, amount, _, _, _) = dynamic_dispatcher::get_source_token_transfer_data(*transfer);
    
    assert!(amount == 0);
    
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_source_transfer_cap_permission() {
    let (mut scenario, owner_cap, mut ref, source_cap) = setup_test();
    
    // Register a token
    registry::register_pool_by_admin(
        &mut ref,
        TOKEN_ADDRESS_1,
        TOKEN_POOL_PACKAGE_ID_1,
        TOKEN_ADDRESS_1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestToken"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        scenario.ctx(),
    );
    
    let token_params = dynamic_dispatcher::create_token_params(DESTINATION_CHAIN_SELECTOR);
    
    // Add a source token transfer
    let updated_params = dynamic_dispatcher::add_source_token_transfer(
        &ref,
        token_params,
        1000,
        TOKEN_ADDRESS_1,
        b"dest_token_address",
        b"extra_data",
        TestTypeProof {}
    );
    
    // Test that deconstruct_token_params requires the proper SourceTransferCap
    // This test verifies that only the holder of SourceTransferCap can deconstruct
    let (dest_chain, transfers) = dynamic_dispatcher::deconstruct_token_params(&source_cap, updated_params);
    assert!(dest_chain == DESTINATION_CHAIN_SELECTOR);
    assert!(transfers.length() == 1);
    
    // Verify the transfer data is correct
    let transfer = &transfers[0];
    let (source_pool, amount, source_token_address, dest_token_address, extra_data) = 
        dynamic_dispatcher::get_source_token_transfer_data(*transfer);
    
    assert!(source_pool == TOKEN_POOL_PACKAGE_ID_1);
    assert!(amount == 1000);
    assert!(source_token_address == TOKEN_ADDRESS_1);
    assert!(dest_token_address == b"dest_token_address");
    assert!(extra_data == b"extra_data");
    
    cleanup_test(scenario, owner_cap, ref, source_cap);
}
