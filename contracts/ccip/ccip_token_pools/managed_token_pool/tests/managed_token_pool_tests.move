#[test_only]
module managed_token_pool::managed_token_pool_tests;

use ccip::offramp_state_helper as offramp_sh;
use ccip::onramp_state_helper as onramp_sh;
use ccip::ownable::OwnerCap as CCIPOwnerCap;
use ccip::rmn_remote;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::token_admin_registry;
use ccip::upgrade_registry;
use managed_token::managed_token::{Self, TokenState, MintCap};
use managed_token::ownable::OwnerCap as TokenOwnerCap;
use managed_token_pool::managed_token_pool::{Self, ManagedTokenPoolState};
use managed_token_pool::ownable::OwnerCap;
use std::bcs;
use std::string;
use std::type_name;
use sui::address;
use sui::clock;
use sui::coin;
use sui::deny_list::{Self, DenyList};
use sui::package;
use sui::test_scenario::{Self, Scenario};

public struct MANAGED_TOKEN_POOL_TESTS has drop {}

const Decimals: u8 = 8;
const DefaultRemoteChain: u64 = 2000;
const DefaultRemotePool: vector<u8> = b"default_remote_pool";
const DefaultRemoteToken: vector<u8> = b"default_remote_token";

const CCIP_ADMIN: address = @0x400;

fun setup_ccip_environment(scenario: &mut Scenario): (CCIPOwnerCap, CCIPObjectRef) {
    scenario.next_tx(@0x0);

    // Create deny list for testing (shared object, will be taken later)
    deny_list::create_for_testing(scenario.ctx());

    scenario.next_tx(CCIP_ADMIN);

    // Create CCIP state object
    state_object::test_init(scenario.ctx());

    // Advance to next transaction to retrieve the created objects
    scenario.next_tx(CCIP_ADMIN);

    // Retrieve the OwnerCap that was transferred to the sender
    let ccip_owner_cap = scenario.take_from_sender<CCIPOwnerCap>();

    // Retrieve the shared CCIPObjectRef
    let mut ccip_ref = scenario.take_shared<CCIPObjectRef>();

    // Initialize required CCIP modules
    rmn_remote::initialize(&mut ccip_ref, &ccip_owner_cap, 1000, scenario.ctx()); // local chain selector = 1000
    upgrade_registry::initialize(&mut ccip_ref, &ccip_owner_cap, scenario.ctx());
    token_admin_registry::initialize(&mut ccip_ref, &ccip_owner_cap, scenario.ctx());
    onramp_sh::test_init(scenario.ctx());
    offramp_sh::test_init(scenario.ctx());

    (ccip_owner_cap, ccip_ref)
}

#[test]
public fun test_type_and_version() {
    let version = managed_token_pool::type_and_version();
    assert!(version == string::utf8(b"ManagedTokenPool 1.6.0"));
}

#[test]
public fun test_initialize_and_basic_functionality() {
    let mut scenario = test_scenario::begin(@managed_token_pool);

    // Setup CCIP environment
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);
    let deny_list = scenario.take_shared<DenyList>();

    // Create managed token
    scenario.next_tx(@managed_token_pool);
    let coin_metadata = {
        let ctx = scenario.ctx();
        let (treasury_cap, coin_metadata) = coin::create_currency(
            MANAGED_TOKEN_POOL_TESTS {},
            Decimals,
            b"MTPT",
            b"ManagedTokenPoolTest",
            b"managed_token_pool_test",
            option::none(),
            ctx,
        );

        managed_token::initialize(
            treasury_cap,
            package::test_claim(MANAGED_TOKEN_POOL_TESTS {}, ctx),
            ctx,
        );
        coin_metadata
    };

    scenario.next_tx(@managed_token_pool);
    {
        let mut token_state = scenario.take_shared<TokenState<MANAGED_TOKEN_POOL_TESTS>>();
        let token_owner_cap = scenario.take_from_sender<TokenOwnerCap<MANAGED_TOKEN_POOL_TESTS>>();

        // Create mint cap for the pool
        managed_token::configure_new_minter(
            &mut token_state,
            &token_owner_cap,
            @managed_token_pool,
            0,
            true,
            scenario.ctx(),
        );

        scenario.return_to_sender(token_owner_cap);
        test_scenario::return_shared(token_state);
    };

    scenario.next_tx(@managed_token_pool);
    {
        // Call test_init to create owner_cap
        managed_token_pool::test_init(scenario.ctx());
    };

    scenario.next_tx(@managed_token_pool);
    {
        let mut owner_cap = scenario.take_from_sender<OwnerCap>();
        let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_POOL_TESTS>>();
        let token_state = scenario.take_shared<TokenState<MANAGED_TOKEN_POOL_TESTS>>();
        let token_owner_cap = scenario.take_from_sender<TokenOwnerCap<MANAGED_TOKEN_POOL_TESTS>>();

        managed_token_pool::initialize_with_managed_token(
            &mut owner_cap,
            &mut ccip_ref,
            &token_state,
            &token_owner_cap,
            &coin_metadata,
            mint_cap,
            @managed_token_pool,
            scenario.ctx(),
        );
        transfer::public_transfer(owner_cap, @managed_token_pool);
        scenario.return_to_sender(token_owner_cap);
        test_scenario::return_shared(token_state);
    };

    scenario.next_tx(@managed_token_pool);
    {
        let pool_state = scenario.take_shared<ManagedTokenPoolState<MANAGED_TOKEN_POOL_TESTS>>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        // Test basic getters
        assert!(managed_token_pool::get_token_decimals(&pool_state) == Decimals);
        let token_address = managed_token_pool::get_token(&pool_state);
        assert!(token_address != @0x0);

        // Test supported chains (should be empty initially)
        let supported_chains = managed_token_pool::get_supported_chains(&pool_state);
        assert!(supported_chains.length() == 0);

        scenario.return_to_sender(owner_cap);
        test_scenario::return_shared(pool_state);
    };

    transfer::public_freeze_object(coin_metadata);
    transfer::public_transfer(ccip_owner_cap, @0x0);
    test_scenario::return_shared(deny_list);
    test_scenario::return_shared(ccip_ref);
    scenario.end();
}

