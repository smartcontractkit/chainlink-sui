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

// Non-zero package address used for the upgrade cap so `authorize_upgrade` does not hit the
// Sui >= 1.73 `0x0` in-progress-upgrade sentinel (the mcms named address resolves to 0x0 in tests).
const TEST_PACKAGE: address = @0xCCCC;

fun create_test_scenario(): Scenario {
    ts::begin(@0xA)
}

fun generate_upgrade_cap(ctx: &mut TxContext): UpgradeCap {
    package::test_publish(mcms_registry::get_multisig_address().to_id(), ctx)
}

// Initializes registry/deployer/account and registers MCMS_DEPLOYER_TEST with the registry, plus an
// UpgradeCap published at `TEST_PACKAGE` in the deployer. Leaves the shared objects returned.
fun setup_registered_package_with_upgrade_cap(scenario: &mut Scenario) {
    {
        let ctx = ts::ctx(scenario);
        mcms_registry::test_init(ctx);
        mcms_deployer::test_init(ctx);
        mcms_account::test_init(ctx);
    };

    ts::next_tx(scenario, @0xA);
    let mut deployer_state = ts::take_shared<DeployerState>(scenario);
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

    let upgrade_cap = package::test_publish(TEST_PACKAGE.to_id(), ctx);
    mcms_deployer::test_register_upgrade_cap(&mut deployer_state, upgrade_cap, ctx);

    transfer::public_transfer(publisher, @0xA);
    ts::return_shared(deployer_state);
    ts::return_shared(registry);
}

// Pre-upgrade: current address == original address, so releasing by the current address succeeds.
#[test]
fun test_release_upgrade_cap_current_before_upgrade() {
    let mut scenario = create_test_scenario();
    setup_registered_package_with_upgrade_cap(&mut scenario);

    {
        ts::next_tx(&mut scenario, @0xA);
        let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
        let registry = ts::take_shared<Registry>(&scenario);

        assert!(mcms_deployer::has_upgrade_cap(&deployer_state, TEST_PACKAGE), 0);

        let upgrade_cap = mcms_deployer::release_upgrade_cap_current(
            &mut deployer_state,
            &registry,
            TEST_PACKAGE,
            MCMS_DEPLOYER_TEST {},
        );

        assert!(!mcms_deployer::has_upgrade_cap(&deployer_state, TEST_PACKAGE), 1);

        transfer::public_transfer(upgrade_cap, @0xA);
        ts::return_shared(deployer_state);
        ts::return_shared(registry);
    };

    ts::end(scenario);
}

// Post-upgrade: `commit_upgrade` re-keys the cap to the CURRENT address.
// `release_upgrade_cap_current` releases it by the current address.
#[test]
fun test_release_upgrade_cap_current_after_upgrade() {
    let mut scenario = create_test_scenario();
    setup_registered_package_with_upgrade_cap(&mut scenario);

    let new_package_address: address;

    // Perform one upgrade so the cap is re-keyed to the current address.
    {
        ts::next_tx(&mut scenario, @0xA);
        let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
        let owner_cap = ts::take_from_sender<mcms_account::OwnerCap>(&scenario);
        let ctx = ts::ctx(&mut scenario);

        let ticket = mcms_deployer::authorize_upgrade(
            &owner_cap,
            &mut deployer_state,
            0,
            vector[],
            TEST_PACKAGE,
            ctx,
        );
        let receipt = package::test_upgrade(ticket);
        new_package_address = receipt.package().to_address();
        mcms_deployer::commit_upgrade(&mut deployer_state, receipt, ctx);

        assert!(!mcms_deployer::has_upgrade_cap(&deployer_state, TEST_PACKAGE), 0);
        assert!(mcms_deployer::has_upgrade_cap(&deployer_state, new_package_address), 1);

        ts::return_to_sender(&scenario, owner_cap);
        ts::return_shared(deployer_state);
    };

    // Release the (re-keyed) cap using its CURRENT address.
    {
        ts::next_tx(&mut scenario, @0xA);
        let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
        let registry = ts::take_shared<Registry>(&scenario);

        let upgrade_cap = mcms_deployer::release_upgrade_cap_current(
            &mut deployer_state,
            &registry,
            new_package_address,
            MCMS_DEPLOYER_TEST {},
        );

        assert!(!mcms_deployer::has_upgrade_cap(&deployer_state, new_package_address), 2);

        transfer::public_transfer(upgrade_cap, @0xA);
        ts::return_shared(deployer_state);
        ts::return_shared(registry);
    };

    ts::end(scenario);
}

// A registered proof cannot release a cap at an address that holds no cap.
#[test]
#[expected_failure(abort_code = mcms::mcms_deployer::EPackageAddressNotRegistered)]
fun test_release_upgrade_cap_current_wrong_address_fails() {
    let mut scenario = create_test_scenario();
    setup_registered_package_with_upgrade_cap(&mut scenario);

    {
        ts::next_tx(&mut scenario, @0xA);
        let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
        let registry = ts::take_shared<Registry>(&scenario);

        // @0xDEAD holds no upgrade cap -> should abort.
        let upgrade_cap = mcms_deployer::release_upgrade_cap_current(
            &mut deployer_state,
            &registry,
            @0xDEAD,
            MCMS_DEPLOYER_TEST {},
        );

        transfer::public_transfer(upgrade_cap, @0xA);
        ts::return_shared(deployer_state);
        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms::mcms_deployer::EPackageAddressNotRegistered)]
fun test_old_release_upgrade_cap_fails_after_upgrade() {
    let mut scenario = create_test_scenario();
    setup_registered_package_with_upgrade_cap(&mut scenario);

    // Perform one upgrade so the cap is re-keyed to the current address.
    {
        ts::next_tx(&mut scenario, @0xA);
        let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
        let owner_cap = ts::take_from_sender<mcms_account::OwnerCap>(&scenario);
        let ctx = ts::ctx(&mut scenario);

        let ticket = mcms_deployer::authorize_upgrade(
            &owner_cap,
            &mut deployer_state,
            0,
            vector[],
            TEST_PACKAGE,
            ctx,
        );
        let receipt = package::test_upgrade(ticket);
        mcms_deployer::commit_upgrade(&mut deployer_state, receipt, ctx);

        ts::return_to_sender(&scenario, owner_cap);
        ts::return_shared(deployer_state);
    };

    // The old function looks up the ORIGINAL address, where the cap no longer lives -> abort.
    {
        ts::next_tx(&mut scenario, @0xA);
        let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
        let registry = ts::take_shared<Registry>(&scenario);

        let upgrade_cap = mcms_deployer::release_upgrade_cap(
            &mut deployer_state,
            &registry,
            MCMS_DEPLOYER_TEST {},
        );

        transfer::public_transfer(upgrade_cap, @0xA);
        ts::return_shared(deployer_state);
        ts::return_shared(registry);
    };

    ts::end(scenario);
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
