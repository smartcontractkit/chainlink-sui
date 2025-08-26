module lock_release_token_pool::lock_release_token_pool;

use std::string::{Self, String};
use std::type_name::{Self, TypeName};

use sui::address;
use sui::clock::Clock;
use sui::coin::{Self, Coin, CoinMetadata, TreasuryCap};
use sui::package::UpgradeCap;

use ccip::eth_abi;
use ccip::onramp_state_helper as onramp_sh;
use ccip::offramp_state_helper as offramp_sh;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::token_admin_registry;

use ccip_token_pool::token_pool::{Self, TokenPoolState};
use ccip_token_pool::ownable::{Self, OwnerCap, OwnableState};

use mcms::bcs_stream;
use mcms::mcms_registry::{Self, Registry, ExecutingCallbackParams};
use mcms::mcms_deployer::{Self, DeployerState};

#[allow(lint(coin_field))]
public struct LockReleaseTokenPoolState<phantom T> has key {
    id: UID,
    token_pool_state: TokenPoolState,
    reserve: Coin<T>,
    /// the rebalancer is the address that can manage liquidity of the token pool
    rebalancer: address,
    ownable_state: OwnableState,
}

const CLOCK_ADDRESS: address = @0x6;

const EInvalidArguments: u64 = 1;
const ETokenPoolBalanceTooLow: u64 = 2;
const EUnauthorized: u64 = 3;
const EInvalidOwnerCap: u64 = 4;
const EInvalidFunction: u64 = 5;

// ================================================================
// |                             Init                             |
// ================================================================

public fun type_and_version(): String {
    string::utf8(b"LockReleaseTokenPool 1.6.0")
}

public fun initialize<T>(
    ref: &mut CCIPObjectRef,
    coin_metadata: &CoinMetadata<T>,
    treasury_cap: &TreasuryCap<T>,
    token_pool_administrator: address,
    rebalancer: address,
    ctx: &mut TxContext,
) {
    let (_, lock_release_token_pool_state_address, _, type_proof_type_name) =
        initialize_internal(coin_metadata, rebalancer, ctx);

    let type_proof_type_name_address = type_proof_type_name.get_address();
    let lock_release_token_pool_package_id = address::from_ascii_bytes(&type_proof_type_name_address.into_bytes());

    token_admin_registry::register_pool(
        ref,
        treasury_cap,
        coin_metadata,
        lock_release_token_pool_package_id,
        string::utf8(b"lock_release_token_pool"),
        token_pool_administrator,
        vector[CLOCK_ADDRESS, lock_release_token_pool_state_address],
        vector[CLOCK_ADDRESS, lock_release_token_pool_state_address],
        TypeProof {},
    );
}

/// this is used by the ccip admin to initialize the token pool
/// it verifies that the tx signer is the CCIP admin
/// it does not require a treasury cap object
public fun initialize_by_ccip_admin<T>(
    ref: &mut CCIPObjectRef,
    ccip_admin_proof: state_object::CCIPAdminProof,
    coin_metadata: &CoinMetadata<T>,
    token_pool_administrator: address,
    rebalancer: address,
    ctx: &mut TxContext,
) {
    let (coin_metadata_address, lock_release_token_pool_state_address, token_type, type_proof_type_name) =
        initialize_internal(coin_metadata, rebalancer, ctx);

    let type_proof_type_name_address = type_proof_type_name.get_address();
    let lock_release_token_pool_package_id = address::from_ascii_bytes(&type_proof_type_name_address.into_bytes());

    token_admin_registry::register_pool_by_admin(
        ref,
        ccip_admin_proof,
        coin_metadata_address,
        lock_release_token_pool_package_id,
        string::utf8(b"lock_release_token_pool"),
        token_type.into_string(),
        token_pool_administrator,
        type_proof_type_name.into_string(),
        vector[CLOCK_ADDRESS, lock_release_token_pool_state_address],
        vector[CLOCK_ADDRESS, lock_release_token_pool_state_address],
        ctx,
    );
}

