/// this module provides the functionality to:
/// 1. store the treasury cap object within the token state
/// 2. store the deny cap object if presented
/// 3. provide the functionality to issue MintCap objects and configure its allowance
/// 4. provide the functionality to mint and burn the token
/// 5. provide the functionality to blocklist and unblocklist addresses
/// 6. provide the functionality to pause and unpause the token
/// 7. provide the functionality to destroy the token
/// 8. provide the functionality to get the owner of the token
/// 9. provide the functionality to get the total supply of the token
module managed_token::managed_token;

use managed_token::mint_allowance::{Self, MintAllowance};
use managed_token::ownable::{Self, OwnerCap, OwnableState};
use mcms::bcs_stream;
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, Registry, ExecutingCallbackParams};
use std::string::{Self, String};
use sui::coin::{
    Self,
    Coin,
    DenyCapV2,
    TreasuryCap,
    deny_list_v2_is_global_pause_enabled_next_epoch as is_paused,
    deny_list_v2_contains_next_epoch as is_blocklisted
};
use sui::deny_list::DenyList;
use sui::event;
use sui::package::UpgradeCap;
use sui::vec_map::{Self, VecMap};

public struct TokenState<phantom T> has key, store {
    id: UID,
    treasury_cap: TreasuryCap<T>,
    deny_cap: Option<DenyCapV2<T>>,
    /// A map of { authorized MintCap ID => its MintAllowance }.
    mint_allowances_map: VecMap<ID, MintAllowance<T>>,
    ownable_state: OwnableState<T>,
}

/// An object representing the ability to mint unlimitedly or up to an allowance.
/// The privilege can be revoked by the token owner.
public struct MintCap<phantom T> has key, store {
    id: UID,
}

// === Events ===

public struct MintCapCreated<phantom T> has copy, drop {
    mint_cap: ID,
}

public struct MinterConfigured<phantom T> has copy, drop {
    mint_cap_owner: address,
    mint_cap: ID,
    allowance: u64,
    is_unlimited: bool,
}

public struct Minted has copy, drop {
    mint_cap: ID,
    minter: address,
    to: address,
    amount: u64,
}

public struct Burnt has copy, drop {
    mint_cap: ID,
    burner: address,
    from: address,
    amount: u64,
}

public struct Blocklisted<phantom T> has copy, drop {
    address: address,
}

public struct Unblocklisted<phantom T> has copy, drop {
    address: address,
}

public struct Paused<phantom T> has copy, drop {}

public struct Unpaused<phantom T> has copy, drop {}

public struct MinterAllowanceIncremented<phantom T> has copy, drop {
    mint_cap: ID,
    allowance_increment: u64,
    new_allowance: u64,
}

public struct MinterUnlimitedAllowanceSet<phantom T> has copy, drop {
    mint_cap: ID,
}

const EDeniedAddress: u64 = 1;
const EDenyCapNotFound: u64 = 2;
const EInsufficientAllowance: u64 = 3;
const EInvalidOwnerCap: u64 = 4;
const EPaused: u64 = 5;
const EUnauthorizedMintCap: u64 = 6;
const EZeroAmount: u64 = 7;
const ECannotIncreaseUnlimitedAllowance: u64 = 8;
const EInvalidFunction: u64 = 9;
const EInvalidStateAddress: u64 = 10;
const EInvalidRegistryAddress: u64 = 11;
const EInvalidDenyListAddress: u64 = 12;

public fun type_and_version(): String {
    string::utf8(b"ManagedToken 1.0.0")
}

public fun initialize<T>(treasury_cap: TreasuryCap<T>, ctx: &mut TxContext) {
    initialize_internal(treasury_cap, option::none(), ctx);
}

public fun initialize_with_deny_cap<T>(
    treasury_cap: TreasuryCap<T>,
    deny_cap: DenyCapV2<T>,
    ctx: &mut TxContext,
) {
    initialize_internal(treasury_cap, option::some(deny_cap), ctx);
}

#[allow(lint(self_transfer))]
fun initialize_internal<T>(
    treasury_cap: TreasuryCap<T>,
    deny_cap: Option<DenyCapV2<T>>,
    ctx: &mut TxContext,
) {
    let (ownable_state, owner_cap) = ownable::new(ctx);

    let state = TokenState<T> {
        id: object::new(ctx),
        treasury_cap,
        deny_cap,
        mint_allowances_map: vec_map::empty(),
        ownable_state,
    };

    transfer::share_object(state);
    transfer::public_transfer(owner_cap, ctx.sender());
}

