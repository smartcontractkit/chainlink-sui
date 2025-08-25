#[test_only]
module ccip::receiver_registry_tests;

use ccip::ownable::OwnerCap;
use ccip::receiver_registry::{Self, ReceiverRegistry};
use ccip::state_object::{Self, CCIPObjectRef};
use std::ascii;
use std::string;
use std::type_name;
use sui::address;
use sui::test_scenario::{Self as ts, Scenario};

public struct RECEIVER_REGISTRY_TESTS has drop {}

public struct TestReceiverProof has drop {}
public struct TestReceiverProof2 has drop {}

const OWNER: address = @0x1000;

// Helper function to get the package ID from a proof type
fun get_package_id_from_proof<ProofType>(): address {
    let proof_tn = type_name::get<ProofType>();
    let address_str = type_name::get_address(&proof_tn);
    address::from_ascii_bytes(&std::ascii::into_bytes(address_str))
}

fun create_test_scenario(addr: address): Scenario {
    ts::begin(addr)
}

fun setup_test(): (Scenario, CCIPObjectRef, OwnerCap) {
    let mut scenario = create_test_scenario(OWNER);
    {
        let ctx = scenario.ctx();
        state_object::test_init(ctx);
    };

    scenario.next_tx(OWNER);
    {
        let ref = scenario.take_shared<CCIPObjectRef>();
        let owner_cap = scenario.take_from_sender<OwnerCap>();
        (scenario, ref, owner_cap)
    }
}

fun cleanup_test(scenario: Scenario, ref: CCIPObjectRef, owner_cap: OwnerCap) {
    // Return the owner cap back to the sender instead of destroying it
    ts::return_to_sender(&scenario, owner_cap);
    // Return the shared object back to the scenario instead of destroying it
    ts::return_shared(ref);
    ts::end(scenario);
}

#[test]
public fun test_initialize() {
    let (mut scenario, mut ref, owner_cap) = setup_test();
    let ctx = scenario.ctx();

    receiver_registry::initialize(&mut ref, &owner_cap, ctx);

    // Verify the registry state was created
    assert!(state_object::contains<ReceiverRegistry>(&ref));

    cleanup_test(scenario, ref, owner_cap);
}

#[test]
#[expected_failure(abort_code = receiver_registry::EAlreadyInitialized)]
public fun test_initialize_already_initialized() {
    let (mut scenario, mut ref, owner_cap) = setup_test();
    let ctx = scenario.ctx();

    receiver_registry::initialize(&mut ref, &owner_cap, ctx);
    // Try to initialize again - should fail
    receiver_registry::initialize(&mut ref, &owner_cap, ctx);

    cleanup_test(scenario, ref, owner_cap);
}

#[test]
public fun test_register_receiver() {
    let (mut scenario, mut ref, owner_cap) = setup_test();
    let ctx = scenario.ctx();

    receiver_registry::initialize(&mut ref, &owner_cap, ctx);

    // Register a receiver
    receiver_registry::register_receiver(&mut ref, TestReceiverProof {});

    // Verify the receiver is registered
    let package_id_1 = get_package_id_from_proof<TestReceiverProof>();
    assert!(receiver_registry::is_registered_receiver(&ref, package_id_1));

    // Get receiver config and verify fields
    let config = receiver_registry::get_receiver_config(&ref, package_id_1);
    let (module_name, proof_typename) = receiver_registry::get_receiver_config_fields(config);

    assert!(module_name == string::utf8(b"receiver_registry_tests"));
    assert!(proof_typename == type_name::into_string(type_name::get<TestReceiverProof>()));

    cleanup_test(scenario, ref, owner_cap);
}

#[test]
#[expected_failure(abort_code = receiver_registry::EAlreadyRegistered)]
public fun test_register_receiver_already_registered() {
    let (mut scenario, mut ref, owner_cap) = setup_test();
    let ctx = scenario.ctx();

    receiver_registry::initialize(&mut ref, &owner_cap, ctx);

    // Register a receiver
    receiver_registry::register_receiver(&mut ref, TestReceiverProof {});

    // Try to register the same receiver again - should fail
    receiver_registry::register_receiver(&mut ref, TestReceiverProof {});

    cleanup_test(scenario, ref, owner_cap);
}

#[test]
#[expected_failure(abort_code = receiver_registry::EAlreadyRegistered)]
public fun test_register_receiver_same_package_different_proof() {
    let (mut scenario, mut ref, owner_cap) = setup_test();
    let ctx = scenario.ctx();

    receiver_registry::initialize(&mut ref, &owner_cap, ctx);

    // Register a receiver with TestReceiverProof
    receiver_registry::register_receiver(&mut ref, TestReceiverProof {});

    // Try to register with TestReceiverProof2 (same package ID) - should fail
    receiver_registry::register_receiver(&mut ref, TestReceiverProof2 {});

    cleanup_test(scenario, ref, owner_cap);
}

