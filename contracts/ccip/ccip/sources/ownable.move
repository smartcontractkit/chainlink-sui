module ccip::ownable;

use mcms::bcs_stream;
use mcms::mcms_registry::{Self, ExecutingCallbackParams, Registry};
use sui::event;

public struct OwnerCap has key, store {
    id: UID,
}

public struct OwnableState has key {
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
const EUnknownFunction: u64 = 3;
const ENoPendingTransfer: u64 = 4;
const ETransferAlreadyAccepted: u64 = 5;
const EOwnerChanged: u64 = 6;
const EProposedOwnerMismatch: u64 = 7;
const ETransferNotAccepted: u64 = 8;

public struct OWNABLE has drop {}

fun init(_witness: OWNABLE, ctx: &mut TxContext) {
    let owner = ctx.sender();

    transfer::share_object(OwnableState {
        id: object::new(ctx),
        owner,
        pending_transfer: option::none(),
    });

    transfer::transfer(OwnerCap { id: object::new(ctx) }, owner);
}

public fun initialize(
    owner_cap: OwnerCap,
    state: &mut OwnableState,
    registry: &mut Registry,
    ctx: &mut TxContext,
) {
    let registry_address = object::id_to_address(&object::id(registry));
    state.owner = registry_address;

    mcms_registry::register_entrypoint(
        registry,
        McmsCallback {},
        option::some(owner_cap),
        ctx,
    );
}

public fun transfer_ownership(
    _: &OwnerCap,
    state: &mut OwnableState,
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

public fun accept_ownership(state: &mut OwnableState, ctx: &mut TxContext) {
    accept_ownership_internal(state, ctx.sender());
}

/// UID is a privileged type that is only accessible by the object owner.
public fun accept_ownership_from_object(state: &mut OwnableState, from: &mut UID) {
    accept_ownership_internal(state, from.to_address());
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
public fun execute_ownership_transfer(
    owner_cap: OwnerCap,
    state: &mut OwnableState,
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
            McmsCallback {},
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

// ================================ MCMS Entrypoint ================================ //

public struct McmsCallback has drop {}

public fun mcms_entrypoint(
    registry: &mut Registry,
    state: &mut OwnableState,
    params: ExecutingCallbackParams, // hot potato
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback, OwnerCap>(
        registry,
        McmsCallback {},
        params,
    );

    let function_bytes = *function.as_bytes();
    let mut stream = bcs_stream::new(data);

    if (function_bytes == b"transfer_ownership") {
        let to = bcs_stream::deserialize_address(&mut stream);
        bcs_stream::assert_is_consumed(&stream);
        transfer_ownership(owner_cap, state, to, ctx);
    } else if (function_bytes == b"accept_ownership") {
        accept_ownership(state, ctx);
    } else if (function_bytes == b"accept_ownership_from_object") {
        let from = bcs_stream::deserialize_address(&mut stream);
        bcs_stream::assert_is_consumed(&stream);
        accept_ownership_internal(state, from);
    } else if (function_bytes == b"execute_ownership_transfer") {
        let to = bcs_stream::deserialize_address(&mut stream);
        bcs_stream::assert_is_consumed(&stream);
        let owner_cap = mcms_registry::release_cap(registry, McmsCallback {});
        // Transfer `OwnerCap` from MCMS to the new owner.
        execute_ownership_transfer(owner_cap, state, registry, to, ctx);
    } else {
        abort EUnknownFunction
    };
}

// =================== Test Functions =================== //

#[test_only]
public fun test_init(ctx: &mut TxContext) {
    init(OWNABLE {}, ctx);
}
