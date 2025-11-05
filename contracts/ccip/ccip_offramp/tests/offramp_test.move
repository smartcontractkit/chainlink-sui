#[test_only]
#[allow(implicit_const_copy)]
module ccip_offramp::offramp_test;

use ccip::fee_quoter::{Self, FeeQuoterCap};
use ccip::offramp_state_helper::{Self as osh, DestTransferCap};
use ccip::receiver_registry;
use ccip::rmn_remote;
use ccip::state_object::{Self, CCIPObjectRef};
use ccip::token_admin_registry;
use ccip::upgrade_registry;
use ccip_offramp::ocr3_base;
use ccip_offramp::offramp::{Self, OffRampState, StaticConfig};
use ccip_offramp::ownable::OwnerCap;
use mcms::mcms_registry::{Self, Registry};
use sui::clock;
use sui::test_scenario::{Self as ts, Scenario};

const OWNER: address = @0x123;
const CHAIN_SELECTOR: u64 = 1000;
const SOURCE_CHAIN_SELECTOR_1: u64 = 2000;
const SOURCE_CHAIN_SELECTOR_2: u64 = 3000;
const PERMISSIONLESS_EXECUTION_THRESHOLD: u32 = 3600; // 1 hour

public struct TestEnv {
    scenario: Scenario,
    state: OffRampState,
    ref: CCIPObjectRef,
    clock: clock::Clock,
}

fun setup(): (TestEnv, OwnerCap, FeeQuoterCap, DestTransferCap) {
    let mut scenario = ts::begin(OWNER);
    let ctx = scenario.ctx();
    let mut clock = clock::create_for_testing(ctx);
    clock.set_for_testing(1_000_000_000);

    // Initialize CCIP components
    state_object::test_init(ctx);
    osh::test_init(ctx);
    offramp::test_init(ctx);

    scenario.next_tx(OWNER);

    let mut ref = ts::take_shared<CCIPObjectRef>(&scenario);
    let state = ts::take_shared<OffRampState>(&scenario);

    let ccip_owner_cap = ts::take_from_sender<ccip::ownable::OwnerCap>(&scenario);
    let owner_cap = ts::take_from_sender<OwnerCap>(&scenario);
    let dest_transfer_cap = ts::take_from_sender<DestTransferCap>(&scenario);

    // Initialize required CCIP components
    upgrade_registry::initialize(&mut ref, &ccip_owner_cap, scenario.ctx());
    token_admin_registry::initialize(&mut ref, &ccip_owner_cap, scenario.ctx());
    rmn_remote::initialize(&mut ref, &ccip_owner_cap, 1000, scenario.ctx());
    receiver_registry::initialize(&mut ref, &ccip_owner_cap, scenario.ctx());
    fee_quoter::initialize(
        &mut ref,
        &ccip_owner_cap,
        1000000, // max_fee_juels_per_msg
        @0x1, // link_token address
        3600, // token_price_staleness_threshold
        vector[], // fee_tokens
        scenario.ctx(),
    );

    scenario.next_tx(OWNER);
    let fee_quoter_cap = ts::take_from_sender<FeeQuoterCap>(&scenario);

    ts::return_to_address(OWNER, ccip_owner_cap);

    let env = TestEnv {
        scenario,
        state,
        ref,
        clock,
    };

    (env, owner_cap, fee_quoter_cap, dest_transfer_cap)
}

fun tear_down(env: TestEnv) {
    let TestEnv { scenario, state, ref, clock } = env;

    ts::return_shared(state);
    ts::return_shared(ref);
    clock.destroy_for_testing();
    ts::end(scenario);
}

fun initialize_offramp(
    env: &mut TestEnv,
    owner_cap: &OwnerCap,
    fee_quoter_cap: FeeQuoterCap,
    dest_transfer_cap: DestTransferCap,
) {
    offramp::initialize(
        &mut env.state,
        owner_cap,
        fee_quoter_cap,
        dest_transfer_cap,
        CHAIN_SELECTOR,
        PERMISSIONLESS_EXECUTION_THRESHOLD,
        vector[SOURCE_CHAIN_SELECTOR_1, SOURCE_CHAIN_SELECTOR_2], // source_chains_selectors
        vector[true, false], // source_chains_is_enabled
        vector[false, true], // source_chains_is_rmn_verification_disabled
        vector[b"onramp_1", b"onramp_2"], // source_chains_on_ramp
        env.scenario.ctx(),
    );
}

