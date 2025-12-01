module lock_release_token_pool::lock_release_token_pool;

use ccip::eth_abi;
use ccip::offramp_state_helper as offramp_sh;
use ccip::onramp_state_helper as onramp_sh;
use ccip::publisher_wrapper;
use ccip::state_object::CCIPObjectRef;
use ccip::token_admin_registry;
use lock_release_token_pool::ownable::{Self, OwnerCap, OwnableState};
use lock_release_token_pool::rate_limiter;
use lock_release_token_pool::token_pool::{Self, TokenPoolState};
use mcms::bcs_stream;
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, Registry, ExecutingCallbackParams};
use std::ascii;
use std::string::{Self, String};
use std::type_name;
use sui::address;
use sui::clock::Clock;
use sui::coin::{Self, Coin, CoinMetadata, TreasuryCap};
use sui::derived_object;
use sui::event;
use sui::package::{Self, UpgradeCap};

public struct LOCK_RELEASE_TOKEN_POOL has drop {}

public struct LockReleaseTokenPoolObject has key {
    id: UID,
}

public struct LockReleaseTokenPoolStatePointer has key, store {
    id: UID,
    lock_release_token_pool_object_id: address,
}

fun init(otw: LOCK_RELEASE_TOKEN_POOL, ctx: &mut TxContext) {
    let (ownable_state, mut owner_cap) = ownable::new(ctx);
    ownable::attach_ownable_state(&mut owner_cap, ownable_state);

    let publisher = package::claim(otw, ctx);
    ownable::attach_publisher(&mut owner_cap, publisher);

    transfer::public_transfer(owner_cap, ctx.sender());
}

#[allow(lint(coin_field))]
public struct LockReleaseTokenPoolState<phantom T> has key {
    id: UID,
    token_pool_state: TokenPoolState,
    reserve: Coin<T>,
    rebalancer_cap_id: ID,
    ownable_state: OwnableState,
}

public struct RebalancerCap<phantom T> has key, store {
    id: UID,
}

public struct RebalancerSet<phantom T> has copy, drop {
    old_rebalancer_cap_id: ID,
    new_rebalancer_cap_id: ID,
}

public struct TokenBucketWrapper has drop, store {
    tokens: u64,
    last_updated: u64,
    is_enabled: bool,
    capacity: u64,
    rate: u64,
}

const CLOCK_ADDRESS: address = @0x6;

const EInvalidArguments: u64 = 1;
const ETokenPoolBalanceTooLow: u64 = 2;
const EInvalidOwnerCap: u64 = 3;
const EInvalidFunction: u64 = 4;
const EInvalidRebalancerCap: u64 = 5;
const EPoolStillRegistered: u64 = 6;
const ERebalancerCapIsInUse: u64 = 7;
const ERebalancerCapNotTransferredOut: u64 = 8;
const EInvalidRebalancer: u64 = 9;
const ERebalancerCapDoesNotExist: u64 = 10;
const ERebalancerCapMismatch: u64 = 11;

// ================================================================
// |                             Init                             |
// ================================================================

public fun type_and_version(): String {
    string::utf8(b"LockReleaseTokenPool 1.6.0")
}

public fun initialize<T>(
    owner_cap: &mut OwnerCap,
    ref: &mut CCIPObjectRef,
    coin_metadata: &CoinMetadata<T>,
    treasury_cap: &TreasuryCap<T>,
    token_pool_administrator: address,
    rebalancer: address,
    ctx: &mut TxContext,
) {
    let coin_metadata_address: address = object::id_to_address(&object::id(coin_metadata));
    let ownable_state = ownable::detach_ownable_state(owner_cap);
    let mut lock_release_token_pool_object = LockReleaseTokenPoolObject { id: object::new(ctx) };
    let lock_release_token_pool_state_pointer = LockReleaseTokenPoolStatePointer {
        id: object::new(ctx),
        lock_release_token_pool_object_id: object::id_address(&lock_release_token_pool_object),
    };
    let rebalancer_cap = RebalancerCap<T> { id: object::new(ctx) };
    let lock_release_token_pool = LockReleaseTokenPoolState<T> {
        id: derived_object::claim(
            &mut lock_release_token_pool_object.id,
            b"LockReleaseTokenPoolState",
        ),
        token_pool_state: token_pool::initialize(
            coin_metadata_address,
            coin_metadata.get_decimals(),
            coin_metadata.get_symbol(),
            vector[],
            ctx,
        ),
        reserve: coin::zero<T>(ctx),
        rebalancer_cap_id: object::id(&rebalancer_cap),
        ownable_state,
    };

    let tn = type_name::with_original_ids<LOCK_RELEASE_TOKEN_POOL>();
    let package_bytes = ascii::into_bytes(tn.address_string());
    let package_id = address::from_ascii_bytes(&package_bytes);

    let publisher_wrapper = publisher_wrapper::create(
        ownable::borrow_publisher(owner_cap),
        TypeProof {},
    );
    let lock_release_token_pool_state_address = object::uid_to_address(&lock_release_token_pool.id);

    token_admin_registry::register_pool(
        ref,
        treasury_cap,
        coin_metadata,
        token_pool_administrator,
        vector[CLOCK_ADDRESS, lock_release_token_pool_state_address],
        vector[CLOCK_ADDRESS, lock_release_token_pool_state_address],
        publisher_wrapper,
        TypeProof {},
    );

    transfer::share_object(lock_release_token_pool);
    transfer::share_object(lock_release_token_pool_object);
    transfer::public_transfer(rebalancer_cap, rebalancer);
    transfer::transfer(lock_release_token_pool_state_pointer, package_id);
}

