#[test_only]
module ccip::token_admin_registry_tests;

use ccip::ownable::OwnerCap;
use ccip::publisher_wrapper;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::token_admin_registry as registry;
use ccip::upgrade_registry;
use mcms::mcms_account;
use mcms::mcms_deployer;
use mcms::mcms_registry::{Self, Registry};
use std::ascii;
use std::bcs;
use std::string;
use std::type_name;
use sui::address;
use sui::coin;
use sui::package;
use sui::test_scenario::{Self as ts, Scenario};

// === Test Witness Types ===

public struct TOKEN_ADMIN_REGISTRY_TESTS has drop {}
public struct TypeProof has drop {}
public struct TypeProof2 has drop {}

// === Constants ===

const DECIMALS: u8 = 8;

// Test addresses
const CCIP_ADMIN: address = @0x1000;
const TOKEN_ADMIN_ADDRESS: address = @0x1;
const TOKEN_ADMIN_ADDRESS_2: address = @0x2;
const RANDOM_USER: address = @0x3;

// === Helper Functions ===

fun create_test_scenario(addr: address): Scenario {
    ts::begin(addr)
}

fun initialize_state_and_registry(scenario: &mut Scenario, admin: address) {
    scenario.next_tx(admin);
    {
        let ctx = scenario.ctx();
        mcms_account::test_init(ctx);
        mcms_registry::test_init(ctx);
        mcms_deployer::test_init(ctx);
        state_object::test_init(ctx);
    };

    scenario.next_tx(admin);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();
        let ctx = scenario.ctx();

        // Initialize upgrade registry first (required by token_admin_registry functions)
        upgrade_registry::initialize(&mut ref, &owner_cap, ctx);
        registry::initialize(&mut ref, &owner_cap, ctx);

        scenario.return_to_sender(owner_cap);
        ts::return_shared(ref);
    };
}

fun create_test_token(
    scenario: &mut Scenario,
): (coin::TreasuryCap<TOKEN_ADMIN_REGISTRY_TESTS>, coin::CoinMetadata<TOKEN_ADMIN_REGISTRY_TESTS>) {
    coin::create_currency(
        TOKEN_ADMIN_REGISTRY_TESTS {},
        DECIMALS,
        b"TEST",
        b"TestToken",
        b"test_token",
        option::none(),
        scenario.ctx(),
    )
}

fun register_test_pool<T>(
    ref: &mut CCIPObjectRef,
    treasury_cap: &coin::TreasuryCap<T>,
    coin_metadata: &coin::CoinMetadata<T>,
    admin: address,
    ctx: &mut TxContext,
) {
    let publisher = package::test_claim(TOKEN_ADMIN_REGISTRY_TESTS {}, ctx);
    let publisher_wrapper = publisher_wrapper::create(&publisher, TypeProof {});

    registry::register_pool(
        ref,
        treasury_cap,
        coin_metadata,
        admin,
        vector<address>[], // lock_or_burn_params
        vector<address>[], // release_or_mint_params
        publisher_wrapper,
        TypeProof {},
    );

    package::burn_publisher(publisher);
}

fun assert_empty_token_config(ref: &CCIPObjectRef, token_address: address) {
    let (
        token_pool_package_id,
        token_pool_module,
        token_type,
        administrator,
        pending_administrator,
        proof,
        _lock_or_burn_params,
        _release_or_mint_params,
    ) = registry::get_token_config_data(ref, token_address);

    assert!(token_pool_package_id == @0x0);
    assert!(token_pool_module == string::utf8(b""));
    assert!(token_type == ascii::string(b""));
    assert!(administrator == @0x0);
    assert!(pending_administrator == @0x0);
    assert!(proof == ascii::string(b""));
}

fun assert_token_config(
    ref: &CCIPObjectRef,
    token_address: address,
    expected_package_id: address,
    expected_module: vector<u8>,
    expected_type: ascii::String,
    expected_admin: address,
    expected_pending_admin: address,
) {
    let (
        token_pool_package_id,
        token_pool_module,
        token_type,
        administrator,
        pending_administrator,
        _proof,
        _lock_or_burn_params,
        _release_or_mint_params,
    ) = registry::get_token_config_data(ref, token_address);

    assert!(token_pool_package_id == expected_package_id);
    assert!(token_pool_module == string::utf8(expected_module));
    assert!(token_type == expected_type);
    assert!(administrator == expected_admin);
    assert!(pending_administrator == expected_pending_admin);
}

