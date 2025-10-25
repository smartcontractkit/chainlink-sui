#[test_only]
module burn_mint_token_pool::burn_mint_token_pool_tests;

use burn_mint_token_pool::burn_mint_token_pool::{Self, BurnMintTokenPoolState};
use burn_mint_token_pool::ownable::OwnerCap;
use ccip::offramp_state_helper as offramp_sh;
use ccip::onramp_state_helper as onramp_sh;
use ccip::rmn_remote;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::token_admin_registry;
use ccip::upgrade_registry;
use std::ascii;
use std::bcs;
use std::string;
use std::type_name;
use sui::clock;
use sui::coin;
use sui::package;
use sui::test_scenario;

public struct BURN_MINT_TOKEN_POOL_TESTS has drop {}

const Decimals: u8 = 8;
const DefaultRemoteChain: u64 = 2000;
const DefaultRemoteToken: vector<u8> = b"default_remote_token";
const DefaultRemotePool: vector<u8> = b"default_remote_pool";

const NewRemoteChain: u64 = 3000;
const NewRemotePool: vector<u8> = b"new_remote_pool";
const NewRemoteToken: vector<u8> = b"new_remote_token";

fun setup_ccip_environment(
    scenario: &mut test_scenario::Scenario,
): (ccip::ownable::OwnerCap, CCIPObjectRef) {
    scenario.next_tx(@burn_mint_token_pool);

    // Create CCIP state object
    state_object::test_init(scenario.ctx());

    // Advance to next transaction to retrieve the created objects
    scenario.next_tx(@burn_mint_token_pool);

    // Retrieve the OwnerCap that was transferred to the sender
    let ccip_owner_cap = scenario.take_from_sender<ccip::ownable::OwnerCap>();

    // Retrieve the shared CCIPObjectRef
    let mut ccip_ref = scenario.take_shared<CCIPObjectRef>();

    // Initialize required CCIP modules
    upgrade_registry::initialize(&mut ccip_ref, &ccip_owner_cap, scenario.ctx());
    token_admin_registry::initialize(&mut ccip_ref, &ccip_owner_cap, scenario.ctx());
    rmn_remote::initialize(&mut ccip_ref, &ccip_owner_cap, 1000, scenario.ctx());

    (ccip_owner_cap, ccip_ref)
}

#[test]
public fun test_type_and_version() {
    let version = burn_mint_token_pool::type_and_version();
    assert!(version == string::utf8(b"BurnMintTokenPool 1.6.0"));
}

#[test]
public fun test_initialize_and_basic_functionality() {
    let mut scenario = test_scenario::begin(@burn_mint_token_pool);

    // Setup CCIP environment
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let ctx = scenario.ctx();

        // Create token
        let (treasury_cap, coin_metadata) = coin::create_currency(
            BURN_MINT_TOKEN_POOL_TESTS {},
            Decimals,
            b"BMTP",
            b"BurnMintTestToken",
            b"burn_mint_test_token",
            option::none(),
            ctx,
        );

        // Initialize burn mint token pool
        burn_mint_token_pool::initialize_by_ccip_admin(
            &mut ccip_ref,
            state_object::create_ccip_admin_proof_for_test(),
            &coin_metadata,
            treasury_cap,
            @0x123, // token_pool_administrator
            package::test_claim(BURN_MINT_TOKEN_POOL_TESTS {}, ctx),
            ctx,
        );

        transfer::public_freeze_object(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);
    test_scenario::return_shared(ccip_ref);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let pool_state = scenario.take_shared<BurnMintTokenPoolState<BURN_MINT_TOKEN_POOL_TESTS>>();
        let owner_cap = scenario.take_from_sender<OwnerCap<BURN_MINT_TOKEN_POOL_TESTS>>();

        // Test basic getters
        let token_address = burn_mint_token_pool::get_token(&pool_state);
        assert!(token_address != @0x0); // Should be a valid address

        let decimals = burn_mint_token_pool::get_token_decimals(&pool_state);
        assert!(decimals == Decimals);

        let supported_chains = burn_mint_token_pool::get_supported_chains(&pool_state);
        assert!(supported_chains.length() == 0); // No chains configured initially

        let allowlist_enabled = burn_mint_token_pool::get_allowlist_enabled(&pool_state);
        assert!(!allowlist_enabled); // Should be disabled by default

        let allowlist = burn_mint_token_pool::get_allowlist(&pool_state);
        assert!(allowlist.length() == 0); // Should be empty initially

        scenario.return_to_sender(owner_cap);
        test_scenario::return_shared(pool_state);
    };

    test_scenario::end(scenario);
}

#[test]
public fun test_chain_configuration_management() {
    let mut scenario = test_scenario::begin(@burn_mint_token_pool);

    // Create CCIP object ref and initialize
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    // Create token and initialize pool
    scenario.next_tx(@burn_mint_token_pool);
    {
        let ctx = scenario.ctx();
        let (treasury_cap, coin_metadata) = coin::create_currency(
            BURN_MINT_TOKEN_POOL_TESTS {},
            Decimals,
            b"BMTP",
            b"BurnMintTestToken",
            b"burn_mint_test_token",
            option::none(),
            ctx,
        );

        burn_mint_token_pool::initialize_by_ccip_admin(
            &mut ccip_ref,
            state_object::create_ccip_admin_proof_for_test(),
            &coin_metadata,
            treasury_cap,
            @0x123,
            package::test_claim(BURN_MINT_TOKEN_POOL_TESTS {}, ctx),
            ctx,
        );

        transfer::public_freeze_object(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);
    test_scenario::return_shared(ccip_ref);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let mut pool_state = scenario.take_shared<
            BurnMintTokenPoolState<BURN_MINT_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap<BURN_MINT_TOKEN_POOL_TESTS>>();

        // Test apply_chain_updates - add chains
        burn_mint_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[], // remove nothing
            vector[DefaultRemoteChain, NewRemoteChain], // add two chains
            vector[vector[DefaultRemotePool], vector[NewRemotePool]], // pool addresses
            vector[DefaultRemoteToken, NewRemoteToken], // token addresses
        );

        // Verify chains were added
        let supported_chains = burn_mint_token_pool::get_supported_chains(&pool_state);
        assert!(supported_chains.length() == 2);
        assert!(burn_mint_token_pool::is_supported_chain(&pool_state, DefaultRemoteChain));
        assert!(burn_mint_token_pool::is_supported_chain(&pool_state, NewRemoteChain));

        // Test remote pool management
        let remote_pools = burn_mint_token_pool::get_remote_pools(&pool_state, DefaultRemoteChain);
        assert!(remote_pools.length() == 1);
        assert!(
            burn_mint_token_pool::is_remote_pool(
                &pool_state,
                DefaultRemoteChain,
                DefaultRemotePool,
            ),
        );

        // Add another remote pool
        burn_mint_token_pool::add_remote_pool(
            &mut pool_state,
            &owner_cap,
            DefaultRemoteChain,
            b"additional_pool",
        );

        let updated_pools = burn_mint_token_pool::get_remote_pools(&pool_state, DefaultRemoteChain);
        assert!(updated_pools.length() == 2);
        assert!(
            burn_mint_token_pool::is_remote_pool(
                &pool_state,
                DefaultRemoteChain,
                b"additional_pool",
            ),
        );

        // Remove a remote pool
        burn_mint_token_pool::remove_remote_pool(
            &mut pool_state,
            &owner_cap,
            DefaultRemoteChain,
            b"additional_pool",
        );

        let final_pools = burn_mint_token_pool::get_remote_pools(&pool_state, DefaultRemoteChain);
        assert!(final_pools.length() == 1);
        assert!(
            !burn_mint_token_pool::is_remote_pool(
                &pool_state,
                DefaultRemoteChain,
                b"additional_pool",
            ),
        );

        // Test get_remote_token
        let remote_token = burn_mint_token_pool::get_remote_token(&pool_state, DefaultRemoteChain);
        assert!(remote_token == DefaultRemoteToken);

        scenario.return_to_sender(owner_cap);
        test_scenario::return_shared(pool_state);
    };

    test_scenario::end(scenario);
}

