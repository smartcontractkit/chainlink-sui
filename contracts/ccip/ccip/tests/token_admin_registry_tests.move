#[test_only]
module ccip::token_admin_registry_tests;

use std::ascii;
use std::string;
use std::type_name;

use sui::coin;
use sui::test_scenario::{Self as ts, Scenario};

use ccip::token_admin_registry as registry;
use ccip::state_object::{Self, OwnerCap, CCIPObjectRef};

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

// Mock pool addresses
const MOCK_TOKEN_POOL_PACKAGE_ID_1: address = @0x1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b;
const MOCK_TOKEN_POOL_PACKAGE_ID_2: address = @0x8a7b6c5d4e3f2a1b0c9d8e7f6a5b4c3d2e1f0a9b8c7d6e5f4a3b2c1d0e9f8a7;

// === Helper Functions ===

fun create_test_scenario(addr: address): Scenario {
    ts::begin(addr)
}

fun initialize_state_and_registry(scenario: &mut Scenario, admin: address) {
    scenario.next_tx(admin);
    {
        let ctx = scenario.ctx();
        state_object::test_init(ctx);
    };

    scenario.next_tx(admin);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();
        let ctx = scenario.ctx();
        
        registry::initialize(&mut ref, &owner_cap, ctx);
        
        scenario.return_to_sender(owner_cap);
        ts::return_shared(ref);
    };
}

fun create_test_token(scenario: &mut Scenario): (coin::TreasuryCap<TOKEN_ADMIN_REGISTRY_TESTS>, coin::CoinMetadata<TOKEN_ADMIN_REGISTRY_TESTS>) {
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

fun create_named_test_token(
    scenario: &mut Scenario, 
    symbol: vector<u8>, 
    name: vector<u8>, 
    description: vector<u8>
): (coin::TreasuryCap<TOKEN_ADMIN_REGISTRY_TESTS>, coin::CoinMetadata<TOKEN_ADMIN_REGISTRY_TESTS>) {
    coin::create_currency(
        TOKEN_ADMIN_REGISTRY_TESTS {},
        DECIMALS,
        symbol,
        name,
        description,
        option::none(),
        scenario.ctx(),
    )
}

fun register_test_pool<T>(
    ref: &mut CCIPObjectRef,
    treasury_cap: &coin::TreasuryCap<T>,
    coin_metadata: &coin::CoinMetadata<T>,
    pool_package_id: address,
    pool_module: vector<u8>,
    admin: address
) {
    registry::register_pool(
        ref,
        treasury_cap,
        coin_metadata,
        pool_package_id,
        string::utf8(pool_module),
        admin,
        vector[], // lock_or_burn_params
        vector[], // release_or_mint_params
        TypeProof {},
    );
}

fun assert_empty_token_config(
    ref: &CCIPObjectRef,
    token_address: address
) {
    let token_config = registry::get_token_config(ref, token_address);
    let (
        token_pool_package_id,
        token_pool_module,
        token_type,
        administrator,
        pending_administrator,
        proof,
        _lock_or_burn_params,
        _release_or_mint_params,
    ) = registry::get_token_config_data(token_config);

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
    expected_pending_admin: address
) {
    let token_config = registry::get_token_config(ref, token_address);
    let (
        token_pool_package_id,
        token_pool_module,
        token_type,
        administrator,
        pending_administrator,
        _proof,
        _lock_or_burn_params,
        _release_or_mint_params,
    ) = registry::get_token_config_data(token_config);

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
            MOCK_TOKEN_POOL_PACKAGE_ID_1,
            b"mock_token_pool",
            TOKEN_ADMIN_ADDRESS
        );

        // Test with registered token
        let pool_address = registry::get_pool(&ref, local_token);
        assert!(pool_address == MOCK_TOKEN_POOL_PACKAGE_ID_1);
        
        let ctx = scenario.ctx();
        transfer::public_transfer(treasury_cap, ctx.sender());
        ts::return_shared(ref);
    };

    transfer::public_freeze_object(coin_metadata);
    ts::end(scenario);
}

