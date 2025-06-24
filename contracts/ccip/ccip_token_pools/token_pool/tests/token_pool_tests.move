#[test_only]
module ccip_token_pool::token_pool_test;

use std::bcs;
use sui::clock;
use sui::coin;
use sui::test_scenario::{Self, Scenario};

use ccip_token_pool::token_pool::{Self, TokenPoolState};

public struct TOKEN_POOL_TEST has drop {}

const Decimals: u8 = 8;
const DefaultRemoteChain: u64 = 2000;
const NewRemoteChain: u64 = 3000;
const DefaultRemoteToken: vector<u8> = b"default_remote_token";
const DefaultRemotePool: vector<u8> = b"default_remote_pool";
const NewRemotePool: vector<u8> = b"new_remote_pool";
const AddAddress1: address = @0x1;
const AddAddress2: address = @0x2;
const RemoveAddress3: address = @0x3;

fun set_up_test(): (Scenario, TokenPoolState) {
    let mut scenario = test_scenario::begin(@ccip_token_pool);
    let ctx = scenario.ctx();

    let (treasury_cap, coin_metadata) = coin::create_currency(
        TOKEN_POOL_TEST {},
        Decimals,
        b"TEST",
        b"TestToken",
        b"test_token",
        option::none(),
        ctx
    );

    let mut state = token_pool::initialize(object::id_to_address(&object::id(&coin_metadata)), Decimals, vector[], ctx);

    // Set state in the pool
    set_up_default_remote_chain(&mut state);

    transfer::public_freeze_object(coin_metadata);
    transfer::public_transfer(treasury_cap, ctx.sender());

    (scenario, state)
}

fun set_up_token_pool_test_with_ccip(): (Scenario, ccip::state_object::OwnerCap, ccip::state_object::CCIPObjectRef, TokenPoolState) {
    let mut scenario = test_scenario::begin(@ccip_token_pool);
    let ctx = scenario.ctx();

    // Create CCIP object ref
    let (owner_cap, mut ref) = ccip::state_object::create(ctx);
    
    // Initialize RMN remote
    ccip::rmn_remote::initialize(&mut ref, &owner_cap, 1000, ctx); // local chain selector = 1000
    
    // Create token pool state
    let (treasury_cap, coin_metadata) = coin::create_currency(
        TOKEN_POOL_TEST {},
        Decimals,
        b"TEST",
        b"TestToken",
        b"test_token",
        option::none(),
        ctx
    );

    let mut state = token_pool::initialize(object::id_to_address(&object::id(&coin_metadata)), Decimals, vector[], ctx);

    // Set up default remote chain
    set_up_default_remote_chain(&mut state);

    transfer::public_freeze_object(coin_metadata);
    transfer::public_transfer(treasury_cap, ctx.sender());

    (scenario, owner_cap, ref, state)
}

fun set_up_default_remote_chain(state: &mut TokenPoolState) {
    token_pool::apply_chain_updates(
        state,
        vector[],
        vector[DefaultRemoteChain],
        vector[vector[DefaultRemotePool]],
        vector[DefaultRemoteToken]
    )
}

#[test]
public fun initialize_correctly_sets_state() {
    let (scenario, state) = set_up_test();

    assert!(token_pool::is_supported_chain(&state, DefaultRemoteChain));

    token_pool::destroy_token_pool(state);
    scenario.end();
}

#[test]
fun add_and_remote_remote_pool_existing_chain() {
    let (scenario, mut state) = set_up_test();

    assert!(!token_pool::is_remote_pool(&state, DefaultRemoteChain, NewRemotePool));
    assert!(token_pool::get_remote_pools(&state, DefaultRemoteChain).length() == 1);

    token_pool::add_remote_pool(&mut state, DefaultRemoteChain, NewRemotePool);

    assert!(token_pool::is_remote_pool(&state, DefaultRemoteChain, NewRemotePool));
    assert!(token_pool::get_remote_pools(&state, DefaultRemoteChain).length() == 2);
    assert!(token_pool::is_supported_chain(&state, DefaultRemoteChain), 1);

    token_pool::remove_remote_pool(&mut state, DefaultRemoteChain, NewRemotePool);
    assert!(!token_pool::is_remote_pool(&state, DefaultRemoteChain, NewRemotePool));
    assert!(token_pool::get_remote_pools(&state, DefaultRemoteChain).length() == 1);

    token_pool::destroy_token_pool(state);
    scenario.end();
}