#[test]
public fun test_allowlist_management() {
    let mut scenario = test_scenario::begin(@burn_mint_token_pool);

    // Setup and initialize
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let ctx = scenario.ctx();
        let (treasury_cap, coin_metadata) = coin::create_currency(
            BURN_MINT_TOKEN_POOL_TESTS {},
            Decimals,
            b"BMTP",
            b"BurnMintTestToken",
            b"burn_mint_test_token",
            option::none(),
            ctx,
        );

        burn_mint_token_pool::initialize_by_ccip_admin(
            &mut ccip_ref,
            state_object::create_ccip_admin_proof_for_test(),
            &coin_metadata,
            treasury_cap,
            @0x123,
            package::test_claim(BURN_MINT_TOKEN_POOL_TESTS {}, ctx),
            ctx,
        );

        transfer::public_freeze_object(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);
    test_scenario::return_shared(ccip_ref);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let pool_state = scenario.take_shared<BurnMintTokenPoolState<BURN_MINT_TOKEN_POOL_TESTS>>();
        let owner_cap = scenario.take_from_sender<OwnerCap<BURN_MINT_TOKEN_POOL_TESTS>>();

        // Initially allowlist should be disabled and empty
        let allowlist_enabled = burn_mint_token_pool::get_allowlist_enabled(&pool_state);
        assert!(!allowlist_enabled);

        let initial_allowlist = burn_mint_token_pool::get_allowlist(&pool_state);
        assert!(initial_allowlist.length() == 0);

        // Note: burn_mint_token_pool doesn't expose set_allowlist_enabled function
        // So we can only test the getter functions for allowlist
        // The allowlist functionality is managed at the token_pool level

        scenario.return_to_sender(owner_cap);
        test_scenario::return_shared(pool_state);
    };

    test_scenario::end(scenario);
}

#[test]
public fun test_rate_limiter_configuration() {
    let mut scenario = test_scenario::begin(@burn_mint_token_pool);

    // Setup and initialize
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let ctx = scenario.ctx();
        let (treasury_cap, coin_metadata) = coin::create_currency(
            BURN_MINT_TOKEN_POOL_TESTS {},
            Decimals,
            b"BMTP",
            b"BurnMintTestToken",
            b"burn_mint_test_token",
            option::none(),
            ctx,
        );

        burn_mint_token_pool::initialize_by_ccip_admin(
            &mut ccip_ref,
            state_object::create_ccip_admin_proof_for_test(),
            &coin_metadata,
            treasury_cap,
            @0x123,
            package::test_claim(BURN_MINT_TOKEN_POOL_TESTS {}, ctx),
            ctx,
        );

        transfer::public_freeze_object(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);
    test_scenario::return_shared(ccip_ref);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let mut pool_state = scenario.take_shared<
            BurnMintTokenPoolState<BURN_MINT_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap<BURN_MINT_TOKEN_POOL_TESTS>>();
        let mut ctx = tx_context::dummy();
        let clock = clock::create_for_testing(&mut ctx);

        // First add a chain to configure rate limiter for
        burn_mint_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[], // remove nothing
            vector[DefaultRemoteChain], // add one chain
            vector[vector[DefaultRemotePool]], // pool addresses
            vector[DefaultRemoteToken], // token addresses
        );

        // Test single chain rate limiter configuration
        burn_mint_token_pool::set_chain_rate_limiter_config(
            &mut pool_state,
            &owner_cap,
            &clock,
            DefaultRemoteChain,
            true, // outbound_is_enabled
            1000, // outbound_capacity
            100, // outbound_rate
            true, // inbound_is_enabled
            2000, // inbound_capacity
            200, // inbound_rate
        );

        // Add another chain for bulk configuration test
        burn_mint_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[], // remove nothing
            vector[NewRemoteChain], // add another chain
            vector[vector[NewRemotePool]], // pool addresses
            vector[NewRemoteToken], // token addresses
        );

        // Test bulk rate limiter configuration
        burn_mint_token_pool::set_chain_rate_limiter_configs(
            &mut pool_state,
            &owner_cap,
            &clock,
            vector[DefaultRemoteChain, NewRemoteChain], // remote_chain_selectors
            vector[true, false], // outbound_is_enableds
            vector[1500, 3000], // outbound_capacities
            vector[150, 300], // outbound_rates
            vector[false, true], // inbound_is_enableds
            vector[2500, 4000], // inbound_capacities
            vector[250, 400], // inbound_rates
        );

        clock.destroy_for_testing();
        transfer::public_transfer(owner_cap, @burn_mint_token_pool);
        test_scenario::return_shared(pool_state);
    };

    test_scenario::end(scenario);
}

#[test]
#[expected_failure(abort_code = burn_mint_token_pool::EInvalidArguments)]
public fun test_invalid_arguments_rate_limiter_configs() {
    let mut scenario = test_scenario::begin(@burn_mint_token_pool);
    let mut ctx = tx_context::dummy();
    let clock = clock::create_for_testing(&mut ctx);

    // Setup and initialize
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let ctx = scenario.ctx();
        let (treasury_cap, coin_metadata) = coin::create_currency(
            BURN_MINT_TOKEN_POOL_TESTS {},
            Decimals,
            b"BMTP",
            b"BurnMintTestToken",
            b"burn_mint_test_token",
            option::none(),
            ctx,
        );

        burn_mint_token_pool::initialize_by_ccip_admin(
            &mut ccip_ref,
            state_object::create_ccip_admin_proof_for_test(),
            &coin_metadata,
            treasury_cap,
            @0x123,
            package::test_claim(BURN_MINT_TOKEN_POOL_TESTS {}, ctx),
            ctx,
        );

        transfer::public_freeze_object(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);
    test_scenario::return_shared(ccip_ref);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let mut pool_state = scenario.take_shared<
            BurnMintTokenPoolState<BURN_MINT_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap<BURN_MINT_TOKEN_POOL_TESTS>>();

        // Test with mismatched vector lengths (should fail with EInvalidArguments)
        burn_mint_token_pool::set_chain_rate_limiter_configs(
            &mut pool_state,
            &owner_cap,
            &clock,
            vector[DefaultRemoteChain, NewRemoteChain], // 2 chains
            vector[true], // 1 outbound_is_enabled (mismatch!)
            vector[1000, 2000], // 2 outbound_capacities
            vector[100, 200], // 2 outbound_rates
            vector[true, false], // 2 inbound_is_enableds
            vector[1500, 2500], // 2 inbound_capacities
            vector[150, 250], // 2 inbound_rates
        );

        scenario.return_to_sender(owner_cap);
        test_scenario::return_shared(pool_state);
    };

    clock.destroy_for_testing();
    test_scenario::end(scenario);
}

