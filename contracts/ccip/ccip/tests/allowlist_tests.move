#[test_only]
module ccip::allowlist_tests;

use ccip::allowlist;
use sui::test_scenario;

fun set_up_test(allowlist: vector<address>, ctx: &mut TxContext): allowlist::AllowlistState {
    allowlist::new(allowlist, ctx)
}

#[test]
public fun init_empty_is_empty_and_disabled() {
    let mut scenario = test_scenario::begin(@0x1);

    let state = set_up_test(vector::empty(), scenario.ctx());

    assert!(!allowlist::get_allowlist_enabled(&state));
    assert!(allowlist::get_allowlist(&state).is_empty());

    // Any address is allowed when the allowlist is disabled
    assert!(allowlist::is_allowed(&state, @0x1111111111111));

    allowlist::destroy_allowlist(state);

    test_scenario::end(scenario);
}

#[test]
public fun init_non_empty_is_non_empty_and_enabled() {
    let mut scenario = test_scenario::begin(@0x1);

    let init_allowlist = vector[@0x1, @0x2];

    let state = set_up_test(init_allowlist, scenario.ctx());

    assert!(allowlist::get_allowlist_enabled(&state));
    assert!(allowlist::get_allowlist(&state).length() == 2);

    // The given addresses are allowed
    assert!(allowlist::is_allowed(&state, init_allowlist[0]));
    assert!(allowlist::is_allowed(&state, init_allowlist[1]));

    assert!(!allowlist::is_allowed(&state, @0x3));

    allowlist::destroy_allowlist(state);

    test_scenario::end(scenario);
}

#[test]
#[expected_failure(abort_code = allowlist::EAllowlistNotEnabled, location = allowlist)]
public fun cannot_add_to_disabled_allowlist() {
    let mut scenario = test_scenario::begin(@0x1);

    let mut state = set_up_test(vector::empty(), scenario.ctx());

    let adds = vector[@0x1];

    allowlist::apply_allowlist_updates(&mut state, vector::empty(), adds);

    allowlist::destroy_allowlist(state);

    test_scenario::end(scenario);
}

#[test]
public fun apply_allowlist_updates_mutates_state() {
    let mut scenario = test_scenario::begin(@0x1);

    let mut state = set_up_test(vector::empty(), scenario.ctx());

    allowlist::set_allowlist_enabled(&mut state, true);

    assert!(allowlist::get_allowlist(&state).is_empty());

    allowlist::apply_allowlist_updates(&mut state, vector::empty(), vector::empty());

    assert!(allowlist::get_allowlist(&state).is_empty());

    let adds = vector[@0x1, @0x2];

    allowlist::apply_allowlist_updates(&mut state, vector::empty(), adds);

    let removes = vector[@0x1];

    allowlist::apply_allowlist_updates(&mut state, removes, vector::empty());

    assert!(allowlist::get_allowlist(&state).length() == 1);
    assert!(allowlist::is_allowed(&state, @0x2));
    assert!(!allowlist::is_allowed(&state, @0x1));

    allowlist::destroy_allowlist(state);

    test_scenario::end(scenario);
}

#[test]
public fun apply_allowlist_updates_removes_before_adds() {
    let mut scenario = test_scenario::begin(@0x1);

    let mut state = set_up_test(vector::empty(), scenario.ctx());
    let account_to_allow = @0x1;

    allowlist::set_allowlist_enabled(&mut state, true);

    let adds_and_removes = vector[account_to_allow];

    allowlist::apply_allowlist_updates(&mut state, vector::empty(), adds_and_removes);

    assert!(allowlist::get_allowlist(&state).length() == 1);
    assert!(allowlist::is_allowed(&state, account_to_allow));

    allowlist::apply_allowlist_updates(&mut state, adds_and_removes, adds_and_removes);

    // Since removes happen before adds, the account should still be allowed
    assert!(allowlist::is_allowed(&state, account_to_allow));

    allowlist::destroy_allowlist(state);

    test_scenario::end(scenario);
}
