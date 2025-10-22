#[test_only]
module mcms::mcms_proof_validation_test;

use mcms::mcms_registry;
use mcms::mcms_test;
use mcms::params;
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
