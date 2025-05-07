module mcms::mcms_account;

use sui::event;

public struct OwnerCap has key, store {
    id: UID,
}

public struct AccountState has key {
    id: UID,
    owner: address,
    owner_cap: OwnerCap,
    pending_owner: address,
}

// =================== Events =================== //

public struct OwnershipTransferRequested has copy, drop {
    from: address,
    to: address,
}

public struct OwnershipTransferred has copy, drop {
    from: address,
    to: address,
}

const EUnauthorized: u64 = 1;
const ECannotTransferToSelf: u64 = 2;
const EMustBeProposedOwner: u64 = 3;

public struct MCMS_ACCOUNT has drop {}

fun init(_witness: MCMS_ACCOUNT, ctx: &mut TxContext) {
    transfer::share_object(AccountState {
        id: object::new(ctx),
        owner: ctx.sender(),
        owner_cap: OwnerCap { id: object::new(ctx) },
        pending_owner: ctx.sender(),
    });
}

public fun transfer_ownership(state: &mut AccountState, to: address, ctx: &mut TxContext) {
    assert!(ctx.sender() == state.owner, EUnauthorized);
    assert!(state.owner != to, ECannotTransferToSelf);

    state.pending_owner = to;

    event::emit(OwnershipTransferRequested { from: state.owner, to });
}

public fun transfer_ownership_from_object(state: &mut AccountState, from: &mut UID, to: address) {
    assert!(from.to_address() == state.owner, EUnauthorized);
    assert!(state.owner != to, ECannotTransferToSelf);

    state.pending_owner = to;

    event::emit(OwnershipTransferRequested { from: state.owner, to });
}

public fun accept_ownership(state: &mut AccountState, ctx: &mut TxContext) {
    let caller = ctx.sender();
    assert!(caller == state.pending_owner, EUnauthorized);

    accept_ownership_internal(state, caller);
}

/// UID is a privileged type that is only accessible by the object owner.
public fun accept_ownership_from_object(state: &mut AccountState, from: &mut UID) {
    let caller = from.to_address();
    assert!(caller == state.pending_owner, EUnauthorized);

    accept_ownership_internal(state, caller);
}

fun accept_ownership_internal(state: &mut AccountState, to: address) {
    assert!(state.pending_owner == to, EMustBeProposedOwner);

    let previous_owner = state.owner;
    state.owner = to;
    state.pending_owner = @0x0;

    event::emit(OwnershipTransferred {
        from: previous_owner,
        to,
    });
}

public fun borrow_owner_cap_as_owner(state: &AccountState, ctx: &mut TxContext): &OwnerCap {
    assert!(ctx.sender() == state.owner, EUnauthorized);
    &state.owner_cap
}

public fun borrow_owner_cap_as_object_owner(state: &AccountState, to: &UID): &OwnerCap {
    assert!(to.to_address() == state.owner, EUnauthorized);
    &state.owner_cap
}

// =================== Test Functions =================== //

#[test_only]
public fun create_for_testing(ctx: &mut TxContext): AccountState {
    let owner_cap = OwnerCap {
        id: object::new(ctx),
    };

    let account_state = AccountState {
        id: object::new(ctx),
        owner: ctx.sender(),
        owner_cap,
        pending_owner: @0x0,
    };

    account_state
}
