#[test_only]
module ccip::nonce_manager_test;

use ccip::nonce_manager::{Self, NonceManagerState};
use ccip::state_object::{Self, CCIPObjectRef, OwnerCap};
use sui::test_scenario::{Self, Scenario};

const SENDER: address = @0x1;

fun set_up_test(): (Scenario, CCIPObjectRef, OwnerCap) {
    let mut scenario = test_scenario::begin(SENDER);
    let ctx = scenario.ctx();

    state_object::test_init(ctx);

    scenario.next_tx(SENDER);

    let ref = scenario.take_shared<CCIPObjectRef>();
    let owner_cap = scenario.take_from_address(SENDER);

    (scenario, ref, owner_cap)
}

fun initialize(ref: &mut CCIPObjectRef, owner_cap: &OwnerCap, ctx: &mut TxContext) {
    nonce_manager::initialize(ref, owner_cap, ctx);
}

fun tear_down_test(scenario: Scenario, ref: CCIPObjectRef, owner_cap: OwnerCap) {
    state_object::destroy_state_object(ref);
    test_scenario::return_to_address(SENDER, owner_cap);
    test_scenario::end(scenario);
}

#[test]
public fun test_initialize() {
    let (mut scenario, mut ref, owner_cap) = set_up_test();
    let ctx = scenario.ctx();
    initialize(&mut ref, &owner_cap, ctx);

    let _state = state_object::borrow<NonceManagerState>(&ref);

    assert!(state_object::contains<NonceManagerState>(&ref));

    tear_down_test(scenario, ref, owner_cap);
}

#[test]
public fun test_get_incremented_outbound_nonce() {
    let (mut scenario, mut ref, owner_cap) = set_up_test();
    initialize(&mut ref, &owner_cap, scenario.ctx());

    let mut nonce = nonce_manager::get_outbound_nonce(&ref, 1, @0x1);
    assert!(nonce == 0);

    scenario.next_tx(SENDER);
    let nonce_manager_cap = scenario.take_from_address(SENDER);
    let mut incremented_nonce = nonce_manager::get_incremented_outbound_nonce(
        &mut ref,
        &nonce_manager_cap,
        1,
        @0x1,
        scenario.ctx(),
    );
    assert!(incremented_nonce == 1);

    nonce = nonce_manager::get_outbound_nonce(&ref, 1, @0x1);
    assert!(nonce == 1);

    incremented_nonce =
        nonce_manager::get_incremented_outbound_nonce(
            &mut ref,
            &nonce_manager_cap,
            1,
            @0x1,
            scenario.ctx(),
        );
    assert!(incremented_nonce == 2);

    nonce = nonce_manager::get_outbound_nonce(&ref, 1, @0x1);
    assert!(nonce == 2);

    tear_down_test(scenario, ref, owner_cap);
    test_scenario::return_to_address(SENDER, nonce_manager_cap);
}
