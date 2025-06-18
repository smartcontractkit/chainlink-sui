module managed_token::managed_token;

use std::string::{Self, String};

use sui::coin::{
    Self, Coin, DenyCapV2, TreasuryCap,
    deny_list_v2_is_global_pause_enabled_next_epoch as is_paused,
    deny_list_v2_contains_next_epoch as is_blocklisted,
};
use sui::deny_list::{DenyList};
use sui::event;
use sui::vec_map::{Self, VecMap};
use sui::package::UpgradeCap;

use managed_token::mint_allowance::{Self, MintAllowance};
use managed_token::ownable::{Self, OwnerCap, OwnableState};

use mcms::mcms_registry::{Self, Registry, ExecutingCallbackParams};
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::bcs_stream;

public struct TokenState<phantom T> has key, store {
    id: UID,
    treasury_cap: TreasuryCap<T>,
    deny_cap: Option<DenyCapV2<T>>,
    /// A map of { authorized MintCap ID => its MintAllowance }.
    mint_allowances_map: VecMap<ID, MintAllowance<T>>,
    ownable_state: OwnableState<T>,
}

/// An object representing the ability to mint up to an allowance
/// specified in the Treasury.
/// The privilege can be revoked by the master minter.
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
    address: address
}

public struct Unblocklisted<phantom T> has copy, drop {
    address: address
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

public fun type_and_version(): String {
    string::utf8(b"ManagedToken 1.0.0")
}

public fun initialize<T: drop>(
    treasury_cap: TreasuryCap<T>,
    ctx: &mut TxContext,
) {
    initialize_internal(treasury_cap, option::none(), ctx);
}

public fun initialize_with_deny_cap<T: drop>(
    treasury_cap: TreasuryCap<T>,
    deny_cap: DenyCapV2<T>,
    ctx: &mut TxContext,
) {
    initialize_internal(treasury_cap, option::some(deny_cap), ctx);
}

#[allow(lint(self_transfer))]
fun initialize_internal<T: drop>(
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

/// Gets the allowance of a MintCap object.
/// Returns 0 if the MintCap object is unauthorized.
public fun mint_allowance<T>(state: &TokenState<T>, mint_cap: ID): (u64, bool) {
    if (!state.is_authorized_mint_cap(mint_cap)) return (0, false);
    state.mint_allowances_map.get(&mint_cap).allowance_info()
}

/// Returns the total amount of Coin<T> in circulation.
public fun total_supply<T>(state: &TokenState<T>): u64 {
    // state.borrow_treasury_cap().total_supply()
    state.treasury_cap.total_supply()
}

/// Checks if a MintCap object is authorized to mint.
public fun is_authorized_mint_cap<T>(state: &TokenState<T>, id: ID): bool {
    state.mint_allowances_map.contains(&id)
}

/// Creates a MintCap object.
/// - Only callable by the master minter.
/// - Only callable if the Treasury object is compatible with this package.
fun create_mint_cap<T>(
    // state: &mut TokenState<T>,
    ctx: &mut TxContext,
): MintCap<T> {
    // treasury.assert_is_compatible();
    // assert!(treasury.roles.master_minter() == ctx.sender(), ENotMasterMinter);
    let mint_cap = MintCap { id: object::new(ctx) };
    event::emit(MintCapCreated<T> {
        mint_cap: object::id(&mint_cap)
    });
    mint_cap
}

/// Convenience function that
/// 1. creates a MintCap
/// 2. transfers the MintCap object to a minter
///
/// - Only callable by the master minter.
/// - Only callable if the Treasury object is compatible with this package.
public fun configure_new_minter<T>(
    state: &mut TokenState<T>,
    _: &OwnerCap<T>,
    minter: address,
    allowance: u64,
    is_unlimited: bool,
    ctx: &mut TxContext,
) {
    let mint_cap = create_mint_cap<T>(ctx);
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
/// - Only callable by the MintCap's controller.
/// - Only callable when not paused.
/// - Only callable if the Treasury object is compatible with this package.
public fun increment_mint_allowance<T>(
    state: &mut TokenState<T>,
    _: &OwnerCap<T>,
    mint_cap_id: ID,
    deny_list: &DenyList,
    allowance_increment: u64,
    _ctx: &TxContext,
) {
    // treasury.assert_is_compatible();

    assert!(!is_paused<T>(deny_list), EPaused);
    assert!(allowance_increment > 0, EZeroAmount);
    assert!(state.is_authorized_mint_cap(mint_cap_id), EUnauthorizedMintCap);

    assert!(state.mint_allowances_map.get(&mint_cap_id).is_unlimited(), ECannotIncreaseUnlimitedAllowance);

    let new_allowance = state.mint_allowances_map.get(&mint_cap_id).value();

    event::emit(MinterAllowanceIncremented<T> {
        mint_cap: mint_cap_id,
        allowance_increment,
        new_allowance,
    });
}

/// Increment allowance for a MintCap
/// - Only callable by the MintCap's controller.
/// - Only callable when not paused.
/// - Only callable if the Treasury object is compatible with this package.
public fun set_unlimited_mint_allowances<T>(
    state: &mut TokenState<T>,
    _: &OwnerCap<T>,
    mint_cap_id: ID,
    deny_list: &DenyList,
    _ctx: &TxContext,
) {
    // treasury.assert_is_compatible();

    assert!(!is_paused<T>(deny_list), EPaused);
    assert!(state.is_authorized_mint_cap(mint_cap_id), EUnauthorizedMintCap);

    state.mint_allowances_map.get_mut(&mint_cap_id).set(0, true);

    event::emit(MinterUnlimitedAllowanceSet<T> {
        mint_cap: mint_cap_id,
    });
}

public fun get_all_mint_caps<T>(
    state: &TokenState<T>,
): vector<ID> {
    state.mint_allowances_map.keys()
}

/// Mints a Coin object with a specified amount (limited to the MintCap's allowance)
/// to a recipient address, increasing the total supply.
/// - Only callable by a minter.
/// - Only callable when not paused.
/// - Only callable if minter is not blocklisted.
/// - Only callable if recipient is not blocklisted.
/// - Only callable if the Treasury object is compatible with this package.
public fun mint_and_transfer<T>(
    state: &mut TokenState<T>,
    mint_cap: &MintCap<T>,
    deny_list: &DenyList,
    amount: u64,
    recipient: address,
    ctx: &mut TxContext
) {
    // treasury.assert_is_compatible();

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
    ctx: &mut TxContext
): Coin<T> {
    // treasury.assert_is_compatible();

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
    assert!(mint_allowance.is_unlimited() || mint_allowance.value() >= amount, EInsufficientAllowance);
    if (!mint_allowance.is_unlimited()) {
        mint_allowance.decrease(amount);
    };
}

/// Burns a Coin object, decreasing the total supply.
/// - Only callable by a minter.
/// - Only callable when not paused.
/// - Only callable if minter is not blocklisted.
/// - Only callable if the Treasury object is compatible with this package.
public fun burn<T>(
    state: &mut TokenState<T>,
    mint_cap: &MintCap<T>,
    deny_list: &DenyList,
    coin: Coin<T>,
    from: address,
    ctx: &TxContext,
) {
    // treasury.assert_is_compatible();

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
        amount
    });
}

public fun blocklist<T>(
    state: &mut TokenState<T>,
    owner_cap: &OwnerCap<T>,
    deny_list: &mut DenyList,
    addr: address,
    ctx: &mut TxContext
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    // treasury.assert_is_compatible();

    if (!is_blocklisted<T>(deny_list, addr)) {
        coin::deny_list_v2_add<T>(deny_list, borrow_deny_cap_mut(state), addr, ctx);
    };
    event::emit(Blocklisted<T> {
        address: addr,
    })
}

/// Unblocklists an address.
/// - Only callable by the blocklister.
/// - Only callable if the Treasury object is compatible with this package.
public fun unblocklist<T>(
    state: &mut TokenState<T>,
    owner_cap: &OwnerCap<T>,
    deny_list: &mut DenyList,
    addr: address,
    ctx: &mut TxContext
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    // treasury.assert_is_compatible();

    if (is_blocklisted<T>(deny_list, addr)) {
        coin::deny_list_v2_remove<T>(deny_list, borrow_deny_cap_mut(state), addr, ctx);
    };
    event::emit(Unblocklisted<T> {
        address: addr
    })
}


/// Triggers stopped state; pause all transfers.
/// - Only callable by the pauser.
/// - Only callable if the Treasury object is compatible with this package.
public fun pause<T>(
    state: &mut TokenState<T>,
    owner_cap: &OwnerCap<T>,
    deny_list: &mut DenyList,
    ctx: &mut TxContext
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    // treasury.assert_is_compatible();

    assert!(state.deny_cap.is_some(), EDenyCapNotFound);

    if (!is_paused<T>(deny_list)) {
        coin::deny_list_v2_enable_global_pause(deny_list, borrow_deny_cap_mut(state), ctx);
    };
    event::emit(Paused<T> {});
}

/// Restores normal state; unpause all transfers.
/// - Only callable by the pauser.
/// - Only callable if the Treasury object is compatible with this package.
public fun unpause<T>(
    state: &mut TokenState<T>,
    owner_cap: &OwnerCap<T>,
    deny_list: &mut DenyList,
    ctx: &mut TxContext,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    // treasury.assert_is_compatible();

    if (is_paused<T>(deny_list)) {
        coin::deny_list_v2_disable_global_pause(deny_list, borrow_deny_cap_mut(state), ctx);
    };
    event::emit(Unpaused<T> {});
}

public fun destroy_managed_token<T>(
    owner_cap: OwnerCap<T>,
    state: TokenState<T>,
    ctx: &mut TxContext,
): (TreasuryCap<T>, Option<DenyCapV2<T>>) {
    assert!(object::id(&owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);

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

    // TODO: instead of returning the treasury cap, we can simply send it to ctx sender
    (treasury_cap, deny_cap)
}

/// Returns an immutable reference of the TreasuryCap.
public fun borrow_treasury_cap<T>(owner_cap: &OwnerCap<T>, state: &TokenState<T>): &TreasuryCap<T> {
    // state.assert_treasury_cap_exists();
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&state.ownable_state), EInvalidOwnerCap);
    &state.treasury_cap
}

