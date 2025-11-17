module burn_mint_token_pool::ownable;

use mcms::mcms_registry::{Self, Registry, PublisherWrapper};
use sui::dynamic_field as df;
use sui::event;
use sui::package::Publisher;

public struct OwnerCap has key, store {
    id: UID,
}

public struct OwnableState has store {
    owner: address,
    pending_transfer: Option<PendingTransfer>,
    owner_cap_id: ID,
}

public struct PendingTransfer has drop, store {
    from: address,
    to: address,
    accepted: bool,
}

public struct OwnableStateKey has copy, drop, store {}

public struct PublisherKey has copy, drop, store {}

// =================== Events =================== //

public struct NewOwnableStateEvent has copy, drop, store {
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

public(package) fun new(ctx: &mut TxContext): (OwnableState, OwnerCap) {
    let owner = ctx.sender();

    let owner_cap = OwnerCap {
        id: object::new(ctx),
    };

    let state = OwnableState {
        owner,
        pending_transfer: option::none(),
        owner_cap_id: object::id(&owner_cap),
    };

    event::emit(NewOwnableStateEvent {
        owner_cap_id: object::id(&owner_cap),
        owner,
    });

    (state, owner_cap)
}

public fun owner_cap_id(state: &OwnableState): ID {
    state.owner_cap_id
}

public fun owner(state: &OwnableState): address {
    state.owner
}

public fun has_pending_transfer(state: &OwnableState): bool {
    state.pending_transfer.is_some()
}

public fun pending_transfer_from(state: &OwnableState): Option<address> {
    state.pending_transfer.map_ref!(|pending_transfer| pending_transfer.from)
}

public fun pending_transfer_to(state: &OwnableState): Option<address> {
    state.pending_transfer.map_ref!(|pending_transfer| pending_transfer.to)
}

public fun pending_transfer_accepted(state: &OwnableState): Option<bool> {
    state.pending_transfer.map_ref!(|pending_transfer| pending_transfer.accepted)
}

public(package) fun attach_ownable_state(owner_cap: &mut OwnerCap, ownable_state: OwnableState) {
    df::add(&mut owner_cap.id, OwnableStateKey {}, ownable_state);
}

public(package) fun detach_ownable_state(owner_cap: &mut OwnerCap): OwnableState {
    df::remove(&mut owner_cap.id, OwnableStateKey {})
}

public(package) fun attach_publisher(owner_cap: &mut OwnerCap, publisher: Publisher) {
    df::add(&mut owner_cap.id, PublisherKey {}, publisher);
}

public(package) fun borrow_publisher(owner_cap: &OwnerCap): &Publisher {
    df::borrow(&owner_cap.id, PublisherKey {})
}

public(package) fun transfer_ownership(
    owner_cap: &OwnerCap,
    state: &mut OwnableState,
    to: address,
    _ctx: &TxContext,
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

public(package) fun accept_ownership(state: &mut OwnableState, ctx: &TxContext) {
    accept_ownership_internal(state, ctx.sender());
}

/// UID is a privileged type that is only accessible by the object owner.
public(package) fun accept_ownership_from_object(
    state: &mut OwnableState,
    from: &UID,
    _ctx: &TxContext,
) {
    accept_ownership_internal(state, from.to_address());
}

public(package) fun mcms_accept_ownership(
    state: &mut OwnableState,
    mcms: address,
    _ctx: &TxContext,
) {
    accept_ownership_internal(state, mcms);
}

fun accept_ownership_internal(state: &mut OwnableState, caller: address) {
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
public(package) fun execute_ownership_transfer(
    owner_cap: OwnerCap,
    state: &mut OwnableState,
    to: address,
    _ctx: &TxContext,
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

    transfer::transfer(owner_cap, to);

    event::emit(OwnershipTransferred { from: current_owner, to: new_owner });
}

#[allow(lint(custom_state_change))]
public(package) fun execute_ownership_transfer_to_mcms<P: drop>(
    owner_cap: OwnerCap,
    state: &mut OwnableState,
    registry: &mut Registry,
    to: address,
    publisher_wrapper: PublisherWrapper<P>,
    proof: P,
    allowed_modules: vector<vector<u8>>,
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

    mcms_registry::register_entrypoint(
        registry,
        publisher_wrapper,
        proof,
        owner_cap,
        allowed_modules,
        ctx,
    );

    event::emit(OwnershipTransferred { from: current_owner, to: new_owner });
}

public fun destroy(state: OwnableState, owner_cap: OwnerCap, _ctx: &mut TxContext) {
    let OwnableState {
        owner: _,
        pending_transfer: _,
        owner_cap_id: state_owner_cap_id,
    } = state;

    let OwnerCap { id: owner_cap_id } = owner_cap;

    assert!(owner_cap_id.uid_to_inner() == state_owner_cap_id, EInvalidOwnerCap);

    object::delete(owner_cap_id);
}

// =================== Test-only functions =================== //

#[test_only]
public fun create_test_owner_cap(ctx: &mut TxContext): OwnerCap {
    OwnerCap { id: object::new(ctx) }
}

#[test_only]
public fun test_destroy_owner_cap(cap: OwnerCap) {
    let OwnerCap { id } = cap;
    object::delete(id);
}
