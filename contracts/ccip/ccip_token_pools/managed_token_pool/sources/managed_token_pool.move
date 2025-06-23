module managed_token_pool::managed_token_pool;

use std::string::{Self, String};
use std::type_name::{Self, TypeName};

use sui::clock::Clock;
use sui::coin::{Coin, CoinMetadata, TreasuryCap};
use sui::deny_list::{DenyList};

use ccip::dynamic_dispatcher as dd;
use ccip::eth_abi;
use ccip::offramp_state_helper as osh;
use ccip::state_object::CCIPObjectRef;
use ccip::token_admin_registry;

use ccip_token_pool::token_pool::{Self, TokenPoolState};

use managed_token::managed_token::{Self, TokenState, MintCap};

public struct OwnerCap has key, store {
    id: UID,
    state_id: ID,
}

public struct ManagedTokenPoolState<phantom T> has key {
    id: UID,
    token_pool_state: TokenPoolState,
    mint_cap: MintCap<T>,
}

const EInvalidCoinMetadata: u64 = 1;
const EInvalidArguments: u64 = 2;
const EInvalidOwnerCap: u64 = 3;

// ================================================================
// |                             Init                             |
// ================================================================

public fun type_and_version(): String {
    string::utf8(b"ManagedTokenPool 1.6.0")
}

#[allow(lint(self_transfer))]
public fun initialize<T: drop>(
    ref: &mut CCIPObjectRef,
    treasury_cap: &TreasuryCap<T>,
    coin_metadata: &CoinMetadata<T>,
    mint_cap: MintCap<T>,
    token_pool_package_id: address,
    token_pool_administrator: address,
    ctx: &mut TxContext,
) {
    let (_, managed_token_state_address, _, _) =
        initialize_internal(coin_metadata, mint_cap, ctx);

    token_admin_registry::register_pool(
        ref,
        treasury_cap,
        coin_metadata,
        token_pool_package_id,
        managed_token_state_address,
        string::utf8(b"managed_token_pool"),
        token_pool_administrator,
        TypeProof {},
    );
}

public fun initialize_by_ccip_admin<T: drop>(
    ref: &mut CCIPObjectRef,
    coin_metadata: &CoinMetadata<T>,
    mint_cap: MintCap<T>,
    token_pool_package_id: address,
    token_pool_administrator: address,
    ctx: &mut TxContext,
) {
    let (coin_metadata_address, managed_token_state_address, token_type_name, type_proof_type_name) =
        initialize_internal(coin_metadata, mint_cap, ctx);

    token_admin_registry::register_pool_by_admin(
        ref,
        coin_metadata_address,
        token_pool_package_id,
        managed_token_state_address,
        string::utf8(b"managed_token_pool"),
        token_type_name.into_string(),
        token_pool_administrator,
        type_proof_type_name.into_string(),
        ctx,
    );
}

#[allow(lint(self_transfer))]
fun initialize_internal<T: drop>(
    coin_metadata: &CoinMetadata<T>,
    mint_cap: MintCap<T>,
    ctx: &mut TxContext,
): (address, address, TypeName, TypeName) {
let coin_metadata_address: address = object::id_to_address(&object::id(coin_metadata));
    assert!(
        coin_metadata_address == @managed_token_coin_metadata,
        EInvalidCoinMetadata
    );

    let managed_token_pool = ManagedTokenPoolState<T> {
        id: object::new(ctx),
        token_pool_state: token_pool::initialize(coin_metadata_address, coin_metadata.get_decimals(), vector[], ctx),
        mint_cap,
    };
    let token_type_name = type_name::get<T>();
    let type_proof_type_name = type_name::get<TypeProof>();
    let managed_token_state_address = object::uid_to_address(&managed_token_pool.id);

    let owner_cap = OwnerCap {
        id: object::new(ctx),
        state_id: object::id(&managed_token_pool),
    };
    transfer::share_object(managed_token_pool);
    transfer::transfer(owner_cap, ctx.sender());

    (coin_metadata_address, managed_token_state_address, token_type_name, type_proof_type_name)
}

public fun add_remote_pool<T>(
    state: &mut ManagedTokenPoolState<T>,
    owner_cap: &OwnerCap,
    remote_chain_selector: u64,
    remote_pool_address: vector<u8>,
) {
    assert!(owner_cap.state_id == object::id(state), EInvalidOwnerCap);
    token_pool::add_remote_pool(
        &mut state.token_pool_state, remote_chain_selector, remote_pool_address
    );
}

public fun remove_remote_pool<T>(
    state: &mut ManagedTokenPoolState<T>,
    owner_cap: &OwnerCap,
    remote_chain_selector: u64,
    remote_pool_address: vector<u8>,
) {
    assert!(owner_cap.state_id == object::id(state), EInvalidOwnerCap);
    token_pool::remove_remote_pool(
        &mut state.token_pool_state, remote_chain_selector, remote_pool_address
    );
}