// === Basic Initialization Tests ===

#[test]
public fun test_initialize() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let ref = scenario.take_shared<CCIPObjectRef>();

        // Verify empty configuration
        assert_empty_token_config(&ref, @0x2);

        ts::return_shared(ref);
    };

    ts::end(scenario);
}

#[test]
public fun test_type_and_version() {
    let version = registry::type_and_version();
    assert!(version == string::utf8(b"TokenAdminRegistry 1.6.0"));
}

#[test]
public fun test_get_pool() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    let (treasury_cap, coin_metadata) = create_test_token(&mut scenario);
    let local_token = object::id_to_address(&object::id(&coin_metadata));

    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        // Test with unregistered token
        let pool_address = registry::get_pool(&ref, local_token);
        assert!(pool_address == @0x0);

        // Register token
        register_test_pool(
            &mut ref,
            &treasury_cap,
            &coin_metadata,
            TOKEN_ADMIN_ADDRESS,
            scenario.ctx(),
        );

        // Test with registered token
        let pool_address = registry::get_pool(&ref, local_token);
        let tn = type_name::with_defining_ids<TypeProof>();
        let expected_package_id = address::from_ascii_bytes(&tn.address_string().into_bytes());
        assert!(pool_address == expected_package_id);

        let ctx = scenario.ctx();
        transfer::public_transfer(treasury_cap, ctx.sender());
        ts::return_shared(ref);
    };

    transfer::public_freeze_object(coin_metadata);
    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = registry::ETokenPoolPackageIdAlreadyRegistered)]
public fun test_register_pool_duplicate_package_id_fails() {
    let mut scenario = create_test_scenario(CCIP_ADMIN);
    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    // First registration with a specific package ID
    scenario.next_tx(CCIP_ADMIN);
    let owner_cap = scenario.take_from_sender<OwnerCap>();
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let ctx = scenario.ctx();

        registry::register_pool_as_owner(
            &owner_cap,
            &mut ref,
            @0xABC1, // coin_metadata_address #1
            @0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA,
            string::utf8(b"dup_pool"),
            ascii::string(b"TypeOne"),
            TOKEN_ADMIN_ADDRESS,
            ascii::string(b"ProofOne"),
            vector[@0x6, @0x1111],
            vector[@0x6, @0x2222],
            ctx,
        );

        // Sanity: mapping from package -> coin metadata should be set
        let mapped = registry::get_pool_local_token(
            &ref,
            @0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA,
        );
        assert!(mapped == @0xABC1);

        ts::return_shared(ref);
    };

    // Second registration with the SAME package ID should fail
    scenario.next_tx(CCIP_ADMIN);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let ctx = scenario.ctx();

        registry::register_pool_as_owner(
            &owner_cap,
            &mut ref,
            @0xABC2, // coin_metadata_address #2
            @0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA, // duplicate package id
            string::utf8(b"dup_pool_2"),
            ascii::string(b"TypeTwo"),
            TOKEN_ADMIN_ADDRESS,
            ascii::string(b"ProofTwo"),
            vector[@0x6, @0x3333],
            vector[@0x6, @0x4444],
            ctx,
        );

        ts::return_shared(ref);
    };

    scenario.return_to_sender(owner_cap);
    ts::end(scenario);
}

// === Registration and Pool Management Tests ===

#[test]
public fun test_register_and_unregister() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    let (treasury_cap, coin_metadata) = create_test_token(&mut scenario);
    let local_token = object::id_to_address(&object::id(&coin_metadata));

    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        register_test_pool(
            &mut ref,
            &treasury_cap,
            &coin_metadata,
            TOKEN_ADMIN_ADDRESS_2,
            scenario.ctx(),
        );

        // Verify registration
        let pool_addresses = registry::get_pools(&ref, vector[local_token]);
        assert!(pool_addresses.length() == 1);
        let tn = type_name::with_defining_ids<TypeProof>();
        let expected_package_id = address::from_ascii_bytes(&tn.address_string().into_bytes());
        assert!(pool_addresses[0] == expected_package_id);
        assert!(registry::is_administrator(&ref, local_token, TOKEN_ADMIN_ADDRESS_2));

        let ctx = scenario.ctx();
        transfer::public_transfer(treasury_cap, ctx.sender());
        ts::return_shared(ref);
    };

    // Unregister the token as the token admin
    scenario.next_tx(TOKEN_ADMIN_ADDRESS_2);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        registry::unregister_pool(&mut ref, local_token, scenario.ctx());
        assert_empty_token_config(&ref, local_token);

        ts::return_shared(ref);
    };

    transfer::public_freeze_object(coin_metadata);
    ts::end(scenario);
}

