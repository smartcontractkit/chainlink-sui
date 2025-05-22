#[test_only]
module ccip::token_admin_registry_tests;

use std::type_name;

use sui::coin;
use sui::test_scenario::{Self as ts, Scenario};

use ccip::token_admin_registry as registry;
use ccip::state_object::{Self, OwnerCap, CCIPObjectRef};

public struct TOKEN_ADMIN_REGISTRY_TESTS has drop {}

public struct TypeProof has drop {}
public struct TypeProof2 has drop {}

const Decimals: u8 = 8;
const CCIP_ADMIN: address = @0x1000;
const TOKEN_ADMIN_ADDRESS: address = @0x1;
const TOKEN_ADMIN_ADDRESS_2: address = @0x2;
const RANDOM_USER: address = @0x3;
const MOCK_TOKEN_POOL_ADDRESS_1: address = @0x1a2b3c4d5e6f7a8b9c0d1e2f3a4b5c6d7e8f9a0b1c2d3e4f5a6b7c8d9e0f1a2b;
const MOCK_TOKEN_POOL_ADDRESS_2: address = @0x8a7b6c5d4e3f2a1b0c9d8e7f6a5b4c3d2e1f0a9b8c7d6e5f4a3b2c1d0e9f8a7;

fun create_test_scenario(addr: address): Scenario {
    ts::begin(addr)
}

#[test]
public fun test_initialize() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    {
        let ctx = scenario.ctx();
        state_object::test_init(ctx);
    };

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        let ctx = scenario.ctx();
        registry::initialize(&mut ref, &owner_cap, ctx);
        let (token_pool_address, administrator, pending_administrator, proof) = registry::get_token_config(&ref, @0x2);
        assert!(proof.is_none());
        assert!(token_pool_address == @0x0);
        assert!(administrator == @0x0);
        assert!(pending_administrator == @0x0);

        scenario.return_to_sender(owner_cap);
        ts::return_shared(ref);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = registry::E_TOKEN_NOT_REGISTERED)]
public fun test_transfer_admin_role_not_registered() {
    let mut scenario = create_test_scenario(CCIP_ADMIN);
    {
        let ctx = scenario.ctx();
        state_object::test_init(ctx);
    };

    scenario.next_tx(CCIP_ADMIN);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        let ctx = scenario.ctx();
        registry::initialize(&mut ref, &owner_cap, ctx);

        registry::transfer_admin_role(&mut ref, @0x2, @0x3, ctx);

        scenario.return_to_sender(owner_cap);
        ts::return_shared(ref);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = registry::E_NOT_ADMINISTRATOR)]
public fun test_register_and_unregister_as_non_admin() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    let (treasury_cap, coin_metadata) = coin::create_currency(
        TOKEN_ADMIN_REGISTRY_TESTS {},
        Decimals,
        b"TEST",
        b"TestToken",
        b"test_token",
        option::none(),
        scenario.ctx(),
    );
    let local_token = object::id_to_address(&object::id(&coin_metadata));
    transfer::public_freeze_object(coin_metadata);

    scenario.next_tx(CCIP_ADMIN);
    {
        let ctx = scenario.ctx();
        state_object::test_init(ctx);
    };

    scenario.next_tx(CCIP_ADMIN);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        let ctx = scenario.ctx();
        registry::initialize(&mut ref, &owner_cap, ctx);
        scenario.return_to_sender(owner_cap);
        ts::return_shared(ref);
    };

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        registry::register_pool(
            &mut ref,
            &treasury_cap,
            local_token,
            MOCK_TOKEN_POOL_ADDRESS_1,
            TOKEN_ADMIN_ADDRESS_2,
            TypeProof {},
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

    ts::end(scenario);
}

