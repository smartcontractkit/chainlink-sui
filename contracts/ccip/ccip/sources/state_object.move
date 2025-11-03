module ccip::state_object;

use ccip::ownable::{Self, OwnerCap, OwnableState};
use mcms::bcs_stream;
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, Registry, ExecutingCallbackParams};
use std::ascii;
use std::string;
use std::type_name;
use sui::address;
use sui::derived_object;
use sui::dynamic_object_field as dof;
use sui::package::{Self, UpgradeCap};

const EModuleAlreadyExists: u64 = 1;
const EModuleDoesNotExist: u64 = 2;
const EInvalidFunction: u64 = 3;
const EInvalidOwnerCap: u64 = 4;
const EPackageIdNotFound: u64 = 5;

public struct CCIPObject has key {
    id: UID,
}

public struct CCIPObjectRef has key, store {
    id: UID,
    package_ids: vector<address>,
    ownable_state: OwnableState,
}

public struct CCIPObjectRefPointer has key, store {
    id: UID,
    ccip_object_id: address,
}

public struct STATE_OBJECT has drop {}

fun init(otw: STATE_OBJECT, ctx: &mut TxContext) {
    let mut ccip_object = CCIPObject { id: object::new(ctx) };
    let (ownable_state, mut owner_cap) = ownable::new(&mut ccip_object.id, ctx);

    let mut ref = CCIPObjectRef {
        id: derived_object::claim(&mut ccip_object.id, b"CCIPObjectRef"),
        package_ids: vector[],
        ownable_state,
    };

    let pointer = CCIPObjectRefPointer {
        id: object::new(ctx),
        ccip_object_id: object::id_address(&ccip_object),
    };

    let tn = type_name::with_original_ids<STATE_OBJECT>();
    let package_bytes = ascii::into_bytes(tn.address_string());
    let package_id = address::from_ascii_bytes(&package_bytes);
    ref.package_ids.push_back(package_id);

    transfer::share_object(ref);
    transfer::share_object(ccip_object);

    let publisher = package::claim(otw, ctx);
    ownable::attach_publisher(&mut owner_cap, publisher);

    transfer::public_transfer(owner_cap, ctx.sender());
    transfer::transfer(pointer, package_id);
}

public fun add_package_id(ref: &mut CCIPObjectRef, owner_cap: &OwnerCap, package_id: address) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&ref.ownable_state), EInvalidOwnerCap);
    ref.package_ids.push_back(package_id);
}

public fun remove_package_id(ref: &mut CCIPObjectRef, owner_cap: &OwnerCap, package_id: address) {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&ref.ownable_state), EInvalidOwnerCap);
    let (found, idx) = ref.package_ids.index_of(&package_id);
    assert!(found, EPackageIdNotFound);
    ref.package_ids.remove(idx);
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

    let tn = type_name::with_defining_ids<T>();
    assert!(!dof::exists_(&ref.id, tn), EModuleAlreadyExists);
    dof::add(&mut ref.id, tn, obj);
}

public(package) fun contains<T>(ref: &CCIPObjectRef): bool {
    let tn = type_name::with_defining_ids<T>();
    dof::exists_(&ref.id, tn)
}

public(package) fun remove<T: key + store>(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    _ctx: &TxContext,
): T {
    assert!(object::id(owner_cap) == ownable::owner_cap_id(&ref.ownable_state), EInvalidOwnerCap);
    let tn = type_name::with_defining_ids<T>();
    assert!(dof::exists_(&ref.id, tn), EModuleDoesNotExist);
    dof::remove(&mut ref.id, tn)
}

public(package) fun borrow<T: key + store>(ref: &CCIPObjectRef): &T {
    let tn = type_name::with_defining_ids<T>();
    dof::borrow(&ref.id, tn)
}

