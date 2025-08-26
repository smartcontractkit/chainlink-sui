module ccip::state_object;

use ccip::ownable::{Self, OwnerCap, OwnableState};
use mcms::bcs_stream;
use mcms::mcms_registry::{Self, Registry, ExecutingCallbackParams};
use std::ascii;
use std::string;
use std::type_name;
use sui::address;
use sui::dynamic_object_field as dof;

const EModuleAlreadyExists: u64 = 1;
const EModuleDoesNotExist: u64 = 2;
const EInvalidFunction: u64 = 3;
const EInvalidOwnerCap: u64 = 4;
const EInvalidRefAddress: u64 = 5;
const EInvalidRegistryAddress: u64 = 6;

public struct CCIPObjectRef has key, store {
    id: UID,
    ownable_state: OwnableState,
}

public struct CCIPObjectRefPointer has key, store {
    id: UID,
    object_ref_id: address,
    owner_cap_id: address,
}

public struct STATE_OBJECT has drop {}

fun init(_witness: STATE_OBJECT, ctx: &mut TxContext) {
    let (ownable_state, owner_cap) = ownable::new(ctx);

    let ref = CCIPObjectRef {
        id: object::new(ctx),
        ownable_state,
    };

    let owner_cap_id = object::id(&owner_cap);

    let pointer = CCIPObjectRefPointer {
        id: object::new(ctx),
        object_ref_id: object::uid_to_address(&ref.id),
        owner_cap_id: object::id_to_address(&owner_cap_id),
    };

    let tn = type_name::get_with_original_ids<STATE_OBJECT>();
    let package_bytes = ascii::into_bytes(tn.get_address());
    let package_id = address::from_ascii_bytes(&package_bytes);

    transfer::share_object(ref);
    transfer::public_transfer(owner_cap, ctx.sender());
    transfer::transfer(pointer, package_id);
}

public fun owner_cap_id(ref: &CCIPObjectRef): ID {
    ref.ownable_state.owner_cap_id()
}

public(package) fun add<T: key + store>(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    obj: T,
    _ctx: &TxContext,
) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&ref.ownable_state), EInvalidOwnerCap);

    let tn = type_name::get<T>();
    assert!(!dof::exists_(&ref.id, tn), EModuleAlreadyExists);
    dof::add(&mut ref.id, tn, obj);
}

public(package) fun contains<T>(ref: &CCIPObjectRef): bool {
    let tn = type_name::get<T>();
    dof::exists_(&ref.id, tn)
}

public(package) fun remove<T: key + store>(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    _ctx: &TxContext,
): T {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&ref.ownable_state), EInvalidOwnerCap);
    let tn = type_name::get<T>();
    assert!(dof::exists_(&ref.id, tn), EModuleDoesNotExist);
    dof::remove(&mut ref.id, tn)
}

public(package) fun borrow<T: key + store>(ref: &CCIPObjectRef): &T {
    let tn = type_name::get<T>();
    dof::borrow(&ref.id, tn)
}

public(package) fun borrow_mut<T: key + store>(ref: &mut CCIPObjectRef): &mut T {
    let tn = type_name::get<T>();
    dof::borrow_mut(&mut ref.id, tn)
}

public fun transfer_ownership(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    to: address,
    ctx: &mut TxContext,
) {
    ownable::transfer_ownership(owner_cap, &mut ref.ownable_state, to, ctx);
}

public fun accept_ownership(ref: &mut CCIPObjectRef, ctx: &mut TxContext) {
    ownable::accept_ownership(&mut ref.ownable_state, ctx);
}

public fun execute_ownership_transfer(
    ref: &mut CCIPObjectRef,
    owner_cap: OwnerCap,
    to: address,
    ctx: &mut TxContext,
) {
    ownable::execute_ownership_transfer(owner_cap, &mut ref.ownable_state, to, ctx);
}

public fun execute_ownership_transfer_to_mcms(
    ref: &mut CCIPObjectRef,
    owner_cap: OwnerCap,
    registry: &mut Registry,
    to: address,
    ctx: &mut TxContext,
) {
    ownable::execute_ownership_transfer_to_mcms(
        owner_cap,
        &mut ref.ownable_state,
        registry,
        to,
        McmsCallback {},
        ctx,
    );
}

public fun owner(ref: &CCIPObjectRef): address {
    ref.ownable_state.owner()
}

public fun has_pending_transfer(ref: &CCIPObjectRef): bool {
    ref.ownable_state.has_pending_transfer()
}

public fun pending_transfer_from(ref: &CCIPObjectRef): Option<address> {
    ref.ownable_state.pending_transfer_from()
}

public fun pending_transfer_to(ref: &CCIPObjectRef): Option<address> {
    ref.ownable_state.pending_transfer_to()
}

public fun pending_transfer_accepted(ref: &CCIPObjectRef): Option<bool> {
    ref.ownable_state.pending_transfer_accepted()
}

// ================================================================
// |                      MCMS Entrypoint                         |
// ================================================================

/// Proof for CCIP admin
public struct CCIPAdminProof has drop {}

public struct McmsCallback has drop {}

fun validate_shared_objects(
    ref: &CCIPObjectRef,
    registry: &Registry,
    stream: &mut bcs_stream::BCSStream,
) {
    let ref_address = bcs_stream::deserialize_address(stream);
    assert!(ref_address == object::id_address(ref), EInvalidRefAddress);
    let registry_address = bcs_stream::deserialize_address(stream);
    assert!(registry_address == object::id_address(registry), EInvalidRegistryAddress);
}

public fun mcms_transfer_ownership(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback, OwnerCap>(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"transfer_ownership"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    validate_shared_objects(ref, registry, &mut stream);

    let to = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    transfer_ownership(ref, owner_cap, to, ctx);
}

public fun mcms_execute_ownership_transfer(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback, OwnerCap>(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"execute_ownership_transfer"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    validate_shared_objects(ref, registry, &mut stream);

    let to = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    let owner_cap = mcms_registry::release_cap(registry, McmsCallback {});
    execute_ownership_transfer(ref, owner_cap, to, ctx);
}

public fun mcms_proof_entrypoint(
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    _ctx: &mut TxContext,
): CCIPAdminProof {
    let (_owner_cap, function, _data) = mcms_registry::get_callback_params<McmsCallback, OwnerCap>(
        registry,
        McmsCallback {},
        params,
    );

    // We validate that the owner cap is registered
    // So we can safely provide a proof that CCIP admin is calling
    assert!(*function.as_bytes() == b"initialize_by_ccip_admin", EInvalidFunction);

    CCIPAdminProof {}
}

// ================================================================
// |                      Test Functions                          |
// ================================================================

#[test_only]
public fun test_init(ctx: &mut TxContext) {
    init(STATE_OBJECT {}, ctx);
}

#[test_only]
public fun pending_transfer(ref: &CCIPObjectRef): (address, address, bool) {
    let from = ownable::pending_transfer_from(&ref.ownable_state);
    let to = ownable::pending_transfer_to(&ref.ownable_state);
    let accepted = ownable::pending_transfer_accepted(&ref.ownable_state);

    (from.get_with_default(@0x0), to.get_with_default(@0x0), accepted.get_with_default(false))
}

#[test_only]
public fun create_ccip_admin_proof_for_test(): CCIPAdminProof {
    CCIPAdminProof {}
}