// ================================================================
// |                 Exposing token_pool functions                |
// ================================================================

/// returns the coin metadata object id of the token
public fun get_token<T>(state: &LockReleaseTokenPoolState<T>): address {
    token_pool::get_token(&state.token_pool_state)
}

public fun get_token_decimals<T>(state: &LockReleaseTokenPoolState<T>): u8 {
    state.token_pool_state.get_local_decimals()
}

public fun get_token_symbol<T>(state: &LockReleaseTokenPoolState<T>): ascii::String {
    state.token_pool_state.get_symbol()
}

public fun get_remote_pools<T>(
    state: &LockReleaseTokenPoolState<T>,
    remote_chain_selector: u64,
): vector<vector<u8>> {
    token_pool::get_remote_pools(&state.token_pool_state, remote_chain_selector)
}

public fun is_remote_pool<T>(
    state: &LockReleaseTokenPoolState<T>,
    remote_chain_selector: u64,
    remote_pool_address: vector<u8>,
): bool {
    token_pool::is_remote_pool(
        &state.token_pool_state,
        remote_chain_selector,
        remote_pool_address,
    )
}

public fun get_remote_token<T>(
    state: &LockReleaseTokenPoolState<T>,
    remote_chain_selector: u64,
): vector<u8> {
    token_pool::get_remote_token(&state.token_pool_state, remote_chain_selector)
}

public fun add_remote_pool<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    owner_cap: &OwnerCap,
    remote_chain_selector: u64,
    remote_pool_address: vector<u8>,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    token_pool::add_remote_pool(
        &mut state.token_pool_state,
        remote_chain_selector,
        remote_pool_address,
    );
}

public fun remove_remote_pool<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    owner_cap: &OwnerCap,
    remote_chain_selector: u64,
    remote_pool_address: vector<u8>,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    token_pool::remove_remote_pool(
        &mut state.token_pool_state,
        remote_chain_selector,
        remote_pool_address,
    );
}

public fun is_supported_chain<T>(
    state: &LockReleaseTokenPoolState<T>,
    remote_chain_selector: u64,
): bool {
    token_pool::is_supported_chain(&state.token_pool_state, remote_chain_selector)
}

public fun get_supported_chains<T>(state: &LockReleaseTokenPoolState<T>): vector<u64> {
    token_pool::get_supported_chains(&state.token_pool_state)
}

public fun apply_chain_updates<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    owner_cap: &OwnerCap,
    remote_chain_selectors_to_remove: vector<u64>,
    remote_chain_selectors_to_add: vector<u64>,
    remote_pool_addresses_to_add: vector<vector<vector<u8>>>,
    remote_token_addresses_to_add: vector<vector<u8>>,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    token_pool::apply_chain_updates(
        &mut state.token_pool_state,
        remote_chain_selectors_to_remove,
        remote_chain_selectors_to_add,
        remote_pool_addresses_to_add,
        remote_token_addresses_to_add,
    );
}

public fun get_allowlist_enabled<T>(state: &LockReleaseTokenPoolState<T>): bool {
    token_pool::get_allowlist_enabled(&state.token_pool_state)
}

public fun get_allowlist<T>(state: &LockReleaseTokenPoolState<T>): vector<address> {
    token_pool::get_allowlist(&state.token_pool_state)
}

