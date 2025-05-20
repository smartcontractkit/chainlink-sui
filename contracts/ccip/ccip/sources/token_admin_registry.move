module ccip::token_admin_registry;

use std::string::{Self, String};
use std::type_name::{Self, TypeName};

use sui::coin::TreasuryCap;
use sui::event;
use sui::linked_table::{Self, LinkedTable};

use ccip::state_object::{Self, CCIPObjectRef, OwnerCap};

public struct TokenAdminRegistryState has key, store {
    id: UID,
    token_configs: LinkedTable<address, TokenConfig>,
}

public struct TokenConfig has store, drop, copy {
    token_pool_address: address,
    administrator: address,
    pending_administrator: address,
    type_proof: Option<TypeName>,
}

public struct PoolSet has copy, drop {
    coin_metadata_address: address,
    previous_pool_address: address,
    new_pool_address: address,
    type_proof: TypeName,
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

const E_NOT_PENDING_ADMINISTRATOR: u64 = 1;
const E_ALREADY_INITIALIZED: u64 = 2;
const E_FUNGIBLE_ASSET_ALREADY_REGISTERED: u64 = 3;
const E_FUNGIBLE_ASSET_NOT_REGISTERED: u64 = 4;
const E_NOT_ADMINISTRATOR: u64 = 5;
const E_TOKEN_ADDRESS_NOT_REGISTERED: u64 = 6;

public fun type_and_version(): String {
    string::utf8(b"TokenAdminRegistry 1.6.0")
}

public fun initialize(
    ref: &mut CCIPObjectRef,
    _: &OwnerCap,
    ctx: &mut TxContext
) {
    assert!(
        !state_object::contains<TokenAdminRegistryState>(ref),
        E_ALREADY_INITIALIZED
    );
    let state = TokenAdminRegistryState {
        id: object::new(ctx),
        token_configs: linked_table::new(ctx),
    };

    state_object::add(ref, state, ctx);
}

public fun get_pools(
    ref: &CCIPObjectRef,
    coin_metadata_addresses: vector<address>
): vector<address>{
    let state = state_object::borrow<TokenAdminRegistryState>(ref);

    let mut token_pool_addresses: vector<address> = vector::empty();
    coin_metadata_addresses.do_ref!(
        |metadata_address| {
            let metadata_address: address = *metadata_address;
            if (state.token_configs.contains(metadata_address)) {
                let token_config = state.token_configs.borrow(metadata_address);
                token_pool_addresses.push_back(token_config.token_pool_address);
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
        token_config.token_pool_address
    } else {
        // returns @0x0 for assets without token pools.
        @0x0
    }
}

public fun get_token_config(
    ref: &CCIPObjectRef, coin_metadata_address: address
): (address, address, address, Option<TypeName>) {
    let state = state_object::borrow<TokenAdminRegistryState>(ref);

    if (state.token_configs.contains(coin_metadata_address)) {
        let token_config = state.token_configs.borrow(coin_metadata_address);
        (
            token_config.token_pool_address,
            token_config.administrator,
            token_config.pending_administrator,
            token_config.type_proof,
        )
    } else {
        (@0x0, @0x0, @0x0, option::none())
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
        assert!(state.token_configs.contains(start_key), E_TOKEN_ADDRESS_NOT_REGISTERED);
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
#[allow(lint(self_transfer))]
public fun register_pool<T, TypeProof: drop>(
    ref: &mut CCIPObjectRef,
    _: &TreasuryCap<T>, // passing in the treasury cap to demonstrate ownership over the token
    coin_metadata_address: address,
    token_pool_address: address,
    _proof: TypeProof,
    ctx: &mut TxContext
) {
    register_pool_internal(
        ref,
        coin_metadata_address,
        token_pool_address,
        _proof,
        ctx,
    );
}

// TODO: validate this can work or we have to re-think the type proof for this function
// only the CCIP owner can call this function.
public fun register_pool_by_admin<TypeProof: drop>(
    ref: &mut CCIPObjectRef,
    coin_metadata_address: address,
    token_pool_address: address,
    _proof: TypeProof, // this needs to be a proof type from the token pool module, not from CCIP admin
    ctx: &mut TxContext,
) {
    assert!(
        ctx.sender() == state_object::get_current_owner(ref),
        E_NOT_ADMINISTRATOR,
    );

    register_pool_internal(
        ref,
        coin_metadata_address,
        token_pool_address,
        _proof,
        ctx
    );
}

fun register_pool_internal<TypeProof: drop>(
    ref: &mut CCIPObjectRef,
    coin_metadata_address: address, // LINK coin metadata
    token_pool_address: address, // a legit LINK source pool
    _proof: TypeProof, // use this proof type to validate the token pool address & token pool module name
    ctx: &TxContext,
) {
    let state = state_object::borrow_mut<TokenAdminRegistryState>(ref);
    assert!(
        !state.token_configs.contains(coin_metadata_address),
        E_FUNGIBLE_ASSET_ALREADY_REGISTERED
    );

    let proof_tn = type_name::get<TypeProof>();
    let token_config = TokenConfig {
        token_pool_address,
        administrator: ctx.sender(),
        pending_administrator: @0x0,
        type_proof: option::some(proof_tn),
    };

    state.token_configs.push_back(coin_metadata_address, token_config);
}

public fun set_pool<TypeProof: drop>(
    ref: &mut CCIPObjectRef,
    coin_metadata_address: address,
    token_pool_address: address,
    _: TypeProof,
    ctx: &mut TxContext
) {
    let state = state_object::borrow_mut<TokenAdminRegistryState>(ref);

    assert!(
        state.token_configs.contains(coin_metadata_address),
        E_FUNGIBLE_ASSET_NOT_REGISTERED
    );

    let token_config = state.token_configs.borrow_mut(coin_metadata_address);

    // the tx signer must be the administrator of the token pool.
    assert!(
        token_config.administrator == ctx.sender(),
        E_NOT_ADMINISTRATOR
    );

    let previous_pool_address = token_config.token_pool_address;
    if (previous_pool_address != token_pool_address) {
        token_config.token_pool_address = token_pool_address;
        token_config.type_proof = option::some(type_name::get<TypeProof>());

        event::emit(
            PoolSet {
                coin_metadata_address,
                previous_pool_address,
                new_pool_address: token_pool_address,
                type_proof: type_name::get<TypeProof>(),
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
        E_FUNGIBLE_ASSET_NOT_REGISTERED
    );

    let token_config = state.token_configs.borrow_mut(coin_metadata_address);

    assert!(
        token_config.administrator == ctx.sender(),
        E_NOT_ADMINISTRATOR
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
        E_FUNGIBLE_ASSET_NOT_REGISTERED
    );

    let token_config = state.token_configs.borrow_mut(coin_metadata_address);

    assert!(
        token_config.pending_administrator == ctx.sender(),
        E_NOT_PENDING_ADMINISTRATOR
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
        E_FUNGIBLE_ASSET_NOT_REGISTERED
    );

    let token_config = state.token_configs.borrow(coin_metadata_address);
    token_config.administrator == administrator
}

#[test_only]
public fun insert_token_addresses_for_test<TypeProof: drop>(
    ref: &mut CCIPObjectRef, token_addresses: vector<address>, _proof: TypeProof
) {
    let state = state_object::borrow_mut<TokenAdminRegistryState>(ref);

    token_addresses.do_ref!(
        |token_address| {
            state.token_configs.push_back(
                *token_address,
                TokenConfig {
                    token_pool_address: @0x0,
                    administrator: @0x0,
                    pending_administrator: @0x0,
                    type_proof: option::some(type_name::get<TypeProof>()),
                }
            );
        }
    );
}
