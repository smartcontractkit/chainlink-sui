/// this module must be used in conjunction with the managed token module
/// it will store the mint cap object within the token pool state
/// the mint cap object is used to mint/burn the token on the managed token module
module managed_token_pool::managed_token_pool;

use ccip::eth_abi;
use ccip::offramp_state_helper as offramp_sh;
use ccip::onramp_state_helper as onramp_sh;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::token_admin_registry;
use ccip_token_pool::ownable::{Self, OwnerCap, OwnableState};
use ccip_token_pool::token_pool::{Self, TokenPoolState};
use managed_token::managed_token::{Self, TokenState, MintCap};
use managed_token::ownable::OwnerCap as ManagedTokenOwnerCap;
use mcms::bcs_stream;
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, Registry, ExecutingCallbackParams};
use std::string::{Self, String};
use std::type_name::{Self, TypeName};
use sui::address;
use sui::clock::Clock;
use sui::coin::{Coin, CoinMetadata};
use sui::deny_list::DenyList;
use sui::package::UpgradeCap;

public struct ManagedTokenPoolState<phantom T> has key {
    id: UID,
    token_pool_state: TokenPoolState,
    mint_cap: MintCap<T>,
    ownable_state: OwnableState,
}

const EInvalidArguments: u64 = 1;
const EInvalidOwnerCap: u64 = 2;
const EInvalidFunction: u64 = 3;
const EInvalidStateAddress: u64 = 4;
const EInvalidRegistryAddress: u64 = 5;

const CLOCK_ADDRESS: address = @0x6;
const DENY_LIST_ADDRESS: address = @0x403;

// ================================================================
// |                             Init                             |
// ================================================================

public fun type_and_version(): String {
    string::utf8(b"ManagedTokenPool 1.6.0")
}

/// Initialize token pool for a managed token
/// This function works with any existing managed token by:
/// 1. Getting the treasury cap reference for registration
/// 2. Creating the token pool state
/// 3. Registering the pool with the token admin registry
/// Note: The mint_cap must be created beforehand through managed_token::configure_new_minter
public fun initialize_with_managed_token<T>(
    ref: &mut CCIPObjectRef,
    managed_token_state: &TokenState<T>,
    owner_cap: &ManagedTokenOwnerCap<T>,
    coin_metadata: &CoinMetadata<T>,
    mint_cap: MintCap<T>,
    token_pool_administrator: address,
    ctx: &mut TxContext,
) {
    // Get treasury cap reference for registration
    let treasury_cap_ref = managed_token::borrow_treasury_cap(managed_token_state, owner_cap);

    // Initialize the token pool
    let (_, managed_token_pool_state_address, _, type_proof_type_name) = initialize_internal(
        coin_metadata,
        mint_cap,
        ctx,
    );

    let type_proof_type_name_address = type_proof_type_name.get_address();
    let managed_token_pool_package_id = address::from_ascii_bytes(
        &type_proof_type_name_address.into_bytes(),
    );

    // Register the pool with the token admin registry
    token_admin_registry::register_pool(
        ref,
        treasury_cap_ref,
        coin_metadata,
        managed_token_pool_package_id,
        string::utf8(b"managed_token_pool"),
        token_pool_administrator,
        vector[
            CLOCK_ADDRESS,
            DENY_LIST_ADDRESS,
            object::id_to_address(&object::id(managed_token_state)),
            managed_token_pool_state_address,
        ],
        vector[
            CLOCK_ADDRESS,
            DENY_LIST_ADDRESS,
            object::id_to_address(&object::id(managed_token_state)),
            managed_token_pool_state_address,
        ],
        TypeProof {},
    );
}

public fun initialize_by_ccip_admin<T>(
    ref: &mut CCIPObjectRef,
    ccip_admin_proof: state_object::CCIPAdminProof,
    coin_metadata: &CoinMetadata<T>,
    mint_cap: MintCap<T>,
    managed_token_state: address,
    token_pool_administrator: address,
    ctx: &mut TxContext,
) {
    let (
        coin_metadata_address,
        managed_token_pool_state_address,
        token_type,
        type_proof_type_name,
    ) = initialize_internal(coin_metadata, mint_cap, ctx);

    let type_proof_type_name_address = type_proof_type_name.get_address();
    let managed_token_pool_package_id = address::from_ascii_bytes(
        &type_proof_type_name_address.into_bytes(),
    );

    token_admin_registry::register_pool_by_admin(
        ref,
        ccip_admin_proof,
        coin_metadata_address,
        managed_token_pool_package_id,
        string::utf8(b"managed_token_pool"),
        token_type.into_string(),
        token_pool_administrator,
        type_proof_type_name.into_string(),
        vector[
            CLOCK_ADDRESS,
            DENY_LIST_ADDRESS,
            managed_token_state,
            managed_token_pool_state_address,
        ],
        vector[
            CLOCK_ADDRESS,
            DENY_LIST_ADDRESS,
            managed_token_state,
            managed_token_pool_state_address,
        ],
        ctx,
    );
}