#[test]
public fun test_chain_configuration_management() {
    let mut scenario = test_scenario::begin(@managed_token_pool);
    let (ccip_owner_cap, ccip_ref, coin_metadata, deny_list) = setup_basic_pool(&mut scenario);

    scenario.next_tx(@managed_token_pool);
    {
        let mut pool_state = scenario.take_shared<
            ManagedTokenPoolState<MANAGED_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        // Test chain updates
        managed_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[],
            vector[DefaultRemoteChain],
            vector[vector[DefaultRemotePool]],
            vector[DefaultRemoteToken],
        );

        // Verify chain was added
        assert!(managed_token_pool::is_supported_chain(&pool_state, DefaultRemoteChain));
        let supported_chains = managed_token_pool::get_supported_chains(&pool_state);
        assert!(supported_chains.length() == 1);
        assert!(supported_chains[0] == DefaultRemoteChain);

        // Test remote pool management
        let new_remote_pool = b"new_remote_pool";
        managed_token_pool::add_remote_pool(
            &mut pool_state,
            &owner_cap,
            DefaultRemoteChain,
            new_remote_pool,
        );

        assert!(
            managed_token_pool::is_remote_pool(&pool_state, DefaultRemoteChain, new_remote_pool),
        );
        let remote_pools = managed_token_pool::get_remote_pools(&pool_state, DefaultRemoteChain);
        assert!(remote_pools.length() == 2);

        // Test remote pool removal
        managed_token_pool::remove_remote_pool(
            &mut pool_state,
            &owner_cap,
            DefaultRemoteChain,
            new_remote_pool,
        );

        assert!(
            !managed_token_pool::is_remote_pool(&pool_state, DefaultRemoteChain, new_remote_pool),
        );
        let remote_pools_after = managed_token_pool::get_remote_pools(
            &pool_state,
            DefaultRemoteChain,
        );
        assert!(remote_pools_after.length() == 1);

        let remote_token = managed_token_pool::get_remote_token(&pool_state, DefaultRemoteChain);
        assert!(remote_token == DefaultRemoteToken);

        scenario.return_to_sender(owner_cap);
        test_scenario::return_shared(pool_state);
    };

    transfer::public_freeze_object(coin_metadata);
    transfer::public_transfer(ccip_owner_cap, @0x0);
    test_scenario::return_shared(ccip_ref);
    test_scenario::return_shared(deny_list);
    scenario.end();
}

#[test]
public fun test_allowlist_management() {
    let mut scenario = test_scenario::begin(@managed_token_pool);
    let (ccip_owner_cap, ccip_ref, coin_metadata, deny_list) = setup_basic_pool(&mut scenario);

    scenario.next_tx(@managed_token_pool);
    {
        let pool_state = scenario.take_shared<ManagedTokenPoolState<MANAGED_TOKEN_POOL_TESTS>>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        // Test initial allowlist state
        assert!(!managed_token_pool::get_allowlist_enabled(&pool_state));
        let initial_allowlist = managed_token_pool::get_allowlist(&pool_state);
        assert!(initial_allowlist.length() == 0);

        scenario.return_to_sender(owner_cap);
        test_scenario::return_shared(pool_state);
    };

    cleanup_test(scenario, ccip_owner_cap, ccip_ref, coin_metadata, deny_list);
}

#[test]
public fun test_rate_limiter_configuration() {
    let mut scenario = test_scenario::begin(@managed_token_pool);
    let (ccip_owner_cap, ccip_ref, coin_metadata, deny_list) = setup_basic_pool(&mut scenario);

    scenario.next_tx(@managed_token_pool);
    {
        let mut pool_state = scenario.take_shared<
            ManagedTokenPoolState<MANAGED_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap>();
        let mut ctx = sui::tx_context::dummy();
        let clock = clock::create_for_testing(&mut ctx);

        // Add chain first
        managed_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[],
            vector[DefaultRemoteChain],
            vector[vector[DefaultRemotePool]],
            vector[DefaultRemoteToken],
        );

        // Test single chain rate limiter config
        managed_token_pool::set_chain_rate_limiter_config(
            &mut pool_state,
            &owner_cap,
            &clock,
            DefaultRemoteChain,
            true,
            1000,
            100,
            true,
            2000,
            200,
        );

        // Test multiple chains rate limiter config
        let chain2 = 3000;
        let chain3 = 4000;

        managed_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[],
            vector[chain2, chain3],
            vector[vector[b"pool2"], vector[b"pool3"]],
            vector[b"token2", b"token3"],
        );

        managed_token_pool::set_chain_rate_limiter_configs(
            &mut pool_state,
            &owner_cap,
            &clock,
            vector[chain2, chain3],
            vector[true, false],
            vector[1500, 2500],
            vector[150, 250],
            vector[false, true],
            vector[3000, 4000],
            vector[300, 400],
        );

        test_scenario::return_shared(pool_state);
        transfer::public_transfer(owner_cap, @managed_token_pool);
        clock.destroy_for_testing();
    };

    cleanup_test(scenario, ccip_owner_cap, ccip_ref, coin_metadata, deny_list);
}