#[test]
public fun test_comprehensive_allowlist_operations() {
    let mut scenario = test_scenario::begin(@burn_mint_token_pool);

    // Setup and initialize
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let ctx = scenario.ctx();
        let (treasury_cap, coin_metadata) = coin::create_currency(
            BURN_MINT_TOKEN_POOL_TESTS {},
            Decimals,
            b"BMTP",
            b"BurnMintTestToken",
            b"burn_mint_test_token",
            option::none(),
            ctx,
        );

        burn_mint_token_pool::initialize_by_ccip_admin(
            &mut ccip_ref,
            state_object::create_ccip_admin_proof_for_test(),
            &coin_metadata,
            treasury_cap,
            @0x123,
            package::test_claim(BURN_MINT_TOKEN_POOL_TESTS {}, ctx),
            ctx,
        );

        transfer::public_freeze_object(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);
    test_scenario::return_shared(ccip_ref);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let pool_state = scenario.take_shared<BurnMintTokenPoolState<BURN_MINT_TOKEN_POOL_TESTS>>();
        let owner_cap = scenario.take_from_sender<OwnerCap<BURN_MINT_TOKEN_POOL_TESTS>>();

        // Initial state: allowlist disabled and empty
        assert!(!burn_mint_token_pool::get_allowlist_enabled(&pool_state));
        let initial_allowlist = burn_mint_token_pool::get_allowlist(&pool_state);
        assert!(initial_allowlist.length() == 0);

        // Note: burn_mint_token_pool doesn't expose set_allowlist_enabled function,
        // so we can only test the getter functions. The allowlist updates require
        // the allowlist to be enabled first, which is not available in this module.

        // Test that allowlist getters work correctly
        assert!(!burn_mint_token_pool::get_allowlist_enabled(&pool_state));
        let allowlist = burn_mint_token_pool::get_allowlist(&pool_state);
        assert!(allowlist.length() == 0);

        scenario.return_to_sender(owner_cap);
        test_scenario::return_shared(pool_state);
    };

    test_scenario::end(scenario);
}

#[test]
public fun test_destroy_token_pool() {
    let mut scenario = test_scenario::begin(@burn_mint_token_pool);

    // Setup and initialize
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let ctx = scenario.ctx();
        let (treasury_cap, coin_metadata) = coin::create_currency(
            BURN_MINT_TOKEN_POOL_TESTS {},
            Decimals,
            b"BMTP",
            b"BurnMintTestToken",
            b"burn_mint_test_token",
            option::none(),
            ctx,
        );

        burn_mint_token_pool::initialize_by_ccip_admin(
            &mut ccip_ref,
            state_object::create_ccip_admin_proof_for_test(),
            &coin_metadata,
            treasury_cap,
            @burn_mint_token_pool,
            package::test_claim(BURN_MINT_TOKEN_POOL_TESTS {}, ctx),
            ctx,
        );

        transfer::public_freeze_object(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let pool_state = scenario.take_shared<BurnMintTokenPoolState<BURN_MINT_TOKEN_POOL_TESTS>>();
        let owner_cap = scenario.take_from_sender<OwnerCap<BURN_MINT_TOKEN_POOL_TESTS>>();

        token_admin_registry::unregister_pool(
            &mut ccip_ref,
            burn_mint_token_pool::get_token(&pool_state),
            scenario.ctx(),
        );

        // Test destroy_token_pool function
        let returned_treasury_cap = burn_mint_token_pool::destroy_token_pool(
            &mut ccip_ref,
            pool_state,
            owner_cap,
            scenario.ctx(),
        );

        // Verify we get back the treasury cap
        transfer::public_transfer(returned_treasury_cap, scenario.ctx().sender());
    };

    test_scenario::return_shared(ccip_ref);
    test_scenario::end(scenario);
}

#[test]
public fun test_comprehensive_rate_limiter_operations() {
    let mut scenario = test_scenario::begin(@burn_mint_token_pool);
    let mut ctx = tx_context::dummy();
    let clock = clock::create_for_testing(&mut ctx);

    // Setup and initialize
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let ctx = scenario.ctx();
        let (treasury_cap, coin_metadata) = coin::create_currency(
            BURN_MINT_TOKEN_POOL_TESTS {},
            Decimals,
            b"BMTP",
            b"BurnMintTestToken",
            b"burn_mint_test_token",
            option::none(),
            ctx,
        );

        burn_mint_token_pool::initialize_by_ccip_admin(
            &mut ccip_ref,
            state_object::create_ccip_admin_proof_for_test(),
            &coin_metadata,
            treasury_cap,
            @0x123,
            package::test_claim(BURN_MINT_TOKEN_POOL_TESTS {}, ctx),
            ctx,
        );

        transfer::public_freeze_object(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);
    test_scenario::return_shared(ccip_ref);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let mut pool_state = scenario.take_shared<
            BurnMintTokenPoolState<BURN_MINT_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap<BURN_MINT_TOKEN_POOL_TESTS>>();

        // First add chains to configure rate limits for
        burn_mint_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[],
            vector[DefaultRemoteChain, NewRemoteChain],
            vector[vector[DefaultRemotePool], vector[NewRemotePool]],
            vector[DefaultRemoteToken, NewRemoteToken],
        );

        // Test single chain rate limiter config
        burn_mint_token_pool::set_chain_rate_limiter_config(
            &mut pool_state,
            &owner_cap,
            &clock,
            DefaultRemoteChain,
            true, // outbound enabled
            1000, // outbound capacity
            100, // outbound rate
            true, // inbound enabled
            2000, // inbound capacity
            200, // inbound rate
        );

        // Test multiple chains rate limiter config with matching vector lengths
        burn_mint_token_pool::set_chain_rate_limiter_configs(
            &mut pool_state,
            &owner_cap,
            &clock,
            vector[DefaultRemoteChain, NewRemoteChain], // 2 chains
            vector[false, true], // 2 outbound_is_enableds
            vector[1500, 2500], // 2 outbound_capacities
            vector[150, 250], // 2 outbound_rates
            vector[true, false], // 2 inbound_is_enableds
            vector[3000, 4000], // 2 inbound_capacities
            vector[300, 400], // 2 inbound_rates
        );

        // Test with empty vectors (should work)
        burn_mint_token_pool::set_chain_rate_limiter_configs(
            &mut pool_state,
            &owner_cap,
            &clock,
            vector[], // empty chains
            vector[], // empty outbound_is_enableds
            vector[], // empty outbound_capacities
            vector[], // empty outbound_rates
            vector[], // empty inbound_is_enableds
            vector[], // empty inbound_capacities
            vector[], // empty inbound_rates
        );

        scenario.return_to_sender(owner_cap);
        test_scenario::return_shared(pool_state);
    };

    clock.destroy_for_testing();
    test_scenario::end(scenario);
}

