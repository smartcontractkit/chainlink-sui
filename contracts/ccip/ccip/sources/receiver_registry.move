module ccip::receiver_registry;

use ccip::ownable::OwnerCap;
use ccip::publisher_wrapper::{Self, PublisherWrapper};
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::upgrade_registry::verify_function_allowed;
use std::ascii;
use std::string::{Self, String};
use std::type_name;
use sui::event;
use sui::linked_table::{Self, LinkedTable};

public struct ReceiverConfig has copy, drop, store {
    module_name: String,
    proof_typename: ascii::String,
}

public struct ReceiverRegistry has key, store {
    id: UID,
    // receiver package id -> receiver config
    receiver_configs: LinkedTable<address, ReceiverConfig>,
}

public struct ReceiverRegistered has copy, drop {
    receiver_package_id: address,
    receiver_module_name: String,
    proof_typename: ascii::String,
}

public struct ReceiverUnregistered has copy, drop {
    receiver_package_id: address,
}

const EAlreadyRegistered: u64 = 1;
const EAlreadyInitialized: u64 = 2;
const EUnknownReceiver: u64 = 3;
const EInvalidOwnerCap: u64 = 4;

const VERSION: u8 = 1;

public fun type_and_version(): String {
    string::utf8(b"ReceiverRegistry 1.6.0")
}

public fun initialize(ref: &mut CCIPObjectRef, owner_cap: &OwnerCap, ctx: &mut TxContext) {
    assert!(object::id(owner_cap) == state_object::owner_cap_id(ref), EInvalidOwnerCap);
    assert!(!state_object::contains<ReceiverRegistry>(ref), EAlreadyInitialized);
    let state = ReceiverRegistry {
        id: object::new(ctx),
        receiver_configs: linked_table::new(ctx),
    };

    state_object::add(ref, owner_cap, state, ctx);
}

public fun register_receiver<ProofType: drop>(
    ref: &mut CCIPObjectRef,
    publisher_wrapper: PublisherWrapper<ProofType>,
    _proof: ProofType,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"receiver_registry"),
        string::utf8(b"register_receiver"),
        VERSION,
    );
    let registry = state_object::borrow_mut<ReceiverRegistry>(ref);
    let proof_typename = type_name::with_defining_ids<ProofType>();
    let receiver_module_name = std::string::from_ascii(type_name::module_string(&proof_typename));
    let receiver_package_id = publisher_wrapper::get_package_address(publisher_wrapper);
    assert!(!registry.receiver_configs.contains(receiver_package_id), EAlreadyRegistered);

    let receiver_config = ReceiverConfig {
        module_name: receiver_module_name,
        proof_typename: proof_typename.into_string(),
    };
    registry.receiver_configs.push_back(receiver_package_id, receiver_config);

    event::emit(ReceiverRegistered {
        receiver_package_id,
        receiver_module_name,
        proof_typename: proof_typename.into_string(),
    });
}

public fun unregister_receiver(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    receiver_package_id: address,
    _: &TxContext,
) {
    verify_function_allowed(
        ref,
        string::utf8(b"receiver_registry"),
        string::utf8(b"unregister_receiver"),
        VERSION,
    );
    assert!(object::id(owner_cap) == state_object::owner_cap_id(ref), EInvalidOwnerCap);

    let registry = state_object::borrow_mut<ReceiverRegistry>(ref);

    assert!(registry.receiver_configs.contains(receiver_package_id), EUnknownReceiver);

    registry.receiver_configs.remove(receiver_package_id);

    event::emit(ReceiverUnregistered {
        receiver_package_id,
    });
}

public fun is_registered_receiver(ref: &CCIPObjectRef, receiver_package_id: address): bool {
    verify_function_allowed(
        ref,
        string::utf8(b"receiver_registry"),
        string::utf8(b"is_registered_receiver"),
        VERSION,
    );
    let registry = state_object::borrow<ReceiverRegistry>(ref);
    registry.receiver_configs.contains(receiver_package_id)
}

public fun get_receiver_config(ref: &CCIPObjectRef, receiver_package_id: address): ReceiverConfig {
    verify_function_allowed(
        ref,
        string::utf8(b"receiver_registry"),
        string::utf8(b"get_receiver_config"),
        VERSION,
    );
    let registry = state_object::borrow<ReceiverRegistry>(ref);

    assert!(registry.receiver_configs.contains(receiver_package_id), EUnknownReceiver);
    *registry.receiver_configs.borrow(receiver_package_id)
}

public fun get_receiver_config_fields(rc: ReceiverConfig): (String, ascii::String) {
    (rc.module_name, rc.proof_typename)
}

// this will return empty string if the receiver is not registered.
public fun get_receiver_info(
    ref: &CCIPObjectRef,
    receiver_package_id: address,
): (String, ascii::String) {
    verify_function_allowed(
        ref,
        string::utf8(b"receiver_registry"),
        string::utf8(b"get_receiver_info"),
        VERSION,
    );
    let registry = state_object::borrow<ReceiverRegistry>(ref);

    if (registry.receiver_configs.contains(receiver_package_id)) {
        let receiver_config = registry.receiver_configs.borrow(receiver_package_id);
        return (receiver_config.module_name, receiver_config.proof_typename)
    };

    (string::utf8(b""), ascii::string(b""))
}
