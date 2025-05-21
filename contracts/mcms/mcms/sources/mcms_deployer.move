module mcms::mcms_deployer;

use mcms::mcms_account::OwnerCap;
use mcms::mcms_registry::{Self, Registry};
use sui::event;
use sui::package::{Self, UpgradeCap, UpgradeTicket, UpgradeReceipt};
use sui::table::{Self, Table};

public struct DeployerState has key {
    id: UID,
    /// Package address -> UpgradeCap
    upgrade_caps: Table<address, UpgradeCap>,
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
    package_address: address,
    old_version: u64,
    new_version: u64,
}

const EPackageAddressNotRegistered: u64 = 1;

public struct MCMS_DEPLOYER has drop {}

fun init(_witness: MCMS_DEPLOYER, ctx: &mut TxContext) {
    let state = DeployerState {
        id: object::new(ctx),
        upgrade_caps: table::new(ctx),
    };

    transfer::share_object(state);
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
        mcms_registry::is_package_registered(registry, package_address),
        EPackageAddressNotRegistered,
    );

    let version = upgrade_cap.version();
    let policy = upgrade_cap.policy();

    state.upgrade_caps.add(package_address, upgrade_cap);

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

/// Commit the upgrade by consuming the `UpgradeTicket`
public fun commit_upgrade(
    state: &mut DeployerState,
    receipt: UpgradeReceipt,
    _ctx: &mut TxContext,
) {
    let package_address = receipt.package().to_address();
    assert!(state.upgrade_caps.contains(package_address), EPackageAddressNotRegistered);

    let cap = state.upgrade_caps.borrow_mut(package_address);
    let old_version = cap.version();

    package::commit_upgrade(cap, receipt);

    event::emit(UpgradeReceiptCommitted {
        package_address,
        old_version,
        new_version: cap.version(),
    });
}

#[test_only]
public fun test_init(ctx: &mut TxContext) {
    init(MCMS_DEPLOYER {}, ctx);
}
