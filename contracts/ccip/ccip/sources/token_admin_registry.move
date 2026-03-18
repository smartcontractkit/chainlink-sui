module ccip::token_admin_registry;

use ccip::ownable::{Self, OwnerCap};
use ccip::publisher_wrapper::{Self, PublisherWrapper};
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::upgrade_registry::verify_function_allowed;
use mcms::bcs_stream;
use mcms::mcms_registry::{Self, Registry, ExecutingCallbackParams};
use std::ascii;
use std::string::{Self, String};
use std::type_name;
use sui::coin::{CoinMetadata, TreasuryCap};
use sui::event;
use sui::linked_table::{Self, LinkedTable};

const VERSION: u8 = 1;

public struct TokenAdminRegistryState has key, store {
    id: UID,
    // coin metadata object id -> token config
    token_configs: LinkedTable<address, TokenConfig>,
    token_pool_package_id_to_coin_metadata: LinkedTable<address, address>,
}

public struct TokenConfig has copy, drop, store {
    token_pool_package_id: address,
    token_pool_module: String,
    // the type of the token, this should be the full type name of the token, e.g. "link_package_id::link::LINK"
    token_type: ascii::String,
    administrator: address,
    pending_administrator: address,
    // type proof of the token pool
    token_pool_type_proof: ascii::String,
    lock_or_burn_params: vector<address>,
    release_or_mint_params: vector<address>,
}

public struct PoolSet has copy, drop {
    coin_metadata_address: address,
    previous_pool_package_id: address,
    new_pool_package_id: address,
    // type proof of the new token pool
    token_pool_type_proof: ascii::String,
    lock_or_burn_params: vector<address>,
    release_or_mint_params: vector<address>,
}

public struct PoolRegistered has copy, drop {
    coin_metadata_address: address,
    token_pool_package_id: address,
    administrator: address,
    // type proof of the token pool
    token_pool_type_proof: ascii::String,
}

public struct PoolUnregistered has copy, drop {
    coin_metadata_address: address,
    previous_pool_address: address,
}

public struct AdministratorTransferRequested has copy, drop {
    coin_metadata_address: address,
    current_admin: address,
    new_admin: address,
}

public struct AdministratorTransferred has copy, drop {
    coin_metadata_address: address,
    new_admin: address,
}

const ENotPendingAdministrator: u64 = 1;
const EAlreadyInitialized: u64 = 2;
const ETokenAlreadyRegistered: u64 = 3;
const ETokenNotRegistered: u64 = 4;
const ENotAdministrator: u64 = 5;
const ETokenAddressNotRegistered: u64 = 6;
const ENotAllowed: u64 = 7;
const EInvalidFunction: u64 = 8;
const EInvalidOwnerCap: u64 = 9;
const ETokenPoolPackageIdAlreadyRegistered: u64 = 10;
const ETokenPoolPackageIdNotRegistered: u64 = 11;

public fun type_and_version(): String {
    string::utf8(b"TokenAdminRegistry 1.6.0")
}

public fun initialize(ref: &mut CCIPObjectRef, owner_cap: &OwnerCap, ctx: &mut TxContext) {
    assert!(object::id(owner_cap) == state_object::owner_cap_id(ref), EInvalidOwnerCap);
    assert!(!state_object::contains<TokenAdminRegistryState>(ref), EAlreadyInitialized);
    let state = TokenAdminRegistryState {
        id: object::new(ctx),
        token_configs: linked_table::new(ctx),
        token_pool_package_id_to_coin_metadata: linked_table::new(ctx),
    };

    state_object::add(ref, owner_cap, state, ctx);
}

public fun get_pools(
    ref: &CCIPObjectRef,
    coin_metadata_addresses: vector<address>,
): vector<address> {
    verify_function_allowed(
        ref,
        string::utf8(b"token_admin_registry"),
        string::utf8(b"get_pools"),
        VERSION,
    );
    let state = state_object::borrow<TokenAdminRegistryState>(ref);

    let mut token_pool_package_ids: vector<address> = vector[];
    coin_metadata_addresses.do_ref!(|metadata_address| {
        let metadata_address: address = *metadata_address;
        if (state.token_configs.contains(metadata_address)) {
            let token_config = state.token_configs.borrow(metadata_address);
            token_pool_package_ids.push_back(token_config.token_pool_package_id);
        } else {
            // returns @0x0 for assets without token pools.
            token_pool_package_ids.push_back(@0x0);
        }
    });

    token_pool_package_ids
}