#[test]
public fun test_register_and_unregister() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    let (treasury_cap, coin_metadata) = coin::create_currency(
        TOKEN_ADMIN_REGISTRY_TESTS {},
        Decimals,
        b"TEST",
        b"TestToken",
        b"test_token",
        option::none(),
        scenario.ctx(),
    );
    let local_token = object::id_to_address(&object::id(&coin_metadata));
    transfer::public_freeze_object(coin_metadata);

    scenario.next_tx(CCIP_ADMIN);
    {
        let ctx = scenario.ctx();
        state_object::test_init(ctx);
    };

    scenario.next_tx(CCIP_ADMIN);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        let ctx = scenario.ctx();
        registry::initialize(&mut ref, &owner_cap, ctx);
        scenario.return_to_sender(owner_cap);
        ts::return_shared(ref);
    };

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        registry::register_pool(
            &mut ref,
            &treasury_cap,
            local_token,
            MOCK_TOKEN_POOL_ADDRESS_1,
            TOKEN_ADMIN_ADDRESS_2, // initial admin
            TypeProof {},
        );

        let pool_addresses = registry::get_pools(&ref, vector[local_token]);
        assert!(pool_addresses.length() == 1);
        assert!(pool_addresses[0] == MOCK_TOKEN_POOL_ADDRESS_1);
        assert!(registry::is_administrator(&ref, local_token, TOKEN_ADMIN_ADDRESS_2));

        let ctx = scenario.ctx();
        transfer::public_transfer(treasury_cap, ctx.sender());
        ts::return_shared(ref);
    };

    // unregister the token as the token admin
    scenario.next_tx(TOKEN_ADMIN_ADDRESS_2);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        registry::unregister_pool(&mut ref, local_token, scenario.ctx());
        let (token_pool_address, administrator, pending_administrator, type_proof_op) = registry::get_token_config(&ref, local_token);

        assert!(token_pool_address == @0x0);
        assert!(administrator == @0x0);
        assert!(pending_administrator == @0x0);
        assert!(type_proof_op.is_none());

        ts::return_shared(ref);
    };

    ts::end(scenario);
}

#[test]
public fun test_register_and_set_pool() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    let (treasury_cap, coin_metadata) = coin::create_currency(
        TOKEN_ADMIN_REGISTRY_TESTS {},
        Decimals,
        b"TEST",
        b"TestToken",
        b"test_token",
        option::none(),
        scenario.ctx(),
    );
    let local_token = object::id_to_address(&object::id(&coin_metadata));
    transfer::public_freeze_object(coin_metadata);

    scenario.next_tx(CCIP_ADMIN);
    {
        let ctx = scenario.ctx();
        state_object::test_init(ctx);
    };

    scenario.next_tx(CCIP_ADMIN);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        let ctx = scenario.ctx();
        registry::initialize(&mut ref, &owner_cap, ctx);
        scenario.return_to_sender(owner_cap);
        ts::return_shared(ref);
    };

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        registry::register_pool(
            &mut ref,
            &treasury_cap,
            local_token,
            MOCK_TOKEN_POOL_ADDRESS_1,
            TOKEN_ADMIN_ADDRESS,
            TypeProof {},
        );

        let pool_addresses = registry::get_pools(&ref, vector[local_token]);
        assert!(pool_addresses.length() == 1);
        assert!(pool_addresses[0] == MOCK_TOKEN_POOL_ADDRESS_1);
        assert!(registry::is_administrator(&ref, local_token, TOKEN_ADMIN_ADDRESS));

        let (token_pool_address, administrator, pending_administrator, proof_op) = registry::get_token_config(&ref, local_token);
        assert!(proof_op.is_some());
        assert!(token_pool_address == MOCK_TOKEN_POOL_ADDRESS_1);
        assert!(administrator == TOKEN_ADMIN_ADDRESS);
        assert!(pending_administrator == @0x0);

        let proof = proof_op.borrow();
        assert!(proof == type_name::get<TypeProof>());

        let ctx = scenario.ctx();

        registry::set_pool(&mut ref, local_token, MOCK_TOKEN_POOL_ADDRESS_2, TypeProof2 {}, ctx);
        registry::transfer_admin_role(&mut ref, local_token, TOKEN_ADMIN_ADDRESS_2, ctx);

        transfer::public_transfer(treasury_cap, ctx.sender());
        ts::return_shared(ref);
    };

    scenario.next_tx(TOKEN_ADMIN_ADDRESS_2);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let (token_pool_address, administrator, pending_administrator, proof_op) = registry::get_token_config(&ref, local_token);
        assert!(token_pool_address == MOCK_TOKEN_POOL_ADDRESS_2);
        assert!(administrator == TOKEN_ADMIN_ADDRESS);
        assert!(pending_administrator == TOKEN_ADMIN_ADDRESS_2);
        assert!(proof_op.is_some());
        let proof = proof_op.borrow();
        // after set_pool, the proof should be TypeProof2 bc it will come from a different pool
        assert!(proof == type_name::get<TypeProof2>());

        registry::accept_admin_role(&mut ref, local_token, scenario.ctx());
        assert!(registry::is_administrator(&ref, local_token, TOKEN_ADMIN_ADDRESS_2));
        ts::return_shared(ref);
    };

    ts::end(scenario);
}

