module lock_release_token_pool::lock_release_token_pool;

use std::string::{Self, String};

use sui::clock::Clock;
use sui::bag::{Self, Bag};
use sui::coin::{Self, Coin, CoinMetadata, TreasuryCap};

use ccip::dynamic_dispatcher as dd;
use ccip::eth_abi;
use ccip::offramp_state_helper as osh;
use ccip::state_object::CCIPObjectRef;
use ccip::token_admin_registry;

use ccip_token_pool::token_pool::{Self, TokenPoolState};

public struct OwnerCap has key, store {
    id: UID,
}

// TODO: ownership model
public struct LockReleaseTokenPoolState has key {
    id: UID,
    // ownable_state: ownable::OwnableState,
    token_pool_state: TokenPoolState,
    coin_store: Bag, // use Bag to avoid type param, but it's also trivial to use a single Coin<T>
}

const COIN_STORE: vector<u8> = b"CoinStore";

const E_INVALID_COIN_METADATA: u64 = 1;
const E_INVALID_ARGUMENTS: u64 = 2;
const E_TOKEN_POOL_BALANCE_TOO_LOW: u64 = 3;

// ================================================================
// |                             Init                             |
// ================================================================

public fun type_and_version(): String {
    string::utf8(b"LockReleaseTokenPool 1.6.0")
}

#[allow(lint(self_transfer))]
public fun initialize<T: drop>(
    ref: &mut CCIPObjectRef,
    coin_metadata: &CoinMetadata<T>,
    treasury_cap: &TreasuryCap<T>,
    ctx: &mut TxContext,
) {
    let coin_metadata_address: address = object::id_to_address(&object::id(coin_metadata));
    assert!(
        coin_metadata_address == @lock_release_local_token,
        E_INVALID_COIN_METADATA
    );

    let mut lock_release_token_pool = LockReleaseTokenPoolState {
        id: object::new(ctx),
        token_pool_state: token_pool::initialize(coin_metadata_address, coin_metadata.get_decimals(), vector[], ctx),
        coin_store: bag::new(ctx),
    };
    lock_release_token_pool.coin_store.add(COIN_STORE, coin::zero<T>(ctx));

    token_admin_registry::register_pool(
        ref,
        treasury_cap,
        coin_metadata_address,
        object::uid_to_address(&lock_release_token_pool.id),
        ctx.sender(),
        TypeProof {},
    );

    let owner_cap = OwnerCap {
        id: object::new(ctx),
    };

    transfer::share_object(lock_release_token_pool);
    transfer::public_transfer(owner_cap, ctx.sender());
}

// ================================================================
// |                 Exposing token_pool functions                |
// ================================================================

// this now returns the address of coin metadata
public fun get_token(state: &LockReleaseTokenPoolState): address {
    token_pool::get_token(&state.token_pool_state)
}

public fun get_router(): address {
    token_pool::get_router()
}

public fun get_token_decimals(state: &LockReleaseTokenPoolState): u8 {
    state.token_pool_state.get_local_decimals()
}

public fun get_remote_pools(
    state: &LockReleaseTokenPoolState,
    remote_chain_selector: u64
): vector<vector<u8>> {
    token_pool::get_remote_pools(&state.token_pool_state, remote_chain_selector)
}

public fun is_remote_pool(
    state: &LockReleaseTokenPoolState,
    remote_chain_selector: u64,
    remote_pool_address: vector<u8>
): bool {
    token_pool::is_remote_pool(
        &state.token_pool_state,
        remote_chain_selector,
        remote_pool_address
    )
}

public fun get_remote_token(
    state: &LockReleaseTokenPoolState, remote_chain_selector: u64
): vector<u8> {
    token_pool::get_remote_token(&state.token_pool_state, remote_chain_selector)
}

public fun add_remote_pool(
    state: &mut LockReleaseTokenPoolState,
    _: &OwnerCap,
    remote_chain_selector: u64,
    remote_pool_address: vector<u8>,
) {
    token_pool::add_remote_pool(
        &mut state.token_pool_state, remote_chain_selector, remote_pool_address
    );
}

