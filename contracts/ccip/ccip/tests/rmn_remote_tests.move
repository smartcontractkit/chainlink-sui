#[test_only]
#[allow(implicit_const_copy)]
module ccip::rmn_remote_test;

use ccip::ownable::OwnerCap;
use ccip::rmn_remote::{Self, RMNRemoteState, CurserCap};
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::upgrade_registry;
use mcms::mcms_account;
use mcms::mcms_deployer;
use mcms::mcms_registry::{Self, Registry};
use std::string;
use std::unit_test;
use sui::bcs;
use sui::test_scenario::{Self, Scenario};

// === Constants ===

// Test addresses and identifiers
const ADMIN_ADDRESS: address = @0x1;
const TEST_CHAIN_SELECTOR: u64 = 1;

// Test data constants
const VALID_DIGEST: vector<u8> = b"00000000000000000000000000000001";
const ZERO_DIGEST: vector<u8> = x"0000000000000000000000000000000000000000000000000000000000000000";
const INVALID_SHORT_DIGEST: vector<u8> = b"000000000000000000000000000000";

// Signer public keys (20 bytes each)
const SIGNER_PUBKEY_1: vector<u8> = b"00000000000000000002";
const SIGNER_PUBKEY_2: vector<u8> = b"00000000000000000003";
const SIGNER_PUBKEY_3: vector<u8> = b"00000000000000000004";
const INVALID_SHORT_PUBKEY: vector<u8> = b"000000000000000000"; // 18 bytes

// Subject identifiers (16 bytes each)
const SUBJECT_1: vector<u8> = b"0000000000000003";
const SUBJECT_2: vector<u8> = b"0000000000000004";
const SUBJECT_U128: vector<u8> = x"00000000000000000000000000000100"; // hex(256)
const GLOBAL_CURSE_SUBJECT: vector<u8> = x"01000000000000000000000000000001";
const INVALID_SHORT_SUBJECT: vector<u8> = b"00003";

// Numerical constants
const F_SIGN_VALUE: u64 = 1;
const F_SIGN_HIGH_VALUE: u64 = 2;
const VERSION_1: u32 = 1;
const U128_VALUE_256: u128 = 256;
const U128_VALUE_100: u128 = 100;

// === Helper Functions ===

fun set_up_test(): (Scenario, OwnerCap, CCIPObjectRef) {
    let mut scenario = test_scenario::begin(ADMIN_ADDRESS);
    let ctx = scenario.ctx();

    state_object::test_init(ctx);

    // Advance to next transaction to retrieve the created objects
    scenario.next_tx(ADMIN_ADDRESS);

    // Retrieve the OwnerCap that was transferred to the sender
    let owner_cap = scenario.take_from_sender<OwnerCap>();

    // Retrieve the shared CCIPObjectRef
    let ref = scenario.take_shared<CCIPObjectRef>();

    (scenario, owner_cap, ref)
}

fun tear_down_test(scenario: Scenario, owner_cap: OwnerCap, ref: CCIPObjectRef) {
    // Return the owner cap back to the sender instead of destroying it
    test_scenario::return_to_sender(&scenario, owner_cap);
    // Return the shared object back to the scenario instead of destroying it
    test_scenario::return_shared(ref);
    test_scenario::end(scenario);
}

fun initialize_rmn_remote(
    ref: &mut CCIPObjectRef,
    owner_cap: &OwnerCap,
    chain_selector: u64,
    ctx: &mut TxContext,
) {
    // Initialize upgrade registry first (required by rmn_remote functions)
    upgrade_registry::initialize(ref, owner_cap, ctx);
    rmn_remote::initialize(ref, owner_cap, chain_selector, ctx);
}

fun setup_basic_config(ref: &mut CCIPObjectRef, owner_cap: &OwnerCap) {
    rmn_remote::set_config(
        ref,
        owner_cap,
        VALID_DIGEST,
        vector[SIGNER_PUBKEY_1, SIGNER_PUBKEY_2, SIGNER_PUBKEY_3],
        vector[0, 1, 2],
        F_SIGN_VALUE,
    );
}

// === Basic Initialization Tests ===

#[test]
public fun test_initialize() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    let _state = state_object::borrow<RMNRemoteState>(&ref);
    assert!(rmn_remote::get_local_chain_selector(&ref) == TEST_CHAIN_SELECTOR);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_type_and_version() {
    // Test the type_and_version function
    let version = rmn_remote::type_and_version();
    assert!(version == std::string::utf8(b"RMNRemote 1.6.0"));
}