#[test]
public fun test_get_all_configured_tokens() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    {
        let ctx = scenario.ctx();
        state_object::test_init(ctx);
    };

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        let ctx = scenario.ctx();
        registry::initialize(&mut ref, &owner_cap, ctx);

        registry::insert_token_addresses_for_test(&mut ref, vector[@0x1, @0x2, @0x3], TypeProof {});

        let (res, next_key, has_more) = registry::get_all_configured_tokens(&ref, @0x0, 0);
        assert!(res.length() == 0);
        assert!(next_key == @0x0);
        assert!(has_more);

        let (res, next_key, has_more) = registry::get_all_configured_tokens(&ref, @0x0, 3);
        assert!(res.length() == 3);
        assert!(vector[@0x1, @0x2, @0x3] == res);
        assert!(next_key == @0x3);
        assert!(!has_more);

        scenario.return_to_sender(owner_cap);
        ts::return_shared(ref);
    };

    ts::end(scenario);
}

#[test]
public fun test_get_all_configured_tokens_edge_cases() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    {
        let ctx = scenario.ctx();
        state_object::test_init(ctx);
    };

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        let ctx = scenario.ctx();
        registry::initialize(&mut ref, &owner_cap, ctx);

        // Test case 1: Empty state
        let (res, next_key, has_more) = registry::get_all_configured_tokens(&ref, @0x0, 1);
        assert!(res.length() == 0);
        assert!(next_key == @0x0);
        assert!(!has_more);

        // Test case 2: Single token
        registry::insert_token_addresses_for_test(&mut ref, vector[@0x1], TypeProof {});
        let (res, _next_key, has_more) = registry::get_all_configured_tokens(&ref, @0x0, 1);
        assert!(res.length() == 1);
        assert!(res[0] == @0x1);
        assert!(!has_more);

        // Test case 3: Start from middle
        registry::insert_token_addresses_for_test(&mut ref, vector[@0x2, @0x3], TypeProof {});
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

        scenario.return_to_sender(owner_cap);
        ts::return_shared(ref);
    };

    ts::end(scenario);
}

#[test]
public fun test_get_all_configured_tokens_pagination() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    {
        let ctx = scenario.ctx();
        state_object::test_init(ctx);
    };

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        let ctx = scenario.ctx();
        registry::initialize(&mut ref, &owner_cap, ctx);

        registry::insert_token_addresses_for_test(&mut ref, vector[@0x1, @0x2, @0x3, @0x4, @0x5], TypeProof {});

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

        scenario.return_to_sender(owner_cap);
        ts::return_shared(ref);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = registry::E_TOKEN_ADDRESS_NOT_REGISTERED)]
public fun test_get_all_configured_tokens_non_existent() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN_ADDRESS);
    {
        let ctx = scenario.ctx();
        state_object::test_init(ctx);
    };

    scenario.next_tx(TOKEN_ADMIN_ADDRESS);
    {
        let mut ref = scenario.take_shared<CCIPObjectRef>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        let ctx = scenario.ctx();
        registry::initialize(&mut ref, &owner_cap, ctx);

        registry::insert_token_addresses_for_test(&mut ref, vector[@0x1, @0x2, @0x3], TypeProof {});

        // Test starting from key between existing tokens
        let (res, _next_key, has_more) = registry::get_all_configured_tokens(&ref, @0x1, 1);
        assert!(res.length() == 1);
        assert!(res[0] == @0x2);
        assert!(has_more);

        // Test starting from non-existent key
        let (res, next_key, has_more) = registry::get_all_configured_tokens(&ref, @0x4, 1);
        assert!(res.length() == 0);
        assert!(next_key == @0x4);
        assert!(!has_more);

        scenario.return_to_sender(owner_cap);
        ts::return_shared(ref);
    };

    ts::end(scenario);
}