#[test]
public fun test_lock_or_burn_functionality() {
    let mut scenario = test_scenario::begin(@managed_token_pool);

    // Setup CCIP environment
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    // Create managed token and pool
    scenario.next_tx(@managed_token_pool);
    onramp_sh::test_init(scenario.ctx()); // Initialize onramp state helper
    let coin_metadata = {
        let ctx = scenario.ctx();
        let (treasury_cap, coin_metadata) = coin::create_currency(
            MANAGED_TOKEN_POOL_TESTS {},
            Decimals,
            b"MTPT",
            b"ManagedTokenPoolTest",
            b"managed_token_pool_test",
            option::none(),
            ctx,
        );

        managed_token::initialize(
            treasury_cap,
            package::test_claim(MANAGED_TOKEN_POOL_TESTS {}, ctx),
            ctx,
        );
        coin_metadata
    };

    scenario.next_tx(@managed_token_pool);
    {
        let mut token_state = scenario.take_shared<TokenState<MANAGED_TOKEN_POOL_TESTS>>();
        let token_owner_cap = scenario.take_from_sender<TokenOwnerCap<MANAGED_TOKEN_POOL_TESTS>>();

        // Create mint caps for the pool and test user
        managed_token::configure_new_minter(
            &mut token_state,
            &token_owner_cap,
            @managed_token_pool,
            0,
            true,
            scenario.ctx(),
        );

        managed_token::configure_new_minter(
            &mut token_state,
            &token_owner_cap,
            @0x456,
            1000000,
            false,
            scenario.ctx(),
        );

        scenario.return_to_sender(token_owner_cap);
        test_scenario::return_shared(token_state);
    };

    scenario.next_tx(@managed_token_pool);
    {
        // Call test_init to create owner_cap
        managed_token_pool::test_init(scenario.ctx());
    };

    scenario.next_tx(@managed_token_pool);
    {
        let mut owner_cap = scenario.take_from_sender<OwnerCap>();
        let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_POOL_TESTS>>();
        let token_state = scenario.take_shared<TokenState<MANAGED_TOKEN_POOL_TESTS>>();
        let token_owner_cap = scenario.take_from_sender<TokenOwnerCap<MANAGED_TOKEN_POOL_TESTS>>();

        managed_token_pool::initialize_with_managed_token(
            &mut owner_cap,
            &mut ccip_ref,
            &token_state,
            &token_owner_cap,
            &coin_metadata,
            mint_cap,
            @managed_token_pool,
            scenario.ctx(),
        );

        transfer::public_transfer(owner_cap, @managed_token_pool);
        scenario.return_to_sender(token_owner_cap);
        test_scenario::return_shared(token_state);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);

    // Set up chain and rate limiting
    scenario.next_tx(@managed_token_pool);
    {
        let mut pool_state = scenario.take_shared<
            ManagedTokenPoolState<MANAGED_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap>();
        let source_transfer_cap = scenario.take_from_sender<onramp_sh::SourceTransferCap>();
        let mut ctx = sui::tx_context::dummy();
        let mut clock = clock::create_for_testing(&mut ctx);

        managed_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[],
            vector[DefaultRemoteChain],
            vector[vector[DefaultRemotePool]],
            vector[DefaultRemoteToken],
        );

        managed_token_pool::set_chain_rate_limiter_config(
            &mut pool_state,
            &owner_cap,
            &clock,
            DefaultRemoteChain,
            true,
            10000,
            1000,
            true,
            10000,
            1000,
        );

        clock.increment_for_testing(1000000);

        // Clean up objects
        clock.destroy_for_testing();
        transfer::public_transfer(source_transfer_cap, @managed_token_pool);
        transfer::public_transfer(owner_cap, @managed_token_pool);
        test_scenario::return_shared(pool_state);
    };

    // Test lock_or_burn - use the address that has the mint cap
    scenario.next_tx(@0x456);
    {
        let mut pool_state = scenario.take_shared<
            ManagedTokenPoolState<MANAGED_TOKEN_POOL_TESTS>,
        >();
        let mut token_state = scenario.take_shared<TokenState<MANAGED_TOKEN_POOL_TESTS>>();
        let user_mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_POOL_TESTS>>();
        let deny_list = scenario.take_shared<DenyList>();
        let mut ctx = sui::tx_context::dummy();
        let mut clock = clock::create_for_testing(&mut ctx);
        clock.increment_for_testing(1000000000);

        // Mint tokens for burning
        let test_coin = managed_token::mint(
            &mut token_state,
            &user_mint_cap,
            &deny_list,
            1000,
            @0x456,
            &mut ctx,
        );

        let initial_coin_value = test_coin.value();
        assert!(initial_coin_value == 1000);

        let mut token_transfer_params = onramp_sh::create_token_transfer_params(
            bcs::to_bytes(&@0x456),
        ); // Use the test user address as token receiver

        // Actually call lock_or_burn function
        managed_token_pool::lock_or_burn<MANAGED_TOKEN_POOL_TESTS>(
            &ccip_ref, // ref parameter
            &mut token_transfer_params,
            test_coin, // c parameter (coin)
            DefaultRemoteChain, // remote_chain_selector parameter
            &clock, // clock parameter
            &deny_list, // deny_list parameter
            &mut token_state, // token_state parameter
            &mut pool_state, // state parameter
            &mut ctx, // context parameter
        );

        // Clean up token params
        let source_transfer_cap = scenario.take_from_address<onramp_sh::SourceTransferCap>(
            @managed_token_pool,
        );

        // Verify transfer data
        let (
            chain_selector,
            token_pool_package_id,
            amount,
            _,
            dest_token_address,
            extra_data,
        ) = onramp_sh::get_source_token_transfer_data(&token_transfer_params);
        assert!(chain_selector == DefaultRemoteChain);
        assert!(token_pool_package_id == @managed_token_pool);
        assert!(amount == initial_coin_value);
        // Note: source_pool should be the token address (coin metadata address), not the package address
        // This is different from the burn mint token pool due to different implementation
        assert!(dest_token_address == DefaultRemoteToken);
        assert!(extra_data.length() > 0); // Should contain encoded decimals

        onramp_sh::deconstruct_token_params(&source_transfer_cap, token_transfer_params);

        clock.destroy_for_testing();
        transfer::public_transfer(source_transfer_cap, @managed_token_pool);
        transfer::public_transfer(user_mint_cap, @0x456);
        test_scenario::return_shared(pool_state);
        test_scenario::return_shared(token_state);
        test_scenario::return_shared(deny_list);
    };

    transfer::public_freeze_object(coin_metadata);
    test_scenario::return_shared(ccip_ref);
    scenario.end();
}

