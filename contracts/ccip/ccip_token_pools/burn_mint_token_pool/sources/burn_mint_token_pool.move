/// this module will store the treasury cap object within the token pool state
/// this will disable any burning/minting of the token outside of the token pool
/// if this is not desired, consider using the lock release token pool or the
/// combination of the managed token and managed token pool
module burn_mint_token_pool::burn_mint_token_pool;

use ccip::bcs_helper;
use ccip::eth_abi;
use ccip::offramp_state_helper as offramp_sh;
use ccip::onramp_state_helper as onramp_sh;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::token_admin_registry;
use ccip_token_pool::ownable::{Self, OwnerCap, OwnableState};
use ccip_token_pool::token_pool::{Self, TokenPoolState};
use mcms::bcs_stream;
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, Registry, ExecutingCallbackParams};
use std::string::{Self, String};
use std::type_name::{Self, TypeName};
use sui::address;
use sui::clock::Clock;
use sui::coin::{Self, Coin, CoinMetadata, TreasuryCap};
use sui::package::UpgradeCap;

public struct BurnMintTokenPoolState<phantom T> has key {
    id: UID,
    token_pool_state: TokenPoolState,
    treasury_cap: TreasuryCap<T>,
    ownable_state: OwnableState,
}

const EInvalidArguments: u64 = 1;
const EInvalidOwnerCap: u64 = 2;
const EInvalidFunction: u64 = 3;

const CLOCK_ADDRESS: address = @0x6;

// ================================================================
// |                             Init                             |
// ================================================================

public fun type_and_version(): String {
    string::utf8(b"BurnMintTokenPool 1.6.0")
}

// coin metadata and decimals can be provided in Move.toml
public fun initialize<T>(
    ref: &mut CCIPObjectRef,
    coin_metadata: &CoinMetadata<T>,
    treasury_cap: TreasuryCap<T>,
    token_pool_administrator: address,
    ctx: &mut TxContext,
) {
    let (_, type_proof_type_name, _, burn_mint_token_pool) = initialize_internal(
        coin_metadata,
        treasury_cap,
        ctx,
    );

    let type_proof_type_name_address = type_proof_type_name.get_address();
    let burn_mint_token_pool_package_id = address::from_ascii_bytes(
        &type_proof_type_name_address.into_bytes(),
    );

    token_admin_registry::register_pool(
        ref,
        &burn_mint_token_pool.treasury_cap,
        coin_metadata,
        burn_mint_token_pool_package_id,
        string::utf8(b"burn_mint_token_pool"),
        token_pool_administrator,
        vector[CLOCK_ADDRESS, object::uid_to_address(&burn_mint_token_pool.id)],
        vector[CLOCK_ADDRESS, object::uid_to_address(&burn_mint_token_pool.id)],
        TypeProof {},
    );

    transfer::share_object(burn_mint_token_pool);
}

public fun initialize_by_ccip_admin<T>(
    ref: &mut CCIPObjectRef,
    ccip_admin_proof: state_object::CCIPAdminProof,
    coin_metadata: &CoinMetadata<T>,
    treasury_cap: TreasuryCap<T>,
    token_pool_administrator: address,
    ctx: &mut TxContext,
) {
    let (
        coin_metadata_address,
        type_proof_type_name,
        token_type,
        burn_mint_token_pool,
    ) = initialize_internal(coin_metadata, treasury_cap, ctx);

    let type_proof_type_name_address = type_proof_type_name.get_address();
    let burn_mint_token_pool_package_id = address::from_ascii_bytes(
        &type_proof_type_name_address.into_bytes(),
    );

    token_admin_registry::register_pool_by_admin(
        ref,
        ccip_admin_proof,
        coin_metadata_address,
        burn_mint_token_pool_package_id,
        string::utf8(b"burn_mint_token_pool"),
        token_type.into_string(),
        token_pool_administrator,
        type_proof_type_name.into_string(),
        vector[CLOCK_ADDRESS, object::uid_to_address(&burn_mint_token_pool.id)],
        vector[CLOCK_ADDRESS, object::uid_to_address(&burn_mint_token_pool.id)],
        ctx,
    );

    transfer::share_object(burn_mint_token_pool);
}

