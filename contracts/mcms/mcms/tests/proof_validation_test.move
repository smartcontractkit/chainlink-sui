#[test_only]
module mcms::proof_validation_test;

use mcms::mcms_registry;
use mcms::mcms_test;
use mcms::params;
use std::bcs;
use std::string;
use std::type_name;

// Mock witness for unregistered package
public struct TestPackageWitness has drop {}

// NOTE: Tests for get_expected_proof_type via timelock_execute_batch are complex
// because they require scheduling operations first. Those code paths are already
// tested via integration tests. The key validation logic is tested via
// get_callback_params tests below which directly exercise the security validations.

/// Test 2b: Mismatched proof types aborts
/// Note: This is hard to test directly because ExecutingCallbackParams is created
/// through timelock_execute_batch which sets the correct expected_proof_type.
/// The mismatch would only occur if there's a bug in get_expected_proof_type.
/// We rely on the other tests to validate get_expected_proof_type is correct.

/// Test 3a: get_callback_params with accept_ownership and McmsProof passes
#[test]
fun test_get_callback_params_accept_ownership_valid() {
    let mut env = mcms_test::setup();

    // Get the package address from TestPackageWitness type
    let witness_type = type_name::with_original_ids<TestPackageWitness>();
    let (unregistered_package, witness_module) = params::get_account_address_and_module_name(witness_type);

    // Create params with accept_ownership and McmsProof
    let params = mcms_registry::test_create_executing_callback_params(
        unregistered_package,
        witness_module,
        string::utf8(b"accept_ownership"),
        bcs::to_bytes(&100u64),
        x"f1",
        0,
        1,
        type_name::with_defining_ids<mcms_registry::McmsProof>(),
    );

    // Call get_callback_params - should succeed
    let (target, module_name, function_name, data) = mcms_registry::get_callback_params(
        mcms_test::env_registry(&mut env),
        params,
        TestPackageWitness {},
    );

    assert!(target == unregistered_package);
    assert!(module_name == witness_module);
    assert!(function_name == string::utf8(b"accept_ownership"));
    assert!(data == bcs::to_bytes(&100u64));

    mcms_test::destroy(env);
}

/// Test 3b: get_callback_params with wrong function name aborts
#[test]
#[expected_failure(abort_code = mcms_registry::EOnlyAcceptOwnershipAllowed, location = mcms_registry)]
fun test_get_callback_params_wrong_function() {
    let mut env = mcms_test::setup();

    // Get the package address from TestPackageWitness type
    let witness_type = type_name::with_original_ids<TestPackageWitness>();
    let (unregistered_package, witness_module) = params::get_account_address_and_module_name(witness_type);

    // Create params with NON-accept_ownership function
    let params = mcms_registry::test_create_executing_callback_params(
        unregistered_package,
        witness_module,
        string::utf8(b"some_other_function"),  // NOT accept_ownership
        bcs::to_bytes(&100u64),
        x"01",
        0,
        1,
        type_name::with_defining_ids<mcms_registry::McmsProof>(),
    );

    // Call get_callback_params - should abort with EOnlyAcceptOwnershipAllowed
    let (_, _, _, _) = mcms_registry::get_callback_params(
        mcms_test::env_registry(&mut env),
        params,
        TestPackageWitness {},
    );

    mcms_test::destroy(env);
}

/// Test 3c: get_callback_params with wrong expected_proof_type aborts
#[test]
#[expected_failure(abort_code = mcms_registry::ENotMcmsAuthorized, location = mcms_registry)]
fun test_get_callback_params_wrong_expected_proof_type() {
    let mut env = mcms_test::setup();

    // Get the package address from TestPackageWitness type
    let witness_type = type_name::with_original_ids<TestPackageWitness>();
    let (unregistered_package, witness_module) = params::get_account_address_and_module_name(witness_type);

    // Create params with accept_ownership but WRONG expected proof type
    let params = mcms_registry::test_create_executing_callback_params(
        unregistered_package,
        witness_module,
        string::utf8(b"accept_ownership"),
        bcs::to_bytes(&100u64),
        x"02",
        0,
        1,
        type_name::with_defining_ids<TestPackageWitness>(),  // WRONG - should be McmsProof
    );

    // Call get_callback_params - should abort with ENotMcmsAuthorized
    let (_, _, _, _) = mcms_registry::get_callback_params(
        mcms_test::env_registry(&mut env),
        params,
        TestPackageWitness {},
    );

    mcms_test::destroy(env);
}

/// Test 3d: get_callback_params validates module name matches proof
#[test]
#[expected_failure(abort_code = mcms_registry::EInvalidModuleName, location = mcms_registry)]
fun test_get_callback_params_module_name_mismatch() {
    let mut env = mcms_test::setup();

    // Get the package address from TestPackageWitness type
    let witness_type = type_name::with_original_ids<TestPackageWitness>();
    let (unregistered_package, _) = params::get_account_address_and_module_name(witness_type);

    // Create params where module_name doesn't match the proof's module
    let params = mcms_registry::test_create_executing_callback_params(
        unregistered_package,
        string::utf8(b"wrong_module"),  // Wrong module name
        string::utf8(b"accept_ownership"),
        bcs::to_bytes(&100u64),
        x"03",
        0,
        1,
        type_name::with_defining_ids<mcms_registry::McmsProof>(),
    );

    // Call get_callback_params - should abort with EInvalidModuleName
    let (_, _, _, _) = mcms_registry::get_callback_params(
        mcms_test::env_registry(&mut env),
        params,
        TestPackageWitness {},
    );

    mcms_test::destroy(env);
}
