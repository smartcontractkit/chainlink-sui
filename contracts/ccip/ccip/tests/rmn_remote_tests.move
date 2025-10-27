#[test_only]
#[allow(implicit_const_copy)]
module ccip::rmn_remote_test;

use ccip::ownable::OwnerCap;
use ccip::rmn_remote::{Self, RMNRemoteState};
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::upgrade_registry;
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
