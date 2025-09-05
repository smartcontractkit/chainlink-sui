#[test_only]
module ccip::upgrade_registry_test;

use ccip::ownable::OwnerCap;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::upgrade_registry;
use std::string;
use sui::test_scenario::{Self, Scenario};

const SENDER_1: address = @0x1;

fun set_up_test(): (Scenario, OwnerCap, CCIPObjectRef) {
    let mut scenario = test_scenario::begin(SENDER_1);
    let ctx = scenario.ctx();

    state_object::test_init(ctx);

    // Advance to next transaction to retrieve the created objects
    scenario.next_tx(SENDER_1);

    // Retrieve the OwnerCap that was transferred to SENDER_1
    let owner_cap = scenario.take_from_sender<OwnerCap>();

    // Retrieve the shared CCIPObjectRef
    let mut ref = scenario.take_shared<CCIPObjectRef>();

    // Initialize the upgrade registry
    upgrade_registry::initialize(&mut ref, &owner_cap, scenario.ctx());

    (scenario, owner_cap, ref)
}

fun tear_down_test(scenario: Scenario, owner_cap: OwnerCap, ref: CCIPObjectRef) {
    // Return the owner cap back to the sender instead of destroying it
    test_scenario::return_to_sender(&scenario, owner_cap);
    // Return the shared object back to the scenario instead of destroying it
    test_scenario::return_shared(ref);
    test_scenario::end(scenario);
}

// =================== Initialization Tests =================== //

#[test]
public fun test_initialize() {
    let (scenario, owner_cap, ref) = set_up_test();

    // Test that we can get empty restrictions initially
    let empty_function_restrictions = upgrade_registry::get_function_restrictions(
        &ref,
        string::utf8(b"test_module"),
        string::utf8(b"test_function"),
    );
    assert!(empty_function_restrictions.is_empty());

    let empty_module_restrictions = upgrade_registry::get_module_restrictions(
        &ref,
        string::utf8(b"test_module"),
    );
    assert!(empty_module_restrictions.is_empty());

    let (package_ids, versions, timestamps) = upgrade_registry::get_package_history(
        &ref,
        string::utf8(b"test_package"),
    );
    assert!(package_ids.is_empty());
    assert!(versions.is_empty());
    assert!(timestamps.is_empty());

    tear_down_test(scenario, owner_cap, ref);
}

// =================== Function Restrictions Tests =================== //

#[test]
public fun test_update_and_get_function_restrictions() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");
    let blocked_versions = vector[1u64, 2u64, 5u64];

    // Update function restrictions
    upgrade_registry::update_function_restrictions(
        &mut ref,
        &owner_cap,
        module_name,
        function_name,
        blocked_versions,
        scenario.ctx(),
    );

    // Get and verify function restrictions
    let retrieved_restrictions = upgrade_registry::get_function_restrictions(
        &ref,
        module_name,
        function_name,
    );
    assert!(retrieved_restrictions == blocked_versions);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_update_function_restrictions_overwrites_existing() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");
    let initial_blocked_versions = vector[1u64, 2u64];
    let updated_blocked_versions = vector[3u64, 4u64, 5u64];

    // Set initial restrictions
    upgrade_registry::update_function_restrictions(
        &mut ref,
        &owner_cap,
        module_name,
        function_name,
        initial_blocked_versions,
        scenario.ctx(),
    );

    // Update with new restrictions
    upgrade_registry::update_function_restrictions(
        &mut ref,
        &owner_cap,
        module_name,
        function_name,
        updated_blocked_versions,
        scenario.ctx(),
    );

    // Verify the restrictions were overwritten
    let retrieved_restrictions = upgrade_registry::get_function_restrictions(
        &ref,
        module_name,
        function_name,
    );
    assert!(retrieved_restrictions == updated_blocked_versions);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_is_function_allowed() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");
    let blocked_versions = vector[1u64, 3u64, 5u64];

    // Initially, all versions should be allowed
    assert!(upgrade_registry::is_function_allowed(&ref, module_name, function_name, 1u64));
    assert!(upgrade_registry::is_function_allowed(&ref, module_name, function_name, 2u64));

    // Set function restrictions
    upgrade_registry::update_function_restrictions(
        &mut ref,
        &owner_cap,
        module_name,
        function_name,
        blocked_versions,
        scenario.ctx(),
    );

    // Test blocked versions
    assert!(!upgrade_registry::is_function_allowed(&ref, module_name, function_name, 1u64));
    assert!(!upgrade_registry::is_function_allowed(&ref, module_name, function_name, 3u64));
    assert!(!upgrade_registry::is_function_allowed(&ref, module_name, function_name, 5u64));

    // Test allowed versions
    assert!(upgrade_registry::is_function_allowed(&ref, module_name, function_name, 2u64));
    assert!(upgrade_registry::is_function_allowed(&ref, module_name, function_name, 4u64));
    assert!(upgrade_registry::is_function_allowed(&ref, module_name, function_name, 6u64));

    tear_down_test(scenario, owner_cap, ref);
}

