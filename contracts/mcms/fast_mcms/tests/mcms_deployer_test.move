#[test_only]
module mcms::mcms_deployer_test;

use mcms::mcms_account;
use mcms::mcms_deployer::{Self, DeployerState};
use mcms::mcms_registry::{Self, Registry};
use sui::package::{Self, UpgradeCap};
use sui::test_scenario::{Self as ts, Scenario};

public struct MCMS_DEPLOYER_TEST has drop {}

public struct TestOwnerCap has key, store {
    id: UID,
}

fun create_test_scenario(): Scenario {
    ts::begin(@0xA)
}

fun generate_upgrade_cap(ctx: &mut TxContext): UpgradeCap {
    package::test_publish(mcms_registry::get_multisig_address().to_id(), ctx)
}

fun init_mcms(scenario: &mut Scenario) {
    let ctx = ts::ctx(scenario);
    mcms_registry::test_init(ctx);
    mcms_deployer::test_init(ctx);
    mcms_account::test_init(ctx);
}

fun register_test_package_with_upgrade_cap(scenario: &mut Scenario): address {
    ts::next_tx(scenario, @0xB);
    let mut deployer_state = ts::take_shared<DeployerState>(scenario);
    let mut registry = ts::take_shared<Registry>(scenario);
    let ctx = ts::ctx(scenario);

    let upgrade_cap = generate_upgrade_cap(ctx);
    let package_address = upgrade_cap.package().to_address();

    ts::next_tx(scenario, @0xA);
    let ctx = ts::ctx(scenario);
    let publisher = package::test_claim(MCMS_DEPLOYER_TEST {}, ctx);
    let owner_cap = TestOwnerCap { id: object::new(ctx) };

    let publisher_wrapper = mcms_registry::create_publisher_wrapper(
        &publisher,
        MCMS_DEPLOYER_TEST {},
    );

    mcms_registry::register_entrypoint<MCMS_DEPLOYER_TEST, TestOwnerCap>(
        &mut registry,
        publisher_wrapper,
        MCMS_DEPLOYER_TEST {},
        owner_cap,
        vector[b"mcms_deployer_test"],
        ctx,
    );

    mcms_deployer::register_upgrade_cap(
        &mut deployer_state,
        &registry,
        upgrade_cap,
        ctx,
    );

    transfer::public_transfer(publisher, @0xA);
    ts::return_shared(deployer_state);
    ts::return_shared(registry);

    package_address
}

