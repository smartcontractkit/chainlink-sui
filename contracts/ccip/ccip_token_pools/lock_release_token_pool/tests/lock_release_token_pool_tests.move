#[test_only]
module lock_release_token_pool::lock_release_token_pool_tests;

use ccip::offramp_state_helper as offramp_sh;
use ccip::onramp_state_helper as onramp_sh;
use ccip::ownable::OwnerCap as CCIPOwnerCap;
use ccip::rmn_remote;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::token_admin_registry;
use ccip::upgrade_registry;
use lock_release_token_pool::lock_release_token_pool::{
    Self,
    LockReleaseTokenPoolState,
    RebalancerCap
};
use lock_release_token_pool::ownable::OwnerCap;
use std::bcs;
use std::string;
use std::type_name;
use sui::address;
use sui::clock;
use sui::coin;
use sui::test_scenario::{Self, Scenario};

public struct LOCK_RELEASE_TOKEN_POOL_TESTS has drop {}

const Decimals: u8 = 8;
const DefaultRemoteChain: u64 = 2000;
const DefaultRemotePool: vector<u8> = b"default_remote_pool";
const DefaultRemoteToken: vector<u8> = b"default_remote_token";

const REBALANCER: address = @0x100;
const TOKEN_ADMIN: address = @0x200;
const CCIP_ADMIN: address = @0x400;

fun create_test_scenario(addr: address): Scenario {
    test_scenario::begin(addr)
}

fun setup_ccip_environment(scenario: &mut Scenario): (CCIPOwnerCap, CCIPObjectRef) {
    scenario.next_tx(CCIP_ADMIN);
    let ctx = scenario.ctx();

    // Create CCIP state object
    state_object::test_init(ctx);

    // Advance to next transaction to retrieve the created objects
    scenario.next_tx(CCIP_ADMIN);

    // Retrieve the OwnerCap that was transferred to the sender
    let ccip_owner_cap = scenario.take_from_sender<CCIPOwnerCap>();

    // Retrieve the shared CCIPObjectRef
    let mut ccip_ref = scenario.take_shared<CCIPObjectRef>();

    // Initialize required CCIP modules
    upgrade_registry::initialize(&mut ccip_ref, &ccip_owner_cap, scenario.ctx());
    rmn_remote::initialize(&mut ccip_ref, &ccip_owner_cap, 1000, scenario.ctx()); // local chain selector = 1000
    token_admin_registry::initialize(&mut ccip_ref, &ccip_owner_cap, scenario.ctx());
    onramp_sh::test_init(scenario.ctx());
    offramp_sh::test_init(scenario.ctx());

    (ccip_owner_cap, ccip_ref)
}

#[test]
public fun test_type_and_version() {
    let version = lock_release_token_pool::type_and_version();
    assert!(version == string::utf8(b"LockReleaseTokenPool 1.6.0"));
}

#[test]
public fun test_initialize_and_basic_functionality() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN);

    // Setup CCIP environment
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(TOKEN_ADMIN);
    {
        let ctx = scenario.ctx();

        // Create test token
        let (treasury_cap, coin_metadata) = coin::create_currency(
            LOCK_RELEASE_TOKEN_POOL_TESTS {},
            Decimals,
            b"TEST",
            b"TestToken",
            b"test_token",
            option::none(),
            ctx,
        );

        // Call test_init to create owner_cap
        lock_release_token_pool::test_init(ctx);

        transfer::public_freeze_object(coin_metadata);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
    };

    scenario.next_tx(TOKEN_ADMIN);
    {
        let mut owner_cap = scenario.take_from_sender<OwnerCap>();
        let coin_metadata = scenario.take_immutable<
            coin::CoinMetadata<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let treasury_cap = scenario.take_from_sender<
            coin::TreasuryCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();

        // Initialize the lock release token pool
        lock_release_token_pool::initialize(
            &mut owner_cap,
            &mut ccip_ref,
            &coin_metadata,
            &treasury_cap,
            TOKEN_ADMIN,
            REBALANCER,
            scenario.ctx(),
        );

        transfer::public_transfer(owner_cap, TOKEN_ADMIN);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
        test_scenario::return_immutable(coin_metadata);
    };

    scenario.next_tx(TOKEN_ADMIN);
    {
        let pool_state = scenario.take_shared<
            LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        // Test basic getters
        assert!(lock_release_token_pool::get_token_decimals(&pool_state) == Decimals);
        // Rebalancer is represented by the RebalancerCap holder; ensure non-zero address
        assert!(lock_release_token_pool::get_rebalancer(&pool_state) != @0x0);
        assert!(
            lock_release_token_pool::get_balance<LOCK_RELEASE_TOKEN_POOL_TESTS>(&pool_state) == 0,
        );

        // Test supported chains (should be empty initially)
        let supported_chains = lock_release_token_pool::get_supported_chains(&pool_state);
        assert!(supported_chains.length() == 0);

        scenario.return_to_sender(owner_cap);
        test_scenario::return_shared(pool_state);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);
    test_scenario::return_shared(ccip_ref);
    test_scenario::end(scenario);
}