/// Returns a mutable reference of the DenyCap.
fun borrow_deny_cap_mut<T>(state: &mut TokenState<T>): &mut DenyCapV2<T> {
    // state.assert_treasury_cap_exists();
    assert!(state.deny_cap.is_some(), EDenyCapNotFound);
    state.deny_cap.borrow_mut()
}

// ================================================================
// |                      Ownable Functions                       |
// ================================================================

public entry fun transfer_ownership<T>(
    state: &mut TokenState<T>,
    owner_cap: &OwnerCap<T>,
    new_owner: address,
    ctx: &mut TxContext,
) {
    ownable::transfer_ownership(owner_cap, &mut state.ownable_state, new_owner, ctx);
}

public entry fun accept_ownership<T>(
    state: &mut TokenState<T>,
    ctx: &mut TxContext,
) {
    ownable::accept_ownership(&mut state.ownable_state, ctx);
}

public fun accept_ownership_from_object<T>(
    state: &mut TokenState<T>,
    from: &mut UID,
    ctx: &mut TxContext,
) {
    ownable::accept_ownership_from_object(&mut state.ownable_state, from, ctx);
}

public fun execute_ownership_transfer<T>(
    owner_cap: OwnerCap<T>,
    ownable_state: &mut OwnableState<T>,
    registry: &mut Registry,
    to: address,
    ctx: &mut TxContext,
) {
    ownable::execute_ownership_transfer(owner_cap, ownable_state, registry, to, ctx);
}