#[test]
#[expected_failure(abort_code = token_pool::EUnknownRemoteChainSelector)]
fun remote_remote_pool_existing_chain() {
    let (scenario, mut state) = set_up_test();

    token_pool::remove_remote_pool(&mut state, NewRemoteChain, NewRemotePool);

    token_pool::destroy_token_pool(state);
    scenario.end();
}

#[test]
fun apply_chain_updates() {
    let (scenario, mut state) = set_up_test();
    let new_remote_chain = 3000;
    let new_remote_token = b"new_remote_token";
    let new_remote_pool_2 = b"new_remote_pool_2";

    assert!(!token_pool::is_supported_chain(&state, new_remote_chain));

    token_pool::apply_chain_updates(
        &mut state,
        vector[],
        vector[new_remote_chain],
        vector[vector[NewRemotePool, new_remote_pool_2]],
        vector[new_remote_token]
    );
    assert!(token_pool::is_supported_chain(&state, new_remote_chain));
    assert!(token_pool::get_remote_pools(&state, new_remote_chain).length() == 2);
    assert!(token_pool::get_remote_token(&state, new_remote_chain) == new_remote_token);

    token_pool::apply_chain_updates(
        &mut state,
        vector[new_remote_chain],
        vector[],
        vector[],
        vector[]
    );
    assert!(!token_pool::is_supported_chain(&state, new_remote_chain));

    token_pool::destroy_token_pool(state);
    scenario.end();
}

#[test]
fun test_calculate_local_amount_same_decimals() {
    // When remote and local decimals are the same, amount should not change
    let remote_amount: u256 = 1000000;
    let remote_decimals: u8 = 8;
    let local_decimals: u8 = 8;

    let local_amount =
        token_pool::calculate_local_amount(
            remote_amount, remote_decimals, local_decimals
        );
    assert!(local_amount == 1000000);
}

#[test]
fun test_calculate_local_amount_more_decimals() {
    // When local has more decimals, amount should increase
    let remote_amount: u256 = 1000000;
    let remote_decimals: u8 = 6; // 6 decimals
    let local_decimals: u8 = 8; // 8 decimals (2 more)

    let local_amount =
        token_pool::calculate_local_amount(
            remote_amount, remote_decimals, local_decimals
        );
    assert!(local_amount == 100000000); // 1000000 * 10^2
}

#[test]
fun test_calculate_local_amount_fewer_decimals() {
    // When local has fewer decimals, amount should decrease
    let remote_amount: u256 = 1000000;
    let remote_decimals: u8 = 8; // 8 decimals
    let local_decimals: u8 = 6; // 6 decimals (2 fewer)

    let local_amount =
        token_pool::calculate_local_amount(
            remote_amount, remote_decimals, local_decimals
        );
    assert!(local_amount == 10000); // 1000000 / 10^2
}

#[test]
#[expected_failure(abort_code = token_pool::EDecimalOverflow)]
fun test_decimal_overflow_protection() {
    // Test for overflow protection - when decimal difference exceeds MAX_SAFE_DECIMAL_DIFF
    let remote_amount: u256 = 1000000;
    let remote_decimals: u8 = 1; // 1 decimal
    let local_decimals: u8 = 100; // 100 decimals (99 more - exceeds the limit of 77)

    let _local_amount =
        token_pool::calculate_local_amount(
            remote_amount, remote_decimals, local_decimals
        );
}