#[allow(lint(self_transfer))]
fun initialize_internal<T>(
    coin_metadata: &CoinMetadata<T>,
    treasury_cap: TreasuryCap<T>,
    ctx: &mut TxContext,
): (address, TypeName, TypeName, BurnMintTokenPoolState<T>) {
    let coin_metadata_address: address = object::id_to_address(&object::id(coin_metadata));
    let (ownable_state, owner_cap) = ownable::new(ctx);

    let burn_mint_token_pool = BurnMintTokenPoolState<T> {
        id: object::new(ctx),
        token_pool_state: token_pool::initialize(
            coin_metadata_address,
            coin_metadata.get_decimals(),
            vector[],
            ctx,
        ),
        treasury_cap,
        ownable_state,
    };
    let type_proof_type_name = type_name::get<TypeProof>();
    let token_type = type_name::get<T>();

    transfer::public_transfer(owner_cap, ctx.sender());

    (coin_metadata_address, type_proof_type_name, token_type, burn_mint_token_pool)
}

// ================================================================
// |                 Exposing token_pool functions                |
// ================================================================

// this now returns the address of coin metadata
public fun get_token<T>(state: &BurnMintTokenPoolState<T>): address {
    token_pool::get_token(&state.token_pool_state)
}

public fun get_token_decimals<T>(state: &BurnMintTokenPoolState<T>): u8 {
    state.token_pool_state.get_local_decimals()
}

public fun get_remote_pools<T>(
    state: &BurnMintTokenPoolState<T>,
    remote_chain_selector: u64,
): vector<vector<u8>> {
    token_pool::get_remote_pools(&state.token_pool_state, remote_chain_selector)
}

public fun is_remote_pool<T>(
    state: &BurnMintTokenPoolState<T>,
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
    state: &BurnMintTokenPoolState<T>,
    remote_chain_selector: u64,
): vector<u8> {
    token_pool::get_remote_token(&state.token_pool_state, remote_chain_selector)
}

public fun add_remote_pool<T>(
    state: &mut BurnMintTokenPoolState<T>,
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
    state: &mut BurnMintTokenPoolState<T>,
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
    state: &BurnMintTokenPoolState<T>,
    remote_chain_selector: u64,
): bool {
    token_pool::is_supported_chain(&state.token_pool_state, remote_chain_selector)
}

public fun get_supported_chains<T>(state: &BurnMintTokenPoolState<T>): vector<u64> {
    token_pool::get_supported_chains(&state.token_pool_state)
}

public fun apply_chain_updates<T>(
    state: &mut BurnMintTokenPoolState<T>,
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

public fun get_allowlist_enabled<T>(state: &BurnMintTokenPoolState<T>): bool {
    token_pool::get_allowlist_enabled(&state.token_pool_state)
}

public fun get_allowlist<T>(state: &BurnMintTokenPoolState<T>): vector<address> {
    token_pool::get_allowlist(&state.token_pool_state)
}

public fun set_allowlist_enabled<T>(
    state: &mut BurnMintTokenPoolState<T>,
    owner_cap: &OwnerCap,
    enabled: bool,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    token_pool::set_allowlist_enabled(&mut state.token_pool_state, enabled);
}

public fun apply_allowlist_updates<T>(
    state: &mut BurnMintTokenPoolState<T>,
    owner_cap: &OwnerCap,
    removes: vector<address>,
    adds: vector<address>,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    token_pool::apply_allowlist_updates(&mut state.token_pool_state, removes, adds);
}

// ================================================================
// |                         Burn/Mint                            |
// ================================================================

public struct TypeProof has drop {}

public fun lock_or_burn<T>(
    ref: &CCIPObjectRef,
    token_transfer_params: &mut onramp_sh::TokenTransferParams,
    c: Coin<T>,
    remote_chain_selector: u64,
    clock: &Clock,
    state: &mut BurnMintTokenPoolState<T>,
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

    coin::burn(&mut state.treasury_cap, c);

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
        TypeProof {},
    )
}

