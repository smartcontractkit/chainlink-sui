#[test_only]
module mcms::mcms_registry_test;

use mcms::mcms_registry::{Self, Registry};
use std::string;
use sui::test_scenario::{Self as ts, Scenario};

public struct TestModuleCap has key, store {
    id: UID,
}

public struct TestModuleWitness has drop {}

const MODULE_NAME: vector<u8> = b"mcms_registry_test";

fun create_test_scenario(): Scenario {
    ts::begin(@0xA)
}

#[test_only]
/// This function acts as a cap gated function.
/// This can be tohught of as calling `set_config` or any admin operation.
fun execute_cap_gated_function(_cap: &TestModuleCap) {}

#[test]
fun test_registry_initialization() {
    let mut scenario = create_test_scenario();

    // Transaction 1: Initialize registry
    {
        let ctx = scenario.ctx();
        mcms_registry::test_init(ctx);
    };

    // Assert that the registry is initialized
    {
        scenario.next_tx(@0xB);
        let registry = scenario.take_shared<Registry>();
        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
fun test_register_entrypoint() {
    let mut scenario = create_test_scenario();

    // Transaction 1: Initialize registry
    {
        let ctx = scenario.ctx();
        mcms_registry::test_init(ctx);
    };

    // Transaction 2: Register a module
    {
        scenario.next_tx(@0xB);

        let mut registry = scenario.take_shared<Registry>();
        let ctx = scenario.ctx();

        let module_cap = TestModuleCap { id: object::new(ctx) };

        mcms_registry::register_entrypoint<TestModuleWitness, TestModuleCap>(
            &mut registry,
            TestModuleWitness {},
            option::some(module_cap),
            ctx,
        );

        let module_name = string::utf8(MODULE_NAME);
        assert!(mcms_registry::is_package_registered(&registry, @mcms));
        assert!(mcms_registry::is_module_registered(&registry, module_name));

        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
fun test_get_callback_params() {
    let mut scenario = create_test_scenario();

    // Transaction 1: Initialize registry
    {
        let ctx = scenario.ctx();
        mcms_registry::test_init(ctx);
    };

    // Transaction 2: Register a module
    {
        scenario.next_tx(@0xB);

        let mut registry = scenario.take_shared<Registry>();
        let ctx = scenario.ctx();

        // Create a module capability
        let module_cap = TestModuleCap { id: object::new(ctx) };

        // Register the module
        mcms_registry::register_entrypoint<TestModuleWitness, TestModuleCap>(
            &mut registry,
            TestModuleWitness {},
            option::some(module_cap),
            ctx,
        );

        ts::return_shared(registry);
    };

    // Transaction 3: Test get_callback_params
    {
        scenario.next_tx(@0xC);

        let mut registry = scenario.take_shared<Registry>();

        // Create callback params
        let params = mcms_registry::test_create_executing_callback_params(
            @mcms,
            string::utf8(MODULE_NAME),
            string::utf8(b"test_function"),
            vector::empty(),
        );

        let (cap, _function_name, _data) = mcms_registry::get_callback_params<
            TestModuleWitness,
            TestModuleCap,
        >(
            &mut registry,
            TestModuleWitness {},
            params,
        );

        // Call function which requires the cap
        execute_cap_gated_function(cap);

        ts::return_shared(registry);
        ts::end(scenario);
    }
}

#[test]
#[expected_failure(abort_code = mcms_registry::EModuleNameMismatch)]
fun test_get_callback_params_with_unregistered_module_name() {
    let mut scenario = create_test_scenario();

    // Transaction 1: Initialize registry
    {
        let ctx = scenario.ctx();
        mcms_registry::test_init(ctx);
    };

    // Transaction 2: Try to use callback without registration
    {
        scenario.next_tx(@0xB);

        let mut registry = scenario.take_shared<Registry>();

        // Create callback params
        let params = mcms_registry::test_create_executing_callback_params(
            @mcms,
            string::utf8(b"mcms_registry_tests"),
            string::utf8(b"test_function"),
            vector::empty(),
        );

        // This should fail because module is not registered
        let (cap, _function_name, _data) = mcms_registry::get_callback_params<
            TestModuleWitness,
            TestModuleCap,
        >(
            &mut registry,
            TestModuleWitness {},
            params,
        );

        execute_cap_gated_function(cap);

        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms_registry::EModuleNameMismatch)]
fun test_get_callback_params_with_wrong_module_name() {
    let mut scenario = create_test_scenario();

    // Transaction 1: Initialize registry
    {
        let ctx = scenario.ctx();
        mcms_registry::test_init(ctx);
    };

    // Transaction 2: Register a module
    {
        scenario.next_tx(@0xB);

        let mut registry = scenario.take_shared<Registry>();
        let ctx = scenario.ctx();

        // Create a module capability
        let module_cap = TestModuleCap { id: object::new(ctx) };

        // Register the module
        mcms_registry::register_entrypoint<TestModuleWitness, TestModuleCap>(
            &mut registry,
            TestModuleWitness {},
            option::some(module_cap),
            ctx,
        );

        // Create callback params with wrong module name
        let params = mcms_registry::test_create_executing_callback_params(
            @mcms,
            string::utf8(b"wrong_module_name"),
            string::utf8(b"test_function"),
            vector::empty(),
        );

        // This should fail because module name doesn't match
        let (_cap, _function_name, _data) = mcms_registry::get_callback_params<
            TestModuleWitness,
            TestModuleCap,
        >(
            &mut registry,
            TestModuleWitness {},
            params,
        );

        ts::return_shared(registry);
    };

    ts::end(scenario);
}