public fun set_allowlist_enabled<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    owner_cap: &OwnerCap,
    enabled: bool,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    token_pool::set_allowlist_enabled(&mut state.token_pool_state, enabled);
}

public fun apply_allowlist_updates<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    owner_cap: &OwnerCap,
    removes: vector<address>,
    adds: vector<address>,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    token_pool::apply_allowlist_updates(&mut state.token_pool_state, removes, adds);
}

// ================================================================
// |                        Lock/Release                          |
// ================================================================

public struct TypeProof has drop {}

public fun lock_or_burn<T: drop>(
    ref: &CCIPObjectRef,
    token_transfer_params: &mut onramp_sh::TokenTransferParams,
    c: Coin<T>,
    remote_chain_selector: u64,
    clock: &Clock,
    state: &mut LockReleaseTokenPoolState<T>,
    ctx: &mut TxContext,
) {
    let amount = c.value();
    let sender = ctx.sender();

    // This function validates various aspects of the lock or burn operation. If any of the validations fail, the transaction will abort.
    let dest_token_address = token_pool::get_remote_token(
        &state.token_pool_state,
        remote_chain_selector,
    );
    token_pool::validate_lock_or_burn(
        ref,
        clock,
        &mut state.token_pool_state,
        sender,
        remote_chain_selector,
        amount,
    );
    coin::join(&mut state.reserve, c);

    let mut extra_data = vector[];
    eth_abi::encode_u8(&mut extra_data, state.token_pool_state.get_local_decimals());

    token_pool::emit_locked_or_burned(&state.token_pool_state, amount, remote_chain_selector);

    onramp_sh::add_token_transfer_param(
        ref,
        token_transfer_params,
        remote_chain_selector,
        amount,
        get_token(state),
        dest_token_address,
        extra_data,
        TypeProof {},
    )
}

/// after releasing the token, this function will mark this particular token transfer as complete
/// and set the local amount of this token transfer according to the balance of coin object.
/// a token pool cannot update token transfer item for other tokens simply by changing the
/// index because each token transfer is protected by a type proof
public fun release_or_mint<T>(
    ref: &CCIPObjectRef,
    receiver_params: &mut offramp_sh::ReceiverParams,
    clock: &Clock,
    state: &mut LockReleaseTokenPoolState<T>,
    ctx: &mut TxContext,
) {
    let (
        token_receiver,
        remote_chain_selector,
        source_amount,
        dest_token_address,
        _,
        source_pool_address,
        source_pool_data,
        _,
    ) = offramp_sh::get_dest_token_transfer_data(receiver_params);

    // local_amount is u64 because the token balance in SUI is u64.
    let local_amount = token_pool::calculate_release_or_mint_amount(
        &state.token_pool_state,
        source_pool_data,
        source_amount,
    );

    token_pool::validate_release_or_mint(
        ref,
        clock,
        &mut state.token_pool_state,
        remote_chain_selector,
        dest_token_address,
        source_pool_address,
        local_amount,
    );

    // split the coin to be released
    assert!(state.reserve.value() >= local_amount, ETokenPoolBalanceTooLow);
    let c: Coin<T> = coin::split(&mut state.reserve, local_amount, ctx);

    token_pool::emit_released_or_minted(
        &state.token_pool_state,
        token_receiver,
        local_amount,
        remote_chain_selector,
    );
    transfer::public_transfer(c, token_receiver);

    offramp_sh::complete_token_transfer(
        ref,
        receiver_params,
        TypeProof {},
    );
}

// ================================================================
// |                    Rate limit config                         |
// ================================================================

public fun set_chain_rate_limiter_configs<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    owner_cap: &OwnerCap,
    clock: &Clock,
    remote_chain_selectors: vector<u64>,
    outbound_is_enableds: vector<bool>,
    outbound_capacities: vector<u64>,
    outbound_rates: vector<u64>,
    inbound_is_enableds: vector<bool>,
    inbound_capacities: vector<u64>,
    inbound_rates: vector<u64>,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    let number_of_chains = remote_chain_selectors.length();

    assert!(
        number_of_chains == outbound_is_enableds.length()
            && number_of_chains == outbound_capacities.length()
            && number_of_chains == outbound_rates.length()
            && number_of_chains == inbound_is_enableds.length()
            && number_of_chains == inbound_capacities.length()
            && number_of_chains == inbound_rates.length(),
        EInvalidArguments,
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
            inbound_rates[i],
        );
        i = i + 1;
    };
}