#[test]
public fun test_edge_cases_and_boundary_conditions() {
    let mut scenario = test_scenario::begin(@burn_mint_token_pool);

    // Setup and initialize
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let ctx = scenario.ctx();
        let (treasury_cap, coin_metadata) = coin::create_currency(
            BURN_MINT_TOKEN_POOL_TESTS {},
            Decimals,
            b"BMTP",
            b"BurnMintTestToken",
            b"burn_mint_test_token",
            option::none(),
            ctx,
        );

        burn_mint_token_pool::initialize_by_ccip_admin(
            &mut ccip_ref,
            state_object::create_ccip_admin_proof_for_test(),
            &coin_metadata,
            treasury_cap,
            @0x123,
            package::test_claim(BURN_MINT_TOKEN_POOL_TESTS {}, ctx),
            ctx,
        );

        transfer::public_freeze_object(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);
    test_scenario::return_shared(ccip_ref);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let mut pool_state = scenario.take_shared<
            BurnMintTokenPoolState<BURN_MINT_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap<BURN_MINT_TOKEN_POOL_TESTS>>();

        // Test operations on empty chain configurations
        let supported_chains = burn_mint_token_pool::get_supported_chains(&pool_state);
        assert!(supported_chains.length() == 0);

        // Test adding and immediately removing the same chain
        burn_mint_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[],
            vector[DefaultRemoteChain],
            vector[vector[DefaultRemotePool]],
            vector[DefaultRemoteToken],
        );

        assert!(burn_mint_token_pool::is_supported_chain(&pool_state, DefaultRemoteChain));

        burn_mint_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[DefaultRemoteChain], // remove the chain we just added
            vector[],
            vector[],
            vector[],
        );

        assert!(!burn_mint_token_pool::is_supported_chain(&pool_state, DefaultRemoteChain));

        // Test adding multiple pools to the same chain
        burn_mint_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[],
            vector[DefaultRemoteChain],
            vector[vector[DefaultRemotePool, b"pool2", b"pool3", b"pool4"]], // 4 pools
            vector[DefaultRemoteToken],
        );

        let multiple_pools = burn_mint_token_pool::get_remote_pools(
            &pool_state,
            DefaultRemoteChain,
        );
        assert!(multiple_pools.length() == 4);
        assert!(
            burn_mint_token_pool::is_remote_pool(
                &pool_state,
                DefaultRemoteChain,
                DefaultRemotePool,
            ),
        );
        assert!(burn_mint_token_pool::is_remote_pool(&pool_state, DefaultRemoteChain, b"pool2"));
        assert!(burn_mint_token_pool::is_remote_pool(&pool_state, DefaultRemoteChain, b"pool3"));
        assert!(burn_mint_token_pool::is_remote_pool(&pool_state, DefaultRemoteChain, b"pool4"));

        // Test removing pools one by one
        burn_mint_token_pool::remove_remote_pool(
            &mut pool_state,
            &owner_cap,
            DefaultRemoteChain,
            b"pool2",
        );
        burn_mint_token_pool::remove_remote_pool(
            &mut pool_state,
            &owner_cap,
            DefaultRemoteChain,
            b"pool4",
        );

        let remaining_pools = burn_mint_token_pool::get_remote_pools(
            &pool_state,
            DefaultRemoteChain,
        );
        assert!(remaining_pools.length() == 2);
        assert!(
            burn_mint_token_pool::is_remote_pool(
                &pool_state,
                DefaultRemoteChain,
                DefaultRemotePool,
            ),
        );
        assert!(burn_mint_token_pool::is_remote_pool(&pool_state, DefaultRemoteChain, b"pool3"));
        assert!(!burn_mint_token_pool::is_remote_pool(&pool_state, DefaultRemoteChain, b"pool2"));
        assert!(!burn_mint_token_pool::is_remote_pool(&pool_state, DefaultRemoteChain, b"pool4"));

        scenario.return_to_sender(owner_cap);
        test_scenario::return_shared(pool_state);
    };

    test_scenario::end(scenario);
}

#[test]
public fun test_lock_or_burn_comprehensive() {
    let mut scenario = test_scenario::begin(@burn_mint_token_pool);

    // Setup CCIP environment
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);
    onramp_sh::test_init(scenario.ctx());

    // Create token and initialize pool
    scenario.next_tx(@burn_mint_token_pool);
    {
        let ctx = scenario.ctx();
        let (mut treasury_cap, coin_metadata) = coin::create_currency(
            BURN_MINT_TOKEN_POOL_TESTS {},
            Decimals,
            b"BMTP",
            b"BurnMintTestToken",
            b"burn_mint_test_token",
            option::none(),
            ctx,
        );

        // Mint some tokens for testing before initializing the pool
        let test_coin = coin::mint(&mut treasury_cap, 1000, ctx); // Small amount to stay within rate limiter
        transfer::public_transfer(test_coin, @0x456); // Transfer to test user

        burn_mint_token_pool::initialize_by_ccip_admin(
            &mut ccip_ref,
            state_object::create_ccip_admin_proof_for_test(),
            &coin_metadata,
            treasury_cap, // treasury_cap is moved here
            @0x123,
            package::test_claim(BURN_MINT_TOKEN_POOL_TESTS {}, ctx),
            ctx,
        );

        transfer::public_freeze_object(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);
    transfer::public_share_object(ccip_ref);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let mut pool_state = scenario.take_shared<
            BurnMintTokenPoolState<BURN_MINT_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap<BURN_MINT_TOKEN_POOL_TESTS>>();
        let ccip_ref = scenario.take_shared<CCIPObjectRef>();
        let source_transfer_cap = scenario.take_from_sender<onramp_sh::SourceTransferCap>();
        let mut ctx = tx_context::dummy();
        let mut clock = clock::create_for_testing(&mut ctx);

        // Configure chain and rate limiter
        burn_mint_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[],
            vector[DefaultRemoteChain],
            vector[vector[DefaultRemotePool]],
            vector[DefaultRemoteToken],
        );

        burn_mint_token_pool::set_chain_rate_limiter_config(
            &mut pool_state,
            &owner_cap,
            &clock,
            DefaultRemoteChain,
            true, // outbound enabled
            10000, // high capacity
            1000, // high rate
            true, // inbound enabled
            10000, // high capacity
            1000, // high rate
        );

        // Advance clock for rate limiter
        clock.increment_for_testing(1000000);

        // Clean up objects
        clock.destroy_for_testing();
        transfer::public_transfer(source_transfer_cap, @burn_mint_token_pool);
        transfer::public_transfer(owner_cap, @burn_mint_token_pool);
        test_scenario::return_shared(pool_state);
        test_scenario::return_shared(ccip_ref);
    };

    // Test lock_or_burn operation
    scenario.next_tx(@0x456); // Switch to test user
    {
        let mut pool_state = scenario.take_shared<
            BurnMintTokenPoolState<BURN_MINT_TOKEN_POOL_TESTS>,
        >();
        let ccip_ref = scenario.take_shared<CCIPObjectRef>();
        let test_coin = scenario.take_from_sender<coin::Coin<BURN_MINT_TOKEN_POOL_TESTS>>();
        let mut ctx = tx_context::dummy();
        let mut clock = clock::create_for_testing(&mut ctx);
        clock.increment_for_testing(1000000000); // Advance more to accumulate enough tokens

        let initial_coin_value = coin::value(&test_coin);
        assert!(initial_coin_value == 1000);

        let mut token_transfer_params = onramp_sh::create_token_transfer_params(
            bcs::to_bytes(&@0x456),
        ); // Use the test user address as token receiver

        // Perform lock_or_burn operation (this burns the coin)
        burn_mint_token_pool::lock_or_burn(
            &ccip_ref,
            &mut token_transfer_params,
            test_coin, // This coin gets burned
            DefaultRemoteChain,
            &clock,
            &mut pool_state,
            &mut ctx,
        );

        // Clean up token params
        let source_transfer_cap = scenario.take_from_address<onramp_sh::SourceTransferCap>(
            @burn_mint_token_pool,
        );

        let (
            remote_chain,
            token_pool_package_id,
            amount,
            source_token_address,
            dest_token_address,
            extra_data,
        ) = onramp_sh::get_source_token_transfer_data(&token_transfer_params);
        assert!(remote_chain == DefaultRemoteChain);

        assert!(amount == initial_coin_value);
        assert!(token_pool_package_id == @burn_mint_token_pool);
        assert!(source_token_address == burn_mint_token_pool::get_token(&pool_state));
        assert!(dest_token_address == DefaultRemoteToken);
        assert!(extra_data.length() > 0); // Should contain encoded decimals

        onramp_sh::deconstruct_token_params(&source_transfer_cap, token_transfer_params);

        clock.destroy_for_testing();
        transfer::public_transfer(source_transfer_cap, @burn_mint_token_pool);
        test_scenario::return_shared(pool_state);
        test_scenario::return_shared(ccip_ref);
    };

    test_scenario::end(scenario);
}