#[allow(lint(self_transfer))]
fun initialize_internal<T>(
    coin_metadata: &CoinMetadata<T>,
    mint_cap: MintCap<T>,
    ctx: &mut TxContext,
): (address, address, TypeName, TypeName) {
    let coin_metadata_address: address = object::id_to_address(&object::id(coin_metadata));
    let (ownable_state, owner_cap) = ownable::new(ctx);

    let managed_token_pool = ManagedTokenPoolState<T> {
        id: object::new(ctx),
        token_pool_state: token_pool::initialize(
            coin_metadata_address,
            coin_metadata.get_decimals(),
            vector[],
            ctx,
        ),
        mint_cap,
        ownable_state,
    };
    let type_proof_type_name = type_name::get<TypeProof>();
    let token_type = type_name::get<T>();
    let managed_token_pool_state_address = object::uid_to_address(&managed_token_pool.id);

    transfer::share_object(managed_token_pool);
    transfer::public_transfer(owner_cap, ctx.sender());

    (coin_metadata_address, managed_token_pool_state_address, token_type, type_proof_type_name)
}

public fun add_remote_pool<T>(
    state: &mut ManagedTokenPoolState<T>,
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
    state: &mut ManagedTokenPoolState<T>,
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
    state: &ManagedTokenPoolState<T>,
    remote_chain_selector: u64,
): bool {
    token_pool::is_supported_chain(&state.token_pool_state, remote_chain_selector)
}

public fun get_supported_chains<T>(state: &ManagedTokenPoolState<T>): vector<u64> {
    token_pool::get_supported_chains(&state.token_pool_state)
}

