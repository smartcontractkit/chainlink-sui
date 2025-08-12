module ccip::token_admin_registry;

use std::ascii;
use std::string::{Self, String};
use std::type_name;

use sui::coin::{CoinMetadata, TreasuryCap};
use sui::event;
use sui::linked_table::{Self, LinkedTable};

use ccip::state_object::{Self, CCIPObjectRef};
use ccip::ownable::OwnerCap;

// TODO: consider add/using a different structure if someone registers too many tokens
// figure out & ask about the vector & map size limit for different structures
public struct TokenAdminRegistryState has key, store {
    id: UID,
    // coin metadata object id -> token config
    token_configs: LinkedTable<address, TokenConfig>,
}

public struct TokenConfig has store, drop, copy {
    token_pool_package_id: address,
    token_pool_module: String,
    // the type of the token
    token_type: ascii::String,
    administrator: address,
    pending_administrator: address,
    // type proof of the token pool
    type_proof: ascii::String,
    lock_or_burn_params: vector<address>,
    release_or_mint_params: vector<address>,
}

public struct PoolSet has copy, drop {
    coin_metadata_address: address,
    previous_pool_package_id: address,
    new_pool_package_id: address,
    // type proof of the new token pool
    type_proof: ascii::String,
    lock_or_burn_params: vector<address>,
    release_or_mint_params: vector<address>,
}

public struct PoolRegistered has copy, drop {
    coin_metadata_address: address,
    token_pool_package_id: address,
    administrator: address,
    // type proof of the token pool
    type_proof: ascii::String,
}

public struct PoolUnregistered has copy, drop {
    coin_metadata_address: address,
    previous_pool_address: address,
}

public struct AdministratorTransferRequested has copy, drop {
    coin_metadata_address: address,
    current_admin: address,
    new_admin: address
}

public struct AdministratorTransferred has copy, drop {
    coin_metadata_address: address,
    new_admin: address
}

const ENotPendingAdministrator: u64 = 1;
const EAlreadyInitialized: u64 = 2;
const ETokenAlreadyRegistered: u64 = 3;
const ETokenNotRegistered: u64 = 4;
const ENotAdministrator: u64 = 5;
const ETokenAddressNotRegistered: u64 = 6;
const ENotAllowed: u64 = 7;

public fun type_and_version(): String {
    string::utf8(b"TokenAdminRegistry 1.6.0")
}

public fun initialize(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    ctx: &mut TxContext
) {
    assert!(
        !state_object::contains<TokenAdminRegistryState>(ref),
        EAlreadyInitialized
    );
    let state = TokenAdminRegistryState {
        id: object::new(ctx),
        token_configs: linked_table::new(ctx),
    };

    state_object::add(ref, owner_cap, state, ctx);
}

public fun get_pools(
    ref: &CCIPObjectRef,
    coin_metadata_addresses: vector<address>
): vector<address>{
    let state = state_object::borrow<TokenAdminRegistryState>(ref);

    let mut token_pool_addresses: vector<address> = vector[];
    coin_metadata_addresses.do_ref!(
        |metadata_address| {
            let metadata_address: address = *metadata_address;
            if (state.token_configs.contains(metadata_address)) {
                let token_config = state.token_configs.borrow(metadata_address);
                token_pool_addresses.push_back(token_config.token_pool_package_id);
            } else {
                // returns @0x0 for assets without token pools.
                token_pool_addresses.push_back(@0x0);
            }
        }
    );

    token_pool_addresses
}

// this function can also take a coin metadata or a coin::zero
// but that requires adding a type parameter to the function
public fun get_pool(ref: &CCIPObjectRef, coin_metadata_address: address): address {
    let state = state_object::borrow<TokenAdminRegistryState>(ref);

    if (state.token_configs.contains(coin_metadata_address)) {
        let token_config = state.token_configs.borrow(coin_metadata_address);
        token_config.token_pool_package_id
    } else {
        // returns @0x0 for assets without token pools.
        @0x0
    }
}

public fun get_token_config(
    ref: &CCIPObjectRef, coin_metadata_address: address
): TokenConfig {
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
            type_proof: ascii::string(b""),
            lock_or_burn_params: vector[],
            release_or_mint_params: vector[],
        }
    }
}

public fun get_token_configs(
    ref: &CCIPObjectRef, coin_metadata_addresses: vector<address>
): vector<TokenConfig> {
    let mut token_configs: vector<TokenConfig> = vector[];

    coin_metadata_addresses.do_ref!(
        |coin_metadata_address| {
            let coin_metadata_address: address = *coin_metadata_address;
            let token_config = get_token_config(ref, coin_metadata_address);
            token_configs.push_back(token_config);
        }
    );

    token_configs
}