#[test]
public fun test_register_pool_by_admin() {
    let mut scenario = create_test_scenario(CCIP_ADMIN);
    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    scenario.next_tx(CCIP_ADMIN);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();
        let ctx = scenario.ctx();
        
        // Register pool as admin (without treasury cap)
        registry::register_pool_by_admin(
            &mut ref,
            &owner_cap,
            @0x123, // coin_metadata_address
            MOCK_TOKEN_POOL_PACKAGE_ID_1, // token_pool_package_id
            string::utf8(b"admin_registered_pool"), // token_pool_module
            ascii::string(b"TestType"),
            TOKEN_ADMIN_ADDRESS, // initial_administrator
            ascii::string(b"AdminProof"), // proof
            vector[], // lock_or_burn_params
            vector[], // release_or_mint_params
            ctx,
        );

        // Verify registration
        let pool_address = registry::get_pool(&ref, @0x123);
        assert!(pool_address == MOCK_TOKEN_POOL_PACKAGE_ID_1);
        assert!(registry::is_administrator(&ref, @0x123, TOKEN_ADMIN_ADDRESS));
        
        scenario.return_to_sender(owner_cap);
        ts::return_shared(ref);
    };

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
            MOCK_TOKEN_POOL_PACKAGE_ID_1,
            b"mock_token_pool",
            TOKEN_ADMIN_ADDRESS_2
        );

        // Verify registration
        let pool_addresses = registry::get_pools(&ref, vector[local_token]);
        assert!(pool_addresses.length() == 1);
        assert!(pool_addresses[0] == MOCK_TOKEN_POOL_PACKAGE_ID_1);
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
public fun test_register_and_set_pool() {
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
            MOCK_TOKEN_POOL_PACKAGE_ID_1,
            b"mock_token_pool",
            TOKEN_ADMIN_ADDRESS
        );

        // Verify initial registration
        let pool_addresses = registry::get_pools(&ref, vector[local_token]);
        assert!(pool_addresses.length() == 1);
        assert!(pool_addresses[0] == MOCK_TOKEN_POOL_PACKAGE_ID_1);
        assert!(registry::is_administrator(&ref, local_token, TOKEN_ADMIN_ADDRESS));

        // Verify detailed configuration
        assert_token_config(
            &ref,
            local_token,
            MOCK_TOKEN_POOL_PACKAGE_ID_1,
            b"mock_token_pool",
            type_name::get<TOKEN_ADMIN_REGISTRY_TESTS>().into_string(),
            TOKEN_ADMIN_ADDRESS,
            @0x0
        );

        let token_config = registry::get_token_config(&ref, local_token);
        let (_, _, token_type, _, _, type_proof, _, _) = registry::get_token_config_data(token_config);
        assert!(token_type == ascii::string(b"0000000000000000000000000000000000000000000000000000000000001000::token_admin_registry_tests::TOKEN_ADMIN_REGISTRY_TESTS"));
        assert!(type_proof == type_name::into_string(type_name::get<TypeProof>()));

        let ctx = scenario.ctx();

        // Update pool configuration
        registry::set_pool(
            &mut ref,
            local_token,
            MOCK_TOKEN_POOL_PACKAGE_ID_2,   
            string::utf8(b"mock_token_pool_2"),
            vector[], // lock_or_burn_params
            vector[], // release_or_mint_params
            TypeProof2 {},
            ctx,
        );

        // Request admin transfer
        registry::transfer_admin_role(&mut ref, local_token, TOKEN_ADMIN_ADDRESS_2, ctx);
        
        transfer::public_transfer(treasury_cap, ctx.sender());
        ts::return_shared(ref);
    };

    scenario.next_tx(TOKEN_ADMIN_ADDRESS_2);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        
        // Verify updated configuration
        assert_token_config(
            &ref,
            local_token,
            MOCK_TOKEN_POOL_PACKAGE_ID_2,
            b"mock_token_pool_2",
            type_name::get<TOKEN_ADMIN_REGISTRY_TESTS>().into_string(),
            TOKEN_ADMIN_ADDRESS,
            TOKEN_ADMIN_ADDRESS_2
        );

        let token_config = registry::get_token_config(&ref, local_token);
        let (_, _, token_type, _, _, type_proof, _, _) = registry::get_token_config_data(token_config);
        assert!(token_type == ascii::string(b"0000000000000000000000000000000000000000000000000000000000001000::token_admin_registry_tests::TOKEN_ADMIN_REGISTRY_TESTS"));
        assert!(type_proof == type_name::into_string(type_name::get<TypeProof2>()));

        // Accept admin role
        registry::accept_admin_role(&mut ref, local_token, scenario.ctx());
        assert!(registry::is_administrator(&ref, local_token, TOKEN_ADMIN_ADDRESS_2));
        
        ts::return_shared(ref);
    };

    transfer::public_freeze_object(coin_metadata);
    ts::end(scenario);
}