// =================== Module Restrictions Tests =================== //

#[test]
public fun test_update_and_get_module_restrictions() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let blocked_versions = vector[1u64, 2u64, 5u64];

    // Update module restrictions
    upgrade_registry::update_module_restrictions(
        &mut ref,
        &owner_cap,
        module_name,
        blocked_versions,
        scenario.ctx(),
    );

    // Get and verify module restrictions
    let retrieved_restrictions = upgrade_registry::get_module_restrictions(&ref, module_name);
    assert!(retrieved_restrictions == blocked_versions);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_update_module_restrictions_overwrites_existing() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let initial_blocked_versions = vector[1u64, 2u64];
    let updated_blocked_versions = vector[3u64, 4u64, 5u64];

    // Set initial restrictions
    upgrade_registry::update_module_restrictions(
        &mut ref,
        &owner_cap,
        module_name,
        initial_blocked_versions,
        scenario.ctx(),
    );

    // Update with new restrictions
    upgrade_registry::update_module_restrictions(
        &mut ref,
        &owner_cap,
        module_name,
        updated_blocked_versions,
        scenario.ctx(),
    );

    // Verify the restrictions were overwritten
    let retrieved_restrictions = upgrade_registry::get_module_restrictions(&ref, module_name);
    assert!(retrieved_restrictions == updated_blocked_versions);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_is_module_allowed() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let blocked_versions = vector[1u64, 3u64, 5u64];

    // Initially, all versions should be allowed
    assert!(upgrade_registry::is_module_allowed(&ref, module_name, 1u64));
    assert!(upgrade_registry::is_module_allowed(&ref, module_name, 2u64));

    // Set module restrictions
    upgrade_registry::update_module_restrictions(
        &mut ref,
        &owner_cap,
        module_name,
        blocked_versions,
        scenario.ctx(),
    );

    // Test blocked versions
    assert!(!upgrade_registry::is_module_allowed(&ref, module_name, 1u64));
    assert!(!upgrade_registry::is_module_allowed(&ref, module_name, 3u64));
    assert!(!upgrade_registry::is_module_allowed(&ref, module_name, 5u64));

    // Test allowed versions
    assert!(upgrade_registry::is_module_allowed(&ref, module_name, 2u64));
    assert!(upgrade_registry::is_module_allowed(&ref, module_name, 4u64));
    assert!(upgrade_registry::is_module_allowed(&ref, module_name, 6u64));

    tear_down_test(scenario, owner_cap, ref);
}

// =================== Function and Module Interaction Tests =================== //

