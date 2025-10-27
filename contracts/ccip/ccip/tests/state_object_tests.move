#[test_only]
module ccip::state_object_test;

use ccip::ownable::{Self, OwnerCap};
use ccip::state_object::{Self, CCIPObjectRef};
use mcms::mcms_account;
use mcms::mcms_deployer;
use mcms::mcms_registry::{Self, Registry};
use std::string;
use sui::address;
use sui::bcs;
use sui::test_scenario::{Self, Scenario};

const SENDER_1: address = @0x1;
const SENDER_2: address = @0x2;

fun set_up_test(): (Scenario, OwnerCap, CCIPObjectRef, TestObject) {
    let mut scenario = test_scenario::begin(SENDER_1);
    let ctx = scenario.ctx();

    state_object::test_init(ctx);

    // Advance to next transaction to retrieve the created objects
    scenario.next_tx(SENDER_1);

    // Retrieve the OwnerCap that was transferred to SENDER_1
    let owner_cap = scenario.take_from_sender<OwnerCap>();

    // Retrieve the shared CCIPObjectRef
    let ref = scenario.take_shared<CCIPObjectRef>();

    let obj = TestObject {
        id: object::new(scenario.ctx()),
    };
    (scenario, owner_cap, ref, obj)
}

fun tear_down_test(scenario: Scenario, owner_cap: OwnerCap, ref: CCIPObjectRef) {
    // Return the owner cap back to the sender instead of destroying it
    test_scenario::return_to_sender(&scenario, owner_cap);
    // Return the shared object back to the scenario instead of destroying it
    test_scenario::return_shared(ref);
    test_scenario::end(scenario);
}

public struct TestObject has key, store {
    id: UID,
}