#[test]
public fun test_chain_configuration_management() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN);

    // Setup CCIP environment
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(TOKEN_ADMIN);
    {
        let ctx = scenario.ctx();

        // Create test token
        let (treasury_cap, coin_metadata) = coin::create_currency(
            LOCK_RELEASE_TOKEN_POOL_TESTS {},
            Decimals,
            b"TEST",
            b"TestToken",
            b"test_token",
            option::none(),
            ctx,
        );

        // Call test_init to create owner_cap
        lock_release_token_pool::test_init(ctx);

        transfer::public_freeze_object(coin_metadata);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
    };

    scenario.next_tx(TOKEN_ADMIN);
    {
        let mut owner_cap = scenario.take_from_sender<OwnerCap>();
        let coin_metadata = scenario.take_immutable<
            coin::CoinMetadata<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let treasury_cap = scenario.take_from_sender<
            coin::TreasuryCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();

        // Initialize pool
        lock_release_token_pool::initialize(
            &mut owner_cap,
            &mut ccip_ref,
            &coin_metadata,
            &treasury_cap,
            TOKEN_ADMIN,
            REBALANCER,
            scenario.ctx(),
        );

        transfer::public_transfer(owner_cap, TOKEN_ADMIN);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
        test_scenario::return_immutable(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);

    scenario.next_tx(TOKEN_ADMIN);
    {
        let mut pool_state = scenario.take_shared<
            LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        // Test chain updates
        lock_release_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[], // remove none
            vector[DefaultRemoteChain], // add one chain
            vector[vector[DefaultRemotePool]], // with one pool
            vector[DefaultRemoteToken], // and token address
        );

        // Verify chain was added
        assert!(lock_release_token_pool::is_supported_chain(&pool_state, DefaultRemoteChain));
        let supported_chains = lock_release_token_pool::get_supported_chains(&pool_state);
        assert!(supported_chains.length() == 1);
        assert!(supported_chains[0] == DefaultRemoteChain);

        // Test remote pool management
        let new_remote_pool = b"new_remote_pool";
        lock_release_token_pool::add_remote_pool(
            &mut pool_state,
            &owner_cap,
            DefaultRemoteChain,
            new_remote_pool,
        );

        assert!(
            lock_release_token_pool::is_remote_pool(
                &pool_state,
                DefaultRemoteChain,
                new_remote_pool,
            ),
        );
        let remote_pools = lock_release_token_pool::get_remote_pools(
            &pool_state,
            DefaultRemoteChain,
        );
        assert!(remote_pools.length() == 2); // original + new

        // Test remote pool removal
        lock_release_token_pool::remove_remote_pool(
            &mut pool_state,
            &owner_cap,
            DefaultRemoteChain,
            new_remote_pool,
        );

        assert!(
            !lock_release_token_pool::is_remote_pool(
                &pool_state,
                DefaultRemoteChain,
                new_remote_pool,
            ),
        );
        let remote_pools_after = lock_release_token_pool::get_remote_pools(
            &pool_state,
            DefaultRemoteChain,
        );
        assert!(remote_pools_after.length() == 1);

        scenario.return_to_sender(owner_cap);
        test_scenario::return_shared(pool_state);
    };

    test_scenario::return_shared(ccip_ref);
    test_scenario::end(scenario);
}

#[test]
public fun test_liquidity_management() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN);

    // Setup CCIP environment
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(TOKEN_ADMIN);
    {
        let ctx = scenario.ctx();

        // Create test token
        let (mut treasury_cap, coin_metadata) = coin::create_currency(
            LOCK_RELEASE_TOKEN_POOL_TESTS {},
            Decimals,
            b"TEST",
            b"TestToken",
            b"test_token",
            option::none(),
            ctx,
        );

        // Call test_init to create owner_cap
        lock_release_token_pool::test_init(ctx);

        // Mint some tokens for testing
        let test_coin = coin::mint(&mut treasury_cap, 1000000, ctx);
        transfer::public_transfer(test_coin, REBALANCER);

        transfer::public_freeze_object(coin_metadata);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
    };

    scenario.next_tx(TOKEN_ADMIN);
    {
        let mut owner_cap = scenario.take_from_sender<OwnerCap>();
        let coin_metadata = scenario.take_immutable<
            coin::CoinMetadata<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let treasury_cap = scenario.take_from_sender<
            coin::TreasuryCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();

        // Initialize pool
        lock_release_token_pool::initialize(
            &mut owner_cap,
            &mut ccip_ref,
            &coin_metadata,
            &treasury_cap,
            TOKEN_ADMIN,
            REBALANCER,
            scenario.ctx(),
        );

        transfer::public_transfer(owner_cap, TOKEN_ADMIN);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
        test_scenario::return_immutable(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);

    // Test provide liquidity
    scenario.next_tx(REBALANCER);
    {
        let mut pool_state = scenario.take_shared<
            LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let liquidity_coin = scenario.take_from_sender<coin::Coin<LOCK_RELEASE_TOKEN_POOL_TESTS>>();
        let rebalancer_cap = scenario.take_from_sender<
            RebalancerCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();

        let initial_balance = lock_release_token_pool::get_balance<LOCK_RELEASE_TOKEN_POOL_TESTS>(
            &pool_state,
        );
        let liquidity_amount = coin::value(&liquidity_coin);

        // Provide liquidity
        lock_release_token_pool::provide_liquidity(
            &mut pool_state,
            &rebalancer_cap,
            liquidity_coin,
            scenario.ctx(),
        );

        // Verify balance increased
        let new_balance = lock_release_token_pool::get_balance<LOCK_RELEASE_TOKEN_POOL_TESTS>(
            &pool_state,
        );
        assert!(new_balance == initial_balance + liquidity_amount);

        transfer::public_transfer(rebalancer_cap, REBALANCER);
        test_scenario::return_shared(pool_state);
    };

    // Test withdraw liquidity
    scenario.next_tx(REBALANCER);
    {
        let mut pool_state = scenario.take_shared<
            LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let rebalancer_cap = scenario.take_from_sender<
            RebalancerCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();

        let withdraw_amount = 500000;
        let initial_balance = lock_release_token_pool::get_balance<LOCK_RELEASE_TOKEN_POOL_TESTS>(
            &pool_state,
        );

        // Withdraw liquidity
        let withdrawn_coin = lock_release_token_pool::withdraw_liquidity<
            LOCK_RELEASE_TOKEN_POOL_TESTS,
        >(
            &mut pool_state,
            &rebalancer_cap,
            withdraw_amount,
            scenario.ctx(),
        );

        // Verify withdrawal
        assert!(coin::value(&withdrawn_coin) == withdraw_amount);
        let new_balance = lock_release_token_pool::get_balance<LOCK_RELEASE_TOKEN_POOL_TESTS>(
            &pool_state,
        );
        assert!(new_balance == initial_balance - withdraw_amount);

        transfer::public_transfer(withdrawn_coin, scenario.ctx().sender());
        transfer::public_transfer(rebalancer_cap, REBALANCER);
        test_scenario::return_shared(pool_state);
    };

    test_scenario::return_shared(ccip_ref);
    test_scenario::end(scenario);
}

