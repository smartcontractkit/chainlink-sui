#[test_only]
module ccip_token_pool::token_pool_test;

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
    let source_pool_data = vector[0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0];

    let _decimal = token_pool::parse_remote_decimals(source_pool_data, 1);
}