#[test]
#[expected_failure(abort_code = token_pool::EInvalidEncodedAmount)]
fun test_local_amount_u64_overflow() {
    let remote_amount: u256 = 0xffffffffffffffffffffffffffffffff;
    let remote_decimals: u8 = 0;
    let local_decimals: u8 = 18;

    let _local_amount =
        token_pool::calculate_local_amount(
            remote_amount, remote_decimals, local_decimals
        );
}

#[test]
fun test_enable_and_update_and_disable_allowlist() {
    let (scenario, mut state) = set_up_test();

    assert!(!token_pool::get_allowlist_enabled(&state));
    token_pool::set_allowlist_enabled(&mut state, true);
    assert!(token_pool::get_allowlist_enabled(&state));

    let removes = vector[];
    let adds = vector[AddAddress1, AddAddress2];
    token_pool::apply_allowlist_updates(
        &mut state,
        removes,
        adds,
    );

    let allowlist = token_pool::get_allowlist(&state);
    assert!(allowlist == adds);

    let removes = vector[AddAddress1, RemoveAddress3];
    let adds = vector[];
    token_pool::apply_allowlist_updates(
        &mut state,
        removes,
        adds,
    );
    let allowlist = token_pool::get_allowlist(&state);
    assert!(allowlist == vector[AddAddress2]);

    token_pool::set_allowlist_enabled(&mut state, false);
    assert!(!token_pool::get_allowlist_enabled(&state));

    token_pool::destroy_token_pool(state);
    scenario.end();
}

#[test]
fun test_parse_remote_decimals() {
    let source_pool_data = vector[0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9];

    let decimal = token_pool::parse_remote_decimals(source_pool_data, 1);
    assert!(decimal == 9);
}

#[test]
#[expected_failure(abort_code = token_pool::EInvalidRemoteChainDecimals)]
fun test_parse_remote_decimals_overflow() {
    // 256
    let source_pool_data = x"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100";

    let _decimal = token_pool::parse_remote_decimals(source_pool_data, 1);
}

#[test]
fun test_set_chain_rate_limiter_config() {
    let mut ctx = tx_context::dummy();
    let clock = clock::create_for_testing(&mut ctx);
    let (scenario, mut state) = set_up_test();

    token_pool::set_chain_rate_limiter_config(
        &clock,
        &mut state,
        DefaultRemoteChain,
        true,
        2000,
        3000,
        true,
        4000,
        5000,
    );

    token_pool::destroy_token_pool(state);
    clock.destroy_for_testing();
    scenario.end();
}

#[test]
fun test_calculate_release_or_mint_amount() {
    let mut ctx = tx_context::dummy();
    let clock = clock::create_for_testing(&mut ctx);
    let (scenario, state) = set_up_test();

    let amount = token_pool::calculate_release_or_mint_amount(
        &state,
        x"0000000000000000000000000000000000000000000000000000000000000010",
        100000000000000,
    );

    // source decimals = 16, local decimals = 8, amount = 100000000000000
    assert!(amount == 1000000);

    token_pool::destroy_token_pool(state);
    clock.destroy_for_testing();
    scenario.end();
}

#[test]
#[expected_failure(abort_code = token_pool::ERemoateChainToAddMismatch)]
fun test_apply_chain_updates_pool_addresses_length_mismatch() {
    let (scenario, mut state) = set_up_test();

    // Test with mismatched lengths: 1 chain selector but 2 pool address vectors
    token_pool::apply_chain_updates(
        &mut state,
        vector[], // remove nothing
        vector[5000], // add 1 chain
        vector[vector[b"pool1"], vector[b"pool2"]], // but 2 pool address vectors
        vector[b"token1"] // and 1 token address
    );

    token_pool::destroy_token_pool(state);
    scenario.end();
}

#[test]
#[expected_failure(abort_code = token_pool::ERemoateChainToAddMismatch)]
fun test_apply_chain_updates_token_addresses_length_mismatch() {
    let (scenario, mut state) = set_up_test();

    // Test with mismatched lengths: 1 chain selector but 2 token addresses
    token_pool::apply_chain_updates(
        &mut state,
        vector[], // remove nothing
        vector[5000], // add 1 chain
        vector[vector[b"pool1"]], // 1 pool address vector
        vector[b"token1", b"token2"] // but 2 token addresses
    );

    token_pool::destroy_token_pool(state);
    scenario.end();
}

