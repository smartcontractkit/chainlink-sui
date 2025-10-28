#[test_only]
module mcms::mcms_proof_validation_test;

use mcms::mcms_registry;
use mcms::mcms_test;
use mcms::params;
use 0x987::mock_cap;
use std::bcs;
use std::string;
use std::type_name;

public struct TestPackageWitness has drop {}

public struct McmsAcceptOwnershipProof has drop {}

#[test]
fun test_get_callback_params_accept_ownership_valid() {
    let mut env = mcms_test::setup();

    let witness_type = type_name::with_original_ids<TestPackageWitness>();
    let (unregistered_package, witness_module) = params::get_account_address_and_module_name(
        witness_type,
    );

    let params = mcms_registry::test_create_executing_callback_params(
        unregistered_package,
        witness_module,
        string::utf8(b"accept_ownership"),
        bcs::to_bytes(&100),
        x"f1",
        0,
        1,
    );

    let data = mcms_registry::get_accept_ownership_data(
        mcms_test::env_registry(&mut env),
        params,
        McmsAcceptOwnershipProof {},
    );

    assert!(data == bcs::to_bytes(&100));

    mcms_test::destroy(env);
}

#[test]
#[
    expected_failure(
        abort_code = mcms_registry::EOnlyAcceptOwnershipAllowed,
        location = mcms_registry,
    ),
]
fun test_get_accept_ownership_data_wrong_function() {
    let mut env = mcms_test::setup();

    let witness_type = type_name::with_original_ids<TestPackageWitness>();
    let (unregistered_package, witness_module) = params::get_account_address_and_module_name(
        witness_type,
    );

    let params = mcms_registry::test_create_executing_callback_params(
        unregistered_package,
        witness_module,
        string::utf8(b"withdraw"), // NOT accept_ownership
        bcs::to_bytes(&100),
        x"01",
        0,
        1,
    );

    // Abort with EOnlyAcceptOwnershipAllowed
    let _ = mcms_registry::get_accept_ownership_data(
        mcms_test::env_registry(&mut env),
        params,
        McmsAcceptOwnershipProof {},
    );

    mcms_test::destroy(env);
}

#[test]
#[expected_failure(abort_code = mcms_registry::EOnlyMcmsAcceptOwnershipProofAllowed, location = mcms_registry)]
fun test_get_accept_ownership_data_wrong_proof_type() {
    let mut env = mcms_test::setup();

    let witness_type = type_name::with_original_ids<TestPackageWitness>();
    let (unregistered_package, witness_module) = params::get_account_address_and_module_name(
        witness_type,
    );

    let params = mcms_registry::test_create_executing_callback_params(
        unregistered_package,
        witness_module,
        string::utf8(b"accept_ownership"),
        bcs::to_bytes(&100),
        x"02",
        0,
        1,
    );

    // Abort with EOnlyMcmsAcceptOwnershipProofAllowed
    let _ = mcms_registry::get_accept_ownership_data(
        mcms_test::env_registry(&mut env),
        params,
        TestPackageWitness {},
    );

    mcms_test::destroy(env);
}

#[test]
#[expected_failure(abort_code = mcms_registry::EInvalidModuleName, location = mcms_registry)]
fun test_get_accept_ownership_data_module_name_mismatch() {
    let mut env = mcms_test::setup();

    let witness_type = type_name::with_original_ids<TestPackageWitness>();
    let (unregistered_package, _) = params::get_account_address_and_module_name(witness_type);

    let params = mcms_registry::test_create_executing_callback_params(
        unregistered_package,
        string::utf8(b"wrong_module"), // Wrong module name
        string::utf8(b"accept_ownership"),
        bcs::to_bytes(&100),
        x"03",
        0,
        1,
    );

    // Abort with EInvalidModuleName
    let _ = mcms_registry::get_accept_ownership_data(
        mcms_test::env_registry(&mut env),
        params,
        McmsAcceptOwnershipProof {},
    );

    mcms_test::destroy(env);
}

