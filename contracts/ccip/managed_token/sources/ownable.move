module managed_token::ownable;

use mcms::mcms_registry::{Self, Registry};
use sui::event;

public struct OwnerCap<phantom T> has key, store {
    id: UID,
}

public struct OwnableState<phantom T> has key, store {
    id: UID,
    owner: address,
    pending_transfer: Option<PendingTransfer>,
    owner_cap_id: ID,
}

public struct PendingTransfer has drop, store {
    from: address,
    to: address,
    accepted: bool,
}

// =================== Events =================== //

public struct NewOwnableStateEvent has copy, drop, store {
    ownable_state_id: ID,
    owner_cap_id: ID,
    owner: address,
}

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

const EInvalidOwnerCap: u64 = 1;
const ECannotTransferToSelf: u64 = 2;
const EMustBeProposedOwner: u64 = 3;
const ENoPendingTransfer: u64 = 4;
const ETransferAlreadyAccepted: u64 = 5;
const EOwnerChanged: u64 = 6;
const EProposedOwnerMismatch: u64 = 7;
const ETransferNotAccepted: u64 = 8;
const ECannotTransferToMcms: u64 = 9;
const EMustTransferToMcms: u64 = 10;

public fun new<T>(ctx: &mut TxContext): (OwnableState<T>, OwnerCap<T>) {
    let owner = ctx.sender();

    let owner_cap = OwnerCap<T> {
        id: object::new(ctx),
    };

    let state = OwnableState {
        id: object::new(ctx),
        owner,
        pending_transfer: option::none(),
        owner_cap_id: object::id(&owner_cap),
    };

    event::emit(NewOwnableStateEvent {
        ownable_state_id: object::id(&state),
        owner_cap_id: object::id(&owner_cap),
        owner,
    });

    (state, owner_cap)
}

public fun owner_cap_id<T>(state: &OwnableState<T>): ID {
    state.owner_cap_id
}

public fun owner<T>(state: &OwnableState<T>): address {
    state.owner
}

public fun has_pending_transfer<T>(state: &OwnableState<T>): bool {
    state.pending_transfer.is_some()
}

public fun pending_transfer_from<T>(state: &OwnableState<T>): Option<address> {
    state.pending_transfer.map_ref!(|pending_transfer| pending_transfer.from)
}

public fun pending_transfer_to<T>(state: &OwnableState<T>): Option<address> {
    state.pending_transfer.map_ref!(|pending_transfer| pending_transfer.to)
}

public fun pending_transfer_accepted<T>(state: &OwnableState<T>): Option<bool> {
    state.pending_transfer.map_ref!(|pending_transfer| pending_transfer.accepted)
}

public fun transfer_ownership<T>(
    owner_cap: &OwnerCap<T>,
    state: &mut OwnableState<T>,
    to: address,
    _ctx: &mut TxContext,
) {
    assert!(object::id(owner_cap) == state.owner_cap_id, EInvalidOwnerCap);
    assert!(state.owner != to, ECannotTransferToSelf);

    state.pending_transfer =
        option::some(PendingTransfer {
            from: state.owner,
            to,
            accepted: false,
        });

    event::emit(OwnershipTransferRequested { from: state.owner, to });
}

public fun accept_ownership<T>(state: &mut OwnableState<T>, ctx: &mut TxContext) {
    accept_ownership_internal(state, ctx.sender());
}

/// UID is a privileged type that is only accessible by the object owner.
public fun accept_ownership_from_object<T>(
    state: &mut OwnableState<T>,
    from: &mut UID,
    _ctx: &mut TxContext,
) {
    accept_ownership_internal(state, from.to_address());
}

public(package) fun mcms_accept_ownership<T>(
    state: &mut OwnableState<T>,
    mcms: address,
    _ctx: &mut TxContext,
) {
    accept_ownership_internal(state, mcms);
}

fun accept_ownership_internal<T>(state: &mut OwnableState<T>, caller: address) {
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
public fun execute_ownership_transfer<T>(
    owner_cap: OwnerCap<T>,
    state: &mut OwnableState<T>,
    to: address,
    _ctx: &mut TxContext,
) {
    assert!(object::id(&owner_cap) == state.owner_cap_id, EInvalidOwnerCap);
    assert!(state.pending_transfer.is_some(), ENoPendingTransfer);

    let pending_transfer = state.pending_transfer.extract();
    let current_owner = state.owner;
    let new_owner = pending_transfer.to;

    // check that the owner has not changed from a direct call to 0x1::transfer::public_transfer,
    // in which case the transfer flow should be restarted.
    assert!(pending_transfer.from == current_owner, EOwnerChanged);
    assert!(new_owner == to, EProposedOwnerMismatch);
    assert!(pending_transfer.accepted, ETransferNotAccepted);

    // Must call `execute_ownership_transfer_to_mcms` instead
    assert!(new_owner != mcms_registry::get_multisig_address(), ECannotTransferToMcms);

    state.owner = to;
    state.pending_transfer = option::none();

    transfer::transfer(owner_cap, to);

    event::emit(OwnershipTransferred { from: current_owner, to: new_owner });
}

#[allow(lint(custom_state_change))]
public fun execute_ownership_transfer_to_mcms<T, P: drop>(
    owner_cap: OwnerCap<T>,
    state: &mut OwnableState<T>,
    registry: &mut Registry,
    to: address,
    proof: P,
    ctx: &mut TxContext,
) {
    assert!(object::id(&owner_cap) == state.owner_cap_id, EInvalidOwnerCap);
    assert!(state.pending_transfer.is_some(), ENoPendingTransfer);

    let pending_transfer = state.pending_transfer.extract();
    let current_owner = state.owner;
    let new_owner = pending_transfer.to;

    // check that the owner has not changed from a direct call to 0x1::transfer::public_transfer,
    // in which case the transfer flow should be restarted.
    assert!(pending_transfer.from == current_owner, EOwnerChanged);
    assert!(new_owner == to, EProposedOwnerMismatch);
    assert!(pending_transfer.accepted, ETransferNotAccepted);
    assert!(to == mcms_registry::get_multisig_address(), EMustTransferToMcms);

    state.owner = to;
    state.pending_transfer = option::none();

    mcms_registry::register_entrypoint(
        registry,
        proof,
        owner_cap,
        ctx,
    );

    event::emit(OwnershipTransferred { from: current_owner, to: new_owner });
}

public fun destroy_ownable_state<T>(state: OwnableState<T>, _ctx: &mut TxContext) {
    let OwnableState {
        id,
        owner: _,
        pending_transfer: _,
        owner_cap_id: _,
    } = state;

    object::delete(id);
}

public fun destroy_owner_cap<T>(owner_cap: OwnerCap<T>, _ctx: &mut TxContext) {
    let OwnerCap { id } = owner_cap;
    object::delete(id);
}