#[test]
fun test_validate_lock_or_burn_success() {
    let (scenario, owner_cap, ref, mut state) = set_up_token_pool_test_with_ccip();
    let mut ctx = tx_context::dummy();
    let mut clock = clock::create_for_testing(&mut ctx);
    token_pool::set_chain_rate_limiter_config(
        &clock,
        &mut state,
        DefaultRemoteChain,
        true,
        3000,
        3000,
        true,
        4000,
        4000,
    );
    // we need to advance the clock in order for the bucket in the rate limiter to be updated
    clock.increment_for_testing(1000000);

    let sender = @0x123;
    let remote_chain_selector = DefaultRemoteChain;
    let local_amount = 1000;

    // Should succeed with valid parameters
    let remote_token = token_pool::validate_lock_or_burn(
        &ref,
        &clock,
        &mut state,
        sender,
        remote_chain_selector,
        local_amount
    );

    assert!(remote_token == DefaultRemoteToken);

    token_pool::destroy_token_pool(state);
    ccip::state_object::destroy_owner_cap(owner_cap);
    ccip::state_object::destroy_state_object(ref);
    clock.destroy_for_testing();
    scenario.end();
}

#[test]
#[expected_failure(abort_code = token_pool::ECursedChain)]
fun test_validate_lock_or_burn_cursed_chain() {
    let (scenario, owner_cap,mut ref, mut state) = set_up_token_pool_test_with_ccip();
    let mut ctx = tx_context::dummy();
    let clock = clock::create_for_testing(&mut ctx);
    
    // Curse the remote chain
    ccip::rmn_remote::curse_multiple(
        &mut ref,
        &owner_cap,
        vector[x"00000000000000000000000000000000"] // curse chain selector 0, but we'll test with DefaultRemoteChain
    );
    
    // Actually curse the DefaultRemoteChain (2000)
    let mut subject_bytes = bcs::to_bytes(&(DefaultRemoteChain as u128));
    subject_bytes.reverse();
    ccip::rmn_remote::curse(&mut ref, &owner_cap, subject_bytes);
    
    let sender = @0x123;
    let local_amount = 1000;

    // Should fail with cursed chain
    let _remote_token = token_pool::validate_lock_or_burn(
        &ref,
        &clock,
        &mut state,
        sender,
        DefaultRemoteChain,
        local_amount
    );

    token_pool::destroy_token_pool(state);
    ccip::state_object::destroy_owner_cap(owner_cap);
    ccip::state_object::destroy_state_object(ref);
    clock.destroy_for_testing();
    scenario.end();
}

#[test]
#[expected_failure(abort_code = token_pool::ENotPublisher)]
fun test_validate_lock_or_burn_allowlist_not_allowed() {
    let (scenario, owner_cap, ref, mut state) = set_up_token_pool_test_with_ccip();
    let mut ctx = tx_context::dummy();
    let clock = clock::create_for_testing(&mut ctx);
    
    // Enable allowlist but don't add sender
    token_pool::set_allowlist_enabled(&mut state, true);
    token_pool::apply_allowlist_updates(&mut state, vector[], vector[@0x456]); // Add different address
    
    let sender = @0x123; // This address is not in allowlist
    let local_amount = 1000;

    // Should fail because sender is not in allowlist
    let _remote_token = token_pool::validate_lock_or_burn(
        &ref,
        &clock,
        &mut state,
        sender,
        DefaultRemoteChain,
        local_amount
    );

    token_pool::destroy_token_pool(state);
    ccip::state_object::destroy_owner_cap(owner_cap);
    ccip::state_object::destroy_state_object(ref);
    clock.destroy_for_testing();
    scenario.end();
}

