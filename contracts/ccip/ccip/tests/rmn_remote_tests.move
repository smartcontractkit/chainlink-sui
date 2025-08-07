#[test_only]
#[allow(implicit_const_copy)]
module ccip::rmn_remote_test;

use ccip::state_object::{Self, CCIPObjectRef};
use ccip::ownable::OwnerCap;
use ccip::rmn_remote::{Self, RMNRemoteState};
use sui::test_scenario::{Self, Scenario};

// === Constants ===

// Test addresses and identifiers
const ADMIN_ADDRESS: address = @0x1;
const OFFRAMP_STATE_ADDRESS: address = @0x123;
const TEST_CHAIN_SELECTOR: u64 = 1;

// Test data constants
const VALID_DIGEST: vector<u8> = b"00000000000000000000000000000001";
const ZERO_DIGEST: vector<u8> = x"0000000000000000000000000000000000000000000000000000000000000000";
const INVALID_SHORT_DIGEST: vector<u8> = b"000000000000000000000000000000";

// Signer public keys (20 bytes each)
const SIGNER_PUBKEY_1: vector<u8> = b"00000000000000000002";
const SIGNER_PUBKEY_2: vector<u8> = b"00000000000000000003";
const SIGNER_PUBKEY_3: vector<u8> = b"00000000000000000004";
const SIGNER_PUBKEY_4: vector<u8> = b"00000000000000000005";
const SIGNER_PUBKEY_5: vector<u8> = b"00000000000000000006";
const INVALID_SHORT_PUBKEY: vector<u8> = b"000000000000000000"; // 18 bytes

// Subject identifiers (16 bytes each)
const SUBJECT_1: vector<u8> = b"0000000000000003";
const SUBJECT_2: vector<u8> = b"0000000000000004";
const SUBJECT_U128: vector<u8> = x"00000000000000000000000000000100"; // hex(256)
const GLOBAL_CURSE_SUBJECT: vector<u8> = x"01000000000000000000000000000001";
const INVALID_SHORT_SUBJECT: vector<u8> = b"00003";

// Merkle root test data
const MERKLE_ROOT_VALUE_1: vector<u8> = b"merkle_root_value_32_bytes_long";
const MERKLE_ROOT_VALUE_2: vector<u8> = b"merkle_root_value_32_bytes_lon2";
const ONRAMP_ADDRESS: vector<u8> = b"onramp_addr";

// Signature test data (64 bytes each)
const VALID_SIGNATURE_1: vector<u8> = b"signature_64_bytes_long_signature_64_bytes_long_signature_64_by";
const VALID_SIGNATURE_2: vector<u8> = b"signature_64_bytes_long_signature_64_bytes_long_signature_64_b2";
const INVALID_SHORT_SIGNATURE: vector<u8> = b"invalid_signature_too_short"; // 28 bytes

// Numerical constants
const F_SIGN_VALUE: u64 = 1;
const F_SIGN_HIGH_VALUE: u64 = 2;
const VERSION_1: u32 = 1;
const CHAIN_SELECTOR_100: u64 = 100;
const CHAIN_SELECTOR_200: u64 = 200;
const SEQ_NR_1: u64 = 1;
const SEQ_NR_2: u64 = 2;
const SEQ_NR_10: u64 = 10;
const SEQ_NR_20: u64 = 20;
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

