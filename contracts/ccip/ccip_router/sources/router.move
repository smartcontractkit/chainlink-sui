module ccip_router::router;

use ccip_router::ownable::{Self, OwnerCap, OwnableState};
use mcms::bcs_stream::{Self, BCSStream};
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, Registry, ExecutingCallbackParams};
use std::ascii;
use std::string::{Self, String};
use std::type_name;
use sui::address;
use sui::derived_object;
use sui::event;
use sui::package::{Self, UpgradeCap};
use sui::vec_map::{Self, VecMap};

public struct ROUTER has drop {}

public struct RouterObject has key {
    id: UID,
}

public struct OnRampSet has copy, drop {
    dest_chain_selector: u64,
    on_ramp_package_id: address,
}

public struct RouterState has key {
    id: UID,
    ownable_state: OwnableState,
    on_ramp_package_ids: VecMap<u64, address>, // dest_chain_selector -> on_ramp_package_id
}

public struct RouterStatePointer has key, store {
    id: UID,
    router_object_id: address,
}

const EParamsLengthMismatch: u64 = 1;
const EOnrampNotFound: u64 = 2;
const EInvalidOwnerCap: u64 = 3;
const EInvalidFunction: u64 = 4;
const EInvalidObjectAddress: u64 = 5;
const EInvalidOnrampAddress: u64 = 6;

fun init(otw: ROUTER, ctx: &mut TxContext) {
    let mut router_object = RouterObject { id: object::new(ctx) };
    let (ownable_state, mut owner_cap) = ownable::new(&mut router_object.id, ctx);

    let router = RouterState {
        id: derived_object::claim(&mut router_object.id, b"RouterState"),
        ownable_state,
        on_ramp_package_ids: vec_map::empty(),
    };

    let router_state_pointer = RouterStatePointer {
        id: object::new(ctx),
        router_object_id: object::id_address(&router_object),
    };

    let tn = type_name::with_original_ids<ROUTER>();
    let package_bytes = ascii::into_bytes(tn.address_string());
    let package_id = address::from_ascii_bytes(&package_bytes);

    transfer::share_object(router);
    transfer::share_object(router_object);

    let publisher = package::claim(otw, ctx);
    ownable::attach_publisher(&mut owner_cap, publisher);

    transfer::public_transfer(owner_cap, ctx.sender());
    transfer::transfer(router_state_pointer, package_id);
}

public(package) fun get_uid(router_object: &mut RouterObject): &mut UID {
    &mut router_object.id
}

public fun type_and_version(): String {
    string::utf8(b"Router 1.6.0")
}

public fun is_chain_supported(router: &RouterState, dest_chain_selector: u64): bool {
    router.on_ramp_package_ids.contains(&dest_chain_selector)
}

// Returns the on ramp package id for the given destination chain selector.
public fun get_on_ramp(router: &RouterState, dest_chain_selector: u64): address {
    assert!(router.on_ramp_package_ids.contains(&dest_chain_selector), EOnrampNotFound);

    *router.on_ramp_package_ids.get(&dest_chain_selector)
}

public fun get_dest_chains(router: &RouterState): vector<u64> {
    router.on_ramp_package_ids.keys()
}

/// Sets the onramp package ids for the given destination chains.
/// This function will overwrite the existing package ids.
/// This function can only be called by the owner of the contract.
/// @param owner_cap The owner capability.
/// @param router The router state.
/// @param dest_chain_selectors The destination chain selectors.
/// @param on_ramp_package_ids The onramp package ids.
public fun set_on_ramps(
    owner_cap: &OwnerCap,
    router: &mut RouterState,
    dest_chain_selectors: vector<u64>,
    on_ramp_package_ids: vector<address>,
) {
    assert!(
        object::id(owner_cap) == ownable::owner_cap_id(&router.ownable_state),
        EInvalidOwnerCap,
    );
    assert!(dest_chain_selectors.length() == on_ramp_package_ids.length(), EParamsLengthMismatch);

    let mut i = 0;
    let selector_len = dest_chain_selectors.length();
    while (i < selector_len) {
        let dest_chain_selector = dest_chain_selectors[i];
        let on_ramp_package_id = on_ramp_package_ids[i];
        assert!(on_ramp_package_id != @0x0, EInvalidOnrampAddress);

        if (router.on_ramp_package_ids.contains(&dest_chain_selector)) {
            router.on_ramp_package_ids.remove(&dest_chain_selector);
        };
        router.on_ramp_package_ids.insert(dest_chain_selector, on_ramp_package_id);
        event::emit(OnRampSet { dest_chain_selector, on_ramp_package_id });
        i = i + 1;
    };
}