public fun is_supported_chain<T>(
    state: &ManagedTokenPoolState<T>,
    remote_chain_selector: u64
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
    remote_token_addresses_to_add: vector<vector<u8>>
) {
    assert!(owner_cap.state_id == object::id(state), EInvalidOwnerCap);
    token_pool::apply_chain_updates(
        &mut state.token_pool_state,
        remote_chain_selectors_to_remove,
        remote_chain_selectors_to_add,
        remote_pool_addresses_to_add,
        remote_token_addresses_to_add
    );
}

public fun get_allowlist_enabled<T>(state: &ManagedTokenPoolState<T>): bool {
    token_pool::get_allowlist_enabled(&state.token_pool_state)
}

public fun get_allowlist<T>(state: &ManagedTokenPoolState<T>): vector<address> {
    token_pool::get_allowlist(&state.token_pool_state)
}

public fun apply_allowlist_updates<T>(
    state: &mut ManagedTokenPoolState<T>,
    owner_cap: &OwnerCap,
    removes: vector<address>,
    adds: vector<address>
) {
    assert!(owner_cap.state_id == object::id(state), EInvalidOwnerCap);
    token_pool::apply_allowlist_updates(&mut state.token_pool_state, removes, adds);
}

// ================================================================
// |                 Exposing token_pool functions                |
// ================================================================

// this now returns the address of coin metadata
public fun get_token<T>(state: &ManagedTokenPoolState<T>): address {
    token_pool::get_token(&state.token_pool_state)
}

public fun get_router(): address {
    token_pool::get_router()
}

public fun get_token_decimals<T>(state: &ManagedTokenPoolState<T>): u8 {
    state.token_pool_state.get_local_decimals()
}

public fun get_remote_pools<T>(
    state: &ManagedTokenPoolState<T>,
    remote_chain_selector: u64
): vector<vector<u8>> {
    token_pool::get_remote_pools(&state.token_pool_state, remote_chain_selector)
}

public fun is_remote_pool<T>(
    state: &ManagedTokenPoolState<T>,
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
    state: &ManagedTokenPoolState<T>, remote_chain_selector: u64
): vector<u8> {
    token_pool::get_remote_token(&state.token_pool_state, remote_chain_selector)
}

// ================================================================
// |                         Burn/Mint                            |
// ================================================================

public struct TypeProof has drop {}

public fun lock_or_burn<T>(
    ref: &CCIPObjectRef,
    clock: &Clock,
    state: &mut ManagedTokenPoolState<T>,
    deny_list: &DenyList,
    token_state: &mut TokenState<T>,
    c: Coin<T>,
    token_params: dd::TokenParams,
    ctx: &mut TxContext
): dd::TokenParams {
    let amount = c.value();
    let sender = ctx.sender();
    let remote_chain_selector = dd::get_destination_chain_selector(&token_params);

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

    dd::add_source_token_transfer(
        ref,
        token_params,
        amount,
        state.token_pool_state.get_token(),
        dest_token_address,
        extra_data,
        TypeProof {},
    )
}

public fun release_or_mint<T>(
    ref: &CCIPObjectRef,
    clock: &Clock,
    pool: &mut ManagedTokenPoolState<T>,
    token_state: &mut TokenState<T>,
    deny_list: &DenyList,
    receiver_params: osh::ReceiverParams,
    index: u64,
    ctx: &mut TxContext
): osh::ReceiverParams {
    let remote_chain_selector = osh::get_source_chain_selector(&receiver_params);
    let (receiver, source_amount, dest_token_address, source_pool_address, source_pool_data, _) = osh::get_token_param_data(&receiver_params, index);
    let local_amount = token_pool::calculate_release_or_mint_amount(
        &pool.token_pool_state,
        source_pool_data,
        source_amount
    );

    token_pool::validate_release_or_mint(
        ref,
        clock,
        &mut pool.token_pool_state,
        remote_chain_selector,
        dest_token_address,
        source_pool_address,
        local_amount
    );

    let c: Coin<T> = managed_token::mint(
        token_state,
        &pool.mint_cap,
        deny_list,
        local_amount,
        receiver,
        ctx,
    );

    token_pool::emit_released_or_minted(
        &mut pool.token_pool_state,
        receiver,
        local_amount,
        remote_chain_selector,
    );

    osh::complete_token_transfer(
        ref,
        receiver_params,
        index,
        c,
        TypeProof {},
    )
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
    inbound_rates: vector<u64>
) {
    assert!(owner_cap.state_id == object::id(state), EInvalidOwnerCap);
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
    state: &mut ManagedTokenPoolState<T>,
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
    assert!(owner_cap.state_id == object::id(state), EInvalidOwnerCap);
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