#[test]
public fun test_rebalancer_management() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN);

    // Setup CCIP environment
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(TOKEN_ADMIN);
    {
        let ctx = scenario.ctx();

        // Create test token
        let (treasury_cap, coin_metadata) = coin::create_currency(
            LOCK_RELEASE_TOKEN_POOL_TESTS {},
            Decimals,
            b"TEST",
            b"TestToken",
            b"test_token",
            option::none(),
            ctx,
        );

        // Call test_init to create owner_cap
        lock_release_token_pool::test_init(ctx);

        transfer::public_freeze_object(coin_metadata);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
    };

    scenario.next_tx(TOKEN_ADMIN);
    {
        let mut owner_cap = scenario.take_from_sender<OwnerCap>();
        let coin_metadata = scenario.take_immutable<
            coin::CoinMetadata<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let treasury_cap = scenario.take_from_sender<
            coin::TreasuryCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();

        // Initialize pool
        lock_release_token_pool::initialize(
            &mut owner_cap,
            &mut ccip_ref,
            &coin_metadata,
            &treasury_cap,
            TOKEN_ADMIN,
            REBALANCER,
            scenario.ctx(),
        );

        transfer::public_transfer(owner_cap, TOKEN_ADMIN);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
        test_scenario::return_immutable(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);

    // Verify that the rebalancer cap was transferred to the rebalancer address
    scenario.next_tx(REBALANCER);
    {
        let rebalancer_cap = scenario.take_from_sender<
            RebalancerCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        // Return the cap to the rebalancer to keep state consistent
        transfer::public_transfer(rebalancer_cap, REBALANCER);
    };
    test_scenario::return_shared(ccip_ref);
    test_scenario::end(scenario);
}

#[test]
public fun test_rate_limiter_configuration() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN);

    // Setup CCIP environment
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(TOKEN_ADMIN);
    {
        let ctx = scenario.ctx();

        // Create test token
        let (treasury_cap, coin_metadata) = coin::create_currency(
            LOCK_RELEASE_TOKEN_POOL_TESTS {},
            Decimals,
            b"TEST",
            b"TestToken",
            b"test_token",
            option::none(),
            ctx,
        );

        // Call test_init to create owner_cap
        lock_release_token_pool::test_init(ctx);

        transfer::public_freeze_object(coin_metadata);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
    };

    scenario.next_tx(TOKEN_ADMIN);
    {
        let mut owner_cap = scenario.take_from_sender<OwnerCap>();
        let coin_metadata = scenario.take_immutable<
            coin::CoinMetadata<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let treasury_cap = scenario.take_from_sender<
            coin::TreasuryCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();

        // Initialize pool
        lock_release_token_pool::initialize(
            &mut owner_cap,
            &mut ccip_ref,
            &coin_metadata,
            &treasury_cap,
            TOKEN_ADMIN,
            REBALANCER,
            scenario.ctx(),
        );

        transfer::public_transfer(owner_cap, TOKEN_ADMIN);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
        test_scenario::return_immutable(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);

    scenario.next_tx(TOKEN_ADMIN);
    {
        let mut pool_state = scenario.take_shared<
            LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap>();
        let ctx = scenario.ctx();
        let clock = clock::create_for_testing(ctx);

        // First add a chain to configure rate limits for
        lock_release_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[],
            vector[DefaultRemoteChain],
            vector[vector[DefaultRemotePool]],
            vector[DefaultRemoteToken],
        );

        // Test single chain rate limiter config
        lock_release_token_pool::set_chain_rate_limiter_config(
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

        // Test multiple chains rate limiter config
        let chain2 = 3000;
        let chain3 = 4000;

        // Add more chains first
        lock_release_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[],
            vector[chain2, chain3],
            vector[vector[b"pool2"], vector[b"pool3"]],
            vector[b"token2", b"token3"],
        );

        // Configure rate limits for multiple chains
        lock_release_token_pool::set_chain_rate_limiter_configs(
            &mut pool_state,
            &owner_cap,
            &clock,
            vector[chain2, chain3],
            vector[true, false], // outbound enabled
            vector[1500, 2500], // outbound capacities
            vector[150, 250], // outbound rates
            vector[false, true], // inbound enabled
            vector[3000, 4000], // inbound capacities
            vector[300, 400], // inbound rates
        );

        scenario.return_to_sender(owner_cap);
        test_scenario::return_shared(pool_state);
        clock.destroy_for_testing();
    };

    test_scenario::return_shared(ccip_ref);
    test_scenario::end(scenario);
}

#[test]
public fun test_allowlist_management() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN);

    // Setup CCIP environment
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(TOKEN_ADMIN);
    {
        let ctx = scenario.ctx();

        // Create test token
        let (treasury_cap, coin_metadata) = coin::create_currency(
            LOCK_RELEASE_TOKEN_POOL_TESTS {},
            Decimals,
            b"TEST",
            b"TestToken",
            b"test_token",
            option::none(),
            ctx,
        );

        // Call test_init to create owner_cap
        lock_release_token_pool::test_init(ctx);

        transfer::public_freeze_object(coin_metadata);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
    };

    scenario.next_tx(TOKEN_ADMIN);
    {
        let mut owner_cap = scenario.take_from_sender<OwnerCap>();
        let coin_metadata = scenario.take_immutable<
            coin::CoinMetadata<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let treasury_cap = scenario.take_from_sender<
            coin::TreasuryCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();

        // Initialize pool
        lock_release_token_pool::initialize(
            &mut owner_cap,
            &mut ccip_ref,
            &coin_metadata,
            &treasury_cap,
            TOKEN_ADMIN,
            REBALANCER,
            scenario.ctx(),
        );

        transfer::public_transfer(owner_cap, TOKEN_ADMIN);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
        test_scenario::return_immutable(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);

    scenario.next_tx(TOKEN_ADMIN);
    {
        let mut pool_state = scenario.take_shared<
            LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        // Test initial allowlist state (should be disabled and empty)
        assert!(!lock_release_token_pool::get_allowlist_enabled(&pool_state));
        let initial_allowlist = lock_release_token_pool::get_allowlist(&pool_state);
        assert!(initial_allowlist.length() == 0);

        // Test enabling allowlist
        lock_release_token_pool::set_allowlist_enabled(&mut pool_state, &owner_cap, true);

        // Test adding addresses to allowlist
        let addresses_to_add = vector[@0x111, @0x222, @0x333];
        lock_release_token_pool::apply_allowlist_updates(
            &mut pool_state,
            &owner_cap,
            vector[], // remove none
            addresses_to_add,
        );

        // Verify addresses were added
        let updated_allowlist = lock_release_token_pool::get_allowlist(&pool_state);
        assert!(updated_allowlist.length() == 3);
        assert!(updated_allowlist == addresses_to_add);

        // Test removing some addresses from allowlist
        let addresses_to_remove = vector[@0x222];
        lock_release_token_pool::apply_allowlist_updates(
            &mut pool_state,
            &owner_cap,
            addresses_to_remove,
            vector[], // add none
        );

        // Verify address was removed
        let final_allowlist = lock_release_token_pool::get_allowlist(&pool_state);
        assert!(final_allowlist.length() == 2);
        assert!(final_allowlist == vector[@0x111, @0x333]);

        scenario.return_to_sender(owner_cap);
        test_scenario::return_shared(pool_state);
    };

    test_scenario::return_shared(ccip_ref);
    test_scenario::end(scenario);
}

