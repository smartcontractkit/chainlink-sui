#[test_only]
module mcms::mcms_deployer_test;

use mcms::mcms_account::{Self, OwnerCap};
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, Registry};
use sui::package::{Self, UpgradeCap};
use sui::test_scenario::{Self as ts, Scenario};

public struct MCMS_DEPLOYER_TEST has drop {}

public struct TestOwnerCap has key, store {
    id: UID,
}

const ADMIN: address = @0xA;
const OTHER: address = @0xB;

fun create_test_scenario(): Scenario {
    ts::begin(ADMIN)
}

fun generate_upgrade_cap(ctx: &mut TxContext): UpgradeCap {
    package::test_publish(mcms_registry::get_multisig_address().to_id(), ctx)
}

/// Boots MCMS account/registry/deployer state and registers `MCMS_DEPLOYER_TEST`
/// as a proof type for `@mcms` with `TestOwnerCap` as its package cap.
fun init_mcms(scenario: &mut Scenario) {
    let ctx_init = ts::ctx(scenario);
    mcms_registry::test_init(ctx_init);
    mcms_deployer::test_init(ctx_init);
    mcms_account::test_init(ctx_init);

    ts::next_tx(scenario, ADMIN);
    let mut registry = ts::take_shared<Registry>(scenario);
    let ctx = ts::ctx(scenario);

    let publisher = package::test_claim(MCMS_DEPLOYER_TEST {}, ctx);
    let publisher_wrapper = mcms_registry::create_publisher_wrapper(
        &publisher,
        MCMS_DEPLOYER_TEST {},
    );
    mcms_registry::register_entrypoint<MCMS_DEPLOYER_TEST, TestOwnerCap>(
        &mut registry,
        publisher_wrapper,
        MCMS_DEPLOYER_TEST {},
        TestOwnerCap { id: object::new(ctx) },
        vector[b"mcms_deployer_test"],
        ctx,
    );
    transfer::public_transfer(publisher, ADMIN);
    ts::return_shared(registry);
}

/// Registers a fresh `UpgradeCap` through the production `register_upgrade_cap`
/// path. Returns `(package_address, cap_id)`.
fun register_cap_through_production(scenario: &mut Scenario): (address, ID) {
    ts::next_tx(scenario, OTHER);
    let mut deployer_state = ts::take_shared<DeployerState>(scenario);
    let registry = ts::take_shared<Registry>(scenario);
    let ctx = ts::ctx(scenario);

    let upgrade_cap = generate_upgrade_cap(ctx);
    let package_address = upgrade_cap.package().to_address();
    let cap_id = object::id(&upgrade_cap);

    mcms_deployer::register_upgrade_cap(
        &mut deployer_state,
        &registry,
        upgrade_cap,
        ctx,
    );

    ts::return_shared(deployer_state);
    ts::return_shared(registry);

    (package_address, cap_id)
}

/// Simulates legacy state: registers an `UpgradeCap` directly in
/// `upgrade_caps`/`cap_to_package` without writing the new `OriginalToCapKey`
/// dynamic field. Mirrors `DeployerState` instances that were populated before
/// the indirection existed.
fun register_cap_legacy(scenario: &mut Scenario): (address, ID) {
    ts::next_tx(scenario, OTHER);
    let mut deployer_state = ts::take_shared<DeployerState>(scenario);
    let registry = ts::take_shared<Registry>(scenario);
    let ctx = ts::ctx(scenario);

    let upgrade_cap = generate_upgrade_cap(ctx);
    let package_address = upgrade_cap.package().to_address();
    let cap_id = object::id(&upgrade_cap);

    mcms_deployer::test_register_upgrade_cap_for_package(
        &mut deployer_state,
        &registry,
        upgrade_cap,
        package_address,
        ctx,
    );

    ts::return_shared(deployer_state);
    ts::return_shared(registry);

    (package_address, cap_id)
}