// === Pool Information Retrieval Tests ===

#[test]
public fun test_get_token_configs() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    
    // Create two test tokens
    let (treasury_cap_1, coin_metadata_1) = create_named_test_token(&mut scenario, b"TEST1", b"TestToken1", b"test_token_1");
    let token_1 = object::id_to_address(&object::id(&coin_metadata_1));
    
    let (treasury_cap_2, coin_metadata_2) = create_named_test_token(&mut scenario, b"TEST2", b"TestToken2", b"test_token_2");
    let token_2 = object::id_to_address(&object::id(&coin_metadata_2));

    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        
        // Register both tokens
        register_test_pool(
            &mut ref,
            &treasury_cap_1,
            &coin_metadata_1,
            MOCK_TOKEN_POOL_PACKAGE_ID_1,
            b"mock_token_pool_1",
            TOKEN_ADMIN_ADDRESS
        );
        
        registry::register_pool(
            &mut ref,
            &treasury_cap_2,
            &coin_metadata_2,
            MOCK_TOKEN_POOL_PACKAGE_ID_2,
            string::utf8(b"mock_token_pool_2"),
            TOKEN_ADMIN_ADDRESS,
            vector[], // lock_or_burn_params
            vector[], // release_or_mint_params
            TypeProof2 {},
        );

        // Test various scenarios
        let token_configs = registry::get_token_configs(&ref, vector[token_1, token_2]);
        let unregistered_token = @0x999;
        let mixed_token_configs = registry::get_token_configs(&ref, vector[token_1, unregistered_token]);
        let empty_token_configs = registry::get_token_configs(&ref, vector[]);
        
        // Verify functionality through get_token_configs results
        assert!(token_configs.length() == 2);
        let (token_pool_package_id_1, _, _, _, _, _, _, _) = registry::get_token_config_data(token_configs[0]);
        let (token_pool_package_id_2, _, _, _, _, _, _, _) = registry::get_token_config_data(token_configs[1]);
        assert!(token_pool_package_id_1 == MOCK_TOKEN_POOL_PACKAGE_ID_1);
        assert!(token_pool_package_id_2 == MOCK_TOKEN_POOL_PACKAGE_ID_2);
        
        // Test with mixed registered/unregistered tokens
        assert!(mixed_token_configs.length() == 2);
        let (mixed_token_pool_package_id_1, _, _, _, _, _, _, _) = registry::get_token_config_data(mixed_token_configs[0]);
        let (mixed_token_pool_package_id_2, _, _, _, _, _, _, _) = registry::get_token_config_data(mixed_token_configs[1]);
        assert!(mixed_token_pool_package_id_1 == MOCK_TOKEN_POOL_PACKAGE_ID_1);
        assert!(mixed_token_pool_package_id_2 == @0x0); // unregistered token
        
        // Test with empty vector
        assert!(empty_token_configs.length() == 0);
        
        let ctx = scenario.ctx();
        transfer::public_transfer(treasury_cap_1, ctx.sender());
        transfer::public_transfer(treasury_cap_2, ctx.sender());
        ts::return_shared(ref);
    };

    transfer::public_freeze_object(coin_metadata_1);
    transfer::public_freeze_object(coin_metadata_2);
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
        
        registry::insert_token_configs_for_test(&mut ref, vector[@0x1, @0x2, @0x3], TypeProof {});

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
        registry::insert_token_configs_for_test(&mut ref, vector[@0x1], TypeProof {});
        let (res, _next_key, has_more) = registry::get_all_configured_tokens(&ref, @0x0, 1);
        assert!(res.length() == 1);
        assert!(res[0] == @0x1);
        assert!(!has_more);

        // Test case 3: Start from middle
        registry::insert_token_configs_for_test(&mut ref, vector[@0x2, @0x3], TypeProof {});
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
        
        registry::insert_token_configs_for_test(&mut ref, vector[@0x1, @0x2, @0x3, @0x4, @0x5], TypeProof {});

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