#[test]
#[expected_failure(abort_code = lock_release_token_pool::EInvalidRebalancerCap)]
public fun test_unauthorized_liquidity_provision() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN);

    // Setup CCIP environment
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(TOKEN_ADMIN);
    {
        let ctx = scenario.ctx();

        // Create test token and initialize pool
        let (mut treasury_cap, coin_metadata) = coin::create_currency(
            LOCK_RELEASE_TOKEN_POOL_TESTS {},
            Decimals,
            b"TEST",
            b"TestToken",
            b"test_token",
            option::none(),
            ctx,
        );

        // Call test_init to create owner_cap
        lock_release_token_pool::test_init(ctx);

        // Mint some tokens for testing
        let test_coin = coin::mint(&mut treasury_cap, 1000000, ctx);
        transfer::public_transfer(test_coin, @0x999); // unauthorized user

        transfer::public_freeze_object(coin_metadata);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
    };

    scenario.next_tx(TOKEN_ADMIN);
    {
        let mut owner_cap = scenario.take_from_sender<OwnerCap>();
        let coin_metadata = scenario.take_immutable<
            coin::CoinMetadata<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let treasury_cap = scenario.take_from_sender<
            coin::TreasuryCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();

        // Initialize pool
        lock_release_token_pool::initialize(
            &mut owner_cap,
            &mut ccip_ref,
            &coin_metadata,
            &treasury_cap,
            TOKEN_ADMIN,
            REBALANCER,
            scenario.ctx(),
        );

        transfer::public_transfer(owner_cap, TOKEN_ADMIN);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
        test_scenario::return_immutable(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);

    // Try to provide liquidity as unauthorized user (should fail)
    scenario.next_tx(@0x999);
    {
        let mut pool_state = scenario.take_shared<
            LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let liquidity_coin = scenario.take_from_sender<coin::Coin<LOCK_RELEASE_TOKEN_POOL_TESTS>>();

        // This should fail with EInvalidRebalancerCap using a fake cap
        let fake_rebalancer_cap = lock_release_token_pool::create_fake_rebalancer_cap(scenario.ctx());
        lock_release_token_pool::provide_liquidity(
            &mut pool_state,
            &fake_rebalancer_cap,
            liquidity_coin,
            scenario.ctx(),
        );

        // Consume the fake cap to avoid linear resource leak
        transfer::public_transfer(fake_rebalancer_cap, @0x0);
        test_scenario::return_shared(pool_state);
    };

    test_scenario::return_shared(ccip_ref);
    test_scenario::end(scenario);
}

#[test]
#[expected_failure(abort_code = lock_release_token_pool::ETokenPoolBalanceTooLow)]
public fun test_withdraw_exceeds_balance() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN);

    // Setup CCIP environment
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(TOKEN_ADMIN);
    {
        let ctx = scenario.ctx();

        // Create test token and initialize pool
        let (mut treasury_cap, coin_metadata) = coin::create_currency(
            LOCK_RELEASE_TOKEN_POOL_TESTS {},
            Decimals,
            b"TEST",
            b"TestToken",
            b"test_token",
            option::none(),
            ctx,
        );

        // Call test_init to create owner_cap
        lock_release_token_pool::test_init(ctx);

        // Mint and provide some liquidity
        let test_coin = coin::mint(&mut treasury_cap, 100000, ctx); // Only 100k tokens
        transfer::public_transfer(test_coin, REBALANCER);

        transfer::public_freeze_object(coin_metadata);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
    };

    scenario.next_tx(TOKEN_ADMIN);
    {
        let mut owner_cap = scenario.take_from_sender<OwnerCap>();
        let coin_metadata = scenario.take_immutable<
            coin::CoinMetadata<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let treasury_cap = scenario.take_from_sender<
            coin::TreasuryCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();

        // Initialize pool
        lock_release_token_pool::initialize(
            &mut owner_cap,
            &mut ccip_ref,
            &coin_metadata,
            &treasury_cap,
            TOKEN_ADMIN,
            REBALANCER,
            scenario.ctx(),
        );

        transfer::public_transfer(owner_cap, TOKEN_ADMIN);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
        test_scenario::return_immutable(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);

    // Provide liquidity
    scenario.next_tx(REBALANCER);
    {
        let mut pool_state = scenario.take_shared<
            LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let liquidity_coin = scenario.take_from_sender<coin::Coin<LOCK_RELEASE_TOKEN_POOL_TESTS>>();

        let rebalancer_cap = scenario.take_from_sender<
            RebalancerCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        lock_release_token_pool::provide_liquidity(
            &mut pool_state,
            &rebalancer_cap,
            liquidity_coin,
            scenario.ctx(),
        );
        transfer::public_transfer(rebalancer_cap, REBALANCER);

        test_scenario::return_shared(pool_state);
    };

    // Try to withdraw more than available (should fail)
    scenario.next_tx(REBALANCER);
    {
        let mut pool_state = scenario.take_shared<
            LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let rebalancer_cap = scenario.take_from_sender<
            RebalancerCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();

        // Try to withdraw 200k tokens when only 100k are available
        let withdrawn_coin = lock_release_token_pool::withdraw_liquidity<
            LOCK_RELEASE_TOKEN_POOL_TESTS,
        >(
            &mut pool_state,
            &rebalancer_cap,
            200000, // More than available
            scenario.ctx(),
        );

        transfer::public_transfer(withdrawn_coin, scenario.ctx().sender());
        transfer::public_transfer(rebalancer_cap, REBALANCER);
        test_scenario::return_shared(pool_state);
    };

    test_scenario::return_shared(ccip_ref);
    test_scenario::end(scenario);
}

