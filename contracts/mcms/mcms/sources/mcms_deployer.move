module mcms::mcms_deployer;

use mcms::bcs_stream;
use mcms::mcms_account::OwnerCap;
use mcms::mcms_registry::{Self, ExecutingCallbackParams, Registry};
use std::type_name;
use sui::address;
use sui::dynamic_field as df;
use sui::event;
use sui::package::{Self, UpgradeCap, UpgradeTicket, UpgradeReceipt};
use sui::table::{Self, Table};

public struct DeployerState has key {
    id: UID,
    /// Package address -> UpgradeCap
    upgrade_caps: Table<address, UpgradeCap>,
    /// UpgradeCap ID -> Package address (For reverse lookup)
    cap_to_package: Table<ID, address>,
}

/// Dynamic-field key binding a proof package's original address to its UpgradeCap ID.
public struct UpgradeCapBindingKey has copy, drop, store {
    original_package_address: address,
}

public struct UpgradeCapRegistered has copy, drop {
    prev_owner: address,
    package_address: address,
    version: u64,
    policy: u8,
}

public struct UpgradeTicketAuthorized has copy, drop {
    package_address: address,
    policy: u8,
    digest: vector<u8>,
}

public struct UpgradeReceiptCommitted has copy, drop {
    old_package_address: address,
    new_package_address: address,
    old_version: u64,
    new_version: u64,
}

const EPackageAddressNotRegistered: u64 = 1;
const EWrongProofType: u64 = 2;
const EPackageAddressMismatch: u64 = 3;
const EUpgradeCapMismatch: u64 = 4;
const EUpgradeCapAlreadyBound: u64 = 5;
const EInvalidMCMSCallbackTarget: u64 = 6;
const EInvalidMCMSCallbackModule: u64 = 7;
const EInvalidMCMSCallbackFunction: u64 = 8;

public struct MCMS_DEPLOYER has drop {}

fun init(_witness: MCMS_DEPLOYER, ctx: &mut TxContext) {
    let state = DeployerState {
        id: object::new(ctx),
        upgrade_caps: table::new(ctx),
        cap_to_package: table::new(ctx),
    };

    transfer::share_object(state);
}

fun upgrade_cap_binding_key(original_package_address: address): UpgradeCapBindingKey {
    UpgradeCapBindingKey { original_package_address }
}

fun has_upgrade_cap_binding(state: &DeployerState, original_package_address: address): bool {
    df::exists_(&state.id, upgrade_cap_binding_key(original_package_address))
}

fun resolve_upgrade_cap_address(state: &DeployerState, original_package_address: address): address {
    if (!has_upgrade_cap_binding(state, original_package_address)) {
        return original_package_address
    };

    let upgrade_cap_id =
        *df::borrow<UpgradeCapBindingKey, ID>(
            &state.id,
            upgrade_cap_binding_key(original_package_address),
        );
    *state.cap_to_package.borrow(upgrade_cap_id)
}

fun add_upgrade_cap_binding(
    state: &mut DeployerState,
    original_package_address: address,
    upgrade_cap_id: ID,
) {
    assert!(!has_upgrade_cap_binding(state, original_package_address), EUpgradeCapAlreadyBound);

    df::add(
        &mut state.id,
        upgrade_cap_binding_key(original_package_address),
        upgrade_cap_id,
    );
}

fun remove_upgrade_cap_binding(
    state: &mut DeployerState,
    original_package_address: address,
    upgrade_cap_id: ID,
) {
    let stored_upgrade_cap_id: ID = df::remove(
        &mut state.id,
        upgrade_cap_binding_key(original_package_address),
    );
    assert!(stored_upgrade_cap_id == upgrade_cap_id, EUpgradeCapMismatch);
}

fun store_upgrade_cap(
    state: &mut DeployerState,
    upgrade_cap: UpgradeCap,
    original_package_address: address,
    current_package_address: address,
    ctx: &mut TxContext,
) {
    let upgrade_cap_id = object::id(&upgrade_cap);
    let version = upgrade_cap.version();
    let policy = upgrade_cap.policy();

    add_upgrade_cap_binding(state, original_package_address, upgrade_cap_id);
    state.cap_to_package.add(upgrade_cap_id, current_package_address);
    state.upgrade_caps.add(current_package_address, upgrade_cap);

    event::emit(UpgradeCapRegistered {
        prev_owner: ctx.sender(),
        package_address: current_package_address,
        version,
        policy,
    });
}