#[test]
public fun test_register() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    let (treasury_cap, coin_metadata) = create_test_token(&mut scenario);
    let local_token = object::id_to_address(&object::id(&coin_metadata));

    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        register_test_pool(
            &mut ref,
            &treasury_cap,
            &coin_metadata,
            TOKEN_ADMIN_ADDRESS,
            scenario.ctx(),
        );

        // Verify initial registration
        let pool_addresses = registry::get_pools(&ref, vector[local_token]);
        assert!(pool_addresses.length() == 1);
        let tn = type_name::with_defining_ids<TypeProof>();
        let expected_package_id = address::from_ascii_bytes(&tn.address_string().into_bytes());
        let expected_module = tn.module_string().into_bytes().to_string();
        assert!(pool_addresses[0] == expected_package_id);
        assert!(registry::is_administrator(&ref, local_token, TOKEN_ADMIN_ADDRESS));

        // Verify detailed configuration
        assert_token_config(
            &ref,
            local_token,
            expected_package_id,
            expected_module.into_bytes(),
            type_name::with_defining_ids<TOKEN_ADMIN_REGISTRY_TESTS>().into_string(),
            TOKEN_ADMIN_ADDRESS,
            @0x0,
        );

        let (_, _, token_type, _, _, type_proof, _, _) = registry::get_token_config_data(
            &ref,
            local_token,
        );
        assert!(
            token_type == ascii::string(b"5ef4b483da6644c84aa78eae4f51a9bfb1fb4554d5134ac98892e931fcbdd6bf::token_admin_registry_tests::TOKEN_ADMIN_REGISTRY_TESTS"),
        );
        assert!(type_proof == type_name::into_string(type_name::with_defining_ids<TypeProof>()));

        let ctx = scenario.ctx();

        // Request admin transfer
        registry::transfer_admin_role(&mut ref, local_token, TOKEN_ADMIN_ADDRESS_2, ctx);

        transfer::public_transfer(treasury_cap, ctx.sender());
        ts::return_shared(ref);
    };

    scenario.next_tx(TOKEN_ADMIN_ADDRESS_2);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        // Since TypeProof and TypeProof2 have the same package ID, the configuration should remain unchanged
        let tn = type_name::with_defining_ids<TypeProof>();
        let expected_package_id = address::from_ascii_bytes(&tn.address_string().into_bytes());
        let expected_module = tn.module_string().into_bytes().to_string();

        // Verify configuration remains unchanged (same package ID means no update)
        assert_token_config(
            &ref,
            local_token,
            expected_package_id,
            expected_module.into_bytes(),
            type_name::with_defining_ids<TOKEN_ADMIN_REGISTRY_TESTS>().into_string(),
            TOKEN_ADMIN_ADDRESS,
            TOKEN_ADMIN_ADDRESS_2,
        );

        let (_, _, token_type, _, _, type_proof, _, _) = registry::get_token_config_data(
            &ref,
            local_token,
        );
        assert!(
            token_type == ascii::string(b"5ef4b483da6644c84aa78eae4f51a9bfb1fb4554d5134ac98892e931fcbdd6bf::token_admin_registry_tests::TOKEN_ADMIN_REGISTRY_TESTS"),
        );
        // Since TypeProof and TypeProof2 have the same package ID, the type proof should remain as TypeProof
        assert!(type_proof == type_name::into_string(type_name::with_defining_ids<TypeProof>()));

        // Accept admin role
        registry::accept_admin_role(&mut ref, local_token, scenario.ctx());
        assert!(registry::is_administrator(&ref, local_token, TOKEN_ADMIN_ADDRESS_2));

        ts::return_shared(ref);
    };

    transfer::public_freeze_object(coin_metadata);
    ts::end(scenario);
}

// === Token Pagination Tests ===

