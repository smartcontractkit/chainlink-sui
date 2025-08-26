#[test_only]
module ccip::state_object_test;

use ccip::ownable::OwnerCap;
use ccip::state_object::{Self, CCIPObjectRef};
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