#[test]
public fun test_get_report_digest_header() {
    // Test the get_report_digest_header function
    let header = rmn_remote::get_report_digest_header();
    // The header should be the keccak256 hash of "RMN_V1_6_ANY2SUI_REPORT"
    assert!(header.length() == 32); // keccak256 produces 32 bytes

    // We can't easily test the exact hash value without keccak256 implementation,
    // but we can verify it's not empty and has correct length
    assert!(header != vector<u8>[]);
}

// === Configuration Management Tests ===

#[test]
public fun test_set_config() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    setup_basic_config(&mut ref, &owner_cap);

    let (version, config) = rmn_remote::get_versioned_config(&ref);

    assert!(version == VERSION_1);

    let (digest, signers, f_sign) = rmn_remote::get_config(&config);
    assert!(digest.length() == VALID_DIGEST.length());
    assert!(signers.length() == 3);
    assert!(f_sign == F_SIGN_VALUE);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_get_config_function() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    setup_basic_config(&mut ref, &owner_cap);

    // Get the config and test the get_config helper function
    let (version, config) = rmn_remote::get_versioned_config(&ref);
    assert!(version == VERSION_1);

    let (digest, signers, f_sign) = rmn_remote::get_config(&config);

    // Verify all config fields
    assert!(digest.length() == VALID_DIGEST.length());
    assert!(signers.length() == 3);
    assert!(f_sign == F_SIGN_VALUE);

    // Note: We can't directly access signer fields without getter functions,
    // but we can verify the length which confirms the structure is correct
    assert!(signers.length() == 3);

    tear_down_test(scenario, owner_cap, ref);
}

// === Curse and Uncurse Tests ===

#[test]
public fun test_curse() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    rmn_remote::curse(&mut ref, &owner_cap, SUBJECT_1);

    let cursed_subjects = rmn_remote::get_cursed_subjects(&ref);
    assert!(cursed_subjects.length() == 1);

    assert!(rmn_remote::is_cursed(&ref, SUBJECT_1));

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_curse_multiple() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    rmn_remote::curse_multiple(
        &mut ref,
        &owner_cap,
        vector[SUBJECT_1, SUBJECT_2],
    );

    let cursed_subjects = rmn_remote::get_cursed_subjects(&ref);
    assert!(cursed_subjects.length() == 2);

    assert!(rmn_remote::is_cursed(&ref, SUBJECT_1));
    assert!(rmn_remote::is_cursed(&ref, SUBJECT_2));

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_uncurse() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    rmn_remote::curse(&mut ref, &owner_cap, SUBJECT_1);
    let mut cursed_subjects = rmn_remote::get_cursed_subjects(&ref);
    assert!(cursed_subjects.length() == 1);
    assert!(rmn_remote::is_cursed(&ref, SUBJECT_1));

    rmn_remote::uncurse(&mut ref, &owner_cap, SUBJECT_1);
    cursed_subjects = rmn_remote::get_cursed_subjects(&ref);
    assert!(cursed_subjects.length() == 0);
    assert!(!rmn_remote::is_cursed(&ref, SUBJECT_1));

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_is_cursed_global() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    rmn_remote::curse(&mut ref, &owner_cap, GLOBAL_CURSE_SUBJECT);

    let cursed_subjects = rmn_remote::get_cursed_subjects(&ref);
    assert!(cursed_subjects.length() == 1);
    assert!(rmn_remote::is_cursed_global(&ref));

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_is_cursed_u128() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    rmn_remote::curse(&mut ref, &owner_cap, SUBJECT_U128);

    assert!(rmn_remote::is_cursed_u128(&ref, U128_VALUE_256));
    assert!(!rmn_remote::is_cursed_u128(&ref, U128_VALUE_100));

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_global_curse_affects_regular_subjects() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);

    // First verify a regular subject is not cursed
    assert!(!rmn_remote::is_cursed(&ref, SUBJECT_1));

    // Curse globally
    rmn_remote::curse(&mut ref, &owner_cap, GLOBAL_CURSE_SUBJECT);

    // Now any subject should be considered cursed due to global curse
    assert!(rmn_remote::is_cursed(&ref, SUBJECT_1));
    assert!(rmn_remote::is_cursed(&ref, SUBJECT_2));
    assert!(rmn_remote::is_cursed_global(&ref));

    // Uncurse globally
    rmn_remote::uncurse(&mut ref, &owner_cap, GLOBAL_CURSE_SUBJECT);

    // Now regular subjects should not be cursed anymore
    assert!(!rmn_remote::is_cursed(&ref, SUBJECT_1));
    assert!(!rmn_remote::is_cursed_global(&ref));

    tear_down_test(scenario, owner_cap, ref);
}