public fun mint_allowance<T>(state: &TokenState<T>, mint_cap: ID): (u64, bool) {
    if (!state.is_authorized_mint_cap(mint_cap)) return (0, false);
    state.mint_allowances_map.get(&mint_cap).allowance_info()
}

/// Returns the total amount of Coin<T> in circulation.
public fun total_supply<T>(state: &TokenState<T>): u64 {
    state.treasury_cap.total_supply()
}

/// Checks if a MintCap object is authorized to mint.
public fun is_authorized_mint_cap<T>(state: &TokenState<T>, id: ID): bool {
    state.mint_allowances_map.contains(&id)
}

/// Convenience function that
/// 1. creates a MintCap and its allowance object
/// 2. transfers the MintCap object to a minter
///
/// - Only callable by the token owner.
public fun configure_new_minter<T>(
    state: &mut TokenState<T>,
    owner_cap: &OwnerCap<T>,
    minter: address,
    allowance: u64,
    is_unlimited: bool,
    ctx: &mut TxContext,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);

    let mint_cap = MintCap<T> { id: object::new(ctx) };
    event::emit(MintCapCreated<T> {
        mint_cap: object::id(&mint_cap),
    });

    let mut mint_allowance = mint_allowance::new<T>();
    mint_allowance.set(allowance, is_unlimited);
    state.mint_allowances_map.insert(object::id(&mint_cap), mint_allowance);

    event::emit(MinterConfigured<T> {
        mint_cap_owner: minter,
        mint_cap: object::id(&mint_cap),
        allowance,
        is_unlimited,
    });

    transfer::transfer(mint_cap, minter);
}

/// Increment allowance for a MintCap
/// - Only callable by the token owner.
/// - Only callable when not paused.
public fun increment_mint_allowance<T>(
    state: &mut TokenState<T>,
    owner_cap: &OwnerCap<T>,
    mint_cap_id: ID,
    deny_list: &DenyList,
    allowance_increment: u64,
    _ctx: &TxContext,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    assert!(!is_paused<T>(deny_list), EPaused);
    assert!(allowance_increment > 0, EZeroAmount);
    assert!(state.is_authorized_mint_cap(mint_cap_id), EUnauthorizedMintCap);

    assert!(
        !state.mint_allowances_map.get(&mint_cap_id).is_unlimited(),
        ECannotIncreaseUnlimitedAllowance,
    );
    state.mint_allowances_map.get_mut(&mint_cap_id).increase(allowance_increment);

    let new_allowance = state.mint_allowances_map.get(&mint_cap_id).value();

    event::emit(MinterAllowanceIncremented<T> {
        mint_cap: mint_cap_id,
        allowance_increment,
        new_allowance,
    });
}

/// Set the unlimited bool for a MintCap's allowance
/// - Only callable by the token owner.
/// - Only callable when not paused.
public fun set_unlimited_mint_allowances<T>(
    state: &mut TokenState<T>,
    owner_cap: &OwnerCap<T>,
    mint_cap_id: ID,
    deny_list: &DenyList,
    is_unlimited: bool,
    _ctx: &TxContext,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    assert!(!is_paused<T>(deny_list), EPaused);
    assert!(state.is_authorized_mint_cap(mint_cap_id), EUnauthorizedMintCap);

    state.mint_allowances_map.get_mut(&mint_cap_id).set(0, is_unlimited);

    event::emit(MinterUnlimitedAllowanceSet<T> {
        mint_cap: mint_cap_id,
    });
}

public fun get_all_mint_caps<T>(state: &TokenState<T>): vector<ID> {
    state.mint_allowances_map.keys()
}

/// Mints a Coin object with a specified amount (limited to the MintCap's allowance)
/// to a recipient address, increasing the total supply.
/// - Only callable by a minter.
/// - Only callable when not paused.
/// - Only callable if minter is not blocklisted.
/// - Only callable if recipient is not blocklisted.
public fun mint_and_transfer<T>(
    state: &mut TokenState<T>,
    mint_cap: &MintCap<T>,
    deny_list: &DenyList,
    amount: u64,
    recipient: address,
    ctx: &mut TxContext,
) {
    validate_mint(state, deny_list, mint_cap, amount, recipient, ctx);

    state.treasury_cap.mint_and_transfer(amount, recipient, ctx);

    event::emit(Minted {
        mint_cap: object::id(mint_cap),
        minter: ctx.sender(),
        to: recipient,
        amount,
    });
}