#[test]
fun test_register_upgrade_cap() {
    let mut scenario = create_test_scenario();

    {
        let ctx = ts::ctx(&mut scenario);
        mcms_registry::test_init(ctx);
        mcms_deployer::test_init(ctx);
        mcms_account::test_init(ctx);
    };

    {
        ts::next_tx(&mut scenario, @0xB);
        let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
        let mut registry = ts::take_shared<Registry>(&scenario);

        let upgrade_cap = generate_upgrade_cap(ts::ctx(&mut scenario));

        ts::next_tx(&mut scenario, @0xA);
        let ctx = ts::ctx(&mut scenario);
        let publisher = package::test_claim(MCMS_DEPLOYER_TEST {}, ctx);

        let publisher_wrapper = mcms_registry::create_publisher_wrapper(
            &publisher,
            MCMS_DEPLOYER_TEST {},
        );

        // First register with MCMS registry
        mcms_registry::register_entrypoint<MCMS_DEPLOYER_TEST, TestOwnerCap>(
            &mut registry,
            publisher_wrapper,
            MCMS_DEPLOYER_TEST {},
            TestOwnerCap { id: object::new(ctx) },
            vector[b"mcms_deployer_test"], // Allowed test module
            ctx,
        );

        // Then register with MCMS deployer
        mcms_deployer::register_upgrade_cap(
            &mut deployer_state,
            &registry,
            upgrade_cap,
            ctx,
        );

        transfer::public_transfer(publisher, @0xA);
        ts::return_shared(deployer_state);
        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
fun test_release_upgrade_cap_at_succeeds() {
    let mut scenario = create_test_scenario();
    init_mcms(&mut scenario);
    let package_address = register_test_package_with_upgrade_cap(&mut scenario);

    ts::next_tx(&mut scenario, @0xA);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let registry = ts::take_shared<Registry>(&scenario);

    assert!(mcms_deployer::has_upgrade_cap(&deployer_state, package_address));

    let upgrade_cap = mcms_deployer::release_upgrade_cap_at(
        &mut deployer_state,
        &registry,
        package_address,
        MCMS_DEPLOYER_TEST {},
    );

    assert!(upgrade_cap.package().to_address() == package_address);
    assert!(!mcms_deployer::has_upgrade_cap(&deployer_state, package_address));

    transfer::public_transfer(upgrade_cap, @0xA);
    ts::return_shared(deployer_state);
    ts::return_shared(registry);
    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms::mcms_deployer::EPackageAddressNotRegistered)]
fun test_release_upgrade_cap_at_fails_after_release_cap() {
    let mut scenario = create_test_scenario();
    init_mcms(&mut scenario);
    let package_address = register_test_package_with_upgrade_cap(&mut scenario);

    ts::next_tx(&mut scenario, @0xA);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let mut registry = ts::take_shared<Registry>(&scenario);

    let _owner_cap = mcms_registry::release_cap<MCMS_DEPLOYER_TEST, TestOwnerCap>(
        &mut registry,
        MCMS_DEPLOYER_TEST {},
    );

    mcms_deployer::release_upgrade_cap_at(
        &mut deployer_state,
        &registry,
        package_address,
        MCMS_DEPLOYER_TEST {},
    );

    ts::return_shared(deployer_state);
    ts::return_shared(registry);
    ts::end(scenario);
}

#[test]
fun test_release_upgrade_cap_at_after_commit_upgrade() {
    let mut scenario = create_test_scenario();
    init_mcms(&mut scenario);
    let old_package_address = register_test_package_with_upgrade_cap(&mut scenario);

    ts::next_tx(&mut scenario, @0xA);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let owner_cap = ts::take_from_sender<mcms_account::OwnerCap>(&scenario);
    let ctx = ts::ctx(&mut scenario);

    let ticket = mcms_deployer::authorize_upgrade(
        &owner_cap,
        &mut deployer_state,
        0,
        vector[],
        old_package_address,
        ctx,
    );
    let receipt = package::test_upgrade(ticket);
    let new_package_address = receipt.package().to_address();

    mcms_deployer::commit_upgrade(&mut deployer_state, receipt, ctx);

    assert!(!mcms_deployer::has_upgrade_cap(&deployer_state, old_package_address));
    assert!(mcms_deployer::has_upgrade_cap(&deployer_state, new_package_address));

    ts::return_to_sender(&scenario, owner_cap);
    ts::return_shared(deployer_state);

    ts::next_tx(&mut scenario, @0xA);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let registry = ts::take_shared<Registry>(&scenario);

    let upgrade_cap = mcms_deployer::release_upgrade_cap_at(
        &mut deployer_state,
        &registry,
        new_package_address,
        MCMS_DEPLOYER_TEST {},
    );

    assert!(upgrade_cap.package().to_address() == new_package_address);

    transfer::public_transfer(upgrade_cap, @0xA);
    ts::return_shared(deployer_state);
    ts::return_shared(registry);
    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms::mcms_deployer::EPackageAddressNotRegistered)]
fun test_release_upgrade_cap_fails_after_commit_upgrade() {
    let mut scenario = create_test_scenario();
    init_mcms(&mut scenario);
    let old_package_address = register_test_package_with_upgrade_cap(&mut scenario);

    ts::next_tx(&mut scenario, @0xA);
    let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
    let owner_cap = ts::take_from_sender<mcms_account::OwnerCap>(&scenario);
    let ctx = ts::ctx(&mut scenario);

    let ticket = mcms_deployer::authorize_upgrade(
        &owner_cap,
        &mut deployer_state,
        0,
        vector[],
        old_package_address,
        ctx,
    );
    let receipt = package::test_upgrade(ticket);
    mcms_deployer::commit_upgrade(&mut deployer_state, receipt, ctx);

    ts::return_to_sender(&scenario, owner_cap);

    let registry = ts::take_shared<Registry>(&scenario);
    mcms_deployer::release_upgrade_cap(
        &mut deployer_state,
        &registry,
        MCMS_DEPLOYER_TEST {},
    );

    ts::return_shared(deployer_state);
    ts::return_shared(registry);
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
        ts::next_tx(&mut scenario, @0xB);
        let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
        let registry = ts::take_shared<Registry>(&scenario);
        let ctx = ts::ctx(&mut scenario);

        let upgrade_cap = generate_upgrade_cap(ctx);

        // This should fail because the package address is not registered
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
