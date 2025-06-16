module managed_token::managed_token;

use std::string::{Self, String};

use sui::coin::{
    Self, Coin, DenyCapV2, TreasuryCap,
    deny_list_v2_is_global_pause_enabled_next_epoch as is_paused,
    deny_list_v2_contains_next_epoch as is_blocklisted,
};
use sui::deny_list::{DenyList};
// use sui::dynamic_object_field as dof;
use sui::event;
use sui::table::{Self, Table};

use managed_token::allowlist::{Self, AllowlistState};
use managed_token::mint_allowance::{Self, MintAllowance};

public struct TokenState<phantom T> has key, store {
    id: UID,
    treasury_cap: TreasuryCap<T>,
    deny_cap: Option<DenyCapV2<T>>,
    /// A map of { controller address => MintCap ID that it controls }.
    controllers: Table<address, ID>,
    /// A map of { authorized MintCap ID => its MintAllowance }.
    mint_allowances: Table<ID, MintAllowance<T>>,
    allowed_minters: AllowlistState,
    allowed_burners: AllowlistState,
}

public struct OwnerCap<phantom T> has key, store {
    id: UID,
    state_id: ID,
}

/// An object representing the ability to mint up to an allowance
/// specified in the Treasury.
/// The privilege can be revoked by the master minter.
public struct MintCap<phantom T> has key, store {
    id: UID,
}

// /// Key for retrieving the `TreasuryCap` stored in a `Treasury<T>` dynamic object field
// public struct TreasuryCapKey has copy, store, drop {}
// /// Key for retrieving `DenyCap` stored in a `Treasury<T>` dynamic object field
// public struct DenyCapKey has copy, store, drop {}

// === Events ===

public struct MintCapCreated<phantom T> has copy, drop {
    mint_cap: ID,
}

public struct ControllerConfigured<phantom T> has copy, drop {
    controller: address,
    mint_cap: ID,
}

public struct ControllerRemoved<phantom T> has copy, drop {
    controller: address,
}

public struct MinterConfigured<phantom T> has copy, drop {
    controller: address,
    mint_cap: ID,
    allowance: u64,
    is_unlimited: bool,
}

public struct MinterRemoved<phantom T> has copy, drop {
    controller: address,
    mint_cap: ID,
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
    controller: address,
    mint_cap: ID,
    allowance_increment: u64,
    new_allowance: u64,
}

const EControllerAlreadyConfigured: u64 = 0;
const EDeniedAddress: u64 = 1;
const EDenyCapNotFound: u64 = 2;
const EInsufficientAllowance: u64 = 3;
const EInvalidOwnerCap: u64 = 4;
const ENotController: u64 = 5;
const EPaused: u64 = 6;
const EUnauthorizedMintCap: u64 = 7;
const EZeroAmount: u64 = 8;
const ECannotIncreaseUnlimitedAllowance: u64 = 9;

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

fun initialize_internal<T: drop>(
    treasury_cap: TreasuryCap<T>,
    deny_cap: Option<DenyCapV2<T>>,
    ctx: &mut TxContext,
) {
    let state = TokenState<T> {
        id: object::new(ctx),
        treasury_cap,
        deny_cap,
        controllers: table::new(ctx),
        mint_allowances: table::new(ctx),
        allowed_minters: allowlist::new(vector[]),
        allowed_burners: allowlist::new(vector[]),
    };
    // dof::add(&mut state.id, TreasuryCapKey {}, treasury_cap);

    let owner_cap = OwnerCap<T> {
        id: object::new(ctx),
        state_id: object::id(&state),
    };

    transfer::share_object(state);
    transfer::public_transfer(owner_cap, ctx.sender());
}

/// Gets the corresponding MintCap ID attached to a controller address.
/// Returns option::none() when input is not a valid controller
public fun get_mint_cap_id<T>(state: &TokenState<T>, controller: address): Option<ID> {
    if (!state.controllers.contains(controller)) return option::none();
    option::some(*state.controllers.borrow(controller))
}

/// Gets the allowance of a MintCap object.
/// Returns 0 if the MintCap object is unauthorized.
public fun mint_allowance<T>(state: &TokenState<T>, mint_cap: ID): (u64, bool) {
    if (!state.is_authorized_mint_cap(mint_cap)) return (0, false);
    state.mint_allowances.borrow(mint_cap).allowance_info()
}

/// Returns the total amount of Coin<T> in circulation.
public fun total_supply<T>(state: &TokenState<T>): u64 {
    // state.borrow_treasury_cap().total_supply()
    state.treasury_cap.total_supply()
}

/// Checks if a MintCap object is authorized to mint.
public fun is_authorized_mint_cap<T>(state: &TokenState<T>, id: ID): bool {
    state.mint_allowances.contains(id)
}