// ================================================================
// |                      Ownable Functions                       |
// ================================================================

public fun owner(state: &RouterState): address {
    ownable::owner(&state.ownable_state)
}

public fun has_pending_transfer(state: &RouterState): bool {
    ownable::has_pending_transfer(&state.ownable_state)
}

public fun pending_transfer_from(state: &RouterState): Option<address> {
    ownable::pending_transfer_from(&state.ownable_state)
}

public fun pending_transfer_to(state: &RouterState): Option<address> {
    ownable::pending_transfer_to(&state.ownable_state)
}

public fun pending_transfer_accepted(state: &RouterState): Option<bool> {
    ownable::pending_transfer_accepted(&state.ownable_state)
}

public fun transfer_ownership(
    state: &mut RouterState,
    owner_cap: &OwnerCap,
    new_owner: address,
    ctx: &mut TxContext,
) {
    ownable::transfer_ownership(owner_cap, &mut state.ownable_state, new_owner, ctx);
}

public fun accept_ownership(state: &mut RouterState, ctx: &mut TxContext) {
    ownable::accept_ownership(&mut state.ownable_state, ctx);
}

public fun accept_ownership_from_object(
    state: &mut RouterState,
    from: &mut UID,
    ctx: &mut TxContext,
) {
    ownable::accept_ownership_from_object(&mut state.ownable_state, from, ctx);
}

public fun mcms_accept_ownership(
    state: &mut RouterState,
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
    bcs_stream::validate_obj_addr(object::id_address(state), &mut stream);
    bcs_stream::assert_is_consumed(&stream);

    let mcms = mcms_registry::get_multisig_address();
    ownable::mcms_accept_ownership(&mut state.ownable_state, mcms, ctx);
}

public fun execute_ownership_transfer(
    owner_cap: OwnerCap,
    state: &mut RouterState,
    to: address,
    ctx: &mut TxContext,
) {
    ownable::execute_ownership_transfer(owner_cap, &mut state.ownable_state, to, ctx);
}

public fun execute_ownership_transfer_to_mcms(
    owner_cap: OwnerCap,
    state: &mut RouterState,
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
        &mut state.ownable_state,
        registry,
        to,
        publisher_wrapper,
        McmsCallback {},
        vector[b"router"],
        ctx,
    );
}

public fun mcms_register_upgrade_cap(
    upgrade_cap: UpgradeCap,
    registry: &mut Registry,
    state: &mut DeployerState,
    ctx: &mut TxContext,
) {
    mcms_deployer::register_upgrade_cap(
        state,
        registry,
        upgrade_cap,
        ctx,
    );
}

// ================================================================
// |                      MCMS Entrypoint                         |
// ================================================================

public struct McmsCallback has drop {}

/// Proof for MCMS Accept Ownership
public struct McmsAcceptOwnershipProof has drop {}

public fun mcms_set_on_ramps(
    state: &mut RouterState,
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
    assert!(function == string::utf8(b"set_on_ramps"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    validate_obj_addrs(
        vector[object::id_address(owner_cap), object::id_address(state)],
        &mut stream,
    );

    let dest_chain_selectors = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_u64(stream),
    );
    let on_ramp_package_ids = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_address(stream),
    );
    bcs_stream::assert_is_consumed(&stream);

    set_on_ramps(
        owner_cap,
        state,
        dest_chain_selectors,
        on_ramp_package_ids,
    );
}

public fun mcms_transfer_ownership(
    state: &mut RouterState,
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
    validate_obj_addrs(
        vector[object::id_address(state), object::id_address(owner_cap)],
        &mut stream,
    );

    let to = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    transfer_ownership(state, owner_cap, to, ctx);
}

public fun mcms_execute_ownership_transfer(
    state: &mut RouterState,
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
    validate_obj_addrs(
        vector[object::id_address(_owner_cap), object::id_address(state)],
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
            McmsCallback {}
        );
        transfer::public_transfer(upgrade_cap, to);
    };

    execute_ownership_transfer(owner_cap, state, to, ctx);
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

fun validate_obj_addrs(addrs: vector<address>, stream: &mut BCSStream) {
    let mut i = 0;
    while (i < addrs.length()) {
        let deserialized_address = bcs_stream::deserialize_address(stream);
        assert!(deserialized_address == addrs[i], EInvalidObjectAddress);
        i = i + 1;
    }
}

// ===================== TESTS =====================

#[test_only]
public fun test_init(ctx: &mut TxContext) {
    init(ROUTER {}, ctx);
}