public fun remove_remote_pool(
    state: &mut LockReleaseTokenPoolState,
    _: &OwnerCap,
    remote_chain_selector: u64,
    remote_pool_address: vector<u8>,
) {
    token_pool::remove_remote_pool(
        &mut state.token_pool_state, remote_chain_selector, remote_pool_address
    );
}

public fun is_supported_chain(
    state: &LockReleaseTokenPoolState,
    remote_chain_selector: u64
): bool {
    token_pool::is_supported_chain(&state.token_pool_state, remote_chain_selector)
}

public fun get_supported_chains(state: &LockReleaseTokenPoolState): vector<u64> {
    token_pool::get_supported_chains(&state.token_pool_state)
}

public fun apply_chain_updates(
    state: &mut LockReleaseTokenPoolState,
    _: &OwnerCap,
    remote_chain_selectors_to_remove: vector<u64>,
    remote_chain_selectors_to_add: vector<u64>,
    remote_pool_addresses_to_add: vector<vector<vector<u8>>>,
    remote_token_addresses_to_add: vector<vector<u8>>
) {
    token_pool::apply_chain_updates(
        &mut state.token_pool_state,
        remote_chain_selectors_to_remove,
        remote_chain_selectors_to_add,
        remote_pool_addresses_to_add,
        remote_token_addresses_to_add
    );
}

public fun get_allowlist_enabled(state: &LockReleaseTokenPoolState): bool {
    token_pool::get_allowlist_enabled(&state.token_pool_state)
}

public fun get_allowlist(state: &LockReleaseTokenPoolState): vector<address> {
    token_pool::get_allowlist(&state.token_pool_state)
}

public fun apply_allowlist_updates(
    state: &mut LockReleaseTokenPoolState,
    _: &OwnerCap,
    removes: vector<address>,
    adds: vector<address>
) {
    token_pool::apply_allowlist_updates(&mut state.token_pool_state, removes, adds);
}

// ================================================================
// |                         Burn/Mint                            |
// ================================================================

public struct TypeProof has drop {}

public fun lock_or_burn<T>(
    ref: &CCIPObjectRef,
    clock: &Clock,
    state: &mut LockReleaseTokenPoolState,
    c: Coin<T>,
    remote_chain_selector: u64,
    token_params: dd::TokenParams,
    ctx: &mut TxContext
): dd::TokenParams {
    let amount = c.value();
    let sender = ctx.sender();

    // This metod validates various aspects of the lock or burn operation. If any of the
    // validations fail, the transaction will abort.
    let dest_token_address =
        token_pool::validate_lock_or_burn(
            ref,
            clock,
            &mut state.token_pool_state,
            sender,
            remote_chain_selector,
            amount,
        );

    let stored_coin: &mut Coin<T> = state.coin_store.borrow_mut(COIN_STORE);
    stored_coin.join(c);

    let mut extra_data = vector[];
    eth_abi::encode_u8(&mut extra_data, state.token_pool_state.get_local_decimals());

    token_pool::emit_locked(&mut state.token_pool_state, amount);

    // update hot potato token params
    dd::add_source_token_transfer(
        ref,
        token_params,
        object::uid_to_address(&state.id),
        amount,
        state.token_pool_state.get_token(),
        dest_token_address,
        extra_data,
        TypeProof {},
    )
}

// TODO: if there are more validations to be done
public fun release_or_mint<T>(
    ref: &CCIPObjectRef,
    clock: &Clock,
    pool: &mut LockReleaseTokenPoolState,
    remote_chain_selector: u64,
    receiver_params: osh::ReceiverParams,
    index: u64,
    ctx: &mut TxContext
): osh::ReceiverParams {
    let (receiver, source_amount, dest_token_address, source_pool_address, source_pool_data) = osh::get_token_param_data(&receiver_params, index);
    let local_decimals = pool.token_pool_state.get_local_decimals();
    let remote_decimals = token_pool::parse_remote_decimals(source_pool_data, local_decimals);
    let local_amount = token_pool::calculate_local_amount(source_amount as u256, remote_decimals, local_decimals);

    token_pool::validate_release_or_mint(
        ref,
        clock,
        &mut pool.token_pool_state,
        remote_chain_selector,
        dest_token_address,
        source_pool_address,
        local_amount
    );

    // split the coin to be released
    let stored_coin: &mut Coin<T> = pool.coin_store.borrow_mut(COIN_STORE);
    assert!(
        stored_coin.value() >= local_amount,
        E_TOKEN_POOL_BALANCE_TOO_LOW
    );
    let c: Coin<T> = stored_coin.split(local_amount, ctx);
    transfer::public_transfer(c, receiver);

    token_pool::emit_released(
        &mut pool.token_pool_state,
        receiver,
        local_amount,
    );

    osh::complete_token_transfer(
        ref,
        receiver_params,
        index,
        local_amount,
        object::uid_to_address(&pool.id),
        TypeProof {},
    )
}