// === Error Condition Tests ===

#[test]
#[expected_failure(abort_code = rmn_remote::EZeroValueNotAllowed)]
public fun test_initialize_zero_chain_selector() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, 0, ctx);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::EAlreadyInitialized)]
public fun test_initialize_already_initialized() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    // Initialize upgrade registry first (required by rmn_remote functions)
    upgrade_registry::initialize(&mut ref, &owner_cap, ctx);
    rmn_remote::initialize(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    // This should fail because rmn_remote is already initialized
    rmn_remote::initialize(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::EInvalidDigestLength)]
public fun test_set_config_invalid_digest_length() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    rmn_remote::set_config(
        &mut ref,
        &owner_cap,
        INVALID_SHORT_DIGEST, // invalid digest length
        vector[SIGNER_PUBKEY_1, SIGNER_PUBKEY_2, SIGNER_PUBKEY_3],
        vector[0, 1, 2],
        F_SIGN_VALUE,
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::EZeroValueNotAllowed)]
public fun test_set_config_zero_digest() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    rmn_remote::set_config(
        &mut ref,
        &owner_cap,
        ZERO_DIGEST, // zero digest
        vector[SIGNER_PUBKEY_1, SIGNER_PUBKEY_2, SIGNER_PUBKEY_3],
        vector[0, 1, 2],
        F_SIGN_VALUE,
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::ENotEnoughSigners)]
public fun test_set_config_not_enough_signers() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    rmn_remote::set_config(
        &mut ref,
        &owner_cap,
        VALID_DIGEST,
        vector[SIGNER_PUBKEY_1, SIGNER_PUBKEY_2, SIGNER_PUBKEY_3],
        vector[0, 1, 2],
        F_SIGN_HIGH_VALUE, // f_sign is 2, but only 3 signers
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::ESignersMismatch)]
public fun test_set_config_signers_mismatch() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    rmn_remote::set_config(
        &mut ref,
        &owner_cap,
        VALID_DIGEST,
        vector[SIGNER_PUBKEY_1, SIGNER_PUBKEY_2],
        vector[0, 1, 2], // 3 signers, but 2 pub keys
        F_SIGN_VALUE,
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::EInvalidSignerOrder)]
public fun test_set_config_invalid_signer_order() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    rmn_remote::set_config(
        &mut ref,
        &owner_cap,
        VALID_DIGEST,
        vector[SIGNER_PUBKEY_1, SIGNER_PUBKEY_2, SIGNER_PUBKEY_3],
        vector[1, 0, 2], // invalid order
        F_SIGN_VALUE,
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::EDuplicateSigner)]
public fun test_set_config_duplicate_signer() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);

    // Try to set config with duplicate signer public keys
    rmn_remote::set_config(
        &mut ref,
        &owner_cap,
        VALID_DIGEST,
        vector[SIGNER_PUBKEY_1, SIGNER_PUBKEY_1, SIGNER_PUBKEY_3], // duplicate!
        vector[0, 1, 2],
        F_SIGN_VALUE,
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::EInvalidPublicKeyLength)]
public fun test_set_config_invalid_public_key_length() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);

    // Try to set config with invalid public key length (not 20 bytes)
    rmn_remote::set_config(
        &mut ref,
        &owner_cap,
        VALID_DIGEST,
        vector[SIGNER_PUBKEY_1, INVALID_SHORT_PUBKEY, SIGNER_PUBKEY_3], // only 18 bytes, should be 20
        vector[0, 1, 2],
        F_SIGN_VALUE,
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::EInvalidSubjectLength)]
public fun test_curse_invalid_subject_length() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    rmn_remote::curse(&mut ref, &owner_cap, INVALID_SHORT_SUBJECT);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::EAlreadyCursed)]
public fun test_curse_already_cursed() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    rmn_remote::curse(&mut ref, &owner_cap, SUBJECT_1);
    rmn_remote::curse(&mut ref, &owner_cap, SUBJECT_1);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::ENotCursed)]