#[allow(lint(self_transfer))]
fun initialize_internal<T>(
    coin_metadata: &CoinMetadata<T>,
    rebalancer: address,
    ctx: &mut TxContext,
): (address, address, TypeName, TypeName) {
    let coin_metadata_address: address = object::id_to_address(&object::id(coin_metadata));
    let (ownable_state, owner_cap) = ownable::new(ctx);

    let mut lock_release_token_pool = LockReleaseTokenPoolState<T> {
        id: object::new(ctx),
        token_pool_state: token_pool::initialize(coin_metadata_address, coin_metadata.get_decimals(), vector[], ctx),
        reserve: coin::zero<T>(ctx),
        rebalancer: @0x0,
        ownable_state,
    };
    set_rebalancer_internal(&mut lock_release_token_pool, rebalancer);
    let type_proof_type_name = type_name::get<TypeProof>();
    let token_type = type_name::get<T>();
    let token_pool_state_address = object::uid_to_address(&lock_release_token_pool.id);

    transfer::share_object(lock_release_token_pool);
    transfer::public_transfer(owner_cap, ctx.sender());

    (coin_metadata_address, token_pool_state_address, token_type, type_proof_type_name)
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

public fun get_remote_pools<T>(
    state: &LockReleaseTokenPoolState<T>,
    remote_chain_selector: u64
): vector<vector<u8>> {
    token_pool::get_remote_pools(&state.token_pool_state, remote_chain_selector)
}

public fun is_remote_pool<T>(
    state: &LockReleaseTokenPoolState<T>,
    remote_chain_selector: u64,
    remote_pool_address: vector<u8>
): bool {
    token_pool::is_remote_pool(
        &state.token_pool_state,
        remote_chain_selector,
        remote_pool_address
    )
}

public fun get_remote_token<T>(
    state: &LockReleaseTokenPoolState<T>, remote_chain_selector: u64
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
        &mut state.token_pool_state, remote_chain_selector, remote_pool_address
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
        &mut state.token_pool_state, remote_chain_selector, remote_pool_address
    );
}

public fun is_supported_chain<T>(
    state: &LockReleaseTokenPoolState<T>,
    remote_chain_selector: u64
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
    remote_token_addresses_to_add: vector<vector<u8>>
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    token_pool::apply_chain_updates(
        &mut state.token_pool_state,
        remote_chain_selectors_to_remove,
        remote_chain_selectors_to_add,
        remote_pool_addresses_to_add,
        remote_token_addresses_to_add
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
    enabled: bool
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    token_pool::set_allowlist_enabled(&mut state.token_pool_state, enabled);
}

public fun apply_allowlist_updates<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    owner_cap: &OwnerCap,
    removes: vector<address>,
    adds: vector<address>
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    token_pool::apply_allowlist_updates(&mut state.token_pool_state, removes, adds);
}

// ================================================================
// |                         Burn/Mint                            |
// ================================================================

public struct TypeProof has drop {}

public fun lock_or_burn<T: drop>(
    ref: &CCIPObjectRef,
    token_transfer_params: &mut onramp_sh::TokenTransferParams,
    c: Coin<T>,
    remote_chain_selector: u64,
    clock: &Clock,
    state: &mut LockReleaseTokenPoolState<T>,
    ctx: &mut TxContext
) {
    let amount = c.value();
    let sender = ctx.sender();

    // This function validates various aspects of the lock or burn operation. If any of the validations fail, the transaction will abort.
    let dest_token_address = token_pool::get_remote_token(&state.token_pool_state, remote_chain_selector);
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

    token_pool::emit_locked_or_burned(&mut state.token_pool_state, amount, remote_chain_selector);

    onramp_sh::add_token_transfer_param(
        ref,
        token_transfer_params,
        remote_chain_selector,
        amount,
        get_token(state),
        dest_token_address,
        extra_data,
        TypeProof {}
    )
}

