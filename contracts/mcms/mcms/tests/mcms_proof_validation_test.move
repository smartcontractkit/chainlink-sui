#[test_only]
module mcms::mcms_proof_validation_test;

use mcms::mcms;
use mcms::mcms_registry;
use mcms::mcms_test;
use mcms::params;
use std::bcs;
use std::string;
use std::type_name;

public struct TestPackageWitness has drop {}

#[test]
#[expected_failure(abort_code = mcms::EWrongProofType)]
fun test_validate_proof_type_registered_package_mismatch() {
    let mut env = mcms_test::setup();

    let mcms_address = mcms_registry::get_multisig_address();

    // Create params with TestPackageWitness as expected type (wrong!)
    // MCMS is registered with McmsProof, not TestPackageWitness
    let params = mcms_registry::test_create_executing_callback_params(
        mcms_address,
        string::utf8(b"mcms_account"),
        string::utf8(b"accept_ownership_as_timelock"),
        bcs::to_bytes(&100),
        x"b1",
        0,
        1,
        type_name::with_defining_ids<TestPackageWitness>(), // WRONG proof type
    );

    // This should fail because expected_proof_type doesn't match registered type
    mcms_test::test_mcms_dispatch_to_account(&mut env, params);

    mcms_test::destroy(env);
}

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
        bcs::to_bytes(&100u64),
        x"f1",
        0,
        1,
        type_name::with_defining_ids<mcms_registry::McmsProof>(),
    );

    let (target, module_name, function_name, data) = mcms_registry::get_callback_params(
        mcms_test::env_registry(&mut env),
        params,
        TestPackageWitness {},
    );

    assert!(target == unregistered_package);
    assert!(module_name == witness_module);
    assert!(function_name == string::utf8(b"accept_ownership"));
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
fun test_get_callback_params_wrong_function() {
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
        type_name::with_defining_ids<mcms_registry::McmsProof>(),
    );

    // Abort with EOnlyAcceptOwnershipAllowed
    let (_, _, _, _) = mcms_registry::get_callback_params(
        mcms_test::env_registry(&mut env),
        params,
        TestPackageWitness {},
    );

    mcms_test::destroy(env);
}

#[test]
#[expected_failure(abort_code = mcms_registry::ENotMcmsAuthorized, location = mcms_registry)]
fun test_get_callback_params_wrong_expected_proof_type() {
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
        type_name::with_defining_ids<TestPackageWitness>(), // WRONG - should be McmsProof
    );

    // Abort with ENotMcmsAuthorized
    let (_, _, _, _) = mcms_registry::get_callback_params(
        mcms_test::env_registry(&mut env),
        params,
        TestPackageWitness {},
    );

    mcms_test::destroy(env);
}

#[test]
#[expected_failure(abort_code = mcms_registry::EInvalidModuleName, location = mcms_registry)]
fun test_get_callback_params_module_name_mismatch() {
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
        type_name::with_defining_ids<mcms_registry::McmsProof>(),
    );

    // Abort with EInvalidModuleName
    let (_, _, _, _) = mcms_registry::get_callback_params(
        mcms_test::env_registry(&mut env),
        params,
        TestPackageWitness {},
    );

    mcms_test::destroy(env);
}

#[test]
#[expected_failure(abort_code = mcms_registry::EPackageIdMismatch, location = mcms_registry)]
fun test_get_callback_params_package_id_mismatch() {
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
        type_name::with_defining_ids<mcms_registry::McmsProof>(), // Correct expected proof type
    );

    // Abort with EPackageIdMismatch
    // because target (@0x9876) doesn't match proof_account_address (TestPackageWitness package)
    let (_, _, _, _) = mcms_registry::get_callback_params(
        mcms_test::env_registry(&mut env),
        params,
        TestPackageWitness {},
    );

    mcms_test::destroy(env);
}

#[test]
fun test_expected_proof_type_accessor() {
    let mut env = mcms_test::setup();

    let witness_type = type_name::with_original_ids<TestPackageWitness>();
    let (package_addr, module_name) = params::get_account_address_and_module_name(witness_type);

    let expected_type = type_name::with_defining_ids<mcms_registry::McmsProof>();
    let params = mcms_registry::test_create_executing_callback_params(
        package_addr,
        module_name,
        string::utf8(b"accept_ownership"),
        bcs::to_bytes(&100),
        x"04",
        0,
        1,
        expected_type,
    );

    assert!(mcms_registry::expected_proof_type(&params) == expected_type);

    let (_, _, _, _) = mcms_registry::get_callback_params(
        mcms_test::env_registry(&mut env),
        params,
        TestPackageWitness {},
    );

    mcms_test::destroy(env);
}

#[test]
fun test_expected_proof_type_mcms_proof() {
    let mut env = mcms_test::setup();

    let witness_type = type_name::with_original_ids<TestPackageWitness>();
    let (package_addr, module_name) = params::get_account_address_and_module_name(witness_type);

    let mcms_type = type_name::with_defining_ids<mcms_registry::McmsProof>();

    let params = mcms_registry::test_create_executing_callback_params(
        package_addr,
        module_name,
        string::utf8(b"accept_ownership"),
        bcs::to_bytes(&1),
        x"07",
        0,
        1,
        mcms_type,
    );

    assert!(mcms_registry::expected_proof_type(&params) == mcms_type);

    let (_, _, _, _) = mcms_registry::get_callback_params(
        mcms_test::env_registry(&mut env),
        params,
        TestPackageWitness {},
    );

    mcms_test::destroy(env);
}