#[test]
public fun test_register_multiple_receivers_same_package() {
    let (mut scenario, mut ref, owner_cap) = setup_test();
    let ctx = scenario.ctx();

    receiver_registry::initialize(&mut ref, &owner_cap, ctx);

    // Register first receiver
    receiver_registry::register_receiver(&mut ref, TestReceiverProof {});

    // Verify both proof types have the same package ID (they're in the same module)
    let package_id_1 = get_package_id_from_proof<TestReceiverProof>();
    let package_id_2 = get_package_id_from_proof<TestReceiverProof2>();
    assert!(package_id_1 == package_id_2);

    // Verify the receiver is registered
    assert!(receiver_registry::is_registered_receiver(&ref, package_id_1));

    // Verify the config contains the first proof type
    let config = receiver_registry::get_receiver_config(&ref, package_id_1);
    let (_, proof_type) = receiver_registry::get_receiver_config_fields(config);

    assert!(proof_type == type_name::into_string(type_name::get<TestReceiverProof>()));

    cleanup_test(scenario, ref, owner_cap);
}

#[test]
public fun test_unregister_receiver() {
    let (mut scenario, mut ref, owner_cap) = setup_test();
    let ctx = scenario.ctx();

    receiver_registry::initialize(&mut ref, &owner_cap, ctx);

    // Register a receiver
    receiver_registry::register_receiver(&mut ref, TestReceiverProof {});

    // Verify it's registered
    let package_id_1 = get_package_id_from_proof<TestReceiverProof>();
    assert!(receiver_registry::is_registered_receiver(&ref, package_id_1));

    // Unregister the receiver
    receiver_registry::unregister_receiver(&mut ref, &owner_cap, package_id_1, ctx);

    // Verify it's no longer registered
    assert!(!receiver_registry::is_registered_receiver(&ref, package_id_1));

    cleanup_test(scenario, ref, owner_cap);
}

#[test]
#[expected_failure(abort_code = receiver_registry::EUnknownReceiver)]
public fun test_unregister_receiver_unknown() {
    let (mut scenario, mut ref, owner_cap) = setup_test();
    let ctx = scenario.ctx();

    receiver_registry::initialize(&mut ref, &owner_cap, ctx);

    // Try to unregister a receiver that was never registered
    let package_id_1 = get_package_id_from_proof<TestReceiverProof>();
    receiver_registry::unregister_receiver(&mut ref, &owner_cap, package_id_1, ctx);

    cleanup_test(scenario, ref, owner_cap);
}

#[test]
public fun test_is_registered_receiver() {
    let (mut scenario, mut ref, owner_cap) = setup_test();
    let ctx = scenario.ctx();

    receiver_registry::initialize(&mut ref, &owner_cap, ctx);

    // Check unregistered receiver
    let package_id_1 = get_package_id_from_proof<TestReceiverProof>();
    assert!(!receiver_registry::is_registered_receiver(&ref, package_id_1));

    // Register receiver
    receiver_registry::register_receiver(&mut ref, TestReceiverProof {});

    // Check registered receiver
    assert!(receiver_registry::is_registered_receiver(&ref, package_id_1));

    // Both proof types have the same package ID since they're in the same module
    let package_id_2 = get_package_id_from_proof<TestReceiverProof2>();
    assert!(package_id_1 == package_id_2);
    assert!(receiver_registry::is_registered_receiver(&ref, package_id_2));

    cleanup_test(scenario, ref, owner_cap);
}

#[test]
#[expected_failure(abort_code = receiver_registry::EUnknownReceiver)]
public fun test_get_receiver_config_unknown() {
    let (mut scenario, mut ref, owner_cap) = setup_test();
    let ctx = scenario.ctx();

    receiver_registry::initialize(&mut ref, &owner_cap, ctx);

    // Try to get config for unregistered receiver
    let package_id_1 = get_package_id_from_proof<TestReceiverProof>();
    let _config = receiver_registry::get_receiver_config(&ref, package_id_1);

    cleanup_test(scenario, ref, owner_cap);
}