public fun get_pool(ref: &CCIPObjectRef, coin_metadata_address: address): address {
    verify_function_allowed(
        ref,
        string::utf8(b"token_admin_registry"),
        string::utf8(b"get_pool"),
        VERSION,
    );
    let state = state_object::borrow<TokenAdminRegistryState>(ref);

    if (state.token_configs.contains(coin_metadata_address)) {
        let token_config = state.token_configs.borrow(coin_metadata_address);
        token_config.token_pool_package_id
    } else {
        // returns @0x0 for assets without token pools.
        @0x0
    }
}

public fun get_token_config_struct(
    ref: &CCIPObjectRef,
    coin_metadata_address: address,
): TokenConfig {
    verify_function_allowed(
        ref,
        string::utf8(b"token_admin_registry"),
        string::utf8(b"get_token_config_struct"),
        VERSION,
    );
    let state = state_object::borrow<TokenAdminRegistryState>(ref);

    if (state.token_configs.contains(coin_metadata_address)) {
        let token_config = state.token_configs.borrow(coin_metadata_address);
        *token_config
    } else {
        TokenConfig {
            token_pool_package_id: @0x0,
            token_pool_module: string::utf8(b""),
            token_type: ascii::string(b""),
            administrator: @0x0,
            pending_administrator: @0x0,
            token_pool_type_proof: ascii::string(b""),
            lock_or_burn_params: vector[],
            release_or_mint_params: vector[],
        }
    }
}

public fun get_pool_local_token(ref: &CCIPObjectRef, token_pool_package_id: address): address {
    verify_function_allowed(
        ref,
        string::utf8(b"token_admin_registry"),
        string::utf8(b"get_pool_local_token"),
        VERSION,
    );
    let state = state_object::borrow<TokenAdminRegistryState>(ref);

    if (state.token_pool_package_id_to_coin_metadata.contains(token_pool_package_id)) {
        *state.token_pool_package_id_to_coin_metadata.borrow(token_pool_package_id)
    } else {
        @0x0
    }
}

public fun get_token_config(
    ref: &CCIPObjectRef,
    coin_metadata_address: address,
): (address, address, address) {
    verify_function_allowed(
        ref,
        string::utf8(b"token_admin_registry"),
        string::utf8(b"get_token_config"),
        VERSION,
    );
    let state = state_object::borrow<TokenAdminRegistryState>(ref);

    if (state.token_configs.contains(coin_metadata_address)) {
        let token_config = state.token_configs.borrow(coin_metadata_address);
        (
            token_config.token_pool_package_id,
            token_config.administrator,
            token_config.pending_administrator,
        )
    } else {
        (@0x0, @0x0, @0x0)
    }
}

public fun get_token_config_data(
    ref: &CCIPObjectRef,
    coin_metadata_address: address,
): (
    address,
    String,
    ascii::String,
    address,
    address,
    ascii::String,
    vector<address>,
    vector<address>,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"token_admin_registry"),
        string::utf8(b"get_token_config_data"),
        VERSION,
    );
    let state = state_object::borrow<TokenAdminRegistryState>(ref);

    if (state.token_configs.contains(coin_metadata_address)) {
        let token_config = state.token_configs.borrow(coin_metadata_address);
        (
            token_config.token_pool_package_id,
            token_config.token_pool_module,
            token_config.token_type,
            token_config.administrator,
            token_config.pending_administrator,
            token_config.token_pool_type_proof,
            token_config.lock_or_burn_params,
            token_config.release_or_mint_params,
        )
    } else {
        (
            @0x0,
            string::utf8(b""),
            ascii::string(b""),
            @0x0,
            @0x0,
            ascii::string(b""),
            vector[],
            vector[],
        )
    }
}