#[test]
public fun test_initialize() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();
    initialize_offramp(&mut env, &owner_cap, fee_quoter_cap, dest_transfer_cap);

    // Test static config
    let static_config = offramp::get_static_config(&env.ref, &env.state);
    let (
        chain_selector,
        rmn_remote,
        token_admin_registry,
        nonce_manager,
    ) = offramp::get_static_config_fields(&env.ref, static_config);
    assert!(chain_selector == CHAIN_SELECTOR);
    assert!(rmn_remote == @ccip);
    assert!(token_admin_registry == @ccip);
    assert!(nonce_manager == @ccip);

    // Test dynamic config
    let dynamic_config = offramp::get_dynamic_config(&env.ref, &env.state);
    let (fee_quoter, threshold) = offramp::get_dynamic_config_fields(&env.ref, dynamic_config);
    assert!(fee_quoter == @ccip);
    assert!(threshold == PERMISSIONLESS_EXECUTION_THRESHOLD);

    // Test source chain configs
    let source_config_1 = offramp::get_source_chain_config(
        &env.ref,
        &env.state,
        SOURCE_CHAIN_SELECTOR_1,
    );
    let (
        router,
        is_enabled,
        min_seq_nr,
        is_rmn_disabled,
        on_ramp,
    ) = offramp::get_source_chain_config_fields(&env.ref, source_config_1);
    assert!(router == @ccip);
    assert!(is_enabled == true);
    assert!(min_seq_nr == 1);
    assert!(is_rmn_disabled == false);
    assert!(on_ramp == b"onramp_1");

    let source_config_2 = offramp::get_source_chain_config(
        &env.ref,
        &env.state,
        SOURCE_CHAIN_SELECTOR_2,
    );
    let (
        router,
        is_enabled,
        min_seq_nr,
        is_rmn_disabled,
        on_ramp,
    ) = offramp::get_source_chain_config_fields(&env.ref, source_config_2);
    assert!(router == @ccip);
    assert!(is_enabled == false);
    assert!(min_seq_nr == 1);
    assert!(is_rmn_disabled == true);
    assert!(on_ramp == b"onramp_2");

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
}

#[test]
public fun test_set_dynamic_config() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();
    initialize_offramp(&mut env, &owner_cap, fee_quoter_cap, dest_transfer_cap);

    // Test initial config
    let initial_config = offramp::get_dynamic_config(&env.ref, &env.state);
    let (_, initial_threshold) = offramp::get_dynamic_config_fields(&env.ref, initial_config);
    assert!(initial_threshold == PERMISSIONLESS_EXECUTION_THRESHOLD);

    // Update config
    let new_threshold = 7200; // 2 hours
    offramp::set_dynamic_config(&env.ref, &mut env.state, &owner_cap, new_threshold);

    // Verify update
    let updated_config = offramp::get_dynamic_config(&env.ref, &env.state);
    let (_, updated_threshold) = offramp::get_dynamic_config_fields(&env.ref, updated_config);
    assert!(updated_threshold == new_threshold);

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
}

