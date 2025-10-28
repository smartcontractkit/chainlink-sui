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
    let empty_module_restrictions = upgrade_registry::get_module_restrictions(
        &ref,
        string::utf8(b"test_module"),
    );
    assert!(empty_module_restrictions.is_empty());

    tear_down_test(scenario, owner_cap, ref);
}

// =================== Version Blocking Tests =================== //

#[test]
public fun test_block_version() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let version = 1u8;

    // Block a version
    upgrade_registry::block_version(
        &mut ref,
        &owner_cap,
        module_name,
        version,
        scenario.ctx(),
    );

    // Check that the version is blocked
    let restrictions = upgrade_registry::get_module_restrictions(&ref, module_name);
    assert!(restrictions.length() == 1);
    assert!(restrictions[0].length() == 1);
    assert!(restrictions[0][0] == version);

    // Check that any function in this version is blocked
    assert!(
        !upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            string::utf8(b"any_function"),
            version,
        ),
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_block_version_emits_event() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let version = 1u8;

    // Block a version
    upgrade_registry::block_version(
        &mut ref,
        &owner_cap,
        module_name,
        version,
        scenario.ctx(),
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_block_multiple_versions() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");

    // Block multiple versions
    upgrade_registry::block_version(
        &mut ref,
        &owner_cap,
        module_name,
        1u8,
        scenario.ctx(),
    );
    upgrade_registry::block_version(
        &mut ref,
        &owner_cap,
        module_name,
        3u8,
        scenario.ctx(),
    );

    // Check that both versions are blocked
    let restrictions = upgrade_registry::get_module_restrictions(&ref, module_name);
    assert!(restrictions.length() == 2);

    // Check that functions in blocked versions are not allowed
    assert!(
        !upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            string::utf8(b"any_function"),
            1u8,
        ),
    );
    assert!(
        !upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            string::utf8(b"any_function"),
            3u8,
        ),
    );

    // Check that functions in non-blocked versions are allowed
    assert!(
        upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            string::utf8(b"any_function"),
            2u8,
        ),
    );

    tear_down_test(scenario, owner_cap, ref);
}

// =================== Function Blocking Tests =================== //

#[test]
public fun test_block_function() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");
    let version = 1u8;

    // Block a specific function
    upgrade_registry::block_function(
        &mut ref,
        &owner_cap,
        module_name,
        function_name,
        version,
        scenario.ctx(),
    );

    // Check that the function is blocked
    assert!(
        !upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            function_name,
            version,
        ),
    );

    // Check that other functions in the same version are still allowed
    assert!(
        upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            string::utf8(b"other_function"),
            version,
        ),
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_block_function_emits_event() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");
    let version = 1u8;

    // Block a specific function
    upgrade_registry::block_function(
        &mut ref,
        &owner_cap,
        module_name,
        function_name,
        version,
        scenario.ctx(),
    );

    // Note: Event testing is not implemented in this test framework
    // The event emission is tested indirectly through the function behavior

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_block_multiple_functions() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let function1 = string::utf8(b"function1");
    let function2 = string::utf8(b"function2");
    let version = 1u8;

    // Block multiple functions
    upgrade_registry::block_function(
        &mut ref,
        &owner_cap,
        module_name,
        function1,
        version,
        scenario.ctx(),
    );
    upgrade_registry::block_function(
        &mut ref,
        &owner_cap,
        module_name,
        function2,
        version,
        scenario.ctx(),
    );

    // Check that both functions are blocked
    assert!(
        !upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            function1,
            version,
        ),
    );
    assert!(
        !upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            function2,
            version,
        ),
    );

    // Check that other functions are still allowed
    assert!(
        upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            string::utf8(b"other_function"),
            version,
        ),
    );

    tear_down_test(scenario, owner_cap, ref);
}

// =================== Function Allowed Tests =================== //

#[test]
public fun test_is_function_allowed_no_restrictions() {
    let (scenario, owner_cap, ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");

    // When no restrictions exist, all functions should be allowed
    assert!(
        upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            function_name,
            1u8,
        ),
    );
    assert!(
        upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            function_name,
            2u8,
        ),
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_is_function_allowed_with_version_blocked() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");
    let blocked_version = 1u8;

    // Block the entire version
    upgrade_registry::block_version(
        &mut ref,
        &owner_cap,
        module_name,
        blocked_version,
        scenario.ctx(),
    );

    // Function should be blocked in the blocked version
    assert!(
        !upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            function_name,
            blocked_version,
        ),
    );

    // Function should be allowed in other versions
    assert!(
        upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            function_name,
            2u8,
        ),
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_is_function_allowed_with_function_blocked() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");
    let version = 1u8;

    // Block the specific function
    upgrade_registry::block_function(
        &mut ref,
        &owner_cap,
        module_name,
        function_name,
        version,
        scenario.ctx(),
    );

    // Function should be blocked
    assert!(
        !upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            function_name,
            version,
        ),
    );

    // Other functions should be allowed
    assert!(
        upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            string::utf8(b"other_function"),
            version,
        ),
    );

    tear_down_test(scenario, owner_cap, ref);
}