fun initialize_rmn_remote(ref: &mut CCIPObjectRef, owner_cap: &OwnerCap, chain_selector: u64, ctx: &mut TxContext) {
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

fun setup_high_threshold_config(ref: &mut CCIPObjectRef, owner_cap: &OwnerCap) {
    rmn_remote::set_config(
        ref,
        owner_cap,
        VALID_DIGEST,
        vector[SIGNER_PUBKEY_1, SIGNER_PUBKEY_2, SIGNER_PUBKEY_3, SIGNER_PUBKEY_4, SIGNER_PUBKEY_5],
        vector[0, 1, 2, 3, 4],
        F_SIGN_HIGH_VALUE, // f_sign = 2, requires at least 3 signatures
    );
}

fun create_basic_verify_params(): (
    address, // off_ramp_state_address
    vector<u64>, // merkle_root_source_chain_selectors
    vector<vector<u8>>, // merkle_root_on_ramp_addresses
    vector<u64>, // merkle_root_min_seq_nrs
    vector<u64>, // merkle_root_max_seq_nrs
    vector<vector<u8>>, // merkle_root_values
    vector<vector<u8>> // signatures
) {
    (
        OFFRAMP_STATE_ADDRESS,
        vector[CHAIN_SELECTOR_100],
        vector[ONRAMP_ADDRESS],
        vector[SEQ_NR_1],
        vector[SEQ_NR_10],
        vector[MERKLE_ROOT_VALUE_1],
        vector[VALID_SIGNATURE_1, VALID_SIGNATURE_2]
    )
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

#[test]
public fun test_get_arm_with_deployed_package() {
    use std::type_name;
    use sui::address;
    
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);

    let tn = type_name::get<CCIPObjectRef>();
    let addr_string = tn.get_address();
    let expected_package_address = address::from_ascii_bytes(&addr_string.into_bytes());

    // Get the arm address using our function
    let arm_address = rmn_remote::get_arm();

    assert!(arm_address != @0x0);
    assert!(arm_address == expected_package_address);

    tear_down_test(scenario, owner_cap, ref);
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

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::EInvalidDigestLength)]
public fun test_set_config_invalid_digest_length() {
    let( mut scenario, owner_cap, mut ref) = set_up_test();
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
    let(mut scenario, owner_cap, mut ref) = set_up_test();
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

#[test]
#[expected_failure(abort_code = rmn_remote::EConfigNotSet)]
public fun test_verify_config_not_set() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    
    // Try to verify without setting config first
    let (
        off_ramp_state_address,
        merkle_root_source_chain_selectors,
        merkle_root_on_ramp_addresses,
        merkle_root_min_seq_nrs,
        merkle_root_max_seq_nrs,
        merkle_root_values,
        signatures
    ) = create_basic_verify_params();
    
    let _result = rmn_remote::verify(
        &ref,
        off_ramp_state_address,
        merkle_root_source_chain_selectors,
        merkle_root_on_ramp_addresses,
        merkle_root_min_seq_nrs,
        merkle_root_max_seq_nrs,
        merkle_root_values,
        signatures
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::EThresholdNotMet)]
public fun test_verify_threshold_not_met() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    setup_high_threshold_config(&mut ref, &owner_cap);
    
    // Try to verify with only 2 signatures (less than f_sign + 1)
    let _result = rmn_remote::verify(
        &ref,
        OFFRAMP_STATE_ADDRESS,
        vector[CHAIN_SELECTOR_100],
        vector[ONRAMP_ADDRESS],
        vector[SEQ_NR_1],
        vector[SEQ_NR_10],
        vector[MERKLE_ROOT_VALUE_1],
        vector[VALID_SIGNATURE_1, VALID_SIGNATURE_2] // only 2 signatures
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::EMerkleRootLengthMismatch)]
public fun test_verify_merkle_root_length_mismatch() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    setup_basic_config(&mut ref, &owner_cap);
    
    // Mismatched array lengths for merkle root components
    let _result = rmn_remote::verify(
        &ref,
        OFFRAMP_STATE_ADDRESS,
        vector[CHAIN_SELECTOR_100, CHAIN_SELECTOR_200], // 2 elements
        vector[ONRAMP_ADDRESS], // 1 element - mismatch!
        vector[SEQ_NR_1, SEQ_NR_2], // 2 elements
        vector[SEQ_NR_10, SEQ_NR_20], // 2 elements
        vector[MERKLE_ROOT_VALUE_1, MERKLE_ROOT_VALUE_2], // 2 elements
        vector[VALID_SIGNATURE_1, VALID_SIGNATURE_2]
    );

    tear_down_test(scenario, owner_cap, ref);
}

#[test]
#[expected_failure(abort_code = rmn_remote::EInvalidSignature)]
public fun test_verify_invalid_signature_length() {
    let (mut scenario, owner_cap, mut ref) = set_up_test();
    let ctx = scenario.ctx();

    initialize_rmn_remote(&mut ref, &owner_cap, TEST_CHAIN_SELECTOR, ctx);
    setup_basic_config(&mut ref, &owner_cap);
    
    // Try to verify with invalid signature length (not 64 bytes)
    let _result = rmn_remote::verify(
        &ref,
        OFFRAMP_STATE_ADDRESS,
        vector[CHAIN_SELECTOR_100],
        vector[ONRAMP_ADDRESS],
        vector[SEQ_NR_1],
        vector[SEQ_NR_10],
        vector[MERKLE_ROOT_VALUE_1],
        vector[INVALID_SHORT_SIGNATURE, VALID_SIGNATURE_2] // only 28 bytes, should be 64
    );

    tear_down_test(scenario, owner_cap, ref);
}