public fun release_or_mint<T>(
    ref: &CCIPObjectRef,
    receiver_params: &mut offramp_sh::ReceiverParams,
    clock: &Clock,
    pool: &mut BurnMintTokenPoolState<T>,
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

    let local_amount = token_pool::calculate_release_or_mint_amount(
        &pool.token_pool_state,
        source_pool_data,
        source_amount,
    );

    token_pool::validate_release_or_mint(
        ref,
        clock,
        &mut pool.token_pool_state,
        remote_chain_selector,
        dest_token_address,
        source_pool_address,
        local_amount,
    );

    let c = coin::mint(
        &mut pool.treasury_cap,
        local_amount,
        ctx,
    );

    token_pool::emit_released_or_minted(
        &mut pool.token_pool_state,
        token_receiver,
        local_amount,
        remote_chain_selector,
    );
    transfer::public_transfer(c, token_receiver);

    offramp_sh::complete_token_transfer(
        ref,
        receiver_params,
        token_receiver,
        dest_token_address,
        TypeProof {},
    );
}

// ================================================================
// |                    Rate limit config                         |
// ================================================================

public fun set_chain_rate_limiter_configs<T>(
    state: &mut BurnMintTokenPoolState<T>,
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
    state: &mut BurnMintTokenPoolState<T>,
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

// destroy the burn mint token pool state and the owner cap, return the treasury cap to the owner
// this should only be called after unregistering the pool from the token admin registry
public fun destroy_token_pool<T>(
    state: BurnMintTokenPoolState<T>,
    owner_cap: OwnerCap,
    ctx: &mut TxContext,
): TreasuryCap<T> {
    assert!(
        object::id(&owner_cap) == ownable::owner_cap_id(&state.ownable_state),
        EInvalidOwnerCap,
    );

    let BurnMintTokenPoolState<T> {
        id: state_id,
        token_pool_state,
        treasury_cap,
        ownable_state,
    } = state;
    token_pool::destroy_token_pool(token_pool_state);
    object::delete(state_id);

    ownable::destroy_ownable_state(ownable_state, ctx);
    ownable::destroy_owner_cap(owner_cap, ctx);

    treasury_cap
}

// ================================================================
// |                      Ownable Functions                       |
// ================================================================

public fun owner<T>(state: &BurnMintTokenPoolState<T>): address {
    ownable::owner(&state.ownable_state)
}

public fun has_pending_transfer<T>(state: &BurnMintTokenPoolState<T>): bool {
    ownable::has_pending_transfer(&state.ownable_state)
}

public fun pending_transfer_from<T>(state: &BurnMintTokenPoolState<T>): Option<address> {
    ownable::pending_transfer_from(&state.ownable_state)
}

public fun pending_transfer_to<T>(state: &BurnMintTokenPoolState<T>): Option<address> {
    ownable::pending_transfer_to(&state.ownable_state)
}

public fun pending_transfer_accepted<T>(state: &BurnMintTokenPoolState<T>): Option<bool> {
    ownable::pending_transfer_accepted(&state.ownable_state)
}

public fun transfer_ownership<T>(
    state: &mut BurnMintTokenPoolState<T>,
    owner_cap: &OwnerCap,
    new_owner: address,
    ctx: &mut TxContext,
) {
    ownable::transfer_ownership(owner_cap, &mut state.ownable_state, new_owner, ctx);
}

public fun accept_ownership<T>(state: &mut BurnMintTokenPoolState<T>, ctx: &mut TxContext) {
    ownable::accept_ownership(&mut state.ownable_state, ctx);
}

public fun accept_ownership_from_object<T>(
    state: &mut BurnMintTokenPoolState<T>,
    from: &mut UID,
    ctx: &mut TxContext,
) {
    ownable::accept_ownership_from_object(&mut state.ownable_state, from, ctx);
}

public fun mcms_accept_ownership<T>(
    state: &mut BurnMintTokenPoolState<T>,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_, _, function, data) = mcms_registry::get_callback_params_for_mcms(
        params,
        McmsCallback<T> {},
    );
    assert!(function == string::utf8(b"mcms_accept_ownership"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_helper::validate_obj_addr(object::id_address(state), &mut stream);

    let mcms = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    ownable::mcms_accept_ownership(&mut state.ownable_state, mcms, ctx);
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
    state: &mut BurnMintTokenPoolState<T>,
    registry: &mut Registry,
    to: address,
    ctx: &mut TxContext,
) {
    ownable::execute_ownership_transfer_to_mcms(
        owner_cap,
        &mut state.ownable_state,
        registry,
        to,
        McmsCallback<T> {},
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

public fun mcms_set_allowlist_enabled<T>(
    state: &mut BurnMintTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback<T>, OwnerCap>(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"set_allowlist_enabled"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_helper::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(registry)],
        &mut stream,
    );
    let enabled = bcs_stream::deserialize_bool(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    set_allowlist_enabled(state, owner_cap, enabled);
}

public fun mcms_apply_allowlist_updates<T>(
    state: &mut BurnMintTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback<T>, OwnerCap>(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"apply_allowlist_updates"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_helper::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(registry)],
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

    apply_allowlist_updates(state, owner_cap, removes, adds);
}

public fun mcms_apply_chain_updates<T>(
    state: &mut BurnMintTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback<T>, OwnerCap>(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"apply_chain_updates"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_helper::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(registry)],
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
        owner_cap,
        remote_chain_selectors_to_remove,
        remote_chain_selectors_to_add,
        remote_pool_addresses_to_add,
        remote_token_addresses_to_add,
    );
}

public fun mcms_add_remote_pool<T>(
    state: &mut BurnMintTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback<T>, OwnerCap>(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"add_remote_pool"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_helper::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(registry)],
        &mut stream,
    );
    let remote_chain_selector = bcs_stream::deserialize_u64(&mut stream);
    let remote_pool_address = bcs_stream::deserialize_vector_u8(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    add_remote_pool(state, owner_cap, remote_chain_selector, remote_pool_address);
}

public fun mcms_remove_remote_pool<T>(
    state: &mut BurnMintTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback<T>, OwnerCap>(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"remove_remote_pool"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_helper::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(registry)],
        &mut stream,
    );
    let remote_chain_selector = bcs_stream::deserialize_u64(&mut stream);
    let remote_pool_address = bcs_stream::deserialize_vector_u8(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    remove_remote_pool(state, owner_cap, remote_chain_selector, remote_pool_address);
}

public fun mcms_set_chain_rate_limiter_configs<T>(
    state: &mut BurnMintTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    clock: &Clock,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback<T>, OwnerCap>(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"set_chain_rate_limiter_configs"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_helper::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(registry)],
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
        owner_cap,
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
    state: &mut BurnMintTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    clock: &Clock,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback<T>, OwnerCap>(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"set_chain_rate_limiter_config"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_helper::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(registry)],
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
        owner_cap,
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

public fun mcms_transfer_ownership<T>(
    state: &mut BurnMintTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback<T>, OwnerCap>(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"transfer_ownership"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_helper::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(registry)],
        &mut stream,
    );

    let to = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    transfer_ownership(state, owner_cap, to, ctx);
}

public fun mcms_execute_ownership_transfer<T>(
    state: &mut BurnMintTokenPoolState<T>,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_owner_cap, function, data) = mcms_registry::get_callback_params<
        McmsCallback<T>,
        OwnerCap,
    >(
        registry,
        McmsCallback<T> {},
        params,
    );
    assert!(function == string::utf8(b"execute_ownership_transfer"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_helper::validate_obj_addrs(
        vector[object::id_address(state), object::id_address(registry)],
        &mut stream,
    );

    let to = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    let owner_cap: OwnerCap = mcms_registry::release_cap(registry, McmsCallback<T> {});
    execute_ownership_transfer(owner_cap, &mut state.ownable_state, to, ctx);
}