#[test]
public fun test_get_all_configured_tokens() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        registry::insert_token_configs_for_test(
            &mut ref,
            TOKEN_ADMIN_ADDRESS,
            vector[@0x1, @0x2, @0x3],
            TypeProof {},
        );

        // Test with max_count = 0
        let (res, next_key, has_more) = registry::get_all_configured_tokens(&ref, @0x0, 0);
        assert!(res.length() == 0);
        assert!(next_key == @0x0);
        assert!(has_more);

        // Test getting all tokens
        let (res, next_key, has_more) = registry::get_all_configured_tokens(&ref, @0x0, 3);
        assert!(res.length() == 3);
        assert!(vector[@0x1, @0x2, @0x3] == res);
        assert!(next_key == @0x3);
        assert!(!has_more);

        ts::return_shared(ref);
    };

    ts::end(scenario);
}

#[test]
public fun test_get_all_configured_tokens_edge_cases() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        // Test case 1: Empty state
        let (res, next_key, has_more) = registry::get_all_configured_tokens(&ref, @0x0, 1);
        assert!(res.length() == 0);
        assert!(next_key == @0x0);
        assert!(!has_more);

        // Test case 2: Single token
        registry::insert_token_configs_for_test(
            &mut ref,
            TOKEN_ADMIN_ADDRESS,
            vector[@0x1],
            TypeProof {},
        );
        let (res, _next_key, has_more) = registry::get_all_configured_tokens(&ref, @0x0, 1);
        assert!(res.length() == 1);
        assert!(res[0] == @0x1);
        assert!(!has_more);

        // Test case 3: Start from middle
        registry::insert_token_configs_for_test(
            &mut ref,
            TOKEN_ADMIN_ADDRESS,
            vector[@0x2, @0x3],
            TypeProof {},
        );
        let (res, _next_key, has_more) = registry::get_all_configured_tokens(&ref, @0x1, 2);
        assert!(res.length() == 2);
        assert!(res[0] == @0x2);
        assert!(res[1] == @0x3);
        assert!(!has_more);

        // Test case 4: Request more than available
        let (res, _next_key, has_more) = registry::get_all_configured_tokens(&ref, @0x0, 5);
        assert!(res.length() == 3);
        assert!(res[0] == @0x1);
        assert!(res[1] == @0x2);
        assert!(res[2] == @0x3);
        assert!(!has_more);

        ts::return_shared(ref);
    };

    ts::end(scenario);
}

#[test]
public fun test_get_all_configured_tokens_pagination() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        registry::insert_token_configs_for_test(
            &mut ref,
            TOKEN_ADMIN_ADDRESS,
            vector[@0x1, @0x2, @0x3, @0x4, @0x5],
            TypeProof {},
        );

        // Test pagination with different chunk sizes
        let mut current_key = @0x0;
        let mut total_tokens = vector[];

        // First page: get 2 tokens
        let (res, next_key, more) = registry::get_all_configured_tokens(&ref, current_key, 2);
        assert!(res.length() == 2);
        assert!(res[0] == @0x1);
        assert!(res[1] == @0x2);
        assert!(more);
        current_key = next_key;
        total_tokens.append(res);

        // Second page: get 2 more tokens
        let (res, next_key, more) = registry::get_all_configured_tokens(&ref, current_key, 2);
        assert!(res.length() == 2);
        assert!(res[0] == @0x3);
        assert!(res[1] == @0x4);
        assert!(more);
        current_key = next_key;
        total_tokens.append(res);

        // Last page: get remaining token
        let (res, _next_key, more) = registry::get_all_configured_tokens(&ref, current_key, 2);
        assert!(res.length() == 1);
        assert!(res[0] == @0x5);
        assert!(!more);
        total_tokens.append(res);

        // Verify we got all tokens in order
        assert!(total_tokens.length() == 5);
        assert!(total_tokens[0] == @0x1);
        assert!(total_tokens[1] == @0x2);
        assert!(total_tokens[2] == @0x3);
        assert!(total_tokens[3] == @0x4);
        assert!(total_tokens[4] == @0x5);

        ts::return_shared(ref);
    };

    ts::end(scenario);
}

// === Error Condition Tests ===

#[test]
#[expected_failure(abort_code = registry::ETokenNotRegistered)]
public fun test_transfer_admin_role_not_registered() {
    let mut scenario = create_test_scenario(CCIP_ADMIN);
    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    scenario.next_tx(CCIP_ADMIN);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        registry::transfer_admin_role(&mut ref, @0x2, @0x3, scenario.ctx());

        ts::return_shared(ref);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = registry::ENotAdministrator)]