/// after releasing the token, this function will mark this particular token transfer as complete
/// and set the local amount of this token transfer according to the balance of coin object.
/// a token pool cannot update token transfer item for other tokens simply by changing the
/// index because each token transfer is protected by a type proof
public fun release_or_mint<T>(
    ref: &CCIPObjectRef,
    receiver_params: &mut offramp_sh::ReceiverParams,
    token_transfer: offramp_sh::DestTokenTransfer,
    clock: &Clock,
    state: &mut LockReleaseTokenPoolState<T>,
    ctx: &mut TxContext
) {
    let (receiver, remote_chain_selector, source_amount, dest_token_address, _, source_pool_address, source_pool_data, _) = offramp_sh::get_dest_token_transfer_data(token_transfer);

    let local_amount = token_pool::calculate_release_or_mint_amount(
        &state.token_pool_state,
        source_pool_data,
        source_amount
    );

    token_pool::validate_release_or_mint(
        ref,
        clock,
        &mut state.token_pool_state,
        remote_chain_selector,
        dest_token_address,
        source_pool_address,
        local_amount
    );

    // split the coin to be released
    assert!(
        state.reserve.value() >= local_amount,
        ETokenPoolBalanceTooLow
    );
    let c: Coin<T> = coin::split(&mut state.reserve, local_amount, ctx);

    token_pool::emit_released_or_minted(
        &mut state.token_pool_state,
        receiver,
        local_amount,
        remote_chain_selector,
    );
    transfer::public_transfer(c, receiver);

    offramp_sh::complete_token_transfer(
        ref,
        receiver_params,
        receiver,
        dest_token_address,
         object::id(state),
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
    inbound_rates: vector<u64>
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
        EInvalidArguments
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
    inbound_rate: u64
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
        inbound_rate
    );
}

// ================================================================
// |                    Liquidity Management                      |
// ================================================================

public fun provide_liquidity<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    c: Coin<T>,
    ctx: &mut TxContext,
) {
    assert!(ctx.sender() == state.rebalancer, EUnauthorized);
    let amount = c.value();

    coin::join(&mut state.reserve, c);

    token_pool::emit_liquidity_added(
        &mut state.token_pool_state, state.rebalancer, amount
    );
}

public fun withdraw_liquidity<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    amount: u64,
    ctx: &mut TxContext,
): Coin<T> {
    assert!(ctx.sender() == state.rebalancer, EUnauthorized);

    assert!(state.reserve.value() >= amount, ETokenPoolBalanceTooLow);

    token_pool::emit_liquidity_removed(&mut state.token_pool_state, state.rebalancer, amount);
    coin::split(&mut state.reserve, amount, ctx)
}

public fun set_rebalancer<T>(
    owner_cap: &OwnerCap,
    state: &mut LockReleaseTokenPoolState<T>,
    rebalancer: address,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    set_rebalancer_internal(state, rebalancer);
}

fun set_rebalancer_internal<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    rebalancer: address,
) {
    token_pool::emit_rebalancer_set(
        &mut state.token_pool_state,
        state.rebalancer,
        rebalancer,
    );
    state.rebalancer = rebalancer;
}

public fun get_rebalancer<T>(state: &LockReleaseTokenPoolState<T>): address {
    state.rebalancer
}

public fun get_balance<T>(state: &LockReleaseTokenPoolState<T>): u64 {
    state.reserve.value()
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

public fun accept_ownership<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    ctx: &mut TxContext,
) {
    ownable::accept_ownership(&mut state.ownable_state, ctx);
}

public fun accept_ownership_from_object<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    from: &mut UID,
    ctx: &mut TxContext,
) {
    ownable::accept_ownership_from_object(&mut state.ownable_state, from, ctx);
}

