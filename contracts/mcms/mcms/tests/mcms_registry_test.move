#[test_only]
module mcms::mcms_registry_test;

use mcms::mcms_registry::{Self, Registry};
use std::string;
use sui::package::{Self, Publisher};
use sui::test_scenario::{Self as ts, Scenario};

public struct TestModuleCap has key, store {
    id: UID,
}

public struct TestModuleWitness has drop {}

public struct DifferentWitness has drop {}

const MODULE_NAME: vector<u8> = b"mcms_registry_test";

fun create_test_scenario(): Scenario {
    ts::begin(@0xA)
}

#[test_only]
fun create_test_publisher(ctx: &mut TxContext): Publisher {
    package::test_claim(TestModuleWitness {}, ctx)
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
        let publisher = create_test_publisher(ctx);

        let publisher_wrapper = mcms_registry::create_publisher_wrapper(
            &publisher,
            TestModuleWitness {},
        );

        mcms_registry::register_entrypoint<TestModuleWitness, TestModuleCap>(
            &mut registry,
            publisher_wrapper,
            TestModuleWitness {},
            module_cap,
            vector[MODULE_NAME], // Allowed test module
            ctx,
        );

        assert!(
            mcms_registry::is_package_registered(
                &registry,
                mcms_registry::get_multisig_address_ascii(),
            ),
        );

        transfer::public_transfer(publisher, @0xA);
        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
fun test_get_accept_ownership_data() {
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
        let publisher = create_test_publisher(ctx);

        let publisher_wrapper = mcms_registry::create_publisher_wrapper(
            &publisher,
            TestModuleWitness {},
        );

        // Register the module
        mcms_registry::register_entrypoint<TestModuleWitness, TestModuleCap>(
            &mut registry,
            publisher_wrapper,
            TestModuleWitness {},
            module_cap,
            vector[MODULE_NAME], // Allowed test module
            ctx,
        );

        transfer::public_transfer(publisher, @0xA);
        ts::return_shared(registry);
    };

    // Transaction 3: Test get_callback_params
    {
        scenario.next_tx(@0xC);

        let mut registry = scenario.take_shared<Registry>();

        // Create callback params
        let params = mcms_registry::test_create_executing_callback_params(
            mcms_registry::get_multisig_address(),
            string::utf8(MODULE_NAME),
            string::utf8(b"test_function"),
            vector::empty(),
            x"0000000000000000000000000000000000000000000000000000000000000001", // batch_id
            0, // sequence_number
            1, // total_in_batch
        );

        let (cap, _function_name, _data) = mcms_registry::get_callback_params_with_caps<
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
#[expected_failure(abort_code = mcms_registry::EPackageNotRegistered)]
fun test_get_accept_ownership_data_with_unregistered_package_cap() {
    let mut scenario = create_test_scenario();

    // Transaction 1: Initialize registry
    {
        let ctx = scenario.ctx();
        mcms_registry::test_init(ctx);
    };

    // Transaction 2: Try to use accept ownership data without registration
    {
        scenario.next_tx(@0xB);

        let mut registry = scenario.take_shared<Registry>();

        // Create accept ownership data params
        let params = mcms_registry::test_create_executing_callback_params(
            mcms_registry::get_multisig_address(),
            string::utf8(MODULE_NAME),
            string::utf8(b"accept_ownership"),
            vector::empty(),
            x"0000000000000000000000000000000000000000000000000000000000000002", // batch_id
            0, // sequence_number
            1, // total_in_batch
        );

        // This should fail because package is not registered
        let (cap, _function_name, _data) = mcms_registry::get_callback_params_with_caps<
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
#[expected_failure(abort_code = mcms_registry::EPackageIdMismatch)]
fun test_get_callback_params_with_wrong_package_name() {
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
        let publisher = create_test_publisher(ctx);

        let publisher_wrapper = mcms_registry::create_publisher_wrapper(
            &publisher,
            TestModuleWitness {},
        );

        // Register the module
        mcms_registry::register_entrypoint<TestModuleWitness, TestModuleCap>(
            &mut registry,
            publisher_wrapper,
            TestModuleWitness {},
            module_cap,
            vector[MODULE_NAME], // Allowed test module
            ctx,
        );

        // Create callback params with wrong target/package ID
        let params = mcms_registry::test_create_executing_callback_params(
            @0x001,
            string::utf8(b"mcms_registry_test"),
            string::utf8(b"test_function"),
            vector::empty(),
            x"0000000000000000000000000000000000000000000000000000000000000003", // batch_id
            0, // sequence_number
            1, // total_in_batch
        );

        // This should fail because package ID doesn't match
        let (_cap, _function_name, _data) = mcms_registry::get_callback_params_with_caps<
            TestModuleWitness,
            TestModuleCap,
        >(
            &mut registry,
            TestModuleWitness {},
            params,
        );

        transfer::public_transfer(publisher, @0xA);
        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
#[allow(implicit_const_copy)]
fun test_add_allowed_modules() {
    let mut scenario = create_test_scenario();

    // Transaction 1: Initialize registry
    {
        let ctx = ts::ctx(&mut scenario);
        mcms_registry::test_init(ctx);
    };

    ts::next_tx(&mut scenario, @0xA);

    // Transaction 2: Register module with initial allowed modules
    {
        let mut registry = ts::take_shared<Registry>(&scenario);
        let ctx = ts::ctx(&mut scenario);
        let module_cap = TestModuleCap { id: object::new(ctx) };
        let publisher = create_test_publisher(ctx);

        let publisher_wrapper = mcms_registry::create_publisher_wrapper(
            &publisher,
            TestModuleWitness {},
        );

        mcms_registry::register_entrypoint<TestModuleWitness, TestModuleCap>(
            &mut registry,
            publisher_wrapper,
            TestModuleWitness {},
            module_cap,
            vector[MODULE_NAME], // Initial allowed module
            ctx,
        );

        transfer::public_transfer(publisher, @0xA);
        ts::return_shared(registry);
    };

    ts::next_tx(&mut scenario, @0xA);

    // Transaction 3: Add new module to allowed list
    {
        let mut registry = ts::take_shared<Registry>(&scenario);
        let ctx = ts::ctx(&mut scenario);

        // Add new module - no ExecutingCallbackParams needed
        mcms_registry::add_allowed_modules(
            &mut registry,
            TestModuleWitness {},
            vector[b"new_module"],
            ctx,
        );

        // Verify the module was added
        let allowed = mcms_registry::get_allowed_modules(
            &registry,
            mcms_registry::get_multisig_address_ascii(),
        );
        assert!(allowed.contains(&MODULE_NAME), 0);
        assert!(allowed.contains(&b"new_module"), 1);

        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms_registry::EModuleAlreadyAllowed)]
fun test_add_module_already_exists() {
    let mut scenario = create_test_scenario();

    // Transaction 1: Initialize registry
    {
        let ctx = ts::ctx(&mut scenario);
        mcms_registry::test_init(ctx);
    };

    ts::next_tx(&mut scenario, @0xA);

    // Transaction 2: Register module
    {
        let mut registry = ts::take_shared<Registry>(&scenario);
        let ctx = ts::ctx(&mut scenario);
        let module_cap = TestModuleCap { id: object::new(ctx) };
        let publisher = create_test_publisher(ctx);

        let publisher_wrapper = mcms_registry::create_publisher_wrapper(
            &publisher,
            TestModuleWitness {},
        );

        mcms_registry::register_entrypoint<TestModuleWitness, TestModuleCap>(
            &mut registry,
            publisher_wrapper,
            TestModuleWitness {},
            module_cap,
            vector[MODULE_NAME],
            ctx,
        );

        transfer::public_transfer(publisher, @0xA);
        ts::return_shared(registry);
    };

    ts::next_tx(&mut scenario, @0xA);

    // Transaction 3: Try to add the same module again (should fail)
    {
        let mut registry = ts::take_shared<Registry>(&scenario);
        let ctx = ts::ctx(&mut scenario);

        // Try to add MODULE_NAME again (already exists)
        mcms_registry::add_allowed_modules(
            &mut registry,
            TestModuleWitness {},
            vector[MODULE_NAME],
            ctx,
        );

        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms_registry::EWrongProofType)]
fun test_add_module_wrong_proof_type() {
    let mut scenario = create_test_scenario();

    // Transaction 1: Initialize registry
    {
        let ctx = ts::ctx(&mut scenario);
        mcms_registry::test_init(ctx);
    };

    ts::next_tx(&mut scenario, @0xA);

    // Transaction 2: Register module with TestModuleWitness
    {
        let mut registry = ts::take_shared<Registry>(&scenario);
        let ctx = ts::ctx(&mut scenario);
        let module_cap = TestModuleCap { id: object::new(ctx) };
        let publisher = create_test_publisher(ctx);

        let publisher_wrapper = mcms_registry::create_publisher_wrapper(
            &publisher,
            TestModuleWitness {},
        );

        mcms_registry::register_entrypoint<TestModuleWitness, TestModuleCap>(
            &mut registry,
            publisher_wrapper,
            TestModuleWitness {},
            module_cap,
            vector[MODULE_NAME],
            ctx,
        );

        transfer::public_transfer(publisher, @0xA);
        ts::return_shared(registry);
    };

    ts::next_tx(&mut scenario, @0xA);

    // Transaction 3: Try to add module with a different unregistered witness type
    {
        let mut registry = ts::take_shared<Registry>(&scenario);
        let ctx = ts::ctx(&mut scenario);

        // This should fail because DifferentWitness package is not registered `EWrongProofType`
        mcms_registry::add_allowed_modules(
            &mut registry,
            DifferentWitness {},
            vector[b"new_module"],
            ctx,
        );

        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms_registry::EPackageNotRegistered)]
fun test_add_module_package_not_registered() {
    let mut scenario = create_test_scenario();

    // Transaction 1: Initialize registry
    {
        let ctx = ts::ctx(&mut scenario);
        mcms_registry::test_init(ctx);
    };

    ts::next_tx(&mut scenario, @0xA);

    // Transaction 2: Try to add module without registering package first
    {
        let mut registry = ts::take_shared<Registry>(&scenario);
        // This should fail because package is not registered
        mcms_registry::add_allowed_modules(
            &mut registry,
            TestModuleWitness {},
            vector[b"new_module"],
            ts::ctx(&mut scenario),
        );

        ts::return_shared(registry);
    };

    ts::end(scenario);
}

// ================================================================
// |         Remove Allowed Modules Tests                        |
// ================================================================

#[test]
#[allow(implicit_const_copy)]
fun test_remove_allowed_modules() {
    let mut scenario = create_test_scenario();

    // Transaction 1: Initialize registry
    {
        let ctx = ts::ctx(&mut scenario);
        mcms_registry::test_init(ctx);
    };

    ts::next_tx(&mut scenario, @0xA);

    // Transaction 2: Register module with initial allowed modules
    {
        let mut registry = ts::take_shared<Registry>(&scenario);
        let ctx = ts::ctx(&mut scenario);
        let module_cap = TestModuleCap { id: object::new(ctx) };
        let publisher = create_test_publisher(ctx);

        let publisher_wrapper = mcms_registry::create_publisher_wrapper(
            &publisher,
            TestModuleWitness {},
        );

        mcms_registry::register_entrypoint<TestModuleWitness, TestModuleCap>(
            &mut registry,
            publisher_wrapper,
            TestModuleWitness {},
            module_cap,
            vector[MODULE_NAME, b"extra_module"], // Register with two modules
            ctx,
        );

        transfer::public_transfer(publisher, @0xA);
        ts::return_shared(registry);
    };

    ts::next_tx(&mut scenario, @0xA);

    // Transaction 3: Remove one module from allowed list
    {
        let mut registry = ts::take_shared<Registry>(&scenario);
        let ctx = ts::ctx(&mut scenario);

        // Remove extra_module
        mcms_registry::remove_allowed_modules(
            &mut registry,
            TestModuleWitness {},
            vector[b"extra_module"],
            ctx,
        );

        // Verify extra_module was removed, but MODULE_NAME still exists
        let allowed = mcms_registry::get_allowed_modules(
            &registry,
            mcms_registry::get_multisig_address_ascii(),
        );
        assert!(allowed.contains(&MODULE_NAME), 0);
        assert!(!allowed.contains(&b"extra_module"), 1);

        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms_registry::EModuleNotInAllowlist)]
fun test_remove_module_not_in_allowlist() {
    let mut scenario = create_test_scenario();

    // Transaction 1: Initialize registry
    {
        let ctx = ts::ctx(&mut scenario);
        mcms_registry::test_init(ctx);
    };

    ts::next_tx(&mut scenario, @0xA);

    // Transaction 2: Register module
    {
        let mut registry = ts::take_shared<Registry>(&scenario);
        let ctx = ts::ctx(&mut scenario);
        let module_cap = TestModuleCap { id: object::new(ctx) };
        let publisher = create_test_publisher(ctx);

        let publisher_wrapper = mcms_registry::create_publisher_wrapper(
            &publisher,
            TestModuleWitness {},
        );

        mcms_registry::register_entrypoint<TestModuleWitness, TestModuleCap>(
            &mut registry,
            publisher_wrapper,
            TestModuleWitness {},
            module_cap,
            vector[MODULE_NAME],
            ctx,
        );

        transfer::public_transfer(publisher, @0xA);
        ts::return_shared(registry);
    };

    ts::next_tx(&mut scenario, @0xA);

    // Transaction 3: Try to remove a module that doesn't exist (should fail)
    {
        let mut registry = ts::take_shared<Registry>(&scenario);
        let ctx = ts::ctx(&mut scenario);

        // Try to remove nonexistent_module (should fail with EModuleNotInAllowlist)
        mcms_registry::remove_allowed_modules(
            &mut registry,
            TestModuleWitness {},
            vector[b"nonexistent_module"],
            ctx,
        );

        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms_registry::EWrongProofType)]
fun test_remove_module_wrong_proof_type() {
    let mut scenario = create_test_scenario();

    // Transaction 1: Initialize registry
    {
        let ctx = ts::ctx(&mut scenario);
        mcms_registry::test_init(ctx);
    };

    ts::next_tx(&mut scenario, @0xA);

    // Transaction 2: Register module with TestModuleWitness
    {
        let mut registry = ts::take_shared<Registry>(&scenario);
        let ctx = ts::ctx(&mut scenario);
        let module_cap = TestModuleCap { id: object::new(ctx) };
        let publisher = create_test_publisher(ctx);

        let publisher_wrapper = mcms_registry::create_publisher_wrapper(
            &publisher,
            TestModuleWitness {},
        );

        mcms_registry::register_entrypoint<TestModuleWitness, TestModuleCap>(
            &mut registry,
            publisher_wrapper,
            TestModuleWitness {},
            module_cap,
            vector[MODULE_NAME],
            ctx,
        );

        transfer::public_transfer(publisher, @0xA);
        ts::return_shared(registry);
    };

    ts::next_tx(&mut scenario, @0xA);

    // Transaction 3: Try to remove module with a different unregistered witness type
    {
        let mut registry = ts::take_shared<Registry>(&scenario);
        let ctx = ts::ctx(&mut scenario);

        // This should fail because DifferentWitness package is not registered `EWrongProofType`
        mcms_registry::remove_allowed_modules(
            &mut registry,
            DifferentWitness {},
            vector[MODULE_NAME],
            ctx,
        );

        ts::return_shared(registry);
    };

    ts::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms_registry::EPackageNotRegistered)]
fun test_remove_module_package_not_registered() {
    let mut scenario = create_test_scenario();

    // Transaction 1: Initialize registry
    {
        let ctx = ts::ctx(&mut scenario);
        mcms_registry::test_init(ctx);
    };

    ts::next_tx(&mut scenario, @0xA);

    // Transaction 2: Try to remove module without registering package first
    {
        let mut registry = ts::take_shared<Registry>(&scenario);
        let ctx = ts::ctx(&mut scenario);

        // This should fail because package is not registered
        mcms_registry::remove_allowed_modules(
            &mut registry,
            TestModuleWitness {},
            vector[MODULE_NAME],
            ctx,
        );

        ts::return_shared(registry);
    };

    ts::end(scenario);
}
