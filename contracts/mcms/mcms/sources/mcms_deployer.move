module mcms::mcms_deployer;

use mcms::mcms_account::OwnerCap;
use mcms::mcms_registry::{Self, Registry};
use std::type_name;
use sui::address;
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

    let package_address = address::from_ascii_bytes(&proof_account_address.into_bytes());
    assert!(state.upgrade_caps.contains(package_address), EPackageAddressNotRegistered);

    let upgrade_cap = state.upgrade_caps.remove(package_address);
    state.cap_to_package.remove(object::id(&upgrade_cap));

    upgrade_cap
}

public fun has_upgrade_cap(state: &DeployerState, package_address: address): bool {
    state.upgrade_caps.contains(package_address)
}

#[test_only]
public fun test_init(ctx: &mut TxContext) {
    init(MCMS_DEPLOYER {}, ctx);
}