public(package) fun borrow_mut<T: key + store>(ref: &mut CCIPObjectRef): &mut T {
    let tn = type_name::with_defining_ids<T>();
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

public fun accept_ownership_from_object(
    ref: &mut CCIPObjectRef,
    from: &mut UID,
    ctx: &mut TxContext,
) {
    ownable::accept_ownership_from_object(&mut ref.ownable_state, from, ctx);
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
    let publisher_wrapper = mcms_registry::create_publisher_wrapper(
        ownable::borrow_publisher(&owner_cap),
        McmsCallback {},
    );

    ownable::execute_ownership_transfer_to_mcms(
        owner_cap,
        &mut ref.ownable_state,
        registry,
        to,
        publisher_wrapper,
        McmsCallback {},
        vector[b"fee_quoter", b"rmn_remote", b"state_object", b"token_admin_registry"],
        ctx,
    );
}

public fun mcms_register_upgrade_cap(
    upgrade_cap: UpgradeCap,
    registry: &mut Registry,
    state: &mut DeployerState,
    ctx: &mut TxContext,
) {
    mcms_deployer::register_upgrade_cap(state, registry, upgrade_cap, ctx);
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

public struct McmsCallback has drop {}

/// Proof for MCMS Accept Ownership
public struct McmsAcceptOwnershipProof has drop {}

public(package) fun mcms_callback(): McmsCallback {
    McmsCallback {}
}

public fun mcms_add_package_id(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback,
        OwnerCap,
    >(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"add_package_id"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(owner_cap)],
        &mut stream,
    );
    let package_id = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);
    add_package_id(ref, owner_cap, package_id);
}

public fun mcms_remove_package_id(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback,
        OwnerCap,
    >(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"remove_package_id"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(owner_cap)],
        &mut stream,
    );
    let package_id = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);
    remove_package_id(ref, owner_cap, package_id);
}

public fun mcms_transfer_ownership(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback,
        OwnerCap,
    >(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"transfer_ownership"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(owner_cap)],
        &mut stream,
    );

    let to = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    transfer_ownership(ref, owner_cap, to, ctx);
}

public fun mcms_accept_ownership(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let data = mcms_registry::get_accept_ownership_data(
        registry,
        params,
        McmsAcceptOwnershipProof {},
    );

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addr(object::id_address(ref), &mut stream);
    bcs_stream::assert_is_consumed(&stream);

    let mcms = mcms_registry::get_multisig_address();
    ownable::mcms_accept_ownership(&mut ref.ownable_state, mcms, ctx);
}

public fun mcms_execute_ownership_transfer(
    ref: &mut CCIPObjectRef,
    registry: &mut Registry,
    deployer_state: &mut DeployerState,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback,
        OwnerCap,
    >(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"execute_ownership_transfer"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(ref), object::id_address(_owner_cap)],
        &mut stream,
    );

    let to = bcs_stream::deserialize_address(&mut stream);
    let package_address = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    let owner_cap = mcms_registry::release_cap(registry, McmsCallback {});

    if (mcms_deployer::has_upgrade_cap(deployer_state, package_address)) {
        let upgrade_cap = mcms_deployer::release_upgrade_cap(
            deployer_state,
            registry,
            McmsCallback {},
        );
        transfer::public_transfer(upgrade_cap, to);
    };

    execute_ownership_transfer(ref, owner_cap, to, ctx);
}

public fun mcms_add_allowed_modules(
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback,
        OwnerCap,
    >(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"add_allowed_modules"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addr(object::id_address(registry), &mut stream);

    let new_module_names = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_vector_u8(stream),
    );
    bcs_stream::assert_is_consumed(&stream);

    mcms_registry::add_allowed_modules(registry, McmsCallback {}, new_module_names, ctx);
}

public fun mcms_remove_allowed_modules(
    registry: &mut Registry,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_owner_cap, function, data) = mcms_registry::get_callback_params_with_caps<
        McmsCallback,
        OwnerCap,
    >(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"remove_allowed_modules"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    bcs_stream::validate_obj_addr(object::id_address(registry), &mut stream);

    let module_names = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_vector_u8(stream),
    );
    bcs_stream::assert_is_consumed(&stream);

    mcms_registry::remove_allowed_modules(registry, McmsCallback {}, module_names, ctx);
}

// ================================================================
// |                      Test Functions                          |
// ================================================================

#[test_only]
public fun test_init(ctx: &mut TxContext) {
    init(STATE_OBJECT {}, ctx);
}

#[test_only]
public fun test_create_mcms_callback(): McmsCallback {
    McmsCallback {}
}

#[test_only]
public fun pending_transfer(ref: &CCIPObjectRef): (address, address, bool) {
    let from = ownable::pending_transfer_from(&ref.ownable_state);
    let to = ownable::pending_transfer_to(&ref.ownable_state);
    let accepted = ownable::pending_transfer_accepted(&ref.ownable_state);

    (from.get_with_default(@0x0), to.get_with_default(@0x0), accepted.get_with_default(false))
}