#[test]
public fun test_apply_source_chain_config_updates() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();
    initialize_offramp(&mut env, &owner_cap, fee_quoter_cap, dest_transfer_cap);

    // Add a new source chain
    let new_chain_selector = 4000;
    offramp::apply_source_chain_config_updates(
        &env.ref,
        &mut env.state,
        &owner_cap,
        vector[new_chain_selector],
        vector[true], // enabled
        vector[false], // rmn verification enabled
        vector[b"onramp_3"],
        env.scenario.ctx(),
    );

    // Verify new chain was added
    let new_config = offramp::get_source_chain_config(&env.ref, &env.state, new_chain_selector);
    let (
        router,
        is_enabled,
        min_seq_nr,
        is_rmn_disabled,
        on_ramp,
    ) = offramp::get_source_chain_config_fields(&env.ref, new_config);
    assert!(router == @ccip);
    assert!(is_enabled == true);
    assert!(min_seq_nr == 1);
    assert!(is_rmn_disabled == false);
    assert!(on_ramp == b"onramp_3");

    // Update existing chain
    offramp::apply_source_chain_config_updates(
        &env.ref,
        &mut env.state,
        &owner_cap,
        vector[SOURCE_CHAIN_SELECTOR_2],
        vector[true], // enable previously disabled chain
        vector[false], // enable rmn verification
        vector[b"onramp_2"], // same onramp
        env.scenario.ctx(),
    );

    // Verify update
    let updated_config = offramp::get_source_chain_config(
        &env.ref,
        &env.state,
        SOURCE_CHAIN_SELECTOR_2,
    );
    let (_, is_enabled, _, is_rmn_disabled, _) = offramp::get_source_chain_config_fields(
        &env.ref,
        updated_config,
    );
    assert!(is_enabled == true);
    assert!(is_rmn_disabled == false);

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
}

#[test]
public fun test_get_all_source_chain_configs() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();
    initialize_offramp(&mut env, &owner_cap, fee_quoter_cap, dest_transfer_cap);

    // Get all configs
    let (chain_selectors, chain_configs) = offramp::get_all_source_chain_configs(
        &env.ref,
        &env.state,
    );

    // Should have 2 chains from initialization
    assert!(chain_selectors.length() == 2);
    assert!(chain_configs.length() == 2);

    // Verify chain selectors are present
    assert!(chain_selectors.contains(&SOURCE_CHAIN_SELECTOR_1));
    assert!(chain_selectors.contains(&SOURCE_CHAIN_SELECTOR_2));

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
}

#[test]
public fun test_set_ocr3_config() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();
    initialize_offramp(&mut env, &owner_cap, fee_quoter_cap, dest_transfer_cap);

    // Set OCR3 config for commit plugin
    let config_digest = b"test_config_digest_32_bytes_long";
    let commit_plugin_type = ocr3_base::ocr_plugin_type_commit();
    let big_f = 1;
    // Need more than 3 * big_f signers, so with big_f = 1, need > 3 signers
    // Signers must be 32-byte ed25519 public keys
    let signers = vector[
        x"1111111111111111111111111111111111111111111111111111111111111111",
        x"2222222222222222222222222222222222222222222222222222222222222222",
        x"3333333333333333333333333333333333333333333333333333333333333333",
        x"4444444444444444444444444444444444444444444444444444444444444444",
        x"5555555555555555555555555555555555555555555555555555555555555555",
    ];
    let transmitters = vector[@0x100, @0x200, @0x300, @0x400, @0x500];

    offramp::set_ocr3_config(
        &env.ref,
        &mut env.state,
        &owner_cap,
        config_digest,
        commit_plugin_type,
        big_f,
        true, // signature verification enabled for commit
        signers,
        transmitters,
    );

    // Verify config was set
    let latest_config = offramp::latest_config_details(&env.state, commit_plugin_type);
    let (
        digest,
        f,
        n,
        sig_verification,
        config_signers,
        config_transmitters,
    ) = offramp::latest_config_digest_fields(latest_config);

    assert!(digest == config_digest);
    assert!(f == big_f);
    assert!(n == 5); // number of signers
    assert!(sig_verification == true);
    assert!(config_signers.length() == 5);
    assert!(config_transmitters.length() == 5);

    // Set OCR3 config for execution plugin
    let execution_plugin_type = ocr3_base::ocr_plugin_type_execution();
    offramp::set_ocr3_config(
        &env.ref,
        &mut env.state,
        &owner_cap,
        config_digest,
        execution_plugin_type,
        big_f,
        false, // signature verification disabled for execution
        signers,
        transmitters,
    );

    // Verify execution config
    let exec_config = offramp::latest_config_details(&env.state, execution_plugin_type);
    let (_, _, _, exec_sig_verification, _, _) = offramp::latest_config_digest_fields(exec_config);
    assert!(exec_sig_verification == false);

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
}