#[test]
public fun test_function_blocked_when_module_blocked() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");
    let blocked_version = 1u64;

    // Block the module version
    upgrade_registry::update_module_restrictions(
        &mut ref,
        &owner_cap,
        module_name,
        vector[blocked_version],
        scenario.ctx(),
    );

    // Function should be blocked even though it has no specific restrictions
    assert!(
        !upgrade_registry::is_function_allowed(&ref, module_name, function_name, blocked_version),
    );

    // Function should be allowed for non-blocked module versions
    assert!(upgrade_registry::is_function_allowed(&ref, module_name, function_name, 2u64));

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_function_blocked_when_both_module_and_function_have_restrictions() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");

    // Block module version 1 and 2
    upgrade_registry::update_module_restrictions(
        &mut ref,
        &owner_cap,
        module_name,
        vector[1u64, 2u64],
        scenario.ctx(),
    );

    // Block function version 3 and 4
    upgrade_registry::update_function_restrictions(
        &mut ref,
        &owner_cap,
        module_name,
        function_name,
        vector[3u64, 4u64],
        scenario.ctx(),
    );

    // Version 1, 2 should be blocked due to module restrictions
    assert!(!upgrade_registry::is_function_allowed(&ref, module_name, function_name, 1u64));
    assert!(!upgrade_registry::is_function_allowed(&ref, module_name, function_name, 2u64));

    // Version 3, 4 should be blocked due to function restrictions
    assert!(!upgrade_registry::is_function_allowed(&ref, module_name, function_name, 3u64));
    assert!(!upgrade_registry::is_function_allowed(&ref, module_name, function_name, 4u64));

    // Version 5 should be allowed (not blocked by either)
    assert!(upgrade_registry::is_function_allowed(&ref, module_name, function_name, 5u64));

    tear_down_test(scenario, owner_cap, ref);
}

// =================== Package History Tests =================== //

#[test]
public fun test_get_package_history_empty() {
    let (scenario, owner_cap, ref) = set_up_test();

    let package_name = string::utf8(b"test_package");
    let (package_ids, versions, timestamps) = upgrade_registry::get_package_history(
        &ref,
        package_name,
    );

    assert!(package_ids.is_empty());
    assert!(versions.is_empty());
    assert!(timestamps.is_empty());

    tear_down_test(scenario, owner_cap, ref);
}

// =================== Edge Cases and Error Conditions =================== //

#[test]
public fun test_empty_blocked_versions_list() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");
    let empty_blocked_versions = vector::empty<u64>();

    // Set empty restrictions
    upgrade_registry::update_function_restrictions(
        &mut ref,
        &owner_cap,
        module_name,
        function_name,
        empty_blocked_versions,
        scenario.ctx(),
    );

    upgrade_registry::update_module_restrictions(
        &mut ref,
        &owner_cap,
        module_name,
        empty_blocked_versions,
        scenario.ctx(),
    );

    // All versions should be allowed
    assert!(upgrade_registry::is_function_allowed(&ref, module_name, function_name, 1u64));
    assert!(upgrade_registry::is_module_allowed(&ref, module_name, 1u64));

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_multiple_modules_and_functions() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module1 = string::utf8(b"module1");
    let module2 = string::utf8(b"module2");
    let function1 = string::utf8(b"function1");
    let function2 = string::utf8(b"function2");

    // Set different restrictions for different modules and functions
    upgrade_registry::update_module_restrictions(
        &mut ref,
        &owner_cap,
        module1,
        vector[1u64],
        scenario.ctx(),
    );

    upgrade_registry::update_function_restrictions(
        &mut ref,
        &owner_cap,
        module2,
        function1,
        vector[2u64],
        scenario.ctx(),
    );

    // Test that restrictions are isolated
    assert!(!upgrade_registry::is_module_allowed(&ref, module1, 1u64));
    assert!(upgrade_registry::is_module_allowed(&ref, module2, 1u64));

    assert!(upgrade_registry::is_function_allowed(&ref, module1, function1, 2u64)); // blocked by module, but version 2 not blocked for module1
    assert!(!upgrade_registry::is_function_allowed(&ref, module2, function1, 2u64)); // blocked by function restriction
    assert!(upgrade_registry::is_function_allowed(&ref, module2, function2, 2u64)); // different function, not blocked

    tear_down_test(scenario, owner_cap, ref);
}