/// Get configured tokens paginated using a start key and limit.
/// Caller should call this on a certain block to ensure you the same state for every call.
///
/// This function retrieves a batch of token addresses from the registry, starting from
/// the token address that comes after the provided start_key.
///
/// @param ref - Reference to the CCIP state object
/// @param start_key - Address to start pagination from (returns tokens AFTER this address)
///                                empty address @0x0 means start from the beginning
/// @param max_count - Maximum number of tokens to return
///
/// @return:
///   - vector<address>: List of token coin metadata addresses (up to max_count)
///   - address: Next key to use for pagination (pass this as start_key in next call)
///   - bool: Whether there are more tokens after this batch
public fun get_all_configured_tokens(
    ref: &CCIPObjectRef,
    start_key: address,
    max_count: u64,
): (vector<address>, address, bool) {
    verify_function_allowed(
        ref,
        string::utf8(b"token_admin_registry"),
        string::utf8(b"get_all_configured_tokens"),
        VERSION,
    );
    let state = state_object::borrow<TokenAdminRegistryState>(ref);

    let mut i = 0;
    let mut results = vector[];
    let mut key = start_key;
    if (key == @0x0) {
        if (state.token_configs.is_empty()) {
            return (results, key, false)
        };
        if (max_count == 0) {
            return (results, key, true)
        };
        key = *state.token_configs.front().borrow();
        results.push_back(key);
        i = 1;
    } else {
        assert!(state.token_configs.contains(start_key), ETokenAddressNotRegistered);
    };

    while (i < max_count) {
        let next_key_opt = state.token_configs.next(key);
        if (next_key_opt.is_none()) {
            return (results, key, false)
        };

        key = *next_key_opt.borrow();
        results.push_back(key);
        i = i + 1;
    };

    // Check if there are more tokens after the last key
    let has_more = state.token_configs.next(key).is_some();
    (results, key, has_more)
}

// ================================================================
// |                       Register Pool                          |
// ================================================================

/// Only the token owner can call this function to register a token pool for the token it owns.
/// The ownership of the token is proven by the presence of the treasury cap.
/// The publisher wrapper proves that the caller owns the token pool package.
public fun register_pool<T, TypeProof: drop>(
    ref: &mut CCIPObjectRef,
    _: &TreasuryCap<T>, // passing in the treasury cap to demonstrate ownership over the token
    coin_metadata: &CoinMetadata<T>,
    initial_administrator: address,
    lock_or_burn_params: vector<address>,
    release_or_mint_params: vector<address>,
    publisher_wrapper: PublisherWrapper<TypeProof>, // Proves ownership over the token pool package.
    _proof: TypeProof,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"token_admin_registry"),
        string::utf8(b"register_pool"),
        VERSION,
    );

    let package_address = publisher_wrapper::get_package_address(publisher_wrapper);
    let proof_tn = type_name::with_defining_ids<TypeProof>();
    let token_pool_module = proof_tn.module_string().into_bytes().to_string();
    let coin_metadata_address = object::id_address(coin_metadata);
    let token_type = type_name::with_defining_ids<T>().into_string();

    register_pool_internal(
        ref,
        coin_metadata_address,
        package_address,
        token_pool_module,
        token_type,
        initial_administrator,
        proof_tn.into_string(),
        lock_or_burn_params,
        release_or_mint_params,
    );
}

/// Only owner of CCIP can call this function to register a token pool.
public fun register_pool_as_owner(
    owner_cap: &OwnerCap,
    ref: &mut CCIPObjectRef,
    coin_metadata_address: address,
    package_address: address,
    token_pool_module: String,
    token_type: ascii::String,
    initial_administrator: address,
    token_pool_type_proof: ascii::String,
    lock_or_burn_params: vector<address>,
    release_or_mint_params: vector<address>,
    _ctx: &mut TxContext,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"token_admin_registry"),
        string::utf8(b"register_pool_as_owner"),
        VERSION,
    );
    assert!(object::id(owner_cap) == state_object::owner_cap_id(ref), EInvalidOwnerCap);

    register_pool_internal(
        ref,
        coin_metadata_address,
        package_address,
        token_pool_module,
        token_type,
        initial_administrator,
        token_pool_type_proof,
        lock_or_burn_params,
        release_or_mint_params,
    );
}