public fun set_chain_rate_limiter_config<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    owner_cap: &OwnerCap,
    clock: &Clock,
    remote_chain_selector: u64,
    outbound_is_enabled: bool,
    outbound_capacity: u64,
    outbound_rate: u64,
    inbound_is_enabled: bool,
    inbound_capacity: u64,
    inbound_rate: u64,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    token_pool::set_chain_rate_limiter_config(
        clock,
        &mut state.token_pool_state,
        remote_chain_selector,
        outbound_is_enabled,
        outbound_capacity,
        outbound_rate,
        inbound_is_enabled,
        inbound_capacity,
        inbound_rate,
    );
}

public fun get_current_inbound_rate_limiter_state<T>(
    clock: &Clock,
    state: &LockReleaseTokenPoolState<T>,
    remote_chain_selector: u64,
): TokenBucketWrapper {
    let token_bucket = token_pool::get_current_inbound_rate_limiter_state(
        &state.token_pool_state,
        clock,
        remote_chain_selector,
    );
    let (tokens, last_updated, is_enabled, capacity, rate) = rate_limiter::get_token_bucket_fields(
        &token_bucket,
    );
    TokenBucketWrapper { tokens, last_updated, is_enabled, capacity, rate }
}

public fun get_current_outbound_rate_limiter_state<T>(
    clock: &Clock,
    state: &LockReleaseTokenPoolState<T>,
    remote_chain_selector: u64,
): TokenBucketWrapper {
    let token_bucket = token_pool::get_current_outbound_rate_limiter_state(
        &state.token_pool_state,
        clock,
        remote_chain_selector,
    );
    let (tokens, last_updated, is_enabled, capacity, rate) = rate_limiter::get_token_bucket_fields(
        &token_bucket,
    );
    TokenBucketWrapper { tokens, last_updated, is_enabled, capacity, rate }
}

// ================================================================
// |                    Liquidity Management                      |
// ================================================================

public fun provide_liquidity<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    rebalancer_cap: &RebalancerCap<T>,
    c: Coin<T>,
    _: &mut TxContext,
) {
    assert!(object::id(rebalancer_cap) == state.rebalancer_cap_id, EInvalidRebalancerCap);
    let amount = c.value();

    coin::join(&mut state.reserve, c);

    token_pool::emit_liquidity_added(
        &state.token_pool_state,
        object::id_to_address(&state.rebalancer_cap_id),
        amount,
    );
}

public fun withdraw_liquidity<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    rebalancer_cap: &RebalancerCap<T>,
    amount: u64,
    ctx: &mut TxContext,
): Coin<T> {
    assert!(object::id(rebalancer_cap) == state.rebalancer_cap_id, EInvalidRebalancerCap);
    assert!(state.reserve.value() >= amount, ETokenPoolBalanceTooLow);

    token_pool::emit_liquidity_removed(
        &state.token_pool_state,
        object::id_to_address(&state.rebalancer_cap_id),
        amount,
    );
    coin::split(&mut state.reserve, amount, ctx)
}

public fun get_balance<T>(state: &LockReleaseTokenPoolState<T>): u64 {
    state.reserve.value()
}

public fun get_rebalancer<T>(state: &LockReleaseTokenPoolState<T>): address {
    object::id_to_address(&state.rebalancer_cap_id)
}

public struct McmsCap<phantom T> has key, store {
    id: UID,
    owner_cap: OwnerCap,
    rebalancer_cap: Option<RebalancerCap<T>>,
}

fun new_mcms_cap<T>(
    owner_cap: OwnerCap,
    rebalancer_cap: Option<RebalancerCap<T>>,
    ctx: &mut TxContext,
): McmsCap<T> {
    McmsCap {
        id: object::new(ctx),
        owner_cap,
        rebalancer_cap,
    }
}

