module mcms::mcms_account;

use mcms::mcms_proof;
use mcms::mcms_registry::{Self, Registry};
use sui::event;

public struct OwnerCap has key, store {
    id: UID,
}

public struct AccountState has key {
    id: UID,
    owner: address,
    pending_transfer: Option<PendingTransfer>,
}

public struct PendingTransfer has drop, store {
    from: address,
    to: address,
    accepted: bool,
}

// =================== Events =================== //

public struct OwnershipTransferRequested has copy, drop {
    from: address,
    to: address,
}

public struct OwnershipTransferAccepted has copy, drop {
    from: address,
    to: address,
}

public struct OwnershipTransferred has copy, drop {
    from: address,
    to: address,
}

const ECannotTransferToSelf: u64 = 1;
const EMustBeProposedOwner: u64 = 2;
const ENoPendingTransfer: u64 = 3;
const EOwnerChanged: u64 = 4;
const EProposedOwnerMismatch: u64 = 5;
const ETransferNotAccepted: u64 = 6;
const ETransferAlreadyAccepted: u64 = 7;

public struct MCMS_ACCOUNT has drop {}

fun init(_witness: MCMS_ACCOUNT, ctx: &mut TxContext) {
    transfer::share_object(AccountState {
        id: object::new(ctx),
        owner: ctx.sender(),
        pending_transfer: option::none(),
    });

    transfer::share_object(OwnerCap { id: object::new(ctx) });
}

public fun transfer_ownership(
    _: &OwnerCap,
    state: &mut AccountState,
    to: address,
    _ctx: &mut TxContext,
) {
    assert!(state.owner != to, ECannotTransferToSelf);

    state.pending_transfer =
        option::some(PendingTransfer {
            from: state.owner,
            to,
            accepted: false,
        });

    event::emit(OwnershipTransferRequested { from: state.owner, to });
}

/// Transfers ownership back to @mcms/timelock.
public entry fun transfer_ownership_to_self(
    owner_cap: &OwnerCap,
    state: &mut AccountState,
    ctx: &mut TxContext,
) {
    transfer_ownership(owner_cap, state, @mcms, ctx);
}

public fun accept_ownership(state: &mut AccountState, ctx: &mut TxContext) {
    accept_ownership_internal(state, ctx.sender());
}

public fun accept_ownership_as_mcms_timelock<T: drop>(
    state: &mut AccountState,
    witness: T,
    _ctx: &mut TxContext,
) {
    mcms_proof::assert_is_mcms_timelock(witness);
    accept_ownership_internal(state, @mcms);
}

/// UID is a privileged type that is only accessible by the object owner.
public fun accept_ownership_from_object(state: &mut AccountState, from: &mut UID) {
    accept_ownership_internal(state, from.to_address());
}

fun accept_ownership_internal(state: &mut AccountState, caller: address) {
    assert!(state.pending_transfer.is_some(), ENoPendingTransfer);

    let pending_transfer = state.pending_transfer.borrow_mut();
    let current_owner = state.owner;

    // check that the owner has not changed from a direct call to 0x1::transfer::public_transfer,
    // in which case the transfer flow should be restarted.
    assert!(current_owner == pending_transfer.from, EOwnerChanged);
    assert!(caller == pending_transfer.to, EMustBeProposedOwner);
    assert!(!pending_transfer.accepted, ETransferAlreadyAccepted);

    pending_transfer.accepted = true;

    event::emit(OwnershipTransferAccepted { from: pending_transfer.from, to: caller });
}

#[allow(lint(custom_state_change))]
public fun execute_ownership_transfer(
    owner_cap: OwnerCap,
    state: &mut AccountState,
    registry: &mut Registry,
    to: address,
    ctx: &mut TxContext,
) {
    assert!(state.pending_transfer.is_some(), ENoPendingTransfer);

    let pending_transfer = state.pending_transfer.extract();
    let current_owner = state.owner;
    let new_owner = pending_transfer.to;

    // check that the owner has not changed from a direct call to 0x1::transfer::public_transfer,
    // in which case the transfer flow should be restarted.
    assert!(pending_transfer.from == current_owner, EOwnerChanged);
    assert!(new_owner == to, EProposedOwnerMismatch);
    assert!(pending_transfer.accepted, ETransferNotAccepted);

    // if the new owner is mcms, we need to add the `OwnerCap` to the registry.
    if (new_owner == @mcms) {
        mcms_registry::register_entrypoint(
            registry,
            mcms_proof::create_mcms_proof(),
            option::some(owner_cap),
            ctx,
        );
    } else {
        transfer::transfer(owner_cap, new_owner);
    };

    state.owner = new_owner;
    state.pending_transfer = option::none();

    event::emit(OwnershipTransferred { from: current_owner, to: new_owner });
}

// =================== Test Functions =================== //

#[test_only]
public fun test_init(ctx: &mut TxContext) {
    init(MCMS_ACCOUNT {}, ctx);
}