#[test]
#[expected_failure(abort_code = token_pool::EUnknownRemoteChainSelector)]
fun test_validate_lock_or_burn_unknown_chain() {
    let (scenario, owner_cap, ref, mut state) = set_up_token_pool_test_with_ccip();
    let mut ctx = tx_context::dummy();
    let clock = clock::create_for_testing(&mut ctx);
    
    let sender = @0x123;
    let unknown_chain = 9999; // Chain not configured
    let local_amount = 1000;

    // Should fail with unknown chain selector
    let _remote_token = token_pool::validate_lock_or_burn(
        &ref,
        &clock,
        &mut state,
        sender,
        unknown_chain,
        local_amount
    );

    token_pool::destroy_token_pool(state);
    ccip::state_object::destroy_owner_cap(owner_cap);
    ccip::state_object::destroy_state_object(ref);
    clock.destroy_for_testing();
    scenario.end();
}

#[test]
fun test_validate_lock_or_burn_with_allowlist_success() {
    let (scenario, owner_cap, ref, mut state) = set_up_token_pool_test_with_ccip();
    let mut ctx = tx_context::dummy();
    let mut clock = clock::create_for_testing(&mut ctx);
    token_pool::set_chain_rate_limiter_config(
        &clock,
        &mut state,
        DefaultRemoteChain,
        true,
        3000,
        3000,
        true,
        4000,
        4000,
    );
    // we need to advance the clock in order for the bucket in the rate limiter to be updated
    clock.increment_for_testing(1000000);
    
    let sender = @0x123;
    
    // Enable allowlist and add sender
    token_pool::set_allowlist_enabled(&mut state, true);
    token_pool::apply_allowlist_updates(&mut state, vector[], vector[sender]);
    
    let local_amount = 1000;

    // Should succeed because sender is in allowlist
    let remote_token = token_pool::validate_lock_or_burn(
        &ref,
        &clock,
        &mut state,
        sender,
        DefaultRemoteChain,
        local_amount
    );

    assert!(remote_token == DefaultRemoteToken);

    token_pool::destroy_token_pool(state);
    ccip::state_object::destroy_owner_cap(owner_cap);
    ccip::state_object::destroy_state_object(ref);
    clock.destroy_for_testing();
    scenario.end();
}

#[test]
fun test_validate_release_or_mint_success() {
    let (scenario, owner_cap, ref, mut state) = set_up_token_pool_test_with_ccip();
    let mut ctx = tx_context::dummy();
    let mut clock = clock::create_for_testing(&mut ctx);
    
    // Set up rate limiter config
    token_pool::set_chain_rate_limiter_config(
        &clock,
        &mut state,
        DefaultRemoteChain,
        true,
        3000,
        3000,
        true,
        4000,
        4000,
    );
    clock.increment_for_testing(1000000);

    let dest_token_address = token_pool::get_token(&state);
    let local_amount = 1000;

    // Should succeed with valid parameters
    token_pool::validate_release_or_mint(
        &ref,
        &clock,
        &mut state,
        DefaultRemoteChain,
        dest_token_address,
        DefaultRemotePool,
        local_amount
    );

    token_pool::destroy_token_pool(state);
    ccip::state_object::destroy_owner_cap(owner_cap);
    ccip::state_object::destroy_state_object(ref);
    clock.destroy_for_testing();
    scenario.end();
}

#[test]
#[expected_failure(abort_code = token_pool::EUnknownToken)]
fun test_validate_release_or_mint_unknown_token() {
    let (scenario, owner_cap, ref, mut state) = set_up_token_pool_test_with_ccip();
    let mut ctx = tx_context::dummy();
    let clock = clock::create_for_testing(&mut ctx);

    let wrong_token_address = @0x999; // Wrong token address
    let local_amount = 1000;

    // Should fail with unknown token
    token_pool::validate_release_or_mint(
        &ref,
        &clock,
        &mut state,
        DefaultRemoteChain,
        wrong_token_address,
        DefaultRemotePool,
        local_amount
    );

    token_pool::destroy_token_pool(state);
    ccip::state_object::destroy_owner_cap(owner_cap);
    ccip::state_object::destroy_state_object(ref);
    clock.destroy_for_testing();
    scenario.end();
}