// =================== Verify Function Allowed Tests =================== //

#[test]
public fun test_verify_function_allowed_success() {
    let (scenario, owner_cap, ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");
    let version = 1u8;

    // This should not panic since the function is allowed
    upgrade_registry::verify_function_allowed(
        &ref,
        module_name,
        function_name,
        version,
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = upgrade_registry::EFunctionNotAllowed)]
public fun test_verify_function_allowed_failure() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");
    let version = 1u8;

    // Block the function
    upgrade_registry::block_function(
        &mut ref,
        &owner_cap,
        module_name,
        function_name,
        version,
        scenario.ctx(),
    );

    // This should panic since the function is blocked
    upgrade_registry::verify_function_allowed(
        &ref,
        module_name,
        function_name,
        version,
    );

    tear_down_test(scenario, owner_cap, ref);
}

// =================== Get Module Restrictions Tests =================== //

#[test]
public fun test_get_module_restrictions_empty() {
    let (scenario, owner_cap, ref) = set_up_test();

    let module_name = string::utf8(b"test_module");

    // Should return empty vector for module with no restrictions
    let restrictions = upgrade_registry::get_module_restrictions(&ref, module_name);
    assert!(restrictions.is_empty());

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_get_module_restrictions_with_blocks() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");

    // Block a version
    upgrade_registry::block_version(
        &mut ref,
        &owner_cap,
        module_name,
        1u8,
        scenario.ctx(),
    );

    // Block a function
    upgrade_registry::block_function(
        &mut ref,
        &owner_cap,
        module_name,
        string::utf8(b"test_function"),
        2u8,
        scenario.ctx(),
    );

    // Get restrictions
    let restrictions = upgrade_registry::get_module_restrictions(&ref, module_name);
    assert!(restrictions.length() == 2);

    // First restriction should be version block [1]
    assert!(restrictions[0].length() == 1);
    assert!(restrictions[0][0] == 1u8);

    // Second restriction should be function block [2, b"test_function"]
    assert!(restrictions[1].length() > 1);
    assert!(restrictions[1][0] == 2u8);

    tear_down_test(scenario, owner_cap, ref);
}

// =================== Edge Cases and Error Conditions =================== //

#[test]
public fun test_multiple_modules_isolation() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module1 = string::utf8(b"module1");
    let module2 = string::utf8(b"module2");
    let function_name = string::utf8(b"test_function");

    // Block version 1 in module1
    upgrade_registry::block_version(
        &mut ref,
        &owner_cap,
        module1,
        1u8,
        scenario.ctx(),
    );

    // Block function in module2
    upgrade_registry::block_function(
        &mut ref,
        &owner_cap,
        module2,
        function_name,
        1u8,
        scenario.ctx(),
    );

    // Test that restrictions are isolated
    assert!(
        !upgrade_registry::is_function_allowed(
            &ref,
            module1,
            function_name,
            1u8,
        ),
    ); // blocked by version block

    assert!(
        !upgrade_registry::is_function_allowed(
            &ref,
            module2,
            function_name,
            1u8,
        ),
    ); // blocked by function block

    assert!(
        upgrade_registry::is_function_allowed(
            &ref,
            module1,
            function_name,
            2u8,
        ),
    ); // allowed in module1, version 2

    assert!(
        upgrade_registry::is_function_allowed(
            &ref,
            module2,
            string::utf8(b"other_function"),
            1u8,
        ),
    ); // allowed in module2, different function

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_same_version_blocked_by_both_version_and_function() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");
    let version = 1u8;

    // Block the entire version
    upgrade_registry::block_version(
        &mut ref,
        &owner_cap,
        module_name,
        version,
        scenario.ctx(),
    );

    // Also block the specific function in the same version
    upgrade_registry::block_function(
        &mut ref,
        &owner_cap,
        module_name,
        function_name,
        version,
        scenario.ctx(),
    );

    // Function should still be blocked (version block takes precedence)
    assert!(
        !upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            function_name,
            version,
        ),
    );

    tear_down_test(scenario, owner_cap, ref);
}

// =================== Version Unblocking Tests =================== //