#[test]
public fun test_add() {
    let (mut scenario, owner_cap, mut ref, obj) = set_up_test();
    let ctx = scenario.ctx();

    state_object::add(&mut ref, &owner_cap, obj, ctx);
    assert!(state_object::contains<TestObject>(&ref));

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_remove() {
    let (mut scenario, owner_cap, mut ref, obj) = set_up_test();
    let ctx = scenario.ctx();

    state_object::add(&mut ref, &owner_cap, obj, ctx);
    assert!(state_object::contains<TestObject>(&ref));

    let obj2: TestObject = state_object::remove<TestObject>(&mut ref, &owner_cap, ctx);
    assert!(!state_object::contains<TestObject>(&ref));

    let TestObject { id } = obj2;
    object::delete(id);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_borrow() {
    let (mut scenario, owner_cap, mut ref, obj) = set_up_test();
    let ctx = scenario.ctx();

    state_object::add(&mut ref, &owner_cap, obj, ctx);
    assert!(state_object::contains<TestObject>(&ref));

    let _obj2: &TestObject = state_object::borrow<TestObject>(&ref);
    assert!(state_object::contains<TestObject>(&ref));

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_borrow_mut() {
    let (mut scenario, owner_cap, mut ref, obj) = set_up_test();
    let ctx = scenario.ctx();

    state_object::add(&mut ref, &owner_cap, obj, ctx);
    assert!(state_object::contains<TestObject>(&ref));

    let _obj2: &mut TestObject = state_object::borrow_mut<TestObject>(&mut ref);
    assert!(state_object::contains<TestObject>(&ref));

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_transfer_ownership() {
    let (mut scenario, owner_cap, mut ref, obj) = set_up_test();
    let ctx = scenario.ctx();

    state_object::add(&mut ref, &owner_cap, obj, ctx);

    let ctx = scenario.ctx();
    let new_owner = SENDER_2;
    state_object::transfer_ownership(&mut ref, &owner_cap, new_owner, ctx);

    let (from, to, accepted) = state_object::pending_transfer(&ref);
    assert!(from == SENDER_1);
    assert!(to == new_owner);
    assert!(!accepted);

    // after transfer, the owner is still the original owner
    let owner = state_object::owner(&ref);
    assert!(owner == SENDER_1);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_accept_and_execute_ownership() {
    let (mut scenario_1, owner_cap, mut ref, obj) = set_up_test();
    let ctx_1 = scenario_1.ctx();
    state_object::add(&mut ref, &owner_cap, obj, ctx_1);

    // tx 1: SENDER_1 transfer ownership to SENDER_2
    // let ctx_1 = scenario_1.ctx();
    let new_owner = SENDER_2;
    state_object::transfer_ownership(&mut ref, &owner_cap, new_owner, ctx_1);
    let (from, to, accepted) = state_object::pending_transfer(&ref);
    assert!(from == SENDER_1);
    assert!(to == new_owner);
    assert!(!accepted);

    test_scenario::end(scenario_1);

    // tx 2: SENDER_2 accepts the ownership transfer
    let mut scenario_2 = test_scenario::begin(new_owner);
    // let accept_cap = test_scenario::take_from_address<state_object::AcceptCap>(&scenario_2, new_owner);
    let ctx_2 = scenario_2.ctx();

    state_object::accept_ownership(&mut ref, ctx_2);
    let (from, to, accepted) = state_object::pending_transfer(&ref);
    assert!(from == SENDER_1);
    assert!(to == new_owner);
    assert!(accepted);
    // after accept, the owner is still the original owner
    let owner_1 = state_object::owner(&ref);
    assert!(owner_1 == SENDER_1);

    test_scenario::end(scenario_2);

    // tx 3: SENDER_1 executes the ownership transfer
    let mut scenario_3 = test_scenario::begin(SENDER_1);
    let ctx_3 = scenario_3.ctx();
    state_object::execute_ownership_transfer(&mut ref, owner_cap, new_owner, ctx_3);
    test_scenario::end(scenario_3);

    let (from, to, accepted) = state_object::pending_transfer(&ref);
    assert!(from == @0x0);
    assert!(to == @0x0);
    assert!(!accepted);
    // after execute, the owner is the new owner
    let owner_2 = state_object::owner(&ref);
    assert!(owner_2 == SENDER_2);

    // tx 4: SENDER_2 can now update the state object
    let mut scenario_4 = test_scenario::begin(SENDER_2);
    let owner_cap_2 = scenario_4.take_from_sender<OwnerCap>();

    let obj2: TestObject = state_object::remove<TestObject>(
        &mut ref,
        &owner_cap_2,
        scenario_4.ctx(),
    );
    assert!(!state_object::contains<TestObject>(&ref));
    let TestObject { id } = obj2;
    object::delete(id);

    // Special cleanup for this test - ownership was transferred, so we transfer owner_cap_2 to dummy address
    transfer::public_transfer(owner_cap_2, @0x0);
    test_scenario::return_shared(ref);
    test_scenario::end(scenario_4);
}

#[test]
#[expected_failure(abort_code = state_object::EInvalidOwnerCap)]
public fun test_add_package_id_with_invalid_owner_cap() {
    let (mut scenario, owner_cap, mut ref, mut obj) = set_up_test();
    let ctx = scenario.ctx();

    // Create a different owner cap using ownable::new
    let (ownable_state, fake_owner_cap) = ownable::new(&mut obj.id, ctx);

    // Try to add a package ID using the fake owner cap - this should fail at line 60
    // The test should abort here with EInvalidOwnerCap
    state_object::add_package_id(&mut ref, &fake_owner_cap, @0x123);

    // This code should never be reached due to the expected failure above
    // Cleanup
    let TestObject { id } = obj;
    object::delete(id);
    transfer::public_transfer(fake_owner_cap, @0x0);

    // Cleanup the original scenario objects
    ownable::destroy(ownable_state, owner_cap, ctx);
    test_scenario::return_shared(ref);
    test_scenario::end(scenario);
}

// ================================================================
// |                  MCMS Integration Tests                      |
// ================================================================

const OWNER: address = @0x123;

fun setup_with_mcms_ownership(): (Scenario, Registry, CCIPObjectRef) {
    let mut scenario = test_scenario::begin(OWNER);

    // Initialize MCMS components
    {
        let ctx = scenario.ctx();
        mcms_account::test_init(ctx);
        mcms_registry::test_init(ctx);
        mcms_deployer::test_init(ctx);
        state_object::test_init(ctx);
    };

    scenario.next_tx(OWNER);

    let mut registry = test_scenario::take_shared<Registry>(&scenario);
    let mut ref = test_scenario::take_shared<CCIPObjectRef>(&scenario);
    let owner_cap = test_scenario::take_from_sender<OwnerCap>(&scenario);

    // Transfer ownership to MCMS (3-step process)
    // Step 1: transfer_ownership
    state_object::transfer_ownership(
        &mut ref,
        &owner_cap,
        mcms_registry::get_multisig_address(),
        scenario.ctx(),
    );

    // Step 2: accept_ownership
    scenario.next_tx(mcms_registry::get_multisig_address());
    state_object::accept_ownership(&mut ref, scenario.ctx());

    // Step 3: execute_ownership_transfer (registers with MCMS)
    state_object::execute_ownership_transfer_to_mcms(
        &mut ref,
        owner_cap,
        &mut registry,
        mcms_registry::get_multisig_address(),
        scenario.ctx(),
    );

    scenario.next_tx(OWNER);

    (scenario, registry, ref)
}

#[test]
fun test_mcms_add_allowed_modules_success() {
    let (mut scenario, mut registry, ref) = setup_with_mcms_ownership();

    // Verify initial allowed modules (should have fee_quoter, rmn_remote, state_object, token_admin_registry)
    let initial_modules = mcms_registry::get_allowed_modules(
        &registry,
        address::to_ascii_string(@ccip),
    );
    assert!(initial_modules.contains(&b"fee_quoter"), 0);
    assert!(initial_modules.contains(&b"rmn_remote"), 1);
    assert!(initial_modules.contains(&b"state_object"), 2);
    assert!(initial_modules.contains(&b"token_admin_registry"), 3);
    assert!(!initial_modules.contains(&b"nonce_manager"), 4); // Should not exist yet

    // Prepare data for mcms_add_allowed_modules
    // Data format: [registry_address][vector<vector<u8>> of module names]
    let mut data = vector::empty<u8>();
    data.append(bcs::to_bytes(&object::id_address(&registry))); // Registry address for validation

    // Serialize vector of module names (vector<vector<u8>>)
    let module_names = vector[b"nonce_manager"];
    data.append(bcs::to_bytes(&module_names));

    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(b"state_object"),
        string::utf8(b"add_allowed_modules"),
        data,
        x"0000000000000000000000000000000000000000000000000000000000000001", // batch_id
        0, // sequence_number
        1, // total_in_batch
    );

    state_object::mcms_add_allowed_modules(
        &mut registry,
        params,
        scenario.ctx(),
    );

    // Verify nonce_manager was added
    let updated_modules = mcms_registry::get_allowed_modules(
        &registry,
        address::to_ascii_string(@ccip),
    );
    assert!(updated_modules.contains(&b"nonce_manager"), 5); // Should exist now

    // Cleanup
    test_scenario::return_shared(registry);
    test_scenario::return_shared(ref);
    test_scenario::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms_registry::EModuleAlreadyAllowed)]
fun test_mcms_add_allowed_modules_already_exists() {
    let (mut scenario, mut registry, ref) = setup_with_mcms_ownership();

    // Try to add "fee_quoter" which already exists in initial allowed modules
    // Data format: [registry_address][vector<vector<u8>> of module names]
    let mut data = vector::empty<u8>();
    data.append(bcs::to_bytes(&object::id_address(&registry))); // Registry address for validation

    // Serialize vector of module names (vector<vector<u8>>)
    let module_names = vector[b"fee_quoter"]; // Already exists
    data.append(bcs::to_bytes(&module_names));

    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(b"state_object"),
        string::utf8(b"add_allowed_modules"),
        data,
        x"0000000000000000000000000000000000000000000000000000000000000001",
        0,
        1,
    );

    // This should fail with EModuleAlreadyAllowed
    state_object::mcms_add_allowed_modules(
        &mut registry,
        params,
        scenario.ctx(),
    );

    // Cleanup (won't be reached due to expected failure)
    test_scenario::return_shared(registry);
    test_scenario::return_shared(ref);
    test_scenario::end(scenario);
}

#[test]
#[expected_failure(abort_code = state_object::EInvalidFunction)]
fun test_mcms_add_allowed_modules_wrong_function_name() {
    let (mut scenario, mut registry, ref) = setup_with_mcms_ownership();

    // Prepare data with correct format
    // Data format: [registry_address][vector<vector<u8>> of module names]
    let mut data = vector::empty<u8>();
    data.append(bcs::to_bytes(&object::id_address(&registry))); // Registry address for validation

    // Serialize vector of module names (vector<vector<u8>>)
    let module_names = vector[b"new_module"];
    data.append(bcs::to_bytes(&module_names));

    // But use wrong function name in params
    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(b"state_object"),
        string::utf8(b"wrong_function"), // Wrong function name!
        data,
        x"0000000000000000000000000000000000000000000000000000000000000001",
        0,
        1,
    );

    // This should fail with EInvalidFunction
    state_object::mcms_add_allowed_modules(
        &mut registry,
        params,
        scenario.ctx(),
    );

    // Cleanup (won't be reached due to expected failure)
    test_scenario::return_shared(registry);
    test_scenario::return_shared(ref);
    test_scenario::end(scenario);
}

// ================================================================
// |         MCMS Remove Allowed Modules Tests                   |
// ================================================================

#[test]
fun test_mcms_remove_allowed_modules_success() {
    let (mut scenario, mut registry, ref) = setup_with_mcms_ownership();

    // First, add a module that we'll later remove
    {
        let mut data = vector::empty<u8>();
        data.append(bcs::to_bytes(&object::id_address(&registry)));
        let module_names = vector[b"nonce_manager"];
        data.append(bcs::to_bytes(&module_names));

        let params = mcms_registry::test_create_executing_callback_params(
            @ccip,
            string::utf8(b"state_object"),
            string::utf8(b"add_allowed_modules"),
            data,
            x"0000000000000000000000000000000000000000000000000000000000000001",
            0,
            1,
        );

        state_object::mcms_add_allowed_modules(
            &mut registry,
            params,
            scenario.ctx(),
        );
    };

    // Verify nonce_manager was added
    let modules_before = mcms_registry::get_allowed_modules(
        &registry,
        address::to_ascii_string(@ccip),
    );
    assert!(modules_before.contains(&b"nonce_manager"), 0);

    // Now remove the module
    {
        let mut data = vector::empty<u8>();
        data.append(bcs::to_bytes(&object::id_address(&registry)));
        let module_names = vector[b"nonce_manager"];
        data.append(bcs::to_bytes(&module_names));

        let params = mcms_registry::test_create_executing_callback_params(
            @ccip,
            string::utf8(b"state_object"),
            string::utf8(b"remove_allowed_modules"),
            data,
            x"0000000000000000000000000000000000000000000000000000000000000002",
            0,
            1,
        );

        state_object::mcms_remove_allowed_modules(
            &mut registry,
            params,
            scenario.ctx(),
        );
    };

    // Verify nonce_manager was removed
    let modules_after = mcms_registry::get_allowed_modules(
        &registry,
        address::to_ascii_string(@ccip),
    );
    assert!(!modules_after.contains(&b"nonce_manager"), 1);

    // Cleanup
    test_scenario::return_shared(registry);
    test_scenario::return_shared(ref);
    test_scenario::end(scenario);
}

#[test]
#[expected_failure(abort_code = mcms_registry::EModuleNotInAllowlist)]
fun test_mcms_remove_allowed_modules_not_in_allowlist() {
    let (mut scenario, mut registry, ref) = setup_with_mcms_ownership();

    // Try to remove a module that doesn't exist
    let mut data = vector::empty<u8>();
    data.append(bcs::to_bytes(&object::id_address(&registry)));
    let module_names = vector[b"nonexistent_module"];
    data.append(bcs::to_bytes(&module_names));

    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(b"state_object"),
        string::utf8(b"remove_allowed_modules"),
        data,
        x"0000000000000000000000000000000000000000000000000000000000000001",
        0,
        1,
    );

    // This should fail with EModuleNotInAllowlist
    state_object::mcms_remove_allowed_modules(
        &mut registry,
        params,
        scenario.ctx(),
    );

    // Cleanup (won't be reached due to expected failure)
    test_scenario::return_shared(registry);
    test_scenario::return_shared(ref);
    test_scenario::end(scenario);
}