// === Pool Configuration Management Tests ===

#[test]
public fun test_set_pool_comprehensive() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    let (treasury_cap, coin_metadata) = create_test_token(&mut scenario);
    let local_token = object::id_to_address(&object::id(&coin_metadata));

    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        
        // Register initial pool
        register_test_pool(
            &mut ref,
            &treasury_cap,
            &coin_metadata,
            MOCK_TOKEN_POOL_PACKAGE_ID_1,
            b"initial_token_pool",
            TOKEN_ADMIN_ADDRESS
        );

        // Verify initial configuration
        assert_token_config(
            &ref,
            local_token,
            MOCK_TOKEN_POOL_PACKAGE_ID_1,
            b"initial_token_pool",
            type_name::get<TOKEN_ADMIN_REGISTRY_TESTS>().into_string(),
            TOKEN_ADMIN_ADDRESS,
            @0x0
        );

        let token_config = registry::get_token_config(&ref, local_token);
        let (_, _, _, _, _, type_proof, _, _) = registry::get_token_config_data(token_config);
        assert!(type_proof == type_name::into_string(type_name::get<TypeProof>()));

        let ctx = scenario.ctx();

        // Test set_pool with different package ID (should update)
        registry::set_pool(
            &mut ref,
            local_token,
            MOCK_TOKEN_POOL_PACKAGE_ID_2,
            string::utf8(b"updated_token_pool"),
            vector[], // lock_or_burn_params
            vector[], // release_or_mint_params
            TypeProof2 {},
            ctx,
        );

        // Verify pool was updated
        assert_token_config(
            &ref,
            local_token,
            MOCK_TOKEN_POOL_PACKAGE_ID_2,
            b"updated_token_pool",
            type_name::get<TOKEN_ADMIN_REGISTRY_TESTS>().into_string(),
            TOKEN_ADMIN_ADDRESS,
            @0x0
        );

        let token_config = registry::get_token_config(&ref, local_token);
        let (_, _, _, _, _, updated_type_proof, _, _) = registry::get_token_config_data(token_config);
        assert!(updated_type_proof == type_name::into_string(type_name::get<TypeProof2>()));

        // Test set_pool with same package ID (should not trigger update)
        registry::set_pool(
            &mut ref,
            local_token,
            MOCK_TOKEN_POOL_PACKAGE_ID_2, // same package ID
            string::utf8(b"should_not_update"),
            vector[], // lock_or_burn_params
            vector[], // release_or_mint_params
            TypeProof {},
            ctx,
        );

        // Verify pool was NOT updated (same package ID means no change)
        assert_token_config(
            &ref,
            local_token,
            MOCK_TOKEN_POOL_PACKAGE_ID_2,
            b"updated_token_pool", // unchanged
            type_name::get<TOKEN_ADMIN_REGISTRY_TESTS>().into_string(),
            TOKEN_ADMIN_ADDRESS,
            @0x0
        );

        let token_config = registry::get_token_config(&ref, local_token);
        let (_, _, _, _, _, final_type_proof, _, _) = registry::get_token_config_data(token_config);
        assert!(final_type_proof == type_name::into_string(type_name::get<TypeProof2>())); // unchanged
        
        transfer::public_transfer(treasury_cap, ctx.sender());
        ts::return_shared(ref);
    };

    transfer::public_freeze_object(coin_metadata);
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
#[expected_failure(abort_code = registry::ENotAllowed)]
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
            MOCK_TOKEN_POOL_PACKAGE_ID_1,
            b"mock_token_pool",
            TOKEN_ADMIN_ADDRESS_2
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
        
        registry::insert_token_configs_for_test(&mut ref, vector[@0x1, @0x2, @0x3], TypeProof {});

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
#[expected_failure(abort_code = registry::ETokenNotRegistered)]
public fun test_set_pool_unregistered_token() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let ctx = scenario.ctx();
        
        // Try to set pool for unregistered token - should fail
        registry::set_pool(
            &mut ref,
            @0x999, // unregistered token
            MOCK_TOKEN_POOL_PACKAGE_ID_1,
            string::utf8(b"test_pool"),
            vector[], // lock_or_burn_params
            vector[], // release_or_mint_params
            TypeProof {},
            ctx,
        );
        
        ts::return_shared(ref);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = registry::ENotAllowed)]