#[test]
public fun test_get_latest_price_sequence_number() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();
    initialize_offramp(&mut env, &owner_cap, fee_quoter_cap, dest_transfer_cap);

    // Initially should be 0
    let initial_seq = offramp::get_latest_price_sequence_number(&env.state);
    assert!(initial_seq == 0);

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
}

#[test]
#[expected_failure(abort_code = offramp::EUnknownSourceChainSelector)]
public fun test_get_execution_state() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();
    initialize_offramp(&mut env, &owner_cap, fee_quoter_cap, dest_transfer_cap);

    // Test getting execution state for an unknown source chain selector
    // This should fail because the source chain selector doesn't exist
    let _execution_state = offramp::get_execution_state(&env.state, 9998, 1);

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
}

#[test]
public fun test_get_source_chain_config_nonexistent() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();
    initialize_offramp(&mut env, &owner_cap, fee_quoter_cap, dest_transfer_cap);

    // Test getting config for non-existent chain
    let nonexistent_chain = 9999;
    let config = offramp::get_source_chain_config(&env.ref, &env.state, nonexistent_chain);
    let (
        router,
        is_enabled,
        min_seq_nr,
        is_rmn_disabled,
        on_ramp,
    ) = offramp::get_source_chain_config_fields(&env.ref, config);

    // Should return default empty config
    assert!(router == @0x0);
    assert!(is_enabled == false);
    assert!(min_seq_nr == 0);
    assert!(is_rmn_disabled == false);
    assert!(on_ramp == vector[]);

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
}

#[test]
public fun test_type_and_version() {
    let version = offramp::type_and_version();
    assert!(version == std::string::utf8(b"OffRamp 1.6.0"));
}

#[test]
public fun test_get_ccip_package_id() {
    let package_id = offramp::get_ccip_package_id();
    assert!(package_id == @ccip);
}

#[test]
#[expected_failure(abort_code = offramp::ESourceChainSelectorsMismatch)]
public fun test_initialize_mismatched_vectors() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();

    // Try to initialize with mismatched vector lengths
    offramp::initialize(
        &mut env.state,
        &owner_cap,
        fee_quoter_cap,
        dest_transfer_cap,
        CHAIN_SELECTOR,
        PERMISSIONLESS_EXECUTION_THRESHOLD,
        vector[SOURCE_CHAIN_SELECTOR_1], // 1 element
        vector[true, false], // 2 elements - mismatch!
        vector[false],
        vector[b"onramp_1"],
        env.scenario.ctx(),
    );

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
}

#[test]
#[expected_failure(abort_code = offramp::EZeroChainSelector)]
public fun test_initialize_zero_chain_selector() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();

    // Try to initialize with zero chain selector
    offramp::initialize(
        &mut env.state,
        &owner_cap,
        fee_quoter_cap,
        dest_transfer_cap,
        CHAIN_SELECTOR,
        PERMISSIONLESS_EXECUTION_THRESHOLD,
        vector[0], // zero chain selector - should fail
        vector[true],
        vector[false],
        vector[b"onramp_1"],
        env.scenario.ctx(),
    );

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
}

#[test]
#[expected_failure(abort_code = offramp::ESignatureVerificationRequiredInCommitPlugin)]
public fun test_ocr3_config_commit_requires_signature_verification() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();
    initialize_offramp(&mut env, &owner_cap, fee_quoter_cap, dest_transfer_cap);

    // Try to set commit plugin without signature verification - should fail
    let signers = vector[
        x"1111111111111111111111111111111111111111111111111111111111111111",
        x"2222222222222222222222222222222222222222222222222222222222222222",
        x"3333333333333333333333333333333333333333333333333333333333333333",
        x"4444444444444444444444444444444444444444444444444444444444444444",
        x"5555555555555555555555555555555555555555555555555555555555555555",
    ];
    offramp::set_ocr3_config(
        &env.ref,
        &mut env.state,
        &owner_cap,
        b"test_config_digest_32_bytes_long",
        ocr3_base::ocr_plugin_type_commit(),
        1,
        false, // signature verification disabled - should fail for commit
        signers,
        vector[@0x100, @0x200, @0x300, @0x400, @0x500],
    );

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
}