#[test]
public fun test_release_or_mint_functionality() {
    let mut scenario = test_scenario::begin(@managed_token_pool);

    scenario.next_tx(@0x0);
    deny_list::create_for_testing(scenario.ctx());

    // Setup CCIP environment
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);

    // Create managed token and pool
    scenario.next_tx(@managed_token_pool);
    offramp_sh::test_init(scenario.ctx()); // Initialize offramp state helper
    let coin_metadata = {
        let ctx = scenario.ctx();
        let (treasury_cap, coin_metadata) = coin::create_currency(
            MANAGED_TOKEN_POOL_TESTS {},
            Decimals,
            b"MTPT",
            b"ManagedTokenPoolTest",
            b"managed_token_pool_test",
            option::none(),
            ctx,
        );

        managed_token::initialize(
            treasury_cap,
            package::test_claim(MANAGED_TOKEN_POOL_TESTS {}, ctx),
            ctx,
        );
        coin_metadata
    };

    scenario.next_tx(@managed_token_pool);
    {
        let mut token_state = scenario.take_shared<TokenState<MANAGED_TOKEN_POOL_TESTS>>();
        let token_owner_cap = scenario.take_from_sender<TokenOwnerCap<MANAGED_TOKEN_POOL_TESTS>>();

        // Create mint cap for the pool
        managed_token::configure_new_minter(
            &mut token_state,
            &token_owner_cap,
            @managed_token_pool,
            0,
            true,
            scenario.ctx(),
        );

        scenario.return_to_sender(token_owner_cap);
        test_scenario::return_shared(token_state);
    };

    scenario.next_tx(@managed_token_pool);
    let coin_metadata_address = object::id_to_address(&object::id(&coin_metadata));
    {
        // Call test_init to create owner_cap
        managed_token_pool::test_init(scenario.ctx());
    };

    scenario.next_tx(@managed_token_pool);
    {
        let mut owner_cap = scenario.take_from_sender<OwnerCap>();
        let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_POOL_TESTS>>();
        let token_state = scenario.take_shared<TokenState<MANAGED_TOKEN_POOL_TESTS>>();
        let token_owner_cap = scenario.take_from_sender<TokenOwnerCap<MANAGED_TOKEN_POOL_TESTS>>();

        managed_token_pool::initialize_with_managed_token(
            &mut owner_cap,
            &mut ccip_ref,
            &token_state,
            &token_owner_cap,
            &coin_metadata,
            mint_cap,
            @managed_token_pool,
            scenario.ctx(),
        );

        scenario.return_to_sender(owner_cap);
        scenario.return_to_sender(token_owner_cap);
        test_scenario::return_shared(token_state);
    };

    transfer::public_transfer(ccip_owner_cap, @0x0);

    // Set up chain and rate limiting
    scenario.next_tx(@managed_token_pool);
    {
        let mut pool_state = scenario.take_shared<
            ManagedTokenPoolState<MANAGED_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap>();
        let dest_transfer_cap = scenario.take_from_sender<offramp_sh::DestTransferCap>();
        let mut ctx = sui::tx_context::dummy();
        let mut clock = clock::create_for_testing(&mut ctx);

        managed_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[],
            vector[DefaultRemoteChain],
            vector[vector[DefaultRemotePool]],
            vector[DefaultRemoteToken],
        );

        managed_token_pool::set_chain_rate_limiter_config(
            &mut pool_state,
            &owner_cap,
            &clock,
            DefaultRemoteChain,
            true,
            10000,
            1000,
            true,
            10000,
            1000,
        );

        clock.increment_for_testing(1000000);

        // Clean up objects
        clock.destroy_for_testing();
        transfer::public_transfer(dest_transfer_cap, @managed_token_pool);
        transfer::public_transfer(owner_cap, @managed_token_pool);
        test_scenario::return_shared(pool_state);
    };

    // Test release_or_mint
    scenario.next_tx(@managed_token_pool);
    {
        let mut pool_state = scenario.take_shared<
            ManagedTokenPoolState<MANAGED_TOKEN_POOL_TESTS>,
        >();
        let mut token_state = scenario.take_shared<TokenState<MANAGED_TOKEN_POOL_TESTS>>();
        let dest_transfer_cap = scenario.take_from_address<offramp_sh::DestTransferCap>(
            @managed_token_pool,
        );
        let deny_list = scenario.take_shared<DenyList>();
        let mut ctx = sui::tx_context::dummy();
        let mut clock = clock::create_for_testing(&mut ctx);
        clock.increment_for_testing(1000000000);

        // Create receiver params for release_or_mint
        let mut receiver_params = offramp_sh::create_receiver_params(
            &dest_transfer_cap,
            DefaultRemoteChain,
        );

        // Add token transfer to receiver params
        let receiver_address = @0x789;
        let source_amount = 5000;
        let source_pool_data = x"0000000000000000000000000000000000000000000000000000000000000008"; // 8 decimals encoded
        let offchain_data = vector[];

        offramp_sh::add_dest_token_transfer(
            &dest_transfer_cap,
            &mut receiver_params,
            receiver_address,
            DefaultRemoteChain, // remote_chain_selector
            source_amount,
            coin_metadata_address,
            @managed_token_pool, // Package address, not token address
            DefaultRemotePool,
            source_pool_data,
            offchain_data,
        );

        // Verify the operation setup
        let source_chain = offramp_sh::get_source_chain_selector(&receiver_params);
        assert!(source_chain == DefaultRemoteChain);

        // Actually call release_or_mint function
        managed_token_pool::release_or_mint(
            &ccip_ref,
            &mut receiver_params,
            &clock,
            &deny_list,
            &mut token_state,
            &mut pool_state,
            &mut ctx,
        );

        // Clean up receiver params
        offramp_sh::deconstruct_receiver_params(&dest_transfer_cap, receiver_params);

        clock.destroy_for_testing();
        transfer::public_transfer(dest_transfer_cap, @managed_token_pool);
        test_scenario::return_shared(pool_state);
        test_scenario::return_shared(token_state);
        test_scenario::return_shared(deny_list);
    };

    // Verify the minted coin was transferred to the receiver
    scenario.next_tx(@0x789);
    {
        let minted_coin = scenario.take_from_sender<coin::Coin<MANAGED_TOKEN_POOL_TESTS>>();

        // Verify the minted amount (should be same as source amount due to same decimals)
        let minted_amount = minted_coin.value();
        assert!(minted_amount == 5000);

        // Transfer the coin back to clean up
        transfer::public_transfer(minted_coin, @managed_token_pool);
    };

    transfer::public_freeze_object(coin_metadata);
    test_scenario::return_shared(ccip_ref);
    scenario.end();
}