/// Sets the rebalancer for the pool when owned by MCMS. This function supports two use cases:
///
/// 1. Setting rebalancer to the MCMS address stores the RebalancerCap inside McmsCap.
///    If MCMS already has rebalancer control, it is a no-op.
///
/// 2. Setting rebalancer to an EOA address creates a new RebalancerCap and transfers it to that address.
///    If MCMS was previously rebalancer, the old cap is destroyed.
///
/// This function can only be called when MCMS owns the pool (via McmsCap).
/// For EOA-to-EOA rebalancer changes, use `set_rebalancer` instead.
public fun mcms_set_rebalancer<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (mcms_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        McmsCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"set_rebalancer"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(&mcms_cap.owner_cap)],
        &mut stream,
    );
    let rebalancer = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    if (rebalancer == mcms_registry::get_multisig_address()) {
        // MCMS already has rebalancer control - verify and no-op
        if (mcms_cap.rebalancer_cap.is_some()) {
            let existing_cap_id = object::id(mcms_cap.rebalancer_cap.borrow());
            assert!(existing_cap_id == state.rebalancer_cap_id, ERebalancerCapMismatch);
        } else {
            // MCMS taking/re-taking rebalancer control
            let new_rebalancer_cap = RebalancerCap<T> {
                id: object::new(ctx),
            };
            let new_rebalancer_cap_id = object::id(&new_rebalancer_cap);
            mcms_cap.rebalancer_cap.fill(new_rebalancer_cap);

            set_rebalancer_cap_id(state, &mcms_cap.owner_cap, new_rebalancer_cap_id, ctx);
        }
    } else {
        // MCMS delegating rebalancer control to an EOA
        set_rebalancer(state, &mcms_cap.owner_cap, rebalancer, ctx);

        // If MCMS previously had rebalancer control, destroy that cap since it's no longer needed
        if (mcms_cap.rebalancer_cap.is_some()) {
            destroy_rebalancer_cap(state, mcms_cap.rebalancer_cap.extract(), ctx);
        }
    }
}

public fun set_rebalancer<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    owner_cap: &OwnerCap,
    rebalancer: address,
    ctx: &mut TxContext,
) {
    // Setting MCMS as rebalancer should be done via `mcms_set_rebalancer`
    assert!(rebalancer != mcms_registry::get_multisig_address(), EInvalidRebalancer);

    let rebalancer_cap = RebalancerCap<T> {
        id: object::new(ctx),
    };

    // Update `LockReleaseTokenPoolState` before sending to rebalancer address
    set_rebalancer_cap_id(
        state,
        owner_cap,
        object::id(&rebalancer_cap),
        ctx,
    );

    transfer::public_transfer(rebalancer_cap, rebalancer);
}

/// Only owner can set rebalancer cap id
fun set_rebalancer_cap_id<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    owner_cap: &OwnerCap,
    new_rebalancer_cap_id: ID,
    _ctx: &mut TxContext,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);

    let old_rebalancer_cap_id = state.rebalancer_cap_id;
    state.rebalancer_cap_id = new_rebalancer_cap_id;

    event::emit(RebalancerSet<T> { old_rebalancer_cap_id, new_rebalancer_cap_id });
}

public fun mcms_destroy_rebalancer_cap<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (mcms_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        McmsCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"destroy_rebalancer_cap"), EInvalidFunction);
    assert!(mcms_cap.rebalancer_cap.is_some(), ERebalancerCapDoesNotExist);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(mcms_cap.rebalancer_cap.borrow())],
        &mut stream,
    );
    bcs_stream::assert_is_consumed(&stream);

    let rebalancer_cap = mcms_cap.rebalancer_cap.extract();
    destroy_rebalancer_cap(state, rebalancer_cap, ctx);
}

/// Clean up old rebalancer caps not in use, anyone can call this function to clean up old rebalancer caps not in use.
public fun destroy_rebalancer_cap<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    rebalancer_cap: RebalancerCap<T>,
    _ctx: &mut TxContext,
) {
    assert!(state.rebalancer_cap_id != object::id(&rebalancer_cap), ERebalancerCapIsInUse);

    let RebalancerCap<T> { id } = rebalancer_cap;
    object::delete(id);
}

// ================================================================
// |                      Ownable Functions                       |
// ================================================================

public fun owner<T>(state: &LockReleaseTokenPoolState<T>): address {
    ownable::owner(&state.ownable_state)
}

public fun has_pending_transfer<T>(state: &LockReleaseTokenPoolState<T>): bool {
    ownable::has_pending_transfer(&state.ownable_state)
}

public fun pending_transfer_from<T>(state: &LockReleaseTokenPoolState<T>): Option<address> {
    ownable::pending_transfer_from(&state.ownable_state)
}

public fun pending_transfer_to<T>(state: &LockReleaseTokenPoolState<T>): Option<address> {
    ownable::pending_transfer_to(&state.ownable_state)
}

