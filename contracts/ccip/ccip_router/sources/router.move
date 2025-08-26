module ccip_router::router;

use ccip_router::ownable::{Self, OwnerCap, OwnableState};
use mcms::bcs_stream;
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, Registry, ExecutingCallbackParams};
use std::string::{Self, String};
use sui::event;
use sui::package::UpgradeCap;
use sui::table::{Self, Table};

public struct ROUTER has drop {}

public struct OnRampSet has copy, drop {
    dest_chain_selector: u64,
    on_ramp_info: OnRampInfo,
}

public struct OnRampInfo has copy, drop, store {
    onramp_address: address,
    onramp_version: vector<u8>,
}

public struct RouterState has key {
    id: UID,
    ownable_state: OwnableState,
    on_ramp_infos: Table<u64, OnRampInfo>,
}

const EParamsLengthMismatch: u64 = 1;
const EOnrampInfoNotFound: u64 = 2;
const EInvalidOnrampVersion: u64 = 3;
const EInvalidOwnerCap: u64 = 4;
const EInvalidFunction: u64 = 5;
const EInvalidStateAddress: u64 = 6;
const EInvalidRegistryAddress: u64 = 7;

fun init(_witness: ROUTER, ctx: &mut TxContext) {
    let (ownable_state, owner_cap) = ownable::new(ctx);

    let router = RouterState {
        id: object::new(ctx),
        ownable_state,
        on_ramp_infos: table::new(ctx),
    };

    transfer::share_object(router);
    transfer::public_transfer(owner_cap, ctx.sender());
}

public fun type_and_version(): String {
    string::utf8(b"Router 1.6.0")
}

public fun is_chain_supported(router: &RouterState, dest_chain_selector: u64): bool {
    router.on_ramp_infos.contains(dest_chain_selector)
}

public fun get_on_ramp_info(router: &RouterState, dest_chain_selector: u64): (address, vector<u8>) {
    assert!(router.on_ramp_infos.contains(dest_chain_selector), EOnrampInfoNotFound);

    let on_ramp_info = *router.on_ramp_infos.borrow(dest_chain_selector);

    (on_ramp_info.onramp_address, on_ramp_info.onramp_version)
}

/// Returns the onRamp versions for the given destination chains.
public fun get_on_ramp_infos(
    router: &RouterState,
    dest_chain_selectors: vector<u64>,
): vector<OnRampInfo> {
    dest_chain_selectors.map!(|dest_chain_selector| {
        if (router.on_ramp_infos.contains(dest_chain_selector)) {
            *router.on_ramp_infos.borrow(dest_chain_selector)
        } else {
            OnRampInfo {
                onramp_address: @0x0,
                onramp_version: vector[],
            }
        }
    })
}

public fun get_on_ramp_version(info: OnRampInfo): vector<u8> {
    info.onramp_version
}

public fun get_on_ramp_address(info: OnRampInfo): address {
    info.onramp_address
}