#[test]
#[expected_failure(abort_code = state_object::EInvalidFunction)]
fun test_mcms_remove_allowed_modules_wrong_function_name() {
    let (mut scenario, mut registry, ref) = setup_with_mcms_ownership();

    // Prepare data with correct format
    let mut data = vector::empty<u8>();
    data.append(bcs::to_bytes(&object::id_address(&registry)));
    let module_names = vector[b"fee_quoter"];
    data.append(bcs::to_bytes(&module_names));

    // But use wrong function name in params
    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(b"state_object"),
        string::utf8(b"wrong_function"), // Wrong function name!
        data,
        x"0000000000000000000000000000000000000000000000000000000000000001",
        0,
        1,
    );

    // This should fail with EInvalidFunction
    state_object::mcms_remove_allowed_modules(
        &mut registry,
        params,
        scenario.ctx(),
    );

    // Cleanup (won't be reached due to expected failure)
    test_scenario::return_shared(registry);
    test_scenario::return_shared(ref);
    test_scenario::end(scenario);
}

// ================================================================
// |         MCMS 3-Step Ownership Transfer Test                 |
// ================================================================

#[test]
fun test_mcms_three_step_ownership_transfer() {
    let (mut scenario, mut registry, mut ref) = setup_with_mcms_ownership();

    // At this point, ownership is with MCMS
    let initial_owner = state_object::owner(&ref);
    assert!(initial_owner == mcms_registry::get_multisig_address());

    let new_owner = SENDER_2;
    scenario.next_tx(OWNER);

    // Step 1: MCMS calls mcms_transfer_ownership to initiate transfer to SENDER_2
    {
        let owner_cap_address = mcms_registry::test_get_cap_address<OwnerCap>(
            &registry,
            @ccip.to_ascii_string(),
        );

        let mut data = vector::empty<u8>();
        data.append(bcs::to_bytes(&object::id_address(&ref)));
        data.append(bcs::to_bytes(&owner_cap_address));
        data.append(bcs::to_bytes(&new_owner));

        let params = mcms_registry::test_create_executing_callback_params(
            @ccip,
            string::utf8(b"state_object"),
            string::utf8(b"transfer_ownership"),
            data,
            x"0000000000000000000000000000000000000000000000000000000000000001", // batch_id
            0, // sequence_number
            1, // total_in_batch
        );

        state_object::mcms_transfer_ownership(
            &mut ref,
            &mut registry,
            params,
            scenario.ctx(),
        );
    };

    // Verify pending transfer was created
    let (from, to, accepted) = state_object::pending_transfer(&ref);
    assert!(from == mcms_registry::get_multisig_address());
    assert!(to == new_owner);
    assert!(!accepted);

    // Owner should still be MCMS
    assert!(state_object::owner(&ref) == mcms_registry::get_multisig_address());

    // Step 2: SENDER_2 (new owner) accepts the ownership transfer
    scenario.next_tx(new_owner);
    {
        state_object::accept_ownership(&mut ref, scenario.ctx());
    };

    // Verify pending transfer was accepted
    let (from2, to2, accepted2) = state_object::pending_transfer(&ref);
    assert!(from2 == mcms_registry::get_multisig_address());
    assert!(to2 == new_owner);
    assert!(accepted2); // Now it's accepted

    // Owner should still be MCMS (not yet executed)
    assert!(state_object::owner(&ref) == mcms_registry::get_multisig_address());

    // Step 3: MCMS calls mcms_execute_ownership_transfer to finalize
    scenario.next_tx(OWNER);
    {
        let owner_cap_address = mcms_registry::test_get_cap_address<OwnerCap>(
            &registry,
            @ccip.to_ascii_string(),
        );
        let mut deployer_state = test_scenario::take_shared<mcms_deployer::DeployerState>(&scenario);

        // Serialize data: [ref_address][owner_cap_address][to_address][package_address]
        let mut data = vector::empty<u8>();
        data.append(bcs::to_bytes(&object::id_address(&ref)));
        data.append(bcs::to_bytes(&owner_cap_address));
        data.append(bcs::to_bytes(&new_owner));
        data.append(bcs::to_bytes(&@ccip));

        let params = mcms_registry::test_create_executing_callback_params(
            @ccip,
            string::utf8(b"state_object"),
            string::utf8(b"execute_ownership_transfer"),
            data,
            x"0000000000000000000000000000000000000000000000000000000000000002", // different batch_id
            0, // sequence_number
            1, // total_in_batch
        );

        state_object::mcms_execute_ownership_transfer(
            &mut ref,
            &mut registry,
            &mut deployer_state,
            params,
            scenario.ctx(),
        );

        test_scenario::return_shared(deployer_state);
    };

    // Verify pending transfer was cleared
    let (from3, to3, accepted3) = state_object::pending_transfer(&ref);
    assert!(from3 == @0x0);
    assert!(to3 == @0x0);
    assert!(!accepted3);

    // Owner should now be SENDER_2
    assert!(state_object::owner(&ref) == new_owner);

    // Step 4: Verify SENDER_2 received the OwnerCap and can use it
    scenario.next_tx(new_owner);
    {
        let owner_cap = scenario.take_from_sender<OwnerCap>();

        // Verify SENDER_2 can perform operations with the OwnerCap
        let test_obj = TestObject {
            id: object::new(scenario.ctx()),
        };
        state_object::add(&mut ref, &owner_cap, test_obj, scenario.ctx());
        assert!(state_object::contains<TestObject>(&ref));

        scenario.return_to_sender(owner_cap);
    };

    test_scenario::return_shared(registry);
    test_scenario::return_shared(ref);
    test_scenario::end(scenario);
}