public fun mint<T>(
    state: &mut TokenState<T>,
    mint_cap: &MintCap<T>,
    deny_list: &DenyList,
    amount: u64,
    recipient: address,
    ctx: &mut TxContext,
): Coin<T> {
    validate_mint(state, deny_list, mint_cap, amount, recipient, ctx);

    let coin: Coin<T> = state.treasury_cap.mint(amount, ctx);

    event::emit(Minted {
        mint_cap: object::id(mint_cap),
        minter: ctx.sender(),
        to: recipient,
        amount,
    });

    coin
}

fun validate_mint<T>(
    state: &mut TokenState<T>,
    deny_list: &DenyList,
    mint_cap: &MintCap<T>,
    amount: u64,
    recipient: address,
    ctx: &TxContext,
) {
    assert!(!is_paused<T>(deny_list), EPaused);
    assert!(!is_blocklisted<T>(deny_list, ctx.sender()), EDeniedAddress);
    assert!(!is_blocklisted<T>(deny_list, recipient), EDeniedAddress);
    let mint_cap_id = object::id(mint_cap);
    assert!(state.is_authorized_mint_cap(mint_cap_id), EUnauthorizedMintCap);
    assert!(amount > 0, EZeroAmount);

    let mint_allowance = state.mint_allowances_map.get_mut(&mint_cap_id);
    assert!(
        mint_allowance.is_unlimited() || mint_allowance.value() >= amount,
        EInsufficientAllowance,
    );
    if (!mint_allowance.is_unlimited()) {
        mint_allowance.decrease(amount);
    };
}

/// Burns a Coin object, decreasing the total supply.
/// - Only callable by a minter.
/// - Only callable when not paused.
/// - Only callable if minter is not blocklisted.
public fun burn<T>(
    state: &mut TokenState<T>,
    mint_cap: &MintCap<T>,
    deny_list: &DenyList,
    coin: Coin<T>,
    from: address,
    ctx: &TxContext,
) {
    assert!(!is_paused<T>(deny_list), EPaused);
    assert!(!is_blocklisted<T>(deny_list, ctx.sender()), EDeniedAddress);
    let mint_cap_id = object::id(mint_cap);
    assert!(state.is_authorized_mint_cap(mint_cap_id), EUnauthorizedMintCap);

    let amount = coin.value();
    assert!(amount > 0, EZeroAmount);

    state.treasury_cap.burn(coin);
    event::emit(Burnt {
        mint_cap: mint_cap_id,
        burner: ctx.sender(),
        from,
        amount,
    });
}

/// Blocklists an address.
/// - Only callable by the token owner.
/// - the address must not be blocklisted already
public fun blocklist<T>(
    state: &mut TokenState<T>,
    owner_cap: &OwnerCap<T>,
    deny_list: &mut DenyList,
    addr: address,
    ctx: &mut TxContext,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);

    if (!is_blocklisted<T>(deny_list, addr)) {
        coin::deny_list_v2_add<T>(deny_list, borrow_deny_cap_mut(state), addr, ctx);
        event::emit(Blocklisted<T> {
            address: addr,
        })
    };
}

/// Unblocklists an address.
/// - Only callable by the token owner.
/// - the address must be blocklisted already
public fun unblocklist<T>(
    state: &mut TokenState<T>,
    owner_cap: &OwnerCap<T>,
    deny_list: &mut DenyList,
    addr: address,
    ctx: &mut TxContext,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);

    if (is_blocklisted<T>(deny_list, addr)) {
        coin::deny_list_v2_remove<T>(deny_list, borrow_deny_cap_mut(state), addr, ctx);
        event::emit(Unblocklisted<T> {
            address: addr,
        })
    };
}

/// Triggers stopped state; pause all transfers.
/// - Only callable by the token owner.
public fun pause<T>(
    state: &mut TokenState<T>,
    owner_cap: &OwnerCap<T>,
    deny_list: &mut DenyList,
    ctx: &mut TxContext,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);

    assert!(state.deny_cap.is_some(), EDenyCapNotFound);

    if (!is_paused<T>(deny_list)) {
        coin::deny_list_v2_enable_global_pause(deny_list, borrow_deny_cap_mut(state), ctx);
        event::emit(Paused<T> {});
    };
}