#[test]
#[expected_failure(abort_code = offramp::ESignatureVerificationNotAllowedInExecutionPlugin)]
public fun test_ocr3_config_execution_forbids_signature_verification() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();
    initialize_offramp(&mut env, &owner_cap, fee_quoter_cap, dest_transfer_cap);

    // Try to set execution plugin with signature verification - should fail
    let signers = vector[
        x"1111111111111111111111111111111111111111111111111111111111111111",
        x"2222222222222222222222222222222222222222222222222222222222222222",
        x"3333333333333333333333333333333333333333333333333333333333333333",
        x"4444444444444444444444444444444444444444444444444444444444444444",
        x"5555555555555555555555555555555555555555555555555555555555555555",
    ];
    offramp::set_ocr3_config(
        &env.ref,
        &mut env.state,
        &owner_cap,
        b"test_config_digest_32_bytes_long",
        ocr3_base::ocr_plugin_type_execution(),
        1,
        true, // signature verification enabled - should fail for execution
        signers,
        vector[@0x100, @0x200, @0x300, @0x400, @0x500],
    );

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
}

#[test]
#[expected_failure(abort_code = offramp::EUnknownSourceChainSelector)]
public fun test_get_execution_state_unknown_chain() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();
    initialize_offramp(&mut env, &owner_cap, fee_quoter_cap, dest_transfer_cap);

    // Try to get execution state for unknown chain
    let _execution_state = offramp::get_execution_state(&env.state, 9999, 1);

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
}

// === Tests for Uncovered Functions ===

#[test]
public fun test_get_ocr3_base() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();
    initialize_offramp(&mut env, &owner_cap, fee_quoter_cap, dest_transfer_cap);

    // Test getting OCR3 base state
    let _ocr3_base = offramp::get_ocr3_base(&env.state);

    // We successfully got a reference to the OCR3 base state
    // Can't test much more without accessing internals, but the function call succeeded

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
}

#[test]
#[expected_failure(abort_code = offramp::EInvalidRoot)]
public fun test_get_merkle_root() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();
    initialize_offramp(&mut env, &owner_cap, fee_quoter_cap, dest_transfer_cap);

    // Create a test merkle root (32 bytes)
    let test_root = x"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef";

    // Try to get a merkle root that hasn't been committed
    // This should trigger EInvalidRoot since the root hasn't been committed
    let _timestamp = offramp::get_merkle_root(&env.state, test_root);

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
}

#[test]
public fun test_config_signers_and_transmitters() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();
    initialize_offramp(&mut env, &owner_cap, fee_quoter_cap, dest_transfer_cap);

    // Set up OCR3 config first
    let config_digest = b"test_config_digest_32_bytes_long";
    let commit_plugin_type = ocr3_base::ocr_plugin_type_commit();
    let big_f = 1;
    let signers = vector[
        x"1111111111111111111111111111111111111111111111111111111111111111",
        x"2222222222222222222222222222222222222222222222222222222222222222",
        x"3333333333333333333333333333333333333333333333333333333333333333",
        x"4444444444444444444444444444444444444444444444444444444444444444",
        x"5555555555555555555555555555555555555555555555555555555555555555",
    ];
    let transmitters = vector[@0x100, @0x200, @0x300, @0x400, @0x500];

    offramp::set_ocr3_config(
        &env.ref,
        &mut env.state,
        &owner_cap,
        config_digest,
        commit_plugin_type,
        big_f,
        true, // signature verification enabled for commit
        signers,
        transmitters,
    );

    // Get the latest config
    let latest_config = offramp::latest_config_details(&env.state, commit_plugin_type);

    // Test config_signers function
    let config_signers = offramp::config_signers(&latest_config);
    assert!(config_signers.length() == 5);
    assert!(
        config_signers[0] == x"1111111111111111111111111111111111111111111111111111111111111111",
    );
    assert!(
        config_signers[4] == x"5555555555555555555555555555555555555555555555555555555555555555",
    );

    // Test config_transmitters function
    let config_transmitters = offramp::config_transmitters(&latest_config);
    assert!(config_transmitters.length() == 5);
    assert!(config_transmitters[0] == @0x100);
    assert!(config_transmitters[4] == @0x500);

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
}

