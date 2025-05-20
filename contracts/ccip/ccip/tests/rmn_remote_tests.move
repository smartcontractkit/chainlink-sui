#[test_only]
module ccip::rmn_remote_test;

use ccip::state_object::{Self, OwnerCap, CCIPObjectRef};
use ccip::rmn_remote::{Self, RMNRemoteState};
use sui::test_scenario::{Self, Scenario};

fun set_up_test(): (Scenario, OwnerCap, CCIPObjectRef) {
    let mut scenario = test_scenario::begin(@0x1);
    let ctx = scenario.ctx();

    let (owner_cap, ref) = state_object::create(ctx);
    (scenario, owner_cap, ref)
}

fun tear_down_test(scenario: Scenario, owner_cap: OwnerCap, ref: CCIPObjectRef) {
    state_object::destroy_owner_cap(owner_cap);
    state_object::destroy_state_object(ref);
    test_scenario::end(scenario);
}

#[test]
public fun test_initialize() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    rmn_remote::initialize(&mut ref, &owner_cap, 1, ctx);
    let _state = state_object::borrow<RMNRemoteState>(&ref);
    assert!(rmn_remote::get_local_chain_selector(&ref) == 1);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::E_ZERO_VALUE_NOT_ALLOWED)]
public fun test_initialize_zero_chain_selector() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    rmn_remote::initialize(&mut ref, &owner_cap, 0, ctx);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::E_ALREADY_INITIALIZED)]
public fun test_initialize_already_initialized() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    rmn_remote::initialize(&mut ref, &owner_cap, 1, ctx);
    rmn_remote::initialize(&mut ref, &owner_cap, 1, ctx);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_set_config() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    rmn_remote::initialize(&mut ref, &owner_cap, 1, ctx);
    rmn_remote::set_config(
        &mut ref,
        &owner_cap,
        b"00000000000000000000000000000001",
        vector[
            b"00000000000000000002",
            b"00000000000000000003",
            b"00000000000000000004"
        ],
        vector[0, 1, 2],
        1,
    );

    let vc = &rmn_remote::get_versioned_config(&ref);
    let (version, config) = rmn_remote::get_version(vc);

    assert!(version == 1);

    let (digest, signers, f_sign) = rmn_remote::get_config(&config);
    assert!(digest == b"00000000000000000000000000000001");
    assert!(signers.length() == 3);
    assert!(f_sign == 1);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::E_INVALID_DIGEST_LENGTH)]
