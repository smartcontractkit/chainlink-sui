#[test_only]
module mcms::mcms_deployer_test;

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
    package::test_publish(@mcms.to_id(), ctx)
}

#[test]
fun test_register_upgrade_cap() {
    let mut scenario = create_test_scenario();

    {
        let ctx = ts::ctx(&mut scenario);
        mcms_registry::test_init(ctx);
        mcms_deployer::test_init(ctx);
    };

    {
        ts::next_tx(&mut scenario, @0xB);
        let mut deployer_state = ts::take_shared<DeployerState>(&scenario);
        let mut registry = ts::take_shared<Registry>(&scenario);
        let ctx = ts::ctx(&mut scenario);

        let upgrade_cap = generate_upgrade_cap(ctx);

        // First register with MCMS registry
        mcms_registry::register_entrypoint<MCMS_DEPLOYER_TEST, TestOwnerCap>(
            &mut registry,
            MCMS_DEPLOYER_TEST {},
            option::some(TestOwnerCap { id: object::new(ctx) }),
            ctx,
        );

        // Then register with MCMS deployer
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