#[test]
#[expected_failure(abort_code = lock_release_token_pool::EInvalidRebalancerCap)]
public fun test_unauthorized_withdrawal() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN);

    // Setup CCIP environment
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(TOKEN_ADMIN);
    {
        let ctx = scenario.ctx();

        // Create test token and initialize pool
        let (mut treasury_cap, coin_metadata) = coin::create_currency(
            LOCK_RELEASE_TOKEN_POOL_TESTS {},
            Decimals,
            b"TEST",
            b"TestToken",
            b"test_token",
            option::none(),
            ctx,
        );

        // Call test_init to create owner_cap
        lock_release_token_pool::test_init(ctx);

        // Mint tokens and provide some liquidity
        let liquidity_coin = coin::mint(&mut treasury_cap, 500000, ctx);
        transfer::public_transfer(liquidity_coin, REBALANCER);

        transfer::public_freeze_object(coin_metadata);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
    };

    scenario.next_tx(TOKEN_ADMIN);
    {
        let mut owner_cap = scenario.take_from_sender<OwnerCap>();
        let coin_metadata = scenario.take_immutable<
            coin::CoinMetadata<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let treasury_cap = scenario.take_from_sender<
            coin::TreasuryCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();

        // Initialize pool
        lock_release_token_pool::initialize(
            &mut owner_cap,
            &mut ccip_ref,
            &coin_metadata,
            &treasury_cap,
            TOKEN_ADMIN,
            REBALANCER,
            scenario.ctx(),
        );

        transfer::public_transfer(owner_cap, TOKEN_ADMIN);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
        test_scenario::return_immutable(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);

    // Provide liquidity as authorized rebalancer
    scenario.next_tx(REBALANCER);
    {
        let mut pool_state = scenario.take_shared<
            LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let liquidity_coin = scenario.take_from_sender<coin::Coin<LOCK_RELEASE_TOKEN_POOL_TESTS>>();

        let rebalancer_cap = scenario.take_from_sender<
            RebalancerCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        lock_release_token_pool::provide_liquidity(
            &mut pool_state,
            &rebalancer_cap,
            liquidity_coin,
            scenario.ctx(),
        );
        transfer::public_transfer(rebalancer_cap, REBALANCER);

        test_scenario::return_shared(pool_state);
    };

    // Try to withdraw as unauthorized user (should fail)
    scenario.next_tx(@0x999); // unauthorized user
    {
        let mut pool_state = scenario.take_shared<
            LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();

        // This should fail with EInvalidRebalancerCap using a fake cap
        let fake_rebalancer_cap = lock_release_token_pool::create_fake_rebalancer_cap(scenario.ctx());
        let withdrawn_coin = lock_release_token_pool::withdraw_liquidity<
            LOCK_RELEASE_TOKEN_POOL_TESTS,
        >(
            &mut pool_state,
            &fake_rebalancer_cap,
            100000,
            scenario.ctx(),
        );

        // Transfer the coin to consume it
        transfer::public_transfer(withdrawn_coin, scenario.ctx().sender());
        // Consume the fake cap to avoid linear resource leak
        transfer::public_transfer(fake_rebalancer_cap, @0x0);

        test_scenario::return_shared(pool_state);
    };

    test_scenario::return_shared(ccip_ref);
    test_scenario::end(scenario);
}

#[test]
public fun test_destroy_token_pool() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN);

    // Setup CCIP environment
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(TOKEN_ADMIN);
    {
        let ctx = scenario.ctx();

        // Create test token and initialize pool
        let (mut treasury_cap, coin_metadata) = coin::create_currency(
            LOCK_RELEASE_TOKEN_POOL_TESTS {},
            Decimals,
            b"TEST",
            b"TestToken",
            b"test_token",
            option::none(),
            ctx,
        );

        // Call test_init to create owner_cap
        lock_release_token_pool::test_init(ctx);

        // Mint tokens and provide some liquidity
        let liquidity_coin = coin::mint(&mut treasury_cap, 500000, ctx);
        transfer::public_transfer(liquidity_coin, REBALANCER);

        transfer::public_freeze_object(coin_metadata);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
    };

    scenario.next_tx(TOKEN_ADMIN);
    {
        let mut owner_cap = scenario.take_from_sender<OwnerCap>();
        let coin_metadata = scenario.take_immutable<
            coin::CoinMetadata<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let treasury_cap = scenario.take_from_sender<
            coin::TreasuryCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();

        // Initialize pool
        lock_release_token_pool::initialize(
            &mut owner_cap,
            &mut ccip_ref,
            &coin_metadata,
            &treasury_cap,
            TOKEN_ADMIN,
            REBALANCER,
            scenario.ctx(),
        );

        transfer::public_transfer(owner_cap, TOKEN_ADMIN);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
        test_scenario::return_immutable(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);

    // Provide some liquidity to test destruction with remaining balance
    scenario.next_tx(REBALANCER);
    {
        let mut pool_state = scenario.take_shared<
            LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let liquidity_coin = scenario.take_from_sender<coin::Coin<LOCK_RELEASE_TOKEN_POOL_TESTS>>();
        let rebalancer_cap = scenario.take_from_sender<
            RebalancerCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();

        lock_release_token_pool::provide_liquidity(
            &mut pool_state,
            &rebalancer_cap,
            liquidity_coin,
            scenario.ctx(),
        );
        transfer::public_transfer(rebalancer_cap, REBALANCER);

        test_scenario::return_shared(pool_state);
    };

    // Test pool destruction
    scenario.next_tx(TOKEN_ADMIN);
    {
        let pool_state = scenario.take_shared<
            LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap>();
        let ctx = scenario.ctx();

        let initial_balance = lock_release_token_pool::get_balance<LOCK_RELEASE_TOKEN_POOL_TESTS>(
            &pool_state,
        );
        assert!(initial_balance > 0); // Should have some liquidity

        token_admin_registry::unregister_pool(
            &mut ccip_ref,
            lock_release_token_pool::get_token(&pool_state),
            ctx,
        );

        // Destroy the pool and get remaining balance
        let remaining_coin = lock_release_token_pool::destroy_token_pool<
            LOCK_RELEASE_TOKEN_POOL_TESTS,
        >(
            &mut ccip_ref,
            pool_state,
            owner_cap,
            ctx,
        );

        // Verify remaining balance matches what was in the pool
        assert!(coin::value(&remaining_coin) == initial_balance);

        // Transfer the remaining balance
        transfer::public_transfer(remaining_coin, ctx.sender());
    };

    test_scenario::return_shared(ccip_ref);
    test_scenario::end(scenario);
}