#[test]
public fun test_get_receiver_config() {
    let (mut scenario, mut ref, owner_cap) = setup_test();
    let ctx = scenario.ctx();

    receiver_registry::initialize(&mut ref, &owner_cap, ctx);

    // Register a receiver
    receiver_registry::register_receiver(&mut ref, TestReceiverProof {});

    // Get the config
    let package_id_1 = get_package_id_from_proof<TestReceiverProof>();
    let config = receiver_registry::get_receiver_config(&ref, package_id_1);
    let (module_name, proof_typename) = receiver_registry::get_receiver_config_fields(config);

    // Verify all fields
    assert!(module_name == string::utf8(b"receiver_registry_tests"));
    assert!(proof_typename == type_name::into_string(type_name::get<TestReceiverProof>()));

    cleanup_test(scenario, ref, owner_cap);
}

#[test]
public fun test_get_receiver_module_and_state() {
    let (mut scenario, mut ref, owner_cap) = setup_test();
    let ctx = scenario.ctx();

    receiver_registry::initialize(&mut ref, &owner_cap, ctx);

    // Test unregistered receiver - should return empty values
    let package_id_1 = get_package_id_from_proof<TestReceiverProof>();
    let (module_name, proof_typename_str) = receiver_registry::get_receiver_info(
        &ref,
        package_id_1,
    );
    assert!(module_name == string::utf8(b""));
    assert!(proof_typename_str == ascii::string(b""));

    // Register a receiver
    receiver_registry::register_receiver(&mut ref, TestReceiverProof {});

    // Test registered receiver - should return actual values
    let (module_name, proof_typename_str) = receiver_registry::get_receiver_info(
        &ref,
        package_id_1,
    );
    assert!(module_name == string::utf8(b"receiver_registry_tests"));
    // The proof typename string should contain the test receiver proof type
    assert!(proof_typename_str == type_name::into_string(type_name::get<TestReceiverProof>()));

    cleanup_test(scenario, ref, owner_cap);
}

#[test]
public fun test_register_receiver_with_zero_state_id() {
    let (mut scenario, mut ref, owner_cap) = setup_test();
    let ctx = scenario.ctx();

    receiver_registry::initialize(&mut ref, &owner_cap, ctx);

    // Register a receiver (stateless receiver)
    receiver_registry::register_receiver(&mut ref, TestReceiverProof {});

    // Verify the receiver is registered
    let package_id_1 = get_package_id_from_proof<TestReceiverProof>();
    let config = receiver_registry::get_receiver_config(&ref, package_id_1);
    let (_, _) = receiver_registry::get_receiver_config_fields(config);

    // Verify get_receiver_info returns correct values
    let (module_name, proof_typename_str) = receiver_registry::get_receiver_info(
        &ref,
        package_id_1,
    );
    assert!(module_name == string::utf8(b"receiver_registry_tests"));
    assert!(proof_typename_str == type_name::into_string(type_name::get<TestReceiverProof>()));

    cleanup_test(scenario, ref, owner_cap);
}

#[test]
public fun test_complete_receiver_lifecycle() {
    let (mut scenario, mut ref, owner_cap) = setup_test();
    let ctx = scenario.ctx();

    receiver_registry::initialize(&mut ref, &owner_cap, ctx);

    // 1. Initially not registered
    let package_id_1 = get_package_id_from_proof<TestReceiverProof>();
    assert!(!receiver_registry::is_registered_receiver(&ref, package_id_1));

    // 2. Register receiver
    receiver_registry::register_receiver(&mut ref, TestReceiverProof {});
    assert!(receiver_registry::is_registered_receiver(&ref, package_id_1));

    // 3. Verify config is correct
    let config = receiver_registry::get_receiver_config(&ref, package_id_1);
    let (module_name, proof_typename) = receiver_registry::get_receiver_config_fields(config);

    assert!(module_name == string::utf8(b"receiver_registry_tests"));
    assert!(proof_typename == type_name::into_string(type_name::get<TestReceiverProof>()));

    // 4. Verify module and proof typename lookup
    let (lookup_module, lookup_proof_typename_str) = receiver_registry::get_receiver_info(
        &ref,
        package_id_1,
    );
    assert!(lookup_module == module_name);
    assert!(lookup_proof_typename_str == proof_typename);

    // 5. Unregister receiver
    receiver_registry::unregister_receiver(&mut ref, &owner_cap, package_id_1, ctx);
    assert!(!receiver_registry::is_registered_receiver(&ref, package_id_1));

    // 6. Verify lookup returns empty values after unregistration
    let (empty_module, empty_proof_typename_str) = receiver_registry::get_receiver_info(
        &ref,
        package_id_1,
    );
    assert!(empty_module == string::utf8(b""));
    assert!(empty_proof_typename_str == ascii::string(b""));

    cleanup_test(scenario, ref, owner_cap);
}

#[test]
public fun test_type_and_version() {
    let version = receiver_registry::type_and_version();
    assert!(version == string::utf8(b"ReceiverRegistry 1.6.0"));
}