public fun apply_chain_updates<T>(
    state: &mut ManagedTokenPoolState<T>,
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

public fun get_allowlist_enabled<T>(state: &ManagedTokenPoolState<T>): bool {
    token_pool::get_allowlist_enabled(&state.token_pool_state)
}

public fun get_allowlist<T>(state: &ManagedTokenPoolState<T>): vector<address> {
    token_pool::get_allowlist(&state.token_pool_state)
}

public fun set_allowlist_enabled<T>(
    state: &mut ManagedTokenPoolState<T>,
    owner_cap: &OwnerCap,
    enabled: bool,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    token_pool::set_allowlist_enabled(&mut state.token_pool_state, enabled);
}

public fun apply_allowlist_updates<T>(
    state: &mut ManagedTokenPoolState<T>,
    owner_cap: &OwnerCap,
    removes: vector<address>,
    adds: vector<address>,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    token_pool::apply_allowlist_updates(&mut state.token_pool_state, removes, adds);
}

// ================================================================
// |                 Exposing token_pool functions                |
// ================================================================

/// returns the coin metadata object id of the token
public fun get_token<T>(state: &ManagedTokenPoolState<T>): address {
    token_pool::get_token(&state.token_pool_state)
}

public fun get_token_decimals<T>(state: &ManagedTokenPoolState<T>): u8 {
    state.token_pool_state.get_local_decimals()
}

public fun get_remote_pools<T>(
    state: &ManagedTokenPoolState<T>,
    remote_chain_selector: u64,
): vector<vector<u8>> {
    token_pool::get_remote_pools(&state.token_pool_state, remote_chain_selector)
}

public fun is_remote_pool<T>(
    state: &ManagedTokenPoolState<T>,
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
    state: &ManagedTokenPoolState<T>,
    remote_chain_selector: u64,
): vector<u8> {
    token_pool::get_remote_token(&state.token_pool_state, remote_chain_selector)
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
    deny_list: &DenyList,
    token_state: &mut TokenState<T>,
    state: &mut ManagedTokenPoolState<T>,
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

    managed_token::burn(
        token_state,
        &state.mint_cap,
        deny_list,
        c,
        sender,
        ctx,
    );

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

/// after releasing the token, this function will mark this particular token transfer as complete
/// and set the local amount of this token transfer according to the balance of coin object.
/// a token pool cannot update token transfer item for other tokens simply by changing the
/// index because each token transfer is protected by a type proof
public fun release_or_mint<T>(
    ref: &CCIPObjectRef,
    receiver_params: &mut offramp_sh::ReceiverParams,
    clock: &Clock,
    deny_list: &DenyList,
    token_state: &mut TokenState<T>,
    state: &mut ManagedTokenPoolState<T>,
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

    let c: Coin<T> = managed_token::mint(
        token_state,
        &state.mint_cap,
        deny_list,
        local_amount,
        token_receiver,
        ctx,
    );

    token_pool::emit_released_or_minted(
        &mut state.token_pool_state,
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
    state: &mut ManagedTokenPoolState<T>,
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
    state: &mut ManagedTokenPoolState<T>,
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

// ================================================================
// |                      Ownable Functions                       |
// ================================================================

public fun owner<T>(state: &ManagedTokenPoolState<T>): address {
    ownable::owner(&state.ownable_state)
}

public fun has_pending_transfer<T>(state: &ManagedTokenPoolState<T>): bool {
    ownable::has_pending_transfer(&state.ownable_state)
}

public fun pending_transfer_from<T>(state: &ManagedTokenPoolState<T>): Option<address> {
    ownable::pending_transfer_from(&state.ownable_state)
}

public fun pending_transfer_to<T>(state: &ManagedTokenPoolState<T>): Option<address> {
    ownable::pending_transfer_to(&state.ownable_state)
}

public fun pending_transfer_accepted<T>(state: &ManagedTokenPoolState<T>): Option<bool> {
    ownable::pending_transfer_accepted(&state.ownable_state)
}

public fun transfer_ownership<T>(
    state: &mut ManagedTokenPoolState<T>,
    owner_cap: &OwnerCap,
    new_owner: address,
    ctx: &mut TxContext,
) {
    ownable::transfer_ownership(owner_cap, &mut state.ownable_state, new_owner, ctx);
}

public fun accept_ownership<T>(state: &mut ManagedTokenPoolState<T>, ctx: &mut TxContext) {
    ownable::accept_ownership(&mut state.ownable_state, ctx);
}

public fun accept_ownership_from_object<T>(
    state: &mut ManagedTokenPoolState<T>,
    from: &mut UID,
    ctx: &mut TxContext,
) {
    ownable::accept_ownership_from_object(&mut state.ownable_state, from, ctx);
}

public fun mcms_accept_ownership<T>(
    state: &mut ManagedTokenPoolState<T>,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_, _, function, data) = mcms_registry::get_callback_params_for_mcms(
        params,
        McmsCallback<T> {},
    );
    assert!(function == string::utf8(b"mcms_accept_ownership"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    let state_address = bcs_stream::deserialize_address(&mut stream);
    assert!(state_address == object::id_address(state), EInvalidStateAddress);

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
    state: &mut ManagedTokenPoolState<T>,
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

fun validate_shared_objects<T>(
    state: &ManagedTokenPoolState<T>,
    registry: &Registry,
    stream: &mut bcs_stream::BCSStream,
) {
    let state_address = bcs_stream::deserialize_address(stream);
    assert!(state_address == object::id_address(state), EInvalidStateAddress);
    let registry_address = bcs_stream::deserialize_address(stream);
    assert!(registry_address == object::id_address(registry), EInvalidRegistryAddress);
}

public fun mcms_set_allowlist_enabled<T>(
    state: &mut ManagedTokenPoolState<T>,
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
    validate_shared_objects(state, registry, &mut stream);

    let enabled = bcs_stream::deserialize_bool(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    set_allowlist_enabled(state, owner_cap, enabled);
}

public fun mcms_apply_allowlist_updates<T>(
    state: &mut ManagedTokenPoolState<T>,
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
    validate_shared_objects(state, registry, &mut stream);

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
    state: &mut ManagedTokenPoolState<T>,
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
    validate_shared_objects(state, registry, &mut stream);

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

public fun mcms_transfer_ownership<T>(
    state: &mut ManagedTokenPoolState<T>,
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
    validate_shared_objects(state, registry, &mut stream);

    let to = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    transfer_ownership(state, owner_cap, to, ctx);
}

public fun mcms_mcms_accept_ownership<T>(
    state: &mut ManagedTokenPoolState<T>,
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
    assert!(function == string::utf8(b"mcms_accept_ownership"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    validate_shared_objects(state, registry, &mut stream);

    let mcms = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    ownable::mcms_accept_ownership(&mut state.ownable_state, mcms, ctx);
}

public fun mcms_execute_ownership_transfer<T>(
    state: &mut ManagedTokenPoolState<T>,
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
    validate_shared_objects(state, registry, &mut stream);

    let to = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    let owner_cap: OwnerCap = mcms_registry::release_cap(registry, McmsCallback<T> {});
    execute_ownership_transfer(owner_cap, &mut state.ownable_state, to, ctx);
}

/// destroy the managed token pool state and the owner cap, return the mint cap to the owner
/// this should only be called after unregistering the pool from the token admin registry
public fun destroy_token_pool<T>(
    state: ManagedTokenPoolState<T>,
    owner_cap: OwnerCap,
    ctx: &mut TxContext,
): MintCap<T> {
    assert!(
        object::id(&owner_cap) == ownable::owner_cap_id(&state.ownable_state),
        EInvalidOwnerCap,
    );

    let ManagedTokenPoolState<T> {
        id: state_id,
        token_pool_state,
        mint_cap,
        ownable_state,
    } = state;
    token_pool::destroy_token_pool(token_pool_state);
    object::delete(state_id);

    // Destroy ownable state and owner cap using helper functions
    ownable::destroy_ownable_state(ownable_state, ctx);
    ownable::destroy_owner_cap(owner_cap, ctx);

    mint_cap
}