#[test]
public fun test_edge_cases_and_getters() {
    let mut scenario = create_test_scenario(TOKEN_ADMIN);

    // Setup CCIP environment
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    scenario.next_tx(TOKEN_ADMIN);
    {
        let ctx = scenario.ctx();

        // Create test token
        let (treasury_cap, coin_metadata) = coin::create_currency(
            LOCK_RELEASE_TOKEN_POOL_TESTS {},
            Decimals,
            b"TEST",
            b"TestToken",
            b"test_token",
            option::none(),
            ctx,
        );

        // Call test_init to create owner_cap
        lock_release_token_pool::test_init(ctx);

        transfer::public_freeze_object(coin_metadata);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
    };

    scenario.next_tx(TOKEN_ADMIN);
    {
        let mut owner_cap = scenario.take_from_sender<OwnerCap>();
        let coin_metadata = scenario.take_immutable<
            coin::CoinMetadata<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let treasury_cap = scenario.take_from_sender<
            coin::TreasuryCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();

        // Initialize pool
        lock_release_token_pool::initialize(
            &mut owner_cap,
            &mut ccip_ref,
            &coin_metadata,
            &treasury_cap,
            TOKEN_ADMIN,
            REBALANCER,
            scenario.ctx(),
        );

        transfer::public_transfer(owner_cap, TOKEN_ADMIN);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
        test_scenario::return_immutable(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);

    scenario.next_tx(TOKEN_ADMIN);
    {
        let mut pool_state = scenario.take_shared<
            LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        // Test applying empty chain updates (should work)
        lock_release_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[], // remove none
            vector[], // add none
            vector[], // no pools
            vector[], // no tokens
        );

        // Verify no chains were added
        let supported_chains = lock_release_token_pool::get_supported_chains(&pool_state);
        assert!(supported_chains.length() == 0);

        // Test getting token address
        let token_address = lock_release_token_pool::get_token(&pool_state);
        assert!(token_address != @0x0); // Should be a valid address

        // Add a chain and test getters
        lock_release_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[],
            vector[DefaultRemoteChain],
            vector[vector[DefaultRemotePool]],
            vector[DefaultRemoteToken],
        );

        // Test getters with existing chain
        let remote_token_after = lock_release_token_pool::get_remote_token(
            &pool_state,
            DefaultRemoteChain,
        );
        assert!(remote_token_after == DefaultRemoteToken);

        let remote_pools = lock_release_token_pool::get_remote_pools(
            &pool_state,
            DefaultRemoteChain,
        );
        assert!(remote_pools.length() == 1);
        assert!(remote_pools[0] == DefaultRemotePool);

        // Test is_remote_pool with existing and non-existing pools
        assert!(
            lock_release_token_pool::is_remote_pool(
                &pool_state,
                DefaultRemoteChain,
                DefaultRemotePool,
            ),
        );
        assert!(
            !lock_release_token_pool::is_remote_pool(
                &pool_state,
                DefaultRemoteChain,
                b"non_existent_pool",
            ),
        );

        scenario.return_to_sender(owner_cap);
        test_scenario::return_shared(pool_state);
    };

    test_scenario::return_shared(ccip_ref);
    test_scenario::end(scenario);
}