#[test]
#[expected_failure(abort_code = managed_token_pool::EInvalidOwnerCap)]
public fun test_invalid_owner_cap_error() {
    let mut scenario = test_scenario::begin(@managed_token_pool);

    // Create first pool
    let (ccip_owner_cap, ccip_ref, coin_metadata, deny_list) = setup_basic_pool(&mut scenario);

    // Create a second, independent pool setup to get a different owner cap
    scenario.next_tx(@0x999);
    let ctx = scenario.ctx();
    state_object::test_init(ctx);

    scenario.next_tx(@0x999);
    let ccip_owner_cap2 = scenario.take_from_sender<CCIPOwnerCap>();
    let mut ccip_ref2 = scenario.take_shared<CCIPObjectRef>();
    upgrade_registry::initialize(&mut ccip_ref2, &ccip_owner_cap2, scenario.ctx());
    token_admin_registry::initialize(&mut ccip_ref2, &ccip_owner_cap2, scenario.ctx());
    rmn_remote::initialize(&mut ccip_ref2, &ccip_owner_cap2, 1000, scenario.ctx());

    let (treasury_cap2, coin_metadata2) = coin::create_currency(
        MANAGED_TOKEN_POOL_TESTS {},
        Decimals,
        b"MTPT2",
        b"ManagedTokenPoolTest2",
        b"managed_token_pool_test2",
        option::none(),
        scenario.ctx(),
    );
    managed_token::initialize(
        treasury_cap2,
        package::test_claim(MANAGED_TOKEN_POOL_TESTS {}, scenario.ctx()),
        scenario.ctx(),
    );

    scenario.next_tx(@0x999);
    {
        let mut token_state2 = scenario.take_shared<TokenState<MANAGED_TOKEN_POOL_TESTS>>();
        let token_owner_cap2 = scenario.take_from_sender<TokenOwnerCap<MANAGED_TOKEN_POOL_TESTS>>();

        managed_token::configure_new_minter(
            &mut token_state2,
            &token_owner_cap2,
            @0x999,
            0,
            true,
            scenario.ctx(),
        );

        scenario.return_to_sender(token_owner_cap2);
        test_scenario::return_shared(token_state2);
    };

    scenario.next_tx(@0x999);
    {
        // Call test_init to create owner_cap
        managed_token_pool::test_init(scenario.ctx());
    };

    scenario.next_tx(@0x999);
    {
        let mut owner_cap = scenario.take_from_sender<OwnerCap>();
        let mint_cap2 = scenario.take_from_sender<MintCap<MANAGED_TOKEN_POOL_TESTS>>();
        let token_state2 = scenario.take_shared<TokenState<MANAGED_TOKEN_POOL_TESTS>>();
        let token_owner_cap2 = scenario.take_from_sender<TokenOwnerCap<MANAGED_TOKEN_POOL_TESTS>>();

        managed_token_pool::initialize_with_managed_token(
            &mut owner_cap,
            &mut ccip_ref2,
            &token_state2,
            &token_owner_cap2,
            &coin_metadata2,
            mint_cap2,
            @0x999,
            scenario.ctx(),
        );

        transfer::public_transfer(owner_cap, @0x999);
        scenario.return_to_sender(token_owner_cap2);
        test_scenario::return_shared(token_state2);
    };

    // Now test with mismatched owner caps
    scenario.next_tx(@managed_token_pool);
    let mut pool_state1 = scenario.take_shared<ManagedTokenPoolState<MANAGED_TOKEN_POOL_TESTS>>();
    let correct_owner_cap = scenario.take_from_sender<OwnerCap>();

    scenario.next_tx(@0x999);
    let wrong_owner_cap = scenario.take_from_sender<OwnerCap>();

    // First add the chain using the correct owner cap so the chain exists
    managed_token_pool::apply_chain_updates(
        &mut pool_state1,
        &correct_owner_cap,
        vector[], // remove none
        vector[DefaultRemoteChain], // add chain
        vector[vector[DefaultRemotePool]], // with pool
        vector[DefaultRemoteToken], // and token
    );

    // Now try to add a remote pool using wrong owner cap - should fail with EInvalidOwnerCap
    managed_token_pool::add_remote_pool(
        &mut pool_state1,
        &wrong_owner_cap, // Wrong owner cap from different pool
        DefaultRemoteChain,
        b"another_pool",
    );

    // This should not be reached due to expected failure
    transfer::public_transfer(correct_owner_cap, @managed_token_pool);
    transfer::public_transfer(wrong_owner_cap, @0x999);
    test_scenario::return_shared(pool_state1);

    // Cleanup (should not be reached)
    transfer::public_freeze_object(coin_metadata);
    transfer::public_freeze_object(coin_metadata2);
    transfer::public_transfer(ccip_owner_cap, @0x0);
    transfer::public_transfer(ccip_owner_cap2, @0x0);
    test_scenario::return_shared(ccip_ref);
    test_scenario::return_shared(ccip_ref2);
    test_scenario::return_shared(deny_list);
    scenario.end();
}