public fun get_token_config_data(token_config: TokenConfig): (address, String, ascii::String, address, address, ascii::String, vector<address>, vector<address>) {
    (
        token_config.token_pool_package_id,
        token_config.token_pool_module,
        token_config.token_type,
        token_config.administrator,
        token_config.pending_administrator,
        token_config.type_proof,
        token_config.lock_or_burn_params,
        token_config.release_or_mint_params,
    )
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
///   - vector<address>: List of token addresses (up to max_count)
///   - address: Next key to use for pagination (pass this as start_key in next call)
///   - bool: Whether there are more tokens after this batch
public fun get_all_configured_tokens(
    ref: &CCIPObjectRef, start_key: address, max_count: u64
): (vector<address>, address, bool) {
    let state = state_object::borrow<TokenAdminRegistryState>(ref);

    let mut i = 0;
    let mut result = vector[];
    let mut key = start_key;
    if (key == @0x0) {
        if (state.token_configs.is_empty()) {
            return (result, key, false)
        };
        if (max_count == 0) {
            return (result, key, true)
        };
        key = *state.token_configs.front().borrow();
        result.push_back(key);
        i = 1;
    } else {
        assert!(state.token_configs.contains(start_key), ETokenAddressNotRegistered);
    };

    while (i < max_count) {
        let next_key_opt = state.token_configs.next(key);
        if (next_key_opt.is_none()) {
            return (result, key, false)
        };

        key = *next_key_opt.borrow();
        result.push_back(key);
        i = i + 1;
    };

    // Check if there are more tokens after the last key
    let has_more = state.token_configs.next(key).is_some();
    (result, key, has_more)
}

// ================================================================
// |                       Register Pool                          |
// ================================================================

// only the token owner with the treasury cap can call this function.
public fun register_pool<T, TypeProof: drop>(
    ref: &mut CCIPObjectRef,
    _: &TreasuryCap<T>, // passing in the treasury cap to demonstrate ownership over the token
    coin_metadata: &CoinMetadata<T>,
    token_pool_package_id: address,
    token_pool_module: String,
    initial_administrator: address,
    lock_or_burn_params: vector<address>,
    release_or_mint_params: vector<address>,
    _proof: TypeProof,
) {
    let coin_metadata_address: address = object::id_to_address(&object::id(coin_metadata));
    let token_type = type_name::get<T>().into_string();
    let proof_tn = type_name::get<TypeProof>();
    register_pool_internal(
        ref,
        coin_metadata_address,
        token_pool_package_id,
        token_pool_module,
        token_type,
        initial_administrator,
        type_name::into_string(proof_tn),
        lock_or_burn_params,
        release_or_mint_params,
    );
}

public fun register_pool_by_admin(
    ref: &mut CCIPObjectRef,
    _: state_object::CCIPAdminProof,
    coin_metadata_address: address,
    token_pool_package_id: address,
    token_pool_module: String,
    token_type: ascii::String,
    initial_administrator: address,
    proof: ascii::String,
    lock_or_burn_params: vector<address>,
    release_or_mint_params: vector<address>,
    _: &mut TxContext,
) {
    register_pool_internal(
        ref,
        coin_metadata_address,
        token_pool_package_id,
        token_pool_module,
        token_type,
        initial_administrator,
        proof,
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
    proof: ascii::String,
    lock_or_burn_params: vector<address>,
    release_or_mint_params: vector<address>,
) {
    let state = state_object::borrow_mut<TokenAdminRegistryState>(ref);
    assert!(
        !state.token_configs.contains(coin_metadata_address),
        ETokenAlreadyRegistered
    );

    let token_config = TokenConfig {
        token_pool_package_id,
        token_pool_module,
        token_type,
        administrator: initial_administrator,
        pending_administrator: @0x0,
        type_proof: proof,
        lock_or_burn_params,
        release_or_mint_params,
    };

    state.token_configs.push_back(coin_metadata_address, token_config);

    event::emit(
        PoolRegistered {
            coin_metadata_address,
            token_pool_package_id,
            administrator: initial_administrator,
            type_proof: proof,
        }
    );
}

public fun unregister_pool(
    ref: &mut CCIPObjectRef,
    coin_metadata_address: address,
    ctx: &mut TxContext,
) {
    let state = state_object::borrow_mut<TokenAdminRegistryState>(ref);

    assert!(
        state.token_configs.contains(coin_metadata_address),
        ETokenNotRegistered
    );

    let token_config = state.token_configs.remove(coin_metadata_address);
    
    assert!(token_config.administrator == ctx.sender(), ENotAllowed);

    let previous_pool_address = token_config.token_pool_package_id;

    event::emit(
        PoolUnregistered {
            coin_metadata_address,
            previous_pool_address,
        }
    );
}

public fun set_pool<TypeProof: drop>(
    ref: &mut CCIPObjectRef,
    coin_metadata_address: address,
    token_pool_package_id: address,
    token_pool_module: String,
    lock_or_burn_params: vector<address>,
    release_or_mint_params: vector<address>,
    _: TypeProof,
    ctx: &mut TxContext,
) {
    let state = state_object::borrow_mut<TokenAdminRegistryState>(ref);

    assert!(
        state.token_configs.contains(coin_metadata_address),
        ETokenNotRegistered
    );

    let token_config = state.token_configs.borrow_mut(coin_metadata_address);

    // the tx signer must be the administrator of the token pool.
    assert!(token_config.administrator == ctx.sender(), ENotAllowed);

    // TODO: sort out the UX here
    // the token pool changes, the package id, state address, module, and type proof will change.
    let previous_pool_package_id = token_config.token_pool_package_id;
    if (previous_pool_package_id != token_pool_package_id) {
        token_config.token_pool_package_id = token_pool_package_id;
        token_config.token_pool_module = token_pool_module;
        token_config.lock_or_burn_params = lock_or_burn_params;
        token_config.release_or_mint_params = release_or_mint_params;
        let proof_tn = type_name::get<TypeProof>();
        let proof_str = type_name::into_string(proof_tn);
        token_config.type_proof = proof_str;

        event::emit(
            PoolSet {
                coin_metadata_address,
                previous_pool_package_id,
                new_pool_package_id: token_pool_package_id,
                type_proof: proof_str,
                lock_or_burn_params,
                release_or_mint_params,
            }
        );
    }
}

public fun transfer_admin_role(
    ref: &mut CCIPObjectRef,
    coin_metadata_address: address,
    new_admin: address,
    ctx: &mut TxContext
) {
    let state = state_object::borrow_mut<TokenAdminRegistryState>(ref);

    assert!(
        state.token_configs.contains(coin_metadata_address),
        ETokenNotRegistered
    );

    let token_config = state.token_configs.borrow_mut(coin_metadata_address);

    assert!(
        token_config.administrator == ctx.sender(),
        ENotAdministrator
    );

    // can be @0x0 to cancel a pending transfer.
    token_config.pending_administrator = new_admin;

    event::emit(
        AdministratorTransferRequested {
            coin_metadata_address,
            current_admin: token_config.administrator,
            new_admin
        }
    );
}

public fun accept_admin_role(
    ref: &mut CCIPObjectRef,
    coin_metadata_address: address,
    ctx: &mut TxContext
) {
    let state = state_object::borrow_mut<TokenAdminRegistryState>(ref);

    assert!(
        state.token_configs.contains(coin_metadata_address),
        ETokenNotRegistered
    );

    let token_config = state.token_configs.borrow_mut(coin_metadata_address);

    assert!(
        token_config.pending_administrator == ctx.sender(),
        ENotPendingAdministrator
    );

    token_config.administrator = token_config.pending_administrator;
    token_config.pending_administrator = @0x0;

    event::emit(
        AdministratorTransferred { coin_metadata_address, new_admin: token_config.administrator }
    );
}

public fun is_administrator(
    ref: &CCIPObjectRef, coin_metadata_address: address, administrator: address
): bool {
    let state = state_object::borrow<TokenAdminRegistryState>(ref);

    assert!(
        state.token_configs.contains(coin_metadata_address),
        ETokenNotRegistered
    );

    let token_config = state.token_configs.borrow(coin_metadata_address);
    token_config.administrator == administrator
}

#[test_only]
public fun insert_token_configs_for_test<TypeProof: drop>(
    ref: &mut CCIPObjectRef,
    coin_metadata_addresses: vector<address>,
    _proof: TypeProof
) {
    let state = state_object::borrow_mut<TokenAdminRegistryState>(ref);
    let mut i = 0;
    while (i < coin_metadata_addresses.length()) {
        let token_config = TokenConfig {
            token_pool_package_id: @0x0,
            token_pool_module: string::utf8(b"TestModule"),
            token_type: ascii::string(b"TestType"),
            administrator: @0x0,
            pending_administrator: @0x0,
            type_proof: ascii::string(b"TestProof"),
            lock_or_burn_params: vector[],
            release_or_mint_params: vector[],
        };
        state.token_configs.push_back(
            coin_metadata_addresses[i],
            token_config,
        );
        i = i + 1;
    }
}
