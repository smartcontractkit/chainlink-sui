module ccip_router::router;

use ccip_router::ownable::{Self, OwnerCap, OwnableState};
use mcms::bcs_stream::{Self, BCSStream};
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, Registry, ExecutingCallbackParams};
use std::string::{Self, String};
use sui::event;
use sui::package::UpgradeCap;
use sui::table::{Self, Table};

public struct ROUTER has drop {}

public struct OnRampSet has copy, drop {
    dest_chain_selector: u64,
    on_ramp: address,
}

public struct RouterState has key {
    id: UID,
    ownable_state: OwnableState,
    on_ramps: Table<u64, address>,
}

const EParamsLengthMismatch: u64 = 1;
const EOnrampNotFound: u64 = 2;
const EInvalidOwnerCap: u64 = 3;
const EInvalidFunction: u64 = 4;
const EInvalidStateAddress: u64 = 5;
const EInvalidObjectAddress: u64 = 6;
const EInvalidOnrampAddress: u64 = 7;

fun init(_witness: ROUTER, ctx: &mut TxContext) {
    let (ownable_state, owner_cap) = ownable::new(ctx);

    let router = RouterState {
        id: object::new(ctx),
        ownable_state,
        on_ramps: table::new(ctx),
    };

    transfer::share_object(router);
    transfer::public_transfer(owner_cap, ctx.sender());
}

public fun type_and_version(): String {
    string::utf8(b"Router 1.6.0")
}

public fun is_chain_supported(router: &RouterState, dest_chain_selector: u64): bool {
    router.on_ramps.contains(dest_chain_selector)
}

public fun get_on_ramp(router: &RouterState, dest_chain_selector: u64): address {
    assert!(router.on_ramps.contains(dest_chain_selector), EOnrampNotFound);

    *router.on_ramps.borrow(dest_chain_selector)
}

/// Sets the onRamp info for the given destination chains.
/// This function will overwrite the existing infos.
/// This function can only be called by the owner of the contract.
/// @param owner_cap The owner capability.
/// @param router The router state.
/// @param dest_chain_selectors The destination chain selectors.
/// @param on_ramp_addresses The onRamp addresses.
public fun set_on_ramps(
    owner_cap: &OwnerCap,
    router: &mut RouterState,
    dest_chain_selectors: vector<u64>,
    on_ramp_addresses: vector<address>,
) {
    assert!(
        object::id(owner_cap) == ownable::owner_cap_id(&router.ownable_state),
        EInvalidOwnerCap,
    );
    assert!(dest_chain_selectors.length() == on_ramp_addresses.length(), EParamsLengthMismatch);

    let mut i = 0;
    let selector_len = dest_chain_selectors.length();
    while (i < selector_len) {
        let dest_chain_selector = dest_chain_selectors[i];
        let on_ramp = on_ramp_addresses[i];
        assert!(on_ramp != @0x0, EInvalidOnrampAddress);

        if (router.on_ramps.contains(dest_chain_selector)) {
            router.on_ramps.remove(dest_chain_selector);
        };
        router.on_ramps.add(dest_chain_selector, on_ramp);
        event::emit(OnRampSet { dest_chain_selector, on_ramp });
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
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_, _, function, data) = mcms_registry::get_callback_params_for_mcms(
        params,
        McmsCallback {},
    );
    assert!(function == string::utf8(b"accept_ownership"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    let state_address = bcs_stream::deserialize_address(&mut stream);
    assert!(state_address == object::id_address(state), EInvalidStateAddress);

    let mcms = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

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
    ownable::execute_ownership_transfer_to_mcms(
        owner_cap,
        &mut state.ownable_state,
        registry,
        to,
        McmsCallback {},
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

public fun mcms_set_on_ramps(
    state: &mut RouterState,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback, OwnerCap>(
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
    let on_ramps = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_address(stream),
    );
    bcs_stream::assert_is_consumed(&stream);

    set_on_ramps(
        owner_cap,
        state,
        dest_chain_selectors,
        on_ramps,
    );
}

public fun mcms_transfer_ownership(
    state: &mut RouterState,
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
    validate_obj_addrs(
        vector[object::id_address(_owner_cap), object::id_address(state)],
        &mut stream,
    );

    let to = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    let owner_cap = mcms_registry::release_cap(registry, McmsCallback {});
    execute_ownership_transfer(owner_cap, state, to, ctx);
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
