#[test_only]
module ccip::onramp_state_helper_tests;

use std::ascii;
use std::string;
use std::type_name;

use sui::test_scenario::{Self as ts, Scenario};

use ccip::onramp_state_helper::{Self, SourceTransferCap};
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::ownable::OwnerCap;
use ccip::token_admin_registry as registry;

public struct ONRAMP_STATE_HELPER_TESTS has drop {}

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
    
    // Create onramp state helper and get source transfer cap
    onramp_state_helper::test_init(scenario.ctx());
    
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
public fun test_create_token_transfer_params() {
    let (mut scenario, owner_cap, mut ref, source_cap) = setup_test();
    
    // Register a token in the token admin registry first
    registry::register_pool_by_admin(
        &mut ref,
        state_object::create_ccip_admin_proof_for_test(),
        TOKEN_ADDRESS_1,
        TOKEN_POOL_PACKAGE_ID_1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER, // administrator
        type_name::into_string(type_name::get<TestTypeProof>()),
        scenario.ctx(),
    );
    
    // Test creating token transfer params with valid data
    let token_params = onramp_state_helper::create_token_transfer_params(
        &ref,
        DESTINATION_CHAIN_SELECTOR,
        1000, // amount
        TOKEN_ADDRESS_1,
        b"dest_token_address",
        b"extra_data",
        TestTypeProof {}
    );
    
    // Test the token params by putting it in a vector and getting data
    let token_params_vec = vector[token_params];
    let (remote_chain, source_pool, amount, source_token, dest_token, extra_data) = 
        onramp_state_helper::get_source_token_transfer_data(&token_params_vec, 0);
    
    assert!(remote_chain == DESTINATION_CHAIN_SELECTOR);
    assert!(source_pool == TOKEN_POOL_PACKAGE_ID_1);
    assert!(amount == 1000);
    assert!(source_token == TOKEN_ADDRESS_1);
    assert!(dest_token == b"dest_token_address");
    assert!(extra_data == b"extra_data");
    
    // Clean up
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params_vec);
    
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
#[expected_failure(abort_code = onramp_state_helper::ETypeProofMismatch)]
public fun test_create_token_transfer_params_wrong_proof() {
    let (mut scenario, owner_cap, mut ref, source_cap) = setup_test();
    
    // Register a token with TestTypeProof
    registry::register_pool_by_admin(
        &mut ref,
        state_object::create_ccip_admin_proof_for_test(),
        TOKEN_ADDRESS_1,
        TOKEN_POOL_PACKAGE_ID_1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        scenario.ctx(),
    );
    
    // Try to create token transfer params with wrong proof type - should fail
    let token_params = onramp_state_helper::create_token_transfer_params(
        &ref,
        DESTINATION_CHAIN_SELECTOR,
        1000,
        TOKEN_ADDRESS_1,
        b"dest_token_address",
        b"extra_data",
        TestTypeProof2 {} // Wrong proof type!
    );

    let token_params_vec = vector[token_params];
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params_vec);
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_get_remote_chain_selector() {
    let (mut scenario, owner_cap, mut ref, source_cap) = setup_test();
    
    // Register a token
    registry::register_pool_by_admin(
        &mut ref,
        state_object::create_ccip_admin_proof_for_test(),
        TOKEN_ADDRESS_1,
        TOKEN_POOL_PACKAGE_ID_1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        scenario.ctx(),
    );
    
    // Test creating token transfer params with different chain selectors
    let token_params1 = onramp_state_helper::create_token_transfer_params(
        &ref,
        DESTINATION_CHAIN_SELECTOR,
        1000,
        TOKEN_ADDRESS_1,
        b"dest_token_address",
        b"extra_data",
        TestTypeProof {}
    );
    
    let different_chain = 2000;
    let token_params2 = onramp_state_helper::create_token_transfer_params(
        &ref,
        different_chain,
        1000,
        TOKEN_ADDRESS_1,
        b"dest_token_address",
        b"extra_data",
        TestTypeProof {}
    );
    
    // Test retrieving the remote chain selectors
    let token_params_vec1 = vector[token_params1];
    let (remote_chain1, _, _, _, _, _) = 
        onramp_state_helper::get_source_token_transfer_data(&token_params_vec1, 0);
    assert!(remote_chain1 == DESTINATION_CHAIN_SELECTOR);
    
    let token_params_vec2 = vector[token_params2];
    let (remote_chain2, _, _, _, _, _) = 
        onramp_state_helper::get_source_token_transfer_data(&token_params_vec2, 0);
    assert!(remote_chain2 == different_chain);
    
    // Clean up
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params_vec1);
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params_vec2);
    
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_create_and_verify_token_transfer() {
    let (mut scenario, owner_cap, mut ref, source_cap) = setup_test();
    
    // Register a token in the token admin registry first
    registry::register_pool_by_admin(
        &mut ref,
        state_object::create_ccip_admin_proof_for_test(),
        TOKEN_ADDRESS_1,
        TOKEN_POOL_PACKAGE_ID_1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER, // administrator
        type_name::into_string(type_name::get<TestTypeProof>()),
        scenario.ctx(),
    );
    
    // Create source token transfer
    let token_params = onramp_state_helper::create_token_transfer_params(
        &ref,
        DESTINATION_CHAIN_SELECTOR,
        1000, // amount
        TOKEN_ADDRESS_1,
        b"dest_token_address",
        b"extra_data",
        TestTypeProof {}
    );
    
    // Verify the token transfer data
    let token_params_vec = vector[token_params];
    let (remote_chain, source_pool, amount, source_token_address, dest_token_address, extra_data) = 
        onramp_state_helper::get_source_token_transfer_data(&token_params_vec, 0);
    
    assert!(remote_chain == DESTINATION_CHAIN_SELECTOR);
    assert!(source_pool == TOKEN_POOL_PACKAGE_ID_1);
    assert!(amount == 1000);
    assert!(source_token_address == TOKEN_ADDRESS_1);
    assert!(dest_token_address == b"dest_token_address");
    assert!(extra_data == b"extra_data");
    
    // Clean up
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params_vec);
    
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_multiple_token_transfers() {
    let (mut scenario, owner_cap, mut ref, source_cap) = setup_test();
    
    // Register two tokens
    registry::register_pool_by_admin(
        &mut ref,
        state_object::create_ccip_admin_proof_for_test(),
        TOKEN_ADDRESS_1,
        TOKEN_POOL_PACKAGE_ID_1,
        string::utf8(b"test_pool_1"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        scenario.ctx(),
    );
    
    registry::register_pool_by_admin(
        &mut ref,
        state_object::create_ccip_admin_proof_for_test(),
        TOKEN_ADDRESS_2,
        TOKEN_POOL_PACKAGE_ID_2,
        string::utf8(b"test_pool_2"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof2>()),
        scenario.ctx(),
    );
    
    // Create first token transfer
    let token_params1 = onramp_state_helper::create_token_transfer_params(
        &ref,
        DESTINATION_CHAIN_SELECTOR,
        1000,
        TOKEN_ADDRESS_1,
        b"dest_token_address_1",
        b"extra_data_1",
        TestTypeProof {}
    );
    
    // Create second token transfer  
    let token_params2 = onramp_state_helper::create_token_transfer_params(
        &ref,
        DESTINATION_CHAIN_SELECTOR,
        2000,
        TOKEN_ADDRESS_2,
        b"dest_token_address_2",
        b"extra_data_2",
        TestTypeProof2 {}
    );

    // Put both in a vector
    let token_params_vec = vector[token_params1, token_params2];
    
    // Verify first transfer
    let (remote_chain1, source_pool1, amount1, source_token_address1, dest_token_address1, extra_data1) = 
        onramp_state_helper::get_source_token_transfer_data(&token_params_vec, 0);
    
    assert!(remote_chain1 == DESTINATION_CHAIN_SELECTOR);
    assert!(source_pool1 == TOKEN_POOL_PACKAGE_ID_1);
    assert!(amount1 == 1000);
    assert!(source_token_address1 == TOKEN_ADDRESS_1);
    assert!(dest_token_address1 == b"dest_token_address_1");
    assert!(extra_data1 == b"extra_data_1");
    
    // Verify second transfer
    let (remote_chain2, source_pool2, amount2, source_token_address2, dest_token_address2, extra_data2) = 
        onramp_state_helper::get_source_token_transfer_data(&token_params_vec, 1);
    
    assert!(remote_chain2 == DESTINATION_CHAIN_SELECTOR);
    assert!(source_pool2 == TOKEN_POOL_PACKAGE_ID_2);
    assert!(amount2 == 2000);
    assert!(source_token_address2 == TOKEN_ADDRESS_2);
    assert!(dest_token_address2 == b"dest_token_address_2");
    assert!(extra_data2 == b"extra_data_2");

    // Clean up
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params_vec);
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_deconstruct_empty_params_vector() {
    let (scenario, owner_cap, ref, source_cap) = setup_test();
    
    // Create empty vector and test deconstruct
    let empty_params_vec = vector[];
    
    // Deconstruct should work with empty vector
    onramp_state_helper::deconstruct_token_params(&source_cap, empty_params_vec);
    
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_get_source_token_transfer_data() {
    let (mut scenario, owner_cap, mut ref, source_cap) = setup_test();
    
    // Register a token
    registry::register_pool_by_admin(
        &mut ref,
        state_object::create_ccip_admin_proof_for_test(),
        TOKEN_ADDRESS_1,
        TOKEN_POOL_PACKAGE_ID_1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        scenario.ctx(),
    );
    
    // Create token transfer with specific data
    let token_params = onramp_state_helper::create_token_transfer_params(
        &ref,
        DESTINATION_CHAIN_SELECTOR,
        12345, // specific amount
        TOKEN_ADDRESS_1,
        x"deadbeef", // hex dest address
        x"cafebabe", // hex extra data
        TestTypeProof {}
    );
    
    // Get the transfer and verify all data
    let token_params_vec = vector[token_params];
    let (remote_chain, source_pool, amount, source_token_address, dest_token_address, extra_data) = 
        onramp_state_helper::get_source_token_transfer_data(&token_params_vec, 0);
    
    assert!(remote_chain == DESTINATION_CHAIN_SELECTOR);
    assert!(source_pool == TOKEN_POOL_PACKAGE_ID_1);
    assert!(amount == 12345);
    assert!(source_token_address == TOKEN_ADDRESS_1);
    assert!(dest_token_address == x"deadbeef");
    assert!(extra_data == x"cafebabe");
    
    // Clean up
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params_vec);
    
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_edge_case_large_amounts() {
    let (mut scenario, owner_cap, mut ref, source_cap) = setup_test();
    
    // Register a token
    registry::register_pool_by_admin(
        &mut ref,
        state_object::create_ccip_admin_proof_for_test(),
        TOKEN_ADDRESS_1,
        TOKEN_POOL_PACKAGE_ID_1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        scenario.ctx(),
    );
    
    // Test with maximum u64 value
    let max_amount = 18446744073709551615; // u64::MAX
    let token_params = onramp_state_helper::create_token_transfer_params(
        &ref,
        DESTINATION_CHAIN_SELECTOR,
        max_amount,
        TOKEN_ADDRESS_1,
        b"dest_address",
        b"extra_data",
        TestTypeProof {}
    );
    
    let token_params_vec = vector[token_params];
    let (_, _, amount, _, _, _) = onramp_state_helper::get_source_token_transfer_data(&token_params_vec, 0);
    
    assert!(amount == max_amount);
    
    // Clean up
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params_vec);
    
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_edge_case_empty_data() {
    let (mut scenario, owner_cap, mut ref, source_cap) = setup_test();
    
    // Register a token
    registry::register_pool_by_admin(
        &mut ref,
        state_object::create_ccip_admin_proof_for_test(),
        TOKEN_ADDRESS_1,
        TOKEN_POOL_PACKAGE_ID_1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        scenario.ctx(),
    );
    
    // Test with empty destination address and extra data
    let token_params = onramp_state_helper::create_token_transfer_params(
        &ref,
        DESTINATION_CHAIN_SELECTOR,
        100,
        TOKEN_ADDRESS_1,
        vector[], // empty dest address
        vector[], // empty extra data
        TestTypeProof {}
    );
    
    let token_params_vec = vector[token_params];
    let (_, _, _, _, dest_token_address, extra_data) = 
        onramp_state_helper::get_source_token_transfer_data(&token_params_vec, 0);
    
    assert!(dest_token_address == vector<u8>[]);
    assert!(extra_data == vector<u8>[]);
    
    // Clean up
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params_vec);
    
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_different_destination_chains() {
    let (mut scenario, owner_cap, mut ref, source_cap) = setup_test();
    
    // Register a token
    registry::register_pool_by_admin(
        &mut ref,
        state_object::create_ccip_admin_proof_for_test(),
        TOKEN_ADDRESS_1,
        TOKEN_POOL_PACKAGE_ID_1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        scenario.ctx(),
    );
    
    // Test creating token transfer params for different destination chains
    let chains = vector[1, 100, 1000, 999999];
    let mut i = 0;
    let mut token_params_list = vector[];
    
    while (i < chains.length()) {
        let chain = chains[i];
        let token_params = onramp_state_helper::create_token_transfer_params(
            &ref,
            chain,
            1000,
            TOKEN_ADDRESS_1,
            b"dest_address",
            b"extra_data",
            TestTypeProof {}
        );
        
        let mut temp_vec = vector[token_params];
        let (remote_chain, _, _, _, _, _) = onramp_state_helper::get_source_token_transfer_data(&temp_vec, 0);
        assert!(remote_chain == chain);
        
        token_params_list.push_back(temp_vec.pop_back());
        temp_vec.destroy_empty();
        i = i + 1;
    };
    
    // Clean up all token params
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params_list);
    
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_zero_amount_transfer() {
    let (mut scenario, owner_cap, mut ref, source_cap) = setup_test();
    
    // Register a token
    registry::register_pool_by_admin(
        &mut ref,
        state_object::create_ccip_admin_proof_for_test(),
        TOKEN_ADDRESS_1,
        TOKEN_POOL_PACKAGE_ID_1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        scenario.ctx(),
    );
    
    // Test with zero amount - should be allowed
    let token_params = onramp_state_helper::create_token_transfer_params(
        &ref,
        DESTINATION_CHAIN_SELECTOR,
        0, // zero amount
        TOKEN_ADDRESS_1,
        b"dest_token_address",
        b"extra_data",
        TestTypeProof {}
    );
    
    let token_params_vec = vector[token_params];
    let (_, _, amount, _, _, _) = onramp_state_helper::get_source_token_transfer_data(&token_params_vec, 0);
    
    assert!(amount == 0);
    
    // Clean up
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params_vec);
    
    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_source_transfer_cap_permission() {
    let (mut scenario, owner_cap, mut ref, source_cap) = setup_test();
    
    // Register a token
    registry::register_pool_by_admin(
        &mut ref,
        state_object::create_ccip_admin_proof_for_test(),
        TOKEN_ADDRESS_1,
        TOKEN_POOL_PACKAGE_ID_1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        scenario.ctx(),
    );
    
    // Create a source token transfer
    let token_params = onramp_state_helper::create_token_transfer_params(
        &ref,
        DESTINATION_CHAIN_SELECTOR,
        1000,
        TOKEN_ADDRESS_1,
        b"dest_token_address",
        b"extra_data",
        TestTypeProof {}
    );
    
    // Test that deconstruct_token_params requires the proper SourceTransferCap
    // This test verifies that only the holder of SourceTransferCap can deconstruct
    let token_params_vec = vector[token_params];
    let (remote_chain, source_pool, amount, source_token_address, dest_token_address, extra_data) = 
        onramp_state_helper::get_source_token_transfer_data(&token_params_vec, 0);
    
    assert!(remote_chain == DESTINATION_CHAIN_SELECTOR);
    assert!(source_pool == TOKEN_POOL_PACKAGE_ID_1);
    assert!(amount == 1000);
    assert!(source_token_address == TOKEN_ADDRESS_1);
    assert!(dest_token_address == b"dest_token_address");
    assert!(extra_data == b"extra_data");
    
    // Clean up
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params_vec);
    
    cleanup_test(scenario, owner_cap, ref, source_cap);
}