#[test]
public fun test_release_or_mint_comprehensive() {
    let mut scenario = test_scenario::begin(@burn_mint_token_pool);

    // Setup CCIP environment
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);
    offramp_sh::test_init(scenario.ctx());

    // Create token and initialize pool
    scenario.next_tx(@burn_mint_token_pool);
    let coin_metadata_address = {
        let ctx = scenario.ctx();
        let (treasury_cap, coin_metadata) = coin::create_currency(
            BURN_MINT_TOKEN_POOL_TESTS {},
            Decimals,
            b"BMTP",
            b"BurnMintTestToken",
            b"burn_mint_test_token",
            option::none(),
            ctx,
        );

        let coin_metadata_address = object::id_to_address(&object::id(&coin_metadata));

        burn_mint_token_pool::initialize_by_ccip_admin(
            &mut ccip_ref,
            state_object::create_ccip_admin_proof_for_test(),
            &coin_metadata,
            treasury_cap,
            @0x123,
            package::test_claim(BURN_MINT_TOKEN_POOL_TESTS {}, ctx),
            ctx,
        );

        transfer::public_freeze_object(coin_metadata);
        coin_metadata_address
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);
    transfer::public_share_object(ccip_ref);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let mut pool_state = scenario.take_shared<
            BurnMintTokenPoolState<BURN_MINT_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap<BURN_MINT_TOKEN_POOL_TESTS>>();
        let ccip_ref = scenario.take_shared<CCIPObjectRef>();
        let dest_transfer_cap = scenario.take_from_sender<offramp_sh::DestTransferCap>();
        let mut ctx = tx_context::dummy();
        let mut clock = clock::create_for_testing(&mut ctx);

        // Configure chain and rate limiter
        burn_mint_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[],
            vector[DefaultRemoteChain],
            vector[vector[DefaultRemotePool]],
            vector[DefaultRemoteToken],
        );

        burn_mint_token_pool::set_chain_rate_limiter_config(
            &mut pool_state,
            &owner_cap,
            &clock,
            DefaultRemoteChain,
            true, // outbound enabled
            10000, // high capacity
            1000, // high rate
            true, // inbound enabled
            10000, // high capacity
            1000, // high rate
        );

        // Advance clock for rate limiter
        clock.increment_for_testing(1000000);

        // Create receiver params for release_or_mint
        let mut receiver_params = offramp_sh::create_receiver_params(
            &dest_transfer_cap,
            DefaultRemoteChain,
        );

        // Add token transfer to receiver params
        let receiver_address = @0x789;
        let source_amount = 5000; // Reduced to stay within rate limiter capacity
        let source_pool_data = x"0000000000000000000000000000000000000000000000000000000000000008"; // 8 decimals encoded
        let offchain_data = vector[];

        offramp_sh::add_dest_token_transfer(
            &dest_transfer_cap,
            &mut receiver_params,
            receiver_address,
            DefaultRemoteChain, // remote_chain_selector
            source_amount,
            coin_metadata_address,
            @burn_mint_token_pool,
            DefaultRemotePool,
            source_pool_data,
            offchain_data,
        );

        // Verify the operation setup
        let source_chain = offramp_sh::get_source_chain_selector(&receiver_params);
        assert!(source_chain == DefaultRemoteChain);

        // Perform release_or_mint operation
        burn_mint_token_pool::release_or_mint(
            &ccip_ref,
            &mut receiver_params,
            &clock,
            &mut pool_state,
            &mut ctx,
        );

        // Clean up receiver params
        offramp_sh::deconstruct_receiver_params(&dest_transfer_cap, receiver_params);

        clock.destroy_for_testing();
        transfer::public_transfer(dest_transfer_cap, @burn_mint_token_pool);
        transfer::public_transfer(owner_cap, @burn_mint_token_pool);
        test_scenario::return_shared(pool_state);
        test_scenario::return_shared(ccip_ref);
    };

    // Verify the minted coin was transferred to the receiver
    scenario.next_tx(@0x789); // Switch to receiver
    {
        let minted_coin = scenario.take_from_sender<coin::Coin<BURN_MINT_TOKEN_POOL_TESTS>>();

        // Verify the minted amount (should be same as source amount due to same decimals)
        let minted_amount = coin::value(&minted_coin);
        assert!(minted_amount == 5000);

        // Transfer the coin back to clean up
        transfer::public_transfer(minted_coin, @burn_mint_token_pool);
    };

    test_scenario::end(scenario);
}

// === Additional Function Tests ===

// Note: The initialize function (with treasury cap) is tested implicitly
// in all other tests that use initialize_by_ccip_admin

#[test]
public fun test_set_allowlist_enabled() {
    let mut scenario = test_scenario::begin(@burn_mint_token_pool);

    // Setup and initialize
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let ctx = scenario.ctx();
        let (treasury_cap, coin_metadata) = coin::create_currency(
            BURN_MINT_TOKEN_POOL_TESTS {},
            Decimals,
            b"BMTP",
            b"BurnMintTestToken",
            b"burn_mint_test_token",
            option::none(),
            ctx,
        );

        burn_mint_token_pool::initialize_by_ccip_admin(
            &mut ccip_ref,
            state_object::create_ccip_admin_proof_for_test(),
            &coin_metadata,
            treasury_cap,
            @0x123,
            package::test_claim(BURN_MINT_TOKEN_POOL_TESTS {}, ctx),
            ctx,
        );

        transfer::public_freeze_object(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);
    test_scenario::return_shared(ccip_ref);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let mut pool_state = scenario.take_shared<
            BurnMintTokenPoolState<BURN_MINT_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap<BURN_MINT_TOKEN_POOL_TESTS>>();

        // Initially allowlist should be disabled
        assert!(!burn_mint_token_pool::get_allowlist_enabled(&pool_state));

        // Enable allowlist
        burn_mint_token_pool::set_allowlist_enabled(&mut pool_state, &owner_cap, true);
        assert!(burn_mint_token_pool::get_allowlist_enabled(&pool_state));

        // Disable allowlist
        burn_mint_token_pool::set_allowlist_enabled(&mut pool_state, &owner_cap, false);
        assert!(!burn_mint_token_pool::get_allowlist_enabled(&pool_state));

        // Enable again for next test
        burn_mint_token_pool::set_allowlist_enabled(&mut pool_state, &owner_cap, true);
        assert!(burn_mint_token_pool::get_allowlist_enabled(&pool_state));

        scenario.return_to_sender(owner_cap);
        test_scenario::return_shared(pool_state);
    };

    test_scenario::end(scenario);
}