/// Restores normal state; unpause all transfers.
/// - Only callable by the token owner.
public fun unpause<T>(
    state: &mut TokenState<T>,
    owner_cap: &OwnerCap<T>,
    deny_list: &mut DenyList,
    ctx: &mut TxContext,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);

    if (is_paused<T>(deny_list)) {
        coin::deny_list_v2_disable_global_pause(deny_list, borrow_deny_cap_mut(state), ctx);
        event::emit(Unpaused<T> {});
    };
}

public fun destroy_managed_token<T>(
    owner_cap: OwnerCap<T>,
    state: TokenState<T>,
    ctx: &mut TxContext,
): (TreasuryCap<T>, Option<DenyCapV2<T>>) {
    assert!(
        object::id(&owner_cap) == ownable::owner_cap_id(&state.ownable_state),
        EInvalidOwnerCap,
    );

    let TokenState<T> {
        id: state_id,
        treasury_cap,
        deny_cap,
        mut mint_allowances_map,
        ownable_state,
    } = state;

    object::delete(state_id);
    let keys = mint_allowances_map.keys();
    let mut i = 0;
    while (i < keys.length()) {
        let (_id, mint_allowance) = mint_allowances_map.remove(&keys[i]);
        mint_allowance.destroy();
        i = i + 1;
    };
    mint_allowances_map.destroy_empty();

    ownable::destroy_ownable_state(ownable_state, ctx);
    ownable::destroy_owner_cap(owner_cap, ctx);

    (treasury_cap, deny_cap)
}

/// Access function to get a reference to the treasury cap
/// This allows other modules to access the treasury cap for registration purposes
public fun borrow_treasury_cap<T>(state: &TokenState<T>, owner_cap: &OwnerCap<T>): &TreasuryCap<T> {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    &state.treasury_cap
}

/// Returns a mutable reference of the DenyCap.
fun borrow_deny_cap_mut<T>(state: &mut TokenState<T>): &mut DenyCapV2<T> {
    assert!(state.deny_cap.is_some(), EDenyCapNotFound);
    state.deny_cap.borrow_mut()
}

#[test_only]
public fun get_ownable_state<T>(state: &mut TokenState<T>): &mut OwnableState<T> {
    &mut state.ownable_state
}

// ================================================================
// |                      Ownable Functions                       |
// ================================================================

public fun owner<T>(state: &TokenState<T>): address {
    ownable::owner(&state.ownable_state)
}

public fun has_pending_transfer<T>(state: &TokenState<T>): bool {
    ownable::has_pending_transfer(&state.ownable_state)
}

public fun pending_transfer_from<T>(state: &TokenState<T>): Option<address> {
    ownable::pending_transfer_from(&state.ownable_state)
}

public fun pending_transfer_to<T>(state: &TokenState<T>): Option<address> {
    ownable::pending_transfer_to(&state.ownable_state)
}

public fun pending_transfer_accepted<T>(state: &TokenState<T>): Option<bool> {
    ownable::pending_transfer_accepted(&state.ownable_state)
}

public fun transfer_ownership<T>(
    state: &mut TokenState<T>,
    owner_cap: &OwnerCap<T>,
    new_owner: address,
    ctx: &mut TxContext,
) {
    ownable::transfer_ownership(owner_cap, &mut state.ownable_state, new_owner, ctx);
}

public fun accept_ownership<T>(state: &mut TokenState<T>, ctx: &mut TxContext) {
    ownable::accept_ownership(&mut state.ownable_state, ctx);
}

public fun accept_ownership_from_object<T>(
    state: &mut TokenState<T>,
    from: &mut UID,
    ctx: &mut TxContext,
) {
    ownable::accept_ownership_from_object(&mut state.ownable_state, from, ctx);
}