/// Performs one `authorize_upgrade` + `test_upgrade` + `commit_upgrade` cycle
/// and returns the new package address that the cap is now keyed under.
fun perform_one_upgrade(scenario: &mut Scenario, current_address: address): address {
    ts::next_tx(scenario, ADMIN);
    let mut deployer_state = ts::take_shared<DeployerState>(scenario);
    let owner_cap = ts::take_from_sender<OwnerCap>(scenario);
    let ctx = ts::ctx(scenario);

    let ticket = mcms_deployer::authorize_upgrade(
        &owner_cap,
        &mut deployer_state,
        0,
        vector[],
        current_address,
        ctx,
    );
    let receipt = package::test_upgrade(ticket);
    let new_address = receipt.package().to_address();
    mcms_deployer::commit_upgrade(&mut deployer_state, receipt, ctx);

    ts::return_to_sender(scenario, owner_cap);
    ts::return_shared(deployer_state);

    new_address
}

#[test]
fun test_register_upgrade_cap() {
    let mut scenario = create_test_scenario();

    init_mcms(&mut scenario);
    let (_, _) = register_cap_through_production(&mut scenario);

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms::mcms_deployer::EPackageAddressNotRegistered)]
fun test_register_upgrade_cap_without_existing_package_fails() {
    let mut scenario = create_test_scenario();

    {
        let ctx = ts::ctx(&mut scenario);
        mcms_registry::test_init(ctx);
        mcms_deployer::test_init(ctx);
    };

    {
        ts::next_tx(&mut scenario, OTHER);
        let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
        let registry = ts::take_shared<Registry>(&scenario);
        let ctx = ts::ctx(&mut scenario);

        let upgrade_cap = generate_upgrade_cap(ctx);

        mcms_deployer::register_upgrade_cap(
            &mut deployer_state,
            &registry,
            upgrade_cap,
            ctx,
        );

        ts::return_shared(deployer_state);
        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
fun test_release_upgrade_cap_succeeds() {
    let mut scenario = create_test_scenario();

    init_mcms(&mut scenario);
    let (_, cap_id) = register_cap_through_production(&mut scenario);

    ts::next_tx(&mut scenario, OTHER);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let registry = ts::take_shared<Registry>(&scenario);

    let upgrade_cap = mcms_deployer::release_upgrade_cap(
        &mut deployer_state,
        &registry,
        MCMS_DEPLOYER_TEST {},
    );
    assert!(object::id(&upgrade_cap) == cap_id);

    transfer::public_transfer(upgrade_cap, OTHER);
    ts::return_shared(deployer_state);
    ts::return_shared(registry);

    ts::end(scenario);
}

#[test]
fun test_release_upgrade_cap_after_upgrade_succeeds() {
    let mut scenario = create_test_scenario();

    init_mcms(&mut scenario);
    let (original_address, cap_id) = register_cap_through_production(&mut scenario);

    let new_address = perform_one_upgrade(&mut scenario, original_address);
    assert!(new_address != original_address);

    ts::next_tx(&mut scenario, OTHER);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let registry = ts::take_shared<Registry>(&scenario);

    assert!(!mcms_deployer::has_upgrade_cap(&deployer_state, original_address));
    assert!(mcms_deployer::has_upgrade_cap(&deployer_state, new_address));

    let upgrade_cap = mcms_deployer::release_upgrade_cap(
        &mut deployer_state,
        &registry,
        MCMS_DEPLOYER_TEST {},
    );
    assert!(object::id(&upgrade_cap) == cap_id);
    assert!(upgrade_cap.package().to_address() == new_address);

    transfer::public_transfer(upgrade_cap, OTHER);
    ts::return_shared(deployer_state);
    ts::return_shared(registry);

    ts::end(scenario);
}

#[test]
fun test_release_upgrade_cap_after_multiple_upgrades_succeeds() {
    let mut scenario = create_test_scenario();

    init_mcms(&mut scenario);
    let (original_address, cap_id) = register_cap_through_production(&mut scenario);

    let mut current = original_address;
    let mut i: u64 = 0;
    while (i < 3) {
        current = perform_one_upgrade(&mut scenario, current);
        i = i + 1;
    };

    ts::next_tx(&mut scenario, OTHER);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let registry = ts::take_shared<Registry>(&scenario);

    let upgrade_cap = mcms_deployer::release_upgrade_cap(
        &mut deployer_state,
        &registry,
        MCMS_DEPLOYER_TEST {},
    );
    assert!(object::id(&upgrade_cap) == cap_id);
    assert!(upgrade_cap.package().to_address() == current);

    transfer::public_transfer(upgrade_cap, OTHER);
    ts::return_shared(deployer_state);
    ts::return_shared(registry);

    ts::end(scenario);
}

#[test]
/// Regression for report 79509 Bug 1. `release_upgrade_cap` no longer reads
/// `Registry`, so calling `mcms_registry::release_cap` first must not block it.
fun test_release_upgrade_cap_after_release_cap_succeeds() {
    let mut scenario = create_test_scenario();

    init_mcms(&mut scenario);
    let (_, cap_id) = register_cap_through_production(&mut scenario);

    ts::next_tx(&mut scenario, OTHER);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let mut registry = ts::take_shared<Registry>(&scenario);

    let owner_cap: TestOwnerCap = mcms_registry::release_cap(
        &mut registry,
        MCMS_DEPLOYER_TEST {},
    );

    let upgrade_cap = mcms_deployer::release_upgrade_cap(
        &mut deployer_state,
        &registry,
        MCMS_DEPLOYER_TEST {},
    );
    assert!(object::id(&upgrade_cap) == cap_id);

    let TestOwnerCap { id } = owner_cap;
    id.delete();
    transfer::public_transfer(upgrade_cap, OTHER);
    ts::return_shared(deployer_state);
    ts::return_shared(registry);

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms::mcms_deployer::EPackageAddressNotRegistered)]
fun test_release_upgrade_cap_without_registration_fails() {
    let mut scenario = create_test_scenario();

    init_mcms(&mut scenario);

    ts::next_tx(&mut scenario, OTHER);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let registry = ts::take_shared<Registry>(&scenario);

    let upgrade_cap = mcms_deployer::release_upgrade_cap(
        &mut deployer_state,
        &registry,
        MCMS_DEPLOYER_TEST {},
    );

    transfer::public_transfer(upgrade_cap, OTHER);
    ts::return_shared(deployer_state);
    ts::return_shared(registry);

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms::mcms_deployer::EPackageAddressNotRegistered)]
/// Legacy-state caps (no `OriginalToCapKey` dynamic field) cannot be released
/// until `backfill_upgrade_cap_records` populates the indirection.
fun test_release_upgrade_cap_on_legacy_state_fails() {
    let mut scenario = create_test_scenario();

    init_mcms(&mut scenario);
    let (_, _) = register_cap_legacy(&mut scenario);

    ts::next_tx(&mut scenario, OTHER);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let registry = ts::take_shared<Registry>(&scenario);

    let upgrade_cap = mcms_deployer::release_upgrade_cap(
        &mut deployer_state,
        &registry,
        MCMS_DEPLOYER_TEST {},
    );

    transfer::public_transfer(upgrade_cap, OTHER);
    ts::return_shared(deployer_state);
    ts::return_shared(registry);

    ts::end(scenario);
}