/// `UpgradeCap` is automatically sent to the initial deployer of the package
/// This function must be called by the owner to register the `UpgradeCap` with MCMS
public fun register_upgrade_cap(
    state: &mut DeployerState,
    registry: &Registry,
    upgrade_cap: UpgradeCap,
    ctx: &mut TxContext,
) {
    let package_address = upgrade_cap.package().to_address();
    // Package must be registered with MCMS
    assert!(
        mcms_registry::is_package_registered(registry, package_address.to_ascii_string()),
        EPackageAddressNotRegistered,
    );

    // Registry entries are keyed by the original package address. Therefore this check also proves
    // that a newly registered cap has not been upgraded yet and current == original.
    store_upgrade_cap(state, upgrade_cap, package_address, package_address, ctx);
}

/// Only MCMS can authorize upgrades
/// `UpgradeTicket` is a "hot potato" which must be consumed after upgrading a package
public fun authorize_upgrade(
    _: &OwnerCap,
    state: &mut DeployerState,
    policy: u8,
    digest: vector<u8>,
    package_address: address,
    _ctx: &mut TxContext,
): UpgradeTicket {
    assert!(state.upgrade_caps.contains(package_address), EPackageAddressNotRegistered);

    let cap = state.upgrade_caps.borrow_mut(package_address);
    event::emit(UpgradeTicketAuthorized {
        package_address,
        policy,
        digest,
    });

    package::authorize_upgrade(cap, policy, digest)
}

/// Commit the upgrade by consuming the `UpgradeReceipt`
public fun commit_upgrade(
    state: &mut DeployerState,
    receipt: UpgradeReceipt,
    _ctx: &mut TxContext,
) {
    let upgrade_cap_id = receipt.cap();
    let new_package_address = receipt.package().to_address();
    let old_package_address = *state.cap_to_package.borrow(upgrade_cap_id);
    assert!(state.upgrade_caps.contains(old_package_address), EPackageAddressNotRegistered);

    let mut cap = state.upgrade_caps.remove(old_package_address);
    state.cap_to_package.remove(object::id(&cap));
    let old_version = cap.version();

    package::commit_upgrade(&mut cap, receipt);
    assert!(cap.package().to_address() == new_package_address, EPackageAddressMismatch);

    let new_version = cap.version();
    state.cap_to_package.add(object::id(&cap), new_package_address);
    state.upgrade_caps.add(new_package_address, cap);

    event::emit(UpgradeReceiptCommitted {
        old_package_address,
        new_package_address,
        old_version,
        new_version,
    });
}

/// Release the upgrade cap for a registered package
/// This must be called before calling `mcms_registry::release_cap` as it relies on registered proof types in registry
public fun release_upgrade_cap<T: drop>(
    state: &mut DeployerState,
    registry: &Registry,
    _proof: T,
): UpgradeCap {
    let proof_type = type_name::with_original_ids<T>();
    let proof_account_address = proof_type.address_string();

    assert!(
        mcms_registry::is_package_registered(registry, proof_account_address),
        EPackageAddressNotRegistered,
    );

    let expected_proof_type = mcms_registry::get_registered_proof_type(
        registry,
        proof_account_address,
    );
    assert!(proof_type == expected_proof_type, EWrongProofType);

    let original_package_address = address::from_ascii_bytes(&proof_account_address.into_bytes());
    let package_address = resolve_upgrade_cap_address(state, original_package_address);
    assert!(state.upgrade_caps.contains(package_address), EPackageAddressNotRegistered);

    let cap = state.upgrade_caps.borrow(package_address);
    let upgrade_cap_id = object::id(cap);
    assert!(cap.package().to_address() == package_address, EPackageAddressMismatch);
    assert!(
        *state.cap_to_package.borrow(upgrade_cap_id) == package_address,
        EPackageAddressMismatch,
    );

    let has_binding = has_upgrade_cap_binding(state, original_package_address);
    if (has_binding) {
        let expected_upgrade_cap_id = df::borrow<UpgradeCapBindingKey, ID>(
            &state.id,
            upgrade_cap_binding_key(original_package_address),
        );
        assert!(*expected_upgrade_cap_id == upgrade_cap_id, EUpgradeCapMismatch);
    };

    let upgrade_cap = state.upgrade_caps.remove(package_address);
    state.cap_to_package.remove(object::id(&upgrade_cap));

    if (has_binding) {
        remove_upgrade_cap_binding(
            state,
            original_package_address,
            object::id(&upgrade_cap),
        );
    };

    upgrade_cap
}