// /// [Package private] Ensures that TreasuryCap exists.
// public(package) fun assert_treasury_cap_exists<T>(state: &TokenState<T>) {
//     assert!(dof::exists_with_type<_, TreasuryCap<T>>(&state.id, TreasuryCapKey {}), ETreasuryCapNotFound);
// }

/// Checks if an address is a mint controller.
fun is_controller<T>(state: &TokenState<T>, controller_addr: address): bool {
    state.controllers.contains(controller_addr)
}

public fun get_allowed_minters<T>(state: &TokenState<T>): vector<address> {
    allowlist::get_allowlist(&state.allowed_minters)
}

public fun get_allowed_burners<T>(state: &TokenState<T>): vector<address> {
    allowlist::get_allowlist(&state.allowed_burners)
}

public fun is_minter_allowed<T>(state: &TokenState<T>, minter: address): bool {
    allowlist::is_allowed(&state.allowed_minters, minter)
}

public fun is_burner_allowed<T>(state: &TokenState<T>, burner: address): bool {
    allowlist::is_allowed(&state.allowed_burners, burner)
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
/// 2. configures the controller for this MintCap object
/// 3. transfers the MintCap object to a minter
///
/// - Only callable by the master minter.
/// - Only callable if the Treasury object is compatible with this package.
public fun configure_new_controller<T>(
    state: &mut TokenState<T>,
    owner_cap: &OwnerCap<T>,
    controller: address,
    minter: address,
    ctx: &mut TxContext,
) {
    let mint_cap = create_mint_cap<T>(ctx);
    configure_controller(state,  owner_cap, controller, object::id(&mint_cap), ctx);
    transfer::transfer(mint_cap, minter)
}

/// Configures a controller of a MintCap object.
/// - Only callable by the master minter.
/// - Only callable if the Treasury object is compatible with this package.
public fun configure_controller<T>(
    state: &mut TokenState<T>,
    _: &OwnerCap<T>,
    controller: address,
    mint_cap_id: ID,
    _ctx: &TxContext,
) {
    // treasury.assert_is_compatible();
    // assert!(treasury.roles.master_minter() == ctx.sender(), ENotMasterMinter);
    assert!(!state.is_controller(controller), EControllerAlreadyConfigured);

    state.controllers.add(controller, mint_cap_id);
    event::emit(ControllerConfigured<T> {
        controller,
        mint_cap: mint_cap_id
    });
}

public fun remove_controller<T>(
    state: &mut TokenState<T>,
    _: &OwnerCap<T>,
    controller: address,
    _ctx: &TxContext,
) {
    // treasury.assert_is_compatible();
    // assert!(treasury.roles.master_minter() == ctx.sender(), ENotMasterMinter);
    assert!(state.is_controller(controller), ENotController);

    state.controllers.remove(controller);

    event::emit(ControllerRemoved<T> {
        controller
    });
}

/// Authorizes a MintCap object to mint and burn, and sets its allowance.
/// - Only callable by the MintCap's controller.
/// - Only callable when not paused.
/// - Only callable if the Treasury object is compatible with this package.
public fun configure_minter<T>(
    state: &mut TokenState<T>,
    _: &OwnerCap<T>,
    deny_list: &DenyList,
    new_allowance: u64,
    is_unlimited: bool,
    ctx: &TxContext
) {
    // treasury.assert_is_compatible();

    assert!(!is_paused<T>(deny_list), EPaused);

    let controller = ctx.sender();
    assert!(state.is_controller(controller), ENotController);

    // if the allowance is unlimited, new_allowance should be 0
    // otherwise, new_allowance should be greater than 0
    assert!(is_unlimited != (new_allowance > 0), EZeroAmount);

    let mint_cap_id = *get_mint_cap_id(state, controller).borrow();
    if (!state.mint_allowances.contains(mint_cap_id)) {
        let mut allowance = mint_allowance::new();
        allowance.set(new_allowance, is_unlimited);
        state.mint_allowances.add(mint_cap_id, allowance);
    } else {
        state.mint_allowances.borrow_mut(mint_cap_id).set(new_allowance, is_unlimited);
    };
    event::emit(MinterConfigured<T> {
        controller,
        mint_cap: mint_cap_id,
        allowance: new_allowance,
        is_unlimited,
    });
}

/// Increment allowance for a MintCap
/// - Only callable by the MintCap's controller.
/// - Only callable when not paused.
/// - Only callable if the Treasury object is compatible with this package.
public fun increment_mint_allowance<T>(
    state: &mut TokenState<T>,
    deny_list: &DenyList,
    allowance_increment: u64,
    ctx: &TxContext,
) {
    // treasury.assert_is_compatible();

    assert!(!is_paused<T>(deny_list), EPaused);
    assert!(allowance_increment > 0, EZeroAmount);

    let controller = ctx.sender();
    assert!(state.is_controller(controller), ENotController);

    let mint_cap_id = *get_mint_cap_id(state, controller).borrow();
    assert!(state.is_authorized_mint_cap(mint_cap_id), EUnauthorizedMintCap);

    assert!(!state.mint_allowances.borrow(mint_cap_id).is_unlimited(), ECannotIncreaseUnlimitedAllowance);
    state.mint_allowances.borrow_mut(mint_cap_id).increase(allowance_increment);
    let new_allowance = state.mint_allowances.borrow(mint_cap_id).value();

    event::emit(MinterAllowanceIncremented<T> {
        controller,
        mint_cap: mint_cap_id,
        allowance_increment,
        new_allowance,
    });
}

/// Deauthorizes a MintCap object.
/// - Only callable by the MintCap's controller.
/// - Only callable if the Treasury object is compatible with this package.
public fun remove_minter<T>(
    state: &mut TokenState<T>,
    ctx: &TxContext
) {
    // treasury.assert_is_compatible();

    let controller = ctx.sender();
    assert!(state.is_controller(controller), ENotController);

    let mint_cap_id = *get_mint_cap_id(state, controller).borrow();
    let mint_allowance = state.mint_allowances.remove(mint_cap_id);
    mint_allowance.destroy();
    event::emit(MinterRemoved<T> {
        controller,
        mint_cap: mint_cap_id
    });
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

    let mint_allowance = state.mint_allowances.borrow_mut(mint_cap_id);
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
    assert!(owner_cap.state_id == object::id(state), EInvalidOwnerCap);
    // treasury.assert_is_compatible();
    // assert!(state.roles.blocklister() == ctx.sender(), ENotBlocklister);

    // assert!(state.deny_cap.is_some(), EDenyCapNotFound);
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
    assert!(owner_cap.state_id == object::id(state), EInvalidOwnerCap);
    // treasury.assert_is_compatible();
    // assert!(state.roles.blocklister() == ctx.sender(), ENotBlocklister);

    // assert!(state.deny_cap.is_some(), EDenyCapNotFound);
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
public  fun pause<T>(
    state: &mut TokenState<T>,
    owner_cap: &OwnerCap<T>,
    deny_list: &mut DenyList,
    ctx: &mut TxContext
) {
    assert!(owner_cap.state_id == object::id(state), EInvalidOwnerCap);
    // treasury.assert_is_compatible();

    // assert!(treasury.roles.pauser() == ctx.sender(), ENotPauser);
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
    assert!(owner_cap.state_id == object::id(state), EInvalidOwnerCap);
    // treasury.assert_is_compatible();
    // assert!(treasury.roles().pauser() == ctx.sender(), ENotPauser);

    if (is_paused<T>(deny_list)) {
        coin::deny_list_v2_disable_global_pause(deny_list, borrow_deny_cap_mut(state), ctx);
    };
    event::emit(Unpaused<T> {});
}

public fun destroy_managed_token<T>(
    owner_cap: OwnerCap<T>,
    state: TokenState<T>,
    _ctx: &mut TxContext,
): (TreasuryCap<T>, Option<DenyCapV2<T>>) {
    assert!(owner_cap.state_id == object::id(&state), EInvalidOwnerCap);

    let TokenState<T> {
        id: state_id,
        treasury_cap,
        deny_cap,
        controllers,
        mint_allowances,
        allowed_minters,
        allowed_burners,
    } = state;

    object::delete(state_id);
    controllers.drop(); // drop a potentially non-empty table
    mint_allowances.drop(); // drop a potentially non-empty table
    allowlist::destroy_allowlist(allowed_minters);
    allowlist::destroy_allowlist(allowed_burners);

    let OwnerCap {
        id: owner_cap_id,
        state_id: _,
    } = owner_cap;
    object::delete(owner_cap_id);

    // TODO: instead of returning the treasury cap, we can simply send it to ctx sender
    (treasury_cap, deny_cap)
}

/// Returns an immutable reference of the TreasuryCap.
public fun borrow_treasury_cap<T>(owner_cap: &OwnerCap<T>, state: &TokenState<T>): &TreasuryCap<T> {
    // state.assert_treasury_cap_exists();
    assert!(owner_cap.state_id == object::id(state), EInvalidOwnerCap);
    &state.treasury_cap
}

/// Returns a mutable reference of the DenyCap.
fun borrow_deny_cap_mut<T>(state: &mut TokenState<T>): &mut DenyCapV2<T> {
    // state.assert_treasury_cap_exists();
    assert!(state.deny_cap.is_some(), EDenyCapNotFound);
    state.deny_cap.borrow_mut()
}