public fun mcms_register_entrypoint<T>(
    registry: &mut Registry,
    state: &mut TokenState<T>,
    owner_cap: OwnerCap<T>,
    ctx: &mut TxContext,
) {
    ownable::set_owner(&owner_cap, &mut state.ownable_state, @mcms, ctx);

    mcms_registry::register_entrypoint(
        registry,
        ownable::mcms_callback(),
        option::some(owner_cap),
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

public fun mcms_entrypoint<T>(
    state: &mut TokenState<T>,
    registry: &mut Registry,
    deny_list: &mut DenyList,
    params: ExecutingCallbackParams, // hot potato
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<
        ownable::McmsCallback,
        OwnerCap<T>,
    >(
        registry,
        ownable::mcms_callback(),
        params,
    );

    let function_bytes = *function.as_bytes();
    let mut stream = bcs_stream::new(data);

    if (function_bytes == b"blocklist") {
        let addr = bcs_stream::deserialize_address(&mut stream);
        bcs_stream::assert_is_consumed(&stream);
        blocklist(state, owner_cap, deny_list, addr, ctx);
    } else if (function_bytes == b"unblocklist") {
        let addr = bcs_stream::deserialize_address(&mut stream);
        bcs_stream::assert_is_consumed(&stream);
        unblocklist(state, owner_cap, deny_list, addr, ctx);
    } else if (function_bytes == b"pause") {
        bcs_stream::assert_is_consumed(&stream);
        pause(state, owner_cap, deny_list, ctx);
    } else if (function_bytes == b"unpause") {
        bcs_stream::assert_is_consumed(&stream);
        unpause(state, owner_cap, deny_list, ctx);
    } else {
        abort EInvalidFunction
    }
}