#[test]
public fun test_lock_or_burn_functionality() {
    let mut scenario = test_scenario::begin(TOKEN_ADMIN);
    let ctx = scenario.ctx();

    // Setup CCIP environment
    state_object::test_init(ctx);

    // Advance to next transaction to retrieve the created objects
    scenario.next_tx(TOKEN_ADMIN);

    // Retrieve the OwnerCap that was transferred to the sender
    let ccip_owner_cap = scenario.take_from_sender<CCIPOwnerCap>();

    // Retrieve the shared CCIPObjectRef
    let mut ccip_ref = scenario.take_shared<CCIPObjectRef>();

    upgrade_registry::initialize(&mut ccip_ref, &ccip_owner_cap, scenario.ctx());
    token_admin_registry::initialize(&mut ccip_ref, &ccip_owner_cap, scenario.ctx());
    rmn_remote::initialize(&mut ccip_ref, &ccip_owner_cap, 1000, scenario.ctx());
    onramp_sh::test_init(scenario.ctx());

    // Create test token and initialize pool in the same transaction
    let ctx = scenario.ctx();
    let (mut treasury_cap, coin_metadata) = coin::create_currency(
        LOCK_RELEASE_TOKEN_POOL_TESTS {},
        Decimals,
        b"TEST",
        b"TestToken",
        b"test_token",
        option::none(),
        ctx,
    );

    lock_release_token_pool::test_init(ctx);

    // Mint some tokens for testing
    let test_coin = coin::mint(&mut treasury_cap, 5000, ctx);
    transfer::public_transfer(test_coin, @0x456); // Transfer to test user

    transfer::public_freeze_object(coin_metadata);
    transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
    transfer::public_transfer(ccip_owner_cap, @0x0);

    scenario.next_tx(TOKEN_ADMIN);
    {
        let mut owner_cap = scenario.take_from_sender<OwnerCap>();
        let coin_metadata = scenario.take_immutable<
            coin::CoinMetadata<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let treasury_cap = scenario.take_from_sender<
            coin::TreasuryCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();

        // Initialize pool
        lock_release_token_pool::initialize(
            &mut owner_cap,
            &mut ccip_ref,
            &coin_metadata,
            &treasury_cap,
            TOKEN_ADMIN,
            REBALANCER,
            scenario.ctx(),
        );

        transfer::public_transfer(owner_cap, TOKEN_ADMIN);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
        test_scenario::return_immutable(coin_metadata);
    };

    scenario.next_tx(TOKEN_ADMIN);
    {
        let mut pool_state = scenario.take_shared<
            LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap>();
        let source_transfer_cap = scenario.take_from_address<onramp_sh::SourceTransferCap>(
            TOKEN_ADMIN,
        );
        let mut ctx = tx_context::dummy();
        let clock = clock::create_for_testing(&mut ctx);

        // Configure chain and rate limiter
        lock_release_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[],
            vector[DefaultRemoteChain],
            vector[vector[DefaultRemotePool]],
            vector[DefaultRemoteToken],
        );

        lock_release_token_pool::set_chain_rate_limiter_config(
            &mut pool_state,
            &owner_cap,
            &clock,
            DefaultRemoteChain,
            true, // outbound enabled
            20000, // high capacity
            2000, // high rate
            true, // inbound enabled
            20000, // high capacity
            2000, // high rate
        );

        clock.destroy_for_testing();
        transfer::public_transfer(source_transfer_cap, TOKEN_ADMIN);
        transfer::public_transfer(owner_cap, TOKEN_ADMIN);
        test_scenario::return_shared(pool_state);
    };

    // Test actual lock_or_burn function call
    scenario.next_tx(@0x456); // Switch to test user
    {
        let mut pool_state = scenario.take_shared<
            LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let test_coin = scenario.take_from_sender<coin::Coin<LOCK_RELEASE_TOKEN_POOL_TESTS>>();
        let mut ctx = tx_context::dummy();
        let mut clock = clock::create_for_testing(&mut ctx);
        clock.increment_for_testing(1000000000); // Advance clock for rate limiter

        let initial_coin_value = coin::value(&test_coin);
        assert!(initial_coin_value == 5000);

        let initial_pool_balance = lock_release_token_pool::get_balance<
            LOCK_RELEASE_TOKEN_POOL_TESTS,
        >(&pool_state);

        let mut token_transfer_params = onramp_sh::create_token_transfer_params(
            bcs::to_bytes(&@0x456),
        ); // Use the test user address as token receiver

        // Call the actual lock_or_burn function
        lock_release_token_pool::lock_or_burn<LOCK_RELEASE_TOKEN_POOL_TESTS>(
            &ccip_ref,
            &mut token_transfer_params,
            test_coin, // This coin gets locked in the pool
            DefaultRemoteChain,
            &clock,
            &mut pool_state,
            &mut ctx,
        );

        // Verify pool balance increased by the locked amount
        let new_pool_balance = lock_release_token_pool::get_balance<LOCK_RELEASE_TOKEN_POOL_TESTS>(
            &pool_state,
        );
        assert!(new_pool_balance == initial_pool_balance + initial_coin_value);

        // Clean up token params
        let source_transfer_cap = scenario.take_from_address<onramp_sh::SourceTransferCap>(
            TOKEN_ADMIN,
        );

        //TOOD: add token package ID to omnramp state helper to continue with this test
        // Calculate the actual package ID from TypeProof (same as initialization)
        let type_proof_type_name = type_name::with_defining_ids<
            lock_release_token_pool::TypeProof,
        >();
        let _type_proof_type_name_address = type_proof_type_name.address_string();
        let actual_package_id = address::from_ascii_bytes(
            &_type_proof_type_name_address.into_bytes(),
        );

        let (
            chain_selector,
            token_pool_package_id,
            amount,
            source_token_address,
            dest_token_address,
            extra_data,
        ) = onramp_sh::get_source_token_transfer_data(&token_transfer_params);
        assert!(chain_selector == DefaultRemoteChain);
        assert!(token_pool_package_id == actual_package_id);
        assert!(amount == initial_coin_value);
        assert!(source_token_address == lock_release_token_pool::get_token(&pool_state));
        assert!(dest_token_address == DefaultRemoteToken);
        assert!(extra_data.length() > 0); // Should contain encoded decimals
        onramp_sh::deconstruct_token_params(&source_transfer_cap, token_transfer_params);

        clock.destroy_for_testing();
        transfer::public_transfer(source_transfer_cap, TOKEN_ADMIN);
        test_scenario::return_shared(pool_state);
    };

    test_scenario::return_shared(ccip_ref);
    test_scenario::end(scenario);
}