public fun test_register_and_unregister_as_non_admin() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    let (treasury_cap, coin_metadata) = create_test_token(&mut scenario);
    let local_token = object::id_to_address(&object::id(&coin_metadata));

    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        register_test_pool(
            &mut ref,
            &treasury_cap,
            &coin_metadata,
            TOKEN_ADMIN_ADDRESS_2,
            scenario.ctx(),
        );

        let ctx = scenario.ctx();
        transfer::public_transfer(treasury_cap, ctx.sender());
        ts::return_shared(ref);
    };

    scenario.next_tx(RANDOM_USER);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        registry::unregister_pool(&mut ref, local_token, scenario.ctx());

        ts::return_shared(ref);
    };

    transfer::public_freeze_object(coin_metadata);
    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = registry::ETokenAddressNotRegistered)]
public fun test_get_all_configured_tokens_non_existent() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        registry::insert_token_configs_for_test(
            &mut ref,
            TOKEN_ADMIN_ADDRESS,
            vector[@0x1, @0x2, @0x3],
            TypeProof {},
        );

        // Test starting from key between existing tokens
        let (res, _next_key, has_more) = registry::get_all_configured_tokens(&ref, @0x1, 1);
        assert!(res.length() == 1);
        assert!(res[0] == @0x2);
        assert!(has_more);

        // Test starting from non-existent key - this should fail
        let (_res, _next_key, _has_more) = registry::get_all_configured_tokens(&ref, @0x4, 1);

        ts::return_shared(ref);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = registry::EAlreadyInitialized)]
public fun test_initialize_already_initialized() {
    let mut scenario = create_test_scenario(CCIP_ADMIN);
    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    // Try to initialize again - should fail
    scenario.next_tx(CCIP_ADMIN);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();
        let ctx = scenario.ctx();

        // This should fail because registry is already initialized
        registry::initialize(&mut ref, &owner_cap, ctx);

        scenario.return_to_sender(owner_cap);
        ts::return_shared(ref);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = registry::ETokenAlreadyRegistered)]
public fun test_register_pool_already_registered() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    let (treasury_cap, coin_metadata) = create_test_token(&mut scenario);

    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        // Register pool first time
        register_test_pool(
            &mut ref,
            &treasury_cap,
            &coin_metadata,
            TOKEN_ADMIN_ADDRESS,
            scenario.ctx(),
        );

        // Try to register the same token again - should fail
        register_test_pool(
            &mut ref,
            &treasury_cap,
            &coin_metadata,
            TOKEN_ADMIN_ADDRESS,
            scenario.ctx(),
        );

        let ctx = scenario.ctx();
        transfer::public_transfer(treasury_cap, ctx.sender());
        ts::return_shared(ref);
    };

    transfer::public_freeze_object(coin_metadata);
    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = registry::ENotAdministrator)]
public fun test_transfer_admin_role_not_administrator() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    let (treasury_cap, coin_metadata) = create_test_token(&mut scenario);
    let local_token = object::id_to_address(&object::id(&coin_metadata));

    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        // Register pool with TOKEN_ADMIN_ADDRESS as admin
        register_test_pool(
            &mut ref,
            &treasury_cap,
            &coin_metadata,
            TOKEN_ADMIN_ADDRESS,
            scenario.ctx(),
        );

        let ctx = scenario.ctx();
        transfer::public_transfer(treasury_cap, ctx.sender());
        ts::return_shared(ref);
    };

    // Try to transfer admin role as non-administrator - should fail
    scenario.next_tx(RANDOM_USER);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        // This should fail because RANDOM_USER is not the administrator
        registry::transfer_admin_role(&mut ref, local_token, TOKEN_ADMIN_ADDRESS_2, scenario.ctx());

        ts::return_shared(ref);
    };

    transfer::public_freeze_object(coin_metadata);
    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = registry::ENotPendingAdministrator)]