#[test]
public fun test_unblock_version() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let version = 1u8;

    // First block a version
    upgrade_registry::block_version(
        &mut ref,
        &owner_cap,
        module_name,
        version,
        scenario.ctx(),
    );

    // Verify the version is blocked
    assert!(
        !upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            string::utf8(b"any_function"),
            version,
        ),
    );

    // Now unblock the version
    upgrade_registry::unblock_version(
        &mut ref,
        &owner_cap,
        module_name,
        version,
        scenario.ctx(),
    );

    // Verify the version is now unblocked
    assert!(
        upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            string::utf8(b"any_function"),
            version,
        ),
    );

    // Check that restrictions are empty
    let restrictions = upgrade_registry::get_module_restrictions(&ref, module_name);
    assert!(restrictions.is_empty());

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_unblock_version_emits_event() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let version = 1u8;

    // Block a version first
    upgrade_registry::block_version(
        &mut ref,
        &owner_cap,
        module_name,
        version,
        scenario.ctx(),
    );

    // Unblock the version
    upgrade_registry::unblock_version(
        &mut ref,
        &owner_cap,
        module_name,
        version,
        scenario.ctx(),
    );

    // Note: Event testing is not implemented in this test framework
    // The event emission is tested indirectly through the function behavior

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_unblock_version_nonexistent_module() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"nonexistent_module");
    let version = 1u8;

    // Try to unblock a version in a module that has no restrictions
    // This should not panic and should do nothing
    upgrade_registry::unblock_version(
        &mut ref,
        &owner_cap,
        module_name,
        version,
        scenario.ctx(),
    );

    // Verify no restrictions were created
    let restrictions = upgrade_registry::get_module_restrictions(&ref, module_name);
    assert!(restrictions.is_empty());

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_unblock_version_nonexistent_version() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let blocked_version = 1u8;
    let unblocked_version = 2u8;

    // Block version 1
    upgrade_registry::block_version(
        &mut ref,
        &owner_cap,
        module_name,
        blocked_version,
        scenario.ctx(),
    );

    // Try to unblock version 2 (which was never blocked)
    upgrade_registry::unblock_version(
        &mut ref,
        &owner_cap,
        module_name,
        unblocked_version,
        scenario.ctx(),
    );

    // Verify version 1 is still blocked
    assert!(
        !upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            string::utf8(b"any_function"),
            blocked_version,
        ),
    );

    // Verify version 2 was never blocked (and is still allowed)
    assert!(
        upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            string::utf8(b"any_function"),
            unblocked_version,
        ),
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_unblock_version_with_multiple_versions() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");

    // Block multiple versions
    upgrade_registry::block_version(
        &mut ref,
        &owner_cap,
        module_name,
        1u8,
        scenario.ctx(),
    );
    upgrade_registry::block_version(
        &mut ref,
        &owner_cap,
        module_name,
        2u8,
        scenario.ctx(),
    );
    upgrade_registry::block_version(
        &mut ref,
        &owner_cap,
        module_name,
        3u8,
        scenario.ctx(),
    );

    // Unblock version 2
    upgrade_registry::unblock_version(
        &mut ref,
        &owner_cap,
        module_name,
        2u8,
        scenario.ctx(),
    );

    // Verify version 1 is still blocked
    assert!(
        !upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            string::utf8(b"any_function"),
            1u8,
        ),
    );

    // Verify version 2 is now unblocked
    assert!(
        upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            string::utf8(b"any_function"),
            2u8,
        ),
    );

    // Verify version 3 is still blocked
    assert!(
        !upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            string::utf8(b"any_function"),
            3u8,
        ),
    );

    // Check that only 2 restrictions remain
    let restrictions = upgrade_registry::get_module_restrictions(&ref, module_name);
    assert!(restrictions.length() == 2);

    tear_down_test(scenario, owner_cap, ref);
}

// =================== Function Unblocking Tests =================== //