public fun pending_transfer_accepted<T>(state: &LockReleaseTokenPoolState<T>): Option<bool> {
    ownable::pending_transfer_accepted(&state.ownable_state)
}

public fun transfer_ownership<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    owner_cap: &OwnerCap,
    new_owner: address,
    ctx: &mut TxContext,
) {
    ownable::transfer_ownership(owner_cap, &mut state.ownable_state, new_owner, ctx);
}

public fun accept_ownership<T>(state: &mut LockReleaseTokenPoolState<T>, ctx: &mut TxContext) {
    ownable::accept_ownership(&mut state.ownable_state, ctx);
}

public fun accept_ownership_from_object<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    from: &mut UID,
    ctx: &mut TxContext,
) {
    ownable::accept_ownership_from_object(&mut state.ownable_state, from, ctx);
}

public fun mcms_accept_ownership<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let data = mcms_registry::get_accept_ownership_data(
        registry,
        params,
        McmsAcceptOwnershipProof<T> {},
    );

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addr(object::id_address(state), &mut stream);
    bcs_stream::assert_is_consumed(&stream);

    let mcms = mcms_registry::get_multisig_address();
    ownable::mcms_accept_ownership(&mut state.ownable_state, mcms, ctx);
}

public fun execute_ownership_transfer<T>(
    owner_cap: OwnerCap,
    state: &mut LockReleaseTokenPoolState<T>,
    to: address,
    ctx: &mut TxContext,
) {
    ownable::execute_ownership_transfer(owner_cap, &mut state.ownable_state, to, ctx);
}

/// Setting rebalancer cap should be called via `mcms_set_rebalancer`
public fun execute_ownership_transfer_to_mcms<T>(
    owner_cap: OwnerCap,
    state: &mut LockReleaseTokenPoolState<T>,
    registry: &mut Registry,
    to: address,
    ctx: &mut TxContext,
) {
    assert!(object::id(&owner_cap) == state.ownable_state.owner_cap_id(), EInvalidOwnerCap);

    let publisher_wrapper = mcms_registry::create_publisher_wrapper(
        ownable::borrow_publisher(&owner_cap),
        McmsCallback<T> {},
    );

    ownable::execute_ownership_and_cap_transfer_to_mcms(
        &mut state.ownable_state,
        registry,
        new_mcms_cap(owner_cap, option::none<RebalancerCap<T>>(), ctx),
        to,
        publisher_wrapper,
        McmsCallback<T> {},
        vector[b"lock_release_token_pool"],
        ctx,
    );
}

public fun mcms_register_upgrade_cap(
    upgrade_cap: UpgradeCap,
    registry: &mut Registry,
    state: &mut DeployerState,
    ctx: &mut TxContext,
) {
    mcms_deployer::register_upgrade_cap(
        state,
        registry,
        upgrade_cap,
        ctx,
    );
}

// ================================================================
// |                      MCMS Entrypoint                         |
// ================================================================

public struct McmsCallback<phantom T> has drop {}

/// Proof for MCMS Accept Ownership
public struct McmsAcceptOwnershipProof<phantom T> has drop {}

public fun mcms_set_allowlist_enabled<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (mcms_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        McmsCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"set_allowlist_enabled"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(&mcms_cap.owner_cap)],
        &mut stream,
    );

    let enabled = bcs_stream::deserialize_bool(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    set_allowlist_enabled(state, &mcms_cap.owner_cap, enabled);
}

public fun mcms_apply_allowlist_updates<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (mcms_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        McmsCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"apply_allowlist_updates"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(&mcms_cap.owner_cap)],
        &mut stream,
    );

    let removes = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_address(stream),
    );
    let adds = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_address(stream),
    );
    bcs_stream::assert_is_consumed(&stream);

    apply_allowlist_updates(state, &mcms_cap.owner_cap, removes, adds);
}

public fun mcms_apply_chain_updates<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (mcms_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        McmsCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"apply_chain_updates"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(&mcms_cap.owner_cap)],
        &mut stream,
    );

    let remote_chain_selectors_to_remove = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_u64(stream),
    );
    let remote_chain_selectors_to_add = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_u64(stream),
    );
    let remote_pool_addresses_to_add = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_vector!(
            stream,
            |stream| bcs_stream::deserialize_vector_u8(stream),
        ),
    );
    let remote_token_addresses_to_add = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_vector_u8(stream),
    );
    bcs_stream::assert_is_consumed(&stream);

    apply_chain_updates(
        state,
        &mcms_cap.owner_cap,
        remote_chain_selectors_to_remove,
        remote_chain_selectors_to_add,
        remote_pool_addresses_to_add,
        remote_token_addresses_to_add,
    );
}