#[test]
#[expected_failure(abort_code = mcms_registry::EPackageIdMismatch, location = mcms_registry)]
fun test_get_accept_ownership_data_package_id_mismatch() {
    let mut env = mcms_test::setup();

    let witness_type = type_name::with_original_ids<TestPackageWitness>();
    let (_, witness_module) = params::get_account_address_and_module_name(witness_type);

    // Create params with WRONG target address that doesn't match TestPackageWitness's address
    let params = mcms_registry::test_create_executing_callback_params(
        @0x9876, // WRONG target - doesn't match TestPackageWitness package address
        witness_module, // Correct module name
        string::utf8(b"accept_ownership"), // Correct function
        bcs::to_bytes(&100),
        x"08",
        0,
        1,
    );

    // Abort with EPackageIdMismatch
    // because target (@0x9876) doesn't match proof_account_address (TestPackageWitness package)
    let _ = mcms_registry::get_accept_ownership_data(
        mcms_test::env_registry(&mut env),
        params,
        McmsAcceptOwnershipProof {},
    );

    mcms_test::destroy(env);
}

#[test]
fun test_get_accept_ownership_data_expected_proof_type_accessor() {
    let mut env = mcms_test::setup();

    let witness_type = type_name::with_original_ids<TestPackageWitness>();
    let (package_addr, module_name) = params::get_account_address_and_module_name(witness_type);

    let params = mcms_registry::test_create_executing_callback_params(
        package_addr,
        module_name,
        string::utf8(b"accept_ownership"),
        bcs::to_bytes(&100),
        x"04",
        0,
        1,
    );

    let _ = mcms_registry::get_accept_ownership_data(
        mcms_test::env_registry(&mut env),
        params,
        McmsAcceptOwnershipProof {},
    );

    mcms_test::destroy(env);
}

#[test]
fun test_get_accept_ownership_data_expected_proof_type_mcms_proof() {
    let mut env = mcms_test::setup();

    let witness_type = type_name::with_original_ids<TestPackageWitness>();
    let (package_addr, module_name) = params::get_account_address_and_module_name(witness_type);

    let params = mcms_registry::test_create_executing_callback_params(
        package_addr,
        module_name,
        string::utf8(b"accept_ownership"),
        bcs::to_bytes(&1),
        x"07",
        0,
        1,
    );

    let _ = mcms_registry::get_accept_ownership_data(
        mcms_test::env_registry(&mut env),
        params,
        McmsAcceptOwnershipProof {},
    );

    mcms_test::destroy(env);
}

#[test]
#[expected_failure(abort_code = mcms_registry::ECapAddressMismatch, location = mcms_registry)]
fun test_register_entrypoint_cap_address_mismatch() {
    // This test verifies that a cap from package B cannot be registered with a PublisherWrapper from package A
    // The assertion `assert!(cap_address == package_address, ECapAddressMismatch)` in register_entrypoint
    // ensures that the cap's package address matches the PublisherWrapper's package address
    //
    // Test setup:
    // - Package A = mcms package at address 0x0
    // - Package B = mock_cap package at address 0x987
    // - We create a PublisherWrapper from Package A (mcms)
    // - We try to register a MockCap from Package B (mock_cap)
    // - This should fail with ECapAddressMismatch

    let mut scenario = sui::test_scenario::begin(@0x1);
    let ctx = scenario.ctx();

    // Create a registry
    mcms_registry::test_init(ctx);

    scenario.next_tx(@0x1);
    let mut registry = scenario.take_shared<mcms_registry::Registry>();

    // Create a Publisher for the mcms package (Package A at address 0x0)
    // using test_claim with TestPackageWitness from the mcms package
    let mcms_publisher = sui::package::test_claim(
        TestPackageWitness {},
        scenario.ctx(),
    );

    // Create a PublisherWrapper using TestPackageWitness from mcms package (0x0)
    let publisher_wrapper = mcms_registry::create_publisher_wrapper(
        &mcms_publisher,
        TestPackageWitness {},
    );

    // Create a MockCap from the mock_cap module (Package B at address 0x987)
    // This is a DIFFERENT package from mcms (0x0)
    let mock_cap_from_different_package = mock_cap::new(scenario.ctx());

    // Try to register the MockCap (from 0x987 package) with publisher_wrapper (from 0x0 package)
    // This should abort with ECapAddressMismatch because:
    // - publisher_wrapper.package_address corresponds to mcms package (0x0)
    // - type_name::with_original_ids<MockCap>().address_string() = "0x987" (mock_cap package)
    // - 0x0 != 0x987, so the assertion fails
    mcms_registry::register_entrypoint(
        &mut registry,
        publisher_wrapper,
        TestPackageWitness {},
        mock_cap_from_different_package,
        vector[b"test_module"],
        scenario.ctx(),
    );

    // Cleanup (this code won't be reached due to expected failure)
    transfer::public_transfer(mcms_publisher, @0x1);
    sui::test_scenario::return_shared(registry);
    scenario.end();
}