public fun test_uncurse_multiple_not_cursed() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);

    // Try to uncurse subjects that were never cursed
    rmn_remote::uncurse_multiple(
        &mut ref,
        &owner_cap,
        vector[SUBJECT_1, SUBJECT_2], // not cursed
    );

    tear_down_test(scenario, owner_cap, ref);
}

// === Upgrade Registry Function Restriction Tests ===

#[test]
#[expected_failure(abort_code = upgrade_registry::EFunctionNotAllowed)]
public fun test_set_config_function_not_allowed() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);

    // Block the set_config function using upgrade registry
    upgrade_registry::block_function(
        &mut ref,
        &owner_cap,
        std::string::utf8(b"rmn_remote"),
        std::string::utf8(b"set_config"),
        1, // block version 1
        ctx,
    );

    // This should fail because the function is blocked by upgrade registry
    rmn_remote::set_config(
        &mut ref,
        &owner_cap,
        VALID_DIGEST,
        vector[SIGNER_PUBKEY_1, SIGNER_PUBKEY_2, SIGNER_PUBKEY_3],
        vector[0, 1, 2],
        F_SIGN_VALUE,
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = upgrade_registry::EFunctionNotAllowed)]
public fun test_curse_function_not_allowed() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    setup_basic_config(&mut ref, &owner_cap);

    // Block the curse function using upgrade registry
    upgrade_registry::block_function(
        &mut ref,
        &owner_cap,
        std::string::utf8(b"rmn_remote"),
        std::string::utf8(b"curse"),
        1, // block version 1
        ctx,
    );

    // This should fail because the function is blocked by upgrade registry
    rmn_remote::curse(&mut ref, &owner_cap, SUBJECT_1);

    tear_down_test(scenario, owner_cap, ref);
}

// ================================================================
// |              Fast Curse via CurserCap Tests                  |
// ================================================================

// === Direct cap-path tests (no MCMS Registry involved) ===

#[test]
public fun test_create_curser_cap_succeeds() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);

    let curser_cap = rmn_remote::create_curser_cap(&mut ref, &owner_cap, ctx);
    unit_test::destroy(curser_cap);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_curse_with_curser_cap_succeeds() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    let curser_cap = rmn_remote::create_curser_cap(&mut ref, &owner_cap, ctx);

    rmn_remote::curse_with_curser_cap(&mut ref, &curser_cap, SUBJECT_1);
    assert!(rmn_remote::is_cursed(&ref, SUBJECT_1));

    let cursed = rmn_remote::get_cursed_subjects(&ref);
    assert!(cursed.length() == 1);

    unit_test::destroy(curser_cap);
    tear_down_test(scenario, owner_cap, ref);
}

#[test]
public fun test_curse_multiple_with_curser_cap_succeeds() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    let curser_cap = rmn_remote::create_curser_cap(&mut ref, &owner_cap, ctx);

    rmn_remote::curse_multiple_with_curser_cap(
        &mut ref,
        &curser_cap,
        vector[SUBJECT_1, SUBJECT_2],
    );
    assert!(rmn_remote::is_cursed(&ref, SUBJECT_1));
    assert!(rmn_remote::is_cursed(&ref, SUBJECT_2));

    unit_test::destroy(curser_cap);
    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::EInvalidSubjectLength)]
public fun test_curse_with_curser_cap_invalid_subject_length() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    let curser_cap = rmn_remote::create_curser_cap(&mut ref, &owner_cap, ctx);

    rmn_remote::curse_with_curser_cap(&mut ref, &curser_cap, INVALID_SHORT_SUBJECT);

    unit_test::destroy(curser_cap);
    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::EAlreadyCursed)]
public fun test_curse_with_curser_cap_already_cursed() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    let curser_cap = rmn_remote::create_curser_cap(&mut ref, &owner_cap, ctx);

    rmn_remote::curse_with_curser_cap(&mut ref, &curser_cap, SUBJECT_1);
    rmn_remote::curse_with_curser_cap(&mut ref, &curser_cap, SUBJECT_1);

    unit_test::destroy(curser_cap);
    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = upgrade_registry::EFunctionNotAllowed)]
public fun test_curse_with_curser_cap_function_not_allowed() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    let curser_cap = rmn_remote::create_curser_cap(&mut ref, &owner_cap, ctx);

    upgrade_registry::block_function(
        &mut ref,
        &owner_cap,
        string::utf8(b"rmn_remote"),
        string::utf8(b"curse_with_curser_cap"),
        1,
        ctx,
    );

    rmn_remote::curse_with_curser_cap(&mut ref, &curser_cap, SUBJECT_1);

    unit_test::destroy(curser_cap);
    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = upgrade_registry::EFunctionNotAllowed)]