public fun test_accept_admin_role_not_pending() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    let (treasury_cap, coin_metadata) = create_test_token(&mut scenario);
    let local_token = object::id_to_address(&object::id(&coin_metadata));

    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        // Register pool with TOKEN_ADMIN_ADDRESS as admin
        register_test_pool(
            &mut ref,
            &treasury_cap,
            &coin_metadata,
            TOKEN_ADMIN_ADDRESS,
            scenario.ctx(),
        );

        // Request admin transfer to TOKEN_ADMIN_ADDRESS_2
        registry::transfer_admin_role(&mut ref, local_token, TOKEN_ADMIN_ADDRESS_2, scenario.ctx());

        let ctx = scenario.ctx();
        transfer::public_transfer(treasury_cap, ctx.sender());
        ts::return_shared(ref);
    };

    // Try to accept admin role as someone who is NOT the pending administrator - should fail
    scenario.next_tx(RANDOM_USER);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        // This should fail because RANDOM_USER is not the pending administrator
        // (TOKEN_ADMIN_ADDRESS_2 is the pending admin, not RANDOM_USER)
        registry::accept_admin_role(&mut ref, local_token, scenario.ctx());

        ts::return_shared(ref);
    };

    transfer::public_freeze_object(coin_metadata);
    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = registry::ENotPendingAdministrator)]
public fun test_accept_admin_role_no_pending_transfer() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    let (treasury_cap, coin_metadata) = create_test_token(&mut scenario);
    let local_token = object::id_to_address(&object::id(&coin_metadata));

    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        // Register pool with TOKEN_ADMIN_ADDRESS as admin
        register_test_pool(
            &mut ref,
            &treasury_cap,
            &coin_metadata,
            TOKEN_ADMIN_ADDRESS,
            scenario.ctx(),
        );

        // NOTE: No admin transfer request made

        let ctx = scenario.ctx();
        transfer::public_transfer(treasury_cap, ctx.sender());
        ts::return_shared(ref);
    };

    // Try to accept admin role when no transfer was requested - should fail
    scenario.next_tx(TOKEN_ADMIN_ADDRESS_2);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        // This should fail because no admin transfer was requested
        // (pending_administrator is @0x0)
        registry::accept_admin_role(&mut ref, local_token, scenario.ctx());

        ts::return_shared(ref);
    };

    transfer::public_freeze_object(coin_metadata);
    ts::end(scenario);
}

// ================================ MCMS Admin Transfer Tests ================================

#[test]
public fun test_mcms_transfer_admin_role() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    let (treasury_cap, coin_metadata) = create_test_token(&mut scenario);
    let local_token = object::id_address(&coin_metadata);
    let mcms = mcms_registry::get_multisig_address();

    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    // Register MCMS capability
    scenario.next_tx(CCIP_ADMIN);
    {
        let owner_cap = scenario.take_from_sender<OwnerCap>();
        let mut registry = scenario.take_shared<Registry>();

        registry::test_mcms_register_entrypoint(owner_cap, &mut registry, scenario.ctx());

        ts::return_shared(registry);
    };

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        register_test_pool(
            &mut ref,
            &treasury_cap,
            &coin_metadata,
            TOKEN_ADMIN_ADDRESS,
            scenario.ctx(),
        );

        let ctx = scenario.ctx();
        transfer::public_transfer(treasury_cap, ctx.sender());
        ts::return_shared(ref);
    };

    // Execute MCMS transfer (called by token administrator)
    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let registry = scenario.take_shared<Registry>();

        registry::transfer_admin_role(&mut ref, local_token, mcms, scenario.ctx());

        // Verify pending transfer is set but admin hasn't changed yet (still TOKEN_ADMIN_ADDRESS)
        assert!(registry::is_administrator(&ref, local_token, TOKEN_ADMIN_ADDRESS));
        assert!(!registry::is_administrator(&ref, local_token, mcms));

        ts::return_shared(ref);
        ts::return_shared(registry);
    };

    transfer::public_freeze_object(coin_metadata);
    ts::end(scenario);
}