#[test]
public fun test_release_or_mint_functionality() {
    let mut scenario = test_scenario::begin(TOKEN_ADMIN);
    let ctx = scenario.ctx();

    // Setup CCIP environment
    state_object::test_init(ctx);

    // Advance to next transaction to retrieve the created objects
    scenario.next_tx(TOKEN_ADMIN);

    // Retrieve the OwnerCap that was transferred to the sender
    let ccip_owner_cap = scenario.take_from_sender<CCIPOwnerCap>();

    // Retrieve the shared CCIPObjectRef
    let mut ccip_ref = scenario.take_shared<CCIPObjectRef>();

    upgrade_registry::initialize(&mut ccip_ref, &ccip_owner_cap, scenario.ctx());
    token_admin_registry::initialize(&mut ccip_ref, &ccip_owner_cap, scenario.ctx());
    rmn_remote::initialize(&mut ccip_ref, &ccip_owner_cap, 1000, scenario.ctx());
    offramp_sh::test_init(scenario.ctx());

    scenario.next_tx(TOKEN_ADMIN);
    {
        let ctx = scenario.ctx();

        // Create test token and initialize pool
        let (mut treasury_cap, coin_metadata) = coin::create_currency(
            LOCK_RELEASE_TOKEN_POOL_TESTS {},
            Decimals,
            b"TEST",
            b"TestToken",
            b"test_token",
            option::none(),
            ctx,
        );

        let _coin_metadata_address = object::id_to_address(&object::id(&coin_metadata));

        // Call test_init to create owner_cap
        lock_release_token_pool::test_init(ctx);

        // Mint tokens and provide liquidity to the pool for release operations
        let liquidity_coin = coin::mint(&mut treasury_cap, 20000, ctx);
        transfer::public_transfer(liquidity_coin, REBALANCER);

        transfer::public_freeze_object(coin_metadata);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
    };

    scenario.next_tx(TOKEN_ADMIN);
    {
        let mut owner_cap = scenario.take_from_sender<OwnerCap>();
        let coin_metadata = scenario.take_immutable<
            coin::CoinMetadata<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let treasury_cap = scenario.take_from_sender<
            coin::TreasuryCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();

        // Initialize pool
        lock_release_token_pool::initialize(
            &mut owner_cap,
            &mut ccip_ref,
            &coin_metadata,
            &treasury_cap,
            TOKEN_ADMIN,
            REBALANCER,
            scenario.ctx(),
        );

        transfer::public_transfer(owner_cap, TOKEN_ADMIN);
        transfer::public_transfer(treasury_cap, TOKEN_ADMIN);
        test_scenario::return_immutable(coin_metadata);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);

    scenario.next_tx(TOKEN_ADMIN);
    {
        let mut pool_state = scenario.take_shared<
            LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap>();
        let dest_transfer_cap = scenario.take_from_sender<offramp_sh::DestTransferCap>();
        let mut ctx = tx_context::dummy();
        let mut clock = clock::create_for_testing(&mut ctx);

        // Configure chain and rate limiter
        lock_release_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[],
            vector[DefaultRemoteChain],
            vector[vector[DefaultRemotePool]],
            vector[DefaultRemoteToken],
        );

        lock_release_token_pool::set_chain_rate_limiter_config(
            &mut pool_state,
            &owner_cap,
            &clock,
            DefaultRemoteChain,
            true, // outbound enabled
            20000, // high capacity
            2000, // high rate
            true, // inbound enabled
            20000, // high capacity
            2000, // high rate
        );

        // Advance clock for rate limiter
        clock.increment_for_testing(1000000000);

        clock.destroy_for_testing();
        transfer::public_transfer(dest_transfer_cap, TOKEN_ADMIN);
        transfer::public_transfer(owner_cap, TOKEN_ADMIN);
        test_scenario::return_shared(pool_state);
    };

    // Provide liquidity to the pool
    scenario.next_tx(REBALANCER);
    {
        let mut pool_state = scenario.take_shared<
            LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let liquidity_coin = scenario.take_from_sender<coin::Coin<LOCK_RELEASE_TOKEN_POOL_TESTS>>();

        let rebalancer_cap = scenario.take_from_sender<
            RebalancerCap<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        lock_release_token_pool::provide_liquidity(
            &mut pool_state,
            &rebalancer_cap,
            liquidity_coin,
            scenario.ctx(),
        );
        transfer::public_transfer(rebalancer_cap, REBALANCER);

        test_scenario::return_shared(pool_state);
    };

    // Proceed with the release_or_mint test
    scenario.next_tx(TOKEN_ADMIN);
    {
        let mut pool_state = scenario.take_shared<
            LockReleaseTokenPoolState<LOCK_RELEASE_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_address<OwnerCap>(TOKEN_ADMIN);
        let dest_transfer_cap = scenario.take_from_address<offramp_sh::DestTransferCap>(
            TOKEN_ADMIN,
        );
        let mut ctx = tx_context::dummy();
        let mut clock = clock::create_for_testing(&mut ctx);

        // Advance clock for rate limiter
        clock.increment_for_testing(1000000000);

        // Get the coin metadata address for the test
        let coin_metadata_address = lock_release_token_pool::get_token(&pool_state);

        // Calculate the actual package ID from TypeProof (same as initialization)
        let type_proof_type_name = type_name::with_defining_ids<
            lock_release_token_pool::TypeProof,
        >();
        let _type_proof_type_name_address = type_proof_type_name.address_string();
        let actual_package_id = address::from_ascii_bytes(
            &_type_proof_type_name_address.into_bytes(),
        );

        // Create receiver params for release_or_mint
        let mut receiver_params = offramp_sh::create_receiver_params(
            &dest_transfer_cap,
            DefaultRemoteChain,
        );

        // Add token transfer to receiver params
        let receiver_address = @0x789;
        let source_amount = 8000;
        let source_pool_data = x"0000000000000000000000000000000000000000000000000000000000000008"; // 8 decimals encoded
        let offchain_data = vector[];

        offramp_sh::add_dest_token_transfer(
            &dest_transfer_cap,
            &mut receiver_params,
            receiver_address,
            DefaultRemoteChain, // remote_chain_selector
            source_amount,
            coin_metadata_address,
            actual_package_id, // Use the dynamically calculated package ID
            DefaultRemotePool,
            source_pool_data,
            offchain_data,
        );

        let initial_pool_balance = lock_release_token_pool::get_balance<
            LOCK_RELEASE_TOKEN_POOL_TESTS,
        >(&pool_state);

        // Verify the operation setup
        let source_chain = offramp_sh::get_source_chain_selector(&receiver_params);
        assert!(source_chain == DefaultRemoteChain);

        // Call the actual release_or_mint function
        lock_release_token_pool::release_or_mint<LOCK_RELEASE_TOKEN_POOL_TESTS>(
            &ccip_ref,
            &mut receiver_params,
            &clock,
            &mut pool_state,
            &mut ctx,
        );

        // Verify pool balance decreased by the released amount
        let new_pool_balance = lock_release_token_pool::get_balance<LOCK_RELEASE_TOKEN_POOL_TESTS>(
            &pool_state,
        );
        assert!((new_pool_balance as u256) == (initial_pool_balance as u256) - source_amount);

        // Clean up receiver params
        offramp_sh::deconstruct_receiver_params(&dest_transfer_cap, receiver_params);

        clock.destroy_for_testing();
        transfer::public_transfer(dest_transfer_cap, TOKEN_ADMIN);
        transfer::public_transfer(owner_cap, TOKEN_ADMIN);
        test_scenario::return_shared(pool_state);
    };

    // Verify the released coin was transferred to the receiver
    scenario.next_tx(@0x789); // Switch to receiver
    {
        let released_coin = scenario.take_from_sender<coin::Coin<LOCK_RELEASE_TOKEN_POOL_TESTS>>();

        // Verify the released amount
        let released_amount = coin::value(&released_coin);
        assert!(released_amount == 8000);

        // Transfer the coin back to clean up
        transfer::public_transfer(released_coin, TOKEN_ADMIN);
    };

    test_scenario::return_shared(ccip_ref);
    test_scenario::end(scenario);
}

// ================================================================
// |             Rate Limiter calculate_refill Tests             |
// ================================================================

#[test]
public fun test_calculate_refill_at_capacity() {
    use lock_release_token_pool::rate_limiter;

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
    use lock_release_token_pool::rate_limiter;

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
    use lock_release_token_pool::rate_limiter;

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
    use lock_release_token_pool::rate_limiter;

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
    use lock_release_token_pool::rate_limiter;

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
    use lock_release_token_pool::rate_limiter;

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
    use lock_release_token_pool::rate_limiter;

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
    use lock_release_token_pool::rate_limiter;

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
    use lock_release_token_pool::rate_limiter;

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
    use lock_release_token_pool::rate_limiter;

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