#[test]
public fun test_apply_allowlist_updates() {
    let mut scenario = test_scenario::begin(@burn_mint_token_pool);

    // Setup and initialize
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let ctx = scenario.ctx();
        let (treasury_cap, coin_metadata) = coin::create_currency(
            BURN_MINT_TOKEN_POOL_TESTS {},
            Decimals,
            b"BMTP",
            b"BurnMintTestToken",
            b"burn_mint_test_token",
            option::none(),
            ctx,
        );

        burn_mint_token_pool::initialize_by_ccip_admin(
            &mut ccip_ref,
            state_object::create_ccip_admin_proof_for_test(),
            &coin_metadata,
            treasury_cap,
            @0x123,
            package::test_claim(BURN_MINT_TOKEN_POOL_TESTS {}, ctx),
            ctx,
        );

        transfer::public_freeze_object(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);
    test_scenario::return_shared(ccip_ref);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let mut pool_state = scenario.take_shared<
            BurnMintTokenPoolState<BURN_MINT_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap<BURN_MINT_TOKEN_POOL_TESTS>>();

        // Enable allowlist first
        burn_mint_token_pool::set_allowlist_enabled(&mut pool_state, &owner_cap, true);
        assert!(burn_mint_token_pool::get_allowlist_enabled(&pool_state));

        // Initially allowlist should be empty
        let initial_allowlist = burn_mint_token_pool::get_allowlist(&pool_state);
        assert!(initial_allowlist.length() == 0);

        // Add addresses to allowlist
        let addresses_to_add = vector[@0x111, @0x222, @0x333];
        burn_mint_token_pool::apply_allowlist_updates(
            &mut pool_state,
            &owner_cap,
            vector[], // no removes
            addresses_to_add,
        );

        // Verify addresses were added
        let updated_allowlist = burn_mint_token_pool::get_allowlist(&pool_state);
        assert!(updated_allowlist.length() == 3);
        assert!(updated_allowlist.contains(&@0x111));
        assert!(updated_allowlist.contains(&@0x222));
        assert!(updated_allowlist.contains(&@0x333));

        // Add more addresses and remove some
        burn_mint_token_pool::apply_allowlist_updates(
            &mut pool_state,
            &owner_cap,
            vector[@0x222], // remove @0x222
            vector[@0x444, @0x555], // add @0x444 and @0x555
        );

        // Verify final state
        let final_allowlist = burn_mint_token_pool::get_allowlist(&pool_state);
        assert!(final_allowlist.length() == 4); // 3 - 1 + 2 = 4
        assert!(final_allowlist.contains(&@0x111));
        assert!(!final_allowlist.contains(&@0x222)); // removed
        assert!(final_allowlist.contains(&@0x333));
        assert!(final_allowlist.contains(&@0x444)); // added
        assert!(final_allowlist.contains(&@0x555)); // added

        // Test removing all addresses
        burn_mint_token_pool::apply_allowlist_updates(
            &mut pool_state,
            &owner_cap,
            vector[@0x111, @0x333, @0x444, @0x555], // remove all
            vector[], // add none
        );

        let empty_allowlist = burn_mint_token_pool::get_allowlist(&pool_state);
        assert!(empty_allowlist.length() == 0);

        scenario.return_to_sender(owner_cap);
        test_scenario::return_shared(pool_state);
    };

    test_scenario::end(scenario);
}