fun register_pool_internal(
    ref: &mut CCIPObjectRef,
    coin_metadata_address: address,
    token_pool_package_id: address,
    token_pool_module: String,
    token_type: ascii::String,
    initial_administrator: address,
    token_pool_type_proof: ascii::String,
    lock_or_burn_params: vector<address>,
    release_or_mint_params: vector<address>,
) {
    let state = state_object::borrow_mut<TokenAdminRegistryState>(ref);
    assert!(!state.token_configs.contains(coin_metadata_address), ETokenAlreadyRegistered);

    let token_config = TokenConfig {
        token_pool_package_id,
        token_pool_module,
        token_type,
        administrator: initial_administrator,
        pending_administrator: @0x0,
        token_pool_type_proof,
        lock_or_burn_params,
        release_or_mint_params,
    };

    state.token_configs.push_back(coin_metadata_address, token_config);
    assert!(
        !state.token_pool_package_id_to_coin_metadata.contains(token_pool_package_id),
        ETokenPoolPackageIdAlreadyRegistered,
    );
    state
        .token_pool_package_id_to_coin_metadata
        .push_back(token_pool_package_id, coin_metadata_address);

    event::emit(PoolRegistered {
        coin_metadata_address,
        token_pool_package_id,
        administrator: initial_administrator,
        token_pool_type_proof,
    });

    event::emit(PoolSet {
        coin_metadata_address,
        previous_pool_package_id: @0x0,
        new_pool_package_id: token_pool_package_id,
        token_pool_type_proof,
        lock_or_burn_params,
        release_or_mint_params,
    });
}

public fun unregister_pool(
    ref: &mut CCIPObjectRef,
    coin_metadata_address: address,
    ctx: &mut TxContext,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"token_admin_registry"),
        string::utf8(b"unregister_pool"),
        VERSION,
    );
    let state = state_object::borrow_mut<TokenAdminRegistryState>(ref);

    assert!(state.token_configs.contains(coin_metadata_address), ETokenNotRegistered);

    let token_config = state.token_configs.borrow(coin_metadata_address);
    assert!(token_config.administrator == ctx.sender(), ENotAdministrator);

    remove_pool_config(state, coin_metadata_address);
}

fun unregister_pool_via_mcms(
    ref: &mut CCIPObjectRef,
    coin_metadata_address: address,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"token_admin_registry"),
        string::utf8(b"unregister_pool"),
        VERSION,
    );
    let state = state_object::borrow_mut<TokenAdminRegistryState>(ref);

    assert!(state.token_configs.contains(coin_metadata_address), ETokenNotRegistered);

    remove_pool_config(state, coin_metadata_address);
}

fun remove_pool_config(
    state: &mut TokenAdminRegistryState,
    coin_metadata_address: address,
) {
    let token_config = state.token_configs.remove(coin_metadata_address);
    let previous_pool_address = token_config.token_pool_package_id;

    assert!(
        state.token_pool_package_id_to_coin_metadata.contains(previous_pool_address),
        ETokenPoolPackageIdNotRegistered,
    );
    state.token_pool_package_id_to_coin_metadata.remove(previous_pool_address);

    event::emit(PoolUnregistered {
        coin_metadata_address,
        previous_pool_address,
    });
}

public fun transfer_admin_role(
    ref: &mut CCIPObjectRef,
    coin_metadata_address: address,
    new_admin: address,
    ctx: &mut TxContext,
) {
    transfer_admin_role_internal(ref, coin_metadata_address, new_admin, ctx.sender());
}

fun transfer_admin_role_internal(
    ref: &mut CCIPObjectRef,
    coin_metadata_address: address,
    new_admin: address,
    caller: address,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"token_admin_registry"),
        string::utf8(b"transfer_admin_role"),
        VERSION,
    );
    let state = state_object::borrow_mut<TokenAdminRegistryState>(ref);

    assert!(state.token_configs.contains(coin_metadata_address), ETokenNotRegistered);

    let token_config = state.token_configs.borrow_mut(coin_metadata_address);

    assert!(token_config.administrator == caller, ENotAdministrator);

    // can be @0x0 to cancel a pending transfer.
    token_config.pending_administrator = new_admin;

    event::emit(AdministratorTransferRequested {
        coin_metadata_address,
        current_admin: token_config.administrator,
        new_admin,
    });
}

