module ccip::nonce_manager;

use ccip::ownable::OwnerCap;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::upgrade_registry::verify_function_allowed;
use std::string::{Self, String};
use sui::table::{Self, Table};

const VERSION: u8 = 1;

// store this cap to onramp
public struct NonceManagerCap has key, store {
    id: UID,
}

public struct NonceManagerState has key, store {
    id: UID,
    // dest chain selector -> sender -> nonce
    outbound_nonces: Table<u64, Table<address, u64>>,
}

const EAlreadyInitialized: u64 = 1;
const EInvalidOwnerCap: u64 = 2;

public fun type_and_version(): String {
    string::utf8(b"NonceManager 1.6.0")
}

#[allow(lint(self_transfer))]
public fun initialize(ref: &mut CCIPObjectRef, owner_cap: &OwnerCap, ctx: &mut TxContext) {
    assert!(object::id(owner_cap) == state_object::owner_cap_id(ref), EInvalidOwnerCap);
    assert!(!state_object::contains<NonceManagerState>(ref), EAlreadyInitialized);

    let state = NonceManagerState {
        id: object::new(ctx),
        outbound_nonces: table::new(ctx),
    };
    let cap = NonceManagerCap {
        id: object::new(ctx),
    };
    state_object::add(ref, owner_cap, state, ctx);
    transfer::transfer(cap, ctx.sender());
}

public fun get_outbound_nonce(ref: &CCIPObjectRef, dest_chain_selector: u64, sender: address): u64 {
    verify_function_allowed(
        ref,
        string::utf8(b"nonce_manager"),
        string::utf8(b"get_outbound_nonce"),
        VERSION,
    );
    let state = state_object::borrow<NonceManagerState>(ref);

    if (!state.outbound_nonces.contains(dest_chain_selector)) {
        return 0
    };

    let dest_chain_nonces = &state.outbound_nonces[dest_chain_selector];
    if (!dest_chain_nonces.contains(sender)) {
        return 0
    };
    dest_chain_nonces[sender]
}

public fun get_incremented_outbound_nonce(
    ref: &mut CCIPObjectRef,
    _: &NonceManagerCap,
    dest_chain_selector: u64,
    sender: address,
    ctx: &mut TxContext,
): u64 {
    verify_function_allowed(
        ref,
        string::utf8(b"nonce_manager"),
        string::utf8(b"get_incremented_outbound_nonce"),
        VERSION,
    );
    let state = state_object::borrow_mut<NonceManagerState>(ref);

    if (!state.outbound_nonces.contains(dest_chain_selector)) {
        state
            .outbound_nonces
            .add(
                dest_chain_selector,
                table::new(ctx),
            );
    };
    let dest_chain_nonces = table::borrow_mut(&mut state.outbound_nonces, dest_chain_selector);
    if (!dest_chain_nonces.contains(sender)) {
        dest_chain_nonces.add(sender, 0);
    };

    let nonce_ref = table::borrow_mut(dest_chain_nonces, sender);
    let incremented_nonce = *nonce_ref + 1;
    *nonce_ref = incremented_nonce;
    incremented_nonce
}