/// Cannot call through `mcms_entrypoint` as owner cap is not registered with MCMS registry
public fun accept_ownership_as_mcms<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
let (_, _, function_name, data) = mcms_registry::get_callback_params_for_mcms(params, McmsCallback<T>{});
    assert!(function_name == string::utf8(b"accept_ownership_as_mcms"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    let mcms = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    ownable::accept_ownership_as_mcms(&mut state.ownable_state, mcms, ctx);
}

public fun execute_ownership_transfer(
    owner_cap: OwnerCap,
    ownable_state: &mut OwnableState,
    to: address,
    ctx: &mut TxContext,
) {
    ownable::execute_ownership_transfer(owner_cap, ownable_state, to, ctx);
}

public fun execute_ownership_transfer_to_mcms<T>(
    owner_cap: OwnerCap,
    state: &mut LockReleaseTokenPoolState<T>,
    registry: &mut Registry,
    to: address,
    ctx: &mut TxContext,
) {
    ownable::execute_ownership_transfer_to_mcms(
        owner_cap,
        &mut state.ownable_state,
        registry,
        to,
        McmsCallback<T>{},
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

public fun mcms_entrypoint<T>(
    state: &mut LockReleaseTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams, // hot potato
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<
        McmsCallback<T>,
        OwnerCap,
    >(
        registry,
        McmsCallback<T>{},
        params,
    );

    let function_bytes = *function.as_bytes();
    let mut stream = bcs_stream::new(data);

    if (function_bytes == b"set_rebalancer") {
        let rebalancer = bcs_stream::deserialize_address(&mut stream);
        bcs_stream::assert_is_consumed(&stream);
        set_rebalancer(owner_cap, state, rebalancer);
    } else if (function_bytes == b"set_allowlist_enabled") {
        let enabled = bcs_stream::deserialize_bool(&mut stream);
        bcs_stream::assert_is_consumed(&stream);
        set_allowlist_enabled(state, owner_cap, enabled);
    } else if (function_bytes == b"apply_allowlist_updates") {
        let removes = bcs_stream::deserialize_vector!(
            &mut stream,
            |stream| bcs_stream::deserialize_address(stream)
        );
        let adds = bcs_stream::deserialize_vector!(
            &mut stream,
            |stream| bcs_stream::deserialize_address(stream)
        );
        bcs_stream::assert_is_consumed(&stream);
        apply_allowlist_updates(state, owner_cap, removes, adds);
    } else if (function_bytes == b"apply_chain_updates") {
        let remote_chain_selectors_to_remove = bcs_stream::deserialize_vector!(
            &mut stream,
            |stream| bcs_stream::deserialize_u64(stream)
        );
        let remote_chain_selectors_to_add = bcs_stream::deserialize_vector!(
            &mut stream,
            |stream| bcs_stream::deserialize_u64(stream)
        );
        let remote_pool_addresses_to_add = bcs_stream::deserialize_vector!(
            &mut stream,
            |stream| bcs_stream::deserialize_vector!(
                stream,
                |stream| bcs_stream::deserialize_vector_u8(stream)
            )
        );
        let remote_token_addresses_to_add = bcs_stream::deserialize_vector!(
            &mut stream,
            |stream| bcs_stream::deserialize_vector_u8(stream)
        );
        bcs_stream::assert_is_consumed(&stream);
        apply_chain_updates(
            state,
            owner_cap,
            remote_chain_selectors_to_remove,
            remote_chain_selectors_to_add,
            remote_pool_addresses_to_add,
            remote_token_addresses_to_add
        );
    } else if (function_bytes == b"transfer_ownership") {
        let to = bcs_stream::deserialize_address(&mut stream);
        bcs_stream::assert_is_consumed(&stream);
        transfer_ownership(state, owner_cap, to, ctx);
    } else if (function_bytes == b"execute_ownership_transfer") {
        let to = bcs_stream::deserialize_address(&mut stream);
        bcs_stream::assert_is_consumed(&stream);
        let owner_cap = mcms_registry::release_cap(registry, McmsCallback<T>{});
        execute_ownership_transfer(owner_cap, &mut state.ownable_state, to, ctx);
    } else {
        abort EInvalidFunction
    };
}

/// destroy the lock release token pool state and the owner cap, return the remaining balance to the owner
/// this should only be called after unregistering the pool from the token admin registry
public fun destroy_token_pool<T>(
    state: LockReleaseTokenPoolState<T>,
    owner_cap: OwnerCap,
    ctx: &mut TxContext,
): Coin<T> {
    assert!(object::id(&owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);

    let LockReleaseTokenPoolState<T> {
        id: state_id,
        token_pool_state,
        reserve,
        rebalancer: _,
        ownable_state,
    } = state;
    token_pool::destroy_token_pool(token_pool_state);
    object::delete(state_id);

    // Destroy ownable state and owner cap using helper functions
    ownable::destroy_ownable_state(ownable_state, ctx);
    ownable::destroy_owner_cap(owner_cap, ctx);

    reserve
}