public fun test_set_pool_unauthorized() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    let (treasury_cap, coin_metadata) = create_test_token(&mut scenario);
    let local_token = object::id_to_address(&object::id(&coin_metadata));

    initialize_state_and_registry(&mut scenario, CCIP_ADMIN);

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        
        // Register pool with TOKEN_ADMIN_ADDRESS_2 as admin
        register_test_pool(
            &mut ref,
            &treasury_cap,
            &coin_metadata,
            MOCK_TOKEN_POOL_PACKAGE_ID_1,
            b"test_pool",
            TOKEN_ADMIN_ADDRESS_2
        );
        
        let ctx = scenario.ctx();
        transfer::public_transfer(treasury_cap, ctx.sender());
        ts::return_shared(ref);
    };

    // Try to set pool as unauthorized user (not admin or owner)
    scenario.next_tx(RANDOM_USER);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let ctx = scenario.ctx();
        
        // Should fail - RANDOM_USER is not the administrator or owner
        registry::set_pool(
            &mut ref,
            local_token,
            MOCK_TOKEN_POOL_PACKAGE_ID_2,
            string::utf8(b"unauthorized_update"),
            vector[], // lock_or_burn_params
            vector[], // release_or_mint_params
            TypeProof2 {},
            ctx,
        );
        
        ts::return_shared(ref);
    };

    transfer::public_freeze_object(coin_metadata);
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
            MOCK_TOKEN_POOL_PACKAGE_ID_1,
            b"first_pool",
            TOKEN_ADMIN_ADDRESS
        );

        // Try to register the same token again - should fail
        register_test_pool(
            &mut ref,
            &treasury_cap,
            &coin_metadata,
            MOCK_TOKEN_POOL_PACKAGE_ID_2,
            b"second_pool",
            TOKEN_ADMIN_ADDRESS
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
            MOCK_TOKEN_POOL_PACKAGE_ID_1,
            b"test_pool",
            TOKEN_ADMIN_ADDRESS
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
            MOCK_TOKEN_POOL_PACKAGE_ID_1,
            b"test_pool",
            TOKEN_ADMIN_ADDRESS
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
            MOCK_TOKEN_POOL_PACKAGE_ID_1,
            b"test_pool",
            TOKEN_ADMIN_ADDRESS
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