public fun accept_admin_role(
    ref: &mut CCIPObjectRef,
    coin_metadata_address: address,
    ctx: &mut TxContext,
) {
    accept_admin_role_internal(ref, coin_metadata_address, ctx.sender());
}

fun accept_admin_role_internal(
    ref: &mut CCIPObjectRef,
    coin_metadata_address: address,
    caller: address,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"token_admin_registry"),
        string::utf8(b"accept_admin_role"),
        VERSION,
    );
    let state = state_object::borrow_mut<TokenAdminRegistryState>(ref);

    assert!(state.token_configs.contains(coin_metadata_address), ETokenNotRegistered);

    let token_config = state.token_configs.borrow_mut(coin_metadata_address);

    assert!(token_config.pending_administrator == caller, ENotPendingAdministrator);

    token_config.administrator = token_config.pending_administrator;
    token_config.pending_administrator = @0x0;

    event::emit(AdministratorTransferred {
        coin_metadata_address,
        new_admin: token_config.administrator,
    });
}

public fun is_pool_registered(ref: &CCIPObjectRef, coin_metadata_address: address): bool {
    verify_function_allowed(
        ref,
        string::utf8(b"token_admin_registry"),
        string::utf8(b"is_pool_registered"),
        VERSION,
    );
    let state = state_object::borrow<TokenAdminRegistryState>(ref);
    state.token_configs.contains(coin_metadata_address)
}

public fun is_administrator(
    ref: &CCIPObjectRef,
    coin_metadata_address: address,
    administrator: address,
): bool {
    verify_function_allowed(
        ref,
        string::utf8(b"token_admin_registry"),
        string::utf8(b"is_administrator"),
        VERSION,
    );
    let state = state_object::borrow<TokenAdminRegistryState>(ref);

    assert!(state.token_configs.contains(coin_metadata_address), ETokenNotRegistered);

    let token_config = state.token_configs.borrow(coin_metadata_address);
    token_config.administrator == administrator
}

// ================================================================
// |                       MCMS Functions                         |
// ================================================================

/// Only callable once validated by MCMS - `ExecutingCallbackParams` is from MCMS.
/// MCMS needs to know the token pool's type proof string to register the token pool.
public fun mcms_register_pool(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        state_object::McmsCallback,
        OwnerCap,
    >(
        registry,
        state_object::mcms_callback(),
        params,
    );
    assert!(function == string::utf8(b"register_pool"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(owner_cap), object::id_address(ref)],
        &mut stream,
    );

    let coin_metadata_address = bcs_stream::deserialize_address(&mut stream);
    let token_pool_package_id = bcs_stream::deserialize_address(&mut stream);
    let token_pool_module = bcs_stream::deserialize_string(&mut stream);
    let token_type_string = bcs_stream::deserialize_string(&mut stream);
    let initial_administrator = bcs_stream::deserialize_address(&mut stream);
    let token_pool_type_proof_string = bcs_stream::deserialize_string(&mut stream);
    let lock_or_burn_params = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_address(stream),
    );
    let release_or_mint_params = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_address(stream),
    );
    bcs_stream::assert_is_consumed(&stream);

    // Convert String to ascii::String
    let token_type = ascii::string(token_type_string.into_bytes());
    let token_pool_type_proof = ascii::string(token_pool_type_proof_string.into_bytes());

    register_pool_as_owner(
        owner_cap,
        ref,
        coin_metadata_address,
        token_pool_package_id,
        token_pool_module,
        token_type,
        initial_administrator,
        token_pool_type_proof,
        lock_or_burn_params,
        release_or_mint_params,
        ctx,
    );
}

public fun mcms_unregister_pool(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    _ctx: &mut TxContext,
) {
    let (_owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        state_object::McmsCallback,
        OwnerCap,
    >(
        registry,
        state_object::mcms_callback(),
        params,
    );
    assert!(function == string::utf8(b"unregister_pool"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref)],
        &mut stream,
    );

    let coin_metadata_address = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    unregister_pool_via_mcms(ref, coin_metadata_address);
}