#[test]
#[expected_failure(abort_code = token_pool::ECursedChain)]
fun test_validate_release_or_mint_cursed_chain() {
    let (scenario, owner_cap, mut ref, mut state) = set_up_token_pool_test_with_ccip();
    let mut ctx = tx_context::dummy();
    let clock = clock::create_for_testing(&mut ctx);
    
    // Curse the remote chain
    let mut subject_bytes = bcs::to_bytes(&(DefaultRemoteChain as u128));
    subject_bytes.reverse();
    ccip::rmn_remote::curse(&mut ref, &owner_cap, subject_bytes);
    
    let dest_token_address = token_pool::get_token(&state);
    let local_amount = 1000;

    // Should fail with cursed chain
    token_pool::validate_release_or_mint(
        &ref,
        &clock,
        &mut state,
        DefaultRemoteChain,
        dest_token_address,
        DefaultRemotePool,
        local_amount
    );

    token_pool::destroy_token_pool(state);
    ccip::state_object::destroy_owner_cap(owner_cap);
    ccip::state_object::destroy_state_object(ref);
    clock.destroy_for_testing();
    scenario.end();
}

#[test]
#[expected_failure(abort_code = token_pool::EUnknownRemotePool)]
fun test_validate_release_or_mint_unknown_remote_pool() {
    let (scenario, owner_cap, ref, mut state) = set_up_token_pool_test_with_ccip();
    let mut ctx = tx_context::dummy();
    let clock = clock::create_for_testing(&mut ctx);
    
    let dest_token_address = token_pool::get_token(&state);
    let unknown_pool_address = b"unknown_pool"; // Pool not configured
    let local_amount = 1000;

    // Should fail with unknown remote pool
    token_pool::validate_release_or_mint(
        &ref,
        &clock,
        &mut state,
        DefaultRemoteChain,
        dest_token_address,
        unknown_pool_address,
        local_amount
    );

    token_pool::destroy_token_pool(state);
    ccip::state_object::destroy_owner_cap(owner_cap);
    ccip::state_object::destroy_state_object(ref);
    clock.destroy_for_testing();
    scenario.end();
}

#[test]
#[expected_failure(abort_code = token_pool::EUnknownRemoteChainSelector)]
fun test_validate_release_or_mint_unknown_chain() {
    let (scenario, owner_cap, ref, mut state) = set_up_token_pool_test_with_ccip();
    let mut ctx = tx_context::dummy();
    let clock = clock::create_for_testing(&mut ctx);
    
    let unknown_chain_selector = 9999; // Chain not configured
    let dest_token_address = token_pool::get_token(&state);
    let local_amount = 1000;

    // Should fail with unknown remote chain selector
    // This will fail at the is_remote_pool check
    token_pool::validate_release_or_mint(
        &ref,
        &clock,
        &mut state,
        unknown_chain_selector,
        dest_token_address,
        DefaultRemotePool,
        local_amount
    );

    token_pool::destroy_token_pool(state);
    ccip::state_object::destroy_owner_cap(owner_cap);
    ccip::state_object::destroy_state_object(ref);
    clock.destroy_for_testing();
    scenario.end();
}