public fun test_set_config_invalid_digest_length() {
    let( mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    rmn_remote::initialize(&mut ref, &owner_cap, 1, ctx);
    rmn_remote::set_config(
        &mut ref,
        &owner_cap,
        b"000000000000000000000000000000", // invalid digest length
        vector[
            b"00000000000000000002",
            b"00000000000000000003",
            b"00000000000000000004"
        ],
        vector[0, 1, 2],
        1,
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::E_ZERO_VALUE_NOT_ALLOWED)]
public fun test_set_config_zero_digest() {
    let(mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    rmn_remote::initialize(&mut ref, &owner_cap, 1, ctx);
    rmn_remote::set_config(
        &mut ref,
        &owner_cap,
        x"0000000000000000000000000000000000000000000000000000000000000000", // zero digest
        vector[
            b"00000000000000000002",
            b"00000000000000000003",
            b"00000000000000000004"
        ],
        vector[0, 1, 2],
        1,
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::E_NOT_ENOUGH_SIGNERS)]
public fun test_set_config_not_enough_signers() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    rmn_remote::initialize(&mut ref, &owner_cap, 1, ctx);
    rmn_remote::set_config(
        &mut ref,
        &owner_cap,
        b"00000000000000000000000000000001",
        vector[
            b"00000000000000000002",
            b"00000000000000000003",
            b"00000000000000000004"
        ],
        vector[0, 1, 2],
        2, // f_sign is 2, but only 3 signers
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::E_SIGNERS_MISMATCH)]
public fun test_set_config_signers_mismatch() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    rmn_remote::initialize(&mut ref, &owner_cap, 1, ctx);
    rmn_remote::set_config(
        &mut ref,
        &owner_cap,
        b"00000000000000000000000000000001",
        vector[
            b"00000000000000000002",
            b"00000000000000000003"
        ],
        vector[0, 1, 2], // 3 signers, but 2 pub keys
        1,
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::E_INVALID_SIGNER_ORDER)]
public fun test_set_config_invalid_signer_order() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    rmn_remote::initialize(&mut ref, &owner_cap, 1, ctx);
    rmn_remote::set_config(
        &mut ref,
        &owner_cap,
        b"00000000000000000000000000000001",
        vector[
            b"00000000000000000002",
            b"00000000000000000003",
            b"00000000000000000004"
        ],
        vector[1, 0, 2], // invalid order
        1,
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_curse() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    rmn_remote::initialize(&mut ref, &owner_cap, 1, ctx);
    rmn_remote::curse(&mut ref, &owner_cap, b"0000000000000003");

    let cursed_subjects = rmn_remote::get_cursed_subjects(&ref);
    assert!(cursed_subjects.length() == 1);

    assert!(rmn_remote::is_cursed(&ref, b"0000000000000003"));

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::E_INVALID_SUBJECT_LENGTH)]
public fun test_curse_invalid_subject_length() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    rmn_remote::initialize(&mut ref, &owner_cap, 1, ctx);
    rmn_remote::curse(&mut ref, &owner_cap, b"00003");

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::E_ALREADY_CURSED)]
public fun test_curse_already_cursed() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    rmn_remote::initialize(&mut ref, &owner_cap, 1, ctx);
    rmn_remote::curse(&mut ref, &owner_cap, b"0000000000000003");
    rmn_remote::curse(&mut ref, &owner_cap, b"0000000000000003");

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_curse_multiple() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    rmn_remote::initialize(&mut ref, &owner_cap, 1, ctx);
    rmn_remote::curse_multiple(
        &mut ref,
        &owner_cap,
        vector[
            b"0000000000000003",
            b"0000000000000004",
        ],
    );

    let cursed_subjects = rmn_remote::get_cursed_subjects(&ref);
    assert!(cursed_subjects.length() == 2);

    assert!(rmn_remote::is_cursed(&ref, b"0000000000000003"));
    assert!(rmn_remote::is_cursed(&ref, b"0000000000000004"));

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_uncurse() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    rmn_remote::initialize(&mut ref, &owner_cap, 1, ctx);
    rmn_remote::curse(&mut ref, &owner_cap, b"0000000000000003");
    let mut cursed_subjects = rmn_remote::get_cursed_subjects(&ref);
    assert!(cursed_subjects.length() == 1);
    assert!(rmn_remote::is_cursed(&ref, b"0000000000000003"));

    rmn_remote::uncurse(&mut ref, &owner_cap, b"0000000000000003");
    cursed_subjects = rmn_remote::get_cursed_subjects(&ref);
    assert!(cursed_subjects.length() == 0);
    assert!(!rmn_remote::is_cursed(&ref, b"0000000000000003"));

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_is_cursed_global() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    rmn_remote::initialize(&mut ref, &owner_cap, 1, ctx);
    rmn_remote::curse(&mut ref, &owner_cap, x"01000000000000000000000000000001");

    let cursed_subjects = rmn_remote::get_cursed_subjects(&ref);
    assert!(cursed_subjects.length() == 1);
    assert!(rmn_remote::is_cursed_global(&ref));

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_is_cursed_u128() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    rmn_remote::initialize(&mut ref, &owner_cap, 1, ctx);
    rmn_remote::curse(&mut ref, &owner_cap, x"00000000000000000000000000000100"); // hex(256)

    assert!(rmn_remote::is_cursed_u128(&ref, 256));
    assert!(!rmn_remote::is_cursed_u128(&ref, 100));

    tear_down_test(scenario, owner_cap, ref);
}