#[test]
public fun test_mcms_accept_admin_role() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    let (treasury_cap, coin_metadata) = create_test_token(&mut scenario);
    let local_token = object::id_address(&coin_metadata);
    let mcms = mcms_registry::get_multisig_address();

    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    // Register MCMS capability
    scenario.next_tx(CCIP_ADMIN);
    {
        let owner_cap = scenario.take_from_sender<OwnerCap>();
        let mut registry = scenario.take_shared<Registry>();

        registry::test_mcms_register_entrypoint(owner_cap, &mut registry, scenario.ctx());

        ts::return_shared(registry);
    };

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        register_test_pool(
            &mut ref,
            &treasury_cap,
            &coin_metadata,
            TOKEN_ADMIN_ADDRESS,
            scenario.ctx(),
        );

        // set pending transfer to MCMS
        registry::transfer_admin_role(&mut ref, local_token, mcms, scenario.ctx());

        let ctx = scenario.ctx();
        transfer::public_transfer(treasury_cap, ctx.sender());
        ts::return_shared(ref);
    };

    // // Execute MCMS accept (called by pending administrator) (TOKEN_ADMIN_ADDRESS_2)
    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let mut registry = scenario.take_shared<Registry>();

        // Create MCMS callback params for accept (need to include object addresses first)
        let mut data = vector::empty<u8>();
        data.append(bcs::to_bytes(&object::id_address(&ref)));
        data.append(bcs::to_bytes(&local_token));

        let params = mcms_registry::test_create_executing_callback_params(
            @ccip,
            string::utf8(b"token_admin_registry"),
            string::utf8(b"accept_admin_role"),
            data,
            x"0000000000000000000000000000000000000000000000000000000000000001",
            0,
            1,
        );

        // Execute MCMS accept
        registry::mcms_accept_admin_role(&mut ref, &mut registry, params, scenario.ctx());

        // Verify admin has changed and no pending transfer (TOKEN_ADMIN_ADDRESS_2 is now the admin)
        assert!(!registry::is_administrator(&ref, local_token, TOKEN_ADMIN_ADDRESS));
        assert!(registry::is_administrator(&ref, local_token, mcms));

        ts::return_shared(ref);
        ts::return_shared(registry);
    };

    transfer::public_freeze_object(coin_metadata);
    ts::end(scenario);
}

#[test]
public fun test_mcms_full_admin_transfer_flow() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    let (treasury_cap, coin_metadata) = create_test_token(&mut scenario);
    let local_token = object::id_address(&coin_metadata);
    let mcms = mcms_registry::get_multisig_address();

    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    // Register MCMS capability
    scenario.next_tx(CCIP_ADMIN);
    {
        let owner_cap = scenario.take_from_sender<OwnerCap>();
        let mut registry = scenario.take_shared<Registry>();

        registry::test_mcms_register_entrypoint(owner_cap, &mut registry, scenario.ctx());

        ts::return_shared(registry);
    };

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        register_test_pool(
            &mut ref,
            &treasury_cap,
            &coin_metadata,
            TOKEN_ADMIN_ADDRESS,
            scenario.ctx(),
        );

        let ctx = scenario.ctx();
        transfer::public_transfer(treasury_cap, ctx.sender());
        ts::return_shared(ref);
    };

    // Step 1: MCMS Transfer (called by token administrator)
    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let registry = scenario.take_shared<Registry>();
        registry::transfer_admin_role(&mut ref, local_token, mcms, scenario.ctx());

        // Verify pending state (TOKEN_ADMIN_ADDRESS is still the admin)
        assert!(registry::is_administrator(&ref, local_token, TOKEN_ADMIN_ADDRESS));
        assert!(!registry::is_administrator(&ref, local_token, mcms));

        ts::return_shared(ref);
        ts::return_shared(registry);
    };

    // Step 2: MCMS Accept
    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let mut registry = scenario.take_shared<Registry>();

        let mut data = vector::empty<u8>();
        data.append(bcs::to_bytes(&object::id_address(&ref)));
        data.append(bcs::to_bytes(&local_token));

        let params = mcms_registry::test_create_executing_callback_params(
            @ccip,
            string::utf8(b"token_admin_registry"),
            string::utf8(b"accept_admin_role"),
            data,
            x"0000000000000000000000000000000000000000000000000000000000000002",
            0,
            1,
        );

        registry::mcms_accept_admin_role(&mut ref, &mut registry, params, scenario.ctx());

        // Verify final state
        assert!(!registry::is_administrator(&ref, local_token, TOKEN_ADMIN_ADDRESS));
        assert!(registry::is_administrator(&ref, local_token, mcms));

        ts::return_shared(ref);
        ts::return_shared(registry);
    };

    transfer::public_freeze_object(coin_metadata);
    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = registry::ENotPendingAdministrator)]