#[test]
fun test_backfill_then_release_succeeds() {
    let mut scenario = create_test_scenario();

    init_mcms(&mut scenario);
    let (original_address, cap_id) = register_cap_legacy(&mut scenario);

    ts::next_tx(&mut scenario, ADMIN);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let registry = ts::take_shared<Registry>(&scenario);
    let owner_cap = ts::take_from_sender<OwnerCap>(&scenario);

    mcms_deployer::backfill_upgrade_cap_records(
        &owner_cap,
        &mut deployer_state,
        &registry,
        vector[original_address],
        vector[cap_id],
    );

    ts::return_to_sender(&scenario, owner_cap);
    ts::return_shared(deployer_state);
    ts::return_shared(registry);

    ts::next_tx(&mut scenario, OTHER);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let registry = ts::take_shared<Registry>(&scenario);

    let upgrade_cap = mcms_deployer::release_upgrade_cap(
        &mut deployer_state,
        &registry,
        MCMS_DEPLOYER_TEST {},
    );
    assert!(object::id(&upgrade_cap) == cap_id);

    transfer::public_transfer(upgrade_cap, OTHER);
    ts::return_shared(deployer_state);
    ts::return_shared(registry);

    ts::end(scenario);
}

#[test]
fun test_backfill_then_release_after_upgrade_succeeds() {
    let mut scenario = create_test_scenario();

    init_mcms(&mut scenario);
    let (original_address, cap_id) = register_cap_legacy(&mut scenario);

    let new_address = perform_one_upgrade(&mut scenario, original_address);

    ts::next_tx(&mut scenario, ADMIN);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let registry = ts::take_shared<Registry>(&scenario);
    let owner_cap = ts::take_from_sender<OwnerCap>(&scenario);

    mcms_deployer::backfill_upgrade_cap_records(
        &owner_cap,
        &mut deployer_state,
        &registry,
        vector[original_address],
        vector[cap_id],
    );

    ts::return_to_sender(&scenario, owner_cap);
    ts::return_shared(deployer_state);
    ts::return_shared(registry);

    ts::next_tx(&mut scenario, OTHER);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let registry = ts::take_shared<Registry>(&scenario);

    let upgrade_cap = mcms_deployer::release_upgrade_cap(
        &mut deployer_state,
        &registry,
        MCMS_DEPLOYER_TEST {},
    );
    assert!(object::id(&upgrade_cap) == cap_id);
    assert!(upgrade_cap.package().to_address() == new_address);

    transfer::public_transfer(upgrade_cap, OTHER);
    ts::return_shared(deployer_state);
    ts::return_shared(registry);

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms::mcms_deployer::EBackfillLengthMismatch)]
fun test_backfill_length_mismatch_fails() {
    let mut scenario = create_test_scenario();

    init_mcms(&mut scenario);
    let (original_address, _) = register_cap_legacy(&mut scenario);

    ts::next_tx(&mut scenario, ADMIN);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let registry = ts::take_shared<Registry>(&scenario);
    let owner_cap = ts::take_from_sender<OwnerCap>(&scenario);

    mcms_deployer::backfill_upgrade_cap_records(
        &owner_cap,
        &mut deployer_state,
        &registry,
        vector[original_address],
        vector[],
    );

    ts::return_to_sender(&scenario, owner_cap);
    ts::return_shared(deployer_state);
    ts::return_shared(registry);

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms::mcms_deployer::EBackfillCapIdNotInState)]
fun test_backfill_unknown_cap_id_fails() {
    let mut scenario = create_test_scenario();

    init_mcms(&mut scenario);
    let (original_address, _) = register_cap_legacy(&mut scenario);

    ts::next_tx(&mut scenario, ADMIN);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let registry = ts::take_shared<Registry>(&scenario);
    let owner_cap = ts::take_from_sender<OwnerCap>(&scenario);

    let bogus_id = object::id_from_address(@0xDEADBEEF);

    mcms_deployer::backfill_upgrade_cap_records(
        &owner_cap,
        &mut deployer_state,
        &registry,
        vector[original_address],
        vector[bogus_id],
    );

    ts::return_to_sender(&scenario, owner_cap);
    ts::return_shared(deployer_state);
    ts::return_shared(registry);

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms::mcms_deployer::EBackfillAlreadyMapped)]
fun test_backfill_already_mapped_fails() {
    let mut scenario = create_test_scenario();

    init_mcms(&mut scenario);
    // Production path already writes the indirection.
    let (original_address, cap_id) = register_cap_through_production(&mut scenario);

    ts::next_tx(&mut scenario, ADMIN);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let registry = ts::take_shared<Registry>(&scenario);
    let owner_cap = ts::take_from_sender<OwnerCap>(&scenario);

    mcms_deployer::backfill_upgrade_cap_records(
        &owner_cap,
        &mut deployer_state,
        &registry,
        vector[original_address],
        vector[cap_id],
    );

    ts::return_to_sender(&scenario, owner_cap);
    ts::return_shared(deployer_state);
    ts::return_shared(registry);

    ts::end(scenario);
}