/// Bind an UpgradeCap registered before proof-to-cap bindings were introduced.
public fun bind_legacy_upgrade_cap(
    state: &mut DeployerState,
    registry: &mut Registry,
    params: ExecutingCallbackParams,
) {
    let (target, module_name, function_name, data) = mcms_registry::get_callback_params_from_mcms(
        registry,
        params,
    );

    assert!(target == mcms_registry::get_multisig_address(), EInvalidMCMSCallbackTarget);
    assert!(*module_name.as_bytes() == b"mcms_deployer", EInvalidMCMSCallbackModule);
    assert!(*function_name.as_bytes() == b"bind_legacy_upgrade_cap", EInvalidMCMSCallbackFunction);

    let stream = &mut bcs_stream::new(data);
    bcs_stream::validate_obj_addrs(
        vector[object::id_address(registry), object::id_address(state)],
        stream,
    );
    let original_package_address = bcs_stream::deserialize_address(stream);
    let current_package_address = bcs_stream::deserialize_address(stream);
    bcs_stream::assert_is_consumed(stream);

    assert!(
        mcms_registry::is_package_registered(
            registry,
            original_package_address.to_ascii_string(),
        ),
        EPackageAddressNotRegistered,
    );
    assert!(state.upgrade_caps.contains(current_package_address), EPackageAddressNotRegistered);

    let upgrade_cap = state.upgrade_caps.borrow(current_package_address);
    let upgrade_cap_id = object::id(upgrade_cap);
    assert!(upgrade_cap.package().to_address() == current_package_address, EPackageAddressMismatch);
    assert!(
        *state.cap_to_package.borrow(upgrade_cap_id) == current_package_address,
        EPackageAddressMismatch,
    );

    add_upgrade_cap_binding(state, original_package_address, upgrade_cap_id);
}

public fun has_upgrade_cap(state: &DeployerState, package_address: address): bool {
    state.upgrade_caps.contains(package_address)
}

/// Return whether the UpgradeCap bound to proof type `T` is held by this deployer.
public fun has_upgrade_cap_for<T>(state: &DeployerState): bool {
    let proof_type = type_name::with_original_ids<T>();
    let proof_account_address = proof_type.address_string();
    let original_package_address = address::from_ascii_bytes(&proof_account_address.into_bytes());
    let package_address = resolve_upgrade_cap_address(state, original_package_address);
    state.upgrade_caps.contains(package_address)
}

#[test_only]
/// Registers `upgrade_cap` under `package_address` while keeping a non-zero `upgrade_cap.package`
/// id from `package::test_publish` for unit tests when `get_multisig_address()` is `@0x0`.
public fun test_register_upgrade_cap_for_package(
    state: &mut DeployerState,
    registry: &Registry,
    upgrade_cap: UpgradeCap,
    package_address: address,
    ctx: &mut TxContext,
) {
    assert!(
        mcms_registry::is_package_registered(registry, package_address.to_ascii_string()),
        EPackageAddressNotRegistered,
    );

    store_upgrade_cap(state, upgrade_cap, package_address, package_address, ctx);
}

#[test_only]
public fun test_init(ctx: &mut TxContext) {
    init(MCMS_DEPLOYER {}, ctx);
}

#[test_only]
/// Register an upgrade cap without requiring MCMS registry check.
/// Needed because in tests @self resolves to 0x0, which Sui >= 1.73
/// rejects in `authorize_upgrade` (uses 0x0 as a sentinel for in-progress upgrades).
public fun test_register_upgrade_cap(
    state: &mut DeployerState,
    upgrade_cap: UpgradeCap,
    ctx: &mut TxContext,
) {
    let package_address = upgrade_cap.package().to_address();
    store_upgrade_cap(state, upgrade_cap, package_address, package_address, ctx);
}

#[test_only]
/// Register a test UpgradeCap under a current address while binding it to a different original
/// package address. This models an already-upgraded package without weakening production checks.
public fun test_register_upgrade_cap_with_original(
    state: &mut DeployerState,
    upgrade_cap: UpgradeCap,
    original_package_address: address,
    ctx: &mut TxContext,
) {
    let current_package_address = upgrade_cap.package().to_address();
    store_upgrade_cap(
        state,
        upgrade_cap,
        original_package_address,
        current_package_address,
        ctx,
    );
}

#[test_only]
/// Register a test UpgradeCap using the pre-binding storage layout.
public fun test_register_upgrade_cap_without_binding(
    state: &mut DeployerState,
    upgrade_cap: UpgradeCap,
    ctx: &mut TxContext,
) {
    let package_address = upgrade_cap.package().to_address();
    let version = upgrade_cap.version();
    let policy = upgrade_cap.policy();

    state.cap_to_package.add(object::id(&upgrade_cap), package_address);
    state.upgrade_caps.add(package_address, upgrade_cap);

    event::emit(UpgradeCapRegistered {
        prev_owner: ctx.sender(),
        package_address,
        version,
        policy,
    });
}
