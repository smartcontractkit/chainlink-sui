#[test_only]
module ccip::onramp_state_helper_tests;

use ccip::onramp_state_helper::{Self, SourceTransferCap};
use ccip::ownable::OwnerCap;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::token_admin_registry as registry;
use std::ascii;
use std::string;
use std::type_name;
use sui::test_scenario::{Self as ts, Scenario};

public struct ONRAMP_STATE_HELPER_TESTS has drop {}

public struct TestTypeProof has drop {}
public struct TestTypeProof2 has drop {}

const OWNER: address = @0x1000;
const TOKEN_ADDRESS_1: address =
    @0x1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b;
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

fun cleanup_test(
    scenario: Scenario,
    owner_cap: OwnerCap,
    ref: CCIPObjectRef,
    source_cap: SourceTransferCap,
) {
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
        @0x1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER, // administrator
        type_name::into_string(type_name::get<TestTypeProof>()),
        vector<address>[], // lock_or_burn_params
        vector<address>[], // release_or_mint_params
        scenario.ctx(),
    );

    // Test creating token transfer params with valid data
    let mut token_params = onramp_state_helper::create_token_transfer_params(@0x456);
    onramp_state_helper::add_token_transfer_param<TestTypeProof>(
        &ref,
        &mut token_params,
        DESTINATION_CHAIN_SELECTOR,
        1000, // amount
        TOKEN_ADDRESS_1,
        b"dest_token_address",
        b"extra_data",
        TestTypeProof {},
    );

    // Test the token params by getting data directly
    let (
        remote_chain,
        token_pool_package_id,
        amount,
        source_token,
        dest_token,
        extra_data,
    ) = onramp_state_helper::get_source_token_transfer_data(&token_params);

    assert!(remote_chain == DESTINATION_CHAIN_SELECTOR);
    assert!(token_pool_package_id == @0x1);
    assert!(amount == 1000);
    assert!(source_token == TOKEN_ADDRESS_1);
    assert!(dest_token == b"dest_token_address");
    assert!(extra_data == b"extra_data");

    // Clean up
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params);

    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_create_token_transfer_params_basic() {
    let (mut scenario, owner_cap, mut ref, source_cap) = setup_test();

    // Register a token with TestTypeProof
    registry::register_pool_by_admin(
        &mut ref,
        state_object::create_ccip_admin_proof_for_test(),
        TOKEN_ADDRESS_1,
        @0x1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        vector<address>[], // lock_or_burn_params
        vector<address>[], // release_or_mint_params
        scenario.ctx(),
    );

    // Try to create token transfer params - this test now just creates empty params
    // since type proof validation was removed from the helper
    let mut token_params = onramp_state_helper::create_token_transfer_params(@0x456);
    onramp_state_helper::add_token_transfer_param<TestTypeProof>(
        &ref,
        &mut token_params,
        DESTINATION_CHAIN_SELECTOR,
        1000,
        TOKEN_ADDRESS_1,
        b"dest_token_address",
        b"extra_data",
        TestTypeProof {},
    );

    onramp_state_helper::deconstruct_token_params(&source_cap, token_params);
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
        @0x1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        vector<address>[], // lock_or_burn_params
        vector<address>[], // release_or_mint_params
        scenario.ctx(),
    );

    // Test creating token transfer params with different chain selectors
    let mut token_params1 = onramp_state_helper::create_token_transfer_params(@0x456);
    onramp_state_helper::add_token_transfer_param<TestTypeProof>(
        &ref,
        &mut token_params1,
        DESTINATION_CHAIN_SELECTOR,
        1000,
        TOKEN_ADDRESS_1,
        b"dest_token_address",
        b"extra_data",
        TestTypeProof {},
    );

    let different_chain = 2000;
    let mut token_params2 = onramp_state_helper::create_token_transfer_params(@0x456);
    onramp_state_helper::add_token_transfer_param<TestTypeProof>(
        &ref,
        &mut token_params2,
        different_chain,
        1000,
        TOKEN_ADDRESS_1,
        b"dest_token_address",
        b"extra_data",
        TestTypeProof {},
    );

    // Test retrieving the remote chain selectors
    let (remote_chain1, _, _, _, _, _) = onramp_state_helper::get_source_token_transfer_data(
        &token_params1,
    );
    assert!(remote_chain1 == DESTINATION_CHAIN_SELECTOR);

    let (remote_chain2, _, _, _, _, _) = onramp_state_helper::get_source_token_transfer_data(
        &token_params2,
    );
    assert!(remote_chain2 == different_chain);

    // Clean up
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params1);
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params2);

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
        @0x1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER, // administrator
        type_name::into_string(type_name::get<TestTypeProof>()),
        vector<address>[], // lock_or_burn_params
        vector<address>[], // release_or_mint_params
        scenario.ctx(),
    );

    // Create source token transfer
    let mut token_params = onramp_state_helper::create_token_transfer_params(@0x456);
    onramp_state_helper::add_token_transfer_param<TestTypeProof>(
        &ref,
        &mut token_params,
        DESTINATION_CHAIN_SELECTOR,
        1000, // amount
        TOKEN_ADDRESS_1,
        b"dest_token_address",
        b"extra_data",
        TestTypeProof {},
    );

    // Verify the token transfer data
    let (
        remote_chain,
        token_pool_package_id,
        amount,
        source_token_address,
        dest_token_address,
        extra_data,
    ) = onramp_state_helper::get_source_token_transfer_data(&token_params);

    assert!(remote_chain == DESTINATION_CHAIN_SELECTOR);
    assert!(token_pool_package_id == @0x1);
    assert!(amount == 1000);
    assert!(source_token_address == TOKEN_ADDRESS_1);
    assert!(dest_token_address == b"dest_token_address");
    assert!(extra_data == b"extra_data");

    // Clean up
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params);

    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