#[test]
public fun test_set_pool() {
    let mut scenario = test_scenario::begin(@burn_mint_token_pool);

    // Setup CCIP environment
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    // Create token and initialize pool normally
    scenario.next_tx(@burn_mint_token_pool);
    let coin_metadata_address = {
        let ctx = scenario.ctx();
        let (treasury_cap, coin_metadata) = coin::create_currency(
            BURN_MINT_TOKEN_POOL_TESTS {},
            Decimals,
            b"BMTP",
            b"BurnMintTestToken",
            b"burn_mint_test_token",
            option::none(),
            ctx,
        );

        let coin_metadata_address = object::id_to_address(&object::id(&coin_metadata));

        // Initialize normally with burn_mint_token_pool
        burn_mint_token_pool::initialize_by_ccip_admin(
            &mut ccip_ref,
            state_object::create_ccip_admin_proof_for_test(),
            &coin_metadata,
            treasury_cap,
            @0x123,
            package::test_claim(BURN_MINT_TOKEN_POOL_TESTS {}, ctx),
            ctx,
        );

        transfer::public_freeze_object(coin_metadata);
        coin_metadata_address
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);
    test_scenario::return_shared(ccip_ref);

    // Verify initial pool registration with burn_mint_token_pool configuration
    scenario.next_tx(@burn_mint_token_pool);
    {
        let pool_state = scenario.take_shared<BurnMintTokenPoolState<BURN_MINT_TOKEN_POOL_TESTS>>();
        let owner_cap = scenario.take_from_sender<OwnerCap<BURN_MINT_TOKEN_POOL_TESTS>>();
        let ccip_ref = scenario.take_shared<CCIPObjectRef>();

        // Verify initial pool registration
        let initial_pool = token_admin_registry::get_pool(&ccip_ref, coin_metadata_address);
        assert!(initial_pool == @burn_mint_token_pool);

        // Get initial configuration
        let (
            initial_package_id,
            initial_module,
            _token_type,
            administrator,
            pending_admin,
            _initial_proof,
            initial_lock_params,
            initial_release_params,
        ) = token_admin_registry::get_token_config_data(&ccip_ref, coin_metadata_address);

        // Verify correct initial registration
        assert!(initial_package_id == @burn_mint_token_pool);
        assert!(initial_module == string::utf8(b"burn_mint_token_pool"));
        assert!(administrator == @0x123);
        assert!(pending_admin == @0x0);
        assert!(initial_lock_params.length() == 2); // CLOCK_ADDRESS and pool state address
        assert!(initial_release_params.length() == 2);

        // Transfer owner_cap to the administrator for the next transaction
        transfer::public_transfer(owner_cap, @0x123);
        test_scenario::return_shared(pool_state);
        test_scenario::return_shared(ccip_ref);
    };

    // Manually update the configuration to a different package ID
    // Must run as the administrator (@0x123) to unregister
    scenario.next_tx(@0x123);
    {
        let mut ccip_ref = scenario.take_shared<CCIPObjectRef>();

        // Use unregister and re-register to change the package ID
        token_admin_registry::unregister_pool(
            &mut ccip_ref,
            coin_metadata_address,
            scenario.ctx(),
        );

        // Register with a different package ID using CCIP admin
        let different_package_id = @0xcafe;
        let different_type_proof = ascii::string(b"0xcafe::different_pool::DifferentTypeProof");
        let different_params = vector[@0x6, @0xfade];

        token_admin_registry::register_pool_by_admin(
            &mut ccip_ref,
            state_object::create_ccip_admin_proof_for_test(),
            coin_metadata_address,
            different_package_id,
            string::utf8(b"different_pool"),
            type_name::into_string(type_name::with_defining_ids<BURN_MINT_TOKEN_POOL_TESTS>()),
            @0x123, // administrator
            different_type_proof,
            different_params,
            different_params,
            scenario.ctx(),
        );

        test_scenario::return_shared(ccip_ref);
    };

    // Verify the different configuration
    scenario.next_tx(@0x123);
    {
        let ccip_ref = scenario.take_shared<CCIPObjectRef>();

        let (
            before_package_id,
            before_module,
            before_token_type,
            _before_administrator,
            _before_pending_admin,
            before_proof,
            before_lock_params,
            before_release_params,
        ) = token_admin_registry::get_token_config_data(&ccip_ref, coin_metadata_address);

        // Verify it's different from burn_mint_token_pool
        assert!(before_package_id == @0xcafe);
        assert!(before_module == string::utf8(b"different_pool"));
        assert!(
            before_token_type == type_name::into_string(type_name::with_defining_ids<BURN_MINT_TOKEN_POOL_TESTS>()),
        );
        assert!(before_proof == ascii::string(b"0xcafe::different_pool::DifferentTypeProof"));
        assert!(before_lock_params == vector[@0x6, @0xfade]);
        assert!(before_release_params == vector[@0x6, @0xfade]);

        test_scenario::return_shared(ccip_ref);
    };

    // Now call set_pool as the administrator to update to the correct burn_mint_token_pool config
    scenario.next_tx(@0x123);
    {
        let mut pool_state = scenario.take_shared<
            BurnMintTokenPoolState<BURN_MINT_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap<BURN_MINT_TOKEN_POOL_TESTS>>();
        let mut ccip_ref = scenario.take_shared<CCIPObjectRef>();

        // Get configuration before set_pool
        let (
            before_package_id,
            before_module,
            before_token_type,
            _before_admin,
            _before_pending,
            before_proof,
            _before_lock,
            _before_release,
        ) = token_admin_registry::get_token_config_data(&ccip_ref, coin_metadata_address);

        assert!(before_package_id == @0xcafe);
        assert!(before_module == string::utf8(b"different_pool"));
        assert!(
            before_token_type == type_name::into_string(type_name::with_defining_ids<BURN_MINT_TOKEN_POOL_TESTS>()),
        );

        // Call set_pool to update to the actual burn_mint_token_pool configuration
        burn_mint_token_pool::set_pool(
            &mut ccip_ref,
            &mut pool_state,
            &owner_cap,
            coin_metadata_address,
            scenario.ctx(),
        );

        // Verify the pool configuration CHANGED after set_pool
        let updated_pool = token_admin_registry::get_pool(&ccip_ref, coin_metadata_address);
        assert!(updated_pool == @burn_mint_token_pool);

        let (
            after_package_id,
            after_module,
            after_token_type,
            after_administrator,
            after_pending_admin,
            after_proof,
            after_lock_params,
            after_release_params,
        ) = token_admin_registry::get_token_config_data(&ccip_ref, coin_metadata_address);

        // Verify the configuration changed to burn_mint_token_pool
        assert!(after_package_id == @burn_mint_token_pool);
        assert!(after_module == string::utf8(b"burn_mint_token_pool"));
        assert!(
            after_token_type == type_name::into_string(type_name::with_defining_ids<BURN_MINT_TOKEN_POOL_TESTS>()),
        );
        assert!(after_administrator == @0x123);
        assert!(after_pending_admin == @0x0);
        assert!(after_lock_params.length() == 2);
        assert!(after_release_params.length() == 2);

        // Verify package ID, module, and type proof actually changed
        assert!(before_package_id != after_package_id);
        assert!(before_module != after_module);
        assert!(before_proof != after_proof);

        // Note: token type remains the same (BURN_MINT_TOKEN_POOL_TESTS) since we're testing the same token
        // but with different pool implementations
        assert!(before_token_type == after_token_type);

        scenario.return_to_sender(owner_cap);
        test_scenario::return_shared(pool_state);
        test_scenario::return_shared(ccip_ref);
    };

    test_scenario::end(scenario);
}

#[test]
public fun test_allowlist_enabled_and_updates_comprehensive() {
    let mut scenario = test_scenario::begin(@burn_mint_token_pool);

    // Setup and initialize
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let ctx = scenario.ctx();
        let (treasury_cap, coin_metadata) = coin::create_currency(
            BURN_MINT_TOKEN_POOL_TESTS {},
            Decimals,
            b"BMTP",
            b"BurnMintTestToken",
            b"burn_mint_test_token",
            option::none(),
            ctx,
        );

        burn_mint_token_pool::initialize_by_ccip_admin(
            &mut ccip_ref,
            state_object::create_ccip_admin_proof_for_test(),
            &coin_metadata,
            treasury_cap,
            @0x123,
            package::test_claim(BURN_MINT_TOKEN_POOL_TESTS {}, ctx),
            ctx,
        );

        transfer::public_freeze_object(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);
    test_scenario::return_shared(ccip_ref);

    scenario.next_tx(@burn_mint_token_pool);
    {
        let mut pool_state = scenario.take_shared<
            BurnMintTokenPoolState<BURN_MINT_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap<BURN_MINT_TOKEN_POOL_TESTS>>();

        // Test initial state
        assert!(!burn_mint_token_pool::get_allowlist_enabled(&pool_state));
        assert!(burn_mint_token_pool::get_allowlist(&pool_state).length() == 0);

        // Test enabling allowlist and adding addresses in sequence
        burn_mint_token_pool::set_allowlist_enabled(&mut pool_state, &owner_cap, true);

        // Test multiple update operations
        burn_mint_token_pool::apply_allowlist_updates(
            &mut pool_state,
            &owner_cap,
            vector[],
            vector[@0xaaa, @0xbbb],
        );

        assert!(burn_mint_token_pool::get_allowlist(&pool_state).length() == 2);

        burn_mint_token_pool::apply_allowlist_updates(
            &mut pool_state,
            &owner_cap,
            vector[@0xaaa],
            vector[@0xccc, @0xddd, @0xeee],
        );

        let current_allowlist = burn_mint_token_pool::get_allowlist(&pool_state);
        assert!(current_allowlist.length() == 4); // 2 - 1 + 3 = 4
        assert!(!current_allowlist.contains(&@0xaaa));
        assert!(current_allowlist.contains(&@0xbbb));
        assert!(current_allowlist.contains(&@0xccc));
        assert!(current_allowlist.contains(&@0xddd));
        assert!(current_allowlist.contains(&@0xeee));

        // Test disabling allowlist (allowlist data should remain but be disabled)
        burn_mint_token_pool::set_allowlist_enabled(&mut pool_state, &owner_cap, false);
        assert!(!burn_mint_token_pool::get_allowlist_enabled(&pool_state));
        assert!(burn_mint_token_pool::get_allowlist(&pool_state).length() == 4); // Data preserved

        // Test re-enabling
        burn_mint_token_pool::set_allowlist_enabled(&mut pool_state, &owner_cap, true);
        assert!(burn_mint_token_pool::get_allowlist_enabled(&pool_state));
        assert!(burn_mint_token_pool::get_allowlist(&pool_state).length() == 4); // Data preserved

        scenario.return_to_sender(owner_cap);
        test_scenario::return_shared(pool_state);
    };

    test_scenario::end(scenario);
}