#[test]
#[expected_failure(abort_code = managed_token_pool::EInvalidArguments)]
public fun test_invalid_arguments_error() {
    let mut scenario = test_scenario::begin(@managed_token_pool);
    let (ccip_owner_cap, ccip_ref, coin_metadata, deny_list) = setup_basic_pool(&mut scenario);

    scenario.next_tx(@managed_token_pool);
    {
        let mut pool_state = scenario.take_shared<
            ManagedTokenPoolState<MANAGED_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap>();
        let mut ctx = sui::tx_context::dummy();
        let clock = clock::create_for_testing(&mut ctx);

        // Test set_chain_rate_limiter_configs with mismatched vector lengths
        // This should fail with EInvalidArguments
        managed_token_pool::set_chain_rate_limiter_configs(
            &mut pool_state,
            &owner_cap,
            &clock,
            vector[DefaultRemoteChain, 3000], // 2 chains
            vector[true], // 1 outbound enabled (mismatch!)
            vector[1000, 2000], // 2 capacities
            vector[100, 200], // 2 rates
            vector[true, false], // 2 inbound enabled
            vector[1500, 2500], // 2 capacities
            vector[150, 250], // 2 rates
        );

        // This should not be reached due to expected failure
        test_scenario::return_shared(pool_state);
        transfer::public_transfer(owner_cap, @managed_token_pool);
        clock.destroy_for_testing();
    };

    cleanup_test(scenario, ccip_owner_cap, ccip_ref, coin_metadata, deny_list);
}

#[test]
public fun test_edge_cases_and_comprehensive_coverage() {
    let mut scenario = test_scenario::begin(@managed_token_pool);
    let (ccip_owner_cap, ccip_ref, coin_metadata, deny_list) = setup_basic_pool(&mut scenario);

    scenario.next_tx(@managed_token_pool);
    {
        let mut pool_state = scenario.take_shared<
            ManagedTokenPoolState<MANAGED_TOKEN_POOL_TESTS>,
        >();
        let owner_cap = scenario.take_from_sender<OwnerCap>();
        let mut ctx = sui::tx_context::dummy();
        let clock = clock::create_for_testing(&mut ctx);

        // Test edge case: empty chain updates (should work)
        managed_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[], // remove none
            vector[], // add none
            vector[], // no pools
            vector[], // no tokens
        );

        // Verify no chains were added
        let supported_chains = managed_token_pool::get_supported_chains(&pool_state);
        assert!(supported_chains.length() == 0);

        // Test edge case: apply_allowlist_updates with empty vectors
        managed_token_pool::apply_allowlist_updates(
            &mut pool_state,
            &owner_cap,
            vector[], // remove none
            vector[], // add none
        );

        // Test adding multiple chains at once
        let chain1 = 1000;
        let chain2 = 2000;
        let chain3 = 3000;

        managed_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[],
            vector[chain1, chain2, chain3],
            vector[
                vector[b"pool1a", b"pool1b"],
                vector[b"pool2a"],
                vector[b"pool3a", b"pool3b", b"pool3c"],
            ],
            vector[b"token1", b"token2", b"token3"],
        );

        // Verify all chains were added
        let supported_chains = managed_token_pool::get_supported_chains(&pool_state);
        assert!(supported_chains.length() == 3);
        assert!(managed_token_pool::is_supported_chain(&pool_state, chain1));
        assert!(managed_token_pool::is_supported_chain(&pool_state, chain2));
        assert!(managed_token_pool::is_supported_chain(&pool_state, chain3));

        // Test remote pool queries
        let pools_chain1 = managed_token_pool::get_remote_pools(&pool_state, chain1);
        assert!(pools_chain1.length() == 2);
        assert!(managed_token_pool::is_remote_pool(&pool_state, chain1, b"pool1a"));
        assert!(managed_token_pool::is_remote_pool(&pool_state, chain1, b"pool1b"));
        assert!(!managed_token_pool::is_remote_pool(&pool_state, chain1, b"nonexistent"));

        let pools_chain3 = managed_token_pool::get_remote_pools(&pool_state, chain3);
        assert!(pools_chain3.length() == 3);

        // Test remote token queries
        let token1 = managed_token_pool::get_remote_token(&pool_state, chain1);
        let token2 = managed_token_pool::get_remote_token(&pool_state, chain2);
        let token3 = managed_token_pool::get_remote_token(&pool_state, chain3);
        assert!(token1 == b"token1");
        assert!(token2 == b"token2");
        assert!(token3 == b"token3");

        // Test removing some chains
        managed_token_pool::apply_chain_updates(
            &mut pool_state,
            &owner_cap,
            vector[chain2], // remove chain2
            vector[], // add none
            vector[], // no pools
            vector[], // no tokens
        );

        // Verify chain2 was removed
        let supported_chains_after = managed_token_pool::get_supported_chains(&pool_state);
        assert!(supported_chains_after.length() == 2);
        assert!(!managed_token_pool::is_supported_chain(&pool_state, chain2));
        assert!(managed_token_pool::is_supported_chain(&pool_state, chain1));
        assert!(managed_token_pool::is_supported_chain(&pool_state, chain3));

        // Test single chain rate limiter with edge values
        managed_token_pool::set_chain_rate_limiter_config(
            &mut pool_state,
            &owner_cap,
            &clock,
            chain1,
            false,
            0,
            0, // disabled outbound with zero values
            true,
            18446744073709551615u64,
            1000000000u64, // enabled inbound with max values
        );

        test_scenario::return_shared(pool_state);
        transfer::public_transfer(owner_cap, @managed_token_pool);
        clock.destroy_for_testing();
    };

    cleanup_test(scenario, ccip_owner_cap, ccip_ref, coin_metadata, deny_list);
}