// === Tests for Uncovered Error Codes ===

#[test]
#[expected_failure(abort_code = offramp::EInvalidRoot)]
public fun test_get_merkle_root_invalid_root() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();
    initialize_offramp(&mut env, &owner_cap, fee_quoter_cap, dest_transfer_cap);

    // Try to get a merkle root that hasn't been committed
    let invalid_root = x"deadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef";
    let _timestamp = offramp::get_merkle_root(&env.state, invalid_root);

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
}

#[test]
#[expected_failure(abort_code = ccip::address::EZeroAddressNotAllowed)]
public fun test_apply_source_chain_config_zero_onramp() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();
    initialize_offramp(&mut env, &owner_cap, fee_quoter_cap, dest_transfer_cap);

    // Try to update with zero/empty onramp address
    // This should fail due to assert_non_zero_address_vector check
    offramp::apply_source_chain_config_updates(
        &env.ref,
        &mut env.state,
        &owner_cap,
        vector[4000], // new chain selector
        vector[true], // enabled
        vector[false], // rmn verification enabled
        vector[vector[]], // Empty onramp address - should fail
        env.scenario.ctx(),
    );

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
}

#[test]
#[expected_failure(abort_code = offramp::EInvalidReportContextLength)]
public fun test_commit_invalid_report_context_length() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();
    initialize_offramp(&mut env, &owner_cap, fee_quoter_cap, dest_transfer_cap);

    let TestEnv { mut scenario, mut state, mut ref, clock } = env;

    // Call commit with invalid report_context length (should be 2, using 1)
    offramp::commit(
        &mut ref,
        &mut state,
        &clock,
        vector[b"invalid"], // length 1 instead of required 2
        vector[],
        vector[],
        scenario.ctx(),
    );

    let env = TestEnv { scenario, state, ref, clock };
    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
}

#[test_only]
public struct TestObj has key {
    id: UID,
}

#[test]
#[expected_failure(abort_code = offramp::EInvalidOwnerCap)]
public fun test_remove_package_id_invalid_owner_cap() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();
    initialize_offramp(&mut env, &owner_cap, fee_quoter_cap, dest_transfer_cap);

    let mut test_obj = TestObj { id: object::new(env.scenario.ctx()) };
    // Create a new owner cap using ownable::new - this will have a different ID
    let (wrong_ownable_state, wrong_owner_cap) = ccip_offramp::ownable::new(
        &mut test_obj.id,
        env.scenario.ctx(),
    );

    // First add a package ID with the correct owner cap
    let test_package_id = @0x999;
    offramp::add_package_id(&mut env.state, &owner_cap, test_package_id);

    // Now try to remove it with the wrong owner cap - should fail at line 320
    offramp::remove_package_id(&mut env.state, &wrong_owner_cap, test_package_id);

    // Clean up the wrong owner cap and ownable state before tear_down
    ccip_offramp::ownable::destroy(wrong_ownable_state, wrong_owner_cap, env.scenario.ctx());

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
    ts::return_shared(test_obj);
}

const EXPLOITER: address = @0x199;