public fun mcms_transfer_admin_role(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    _ctx: &mut TxContext,
) {
    let (_owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        state_object::McmsCallback,
        OwnerCap,
    >(
        registry,
        state_object::mcms_callback(),
        params,
    );
    assert!(function == string::utf8(b"transfer_admin_role"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref)],
        &mut stream,
    );

    let coin_metadata_address = bcs_stream::deserialize_address(&mut stream);
    let new_admin = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    transfer_admin_role_internal(
        ref,
        coin_metadata_address,
        new_admin,
        mcms_registry::get_multisig_address(),
    );
}

public fun mcms_accept_admin_role(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    _ctx: &mut TxContext,
) {
    let (_owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        state_object::McmsCallback,
        OwnerCap,
    >(
        registry,
        state_object::mcms_callback(),
        params,
    );
    assert!(function == string::utf8(b"accept_admin_role"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref)],
        &mut stream,
    );

    let coin_metadata_address = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    accept_admin_role_internal(ref, coin_metadata_address, mcms_registry::get_multisig_address());
}

#[test_only]
public fun insert_token_configs_for_test<TypeProof: drop>(
    ref: &mut CCIPObjectRef,
    administrator: address,
    coin_metadata_addresses: vector<address>,
    _proof: TypeProof,
) {
    let state = state_object::borrow_mut<TokenAdminRegistryState>(ref);
    let mut i = 0;
    while (i < coin_metadata_addresses.length()) {
        let token_config = TokenConfig {
            token_pool_package_id: @0x0,
            token_pool_module: string::utf8(b"TestModule"),
            token_type: ascii::string(b"TestType"),
            administrator,
            pending_administrator: @0x0,
            token_pool_type_proof: ascii::string(b"TestProof"),
            lock_or_burn_params: vector[],
            release_or_mint_params: vector[],
        };
        state
            .token_configs
            .push_back(
                coin_metadata_addresses[i],
                token_config,
            );
        i = i + 1;
    }
}

#[test_only]
public fun transfer_admin_role_internal_for_test(
    ref: &mut CCIPObjectRef,
    coin_metadata_address: address,
    new_admin: address,
    caller: address,
) {
    transfer_admin_role_internal(ref, coin_metadata_address, new_admin, caller);
}

#[test_only]
public fun accept_admin_role_internal_for_test(
    ref: &mut CCIPObjectRef,
    coin_metadata_address: address,
    caller: address,
) {
    accept_admin_role_internal(ref, coin_metadata_address, caller);
}

#[test_only]
public fun has_pending_admin_transfer(ref: &CCIPObjectRef, coin_metadata_address: address): bool {
    let state = state_object::borrow<TokenAdminRegistryState>(ref);
    assert!(state.token_configs.contains(coin_metadata_address), ETokenNotRegistered);

    let token_config = state.token_configs.borrow(coin_metadata_address);
    token_config.pending_administrator != @0x0
}

#[test_only]
public fun get_pending_admin_transfer(
    ref: &CCIPObjectRef,
    coin_metadata_address: address,
): (address, address) {
    let state = state_object::borrow<TokenAdminRegistryState>(ref);
    assert!(state.token_configs.contains(coin_metadata_address), ETokenNotRegistered);

    let token_config = state.token_configs.borrow(coin_metadata_address);
    (token_config.administrator, token_config.pending_administrator)
}

#[test_only]
public fun test_mcms_register_entrypoint(
    owner_cap: OwnerCap,
    registry: &mut Registry,
    ctx: &mut TxContext,
) {
    let publisher = ownable::borrow_publisher(&owner_cap);
    let publisher_wrapper = mcms_registry::create_publisher_wrapper(
        publisher,
        state_object::mcms_callback(),
    );
    mcms_registry::register_entrypoint(
        registry,
        publisher_wrapper,
        state_object::mcms_callback(),
        owner_cap,
        vector[b"fee_quoter", b"rmn_remote", b"state_object", b"token_admin_registry"], // Allowed CCIP modules
        ctx,
    );
}