public fun test_mcms_accept_admin_role_no_pending_transfer_fails() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    let (treasury_cap, coin_metadata) = create_test_token(&mut scenario);
    let local_token = object::id_address(&coin_metadata);

    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);
    scenario.next_tx(CCIP_ADMIN);
    {
        let owner_cap = scenario.take_from_sender<OwnerCap>();
        let mut registry = scenario.take_shared<Registry>();

        registry::test_mcms_register_entrypoint(owner_cap, &mut registry, scenario.ctx());

        ts::return_shared(registry);
    };

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        register_test_pool(
            &mut ref,
            &treasury_cap,
            &coin_metadata,
            TOKEN_ADMIN_ADDRESS,
            scenario.ctx(),
        );

        let ctx = scenario.ctx();
        transfer::public_transfer(treasury_cap, ctx.sender());
        ts::return_shared(ref);
    };

    // Try to accept without pending transfer - should fail (TOKEN_ADMIN_ADDRESS is still the admin)
    scenario.next_tx(CCIP_ADMIN);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let mut registry = scenario.take_shared<Registry>();

        let mut data = vector::empty<u8>();
        data.append(bcs::to_bytes(&object::id_address(&ref)));
        data.append(bcs::to_bytes(&local_token));

        let params = mcms_registry::test_create_executing_callback_params(
            @ccip,
            string::utf8(b"token_admin_registry"),
            string::utf8(b"accept_admin_role"),
            data,
            x"0000000000000000000000000000000000000000000000000000000000000003",
            0,
            1,
        );

        registry::mcms_accept_admin_role(&mut ref, &mut registry, params, scenario.ctx());

        ts::return_shared(ref);
        ts::return_shared(registry);
    };

    transfer::public_freeze_object(coin_metadata);
    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = registry::ETokenNotRegistered)]
public fun test_mcms_transfer_admin_role_token_not_registered_fails() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    let mcms = mcms_registry::get_multisig_address();

    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    // Register MCMS capability
    scenario.next_tx(CCIP_ADMIN);
    {
        let owner_cap = scenario.take_from_sender<OwnerCap>();
        let mut registry = scenario.take_shared<Registry>();

        registry::test_mcms_register_entrypoint(owner_cap, &mut registry, scenario.ctx());

        ts::return_shared(registry);
    };

    // Try to transfer admin for unregistered token - should fail
    scenario.next_tx(CCIP_ADMIN);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let mut registry = scenario.take_shared<Registry>();

        let mut data = vector::empty<u8>();
        data.append(bcs::to_bytes(&object::id_address(&ref)));
        data.append(bcs::to_bytes(&@0x999)); // unregistered token
        data.append(bcs::to_bytes(&mcms));

        let params = mcms_registry::test_create_executing_callback_params(
            @ccip,
            string::utf8(b"token_admin_registry"),
            string::utf8(b"transfer_admin_role"),
            data,
            x"0000000000000000000000000000000000000000000000000000000000000004",
            0,
            1,
        );

        registry::mcms_transfer_admin_role(&mut ref, &mut registry, params, scenario.ctx());

        ts::return_shared(ref);
        ts::return_shared(registry);
    };

    ts::end(scenario);
}

// === Upgrade Registry Function Restriction Tests ===

#[test]
#[expected_failure(abort_code = upgrade_registry::EFunctionNotAllowed)]
public fun test_register_pool_function_not_allowed() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    let (treasury_cap, coin_metadata) = create_test_token(&mut scenario);

    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    scenario.next_tx(CCIP_ADMIN);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();
        let ctx = scenario.ctx();

        // Block the register_pool function using upgrade registry
        upgrade_registry::block_function(
            &mut ref,
            &owner_cap,
            string::utf8(b"token_admin_registry"),
            string::utf8(b"register_pool"),
            1, // block version 1
            ctx,
        );

        scenario.return_to_sender(owner_cap);
        ts::return_shared(ref);
    };

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();

        // This should fail because the function is blocked by upgrade registry
        register_test_pool(
            &mut ref,
            &treasury_cap,
            &coin_metadata,
            TOKEN_ADMIN_ADDRESS,
            scenario.ctx(),
        );

        let ctx = scenario.ctx();
        transfer::public_transfer(treasury_cap, ctx.sender());
        ts::return_shared(ref);
    };

    transfer::public_freeze_object(coin_metadata);
    ts::end(scenario);
}
