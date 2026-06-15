module mcms::mcms_deployer;

use mcms::mcms_account::OwnerCap;
use mcms::mcms_registry::{Self, Registry};
use std::type_name::{Self, TypeName};
use sui::address;
use sui::dynamic_field as df;
use sui::event;
use sui::package::{Self, UpgradeCap, UpgradeTicket, UpgradeReceipt};
use sui::table::{Self, Table};

public struct DeployerState has key {
    id: UID,
    /// Package address -> UpgradeCap. The key is rekeyed by `commit_upgrade`
    /// to the post-upgrade package address.
    upgrade_caps: Table<address, UpgradeCap>,
    /// UpgradeCap ID -> current package address.
    cap_to_package: Table<ID, address>,
}

/// Dynamic-field key on `DeployerState.id` providing a stable mapping from the
/// original publish address of a registered package (where its proof type
/// lives) to the corresponding `UpgradeCap` ID, regardless of how many times
/// the package has been upgraded.
public struct OriginalToCapKey has copy, drop, store { original: address }

/// Value stored under `OriginalToCapKey`. Holds the `UpgradeCap` object ID
/// (stable across upgrades) and the proof type registered for the original
/// package address (so `release_upgrade_cap` does not need to read `Registry`).
public struct OriginalToCapRecord has copy, drop, store {
    cap_id: ID,
    proof_type: TypeName,
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
const EBackfillLengthMismatch: u64 = 3;
const EBackfillCapIdNotInState: u64 = 4;
const EBackfillAlreadyMapped: u64 = 5;

public struct MCMS_DEPLOYER has drop {}

fun init(_witness: MCMS_DEPLOYER, ctx: &mut TxContext) {
    let state = DeployerState {
        id: object::new(ctx),
        upgrade_caps: table::new(ctx),
        cap_to_package: table::new(ctx),
    };

    transfer::share_object(state);
}

/// `UpgradeCap` is automatically sent to the initial deployer of the package
/// This function must be called by the owner to register the `UpgradeCap` with MCMS.
/// Records a stable `original publish address -> (cap_id, proof_type)` mapping so
/// that the cap can be released after any number of upgrades without consulting
/// the `Registry`.
public fun register_upgrade_cap(
    state: &mut DeployerState,
    registry: &Registry,
    upgrade_cap: UpgradeCap,
    ctx: &mut TxContext,
) {
    let package_address = upgrade_cap.package().to_address();
    // Package must be registered with MCMS. At registration time
    // `upgrade_cap.package()` is necessarily the original publish address
    // because no `commit_upgrade` has rekeyed it yet.
    assert!(
        mcms_registry::is_package_registered(registry, package_address.to_ascii_string()),
        EPackageAddressNotRegistered,
    );

    let version = upgrade_cap.version();
    let policy = upgrade_cap.policy();
    let cap_id = object::id(&upgrade_cap);
    let proof_type = mcms_registry::get_registered_proof_type(
        registry,
        package_address.to_ascii_string(),
    );

    state.cap_to_package.add(cap_id, package_address);
    state.upgrade_caps.add(package_address, upgrade_cap);
    df::add(
        &mut state.id,
        OriginalToCapKey { original: package_address },
        OriginalToCapRecord { cap_id, proof_type },
    );

    event::emit(UpgradeCapRegistered {
        prev_owner: ctx.sender(),
        package_address,
        version,
        policy,
    });
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
    let new_package_address = receipt.package().to_address();
    let old_package_address = *state.cap_to_package.borrow(receipt.cap());
    assert!(state.upgrade_caps.contains(old_package_address), EPackageAddressNotRegistered);

    let mut cap = state.upgrade_caps.remove(old_package_address);
    state.cap_to_package.remove(object::id(&cap));
    let old_version = cap.version();

    package::commit_upgrade(&mut cap, receipt);

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

/// Release the upgrade cap for a registered package.
///
/// Authorization is bound to the proof type `T`: the caller must construct an
/// instance of the exact drop type that was registered for the package at its
/// original publish address. The cap is resolved via the stable
/// `OriginalToCapKey -> OriginalToCapRecord` dynamic field, so this function
/// works correctly after any number of `commit_upgrade` calls and does not
/// depend on `Registry` state. The `registry` parameter is retained for
/// signature compatibility but is no longer consulted, so call ordering
/// relative to `mcms_registry::release_cap` is irrelevant.
public fun release_upgrade_cap<T: drop>(
    state: &mut DeployerState,
    _registry: &Registry,
    _proof: T,
): UpgradeCap {
    let proof_type = type_name::with_original_ids<T>();
    let original_addr = address::from_ascii_bytes(
        &proof_type.address_string().into_bytes(),
    );
    let key = OriginalToCapKey { original: original_addr };
    assert!(df::exists_(&state.id, key), EPackageAddressNotRegistered);

    let OriginalToCapRecord { cap_id, proof_type: registered_proof_type } =
        df::remove(&mut state.id, key);
    assert!(proof_type == registered_proof_type, EWrongProofType);

    let current_addr = *state.cap_to_package.borrow(cap_id);
    let upgrade_cap = state.upgrade_caps.remove(current_addr);
    state.cap_to_package.remove(cap_id);

    upgrade_cap
}

/// One-time migration entry point for `DeployerState` instances that registered
/// upgrade caps before the stable `OriginalToCapKey` indirection existed.
///
/// Pairs `(original_addrs[i], cap_ids[i])` must satisfy:
/// - `cap_ids[i]` is currently tracked by `state.cap_to_package` (i.e. the cap
///   is held by `DeployerState`).
/// - `original_addrs[i]` has no existing `OriginalToCapKey` mapping yet.
/// - `original_addrs[i]` is still registered in `Registry` so the proof type
///   can be snapshotted from `mcms_registry::get_registered_proof_type`.
///
/// Gated by `&OwnerCap` so it can only be invoked through an MCMS-authorized
/// transaction.
public fun backfill_upgrade_cap_records(
    _: &OwnerCap,
    state: &mut DeployerState,
    registry: &Registry,
    original_addrs: vector<address>,
    cap_ids: vector<ID>,
) {
    let n = original_addrs.length();
    assert!(n == cap_ids.length(), EBackfillLengthMismatch);
    let mut i = 0;
    while (i < n) {
        let original = original_addrs[i];
        let cap_id = cap_ids[i];
        assert!(state.cap_to_package.contains(cap_id), EBackfillCapIdNotInState);
        let key = OriginalToCapKey { original };
        assert!(!df::exists_(&state.id, key), EBackfillAlreadyMapped);
        let proof_type = mcms_registry::get_registered_proof_type(
            registry,
            original.to_ascii_string(),
        );
        df::add(&mut state.id, key, OriginalToCapRecord { cap_id, proof_type });
        i = i + 1;
    };
}

public fun has_upgrade_cap(state: &DeployerState, package_address: address): bool {
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

#[test_only]
public fun test_init(ctx: &mut TxContext) {
    init(MCMS_DEPLOYER {}, ctx);
}