#[test]
fun test_validate_release_or_mint_with_different_pool() {
    let (scenario, owner_cap, ref, mut state) = set_up_token_pool_test_with_ccip();
    let mut ctx = tx_context::dummy();
    let mut clock = clock::create_for_testing(&mut ctx);
    
    // Add another remote pool to the default chain
    let additional_pool = b"additional_pool";
    token_pool::add_remote_pool(&mut state, DefaultRemoteChain, additional_pool);
    
    // Set up rate limiter config
    token_pool::set_chain_rate_limiter_config(
        &clock,
        &mut state,
        DefaultRemoteChain,
        true,
        3000,
        3000,
        true,
        4000,
        4000,
    );
    clock.increment_for_testing(1000000);

    let dest_token_address = token_pool::get_token(&state);
    let local_amount = 1000;

    // Should succeed with the additional pool
    token_pool::validate_release_or_mint(
        &ref,
        &clock,
        &mut state,
        DefaultRemoteChain,
        dest_token_address,
        additional_pool,
        local_amount
    );

    token_pool::destroy_token_pool(state);
    ccip::state_object::destroy_owner_cap(owner_cap);
    ccip::state_object::destroy_state_object(ref);
    clock.destroy_for_testing();
    scenario.end();
}

#[test]
#[expected_failure(abort_code = ccip_token_pool::rate_limiter::ETokenMaxCapacityExceeded)]
fun test_validate_lock_or_burn_rate_limit_max_capacity_exceeded() {
    let (scenario, owner_cap, ref, mut state) = set_up_token_pool_test_with_ccip();
    let mut ctx = tx_context::dummy();
    let mut clock = clock::create_for_testing(&mut ctx);
    
    // Set up rate limiter with small capacity
    let small_capacity = 500;
    token_pool::set_chain_rate_limiter_config(
        &clock,
        &mut state,
        DefaultRemoteChain,
        true, // outbound enabled
        small_capacity, // small outbound capacity
        1000, // outbound rate
        true,
        4000,
        4000,
    );
    clock.increment_for_testing(1000000);
    
    let sender = @0x123;
    let local_amount = 1000; // Request more than capacity (500)

    // Should fail because requested amount exceeds max capacity
    let _remote_token = token_pool::validate_lock_or_burn(
        &ref,
        &clock,
        &mut state,
        sender,
        DefaultRemoteChain,
        local_amount
    );

    token_pool::destroy_token_pool(state);
    ccip::state_object::destroy_owner_cap(owner_cap);
    ccip::state_object::destroy_state_object(ref);
    clock.destroy_for_testing();
    scenario.end();
}

#[test]
#[expected_failure(abort_code = ccip_token_pool::rate_limiter::ETokenRateLimitReached)]
fun test_validate_lock_or_burn_rate_limit_reached() {
    let (scenario, owner_cap, ref, mut state) = set_up_token_pool_test_with_ccip();
    let mut ctx = tx_context::dummy();
    let mut clock = clock::create_for_testing(&mut ctx);
    
    // Set up rate limiter with low rate and capacity
    let capacity = 1000;
    let rate = 1; // Very low rate (1 token per second)
    token_pool::set_chain_rate_limiter_config(
        &clock,
        &mut state,
        DefaultRemoteChain,
        true, // outbound enabled
        capacity,
        rate,
        true,
        4000,
        4000,
    );
    
    // Advance clock to allow some tokens to accumulate
    clock.increment_for_testing(100_000); // 100 seconds = 100 tokens at rate 1/sec
    
    let sender = @0x123;
    
    // First request should succeed (consume 50 tokens)
    let _remote_token1 = token_pool::validate_lock_or_burn(
        &ref,
        &clock,
        &mut state,
        sender,
        DefaultRemoteChain,
        50
    );
    
    // Second request should succeed (consume another 50 tokens, total 100)
    let _remote_token2 = token_pool::validate_lock_or_burn(
        &ref,
        &clock,
        &mut state,
        sender,
        DefaultRemoteChain,
        50
    );
    
    // Third request should fail - trying to consume more tokens than available
    // Available tokens should be 0 after consuming 100 tokens
    let _remote_token3 = token_pool::validate_lock_or_burn(
        &ref,
        &clock,
        &mut state,
        sender,
        DefaultRemoteChain,
        1 // Even 1 token should fail
    );

    token_pool::destroy_token_pool(state);
    ccip::state_object::destroy_owner_cap(owner_cap);
    ccip::state_object::destroy_state_object(ref);
    clock.destroy_for_testing();
    scenario.end();
}
