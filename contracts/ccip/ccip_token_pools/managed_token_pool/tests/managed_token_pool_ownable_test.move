#[test_only]
module managed_token_pool::managed_token_pool_ownable_test;

use ccip_token_pool::ownable;

#[test]
public fun test_ownable_module_functions_exist() {
    // Test that we can create ownable state and cap
    let mut scenario = sui::test_scenario::begin(@0x123);
    let (ownable_state, owner_cap) = ownable::new(scenario.ctx());

    // Test query functions
    let _owner = ownable::owner(&ownable_state);
    let _owner_cap_id = ownable::owner_cap_id(&ownable_state);
    let _has_pending = ownable::has_pending_transfer(&ownable_state);
    let _pending_from = ownable::pending_transfer_from(&ownable_state);
    let _pending_to = ownable::pending_transfer_to(&ownable_state);
    let _pending_accepted = ownable::pending_transfer_accepted(&ownable_state);

    // Clean up
    ownable::destroy_ownable_state(ownable_state, scenario.ctx());
    ownable::destroy_owner_cap(owner_cap, scenario.ctx());
    scenario.end();
}