public fun test_create_curser_cap_function_not_allowed() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);

    upgrade_registry::block_function(
        &mut ref,
        &owner_cap,
        string::utf8(b"rmn_remote"),
        string::utf8(b"create_curser_cap"),
        1,
        ctx,
    );

    let curser_cap = rmn_remote::create_curser_cap(&mut ref, &owner_cap, ctx);
    unit_test::destroy(curser_cap);

    tear_down_test(scenario, owner_cap, ref);
}

// === MCMS callback tests ===

const BATCH_ID_1: vector<u8> = x"0000000000000000000000000000000000000000000000000000000000000001";
const BATCH_ID_2: vector<u8> = x"0000000000000000000000000000000000000000000000000000000000000002";
const BATCH_ID_3: vector<u8> = x"0000000000000000000000000000000000000000000000000000000000000003";
const BATCH_ID_4: vector<u8> = x"0000000000000000000000000000000000000000000000000000000000000004";

/// Initializes CCIP and a single MCMS Registry (used here as the FAST Registry)
/// pre-populated with a fresh `CurserCap` for `@ccip` and
/// `allowed_modules = [b"rmn_remote"]`. Used to exercise the fast-path
/// callbacks (`mcms_curse_with_curser_cap` etc.) end-to-end.
fun setup_fast_registry_with_cap(): (Scenario, OwnerCap, CCIPObjectRef, Registry, ID) {
    let mut scenario = test_scenario::begin(ADMIN_ADDRESS);
    {
        let ctx = scenario.ctx();
        mcms_account::test_init(ctx);
        mcms_registry::test_init(ctx);
        mcms_deployer::test_init(ctx);
        state_object::test_init(ctx);
    };
    scenario.next_tx(ADMIN_ADDRESS);

    let mut fast_registry = test_scenario::take_shared<Registry>(&scenario);
    let mut ref = test_scenario::take_shared<CCIPObjectRef>(&scenario);
    let owner_cap = test_scenario::take_from_sender<OwnerCap>(&scenario);

    upgrade_registry::initialize(&mut ref, &owner_cap, scenario.ctx());
    rmn_remote::initialize(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, scenario.ctx());

    let curser_cap = rmn_remote::create_curser_cap(&mut ref, &owner_cap, scenario.ctx());
    let curser_cap_id = object::id_address(&curser_cap);

    let publisher_wrapper = mcms_registry::create_publisher_wrapper(
        ccip::ownable::borrow_publisher(&owner_cap),
        state_object::test_create_mcms_callback(),
    );
    mcms_registry::register_entrypoint(
        &mut fast_registry,
        publisher_wrapper,
        state_object::test_create_mcms_callback(),
        curser_cap,
        vector[b"rmn_remote"],
        scenario.ctx(),
    );

    scenario.next_tx(ADMIN_ADDRESS);

    (scenario, owner_cap, ref, fast_registry, sui::object::id_from_address(curser_cap_id))
}

fun tear_down_fast_registry(
    scenario: Scenario,
    owner_cap: OwnerCap,
    ref: CCIPObjectRef,
    fast_registry: Registry,
) {
    test_scenario::return_to_sender(&scenario, owner_cap);
    test_scenario::return_shared(ref);
    test_scenario::return_shared(fast_registry);
    test_scenario::end(scenario);
}

#[test]
public fun test_mcms_curse_with_curser_cap_succeeds() {
    let (mut scenario, owner_cap, mut ref, mut fast_registry, curser_cap_id) =
        setup_fast_registry_with_cap();

    let mut data = vector::empty<u8>();
    data.append(bcs::to_bytes(&object::id_address(&ref)));
    data.append(bcs::to_bytes(&sui::object::id_to_address(&curser_cap_id)));
    data.append(bcs::to_bytes(&SUBJECT_1));

    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(b"rmn_remote"),
        string::utf8(b"curse_with_curser_cap"),
        data,
        BATCH_ID_1,
        0,
        1,
    );

    rmn_remote::mcms_curse_with_curser_cap(&mut ref, &mut fast_registry, params);

    assert!(rmn_remote::is_cursed(&ref, SUBJECT_1));

    tear_down_fast_registry(scenario, owner_cap, ref, fast_registry);
}