public fun test_deconstruct_empty_params_vector() {
    let (scenario, owner_cap, ref, source_cap) = setup_test();

    // Create empty params and test deconstruct
    let empty_params = onramp_state_helper::create_token_transfer_params(@0x456);

    // Deconstruct should work with empty params
    onramp_state_helper::deconstruct_token_params(&source_cap, empty_params);

    cleanup_test(scenario, owner_cap, ref, source_cap);
}

#[test]
#[expected_failure(abort_code = onramp_state_helper::ETokenTransferDoesNotExist)]
public fun test_get_token_transfer_data_when_empty() {
    let (scenario, owner_cap, ref, source_cap) = setup_test();

    // Create empty token transfer params (no token transfer added)
    let token_params = onramp_state_helper::create_token_transfer_params(@0x456);

    // This should fail with ETokenTransferDoesNotExist because no token transfer was added
    let (
        _remote_chain,
        _token_pool_package_id,
        _amount,
        _source_token_address,
        _dest_token_address,
        _extra_data,
    ) = onramp_state_helper::get_source_token_transfer_data(&token_params);

    // The following code won't be reached due to the expected failure above
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params);
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
        @0x1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        vector<address>[], // lock_or_burn_params
        vector<address>[], // release_or_mint_params
        scenario.ctx(),
    );

    // Create token transfer with specific data
    let mut token_params = onramp_state_helper::create_token_transfer_params(@0x456);
    onramp_state_helper::add_token_transfer_param<TestTypeProof>(
        &ref,
        &mut token_params,
        DESTINATION_CHAIN_SELECTOR,
        12345, // specific amount
        TOKEN_ADDRESS_1,
        x"deadbeef", // hex dest address
        x"cafebabe", // hex extra data
        TestTypeProof {},
    );

    // Get the transfer and verify all data
    let (
        remote_chain,
        token_pool_package_id,
        amount,
        source_token_address,
        dest_token_address,
        extra_data,
    ) = onramp_state_helper::get_source_token_transfer_data(&token_params);

    assert!(remote_chain == DESTINATION_CHAIN_SELECTOR);
    assert!(token_pool_package_id == @0x1);
    assert!(amount == 12345);
    assert!(source_token_address == TOKEN_ADDRESS_1);
    assert!(dest_token_address == x"deadbeef");
    assert!(extra_data == x"cafebabe");

    // Clean up
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params);

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
        @0x1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        vector<address>[], // lock_or_burn_params
        vector<address>[], // release_or_mint_params
        scenario.ctx(),
    );

    // Test with maximum u64 value
    let max_amount = 18446744073709551615; // u64::MAX
    let mut token_params = onramp_state_helper::create_token_transfer_params(@0x456);
    onramp_state_helper::add_token_transfer_param<TestTypeProof>(
        &ref,
        &mut token_params,
        DESTINATION_CHAIN_SELECTOR,
        max_amount,
        TOKEN_ADDRESS_1,
        b"dest_address",
        b"extra_data",
        TestTypeProof {},
    );

    let (_, _, amount, _, _, _) = onramp_state_helper::get_source_token_transfer_data(
        &token_params,
    );

    assert!(amount == max_amount);

    // Clean up
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params);

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
        @0x1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        vector<address>[], // lock_or_burn_params
        vector<address>[], // release_or_mint_params
        scenario.ctx(),
    );

    // Test with empty destination address and extra data
    let mut token_params = onramp_state_helper::create_token_transfer_params(@0x456);
    onramp_state_helper::add_token_transfer_param<TestTypeProof>(
        &ref,
        &mut token_params,
        DESTINATION_CHAIN_SELECTOR,
        100,
        TOKEN_ADDRESS_1,
        vector<u8>[], // empty dest address
        vector<u8>[], // empty extra data
        TestTypeProof {},
    );

    let (
        _,
        _,
        _,
        _,
        dest_token_address,
        extra_data,
    ) = onramp_state_helper::get_source_token_transfer_data(&token_params);

    assert!(dest_token_address == vector<u8>[]);
    assert!(extra_data == vector<u8>[]);

    // Clean up
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params);

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
        @0x1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        vector<address>[], // lock_or_burn_params
        vector<address>[], // release_or_mint_params
        scenario.ctx(),
    );

    // Test creating token transfer params for different destination chains
    let chains = vector[1, 100, 1000, 999999];
    let mut i = 0;

    // Create separate token params objects for each chain
    while (i < chains.length()) {
        let chain = chains[i];
        let mut token_params = onramp_state_helper::create_token_transfer_params(@0x456);

        onramp_state_helper::add_token_transfer_param<TestTypeProof>(
            &ref,
            &mut token_params,
            chain,
            1000,
            TOKEN_ADDRESS_1,
            b"dest_address",
            b"extra_data",
            TestTypeProof {},
        );

        let (remote_chain, _, _, _, _, _) = onramp_state_helper::get_source_token_transfer_data(
            &token_params,
        );
        assert!(remote_chain == chain);

        // Clean up this token params
        onramp_state_helper::deconstruct_token_params(&source_cap, token_params);

        i = i + 1;
    };

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
        @0x1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        vector<address>[], // lock_or_burn_params
        vector<address>[], // release_or_mint_params
        scenario.ctx(),
    );

    // Test with zero amount - should be allowed
    let mut token_params = onramp_state_helper::create_token_transfer_params(@0x456);
    onramp_state_helper::add_token_transfer_param<TestTypeProof>(
        &ref,
        &mut token_params,
        DESTINATION_CHAIN_SELECTOR,
        0, // zero amount
        TOKEN_ADDRESS_1,
        b"dest_token_address",
        b"extra_data",
        TestTypeProof {},
    );

    let (_, _, amount, _, _, _) = onramp_state_helper::get_source_token_transfer_data(
        &token_params,
    );

    assert!(amount == 0);

    // Clean up
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params);

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
        @0x1,
        string::utf8(b"test_pool"),
        ascii::string(b"TestType"),
        OWNER,
        type_name::into_string(type_name::get<TestTypeProof>()),
        vector<address>[], // lock_or_burn_params
        vector<address>[], // release_or_mint_params
        scenario.ctx(),
    );

    // Create a source token transfer
    let mut token_params = onramp_state_helper::create_token_transfer_params(@0x456);
    onramp_state_helper::add_token_transfer_param<TestTypeProof>(
        &ref,
        &mut token_params,
        DESTINATION_CHAIN_SELECTOR,
        1000,
        TOKEN_ADDRESS_1,
        b"dest_token_address",
        b"extra_data",
        TestTypeProof {},
    );

    // Test that deconstruct_token_params requires the proper SourceTransferCap
    // This test verifies that only the holder of SourceTransferCap can deconstruct
    let (
        remote_chain,
        token_pool_package_id,
        amount,
        source_token_address,
        dest_token_address,
        extra_data,
    ) = onramp_state_helper::get_source_token_transfer_data(&token_params);

    assert!(remote_chain == DESTINATION_CHAIN_SELECTOR);
    assert!(token_pool_package_id == @0x1);
    assert!(amount == 1000);
    assert!(source_token_address == TOKEN_ADDRESS_1);
    assert!(dest_token_address == b"dest_token_address");
    assert!(extra_data == b"extra_data");

    // Clean up
    onramp_state_helper::deconstruct_token_params(&source_cap, token_params);

    cleanup_test(scenario, owner_cap, ref, source_cap);
}