public fun mcms_add_remote_pool<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (mcms_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        McmsCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"add_remote_pool"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(&mcms_cap.owner_cap)],
        &mut stream,
    );
    let remote_chain_selector = bcs_stream::deserialize_u64(&mut stream);
    let remote_pool_address = bcs_stream::deserialize_vector_u8(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    add_remote_pool(state, &mcms_cap.owner_cap, remote_chain_selector, remote_pool_address);
}

public fun mcms_remove_remote_pool<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (mcms_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        McmsCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"remove_remote_pool"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(&mcms_cap.owner_cap)],
        &mut stream,
    );
    let remote_chain_selector = bcs_stream::deserialize_u64(&mut stream);
    let remote_pool_address = bcs_stream::deserialize_vector_u8(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    remove_remote_pool(state, &mcms_cap.owner_cap, remote_chain_selector, remote_pool_address);
}

public fun mcms_set_chain_rate_limiter_configs<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    clock: &Clock,
) {
    let (mcms_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        McmsCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"set_chain_rate_limiter_configs"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[
            object::id_address(state),
            object::id_address(&mcms_cap.owner_cap),
            object::id_address(clock),
        ],
        &mut stream,
    );

    let remote_chain_selectors = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_u64(stream),
    );
    let outbound_is_enableds = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_bool(stream),
    );
    let outbound_capacities = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_u64(stream),
    );
    let outbound_rates = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_u64(stream),
    );
    let inbound_is_enableds = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_bool(stream),
    );
    let inbound_capacities = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_u64(stream),
    );
    let inbound_rates = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_u64(stream),
    );
    bcs_stream::assert_is_consumed(&stream);

    set_chain_rate_limiter_configs(
        state,
        &mcms_cap.owner_cap,
        clock,
        remote_chain_selectors,
        outbound_is_enableds,
        outbound_capacities,
        outbound_rates,
        inbound_is_enableds,
        inbound_capacities,
        inbound_rates,
    );
}

public fun mcms_set_chain_rate_limiter_config<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    clock: &Clock,
) {
    let (mcms_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        McmsCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"set_chain_rate_limiter_config"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[
            object::id_address(state),
            object::id_address(&mcms_cap.owner_cap),
            object::id_address(clock),
        ],
        &mut stream,
    );
    let remote_chain_selector = bcs_stream::deserialize_u64(&mut stream);
    let outbound_is_enabled = bcs_stream::deserialize_bool(&mut stream);
    let outbound_capacity = bcs_stream::deserialize_u64(&mut stream);
    let outbound_rate = bcs_stream::deserialize_u64(&mut stream);
    let inbound_is_enabled = bcs_stream::deserialize_bool(&mut stream);
    let inbound_capacity = bcs_stream::deserialize_u64(&mut stream);
    let inbound_rate = bcs_stream::deserialize_u64(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    set_chain_rate_limiter_config(
        state,
        &mcms_cap.owner_cap,
        clock,
        remote_chain_selector,
        outbound_is_enabled,
        outbound_capacity,
        outbound_rate,
        inbound_is_enabled,
        inbound_capacity,
        inbound_rate,
    );
}

public fun mcms_provide_liquidity<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    registry: &mut Registry,
    coin: Coin<T>,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (mcms_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        McmsCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"provide_liquidity"), EInvalidFunction);
    assert!(mcms_cap.rebalancer_cap.is_some(), ERebalancerCapDoesNotExist);

    let rebalancer_cap = mcms_cap.rebalancer_cap.borrow();

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[
            object::id_address(state),
            object::id_address(rebalancer_cap),
            object::id_address(&coin),
        ],
        &mut stream,
    );
    bcs_stream::assert_is_consumed(&stream);

    provide_liquidity(state, rebalancer_cap, coin, ctx);
}

public fun mcms_withdraw_liquidity<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (mcms_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        McmsCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"withdraw_liquidity"), EInvalidFunction);
    assert!(mcms_cap.rebalancer_cap.is_some(), ERebalancerCapDoesNotExist);

    let rebalancer_cap = mcms_cap.rebalancer_cap.borrow();

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(rebalancer_cap)],
        &mut stream,
    );

    let amount = bcs_stream::deserialize_u64(&mut stream);
    let to = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    let coin = withdraw_liquidity(state, rebalancer_cap, amount, ctx);
    transfer::public_transfer(coin, to);
}