public fun mcms_accept_ownership<T>(
    state: &mut TokenState<T>,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_, _, function_name, data) = mcms_registry::get_callback_params_for_mcms(
        params,
        McmsCallback {},
    );
    assert!(function_name == string::utf8(b"mcms_accept_ownership"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    let state_address = bcs_stream::deserialize_address(&mut stream);
    assert!(state_address == object::id_address(state), EInvalidStateAddress);

    let mcms = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    ownable::mcms_accept_ownership(&mut state.ownable_state, mcms, ctx);
}

public fun execute_ownership_transfer<T>(
    owner_cap: OwnerCap<T>,
    state: &mut TokenState<T>,
    to: address,
    ctx: &mut TxContext,
) {
    ownable::execute_ownership_transfer(owner_cap, &mut state.ownable_state, to, ctx);
}

public fun execute_ownership_transfer_to_mcms<T>(
    owner_cap: OwnerCap<T>,
    state: &mut TokenState<T>,
    registry: &mut Registry,
    to: address,
    ctx: &mut TxContext,
) {
    ownable::execute_ownership_transfer_to_mcms(
        owner_cap,
        &mut state.ownable_state,
        registry,
        to,
        McmsCallback {},
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

public struct McmsCallback has drop {}

fun validate_shared_objects<T>(
    state: &TokenState<T>,
    registry: &Registry,
    deny_list: &DenyList,
    stream: &mut bcs_stream::BCSStream,
) {
    let state_address = bcs_stream::deserialize_address(stream);
    assert!(state_address == object::id_address(state), EInvalidStateAddress);
    let registry_address = bcs_stream::deserialize_address(stream);
    assert!(registry_address == object::id_address(registry), EInvalidRegistryAddress);
    let deny_list_address = bcs_stream::deserialize_address(stream);
    assert!(deny_list_address == object::id_address(deny_list), EInvalidDenyListAddress);
}

public fun mcms_configure_new_minter<T>(
    state: &mut TokenState<T>,
    registry: &mut Registry,
    deny_list: &mut DenyList,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback, OwnerCap<T>>(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"configure_new_minter"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    validate_shared_objects(state, registry, deny_list, &mut stream);

    let minter = bcs_stream::deserialize_address(&mut stream);
    let allowance = bcs_stream::deserialize_u64(&mut stream);
    let is_unlimited = bcs_stream::deserialize_bool(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    configure_new_minter(state, owner_cap, minter, allowance, is_unlimited, ctx);
}

public fun mcms_increment_mint_allowance<T>(
    state: &mut TokenState<T>,
    registry: &mut Registry,
    deny_list: &mut DenyList,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback, OwnerCap<T>>(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"increment_mint_allowance"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    validate_shared_objects(state, registry, deny_list, &mut stream);

    let mint_cap_address = bcs_stream::deserialize_address(&mut stream);
    let allowance_increment = bcs_stream::deserialize_u64(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    increment_mint_allowance(
        state,
        owner_cap,
        mint_cap_address.to_id(),
        deny_list,
        allowance_increment,
        ctx,
    );
}

public fun mcms_set_unlimited_mint_allowances<T>(
    state: &mut TokenState<T>,
    registry: &mut Registry,
    deny_list: &mut DenyList,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback, OwnerCap<T>>(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"set_unlimited_mint_allowances"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    validate_shared_objects(state, registry, deny_list, &mut stream);

    let mint_cap_address = bcs_stream::deserialize_address(&mut stream);
    let is_unlimited = bcs_stream::deserialize_bool(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    set_unlimited_mint_allowances(
        state,
        owner_cap,
        mint_cap_address.to_id(),
        deny_list,
        is_unlimited,
        ctx,
    );
}

public fun mcms_blocklist<T>(
    state: &mut TokenState<T>,
    registry: &mut Registry,
    deny_list: &mut DenyList,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback, OwnerCap<T>>(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"blocklist"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    validate_shared_objects(state, registry, deny_list, &mut stream);

    let addr = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    blocklist(state, owner_cap, deny_list, addr, ctx);
}

public fun mcms_unblocklist<T>(
    state: &mut TokenState<T>,
    registry: &mut Registry,
    deny_list: &mut DenyList,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback, OwnerCap<T>>(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"unblocklist"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    validate_shared_objects(state, registry, deny_list, &mut stream);

    let addr = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    unblocklist(state, owner_cap, deny_list, addr, ctx);
}

public fun mcms_pause<T>(
    state: &mut TokenState<T>,
    registry: &mut Registry,
    deny_list: &mut DenyList,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback, OwnerCap<T>>(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"pause"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    validate_shared_objects(state, registry, deny_list, &mut stream);
    bcs_stream::assert_is_consumed(&stream);

    pause(state, owner_cap, deny_list, ctx);
}

public fun mcms_unpause<T>(
    state: &mut TokenState<T>,
    registry: &mut Registry,
    deny_list: &mut DenyList,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback, OwnerCap<T>>(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"unpause"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    validate_shared_objects(state, registry, deny_list, &mut stream);
    bcs_stream::assert_is_consumed(&stream);

    unpause(state, owner_cap, deny_list, ctx);
}