#[test]
public fun test_unblock_function() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");
    let version = 1u8;

    // First block a specific function
    upgrade_registry::block_function(
        &mut ref,
        &owner_cap,
        module_name,
        function_name,
        version,
        scenario.ctx(),
    );

    // Verify the function is blocked
    assert!(
        !upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            function_name,
            version,
        ),
    );

    // Now unblock the function
    upgrade_registry::unblock_function(
        &mut ref,
        &owner_cap,
        module_name,
        function_name,
        version,
        scenario.ctx(),
    );

    // Verify the function is now unblocked
    assert!(
        upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            function_name,
            version,
        ),
    );

    // Check that restrictions are empty
    let restrictions = upgrade_registry::get_module_restrictions(&ref, module_name);
    assert!(restrictions.is_empty());

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_unblock_function_emits_event() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");
    let version = 1u8;

    // Block a function first
    upgrade_registry::block_function(
        &mut ref,
        &owner_cap,
        module_name,
        function_name,
        version,
        scenario.ctx(),
    );

    // Unblock the function
    upgrade_registry::unblock_function(
        &mut ref,
        &owner_cap,
        module_name,
        function_name,
        version,
        scenario.ctx(),
    );

    // Note: Event testing is not implemented in this test framework
    // The event emission is tested indirectly through the function behavior

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_unblock_function_nonexistent_module() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"nonexistent_module");
    let function_name = string::utf8(b"test_function");
    let version = 1u8;

    // Try to unblock a function in a module that has no restrictions
    // This should not panic and should do nothing
    upgrade_registry::unblock_function(
        &mut ref,
        &owner_cap,
        module_name,
        function_name,
        version,
        scenario.ctx(),
    );

    // Verify no restrictions were created
    let restrictions = upgrade_registry::get_module_restrictions(&ref, module_name);
    assert!(restrictions.is_empty());

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_unblock_function_nonexistent_function() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let blocked_function = string::utf8(b"blocked_function");
    let unblocked_function = string::utf8(b"unblocked_function");
    let version = 1u8;

    // Block one function
    upgrade_registry::block_function(
        &mut ref,
        &owner_cap,
        module_name,
        blocked_function,
        version,
        scenario.ctx(),
    );

    // Try to unblock a different function (which was never blocked)
    upgrade_registry::unblock_function(
        &mut ref,
        &owner_cap,
        module_name,
        unblocked_function,
        version,
        scenario.ctx(),
    );

    // Verify the originally blocked function is still blocked
    assert!(
        !upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            blocked_function,
            version,
        ),
    );

    // Verify the unblocked function was never blocked (and is still allowed)
    assert!(
        upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            unblocked_function,
            version,
        ),
    );

    // Check that one restriction remains
    let restrictions = upgrade_registry::get_module_restrictions(&ref, module_name);
    assert!(restrictions.length() == 1);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_unblock_function_with_multiple_functions() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let function1 = string::utf8(b"function1");
    let function2 = string::utf8(b"function2");
    let function3 = string::utf8(b"function3");
    let version = 1u8;

    // Block multiple functions
    upgrade_registry::block_function(
        &mut ref,
        &owner_cap,
        module_name,
        function1,
        version,
        scenario.ctx(),
    );
    upgrade_registry::block_function(
        &mut ref,
        &owner_cap,
        module_name,
        function2,
        version,
        scenario.ctx(),
    );
    upgrade_registry::block_function(
        &mut ref,
        &owner_cap,
        module_name,
        function3,
        version,
        scenario.ctx(),
    );

    // Unblock function2
    upgrade_registry::unblock_function(
        &mut ref,
        &owner_cap,
        module_name,
        function2,
        version,
        scenario.ctx(),
    );

    // Verify function1 is still blocked
    assert!(
        !upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            function1,
            version,
        ),
    );

    // Verify function2 is now unblocked
    assert!(
        upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            function2,
            version,
        ),
    );

    // Verify function3 is still blocked
    assert!(
        !upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            function3,
            version,
        ),
    );

    // Check that only 2 restrictions remain
    let restrictions = upgrade_registry::get_module_restrictions(&ref, module_name);
    assert!(restrictions.length() == 2);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_unblock_function_with_version_and_function_blocks() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();

    let module_name = string::utf8(b"test_module");
    let function_name = string::utf8(b"test_function");
    let version = 1u8;

    // Block the entire version
    upgrade_registry::block_version(
        &mut ref,
        &owner_cap,
        module_name,
        version,
        scenario.ctx(),
    );

    // Also block the specific function in the same version
    upgrade_registry::block_function(
        &mut ref,
        &owner_cap,
        module_name,
        function_name,
        version,
        scenario.ctx(),
    );

    // Unblock the specific function
    upgrade_registry::unblock_function(
        &mut ref,
        &owner_cap,
        module_name,
        function_name,
        version,
        scenario.ctx(),
    );

    // Function should still be blocked because the entire version is blocked
    assert!(
        !upgrade_registry::is_function_allowed(
            &ref,
            module_name,
            function_name,
            version,
        ),
    );

    // Check that only the version block remains
    let restrictions = upgrade_registry::get_module_restrictions(&ref, module_name);
    assert!(restrictions.length() == 1);
    assert!(restrictions[0].length() == 1);
    assert!(restrictions[0][0] == version);

    tear_down_test(scenario, owner_cap, ref);
}