// ================================================================
// |             Rate Limiter calculate_refill Tests             |
// ================================================================

#[test]
public fun test_calculate_refill_at_capacity() {
    use burn_mint_token_pool::rate_limiter;

    // Create a bucket that's already at capacity
    let bucket = rate_limiter::test_create_token_bucket(
        1000, // tokens (at capacity)
        0, // last_updated
        true, // is_enabled
        1000, // capacity
        10, // rate
    );

    // Test with various time differences
    assert!(rate_limiter::test_calculate_refill(&bucket, 0) == 1000);
    assert!(rate_limiter::test_calculate_refill(&bucket, 10) == 1000);
    assert!(rate_limiter::test_calculate_refill(&bucket, 100) == 1000);
}

#[test]
public fun test_calculate_refill_above_capacity() {
    use burn_mint_token_pool::rate_limiter;

    // Create a bucket with tokens exceeding capacity (edge case)
    let bucket = rate_limiter::test_create_token_bucket(
        1500, // tokens (above capacity)
        0, // last_updated
        true, // is_enabled
        1000, // capacity
        10, // rate
    );

    // Should return capacity when tokens > capacity
    assert!(rate_limiter::test_calculate_refill(&bucket, 10) == 1000);
}

#[test]
public fun test_calculate_refill_zero_rate() {
    use burn_mint_token_pool::rate_limiter;

    // Create a bucket with rate = 0
    let bucket = rate_limiter::test_create_token_bucket(
        500, // tokens
        0, // last_updated
        true, // is_enabled
        1000, // capacity
        0, // rate (zero)
    );

    // Should return current tokens when rate is 0
    assert!(rate_limiter::test_calculate_refill(&bucket, 0) == 500);
    assert!(rate_limiter::test_calculate_refill(&bucket, 10) == 500);
    assert!(rate_limiter::test_calculate_refill(&bucket, 100) == 500);
}

#[test]
public fun test_calculate_refill_normal_refill() {
    use burn_mint_token_pool::rate_limiter;

    // Create a bucket with normal values
    let bucket = rate_limiter::test_create_token_bucket(
        100, // tokens
        0, // last_updated
        true, // is_enabled
        1000, // capacity
        10, // rate (10 tokens per second)
    );

    // Test various time differences
    // After 10 seconds: 100 + (10 * 10) = 200
    assert!(rate_limiter::test_calculate_refill(&bucket, 10) == 200);

    // After 50 seconds: 100 + (50 * 10) = 600
    assert!(rate_limiter::test_calculate_refill(&bucket, 50) == 600);

    // After 89 seconds: 100 + (89 * 10) = 990
    assert!(rate_limiter::test_calculate_refill(&bucket, 89) == 990);
}

#[test]
public fun test_calculate_refill_saturate_at_capacity() {
    use burn_mint_token_pool::rate_limiter;

    // Create a bucket that will exceed capacity with refill
    let bucket = rate_limiter::test_create_token_bucket(
        100, // tokens
        0, // last_updated
        true, // is_enabled
        1000, // capacity
        10, // rate (10 tokens per second)
    );

    // Time diff that would exceed capacity: 100 + (100 * 10) = 1100 > 1000
    // Should saturate to capacity (1000)
    assert!(rate_limiter::test_calculate_refill(&bucket, 100) == 1000);

    // Even larger time diff should still saturate to capacity
    assert!(rate_limiter::test_calculate_refill(&bucket, 1000) == 1000);
}

#[test]
public fun test_calculate_refill_exact_capacity() {
    use burn_mint_token_pool::rate_limiter;

    // Create a bucket where refill will exactly reach capacity
    let bucket = rate_limiter::test_create_token_bucket(
        100, // tokens
        0, // last_updated
        true, // is_enabled
        1000, // capacity
        10, // rate (10 tokens per second)
    );

    // Exactly 90 seconds needed to reach capacity: 100 + (90 * 10) = 1000
    assert!(rate_limiter::test_calculate_refill(&bucket, 90) == 1000);
}

#[test]
public fun test_calculate_refill_zero_tokens() {
    use burn_mint_token_pool::rate_limiter;

    // Create a bucket starting with zero tokens
    let bucket = rate_limiter::test_create_token_bucket(
        0, // tokens (empty)
        0, // last_updated
        true, // is_enabled
        1000, // capacity
        10, // rate
    );

    // Test refill from empty
    assert!(rate_limiter::test_calculate_refill(&bucket, 0) == 0);
    assert!(rate_limiter::test_calculate_refill(&bucket, 10) == 100);
    assert!(rate_limiter::test_calculate_refill(&bucket, 50) == 500);
    assert!(rate_limiter::test_calculate_refill(&bucket, 100) == 1000);
    assert!(rate_limiter::test_calculate_refill(&bucket, 200) == 1000); // saturate
}

#[test]
public fun test_calculate_refill_high_rate() {
    use burn_mint_token_pool::rate_limiter;

    // Create a bucket with high rate
    let bucket = rate_limiter::test_create_token_bucket(
        0, // tokens
        0, // last_updated
        true, // is_enabled
        10000, // capacity
        1000, // rate (high rate)
    );

    // After 5 seconds: 0 + (5 * 1000) = 5000
    assert!(rate_limiter::test_calculate_refill(&bucket, 5) == 5000);

    // After 10 seconds: should saturate to capacity
    assert!(rate_limiter::test_calculate_refill(&bucket, 10) == 10000);
}

#[test]
public fun test_calculate_refill_small_capacity() {
    use burn_mint_token_pool::rate_limiter;

    // Create a bucket with small capacity
    let bucket = rate_limiter::test_create_token_bucket(
        5, // tokens
        0, // last_updated
        true, // is_enabled
        10, // capacity (small)
        1, // rate
    );

    // After 3 seconds: 5 + (3 * 1) = 8
    assert!(rate_limiter::test_calculate_refill(&bucket, 3) == 8);

    // After 5 seconds: should be exactly at capacity
    assert!(rate_limiter::test_calculate_refill(&bucket, 5) == 10);

    // After 10 seconds: should saturate to capacity
    assert!(rate_limiter::test_calculate_refill(&bucket, 10) == 10);
}

#[test]
public fun test_calculate_refill_edge_case_one_below_capacity() {
    use burn_mint_token_pool::rate_limiter;

    // Create a bucket one token below capacity
    let bucket = rate_limiter::test_create_token_bucket(
        999, // tokens (one below capacity)
        0, // last_updated
        true, // is_enabled
        1000, // capacity
        10, // rate
    );

    // After 0 seconds: stays at 999
    assert!(rate_limiter::test_calculate_refill(&bucket, 0) == 999);

    // After 1 second: 999 + (1 * 10) = 1009, but should saturate to 1000
    assert!(rate_limiter::test_calculate_refill(&bucket, 1) == 1000);
}