#[test]
#[expected_failure(abort_code = mcms_registry::EProofTypeNotRegistered)]
public fun test_release_cap_should_fail_with_proof_type_not_registered() {
    let (mut env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();
    initialize_offramp(&mut env, &owner_cap, fee_quoter_cap, dest_transfer_cap);

    mcms_registry::test_init(env.scenario.ctx());

    env.scenario.next_tx(OWNER);

    // -----------------------------------------------------
    // Transfer the OffRamp ownership to MCMS

    let mcms = mcms_registry::get_multisig_address();
    offramp::transfer_ownership(&env.ref, &mut env.state, &owner_cap, mcms, env.scenario.ctx());
    assert!(offramp::pending_transfer_from(&env.state).extract() == OWNER);
    assert!(offramp::pending_transfer_to(&env.state).extract() == mcms);
    assert!(!offramp::pending_transfer_accepted(&env.state).extract());

    assert!(offramp::owner(&env.state) == OWNER);

    let TestEnv {
        scenario,
        state,
        ref,
        clock,
    } = env;

    ts::end(scenario);

    // -----------------------------------------------------
    // MCMS accepts ownership

    let scenario_2 = ts::begin(mcms);

    let mut env = TestEnv {
        scenario: scenario_2,
        state,
        ref,
        clock,
    };

    offramp::accept_ownership(&env.ref, &mut env.state, env.scenario.ctx());
    assert!(offramp::pending_transfer_from(&env.state).extract() == OWNER);
    assert!(offramp::pending_transfer_to(&env.state).extract() == mcms);
    assert!(offramp::pending_transfer_accepted(&env.state).extract());

    assert!(offramp::owner(&env.state) == OWNER);

    let TestEnv {
        scenario,
        state,
        ref,
        clock,
    } = env;

    ts::end(scenario);

    // -----------------------------------------------------
    // Execute the ownership transfer

    let scenario_3 = ts::begin(OWNER);

    let mut env = TestEnv {
        scenario: scenario_3,
        state,
        ref,
        clock,
    };

    let mut registry = ts::take_shared<Registry>(&env.scenario);
    let cap_id = object::id(&owner_cap);

    offramp::execute_ownership_transfer_to_mcms(
        &env.ref,
        owner_cap,
        &mut env.state,
        &mut registry,
        mcms,
        env.scenario.ctx(),
    );

    assert!(offramp::owner(&env.state) == mcms);

    let TestEnv {
        scenario,
        state,
        ref,
        clock,
    } = env;

    ts::end(scenario);

    // -----------------------------------------------------
    // Attempt to extract OwnerCap with wrong witness (StaticConfig instead of McmsCallback)
    // This SHOULD FAIL with EProofTypeNotRegistered

    let scenario_4 = ts::begin(EXPLOITER);

    let env = TestEnv {
        scenario: scenario_4,
        state,
        ref,
        clock,
    };

    let static_config = offramp::get_static_config(&env.ref, &env.state);

    // This will abort with EProofTypeNotRegistered because StaticConfig != McmsCallback
    let extracted_owner_cap = mcms_registry::release_cap<StaticConfig, OwnerCap>(
        &mut registry,
        static_config,
    );

    // Code below will never execute due to abort above
    let extracted_cap_id = object::id(&extracted_owner_cap);
    assert!(cap_id == extracted_cap_id);

    ts::return_shared(registry);
    tear_down(env);
    transfer::public_transfer(extracted_owner_cap, EXPLOITER);
}

#[test]
public fun test_calculate_message_hash() {
    let (env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();
    
    // Expected hash values
    let expected_hash_no_tokens = x"9f9be87e216efa0b1571131d9295e3802c5c9a3d6e369d230c72520a2e854a9e";
    let expected_hash_with_tokens = x"d183d22cb0b713da1b6b42d9c35cc9e1268257ff703c6579d6aa68fdfb1ff4b2";
    
    // Test with no tokens
    let message_id = x"1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef";
    let source_chain_selector = 123456789u64;
    let dest_chain_selector = 987654321u64;
    let sequence_number = 42u64;
    let nonce = 0u64; // Must be 0 for out of order execution
    let sender = x"8765432109fedcba8765432109fedcba87654321";
    let receiver = @0x1234;
    let on_ramp = b"source-onramp-address";
    let data = b"sample message data";
    let gas_limit = 500000u256;
    let token_receiver = @0x0; // No tokens
    let source_pool_addresses = vector<vector<u8>>[];
    let dest_token_addresses = vector<address>[];
    let dest_gas_amounts = vector<u32>[];
    let extra_datas = vector<vector<u8>>[];
    let amounts = vector<u256>[];

    let hash_no_tokens = offramp::calculate_message_hash(
        &env.ref,
        message_id,
        source_chain_selector,
        dest_chain_selector,
        sequence_number,
        nonce,
        sender,
        receiver,
        on_ramp,
        data,
        gas_limit,
        token_receiver,
        source_pool_addresses,
        dest_token_addresses,
        dest_gas_amounts,
        extra_datas,
        amounts,
    );

    // Verify hash is 32 bytes
    assert!(hash_no_tokens.length() == 32, 0);
    // Verify hash matches expected value
    assert!(hash_no_tokens == expected_hash_no_tokens, 1);

    // Test with tokens
    let token_receiver = @0x5678;
    let source_pool_addresses = vector[
        x"abcdef1234567890abcdef1234567890abcdef12",
        x"123456789abcdef123456789abcdef123456789a",
    ];
    let dest_token_addresses = vector[@0x5678, @0x9abc];
    let dest_gas_amounts = vector[10000u32, 20000u32];
    let extra_datas = vector[x"00112233", x"ffeeddcc"];
    let amounts = vector[1000000u256, 5000000u256];

    let hash_with_tokens = offramp::calculate_message_hash(
        &env.ref,
        message_id,
        source_chain_selector,
        dest_chain_selector,
        sequence_number,
        nonce,
        sender,
        receiver,
        on_ramp,
        data,
        gas_limit,
        token_receiver,
        source_pool_addresses,
        dest_token_addresses,
        dest_gas_amounts,
        extra_datas,
        amounts,
    );

    // Verify hash is 32 bytes
    assert!(hash_with_tokens.length() == 32, 2);
    // Verify hash matches expected value
    assert!(hash_with_tokens == expected_hash_with_tokens, 3);
    
    // Hashes should be different when tokens are included
    assert!(hash_no_tokens != hash_with_tokens, 4);

    // Test that changing any parameter changes the hash
    let different_sequence = offramp::calculate_message_hash(
        &env.ref,
        message_id,
        source_chain_selector,
        dest_chain_selector,
        sequence_number + 1, // Changed
        nonce,
        sender,
        receiver,
        on_ramp,
        data,
        gas_limit,
        token_receiver,
        source_pool_addresses,
        dest_token_addresses,
        dest_gas_amounts,
        extra_datas,
        amounts,
    );
    assert!(hash_with_tokens != different_sequence, 5);

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
    transfer::public_transfer(fee_quoter_cap, OWNER);
    transfer::public_transfer(dest_transfer_cap, OWNER);
}

#[test]
public fun test_calculate_metadata_hash() {
    let (env, owner_cap, fee_quoter_cap, dest_transfer_cap) = setup();

    // Expected hash values
    let expected_metadata_hash = x"b62ec658417caa5bcc6ff1d8c45f8b1cb52e1b0ed71603a04b250b107ed836d9";
    let expected_metadata_hash_different_source = x"89da72ab93f7bd546d60b58a1e1b5f628fd456fe163614ff1e31a2413ca1b55a";

    let source_chain_selector = 123456789u64;
    let dest_chain_selector = 987654321u64;
    let on_ramp = b"source-onramp-address";

    let metadata_hash = offramp::calculate_metadata_hash(
        &env.ref,
        source_chain_selector,
        dest_chain_selector,
        on_ramp,
    );

    // Verify hash is 32 bytes
    assert!(metadata_hash.length() == 32, 0);
    // Verify hash matches expected value
    assert!(metadata_hash == expected_metadata_hash, 1);

    // Test that changing source chain selector produces different hash
    let metadata_hash_different_source = offramp::calculate_metadata_hash(
        &env.ref,
        source_chain_selector + 1,
        dest_chain_selector,
        on_ramp,
    );
    
    assert!(metadata_hash != metadata_hash_different_source, 2);
    // Verify the different hash matches expected value
    assert!(metadata_hash_different_source == expected_metadata_hash_different_source, 3);

    // Test that changing destination chain selector produces different hash
    let metadata_hash_different_dest = offramp::calculate_metadata_hash(
        &env.ref,
        source_chain_selector,
        dest_chain_selector + 1,
        on_ramp,
    );
    assert!(metadata_hash != metadata_hash_different_dest, 4);

    // Test that changing on_ramp produces different hash
    let metadata_hash_different_onramp = offramp::calculate_metadata_hash(
        &env.ref,
        source_chain_selector,
        dest_chain_selector,
        b"different-onramp-address",
    );
    assert!(metadata_hash != metadata_hash_different_onramp, 5);

    tear_down(env);
    ts::return_to_address(OWNER, owner_cap);
    transfer::public_transfer(fee_quoter_cap, OWNER);
    transfer::public_transfer(dest_transfer_cap, OWNER);
}