#[test]
public fun test_mcms_curse_multiple_with_curser_cap_succeeds() {
    let (mut scenario, owner_cap, mut ref, mut fast_registry, curser_cap_id) =
        setup_fast_registry_with_cap();

    let subjects = vector[SUBJECT_1, SUBJECT_2];

    let mut data = vector::empty<u8>();
    data.append(bcs::to_bytes(&object::id_address(&ref)));
    data.append(bcs::to_bytes(&sui::object::id_to_address(&curser_cap_id)));
    data.append(bcs::to_bytes(&subjects));

    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(b"rmn_remote"),
        string::utf8(b"curse_multiple_with_curser_cap"),
        data,
        BATCH_ID_2,
        0,
        1,
    );

    rmn_remote::mcms_curse_multiple_with_curser_cap(&mut ref, &mut fast_registry, params);

    assert!(rmn_remote::is_cursed(&ref, SUBJECT_1));
    assert!(rmn_remote::is_cursed(&ref, SUBJECT_2));

    tear_down_fast_registry(scenario, owner_cap, ref, fast_registry);
}

#[test]
#[expected_failure(abort_code = rmn_remote::EInvalidFunction)]
public fun test_mcms_curse_with_curser_cap_wrong_function_name() {
    let (mut scenario, owner_cap, mut ref, mut fast_registry, curser_cap_id) =
        setup_fast_registry_with_cap();

    let mut data = vector::empty<u8>();
    data.append(bcs::to_bytes(&object::id_address(&ref)));
    data.append(bcs::to_bytes(&sui::object::id_to_address(&curser_cap_id)));
    data.append(bcs::to_bytes(&SUBJECT_1));

    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(b"rmn_remote"),
        string::utf8(b"uncurse"),
        data,
        BATCH_ID_3,
        0,
        1,
    );

    rmn_remote::mcms_curse_with_curser_cap(&mut ref, &mut fast_registry, params);

    tear_down_fast_registry(scenario, owner_cap, ref, fast_registry);
}

// === Slow-MCMS bootstrap tests (two Registries) ===

const FAST_RECIPIENT: address = @0xF45CA;

/// Initializes CCIP, transfers `OwnerCap` into a SLOW MCMS Registry, then
/// shares a second empty Registry to act as the FAST instance.
fun setup_slow_and_fast_registries(): (Scenario, CCIPObjectRef, Registry /* slow */, Registry /* fast */) {
    let mut scenario = test_scenario::begin(ADMIN_ADDRESS);
    {
        let ctx = scenario.ctx();
        mcms_account::test_init(ctx);
        mcms_registry::test_init(ctx);
        mcms_deployer::test_init(ctx);
        state_object::test_init(ctx);
    };
    scenario.next_tx(ADMIN_ADDRESS);

    let mut slow_registry = test_scenario::take_shared<Registry>(&scenario);
    let mut ref = test_scenario::take_shared<CCIPObjectRef>(&scenario);
    let owner_cap = test_scenario::take_from_sender<OwnerCap>(&scenario);

    upgrade_registry::initialize(&mut ref, &owner_cap, scenario.ctx());
    rmn_remote::initialize(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, scenario.ctx());

    state_object::transfer_ownership(
        &mut ref,
        &owner_cap,
        mcms_registry::get_multisig_address(),
        scenario.ctx(),
    );
    scenario.next_tx(mcms_registry::get_multisig_address());
    state_object::accept_ownership(&mut ref, scenario.ctx());

    state_object::execute_ownership_transfer_to_mcms(
        &mut ref,
        owner_cap,
        &mut slow_registry,
        @mcms,
        scenario.ctx(),
    );

    scenario.next_tx(ADMIN_ADDRESS);
    {
        let ctx = scenario.ctx();
        mcms_registry::test_init(ctx);
    };
    scenario.next_tx(ADMIN_ADDRESS);

    let fast_registry = test_scenario::take_shared<Registry>(&scenario);

    (scenario, ref, slow_registry, fast_registry)
}

fun tear_down_slow_and_fast(
    scenario: Scenario,
    ref: CCIPObjectRef,
    slow_registry: Registry,
    fast_registry: Registry,
) {
    test_scenario::return_shared(ref);
    test_scenario::return_shared(slow_registry);
    test_scenario::return_shared(fast_registry);
    test_scenario::end(scenario);
}