#[test]
public fun test_initialize_with_managed_token_function() {
    let mut scenario = test_scenario::begin(@managed_token_pool);

    // Setup CCIP environment
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(&mut scenario);
    let deny_list = scenario.take_shared<DenyList>();

    // Create managed token first
    scenario.next_tx(@managed_token_pool);
    let coin_metadata = {
        let ctx = scenario.ctx();
        let (treasury_cap, coin_metadata) = coin::create_currency(
            MANAGED_TOKEN_POOL_TESTS {},
            Decimals,
            b"MTPT",
            b"ManagedTokenPoolTest",
            b"managed_token_pool_test",
            option::none(),
            ctx,
        );

        managed_token::initialize(
            treasury_cap,
            package::test_claim(MANAGED_TOKEN_POOL_TESTS {}, ctx),
            ctx,
        );
        coin_metadata
    };

    // Create mint cap for the token pool
    scenario.next_tx(@managed_token_pool);
    {
        let mut token_state = scenario.take_shared<TokenState<MANAGED_TOKEN_POOL_TESTS>>();
        let token_owner_cap = scenario.take_from_sender<TokenOwnerCap<MANAGED_TOKEN_POOL_TESTS>>();

        // Create mint cap for the pool
        managed_token::configure_new_minter(
            &mut token_state,
            &token_owner_cap,
            @managed_token_pool,
            0,
            true,
            scenario.ctx(),
        );

        scenario.return_to_sender(token_owner_cap);
        test_scenario::return_shared(token_state);
    };

    // Test the new initialize_with_managed_token function
    scenario.next_tx(@managed_token_pool);
    {
        // Call test_init to create owner_cap
        managed_token_pool::test_init(scenario.ctx());
    };

    scenario.next_tx(@managed_token_pool);
    {
        let mut owner_cap = scenario.take_from_sender<OwnerCap>();
        let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_POOL_TESTS>>();
        let token_state = scenario.take_shared<TokenState<MANAGED_TOKEN_POOL_TESTS>>();
        let token_owner_cap = scenario.take_from_sender<TokenOwnerCap<MANAGED_TOKEN_POOL_TESTS>>();

        // Use the new function that takes managed token state and owner cap directly
        managed_token_pool::initialize_with_managed_token(
            &mut owner_cap,
            &mut ccip_ref,
            &token_state,
            &token_owner_cap,
            &coin_metadata,
            mint_cap,
            @managed_token_pool, // token pool administrator
            scenario.ctx(),
        );

        transfer::public_transfer(owner_cap, @managed_token_pool);
        scenario.return_to_sender(token_owner_cap);
        test_scenario::return_shared(token_state);
    };

    // Verify the pool was initialized correctly
    scenario.next_tx(@managed_token_pool);
    {
        let pool_state = scenario.take_shared<ManagedTokenPoolState<MANAGED_TOKEN_POOL_TESTS>>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        // Test basic getters to ensure initialization worked
        assert!(managed_token_pool::get_token_decimals(&pool_state) == Decimals);
        let token_address = managed_token_pool::get_token(&pool_state);
        assert!(token_address != @0x0);

        // Verify owner is correct
        assert!(managed_token_pool::owner(&pool_state) == @managed_token_pool);

        // Test supported chains (should be empty initially)
        let supported_chains = managed_token_pool::get_supported_chains(&pool_state);
        assert!(supported_chains.length() == 0);

        // Test allowlist (should be disabled initially)
        assert!(!managed_token_pool::get_allowlist_enabled(&pool_state));
        let allowlist = managed_token_pool::get_allowlist(&pool_state);
        assert!(allowlist.length() == 0);

        scenario.return_to_sender(owner_cap);
        test_scenario::return_shared(pool_state);
    };

    // Verify the token pool is registered in token admin registry
    scenario.next_tx(@managed_token_pool);
    {
        // Calculate the actual package ID from TypeProof (same as initialization)
        let type_proof_type_name = type_name::with_defining_ids<managed_token_pool::TypeProof>();
        let type_proof_type_name_address = type_proof_type_name.address_string();
        let actual_package_id = address::from_ascii_bytes(
            &type_proof_type_name_address.into_bytes(),
        );

        let coin_metadata_address = object::id_to_address(&object::id(&coin_metadata));
        let pool_address = token_admin_registry::get_pool(&ccip_ref, coin_metadata_address);
        assert!(pool_address == actual_package_id); // Should match the dynamically calculated package id

        let (
            pool_package_id,
            pool_module,
            token_type,
            admin,
            pending_admin,
            type_proof,
            _lock_or_burn_params,
            _release_or_mint_params,
        ) = token_admin_registry::get_token_config_data(
            &ccip_ref,
            coin_metadata_address,
        );

        assert!(pool_package_id == actual_package_id);
        assert!(pool_module == string::utf8(b"managed_token_pool"));
        assert!(
            token_type == type_name::with_defining_ids<MANAGED_TOKEN_POOL_TESTS>().into_string(),
        );
        assert!(admin == @managed_token_pool);
        assert!(pending_admin == @0x0);
        // type_proof should be the TypeProof type name - we just check it's not empty
        assert!(type_proof.length() > 0);
    };

    cleanup_test(scenario, ccip_owner_cap, ccip_ref, coin_metadata, deny_list);
}

