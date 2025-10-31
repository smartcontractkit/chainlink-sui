#[test_only]
module ccip::publisher_wrapper_tests;

use ccip::publisher_wrapper;
use ccip::state_object;
use sui::address;
use sui::package;
use sui::test_scenario as ts;

/// OTW for test module
public struct PUBLISHER_WRAPPER_TESTS has drop {}

public struct TestTypeProof has drop {}
public struct TestTypeProof2 has drop {}

const ADMIN: address = @0xADDD;

#[test]
/// Test that publisher wrapper can be created with a valid type from the same module
public fun test_create_publisher_wrapper_with_valid_type() {
    let mut scenario = ts::begin(ADMIN);
    let ctx = scenario.ctx();

    let publisher = package::test_claim(PUBLISHER_WRAPPER_TESTS {}, ctx);
    let expected_package_address = address::from_ascii_bytes(publisher.package().as_bytes());

    // Create publisher wrapper with TestTypeProof - should succeed
    let publisher_wrapper = publisher_wrapper::create(&publisher, TestTypeProof {});

    let extracted_address = publisher_wrapper::destroy(publisher_wrapper);
    assert!(extracted_address == expected_package_address);

    package::burn_publisher(publisher);
    scenario.end();
}

#[test]
/// Test that multiple proof types from the SAME module can create publisher wrapper
/// This verifies that any type from the correct module works
public fun test_publisher_wrapper_accepts_multiple_types_from_same_module() {
    let mut scenario = ts::begin(ADMIN);
    let ctx = scenario.ctx();

    let publisher = package::test_claim(PUBLISHER_WRAPPER_TESTS {}, ctx);
    let expected_package_address = address::from_ascii_bytes(publisher.package().as_bytes());

    // TestTypeProof is from publisher_wrapper_tests module - should succeed
    let wrapper1 = publisher_wrapper::create(&publisher, TestTypeProof {});

    let extracted_address = publisher_wrapper::destroy(wrapper1);
    assert!(extracted_address == expected_package_address);

    // TestTypeProof2 is also from publisher_wrapper_tests module - should also succeed
    let wrapper2 = publisher_wrapper::create(&publisher, TestTypeProof2 {});
    let extracted_address2 = publisher_wrapper::destroy(wrapper2);

    assert!(extracted_address2 == expected_package_address);

    package::burn_publisher(publisher);
    scenario.end();
}

#[test]
/// Test that the package address extracted from publisher wrapper matches the publisher's package
public fun test_publisher_wrapper_extracts_correct_package_address() {
    let mut scenario = ts::begin(ADMIN);
    let ctx = scenario.ctx();

    let publisher = package::test_claim(PUBLISHER_WRAPPER_TESTS {}, ctx);
    let expected_package_address = address::from_ascii_bytes(publisher.package().as_bytes());

    let wrapper = publisher_wrapper::create(&publisher, TestTypeProof {});
    let extracted_address = publisher_wrapper::destroy(wrapper);
    assert!(extracted_address == expected_package_address);

    package::burn_publisher(publisher);
    scenario.end();
}

// ================================================================
// |                   Error Case Tests                           |
// ================================================================

#[test]
#[expected_failure(abort_code = publisher_wrapper::EProofNotAtPublisherAddressAndModule)]
/// Test that creating a publisher wrapper with a proof type from a different module fails
/// Only types from the same module as the publisher can be used
public fun test_create_publisher_wrapper_with_type_from_different_module_fails() {
    let mut scenario = ts::begin(ADMIN);
    let ctx = scenario.ctx();

    let publisher = package::test_claim(PUBLISHER_WRAPPER_TESTS {}, ctx);

    // Try to create publisher wrapper with McmsCallback which is from ccip::state_object module
    // This should fail because McmsCallback is not from publisher_wrapper_tests module
    let publisher_wrapper = publisher_wrapper::create(
        &publisher,
        state_object::test_create_mcms_callback(),
    );

    // Should never reach here
    publisher_wrapper::destroy(publisher_wrapper);
    package::burn_publisher(publisher);
    scenario.end();
}