#[test]
public fun test_mcms_create_curser_cap_succeeds() {
    let (mut scenario, mut ref, mut slow_registry, fast_registry) = setup_slow_and_fast_registries();

    let owner_cap_address = test_owner_cap_address_in_registry(&slow_registry);

    let mut data = vector::empty<u8>();
    data.append(bcs::to_bytes(&object::id_address(&ref)));
    data.append(bcs::to_bytes(&owner_cap_address));

    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(b"rmn_remote"),
        string::utf8(b"create_curser_cap"),
        data,
        BATCH_ID_1,
        0,
        1,
    );

    let curser_cap = rmn_remote::mcms_create_curser_cap(
        &mut ref,
        &mut slow_registry,
        params,
        scenario.ctx(),
    );

    sui::transfer::public_transfer(curser_cap, FAST_RECIPIENT);

    tear_down_slow_and_fast(scenario, ref, slow_registry, fast_registry);
}

#[test]
public fun test_mcms_register_curser_cap_succeeds() {
    let (mut scenario, mut ref, mut slow_registry, mut fast_registry) =
        setup_slow_and_fast_registries();

    let owner_cap_address = test_owner_cap_address_in_registry(&slow_registry);

    // Op 0: mcms_create_curser_cap
    let mut create_data = vector::empty<u8>();
    create_data.append(bcs::to_bytes(&object::id_address(&ref)));
    create_data.append(bcs::to_bytes(&owner_cap_address));

    let create_params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(b"rmn_remote"),
        string::utf8(b"create_curser_cap"),
        create_data,
        BATCH_ID_2,
        0,
        2,
    );

    let curser_cap = rmn_remote::mcms_create_curser_cap(
        &mut ref,
        &mut slow_registry,
        create_params,
        scenario.ctx(),
    );

    // Op 1: mcms_register_curser_cap
    let mut register_data = vector::empty<u8>();
    register_data.append(bcs::to_bytes(&object::id_address(&ref)));
    register_data.append(bcs::to_bytes(&owner_cap_address));
    register_data.append(bcs::to_bytes(&object::id_address(&fast_registry)));

    let register_params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(b"rmn_remote"),
        string::utf8(b"register_curser_cap"),
        register_data,
        BATCH_ID_2,
        1,
        2,
    );

    rmn_remote::mcms_register_curser_cap(
        &mut ref,
        &mut slow_registry,
        &mut fast_registry,
        register_params,
        curser_cap,
        scenario.ctx(),
    );

    let allowed = mcms_registry::get_allowed_modules(
        &fast_registry,
        sui::address::to_ascii_string(@ccip),
    );
    assert!(allowed == vector[b"rmn_remote"]);

    tear_down_slow_and_fast(scenario, ref, slow_registry, fast_registry);
}

#[test]
public fun test_mcms_mint_and_register_curser_cap_succeeds() {
    let (mut scenario, mut ref, mut slow_registry, mut fast_registry) =
        setup_slow_and_fast_registries();

    let owner_cap_address = test_owner_cap_address_in_registry(&slow_registry);

    let mut data = vector::empty<u8>();
    data.append(bcs::to_bytes(&object::id_address(&ref)));
    data.append(bcs::to_bytes(&owner_cap_address));
    data.append(bcs::to_bytes(&object::id_address(&fast_registry)));

    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(b"rmn_remote"),
        string::utf8(b"mint_and_register_curser_cap"),
        data,
        BATCH_ID_4,
        0,
        1,
    );

    rmn_remote::mcms_mint_and_register_curser_cap(
        &mut ref,
        &mut slow_registry,
        &mut fast_registry,
        params,
        scenario.ctx(),
    );

    let allowed = mcms_registry::get_allowed_modules(
        &fast_registry,
        sui::address::to_ascii_string(@ccip),
    );
    assert!(allowed == vector[b"rmn_remote"]);

    tear_down_slow_and_fast(scenario, ref, slow_registry, fast_registry);
}

#[test]
#[expected_failure(abort_code = mcms_registry::EPackageCapNotRegistered)]
public fun test_mcms_curse_with_curser_cap_unregistered_registry() {
    let (mut scenario, owner_cap, mut ref, mut fast_registry, curser_cap_id) =
        setup_fast_registry_with_cap();

    let released_cap = mcms_registry::release_cap<
        state_object::McmsCallback,
        rmn_remote::CurserCap,
    >(
        &mut fast_registry,
        state_object::test_create_mcms_callback(),
    );
    unit_test::destroy(released_cap);

    let mut data = vector::empty<u8>();
    data.append(bcs::to_bytes(&object::id_address(&ref)));
    data.append(bcs::to_bytes(&sui::object::id_to_address(&curser_cap_id)));
    data.append(bcs::to_bytes(&SUBJECT_1));

    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(b"rmn_remote"),
        string::utf8(b"curse_with_curser_cap"),
        data,
        BATCH_ID_1,
        0,
        1,
    );

    rmn_remote::mcms_curse_with_curser_cap(&mut ref, &mut fast_registry, params);

    tear_down_fast_registry(scenario, owner_cap, ref, fast_registry);
}