// Helper functions
fun setup_basic_pool(
    scenario: &mut test_scenario::Scenario,
): (CCIPOwnerCap, CCIPObjectRef, coin::CoinMetadata<MANAGED_TOKEN_POOL_TESTS>, DenyList) {
    let (ccip_owner_cap, mut ccip_ref) = setup_ccip_environment(scenario);

    scenario.next_tx(@managed_token_pool);
    let coin_metadata = {
        let ctx = scenario.ctx();
        let (treasury_cap, coin_metadata) = coin::create_currency(
            MANAGED_TOKEN_POOL_TESTS {},
            Decimals,
            b"MTPT",
            b"ManagedTokenPoolTest",
            b"managed_token_pool_test",
            option::none(),
            ctx,
        );

        managed_token::initialize(
            treasury_cap,
            package::test_claim(MANAGED_TOKEN_POOL_TESTS {}, ctx),
            ctx,
        );
        coin_metadata
    };

    scenario.next_tx(@managed_token_pool);
    {
        let mut token_state = scenario.take_shared<TokenState<MANAGED_TOKEN_POOL_TESTS>>();
        let token_owner_cap = scenario.take_from_sender<TokenOwnerCap<MANAGED_TOKEN_POOL_TESTS>>();

        managed_token::configure_new_minter(
            &mut token_state,
            &token_owner_cap,
            @managed_token_pool,
            0,
            true,
            scenario.ctx(),
        );

        managed_token::configure_new_minter(
            &mut token_state,
            &token_owner_cap,
            @0x456,
            1000000,
            false,
            scenario.ctx(),
        );

        scenario.return_to_sender(token_owner_cap);
        test_scenario::return_shared(token_state);
    };

    scenario.next_tx(@managed_token_pool);
    {
        // Call test_init to create owner_cap
        managed_token_pool::test_init(scenario.ctx());
    };

    scenario.next_tx(@managed_token_pool);
    {
        let mut owner_cap = scenario.take_from_sender<OwnerCap>();
        let mint_cap = scenario.take_from_sender<MintCap<MANAGED_TOKEN_POOL_TESTS>>();
        let token_state = scenario.take_shared<TokenState<MANAGED_TOKEN_POOL_TESTS>>();
        let token_owner_cap = scenario.take_from_sender<TokenOwnerCap<MANAGED_TOKEN_POOL_TESTS>>();

        managed_token_pool::initialize_with_managed_token(
            &mut owner_cap,
            &mut ccip_ref,
            &token_state,
            &token_owner_cap,
            &coin_metadata,
            mint_cap,
            @managed_token_pool,
            scenario.ctx(),
        );

        transfer::public_transfer(owner_cap, @managed_token_pool);
        scenario.return_to_sender(token_owner_cap);
        test_scenario::return_shared(token_state);
    };

    let deny_list = scenario.take_shared<DenyList>();

    (ccip_owner_cap, ccip_ref, coin_metadata, deny_list)
}

fun cleanup_test(
    scenario: test_scenario::Scenario,
    ccip_owner_cap: CCIPOwnerCap,
    ccip_ref: CCIPObjectRef,
    coin_metadata: coin::CoinMetadata<MANAGED_TOKEN_POOL_TESTS>,
    deny_list: DenyList,
) {
    transfer::public_freeze_object(coin_metadata);
    transfer::public_transfer(ccip_owner_cap, @0x0);
    test_scenario::return_shared(ccip_ref);
    test_scenario::return_shared(deny_list);
    scenario.end();
}

// ================================================================
// |             Rate Limiter calculate_refill Tests             |
// ================================================================

#[test]
public fun test_calculate_refill_at_capacity() {
    use managed_token_pool::rate_limiter;

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
    use managed_token_pool::rate_limiter;

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
    use managed_token_pool::rate_limiter;

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
    use managed_token_pool::rate_limiter;

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
    use managed_token_pool::rate_limiter;

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
    use managed_token_pool::rate_limiter;

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
    use managed_token_pool::rate_limiter;

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
    use managed_token_pool::rate_limiter;

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
    use managed_token_pool::rate_limiter;

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
    use managed_token_pool::rate_limiter;

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