/// Sets the onRamp info for the given destination chains.
/// This function will overwrite the existing infos.
/// This function can only be called by the owner of the contract.
/// @param owner_cap The owner capability.
/// @param router The router state.
/// @param dest_chain_selectors The destination chain selectors.
/// @param on_ramp_addresses The onRamp addresses.
/// @param on_ramp_versions The onRamp versions, the inner vector must be of length 0 or 3. 0 indicates
/// the destination chain is no longer supported. Length 3 encodes the version of the onRamp contract.
public fun set_on_ramp_infos(
    owner_cap: &OwnerCap,
    router: &mut RouterState,
    dest_chain_selectors: vector<u64>,
    on_ramp_addresses: vector<address>,
    on_ramp_versions: vector<vector<u8>>,
) {
    assert!(
        object::id(owner_cap) == ownable::owner_cap_id(&router.ownable_state),
        EInvalidOwnerCap,
    );
    assert!(dest_chain_selectors.length() == on_ramp_addresses.length(), EParamsLengthMismatch);
    assert!(dest_chain_selectors.length() == on_ramp_versions.length(), EParamsLengthMismatch);

    let mut i = 0;
    let selector_len = dest_chain_selectors.length();
    while (i < selector_len) {
        let dest_chain_selector = dest_chain_selectors[i];
        let version = on_ramp_versions[i];

        if (version.length() == 0) {
            if (router.on_ramp_infos.contains(dest_chain_selector)) {
                router.on_ramp_infos.remove(dest_chain_selector);
            };
            event::emit(OnRampSet {
                dest_chain_selector,
                on_ramp_info: OnRampInfo {
                    onramp_address: @0x0,
                    onramp_version: vector[],
                },
            });
        } else {
            assert!(version.length() == 3, EInvalidOnrampVersion);
            if (router.on_ramp_infos.contains(dest_chain_selector)) {
                router.on_ramp_infos.remove(dest_chain_selector);
            };

            let info = OnRampInfo {
                onramp_address: on_ramp_addresses[i],
                onramp_version: on_ramp_versions[i],
            };
            router.on_ramp_infos.add(dest_chain_selector, info);

            event::emit(OnRampSet { dest_chain_selector, on_ramp_info: info });
        };
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

/// Cannot call through `mcms_entrypoint` as owner cap is not registered with MCMS registry
public fun accept_ownership_as_mcms(
    state: &mut RouterState,
    params: ExecutingCallbackParams,
    ctx: &mut TxContext,
) {
    let (_, _, function_name, data) = mcms_registry::get_callback_params_for_mcms(
        params,
        McmsCallback {},
    );
    assert!(function_name == string::utf8(b"accept_ownership_as_mcms"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    let mcms = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    ownable::accept_ownership_as_mcms(&mut state.ownable_state, mcms, ctx);
}

public fun execute_ownership_transfer(
    owner_cap: OwnerCap,
    ownable_state: &mut OwnableState,
    to: address,
    ctx: &mut TxContext,
) {
    ownable::execute_ownership_transfer(owner_cap, ownable_state, to, ctx);
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

fun validate_shared_objects(
    state: &RouterState,
    registry: &Registry,
    stream: &mut bcs_stream::BCSStream,
) {
    let state_address = bcs_stream::deserialize_address(stream);
    assert!(state_address == object::id_address(state), EInvalidStateAddress);
    let registry_address = bcs_stream::deserialize_address(stream);
    assert!(registry_address == object::id_address(registry), EInvalidRegistryAddress);
}

public fun mcms_set_on_ramp_infos(
    state: &mut RouterState,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (owner_cap, function, data) = mcms_registry::get_callback_params<McmsCallback, OwnerCap>(
        registry,
        McmsCallback {},
        params,
    );
    assert!(function == string::utf8(b"set_on_ramp_infos"), EInvalidFunction);

    let mut stream = bcs_stream::new(data);
    validate_shared_objects(state, registry, &mut stream);

    let dest_chain_selectors = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_u64(stream),
    );
    let on_ramp_addresses = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_address(stream),
    );
    let on_ramp_versions = bcs_stream::deserialize_vector!(
        &mut stream,
        |stream| bcs_stream::deserialize_vector!(
            stream,
            |stream| bcs_stream::deserialize_u8(stream),
        ),
    );
    bcs_stream::assert_is_consumed(&stream);

    set_on_ramp_infos(
        owner_cap,
        state,
        dest_chain_selectors,
        on_ramp_addresses,
        on_ramp_versions,
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
    validate_shared_objects(state, registry, &mut stream);

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
    validate_shared_objects(state, registry, &mut stream);

    let to = bcs_stream::deserialize_address(&mut stream);
    bcs_stream::assert_is_consumed(&stream);

    let owner_cap = mcms_registry::release_cap(registry, McmsCallback {});
    execute_ownership_transfer(owner_cap, &mut state.ownable_state, to, ctx);
}

// ===================== TESTS =====================

#[test_only]
public fun test_init(ctx: &mut TxContext) {
    init(ROUTER {}, ctx);
}