#[test]
public fun test_mcms_curse_and_uncurse_via_owner_cap() {
    let (mut scenario, mut ref, mut slow_registry, fast_registry) = setup_slow_and_fast_registries();

    let owner_cap_address = test_owner_cap_address_in_registry(&slow_registry);

    let mut curse_data = vector::empty<u8>();
    curse_data.append(bcs::to_bytes(&object::id_address(&ref)));
    curse_data.append(bcs::to_bytes(&owner_cap_address));
    curse_data.append(bcs::to_bytes(&vector[SUBJECT_1]));

    let curse_params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(b"rmn_remote"),
        string::utf8(b"curse_multiple"),
        curse_data,
        BATCH_ID_3,
        0,
        1,
    );

    rmn_remote::mcms_curse_multiple(&mut ref, &mut slow_registry, curse_params);
    assert!(rmn_remote::is_cursed(&ref, SUBJECT_1));

    let mut uncurse_data = vector::empty<u8>();
    uncurse_data.append(bcs::to_bytes(&object::id_address(&ref)));
    uncurse_data.append(bcs::to_bytes(&owner_cap_address));
    uncurse_data.append(bcs::to_bytes(&vector[SUBJECT_1]));

    let uncurse_params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(b"rmn_remote"),
        string::utf8(b"uncurse_multiple"),
        uncurse_data,
        BATCH_ID_3,
        0,
        1,
    );

    rmn_remote::mcms_uncurse_multiple(&mut ref, &mut slow_registry, uncurse_params);
    assert!(!rmn_remote::is_cursed(&ref, SUBJECT_1));

    tear_down_slow_and_fast(scenario, ref, slow_registry, fast_registry);
}

#[test]
#[expected_failure]
public fun test_mcms_uncurse_on_fast_registry_aborts() {
    let (mut scenario, owner_cap, mut ref, mut fast_registry, curser_cap_id) =
        setup_fast_registry_with_cap();

    let mut data = vector::empty<u8>();
    data.append(bcs::to_bytes(&object::id_address(&ref)));
    data.append(bcs::to_bytes(&sui::object::id_to_address(&curser_cap_id)));
    data.append(bcs::to_bytes(&vector[SUBJECT_1]));

    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(b"rmn_remote"),
        string::utf8(b"uncurse_multiple"),
        data,
        BATCH_ID_3,
        0,
        1,
    );

    rmn_remote::mcms_uncurse_multiple(&mut ref, &mut fast_registry, params);

    tear_down_fast_registry(scenario, owner_cap, ref, fast_registry);
}

#[test]
public fun test_global_curse_via_curser_cap_fast_path() {
    let (mut scenario, owner_cap, mut ref, mut fast_registry, curser_cap_id) =
        setup_fast_registry_with_cap();

    let mut data = vector::empty<u8>();
    data.append(bcs::to_bytes(&object::id_address(&ref)));
    data.append(bcs::to_bytes(&sui::object::id_to_address(&curser_cap_id)));
    data.append(bcs::to_bytes(&x"01000000000000000000000000000001"));

    let params = mcms_registry::test_create_executing_callback_params(
        @ccip,
        string::utf8(b"rmn_remote"),
        string::utf8(b"curse_with_curser_cap"),
        data,
        BATCH_ID_1,
        0,
        1,
    );

    rmn_remote::mcms_curse_with_curser_cap(&mut ref, &mut fast_registry, params);
    assert!(rmn_remote::is_cursed_global(&ref));

    tear_down_fast_registry(scenario, owner_cap, ref, fast_registry);
}

/// Helper to read the address of the `OwnerCap` stored inside a slow MCMS
/// Registry, so callers can populate `validate_obj_addrs` data correctly.
fun test_owner_cap_address_in_registry(registry: &Registry): address {
    mcms_registry::test_get_cap_address<OwnerCap>(
        registry,
        sui::address::to_ascii_string(@ccip),
    )
}