public fun mcms_transfer_ownership<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (mcms_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        McmsCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"transfer_ownership"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(&mcms_cap.owner_cap)],
        &mut stream,
    );

    let to = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    transfer_ownership(state, &mcms_cap.owner_cap, to, ctx);
}

public fun mcms_execute_ownership_transfer<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    registry: &mut Registry,
    deployer_state: &mut DeployerState,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (mcms_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        McmsCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"execute_ownership_transfer"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(&mcms_cap.owner_cap), object::id_address(state)],
        &mut stream,
    );

    let to = bcs_stream::deserialize_address(&mut stream);
    let package_address = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    let McmsCap<T> { id, owner_cap, rebalancer_cap } = mcms_registry::release_cap(
        registry,
        McmsCallback<T> {},
    );
    assert!(rebalancer_cap.is_none(), ERebalancerCapNotTransferredOut);

    rebalancer_cap.destroy_none();
    object::delete(id);

    if (mcms_deployer::has_upgrade_cap(deployer_state, package_address)) {
        let upgrade_cap = mcms_deployer::release_upgrade_cap(
            deployer_state,
            registry,
            McmsCallback<T> {},
        );
        transfer::public_transfer(upgrade_cap, to);
    };

    execute_ownership_transfer(owner_cap, state, to, ctx);
}

public fun mcms_add_allowed_modules<T>(
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_mcms_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        McmsCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"add_allowed_modules"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addr(object::id_address(registry), &mut stream);

    let new_module_names = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_vector_u8(stream),
    );
    bcs_stream::assert_is_consumed(&stream);

    mcms_registry::add_allowed_modules(registry, McmsCallback<T> {}, new_module_names, ctx);
}

public fun mcms_remove_allowed_modules<T>(
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_mcms_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        McmsCap<T>,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"remove_allowed_modules"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addr(object::id_address(registry), &mut stream);

    let module_names = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_vector_u8(stream),
    );
    bcs_stream::assert_is_consumed(&stream);

    mcms_registry::remove_allowed_modules(registry, McmsCallback<T> {}, module_names, ctx);
}

/// destroy the lock release token pool state and the owner cap, return the remaining balance to the owner
/// this should only be called after unregistering the pool from the token admin registry
public fun destroy_token_pool<T>(
    ref: &mut CCIPObjectRef,
    state: LockReleaseTokenPoolState<T>,
    owner_cap: OwnerCap,
    ctx: &mut TxContext,
): Coin<T> {
    assert!(
        object::id(&owner_cap) == ownable::owner_cap_id(&state.ownable_state),
        EInvalidOwnerCap,
    );
    assert!(
        !token_admin_registry::is_pool_registered(ref, get_token(&state)),
        EPoolStillRegistered,
    );

    let LockReleaseTokenPoolState<T> {
        id: state_id,
        token_pool_state,
        reserve,
        rebalancer_cap_id: _,
        ownable_state,
    } = state;
    token_pool::destroy_token_pool(token_pool_state);
    object::delete(state_id);

    // Destroy ownable state and owner cap using helper functions
    ownable::destroy(ownable_state, owner_cap, ctx);

    reserve
}

public fun mcms_destroy_token_pool<T>(
    ref: &mut CCIPObjectRef,
    state: LockReleaseTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback<T>,
        OwnerCap,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"destroy_token_pool"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(&state)],
        &mut stream,
    );

    let to = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    let owner_cap = mcms_registry::release_cap<McmsCallback<T>, OwnerCap>(
        registry,
        McmsCallback<T> {},
    );

    let reserve = destroy_token_pool(ref, state, owner_cap, ctx);
    transfer::public_transfer(reserve, to);
}

#[test_only]
public fun test_mcms_callback<T>(): McmsCallback<T> {
    McmsCallback<T> {}
}

#[test_only]
public fun create_fake_rebalancer_cap<T>(ctx: &mut TxContext): RebalancerCap<T> {
    RebalancerCap {
        id: object::new(ctx),
    }
}

#[test_only]
public fun mcms_rebalancer_cap_address<T>(mcms_cap: &McmsCap<T>): address {
    object::id_address(mcms_cap.rebalancer_cap.borrow())
}

#[test_only]
public fun test_init(ctx: &mut TxContext) {
    init(LOCK_RELEASE_TOKEN_POOL {}, ctx);
}