// ================================================================
// |                    Rate limit config                         |
// ================================================================

public fun set_chain_rate_limiter_configs(
    state: &mut LockReleaseTokenPoolState,
    _: &OwnerCap,
    clock: &Clock,
    remote_chain_selectors: vector<u64>,
    outbound_is_enableds: vector<bool>,
    outbound_capacities: vector<u64>,
    outbound_rates: vector<u64>,
    inbound_is_enableds: vector<bool>,
    inbound_capacities: vector<u64>,
    inbound_rates: vector<u64>
) {
    let number_of_chains = remote_chain_selectors.length();

    assert!(
        number_of_chains == outbound_is_enableds.length()
            && number_of_chains == outbound_capacities.length()
            && number_of_chains == outbound_rates.length()
            && number_of_chains == inbound_is_enableds.length()
            && number_of_chains == inbound_capacities.length()
            && number_of_chains == inbound_rates.length(),
        E_INVALID_ARGUMENTS
    );

    let mut i = 0;
    while (i < number_of_chains) {
        token_pool::set_chain_rate_limiter_config(
            clock,
            &mut state.token_pool_state,
            remote_chain_selectors[i],
            outbound_is_enableds[i],
            outbound_capacities[i],
            outbound_rates[i],
            inbound_is_enableds[i],
            inbound_capacities[i],
            inbound_rates[i]
        );
        i = i + 1;
    };
}

public fun set_chain_rate_limiter_config(
    state: &mut LockReleaseTokenPoolState,
    _: &OwnerCap,
    clock: &Clock,
    remote_chain_selector: u64,
    outbound_is_enabled: bool,
    outbound_capacity: u64,
    outbound_rate: u64,
    inbound_is_enabled: bool,
    inbound_capacity: u64,
    inbound_rate: u64
) {
    token_pool::set_chain_rate_limiter_config(
        clock,
        &mut state.token_pool_state,
        remote_chain_selector,
        outbound_is_enabled,
        outbound_capacity,
        outbound_rate,
        inbound_is_enabled,
        inbound_capacity,
        inbound_rate
    );
}

public fun get_balance<T>(state: &LockReleaseTokenPoolState): u64 {
    let stored_coin: &Coin<T> = state.coin_store.borrow(COIN_STORE);
    stored_coin.value()
}

public fun add_liquidity<T>(
    state: &mut LockReleaseTokenPoolState,
    c: Coin<T>
) {
    let stored_coin: &mut Coin<T> = state.coin_store.borrow_mut(COIN_STORE);
    stored_coin.join(c);
}

public fun remove_liquidity<T>(
    state: &mut LockReleaseTokenPoolState,
    _: &OwnerCap,
    amount: u64,
    ctx: &mut TxContext
): Coin<T> {
    let stored_coin: &mut Coin<T> = state.coin_store.borrow_mut(COIN_STORE);
    assert!(stored_coin.value() >= amount, E_TOKEN_POOL_BALANCE_TOO_LOW);
    stored_coin.split(amount, ctx)
}

#[allow(lint(self_transfer))]
public fun remove_liquidity_and_transfer<T>(
    state: &mut LockReleaseTokenPoolState,
    _: &OwnerCap,
    amount: u64,
    ctx: &mut TxContext
) {
    let stored_coin: &mut Coin<T> = state.coin_store.borrow_mut(COIN_STORE);
    assert!(stored_coin.value() >= amount, E_TOKEN_POOL_BALANCE_TOO_LOW);
    transfer::public_transfer(stored_coin.split(amount, ctx), ctx.sender());
}
